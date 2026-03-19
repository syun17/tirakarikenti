from fastapi import FastAPI, UploadFile, File
import uvicorn
import shutil
import os
import base64
import io

import torch
import clip
from PIL import Image
import matplotlib.pyplot as plt
import numpy as np
from ultralytics import YOLO
import time

device = "cpu"
model, preprocess = clip.load("ViT-L/14", device=device)

# ===== 改良プロンプト =====

clean_labels = [
    "a very clean minimalist bedroom with no clutter",
    "a spotless clean bedroom",
    "a perfectly organized tidy bedroom",
    "a neat and organized bedroom",
]

normal_labels = [
    "a normal lived in bedroom",
    "a typical bedroom someone lives in",
    "an average bedroom with some items",
    "a slightly cluttered bedroom"
]

messy_labels = [
    "a messy bedroom with clothes on the floor",
    "a cluttered bedroom",
    "a disorganized bedroom",
    "a bedroom with many things scattered around"
]

very_messy_labels = [
    "a very messy cluttered bedroom",
    "an extremely messy dirty bedroom",
    "a chaotic messy room",
    "a dirty cluttered bedroom"
]

labels = clean_labels + normal_labels + messy_labels + very_messy_labels

text = clip.tokenize(labels).to(device)
with torch.no_grad():
    text_features = model.encode_text(text)
    text_features /= text_features.norm(dim=-1, keepdim=True)

yolo_model = YOLO("/app/yolo11l.pt")

app = FastAPI()


UPLOAD_DIR = "uploads"
os.makedirs(UPLOAD_DIR, exist_ok=True)

@app.post("/analyze/")
async def analyze(file: UploadFile = File(...)):

    # ===== 画像をメモリ上で読み込み =====
    contents:bytes = await file.read()
    if len(contents) > 5 * 1024 * 1024:
        return {"error": "File too large"}

    image_pil = Image.open(io.BytesIO(contents)).convert("RGB")

    # ===== CLIP用前処理 =====
    image = preprocess(image_pil).unsqueeze(0).to(device)

    # ===== 推論 =====
    with torch.no_grad():

        image_features = model.encode_image(image)

        # 正規化（CLIP公式推奨）
        image_features /= image_features.norm(dim=-1, keepdim=True)

        logits = (100.0 * image_features @ text_features.T).softmax(dim=-1)

    scores = logits[0].cpu().numpy()

    # ===== カテゴリ平均 =====
    clean_score = np.mean(scores[0:4])
    normal_score = np.mean(scores[4:8])
    messy_score = np.mean(scores[8:12])
    very_messy_score = np.mean(scores[12:16])

    clean_ratio = clean_score / (clean_score + normal_score + messy_score + very_messy_score)
    normal_ratio = normal_score / (clean_score + normal_score + messy_score + very_messy_score)
    messy_ratio = messy_score / (clean_score + normal_score + messy_score + very_messy_score)
    very_messy_ratio = very_messy_score / (clean_score + normal_score + messy_score + very_messy_score)

    # ===== Messinessスコア（改良） =====
    messiness_score = (normal_ratio*1 + messy_ratio*2 + very_messy_ratio*4) / 4 * 100

    # ---------- YOLO inference ----------

    img_np = np.array(image_pil)
    results = yolo_model(img_np,imgsz=640,augment=True ,device="cpu")

    objects = {}

    for box in results[0].boxes:
        cls = int(box.cls[0])
        conf = float(box.conf[0])
        name = yolo_model.names[cls]

        if name in ["bottle","cup","wine glass","bowl","can"] and conf > 0.5:
            objects[name] = objects.get(name, 0) + 1
    


    return {
        "image": file.filename,
        "score": round(float(messiness_score), 1),
        "yolo": objects
    }


@app.get("/")
async def root():
    return {"message": "Hello World"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, log_level="debug")
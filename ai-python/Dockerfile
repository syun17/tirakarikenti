FROM python:3.10-slim

# システム依存パッケージのインストール
RUN apt-get update && apt-get install -y \
    git \
    wget \
    libgl1 \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender-dev \
    libgomp1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Pythonパッケージのインストール
RUN pip install --no-cache-dir --upgrade pip

# PyTorch (CPU版)
RUN pip install --no-cache-dir torch torchvision --index-url https://download.pytorch.org/whl/cpu

# CLIP
RUN pip install --no-cache-dir git+https://github.com/openai/CLIP.git

# その他の依存パッケージ
RUN pip install --no-cache-dir \
    fastapi \
    uvicorn[standard] \
    python-multipart \
    ultralytics \
    Pillow \
    numpy \
    matplotlib

# YOLOモデルの事前ダウンロード
RUN python -c "from ultralytics import YOLO; YOLO('yolo11l.pt')"

# アプリケーションのコピー
COPY main.py .

# アップロードディレクトリの作成
RUN mkdir -p uploads

EXPOSE 8000

CMD ["python", "main.py"]

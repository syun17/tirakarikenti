"use client";

import { useState, useCallback } from "react";
import { PhotoUploader } from "@/components/photo-uploader";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { Save, Sparkles, Wand2 } from "lucide-react";
import { Navigation } from "@/components/navigation";
import { useRouter } from "next/navigation";

export default function AnalysisPage() {
  const [imageDataUrl, setImageDataUrl] = useState<string | null>(null);
  const [score, setScore] = useState(50);
  const [note, setNote] = useState("");
  const [isSaving, setIsSaving] = useState(false);
  const { toast } = useToast();

  const router = useRouter();

  const handleAnalyze = useCallback(async () => {
    if (!imageDataUrl) {
      toast({
        title: "写真が必要です",
        description: "分析する前に写真をアップロードしてください。",
        variant: "destructive",
      });
      return;
    }

    setIsSaving(true);

    try {
      // Base64 to Blob
      const res = await fetch(imageDataUrl);
      const blob = await res.blob();
      
      const formData = new FormData();
      formData.append("image", blob, "image.jpg");
      formData.append("room_id", "1"); // デフォルトの部屋ID

      const response = await fetch("http://localhost:8080/analyze-python", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("分析に失敗しました");
      }

      const result = await response.json();
      
      setScore(Math.round(result.score));
      
      // YOLO結果をメモ形式にして表示
      let yoloNote = "";
      for (const [name, count] of Object.entries(result.yolo)) {
        if (yoloNote) yoloNote += ", ";
        yoloNote += `${name}: ${count}`;
      }
      setNote(yoloNote);

      toast({
        title: "分析が完了しました",
        description: `片付け度スコア: ${Math.round(result.score)}点`,
      });

      // Reset image and redirect or wait
      setIsSaving(false);

    } catch (error) {
      console.error(error);
      toast({
        title: "エラーが発生しました",
        description: "バックエンドとの通信に失敗しました。",
        variant: "destructive",
      });
      setIsSaving(false);
    }
  }, [imageDataUrl, toast, router]);

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      <main className="container max-w-2xl mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-foreground">分析</h1>
          <p className="text-muted-foreground mt-1">
            部屋の写真を撮って片付け度を記録しましょう
          </p>
        </div>

        <div className="space-y-6">
          {/* Photo Upload */}
          <section>
            <h2 className="text-sm font-medium text-muted-foreground mb-3 flex items-center gap-2">
              <Sparkles className="h-4 w-4" />
              写真をアップロード
            </h2>
            <PhotoUploader
              imageDataUrl={imageDataUrl}
              onImageChange={setImageDataUrl}
            />
          </section>

          {/* Results Reveal (Optional, since we redirect) */}
          {score !== 50 && (
            <Card className="bg-primary/5 border-primary/20">
              <CardContent className="pt-6">
                <div className="text-center">
                  <p className="text-sm font-medium text-muted-foreground mb-1">
                    分析結果
                  </p>
                  <p className="text-4xl font-bold text-primary mb-2">
                    {score}点
                  </p>
                  <p className="text-sm text-foreground">{note}</p>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Analyze Button */}
          <Button
            size="lg"
            className="w-full"
            onClick={handleAnalyze}
            disabled={isSaving || !imageDataUrl}
          >
            <Wand2 className="h-4 w-4 mr-2" />
            {isSaving ? "分析中..." : "AIで分析して保存"}
          </Button>
        </div>
      </main>
    </div>
  );
}

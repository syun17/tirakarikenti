"use client";

import { useCallback, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Upload, ImageIcon, X } from "lucide-react";
import { Button } from "@/components/ui/button";

interface PhotoUploaderProps {
  imageDataUrl: string | null;
  onImageChange: (dataUrl: string | null) => void;
}

export function PhotoUploader({
  imageDataUrl,
  onImageChange,
}: PhotoUploaderProps) {
  const [isDragging, setIsDragging] = useState(false);

  const handleFile = useCallback(
    (file: File) => {
      if (!file.type.startsWith("image/")) return;

      const reader = new FileReader();
      reader.onload = (e) => {
        const result = e.target?.result as string;
        onImageChange(result);
      };
      reader.readAsDataURL(file);
    },
    [onImageChange]
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragging(false);
      const file = e.dataTransfer.files[0];
      if (file) handleFile(file);
    },
    [handleFile]
  );

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (file) handleFile(file);
    },
    [handleFile]
  );

  const handleRemove = useCallback(() => {
    onImageChange(null);
  }, [onImageChange]);

  if (imageDataUrl) {
    return (
      <Card className="overflow-hidden">
        <CardContent className="p-0 relative">
          <img
            src={imageDataUrl}
            alt="アップロードされた写真"
            className="w-full aspect-[4/3] object-cover"
          />
          <Button
            variant="destructive"
            size="icon"
            className="absolute top-3 right-3"
            onClick={handleRemove}
          >
            <X className="h-4 w-4" />
            <span className="sr-only">写真を削除</span>
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card
      className={`border-2 border-dashed transition-colors ${
        isDragging ? "border-primary bg-primary/5" : "border-muted-foreground/25"
      }`}
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
    >
      <CardContent className="p-8">
        <label className="flex flex-col items-center justify-center gap-4 cursor-pointer">
          <div className="p-4 rounded-full bg-muted">
            {isDragging ? (
              <Upload className="h-8 w-8 text-primary" />
            ) : (
              <ImageIcon className="h-8 w-8 text-muted-foreground" />
            )}
          </div>
          <div className="text-center">
            <p className="font-medium text-foreground">
              写真をドラッグ&ドロップ
            </p>
            <p className="text-sm text-muted-foreground mt-1">
              またはクリックして選択
            </p>
          </div>
          <input
            type="file"
            accept="image/*"
            className="hidden"
            onChange={handleInputChange}
          />
        </label>
      </CardContent>
    </Card>
  );
}

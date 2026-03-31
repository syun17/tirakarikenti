"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import type { TidyRecord } from "@/lib/types";

interface RecordCardProps {
  record: TidyRecord;
  onDelete: (id: string) => void;
}

function getScoreColor(score: number): string {
  if (score < 40) return "bg-score-low";
  if (score < 70) return "bg-score-medium";
  return "bg-score-high";
}

export function RecordCard({ record, onDelete }: RecordCardProps) {
  const date = new Date(record.date);
  const formattedDate = date.toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
  const formattedTime = date.toLocaleTimeString("ja-JP", {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <Card className="overflow-hidden">
      <CardContent className="p-0">
        <div className="flex gap-3 p-3">
          {/* Info */}
          <div className="flex-1 min-w-0 flex flex-col justify-between">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <span
                  className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold text-primary-foreground ${getScoreColor(
                    record.score
                  )}`}
                >
                  {record.score}点
                </span>
                <span className="text-xs text-muted-foreground">
                  {formattedDate}
                </span>
              </div>
              <p className="text-xs text-muted-foreground">{formattedTime}</p>
              {record.note && (
                <p className="text-sm text-foreground mt-1 whitespace-pre-wrap">
                  {record.note}
                </p>
              )}
            </div>
          </div>

          {/* Delete Button */}
          <div className="flex-shrink-0 self-start">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-muted-foreground hover:text-destructive"
              onClick={() => onDelete(record.id)}
            >
              <Trash2 className="h-4 w-4" />
              <span className="sr-only">削除</span>
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

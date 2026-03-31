"use client";

import { Slider } from "@/components/ui/slider";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface ScoreSliderProps {
  value: number;
  onChange: (value: number) => void;
}

function getScoreColor(score: number): string {
  if (score < 40) return "text-score-low";
  if (score < 70) return "text-score-medium";
  return "text-score-high";
}

function getScoreLabel(score: number): string {
  if (score < 20) return "とても散らかっている";
  if (score < 40) return "散らかっている";
  if (score < 60) return "普通";
  if (score < 80) return "きれい";
  return "とてもきれい";
}

function getScoreBgColor(score: number): string {
  if (score < 40) return "bg-score-low";
  if (score < 70) return "bg-score-medium";
  return "bg-score-high";
}

export function ScoreSlider({ value, onChange }: ScoreSliderProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-medium flex items-center justify-between">
          <span>片付け度スコア</span>
          <span className={`text-2xl font-bold ${getScoreColor(value)}`}>
            {value}
            <span className="text-sm font-normal text-muted-foreground ml-1">
              / 100
            </span>
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center gap-3">
          <span className="text-xs text-muted-foreground w-8">0</span>
          <Slider
            value={[value]}
            onValueChange={(v) => onChange(v[0])}
            max={100}
            min={0}
            step={1}
            className="flex-1"
          />
          <span className="text-xs text-muted-foreground w-8 text-right">
            100
          </span>
        </div>
        <div className="flex items-center justify-center">
          <span
            className={`px-3 py-1 rounded-full text-sm font-medium text-primary-foreground ${getScoreBgColor(
              value
            )}`}
          >
            {getScoreLabel(value)}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}

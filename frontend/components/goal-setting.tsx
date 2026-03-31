"use client";

import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Slider } from "@/components/ui/slider";
import { Target } from "lucide-react";

interface GoalSettingProps {
  goalScore: number;
  onGoalChange: (goal: number) => void;
}

export function GoalSetting({ goalScore, onGoalChange }: GoalSettingProps) {
  const [open, setOpen] = useState(false);
  const [tempGoal, setTempGoal] = useState(goalScore);

  const handleOpen = (isOpen: boolean) => {
    if (isOpen) {
      setTempGoal(goalScore);
    }
    setOpen(isOpen);
  };

  const handleSave = () => {
    onGoalChange(tempGoal);
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <Target className="h-4 w-4" />
          目標: {goalScore}点
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>目標スコアを設定</DialogTitle>
        </DialogHeader>
        <div className="space-y-6 py-4">
          <div className="text-center">
            <span className="text-4xl font-bold text-primary">{tempGoal}</span>
            <span className="text-lg text-muted-foreground ml-1">点</span>
          </div>
          <div className="flex items-center gap-3 px-2">
            <span className="text-sm text-muted-foreground">0</span>
            <Slider
              value={[tempGoal]}
              onValueChange={(v) => setTempGoal(v[0])}
              max={100}
              min={0}
              step={5}
              className="flex-1"
            />
            <span className="text-sm text-muted-foreground">100</span>
          </div>
          <Button onClick={handleSave} className="w-full">
            保存
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

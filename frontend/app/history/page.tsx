"use client";

import { useState, useEffect, useCallback } from "react";
import { Navigation } from "@/components/navigation";
import { TidyChart } from "@/components/tidy-chart";
import { RecordCard } from "@/components/record-card";
import { GoalSetting } from "@/components/goal-setting";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  getRecords,
  deleteRecord,
  getSettings,
  saveSettings,
  filterRecordsByPeriod,
} from "@/lib/storage";
import { useToast } from "@/hooks/use-toast";
import type { TidyRecord, PeriodFilter, AppSettings } from "@/lib/types";
import { Calendar, TrendingUp } from "lucide-react";

export default function HistoryPage() {
  const [records, setRecords] = useState<TidyRecord[]>([]);
  const [settings, setSettings] = useState<AppSettings>({ goalScore: 80 });
  const [period, setPeriod] = useState<PeriodFilter>("week");
  const [isLoading, setIsLoading] = useState(true);
  const { toast } = useToast();

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const fetchedRecords = await getRecords();
      setRecords(fetchedRecords);
      setSettings(getSettings());
      setIsLoading(false);
    };
    loadData();
  }, []);

  const filteredRecords = filterRecordsByPeriod(records, period);
  const sortedRecords = [...filteredRecords].sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
  );

  const handleDelete = useCallback(
    async (id: string) => {
      await deleteRecord(id);
      const updatedRecords = await getRecords();
      setRecords(updatedRecords);
      toast({
        title: "記録を削除しました",
      });
    },
    [toast]
  );

  const handleGoalChange = useCallback(
    (goalScore: number) => {
      const newSettings = { ...settings, goalScore };
      saveSettings(newSettings);
      setSettings(newSettings);
      toast({
        title: "目標を更新しました",
        description: `新しい目標: ${goalScore}点`,
      });
    },
    [settings, toast]
  );

  // Calculate stats
  const avgScore =
    filteredRecords.length > 0
      ? Math.round(
          filteredRecords.reduce((sum, r) => sum + r.score, 0) /
            filteredRecords.length
        )
      : 0;

  const latestScore = sortedRecords[0]?.score ?? 0;

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <Navigation />
        <main className="container max-w-2xl mx-auto px-4 py-8">
          <div className="animate-pulse space-y-4">
            <div className="h-8 w-32 bg-muted rounded" />
            <div className="h-[300px] bg-muted rounded-lg" />
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      <main className="container max-w-2xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-start justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-foreground">履歴</h1>
            <p className="text-muted-foreground mt-1">
              片付け度の推移と記録を確認
            </p>
          </div>
          <GoalSetting
            goalScore={settings.goalScore}
            onGoalChange={handleGoalChange}
          />
        </div>

        {/* Stats */}
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="bg-card rounded-lg border p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <TrendingUp className="h-4 w-4" />
              <span className="text-xs font-medium">平均スコア</span>
            </div>
            <p className="text-2xl font-bold text-foreground">
              {avgScore}
              <span className="text-sm font-normal text-muted-foreground ml-1">
                点
              </span>
            </p>
          </div>
          <div className="bg-card rounded-lg border p-4">
            <div className="flex items-center gap-2 text-muted-foreground mb-1">
              <Calendar className="h-4 w-4" />
              <span className="text-xs font-medium">最新スコア</span>
            </div>
            <p className="text-2xl font-bold text-foreground">
              {latestScore}
              <span className="text-sm font-normal text-muted-foreground ml-1">
                点
              </span>
            </p>
          </div>
        </div>

        {/* Period Filter */}
        <Tabs
          value={period}
          onValueChange={(v) => setPeriod(v as PeriodFilter)}
          className="mb-6"
        >
          <TabsList className="w-full">
            <TabsTrigger value="week" className="flex-1">
              週間
            </TabsTrigger>
            <TabsTrigger value="month" className="flex-1">
              月間
            </TabsTrigger>
            <TabsTrigger value="all" className="flex-1">
              全期間
            </TabsTrigger>
          </TabsList>
        </Tabs>

        {/* Chart */}
        <div className="mb-6">
          <TidyChart
            records={filteredRecords}
            goalScore={settings.goalScore}
          />
        </div>

        {/* Records List */}
        <section>
          <h2 className="text-base font-semibold text-foreground mb-4">
            記録一覧
            <span className="text-sm font-normal text-muted-foreground ml-2">
              ({sortedRecords.length}件)
            </span>
          </h2>
          {sortedRecords.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              <p>この期間の記録はありません</p>
              <p className="text-sm mt-1">
                分析ページから新しい記録を追加してください
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {sortedRecords.map((record) => (
                <RecordCard
                  key={record.id}
                  record={record}
                  onDelete={handleDelete}
                />
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

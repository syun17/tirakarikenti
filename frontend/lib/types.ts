export interface TidyRecord {
  id: string;
  date: string; // ISO形式
  score: number; // 1-100
  note?: string; // オプションのメモ
}

export interface AppSettings {
  goalScore: number; // 目標スコア (1-100)
}

export type PeriodFilter = "week" | "month" | "all";

import type { TidyRecord, AppSettings } from "./types";

const RECORDS_KEY = "tidy-records";
const SETTINGS_KEY = "tidy-settings";

const DEFAULT_SETTINGS: AppSettings = {
  goalScore: 80,
};

const API_BASE = "http://localhost:8080";

// Records
export async function getRecords(): Promise<TidyRecord[]> {
  try {
    const response = await fetch(`${API_BASE}/rooms/history`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ room_id: 1 }), // デフォルトの部屋ID
    });
    if (!response.ok) return [];
    
    interface BackendImg {
      id: number;
      score: number;
      note: string;
      created_at: string;
      room_id: number;
    }
    
    const data = (await response.json()) as BackendImg[];
    return data.map((img) => ({
      id: img.id.toString(),
      date: img.created_at,
      score: Math.round(img.score),
      note: img.note,
    }));
  } catch (error) {
    console.error("Failed to fetch records:", error);
    return [];
  }
}

export async function deleteRecord(id: string): Promise<void> {
  try {
    await fetch(`${API_BASE}/images/${id}`, {
      method: "DELETE",
    });
  } catch (error) {
    console.error("Failed to delete record:", error);
  }
}

// saveRecordは分析ボタンのフローで直接APIを叩くため、ここでは廃止または更新
export async function saveRecord(record: Omit<TidyRecord, "id" | "date">): Promise<void> {
  // 実際には AnalyzePythonHandler を通じて保存されるため、
  // この関数は現在は直接使用されないか、ラップする形になる。
}

// Settings
export function getSettings(): AppSettings {
  if (typeof window === "undefined") return DEFAULT_SETTINGS;
  const data = localStorage.getItem(SETTINGS_KEY);
  if (!data) return DEFAULT_SETTINGS;
  try {
    return JSON.parse(data) as AppSettings;
  } catch {
    return DEFAULT_SETTINGS;
  }
}

export function saveSettings(settings: AppSettings): void {
  localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings));
}

// Utility
export function generateId(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
}

export function filterRecordsByPeriod(
  records: TidyRecord[],
  period: "week" | "month" | "all"
): TidyRecord[] {
  if (period === "all") return records;

  const now = new Date();
  const cutoff = new Date();

  if (period === "week") {
    cutoff.setDate(now.getDate() - 7);
  } else if (period === "month") {
    cutoff.setMonth(now.getMonth() - 1);
  }

  return records.filter((r) => new Date(r.date) >= cutoff);
}

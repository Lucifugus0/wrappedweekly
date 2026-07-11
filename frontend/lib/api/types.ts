export type User = {
  id: string;
  email: string;
  name: string;
  created_at: string;
};

export type ActivityCategory = "workout" | "reading" | "coding" | "spending";

export type Activity = {
  id: string;
  category: ActivityCategory;
  value: number;
  note: string | null;
  occurred_at: string;
  created_at: string;
  updated_at: string;
};

export type ActivityListResponse = {
  items: Activity[];
  page: number;
  size: number;
  total: number;
};

export type CategoryTotal = {
  category: ActivityCategory;
  total: number;
  count: number;
};

export type DayTotal = {
  date: string;
  total: number;
};

export type RecapStats = {
  week_start: string;
  week_end: string;
  totals_by_category: CategoryTotal[];
  top_category: ActivityCategory | null;
  most_productive_day: string | null;
  daily_breakdown: DayTotal[];
  total_activities: number;
  total_value: number;
  prev_week_total_value: number;
  change_vs_prev_week_pct: number | null;
};

export type Recap = {
  id: string;
  slug: string;
  week_start: string;
  week_end: string;
  stats: RecapStats;
  narrative: string;
  created_at: string;
};

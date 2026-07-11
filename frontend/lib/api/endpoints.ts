import { api } from "./client";
import type {
  Activity,
  ActivityListResponse,
  Recap,
  RecapStats,
  User,
} from "./types";

export type RegisterInput = { email: string; password: string; name: string };
export type LoginInput = { email: string; password: string };
export type ActivityInput = {
  category: string;
  value: number;
  note?: string | null;
  occurred_at: string;
};

export const authApi = {
  register: (input: RegisterInput) => api.post<User>("/auth/register", input),
  login: (input: LoginInput) => api.post<{ user: User }>("/auth/login", input),
  logout: () => api.post<null>("/auth/logout"),
  me: () => api.get<User>("/auth/me"),
};

export type ActivityListFilter = {
  category?: string;
  from?: string;
  to?: string;
};

export const activityApi = {
  list: (page = 1, size = 20, filter: ActivityListFilter = {}) => {
    const params = new URLSearchParams({ page: String(page), size: String(size) });
    if (filter.category) params.set("category", filter.category);
    if (filter.from) params.set("from", filter.from);
    if (filter.to) params.set("to", filter.to);
    return api.get<ActivityListResponse>(`/activities?${params.toString()}`);
  },
  get: (id: string) => api.get<Activity>(`/activities/${id}`),
  create: (input: ActivityInput) => api.post<Activity>("/activities", input),
  update: (id: string, input: ActivityInput) => api.put<Activity>(`/activities/${id}`, input),
  delete: (id: string) => api.delete<null>(`/activities/${id}`),
};

export const recapApi = {
  generate: (weekOf?: string, force = false) =>
    api.post<Recap>("/recaps/generate", { week_of: weekOf, force }),
  list: () => api.get<Recap[]>("/recaps"),
  get: (id: string) => api.get<Recap>(`/recaps/${id}`),
};

export const dashboardApi = {
  summary: () => api.get<RecapStats>("/dashboard/summary"),
};

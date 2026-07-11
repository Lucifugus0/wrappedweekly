// Central fetch wrapper for all backend calls from Client Components.
// Server Components (e.g. the public recap page) fetch directly with an
// absolute URL instead — see app/w/[slug]/page.tsx.

export class ApiError extends Error {
  status: number;
  fieldErrors?: Record<string, string>;

  constructor(status: number, message: string, fieldErrors?: Record<string, string>) {
    super(message);
    this.status = status;
    this.fieldErrors = fieldErrors;
  }
}

type Envelope<T> = {
  data: T;
  message: string;
};

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "/api/v1";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  const body = (await res.json().catch(() => null)) as Envelope<T> | null;

  if (!res.ok) {
    const message = body?.message || `Request gagal (${res.status})`;
    const fieldErrors =
      body && typeof body.data === "object" && body.data !== null && !Array.isArray(body.data)
        ? (body.data as unknown as Record<string, string>)
        : undefined;
    throw new ApiError(res.status, message, fieldErrors);
  }

  return (body?.data ?? (undefined as unknown)) as T;
}

export const api = {
  get: <T>(path: string) => request<T>(path, { method: "GET" }),
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, { method: "POST", body: body ? JSON.stringify(body) : undefined }),
  put: <T>(path: string, body?: unknown) =>
    request<T>(path, { method: "PUT", body: body ? JSON.stringify(body) : undefined }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
};

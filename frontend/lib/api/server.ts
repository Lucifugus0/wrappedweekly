// Server-side fetch helper for Server Components. Unlike lib/api/client.ts
// (used by Client Components via relative "/api/v1" + Next.js rewrites),
// code running on the server has no browser origin to resolve a relative URL
// against, so it must call the backend directly with an absolute URL.
const BACKEND_URL = process.env.BACKEND_INTERNAL_URL || "http://localhost:8080";

type Envelope<T> = { data: T; message: string };

export async function serverGet<T>(path: string): Promise<T | null> {
  const res = await fetch(`${BACKEND_URL}/api/v1${path}`, {
    cache: "no-store",
  });

  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`Backend request failed: ${res.status}`);

  const body = (await res.json()) as Envelope<T>;
  return body.data;
}

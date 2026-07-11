// <input type="datetime-local"> gives/accepts "YYYY-MM-DDTHH:mm" with no
// timezone info, interpreted as the browser's local time. The backend expects
// RFC3339. We convert through a real Date so the local-time intent survives.

export function toDatetimeLocalValue(isoString: string): string {
  const d = new Date(isoString);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function fromDatetimeLocalValue(value: string): string {
  return new Date(value).toISOString();
}

export function nowAsDatetimeLocalValue(): string {
  return toDatetimeLocalValue(new Date().toISOString());
}

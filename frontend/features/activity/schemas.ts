import { z } from "zod";

export const activityCategories = ["workout", "reading", "coding", "spending"] as const;

// value stays a string through react-hook-form (native number inputs emit
// strings); it's parsed to a number only at submit time in the dialog. This
// avoids the z.coerce input/output type mismatch that trips up the resolver's
// generic inference.
export const activitySchema = z.object({
  category: z.enum(activityCategories, { message: "Pilih kategori" }),
  value: z
    .string()
    .min(1, "Value wajib diisi")
    .refine((v) => !Number.isNaN(Number(v)) && Number(v) >= 0, "Value harus angka >= 0"),
  note: z.string().trim().optional(),
  occurred_at: z.string().min(1, "Tanggal wajib diisi"),
});
export type ActivityFormValues = z.infer<typeof activitySchema>;

export const categoryLabels: Record<(typeof activityCategories)[number], string> = {
  workout: "Olahraga",
  reading: "Membaca",
  coding: "Ngoding",
  spending: "Pengeluaran",
};

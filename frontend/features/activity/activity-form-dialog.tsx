"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { ApiError } from "@/lib/api/client";
import type { Activity } from "@/lib/api/types";
import { useCreateActivity, useUpdateActivity } from "./use-activities";
import { activityCategories, activitySchema, categoryLabels, type ActivityFormValues } from "./schemas";
import { fromDatetimeLocalValue, nowAsDatetimeLocalValue, toDatetimeLocalValue } from "./datetime";

export function ActivityFormDialog({
  open,
  onOpenChange,
  activity,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  activity?: Activity | null;
}) {
  const isEdit = !!activity;
  const createMutation = useCreateActivity();
  const updateMutation = useUpdateActivity();
  const mutation = isEdit ? updateMutation : createMutation;

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<ActivityFormValues>({
    resolver: zodResolver(activitySchema),
    defaultValues: {
      category: "coding",
      value: "0",
      note: "",
      occurred_at: nowAsDatetimeLocalValue(),
    },
  });

  useEffect(() => {
    if (open) {
      reset(
        activity
          ? {
              category: activity.category,
              value: String(activity.value),
              note: activity.note ?? "",
              occurred_at: toDatetimeLocalValue(activity.occurred_at),
            }
          : {
              category: "coding",
              value: "0",
              note: "",
              occurred_at: nowAsDatetimeLocalValue(),
            }
      );
    }
  }, [open, activity, reset]);

  const category = watch("category");

  const onSubmit = (values: ActivityFormValues) => {
    const input = {
      category: values.category,
      value: Number(values.value),
      note: values.note || null,
      occurred_at: fromDatetimeLocalValue(values.occurred_at),
    };

    const onError = (error: unknown) => {
      toast.error(error instanceof ApiError ? error.message : "Gagal menyimpan aktivitas");
    };
    const onSuccess = () => {
      toast.success(isEdit ? "Aktivitas berhasil diperbarui" : "Aktivitas berhasil dicatat");
      onOpenChange(false);
    };

    if (isEdit && activity) {
      updateMutation.mutate({ id: activity.id, input }, { onSuccess, onError });
    } else {
      createMutation.mutate(input, { onSuccess, onError });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit Aktivitas" : "Catat Aktivitas"}</DialogTitle>
        </DialogHeader>
        <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <Label>Kategori</Label>
            <Select value={category} onValueChange={(v) => setValue("category", v as ActivityFormValues["category"])}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Pilih kategori" />
              </SelectTrigger>
              <SelectContent>
                {activityCategories.map((c) => (
                  <SelectItem key={c} value={c}>
                    {categoryLabels[c]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.category && <p className="text-sm text-destructive">{errors.category.message}</p>}
          </div>

          <div className="space-y-2">
            <Label htmlFor="value">Value (menit / halaman / rupiah, dst.)</Label>
            <Input id="value" type="number" step="any" {...register("value")} />
            {errors.value && <p className="text-sm text-destructive">{errors.value.message}</p>}
          </div>

          <div className="space-y-2">
            <Label htmlFor="occurred_at">Waktu</Label>
            <Input id="occurred_at" type="datetime-local" {...register("occurred_at")} />
            {errors.occurred_at && (
              <p className="text-sm text-destructive">{errors.occurred_at.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="note">Catatan (opsional)</Label>
            <Textarea id="note" {...register("note")} />
          </div>

          <DialogFooter>
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? "Menyimpan..." : "Simpan"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

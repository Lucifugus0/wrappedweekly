"use client";

import { useState } from "react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { ApiError } from "@/lib/api/client";
import type { Activity } from "@/lib/api/types";
import { useActivities, useDeleteActivity } from "@/features/activity/use-activities";
import { ActivityFormDialog } from "@/features/activity/activity-form-dialog";
import { activityCategories, categoryLabels } from "@/features/activity/schemas";

export default function ActivitiesPage() {
  const [page, setPage] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingActivity, setEditingActivity] = useState<Activity | null>(null);
  const [categoryFilter, setCategoryFilter] = useState<string>("all");
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");

  const filter = {
    category: categoryFilter === "all" ? undefined : categoryFilter,
    from: fromDate ? new Date(fromDate).toISOString() : undefined,
    // "to" is exclusive on the backend, so bump to the start of the next day
    // to make the date picker feel inclusive of the selected end date.
    to: toDate ? new Date(new Date(toDate).getTime() + 24 * 60 * 60 * 1000).toISOString() : undefined,
  };

  const { data, isLoading, isError, refetch } = useActivities(page, 20, filter);
  const deleteMutation = useDeleteActivity();

  const resetToFirstPage = <T,>(setter: (v: T) => void) => (v: T) => {
    setter(v);
    setPage(1);
  };

  const openCreate = () => {
    setEditingActivity(null);
    setDialogOpen(true);
  };
  const openEdit = (activity: Activity) => {
    setEditingActivity(activity);
    setDialogOpen(true);
  };

  const handleDelete = (id: string) => {
    if (!confirm("Hapus aktivitas ini?")) return;
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success("Aktivitas dihapus"),
      onError: (error) =>
        toast.error(error instanceof ApiError ? error.message : "Gagal menghapus aktivitas"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Aktivitas</h1>
        <Button onClick={openCreate}>+ Catat Aktivitas</Button>
      </div>

      <Card>
        <CardContent className="flex flex-wrap items-end gap-4 py-4">
          <div className="space-y-1.5">
            <Label>Kategori</Label>
            <Select
              value={categoryFilter}
              onValueChange={(v) => resetToFirstPage(setCategoryFilter)(v ?? "all")}
            >
              <SelectTrigger className="h-9 w-40">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Semua kategori</SelectItem>
                {activityCategories.map((c) => (
                  <SelectItem key={c} value={c}>
                    {categoryLabels[c]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="from-date">Dari tanggal</Label>
            <Input
              id="from-date"
              type="date"
              className="h-9 w-40"
              value={fromDate}
              onChange={(e) => resetToFirstPage(setFromDate)(e.target.value)}
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="to-date">Sampai tanggal</Label>
            <Input
              id="to-date"
              type="date"
              className="h-9 w-40"
              value={toDate}
              onChange={(e) => resetToFirstPage(setToDate)(e.target.value)}
            />
          </div>
          {(categoryFilter !== "all" || fromDate || toDate) && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                setCategoryFilter("all");
                setFromDate("");
                setToDate("");
                setPage(1);
              }}
            >
              Reset filter
            </Button>
          )}
        </CardContent>
      </Card>

      {isLoading && (
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-16 w-full" />
          ))}
        </div>
      )}

      {isError && !isLoading && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
            <p className="text-muted-foreground">Gagal memuat aktivitas.</p>
            <Button variant="outline" onClick={() => refetch()}>
              Coba lagi
            </Button>
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && data && data.items.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
            <p className="text-muted-foreground">Belum ada aktivitas tercatat minggu ini.</p>
            <Button onClick={openCreate}>Catat aktivitas pertamamu</Button>
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && data && data.items.length > 0 && (
        <div className="space-y-2">
          {data.items.map((activity) => (
            <Card key={activity.id}>
              <CardContent className="flex items-center justify-between py-4">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">{categoryLabels[activity.category]}</Badge>
                    <span className="font-medium">{activity.value}</span>
                    <span className="text-sm text-muted-foreground">
                      {new Date(activity.occurred_at).toLocaleString("id-ID")}
                    </span>
                  </div>
                  {activity.note && <p className="text-sm text-muted-foreground">{activity.note}</p>}
                </div>
                <div className="flex gap-2">
                  <Button size="sm" variant="outline" onClick={() => openEdit(activity)}>
                    Edit
                  </Button>
                  <Button size="sm" variant="destructive" onClick={() => handleDelete(activity.id)}>
                    Hapus
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {data && data.total > data.size && (
        <div className="flex justify-center gap-2">
          <Button variant="outline" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
            Sebelumnya
          </Button>
          <Button
            variant="outline"
            disabled={page * data.size >= data.total}
            onClick={() => setPage((p) => p + 1)}
          >
            Berikutnya
          </Button>
        </div>
      )}

      <ActivityFormDialog open={dialogOpen} onOpenChange={setDialogOpen} activity={editingActivity} />
    </div>
  );
}

"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

import { useDashboardSummary } from "@/features/dashboard/use-dashboard";
import { DailyChart } from "@/features/dashboard/daily-chart";
import { categoryLabels } from "@/features/activity/schemas";

export default function DashboardPage() {
  const { data: stats, isLoading, isError, refetch } = useDashboardSummary();

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-1/3" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (isError || !stats) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
          <p className="text-muted-foreground">Gagal memuat ringkasan dashboard.</p>
          <Button variant="outline" onClick={() => refetch()}>
            Coba lagi
          </Button>
        </CardContent>
      </Card>
    );
  }

  const isEmpty = stats.total_activities === 0;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold">Dashboard</h1>
        <p className="text-sm text-muted-foreground">
          Minggu ini: {new Date(stats.week_start).toLocaleDateString("id-ID")} -{" "}
          {new Date(stats.week_end).toLocaleDateString("id-ID")}
        </p>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm text-muted-foreground">Total Aktivitas</CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-semibold">{stats.total_activities}</CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm text-muted-foreground">Total Value</CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-semibold">{stats.total_value}</CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm text-muted-foreground">Perubahan vs Minggu Lalu</CardTitle>
          </CardHeader>
          <CardContent className="text-2xl font-semibold">
            {stats.change_vs_prev_week_pct === null
              ? "N/A"
              : `${stats.change_vs_prev_week_pct > 0 ? "+" : ""}${stats.change_vs_prev_week_pct.toFixed(0)}%`}
          </CardContent>
        </Card>
      </div>

      {isEmpty ? (
        <Card>
          <CardContent className="py-10 text-center text-muted-foreground">
            Belum ada aktivitas minggu ini. Mulai catat di halaman Aktivitas.
          </CardContent>
        </Card>
      ) : (
        <>
          <Card>
            <CardHeader>
              <CardTitle>Aktivitas per Hari</CardTitle>
            </CardHeader>
            <CardContent>
              <DailyChart data={stats.daily_breakdown} />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Komposisi Kategori</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-2">
              {stats.totals_by_category.map((ct) => (
                <Badge key={ct.category} variant="secondary">
                  {categoryLabels[ct.category]}: {ct.total} ({ct.count}x)
                </Badge>
              ))}
              {stats.top_category && (
                <p className="mt-2 w-full text-sm text-muted-foreground">
                  Kategori paling aktif: <strong>{categoryLabels[stats.top_category]}</strong>
                </p>
              )}
              {stats.most_productive_day && (
                <p className="w-full text-sm text-muted-foreground">
                  Hari paling produktif:{" "}
                  <strong>{new Date(stats.most_productive_day).toLocaleDateString("id-ID")}</strong>
                </p>
              )}
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}

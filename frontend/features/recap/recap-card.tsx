import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { categoryLabels } from "@/features/activity/schemas";
import type { Recap } from "@/lib/api/types";

// Plain server-renderable component (no "use client", no hooks) so it can be
// reused by both the authenticated recap list and the public SSR share page.
export function RecapCard({ recap }: { recap: Recap }) {
  const { stats } = recap;

  return (
    <Card>
      <CardHeader>
        <CardTitle>
          Recap {new Date(recap.week_start).toLocaleDateString("id-ID")} -{" "}
          {new Date(stats.week_end).toLocaleDateString("id-ID")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="whitespace-pre-line text-sm leading-relaxed">{recap.narrative}</p>

        <div className="flex flex-wrap gap-2">
          {stats.totals_by_category.map((ct) => (
            <Badge key={ct.category} variant="secondary">
              {categoryLabels[ct.category]}: {ct.total}
            </Badge>
          ))}
        </div>

        <dl className="grid grid-cols-2 gap-3 text-sm sm:grid-cols-4">
          <div>
            <dt className="text-muted-foreground">Total Aktivitas</dt>
            <dd className="font-medium">{stats.total_activities}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">Total Value</dt>
            <dd className="font-medium">{stats.total_value}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">Kategori Teratas</dt>
            <dd className="font-medium">
              {stats.top_category ? categoryLabels[stats.top_category] : "-"}
            </dd>
          </div>
          <div>
            <dt className="text-muted-foreground">vs Minggu Lalu</dt>
            <dd className="font-medium">
              {stats.change_vs_prev_week_pct === null
                ? "N/A"
                : `${stats.change_vs_prev_week_pct > 0 ? "+" : ""}${stats.change_vs_prev_week_pct.toFixed(0)}%`}
            </dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  );
}

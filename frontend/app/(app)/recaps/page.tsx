"use client";

import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

import { ApiError } from "@/lib/api/client";
import { useGenerateRecap, useRecaps } from "@/features/recap/use-recaps";
import { RecapCard } from "@/features/recap/recap-card";

export default function RecapsPage() {
  const { data: recaps, isLoading, isError, refetch } = useRecaps();
  const generateMutation = useGenerateRecap();

  const handleGenerate = () => {
    generateMutation.mutate(undefined, {
      onSuccess: () => toast.success("Recap minggu ini berhasil dibuat"),
      onError: (error) =>
        toast.error(error instanceof ApiError ? error.message : "Gagal membuat recap"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Weekly Recap</h1>
        <Button onClick={handleGenerate} disabled={generateMutation.isPending}>
          {generateMutation.isPending ? "Membuat recap..." : "Generate Recap Minggu Ini"}
        </Button>
      </div>

      {isLoading && (
        <div className="space-y-2">
          <Skeleton className="h-40 w-full" />
        </div>
      )}

      {isError && !isLoading && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
            <p className="text-muted-foreground">Gagal memuat daftar recap.</p>
            <Button variant="outline" onClick={() => refetch()}>
              Coba lagi
            </Button>
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && recaps && recaps.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-10 text-center">
            <p className="text-muted-foreground">
              Belum ada recap. Catat aktivitas lalu klik &quot;Generate Recap Minggu Ini&quot;.
            </p>
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && recaps && recaps.length > 0 && (
        <div className="space-y-4">
          {recaps.map((recap) => (
            <div key={recap.id} className="space-y-2">
              <RecapCard recap={recap} />
              <div className="flex items-center gap-2 text-sm">
                <span className="text-muted-foreground">Link publik:</span>
                <code className="rounded bg-muted px-2 py-1">/w/{recap.slug}</code>
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => {
                    navigator.clipboard.writeText(`${window.location.origin}/w/${recap.slug}`);
                    toast.success("Link disalin");
                  }}
                >
                  Salin Link
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

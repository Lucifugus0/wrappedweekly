import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { RecapCard } from "@/features/recap/recap-card";
import { serverGet } from "@/lib/api/server";
import type { Recap } from "@/lib/api/types";

type Props = {
  params: Promise<{ slug: string }>;
};

async function getRecap(slug: string): Promise<Recap | null> {
  return serverGet<Recap>(`/recaps/public/${slug}`);
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const recap = await getRecap(slug);

  if (!recap) {
    return { title: "Recap tidak ditemukan — Wrapped Weekly" };
  }

  const title = `Wrapped Weekly — ${new Date(recap.week_start).toLocaleDateString("id-ID")}`;
  const description = recap.narrative.slice(0, 160);

  return {
    title,
    description,
    openGraph: {
      title,
      description,
      type: "website",
      // Dynamic per-recap OG image is out of scope for this MVP (see
      // docs/PROJECT_PLAN.md "Yang Sengaja Tidak Dikerjakan") — a single
      // static share image is used for all recaps instead.
      images: ["/og-default.png"],
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
    },
  };
}

// Server Component — this page is rendered on the server on every request
// (cache: "no-store" in serverGet) so crawlers and link-preview bots receive
// fully-formed HTML with the OG meta tags already in <head>, without needing
// to execute client-side JS.
export default async function PublicRecapPage({ params }: Props) {
  const { slug } = await params;
  const recap = await getRecap(slug);

  if (!recap) {
    notFound();
  }

  return (
    <div className="mx-auto min-h-screen max-w-2xl space-y-6 p-4 py-10">
      <div className="text-center">
        <p className="text-sm text-muted-foreground">Wrapped Weekly — Recap Publik</p>
      </div>
      <RecapCard recap={recap} />
      <p className="text-center text-xs text-muted-foreground">
        Buat recap mingguanmu sendiri di Wrapped Weekly.
      </p>
    </div>
  );
}

"use client";

import dynamic from "next/dynamic";
import { Button } from "@/components/ui/button";

export const ThemeToggleLoader = dynamic(
  () => import("./theme-toggle").then((mod) => mod.ThemeToggle),
  {
    ssr: false,
    loading: () => <Button variant="ghost" size="icon" disabled className="size-9" />,
  }
);

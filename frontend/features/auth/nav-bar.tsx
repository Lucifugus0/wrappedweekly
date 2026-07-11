"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { ThemeToggleLoader } from "@/components/theme-toggle-loader";
import { authApi } from "@/lib/api/endpoints";
import { useCurrentUser } from "./use-auth";

const links = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/activities", label: "Aktivitas" },
  { href: "/recaps", label: "Recap" },
];

export function NavBar() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data: user } = useCurrentUser();

  const logoutMutation = useMutation({
    mutationFn: authApi.logout,
    onSuccess: () => {
      queryClient.setQueryData(["auth", "me"], null);
      router.push("/login");
    },
  });

  return (
    <nav className="border-b">
      <div className="mx-auto flex max-w-4xl items-center justify-between p-4">
        <div className="flex items-center gap-6">
          <span className="font-semibold">Wrapped Weekly</span>
          <div className="flex gap-4 text-sm text-muted-foreground">
            {links.map((link) => (
              <Link key={link.href} href={link.href} className="hover:text-foreground">
                {link.label}
              </Link>
            ))}
          </div>
        </div>
        <div className="flex items-center gap-3">
          {user && <span className="text-sm text-muted-foreground">{user.name}</span>}
          <ThemeToggleLoader />
          <Button
            variant="outline"
            size="sm"
            onClick={() => logoutMutation.mutate()}
            disabled={logoutMutation.isPending}
          >
            Keluar
          </Button>
        </div>
      </div>
    </nav>
  );
}

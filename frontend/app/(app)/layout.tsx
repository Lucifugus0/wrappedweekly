import { AuthGuard } from "@/features/auth/auth-guard";
import { NavBar } from "@/features/auth/nav-bar";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <NavBar />
      <main className="mx-auto w-full max-w-4xl flex-1 p-4">{children}</main>
    </AuthGuard>
  );
}

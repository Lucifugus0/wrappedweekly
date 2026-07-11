"use client";

import { useQuery } from "@tanstack/react-query";
import { authApi } from "@/lib/api/endpoints";
import { ApiError } from "@/lib/api/client";

export function useCurrentUser() {
  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: authApi.me,
    retry: false,
    // 401 just means "not logged in" — not a real error to surface to the user.
    throwOnError: (error) => !(error instanceof ApiError && error.status === 401),
  });
}

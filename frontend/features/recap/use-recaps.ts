"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { recapApi } from "@/lib/api/endpoints";

export function useRecaps() {
  return useQuery({
    queryKey: ["recaps"],
    queryFn: recapApi.list,
  });
}

export function useGenerateRecap() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => recapApi.generate(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["recaps"] });
    },
  });
}

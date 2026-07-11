"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { activityApi, type ActivityInput, type ActivityListFilter } from "@/lib/api/endpoints";

export function useActivities(page: number, size = 20, filter: ActivityListFilter = {}) {
  return useQuery({
    queryKey: ["activities", page, size, filter],
    queryFn: () => activityApi.list(page, size, filter),
  });
}

export function useCreateActivity() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: ActivityInput) => activityApi.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["activities"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });
}

export function useUpdateActivity() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: ActivityInput }) =>
      activityApi.update(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["activities"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });
}

export function useDeleteActivity() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => activityApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["activities"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });
}

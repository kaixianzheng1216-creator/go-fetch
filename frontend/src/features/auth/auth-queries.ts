import { useQuery } from "@tanstack/react-query"

import { ApiError, api, type User } from "@/lib/api"

export const authQueryKeys = {
  me: ["auth", "me"] as const,
}

async function getCurrentUser(): Promise<User | null> {
  try {
    return await api.me()
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null
    }
    throw error
  }
}

export function useCurrentUser() {
  return useQuery({
    queryKey: authQueryKeys.me,
    queryFn: getCurrentUser,
  })
}

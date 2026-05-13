import { useQuery } from "@tanstack/react-query"

import { api } from "@/lib/api"

export const websiteQueryKeys = {
  all: ["websites"] as const,
  lists: () => [...websiteQueryKeys.all, "list"] as const,
  detail: (websiteID: string) => [...websiteQueryKeys.all, "detail", websiteID] as const,
}

export function useWebsitesQuery() {
  return useQuery({
    queryKey: websiteQueryKeys.lists(),
    queryFn: api.websites,
  })
}

export function useWebsiteQuery(websiteID: string | undefined) {
  return useQuery({
    queryKey: websiteID ? websiteQueryKeys.detail(websiteID) : websiteQueryKeys.detail(""),
    queryFn: () => api.website(websiteID!),
    enabled: Boolean(websiteID),
  })
}

export function useWebsiteDashboardQuery(websiteID: string | undefined) {
  return useQuery({
    queryKey: websiteID
      ? [...websiteQueryKeys.detail(websiteID), "dashboard"]
      : ["websites", "dashboard", ""],
    queryFn: async () => {
      const [website, stats, pageviews, paths, referrers, browsers, events] = await Promise.all([
        api.website(websiteID!),
        api.websiteStats(websiteID!),
        api.websitePageviews(websiteID!),
        api.websiteMetrics(websiteID!, "path"),
        api.websiteMetrics(websiteID!, "referrer"),
        api.websiteMetrics(websiteID!, "browser"),
        api.websiteMetrics(websiteID!, "event"),
      ])

      return {
        website,
        stats,
        pageviews,
        metrics: { paths, referrers, browsers, events },
      }
    },
    enabled: Boolean(websiteID),
  })
}

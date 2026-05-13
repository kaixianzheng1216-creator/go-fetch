import createClient from "openapi-fetch"

import type { components, operations, paths } from "@/lib/api-types"

export type User = components["schemas"]["User"]
export type WebsiteRequest = components["schemas"]["WebsiteRequest"]
export type PageviewPoint = components["schemas"]["PageviewPoint"]
export type MetricRow = components["schemas"]["MetricRow"]
export type MetricType = operations["websiteMetrics"]["parameters"]["query"]["type"]

type ErrorResponse = {
  detail?: string
  title?: string
  message?: string
  error?: { message?: string }
}
type ApiResult<T> = {
  data?: T
  error?: ErrorResponse
  response: Response
}

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

const client = createClient<paths>({ credentials: "same-origin" })

async function unwrap<T>(request: Promise<ApiResult<T>>): Promise<T> {
  const { data, error, response } = await request

  if (error) {
    throw new ApiError(
      response.status,
      error.detail ?? error.error?.message ?? error.message ?? error.title ?? response.statusText,
    )
  }
  if (data === undefined) {
    throw new ApiError(response.status, response.statusText)
  }
  return data
}

export const api = {
  me: () => unwrap(client.GET("/api/me")),
  login: (username: string, password: string) =>
    unwrap(
      client.POST("/api/login", {
        body: { username, password },
      }),
    ),
  logout: () => unwrap(client.POST("/api/logout")),
  websites: () => unwrap(client.GET("/api/websites")),
  createWebsite: (input: WebsiteRequest) =>
    unwrap(
      client.POST("/api/websites", {
        body: input,
      }),
    ),
  website: (websiteID: string) =>
    unwrap(
      client.GET("/api/websites/{websiteID}", {
        params: { path: { websiteID } },
      }),
    ),
  updateWebsite: (websiteID: string, input: WebsiteRequest) =>
    unwrap(
      client.PATCH("/api/websites/{websiteID}", {
        params: { path: { websiteID } },
        body: input,
      }),
    ),
  deleteWebsite: (websiteID: string) =>
    unwrap(
      client.DELETE("/api/websites/{websiteID}", {
        params: { path: { websiteID } },
      }),
    ),
  websiteStats: (websiteID: string) =>
    unwrap(
      client.GET("/api/websites/{websiteID}/stats", {
        params: { path: { websiteID } },
      }),
    ),
  websitePageviews: (websiteID: string) =>
    unwrap(
      client.GET("/api/websites/{websiteID}/pageviews", {
        params: { path: { websiteID }, query: { unit: "hour" } },
      }),
    ),
  websiteMetrics: (websiteID: string, type: MetricType) =>
    unwrap(
      client.GET("/api/websites/{websiteID}/metrics", {
        params: { path: { websiteID }, query: { type, limit: 8 } },
      }),
    ),
}

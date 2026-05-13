import { useMemo, useState } from "react"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import {
  Activity,
  ArrowLeft,
  BarChart3,
  Clock3,
  Copy,
  MousePointer2,
  Pencil,
  RefreshCw,
  Trash2,
  Users,
} from "lucide-react"
import { useNavigate, useParams } from "react-router"

import { ErrorBanner } from "@/components/feedback/error-banner"
import { PageHeader } from "@/components/layout/page-header"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Textarea } from "@/components/ui/textarea"
import { MetricList } from "@/features/websites/components/metric-list"
import { PageviewChart } from "@/features/websites/components/pageview-chart"
import { DetailSkeleton } from "@/features/websites/components/skeletons"
import { StatCard } from "@/features/websites/components/stat-card"
import { useWebsiteDashboardQuery, websiteQueryKeys } from "@/features/websites/website-queries"
import { api } from "@/lib/api"

export function WebsiteDetailPage() {
  const navigate = useNavigate()
  const { websiteID } = useParams<{ websiteID: string }>()
  const queryClient = useQueryClient()
  const dashboardQuery = useWebsiteDashboardQuery(websiteID)
  const [copied, setCopied] = useState(false)
  const deleteWebsite = useMutation({
    mutationFn: () => api.deleteWebsite(websiteID!),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: websiteQueryKeys.lists() })
      navigate("/websites")
    },
  })

  const website = dashboardQuery.data?.website
  const stats = dashboardQuery.data?.stats
  const pageviews = dashboardQuery.data?.pageviews ?? []
  const metrics = dashboardQuery.data?.metrics ?? {
    paths: [],
    referrers: [],
    browsers: [],
    events: [],
  }

  const snippet = useMemo(() => {
    if (!website) return ""
    return `<script defer src="${window.location.origin}/script.js" data-website-id="${website.id}"></script>`
  }, [website])

  async function copySnippet() {
    if (!snippet) return

    await navigator.clipboard.writeText(snippet)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1400)
  }

  if (!websiteID) {
    return <ErrorBanner message="Website id is missing." />
  }

  if (dashboardQuery.isPending) return <DetailSkeleton />

  if (dashboardQuery.isError || !website || !stats) {
    return (
      <div className="grid gap-6">
        <PageHeader
          eyebrow="Overview"
          title="Website"
          description="The selected website could not be loaded."
          actions={
            <Button variant="outline" onClick={() => navigate("/websites")}>
              <ArrowLeft />
              Back
            </Button>
          }
        />
        <ErrorBanner
          message={
            dashboardQuery.error instanceof Error
              ? dashboardQuery.error.message
              : "Website not found"
          }
        />
      </div>
    )
  }

  return (
    <div className="grid gap-6">
      <PageHeader
        eyebrow={website.domain || "No domain set"}
        title={website.name}
        description={`Website id: ${website.id}`}
        actions={
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" onClick={() => navigate(`/websites/${website.id}/edit`)}>
              <Pencil />
              Edit
            </Button>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="destructive">
                  <Trash2 />
                  Delete
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete website?</AlertDialogTitle>
                  <AlertDialogDescription>
                    {website.name} will be removed from this dashboard. Existing analytics rows are
                    not shown after deletion.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    disabled={deleteWebsite.isPending}
                    onClick={() => deleteWebsite.mutate()}
                  >
                    Delete
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        }
      />

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
        <StatCard label="Pageviews" value={stats.pageviews} icon={MousePointer2} tone="primary" />
        <StatCard label="Visitors" value={stats.visitors} icon={Users} tone="success" />
        <StatCard label="Visits" value={stats.visits} icon={Activity} tone="neutral" />
        <StatCard label="Bounces" value={stats.bounces} icon={BarChart3} tone="danger" />
        <StatCard
          label="Avg visit"
          value={`${stats.avgVisitSeconds}s`}
          icon={Clock3}
          tone="warning"
        />
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_360px]">
        <Card>
          <CardHeader className="flex-row items-center justify-between space-y-0">
            <div>
              <CardTitle>Pageviews</CardTitle>
              <CardDescription>Hourly pageviews and visitors.</CardDescription>
            </div>
            <Button variant="outline" size="sm" onClick={() => void dashboardQuery.refetch()}>
              <RefreshCw className={dashboardQuery.isFetching ? "animate-spin" : ""} />
              Refresh
            </Button>
          </CardHeader>
          <CardContent>
            <PageviewChart points={pageviews} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Tracking code</CardTitle>
            <CardDescription>Paste this into the tracked website.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-3">
            <Textarea value={snippet} readOnly className="min-h-28 resize-none font-mono text-xs" />
            <Button variant="outline" onClick={() => void copySnippet()}>
              <Copy />
              {copied ? "Copied" : "Copy snippet"}
            </Button>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Breakdowns</CardTitle>
          <CardDescription>Top dimensions for the selected website.</CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="paths">
            <TabsList className="grid w-full grid-cols-2 md:grid-cols-4">
              <TabsTrigger value="paths">Paths</TabsTrigger>
              <TabsTrigger value="referrers">Referrers</TabsTrigger>
              <TabsTrigger value="browsers">Browsers</TabsTrigger>
              <TabsTrigger value="events">Events</TabsTrigger>
            </TabsList>
            <TabsContent value="paths">
              <MetricList rows={metrics.paths || []} />
            </TabsContent>
            <TabsContent value="referrers">
              <MetricList rows={metrics.referrers || []} />
            </TabsContent>
            <TabsContent value="browsers">
              <MetricList rows={metrics.browsers || []} />
            </TabsContent>
            <TabsContent value="events">
              <MetricList rows={metrics.events || []} />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  )
}

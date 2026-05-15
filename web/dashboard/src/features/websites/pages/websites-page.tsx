import { MoreHorizontal, Plus, RefreshCw } from "lucide-react"
import { useNavigate } from "react-router"

import { ErrorBanner } from "@/components/feedback/error-banner"
import { PageHeader } from "@/components/layout/page-header"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { TableSkeleton } from "@/features/websites/components/skeletons"
import { useWebsitesQuery } from "@/features/websites/website-queries"
import { formatDate } from "@/lib/format"

export function WebsitesPage() {
  const navigate = useNavigate()
  const websitesQuery = useWebsitesQuery()
  const websites = websitesQuery.data ?? []

  return (
    <div className="grid gap-6">
      <PageHeader
        eyebrow="Inventory"
        title="Websites"
        description="Create tracking targets and open their analytics workspace."
        actions={
          <Button onClick={() => navigate("/websites/new")}>
            <Plus />
            New website
          </Button>
        }
      />

      <Card>
        <CardHeader className="flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle>Configured websites</CardTitle>
            <CardDescription>{websites.length} total</CardDescription>
          </div>
          <Button variant="outline" size="sm" onClick={() => void websitesQuery.refetch()}>
            <RefreshCw className={websitesQuery.isFetching ? "animate-spin" : ""} />
            Refresh
          </Button>
        </CardHeader>
        <CardContent>
          {websitesQuery.isError && (
            <ErrorBanner
              message={
                websitesQuery.error instanceof Error
                  ? websitesQuery.error.message
                  : "加载站点列表失败"
              }
            />
          )}
          {websitesQuery.isPending ? (
            <TableSkeleton />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Domain</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-14"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {websites.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} className="text-muted-foreground h-28 text-center">
                      No websites yet.
                    </TableCell>
                  </TableRow>
                ) : (
                  websites.map((website) => (
                    <TableRow key={website.id}>
                      <TableCell>
                        <button
                          className="hover:text-primary font-medium"
                          onClick={() => navigate(`/websites/${website.id}`)}
                        >
                          {website.name}
                        </button>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {website.domain || "No domain set"}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {formatDate(website.createdAt)}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              variant="ghost"
                              size="icon"
                              aria-label={`Actions for ${website.name}`}
                            >
                              <MoreHorizontal />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => navigate(`/websites/${website.id}`)}>
                              Open dashboard
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => navigate(`/websites/${website.id}/edit`)}
                            >
                              Edit website
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

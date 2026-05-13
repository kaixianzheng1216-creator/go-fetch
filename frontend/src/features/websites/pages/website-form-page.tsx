import { useEffect } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useForm } from "react-hook-form"
import { useNavigate, useParams } from "react-router"
import { ArrowLeft, Loader2 } from "lucide-react"

import { ErrorBanner } from "@/components/feedback/error-banner"
import { PageHeader } from "@/components/layout/page-header"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Skeleton } from "@/components/ui/skeleton"
import { useWebsiteQuery, websiteQueryKeys } from "@/features/websites/website-queries"
import { websiteSchema, type WebsiteFormValues } from "@/features/websites/website-schema"
import { api } from "@/lib/api"

export function WebsiteFormPage({ mode }: { mode: "create" | "edit" }) {
  const navigate = useNavigate()
  const { websiteID } = useParams<{ websiteID: string }>()
  const queryClient = useQueryClient()
  const websiteQuery = useWebsiteQuery(mode === "edit" ? websiteID : undefined)
  const backPath = mode === "edit" && websiteID ? `/websites/${websiteID}` : "/websites"
  const form = useForm<WebsiteFormValues>({
    resolver: zodResolver(websiteSchema),
    defaultValues: {
      name: "",
      domain: "",
    },
  })

  useEffect(() => {
    if (!websiteQuery.data) return
    form.reset({
      name: websiteQuery.data.name,
      domain: websiteQuery.data.domain,
    })
  }, [form, websiteQuery.data])

  const saveWebsite = useMutation({
    mutationFn: (values: WebsiteFormValues) => {
      const input = { name: values.name.trim(), domain: values.domain.trim() }
      return mode === "create" ? api.createWebsite(input) : api.updateWebsite(websiteID!, input)
    },
    onSuccess: async (website) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: websiteQueryKeys.lists() }),
        queryClient.invalidateQueries({ queryKey: websiteQueryKeys.detail(website.id) }),
      ])
      navigate(`/websites/${website.id}`)
    },
  })

  if (mode === "edit" && !websiteID) {
    return <ErrorBanner message="Website id is missing." />
  }

  return (
    <div className="grid gap-6">
      <PageHeader
        eyebrow="Configuration"
        title={mode === "create" ? "New website" : "Edit website"}
        description="Set the dashboard name and optional production domain."
        actions={
          <Button variant="outline" onClick={() => navigate(backPath)}>
            <ArrowLeft />
            Back
          </Button>
        }
      />

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>Website details</CardTitle>
          <CardDescription>
            The generated website id is used by the tracking script.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {websiteQuery.isPending && mode === "edit" ? (
            <div className="grid gap-4">
              <Skeleton className="h-9" />
              <Skeleton className="h-9" />
            </div>
          ) : (
            <form
              className="grid gap-4"
              onSubmit={form.handleSubmit((values) => saveWebsite.mutate(values))}
            >
              <div className="grid gap-2">
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  maxLength={100}
                  aria-invalid={Boolean(form.formState.errors.name)}
                  {...form.register("name")}
                />
                {form.formState.errors.name && (
                  <p className="text-destructive text-sm">{form.formState.errors.name.message}</p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="domain">Domain</Label>
                <Input
                  id="domain"
                  maxLength={500}
                  placeholder="example.com"
                  aria-invalid={Boolean(form.formState.errors.domain)}
                  {...form.register("domain")}
                />
                {form.formState.errors.domain && (
                  <p className="text-destructive text-sm">{form.formState.errors.domain.message}</p>
                )}
              </div>
              {websiteQuery.isError && (
                <ErrorBanner
                  message={
                    websiteQuery.error instanceof Error
                      ? websiteQuery.error.message
                      : "Failed to load website"
                  }
                />
              )}
              {saveWebsite.error && (
                <ErrorBanner
                  message={
                    saveWebsite.error instanceof Error
                      ? saveWebsite.error.message
                      : "Failed to save website"
                  }
                />
              )}
              <div className="flex justify-end gap-2">
                <Button type="button" variant="outline" onClick={() => navigate(backPath)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={saveWebsite.isPending}>
                  {saveWebsite.isPending && <Loader2 className="animate-spin" />}
                  Save
                </Button>
              </div>
            </form>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

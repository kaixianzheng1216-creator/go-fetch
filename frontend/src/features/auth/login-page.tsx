import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useForm } from "react-hook-form"
import { Navigate, useLocation, useNavigate } from "react-router"
import { BarChart3, Loader2 } from "lucide-react"

import { BootScreen } from "@/components/feedback/boot-screen"
import { ErrorBanner } from "@/components/feedback/error-banner"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { authQueryKeys, useCurrentUser } from "@/features/auth/auth-queries"
import { loginSchema, type LoginFormValues } from "@/features/auth/login-schema"
import { api } from "@/lib/api"

type LoginLocationState = {
  from?: {
    pathname?: string
  }
}

export function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const queryClient = useQueryClient()
  const currentUser = useCurrentUser()
  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  })
  const from = (location.state as LoginLocationState | null)?.from?.pathname || "/websites"

  const login = useMutation({
    mutationFn: (values: LoginFormValues) => api.login(values.username.trim(), values.password),
    onSuccess: (result) => {
      queryClient.setQueryData(authQueryKeys.me, result.user)
      navigate(from, { replace: true })
    },
  })

  if (currentUser.isPending) return <BootScreen />
  if (currentUser.isError) throw currentUser.error
  if (currentUser.data) return <Navigate to="/websites" replace />

  return (
    <main className="bg-background grid min-h-screen lg:grid-cols-[minmax(0,1fr)_480px]">
      <section className="hidden min-h-screen flex-col justify-between border-r bg-neutral-950 px-10 py-8 text-neutral-50 lg:flex">
        <div className="flex items-center gap-3">
          <div className="grid size-10 place-items-center rounded-lg bg-neutral-50 text-neutral-950">
            <BarChart3 className="size-5" />
          </div>
          <div>
            <div className="font-semibold">go-fetch</div>
            <div className="text-xs text-neutral-400">analytics MVP</div>
          </div>
        </div>

        <div className="max-w-3xl">
          <Badge variant="secondary" className="mb-6 bg-neutral-800 text-neutral-100">
            Self-hosted analytics
          </Badge>
          <h1 className="text-5xl leading-none font-semibold">
            Small stack, useful signal, clear ownership.
          </h1>
          <p className="mt-5 max-w-xl text-sm leading-6 text-neutral-400">
            Track pageviews, visitors, referrals, browsers, and custom events with a Go backend and
            a React dashboard.
          </p>
        </div>

        <div className="grid grid-cols-3 gap-3 text-sm text-neutral-400">
          <div className="rounded-lg border border-neutral-800 p-4">Pageviews</div>
          <div className="rounded-lg border border-neutral-800 p-4">Visitors</div>
          <div className="rounded-lg border border-neutral-800 p-4">Events</div>
        </div>
      </section>

      <section className="grid min-h-screen place-items-center p-6">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-2xl">Sign in</CardTitle>
            <CardDescription>Use the administrator account.</CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className="grid gap-4"
              onSubmit={form.handleSubmit((values) => login.mutate(values))}
            >
              <div className="grid gap-2">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  autoComplete="username"
                  aria-invalid={Boolean(form.formState.errors.username)}
                  {...form.register("username")}
                />
                {form.formState.errors.username && (
                  <p className="text-destructive text-sm">
                    {form.formState.errors.username.message}
                  </p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  autoComplete="current-password"
                  aria-invalid={Boolean(form.formState.errors.password)}
                  {...form.register("password")}
                />
                {form.formState.errors.password && (
                  <p className="text-destructive text-sm">
                    {form.formState.errors.password.message}
                  </p>
                )}
              </div>
              {login.error && (
                <ErrorBanner
                  message={login.error instanceof Error ? login.error.message : "Login failed"}
                />
              )}
              <Button type="submit" disabled={login.isPending}>
                {login.isPending && <Loader2 className="animate-spin" />}
                Sign in
              </Button>
            </form>
          </CardContent>
        </Card>
      </section>
    </main>
  )
}

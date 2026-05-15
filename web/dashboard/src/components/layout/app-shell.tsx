import type { ReactNode } from "react"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { BarChart3, Globe2, LogOut } from "lucide-react"
import { useLocation, useNavigate } from "react-router"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { authQueryKeys } from "@/features/auth/auth-queries"
import { api, type User } from "@/lib/api"

export function AppShell({ user, children }: { user: User; children: ReactNode }) {
  const location = useLocation()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const logout = useMutation({
    mutationFn: api.logout,
    onSettled: () => {
      queryClient.removeQueries()
      queryClient.setQueryData(authQueryKeys.me, null)
      navigate("/login", { replace: true })
    },
  })

  return (
    <main className="bg-background min-h-screen">
      <header className="bg-background/95 sticky top-0 z-40 border-b backdrop-blur">
        <div className="mx-auto flex h-16 w-full max-w-7xl items-center justify-between px-4 sm:px-6">
          <button
            className="flex items-center gap-3 text-left"
            onClick={() => navigate("/websites")}
          >
            <span className="bg-primary text-primary-foreground grid size-9 place-items-center rounded-lg">
              <BarChart3 className="size-4" />
            </span>
            <span>
              <span className="block text-sm leading-none font-semibold">go-fetch</span>
              <span className="text-muted-foreground block text-xs">analytics</span>
            </span>
          </button>

          <nav className="flex items-center gap-2">
            <Button
              variant={location.pathname === "/websites" ? "secondary" : "ghost"}
              size="sm"
              onClick={() => navigate("/websites")}
            >
              <Globe2 />
              Websites
            </Button>
            <Separator className="hidden h-6 w-px sm:block" />
            <Badge variant="outline" className="hidden sm:inline-flex">
              {user.username}
            </Badge>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => logout.mutate()}
                  disabled={logout.isPending}
                  aria-label="Log out"
                >
                  <LogOut />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Log out</TooltipContent>
            </Tooltip>
          </nav>
        </div>
      </header>

      <div className="mx-auto w-full max-w-7xl px-4 py-6 sm:px-6">{children}</div>
    </main>
  )
}

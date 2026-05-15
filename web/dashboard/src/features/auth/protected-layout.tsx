import { Navigate, Outlet, useLocation } from "react-router"

import { BootScreen } from "@/components/feedback/boot-screen"
import { AppShell } from "@/components/layout/app-shell"
import { useCurrentUser } from "@/features/auth/auth-queries"

export function ProtectedLayout() {
  const location = useLocation()
  const currentUser = useCurrentUser()

  if (currentUser.isPending) return <BootScreen />
  if (currentUser.isError) throw currentUser.error

  if (!currentUser.data) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return (
    <AppShell user={currentUser.data}>
      <Outlet />
    </AppShell>
  )
}

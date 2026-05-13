import { lazy, Suspense } from "react"
import { Navigate, Route, Routes } from "react-router"

import { BootScreen } from "@/components/feedback/boot-screen"
import { ProtectedLayout } from "@/features/auth/protected-layout"

const LoginPage = lazy(() =>
  import("@/features/auth/login-page").then((module) => ({ default: module.LoginPage })),
)
const WebsitesPage = lazy(() =>
  import("@/features/websites/pages/websites-page").then((module) => ({
    default: module.WebsitesPage,
  })),
)
const WebsiteFormPage = lazy(() =>
  import("@/features/websites/pages/website-form-page").then((module) => ({
    default: module.WebsiteFormPage,
  })),
)
const WebsiteDetailPage = lazy(() =>
  import("@/features/websites/pages/website-detail-page").then((module) => ({
    default: module.WebsiteDetailPage,
  })),
)

export function App() {
  return (
    <Suspense fallback={<BootScreen />}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<ProtectedLayout />}>
          <Route path="/" element={<Navigate to="/websites" replace />} />
          <Route path="/websites" element={<WebsitesPage />} />
          <Route path="/websites/new" element={<WebsiteFormPage mode="create" />} />
          <Route path="/websites/:websiteID" element={<WebsiteDetailPage />} />
          <Route path="/websites/:websiteID/edit" element={<WebsiteFormPage mode="edit" />} />
        </Route>
        <Route path="*" element={<Navigate to="/websites" replace />} />
      </Routes>
    </Suspense>
  )
}

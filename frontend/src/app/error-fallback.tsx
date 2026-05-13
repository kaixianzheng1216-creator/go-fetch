import type { FallbackProps } from "react-error-boundary"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export function ErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  return (
    <main className="bg-background grid min-h-screen place-items-center p-6">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Something went wrong</CardTitle>
          <CardDescription>The dashboard hit an unexpected runtime error.</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <pre className="bg-muted text-muted-foreground max-h-48 overflow-auto rounded-md p-3 text-xs">
            {error instanceof Error ? error.message : "Unknown error"}
          </pre>
          <Button onClick={resetErrorBoundary}>Reload view</Button>
        </CardContent>
      </Card>
    </main>
  )
}

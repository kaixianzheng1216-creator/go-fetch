import type { FallbackProps } from "react-error-boundary"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export function ErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  return (
    <main className="bg-background grid min-h-screen place-items-center p-6">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>页面出错了</CardTitle>
          <CardDescription>控制台运行时发生了未预期的错误。</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <pre className="bg-muted text-muted-foreground max-h-48 overflow-auto rounded-md p-3 text-xs">
            {error instanceof Error ? error.message : "未知错误"}
          </pre>
          <Button onClick={resetErrorBoundary}>重新加载</Button>
        </CardContent>
      </Card>
    </main>
  )
}

import { Loader2 } from "lucide-react"

export function BootScreen() {
  return (
    <main className="bg-background grid min-h-screen place-items-center">
      <div className="bg-card text-muted-foreground flex items-center gap-3 rounded-lg border px-4 py-3 text-sm shadow-sm">
        <Loader2 className="text-primary size-4 animate-spin" />
        Loading workspace
      </div>
    </main>
  )
}

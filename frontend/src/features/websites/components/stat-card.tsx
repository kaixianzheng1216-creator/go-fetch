import type { LucideIcon } from "lucide-react"

import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"

export function StatCard({
  label,
  value,
  icon: Icon,
  tone,
}: {
  label: string
  value: number | string
  icon: LucideIcon
  tone: "primary" | "success" | "warning" | "danger" | "neutral"
}) {
  return (
    <Card>
      <CardContent className="p-4">
        <div className="mb-4 flex items-center justify-between">
          <span className="text-muted-foreground text-xs font-semibold">{label}</span>
          <span
            className={cn(
              "grid size-8 place-items-center rounded-md",
              tone === "primary" && "bg-sky-500/10 text-sky-700",
              tone === "success" && "bg-emerald-500/10 text-emerald-700",
              tone === "warning" && "bg-amber-500/10 text-amber-700",
              tone === "danger" && "bg-red-500/10 text-red-700",
              tone === "neutral" && "bg-secondary text-secondary-foreground",
            )}
          >
            <Icon className="size-4" />
          </span>
        </div>
        <div className="text-3xl font-semibold">{value}</div>
      </CardContent>
    </Card>
  )
}

import type { MetricRow } from "@/lib/api"

export function MetricList({ rows }: { rows: MetricRow[] }) {
  const max = Math.max(1, ...rows.map((row) => row.views))

  if (rows.length === 0) {
    return (
      <div className="text-muted-foreground mt-4 rounded-md border border-dashed py-10 text-center text-sm">
        No data.
      </div>
    )
  }

  return (
    <div className="mt-4 grid gap-3">
      {rows.map((row) => (
        <div key={row.name} className="grid gap-1">
          <div className="flex items-center justify-between gap-3 text-sm">
            <span className="truncate font-medium">{row.name}</span>
            <span className="text-muted-foreground">{row.views}</span>
          </div>
          <div className="bg-muted h-2 rounded-full">
            <div
              className="bg-primary h-2 rounded-full"
              style={{
                width: `${Math.max(4, Math.round((row.views / max) * 100))}%`,
              }}
            />
          </div>
        </div>
      ))}
    </div>
  )
}

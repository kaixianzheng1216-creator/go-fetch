import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts"

import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart"
import type { PageviewPoint } from "@/lib/api"

const pageviewChartConfig = {
  views: {
    label: "Pageviews",
    color: "var(--chart-2)",
  },
  visitors: {
    label: "Visitors",
    color: "var(--chart-3)",
  },
} satisfies ChartConfig

export function PageviewChart({ points }: { points: PageviewPoint[] }) {
  if (points.length === 0) {
    return (
      <div className="text-muted-foreground grid h-72 place-items-center rounded-md border border-dashed text-sm">
        No pageviews yet.
      </div>
    )
  }

  return (
    <ChartContainer config={pageviewChartConfig} className="aspect-auto h-72 w-full">
      <BarChart accessibilityLayer data={points}>
        <CartesianGrid vertical={false} />
        <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={10} minTickGap={18} />
        <YAxis width={36} tickLine={false} axisLine={false} />
        <ChartTooltip content={<ChartTooltipContent />} />
        <Bar dataKey="views" fill="var(--color-views)" radius={4} />
        <Bar dataKey="visitors" fill="var(--color-visitors)" radius={4} />
      </BarChart>
    </ChartContainer>
  )
}

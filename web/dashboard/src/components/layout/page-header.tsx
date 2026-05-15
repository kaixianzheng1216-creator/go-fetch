import type { ReactNode } from "react"

export function PageHeader({
  eyebrow,
  title,
  description,
  actions,
}: {
  eyebrow: string
  title: string
  description: string
  actions?: ReactNode
}) {
  return (
    <div className="flex flex-col gap-4 border-b pb-5 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <div className="text-muted-foreground text-xs font-semibold uppercase">{eyebrow}</div>
        <h1 className="mt-2 text-3xl font-semibold">{title}</h1>
        <p className="text-muted-foreground mt-2 max-w-2xl text-sm">{description}</p>
      </div>
      {actions}
    </div>
  )
}

import { Skeleton } from "@/components/ui/skeleton"

export function DetailSkeleton() {
  return (
    <div className="grid gap-6">
      <div className="grid gap-3 border-b pb-5">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-10 w-72" />
        <Skeleton className="h-4 w-96 max-w-full" />
      </div>
      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
        {Array.from({ length: 5 }).map((_, index) => (
          <Skeleton key={index} className="h-32" />
        ))}
      </section>
      <Skeleton className="h-96" />
    </div>
  )
}

export function TableSkeleton() {
  return (
    <div className="grid gap-2">
      {Array.from({ length: 5 }).map((_, index) => (
        <Skeleton key={index} className="h-11" />
      ))}
    </div>
  )
}

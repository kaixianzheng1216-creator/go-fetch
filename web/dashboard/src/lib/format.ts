export function formatDate(value: string) {
  if (!value) return ""
  return new Intl.DateTimeFormat("zh-CN", { dateStyle: "medium" }).format(new Date(value))
}

import { z } from "zod"

export const websiteSchema = z.object({
  name: z.string().trim().min(1, "请输入网站名称").max(100, "网站名称不能超过 100 个字符"),
  domain: z.string().trim().max(500, "域名不能超过 500 个字符"),
})

export type WebsiteFormValues = z.infer<typeof websiteSchema>

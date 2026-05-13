import { z } from "zod"

export const websiteSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(100, "Name is too long"),
  domain: z.string().trim().max(500, "Domain is too long"),
})

export type WebsiteFormValues = z.infer<typeof websiteSchema>

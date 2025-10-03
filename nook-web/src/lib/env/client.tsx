import z from "zod";

export const clientEnvSchema = z.object({
  NEXT_PUBLIC_API_PUBLIC_BASE_URL: z.string().default("http://localhost:8080"),
  NEXT_PUBLIC_DOCS_LINK: z.string().default("http://localhost:4000"),
  NEXT_PUBLIC_XENON_LINK: z.string().default("https://xenon.bot"),
  NEXT_PUBLIC_GITHUB_LINK: z
    .string()
    .default("https://github.com/merlinfuchs/nook"),
  NEXT_PUBLIC_DISCORD_LINK: z.string().default("https://discord.gg"),
  NEXT_PUBLIC_CONTACT_EMAIL: z.string().default("contact@xenon.bot"),
  NEXT_PUBLIC_PADDLE_ENVIRONMENT: z
    .enum(["sandbox", "production"])
    .default("sandbox"),
  NEXT_PUBLIC_PADDLE_AUTH_TOKEN: z.string().default(""),
  NEXT_PUBLIC_HIGHLIGHT_PROJECT_ID: z.string().default(""),
});

export default clientEnvSchema.parse({
  NEXT_PUBLIC_API_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_API_PUBLIC_BASE_URL,
  NEXT_PUBLIC_DOCS_LINK: process.env.NEXT_PUBLIC_DOCS_LINK,
  NEXT_PUBLIC_XENON_LINK: process.env.NEXT_PUBLIC_XENON_LINK,
  NEXT_PUBLIC_DISCORD_LINK: process.env.NEXT_PUBLIC_DISCORD_LINK,
  NEXT_PUBLIC_CONTACT_EMAIL: process.env.NEXT_PUBLIC_CONTACT_EMAIL,
  NEXT_PUBLIC_PADDLE_ENVIRONMENT: process.env.NEXT_PUBLIC_PADDLE_ENVIRONMENT,
  NEXT_PUBLIC_PADDLE_AUTH_TOKEN: process.env.NEXT_PUBLIC_PADDLE_AUTH_TOKEN,
  NEXT_PUBLIC_HIGHLIGHT_PROJECT_ID:
    process.env.NEXT_PUBLIC_HIGHLIGHT_PROJECT_ID,
});

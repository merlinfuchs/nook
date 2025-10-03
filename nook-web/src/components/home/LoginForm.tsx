import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import discordLogo from "@/assets/discord.svg";
import Image from "next/image";
import env from "@/lib/env/client";

export default function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const loginUrl = env.NEXT_PUBLIC_API_PUBLIC_BASE_URL + "/v1/auth/login";

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            Login with your Discord account to access your dashboard.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form>
            <div className="grid gap-6">
              <div className="flex flex-col gap-4">
                <Button variant="outline" className="w-full" asChild>
                  <a href={loginUrl}>
                    <Image
                      src={discordLogo}
                      alt="Discord Logo"
                      className="size-4"
                      height={64}
                      width={64}
                    />
                    Login with Discord
                  </a>
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
        You will be redirected to the Discord login page.
      </div>
    </div>
  );
}

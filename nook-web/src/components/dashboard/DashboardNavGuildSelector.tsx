import { useIsMobile } from "@/hooks/use-mobile";
import { getGuildIconUrl } from "@/lib/discord/cdn";
import { useGuild, useGuilds } from "@/lib/hooks/api";
import { abbreviateName } from "@/lib/utils";
import { CogIcon } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useMemo } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { SidebarMenuButton, SidebarMenuItem } from "../ui/sidebar";

export default function DashboardNavGuildSelector() {
  const isMobile = useIsMobile();
  const router = useRouter();

  const guild = useGuild(undefined, (res) => {
    if (!res.success) {
      if (res.error.code === "unauthorized") {
        router.push("/login");
      } else if (
        res.error.code === "unknown_guild" ||
        res.error.code === "missing_access"
      ) {
        router.push("/guilds");
      }
    }
  });

  const guilds = useGuilds();

  const filteredGuilds = useMemo(() => {
    return guilds?.filter((g) => g.access && g.bot && g.id !== guild?.id);
  }, [guilds, guild]);

  const setGuild = useCallback(
    (guildId: string) => {
      const currentPath = window.location.pathname;
      const newPath = currentPath.replace(
        /\/dashboard\/[^\/]+/,
        `/dashboard/${guildId}`
      );
      router.push(newPath);
    },
    [router]
  );

  return (
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <SidebarMenuButton size="lg" asChild>
            <a href="#">
              <Avatar>
                {guild?.icon && (
                  <AvatarImage src={getGuildIconUrl(guild?.id, guild?.icon)} />
                )}
                <AvatarFallback>
                  {abbreviateName(guild?.name ?? "")}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{guild?.name}</span>
                <span className="truncate text-xs text-muted-foreground">
                  {guild?.id}
                </span>
              </div>
            </a>
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
          align="start"
          side={isMobile ? "bottom" : "right"}
          sideOffset={4}
        >
          <DropdownMenuLabel className="text-xs text-muted-foreground">
            Apps
          </DropdownMenuLabel>
          {filteredGuilds?.map((guild) => (
            <DropdownMenuItem
              key={guild!.id}
              onClick={() => setGuild(guild!.id)}
              className="gap-2 p-2"
            >
              <img
                src={getGuildIconUrl(guild!.id, guild!.icon ?? "")}
                className="size-6 shrink-0 rounded-sm border"
                alt=""
              />
              {guild!.name}
            </DropdownMenuItem>
          ))}
          <DropdownMenuSeparator />
          <DropdownMenuItem className="gap-2 p-2 cursor-pointer" asChild>
            <Link href="/dashboard">
              <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                <CogIcon className="size-4" />
              </div>
              <div className="font-medium">Manage servers</div>
            </Link>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  );
}

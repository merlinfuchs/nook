"use client";

import {
  LayoutDashboardIcon,
  SettingsIcon,
  WandSparklesIcon,
} from "lucide-react";

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useGuildId } from "@/lib/hooks/params";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useCallback, useMemo } from "react";

export default function DashboardNavMain() {
  const guildId = useGuildId();
  const pathname = usePathname();

  const isActive = useCallback(
    (path: string, exact = false) => {
      if (exact) {
        return pathname === path;
      }

      return pathname.startsWith(path);
    },
    [pathname]
  );

  const items = useMemo(() => {
    const dashboardPath = `/dashboard/${guildId}`;

    return [
      {
        title: "Dashboard",
        url: dashboardPath,
        active: isActive(dashboardPath, true),
        icon: LayoutDashboardIcon,
      },
      {
        title: "Personalize",
        url: dashboardPath + "/personalize",
        active: isActive(dashboardPath + "/personalize"),
        icon: WandSparklesIcon,
      },
      {
        title: "Settings",
        url: dashboardPath + "/settings",
        active: isActive(dashboardPath + "/settings"),
        icon: SettingsIcon,
      },
    ];
  }, [isActive, guildId]);

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Server</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => (
          <SidebarMenuItem key={item.title}>
            <SidebarMenuButton
              asChild
              tooltip={item.title}
              isActive={item.active}
            >
              <Link href={item.url}>
                <item.icon />
                <span>{item.title}</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}

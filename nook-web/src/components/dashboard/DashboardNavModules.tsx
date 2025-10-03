"use client";

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useModules } from "@/lib/hooks/api";
import { useGuildId } from "@/lib/hooks/params";
import { DynamicIcon } from "lucide-react/dynamic";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useCallback, useMemo } from "react";

export default function DashboardNavModules() {
  const guildId = useGuildId();
  const modules = useModules();

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

    return modules
      ?.filter((m) => !m.metadata.internal)
      .map((module) => ({
        name: module.metadata.name,
        url: `${dashboardPath}/modules/${module.id}`,
        icon: module.metadata.icon,
        active: isActive(`${dashboardPath}/modules/${module.id}`),
      }));
  }, [isActive, guildId, modules]);

  return (
    <SidebarGroup className="group-data-[collapsible=icon]:hidden">
      <SidebarGroupLabel>Modules</SidebarGroupLabel>
      <SidebarMenu>
        {items?.map((module) => (
          <SidebarMenuItem key={module.name}>
            <SidebarMenuButton asChild isActive={module.active}>
              <Link href={module.url}>
                <DynamicIcon name={module.icon as any} />
                <span>{module.name}</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}

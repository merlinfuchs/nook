"use client";

import * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
} from "@/components/ui/sidebar";
import DashboardNavGuildSelector from "./DashboardNavGuildSelector";
import DashboardNavMain from "./DashboardNavMain";
import DashboardNavModules from "./DashboardNavModules";
import DashboardNavSecondary from "./DashboardNavSecondary";
import DashboardNavUser from "./DashboardNavUser";

export default function DashboardSidebar({
  ...props
}: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar variant="inset" collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <DashboardNavGuildSelector />
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <DashboardNavMain />
        <DashboardNavModules />
        <DashboardNavSecondary className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <DashboardNavUser />
      </SidebarFooter>
    </Sidebar>
  );
}

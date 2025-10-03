"use client";

import { useModules } from "@/lib/hooks/api";
import { useMemo } from "react";
import ModuleListEntry from "./ModuleListEntry";
import QuickstartCard from "./QuickstartCard";

export default function ModuleList() {
  const modules = useModules();

  const filteredModules = useMemo(() => {
    return modules?.filter((m) => !m.metadata.internal);
  }, [modules]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-5">
      <QuickstartCard />
      {filteredModules?.map((mod) => (
        <ModuleListEntry {...mod} key={mod.id} />
      ))}
    </div>
  );
}

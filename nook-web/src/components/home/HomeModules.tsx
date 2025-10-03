"use client";

import { useGlobalModules } from "@/lib/hooks/api";
import { DynamicIcon } from "lucide-react/dynamic";
import { useMemo } from "react";

export default function HomeModules() {
  const modules = useGlobalModules();

  const filteredModules = useMemo(() => {
    if (!modules) return [];
    return modules.filter((m) => !m.metadata.internal);
  }, [modules]);

  return (
    <div id="modules" className="w-full py-12 xs:py-20 px-6">
      <h2 className="text-3xl xs:text-4xl md:text-5xl leading-[1.15]! font-bold tracking-tight text-center">
        A Module For Every Need
      </h2>
      <p className="mt-3 xs:text-lg text-center text-muted-foreground">
        Nook is designed to be fully modular, so you can pick and choose the
        features you need.
      </p>

      <div className="w-full max-w-(--breakpoint-lg) mx-auto mt-10 sm:mt-16 grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredModules.map((module) => (
          <div
            key={module.metadata.name}
            className="flex flex-col bg-background border rounded-xl py-6 px-5"
          >
            <div className="mb-3 h-10 w-10 flex items-center justify-center bg-muted rounded-full">
              <DynamicIcon
                name={module.metadata.icon as any}
                className="h-6 w-6"
              />
            </div>
            <span className="text-lg font-semibold">
              {module.metadata.name}
            </span>
            <p className="mt-1 text-foreground/80 text-[15px]">
              {module.metadata.description}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}

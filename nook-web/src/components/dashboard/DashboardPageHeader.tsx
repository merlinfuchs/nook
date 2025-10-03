"use client";

import { cn } from "@/lib/utils";
import { DynamicIcon } from "lucide-react/dynamic";
import { useEffect } from "react";
import { Separator } from "../ui/separator";

export default function DashboardPageHeader({
  icon,
  title,
  description,
  children,
  tabs,
  tab,
  onTabChange,
}: {
  icon?: string;
  title: string;
  description: string;
  children?: React.ReactNode;

  tabs?: {
    label: string;
    value: string;
    disabled?: boolean;
  }[];
  tab?: string;
  onTabChange?: (tab: string) => void;
}) {
  useEffect(() => {
    if (!tabs?.length) return;

    const activeTab = tabs.find((t) => t.value === tab);
    if (!activeTab || activeTab.disabled) {
      const newTab = tabs.find((t) => !t.disabled);
      if (newTab) {
        onTabChange?.(newTab.value);
      }
    }
  }, [tab, tabs]);

  return (
    <div className="flex-auto">
      <div className="flex gap-5 flex-col lg:flex-row justify-between lg:items-end">
        <div>
          <div className="flex items-center gap-2 mb-2">
            {icon && (
              <DynamicIcon
                name={icon as any}
                className="size-8 text-muted-foreground"
              />
            )}
            <h1 className="text-3xl font-bold">{title}</h1>
          </div>
          <p className="text-muted-foreground text-lg font-light">
            {description}
          </p>
        </div>
        <div className="flex items-center gap-5">{children}</div>
      </div>
      {tabs ? (
        <div className="mt-8">
          <div className="flex items-center gap-5">
            {tabs.map((t) => (
              <button
                key={t.value}
                onClick={() => onTabChange?.(t.value)}
                className={cn(
                  "text-lg font-medium pb-4 text-muted-foreground border-b-2 border-transparent px-3 hover:text-foreground disabled:opacity-50 disabled:hover:text-muted-foreground disabled:cursor-not-allowed",
                  t.value === tab && "text-foreground border-foreground"
                )}
                disabled={t.disabled}
              >
                {t.label}
              </button>
            ))}
          </div>
          <Separator className="mb-10" />
        </div>
      ) : (
        <Separator className="mt-6 mb-10" />
      )}
    </div>
  );
}

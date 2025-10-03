"use client";

import { useUser } from "@/lib/hooks/api";
import { useRouter } from "next/navigation";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  useUser(undefined, (res) => {
    if (res.success === false) {
      router.push(`/login`);
    }
  });

  return <div>{children}</div>;
}

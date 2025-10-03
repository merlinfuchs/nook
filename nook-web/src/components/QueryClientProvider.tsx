"use client";

import { QueryClientProvider as TanstackQueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/lib/api/client";

export default function QueryClientProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <TanstackQueryClientProvider client={queryClient}>
      {children}
    </TanstackQueryClientProvider>
  );
}

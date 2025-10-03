import { Button } from "@/components/ui/button";
import { ArrowUpRight, Forward } from "lucide-react";
import Link from "next/link";
import env from "@/lib/env/client";

export default function HomeCTABanner() {
  return (
    <div className="px-6">
      <div className="dark:border relative overflow-hidden my-20 w-full dark bg-background text-foreground max-w-(--breakpoint-lg) mx-auto rounded-2xl py-10 md:py-16 px-6 md:px-14">
        <div className="relative z-0 flex flex-col gap-3">
          <h3 className="text-3xl md:text-4xl font-semibold">
            Ready to Elevate Your Experience?
          </h3>
          <p className="mt-2 text-base md:text-lg">
            Take your Discord experience to the next level with Nook. Sign up
            today and start exploring!
          </p>
        </div>
        <div className="relative z-0 mt-14 flex flex-col sm:flex-row gap-4">
          <Button size="lg" asChild>
            <Link href="/dashboard">
              Get Started <ArrowUpRight className="h-5! w-5!" />
            </Link>
          </Button>
          <Button size="lg" variant="outline" asChild>
            <Link href={env.NEXT_PUBLIC_DISCORD_LINK} target="_blank">
              Discord Server <Forward className="h-5! w-5!" />
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}

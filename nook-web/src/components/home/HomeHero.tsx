import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import env from "@/lib/env/client";
import { ArrowUpRight, Forward } from "lucide-react";
import Link from "next/link";

export default function HomeHero() {
  return (
    <div className="min-h-[calc(100vh-6rem)] flex flex-col items-center py-20 px-6">
      <div className="md:mt-6 flex items-center justify-center">
        <div className="text-center max-w-2xl flex flex-col items-center">
          <Badge className="py-1 px-2">Nook is now live! 🚀</Badge>
          <h1 className="mt-6 max-w-[20ch] text-3xl xs:text-4xl sm:text-5xl md:text-6xl font-bold leading-[1.2]! tracking-tight">
            The friendliest way to start your Discord
          </h1>
          <p className="mt-6 max-w-[60ch] xs:text-lg">
            Whether you’re creating your very first server or helping friends
            get started, Nook makes Discord feel simple. It’s lightweight,
            beginner-friendly, and always ready to help.
          </p>
          <div className="mt-12 flex flex-col sm:flex-row items-center sm:justify-center gap-4">
            <Button
              size="lg"
              className="w-full sm:w-auto rounded-full text-base"
              asChild
            >
              <Link href="/dashboard">
                Get Started <ArrowUpRight className="h-5! w-5!" />
              </Link>
            </Button>
            <Button
              variant="outline"
              size="lg"
              className="w-full sm:w-auto rounded-full text-base shadow-none"
              asChild
            >
              <Link href={env.NEXT_PUBLIC_DISCORD_LINK} target="_blank">
                Discord Server <Forward className="h-5! w-5!" />
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

"use client";

import { Button } from "@/components/ui/button";
import HomeNavMenu from "./HomeNavMenu";
import HomeNavSheet from "./HomeNavSheet";
import logo from "@/assets/logo.svg";
import Image from "next/image";
import { useUser } from "@/lib/hooks/api";
import Link from "next/link";
import env from "@/lib/env/client";

export default function HomeNavBar() {
  const user = useUser();

  const inviteUrl = env.NEXT_PUBLIC_API_PUBLIC_BASE_URL + "/v1/auth/invite";

  return (
    <nav className="fixed z-10 top-6 inset-x-4 h-14 xs:h-16 bg-background/50 backdrop-blur-xs border dark:border-slate-700/70 max-w-(--breakpoint-xl) mx-auto rounded-full">
      <div className="h-full flex items-center justify-between mx-auto px-4">
        <div className="flex items-center gap-2">
          <Image
            src={logo}
            alt="Nook Logo"
            className="w-10 h-10"
            height={128}
            width={128}
          />
          <div className="text-xl font-bold text-primary">Nook</div>
        </div>

        {/* Desktop Menu */}
        <HomeNavMenu className="hidden md:block" />

        <div className="flex items-center gap-3">
          <Button variant="outline" className="hidden sm:inline-flex" asChild>
            <a href={inviteUrl}>Add Bot</a>
          </Button>
          {user ? (
            <Button asChild>
              <Link href="/dashboard">Dashboard</Link>
            </Button>
          ) : (
            <Button asChild>
              <Link href="/login">Sign In</Link>
            </Button>
          )}

          {/* Mobile Menu */}
          <div className="md:hidden">
            <HomeNavSheet />
          </div>
        </div>
      </div>
    </nav>
  );
}

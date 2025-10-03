import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { useUser } from "@/lib/hooks/api";
import { Menu } from "lucide-react";
import Link from "next/link";
import HomeNavMenu from "./HomeNavMenu";
import env from "@/lib/env/client";
import logo from "@/assets/logo.svg";
import Image from "next/image";

export default function HomeNavSheet() {
  const user = useUser();

  const inviteUrl = env.NEXT_PUBLIC_API_PUBLIC_BASE_URL + "/v1/auth/invite";

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="outline" size="icon" className="rounded-full">
          <Menu />
        </Button>
      </SheetTrigger>
      <SheetContent>
        <SheetHeader>
          <div className="flex space-x-2 items-center">
            <Image
              src={logo}
              alt="Nook Logo"
              className="size-6"
              height={128}
              width={128}
            />
            <SheetTitle>Nook</SheetTitle>
          </div>
        </SheetHeader>

        <div className="px-3">
          <HomeNavMenu orientation="vertical" />
        </div>

        <SheetFooter>
          <div className="space-y-4">
            <Button variant="outline" className="w-full sm:hidden">
              <a href={inviteUrl}>Add Bot</a>
            </Button>
            {user ? (
              <Button className="w-full" asChild>
                <Link href="/dashboard">Dashboard</Link>
              </Button>
            ) : (
              <Button className="w-full" asChild>
                <Link href="/login">Sign In</Link>
              </Button>
            )}
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}

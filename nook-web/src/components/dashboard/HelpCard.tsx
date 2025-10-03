import env from "@/lib/env/client";
import { ArrowRightIcon, LifeBuoyIcon } from "lucide-react";
import Link from "next/link";
import { Button } from "../ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";

export default function HelpCard() {
  return (
    <Card className="border-primary border-2 p-0 flex flex-col">
      <CardHeader className="mb-3 p-5 pb-0 flex-auto">
        <div className="flex items-center gap-5">
          <div className="size-12 border rounded-lg flex items-center justify-center flex-none">
            <LifeBuoyIcon className="size-7" />
          </div>
          <div>
            <CardDescription className="mb-0.5 uppercase text-xs">
              Help Center
            </CardDescription>
            <CardTitle className="text-xl">New around here?</CardTitle>
          </div>
        </div>
      </CardHeader>
      <CardFooter className="gap-3 flex-col md:flex-row p-5 pt-1">
        <Button asChild>
          <Link
            href={env.NEXT_PUBLIC_DISCORD_LINK}
            target="_blank"
            className="w-full md:w-auto"
          >
            <ArrowRightIcon className="size-4 mr-2" />
            Join the Discord server
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

"use client";

import { GuildWire } from "@/lib/types/guild.gen";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { getGuildIconUrl } from "@/lib/discord/cdn";
import { abbreviateName } from "@/lib/utils";
import Link from "next/link";
import { Button } from "../ui/button";
import env from "@/lib/env/client";

export default function GuildListEntry({ guild }: { guild: GuildWire }) {
  const inviteUrl =
    env.NEXT_PUBLIC_API_PUBLIC_BASE_URL +
    "/v1/auth/invite?guild_id=" +
    guild.id;

  return (
    <Card>
      <CardHeader className="flex flex-row gap-4 flex-auto">
        <Avatar className="size-10">
          {guild.icon && (
            <AvatarImage src={getGuildIconUrl(guild.id, guild.icon)} alt="" />
          )}
          <AvatarFallback className="text-xl">
            {abbreviateName(guild.name)}
          </AvatarFallback>
        </Avatar>
        <div className="flex-auto truncate">
          <CardTitle className="text-base truncate">{guild.name}</CardTitle>
          <CardDescription className="text-sm">{guild.id}</CardDescription>
        </div>
      </CardHeader>
      <CardFooter>
        {guild.bot ? (
          <Button asChild>
            <Link href={`/dashboard/${guild.id}`}>Open Dashboard</Link>
          </Button>
        ) : (
          <Button variant="outline" asChild>
            <Link href={inviteUrl}>Add Bot</Link>
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}

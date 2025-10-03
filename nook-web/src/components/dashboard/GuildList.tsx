"use client";

import { useGuilds } from "@/lib/hooks/api";
import GuildListEntry from "./GuildListEntry";
import { useMemo } from "react";

export default function GuildList() {
  const guilds = useGuilds();

  const filteredGuilds = useMemo(() => {
    return guilds
      ?.filter((guild) => guild.access)
      .sort((a, b) => Number(b.bot) - Number(a.bot));
  }, [guilds]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      {filteredGuilds?.map((guild) => (
        <GuildListEntry guild={guild} key={guild.id} />
      ))}
    </div>
  );
}

"use client";

import { useGuildId } from "@/lib/hooks/params";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Input } from "../ui/input";
import { useGuildSettings } from "@/lib/hooks/api";
import { useGuildSettingsUpdateMutation } from "@/lib/api/mutations";

export default function GuildSettingsPrefix() {
  const guildId = useGuildId();

  const guildSettings = useGuildSettings();
  const updateMutation = useGuildSettingsUpdateMutation(guildId);

  const [commandPrefix, setCommandPrefix] = useState("");

  const onSubmit = useCallback(() => {
    updateMutation.mutate(
      { command_prefix: commandPrefix },
      {
        onSuccess: (res) => {
          if (res.success) {
            toast.success("Settings have been saved!");
          } else {
            toast.error(`Failed to save settings: ${res.error.message}`);
          }
        },
      }
    );
  }, [updateMutation, commandPrefix]);

  useEffect(() => {
    if (guildSettings) {
      setCommandPrefix(guildSettings.command_prefix ?? "");
    }
  }, [guildSettings]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Command Prefix</CardTitle>
        <CardDescription>
          Configure the command prefix for your server here.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <Input
          value={commandPrefix}
          placeholder={guildSettings?.default.command_prefix}
          onChange={(e) => setCommandPrefix(e.target.value)}
        />
      </CardContent>
      <CardFooter>
        <Button onClick={onSubmit}>Update Prefix</Button>
      </CardFooter>
    </Card>
  );
}

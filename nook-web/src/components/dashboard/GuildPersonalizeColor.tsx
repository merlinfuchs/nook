"use client";

import { useGuildSettingsUpdateMutation } from "@/lib/api/mutations";
import { useGuildSettings } from "@/lib/hooks/api";
import { useGuildId } from "@/lib/hooks/params";
import { useCallback, useMemo } from "react";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import GuildColorSchemeSelector from "./GuildColorSchemeSelector";

export default function GuildPersonalizeColor() {
  const guildId = useGuildId();

  const guildSettings = useGuildSettings();
  const updateMutation = useGuildSettingsUpdateMutation(guildId);

  const colorScheme = useMemo(() => {
    return guildSettings?.color_scheme ?? guildSettings?.default.color_scheme;
  }, [guildSettings]);

  const setColorScheme = useCallback(
    (colorScheme: string) => {
      updateMutation.mutate(
        { color_scheme: colorScheme },
        {
          onSuccess: (res) => {
            if (res.success) {
              toast.success("Color scheme updated!");
            } else {
              toast.error(
                `Failed to update color scheme: ${res.error.message}`
              );
            }
          },
        }
      );
    },
    [updateMutation]
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Color Scheme</CardTitle>
        <CardDescription>
          Chose a color scheme. This is used for response messages across all
          modules.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <GuildColorSchemeSelector
          value={colorScheme ?? null}
          onChange={setColorScheme}
        />
      </CardContent>
    </Card>
  );
}

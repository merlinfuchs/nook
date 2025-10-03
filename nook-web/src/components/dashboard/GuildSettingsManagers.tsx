"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { XIcon } from "lucide-react";
import {
  useGuildManagerAddMutation,
  useGuildManagerRemoveMutation,
} from "@/lib/api/mutations";
import { useGuildId } from "@/lib/hooks/params";
import { useGuildManagers } from "@/lib/hooks/api";
import { useCallback, useState } from "react";
import { toast } from "sonner";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { getUserAvatarUrl } from "@/lib/discord/cdn";
import { abbreviateName } from "@/lib/utils";

export default function GuildSettingsManagers() {
  const guildId = useGuildId();

  const managers = useGuildManagers();

  const addMutation = useGuildManagerAddMutation(guildId);
  const removeMutation = useGuildManagerRemoveMutation(guildId);

  const [newUserID, setNewUserID] = useState("");

  const addManager = useCallback(() => {
    addMutation.mutate(
      {
        user_id: newUserID,
        role: "admin",
      },
      {
        onSuccess: (res) => {
          if (res.success) {
            setNewUserID("");
            toast.success("Manager added!");
          } else {
            toast.error(`Failed to add manager: ${res.error.message}`);
          }
        },
      }
    );
  }, [addMutation, newUserID, setNewUserID]);

  const removeManager = useCallback(
    (userID: string) => {
      removeMutation.mutate(userID, {
        onSuccess: (res) => {
          if (res.success) {
            toast.success("Manager removed!");
          } else {
            toast.error(`Failed to remove manager: ${res.error.message}`);
          }
        },
      });
    },
    [removeMutation]
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Bot Admins</CardTitle>
        <CardDescription>
          Users with{" "}
          <span className="font-bold">administrator permissions</span> will
          always be able to manage the bot and access the server&apos;s
          dashboard. You can add additional users below that will have access
          regardless of their permissions.
        </CardDescription>
      </CardHeader>

      {!!managers?.length && (
        <CardContent>
          <div className="flex flex-wrap gap-3">
            {managers?.map((manager) => (
              <div
                className="p-1 rounded-full border flex items-center gap-2"
                key={manager.user.id}
              >
                <Avatar className="size-7">
                  <AvatarImage
                    src={getUserAvatarUrl(manager.user.id, manager.user.avatar)}
                  />
                  <AvatarFallback>
                    {abbreviateName(manager.user.display_name)}
                  </AvatarFallback>
                </Avatar>
                <div className="text-sm font-bold">
                  {manager.user.display_name}
                </div>
                <button
                  className="p-1.5 rounded-full hover:bg-muted"
                  onClick={() => removeManager(manager.user.id)}
                >
                  <XIcon className="size-4" />
                </button>
              </div>
            ))}
          </div>
        </CardContent>
      )}
      <CardFooter>
        <div className="flex gap-3">
          <Input
            placeholder="User ID"
            value={newUserID}
            onChange={(e) => setNewUserID(e.target.value)}
          />
          <Button variant="outline" onClick={addManager}>
            Add User
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}

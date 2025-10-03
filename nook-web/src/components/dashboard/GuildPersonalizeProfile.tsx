"use client";

import { useGuildProfileUpdateMutation } from "@/lib/api/mutations";
import { useGuildProfile } from "@/lib/hooks/api";
import { useGuildId } from "@/lib/hooks/params";
import { abbreviateName } from "@/lib/utils";
import { TrashIcon, UploadIcon } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { toast } from "sonner";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
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
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";

export default function GuildPersonalizeProfile() {
  const guildId = useGuildId();

  const profile = useGuildProfile();
  const updateMutation = useGuildProfileUpdateMutation(guildId);

  const [name, setName] = useState("");
  const [bio, setBio] = useState("");
  const [newAvatar, setNewAvatar] = useState<string | null>(null);

  const avatarInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setName(profile?.custom_name ?? "");
    setBio(profile?.custom_bio ?? "");
  }, [profile]);

  const onAvatarChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (file) {
        // Check file size (1MB = 1024 * 1024 bytes)
        if (file.size > 1024 * 1024) {
          toast("File size must be less than 1MB");
          return;
        }

        // Encode file as data URL
        const reader = new FileReader();
        reader.onload = (event) => {
          const dataUrl = event.target?.result as string;
          setNewAvatar(dataUrl);
        };
        reader.onerror = () => {
          toast("Error reading file");
        };
        reader.readAsDataURL(file);
      }
    },
    []
  );

  const onSubmit = useCallback(() => {
    updateMutation.mutate(
      {
        name,
        bio,
        avatar: newAvatar ?? undefined,
      },
      {
        onSuccess: (res) => {
          if (res.success) {
            toast.success("Profile updated!");
            setNewAvatar(null);
          } else {
            toast.error(`Failed to update profile: ${res.error.message}`);
          }
        },
      }
    );
  }, [updateMutation, name, bio, newAvatar]);

  const avatarSrc = useMemo(() => {
    if (newAvatar === "") {
      return profile?.default_avatar_url;
    }

    if (newAvatar) {
      return newAvatar;
    }

    if (profile?.custom_avatar_url) {
      return profile.custom_avatar_url;
    }

    return profile?.default_avatar_url;
  }, [newAvatar, profile]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Bot Profile</CardTitle>
        <CardDescription>
          Chose a color scheme. This is used for response messages across all
          modules.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <div className="flex gap-8">
          <div className="flex flex-col items-center gap-3">
            <Avatar className="size-32 cursor-pointer relative overflow-visible">
              {<AvatarImage src={avatarSrc} className="rounded-full" alt="" />}
              <AvatarFallback>
                {abbreviateName(profile?.custom_name ?? "")}
              </AvatarFallback>

              {profile?.custom_avatar_url && newAvatar === null && (
                <Button
                  variant="outline"
                  size="icon"
                  className="absolute top-0 left-0 bg-muted!"
                  onClick={() => setNewAvatar("")}
                >
                  <TrashIcon className="size-4" />
                </Button>
              )}
            </Avatar>
            <div>
              <input
                type="file"
                accept="image/png, image/jpeg"
                className="hidden"
                onChange={onAvatarChange}
                ref={avatarInputRef}
              />
              <Button
                variant="outline"
                className="flex items-center gap-2"
                onClick={() => avatarInputRef.current?.click()}
              >
                <UploadIcon className="size-4" />
                <div>Upload Image</div>
              </Button>
            </div>
          </div>
          <div className="flex flex-col gap-3 flex-auto">
            <div>
              <Label className="mb-2 font-bold text-base">Name</Label>
              <Input
                placeholder={profile?.default_name}
                className="max-w-48"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div>
              <Label className="mb-2 font-bold text-base">Bio</Label>
              <Textarea
                className="resize-none"
                rows={5}
                placeholder={profile?.default_bio}
                value={bio}
                onChange={(e) => setBio(e.target.value)}
              />
            </div>
          </div>
        </div>
      </CardContent>

      <CardFooter>
        <Button onClick={onSubmit}>Update Profile</Button>
      </CardFooter>
    </Card>
  );
}

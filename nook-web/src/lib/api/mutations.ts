import { useMutation, useQueryClient } from "@tanstack/react-query";
import { SubscriptionManageResponseWire } from "../types/billing.gen";
import {
  GuildManagerAddAddResponseWire,
  GuildManagerCreateAddWire,
  GuildManagerRemoveResponseWire,
  GuildProfileUpdateRequestWire,
  GuildProfileUpdateResponseWire,
  GuildSettingsGetResponseWire,
  GuildSettingsUpdateRequestWire,
  GuildSettingsUpdateResponseWire,
} from "../types/guild.gen";
import {
  ModuleConfigureRequestWire,
  ModuleConfigureResponseWire,
  ModuleGetResponseWire,
} from "../types/manage.gen";
import client, { apiRequest } from "./client";
import { APIResponse } from "./response";

export function useBillingSubscriptionManageMutation(
  guildId: string,
  subscriptionID: string
) {
  return useMutation({
    mutationFn: () =>
      apiRequest<SubscriptionManageResponseWire>(
        `/v1/guilds/${guildId}/billing/subscriptions/${subscriptionID}/manage`,
        {
          method: "POST",
        }
      ),
  });
}

export function useAuthLogoutMutation() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: () => apiRequest<void>("/v1/auth/logout", { method: "POST" }),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: [] });
    },
  });
}

export function useGuildSettingsUpdateMutation(guildId: string) {
  return useMutation({
    mutationFn: (data: GuildSettingsUpdateRequestWire) =>
      apiRequest<GuildSettingsUpdateResponseWire>(
        `/v1/guilds/${guildId}/settings`,
        {
          method: "PUT",
          body: JSON.stringify(data),
          headers: { "Content-Type": "application/json" },
        }
      ),
    onSuccess: (res) => {
      if (res.success) {
        client.setQueryData(
          ["guilds", guildId, "settings"],
          res satisfies APIResponse<GuildSettingsGetResponseWire>
        );
      }
    },
  });
}

export function useModuleConfigureMutation(guildId: string, moduleId: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (data: ModuleConfigureRequestWire) =>
      apiRequest<ModuleConfigureResponseWire>(
        `/v1/guilds/${guildId}/modules/${moduleId}`,
        {
          method: "PUT",
          body: JSON.stringify(data),
          headers: { "Content-Type": "application/json" },
        }
      ),
    onSuccess: (res) => {
      if (res.success) {
        client.setQueryData(
          ["guilds", guildId, "modules", moduleId],
          res satisfies APIResponse<ModuleGetResponseWire>
        );
      }
    },
  });
}

export function useGuildManagerAddMutation(guildId: string) {
  return useMutation({
    mutationFn: (data: GuildManagerCreateAddWire) =>
      apiRequest<GuildManagerAddAddResponseWire>(
        `/v1/guilds/${guildId}/managers`,
        {
          method: "POST",
          body: JSON.stringify(data),
          headers: { "Content-Type": "application/json" },
        }
      ),
    onSuccess: (res) => {
      if (res.success) {
        client.invalidateQueries({ queryKey: ["guilds", guildId, "managers"] });
      }
    },
  });
}

export function useGuildManagerRemoveMutation(guildId: string) {
  return useMutation({
    mutationFn: (userID: string) =>
      apiRequest<GuildManagerRemoveResponseWire>(
        `/v1/guilds/${guildId}/managers/${userID}`,
        {
          method: "DELETE",
        }
      ),
    onSuccess: (res) => {
      if (res.success) {
        client.invalidateQueries({ queryKey: ["guilds", guildId, "managers"] });
      }
    },
  });
}

export function useGuildProfileUpdateMutation(guildId: string) {
  return useMutation({
    mutationFn: (data: GuildProfileUpdateRequestWire) =>
      apiRequest<GuildProfileUpdateResponseWire>(
        `/v1/guilds/${guildId}/profile`,
        {
          method: "PUT",
          body: JSON.stringify(data),
          headers: { "Content-Type": "application/json" },
        }
      ),
    onSuccess: (res) => {
      if (res.success) {
        client.invalidateQueries({ queryKey: ["guilds", guildId, "profile"] });
      }
    },
  });
}

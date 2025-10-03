import { useQuery } from "@tanstack/react-query";
import { apiRequest } from "./client";
import { UserGetResponseWire } from "../types/user.gen";
import {
  ChannelGetResponseWire,
  ChannelListResponseWire,
  GuildGetResponseWire,
  GuildListResponseWire,
  GuildManagerListResponseWire,
  GuildProfileGetResponseWire,
  GuildSettingsGetResponseWire,
  RoleGetResponseWire,
  RoleListResponseWire,
} from "../types/guild.gen";
import {
  ModuleListResponseWire,
  ModuleGetResponseWire,
} from "../types/manage.gen";
import {
  BillingFeaturesGetResponseWire,
  BillingPlanListResponseWire,
  SubscriptionListResponseWire,
} from "../types/billing.gen";

export function useUserQuery(userId = "@me") {
  return useQuery({
    queryKey: ["users", userId],
    queryFn: () => apiRequest<UserGetResponseWire>(`/v1/users/${userId}`),
    enabled: !!userId,
  });
}

export function useGuildsQuery() {
  return useQuery({
    queryKey: ["guilds"],
    queryFn: () => apiRequest<GuildListResponseWire>(`/v1/guilds`),
  });
}

export function useGuildQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId],
    queryFn: () => apiRequest<GuildGetResponseWire>(`/v1/guilds/${guildId}`),
    enabled: !!guildId,
  });
}

export function useGuildRolesQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "roles"],
    queryFn: () =>
      apiRequest<RoleListResponseWire>(`/v1/guilds/${guildId}/roles`),
    enabled: !!guildId,
  });
}

export function useGuildRoleQuery(guildId: string, roleId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "roles", roleId],
    queryFn: () =>
      apiRequest<RoleGetResponseWire>(`/v1/guilds/${guildId}/roles/${roleId}`),
    enabled: !!guildId && !!roleId,
  });
}

export function useGuildChannelsQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "channels"],
    queryFn: () =>
      apiRequest<ChannelListResponseWire>(`/v1/guilds/${guildId}/channels`),
    enabled: !!guildId,
  });
}

export function useGuildChannelQuery(guildId: string, channelId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "channels", channelId],
    queryFn: () =>
      apiRequest<ChannelGetResponseWire>(
        `/v1/guilds/${guildId}/channels/${channelId}`
      ),
    enabled: !!guildId && !!channelId,
  });
}

export function useGuildSettingsQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "settings"],
    queryFn: () =>
      apiRequest<GuildSettingsGetResponseWire>(
        `/v1/guilds/${guildId}/settings`
      ),
    enabled: !!guildId,
  });
}

export function useBillingPlansQuery() {
  return useQuery({
    queryKey: ["billing", "plans"],
    queryFn: () => apiRequest<BillingPlanListResponseWire>(`/v1/billing/plans`),
  });
}

export function useBillingFeaturesQuery(guildId: string) {
  return useQuery({
    queryKey: ["billing", "guildss", guildId, "features"],
    queryFn: () =>
      apiRequest<BillingFeaturesGetResponseWire>(
        `/v1/guilds/${guildId}/billing/features`
      ),
    enabled: !!guildId,
  });
}

export function useBillingSubscriptionsQuery(guildId: string) {
  return useQuery({
    queryKey: ["billing", "guildss", guildId, "subscriptions"],
    queryFn: () =>
      apiRequest<SubscriptionListResponseWire>(
        `/v1/guilds/${guildId}/billing/subscriptions`
      ),
    enabled: !!guildId,
  });
}

export function useGlobalModulesQuery() {
  return useQuery({
    queryKey: ["modules"],
    queryFn: () => apiRequest<ModuleListResponseWire>(`/v1/modules`),
  });
}

export function useModulesQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "modules"],
    queryFn: () =>
      apiRequest<ModuleListResponseWire>(`/v1/guilds/${guildId}/modules`),
    enabled: !!guildId,
  });
}

export function useModuleQuery(guildId: string, moduleId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "modules", moduleId],
    queryFn: () =>
      apiRequest<ModuleGetResponseWire>(
        `/v1/guilds/${guildId}/modules/${moduleId}`
      ),
    enabled: !!guildId && !!moduleId,
  });
}

export function useGuildManagersQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "managers"],
    queryFn: () =>
      apiRequest<GuildManagerListResponseWire>(
        `/v1/guilds/${guildId}/managers`
      ),
    enabled: !!guildId,
  });
}

export function useGuildProfileQuery(guildId: string) {
  return useQuery({
    queryKey: ["guilds", guildId, "profile"],
    queryFn: () =>
      apiRequest<GuildProfileGetResponseWire>(`/v1/guilds/${guildId}/profile`),
    enabled: !!guildId,
  });
}

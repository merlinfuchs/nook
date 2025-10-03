import { useEffect } from "react";
import {
  useBillingFeaturesQuery,
  useBillingPlansQuery,
  useBillingSubscriptionsQuery,
  useGlobalModulesQuery,
  useGuildChannelQuery,
  useGuildChannelsQuery,
  useGuildManagersQuery,
  useGuildProfileQuery,
  useGuildQuery,
  useGuildRoleQuery,
  useGuildRolesQuery,
  useGuildSettingsQuery,
  useGuildsQuery,
  useModuleQuery,
  useModulesQuery,
  useUserQuery,
} from "../api/queries";
import { APIResponse } from "../api/response";
import {
  BillingFeaturesGetResponseWire,
  BillingPlanListResponseWire,
  SubscriptionListResponseWire,
} from "../types/billing.gen";
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
import { UserGetResponseWire } from "../types/user.gen";
import { useGuildId, useModuleId } from "./params";
import {
  ModuleGetResponseWire,
  ModuleListResponseWire,
} from "../types/manage.gen";

export function useResponseData<T>(
  {
    data,
  }: {
    data?: APIResponse<T>;
  },
  callback?: (res: APIResponse<T>) => void
): T | undefined {
  useEffect(() => {
    if (data !== undefined && callback) {
      callback(data);
    }
  }, [data, callback]);

  return data?.success ? data.data : undefined;
}

export function useUser(
  userId?: string,
  callback?: (res: APIResponse<UserGetResponseWire>) => void
) {
  const query = useUserQuery(userId);
  return useResponseData(query, callback);
}

export function useGuilds(
  callback?: (res: APIResponse<GuildListResponseWire>) => void
) {
  const query = useGuildsQuery();
  return useResponseData(query, callback);
}

export function useGuild(
  guildId?: string,
  callback?: (res: APIResponse<GuildGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();

  const query = useGuildQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useGuildRoles(
  guildId?: string,
  callback?: (res: APIResponse<RoleListResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildRolesQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useGuildRole(
  roleId: string,
  guildId?: string,
  callback?: (res: APIResponse<RoleGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildRoleQuery(guildId ?? currentGuildId, roleId);
  return useResponseData(query, callback);
}

export function useGuildChannels(
  guildId?: string,
  callback?: (res: APIResponse<ChannelListResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildChannelsQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useGuildChannel(
  channelId: string,
  guildId?: string,
  callback?: (res: APIResponse<ChannelGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildChannelQuery(guildId ?? currentGuildId, channelId);
  return useResponseData(query, callback);
}

export function useGuildSettings(
  guildId?: string,
  callback?: (res: APIResponse<GuildSettingsGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildSettingsQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useBillingPlans(
  callback?: (res: APIResponse<BillingPlanListResponseWire>) => void
) {
  const query = useBillingPlansQuery();
  return useResponseData(query, callback);
}

export function useBillingFeatures(
  guildId?: string,
  callback?: (res: APIResponse<BillingFeaturesGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useBillingFeaturesQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useBillingSubscriptions(
  guildId?: string,
  callback?: (res: APIResponse<SubscriptionListResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useBillingSubscriptionsQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useGlobalModules(
  callback?: (res: APIResponse<ModuleListResponseWire>) => void
) {
  const query = useGlobalModulesQuery();
  return useResponseData(query, callback);
}

export function useModules(
  guildId?: string,
  callback?: (res: APIResponse<ModuleListResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useModulesQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useModule(
  guildId?: string,
  moduleId?: string,
  callback?: (res: APIResponse<ModuleGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const currentModuleId = useModuleId();
  const query = useModuleQuery(
    guildId ?? currentGuildId,
    moduleId ?? currentModuleId
  );
  return useResponseData(query, callback);
}

export function useGuildManagers(
  guildId?: string,
  callback?: (res: APIResponse<GuildManagerListResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildManagersQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function useGuildProfile(
  guildId?: string,
  callback?: (res: APIResponse<GuildProfileGetResponseWire>) => void
) {
  const currentGuildId = useGuildId();
  const query = useGuildProfileQuery(guildId ?? currentGuildId);
  return useResponseData(query, callback);
}

export function getGuildIconUrl(guildId: string, iconHash: string) {
  return `https://cdn.discordapp.com/icons/${guildId}/${iconHash}.png`;
}

export function getUserAvatarUrl(userId: string, avatarHash?: string | null) {
  if (!avatarHash) {
    const index = (BigInt(userId ?? "0") >> BigInt(22)) % BigInt(6);
    return `https://cdn.discordapp.com/embed/avatars/${index}.png`;
  }

  return `https://cdn.discordapp.com/avatars/${userId}/${avatarHash}.png`;
}

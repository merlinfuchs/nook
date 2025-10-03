package guild

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

func (m *GuildModule) tryUpsertUser(ctx context.Context, userID common.ID) (*model.User, error) {
	discordUser, err := m.rest.GetUser(userID, rest.WithCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if discordUser.Bot {
		return nil, fmt.Errorf("user is a bot")
	}

	displayName := discordUser.Username
	if discordUser.GlobalName != nil {
		displayName = *discordUser.GlobalName
	}

	user := &model.User{
		ID:            userID,
		Username:      discordUser.Username,
		Discriminator: discordUser.Discriminator,
		DisplayName:   displayName,
		Avatar:        common.PtrToNullString(discordUser.Avatar),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	return m.userStore.UpsertUser(ctx, user)
}

func (m *GuildModule) guildSettingsToWire(settings model.GuildSettings) *GuildSettingsWire {
	defaultSettings := m.guildSettingsManager.DefaultSettings()

	return &GuildSettingsWire{
		CommandPrefix: settings.CommandPrefix,
		ColorScheme:   settings.ColorScheme,
		Default: DefaultGuildSettingsWire{
			CommandPrefix: defaultSettings.CommandPrefix,
			ColorScheme:   defaultSettings.ColorScheme,

			ResolvedColor: fmt.Sprintf("#%06x", defaultSettings.Color()),
		},
	}
}

func (m *GuildModule) guildCurrentMember(ctx context.Context, guildID common.ID) (*discord.Member, error) {
	var res *discord.Member
	err := m.rest.Do(
		// TODO: This is a hack to get the member, we should use the proper endpoint
		rest.UpdateMember.Compile(nil, guildID, "@me"),
		nil,
		&res,
		rest.WithCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	res.GuildID = guildID
	return res, nil
}

func (m *GuildModule) updateGuildCurrentMember(
	ctx context.Context,
	guildID common.ID,
	name string,
	bio string,
	avatar *string,
) (*discord.Member, error) {
	var res *discord.Member
	req := struct {
		Name   string  `json:"nick"`
		Bio    string  `json:"bio"`
		Avatar *string `json:"avatar,omitempty"`
	}{
		Name:   name,
		Bio:    bio,
		Avatar: avatar,
	}

	err := m.rest.Do(
		rest.UpdateMember.Compile(nil, guildID, "@me"),
		req,
		&res,
		rest.WithCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	res.GuildID = guildID
	return res, nil
}

func memberToProfileWire(member *discord.Member) *GuildProfileWire {
	res := GuildProfileWire{
		UserID:      member.User.ID,
		DefaultName: member.User.Username,
	}
	if member.User.Avatar != nil {
		res.DefaultAvatarURL = *member.User.AvatarURL()
	}
	if member.User.GlobalName != nil {
		res.DefaultName = *member.User.GlobalName
	}

	if member.Avatar != nil {
		res.CustomAvatarURL = *member.AvatarURL()
	}
	if member.Nick != nil {
		res.CustomName = *member.Nick
	}

	return &res
}

package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/api/session"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/store"
)

const StateCookieName = "nook-oauth-state"
const RedirectCookieName = "nook-oauth-redirect"

type AuthModuleConfig struct {
	ClientID      snowflake.ID
	ClientSecret  string
	APIURL        string
	AppURL        string
	SecureCookies bool
	StrictCookies bool
}

type AuthModule struct {
	client         oauth2.Client
	userStore      store.UserStore
	sessionStore   store.SessionStore
	sessionManager *session.SessionManager

	callbackURI   string
	redirectURI   string
	secureCookies bool
}

func NewAuthModule(
	cfg AuthModuleConfig,
	userStore store.UserStore,
	sessionStore store.SessionStore,
	sessionManager *session.SessionManager,
) *AuthModule {
	client := oauth2.New(cfg.ClientID, cfg.ClientSecret, oauth2.WithRestClientConfigOpts())

	callbackURI := fmt.Sprintf("%s/auth/callback", cfg.APIURL)
	redirectURI := cfg.AppURL

	return &AuthModule{
		client:         client,
		userStore:      userStore,
		sessionStore:   sessionStore,
		sessionManager: sessionManager,

		callbackURI:   callbackURI,
		redirectURI:   redirectURI,
		secureCookies: cfg.SecureCookies,
	}
}

func (m *AuthModule) ModuleID() string {
	return "auth"
}

func (m *AuthModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Authentication",
		Description:    "Handles authentication for the API.",
		Icon:           "key-round",
		Internal:       true,
		DefaultEnabled: true,
	}
}

func (m *AuthModule) Endpoints(mx api.HandlerGroup) {
	mx.Get("/auth/login", func(c *api.Context) error {
		url, state := m.client.GenerateAuthorizationURLState(oauth2.AuthorizationURLParams{
			RedirectURI: m.callbackURI,
			Scopes: []discord.OAuth2Scope{
				discord.OAuth2ScopeIdentify,
				discord.OAuth2ScopeGuilds,
			},
		})

		if redirectPath := c.Query("redirect"); redirectPath != "" {
			c.SetCookie(&http.Cookie{
				Name:     RedirectCookieName,
				Value:    redirectPath,
				HttpOnly: true,
				Secure:   m.secureCookies,
			})
		} else {
			c.DeleteCookie(RedirectCookieName)
		}

		c.SetCookie(&http.Cookie{
			Name:     StateCookieName,
			Value:    state,
			HttpOnly: true,
			Secure:   m.secureCookies,
		})

		c.Redirect(url, http.StatusTemporaryRedirect)
		return nil
	})

	mx.Get("/auth/invite", func(c *api.Context) error {
		rawGuildID := c.Query("guild_id")
		var guildID snowflake.ID
		if rawGuildID != "" {
			gid, err := common.ParseID(rawGuildID)
			if err != nil {
				return fmt.Errorf("invalid guild_id: %w", err)
			}
			guildID = gid
		}

		url, state := m.client.GenerateAuthorizationURLState(oauth2.AuthorizationURLParams{
			RedirectURI: m.callbackURI,
			Scopes: []discord.OAuth2Scope{
				discord.OAuth2ScopeIdentify,
				discord.OAuth2ScopeGuilds,
				discord.OAuth2ScopeBot,
				discord.OAuth2ScopeApplicationsCommands,
			},
			GuildID: guildID,
		})

		if redirectPath := c.Query("redirect"); redirectPath != "" {
			c.SetCookie(&http.Cookie{
				Name:     RedirectCookieName,
				Value:    redirectPath,
				HttpOnly: true,
				Secure:   m.secureCookies,
			})
		} else {
			c.DeleteCookie(RedirectCookieName)
		}

		c.SetCookie(&http.Cookie{
			Name:     StateCookieName,
			Value:    state,
			HttpOnly: true,
			Secure:   m.secureCookies,
		})

		c.Redirect(url, http.StatusTemporaryRedirect)
		return nil
	})

	mx.Get("/auth/callback", func(c *api.Context) error {
		queryCode := c.Query("code")
		queryState := c.Query("state")
		queryError := c.Query("error")
		guildID := c.Query("guild_id")

		if queryError != "" {
			c.Redirect(m.redirectURI, http.StatusTemporaryRedirect)
			return nil
		}

		state := c.Cookie(StateCookieName)
		if state == "" {
			return fmt.Errorf("state cookie not found")
		}

		if state != queryState {
			return fmt.Errorf("state mismatch")
		}

		discordSession, _, err := m.client.StartSession(queryCode, queryState, rest.WithCtx(c.Context()))
		if err != nil {
			return fmt.Errorf("failed to exchange code for session: %w", err)
		}

		discordUser, err := m.client.GetUser(discordSession, rest.WithCtx(c.Context()))
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		displayName := discordUser.Username
		if discordUser.GlobalName != nil {
			displayName = *discordUser.GlobalName
		}

		user, err := m.userStore.UpsertUser(c.Context(), &model.User{
			ID:            common.ID(discordUser.ID),
			Username:      discordUser.Username,
			Discriminator: discordUser.Discriminator,
			DisplayName:   displayName,
			Avatar:        common.PtrToNullString(discordUser.Avatar),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert user: %w", err)
		}

		discordGuilds, err := m.client.GetGuilds(discordSession, rest.WithCtx(c.Context()))
		if err != nil {
			return fmt.Errorf("failed to get user guilds: %w", err)
		}

		guilds := make([]model.SessionGuild, len(discordGuilds))
		for i, guild := range discordGuilds {
			guilds[i] = model.SessionGuild{
				ID:          common.ID(guild.ID),
				Name:        guild.Name,
				Icon:        common.PtrToNullString(guild.Icon),
				Owner:       guild.Owner,
				Permissions: guild.Permissions,
			}
		}

		scopes := make([]string, len(discordSession.Scopes))
		for i, scope := range discordSession.Scopes {
			scopes[i] = string(scope)
		}

		_, _, err = m.sessionManager.CreateSessionCookie(c, session.CreateSessionOpts{
			UserID:         user.ID,
			TokenAccess:    discordSession.AccessToken,
			TokenRefresh:   discordSession.RefreshToken,
			TokenScopes:    scopes,
			TokenExpiresAt: discordSession.Expiration,
			Guilds:         guilds,
		})
		if err != nil {
			return fmt.Errorf("failed to create session cookie: %w", err)
		}

		c.Redirect(m.getRedirectURI(c, guildID), http.StatusTemporaryRedirect)
		return nil
	})

	mx.Post("/auth/logout", api.Typed(func(c *api.Context) (*api.Empty, error) {
		if err := m.sessionManager.DeleteSession(c); err != nil {
			return nil, fmt.Errorf("failed to delete session: %w", err)
		}
		return &api.Empty{}, nil
	}))
}

func (m *AuthModule) Router() module.Router[AuthConfig] {
	return module.NewRouter[AuthConfig]()
}

func (m *AuthModule) getRedirectURI(c *api.Context, guildID string) string {
	redirectURI := m.redirectURI
	if redirectPath := c.Query("redirect"); redirectPath != "" {
		redirectURI += redirectPath
	} else {
		if guildID != "" {
			redirectURI += "/dashboard/" + guildID + "?quickstart"
		} else {
			redirectURI += "/dashboard"
		}
	}
	return redirectURI
}

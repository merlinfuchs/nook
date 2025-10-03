package logging

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/jellydator/ttlcache/v3"
	"github.com/merlinfuchs/nook/nook-service/module"
)

const ModuleID = "logging"

type LoggingModule struct {
	messageCache      *ttlcache.Cache[messageCacheKey, discord.Message]
	auditLogCollector *auditLogCollector
}

func NewLoggingModule() *LoggingModule {
	m := &LoggingModule{
		messageCache:      ttlcache.New(ttlcache.WithTTL[messageCacheKey, discord.Message](time.Minute * 3)),
		auditLogCollector: newAuditLogCollector(),
	}

	go m.messageCache.Start()

	return m
}

func (m *LoggingModule) ModuleID() string {
	return ModuleID
}

func (m *LoggingModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Logging",
		Description:    "Log your server activity like bans, kicks, joins, leaves, etc. to a channel. ",
		Icon:           "scroll-text",
		Internal:       false,
		DefaultEnabled: false,
		DefaultConfig:  defaultConfig,
		ConfigSchema:   configSchema,
		ConfigUISchema: configUISchema,
	}
}

func (m *LoggingModule) Router() module.Router[LoggingConfig] {
	return module.NewRouter[LoggingConfig]().
		Handle(m.auditLogCollector).
		Handle(module.ListenerFunc(m.handleMessageCreate)).
		Handle(module.ListenerFunc(m.handleMessageDelete)).
		Handle(module.ListenerFunc(m.handleMessageEdit)).
		Handle(module.ListenerFunc(m.handleMemberJoin)).
		Handle(module.ListenerFunc(m.handleMemberLeave)).
		Handle(module.ListenerFunc(m.handleChannelCreate)).
		Handle(module.ListenerFunc(m.handleChannelUpdate)).
		Handle(module.ListenerFunc(m.handleChannelDelete)).
		Handle(module.ListenerFunc(m.handleRoleCreate)).
		Handle(module.ListenerFunc(m.handleRoleUpdate)).
		Handle(module.ListenerFunc(m.handleRoleDelete)).
		Handle(module.ListenerFunc(m.handleGuildUpdate)).
		Handle(module.ListenerFunc(m.handleEmojiCreate)).
		Handle(module.ListenerFunc(m.handleEmojiUpdate)).
		Handle(module.ListenerFunc(m.handleEmojiDelete)).
		Handle(module.ListenerFunc(m.handleStickerCreate)).
		Handle(module.ListenerFunc(m.handleStickerUpdate)).
		Handle(module.ListenerFunc(m.handleStickerDelete))
}

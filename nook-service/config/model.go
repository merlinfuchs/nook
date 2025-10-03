package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/merlinfuchs/nook/nook-service/common"
)

type Config struct {
	Logging    LoggingConfig    `toml:"logging"`
	Database   DatabaseConfig   `toml:"database"`
	API        APIConfig        `toml:"api"`
	App        AppConfig        `toml:"app"`
	UserLimits UserLimitsConfig `toml:"user_limits"`
	Discord    DiscordConfig    `toml:"discord"`
	Billing    BillingConfig    `toml:"billing"`
	Defaults   DefaultsConfig   `toml:"defaults"`
}

func (cfg *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

func LoadConfig(basePath string) (*Config, error) {
	return loadConfig[*Config](basePath)
}

type DatabaseConfig struct {
	Postgres PostgresConfig `toml:"postgres"`
	S3       S3Config       `toml:"s3"`
}

type LoggingConfig struct {
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxAge     int    `toml:"max_age"`
	MaxBackups int    `toml:"max_backups"`
}

type PostgresConfig struct {
	Host     string `toml:"host" validate:"required"`
	Port     int    `toml:"port" validate:"required"`
	DBName   string `toml:"db_name" validate:"required"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password"`
}

type S3Config struct {
	Endpoint        string `toml:"endpoint" validate:"required"`
	AccessKeyID     string `toml:"access_key_id" validate:"required"`
	SecretAccessKey string `toml:"secret_access_key" validate:"required"`
	Secure          bool   `toml:"secure"`
	SSECKey         string `toml:"ssec_key"`
}

type APIConfig struct {
	Host          string `toml:"host" validate:"required"`
	Port          int    `toml:"port" validate:"required"`
	PublicBaseURL string `toml:"public_base_url" validate:"required"`
	SecureCookies bool   `toml:"secure_cookies"`
	StrictCookies bool   `toml:"strict_cookies"`
}

type AppConfig struct {
	PublicBaseURL string `toml:"public_base_url" validate:"required"`
}

type DiscordConfig struct {
	ClientID     common.ID `toml:"client_id" validate:"required"`
	ClientSecret string    `toml:"client_secret" validate:"required"`
	BotToken     string    `toml:"bot_token" validate:"required"`
	PublicKey    string    `toml:"public_key" validate:"required"`
}

type DefaultsConfig struct {
	CommandPrefix string `toml:"command_prefix" validate:"required"`
	ColorScheme   string `toml:"color_scheme" validate:"required"`
}

type UserLimitsConfig struct {
}

type BillingConfig struct {
	PaddleWebhookSecret string `toml:"paddle_webhook_secret"`
	PaddleAPIKey        string `toml:"paddle_api_key"`
	PaddleEnvironment   string `toml:"paddle_environment" validate:"required"`

	Plans []BillingPlanConfig `toml:"plans" validate:"required"`
}

type BillingPlanConfig struct {
	ID          string `toml:"id" validate:"required"`
	Title       string `toml:"title" validate:"required"`
	Description string `toml:"description" validate:"required"`
	Default     bool   `toml:"default"`
	Popular     bool   `toml:"popular"`
	Hidden      bool   `toml:"hidden"`

	PaddleMonthlyPriceID string `toml:"paddle_monthly_price_id"`
	PaddleYearlyPriceID  string `toml:"paddle_yearly_price_id"`

	DiscordRoleID string `toml:"discord_role_id"`

	FeatureBasicAccess bool `toml:"feature_basic_access"`
}

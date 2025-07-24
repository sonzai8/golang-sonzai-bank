package utils

import (
	"github.com/spf13/viper"
	"time"
)

// Config store all configuration of the application.
// the values a ready by viper
type Config struct {
	DbDriver             DatabaseConfig `mapstructure:"DB_DRIVER"`
	RedisConfig          RedisConfig    `mapstructure:"REDIS"`
	AppConfig            AppConfig      `mapstructure:"APP_CONFIG"`
	EmailConfig          EmailConfig    `mapstructure:"EMAIL_CONFIG"`
	TokenSymmetricKey    string         `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration  `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration  `mapstructure:"REFRESH_TOKEN_DURATION"`
}

type DatabaseConfig struct {
	Host    string `mapstructure:"DB_HOST"`
	Port    string `mapstructure:"DB_PORT"`
	User    string `mapstructure:"DB_USER"`
	Pass    string `mapstructure:"DB_PASSWORD"`
	Name    string `mapstructure:"DB_NAME"`
	SSLMode string `mapstructure:"DB_SSLMODE"`
}

type AppConfig struct {
	Environment  string `mapstructure:"ENVIRONMENT"`
	HttpPort     string `mapstructure:"HTTP_APP_PORT"`
	GrpcPort     string `mapstructure:"GRPC_APP_PORT"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
}

type RedisConfig struct {
	Address string `mapstructure:"ADDRESS"`
}

type EmailConfig struct {
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

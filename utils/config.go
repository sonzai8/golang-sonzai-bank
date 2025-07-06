package utils

import "github.com/spf13/viper"

// Config store all configuration of the application.
// the values a ready by viper
type Config struct {
	DbDriver  DatabaseConfig `mapstructure:"DB_DRIVER"`
	AppConfig AppConfig      `mapstructure:"APP_CONFIG"`
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
	Port string `mapstructure:"APP_PORT"`
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

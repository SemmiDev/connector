package g_learning_connector

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	DBConnection   string        `mapstructure:"DB_CONNECTION"`
	DBHost         string        `mapstructure:"DB_HOST"`
	DBPort         string        `mapstructure:"DB_PORT"`
	DBDatabase     string        `mapstructure:"DB_DATABASE"`
	DBUsername     string        `mapstructure:"DB_USERNAME"`
	DBPassword     string        `mapstructure:"DB_PASSWORD"`
	DBPoolIdle     int           `mapstructure:"DB_POOL_IDLE"`
	DBPoolMax      int           `mapstructure:"DB_POOL_MAX"`
	DBPoolLifetime time.Duration `mapstructure:"DB_POOL_LIFETIME"`
}

func NewConfig() (*Config, error) {
	viperConfig := viper.New()

	viperConfig.SetConfigFile("app.env")
	viperConfig.AddConfigPath(".")

	// read from environment variables
	viperConfig.AutomaticEnv()

	err := viperConfig.ReadInConfig()
	if err != nil {
		// if err is not the file not found, so return immedietly
		// and assume config load from environment variable
		if !strings.Contains(err.Error(), "no such file or directory") {
			return nil, err
		}
	}

	var config Config
	err = viperConfig.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

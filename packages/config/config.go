package config

import (
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		PgDSN         string `mapstructure:"PG_DSN"`
		ServerAddress string `mapstructure:"SERVER_ADDRESS"`
		NATSConfig    `mapstructure:",squash"`
	}
	NATSConfig struct {
		URL        string `mapstructure:"NATS_URL"`
		ClusterID  string `mapstructure:"NATS_CLUSTER_ID"`
		ClientID   string `mapstructure:"NATS_CLIENT_ID"`
		Subject    string `mapstructure:"NATS_SUBJECT"`
		QueueGroup string `mapstructure:"NATS_QUEUE_GROUP"`
		Durable    string `mapstructure:"NATS_DURABLE"`
	}
)

var (
	once   sync.Once
	config *Config
)

// Load reads configuration from file or environment variables.
func Load(path string) (err error) {
	once.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.Unmarshal(&config)
	})

	return
}

// Get .
func Get() *Config {
	return config
}

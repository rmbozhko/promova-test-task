package util

import (
	"github.com/spf13/viper"
	"log"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	MigrationURL  string `mapstructure:"MIGRATION_URL"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from a file or environment variables.
func LoadConfig(pathToConfigFile string) (config Config, err error) {
	viper.AddConfigPath(pathToConfigFile)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// overrides read values from config file with env vars if such exist
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	if err = viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("Loaded config: %+v", config)
	return
}

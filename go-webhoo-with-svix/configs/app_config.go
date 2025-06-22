package configs

import (
	"os"

	"github.com/spf13/viper"
)

// get svix_key from config file using viper
func InitConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs"
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	return err
}

func GetSvixKey() string {
	return viper.GetString("svix_key")
}

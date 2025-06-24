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

func GetSvixURL() string {
	return viper.GetString("svix_app_url")
}

// svix_app_id
func GetSvixAppID() string {
	return viper.GetString("svix_app_id")
}

// svix_app_name
func GetSvixAppName() string {
	return viper.GetString("svix_app_name")
}

// svix_signing_key
func GetSvixSigningKey() string {
	return viper.GetString("svix_signing_key")
}

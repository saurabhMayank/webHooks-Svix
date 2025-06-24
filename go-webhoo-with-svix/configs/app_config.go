package configs

import (
	"log"
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
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return err
	}
	log.Println("Config file loaded successfully")
	return nil
}

func GetSvixKey() string {
	svixKey := viper.GetString("svix_key")
	if svixKey == "" {
		log.Println("Warning: 'svix_key' not found in configuration, using default value")
		svixKey = "testsk_EhivPK7Gmn8ka4LLSTj5VhMr_3oWyZnS.eu" // Default value
	}
	return svixKey
}

func GetSvixURL() string {
	svixURL := viper.GetString("svix_app_url")
	if svixURL == "" {
		log.Println("Warning: 'svix_app_url' not found in configuration, using default value")
		svixURL = "https://api.eu.svix.com/api/v1/app/" // Default value
	}
	return svixURL
}

// svix_app_id
func GetSvixAppID() string {
	svixAppID := viper.GetString("svix_app_id")
	if svixAppID == "" {
		log.Println("Warning: 'svix_app_id' not found in configuration, using default value")
		svixAppID = "app_2yVDU5vouhEhfuW6OzOPdaTBvqI" // Default value
	}
	return svixAppID
}

// svix_app_name
func GetSvixAppName() string {
	svixAppName := viper.GetString("svix_app_name")
	if svixAppName == "" {
		log.Println("Warning: 'svix_app_name' not found in configuration, using default value")
		svixAppName = "SvixWebhookTestApp" // Default value
	}
	return svixAppName
}

// svix_signing_key
func GetSvixSigningKey() string {
	svixSigningKey := viper.GetString("svix_signing_key")
	if svixSigningKey == "" {
		log.Println("Warning: 'svix_signing_key' not found in configuration, using default value")
		svixSigningKey = "whsec_aaQUHgPtuFI7uBoVpZ6j2QcuWAlmpfR6" // Default value
	}
	return svixSigningKey
}

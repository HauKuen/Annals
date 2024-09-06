package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	AppMode  string
	HttpPort string
	JwtKey   string

	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
)

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	AppMode = viper.GetString("server.app_mode")
	HttpPort = viper.GetString("server.http_port")
	JwtKey = viper.GetString("server.jwt_key")
	Host = viper.GetString("mysql.host")
	Port = viper.GetString("mysql.port")
	User = viper.GetString("mysql.user")
	Password = viper.GetString("mysql.password")
	Dbname = viper.GetString("mysql.db_name")

	return validateConfig()
}

func validateConfig() error {
	requiredConfigs := map[string]string{
		"AppMode":  AppMode,
		"HttpPort": HttpPort,
		"JwtKey":   JwtKey,
		"Host":     Host,
		"Port":     Port,
		"User":     User,
		"Password": Password,
		"Dbname":   Dbname,
	}

	var missingConfigs []string
	for name, value := range requiredConfigs {
		if value == "" {
			missingConfigs = append(missingConfigs, name)
		}
	}

	if len(missingConfigs) > 0 {
		return fmt.Errorf("missing required configurations: %v", missingConfigs)
	}

	return nil
}

func init() {
	if err := LoadConfig(); err != nil {
		panic(err)
	}
}

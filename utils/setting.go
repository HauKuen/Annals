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

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	AppMode = viper.GetString("server.app_mode")
	HttpPort = viper.GetString("server.http_port")
	JwtKey = viper.GetString("server.jwt_key")
	Host = viper.GetString("mysql.host")
	Port = viper.GetString("mysql.port")
	User = viper.GetString("mysql.user")
	Password = viper.GetString("mysql.password")
	Dbname = viper.GetString("mysql.db_name")

	InitLogger()

	Log.Info("Configuration loaded successfully")
}

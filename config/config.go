package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Host string
	Port string
	User string
	Password string
	Database string
)


func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	Host = viper.GetString("host")
	if Host == "" {
		logrus.Fatal("host empty")
	}
	Port = viper.GetString("port")
	if Port == "" {
		logrus.Fatal("port empty")
	}
	User = viper.GetString("user")
	if User == "" {
		logrus.Fatal("user empty")
	}
	Password = viper.GetString("password")
	if Password == "" {
		logrus.Fatal("port empty")
	}
	Database = viper.GetString("database")
	if Database == "" {
		logrus.Fatal("database empty")
	}
}

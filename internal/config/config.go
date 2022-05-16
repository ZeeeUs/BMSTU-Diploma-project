package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type DbConfig struct {
	DbName     string `mapstructure:"POSTGRES_DB"`
	DbPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DbUser     string `mapstructure:"POSTGRES_USER"`
	DbHostName string `mapstructure:"DB_HOSTNAME"`
	DbPort     int    `mapstructure:"DB_PORT"`
}

type ServerConfig struct {
	Addr string
}

type TimeoutsConfig struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

type Config struct {
	DbConfig     DbConfig
	ServerConfig ServerConfig
	Timeouts     TimeoutsConfig
}

func SetupSource(cfgName string) {
	viper.SetConfigName(cfgName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
}

func NewConfig() *Config {
	SetupSource("config")

	return &Config{
		DbConfig: DbConfig{
			DbName:     viper.GetString("database.postgres_db"),
			DbPassword: viper.GetString("database.postgres_password"),
			DbUser:     viper.GetString("database.postgres_user"),
			DbHostName: viper.GetString("database.db_hostname"),
			DbPort:     viper.GetInt("database.port"),
		},
		ServerConfig: ServerConfig{
			Addr: viper.GetString("server.addr"),
		},
		Timeouts: TimeoutsConfig{
			WriteTimeout:   5 * time.Second,
			ReadTimeout:    5 * time.Second,
			ContextTimeout: time.Second * 2,
		},
	}
}

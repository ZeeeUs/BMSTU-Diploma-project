package main

import (
	"github.com/spf13/viper"
)

type DbConfig struct {
	DbName     string `mapstructure:"POSTGRES_DB"`
	DbPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DbUser     string `mapstructure:"POSTGRES_USER"`
	DbHostName string `mapstructure:"DB_HOSTNAME"`
	DbPort     int    `mapstructure:"DB_PORT"`
}

type Config struct {
	DbConfig DbConfig
}

func LoadConfig() (config Config, err error) {
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return
}

//func newConfig() *Config {
//	return &Config{
//		DbConfig: DbConfig,
//	}
//}

//func BuildConnString() string {
//	var config DbConfig
//	// urlExample := "postgres://username:password@localhost:5432/database_name"
//	connStr := fmt.Sprintf(
//		"postgresql://%s:%s@%s:%d/%s",
//		config.DbUser, config.DbPassword, config.DbHostName, config.DbPort, config.DbName,
//	)
//
//	return connStr
//}

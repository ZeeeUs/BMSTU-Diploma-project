package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func SetupSource(cfgName string) {
	viper.SetConfigName(cfgName)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	//viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
}

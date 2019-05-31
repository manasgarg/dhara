package main

import (
	"fmt"

	"github.com/manasgarg/dhara/apis"
	"github.com/manasgarg/dhara/stream"
	"github.com/manasgarg/dhara/utils"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("dhara")
	viper.AddConfigPath("$HOME/.dhara")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("dhara")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		utils.SLogger.Fatalw("Error in reading config", "err", err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	utils.InitLogger()
	stream.InitConfig()

	initConfig()

	apis.StartHTTPServer("0.0.0.0", 5678)
}

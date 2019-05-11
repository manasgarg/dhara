package main

import (
	"github.com/manasgarg/dhara/apis"
	"github.com/manasgarg/dhara/utils"
)

func main() {
	utils.InitLogger()
	apis.StartHTTPServer("0.0.0.0", 5678)
}

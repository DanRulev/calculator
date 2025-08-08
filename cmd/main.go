package main

import (
	"calculator-go/internal/app"
	"log"
)

const (
	CONFIG_PATH = "configs"
	CONFIG_FILE = "config"
)

func main() {
	log.Println("start")
	app.Start(CONFIG_PATH, CONFIG_FILE)
}

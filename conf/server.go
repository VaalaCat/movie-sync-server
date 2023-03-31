package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Setting struct {
	Port        string
	Addr        string
	AllowOrigin string
}

var ServerSetting Setting

func init() {
	// load from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	if port == "" {
		port = "3000"
	}
	if allowOrigin == "" {
		allowOrigin = "*"
	}

	ServerSetting = Setting{
		Port:        port,
		AllowOrigin: allowOrigin,
	}
}

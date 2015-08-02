package main

import (
	"flag"
	"helios/cors"
	"helios/github"
	"helios/helios"
	"helios/slack"
	"helios/static"
	"helios/weather"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	port          string
	usersFilename string // User Access Tokens file path
)

func flags() {
	flag.StringVar(&port, "port", "8989", "Port to run the server on")
	flag.StringVar(&usersFilename, "users", "users.csv", "Path to UAT (User Access Tokens) location. Used to store user github tokens")

	// Use env variables if they are defined
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	flag.Parse()
}

func main() {
	// Initialize command line args
	flags()

	h := helios.New()

	h.Use(cors.Service())
	h.Use(static.Service())
	h.Use(weather.Service())
	h.Use(github.Service())
	h.Use(slack.Service())

	// Initialize helios
	h.Run(port)
}

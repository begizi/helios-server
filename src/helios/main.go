package main

import (
	"flag"
	"helios/cors"
	"helios/github"
	"helios/helios"
	"helios/static"
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

	h.Use(cors.Plugin())
	h.Use(static.Plugin())
	h.Use(github.Plugin())

	// Initialize helios
	h.Run(port)
}

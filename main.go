package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	githubProvider "github.com/markbates/goth/providers/github"
	"github.com/tommy351/gin-cors"
)

const TIME_FORMAT = "2006-01-02T15:04:05Z"

var (
	githubKey     string
	githubSecret  string
	host          string
	port          string
	publicDir     string
	usersFilename string // User Access Tokens file path
	Users         []User
	LastEvent     Event
	usersFile     *os.File
)

var eventChan = make(chan []github.Event)

type Event struct {
	sync.RWMutex
	EventTime time.Time
}

func flags() {
	flag.StringVar(&publicDir, "public", "", "Path to your public assets folder")
	flag.StringVar(&port, "port", "8989", "Port to run the server on")
	flag.StringVar(&host, "host", "locahost", "Host the server is running on")
	flag.StringVar(&githubKey, "github-key", "", "Github oauth application key")
	flag.StringVar(&githubSecret, "github-secret", "", "Github oauth application secret")
	flag.StringVar(&usersFilename, "users", "users.csv", "Path to UAT (User Access Tokens) location. Used to store user github tokens")

	// Use env variables if they are defined
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	if len(os.Getenv("GITHUB_KEY")) > 0 {
		githubKey = os.Getenv("GITHUB_KEY")
	}

	if len(os.Getenv("GITHUB_SECRET")) > 0 {
		githubSecret = os.Getenv("GITHUB_SECRET")
	}

	flag.Parse()
}

func authorization() {
	// Setup Goth Authentication
	goth.UseProviders(
		githubProvider.New(githubKey, githubSecret, fmt.Sprintf("http://localhost:8989/auth/github/callback"), "repo", "user:email"),
	)
}

func loadUsersCSV() {
	var err error
	// Open and parse existing users from the uat file
	usersFile, err = os.OpenFile(usersFilename, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		log.Fatalf("Failed to open users file", err)
	}

	csvReader := csv.NewReader(usersFile)
	rawCSV, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file", err)
	}

	for _, row := range rawCSV {
		u := User{row[0], row[1]}
		Users = append(Users, u)
	}
}

func main() {
	// Initialize command line args
	flags()

	// Set github provider
	authorization()

	// Setup Socket Server
	socketServer, err := setupSocketIO()
	if err != nil {
		log.Fatalf("Failed to set up socket server", err)
	}

	err = startSocketPusher(socketServer, eventChan)
	if err != nil {
		log.Fatalf("Failed to start socket pusher go routine", err)
	}

	// Load registered users from csv
	loadUsersCSV()
	defer usersFile.Close()

	// Set the initial last event time to now
	LastEvent.EventTime = time.Now()

	// Start existing users go routines
	startExistingUsers(eventChan)

	// Create Engine Instance
	r := gin.Default()

	// Engine settings
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	r.HandleMethodNotAllowed = true

	// Middleware
	r.Use(cors.Middleware(cors.Options{AllowCredentials: true}))
	r.Use(static.Serve("/", static.LocalFile(publicDir, false)))

	// Handlers
	setupRoutes(r)

	// Start
	r.Run(":" + port)
}

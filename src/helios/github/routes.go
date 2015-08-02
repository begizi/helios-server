package github

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func providerAuth(c *gin.Context) {
	gothic.GetProviderName = func(req *http.Request) (string, error) { return "github", nil }
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func providerCallback(c *gin.Context) {
	var user User

	// Run user auth using the gothic library
	githubUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Fatalf("Failed to create user from callback", err)
	}

	user.Username = githubUser.RawData["login"].(string)
	user.AccessToken = githubUser.AccessToken

	// If the user doesn't exist yet
	if _, ok := Users[user.Username]; !ok {
		userFile, err := os.OpenFile("users.csv", os.O_APPEND|os.O_WRONLY, 0644)
		defer userFile.Close()

		_, err = userFile.WriteString(fmt.Sprintf("%s,%s\n", user.Username, user.AccessToken))
		if err != nil {
			log.Fatalf("Failed to write new users to CSV")
		}

		Users[user.Username] = user
		// startUser(user)

	} else {
		fmt.Println("User Already Exists")
	}

	c.JSON(200, user)
}

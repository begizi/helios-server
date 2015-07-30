package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func setupRoutes(r *gin.Engine) {
	//Oauth Authenticaton and Callbacks
	r.GET("/auth/github/callback", providerCallback)
	r.GET("/auth/github", providerAuth)
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

	c.JSON(200, user)
}

func providerAuth(c *gin.Context) {
	gothic.GetProviderName = func(req *http.Request) (string, error) { return "github", nil }
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

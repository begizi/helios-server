package github

import (
	"sync"
	"time"

	"golang.org/x/oauth2"
)

type User struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type Event struct {
	sync.RWMutex
	EventTime time.Time
}

type tokenSource struct {
	token *oauth2.Token
}

func (t tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

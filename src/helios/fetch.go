package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	token *oauth2.Token
}

func (t tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func startExistingUsers() {
	fmt.Println("Starting go routines")
	for _, u := range Users {
		startUser(u)
	}
}

func startUser(u User) {
	go userRoutine(u, eventChan)
}

func userRoutine(u User, c chan<- []github.Event) error {

	ts := tokenSource{
		&oauth2.Token{
			AccessToken: u.AccessToken,
		},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	//List Options Page, PerPage
	opts := github.ListOptions{1, 1}

	for {
		events, resp, err := client.Activity.ListEventsPerformedByUser(u.Username, false, &opts)
		if err != nil {
			log.Fatalf("Problem retrieving events for user", u.Username, err.Error())
		}

		newEventTime := *events[0].CreatedAt

		fmt.Println(newEventTime)

		// read in last event and compare to new event time
		LastEvent.Lock()
		dur := LastEvent.EventTime.Sub(newEventTime)
		if dur.Seconds() > 0.0 {
			LastEvent.EventTime = newEventTime
		}
		LastEvent.Unlock()

		c <- events

		// Wait as long as the X-Poll-Interval header says to
		interval, err := strconv.ParseInt(resp.Header["X-Poll-Interval"][0], 10, 8)
		if err != nil {
			// if strconv failed for whatever reason, use the default X-Poll-Interval value of 60
			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	panic("Shouldn't be here")
}

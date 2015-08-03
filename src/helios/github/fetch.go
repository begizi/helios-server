package github

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	log "gopkg.in/inconshreveable/log15.v2"
)

func startExistingUsers() {
	fmt.Println("Starting go routines")
	for _, u := range Users {
		startUser(u)
	}
}

func startUser(u User) {
	go userRoutine(u, EventChan)
}

func userRoutine(u User, c chan<- interface{}) error {

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
			log.Warn("Problem retrieving events for user", "username", u.Username, "error", err.Error())
		}

		newEventTime := LastEvent.EventTime

		if len(events) > 0 {
			newEventTime = *events[0].CreatedAt
		}

		// read in last event and compare to new event time
		LastEvent.Lock()
		dur := LastEvent.EventTime.Sub(newEventTime)
		if dur.Seconds() > 0.0 {
			LastEvent.EventTime = newEventTime
			c <- events[0]
		}
		LastEvent.Unlock()

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

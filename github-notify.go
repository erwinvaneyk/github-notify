package main

import (
	"fmt"
	"github.com/erwinvaneyk/go-pushbullet"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"os"
	"strconv"
	"strings"
	"time"
)

var notRead []github.Notification

func authGithub(apikey string) *github.Client {
	if apikey == "" {
		panic("No API-key set for Github!")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apikey},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	return client
}

func authPushbullet(apikey string) *pushbullet.Client {
	pb := pushbullet.New(apikey)
	return pb
}

func retrieveNotifications(client *github.Client) []github.Notification {
	opts := &github.NotificationListOptions{All: false, Participating: true}
	var results []github.Notification
	notes, _, _ := client.Activity.ListNotifications(opts)

	// Clean up notRead
	for i := len(notRead) - 1; i >= 0; i-- {
		if(!hasNotification(notes, notRead[i])) {
			notRead = append(notRead[:i], notRead[i+1:]...)
		}
	}

	for _, note := range notes {
		if *note.Reason == "mention" {
			if(!hasNotification(notRead, note)) {
				fmt.Printf("event: %s - %s (%s) in (%s)\n", *note.Reason, *note.Subject.Title, *note.ID, *note.Repository.Name)
				results = append(results, note)
			}
		}
	}
	return results
}

func pushMentionToPushBullet(pb *pushbullet.Client, mention github.Notification) {
	// URL rewrite hack
	url := *mention.Subject.URL
	url = strings.Replace(url, "api.", "", 1)
	url = strings.Replace(url, "pulls", "pull", 1)
	url = strings.Replace(url, "repos/", "", 1)
	pb.PushLink("", "mentioned in "+*mention.Repository.Name, url, *mention.Subject.Title)

	// Remember that the user is already notified
	notRead = append(notRead, mention)
}

func main() {
	client := authGithub(os.Getenv("GITHUB_API_KEY"))
	pb := authPushbullet(os.Getenv("PUSHBULLET_API_KEY"))
	interval, err := strconv.Atoi(os.Getenv("CHECK_INTERVAL"))
	if err != nil {
		interval = 300
		fmt.Printf("Invalid or no interval defined using default: %d seconds\n", interval)
	}

	for {
		println("Checking notifications...")
		notes := retrieveNotifications(client)
		// if there are mentions push one to pushbullet
		if len(notes) > 0 {
			fmt.Printf("Found %d mentions; pushing oldest one.\n", len(notes))
			pushMentionToPushBullet(pb, notes[len(notes)-1])
		} else {
			fmt.Println("No mentions found.")
		}
		fmt.Printf("Number of unread notifications: %d\n", len(notRead))
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func hasNotification(slice []github.Notification, item github.Notification) bool {
	for _, sliceItem := range slice {
		if(*sliceItem.ID == *item.ID) {
			return true;
		}
	}
	return false;
}
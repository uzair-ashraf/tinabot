package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	d             *discordgo.Session
	Version       string
	TwitterApiKey string
	TwitterSecret string
	BotToken      string
	tinaTweets    []string
	danID         = "498665578001858560"
	isPaused      bool
)

func init() {
	var err error
	d, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		panic(fmt.Errorf("invalid bot parameters: %v", err))
	}
}

func main() {
	tinaTweets = []string{}
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     TwitterApiKey,
		ClientSecret: TwitterSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)
	// Create Bot Handler
	d.AddHandler(messageCreate)

	// Open discord stream
	err := d.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	log.Println("tinabot is ready to send messages.")

	for {
		tinaTweetScraper(client)
		log.Println("Process sleeping for 60 minutes")
		time.Sleep(time.Minute * 60)
	}
}

func tinaTweetScraper(client *twitter.Client) {
	log.Println("Fetching Tina Tweets")
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          3098444432,
		ScreenName:      "tinaqtart",
		Count:           200,
		ExcludeReplies:  twitter.Bool(true),
		IncludeRetweets: twitter.Bool(false),
		TrimUser:        twitter.Bool(true),
	})
	if err != nil {
		log.Println(fmt.Errorf("Could not get tina tweets Error: %v", err))
		return
	}
	tinaTweetsNew := []string{}
	for _, tweet := range tweets {
		if tweet.ExtendedEntities == nil {
			continue
		}
		for _, pics := range tweet.ExtendedEntities.Media {
			if pics.Type == "photo" {
				if pics.MediaURLHttps != "" {
					tinaTweetsNew = append(tinaTweetsNew, pics.MediaURLHttps)
				}
			}
		}
	}
	tinaTweets = tinaTweetsNew
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	args := strings.Split(m.Content, " ")
	command := args[0]
	if containsTina(args) {
		handleTina(s, m)
		return
	}
	if command != "!tina" || len(args) == 1 {
		return
	}
	if args[1] == "help" {
		handleHelp(s, m)
		return
	}
	if args[1] == "shutup" {
		handleStop(s, m)
		return
	}
	if args[1] == "continue" {
		handleStart(s, m)
		return
	}

}

func handleTina(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isPaused {
		return
	}
	image := getRandomTinaPic(tinaTweets)
	s.ChannelMessageSend(m.ChannelID,
		"Follow @tinaqtart on Twitter"+
			"\n"+
			image,
	)
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID,
		"Welcome to the Tina Discord Bot\n"+
			"This application was developed by Shadoki#0001\n"+
			"Verison: "+Version+"\n"+
			"!tina help command shows you the available commands\n"+
			"!tina continue command lets you continue tinabot listening\n"+
			"!tina shutup command lets you silence tinabot",
	)
}

func handleStart(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == danID {
		s.ChannelMessageSend(m.ChannelID, "No Dan.")
		return
	}
	if isPaused {
		s.ChannelMessageSend(m.ChannelID, "Okie I'll be listening.")
		isPaused = false
	} else {
		s.ChannelMessageSend(m.ChannelID, "Bro I'm already working what do you want from me.")
	}
}

func handleStop(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == danID {
		s.ChannelMessageSend(m.ChannelID, "Fuck off Dan don't tell me to shutup.")
		return
	}
	if isPaused {
		s.ChannelMessageSend(m.ChannelID, "Bro I'm already quiet stop.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Fine I'll be quiet.")
		isPaused = true
	}
}

func containsTina(s []string) bool {
	tinaMap := map[string]bool{
		"tina":   true,
		"tina!":  true,
		"tina?":  true,
		"tina.":  true,
		"tina's": true,
	}
	for _, a := range s {
		if _, ok := tinaMap[strings.ToLower(a)]; ok {
			return true
		}
	}
	return false
}

func getRandomTinaPic(s []string) string {
	randomIndex := rand.Intn(len(s))
	pick := s[randomIndex]
	return pick
}

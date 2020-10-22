package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

const ApiKey = "~Api Key~"
const ApiKeySecret = "~Api Key Secret~"
const AccessToken = "~Access Token~"
const AccessTokenSecret = "~Access Token Secret~"
const AspEnvVariableName = "ASPNETCORE_PORT"

type NewsTweet struct {
	Source string `json:"source"`
	Title  string `json:"title"`
	Link   string `json:"link"`
}

func main() {

	port := "8000"
	r := mux.NewRouter()
	r.HandleFunc("/", TweetHandler).Methods("POST")

	fmt.Printf("Server listening on port 8000... ")

	if os.Getenv(AspEnvVariableName) != "" {
		port = os.Getenv(AspEnvVariableName)
	}

	server := http.ListenAndServe(":"+port, r)

	log.Fatal(server)

}

func TweetHandler(w http.ResponseWriter, r *http.Request) {

	var newsBody NewsTweet

	// Decode the body
	err := json.NewDecoder(r.Body).Decode(&newsBody)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// Set up keys and secrets
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", ApiKey, "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", ApiKeySecret, "Twitter Consumer Secret")
	accessToken := flags.String("access-token", AccessToken, "Twitter Access Token")
	accessSecret := flags.String("access-secret", AccessTokenSecret, "Twitter Access Secret")

	flags.Parse(os.Args[1:])

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Format the tweets body
	formattedTweet := buildTweet(newsBody)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update(formattedTweet, nil)

	if err != nil {
		log.Fatalf(err.Error())
	}

	createdAt := tweet.CreatedAt
	fmt.Fprintf(w, "Tweet succesful %s - %s", resp.Status, createdAt)

}

func buildTweet(tweet NewsTweet) string {

	source := FormatSourceName(tweet)

	return fmt.Sprintf("%s\n%s\n#%s %s",
		strings.ToUpper(tweet.Source),
		tweet.Title,
		source,
		tweet.Link)

}

func FormatSourceName(tweet NewsTweet) string {

	source := ""

	if strings.ContainsAny(tweet.Source, "-") {
		source = strings.ReplaceAll(tweet.Source, "-", "")
	} else {
		source = strings.ReplaceAll(tweet.Source, " ", "")
	}
	return source
}

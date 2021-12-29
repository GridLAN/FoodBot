package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used at runtime
var (
	Token string = os.Getenv("TOKEN")
)

// Meal is a struct that holds the meal data
type Meal struct {
	Meals []struct {
		Name    string `json:"strMeal"`
		Thumb   string `json:"strMealThumb"`
		Youtube string `json:"strYoutube"`
	} `json:"meals"`
}

func randomMealFunc() (randomMeal Meal) {

	var (
		mealUrl string = "https://www.themealdb.com/api/json/v1/1/random.php"
	)

	// Create a new HTTP client with a timeout of 1 second
	var mealClient = http.Client{Timeout: 10 * time.Second}

	// Build a GET request to the meal API endpoint
	req, err := http.NewRequest(http.MethodGet, mealUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Discord FoodBot")

	// Send the request and store the response
	res, getErr := mealClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Close body
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Read body
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// Unmarshal body
	randomMeal = Meal{}
	jsonErr := json.Unmarshal(body, &randomMeal)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return
}

func main() {
	// Check for token
	if Token == "" {
		log.Fatal("Token cannot be empty")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "food help" reply with a list of commands
	if m.Content == ".food help" || m.Content == ".Food help" {

		s.ChannelMessageSend(m.ChannelID, "Food Bot Commands:\n'.food help' - This message\n.'food random' - Random meal")
	}

	// If the message is "food" reply with a random dish to delight!
	if m.Content == ".food" || m.Content == ".Food" {

		randomMeal := randomMealFunc()

		s.ChannelMessageSend(m.ChannelID, ""+randomMeal.Meals[0].Name+"\n"+randomMeal.Meals[0].Thumb)
	}

}

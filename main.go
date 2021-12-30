package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used at runtime
var (
	Token            string = os.Getenv("TOKEN")
	MealAPI          string = "https://www.themealdb.com/api/json/v1/1"
	MealCategoryList []string
	MealOriginList   []string
)

// Meal is a struct that holds the meal data
type Meal struct {
	Meals []struct {
		Name    string `json:"strMeal"`
		Thumb   string `json:"strMealThumb"`
		Youtube string `json:"strYoutube"`
	} `json:"meals"`
}

type MealAPIMetadata struct {
	Data []struct {
		Origin   string `json:"strArea"`
		Category string `json:"strCategory"`
	} `json:"meals"`
}

func getMeal(MealAPI string) (randomMeal Meal) {

	// Create a new HTTP client with a timeout of 1 second
	var mealClient = http.Client{Timeout: 10 * time.Second}

	// Build a GET request to the meal API endpoint
	req, err := http.NewRequest(http.MethodGet, MealAPI, nil)
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

func getMealAPIMetadata(MealAPI string) (mealAPIMetadata MealAPIMetadata) {

	// Create a new HTTP client with a timeout of 1 second
	var mealClient = http.Client{Timeout: 10 * time.Second}

	// Build a GET request to the meal API endpoint
	req, err := http.NewRequest(http.MethodGet, MealAPI, nil)
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
	mealAPIMetadata = MealAPIMetadata{}
	jsonErr := json.Unmarshal(body, &mealAPIMetadata)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return
}

func init() {
	// Populate meal categories
	// Get a list of meal categories from the API
	mealCategories := getMealAPIMetadata(MealAPI + "/list.php?c=list")

	// Construct a list of meal categories
	for _, mealCategory := range mealCategories.Data {
		MealCategoryList = append(MealCategoryList, mealCategory.Category)
	}

	// Populate meal origins
	// Get a list of meal origins from the API
	mealOrigins := getMealAPIMetadata(MealAPI + "/list.php?a=list")

	// Construct a list of meal origins
	for _, mealOrigin := range mealOrigins.Data {
		MealOriginList = append(MealOriginList, mealOrigin.Origin)
	}
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

	// Print MealCategoryList
	fmt.Println(MealCategoryList)

	// Print MealOriginsList
	fmt.Println(MealOriginList)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Convert all messages content to lowercase
	m.Content = strings.ToLower(m.Content)

	// Switch block for message content and response
	switch m.Content {
	case "food help":
		s.ChannelMessageSend(m.ChannelID, "```It's time to eat, fatass!\n\nfood - Random meal\nfood categories - List of meal categories\nfood origins - List of meal origins```")

	case "food":
		// Set API endpoint to random
		apiEndpoint := "/random.php"

		// Get a random meal from the API
		randomMeal := getMeal(MealAPI + apiEndpoint)

		// Send the message to the channel
		s.ChannelMessageSend(m.ChannelID, ""+randomMeal.Meals[0].Name+"\n"+randomMeal.Meals[0].Thumb)

	case "food origins":
		// Set API endpoint to meal areas
		apiEndpoint := "/list.php?a=list"

		// Get a list of meal areas from the API
		mealOrigins := getMealAPIMetadata(MealAPI + apiEndpoint)

		// Construct a list of meal areas
		var mealOriginList []string
		for _, mealArea := range mealOrigins.Data {
			mealOriginList = append(mealOriginList, mealArea.Origin)
		}

		// Send the message to the channel
		s.ChannelMessageSend(m.ChannelID, "```"+"Available Meal Origins:\n\n"+strings.Join(mealOriginList, "\n")+"```")

	case "food categories":
		// Set API endpoint to meal categories
		apiEndpoint := "/list.php?c=list"

		// Get a list of meal categories from the API
		mealCategories := getMealAPIMetadata(MealAPI + apiEndpoint)

		// Construct a list of meal categories
		var mealCategoryList []string
		for _, mealCategory := range mealCategories.Data {
			mealCategoryList = append(mealCategoryList, mealCategory.Category)
		}

		// Send the message to the channel
		s.ChannelMessageSend(m.ChannelID, "```"+"Available Meal Categories:\n\n"+strings.Join(mealCategoryList, "\n")+"```")
	}

}

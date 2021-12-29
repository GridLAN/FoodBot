package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var mealUrl = "https://www.themealdb.com/api/json/v1/1/random.php"

// Meal is a struct that holds the meal data
type Meal struct {
	Meals []struct {
		Name     string `json:"strMeal"`
		Category string `json:"strCategory"`
		Area     string `json:"strArea"`
		Thumb    string `json:"strMealThumb"`
		Youtube  string `json:"strYoutube"`
	} `json:"meals"`
}

func main() {

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
	randomMeal := Meal{}
	jsonErr := json.Unmarshal(body, &randomMeal)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// Print & return meal name
	fmt.Println(randomMeal.Meals[0].Name)

	// Print & return meal category
	fmt.Println(randomMeal.Meals[0].Category)

	// Print & return meal Area
	fmt.Println(randomMeal.Meals[0].Area)

	// Print & return meal Thumb
	fmt.Println(randomMeal.Meals[0].Thumb)

	// Print & return meal Youtube link
	fmt.Println(randomMeal.Meals[0].Youtube)
}

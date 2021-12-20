package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type Foo struct {
	Bar string
}

func main() {
	var foo1 string
	getJson("https://www.themealdb.com/api/json/v1/1/random.php", foo1)
	println(foo1)

	// alternately:

	// foo2 := Foo{}
	// getJson("http://example.com", &foo2)
	// println(foo2.Bar)
}

// // Send HTTP GET request
// resp, err := http.Get("https://www.themealdb.com/api/json/v1/1/random.php")
// if err != nil {
// 	log.Fatalln(err)
// }

// // Read the response body
// bodyBytes, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	log.Fatalln(err)
// }

// // Convert the bodyBytes to type string
// sb := string(bodyBytes)
// log.Printf(sb)

// var responseObject MealResponse
// json.Unmarshal(bodyBytes, &responseObject)

// fmt.Println(responseObject)

// var meal Meal
// json.Unmarshal([]byte(body), &meal)

// fmt.Printf("Meal: %s", meal.Meal)

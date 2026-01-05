package main

import (
	"fmt"
	"math/rand/v2"
	"net/url"
)

func random_number_within_range(min, max int) int {
	return min + rand.IntN(max-min+1)
}

func generateRandomString() string {
	var alphabet = [35]string{"a", "b", "c", "d", "e", "f", "g", "h", "i",
		"j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "6", "7", "8", "9"}
	random_string := ""
	min := 0
	max := len(alphabet) - 1
	for i := 0; i < 7; i++ {
		random_index := random_number_within_range(min, max)
		random_string += alphabet[random_index]
	}

	return random_string

}

func main() {
	fmt.Println("URL Shortener Service")

	u, err := url.Parse("https://www.example.com/path/to/resource?param1=value1&param2=value2#fragment")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Type  %T", u)
	fmt.Println("Current URL: ", u)
	fmt.Println("Protocol: ", u.Scheme)
	fmt.Println("Hostname: ", u.Hostname())
	fmt.Println("Path: ", u.Path)
	fmt.Println("Raw query: ", u.RawQuery)
	fmt.Println("Fragment: ", u.Fragment)

	get_random_string := generateRandomString()

	fmt.Println(get_random_string)

	construct_new_url := url.URL{
		Scheme: "http",
		Host:   "localhost:8000",
		Path:   get_random_string,
	}

	fmt.Println(construct_new_url.String())


	

}

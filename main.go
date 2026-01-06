package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var defaultScheme string = "http"
var defaultHost string = "localhost:8000"

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

func getHostNameandScheme() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file")
	}
	if hostname := os.Getenv("APP_HOST_NAME"); hostname != "" {
		defaultHost = hostname
	}
	if scheme := os.Getenv("APP_SCHEME"); scheme != "" {
		defaultScheme = scheme
	}

	return defaultHost, defaultScheme
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("URL shortner failed due to ", r)
			os.Exit(1)
		}
	}()
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
	host, scheme := getHostNameandScheme()
	fmt.Println(get_random_string)

	construct_new_url := url.URL{
		Scheme: host,
		Host:   scheme,
		Path:   get_random_string,
	}
	http.HandleFunc("/"+get_random_string, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	})
	fmt.Println(construct_new_url.String())

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}

}

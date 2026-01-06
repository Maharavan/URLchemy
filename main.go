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

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func random_number_within_range(min, max int) int {
	return min + rand.IntN(max-min+1)
}

func base62Encoder(min, max int) string {
	number := random_number_within_range(min, max)
	base62 := ""

	for number > 0 {
		remainder := number % 62
		base62 += string(base62Digits[remainder])
		number /= 62
	}
	return base62
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

	get_random_string := base62Encoder(10000, 80000)
	host, scheme := getHostNameandScheme()
	fmt.Println(host, scheme, get_random_string)
	construct_new_url := url.URL{
		Scheme: scheme,
		Host:   host,
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

package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomBytes() []byte {

	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return b
}

func base62Encoder() string {
	number := generateRandomBytes()

	for i := range number {
		number[i] = base62Digits[int(number[i])%62]
	}

	return string(number)
}

func getHostNameandScheme() (string, string) {
	var defaultScheme string = "http"
	var defaultHost string = "localhost:8000"

	godotenv.Load()

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

	get_random_string := base62Encoder()
	host, scheme := getHostNameandScheme()
	construct_new_url := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   get_random_string,
	}
	http.HandleFunc("/"+get_random_string, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	})
	fmt.Println("Redirected URL:", construct_new_url.String())

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}

}

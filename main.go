package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Payload struct {
	URL string `json:"url"`
}

var (
	inmemory = make(map[string]string)
	mu       sync.RWMutex
)

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
	var defaultHost string = "127.0.0.1:8000"

	if hostname := os.Getenv("APP_HOST_NAME"); hostname != "" {
		defaultHost = hostname
	}
	if scheme := os.Getenv("APP_SCHEME"); scheme != "" {
		defaultScheme = scheme
	}

	return defaultHost, defaultScheme
}

func retrievelongurl(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("URL shortner failed due to ", r)
			panic(r)
		}
	}()

	var p Payload

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	dummy_url := p.URL
	u, err := url.Parse(dummy_url)
	if err != nil {
		panic(err)
	}

	var get_random_string string

	fmt.Printf("Type  %T", u)
	fmt.Println("Current URL: ", u)
	fmt.Println("Protocol: ", u.Scheme)
	fmt.Println("Hostname: ", u.Hostname())
	fmt.Println("Path: ", u.Path)
	fmt.Println("Raw query: ", u.RawQuery)
	fmt.Println("Fragment: ", u.Fragment)
	hostname, scheme := getHostNameandScheme()
	for {
		get_random_string = base62Encoder()
		if _, ok := inmemory[get_random_string]; !ok {
			mu.Lock()
			inmemory[get_random_string] = u.String()
			mu.Unlock()
			break
		}

	}

	construct_new_url := url.URL{
		Scheme: scheme,
		Host:   hostname,
		Path:   get_random_string,
	}
	url_data := Payload{
		URL: construct_new_url.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(url_data)
}

func rerouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println(r.URL)
	random := strings.TrimPrefix(r.URL.Path, "/")
	if random == "" {
		http.Error(w, fmt.Sprintf("%v", inmemory), http.StatusNotFound)
		return
	}

	mu.RLock()

	longurl, ok := inmemory[random]
	mu.RUnlock()

	if !ok {
		http.Error(w, "Invalid data", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longurl, http.StatusMovedPermanently)

}

func main() {
	godotenv.Load()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /longurl", retrievelongurl)
	mux.HandleFunc("GET /", rerouter)
	log.Print("URL Shortner Service.....")
	log.Print("Connecting to Port")
	port := ":8000"

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}

}

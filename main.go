package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Payload struct {
	URL string `json:"url"`
}

type RedisCache struct {
	client *redis.Client
}

var (
	inmemory = make(map[string]string)
	mu       sync.RWMutex
	cache    *RedisCache
	ctx      = context.Background()
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
			log.Println("URL shortner failed due to ", r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	log.Printf("Type  %T", u)
	log.Println("Current URL: ", u)
	log.Println("Protocol: ", u.Scheme)
	log.Println("Hostname: ", u.Hostname())
	log.Println("Path: ", u.Path)
	log.Println("Raw query: ", u.RawQuery)
	log.Println("Fragment: ", u.Fragment)
	hostname, scheme := getHostNameandScheme()
	for {
		get_random_string = base62Encoder()
		if _, err := cache.client.Get(ctx, get_random_string).Result(); err != nil {
			mu.Lock()
			if err := cache.client.Set(ctx, get_random_string, u.String(), 0).Err(); err != nil {
				log.Fatal(err)
			}
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
	log.Println(r.URL)
	random := strings.TrimPrefix(r.URL.Path, "/")
	if random == "" {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	mu.RLock()

	longurl, err := cache.client.Get(ctx, random).Result()
	mu.RUnlock()

	if err != nil {
		http.Error(w, "Invalid data", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longurl, http.StatusMovedPermanently)

}

func redisconnection() {

	cache = &RedisCache{client: redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "",
		DB:       0,
	})}

	status, err := cache.client.Ping(ctx).Result()

	if err != nil {
		log.Fatal("Error connecting to redis", err)
	}

	log.Println("Connected to Redis:", status)
}

func main() {
	godotenv.Load()
	redisconnection()
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

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	shortURLs map[string]string
	mutex     sync.Mutex
	prefix    = "https://swift.ly/" // Overridden using the SWIFTLY_PREFIX environment variable
	port      = "8080"              // Overridden using the SWIFTLY_PORT environment variable
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func main() {
	shortURLs = make(map[string]string)

	// environment variable overrides
	if prefixEnv := os.Getenv("SWIFTLY_PREFIX"); prefixEnv != "" {
		prefix = prefixEnv
	}
	if portEnv := os.Getenv("SWIFTLY_PORT"); portEnv != "" {
		port = portEnv
	}

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)

	log.Println("SwiftLink is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, ok := shortURLs[shortURL]
	if ok {
		log.Printf("Redirecting %s to %s\n", shortURL, longURL)
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		log.Printf("Short URL %s not found\n", shortURL)
		http.NotFound(w, r)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the ShortenRequest parameter from the json body
	var req ShortenRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(req.URL)

	mutex.Lock()
	shortURLs[shortURL] = req.URL
	mutex.Unlock()
	log.Printf("Updated Short URL: %s --> Long URL: %s\n", shortURL, req.URL)

	// Return the short URL in the response

	resp := ShortenResponse{ShortURL: prefix + shortURL}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func generateShortURL(longURL string) string {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the longURL already exists in the map
	for k, v := range shortURLs {
		if v == longURL {
			return k
		}
	}

	// Generate a new hash that is not in the map
	uniqueHash := generateUniqueHash(longURL)

	return uniqueHash
}

func generateUniqueHash(longURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(longURL))
	sha := hex.EncodeToString(hasher.Sum(nil))

	return sha[:8] // Take the first 8 characters as the short URL
}

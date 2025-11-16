package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var store = NewScoreStore()

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]interface{}{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(resp)
}

func addScoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := os.Getenv("SECRET")
	headerSecret := r.Header.Get("android-secret")

	if headerSecret != secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type reqBody struct {
		Player string `json:"player"`
		Score  int    `json:"score"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.Player == "" {
		http.Error(w, "Player is required", http.StatusBadRequest)
		return
	}

	if body.Score < 0 {
		http.Error(w, "Score must be positive", http.StatusBadRequest)
		return
	}

	store.SetScore(body.Player, body.Score)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"player": body.Player,
		"score":  body.Score,
	})
}

func getRatingsHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	result := store.GetRating(limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("No .env file found")
	}

	http.HandleFunc("/health", health)

	http.HandleFunc("/rating", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addScoreHandler(w, r)
		} else if r.Method == http.MethodGet {
			getRatingsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :3000")
	http.ListenAndServe(":3000", nil)
}

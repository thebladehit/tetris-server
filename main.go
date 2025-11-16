package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var store = NewRatingStore()

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]interface{}{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(resp)
}

func addRatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type reqBody struct {
		Player string `json:"player"`
		Rating int    `json:"rating"`
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

	store.SetRating(body.Player, body.Rating)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"player": body.Player,
		"rating": body.Rating,
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
	http.HandleFunc("/health", health)

	store.SetRating("Bogdan", 1500)
	store.SetRating("Anna", 1600)
	store.SetRating("Mike", 1400)
	store.SetRating("Sara", 1700)

	http.HandleFunc("/rating", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addRatingHandler(w, r)
		} else if r.Method == http.MethodGet {
			getRatingsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :3000")
	http.ListenAndServe(":3000", nil)
}

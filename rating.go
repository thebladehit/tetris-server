package main

import (
	"sort"
	"sync"
)

type RatingStore struct {
	mu      sync.RWMutex
	ratings map[string]int
}

func NewRatingStore() *RatingStore {
	return &RatingStore{
		ratings: make(map[string]int),
	}
}

func (s *RatingStore) SetRating(player string, rating int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ratings[player] = rating
}

func (s *RatingStore) GetRating(limit int) []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]map[string]interface{}, 0, len(s.ratings))
	for player, rating := range s.ratings {
		list = append(list, map[string]interface{}{
			"player": player,
			"rating": rating,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["rating"].(int) > list[j]["rating"].(int)
	})

	if limit > len(list) {
		limit = len(list)
	}

	return list[:limit]
}

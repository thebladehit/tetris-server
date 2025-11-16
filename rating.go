package main

import (
	"sort"
	"sync"
)

type ScoreStore struct {
	mu     sync.RWMutex
	scores map[string]int
}

func NewScoreStore() *ScoreStore {
	return &ScoreStore{
		scores: make(map[string]int),
	}
}

func (s *ScoreStore) SetScore(player string, rating int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores[player] = rating
}

func (s *ScoreStore) GetRating(limit int) []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]map[string]interface{}, 0, len(s.scores))
	for player, rating := range s.scores {
		list = append(list, map[string]interface{}{
			"player": player,
			"score":  rating,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["score"].(int) > list[j]["score"].(int)
	})

	if limit > len(list) {
		limit = len(list)
	}

	return list[:limit]
}

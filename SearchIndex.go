package arn

import (
	"sort"
	"strings"

	"github.com/aerogo/flow"
)

// SearchIndex ...
type SearchIndex struct {
	TextToID map[string]string `json:"textToId"`
}

// NewSearchIndex ...
func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		TextToID: make(map[string]string),
	}
}

// GetSearchIndex ...
func GetSearchIndex(id string) (*SearchIndex, error) {
	obj, err := DB.Get("SearchIndex", id)
	return obj.(*SearchIndex), err
}

// Search is a fuzzy search.
func Search(term string, maxUsers, maxAnime int) ([]*User, []*Anime) {
	term = strings.ToLower(term)

	if term == "" {
		return nil, nil
	}

	var userResults []*User
	var animeResults []*Anime

	type SearchItem struct {
		text       string
		similarity float64
	}

	// Search everything in parallel
	flow.Parallel(func() {
		// Search userResults
		var user *User

		userSearchIndex, err := GetSearchIndex("User")

		if err != nil {
			return
		}

		textToID := userSearchIndex.TextToID

		// Search items
		items := make([]*SearchItem, 0)

		for name := range textToID {
			s := StringSimilarity(term, name)

			if s < MinimumStringSimilarity {
				continue
			}

			items = append(items, &SearchItem{
				text:       name,
				similarity: s,
			})
		}

		// Sort
		sort.Slice(items, func(i, j int) bool {
			return items[i].similarity > items[j].similarity
		})

		// Limit
		if len(items) >= maxUsers {
			items = items[:maxUsers]
		}

		// Fetch data
		for _, item := range items {
			user, err = GetUser(textToID[item.text])

			if err != nil {
				continue
			}

			userResults = append(userResults, user)
		}
	}, func() {
		// Search anime
		var anime *Anime

		animeSearchIndex, err := GetSearchIndex("Anime")

		if err != nil {
			return
		}

		textToID := animeSearchIndex.TextToID

		// Search items
		items := make([]*SearchItem, 0)

		for name := range textToID {
			s := StringSimilarity(term, name)

			if s < MinimumStringSimilarity {
				continue
			}

			items = append(items, &SearchItem{
				text:       name,
				similarity: s,
			})
		}

		// Sort
		sort.Slice(items, func(i, j int) bool {
			return items[i].similarity > items[j].similarity
		})

		// Limit
		if len(items) >= maxAnime {
			items = items[:maxAnime]
		}

		// Fetch data
		for _, item := range items {
			anime, err = GetAnime(textToID[item.text])

			if err != nil {
				continue
			}

			animeResults = append(animeResults, anime)
		}
	})

	return userResults, animeResults
}

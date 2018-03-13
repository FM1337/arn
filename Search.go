package arn

import (
	"sort"
	"strings"

	"github.com/aerogo/flow"
)

// MinimumStringSimilarity is the minimum JaroWinkler distance we accept for search results.
const MinimumStringSimilarity = 0.89

// popularityDamping reduces the factor of popularity in search results.
const popularityDamping = 0.048

// SearchResult ...
type SearchResult struct {
	obj        interface{}
	similarity float64
}

// Search is a fuzzy search.
func Search(term string, maxUsers, maxAnime, maxPosts, maxThreads, maxTracks int, maxCharacters int) ([]*User, []*Anime, []*Post, []*Thread, []*SoundTrack, []*Character) {
	if term == "" {
		return nil, nil, nil, nil, nil, nil
	}

	var userResults []*User
	var animeResults []*Anime
	var postResults []*Post
	var threadResults []*Thread
	var trackResults []*SoundTrack
	var characterResults []*Character

	flow.Parallel(func() {
		userResults = SearchUsers(term, maxUsers)
	}, func() {
		animeResults = SearchAnime(term, maxAnime)
	}, func() {
		postResults = SearchPosts(term, maxPosts)
	}, func() {
		threadResults = SearchThreads(term, maxThreads)
	}, func() {
		trackResults = SearchSoundTracks(term, maxTracks)
	}, func() {
		characterResults = SearchCharacters(term, maxCharacters)
	})

	return userResults, animeResults, postResults, threadResults, trackResults, characterResults
}

// SearchCharacters searches all characters.
func SearchCharacters(originalTerm string, maxLength int) []*Character {
	if maxLength == 0 {
		return nil
	}

	term := RemoveSpecialCharacters(strings.ToLower(originalTerm))

	var results []*SearchResult

	for character := range StreamCharacters() {
		if character.ID == originalTerm {
			return []*Character{character}
		}

		if character.Image == "" {
			continue
		}

		text := strings.ToLower(character.Name)

		if text == term {
			results = append(results, &SearchResult{
				obj:        character,
				similarity: 1000,
			})
			continue
		}

		for index, char := range text {
			if char == ' ' {
				firstName := text[:index]
				lastName := text[index+1:]

				if firstName == term {
					results = append(results, &SearchResult{
						obj:        character,
						similarity: 10.0,
					})
					break
				}

				if lastName == term {
					results = append(results, &SearchResult{
						obj:        character,
						similarity: 1,
					})
					break
				}
			}
		}
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		similarityA := results[i].similarity
		similarityB := results[j].similarity

		if similarityA == similarityB {
			characterA := results[i].obj.(*Character)
			characterB := results[j].obj.(*Character)

			if characterA.Name == characterB.Name {
				return characterA.ID < characterB.ID
			}

			return characterA.Name < characterB.Name
		}

		return similarityA > similarityB
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*Character, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*Character)
	}

	return final
}

// SearchSoundTracks searches all soundtracks.
func SearchSoundTracks(originalTerm string, maxLength int) []*SoundTrack {
	term := RemoveSpecialCharacters(strings.ToLower(originalTerm))

	var results []*SearchResult

	for track := range StreamSoundTracks() {
		if track.ID == originalTerm {
			return []*SoundTrack{track}
		}

		text := strings.ToLower(track.Title)

		// Similarity check
		similarity := AdvancedStringSimilarity(term, text)

		if similarity >= MinimumStringSimilarity {
			results = append(results, &SearchResult{
				obj:        track,
				similarity: similarity,
			})
			continue
		}
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*SoundTrack, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*SoundTrack)
	}

	return final
}

// SearchPosts searches all posts.
func SearchPosts(originalTerm string, maxLength int) []*Post {
	term := RemoveSpecialCharacters(strings.ToLower(originalTerm))

	var results []*Post

	for post := range StreamPosts() {
		if post.ID == originalTerm {
			return []*Post{post}
		}

		text := strings.ToLower(post.Text)

		if !strings.Contains(text, term) {
			continue
		}

		results = append(results, post)
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Created > results[j].Created
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	return results
}

// SearchThreads searches all threads.
func SearchThreads(originalTerm string, maxLength int) []*Thread {
	term := RemoveSpecialCharacters(strings.ToLower(originalTerm))

	var results []*Thread

	for thread := range StreamThreads() {
		if thread.ID == originalTerm {
			return []*Thread{thread}
		}

		text := strings.ToLower(thread.Text)

		if strings.Contains(text, term) {
			results = append(results, thread)
			continue
		}

		text = strings.ToLower(thread.Title)

		if strings.Contains(text, term) {
			results = append(results, thread)
			continue
		}
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Created > results[j].Created
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	return results
}

// SearchUsers searches all users.
func SearchUsers(originalTerm string, maxLength int) []*User {
	term := RemoveSpecialCharacters(strings.ToLower(originalTerm))

	var results []*SearchResult

	for user := range StreamUsers() {
		if user.ID == originalTerm {
			return []*User{user}
		}

		text := strings.ToLower(user.Nick)

		// Similarity check
		similarity := AdvancedStringSimilarity(term, text)

		if similarity < MinimumStringSimilarity {
			continue
		}

		results = append(results, &SearchResult{
			obj:        user,
			similarity: similarity,
		})
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*User, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*User)
	}

	return final
}

// SearchAnime searches all anime.
func SearchAnime(originalTerm string, maxLength int) []*Anime {
	term := RemoveSpecialCharacters(strings.ToLower(originalTerm))

	var results []*SearchResult

	check := func(text string) float64 {
		return AdvancedStringSimilarity(term, RemoveSpecialCharacters(strings.ToLower(text)))
	}

	add := func(anime *Anime, similarity float64) {
		similarity += float64(anime.Popularity.Total()) * popularityDamping

		results = append(results, &SearchResult{
			obj:        anime,
			similarity: similarity,
		})
	}

	for anime := range StreamAnime() {
		if anime.ID == originalTerm {
			return []*Anime{anime}
		}

		// Canonical title
		similarity := check(anime.Title.Canonical)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}

		// English
		similarity = check(anime.Title.English)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}

		// Synonyms
		for _, synonym := range anime.Title.Synonyms {
			similarity := check(synonym)

			if similarity >= MinimumStringSimilarity {
				add(anime, similarity)
				goto nextAnime
			}
		}

		// Japanese
		similarity = check(anime.Title.Japanese)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}

	nextAnime:
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*Anime, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*Anime)
	}

	return final
}

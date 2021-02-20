package marvel

import (
	"time"
)

// Servicer is the interface for the service layer containing functionality
// for fetching and storing Marvel characters.
type Servicer interface {
	GetAllCharacterIds() ([]int, error)
	GetCharacter(id int) (*Character, error)
	ReloadCache() error
}

// Service is the concrete implementation of Servicer. It uses a cache to
// reduce calls to Marvel's API.
type Service struct {
	client MarvelDataFetcher
	cache  CharacterCache
}

func NewService(client MarvelDataFetcher, cache CharacterCache) *Service {
	return &Service{
		client: client,
		cache:  cache,
	}
}

// GetAllCharacterIds returns the character IDs of all Marvel characters.
func (s *Service) GetAllCharacterIds() ([]int, error) {
	cachedCharIds, cachedLatestModified := s.cache.GetCharacterIds()

	var latestModified *time.Time
	if !cachedLatestModified.IsZero() {
		latestModified = &cachedLatestModified
	}

	characters, err := s.client.GetAllCharacters(latestModified)
	if err != nil {
		return nil, err
	}

	if len(characters) > 0 {
		// s.client.GetAllCharacters() already returns the latest modified character as the first element
		// if we really want to be safe, we can add a simple logic here to get the most recent `Modified`
		newLatestModified, err := time.Parse(dateFormatMarvelApi, characters[0].Modified)
		if err != nil {
			return nil, err
		}

		for _, c := range characters {
			cachedCharIds.Add(c.Id)
		}

		s.cache.SetCharacterIds(cachedCharIds, newLatestModified)
	}

	return cachedCharIds.ToSlice(), nil
}

// GetCharacter returns information about a specific character, given the character's ID.
func (s *Service) GetCharacter(id int) (*Character, error) {
	charData, err := s.client.GetCharacter(id)
	if err != nil {
		return nil, err
	}

	return &Character{
		Id:          charData.Id,
		Name:        charData.Name,
		Description: charData.Description,
	}, nil
}

// ReloadCache fetches all character IDs from Marvel's API and stores
// them in a cache, along with the latest modified time.
func (s *Service) ReloadCache() error {
	characters, err := s.client.GetAllCharacters(nil)
	if err != nil {
		return err
	}

	if len(characters) == 0 {
		return nil
	}

	// s.client.GetAllCharacters() already returns the latest modified character as the first element
	// if we really want to be safe, we can add a simple logic here to get the most recent `Modified`
	latestModified, err := time.Parse(dateFormatMarvelApi, characters[0].Modified)
	if err != nil {
		return err
	}

	charIds := NewIntSet()
	for _, c := range characters {
		charIds.Add(c.Id)
	}

	s.cache.SetCharacterIds(*charIds, latestModified)
	return nil
}

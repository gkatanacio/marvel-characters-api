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

// Service is the concrete implementation of Servicer. It keeps track of
// the timestamp of the latest modified character, along with a cache to
// reduce calls to Marvel's API.
type Service struct {
	client         MarvelDataFetcher
	cache          Cache
	latestModified *time.Time
}

func NewService(client MarvelDataFetcher, cache Cache) *Service {
	return &Service{
		client: client,
		cache:  cache,
	}
}

// GetAllCharacterIds returns the character IDs of all Marvel characters.
func (s *Service) GetAllCharacterIds() ([]int, error) {
	characters, err := s.client.GetAllCharacters(s.latestModified)
	if err != nil {
		return nil, err
	}

	cachedCharIds := s.cache.GetCharacterIds()

	if len(characters) > 0 {
		// s.client.GetAllCharacters() already returns the latest modified character as the first element
		// if we really want to be safe, we can add a simple logic here to get the most recent `Modified`
		latestModified, err := time.Parse(dateFormatMarvelApi, characters[0].Modified)
		if err != nil {
			return nil, err
		}
		s.latestModified = &latestModified

		for _, c := range characters {
			cachedCharIds.Add(c.Id)
		}

		s.cache.SetCharacterIds(cachedCharIds)
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
// them in a cache.
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
	s.latestModified = &latestModified

	charIds := NewIntSet()
	for _, c := range characters {
		charIds.Add(c.Id)
	}

	s.cache.SetCharacterIds(*charIds)
	return nil
}

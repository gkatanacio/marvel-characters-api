package marvel_test

import (
	"testing"
	"time"

	"github.com/gkatanacio/marvel-characters-api/internal/errs"
	"github.com/gkatanacio/marvel-characters-api/internal/marvel"
	"github.com/gkatanacio/marvel-characters-api/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Service_GetAllCharacterIds_NewCharsFetched(t *testing.T) {
	// given
	var latestModified *time.Time

	clientMock := new(mocks.MarvelDataFetcher)
	clientMock.On("GetAllCharacters", latestModified).Return([]*marvel.MarvelApiCharacterData{
		{
			Id:       1009282,
			Name:     "Doctor Strange",
			Modified: "2020-07-21T10:33:36-0400",
		},
		{
			Id:       1009187,
			Name:     "Black Panther",
			Modified: "2018-06-19T16:39:46-0400",
		},
	}, nil)

	service := marvel.NewService(clientMock, marvel.NewInMemCache())

	// when
	charIds, err := service.GetAllCharacterIds()

	// then
	assert.NoError(t, err)
	assert.NotNil(t, charIds)
	assert.ElementsMatch(t, []int{1009282, 1009187}, charIds)
	clientMock.AssertExpectations(t)
}

func Test_Service_GetAllCharacterIds_NoCharsFetched(t *testing.T) {
	// given
	latestModified := time.Now()

	clientMock := new(mocks.MarvelDataFetcher)
	clientMock.On("GetAllCharacters", &latestModified).Return(nil, nil)

	cachedIds := marvel.NewIntSet()
	cachedIds.Add(1009351)
	cachedIds.Add(1011490)
	cachedIds.Add(1011001)
	cache := marvel.NewInMemCache()
	cache.SetCharacterIds(*cachedIds, latestModified)
	service := marvel.NewService(clientMock, cache)

	// when
	charIds, err := service.GetAllCharacterIds()

	// then
	assert.NoError(t, err)
	assert.NotNil(t, charIds)
	assert.ElementsMatch(t, []int{1009351, 1011490, 1011001}, charIds)
	clientMock.AssertExpectations(t)
}

func Test_Service_GetCharacter_Found(t *testing.T) {
	// given
	charId := 1009610

	clientMock := new(mocks.MarvelDataFetcher)
	clientMock.On("GetCharacter", charId).Return(&marvel.MarvelApiCharacterData{
		Id:          charId,
		Name:        "Spider-Man",
		Description: "Bitten by a radioactive spider, high school student Peter Parker gained the speed, strength and powers of a spider.",
		Modified:    "2020-07-21T10:30:10-0400",
	}, nil)

	service := marvel.NewService(clientMock, marvel.NewInMemCache())

	// when
	character, err := service.GetCharacter(charId)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, character)
	assert.Equal(t, marvel.Character{
		Id:          character.Id,
		Name:        character.Name,
		Description: character.Description,
	}, *character)
	clientMock.AssertExpectations(t)
}

func Test_Service_GetCharacter_NotFound(t *testing.T) {
	// given
	charId := 9111111

	clientMock := new(mocks.MarvelDataFetcher)
	clientMock.On("GetCharacter", charId).Return(nil, errs.NewNotFound("no results"))

	service := marvel.NewService(clientMock, marvel.NewInMemCache())

	// when
	_, err := service.GetCharacter(charId)

	// then
	assert.Error(t, err)
	assert.IsType(t, new(errs.NotFound), err)
	clientMock.AssertExpectations(t)
}

func Test_Service_ReloadCache_HappyPath(t *testing.T) {
	// given
	var nilTime *time.Time

	clientMock := new(mocks.MarvelDataFetcher)
	clientMock.On("GetAllCharacters", nilTime).Return([]*marvel.MarvelApiCharacterData{
		{
			Id:       1009282,
			Name:     "Doctor Strange",
			Modified: "2020-07-21T10:33:36-0400",
		},
	}, nil)

	service := marvel.NewService(clientMock, marvel.NewInMemCache())

	// when
	err := service.ReloadCache()

	// then
	assert.NoError(t, err)
	clientMock.AssertExpectations(t)
}

package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gkatanacio/marvel-characters-api/internal/handlers"
	"github.com/gkatanacio/marvel-characters-api/internal/marvel"
	"github.com/gkatanacio/marvel-characters-api/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllCharactersHandler_Handle_HappyPath(t *testing.T) {
	// given
	marvelServiceMock := new(mocks.Servicer)
	marvelServiceMock.On("GetAllCharacterIds").Return([]int{1009351, 1011490, 1011001, 1009595}, nil)

	handler := handlers.NewGetAllCharactersHandler(marvelServiceMock)

	rr := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/characters", nil)

	// when
	handler.Handle(rr, req)

	// then
	assert.Equal(t, http.StatusOK, rr.Code)
	marvelServiceMock.AssertExpectations(t)
}

func Test_GetCharacterInfoHandler_Handle_HappyPath(t *testing.T) {
	// given
	charId := 1009351

	marvelServiceMock := new(mocks.Servicer)
	marvelServiceMock.On("GetCharacter", charId).Return(&marvel.Character{
		Id:          charId,
		Name:        "Hulk",
		Description: "An all too often misunderstood hero, the angrier the Hulk gets, the stronger the Hulk gets.",
	}, nil)

	handler := handlers.NewGetCharacterInfoHandler(marvelServiceMock)

	rr := httptest.NewRecorder()

	strCharId := strconv.Itoa(charId)
	req, _ := http.NewRequest(http.MethodGet, "/characters/"+strCharId, nil)
	req = mux.SetURLVars(req, map[string]string{"id": strCharId})

	// when
	handler.Handle(rr, req)

	// then
	assert.Equal(t, http.StatusOK, rr.Code)
	marvelServiceMock.AssertExpectations(t)
}

func Test_GetCharacterInfoHandler_Handle_InvalidId(t *testing.T) {
	// given
	charId := "123abc"

	handler := handlers.NewGetCharacterInfoHandler(nil)

	rr := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/characters/"+charId, nil)
	req = mux.SetURLVars(req, map[string]string{"id": charId})

	// when
	handler.Handle(rr, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

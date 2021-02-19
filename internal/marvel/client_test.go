package marvel_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gkatanacio/marvel-characters-api/internal/errs"
	"github.com/gkatanacio/marvel-characters-api/internal/marvel"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func testServer(path string, status int, testResponse string) *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintln(w, testResponse)
	}))
	return httptest.NewServer(r)
}

func testCfg(baseUrl string) *marvel.Config {
	return &marvel.Config{
		ApiBaseUrl:    baseUrl,
		ApiKeyPublic:  "1234",
		ApiKeyPrivate: "5678",
	}
}

func Test_Client_GetAllCharacters_HappyPath(t *testing.T) {
	// given
	testResponse := `
	{
		"code": 200,
		"status": "Ok",
		"copyright": "© 2021 MARVEL",
		"data": {
			"offset": 0,
			"limit": 100,
			"total": 1,
			"count": 1,
			"results": [
				{
					"id": 1011334,
					"name": "3-D Man",
					"description": "",
					"modified": "2014-04-29T14:18:17-0400"
				}
			]
		}
	}
	`
	ts := testServer("/v1/public/characters", http.StatusOK, testResponse)
	defer ts.Close()

	cfg := testCfg(ts.URL)

	client := marvel.NewClient(cfg)

	// when
	chars, err := client.GetAllCharacters(nil)

	// then
	assert.NoError(t, err)
	assert.Len(t, chars, 1)
	assert.Equal(t, marvel.MarvelApiCharacterData{
		Id:          1011334,
		Name:        "3-D Man",
		Description: "",
		Modified:    "2014-04-29T14:18:17-0400",
	}, *chars[0])
}

func Test_Client_GetCharacter_HappyPath(t *testing.T) {
	// given
	testResponse := `
	{
		"code": 200,
		"status": "Ok",
		"copyright": "© 2021 MARVEL",
		"data": {
			"offset": 0,
			"limit": 1,
			"total": 1,
			"count": 1,
			"results": [
				{
					"id": 1009351,
					"name": "Hulk",
					"description": "An all too often misunderstood hero, the angrier the Hulk gets, the stronger the Hulk gets.",
					"modified": "2020-07-21T10:35:15-0400"
				}
			]
		}
	}
	`
	ts := testServer("/v1/public/characters/{characterId}", http.StatusOK, testResponse)
	defer ts.Close()

	cfg := testCfg(ts.URL)

	client := marvel.NewClient(cfg)

	// when
	charData, err := client.GetCharacter(1009351)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, charData)
	assert.Equal(t, marvel.MarvelApiCharacterData{
		Id:          1009351,
		Name:        "Hulk",
		Description: "An all too often misunderstood hero, the angrier the Hulk gets, the stronger the Hulk gets.",
		Modified:    "2020-07-21T10:35:15-0400",
	}, *charData)
}

func Test_Client_GetCharacter_404(t *testing.T) {
	// given
	testResponse := `
	{
		"code": 404,
		"status": "We couldn't find that character"
	}
	`
	ts := testServer("/v1/public/characters/{characterId}", http.StatusNotFound, testResponse)
	defer ts.Close()

	cfg := testCfg(ts.URL)

	client := marvel.NewClient(cfg)

	// when
	_, err := client.GetCharacter(9111111)

	// then
	assert.Error(t, err)
	assert.IsType(t, new(errs.NotFound), err)
}

package marvel

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gkatanacio/marvel-characters-api/internal/errs"
	"golang.org/x/sync/errgroup"
)

// MarvelDataFetcher is the interface to abstract the actual
// calls to Marvel's API.
type MarvelDataFetcher interface {
	GetAllCharacters(modifiedSince *time.Time) ([]*MarvelApiCharacterData, error)
	GetCharacter(id int) (*MarvelApiCharacterData, error)
}

// Client is the concrete implementation of MarvelDataFetcher.
type Client struct {
	cfg        *Config
	httpClient *http.Client
}

func NewClient(cfg *Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}

// GetAllCharacters fetches all characters modified since the optionally
// provided `modifiedSince` timestamp. If `modifiedSince` is nil, GetAllCharacters fetches all
// the Marvel characters. The MarvelApiCharacterData with the most recent
// MarvelApiCharacterData.Modified is set as the first element in the returned slice.
func (c *Client) GetAllCharacters(modifiedSince *time.Time) ([]*MarvelApiCharacterData, error) {
	var characters []*MarvelApiCharacterData

	batchSize := 100

	qp := map[string]string{
		"limit":   strconv.Itoa(batchSize),
		"orderBy": "-modified",
	}
	if modifiedSince != nil {
		qp["modifiedSince"] = modifiedSince.Format(dateFormatMarvelApi)
	}

	marvelApiResp, err := c.httpGet("/v1/public/characters", qp)
	if err != nil {
		return nil, err
	}

	for _, r := range marvelApiResp.Data.Results {
		char, err := resultToMarvelApiCharacterData(r)
		if err != nil {
			return nil, err
		}
		characters = append(characters, char)
	}

	// fetch remaining characters (asynchronously) if needed
	if marvelApiResp.Data.Total > batchSize {
		eg := new(errgroup.Group)
		remainingChars := make(chan *MarvelApiCharacterData)

		for i := batchSize; i < marvelApiResp.Data.Total; i += batchSize {
			offset := i
			eg.Go(func() error {
				qpCopy := make(map[string]string)
				for k, v := range qp {
					qpCopy[k] = v
				}
				qpCopy["offset"] = strconv.Itoa(offset)

				remainingResp, err := c.httpGet("/v1/public/characters", qpCopy)
				if err != nil {
					return err
				}

				for _, r := range remainingResp.Data.Results {
					char, err := resultToMarvelApiCharacterData(r)
					if err != nil {
						return err
					}
					remainingChars <- char
				}

				return nil
			})
		}

		go func() {
			eg.Wait()
			close(remainingChars)
		}()

		for rc := range remainingChars {
			characters = append(characters, rc)
		}

		if err := eg.Wait(); err != nil {
			return nil, err
		}
	}

	return characters, nil
}

// GetCharacter fetches the character's data, given a character ID.
func (c *Client) GetCharacter(id int) (*MarvelApiCharacterData, error) {
	qp := map[string]string{
		"limit": "1",
	}

	marvelApiResp, err := c.httpGet(fmt.Sprintf("/v1/public/characters/%d", id), qp)
	if err != nil {
		return nil, err
	}

	// safe to assume Results[0] has value because no Results will throw 404 error in httpGet()
	return resultToMarvelApiCharacterData(marvelApiResp.Data.Results[0])
}

func (c *Client) httpGet(path string, additionalQueryParams map[string]string) (*MarvelApiResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.cfg.ApiBaseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	qp := c.authParams()
	for k, v := range additionalQueryParams {
		qp[k] = v
	}

	addQueryParams(req, qp)

	log.Println(req.URL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := new(MarvelApiErrResponse)
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, err
		}
		log.Printf(`error from marvel api: [%d] {status: "%s", message: "%s"}`, resp.StatusCode, errResp.Status, errResp.Message)

		if resp.StatusCode == http.StatusNotFound {
			return nil, errs.NewNotFound("no results")
		}

		return nil, errs.NewBadGateway("error response from marvel api")
	}

	marvelApiResp := new(MarvelApiResponse)
	if err := json.Unmarshal(respBody, &marvelApiResp); err != nil {
		return nil, err
	}

	return marvelApiResp, nil
}

func (c *Client) authParams() map[string]string {
	ts := getTs()
	hash := computeHash(ts, c.cfg.ApiKeyPrivate, c.cfg.ApiKeyPublic)

	return map[string]string{
		"apikey": c.cfg.ApiKeyPublic,
		"ts":     ts,
		"hash":   hash,
	}
}

func addQueryParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()
}

func getTs() string {
	now := time.Now()
	return strconv.Itoa(int(now.Unix()))
}

func computeHash(ts, privateKey, publicKey string) string {
	s := ts + privateKey + publicKey
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func resultToMarvelApiCharacterData(result interface{}) (*MarvelApiCharacterData, error) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	marvelApiCharacterData := new(MarvelApiCharacterData)
	if err := json.Unmarshal(jsonData, marvelApiCharacterData); err != nil {
		return nil, err
	}

	return marvelApiCharacterData, nil
}

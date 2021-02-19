package marvel

const (
	dateFormatMarvelApi = "2006-01-02T15:04:05-0700"
)

type MarvelApiResponse struct {
	Data struct {
		Offset  int           `json:"offset"`
		Limit   int           `json:"limit"`
		Total   int           `json:"total"`
		Count   int           `json:"count"`
		Results []interface{} `json:"results"`
	} `json:"data"`
}

type MarvelApiErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MarvelApiCharacterData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Modified    string `json:"modified"`
}

type Character struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

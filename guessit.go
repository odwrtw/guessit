package guessit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Errors
var (
	ErrInvalidRequest = errors.New("invalid guessit request")
	ErrServer         = errors.New("guessit server error")
)

// APIendpoint
const APIendpoint = "http://guessit.quimbo.fr/guess/"

// Types
const (
	Episode = "episode"
	Movie   = "movie"
)

// Response from the API
type Response struct {
	Episode  int    `json:"episodeNumber"`
	Quality  string `json:"screenSize"`
	Season   int    `json:"season"`
	ShowName string `json:"series"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Year     int    `json:"year"`
}

// Guess calls the guessit API to get the response
func Guess(filename string) (*Response, error) {
	// Guess it
	resp, err := http.Get(APIendpoint + url.QueryEscape(filename))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check http status
	switch resp.StatusCode {
	case http.StatusOK:
		// All good
	case http.StatusBadRequest:
		return nil, ErrInvalidRequest
	default:
		return nil, ErrServer
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	var response *Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

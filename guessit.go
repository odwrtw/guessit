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
	ErrInvalidRequest = errors.New("guessit: invalid request")
	ErrServer         = errors.New("guessit: server error")
)

// APIendpoint represents the default API endpoint
const APIendpoint = "http://guessit.quimbo.fr/guess/"

// Response from the API
type Response struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Episode      int    `json:"episode"`
	Season       int    `json:"season"`
	Year         int    `json:"year"`
	Quality      string `json:"screenSize"`
	ReleaseGroup string `json:"releaseGroup"`
	AudioCodec   string `json:"audio_codec"`
	VideoCodec   string `json:"video_codec"`
	Container    string `json:"container"`
	Format       string `json:"format"`
	MimeType     string `json:"mimetype"`
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

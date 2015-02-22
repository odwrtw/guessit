package guessit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// APIendpoint
const APIendpoint = "http://guessit.io/guess"

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
	// Generate URL
	u, err := url.Parse(APIendpoint)
	if err != nil {
		return nil, err
	}
	urlValues := &url.Values{}
	urlValues.Add("filename", filename)
	u.RawQuery = urlValues.Encode()
	fmt.Println(u)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

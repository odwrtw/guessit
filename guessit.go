package guessit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Errors
var (
	ErrInvalidRequest = errors.New("guessit: invalid request")
	ErrServer         = errors.New("guessit: server error")
)

// Request represents the request params
type Request struct {
	Name string `json:"name"`
}

// Response from the API
type Response struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Episode      int    `json:"episode"`
	Season       int    `json:"season"`
	Year         int    `json:"year"`
	Quality      string `json:"screen_size"`
	ReleaseGroup string `json:"release_group"`
	AudioCodec   string `json:"audio_codec"`
	VideoCodec   string `json:"video_codec"`
	Container    string `json:"container"`
	Format       string `json:"format"`
	MimeType     string `json:"mimetype"`
}

// Client represents a guessit client
type Client struct {
	endpoint string
}

// New returns a new client
func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// Guess calls the guessit API to get the response
func (c *Client) Guess(filename string) (*Response, error) {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	data, err := json.Marshal(Request{Name: filename})
	if err != nil {
		return nil, err
	}

	// Guess it
	buf := bytes.NewBuffer(data)
	resp, err := httpClient.Post(c.endpoint, "application/json", buf)
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

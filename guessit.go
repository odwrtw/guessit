package guessit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Errors
var (
	ErrInvalidRequest = errors.New("guessit: invalid request")
	ErrServerError    = errors.New("guessit: server error")
)

var defaultTimeout = 10 * time.Second

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
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return c.GuessWithContext(ctx, filename)
}

// GuessWithContext guess with a context
func (c *Client) GuessWithContext(ctx context.Context, filename string) (*Response, error) {
	r := struct {
		Name string `json:"name"`
	}{Name: filename}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// All good
	case http.StatusBadRequest:
		return nil, ErrInvalidRequest
	default:
		return nil, ErrServerError
	}

	response := &Response{}
	defer resp.Body.Close()
	return response, json.NewDecoder(resp.Body).Decode(response)
}

package guessit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var responseMap = map[string][]byte{
	"Mr.Robot.S04E01.401.Unauthorized.720p.AMZN.WEB-DL.DDP5.1.H.264-NTG[eztv].mkv": []byte(`
	{
	  "audio_channels": "5.1",
	  "audio_codec": "Dolby Digital Plus",
	  "container": "mkv",
	  "episode": 1,
	  "episode_title": "401 Unauthorized",
	  "release_group": "NTG[eztv]",
	  "screen_size": "720p",
	  "season": 4,
	  "source": "Web",
	  "streaming_service": "Amazon Prime",
	  "title": "Mr Robot",
	  "type": "episode",
	  "video_codec": "H.264"
	}`),
	"The.Matrix.Reloaded.2003.1080p.BrRip.x264.YIFY.mp4": []byte(`
	{
	  "container": "mp4",
	  "mimetype": "video/mp4",
	  "other": [
		"Reencoded",
		"Rip"
	  ],
	  "release_group": "YIFY",
	  "screen_size": "1080p",
	  "source": "Blu-ray",
	  "title": "The Matrix Reloaded",
	  "type": "movie",
	  "video_codec": "H.264",
	  "year": 2003
	}`),
}

type testServer struct {
	*testing.T
}

func (ts *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, ok := responseMap[params.Name]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "applicationjson")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func TestGuessit(t *testing.T) {
	server := httptest.NewServer(&testServer{T: t})
	defer server.Close()
	client := New(server.URL)

	tt := []struct {
		name          string
		expectedError error
		expected      *Response
	}{
		{
			name: "Mr.Robot.S04E01.401.Unauthorized.720p.AMZN.WEB-DL.DDP5.1.H.264-NTG[eztv].mkv",
			expected: &Response{
				Type:         "episode",
				Title:        "Mr Robot",
				Episode:      1,
				Season:       4,
				Quality:      "720p",
				ReleaseGroup: "NTG[eztv]",
				AudioCodec:   "Dolby Digital Plus",
				VideoCodec:   "H.264",
				Container:    "mkv",
			},
		},
		{
			name: "The.Matrix.Reloaded.2003.1080p.BrRip.x264.YIFY.mp4",
			expected: &Response{
				Type:         "movie",
				Title:        "The Matrix Reloaded",
				Year:         2003,
				Quality:      "1080p",
				ReleaseGroup: "YIFY",
				VideoCodec:   "H.264",
				Container:    "mp4",
				MimeType:     "video/mp4",
			},
		},
		{
			name:          "yolo",
			expectedError: ErrServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := client.Guess(tc.name)
			if err != tc.expectedError {
				t.Fatalf("expected error %q, got %q", tc.expectedError.Error(), err.Error())
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("invalid result")
			}
		})
	}
}

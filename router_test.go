package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWelcomePageHandler(t *testing.T) {
	assert := assert.New(t)

	mux := createMuxRouter()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/welcome", nil)
	mux.ServeHTTP(res, req)

	assert.Equal(res.Code, 200)
	assert.Contains(res.Body.String(), "mjpeg streamer demo app written in GOLANG")

}

func TestMJPEGStreamerHandler(t *testing.T) {
	assert := assert.New(t)
	streamSourceFiles = []string{
		"./video/test/video-0001.jpeg",
	}
	mux := createMuxRouter()
	res := httptest.NewRecorder()
	// ask only for n frames, so that stream would finish rather quickly
	n := 5
	url := fmt.Sprintf("/mjpeg?frames=%d", n)
	req, _ := http.NewRequest("GET", url, nil)
	mux.ServeHTTP(res, req)

	assert.Equal(res.Code, 200)
	//assert.Contains(res.Body.String(), "mjpeg streamer demo app written in GOLANG")
	parsed := parseHTTPStream(res.Body.String(), res.HeaderMap)
	assert.Equal(parsed.numberOfFrames, n)

}

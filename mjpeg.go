package main

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

type MJPEGStream struct {
	frames         [][]byte
	numberOfFrames int
}

func imageStreamGenerator(c chan string, controlC chan string, files []string) {
	i := 0
	for {
		select {
		case c <- files[i]:
			i++
			if i >= len(files) {
				i = 0
			}
		case cmd := <-controlC:
			log.Printf("WEHEEJ: quiting streamgenerator, received: %s", cmd)
			return

		}
	}
}

func parseHTTPStream(body string, headers http.Header) (parsedStream MJPEGStream) {
	contentType := headers.Get("Content-Type")
	boundaryRe := regexp.MustCompile(`^multipart/x-mixed-replace;boundary=(.*)$`)
	boundary := boundaryRe.FindStringSubmatch(contentType)[1]

	bodyReader := strings.NewReader(body)
	reader := multipart.NewReader(bodyReader, boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			log.Printf("Multipart reader finished. Error: %v", err)
			break
		}
		frame, err := ioutil.ReadAll(part)
		if err != nil {
			log.Fatal(err)
		}
		parsedStream.frames = append(parsedStream.frames, frame)
		parsedStream.numberOfFrames++
	}
	return parsedStream
}

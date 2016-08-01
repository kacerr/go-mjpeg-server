package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

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

func mjpeg(responseWriter http.ResponseWriter, request *http.Request) {
	var source = make(chan []byte)

	log.Printf("Start request %s", request.URL)

	mimeWriter := multipart.NewWriter(responseWriter)

	log.Printf("Boundary: %s", mimeWriter.Boundary())

	contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
	responseWriter.Header().Add("Content-Type", contentType)

	for {
		frameStartTime := time.Now()
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")

		partWriter, partErr := mimeWriter.CreatePart(partHeader)
		if nil != partErr {
			log.Printf(partErr.Error())
			break
		}

		snapshot := <-source
		if _, writeErr := partWriter.Write(snapshot); nil != writeErr {
			log.Printf(writeErr.Error())
		}
		frameEndTime := time.Now()

		frameDuration := frameEndTime.Sub(frameStartTime)
		fps := float64(time.Second) / float64(frameDuration)
		log.Printf("Frame time: %s (%.2f)", frameDuration, fps)
	}

	log.Printf("Success request")
}

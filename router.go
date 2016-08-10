package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"time"
)

func createMuxRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/debug/vars", http.DefaultServeMux)
	mux.HandleFunc("/mjpeg", handlerServeStream)
	mux.HandleFunc("/welcome", handleWelcome)

	return mux
}

func handleWelcome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome. \n This is mjpeg streamer demo app written in GOLANG"))
}

func handlerServeStream(w http.ResponseWriter, r *http.Request) {
	var image []byte
	var file string
	var framesToSend int
	var err error

	r.ParseForm()
	if frames, keyExist := r.Form["frames"]; keyExist {
		framesToSend, err = strconv.Atoi(frames[0])
		if err != nil {
			framesToSend = -1
		}
	}

	streamsInProgress.Add(1)
	defer streamsInProgress.Add(-1)

	imageStreamChannel := make(chan string)
	controlChannel := make(chan string)

	// quit imageStreamGenerator go routine when request is ended
	defer func() {
		controlChannel <- "stop"
	}()

	go imageStreamGenerator(imageStreamChannel, controlChannel, streamSourceFiles)

	mimeWriter := multipart.NewWriter(w)

	log.Printf("Boundary: %s", mimeWriter.Boundary())

	contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	for {
		//frameStartTime := time.Now()
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")
		partHeader.Add("X-Current-Time", time.Now().Format(time.RFC1123))

		partWriter, partErr := mimeWriter.CreatePart(partHeader)
		if nil != partErr {
			log.Printf(partErr.Error())
			break
		}

		file = <-imageStreamChannel
		image, _ = ioutil.ReadFile(file)
		if _, writeErr := partWriter.Write(image); nil != writeErr {
			log.Printf(writeErr.Error())
		}
		time.Sleep(100 * time.Millisecond)

		// if number of frames that are supposed to be send was set then decrease the counter and exit on 0
		if framesToSend > 0 {
			log.Printf("Still %d frames until finshed \n", framesToSend)
			framesToSend--
			if framesToSend == 0 {
				// we need to end multipart stream in a correct fashion
				endOfStream := fmt.Sprintf("\r\n--%s--\r\n", mimeWriter.Boundary())
				partWriter.Write([]byte(endOfStream))
				break
			}

		}

	}
}

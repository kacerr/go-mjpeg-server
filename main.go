package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {

	var files []string
	// router
	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Yes, this is status !"))
		//io.WriteString(w, "Yes, this is status !")
	})

	mux.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		// display current path
		_, filename, _, _ := runtime.Caller(1)
		w.Write([]byte(filename + "\n"))

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(dir + "\n"))

		// load list of files: video/*.jpeg
		files, _ = filepath.Glob(dir + "/video/*.jpeg")
		// print list of files
		for i, file := range files {
			w.Write([]byte(file + " "))
			if math.Mod(float64(i), 10) == 0 {
				w.Write([]byte("\n"))
			}
		}

		//w.Write([]byte("Yes, this is status !"))
		//io.WriteString(w, "Yes, this is status !")
	})

	mux.HandleFunc("/serve", func(w http.ResponseWriter, r *http.Request) {
		var file string
		var image []byte

		imageStreamChannel := make(chan string)
		controlChannel := make(chan string)

		// quit imageStreamGenerator go routine when request is ended
		defer func() {
			controlChannel <- "stop"
		}()

		/*
			for _, file := range files {
				log.Println(file)
			}
		*/

		// start go routine that generates stream of filenames that could be served
		go imageStreamGenerator(imageStreamChannel, controlChannel, files)

		mimeWriter := multipart.NewWriter(w)

		log.Printf("Boundary: %s", mimeWriter.Boundary())

		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		w.Header().Add("Content-Type", contentType)
		w.Header().Add("Access-Control-Allow-Origin", "*")

		for {
			//frameStartTime := time.Now()
			partHeader := make(textproto.MIMEHeader)
			partHeader.Add("Content-Type", "image/jpeg")

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

			/*
				frameEndTime := time.Now()
				frameDuration := frameEndTime.Sub(frameStartTime)
				fps := float64(time.Second) / float64(frameDuration)
				log.Printf("Frame time: %s (%.2f)", frameDuration, fps)
			*/
		}

		/*
			file := <-imageStreamChannel
			w.Write([]byte(file + "\n"))
		*/
	})

	log.Println("Server is up and listening !!!")
	err := http.ListenAndServe(":3003", mux)
	log.Fatalf("FATAL: %v", err)
}

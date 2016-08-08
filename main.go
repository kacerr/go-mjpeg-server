package main

import (
	"expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var streamSourceFiles []string
var streamsInProgress = expvar.NewInt("streamsInProgress")

func main() {

	sourceType := flag.String("source-type", "directory", "Set image stream source type (only directory is available at the moment)")
	sourcePath := flag.String("source-path", "./video", "Image stream source path")
	listenPort := flag.Int("port", 3003, "Server TCP port")
	flag.Parse()

	if *sourceType == "directory" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		streamSourceFiles, _ = filepath.Glob(dir + "/" + *sourcePath + "/*.jpeg")
		if len(streamSourceFiles) == 0 {
			fmt.Println("ERROR: \n")
			fmt.Printf("  No files to be streamed found in path: %s \n\n", dir+"/"+*sourcePath+"/*.jpeg")
			os.Exit(1)
		}
		fmt.Printf("%d files have been found in source path: %s \n\n", len(streamSourceFiles), dir+"/"+*sourcePath+"/*.jpeg")
	} else {
		fmt.Println("ERROR: \n")
		fmt.Println("  We are sorry, only supported source type at the moment is directory \n\n")
		os.Exit(1)
	}
	//var files []string
	// router
	mux := createMuxRouter()

	listenURL := fmt.Sprintf(":%d", *listenPort)
	log.Println("Server is going up and starts listening on address: " + listenURL)
	err := http.ListenAndServe(":3003", mux)
	//err := http.ListenAndServe(listenURL, mux)
	if err != nil {
		log.Fatalf("FATAL: %v", err)
	}
}

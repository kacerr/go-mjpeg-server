package main

import "log"

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

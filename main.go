package main

import (
	"./src/crawler"
	"os"
	"strconv"
)

func main() {
	domen := os.Args[1]

	var wCount int
	if len(os.Args) < 3 {
		wCount = 1
	} else {
		wCount, _ = strconv.Atoi(os.Args[2])
	}

	cr := crawler.NewCrawler(domen)
	cr.Run(wCount)
}

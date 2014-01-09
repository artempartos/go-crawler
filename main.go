package main

import (
	"./src/crawler"
	"os"
)

func main() {
	domen := os.Args[1]
	cr := crawler.NewCrawler(domen)
	cr.Run()
}

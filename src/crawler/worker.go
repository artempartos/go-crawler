package crawler

import (
	"fmt"
	"github.com/opesun/goquery"
	"net/http"
)

type Worker struct {
	in   linkChan
	out  responseChan
	free workerChan
}

func NewWorker(free workerChan, out responseChan) *Worker {
	in := make(linkChan)
	return &Worker{in, out, free}
}

func (w *Worker) Run() {
	go func() {
		for {
			w.free <- w
			link := <-w.in
			resp := w.process(link)
			w.out <- resp
		}
	}()
}

func (w *Worker) Process(link string) {
	w.in <- link
}

func (w *Worker) process(link string) CrawlerResponse {
	response, err := http.Get(link)

	if err == nil {
		fmt.Println("processing... ", link, response.Status)
		defer response.Body.Close()
	}

	switch {
	case err != nil || response.StatusCode == 404:
		return CrawlerResponse{success: false, current: link}
	default:
		return ProcessHttpResponse(response, link)
	}

}

func ProcessHttpResponse(resp *http.Response, link string) CrawlerResponse {
	x, err := goquery.Parse(resp.Body)
	if err == nil {
		links := x.Find("a").Attrs("href")
		return CrawlerResponse{success: true, current: link, links: links}
	} else {
		return CrawlerResponse{success: false, current: link}
	}
}

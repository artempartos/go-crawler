package crawler

import (
	"fmt"
	"github.com/opesun/goquery"
)

type Worker struct {
	in    linkChan
	out   responseChan
	free  workerChan
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

func (w *Worker) process(link string) Response {
	fmt.Println("processing... ", link)
	x, err := goquery.ParseUrl(link)

	if err == nil {
		links := x.Find("a").Attrs("href")
		return Response{success: true, current: link, links: links}
	} else {
		return Response{success: false, current: link}
	}
}

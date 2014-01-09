package crawler

import (
	"fmt"
	"github.com/opesun/goquery"
	"log"
)

type Worker struct {
	outToMaster     ResponseChan
	in              StringChan
	outToController WorkerChan
}

func NewWorker(channelToController WorkerChan, RespChannel ResponseChan) *Worker {
    StrIn := make(StringChan)
	return &Worker{outToMaster: RespChannel, outToController: channelToController, in: StrIn}
}

func (w *Worker) Run() {
	go func() {
		for {
		    log.Println("worker send self to controller")
		    w.outToController <- w

		    log.Println("worker wait link")
			link := <-w.in
			resp := w.process(link)

			log.Println("worker send response")
			w.outToMaster <- resp
		}
	}()
}

func (w *Worker) Process(link string) {
	w.in <- link
}

func (w *Worker) process(link string) Response {
	fmt.Println("processing... ", "/" + link)

	x, err := goquery.ParseUrl(Domen + "/" + link)

	if err == nil {
		links := x.Find("a").Attrs("href")
		return Response{success: true, current: link, links: links}
	} else {
		return Response{success: false, current: link}
	}
}

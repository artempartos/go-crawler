package crawler

import (
	"fmt"
	"github.com/opesun/goquery"
)

type Worker struct {
	outToMaster     ResponseChan
	in              StringChan
	workerChan WorkerChan
}

func NewWorker(RespChannel ResponseChan, workerChan WorkerChan) *Worker {
    StrIn := make(StringChan)
	return &Worker{outToMaster: RespChannel, workerChan: workerChan, in: StrIn}
}

func (w *Worker) Run() {
	go func() {
		for {
		    w.workerChan <- w

			link := <-w.in
			resp := w.process(link)

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

package crawler

import (
	"fmt"
	"regexp"
)

var Domen string
var WorkersCount int = 20

type Crawler struct {
	domen       string
	links       []string
	in_channel  ResponseChan
	result      map[string]string
	workers     WorkerChan
}

type Response struct {
	success bool
	links   []string
	current string
}

type ResponseChan chan Response
type StringChan chan string
type WorkerChan chan *Worker

func NewCrawler(domen string) *Crawler {
	in := make(ResponseChan)
	workers := make(WorkerChan, WorkersCount )
	result := make(map[string]string)
	Domen = domen
	return &Crawler{domen: domen, in_channel: in, result: result, workers: workers}
}

func (c *Crawler) Run() {

	for i := 0; i < WorkersCount; i++ {
        w := NewWorker(c.in_channel, c.workers)
        w.Run()
    }

    root := ""
    c.AddToQueue(root)

    for {
        hasWork := len(c.links) > 0
        select {
        case response := <-c.in_channel:
            c.ResponseProcess(response)
        case worker := <- c.workers:
            if hasWork {
                link := c.links[0]
                c.links = c.links[1:]
                worker.Process(link)
            } else {
                c.workers <- worker
            }
        }
    }
}

func (c *Crawler) ResponseProcess(response Response) {
	if response.success {

		c.result[response.current] = "ok"
		for _, link := range response.links {
			c.LinkProcess(link)
		}

	} else {
		c.result[response.current] = "fail"
	}
	PrintResponse(c.result)
}

func (c *Crawler) LinkProcess(link string) {
	isAbsolute, _ := regexp.MatchString("http*", link)

	if isAbsolute {
		sameDomen, _ := regexp.MatchString("*"+Domen+"*", link)
		if sameDomen {
			reg := regexp.MustCompile(Domen+"/")
			relative := reg.ReplaceAllString(link, "")
			c.AddToQueue(relative)
		} else {
			c.result[link] = "anotherDomen"
		}

	} else {
		isAnchor, _ := regexp.MatchString("#", link)
		if isAnchor {
			c.result[link] = "Anchor"
		} else {
			c.AddToQueue(link)
		}
	}
}

func (c *Crawler) AddToQueue(link string) {
	_, ok := c.result[link]
	if !ok {
		c.result[link] = "inQueue"
        c.links = append(c.links, link)
	}
}

func PrintResponse(result map[string]string) {
	var ok, failed, domen, queue, anchor int
	for _, v := range result {
		if v == "ok" {
			ok++
		}
		if v == "fail" {
			failed++
		}
		if v == "anotherDomen" {
			domen++
		}
		if v == "Anchor" {
			anchor++
		}
		if v == "inQueue" {
			queue++
		}

	}
	fmt.Printf("\tok: %v, queue: %v, domen: %v, anchor: %v, failed: %v\n", ok, queue, domen, anchor, failed)
}

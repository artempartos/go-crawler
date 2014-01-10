package crawler

import (
	"fmt"
	"regexp"
)

type responseChan chan CrawlerResponse
type linkChan chan string
type workerChan chan *Worker
type resultMap map[string]string

type Crawler struct {
	domen   string
	in      responseChan
	out     linkChan
	workers workerChan
	result  resultMap
	links   []string
}

type CrawlerResponse struct {
	success bool
	links   []string
	current string
}

func NewCrawler(domen string) *Crawler {
	in := make(responseChan)
	out := make(linkChan)
	workers := make(workerChan)
	result := make(resultMap)
	return &Crawler{domen, in, out, workers, result, []string{}}
}

func (c *Crawler) Run(workersCount int) {
	c.runWorkers(workersCount)
	c.PushLink(c.domen)

	for {
		if c.HasLink() {
			select {
			case response := <-c.in:
				c.ResponseProcess(response)
			case worker := <-c.workers:
				worker.Process(c.PopLink())
			}
		} else {
			response := <-c.in
			c.ResponseProcess(response)
		}
	}
}

func (c *Crawler) runWorkers(count int) {
	for i := 0; i < count; i++ {
		w := NewWorker(c.workers, c.in)
		w.Run()
	}
}

func (c *Crawler) PopLink() string {
	var link string
	link, c.links = c.links[len(c.links)-1], c.links[:len(c.links)-1]
	return link
}

func (c *Crawler) HasLink() bool {
	return len(c.links) > 0
}

func (c *Crawler) ResponseProcess(response CrawlerResponse) {
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
		sameDomen, _ := regexp.MatchString(c.domen+"*", link)
		if sameDomen {
			c.PushLink(link)
		} else {
			c.result[link] = "anotherDomen"
		}
	} else {
		isAnchor, _ := regexp.MatchString("#", link)
		if isAnchor {
			c.result[link] = "Anchor"
		} else {
			absolute := c.domen + "/" + link
			c.PushLink(absolute)
		}
	}
}

func (c *Crawler) PushLink(link string) {
	_, ok := c.result[link]
	if !ok {
		c.result[link] = "inQueue"
		c.links = append(c.links, link)
	}
}

func PrintResponse(result map[string]string) {
	var ok, failed, domen, queue, anchor int
	for _, v := range result {
		switch v {
		case "ok":
			ok++
		case "fail":
			failed++
		case "anotherDomen":
			domen++
		case "Anchor":
			anchor++
		case "inQueue":
			queue++
		}
	}
	fmt.Printf("\tok: %v, queue: %v, domen: %v, anchor: %v, failed: %v\n", ok, queue, domen, anchor, failed)
}

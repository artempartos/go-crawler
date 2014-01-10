package crawler

import (
	"fmt"
	"net/url"
	"regexp"
)

type responseChan chan CrawlerResponse
type linkChan chan string
type workerChan chan *Worker
type resultMap map[string]string
type Crawler struct {
	domen   *url.URL
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
	domen_link, _ := url.Parse(domen)
	return &Crawler{domen_link, in, out, workers, result, []string{}}
}

func (c *Crawler) Run(workersCount int) {
	c.runWorkers(workersCount)
	c.PushLink(c.domen.String())
	c.RunLoop()
}

func (c *Crawler) RunLoop() {
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

func (c *Crawler) PushLink(link string) {
	_, ok := c.result[link]
	if !ok {
		c.result[link] = "inQueue"
		c.links = append(c.links, link)
	}
}

func (c *Crawler) HasLink() bool {
	return len(c.links) > 0
}

func (c *Crawler) ResponseProcess(response CrawlerResponse) {
	if response.success {
		c.result[response.current] = "ok"
		for _, link_string := range response.links {
			link, _ := url.Parse(link_string)
			c.LinkProcess(link)
		}
	} else {
		c.result[response.current] = "fail"
	}
	PrintResponse(c.result)
}

func (c *Crawler) LinkProcess(url *url.URL) {

	if url.IsAbs() {
		if url.Host == c.domen.Host {
			c.PushLink(url.String())
		} else {
			c.result[url.String()] = "anotherDomen"
		}
	} else {
		isAnchor, _ := regexp.MatchString("#", url.String())
		if isAnchor {
			c.result[url.String()] = "Anchor"
		} else {
			absolute := c.domen.String() + "/" + url.String()
			c.PushLink(absolute)
		}
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

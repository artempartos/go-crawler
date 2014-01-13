package crawler

import (
	"github.com/wsxiaoys/terminal/color"
	"net/url"
)

type responseChan chan CrawlerResponse
type linkChan chan string
type workerChan chan *Worker
type resultMap map[string]bool

type Crawler struct {
	domen   *url.URL
	in      responseChan
	out     linkChan
	workers workerChan
	result  resultMap
	links   []string
	count   int
}

type CrawlerResponse struct {
	success bool
	Status  string
	links   []string
	current string
}

func NewCrawler(domen string) *Crawler {
	in := make(responseChan)
	out := make(linkChan)
	workers := make(workerChan)
	result := make(resultMap)
	domen_link, _ := url.Parse(domen)
	return &Crawler{domen_link, in, out, workers, result, []string{}, 0}
}

func (c *Crawler) Run(workersCount int) {
	c.runWorkers(workersCount)
	c.PushLink(c.domen.String())
	c.RunLoop()
}

func (c *Crawler) RunLoop() {
	for {
		var wCh workerChan
		if c.HasLink() {
			wCh = c.workers
		}

		select {
		case response := <-c.in:
			c.count++
			c.ResponseProcess(response)
			var col string
			if response.success {
				col = "@g"
			} else {
				col = "@r"
			}
			color.Println(col, "#", c.count, response.current, response.Status)
		case worker := <-wCh:
			l := c.PopLink()
			c.result[l] = true
			worker.Process(l)
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
		c.links = append(c.links, link)
	}
}

func (c *Crawler) HasLink() bool {
	return len(c.links) > 0
}

func (c *Crawler) ResponseProcess(response CrawlerResponse) {
	if response.success {
		for _, link_string := range response.links {
			c.LinkProcess(link_string)
		}
	}
}

func (c *Crawler) LinkProcess(link_string string) {
	link, err := NewLink(link_string)

	if err == nil {
		switch {
		case link.isSameHost(c.domen):
			c.PushLink(link.Unify())
		case link.isRelative():
			c.PushLink(link.withHost(c.domen))
		}
	}
}

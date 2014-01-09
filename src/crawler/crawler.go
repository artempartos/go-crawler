package crawler

import (
	"fmt"
	"regexp"
)

var Domen string

type Crawler struct {
	domen       string
	in_channel  ResponseChan
	out_channel StringChan
	result      map[string]string
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
	out := make(StringChan, 50)
	result := make(map[string]string)
	Domen = domen
	return &Crawler{domen: domen, in_channel: in, out_channel: out, result: result}
}

func (c *Crawler) Run() {
	controller := NewController(c.in_channel, c.out_channel, 10)
	controller.Run()
	root := ""
	c.AddToQueue(root)
	//TODO fix on select for read/write
	for {
		response := <-c.in_channel
		c.ResponseProcess(response)
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
		c.out_channel <- link
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
	if ok == 84 {
		fmt.Println(result)
	}
}

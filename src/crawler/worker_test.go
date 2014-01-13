package crawler

import (
	. "launchpad.net/gocheck"
)

func (s *MySuite) TestWorkerChan(c *C) {
	workerChan := make(workerChan)
	worker := NewWorker(workerChan, nil)

	worker.Run()
	w := <-workerChan
	c.Assert(w, Equals, worker)
}

func (s *MySuite) TestWorkerProcessing(c *C) {
	respChan := make(responseChan)
	workerChan := make(workerChan)
	worker := NewWorker(workerChan, respChan)

	worker.Run()
	w := <-workerChan
	w.Process("http://nox73.ru/")

	response := <-respChan
	c.Assert(response.success, Equals, true)
}

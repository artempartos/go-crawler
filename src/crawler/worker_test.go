package crawler

import (
	. "launchpad.net/gocheck"
)

func (s *MySuite) TestWorker(c *C) {
	chanMaster := make(ResponseChan)
	chanController := make(WorkerChan)
	worker := NewWorker(chanController, chanMaster)
	worker.Run()
	w := <-chanController
	c.Assert(w, Equals, worker)
}

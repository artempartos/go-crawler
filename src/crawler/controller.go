package crawler

type Controller struct {
	channelForResponse ResponseChan
	channelForWorkers  WorkerChan
	inChannel          StringChan
	workersCount       int
}

func NewController(chanMaster ResponseChan, linkChan StringChan, workCount int) *Controller {
	ForWorkers := make(WorkerChan, workCount)
	return &Controller{channelForResponse: chanMaster, inChannel: linkChan, workersCount: workCount, channelForWorkers: ForWorkers}
}

func (cont *Controller) Run() {
	for i := 0; i < cont.workersCount; i++ {
		w := NewWorker(cont.channelForWorkers, cont.channelForResponse)
		w.Run()
	}

	go func() {
		for {
			worker := <-cont.channelForWorkers
			link := <-cont.inChannel
			worker.Process(link)
		}
	}()
}

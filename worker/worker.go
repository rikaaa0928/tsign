package worker

type workerData struct {
	act *account
	d   *data
	f   func(showData ShowData)
}

func NewWorker() *Worker {
	return &Worker{c: make(chan *workerData)}
}

type Worker struct {
	c chan *workerData
}

func (w *workerData) Do() {
	s := Sign(w.act, w.d)
	if w.f != nil {
		w.f(s)
	}
}

func (w *Worker) SyncGo() {
	for {
		d := <-w.c
		go d.Do()
	}
}

func (w *Worker) AsyncGo() {
	go w.SyncGo()
}

func (w *Worker) AsyncNotify(d *workerData) {
	go w.SyncNotify(d)
}

func (w *Worker) SyncNotify(d *workerData) {
	w.c <- d
}

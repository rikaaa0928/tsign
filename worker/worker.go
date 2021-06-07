package worker

type workerData struct {
	act *account
	d   *data
	f   func(showData ShowData)
}

type worker struct {
	c chan *workerData
}

func (w *workerData) Do() {
	s := Sign(w.act, w.d)
	if w.f != nil {
		w.f(s)
	}
}

func (w *worker) SyncGo() {
	for {
		d := <-w.c
		d.Do()
	}
}

func (w *worker) AsyncGo() {
	go w.SyncGo()
}

func (w *worker) AsyncNotify(d *workerData) {
	go w.SyncNotify(d)
}

func (w *worker) SyncNotify(d *workerData) {
	w.c <- d
}

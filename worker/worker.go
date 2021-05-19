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
	w.f(s)
}

func (w *worker) Go() {
	for {
		d := <-w.c
		d.Do()
	}
}

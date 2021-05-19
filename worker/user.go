package worker

import "sync/atomic"

type user struct {
	a       *account
	dataMap atomic.Value
	showMap atomic.Value
	c       chan *workerData
}

func NewWorker(a *account) *user {
	w := new(user)
	w.a = a
	dm := make(map[string]*data)
	l, err := a.GetList()
	if err == nil {
		for _, v := range l {
			dm[v.name] = v
		}
	}
	w.dataMap.Store(dm)
	sm := make(map[string]ShowData, len(dm))
	w.showMap.Store(sm)
	return w
}

package worker

import (
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := &Worker{c: make(chan *workerData)}
	w.AsyncGo()
	u := NewUserManager("../test")
	wg := sync.WaitGroup{}
	for _, v := range u.UserMap() {
		wg.Add(1)
		dm := v.DataMap()
		a := v.a
		v := v
		go func() {
			for _, d := range dm {
				w.AsyncNotify(&workerData{d: d, act: a, f: func(showData ShowData) {
					v.UpdateShowData(showData.Name, showData)
					t.Log(showData)
				}})
				time.Sleep(time.Second)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(time.Second)
}

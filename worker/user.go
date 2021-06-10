package worker

import (
	"io/ioutil"
	"log"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type user struct {
	a       *account
	dataMap atomic.Value
	showMap *sync.Map
	c       chan *workerData
}

func NewUser(a *account) *user {
	u := &user{a: a, showMap: &sync.Map{}}
	dm := make(map[string]*data)
	l, err := a.GetList()
	if err == nil {
		for _, v := range l {
			dm[v.unicodeName] = v
			u.showMap.Store(ToUtf8(v.name), ShowData{Name: ToUtf8(v.name), Exp: v.exp, Stat: "waiting"})
		}
	}
	u.dataMap.Store(dm)
	return u
}

func (u *user) RefreshData() {
	dm := u.dataMap.Load().(map[string]*data)
	l, err := u.a.GetList()
	if err == nil {
		for _, v := range l {
			if _, ok := dm[v.unicodeName]; !ok {
				dm[v.unicodeName] = v
				u.showMap.Store(ToUtf8(v.name), ShowData{Name: ToUtf8(v.name), Exp: v.exp, Stat: "waiting"})
			}
		}
	}
	u.dataMap.Store(dm)
}

func (u *user) DataMap() map[string]*data {
	return u.dataMap.Load().(map[string]*data)
}

func (u *user) ShowMap() *sync.Map {
	return u.showMap
}

func (u *user) UpdateShowData(key string, showData ShowData) {
	u.showMap.Store(key, showData)
}

func (u *user) FeedWorker(w *worker) {
	dm := u.dataMap.Load().(map[string]*data)
	dur := float64(time.Hour*20) / (float64(time.Second) * float64(len(dm)*4))
	if dur < 1.0 {
		dur = 1.0
	}
	d := time.Duration(dur * float64(time.Second))
	log.Printf("dur %v, for user: %v, len(dataMap): %v", d, u.a.name, len(dm))
	ticker := time.Tick(d)
	for _, v := range dm {
		<-ticker
		w.AsyncNotify(&workerData{d: v, act: u.a, f: func(showData ShowData) {
			u.UpdateShowData(showData.Name, showData)
		}})
	}
}

type UserMgr struct {
	m    map[string]*user
	path string
}

func (u *UserMgr) Refresh() {
	dir, err := ioutil.ReadDir(u.path)
	if err != nil {
		log.Fatalln(err)
		return
	}
	wg := sync.WaitGroup{}
	for _, f := range dir {
		wg.Add(1)
		f := f
		go func() {
			defer func() {
				wg.Done()
			}()
			if f.IsDir() {
				return
			}
			a := NewAcount(path.Join(u.path, f.Name()))
			if a == nil || !a.valid {
				return
			}
			u.m[a.name] = NewUser(a)
		}()
	}
	wg.Wait()
}

func (u *UserMgr) UserMap() map[string]*user {
	return u.m
}

func NewUserManager(p string) (u *UserMgr) {
	u = &UserMgr{m: make(map[string]*user), path: p}
	u.Refresh()
	return
}

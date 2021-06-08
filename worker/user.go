package worker

import (
	"io/ioutil"
	"path"
	"sync"
	"sync/atomic"
)

type user struct {
	a       *account
	dataMap atomic.Value
	showMap sync.Map
	c       chan *workerData
}

func NewUser(a *account) *user {
	u := new(user)
	u.a = a
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

func (u *user) DataMap() map[string]*data {
	return u.dataMap.Load().(map[string]*data)
}

func (u *user) ShowMap() sync.Map {
	return u.showMap
}

func (u *user) UpdateShowData(key string, showData ShowData) {
	u.showMap.Store(key, showData)
}

type userMgr struct {
	m    map[string]*user
	path string
}

func (u *userMgr) Refresh() {
	dir, err := ioutil.ReadDir(u.path)
	if err != nil {
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

func (u *userMgr) UserMap() map[string]*user {
	return u.m
}

func NewUserManager(p string) (u *userMgr) {
	u = &userMgr{m: make(map[string]*user), path: p}
	u.Refresh()
	return
}

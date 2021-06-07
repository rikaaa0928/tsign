package worker

import (
	"io/ioutil"
	"path"
	"sync/atomic"
)

type user struct {
	a       *account
	dataMap atomic.Value
	showMap atomic.Value
	c       chan *workerData
}

func NewUser(a *account) *user {
	u := new(user)
	u.a = a
	dm := make(map[string]*data)
	l, err := a.GetList()
	if err == nil {
		for _, v := range l {
			dm[v.name] = v
		}
	}
	u.dataMap.Store(dm)
	sm := make(map[string]ShowData, len(dm))
	u.showMap.Store(sm)
	return u
}

func (u *user) DataMap() map[string]*data {
	return u.dataMap.Load().(map[string]*data)
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
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		a := NewAcount(path.Join(u.path, f.Name()))
		if a == nil || !a.valid {
			continue
		}
		u.m[a.name] = NewUser(a)
	}
}

func (u *userMgr) UserMap() map[string]*user {
	return u.m
}

func NewUserManager(p string) (u *userMgr) {
	u = &userMgr{m: make(map[string]*user), path: p}
	u.Refresh()
	return
}

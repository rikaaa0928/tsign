package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/rikaaa0928/tsign/worker"
)

type web struct {
	uma *atomic.Value
}

func (h *web) HandleIndex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	defer func() {
		err := recover()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%+v", err)))
		}
	}()
	umg := h.uma.Load().(*worker.UserMgr)
	um := umg.UserMap()
	//getTime
	ti := ""
	location := ""
	utc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("err: \n", err.Error())
		ti = err.Error()
	} else {
		tt := time.Now().In(utc)
		location = tt.Location().String()
		ti = tt.Format("2006-01-02 15:04:05")
	}

	//Creat menu
	menuStr := ""
	isUser := false
	keyList := make([]string, 0)
	for k := range um {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, x := range keyList {
		m := menuT{Name: x, Url: "?n=" + x}
		buf := new(bytes.Buffer)
		tplContent := ""
		if len(r.Form["n"]) > 0 && r.Form["n"][0] == x {
			tplContent = menuSelect
			isUser = true
		} else {
			tplContent = menu
		}
		t, e := template.New("menu").Parse(tplContent)
		if e != nil {
			log.Println(e)
			continue
		}
		e = t.Execute(buf, m)
		if e != nil {
			log.Println(e)
			continue
		}
		menuStr += buf.String()
	}
	in := indexT{Menu: menuStr, Time: ti, Location: location}

	//Create Panel
	in.Users = len(um)
	in.Errors = 0
	in.Counts = 0
	for _, x1 := range um {
		x1.ShowMap().Range(func(key, value interface{}) bool {
			in.Counts++
			if value.(worker.ShowData).Done || value.(worker.ShowData).Tried == 0 {
				return true
			}
			in.Errors++
			return true
		})
	}
	//Create body
	if isUser {
		in.Body = userBody(umg, r.Form["n"][0])
	} else {
		in.Body = indexBody(umg)
	}

	t, e := template.New("index").Parse(index)
	if e != nil {
		log.Println(e)
		fmt.Fprintln(w, "error:"+e.Error())
		return
	}
	t.Execute(w, in)
}

func Start(um *atomic.Value, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", (&web{um}).HandleIndex)
	mux.HandleFunc("/template/", func(writer http.ResponseWriter, request *http.Request) {
		if request.RequestURI == "/template/w3.css" {
			writer.Header().Add("content-type", "text/css; charset=utf-8")
			writer.Write(w3)
		}
	})
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

func indexBody(umgr *worker.UserMgr) (b string) {
	str := ""
	keyList := make([]string, 0)
	for k := range umgr.UserMap() {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		v := umgr.UserMap()[k]
		if !v.Account().IsValid() {
			str += makeListI(k, "Error !", "fa-user")
		} else {
			str += makeListI(k, "Fine !", "fa-user")
		}
	}
	b = makeList("UserList", str)

	str = ""
	keyList = make([]string, 0)
	for k := range umgr.UserMap() {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		v := umgr.UserMap()[k]
		n := 0
		m := 0
		v.ShowMap().Range(func(_, value interface{}) bool {
			if value.(worker.ShowData).Done {
				n++
			}
			m++
			return true
		})
		str += makeProgressI(float64(n)/float64(m)*100, strconv.Itoa(n)+"/"+strconv.Itoa(m), k)
	}
	b += makeProgress("Finished", str)
	return
}

func userBody(umgr *worker.UserMgr, user string) (b string) {
	str := ""
	b = ""
	for k, v := range umgr.UserMap() {
		if k == user {
			var list sortList
			v.ShowMap().Range(func(key, value interface{}) bool {
				fk := sortTy{Exp: value.(worker.ShowData).Exp, Name: key.(string)}
				list = append(list, &fk)
				return true
			})
			sort.Sort(list)
			for _, tb := range list {
				st, ok := v.ShowMap().Load(tb.Name)
				if !ok {
					continue
				}
				if !st.(worker.ShowData).Done && st.(worker.ShowData).Tried > 1 {
					str += makeListI(tb.Name, st.(worker.ShowData).Stat, "fa-comment")
				} else {
					str += makeListI(tb.Name, fmt.Sprintf("%s; Exp: %d!", st.(worker.ShowData).Stat, st.(worker.ShowData).Exp), "fa-comment")
				}
			}
		}
	}
	b = makeList("TiebaList", str)
	return
}

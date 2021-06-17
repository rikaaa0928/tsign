package web

import (
	"bytes"
	"log"
	"text/template"
)

type menuT struct {
	Name string
	Url  string
}

type indexT struct {
	Location string
	Time     string
	Menu     string
	Errors   int
	Counts   int
	Users    int
	Body     string
}

type listT struct {
	Name   string
	Inside string
}

type listIT struct {
	X string
	Y string
	T string
}
type proT struct {
	Title  string
	Inside string
}

type proIT struct {
	P float64
	N string
	T string
}

func makeListI(x, y, tp string) (b string) {
	a := listIT{X: x, Y: y, T: tp}
	buf := new(bytes.Buffer)
	t, e := template.New("listInside").Parse(listInside)
	if e != nil {
		log.Println(e)
	}
	e = t.Execute(buf, a)
	if e != nil {
		log.Println(e)
	}
	b = buf.String()
	return
}
func makeList(name, inside string) (b string) {
	a := listT{Name: name, Inside: inside}
	buf := new(bytes.Buffer)
	t, e := template.New("list").Parse(list)
	if e != nil {
		log.Println(e)
		return ""
	}
	e = t.Execute(buf, a)
	if e != nil {
		log.Println(e)
		return ""
	}
	b = buf.String()
	return
}
func makeProgressI(p float64, n string, ti string) (b string) {
	a := proIT{P: p, N: n, T: ti}
	buf := new(bytes.Buffer)
	t, e := template.New("proInside").Parse(proInside)
	if e != nil {
		log.Println(e)
	}
	e = t.Execute(buf, a)
	if e != nil {
		log.Println(e)
	}
	b = buf.String()
	return
}
func makeProgress(title, inside string) (b string) {
	a := proT{Title: title, Inside: inside}
	buf := new(bytes.Buffer)
	t, e := template.New("progress").Parse(progress)
	if e != nil {
		log.Println(e)
		return ""
	}
	e = t.Execute(buf, a)
	if e != nil {
		log.Println(e)
		return ""
	}
	b = buf.String()
	return
}

type sortTy struct {
	Exp  int
	Name string
}
type sortList []*sortTy

func (I sortList) Len() int {
	return len(I)
}
func (I sortList) Less(i, j int) bool {
	return I[i].Exp > I[j].Exp
}
func (I sortList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

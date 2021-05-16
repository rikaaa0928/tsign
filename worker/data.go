package worker

import (
	"fmt"
	"regexp"
	"strconv"
)

type data struct {
	id          int
	name        string
	unicodeName string
	exp         int
	tried       int
}

func NewData(html string) *data {
	d := data{}
	exp := regexp.MustCompile("<a href=\"/f\\?kw=(.*?)\" title=\"(.*?)\"")
	names := exp.FindStringSubmatch(html)
	if names == nil {
		return nil
	}
	d.unicodeName = names[1]
	d.name = names[2]
	exp = regexp.MustCompile("<a class=\"cur_exp\".+?>(\\d+)</a>")
	d.exp, _ = strconv.Atoi(exp.FindStringSubmatch(html)[1])
	exp = regexp.MustCompile("balvid=\"(\\d+)\"")
	d.id, _ = strconv.Atoi(exp.FindStringSubmatch(html)[1])
	return &d
}

type ShowData struct {
	Name  string
	Exp   int
	Tried int
	Stat  string
}

func (d ShowData) String() string {
	return fmt.Sprintf("%s: %s. exp: %d, tried %d times.", d.Name, d.Stat, d.Exp, d.Tried)
}

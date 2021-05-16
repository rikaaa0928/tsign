package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type User struct {
	cookieJar *cookiejar.Jar
	valid     bool
	name      string
}

func NewUser(file string) *User {
	u := &User{}
	var err error
	u.name = strings.TrimSuffix(path.Base(file), ".ts")

	u.cookieJar, err = cookiejar.New(nil)
	if err != nil {
		u.valid = false
		return u
	}
	cookies := make([]*http.Cookie, 0)
	if _, err := os.Stat(file); err == nil {
		rawCookie, err := ioutil.ReadFile(file)
		if err != nil {
			u.valid = false
			return u
		}
		rawCookie = bytes.Trim(rawCookie, "\xef\xbb\xbf")
		rawCookieList := strings.Split(strings.Replace(string(rawCookie), "\r\n", "\n", -1), "\n")
		for _, rawCookieLine := range rawCookieList {
			rawCookieInfo := strings.SplitN(rawCookieLine, "=", 2)
			if len(rawCookieInfo) < 2 {
				continue
			}
			cookies = append(cookies, &http.Cookie{
				Name:   rawCookieInfo[0],
				Value:  rawCookieInfo[1],
				Domain: ".baidu.com",
			})
		}
		log.Printf("Verifying imported cookies from %s...", file)
		URL, _ := url.Parse("http://baidu.com")
		u.cookieJar.SetCookies(URL, cookies)
		u.valid = u.IsLogIn()
		log.Printf("%v\n", u.valid)
	}
	return u
}

func (u User) IsValid() bool {
	return u.valid
}

func (u User) Cookie() *cookiejar.Jar {
	return u.cookieJar
}

func (u *User) IsLogIn() (r bool) {
	defer func() {
		u.valid = r
	}()
	body, err := Fetch("http://tieba.baidu.com/dc/common/tbs", nil, u.cookieJar)
	if err != nil {
		return false
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return false
	}
	v, ok := m["is_login"]
	if !ok {
		return false
	}
	_, ok = v.(float64)
	if ok {
		return true
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err = strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func (u *User) GetList() ([]*data, error) {
	pn := 0
	list := make([]*data, 0)
	for {
		pn++
		urlStr := "http://tieba.baidu.com/f/like/mylike?pn=" + fmt.Sprintf("%d", pn)
		body, fetchErr := Fetch(urlStr, nil, u.cookieJar)
		if fetchErr != nil {
			return nil, fetchErr
		}
		reg := regexp.MustCompile("<tr><td>.+?</tr>")
		allTr := reg.FindAllString(string(body), -1)
		for _, line := range allTr {
			nData := NewData(line)
			if nData == nil {
				continue
			}
			list = append(list, nData)
		}
		if allTr == nil {
			break
		}
	}
	return list, nil
}

func (u User) GetTBS() string {
	body, err := Fetch("http://tieba.baidu.com/dc/common/tbs", nil, u.cookieJar)
	if err != nil {
		return ""
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return ""
	}
	v, ok := m["tbs"]
	if !ok {
		return ""
	}
	str, ok := v.(string)
	if !ok {
		return ""
	}
	return str
}

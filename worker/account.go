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

type account struct {
	cookieJar *cookiejar.Jar
	valid     bool
	name      string
}

func NewAcount(file string) *account {
	u := &account{}
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

func (a account) IsValid() bool {
	return a.valid
}

func (a account) Cookie() *cookiejar.Jar {
	return a.cookieJar
}

func (a *account) IsLogIn() (r bool) {
	defer func() {
		a.valid = r
	}()
	body, err := Fetch("http://tieba.baidu.com/dc/common/tbs", nil, a.cookieJar)
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
	f, ok := v.(float64)
	if ok {
		return int(f) == 1
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return i == 1
}

func (a *account) GetList() ([]*data, error) {
	pn := 0
	list := make([]*data, 0)
	for {
		pn++
		urlStr := "http://tieba.baidu.com/f/like/mylike?pn=" + fmt.Sprintf("%d", pn)
		body, fetchErr := Fetch(urlStr, nil, a.cookieJar)
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

func (a account) GetTBS() string {
	body, err := Fetch("http://tieba.baidu.com/dc/common/tbs", nil, a.cookieJar)
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

func (a account) GetCookie(name string) string {
	cookieUrl, _ := url.Parse("http://tieba.baidu.com")
	cookies := a.cookieJar.Cookies(cookieUrl)
	for _, cookie := range cookies {
		if name == cookie.Name {
			return cookie.Value
		}
	}
	return ""
}

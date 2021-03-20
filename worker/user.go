package worker

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
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
		u.valid = GetLoginStatus(u.cookieJar)
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

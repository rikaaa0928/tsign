package functions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type SignData struct {
	Cookies   []*http.Cookie
	ID        int
	Name      string
	cookieJar *cookiejar.Jar
}

func (u *SignData) init() (err error) {
	u.cookieJar, err = cookiejar.New(nil)
	if err != nil {
		return
	}
	URL, _ := url.Parse("http://baidu.com")
	u.cookieJar.SetCookies(URL, u.Cookies)
	return
}

func (a SignData) GetTBS() string {
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

func (a SignData) GetCookie(name string) string {
	cookieUrl, _ := url.Parse("http://tieba.baidu.com")
	cookies := a.cookieJar.Cookies(cookieUrl)
	for _, cookie := range cookies {
		if name == cookie.Name {
			return cookie.Value
		}
	}
	return ""
}

func Fetch(targetUrl string, postData map[string]string, ptrCookieJar *cookiejar.Jar) ([]byte, error) {
	var request *http.Request
	httpClient := &http.Client{
		Jar: ptrCookieJar,
	}
	if nil == postData {
		request, _ = http.NewRequest("GET", targetUrl, nil)
	} else {
		postParams := url.Values{}
		for key, value := range postData {
			postParams.Set(key, value)
		}
		postDataStr := postParams.Encode()
		postDataBytes := []byte(postDataStr)
		postBytesReader := bytes.NewReader(postDataBytes)
		request, _ = http.NewRequest("POST", targetUrl, postBytesReader)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	response, fetchError := httpClient.Do(request)
	if fetchError != nil {
		return nil, fetchError
	}
	defer response.Body.Close()
	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}
	return body, nil
}

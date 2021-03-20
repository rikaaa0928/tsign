package worker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

func GetLoginStatus(ptrCookieJar *cookiejar.Jar) bool {
	body, err := Fetch("http://tieba.baidu.com/dc/common/tbs", nil, ptrCookieJar)
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

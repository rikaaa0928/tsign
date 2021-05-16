package worker

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

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

func Sign(u *User, d data) (r ShowData) {
	postData := make(map[string]string)
	postData["BDUSS"] = GetCookie(ptrCookieJar, "BDUSS")
	postData["_client_id"] = "03-00-DA-59-05-00-72-96-06-00-01-00-04-00-4C-43-01-00-34-F4-02-00-BC-25-09-00-4E-36"
	postData["_client_type"] = "4"
	postData["_client_version"] = "1.2.1.17"
	postData["_phone_imei"] = "540b43b59d21b7a4824e1fd31b08e9a6"
	postData["fid"] = fmt.Sprintf("%d", tieba.TiebaId)
	postData["kw"] = tieba.Name
	postData["net_type"] = "3"
	postData["tbs"] = getTbs(ptrCookieJar)

	var keys []string
	for key := range postData {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))

	sign_str := ""
	for _, key := range keys {
		sign_str += fmt.Sprintf("%s=%s", key, postData[key])
	}
	sign_str += "tiebaclient!!!"

	MD5 := md5.New()
	MD5.Write([]byte(sign_str))
	MD5Result := MD5.Sum(nil)
	signValue := make([]byte, 32)
	hex.Encode(signValue, MD5Result)
	postData["sign"] = strings.ToUpper(string(signValue))

	body, fetchErr := Fetch("http://c.tieba.baidu.com/c/c/forum/sign", postData, u.cookieJar)
	if fetchErr != nil {
		r.Stat = fetchErr.Error()
		return
	}
	json, parseErr := NewJson([]byte(body))
	if parseErr != nil {
		return 1, parseErr.Error(), 0
	}
	if _exp, succeed := json.Get("user_info").CheckGet("sign_bonus_point"); succeed {
		exp, _ := strconv.Atoi(_exp.MustString())
		return 2, fmt.Sprintf("签到成功，获得经验值 %d", exp), exp
	}
	switch json.Get("error_code").MustString() {
	case "340010":
		fallthrough
	case "160002":
		fallthrough
	case "3":
		return 2, "你已经签到过了", 0
	case "1":
		fallthrough
	case "340008": // 黑名单
		fallthrough
	case "340006": // 被封啦
		fallthrough
	case "160004":
		return -1, fmt.Sprintf("ERROR-%s: %s", json.Get("error_code").MustString(), json.Get("error_msg").MustString()), 0
	case "160003":
		fallthrough
	case "160008":
		fallthrough
	default:
		return 1, fmt.Sprintf("ERROR-%s: %s", json.Get("error_code").MustString(), json.Get("error_msg").MustString()), 0
	}
	return -255, "", 0
	return
}

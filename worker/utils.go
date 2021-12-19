package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/rikaaa0928/tsign/functions"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

func Sign(a *account, d *data) (r ShowData) {
	defer func(d *data) {
		r.Done = d.done
	}(d)
	d.tried++
	r.Exp = d.exp
	r.Tried = d.tried
	r.Name = ToUtf8(d.name)
	//postData := make(map[string]string)
	//postData["BDUSS"] = a.GetCookie("BDUSS")
	//postData["_client_id"] = "03-00-DA-59-05-00-72-96-06-00-01-00-04-00-4C-43-01-00-34-F4-02-00-BC-25-09-00-4E-36"
	//postData["_client_type"] = "4"
	//postData["_client_version"] = "1.2.1.17"
	//postData["_phone_imei"] = "540b43b59d21b7a4824e1fd31b08e9a6"
	//postData["fid"] = fmt.Sprintf("%d", d.id)
	//postData["kw"] = d.name
	//postData["net_type"] = "3"
	//postData["tbs"] = a.GetTBS()
	//
	//var keys []string
	//for key := range postData {
	//	keys = append(keys, key)
	//}
	//sort.Sort(sort.StringSlice(keys))
	//
	//sign_str := ""
	//for _, key := range keys {
	//	sign_str += fmt.Sprintf("%s=%s", key, postData[key])
	//}
	//sign_str += "tiebaclient!!!"
	//MD5 := md5.New()
	//MD5.Write([]byte(sign_str))
	//MD5Result := MD5.Sum(nil)
	//signValue := make([]byte, 32)
	//hex.Encode(signValue, MD5Result)
	//postData["sign"] = strings.ToUpper(string(signValue))
	//
	//body, fetchErr := Fetch("http://c.tieba.baidu.com/c/c/forum/sign", postData, a.cookieJar)
	//
	//fmt.Printf("%d,%d,%s\n", len(d.name), len(ToUtf8(d.name)), ToUtf8(d.name))
	post := functions.SignData{
		Cookies: a.Cookies,
		ID:      d.id,
		Name:    ToUtf8(d.name),
	}
	postBytes, err := json.Marshal(post)
	if err != nil {
		r.Stat = err.Error()
		return
	}
	pData := make(map[string]string)
	pData["data"] = string(postBytes)

	body, fetchErr := Fetch(os.Getenv("SIGN_FUNC"), pData, a.cookieJar)
	if fetchErr != nil {
		r.Stat = fetchErr.Error()
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	//fmt.Println(string(body))
	//json, parseErr := NewJson([]byte(body))
	if err != nil {
		r.Stat = err.Error()
		return
	}
	o1, ok := m["user_info"]
	if ok {
		m1, ok := o1.(map[string]interface{})
		if !ok {
			r.Stat = "error response structure user_info map: " + string(body)
			return
		}
		expInt, ok := m1["sign_bonus_point"]
		if ok {
			expStr, ok := expInt.(string)
			if !ok {
				r.Stat = "error response structure sign_bonus_point string: " + string(body)
				return
			}
			_, err = strconv.Atoi(expStr)
			if err != nil {
				r.Stat = "error response structure sign_bonus_point int: " + string(body)
				return
			}
			d.done = true
			r.Stat = "done and get exp: " + expStr
			return
		}
	}
	errCode, ok := m["error_code"]
	if !ok {
		r.Stat = "error response structure: " + string(body)
		return
	}
	errCodeStr, ok := errCode.(string)
	if ok && errCodeStr == "160002" {
		r.Stat = "done."
		d.done = true
		return
	}
	r.Stat = fmt.Sprintf("error code: %s, error msg: %s", m["error_code"], m["error_msg"])
	return
}

func ToUtf8(gbkString string) string {
	I := bytes.NewReader([]byte(gbkString))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, _ := ioutil.ReadAll(O)
	return string(d)
}

func Unicode2ZHCN(unicode string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(unicode), `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}

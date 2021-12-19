package functions

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// HelloWorld writes "Hello, World!" to the HTTP response.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	str := r.FormValue("data")
	data := &SignData{}
	err = json.Unmarshal([]byte(str), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tn := data.Name
	data.Name = Utf8ToGbk(tn)
	//fmt.Printf("%d,%d,%s,%s\n", len(data.Name), len(tn), GbkToUtf8(data.Name), tn)
	err = data.init()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postData := make(map[string]string)
	postData["BDUSS"] = data.GetCookie("BDUSS")
	postData["_client_id"] = "03-00-DA-59-05-00-72-96-06-00-01-00-04-00-4C-43-01-00-34-F4-02-00-BC-25-09-00-4E-36"
	postData["_client_type"] = "4"
	postData["_client_version"] = "1.2.1.17"
	postData["_phone_imei"] = "540b43b59d21b7a4824e1fd31b08e9a6"
	postData["fid"] = fmt.Sprintf("%d", data.ID)
	postData["kw"] = data.Name
	postData["net_type"] = "3"
	postData["tbs"] = data.GetTBS()

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

	body, err := Fetch("http://c.tieba.baidu.com/c/c/forum/sign", postData, data.cookieJar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))
	fmt.Fprint(w, string(body))
}

func GbkToUtf8(gbkString string) string {
	I := bytes.NewReader([]byte(gbkString))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, _ := ioutil.ReadAll(O)
	return string(d)
}

func Utf8ToGbk(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	d, _ := ioutil.ReadAll(reader)
	return string(d)
}

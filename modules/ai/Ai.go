package ai

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Ai struct {
	app_id string
	app_key string
	session string
	question string
	time_stamp string
	nonce_str string
}

type Answer struct {
	ret int
	msg string
	data map[string]string
}

func (this *Ai) SetAppInfo(app_id, app_key string)  {
	this.app_id = app_id
	this.app_key = app_key
	this.time_stamp = strconv.FormatInt(time.Now().Unix(), 10)
	this.nonce_str = GetRandomString(16)
}

func (this *Ai) NlpTextchat(httpClient *http.Client, question string, session string) {
	this.session = session
	this.question = question
	var requestURL = "https://api.ai.qq.com/fcgi-bin/nlp/nlp_textchat"
	this.httpPostJSON(httpClient, requestURL, this.getRequestBody())

}

func (this *Ai) httpPostJSON(clt *http.Client, url string, requestBody []byte/*, response interface{}*/) error {
	httpResp, err := clt.Post(url, "", bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	fmt.Println(string(body))
	return err
	//return json.NewDecoder(httpResp.Body).Decode(response)
}

func (this *Ai) getRequestBody() []byte{
	var newMp = make([]string, 0)
	t := reflect.TypeOf(Ai{})
	for k := 0; k < t.NumField(); k++ {
		key := t.Field(k).Name
		if key == "app_key" {
			continue
		}
		newMp = append(newMp, key)
	}
	sort.Strings(newMp)
	buf := make([]byte, 0, 150)
	for _, v := range newMp {
		query := url.Values{}
		query.Add(v, reflect.ValueOf(*this).FieldByName(v).String())
		buf = append(buf, (query.Encode())+"&"...)
	}
	requestBody := string(buf)
	buf = append(buf, ("app_key=" + this.app_key)...)
	hashsum := md5.Sum(buf)
	sign := strings.ToUpper(hex.EncodeToString(hashsum[:]))
	requestBody += "&sign=" + sign
	fmt.Println(requestBody)
	return []byte(requestBody)
}

func GetRandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
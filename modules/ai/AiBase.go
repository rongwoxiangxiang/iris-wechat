package ai

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"iris/config"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AiBase struct {
	app_id string
	app_key string
	session string
	question string
	time_stamp string
	nonce_str string
}

func (this *AiBase) SetAiServer(){
	this.app_id = config.GetConfigs().OtherConfigs.QQAiServer.AppId
	this.app_key = config.GetConfigs().OtherConfigs.QQAiServer.AppKey
}


func (this *AiBase) httpPostJSON(url string, requestBody string, response interface{}) error {
	httpResp, err := (&http.Client{}).Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}
	return json.NewDecoder(httpResp.Body).Decode(response)
}

func (this *AiBase) getRequestBody() string{

	this.nonce_str = GetRandomString(16)
	this.time_stamp = strconv.FormatInt(time.Now().Unix(), 10)

	var newMp = make([]string, 0)
	t := reflect.TypeOf(AiBase{})
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
	requestBody += "sign=" + strings.ToUpper(hex.EncodeToString(hashsum[:]))
	return requestBody
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
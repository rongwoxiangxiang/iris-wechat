package ai

import (
	"net/http"
)

type Ai struct {
	AiBase
}

type Answer struct {
	ErrCode int `json:"ret"`
	Msg string `json:"msg"`
	DataJson `json:"data"`
}

type DataJson struct {
	Session string `json:"session"`
	AnswerData string `json:"answer"`
}

//智能闲聊
func (this *Ai) NlpTextchat(httpClient *http.Client, question string, session string) (answer Answer) {
	this.session = session
	this.question = question
	this.httpPostJSON("https://api.ai.qq.com/fcgi-bin/nlp/nlp_textchat",
		this.getRequestBody(),
		&answer)
	return
}
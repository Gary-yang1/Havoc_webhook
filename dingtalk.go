package main

import (
	"crypto/hmac"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Settings struct {
	AccessToken string `json:"access_token"`
	Secret      string `json:"Secret"`
}

type Robot struct {
	Webhook string
	Secret  string
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Embeds struct {
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
}

type Message struct {
	username   string   `json:"username"`
	avatar_url string   `json:"avatar_url"`
	Embeds     []Embeds `json:"embeds"`
}

func (r *Robot) sign(timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, r.Secret)
	h := hmac.New(sha256.New, []byte(r.Secret))
	h.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sign
}

//go:embed settings.json
var Json []byte

func dingtalk(message Message) {
	var settings Settings
	err := json.Unmarshal(Json, &settings)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	robot := Robot{
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=" + settings.AccessToken,
		Secret:  settings.Secret,
	}
	timestamp := time.Now().Unix() * 1000 // 时间戳
	signs := robot.sign(timestamp)        // 签名
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", robot.Webhook, timestamp, signs)
	var demondID, user, pc, os, iip, callbacktime string

	fmt.Println(message)

	for _, field := range message.Embeds[0].Fields {
		switch field.Name {
		case "Agent ID":
			demondID = field.Value
		case "Username":
			user = field.Value
		case "Hostname":
			pc = field.Value
		case "OS Version":
			os = field.Value
		case "Internal IP":
			iip = field.Value
		case "First Callback":
			callbacktime = field.Value
		}
	}
	title := "Havoc上线提醒"
	text := fmt.Sprintf(
		"## Havoc有新主机上线，请注意\n\n 主机ID: %s\n\n用户: %s\n\n计算机名称: %s\n\n操作系统: %s\n\n内网IP: %s\n\n上线时间: %s",
		demondID, user, pc, os, iip, callbacktime,
	)
	requestData := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}
	data, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

}

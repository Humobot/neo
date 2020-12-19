package xunfei

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "github.com/hajimehoshi/oto"
	// "github.com/tosone/minimp3"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"errors"
	"strings"
	// "sync"
	"time"

	"humobot/neo/pkg/websocket"
)

/**
 * 语音听写流式 WebAPI 接口调用示例 接口文档（必看）：https://doc.xfyun.cn/rest_api/语音听写（流式版）.html
 * webapi 听写服务参考帖子（必看）：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=38947&extra=
 * 语音听写流式WebAPI 服务，热词使用方式：登陆开放平台https://www.xfyun.cn/后，找到控制台--我的应用---语音听写---服务管理--上传热词
 * 注意：热词只能在识别的时候会增加热词的识别权重，需要注意的是增加相应词条的识别率，但并不是绝对的，具体效果以您测试为准。
 * 错误码链接：https://www.xfyun.cn/document/error-code （code返回错误码时必看）
 */

type Xunfei struct {
	hostUrl string
    host string
    appid string
    apiSecret string
	apiKey string
}

const (
	STATUS_FIRST_FRAME    = 0
	STATUS_CONTINUE_FRAME = 1
	STATUS_LAST_FRAME     = 2
)

type RespData struct {
	Sid     string `json:"sid"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Audio  string `json:"audio,omitempty"`
	Ced    int    `json:"ced,omitempty"`
	Status int    `json:"status,omitempty"`
}


func NewXunfei(hostUrl, host, appid, apiSecret, apiKey string) *Xunfei {
	return &Xunfei{
		hostUrl: hostUrl,
		host: host,
    	appid: appid,
    	apiSecret: apiSecret,
		apiKey: apiKey,
	}
}

func (x *Xunfei) makeAuthURL(uri string) string {
	return assembleAuthUrl(x.hostUrl + uri, x.apiKey, x.apiSecret)
}

func (x *Xunfei) TtsOnline(word string, filePath string) error {
	client := websocket.NewClient(x.makeAuthURL("/v2/tts"))

	frameData := map[string]interface{}{
		"common": map[string]interface{}{
			"app_id": x.appid,
		},
		"business": map[string]interface{}{
			"vcn":   "xiaoyan",
			"aue":   "lame",
			"speed": 50,
			"tte":   "UTF8",
			"sfl":   1,
		},
		"data": map[string]interface{}{
			"status":   STATUS_LAST_FRAME,
			"encoding": "UTF8",
			"text":     base64.StdEncoding.EncodeToString([]byte(word)),
		},
	}

	audioFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		audioFile.Close()
		panic(err)
	}

	client.Send(frameData, func(msg []byte, done chan string) {
		// now := time.Now()
		var err error
		var resp = RespData{}
		json.Unmarshal(msg, &resp)
		//fmt.Println(string(msg))
		//fmt.Println(resp.Data.Audio, resp.Sid)
		if resp.Code != 0 {
			// fmt.Println(resp.Code, resp.Message, time.Since(now))
			done <- "done"
			err = errors.New(resp.Message)
		}
		//decoder.Decode(&resp.Data.Audio)

		audiobytes, err := base64.StdEncoding.DecodeString(resp.Data.Audio)
		if err != nil {
			return
		}
		_, err = audioFile.Write(audiobytes)
		if err != nil {
			return
		}

		if resp.Data.Status == 2 {
			//fmt.Println("final:",decoder.String())
			// fmt.Println(resp.Code, resp.Message, time.Since(now))
			audioFile.Close()
			done <- "done"
		}

	})

	return err
}

//创建鉴权url  apikey 即 hmac username. thk copy from https://github.com/lingjiao0710/iat_ws_go_demo/blob/5de5296915c31aa2600fb05b7f81b385aa94c294/ttsonline/main.go
func assembleAuthUrl(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	fmt.Println(sgin)
	//签名结果
	sha := hmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func hmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
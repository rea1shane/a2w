package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	myHttp "github.com/rea1shane/gooooo/http"
	"github.com/rea1shane/gooooo/log"
	myTime "github.com/rea1shane/gooooo/time"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Notification Alertmanager 发送的告警通知
type Notification struct {
	Receiver string  `json:"receiver"`
	Status   string  `json:"status"`
	Alerts   []Alert `json:"alerts"`

	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// Alert 告警实例
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

const (
	webhookUrl     = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	okMsg          = `{"errcode":0,"errmsg":"ok"}`
	markdownMaxLen = 4096 // markdownMaxLen 企业微信 Markdown 消息体最大长度为 4096
)

var (
	tmplPath, tmplName string
)

func main() {
	port := flag.Int("port", 5001, "监听端口")
	flag.StringVar(&tmplPath, "template", "./templates/base.tmpl", "模板文件")
	flag.Parse()

	split := strings.Split(tmplPath, "/")
	tmplName = split[len(split)-1]

	logger := logrus.New()
	formatter := log.NewFormatter()
	formatter.FieldsOrder = []string{"StatusCode", "Latency"}
	logger.SetFormatter(formatter)

	app := myHttp.NewHandler(logger, 0)

	app.GET("/", health)
	app.POST("/send", send)

	app.Run(fmt.Sprintf("0.0.0.0:%d", *port))
}

// health 健康检查
func health(c *gin.Context) {
	c.Writer.WriteString("ok")
}

// send 发送消息
func send(c *gin.Context) {
	// 获取 bot key
	key := c.Query("key")

	// 解析 Alertmanager 消息
	decoder := json.NewDecoder(c.Request.Body)
	var notification *Notification
	if err := decoder.Decode(&notification); err != nil {
		e := c.Error(err)
		e.Meta = "解析 Alertmanager 消息失败"
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// 填充模板
	var tfm = make(template.FuncMap)
	tfm["timeFormat"] = timeFormat
	tfm["timeDuration"] = timeDuration
	tfm["timeFromNow"] = timeFromNow
	tmpl := template.Must(template.New(tmplName).Funcs(tfm).ParseFiles(tmplPath))
	var content bytes.Buffer
	if err := tmpl.Execute(&content, notification); err != nil {
		e := c.Error(err)
		e.Meta = "填充模板失败"
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 消息分段
	// 为了解决企业微信 Markdown 消息体长度限制问题
	var msgs []string
	if content.Len() <= markdownMaxLen {
		msgs = append(msgs, content.String())
	} else {
		// 分段消息标识头
		snippetHeader := `<font color="comment">**(%d/%d)**</font>`

		// 单条分段最大长度，预留一些空间用于添加分段头和容错
		snippetMaxLen := markdownMaxLen - 128

		// 消息切割
		// 企业微信中，连续至少三个的换行符才被视为两个换行符
		fragments := strings.Split(content.String(), "\n\n\n")

		var snippetBuilder strings.Builder
		snippetBuilder.Grow(markdownMaxLen)

		// 拼接消息
		for _, fragment := range fragments {
			// 切割后的单条消息都过长
			if len(fragment) > snippetMaxLen {
				e := c.Error(errors.New(fmt.Sprintf("切割后的消息长度 %d 仍超出片段长度限制 %d", len(fragment), snippetMaxLen)))
				e.Meta = "分段消息失败"
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}

			// 拼接消息后超出限制长度
			if snippetBuilder.Len()+len(fragment) > snippetMaxLen {
				msgs = append(msgs, snippetBuilder.String())
				snippetBuilder.Reset()
				snippetBuilder.Grow(markdownMaxLen)
			}

			snippetBuilder.WriteString("\n\n\n")
			snippetBuilder.WriteString(fragment)
		}

		msgs = append(msgs, snippetBuilder.String())

		// 添加分段头
		for index, snippet := range msgs {
			snippetBuilder.Reset()
			snippetBuilder.Grow(markdownMaxLen)
			snippetBuilder.WriteString(fmt.Sprintf(snippetHeader, index+1, len(msgs)))
			snippetBuilder.WriteString(snippet)
			msgs[index] = snippetBuilder.String()
		}
	}

	for _, msg := range msgs {
		// 请求企业微信
		postBody, _ := json.Marshal(map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]interface{}{
				"content": msg,
			},
		})
		postBodyBuffer := bytes.NewBuffer(postBody)
		wecomResp, err := http.Post(webhookUrl+key, "application/json", postBodyBuffer)
		if err != nil {
			e := c.Error(err)
			e.Meta = "发起企业微信请求失败"
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 处理请求结果
		wecomRespBody, _ := io.ReadAll(wecomResp.Body)
		wecomResp.Body.Close()
		if wecomResp.StatusCode != http.StatusOK || string(wecomRespBody) != okMsg {
			e := c.Error(errors.New(string(wecomRespBody)))
			e.Meta = "请求企业微信失败，HTTP Code: " + strconv.Itoa(wecomResp.StatusCode)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// timeFormat 格式化时间
func timeFormat(t time.Time) string {
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

// timeDuration 计算结束时间距开始时间的时间差
func timeDuration(startTime, endTime time.Time) string {
	duration := endTime.Sub(startTime)
	return myTime.FormatDuration(duration)
}

// timeFromNow 计算当前时间距开始时间地时间差
func timeFromNow(startTime time.Time) string {
	return timeDuration(startTime, time.Now())
}

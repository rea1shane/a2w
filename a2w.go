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
	"io/ioutil"
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
	webhookUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
)

var (
	tmpl, tmplName string
)

func main() {
	port := flag.Int("port", 5001, "监听端口")

	flag.StringVar(&tmpl, "template", "./templates/base.tmpl", "模板文件")
	split := strings.Split(tmpl, "/")
	tmplName = split[len(split)-1]

	logger := logrus.New()
	formatter := log.GetFormatter()
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
	tmpl := template.Must(template.New(tmplName).Funcs(tfm).ParseFiles(tmpl))
	var content bytes.Buffer
	if err := tmpl.Execute(&content, notification); err != nil {
		e := c.Error(err)
		e.Meta = "填充模板失败"
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 请求企业微信
	postBody, _ := json.Marshal(map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content.String(),
		},
	})
	postBodyBuffer := bytes.NewBuffer(postBody)
	wecomResp, err := http.Post(webhookUrl+key, "application/json", postBodyBuffer)
	if err != nil {
		e := c.Error(err)
		e.Meta = "发起企业微信请求错误"
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer wecomResp.Body.Close()

	// 处理请求结果
	wecomRespBody, _ := ioutil.ReadAll(wecomResp.Body)
	if wecomResp.StatusCode != http.StatusOK {
		e := c.Error(errors.New(string(wecomRespBody)))
		e.Meta = "请求企业微信失败，HTTP Code: " + strconv.Itoa(wecomResp.StatusCode)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
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

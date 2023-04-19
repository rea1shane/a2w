package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rea1shane/gooooo/log"
	myTime "github.com/rea1shane/gooooo/time"
	"io/ioutil"
	"net/http"
	"strconv"
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
	logger = log.GetLogger()
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.GET("/", health)
	app.POST("/send", send)

	(&http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9099),
		Handler: app,
	}).ListenAndServe()
}

// health 健康检查
func health(c *gin.Context) {
	c.Writer.WriteString("ok")
}

// send 发送消息
func send(c *gin.Context) {
	// 获取 bot key
	key := c.Query("key")

	// 解析 Alertmanager 通知
	decoder := json.NewDecoder(c.Request.Body)
	var notification *Notification
	if err := decoder.Decode(&notification); err != nil {
		logger.Error("解析 Alertmanager 消息错误: " + err.Error())
		return
	}

	// 填充模板
	var tfm = make(template.FuncMap)
	tfm["timeFormat"] = timeFormat
	tfm["timeDuration"] = timeDuration
	tfm["timeFromNow"] = timeFromNow
	tmpl := template.Must(template.New("message.tmpl").Funcs(tfm).ParseFiles("./templates/message.tmpl"))
	var content bytes.Buffer
	if err := tmpl.Execute(&content, notification); err != nil {
		logger.Error("填充模板错误: " + err.Error())
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
		logger.Error("发起企业微信请求错误: " + err.Error())
		return
	}
	defer wecomResp.Body.Close()

	// 处理请求结果
	wecomRespBody, _ := ioutil.ReadAll(wecomResp.Body)
	if wecomResp.StatusCode != http.StatusOK {
		logger.Error("请求企业微信失败，HTTP Code: ", strconv.Itoa(wecomResp.StatusCode), " 返回体: ", string(wecomRespBody))
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write(wecomRespBody)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(wecomRespBody)
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

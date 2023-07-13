# Alertmanager to WeCom

[![docker-image-ci](https://github.com/rea1shane/a2w/actions/workflows/docker-image-ci.yml/badge.svg)](https://github.com/rea1shane/a2w/actions/workflows/docker-image-ci.yml)

通过 [企业微信机器人](https://developer.work.weixin.qq.com/document/path/91770) 发送 [Alertmanager](https://github.com/prometheus/alertmanager) 通知。

## 使用方式

1. 运行本项目：

   ```shell
   make run
   ```

   或是部署在 Docker 中：

   ```shell
   make docker-run
   ```

   或是直接使用 [Docker Hub](https://hub.docker.com/r/rea1shane/a2w) 中已发布的镜像：

   ```shell
   docker run --name a2w -d -p 5001:5001 rea1shane/a2w
   ```

   查看项目使用说明：

   ```shell
   make help
   ```

1. 在企业微信中创建机器人，在机器人的“webhook 地址”中获取 `key` 值，webhook 样式为：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key={key}`。
1. 修改 Alertmanager 配置文件：

   ```yaml
   route:
     receiver: "a2w"
   receivers:
     - name: "a2w"
       webhook_configs:
         - url: "http://{a2w_address}/send?key={key}"
   ```

## 消息模板

消息模板决定了企业微信机器人发出的消息格式，修改 `Makefile` 中的 `TEMPLATE` 变量的值来选择模板。

模板的使用注意事项请看同路径下的同名 Markdown 文件。

因为企业微信机器人接口限制单条消息的最大长度为 4096，所以本软件会对大于此限制的长消息进行分段。如果你使用自定义模板，请在想要分段的地方留一个空行（在企业微信中，至少三个连续的 `\n` 才被认为是一个空行），以便本软件对消息进行正确的分段。

## 构建项目

编译二进制文件：

```shell
make build
```

或是构建 Docker 镜像：

```shell
make docker-build
```

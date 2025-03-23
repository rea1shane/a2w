# Alertmanager to WeCom

[![docker-image-ci](https://github.com/rea1shane/a2w/actions/workflows/docker-image-ci.yml/badge.svg)](https://github.com/rea1shane/a2w/actions/workflows/docker-image-ci.yml)

通过 [企业微信机器人](https://developer.work.weixin.qq.com/document/path/91770) 发送 [Alertmanager](https://github.com/prometheus/alertmanager) 通知。查看 [演示](https://github.com/rea1shane/a2w/tree/main/demo)。

## 使用方式

1. 运行 A2W：

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

   或是通过 Helm 部署在 Kubernetes 中：

   ```shell
   make helm-install
   ```

   查看项目使用说明：

   ```shell
   make help
   ```

1. 在企业微信中创建机器人，在机器人的 “webhook 地址” 中获取 `key` 值，webhook 样式为：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key={key}`。
1. 修改 Alertmanager 配置文件：

   ```yaml
   route:
     receiver: "a2w"
   receivers:
     - name: "a2w"
       webhook_configs:
         - url: "http://{a2w_address}/send?key={key}"
   ```

## 功能说明

### 时区

A2W 使用本地时区。Docker 镜像的默认时区为 `Asia/Shanghai`。如果想要修改 Docker 镜像的时区，可以在 `docker run` 时使用参数 `-e TZ={TIME_ZONE}` 进行指定，如：

```shell
docker run --name a2w -d -p 5001:5001 -e TZ=Asia/Tokyo rea1shane/a2w
```

### 消息模板

消息模板决定了企业微信机器人发出的消息格式，在启动 A2W 时通过 `--template` 指定模板文件所在目录。默认模板的使用说明见 [文档](https://github.com/rea1shane/a2w/blob/main/templates/base.md)。

> [!NOTE]
>
> 因为企业微信机器人接口限制单条消息的最大长度为 4096，所以本软件会对大于此限制的长消息进行分段。如果你使用自定义模板，请在想要分段的地方留一个空行（在企业微信中，至少三个连续的 `\n` 才被认为是一个空行），以便本软件对消息进行正确的分段。

使用 tmpl URL Query 可以指定要使用的其他模板（以便适应想要让不同的告警使用不同的模板场景）：

```yaml
receivers:
  - name: "a2w"
    webhook_configs:
      - url: "http://{a2w_address}/send?tmpl=base_two&key={key}"
```

> [!Attention]
>
> 模板文件必须以 .tmpl 作为后缀，前缀要符合 URL 字符规范 [RFC 3986, Uniform Resource Identifier(URI): Generic Syntax](https://www.rfc-editor.org/rfc/rfc3986.html)。不要使用中文作为文件名。

### 用户提醒

A2W 支持用户提醒功能，修改 Alertmanager 中的配置如下：

```yaml
receivers:
  - name: "a2w"
    webhook_configs:
      - url: "http://{a2w_address}/send?key={key}&mention=user1&mention=user2"
```

当 A2W 接收到该 receiver 的通知时，会在发送至企业微信机器人的消息尾部添加一个 @ 列表，包含 url 中指定的所有用户。

## 构建项目

编译二进制文件：

```shell
make build
```

或是构建 Docker 镜像：

```shell
make docker-build
```

## 贡献

- [Bryan Chen](https://github.com/DesireWithin) 提出设置容器默认时区为东八时区。([#4](https://github.com/rea1shane/a2w/issues/4))
- [DesistDaydream](https://github.com/DesistDaydream) 提供用户提醒功能。([#7](https://github.com/rea1shane/a2w/pull/7))

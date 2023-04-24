# Alertmanager to WeCom

通过 [企业微信机器人](https://developer.work.weixin.qq.com/document/path/91770) 向企业微信发送 [Alertmanager](https://github.com/prometheus/alertmanager) 通知。

## 使用方式

1. 运行本项目：

   ```shell
   make run
   ```

   或是在 Docker 中部署：

   ```shell
   make docker-run
   ```

1. 在企业微信中创建机器人，在机器人的“webhook 地址”中获取 `key` 值，webhook 样式为：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key={key}`。
1. 修改 Alertmanager 配置文件：

   ```yaml
   route:
       receiver: 'a2w'
   receivers:
       - name: 'a2w'
           webhook_configs:
               - url: 'http://{a2w_host}:9099/send?key={key}'
   ```

## 消息模板

消息模板决定了企业微信机器人发出的消息格式。

### `base.tmpl`

若要使用该消息模板，在告警规则定义中必须包含：

- `labels.level`：告警规则等级。

可以选择包含：

- `annotations.current`：当前状态的表达式结果值，可以通过 `{{ $value }}` 获取。
- `annotations.labels`：可以定位到该告警实例的标签列表。

### `multiple-cluster.tmpl`

使用该模板与使用 [`base.tmpl`](#basetmpl) 相比，多了一个必选标签：

- `labels.cluster`：集群名称。

## 构建项目

编译二进制文件：

```shell
make build
```

或是构建 Docker 镜像：

```shell
make docker-build
```

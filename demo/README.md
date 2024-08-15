# Demo

本文件夹中的内容不属于 A2W 功能的一部分，只是用来测试或演示 A2W 的。该路径下的 Docker Compose 会启动 Prometheus、Alertmanager、A2W 三个服务。

在 Prometheus 中注册了三个不存在的 exporter，导致 Prometheus 产生告警，该告警最终会被转发到企业微信机器人中。 如果想要测试告警恢复效果，在告警产生后移除 `prometheus/targets.json` 中不存在的 exporter 即可。

## 使用方式

在启动 Docker Compose 之前，需要将 `alertmanager/alertmanager.yml` 中的 `<YOUR_KEY>` 替换成你的企业微信机器人的 key。然后执行：

```shell
docker compose up -d
```

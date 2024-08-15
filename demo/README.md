# Demo

本文件夹中的内容不属于 A2W 功能的一部分，只是用来测试或演示 A2W 的。该路径下的 Docker Compose 会启动 Prometheus、Alertmanager、A2W 三个服务。

## 使用方式

在启动 Docker Compose 之前，需要将 `alertmanager/alertmanager.yml` 中的 `<YOUR_KEY>` 替换成你的企业微信机器人的 key。然后执行：

```shell
docker compose up -d
```

调整 `prometheus/rules/rule.yml` 中告警的表达式，在企业微信中观察告警触发和恢复后的消息效果。

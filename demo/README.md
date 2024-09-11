# 演示

该路径下的 Docker Compose 会启动 Prometheus、Alertmanager、A2W 三个服务。

在启动 Docker Compose 之前，需要将 `alertmanager/alertmanager.yml` 中的 `<YOUR_KEY>` 替换成你的企业微信机器人的 key。然后执行：

```shell
docker compose up -d
```

调整 `prometheus/rules/rule.yml` 中告警的表达式，在企业微信中观察告警触发和恢复后的消息效果。

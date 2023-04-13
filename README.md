# Alertmanager to WeCom bot

## templates

### `message.tmpl`

在告警规则定义中，必须包含：

- `labels.level`：告警规则等级。
- `annotations.labels`：可以定位到该告警实例的标签列表。

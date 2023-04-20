# Alertmanager to WeCom bot

## 消息模板

### `base.tmpl`

在告警规则定义中，必须包含：

- `labels.level`：告警规则等级。
- `annotations.labels`：可以定位到该告警实例的标签列表。

可以选择包含：

- `annotations.current`：当前状态的表达式结果值，可以通过 `{{ $value }}` 获取。

# `base.tmpl`

Prometheus rule 中可以选择包含：

- `labels.level`：告警规则等级。
- `annotations.current`：当前状态的表达式结果值，可以通过 `{{ $value }}` 获取。
- `annotations.labels`：可以定位到该告警实例的标签列表。
- `annotations.mentions`：通知人员列表。

示例：

```yaml
groups:
  - name: Node
    rules:
      - alert: High node load
        expr: node_load15 / count without(cpu, mode) (node_cpu_seconds_total{mode="idle"}) * 100 > 100
        for: 5m
        labels:
          level: critical
        annotations:
          current: Average load_15 per core is {{ $value }}%
          labels: instance="{{ $labels.instance }}"

      - alert: High disk usage
        expr: (1 - node_filesystem_free_bytes / node_filesystem_size_bytes) * 100 > 85
        for: 15m
        annotations:
          current: "{{ $value }}% usage"
          labels: instance="{{ $labels.instance }}", device="{{ $labels.device }}", fstype="{{ $labels.fstype }}", mountpoint="{{ $labels.mountpoint }}"
          mentions: <@user1>, <@user2>
```

# 说明

本文件夹中的内容用于启动一个 Docker Compose 来测试 A2W，不是 A2W 程序的一部分。

## 使用方式

在启动 Docker Compose 之前，需要将 `alertmanager/alertmanager.yml` 中的 `<YOUR_KEY>` 替换成你的企业微信机器人的 key。然后执行：

```shell
docker compose up -d
```

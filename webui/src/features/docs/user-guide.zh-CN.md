---
title: Resin 使用说明
language: zh-CN
version: v1.3.0
auth: none
---

# Resin 使用说明

本文描述当前版本的可调用地址。不需要登录。机器可读副本：

- 中文：`/ui/user-guide.zh-CN.md`
- English: `/ui/user-guide.en.md`
- 索引：`/llms.txt`

## 版本

- 侧边栏版本在编译时写入。
- `docker compose up -d --build` 且未设置 `RESIN_VERSION` 时，版本取仓库根目录 `VERSION`，当前为 `v1.3.0`。
- 设置 `RESIN_VERSION` 可覆盖这一次构建。
- 发布包使用 Git 标签作为版本。

## 订阅

- 远程订阅字段 `user_agent`：留空使用 `clash.meta`。最长 256 字符，不能含控制字符。本地订阅不能设置。
- 字段 `probe_interval`：空、`0` 或 `0s` 使用全局探测间隔。否则为 Go duration，且 `>= 10s`。一节点属多个订阅时取最短正间隔。
- `POST /api/v1/subscriptions/batch`：正文 `text` 每行一个 `http` 或 `https` 链接。空行和 `#` 开头忽略。`name_regex` 有捕获组时用第一组，否则名称取主机名倒数第二段。重名追加 `-2`、`-3`。
- 纯文本代理行：`IP:PORT`、`IP:PORT:USER:PASS`。也支持 vmess、vless、trojan、ss、http、socks、sing-box、Clash。

## 平台接入

以下地址只包含该平台当前健康且已有出口 IP 的节点。把 `{host}`、`{proxy_token}`、`{platform}` 换成实际值。未启用代理令牌时，路径第一段可以是任意非空占位，例如 `public`。

### 正向代理

- 健康节点订阅，默认 base64 URI：

`GET http://{host}/{proxy_token}/api/v1/healthy-subscription?platform={platform}&format=uri`

- 直接接入仍走平台路由，只分配该平台可路由节点：

`http://{platform}.{account}:{proxy_token}@{host}`

`socks5h://{platform}.{account}:{proxy_token}@{host}`

### 反向代理

- 健康节点订阅，sing-box JSON：

`GET http://{host}/{proxy_token}/api/v1/healthy-subscription?platform={platform}&format=sing-box`

- 单次请求地址：

`http://{host}/{proxy_token}/{platform}.{account}/{protocol}/{target}`

## 平台规则

- `ipv4_only=true`：不分配 IPv6 出口、IPv6 服务器，以及 `domain_strategy=ipv6_only` 的节点。
- 节点名规则按行：普通行满足其一，`*` 必须匹配，`!` 排除。
- 租约节点失效后，优先换成同一国家的可路由节点；没有时再按分配策略选择。

## 熔断

- `POST /api/v1/nodes/{hash}/actions/open-circuit`
- `POST /{proxy_token}/api/v1/nodes/{hash}/actions/open-circuit`
- 成功响应：`{"status":"ok"}`

## 请求日志与登录

- `DELETE /api/v1/request-logs` 清除全部请求日志，包括尚未落盘的队列。不可撤销。需要管理员令牌。
- 管理接口返回 `401` 时，Web UI 删除本地管理员令牌并回到 `/ui/login`，`next` 保留原路径。

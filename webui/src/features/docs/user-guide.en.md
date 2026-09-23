---
title: Resin user guide
language: en
version: v1.3.0
auth: none
---

# Resin user guide

Callable URLs for the current version. No login is required. Machine-readable copies:

- 中文: `/ui/user-guide.zh-CN.md`
- English: `/ui/user-guide.en.md`
- Index: `/llms.txt`

## Version

- The sidebar version is compiled into the binary.
- `docker compose up -d --build` uses the repository `VERSION` file when `RESIN_VERSION` is unset. Current value: `v1.3.0`.
- Set `RESIN_VERSION` to override one build.
- Release artifacts use the Git tag.

## Subscriptions

- Remote field `user_agent`: blank uses `clash.meta`. Maximum 256 characters. No control characters. Local subscriptions cannot set it.
- Field `probe_interval`: blank, `0`, or `0s` uses the global interval. Otherwise a Go duration `>= 10s`. A node in several subscriptions uses the shortest positive interval.
- `POST /api/v1/subscriptions/batch`: `text` is one `http` or `https` URL per line. Blank lines and `#` lines are ignored. `name_regex` uses its first capture group; otherwise the name is the hostname label before the public suffix. Duplicates get `-2`, `-3`.
- Plain proxy lines: `IP:PORT` and `IP:PORT:USER:PASS`. Also vmess, vless, trojan, ss, http, socks, sing-box, and Clash.

## Platform access

These subscription URLs contain only nodes that are healthy for that platform and already have an egress IP. Replace `{host}`, `{proxy_token}`, and `{platform}`. If the proxy token is disabled, the first path segment may be any non-empty placeholder such as `public`.

### Forward proxy

- Healthy-node subscription, base64 URI list:

`GET http://{host}/{proxy_token}/api/v1/healthy-subscription?platform={platform}&format=uri`

- Direct access still uses platform routing and only routable nodes:

`http://{platform}.{account}:{proxy_token}@{host}`

`socks5h://{platform}.{account}:{proxy_token}@{host}`

### Reverse proxy

- Healthy-node subscription, sing-box JSON:

`GET http://{host}/{proxy_token}/api/v1/healthy-subscription?platform={platform}&format=sing-box`

- One request:

`http://{host}/{proxy_token}/{platform}.{account}/{protocol}/{target}`

## Platform rules

- `ipv4_only=true`: do not assign IPv6 exits, IPv6 servers, or `domain_strategy=ipv6_only` nodes.
- Node-name rules are line oriented: a plain line matches any, `*` is required, `!` excludes.
- After a leased node fails, replacement prefers the same country. Otherwise the allocation policy chooses.

## Circuit breaker

- `POST /api/v1/nodes/{hash}/actions/open-circuit`
- `POST /{proxy_token}/api/v1/nodes/{hash}/actions/open-circuit`
- Success body: `{"status":"ok"}`

## Request logs and login

- `DELETE /api/v1/request-logs` deletes every request log, including the unflushed queue. It cannot be undone. Admin token required.
- A management API `401` clears the stored admin token and returns the Web UI to `/ui/login`, preserving the previous path in `next`.

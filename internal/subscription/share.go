package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// EncodeSingBoxSubscription renders healthy nodes as a sing-box subscription document.
func EncodeSingBoxSubscription(raws []json.RawMessage) ([]byte, error) {
	outbounds := make([]json.RawMessage, 0, len(raws))
	for _, raw := range raws {
		if len(raw) == 0 {
			continue
		}
		outbounds = append(outbounds, raw)
	}
	return json.Marshal(map[string]any{"outbounds": outbounds})
}

// EncodeURISubscription renders importable proxy URIs, base64-encoded one line each.
func EncodeURISubscription(raws []json.RawMessage) ([]byte, error) {
	lines := make([]string, 0, len(raws))
	for _, raw := range raws {
		line, ok := outboundToURI(raw)
		if !ok {
			continue
		}
		lines = append(lines, line)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(strings.Join(lines, "\n")))
	return []byte(encoded), nil
}

func outboundToURI(raw json.RawMessage) (string, bool) {
	var outbound map[string]any
	if err := json.Unmarshal(raw, &outbound); err != nil {
		return "", false
	}
	kind, _ := outbound["type"].(string)
	server, _ := outbound["server"].(string)
	port := jsonPort(outbound["server_port"])
	tag, _ := outbound["tag"].(string)
	if server == "" || port <= 0 {
		return "", false
	}
	host := hostForURI(server)
	fragment := url.QueryEscape(tag)
	switch kind {
	case "shadowsocks":
		method, _ := outbound["method"].(string)
		password, _ := outbound["password"].(string)
		if method == "" || password == "" {
			return "", false
		}
		user := base64.RawURLEncoding.EncodeToString([]byte(method + ":" + password))
		return fmt.Sprintf("ss://%s@%s:%d#%s", user, host, port, fragment), true
	case "trojan":
		password, _ := outbound["password"].(string)
		if password == "" {
			return "", false
		}
		return fmt.Sprintf("trojan://%s@%s:%d#%s", url.QueryEscape(password), host, port, fragment), true
	case "vless":
		id, _ := outbound["uuid"].(string)
		if id == "" {
			return "", false
		}
		return fmt.Sprintf("vless://%s@%s:%d?encryption=none#%s", id, host, port, fragment), true
	case "vmess":
		id, _ := outbound["uuid"].(string)
		if id == "" {
			return "", false
		}
		payload, err := json.Marshal(map[string]any{
			"v": "2", "ps": tag, "add": server, "port": port, "id": id,
			"aid": 0, "net": "tcp", "type": "none", "tls": "",
		})
		if err != nil {
			return "", false
		}
		return "vmess://" + base64.StdEncoding.EncodeToString(payload), true
	case "http", "socks":
		username, _ := outbound["username"].(string)
		password, _ := outbound["password"].(string)
		scheme := kind
		if kind == "socks" {
			scheme = "socks5"
		}
		if username != "" || password != "" {
			return fmt.Sprintf("%s://%s:%s@%s:%d#%s", scheme, url.QueryEscape(username), url.QueryEscape(password), host, port, fragment), true
		}
		return fmt.Sprintf("%s://%s:%d#%s", scheme, host, port, fragment), true
	default:
		return "", false
	}
}

func hostForURI(server string) string {
	if strings.Contains(server, ":") && !strings.HasPrefix(server, "[") {
		return "[" + server + "]"
	}
	return server
}

func jsonPort(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case json.Number:
		parsed, err := typed.Int64()
		if err != nil {
			return 0
		}
		return int(parsed)
	default:
		return 0
	}
}

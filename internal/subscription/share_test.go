package subscription

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestEncodeURISubscription_Shadowsocks(t *testing.T) {
	raw := json.RawMessage(`{"type":"shadowsocks","tag":"hk","server":"1.2.3.4","server_port":443,"method":"aes-256-gcm","password":"secret"}`)
	body, err := EncodeURISubscription([]json.RawMessage{raw})
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(decoded), "ss://") || !strings.Contains(string(decoded), "1.2.3.4:443") {
		t.Fatalf("uri: %s", decoded)
	}
}

func TestEncodeSingBoxSubscription(t *testing.T) {
	raw := json.RawMessage(`{"type":"http","server":"1.2.3.4","server_port":8080}`)
	body, err := EncodeSingBoxSubscription([]json.RawMessage{raw})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"outbounds"`) {
		t.Fatalf("body: %s", body)
	}
}

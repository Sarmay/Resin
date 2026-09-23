package api

import (
	"net/http"
	"testing"
)

func TestAPIContract_SubscriptionLocalCreateValidation(t *testing.T) {
	srv, _, _ := newControlPlaneTestServer(t)

	rec := doJSONRequest(t, srv, http.MethodPost, "/api/v1/subscriptions", map[string]any{
		"name":                    "sub-local",
		"source_type":             "local",
		"incremental_alive_nodes": true,
	}, true)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("create local without content status: got %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	assertErrorCode(t, rec, "INVALID_ARGUMENT")

	rec = doJSONRequest(t, srv, http.MethodPost, "/api/v1/subscriptions", map[string]any{
		"name":        "sub-local",
		"source_type": "local",
		"content":     "vmess://example",
		"url":         "https://example.com/sub",
	}, true)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("create local with url status: got %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	assertErrorCode(t, rec, "INVALID_ARGUMENT")
}

func TestAPIContract_SubscriptionUserAgent(t *testing.T) {
	srv, _, _ := newControlPlaneTestServer(t)

	rec := doJSONRequest(t, srv, http.MethodPost, "/api/v1/subscriptions", map[string]any{
		"name":       "sub-remote",
		"url":        "https://example.com/sub",
		"user_agent": "FlClash/0.8.92",
	}, true)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create remote with user agent status: got %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	created := decodeJSONMap(t, rec)
	if got, _ := created["user_agent"].(string); got != "FlClash/0.8.92" {
		t.Fatalf("user_agent: got %q, want %q", got, "FlClash/0.8.92")
	}
	subID, _ := created["id"].(string)

	rec = doJSONRequest(t, srv, http.MethodPatch, "/api/v1/subscriptions/"+subID, map[string]any{
		"user_agent": "bad\r\nagent",
	}, true)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("patch invalid user agent status: got %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	assertErrorCode(t, rec, "INVALID_ARGUMENT")
}

func TestAPIContract_SubscriptionSourceTypeReadOnlyOnPatch(t *testing.T) {
	srv, _, _ := newControlPlaneTestServer(t)

	createRec := doJSONRequest(t, srv, http.MethodPost, "/api/v1/subscriptions", map[string]any{
		"name": "sub-remote",
		"url":  "https://example.com/sub",
	}, true)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create remote subscription status: got %d, want %d, body=%s", createRec.Code, http.StatusCreated, createRec.Body.String())
	}
	body := decodeJSONMap(t, createRec)
	subID, _ := body["id"].(string)
	if subID == "" {
		t.Fatalf("create remote subscription missing id: body=%s", createRec.Body.String())
	}

	rec := doJSONRequest(t, srv, http.MethodPatch, "/api/v1/subscriptions/"+subID, map[string]any{
		"source_type": "local",
	}, true)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("patch source_type status: got %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	assertErrorCode(t, rec, "INVALID_ARGUMENT")
}

package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/Resinat/Resin/internal/node"
	"github.com/Resinat/Resin/internal/subscription"
)

// RenderHealthySubscription returns a subscription document containing currently healthy nodes.
// platformName limits the result to that platform's routable view. Empty includes every healthy node.
func (s *ControlPlaneService) RenderHealthySubscription(format, platformName string) ([]byte, string, error) {
	if s == nil || s.Pool == nil {
		return nil, "", internal("healthy subscription", fmt.Errorf("pool is not configured"))
	}
	platformName = strings.TrimSpace(platformName)
	var view interface{ Contains(node.Hash) bool }
	if platformName != "" {
		plat, ok := s.Pool.GetPlatformByName(platformName)
		if !ok {
			return nil, "", notFound("platform not found")
		}
		view = plat.View()
	}
	raws := make([]json.RawMessage, 0)
	s.Pool.RangeNodes(func(hash node.Hash, entry *node.NodeEntry) bool {
		if view != nil && !view.Contains(hash) {
			return true
		}
		if entry == nil || !entry.IsHealthy() || !entry.GetEgressIP().IsValid() || len(entry.RawOptions) == 0 {
			return true
		}
		raws = append(raws, append(json.RawMessage(nil), entry.RawOptions...))
		return true
	})

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "uri":
		body, err := subscription.EncodeURISubscription(raws)
		if err != nil {
			return nil, "", internal("encode healthy subscription", err)
		}
		return body, "text/plain; charset=utf-8", nil
	case "sing-box", "singbox":
		body, err := subscription.EncodeSingBoxSubscription(raws)
		if err != nil {
			return nil, "", internal("encode healthy subscription", err)
		}
		return body, "application/json", nil
	default:
		return nil, "", invalidArg("format: must be uri or sing-box")
	}
}

// BatchSubscriptionError reports one rejected line from a batch import.
type BatchSubscriptionError struct {
	Line    int    `json:"line"`
	URL     string `json:"url"`
	Message string `json:"message"`
}

// BatchCreateSubscriptionsRequest imports many remote subscription URLs.
type BatchCreateSubscriptionsRequest struct {
	Text           string  `json:"text"`
	NameRegex      *string `json:"name_regex"`
	UpdateInterval *string `json:"update_interval"`
	UserAgent      *string `json:"user_agent"`
	ProbeInterval  *string `json:"probe_interval"`
	Enabled        *bool   `json:"enabled"`
	Ephemeral      *bool   `json:"ephemeral"`
}

// BatchCreateSubscriptionsResponse lists subscriptions created and lines rejected.
type BatchCreateSubscriptionsResponse struct {
	Created []SubscriptionResponse   `json:"created"`
	Errors  []BatchSubscriptionError `json:"errors"`
}

// CreateSubscriptionsBatch creates one remote subscription per non-empty line.
// Invalid lines are reported without failing the whole batch. Refresh stays on the scheduler.
func (s *ControlPlaneService) CreateSubscriptionsBatch(req BatchCreateSubscriptionsRequest) (*BatchCreateSubscriptionsResponse, error) {
	if strings.TrimSpace(req.Text) == "" {
		return nil, invalidArg("text is required")
	}
	var namePattern *regexp.Regexp
	if req.NameRegex != nil && strings.TrimSpace(*req.NameRegex) != "" {
		compiled, err := regexp.Compile(strings.TrimSpace(*req.NameRegex))
		if err != nil {
			return nil, invalidArg("name_regex: " + err.Error())
		}
		namePattern = compiled
	}

	used := map[string]struct{}{}
	if s.SubMgr != nil {
		s.SubMgr.Range(func(_ string, sub *subscription.Subscription) bool {
			used[sub.Name()] = struct{}{}
			return true
		})
	}

	resp := &BatchCreateSubscriptionsResponse{
		Created: []SubscriptionResponse{},
		Errors:  []BatchSubscriptionError{},
	}
	for index, line := range strings.Split(req.Text, "\n") {
		raw := strings.TrimSpace(line)
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}
		name, err := subscriptionNameForURL(raw, namePattern, used)
		if err != nil {
			resp.Errors = append(resp.Errors, BatchSubscriptionError{Line: index + 1, URL: raw, Message: err.Error()})
			continue
		}
		created, createErr := s.CreateSubscription(CreateSubscriptionRequest{
			Name:           &name,
			URL:            &raw,
			UpdateInterval: req.UpdateInterval,
			UserAgent:      req.UserAgent,
			ProbeInterval:  req.ProbeInterval,
			Enabled:        req.Enabled,
			Ephemeral:      req.Ephemeral,
		})
		if createErr != nil {
			delete(used, name)
			resp.Errors = append(resp.Errors, BatchSubscriptionError{Line: index + 1, URL: raw, Message: createErr.Error()})
			continue
		}
		resp.Created = append(resp.Created, *created)
	}
	return resp, nil
}

func subscriptionNameForURL(rawURL string, pattern *regexp.Regexp, used map[string]struct{}) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", fmt.Errorf("url must be an absolute http(s) URL")
	}
	base := ""
	if pattern != nil {
		match := pattern.FindStringSubmatch(rawURL)
		switch {
		case len(match) > 1 && strings.TrimSpace(match[1]) != "":
			base = strings.TrimSpace(match[1])
		case len(match) > 0 && strings.TrimSpace(match[0]) != "":
			base = strings.TrimSpace(match[0])
		}
	}
	if base == "" {
		host := parsed.Hostname()
		labels := strings.Split(host, ".")
		if len(labels) >= 2 && labels[len(labels)-2] != "" {
			base = labels[len(labels)-2]
		} else if host != "" {
			base = host
		} else {
			base = "sub"
		}
	}
	base = sanitizeSubscriptionName(base)
	name := base
	if _, exists := used[name]; exists {
		for i := 2; ; i++ {
			candidate := fmt.Sprintf("%s-%d", base, i)
			if _, taken := used[candidate]; !taken {
				name = candidate
				break
			}
		}
	}
	used[name] = struct{}{}
	return name, nil
}

func sanitizeSubscriptionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "sub"
	}
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}

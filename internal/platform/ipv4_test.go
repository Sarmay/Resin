package platform

import (
	"net/netip"
	"testing"

	"github.com/Resinat/Resin/internal/node"
)

func TestPlatform_EvaluateNode_IPv4OnlySkipsIPv6Egress(t *testing.T) {
	p := NewPlatform("p1", "Test", nil, nil)
	p.IPv4Only = true
	h := makeHash(`{"type":"ss","server":"example.com"}`)
	entry := makeFullyRoutableEntry(h, "sub1")
	entry.SetEgressIP(netip.MustParseAddr("2001:db8::1"))

	p.FullRebuild(func(fn func(node.Hash, *node.NodeEntry) bool) {
		fn(h, entry)
	}, alwaysLookup, usGeoLookup)

	if p.View().Size() != 0 {
		t.Fatal("ipv4-only platform should skip an IPv6 exit")
	}
}

func TestPlatform_EvaluateNode_IPv4OnlyKeepsIPv4Egress(t *testing.T) {
	p := NewPlatform("p1", "Test", nil, nil)
	p.IPv4Only = true
	h := makeHash(`{"type":"ss","server":"example.com"}`)
	entry := makeFullyRoutableEntry(h, "sub1")

	p.FullRebuild(func(fn func(node.Hash, *node.NodeEntry) bool) {
		fn(h, entry)
	}, alwaysLookup, usGeoLookup)

	if p.View().Size() != 1 {
		t.Fatal("ipv4-only platform should keep an IPv4 exit")
	}
}

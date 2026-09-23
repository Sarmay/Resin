package routing

import (
	"testing"
	"time"

	"github.com/Resinat/Resin/internal/model"
	"github.com/Resinat/Resin/internal/platform"
)

func TestRouteSticky_ReplacementPrefersPreviousCountry(t *testing.T) {
	pool := newRouterTestPool()
	plat := platform.NewPlatform("plat", "plat", nil, nil)
	plat.StickyTTLNs = int64(time.Hour)
	pool.addPlatform(plat)

	us1Hash, us1 := newRoutableEntry(t, `{"type":"ss","server":"1.1.1.1"}`, "1.1.1.1")
	us1.SetEgressRegion("us")
	us2Hash, us2 := newRoutableEntry(t, `{"type":"ss","server":"2.2.2.2"}`, "2.2.2.2")
	us2.SetEgressRegion("us")
	_, jp := newRoutableEntry(t, `{"type":"ss","server":"3.3.3.3"}`, "3.3.3.3")
	jp.SetEgressRegion("jp")
	jpHash := jp.Hash
	pool.addEntry(us1Hash, us1)
	pool.addEntry(us2Hash, us2)
	pool.addEntry(jpHash, jp)
	pool.rebuildPlatformView(plat)

	router := newTestRouter(pool, nil)
	state := router.ensurePlatformState(plat.ID)
	now := time.Now()
	if err := router.UpsertLease(model.Lease{
		PlatformID:     plat.ID,
		Account:        "acct",
		NodeHash:       us1Hash.Hex(),
		EgressIP:       "1.1.1.1",
		CreatedAtNs:    now.UnixNano(),
		ExpiryNs:       now.Add(time.Hour).UnixNano(),
		LastAccessedNs: now.UnixNano(),
	}); err != nil {
		t.Fatal(err)
	}
	us1.CircuitOpenSince.Store(time.Now().UnixNano())
	pool.rebuildPlatformView(plat)

	result, err := router.routeSticky(plat, state, "acct", "example.com", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if result.NodeHash != us2Hash {
		t.Fatalf("replacement hash = %s, want same-country %s", result.NodeHash.Hex(), us2Hash.Hex())
	}
}

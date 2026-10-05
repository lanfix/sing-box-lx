package rule

import (
	"context"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
)

// lx: SPEC 113 — `source_mac_address` and `source_hostname` in headless rules,
// so a remote or local rule-set can change the device list without reloading
// the whole config. The neighbor resolver is created before rule-sets are
// loaded and a rule-set cannot start it, so such a rule requires the router to
// look up neighbors already (`route.find_neighbor` or a neighbor route/DNS
// rule); otherwise the rule-set is rejected instead of silently never matching.
func appendNeighborHeadlessItems(ctx context.Context, rule *DefaultHeadlessRule, options option.DefaultHeadlessRule) error {
	if len(options.SourceMACAddress) == 0 && len(options.SourceHostname) == 0 {
		return nil
	}
	router := service.FromContext[adapter.Router](ctx)
	if router != nil && !router.NeedFindNeighbor() {
		return E.New("`source_mac_address` and `source_hostname` rule items in rule-set require `route.find_neighbor`")
	}
	if len(options.SourceMACAddress) > 0 {
		item := NewSourceMACAddressItem(options.SourceMACAddress)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	if len(options.SourceHostname) > 0 {
		item := NewSourceHostnameItem(options.SourceHostname)
		rule.items = append(rule.items, item)
		rule.allItems = append(rule.allItems, item)
	}
	return nil
}

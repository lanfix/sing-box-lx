package rule

import (
	"context"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
)

// lx: SPEC 113 — neighbor items in headless rules. The resolver starts before
// rule-sets load, so without `route.find_neighbor` the rule-set is rejected.
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

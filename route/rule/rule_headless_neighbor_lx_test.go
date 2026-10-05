package rule

import (
	"context"
	"net"
	"net/netip"
	"testing"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/service"

	"github.com/stretchr/testify/require"
)

type neighborTestRouter struct {
	adapter.Router
	needFindNeighbor bool
}

func (r *neighborTestRouter) NeedFindNeighbor() bool {
	return r.needFindNeighbor
}

func neighborTestContext(needFindNeighbor bool) context.Context {
	return service.ContextWith[adapter.Router](context.Background(), &neighborTestRouter{needFindNeighbor: needFindNeighbor})
}

func neighborTestRuleSet(t *testing.T, ctx context.Context, config string) (adapter.RuleSet, error) {
	t.Helper()
	var plainOptions option.PlainRuleSetCompat
	err := json.UnmarshalContext(ctx, []byte(config), &plainOptions)
	require.NoError(t, err)
	return NewLocalRuleSet(ctx, log.NewNOPFactory().Logger(), "devices", option.RuleSet{
		Type:          C.RuleSetTypeInline,
		InlineOptions: plainOptions.Options,
	})
}

func neighborTestMetadata(t *testing.T, source string, mac string, hostname string) *adapter.InboundContext {
	t.Helper()
	metadata := &adapter.InboundContext{
		Source:         M.SocksaddrFrom(netip.MustParseAddr(source), 50000),
		SourceHostname: hostname,
	}
	if mac != "" {
		hardwareAddr, err := net.ParseMAC(mac)
		require.NoError(t, err)
		metadata.SourceMACAddress = hardwareAddr
	}
	return metadata
}

func TestNeighborHeadlessRuleSourceFormat(t *testing.T) {
	t.Parallel()
	ruleSet, err := json.UnmarshalExtended[option.PlainRuleSetCompat]([]byte(`{
		"version": 5,
		"rules": [{"source_mac_address": ["AA:BB:CC:DD:EE:01", "aa-bb-cc-dd-ee-02"], "source_hostname": "tv"}]
	}`))
	require.NoError(t, err)
	plainRuleSet, err := ruleSet.Upgrade()
	require.NoError(t, err)
	require.Len(t, plainRuleSet.Rules, 1)
	defaultRule := plainRuleSet.Rules[0].DefaultOptions
	require.Equal(t, []string{"AA:BB:CC:DD:EE:01", "aa-bb-cc-dd-ee-02"}, []string(defaultRule.SourceMACAddress))
	require.Equal(t, []string{"tv"}, []string(defaultRule.SourceHostname))
	require.True(t, defaultRule.IsValid())
}

func TestNeighborHeadlessRuleMatchMACAddress(t *testing.T) {
	t.Parallel()
	ruleSet, err := neighborTestRuleSet(t, neighborTestContext(true), `{
		"version": 5,
		"rules": [{"source_mac_address": ["AA:BB:CC:DD:EE:01", "aa-bb-cc-dd-ee-02"]}]
	}`)
	require.NoError(t, err)
	require.True(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.10", "aa:bb:cc:dd:ee:01", "")))
	require.True(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.11", "AA:BB:CC:DD:EE:02", "")))
	require.False(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.12", "aa:bb:cc:dd:ee:03", "")))
	require.False(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.13", "", "")))
}

func TestNeighborHeadlessRuleMatchHostname(t *testing.T) {
	t.Parallel()
	ruleSet, err := neighborTestRuleSet(t, neighborTestContext(true), `{
		"version": 5,
		"rules": [{"source_hostname": ["tv", "laptop"]}]
	}`)
	require.NoError(t, err)
	require.True(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.10", "", "tv")))
	require.False(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.11", "", "phone")))
	require.False(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.12", "", "")))
}

func TestNeighborHeadlessRuleCombinesWithOtherItems(t *testing.T) {
	t.Parallel()
	ruleSet, err := neighborTestRuleSet(t, neighborTestContext(true), `{
		"version": 5,
		"rules": [
			{"source_mac_address": "aa:bb:cc:dd:ee:01", "source_ip_cidr": "192.168.1.0/24"},
			{"type": "logical", "mode": "and", "rules": [
				{"source_mac_address": "aa:bb:cc:dd:ee:01", "invert": true},
				{"source_ip_cidr": "10.0.0.0/8"}
			]}
		]
	}`)
	require.NoError(t, err)
	require.True(t, ruleSet.Match(neighborTestMetadata(t, "192.168.1.10", "aa:bb:cc:dd:ee:01", "")))
	require.False(t, ruleSet.Match(neighborTestMetadata(t, "192.168.2.10", "aa:bb:cc:dd:ee:01", "")))
	require.True(t, ruleSet.Match(neighborTestMetadata(t, "10.0.0.5", "aa:bb:cc:dd:ee:02", "")))
	require.False(t, ruleSet.Match(neighborTestMetadata(t, "10.0.0.5", "aa:bb:cc:dd:ee:01", "")))
}

func TestNeighborHeadlessRuleRequiresFindNeighbor(t *testing.T) {
	t.Parallel()
	for _, config := range []string{
		`{"version": 5, "rules": [{"source_mac_address": "aa:bb:cc:dd:ee:01"}]}`,
		`{"version": 5, "rules": [{"source_hostname": "tv"}]}`,
		`{"version": 5, "rules": [{"type": "logical", "mode": "or", "rules": [
			{"domain": "example.com"},
			{"source_mac_address": "aa:bb:cc:dd:ee:01"}
		]}]}`,
	} {
		_, err := neighborTestRuleSet(t, neighborTestContext(false), config)
		require.ErrorContains(t, err, "route.find_neighbor", config)
	}
	_, err := neighborTestRuleSet(t, neighborTestContext(false), `{"version": 5, "rules": [{"domain": "example.com"}]}`)
	require.NoError(t, err)
}

func TestNeighborHeadlessRuleWithoutRouter(t *testing.T) {
	t.Parallel()
	rule, err := NewDefaultHeadlessRule(context.Background(), option.DefaultHeadlessRule{
		SourceMACAddress: []string{"aa:bb:cc:dd:ee:01"},
	})
	require.NoError(t, err)
	require.True(t, rule.Match(neighborTestMetadata(t, "192.168.1.10", "aa:bb:cc:dd:ee:01", "")))
}

package srs

import (
	"bytes"
	"testing"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"

	"github.com/stretchr/testify/require"
)

// lx: SPEC 113 — neighbor items have no binary encoding: compiling them must
// fail instead of silently dropping the condition and widening the rule.
func TestWriteRejectsNeighborItems(t *testing.T) {
	t.Parallel()
	for _, rule := range []option.DefaultHeadlessRule{
		{SourceMACAddress: []string{"aa:bb:cc:dd:ee:01"}},
		{SourceHostname: []string{"tv"}},
	} {
		err := Write(&bytes.Buffer{}, option.PlainRuleSet{
			Rules: []option.HeadlessRule{{Type: C.RuleTypeDefault, DefaultOptions: rule}},
		}, C.RuleSetVersionCurrent)
		require.ErrorContains(t, err, "only supported in source rule-set")
		err = Write(&bytes.Buffer{}, option.PlainRuleSet{
			Rules: []option.HeadlessRule{{Type: C.RuleTypeLogical, LogicalOptions: option.LogicalHeadlessRule{
				Mode:  C.LogicalTypeOr,
				Rules: []option.HeadlessRule{{Type: C.RuleTypeDefault, DefaultOptions: rule}},
			}}},
		}, C.RuleSetVersionCurrent)
		require.ErrorContains(t, err, "only supported in source rule-set")
	}
	err := Write(&bytes.Buffer{}, option.PlainRuleSet{
		Rules: []option.HeadlessRule{{Type: C.RuleTypeDefault, DefaultOptions: option.DefaultHeadlessRule{
			Domain: []string{"example.com"},
		}}},
	}, C.RuleSetVersionCurrent)
	require.NoError(t, err)
}

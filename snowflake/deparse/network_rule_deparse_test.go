package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// NETWORK RULE DDL round-trips (gap-network-rule)
//
// assertRoundTrip compares ASTs (Loc-stripped), so cosmetic normalizations the
// deparser applies — uppercasing option names/keys, normalizing whitespace —
// are expected and verified to be AST-equivalent.
// ---------------------------------------------------------------------------

func TestDeparse_CreateNetworkRule(t *testing.T) {
	assertRoundTrip(t, "CREATE NETWORK RULE corporate_network TYPE = AWSVPCEID VALUE_LIST = ('vpce-123abc3420c1931') MODE = INTERNAL_STAGE COMMENT = 'corporate privatelink endpoint'")
	assertRoundTrip(t, "CREATE NETWORK RULE cloud_network TYPE = IPV4 VALUE_LIST = ('47.88.25.32/27') COMMENT = 'cloud egress ip range'")
	assertRoundTrip(t, "CREATE NETWORK RULE gcp_rule TYPE = GCPPSCID MODE = INGRESS VALUE_LIST = ('31618973889077266')")
	assertRoundTrip(t, "CREATE NETWORK RULE external_access_rule TYPE = HOST_PORT MODE = EGRESS VALUE_LIST = ('example.com', 'example.com:443')")
	assertRoundTrip(t, "CREATE OR REPLACE NETWORK RULE ext_db.network_rules.azure_rule MODE = EGRESS TYPE = PRIVATE_HOST_PORT VALUE_LIST = ('externalaccessdemo.database.windows.net')")
	assertRoundTrip(t, "CREATE NETWORK RULE IF NOT EXISTS r TYPE = IPV4 VALUE_LIST = ('0.0.0.0/0')")
}

func TestDeparse_AlterNetworkRule(t *testing.T) {
	assertRoundTrip(t, "ALTER NETWORK RULE r SET VALUE_LIST = ('1.2.3.4', '5.6.7.8')")
	assertRoundTrip(t, "ALTER NETWORK RULE r SET VALUE_LIST = ('1.2.3.4') COMMENT = 'updated'")
	assertRoundTrip(t, "ALTER NETWORK RULE IF EXISTS r SET COMMENT = 'just a comment'")
	assertRoundTrip(t, "ALTER NETWORK RULE r UNSET VALUE_LIST")
	assertRoundTrip(t, "ALTER NETWORK RULE r UNSET VALUE_LIST, COMMENT")
}

func TestDeparse_DropNetworkRule(t *testing.T) {
	assertRoundTrip(t, "DROP NETWORK RULE r")
	assertRoundTrip(t, "DROP NETWORK RULE IF EXISTS db.sch.r")
}

// TestDeparse_NetworkRule_ExactString pins the exact deparsed output for a
// representative CREATE and ALTER, guarding the keyword spelling and spacing.
func TestDeparse_NetworkRule_ExactString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{
			in:   "CREATE OR REPLACE NETWORK RULE r TYPE = ipv4 VALUE_LIST = ('1.2.3.4')",
			want: "CREATE OR REPLACE NETWORK RULE r TYPE = IPV4 VALUE_LIST = ('1.2.3.4')",
		},
		{
			in:   "ALTER NETWORK RULE r UNSET value_list, comment",
			want: "ALTER NETWORK RULE r UNSET VALUE_LIST, COMMENT",
		},
		{
			in:   "DROP NETWORK RULE IF EXISTS r",
			want: "DROP NETWORK RULE IF EXISTS r",
		},
	}
	for _, c := range cases {
		f, err := parser.Parse(c.in)
		require.NoError(t, err, "parse %q", c.in)
		out, err := deparse.DeparseFile(f)
		require.NoError(t, err, "deparse %q", c.in)
		require.Equal(t, c.want, out)
	}
}

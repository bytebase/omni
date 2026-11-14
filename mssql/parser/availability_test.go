package parser

import "testing"

// TestParseAvailabilityGroupRoutingURL: a routing URL is a 'TCP://host:port'
// string. SQL Server 2022 refuses NONE there (error 153), though the
// documentation for later versions shows it.
func TestParseAvailabilityGroupRoutingURL(t *testing.T) {
	for _, tt := range []struct {
		sql string
		ok  bool
	}{
		{"ALTER AVAILABILITY GROUP a MODIFY REPLICA ON 'a' WITH (SECONDARY_ROLE (READ_ONLY_ROUTING_URL = 'TCP://a:1'))", true},
		{"ALTER AVAILABILITY GROUP a MODIFY REPLICA ON 'a' WITH (PRIMARY_ROLE (READ_WRITE_ROUTING_URL = 'TCP://a:1'))", true},
		{"ALTER AVAILABILITY GROUP a MODIFY REPLICA ON 'a' WITH (PRIMARY_ROLE (READ_ONLY_ROUTING_LIST = NONE))", true},
		{"ALTER AVAILABILITY GROUP a MODIFY REPLICA ON 'a' WITH (SECONDARY_ROLE (READ_ONLY_ROUTING_URL = NONE))", false},
		{"ALTER AVAILABILITY GROUP a MODIFY REPLICA ON 'a' WITH (PRIMARY_ROLE (READ_WRITE_ROUTING_URL = NONE))", false},
		{"ALTER AVAILABILITY GROUP a MODIFY REPLICA ON 'a' WITH (SECONDARY_ROLE (READ_ONLY_ROUTING_URL = 'a'))", false},
	} {
		_, err := Parse(tt.sql)
		if (err == nil) != tt.ok {
			t.Errorf("Parse(%q): err = %v, want ok = %v", tt.sql, err, tt.ok)
		}
	}
}

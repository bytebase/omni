package driver

import (
	"context"
	"errors"
	"maps"
	"net"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	omnidriver "github.com/bytebase/omni/driver"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		banner string
		first  int
		second int
	}{
		{banner: "12.1.0.2.0", first: 12, second: 1},
		{banner: "12.1.0.", first: 12, second: 1},
		{banner: "23.26.3.0.0", first: 23, second: 26},
	}
	for _, test := range tests {
		version, err := parseVersion(test.banner)
		if err != nil {
			t.Fatalf("%s: %v", test.banner, err)
		}
		if version.First != test.first || version.Second != test.second {
			t.Errorf("%s parsed as %+v, want %d.%d", test.banner, version, test.first, test.second)
		}
	}
	if _, err := parseVersion("Oracle Database 23ai"); err == nil {
		t.Error("a version string without a release parsed")
	}
}

func TestConnectionOptions(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
		useTLS bool
		want   map[string]string
	}{
		{
			name: "no parameters",
			want: map[string]string{"CONNECTION TIMEOUT": "10"},
		},
		{
			name:   "the service stays out of the options",
			params: map[string]string{"SERVICE_NAME": "ORCLPDB1"},
			want:   map[string]string{"CONNECTION TIMEOUT": "10"},
		},
		{
			name:   "a SID passes through in its own spelling",
			params: map[string]string{"Sid": "ORCL"},
			want:   map[string]string{"CONNECTION TIMEOUT": "10", "Sid": "ORCL"},
		},
		{
			name:   "dial timeout in the exact spelling",
			params: map[string]string{"CONNECTION TIMEOUT": "5"},
			want:   map[string]string{"CONNECTION TIMEOUT": "5"},
		},
		{
			name:   "dial timeout in lowercase",
			params: map[string]string{"connection timeout": "5"},
			want:   map[string]string{"connection timeout": "5"},
		},
		{
			name:   "dial timeout alias",
			params: map[string]string{"CONNECT TIMEOUT": "5"},
			want:   map[string]string{"CONNECT TIMEOUT": "5"},
		},
		{
			name:   "dial timeout turned off",
			params: map[string]string{"CONNECTION TIMEOUT": "0"},
			want:   map[string]string{"CONNECTION TIMEOUT": "0"},
		},
		{
			name:   "read timeout is a different option",
			params: map[string]string{"TIMEOUT": "30"},
			want:   map[string]string{"TIMEOUT": "30", "CONNECTION TIMEOUT": "10"},
		},
		{
			name: "other options pass through in their own spelling",
			params: map[string]string{
				"SERVER":     "standby-1.example.com:1521, standby-2.example.com:1521",
				"WALLET":     "/opt/oracle/wallet",
				"AUTH TYPE":  "TCPS",
				"ssl verify": "false",
			},
			want: map[string]string{
				"CONNECTION TIMEOUT": "10",
				"SERVER":             "standby-1.example.com:1521, standby-2.example.com:1521",
				"WALLET":             "/opt/oracle/wallet",
				"AUTH TYPE":          "TCPS",
				"ssl verify":         "false",
			},
		},
		{
			name:   "TLS sets SSL",
			useTLS: true,
			want:   map[string]string{"CONNECTION TIMEOUT": "10", "SSL": "TRUE"},
		},
		{
			name:   "TLS replaces an SSL option in any spelling",
			params: map[string]string{"ssl": "FALSE"},
			useTLS: true,
			want:   map[string]string{"CONNECTION TIMEOUT": "10", "SSL": "TRUE"},
		},
		{
			name:   "TLS passes the other TLS options through",
			params: map[string]string{"SSL VERIFY": "FALSE", "WALLET": "/opt/oracle/wallet"},
			useTLS: true,
			want: map[string]string{
				"CONNECTION TIMEOUT": "10",
				"SSL":                "TRUE",
				"SSL VERIFY":         "FALSE",
				"WALLET":             "/opt/oracle/wallet",
			},
		},
		{
			name:   "without TLS an SSL option passes through",
			params: map[string]string{"ssl": "TRUE"},
			want:   map[string]string{"CONNECTION TIMEOUT": "10", "ssl": "TRUE"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			params := maps.Clone(test.params)
			got := connectionOptions(test.params, test.useTLS)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("options (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(params, test.params); diff != "" {
				t.Errorf("connectionOptions changed its input (-before +after):\n%s", diff)
			}
		})
	}
}

// TestOpenDialsThroughConfig pins that a caller's Dial reaches the client
// instead of a direct connection: the test's dialer fails, and so does the
// connection.
func TestOpenDialsThroughConfig(t *testing.T) {
	dialed := false
	_, err := Open(context.Background(), omnidriver.Config{
		Host:     "primary.example.com",
		Port:     "1521",
		User:     "scott",
		Password: "tiger",
		Database: "APP",
		Params:   map[string]string{"SERVICE_NAME": "ORCLPDB1"},
		Dial: func(context.Context, string, string) (net.Conn, error) {
			dialed = true
			return nil, errors.New("the test's dialer refuses")
		},
	})
	if !dialed {
		t.Fatalf("the driver dialed directly: %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), "the test's dialer refuses") {
		t.Errorf("error = %v, want the dialer's", err)
	}
}

func TestServiceName(t *testing.T) {
	for _, test := range []struct {
		name   string
		params map[string]string
		want   string
	}{
		{name: "named", params: map[string]string{"SERVICE_NAME": "ORCLPDB1"}, want: "ORCLPDB1"},
		{name: "named in another spelling", params: map[string]string{"service_name": "ORCLPDB1"}, want: "ORCLPDB1"},
		{
			// go-ora reads the SID in place of the service.
			name:   "a SID stands in for the service",
			params: map[string]string{"SID": "ORCL"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := serviceName(test.params)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("service = %q, want %q", got, test.want)
			}
		})
	}
	if _, err := serviceName(nil); err == nil {
		t.Error("an endpoint with no service and no SID was accepted")
	}
	if _, err := serviceName(map[string]string{"SERVICE_NAME": " "}); err == nil {
		t.Error("an empty service was accepted")
	}
}

func TestIsInsertFirst(t *testing.T) {
	for _, test := range []struct {
		statement string
		want      bool
	}{
		{statement: "INSERT FIRST WHEN id < 10 THEN INTO t1 VALUES (id) SELECT id FROM t", want: true},
		{statement: "insert first when id < 10 then into t1 values (id) select id from t", want: true},
		{statement: "/* a script's own note */ INSERT FIRST WHEN 1 = 1 THEN INTO t VALUES (1) SELECT 1 FROM DUAL", want: true},
		{statement: "INSERT ALL INTO t1 VALUES (id) INTO t2 VALUES (id) SELECT id FROM t"},
		{statement: "INSERT INTO t VALUES (1)"},
		{statement: "UPDATE t SET a = 1"},
	} {
		if got := isInsertFirst(test.statement); got != test.want {
			t.Errorf("isInsertFirst(%q) = %v, want %v", test.statement, got, test.want)
		}
	}
}

func TestWithDynamicSamplingHint(t *testing.T) {
	const hint = "OPT_PARAM('optimizer_dynamic_sampling' 11)"
	for _, test := range []struct {
		statement string
		want      string
	}{
		{statement: "UPDATE t SET a = 1", want: "UPDATE /*+ " + hint + " */ t SET a = 1"},
		{statement: "  delete from t", want: "  delete /*+ " + hint + " */ from t"},
		{
			// Only a hint block joins the hint: a plain comment is not one, so
			// the hint goes before it.
			statement: "INSERT /* move */ FIRST WHEN id < 10 THEN INTO t1 VALUES (id) SELECT id FROM t",
			want:      "INSERT /*+ " + hint + " */ /* move */ FIRST WHEN id < 10 THEN INTO t1 VALUES (id) SELECT id FROM t",
		},
		{
			statement: "MERGE /*+ PARALLEL(4) */ INTO t2 USING t ON (t2.id = t.id) WHEN MATCHED THEN UPDATE SET t2.v = t.v",
			want:      "MERGE /*+ " + hint + "  PARALLEL(4) */ INTO t2 USING t ON (t2.id = t.id) WHEN MATCHED THEN UPDATE SET t2.v = t.v",
		},
		{
			// Other statements are planned as they are.
			statement: "SELECT * FROM t",
			want:      "SELECT * FROM t",
		},
	} {
		if got := withDynamicSamplingHint(test.statement); got != test.want {
			t.Errorf("withDynamicSamplingHint(%q) = %q, want %q", test.statement, got, test.want)
		}
	}
}

func TestSkipSpaceAndComments(t *testing.T) {
	for _, test := range []struct {
		statement string
		want      int
	}{
		{statement: "UPDATE t SET a = 1", want: 0},
		{statement: " \t\nUPDATE t", want: 3},
		{statement: "-- a note\nUPDATE t", want: 10},
		{statement: "/* a note */ UPDATE t", want: 13},
		{statement: "-- an unterminated note", want: 23},
		{statement: "/* an unterminated note", want: 23},
	} {
		if got := skipSpaceAndComments(test.statement); got != test.want {
			t.Errorf("skipSpaceAndComments(%q) = %d, want %d", test.statement, got, test.want)
		}
	}
}

func TestPlanJSON(t *testing.T) {
	columns := []string{"ID", "OPERATION", "OBJECT_NAME", "FILTER_PREDICATES"}
	rows := [][]any{
		{[]byte("0"), []byte("SELECT STATEMENT"), nil, nil},
		{[]byte("1"), []byte("TABLE ACCESS"), []byte("T"), []byte("a\n= 1")},
	}
	got, err := planJSON(columns, rows)
	if err != nil {
		t.Fatal(err)
	}
	const want = `[{"id":"0","operation":"SELECT STATEMENT"},{"filter_predicates":"a\n= 1","id":"1","object_name":"T","operation":"TABLE ACCESS"}]`
	if got != want {
		t.Errorf("plan = %s, want %s", got, want)
	}
	// A row shorter than the columns names no value for the rest.
	if _, err := planJSON(columns, [][]any{{[]byte("0")}}); err != nil {
		t.Errorf("a ragged row failed: %v", err)
	}
}

package driver

import (
	"context"
	"crypto/tls"
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

func baseConfig() driver.Config {
	return driver.Config{User: "root", Host: "db.example.com", Port: "9030"}
}

func TestConnConfig(t *testing.T) {
	mc, err := connConfig(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	if !mc.MultiStatements || !mc.InterpolateParams || mc.MaxAllowedPacket != 0 || mc.Net != "tcp" || mc.Addr != "db.example.com:9030" {
		t.Fatalf("unexpected defaults: %+v", mc)
	}

	// Special characters in credentials and Unix socket path never inject DSN parameters.
	unixCfg := driver.Config{
		User:     "app?allowAllFiles=true",
		Password: "p@ss/word?tls=false&allowAllFiles=true",
		Host:     "/var/run/mysqld/mysqld.sock",
		Database: "sales",
		Params:   map[string]string{"timeout": "10s", "allowCleartextPasswords": "true"},
	}
	mcUnix, err := connConfig(unixCfg)
	if err != nil {
		t.Fatal(err)
	}
	if mcUnix.Net != "unix" || mcUnix.User != unixCfg.User || mcUnix.Passwd != unixCfg.Password || mcUnix.AllowAllFiles {
		t.Fatalf("unexpected unix config: %+v", mcUnix)
	}

	for _, key := range []string{"allowAllFiles", "allowallfiles", "AllowAllFiles", " allowAllFiles "} {
		bad := baseConfig()
		bad.Params = map[string]string{key: "true"}
		if _, err := connConfig(bad); err == nil {
			t.Errorf("accepted %q", key)
		}
	}

	// TLS cannot be downgraded by Params; nil cfg.TLS with tls=true infers host ServerName.
	tlsCfg := baseConfig()
	tlsCfg.Params = map[string]string{"tls": "false", "allowFallbackToPlaintext": "true"}
	tlsCfg.TLS = func(host string) (*tls.Config, error) { return &tls.Config{ServerName: host}, nil }
	mcTLS, err := connConfig(tlsCfg)
	if err != nil || mcTLS.TLS == nil || mcTLS.TLS.ServerName != "db.example.com" || mcTLS.AllowFallbackToPlaintext {
		t.Fatalf("TLS config = %+v, err = %v", mcTLS, err)
	}

	paramTLS := baseConfig()
	paramTLS.Params = map[string]string{"tls": "true"}
	mcParam, err := connConfig(paramTLS)
	if err != nil || mcParam.TLS == nil || mcParam.TLS.ServerName != "db.example.com" {
		t.Fatalf("Params tls=true = %+v, err = %v", mcParam, err)
	}

	// Open does not reach the server.
	dead := baseConfig()
	dead.Host, dead.Port = "127.0.0.1", "1"
	d, err := Open(context.Background(), dead)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Ping(context.Background()); err == nil {
		t.Error("Ping succeeded on closed port")
	}
	_ = d.Close()
}

// sanitizeCommentToken keeps an application name inside the app comment it is
// written into, even when removing its terminator would join another.
func TestSanitizeCommentToken(t *testing.T) {
	for _, name := range []string{"**//", "*/ SELECT 2; -- ", "a*/b*/c"} {
		if got := sanitizeCommentToken(name); strings.Contains(got, "*/") {
			t.Errorf("sanitizeCommentToken(%q) = %q still ends the comment", name, got)
		}
	}
	if got := sanitizeCommentToken("bytebase"); got != "bytebase" {
		t.Errorf("sanitizeCommentToken(bytebase) = %q", got)
	}
}

// A DELIMITER script is cut by its own terminator, and the directives are not
// sent: the server knows neither.
func TestSplitExecuteCommandsDelimiter(t *testing.T) {
	tests := []struct {
		script string
		want   []string
	}{
		{
			script: "DELIMITER $$\nCREATE TABLE t (a INT)$$\nDELIMITER ;\nINSERT INTO t VALUES (1);",
			want:   []string{"CREATE TABLE t (a INT)", "INSERT INTO t VALUES (1);"},
		},
		{
			script: "  delimiter //\nSELECT 1 // \nSELECT 2 //",
			want:   []string{"SELECT 1", "SELECT 2"},
		},
		{
			script: "DELIMITER ;\nSELECT 'a;b';\nSELECT 2;",
			want:   []string{"SELECT 'a;b';", "SELECT 2;"},
		},
		{
			// "1--2" is a subtraction, not a comment: the engines' lexers require
			// a space after "--".
			script: "DELIMITER $$\nSELECT 1--2$$\nSELECT 3$$",
			want:   []string{"SELECT 1--2", "SELECT 3"},
		},
		{
			// A comment before a directive belongs to no statement, so it does
			// not hide the directive that follows it.
			script: "-- migration header\nDELIMITER $$\nSELECT 1$$ -- first query\nDELIMITER ;\nSELECT 2;",
			want:   []string{"SELECT 1", "SELECT 2;"},
		},
		{
			// '#' and '//' open no line comment, so a directive's terminator on
			// that line still ends the statement: the text keeps both markers.
			script: "DELIMITER $$\nSELECT 1 # x$$\nSELECT 2 // y$$",
			want:   []string{"SELECT 1 # x", "SELECT 2 // y"},
		},
	}
	for _, tc := range tests {
		commands := splitExecuteCommands(tc.script)
		var got []string
		for _, cmd := range commands {
			got = append(got, cmd.Text)
			if cmd.ByteEnd > len(tc.script) || cmd.ByteStart < 0 || cmd.ByteStart >= cmd.ByteEnd {
				t.Errorf("%q: command range [%d:%d] is outside the script", tc.script, cmd.ByteStart, cmd.ByteEnd)
			}
		}
		if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
			t.Errorf("splitExecuteCommands(%q) = %q, want %q", tc.script, got, tc.want)
		}
	}
}

// A zero-byte value is bytes, not SQL NULL: the slice is non-nil.
func TestConvertValueEmptyBytes(t *testing.T) {
	empty := []byte{}
	if v := convertValue("BIT", nil, &empty); v.Kind != driver.Bytes || len(v.Bytes) != 0 {
		t.Errorf("empty bytes decoded as %v", v)
	}
	var null []byte
	if v := convertValue("BIT", nil, &null); v.Kind != driver.Null {
		t.Errorf("nil bytes decoded as %v", v)
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version  string
		want     string
		wantRest string
	}{
		{
			version:  "5.7.99",
			want:     "5.7.99",
			wantRest: "",
		},
		{
			version:  "5.7.22-log",
			want:     "5.7.22",
			wantRest: "-log",
		},
		{
			version:  "5.6.29_ddm_3.0.1.7",
			want:     "5.6.29",
			wantRest: "_ddm_3.0.1.7",
		},
		{
			version:  "10.4.7-MariaDB",
			want:     "10.4.7",
			wantRest: "-MariaDB",
		},
	}

	for _, tc := range tests {
		version, rest, err := parseVersion(tc.version)
		if err != nil {
			t.Fatalf("%q: %v", tc.version, err)
		}
		if version != tc.want || rest != tc.wantRest {
			t.Errorf("%q -> (%q, %q), want (%q, %q)", tc.version, version, rest, tc.want, tc.wantRest)
		}
	}
}

func TestContainsDelimiterDirective(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{"basic directive", "DELIMITER ;;\nCREATE PROCEDURE p() BEGIN SELECT 1; END;;\nDELIMITER ;", true},
		{"indented directive", "  DELIMITER //\nSELECT 1//\nDELIMITER ;", true},
		{"tab-indented directive", "\tDELIMITER $$\n", true},
		{"string literal not a directive", "INSERT INTO t VALUES ('DELIMITER ');\nSELECT 1;", false},
		{"comment not a directive", "-- DELIMITER //\nSELECT 1;", false},
		{"column name not a directive", "SELECT delimiter FROM t;", false},
		{"no delimiter at all", "SELECT 1;\nSELECT 2;", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsDelimiterDirective(tt.sql); got != tt.want {
				t.Errorf("containsDelimiterDirective = %v, want %v", got, tt.want)
			}
		})
	}
}

// BYT-9689: Doris reports materialized views in information_schema.tables with
// TABLE_TYPE='BASE TABLE', and some setups with 'MATERIALIZED VIEW', so the table type
// alone cannot identify one. Membership in the materialized-view set is the
// authoritative signal, whatever type the row carries.
func TestIsMaterializedView(t *testing.T) {
	mvMap := map[string]*metadatapb.MaterializedViewMetadata{
		"mv_async": {Name: "mv_async"},
	}

	for _, tableType := range []string{baseTableType, materializedViewType, viewTableType} {
		if !isMaterializedView(tableType, "mv_async", mvMap) {
			t.Errorf("%s in the materialized-view set is not a materialized view", tableType)
		}
	}

	// A name absent from the MV set is never a materialized view.
	if isMaterializedView(viewTableType, "v_regular", mvMap) {
		t.Error("a view is a materialized view")
	}
	if isMaterializedView(baseTableType, "sales", mvMap) {
		t.Error("a table is a materialized view")
	}
}

// #20566 review (P2): SyncInstance's system-database exclusion must match the parser's
// system set. Doris has no `sys` metadatabase, so excluding it would risk hiding a user
// database named sys.
func TestSystemDatabaseExclusion(t *testing.T) {
	got := systemDatabaseExclusion()
	for _, want := range []string{"'information_schema'", "'_statistics_'"} {
		if !strings.Contains(got, want) {
			t.Errorf("system database exclusion %q does not hold %s", got, want)
		}
	}
	if strings.Contains(got, "'sys'") {
		t.Errorf("system database exclusion %q hides a user database named sys", got)
	}
}

func TestExecuteEmptyScript(t *testing.T) {
	for _, script := range []string{"", "   \n\t ", "-- only a comment\n;"} {
		log := &stepLogger{}
		affected, err := (&Driver{}).Execute(context.Background(), script, driver.ExecuteOptions{Log: log})
		if err != nil {
			t.Fatalf("Execute(%q) err = %v", script, err)
		}
		if affected != 0 {
			t.Errorf("Execute(%q) affected = %d, want 0", script, affected)
		}
		if len(log.steps) != 0 {
			t.Errorf("Execute(%q) logged %v, want no log entries", script, log.steps)
		}
	}
}

package driver

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/driver"
)

// TestInitConfig pins the configuration the statements after CREATE DATABASE run
// with: the new database, and a dbname Params that would override it removed
// without changing the caller's map.
func TestInitConfig(t *testing.T) {
	params := map[string]string{"dbname": "dev", "search_path": "app"}
	cfg := initConfig(driver.Config{Database: "dev", Params: params}, "newdb")
	if cfg.Database != "newdb" {
		t.Errorf("Database = %q, want newdb", cfg.Database)
	}
	if _, ok := cfg.Params["dbname"]; ok {
		t.Error("the dbname parameter survived, so the initialization would run in dev")
	}
	if cfg.Params["search_path"] != "app" {
		t.Errorf("search_path = %q, want it kept", cfg.Params["search_path"])
	}
	if params["dbname"] != "dev" {
		t.Error("the caller's Params map was changed")
	}
}

// TestRefuseImplicitCommit pins which statements a transactional script may not
// hold: the ones Redshift commits on its own, so a rollback would not undo them.
func TestRefuseImplicitCommit(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{statement: "INSERT INTO facts VALUES (1)"},
		{statement: "CREATE TABLE t (a int)"},
		{statement: "TRUNCATE facts", want: true},
		{statement: "vacuum full facts", want: true},
		{statement: "CREATE DATABASE other", want: true},
		{statement: "DROP DATABASE other", want: true},
		{statement: "COMMIT", want: true},
		{statement: "END", want: true},
		{statement: "BEGIN", want: true},
	} {
		err := refuseImplicitCommit(tc.statement)
		if tc.want && err == nil {
			t.Errorf("refuseImplicitCommit(%q) accepted a statement Redshift commits on its own", tc.statement)
		}
		if !tc.want && err != nil {
			t.Errorf("refuseImplicitCommit(%q) = %v, want no error", tc.statement, err)
		}
		if tc.want && err != nil && !strings.Contains(err.Error(), "autocommit") {
			t.Errorf("refuseImplicitCommit(%q) = %v, want the autocommit hint", tc.statement, err)
		}
	}
}

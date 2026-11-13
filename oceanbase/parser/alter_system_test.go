package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestAlterSystemAction pins the words an AlterSystemStmt is named by. Whether
// each form parses is held to OceanBase by its rows in testdata/compat.tsv
// (rule alter_system_stmt and the rules under it).
func TestAlterSystemAction(t *testing.T) {
	for _, tc := range []struct{ sql, action string }{
		{"ALTER SYSTEM MAJOR FREEZE", "MAJOR FREEZE"},
		{"ALTER SYSTEM MAJOR FREEZE TENANT = t1, t2 REBUILD COLUMN GROUP", "MAJOR FREEZE"},
		{"ALTER SYSTEM MINOR FREEZE TENANT = ALL SERVER = ('127.0.0.1:2882') ZONE = z1", "MINOR FREEZE"},
		{"ALTER SYSTEM SET memory_limit = '8G' SCOPE = BOTH", "SET"},
		{"ALTER SYSTEM memory_limit = '8G'", "SET"},
		{"alter system set a = 1 comment 'c', b = -1.5 tenant = t1", "SET"},
		{"ALTER SYSTEM RESET memory_limit SCOPE = SPFILE", "RESET"},
		{"ALTER SYSTEM FLUSH PLAN CACHE TENANT = t1 GLOBAL", "FLUSH"},
		{"ALTER SYSTEM ADD SERVER '127.0.0.1:2882', '127.0.0.2:2882' ZONE 'z1'", "ADD SERVER"},
		{"ALTER SYSTEM CANCEL DELETE SERVER '127.0.0.1:2882'", "CANCEL DELETE SERVER"},
		{"ALTER SYSTEM BACKUP INCREMENTAL DATABASE PLUS ARCHIVELOG", "BACKUP INCREMENTAL DATABASE"},
		{"ALTER SYSTEM BACKUP TENANT = t1 TO 'file:///b'", "BACKUP"},
		{"ALTER SYSTEM SWITCHOVER TO STANDBY TENANT = t1 VERIFY", "SWITCHOVER TO STANDBY"},
		{"ALTER SYSTEM RECOVER TABLE d.t, d.* TO TENANT t1 FROM 'file:///b' WITH 'pool_list=p' REMAP TABLE d.t : t2", "RECOVER TABLE"},
		{"ALTER DISKGROUP g ADD DISK 'd' SERVER '127.0.0.1:2882'", "DISKGROUP ADD DISK"},
		{"SET ENCRYPTION ON IDENTIFIED BY 'p' ONLY", "SET ENCRYPTION"},
		{"SET DECRYPTION IDENTIFIED BY 'p', 'q'", "SET DECRYPTION"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		stmt, ok := list.Items[0].(*nodes.AlterSystemStmt)
		if len(list.Items) != 1 || !ok {
			t.Errorf("Parse(%q) = %T, want one *ast.AlterSystemStmt", tc.sql, list.Items[0])
			continue
		}
		if stmt.Action != tc.action {
			t.Errorf("Parse(%q): Action = %q, want %q", tc.sql, stmt.Action, tc.action)
		}
		if got := tc.sql[stmt.Loc.Start:stmt.Loc.End]; got != tc.sql {
			t.Errorf("Parse(%q): Loc spans %q", tc.sql, got)
		}
	}
}

// TestAlterSystemRejected holds a few statements no form takes, and where the
// error points: at the furthest token a form refused.
func TestAlterSystemRejected(t *testing.T) {
	for _, tc := range []struct {
		sql string
		pos int
	}{
		{"ALTER SYSTEM", 12},
		{"ALTER SYSTEM MAJOR", 18},
		{"ALTER SYSTEM MAJOR FREEZE zzz", 26},
		{"ALTER SYSTEM SET status = 1", 17},         // a keyword is no parameter's name
		{"ALTER SYSTEM FLUSH CLOG CACHE", 19},       // CLOG is no keyword of the lexer's
		{"ALTER SYSTEM SET a = 1 SCOPE = DISK", 31}, // MEMORY, SPFILE or BOTH
		{"SET ENCRYPTION ON IDENTIFIED BY 'p'", 35},
	} {
		_, err := Parse(tc.sql)
		if err == nil {
			t.Errorf("Parse(%q) was accepted", tc.sql)
			continue
		}
		if pe, ok := err.(*ParseError); !ok || pe.Position != tc.pos {
			t.Errorf("Parse(%q): error %v, want a ParseError at %d", tc.sql, err, tc.pos)
		}
	}
}

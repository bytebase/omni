package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestKillBackupResetTree pins what MariaDB's KILL, BACKUP, RESET and SHUTDOWN
// forms leave in the tree.
func TestKillBackupResetTree(t *testing.T) {
	parseOne := func(sql string) ast.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		if len(list.Items) != 1 {
			t.Fatalf("Parse(%q) gave %d statements, want 1", sql, len(list.Items))
		}
		return list.Items[0]
	}

	if k, ok := parseOne("KILL SOFT QUERY ID 7").(*ast.KillStmt); !ok || k.Signal != "SOFT" || !k.Query || !k.QueryID || k.ConnectionID == nil {
		t.Errorf("KILL SOFT QUERY ID 7 = %+v", k)
	}
	if k, ok := parseOne("KILL HARD USER 'u'@'h'").(*ast.KillStmt); !ok || k.Signal != "HARD" || k.User != "u@h" || k.ConnectionID != nil {
		t.Errorf("KILL HARD USER = %+v", k)
	}
	if b, ok := parseOne("BACKUP LOCK db.t").(*ast.BackupStmt); !ok || b.Lock == nil || b.Lock.Schema != "db" || b.Lock.Name != "t" {
		t.Errorf("BACKUP LOCK = %+v", b)
	}
	if b, ok := parseOne("BACKUP STAGE block_ddl").(*ast.BackupStmt); !ok || b.Stage != "BLOCK_DDL" {
		t.Errorf("BACKUP STAGE = %+v", b)
	}
	if b, ok := parseOne("BACKUP UNLOCK").(*ast.BackupStmt); !ok || !b.Unlock {
		t.Errorf("BACKUP UNLOCK = %+v", b)
	}
	if r, ok := parseOne("RESET SLAVE 'c' ALL").(*ast.ResetReplicaStmt); !ok || r.Channel != "c" || !r.All {
		t.Errorf("RESET SLAVE 'c' ALL = %+v", r)
	}
	if f, ok := parseOne("RESET SLAVE ALL, MASTER").(*ast.FlushStmt); !ok || len(f.Options) != 2 || f.Options[0] != "SLAVE ALL" || f.Options[1] != "MASTER" {
		t.Errorf("RESET SLAVE ALL, MASTER = %+v", f)
	}
	if s, ok := parseOne("SHUTDOWN WAIT FOR ALL REPLICAS").(*ast.ShutdownStmt); !ok || !s.WaitForAllSlaves {
		t.Errorf("SHUTDOWN WAIT FOR ALL REPLICAS = %+v", s)
	}
}

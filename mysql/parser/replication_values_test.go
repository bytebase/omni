package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// TestReplicationHexValues: a 0x number is ulong_num's HEX_NUM, so the AST
// holds its value, not zero.
func TestReplicationHexValues(t *testing.T) {
	list, err := Parse("CHANGE MASTER TO IGNORE_SERVER_IDS = (0x10, 2)")
	if err != nil {
		t.Fatal(err)
	}
	if ids := list.Items[0].(*nodes.ChangeReplicationSourceStmt).Options[0].IDs; len(ids) != 2 || ids[0] != 16 || ids[1] != 2 {
		t.Errorf("IGNORE_SERVER_IDS = %v, want [16 2]", ids)
	}
	list, err = Parse("START REPLICA UNTIL RELAY_LOG_FILE = 'r', RELAY_LOG_POS = 0x10")
	if err != nil {
		t.Fatal(err)
	}
	if pos := list.Items[0].(*nodes.StartReplicaStmt).UntilPos; pos != 16 {
		t.Errorf("UntilPos = %d, want 16", pos)
	}
}

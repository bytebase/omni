package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// A trigger named with its database, as OceanBase 4.3.5 reads it, every verdict measured by
// running the statement on oceanbase-ce:4.3.5.6 on 2026-11-13, in a database that holds the
// table t. A statement that parses may still be refused: a database that is missing answers
// 1049, and that is no verdict on the syntax.
//
// Every name in CREATE TRIGGER is an sp_name, `ident '.' ident | ident`: the trigger's, its
// table's, and the one FOLLOWS or PRECEDES names. The migration generator writes the trigger's
// database, and OceanBase stores the statement as it was written, so a trigger it made did
// not read back until this parser took the qualified name.
//
// Ref: plsql_trigger_source, trigger_definition and sp_name in OceanBase's
// src/pl/parser/pl_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 1165, 1180 and 737).
var triggerNameCases = inDatabase(triggerNameDatabase, []parserFormCase{
	{"CREATE TRIGGER DB.n01 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE TRIGGER `DB`.`n02` BEFORE INSERT ON `DB`.`t` FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE TRIGGER DB . n03 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE TRIGGER IF NOT EXISTS DB.n04 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE TRIGGER n05 BEFORE INSERT ON DB.t FOR EACH ROW SET NEW.a = 1", true},
	{"CREATE TRIGGER DB.n06 AFTER UPDATE ON t FOR EACH ROW BEGIN SET @x = 1; END", true},
	// The trigger it comes after or before is named the same way.
	{"CREATE TRIGGER DB.n07 BEFORE INSERT ON t FOR EACH ROW FOLLOWS DB.n01 SET NEW.a = 1", true},
	{"CREATE TRIGGER DB.n08 BEFORE INSERT ON t FOR EACH ROW PRECEDES n01 SET NEW.a = 1", true},
	{"CREATE TRIGGER DB.n09 BEFORE INSERT ON t FOR EACH ROW FOLLOWS `DB`.`n01` SET NEW.a = 1", true},
	// A database that is missing is found after parsing (1049).
	{"CREATE TRIGGER omni_t_no_such_database.n10 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", true},
	{"DROP TRIGGER DB.n05", true},
	{"DROP TRIGGER IF EXISTS DB.n11", true},

	// A name has two parts at most, and neither is a string.
	{"CREATE TRIGGER a.DB.n12 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE TRIGGER DB. BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE TRIGGER .n13 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE TRIGGER 'DB'.n14 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE TRIGGER DB.'n15' BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE TRIGGER DB.n16 BEFORE INSERT ON a.DB.t FOR EACH ROW SET NEW.a = 1", false},
	{"CREATE TRIGGER DB.n17 BEFORE INSERT ON t FOR EACH ROW FOLLOWS a.DB.n01 SET NEW.a = 1", false},
	{"CREATE TRIGGER DB.n18 BEFORE INSERT ON t FOR EACH ROW FOLLOWS SET NEW.a = 1", false},
	{"DROP TRIGGER a.DB.n01", false},
})

// triggerNameDatabase is the database TestTriggerNameFormsOceanBaseOracle runs in.
var triggerNameDatabase = testcontainer.DatabaseName("trigger_name_forms")

// inDatabase writes database for the word DB in each case, so a form names the database
// its test is given.
func inDatabase(database string, cases []parserFormCase) []parserFormCase {
	out := make([]parserFormCase, len(cases))
	for i, c := range cases {
		out[i] = parserFormCase{strings.ReplaceAll(c.sql, "DB", database), c.accept}
	}
	return out
}

// TestTriggerNameForms holds omni to the recorded verdicts without a server.
func TestTriggerNameForms(t *testing.T) {
	for _, c := range triggerNameCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestTriggerNameTree holds what a qualified name leaves in the tree: the database beside
// the name, for the trigger and for the one it follows.
func TestTriggerNameTree(t *testing.T) {
	for _, c := range []struct {
		sql                                  string
		schema, name, orderSchema, orderName string
	}{
		{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1", "", "tr", "", ""},
		{"CREATE TRIGGER d.tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1", "d", "tr", "", ""},
		{"CREATE TRIGGER `d`.`tr` BEFORE INSERT ON `d`.`t` FOR EACH ROW FOLLOWS d.other SET @a = 1", "d", "tr", "d", "other"},
		{"CREATE DEFINER = root@% TRIGGER d.tr BEFORE INSERT ON t FOR EACH ROW PRECEDES other SET @a = 1", "d", "tr", "", "other"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			stmt := list.Items[0].(*nodes.CreateTriggerStmt)
			if stmt.Schema != c.schema || stmt.Name != c.name {
				t.Errorf("trigger is %q.%q, want %q.%q", stmt.Schema, stmt.Name, c.schema, c.name)
			}
			var orderSchema, orderName string
			if stmt.Order != nil {
				orderSchema, orderName = stmt.Order.Schema, stmt.Order.TriggerName
			}
			if orderSchema != c.orderSchema || orderName != c.orderName {
				t.Errorf("order names %q.%q, want %q.%q", orderSchema, orderName, c.orderSchema, c.orderName)
			}
		})
	}
}

package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestSecurityPolicyShapes holds what CREATE, ALTER and DROP SECURITY POLICY
// require, measured on the pinned server under SET PARSEONLY ON.
//
// Grammar: TSql160.g:7639 createSecurityPolicyStatement_SECURITY,
// TSql160.g:8711 alterSecurityPolicyStatement_SECURITY,
// TSql160.g:9015 dropSecurityPolicyStatement_SECURITY.
func TestSecurityPolicyShapes(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		sql := "CREATE SECURITY POLICY rls.pol ADD BLOCK PREDICATE rls.fn(a) ON dbo.t AFTER INSERT, ADD FILTER PREDICATE rls.fn2() ON dbo.t WITH (STATE = ON, SCHEMABINDING = OFF) NOT FOR REPLICATION"
		list := ParseAndCheck(t, sql)
		stmt, ok := list.Items[0].(*ast.SecurityPolicyStmt)
		if !ok {
			t.Fatalf("got %T, want *ast.SecurityPolicyStmt", list.Items[0])
		}
		if stmt.Name == nil || stmt.Name.Object != "pol" {
			t.Errorf("Name = %v, want rls.pol", stmt.Name)
		}
		if !stmt.NotForReplication {
			t.Error("NotForReplication = false, want true")
		}
		if stmt.StateOn == nil || !*stmt.StateOn {
			t.Errorf("StateOn = %v, want true", stmt.StateOn)
		}
		if stmt.SchemaBinding == nil || *stmt.SchemaBinding {
			t.Errorf("SchemaBinding = %v, want false", stmt.SchemaBinding)
		}
		if stmt.Predicates == nil || stmt.Predicates.Len() != 2 {
			t.Fatalf("Predicates = %v, want two", stmt.Predicates)
		}
		block, ok := stmt.Predicates.Items[0].(*ast.SecurityPredicate)
		if !ok {
			t.Fatalf("predicate = %T, want *ast.SecurityPredicate", stmt.Predicates.Items[0])
		}
		if block.PredicateType != "BLOCK" || block.BlockDMLOp != "AFTER INSERT" {
			t.Errorf("predicate = %q, %q, want BLOCK, AFTER INSERT", block.PredicateType, block.BlockDMLOp)
		}
		checkLocation(t, sql, "SecurityPolicyStmt", stmt.Loc)
		checkLocation(t, sql, "SecurityPredicate", block.Loc)
	})

	// The body of an ALTER is one alternative, and its WITH takes STATE alone:
	// SCHEMABINDING is the CREATE's (Msg 102 at the option), a second option is
	// Msg 102 at the comma.
	for _, sql := range []string{
		"CREATE SECURITY POLICY a",
		"CREATE SECURITY POLICY a NOT FOR REPLICATION",
		"CREATE SECURITY POLICY a WITH (STATE = ON)",
		"CREATE SECURITY POLICY a WITH (SCHEMABINDING = ON, STATE = OFF)",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() ON a",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() ON a, ADD FILTER PREDICATE a() ON a",
		"CREATE SECURITY POLICY a ADD BLOCK PREDICATE a() ON a BEFORE DELETE",
		"CREATE SECURITY POLICY a ADD BLOCK PREDICATE a() ON a AFTER UPDATE",
		"ALTER SECURITY POLICY a WITH (STATE = OFF)",
		"ALTER SECURITY POLICY a ADD FILTER PREDICATE a() ON a",
		"ALTER SECURITY POLICY a ALTER BLOCK PREDICATE a() ON a BEFORE DELETE",
		"ALTER SECURITY POLICY a DROP BLOCK PREDICATE ON a, ADD FILTER PREDICATE a() ON a",
		"ALTER SECURITY POLICY a ADD NOT FOR REPLICATION",
		"ALTER SECURITY POLICY a DROP NOT FOR REPLICATION",
		"DROP SECURITY POLICY a",
		"DROP SECURITY POLICY IF EXISTS a",
	} {
		ParseAndCheck(t, sql)
	}

	// Each row is a statement the pinned server refuses under SET PARSEONLY ON,
	// and a near-miss of a rule in testdata/compat.tsv.
	for _, sql := range []string{
		// The name the grammar names is missing.
		"CREATE SECURITY POLICY",
		"ALTER SECURITY POLICY ADD NOT FOR REPLICATION",
		"DROP SECURITY POLICY",
		// The predicate's type is FILTER or BLOCK.
		"CREATE SECURITY POLICY a ADD PREDICATE a() ON a",
		"ALTER SECURITY POLICY a ADD PREDICATE a() ON a",
		"ALTER SECURITY POLICY a ALTER PREDICATE a() ON a",
		"ALTER SECURITY POLICY a DROP PREDICATE ON a",
		// A FILTER takes no BEFORE/AFTER, and only AFTER INSERT, AFTER UPDATE,
		// BEFORE UPDATE and BEFORE DELETE exist.
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() ON a AFTER INSERT",
		"CREATE SECURITY POLICY a ADD BLOCK PREDICATE a() ON a BEFORE INSERT",
		"CREATE SECURITY POLICY a ADD BLOCK PREDICATE a() ON a AFTER DELETE",
		"CREATE SECURITY POLICY a ADD BLOCK PREDICATE a() ON a AFTER",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a ON a",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() a",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE ON a",
		// A predicate list is ADD elements in a CREATE and ADD, ALTER or DROP in
		// an ALTER, each joined by a comma.
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() ON a, DROP FILTER PREDICATE ON a",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() ON a ADD FILTER PREDICATE a() ON a",
		"CREATE SECURITY POLICY a ADD FILTER PREDICATE a() ON a,",
		// The option list is STATE and SCHEMABINDING, each at most once, and an
		// ALTER holds STATE alone.
		"CREATE SECURITY POLICY a WITH ()",
		"CREATE SECURITY POLICY a WITH (STATE ON)",
		"CREATE SECURITY POLICY a WITH (STATE =)",
		"CREATE SECURITY POLICY a WITH (STATE = ON",
		"CREATE SECURITY POLICY a WITH (STATE = ON, STATE = OFF)",
		"CREATE SECURITY POLICY a WITH (SCHEMABINDING = ON, SCHEMABINDING = OFF)",
		"CREATE SECURITY POLICY a WITH (STATEMENT = ON)",
		"ALTER SECURITY POLICY a WITH (SCHEMABINDING = ON)",
		"ALTER SECURITY POLICY a WITH (STATE = ON, STATE = OFF)",
		"ALTER SECURITY POLICY a WITH (SCHEME = ON)",
		"ALTER SECURITY POLICY a WITH (",
		"ALTER SECURITY POLICY a WITH (STATE = ON) WITH (STATE = OFF)",
		"ALTER SECURITY POLICY a WITH STATE = ON)",
		// NOT FOR REPLICATION's words, and an ALTER's ADD or DROP before them.
		"CREATE SECURITY POLICY a NOT REPLICATION",
		"CREATE SECURITY POLICY a NOT FOR",
		"ALTER SECURITY POLICY a NOT FOR REPLICATION",
		"ALTER SECURITY POLICY a ADD NOT REPLICATION",
		"ALTER SECURITY POLICY a ADD NOT FOR",
		"ALTER SECURITY POLICY a DROP NOT FOR",
		// IF takes EXISTS.
		"DROP SECURITY POLICY IF a",
		"DROP SECURITY POLICY IF",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): accepted, want rejected", sql)
		} else if _, ok := err.(*ParseError); !ok {
			t.Errorf("Parse(%q): %T, want *ParseError", sql, err)
		}
	}
}

// TestSensitivityClassificationShapes holds what ADD and DROP SENSITIVITY
// CLASSIFICATION require, measured on the pinned server under SET PARSEONLY
// ON: the option names and their value types, the rank's words, a repeated
// option, an empty string, and a column named [schema_name.]table_name.column_name.
//
// Grammar: TSql160.g:8819 addSensitivityClassificationStatement_SENSITIVITY,
// TSql160.g:9021 dropSensitivityClassificationStatement_SENSITIVITY and
// TSql160.g:4317 sensitivityClassificationOption.
func TestSensitivityClassificationShapes(t *testing.T) {
	t.Run("options", func(t *testing.T) {
		sql := "ADD SENSITIVITY CLASSIFICATION TO dbo.t.c WITH (LABEL = 'a', LABEL_ID = 'b', INFORMATION_TYPE = 'c', INFORMATION_TYPE_ID = 'd', RANK = high)"
		list := ParseAndCheck(t, sql)
		stmt, ok := list.Items[0].(*ast.SensitivityClassificationStmt)
		if !ok {
			t.Fatalf("got %T, want *ast.SensitivityClassificationStmt", list.Items[0])
		}
		if stmt.Columns == nil || stmt.Columns.Len() != 1 {
			t.Fatalf("Columns = %v, want one", stmt.Columns)
		}
		column, ok := stmt.Columns.Items[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("column = %T, want *ast.TableRef", stmt.Columns.Items[0])
		}
		// A three-part name is schema.table.column, which a TableRef holds as its
		// database, schema and object.
		if column.Database != "dbo" || column.Schema != "t" || column.Object != "c" {
			t.Errorf("column = %+v, want dbo.t.c", column)
		}
		if stmt.Options == nil || stmt.Options.Len() != 5 {
			t.Fatalf("Options = %v, want five", stmt.Options)
		}
		rank, ok := stmt.Options.Items[4].(*ast.SensitivityOption)
		if !ok {
			t.Fatalf("option = %T, want *ast.SensitivityOption", stmt.Options.Items[4])
		}
		if rank.Key != "RANK" || rank.Value != "high" {
			t.Errorf("option = %q = %q, want RANK = high", rank.Key, rank.Value)
		}
		checkLocation(t, sql, "SensitivityClassificationStmt", stmt.Loc)
		checkLocation(t, sql, "SensitivityOption", rank.Loc)
	})

	for _, sql := range []string{
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = 'a')",
		"ADD SENSITIVITY CLASSIFICATION TO a.a.a WITH (LABEL = 'a')",
		"ADD SENSITIVITY CLASSIFICATION TO a.a, b.b WITH (RANK = NONE)",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (RANK = CRITICAL)",
		"DROP SENSITIVITY CLASSIFICATION FROM a.a",
		"DROP SENSITIVITY CLASSIFICATION FROM a.a.a, b.b",
	} {
		ParseAndCheck(t, sql)
	}

	// Each row is a statement the pinned server refuses under SET PARSEONLY ON,
	// and a near-miss of a rule in testdata/compat.tsv.
	for _, sql := range []string{
		// Both words, and the list's TO or FROM, are required.
		"ADD SENSITIVITY TO a.a WITH (LABEL = 'a')",
		"ADD SENSITIVITY CLASSIFICATION a.a WITH (LABEL = 'a')",
		"ADD SENSITIVITY CLASSIFICATION TO WITH (LABEL = 'a')",
		"DROP SENSITIVITY FROM a.a",
		"DROP SENSITIVITY CLASSIFICATION a.a",
		"DROP SENSITIVITY CLASSIFICATION FROM",
		// A column takes a table name, and no database part: one part is Msg
		// 16103, four are Msg 16110 and an empty part is Msg 1004.
		"ADD SENSITIVITY CLASSIFICATION TO a WITH (LABEL = 'a')",
		"DROP SENSITIVITY CLASSIFICATION FROM a",
		"ADD SENSITIVITY CLASSIFICATION TO a.a.a.a WITH (LABEL = 'a')",
		"DROP SENSITIVITY CLASSIFICATION FROM a..a",
		"ADD SENSITIVITY CLASSIFICATION TO a.a, WITH (LABEL = 'a')",
		"ADD SENSITIVITY CLASSIFICATION TO a.a,",
		"DROP SENSITIVITY CLASSIFICATION FROM a.a,",
		// WITH, its parentheses and one option are required.
		"ADD SENSITIVITY CLASSIFICATION TO a.a",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH ()",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = 'a'",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = 'a',)",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL)",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL =)",
		// The option name is one of the five the statement knows (Msg 16101).
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (a = 'a')",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = 'a', LABELS = 'b')",
		// A string is the value of every option but RANK (Msg 16115), RANK
		// takes one of five words (Msgs 16113, 16114), and no string is empty
		// (Msg 16106).
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = HIGH)",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL_ID = a)",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (RANK = 'a')",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (RANK = HIGHER)",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = '')",
		// An option is named once (Msg 16104).
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (LABEL = 'a', LABEL = 'b')",
		"ADD SENSITIVITY CLASSIFICATION TO a.a WITH (RANK = HIGH, RANK = LOW)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): accepted, want rejected", sql)
		} else if _, ok := err.(*ParseError); !ok {
			t.Errorf("Parse(%q): %T, want *ParseError", sql, err)
		}
	}
}

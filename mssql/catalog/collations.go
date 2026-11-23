package catalog

import (
	_ "embed"
	"strings"
	"sync"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
)

// collationsText is every collation SQL Server 2022 lists: the names
// sys.fn_helpcollations() returns on the pinned server, one per line after
// the comment lines. TestCollations writes it from the server and holds it, and
// lookupCollation, to what the server accepts.
//
//go:embed collations.txt
var collationsText string

// collationNames maps each collation's lower-cased name to the spelling SQL
// Server lists and stores.
var collationNames = sync.OnceValue(func() map[string]string {
	names := map[string]string{}
	for _, line := range strings.Split(collationsText, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		names[strings.ToLower(line)] = line
	}
	return names
})

// lookupCollation returns the spelling SQL Server 2022 stores for the
// collation name, which it matches without regard to case. SQL Server 2022
// accepts no version-less, _SC or _UTF8 form, and no reordering of a name's
// flags, that sys.fn_helpcollations() does not list (TestCollations).
func lookupCollation(name string) (string, bool) {
	canonical, ok := collationNames()[strings.ToLower(name)]
	return canonical, ok
}

// A COLLATE clause naming a collation SQL Server does not know fails the
// statement (448). collations.txt holds every collation the pinned SQL Server
// 2022 lists, but Bytebase runs the catalog, and mssql/review's WALK_THROUGH,
// on every version, and others know more: SQL Server 2025 adds 32 _160
// collations, such as Chinese_Simplified_Pinyin_160_CI_AS_SC_UTF8, and Azure
// SQL may add more. So a name off the list is refused only when no version
// can know it: one with no comparison style the collation grammar gives every
// collation, a case and an accent sensitivity or BIN or BIN2, after its
// designator, or one whose KS, WS, and VSS are out of the grammar's order.
//
// docs: https://learn.microsoft.com/en-us/sql/t-sql/statements/windows-collation-name-transact-sql
// docs: https://learn.microsoft.com/en-us/sql/t-sql/statements/sql-server-collation-name-transact-sql
//
// sys.fn_helpcollations() on SQL Server 2017 (14.0.3570.2) and 2019
// (15.0.4510.2) lists a subset of 2022's names, and on 2025 (17.0.5065.9)
// 2022's and the 32 _160 names, all of that shape. Each of them takes
// Compatibility_52_409_30003, which none lists, and refuses
// Compatibility_198_30A_30001, Latin1_General_CI_AS_WS_KS,
// Latin1_General_CI_AS_VSS_WS, Latin1_General_CS, and Latin1_General
// (measured, bytebase/omni#1705), so a Compatibility_ name is taken as one
// some version knows.

// collationSensitivities are the optional sensitivities of a Windows
// collation's comparison style, in the order the grammar writes them.
var collationSensitivities = []string{"KS", "WS", "VSS"}

// compatibilityPrefix begins the name of a compatibility collation, which
// SQL Server stores in this case however it is written (measured).
const compatibilityPrefix = "Compatibility_"

// collationShaped reports whether some SQL Server version may know the
// collation name: a Compatibility_ name, or a designator followed by a
// comparison style, CI or CS and AI or AS, or BIN or BIN2. A segment after
// the style the catalog does not know, such as a future flag, may be
// anything; KS, WS, and VSS are each written once and in that order.
func collationShaped(name string) bool {
	if isCompatibilityCollation(name) {
		return true
	}
	parts := strings.Split(strings.ToUpper(name), "_")
	for i := 1; i < len(parts); i++ {
		var rest []string
		switch {
		case parts[i] == "BIN" || parts[i] == "BIN2":
			rest = parts[i+1:]
		case (parts[i] == "CI" || parts[i] == "CS") && i+1 < len(parts) && (parts[i+1] == "AI" || parts[i+1] == "AS"):
			rest = parts[i+2:]
		default:
			continue
		}
		next := 0
		ordered := true
		for _, p := range rest {
			for j, s := range collationSensitivities {
				if p != s {
					continue
				}
				if j < next {
					ordered = false
				}
				next = j + 1
			}
		}
		if ordered {
			return true
		}
	}
	return false
}

func isCompatibilityCollation(name string) bool {
	return len(name) > len(compatibilityPrefix) && strings.EqualFold(name[:len(compatibilityPrefix)], compatibilityPrefix)
}

// isDefaultCollation reports database_default or catalog_default, which a
// column or an expression may name in place of a collation, though a database
// may not.
func isDefaultCollation(name string) bool {
	return strings.EqualFold(name, "database_default") || strings.EqualFold(name, "catalog_default")
}

// UnknownCollation reports a collation name no SQL Server version knows,
// matched without regard to case: one neither SQL Server 2022 lists nor
// another version may know. database_default and catalog_default name a
// collation where defaultAllowed is set: a column's or an expression's, not a
// database's or a partition function's (measured). The catalog refuses
// exactly these names (Msg 448), and so does mssql/review's WALK_THROUGH.
func UnknownCollation(name string, defaultAllowed bool) bool {
	switch {
	case name == "":
		return false
	case defaultAllowed && isDefaultCollation(name):
		return false
	}
	_, listed := lookupCollation(name)
	return !listed && !collationShaped(name)
}

// checkCollation refuses a collation name no SQL Server version knows
// (UnknownCollation, Msg 448).
func checkCollation(name string, defaultAllowed bool) error {
	if UnknownCollation(name, defaultAllowed) {
		return mssqlErr(ErrInvalidCollation, "Invalid collation '%s'.", name)
	}
	return nil
}

// storedCollation is the collation SQL Server stores for a column written
// with COLLATE name: the listed spelling of a known collation, a
// compatibility collation with its prefix spelled as SQL Server stores it, or
// name as written otherwise (database_default, a collation another version
// lists, or a name a lenient load kept).
func storedCollation(name string) string {
	if canonical, ok := lookupCollation(name); ok {
		return canonical
	}
	if isCompatibilityCollation(name) {
		return compatibilityPrefix + name[len(compatibilityPrefix):]
	}
	return name
}

// checkDeclaredCollations refuses an unknown collation that a statement
// declares for a column, a database or a partition function's input, anywhere
// in the statement, a program's body included. SQL Server checks these names
// while it parses, so the refusal (Msg 448) comes before any name in the
// statement is resolved: ALTER TABLE nope ADD a INT COLLATE a and CREATE TABLE
// t (c nope COLLATE a) are both 448 (measured). A COLLATE operator in an
// expression is checked where the expression is bound instead.
func checkDeclaredCollations(node ast.Node) error {
	var err error
	ast.Inspect(node, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		switch n := n.(type) {
		case *ast.ColumnDef:
			err = checkCollation(n.Collation, true)
		case *ast.AlterTableAction:
			if n.Type == ast.ATAlterColumn {
				err = checkCollation(n.Collation, true)
			}
		case *ast.RowsetColumnDef:
			err = checkCollation(n.Collation, true)
		case *ast.CreatePartitionFunctionStmt:
			err = checkCollation(n.Collation, false)
		case *ast.CreateDatabaseStmt:
			err = checkCollation(n.Collation, false)
		case *ast.AlterDatabaseStmt:
			if n.Action == "COLLATE" {
				err = checkCollation(n.TargetName, false)
			}
		}
		return err == nil
	})
	return err
}

// checkExprCollations checks each COLLATE operator in a table's expression (a
// column's DEFAULT, a computed column, a CHECK constraint) as SQL Server binds
// it: the operand's type first (Msg 447), then the name (Msg 448). The catalog
// does not bind these expressions otherwise, so an operand's type is known
// only for a literal or a column of the table, through class.
func checkExprCollations(expr ast.Node, class func(column string) qtype) error {
	if expr == nil {
		return nil
	}
	var err error
	ast.Inspect(expr, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		e, ok := n.(*ast.CollateExpr)
		if !ok {
			return true
		}
		// The operand is bound before the operator.
		if err = checkExprCollations(e.Expr, class); err != nil {
			return false
		}
		if t := operandClass(e.Expr, class); t != qtString && t != qtUnknown {
			err = mssqlErr(ErrCollateInvalidType, "Expression type %s is invalid for COLLATE clause.", t.typeName())
			return false
		}
		err = checkCollation(e.Collation, true)
		return false
	})
	return err
}

func operandClass(e ast.Node, class func(column string) qtype) qtype {
	switch e := e.(type) {
	case *ast.ParenExpr:
		return operandClass(e.Expr, class)
	case *ast.Literal:
		switch e.Type {
		case ast.LitString:
			return qtString
		case ast.LitInteger, ast.LitFloat, ast.LitMoney:
			return qtNumber
		}
	case *ast.ColumnRef:
		if e.Table == "" && class != nil {
			return class(e.Column)
		}
	}
	return qtUnknown
}

// columnDefExprs are the expressions a column definition holds: its DEFAULT,
// its computed expression, and its constraints' DEFAULT and CHECK.
func columnDefExprs(col *ast.ColumnDef) []ast.Node {
	var exprs []ast.Node
	if col.DefaultExpr != nil {
		exprs = append(exprs, col.DefaultExpr)
	}
	if col.Computed != nil && col.Computed.Expr != nil {
		exprs = append(exprs, col.Computed.Expr)
	}
	if col.Constraints != nil {
		for _, item := range col.Constraints.Items {
			if cd, ok := item.(*ast.ConstraintDef); ok && cd.Expr != nil {
				exprs = append(exprs, cd.Expr)
			}
		}
	}
	return exprs
}

// checkCreateTableExprCollations checks the COLLATE operators in a CREATE
// TABLE's column and table constraint expressions.
func checkCreateTableExprCollations(n *ast.CreateTableStmt) error {
	var cols []columnClass
	var exprs []ast.Node
	if n.Columns != nil {
		for _, item := range n.Columns.Items {
			if col, ok := item.(*ast.ColumnDef); ok {
				cols = append(cols, columnDefClass(col))
				exprs = append(exprs, columnDefExprs(col)...)
			}
		}
	}
	if n.Constraints != nil {
		for _, item := range n.Constraints.Items {
			if cd, ok := item.(*ast.ConstraintDef); ok && cd.Expr != nil {
				exprs = append(exprs, cd.Expr)
			}
		}
	}
	class := classOf(cols)
	for _, expr := range exprs {
		if err := checkExprCollations(expr, class); err != nil {
			return err
		}
	}
	return nil
}

// tableColumnClass is the class of a column of table, or of added, the column
// an ALTER TABLE ADD defines.
func tableColumnClass(table *metadata.TableMetadata, added *ast.ColumnDef) func(string) qtype {
	cols := make([]columnClass, 0, len(table.Columns)+1)
	for _, col := range table.Columns {
		t := typeClass(col.Type)
		if isComputedColumnType(col.Type) {
			t = qtUnknown
		}
		cols = append(cols, columnClass{col.Name, t})
	}
	if added != nil {
		cols = append(cols, columnDefClass(added))
	}
	return classOf(cols)
}

// columnClass is a column's name and the class of its type.
type columnClass struct {
	name string
	typ  qtype
}

func columnDefClass(col *ast.ColumnDef) columnClass {
	if col.Computed != nil {
		return columnClass{col.Name, qtUnknown}
	}
	return columnClass{col.Name, dataTypeClass(col.DataType)}
}

// classOf resolves a column name as SQL Server does in a database of either
// case sensitivity: a column of exactly that name, or else the one column
// whose name matches it ignoring case. A name that several columns match
// ignoring case, [c] and [C] in a case-sensitive database, has no class.
func classOf(cols []columnClass) func(string) qtype {
	return func(name string) qtype {
		for _, col := range cols {
			if col.name == name {
				return col.typ
			}
		}
		t, found := qtUnknown, false
		for _, col := range cols {
			if strings.EqualFold(col.name, name) {
				if found {
					return qtUnknown
				}
				t, found = col.typ, true
			}
		}
		return t
	}
}

package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	mssqlcatalog "github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// A COLLATE clause naming a collation SQL Server does not know fails the
// statement (448). mssqlcatalog.KnownCollation knows every collation of the
// pinned SQL Server 2022, but the walk reviews other versions too, and
// they know others: SQL Server 2025 adds 32 _160 collations, such as
// Chinese_Simplified_Pinyin_160_CI_AS_SC_UTF8, and Azure SQL may add more.
// So the walk refuses only a name no version can know: one with no
// comparison style the collation grammar gives every collation, a case and
// an accent sensitivity or BIN or BIN2, after its designator, or one whose
// KS, WS, and VSS are out of the grammar's order.
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
// (measured), so a Compatibility_ name is taken as one some version knows.

// collationSensitivities are the optional sensitivities of a Windows
// collation's comparison style, in the order the grammar writes them.
var collationSensitivities = []string{"KS", "WS", "VSS"}

// collationShaped reports whether some SQL Server version may know the
// collation name: a Compatibility_ name, or a designator followed by a
// comparison style, CI or CS and AI or AS, or BIN or BIN2. A segment after
// the style the walk does not know, such as a future flag, may be anything;
// KS, WS, and VSS are each written once and in that order.
func collationShaped(name string) bool {
	if len(name) > len("Compatibility_") && strings.EqualFold(name[:len("Compatibility_")], "Compatibility_") {
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

// unknownCollation reports a collation name no SQL Server version knows.
// database_default and catalog_default name a collation where defaultAllowed
// is set: a column's or an expression's, not a database's or a partition
// function's.
func unknownCollation(name string, defaultAllowed bool) bool {
	switch {
	case name == "":
		return false
	case defaultAllowed && mssqlcatalog.IsDefaultCollation(name):
		return false
	}
	return !mssqlcatalog.KnownCollation(name) && !collationShaped(name)
}

func msgCollation(name string) string {
	return "Collation " + review.Code(name) + " doesn't exist."
}

// declaredCollation fails a statement that declares an unknown collation
// for a column, a database, or a partition function's input, anywhere in
// it, a module's body included. SQL Server checks these names as it parses
// the batch, before it resolves any name: ALTER TABLE nope ADD a INT
// COLLATE a and CREATE PROCEDURE p AS ALTER TABLE nope ADD c NVARCHAR(10)
// COLLATE a are both 448 (measured, mssql/catalog's simulate.tsv). These
// are the clauses mssql/catalog's checkDeclaredCollations checks.
func (b *builder) declaredCollation(n ast.Node) bool {
	failed := false
	check := func(name string, defaultAllowed bool, loc ast.Loc) {
		if unknownCollation(name, defaultAllowed) {
			b.fail(rangeOf(loc), review.ObjectNotFound, msgCollation(name))
			failed = true
		}
	}
	ast.Inspect(n, func(m ast.Node) bool {
		switch v := m.(type) {
		case *ast.ColumnDef:
			check(v.Collation, true, v.Loc)
		case *ast.AlterTableAction:
			if v.Type == ast.ATAlterColumn {
				check(v.Collation, true, v.Loc)
			}
		case *ast.RowsetColumnDef:
			check(v.Collation, true, v.Loc)
		case *ast.CreatePartitionFunctionStmt:
			check(v.Collation, false, v.Loc)
		case *ast.CreateDatabaseStmt:
			check(v.Collation, false, v.Loc)
		case *ast.AlterDatabaseStmt:
			if v.Action == "COLLATE" {
				check(v.TargetName, false, v.Loc)
			}
		}
		return !failed
	})
	return failed
}

// collateFails fails the statement at the first COLLATE operator in expr
// that collateFail fails.
func (b *builder) collateFails(expr ast.Node, char func(ast.ExprNode) bool) bool {
	if expr == nil {
		return false
	}
	failed := false
	ast.Inspect(expr, func(m ast.Node) bool {
		if e, ok := m.(*ast.CollateExpr); ok && b.collateFail(e, char) {
			failed = true
		}
		return !failed
	})
	return failed
}

// collateFail fails the statement at a COLLATE operator whose operand is
// certainly a character string, as char tells, and whose collation no
// version knows (448). SQL Server binds the operand first, so an operand
// of another type is 447 instead, and one whose type the walk does not
// know says nothing. The caller checks only an expression SQL Server
// certainly binds, after what it checks first.
func (b *builder) collateFail(e *ast.CollateExpr, char func(ast.ExprNode) bool) bool {
	if !unknownCollation(e.Collation, true) || !charOperand(e.Expr, char) {
		return false
	}
	b.fail(rangeOf(e.Loc), review.ObjectNotFound, msgCollation(e.Collation))
	return true
}

// charOperand reports whether a COLLATE operand is certainly a character
// string: a string literal, another COLLATE of one, or what char vouches
// for.
func charOperand(e ast.ExprNode, char func(ast.ExprNode) bool) bool {
	switch v := e.(type) {
	case *ast.ParenExpr:
		return charOperand(v.Expr, char)
	case *ast.CollateExpr:
		return charOperand(v.Expr, char)
	case *ast.Literal:
		return v.Type == ast.LitString
	}
	return char != nil && char(e)
}

// charTypes are the system character types a collation applies to, as a
// type is named (lowercased, without its length).
var charTypes = map[string]bool{
	"char": true, "varchar": true, "nchar": true, "nvarchar": true,
	"text": true, "ntext": true, "sysname": true,
	"character": true, "char varying": true, "character varying": true,
	"national char": true, "national character": true,
	"national char varying": true, "national character varying": true,
	"national text": true,
}

// charDataType reports whether a declared type is a system character type.
func charDataType(dt *ast.DataType) bool {
	return dt != nil && isSystemBaseType(dt) && charTypes[strings.ToLower(strings.TrimSpace(dt.Name))]
}

// pathColType returns the lowercased base type of the synced column of
// exactly that spelling of the table at path, "" when the walk does not
// know it: the snapshot is not trusted, or the change may have altered
// the table. The caller knows the path means the table (knownTable).
func (b *builder) pathColType(path, col string) string {
	s := b.w.snap
	if s == nil || b.w.target.OtherSpecs || b.w.allTablesDirty || b.w.dirtyTables[path] || s.ambiguous[path] {
		return ""
	}
	st := s.tables[path]
	if st == nil {
		return ""
	}
	return st.colTypes[col]
}

// defColumnType returns the type of a column of a table definition, as a
// type is named (charTypes): one of cols spelled exactly so, else one typ
// knows, else "". A spelling that differs only in case means the same
// column or none.
func defColumnType(cols []*ast.ColumnDef, typ func(col string) string, col string) string {
	for _, c := range cols {
		if c.Name != col {
			continue
		}
		if c.Computed != nil || c.DataType == nil || !isSystemBaseType(c.DataType) {
			return ""
		}
		return strings.ToLower(strings.TrimSpace(c.DataType.Name))
	}
	if typ != nil {
		return typ(col)
	}
	return ""
}

// systemTypes reports whether every column a table definition declares has
// a system type, which SQL Server finds before it binds the definition's
// expressions; another may not exist (2715).
func systemTypes(d *tableDef) bool {
	for _, c := range d.columns {
		if c.Computed == nil && !isSystemBaseType(c.DataType) {
			return false
		}
	}
	return true
}

// exprs are the expressions a table definition binds: each column's
// DEFAULT and computed expression, and each constraint's DEFAULT or CHECK.
func (d *tableDef) exprs() []ast.Node {
	var exprs []ast.Node
	add := func(e ast.ExprNode) {
		if e != nil {
			exprs = append(exprs, e)
		}
	}
	for _, col := range d.columns {
		add(col.DefaultExpr)
		if col.Computed != nil {
			add(col.Computed.Expr)
		}
	}
	for _, cd := range d.constraints {
		add(cd.Expr)
	}
	return exprs
}

// tableCollateFails checks the COLLATE operators of a table definition's
// expressions (collateFails), each of whose columns is one of the
// definition's or one typ gives the type of; an expression naming a
// column the walk does not know may fail on it first (207).
func (b *builder) tableCollateFails(d *tableDef, typ func(col string) string) bool {
	unqualified := func(c *ast.ColumnRef) bool {
		return c.Table == "" && c.Schema == "" && c.Database == "" && c.Server == ""
	}
	char := func(e ast.ExprNode) bool {
		c, ok := e.(*ast.ColumnRef)
		return ok && unqualified(c) && charTypes[defColumnType(d.columns, typ, c.Column)]
	}
	for _, e := range d.exprs() {
		known := true
		ast.Inspect(e, func(m ast.Node) bool {
			if c, ok := m.(*ast.ColumnRef); ok && (!unqualified(c) || defColumnType(d.columns, typ, c.Column) == "") {
				known = false
			}
			return known
		})
		if known && b.collateFails(e, char) {
			return true
		}
	}
	return false
}

// columnType gives the synced type of a column of the table r means.
func (b *builder) columnType(r objRef) func(col string) string {
	return func(col string) string { return b.knownColType(r, col) }
}

// constraintBinds reports whether SQL Server certainly binds the
// expression of a CHECK or DEFAULT that ALTER TABLE ... ADD CONSTRAINT
// adds to a table the walk knows: a CHECK once the table is found, and a
// DEFAULT once the column it is FOR is found without one (1752, 1781
// first; measured, mssql/catalog).
func (b *builder) constraintBinds(r objRef, path string, known bool, c *ast.ConstraintDef) bool {
	if !known {
		return false
	}
	switch c.Type {
	case ast.ConstraintCheck:
		return true
	case ast.ConstraintDefault:
		if c.Columns == nil || len(c.Columns.Items) != 1 {
			return false
		}
		s, ok := c.Columns.Items[0].(*ast.String)
		if !ok {
			return false
		}
		cn, ok := keyed(s.Str)
		if !ok || b.see(spaceColumn, path, cn).fact != walk.Present {
			return false
		}
		st := b.syncedTable(r)
		return st != nil && st.colTypes[s.Str] != "" && st.defaults[s.Str] == ""
	}
	return false
}

// typedVar is a parameter or variable and its declared type.
type typedVar struct {
	name string
	typ  *ast.DataType
}

// typedParams lists a module's parameters with their types.
func typedParams(params *ast.List) []typedVar {
	var out []typedVar
	if params == nil {
		return nil
	}
	for _, item := range params.Items {
		if p, ok := item.(*ast.ParamDef); ok {
			out = append(out, typedVar{p.Name, p.DataType})
		}
	}
	return out
}

// charOperand reports whether a COLLATE operand of a bound query is
// certainly a character string: a variable every declaration it may mean
// declares so, or a column of the block's one source whose synced type is
// a character type. An ORDER BY name the select list may give is not
// known.
func (q *binder) charOperand(e ast.ExprNode, sc *qscope) bool {
	switch v := e.(type) {
	case *ast.VariableRef:
		found := false
		for _, d := range q.vars {
			if differ(strings.TrimPrefix(d.name, "@"), strings.TrimPrefix(v.Name, "@")) {
				continue
			}
			if !charDataType(d.typ) {
				return false
			}
			found = true
		}
		return found
	case *ast.ColumnRef:
		if sc == nil || sc.opaque || len(sc.sources) != 1 || v.Schema != "" || v.Database != "" || v.Server != "" {
			return false
		}
		src := sc.sources[0]
		if src.path == "" || v.Table != "" && differ(src.exposed, v.Table) {
			return false
		}
		if sc.order && v.Table == "" && hasName(sc.aliases, v.Column) {
			return false
		}
		return charTypes[q.b.pathColType(src.path, v.Column)]
	}
	return false
}

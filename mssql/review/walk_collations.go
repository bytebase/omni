package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	mssqlcatalog "github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// A COLLATE clause naming a collation SQL Server does not know fails the
// statement (448). The walk reviews every SQL Server version, so it refuses
// only a name no version can know, by the rule mssql/catalog refuses with
// (mssqlcatalog.UnknownCollation).

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
		if mssqlcatalog.UnknownCollation(name, defaultAllowed) {
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
// that collateFail fails, with types.
func (b *builder) collateFails(expr ast.Node, typ func(ast.ExprNode) string) bool {
	if expr == nil {
		return false
	}
	failed := false
	ast.Inspect(expr, func(m ast.Node) bool {
		if e, ok := m.(*ast.CollateExpr); ok && b.collateFail(e, typ, true) {
			failed = true
		}
		return !failed
	})
	return failed
}

// collateFail fails the statement at a COLLATE operator whose operand is
// certainly of a type no collation applies to (447), or certainly a
// character string, as typ tells, and whose collation no version knows
// (448). SQL Server binds the operand first, so its type fails before the
// collation, and an operand whose type the walk does not know says
// nothing. The caller checks only an expression SQL Server certainly
// binds, after what it checks first, and the operand's type only where
// types says that is measured.
func (b *builder) collateFail(e *ast.CollateExpr, typ func(ast.ExprNode) string, types bool) bool {
	if t := operandType(e.Expr, typ); types && nonCharTypes[t] {
		b.fail(rangeOf(e.Loc), review.UnsupportedForm, "SQL Server doesn't accept `COLLATE` on an expression of type "+review.Code(t)+".")
		return true
	}
	if !mssqlcatalog.UnknownCollation(e.Collation, true) || !charOperand(e.Expr, typ) {
		return false
	}
	b.fail(rangeOf(e.Loc), review.ObjectNotFound, msgCollation(e.Collation))
	return true
}

// charOperand reports whether a COLLATE operand is certainly a character
// string: a string literal, another COLLATE of one, or what typ vouches
// for.
func charOperand(e ast.ExprNode, typ func(ast.ExprNode) string) bool {
	switch v := e.(type) {
	case *ast.ParenExpr:
		return charOperand(v.Expr, typ)
	case *ast.CollateExpr:
		return charOperand(v.Expr, typ)
	case *ast.Literal:
		return v.Type == ast.LitString
	}
	return typ != nil && charTypes[typ(e)]
}

// operandType returns the lowercased system type an expression certainly
// has: a number literal's, or what typ tells of a column or a variable;
// "" when the walk does not know it, or for a character string.
func operandType(e ast.ExprNode, typ func(ast.ExprNode) string) string {
	switch v := e.(type) {
	case *ast.ParenExpr:
		return operandType(v.Expr, typ)
	case *ast.CollateExpr:
		return ""
	case *ast.Literal:
		switch v.Type {
		case ast.LitInteger:
			if v.Ival >= -1<<31 && v.Ival < 1<<31 {
				return "int"
			}
			return "numeric"
		case ast.LitMoney:
			return "money"
		}
		return ""
	case *ast.ColumnRef, *ast.VariableRef:
		if typ != nil {
			return typ(e)
		}
	}
	return ""
}

// nonCharTypes are system types no collation applies to (447), and that
// TRIM and the like refuse.
var nonCharTypes = map[string]bool{
	"bit": true, "tinyint": true, "smallint": true, "int": true, "bigint": true,
	"decimal": true, "numeric": true, "money": true, "smallmoney": true, "float": true, "real": true,
	"date": true, "time": true, "datetime": true, "datetime2": true, "datetimeoffset": true, "smalldatetime": true,
	"uniqueidentifier": true, "binary": true, "varbinary": true,
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

// pathColType returns the lowercased base type of the synced column of
// exactly that spelling of the table at path, "" when the walk does not
// know it: the snapshot is not trusted, or the change may have altered
// the table. The caller knows the path means the table (knownTable).
func (b *builder) pathColType(path, col string) string {
	if st := b.pathTable(path); st != nil {
		return st.colTypes[col]
	}
	return ""
}

// pathTable returns the synced table at path, nil when the snapshot is
// not trusted or the change may have altered the table. The caller knows
// the path means the table (knownTable).
func (b *builder) pathTable(path string) *snapTable {
	s := b.w.snap
	if s == nil || b.w.target.OtherSpecs || b.w.allTablesDirty || b.w.dirtyTables[path] || s.ambiguous[path] {
		return nil
	}
	return s.tables[path]
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
	colType := func(e ast.ExprNode) string {
		if c, ok := e.(*ast.ColumnRef); ok && unqualified(c) {
			return defColumnType(d.columns, typ, c.Column)
		}
		return ""
	}
	for _, e := range d.exprs() {
		known := true
		ast.Inspect(e, func(m ast.Node) bool {
			if c, ok := m.(*ast.ColumnRef); ok && (!unqualified(c) || defColumnType(d.columns, typ, c.Column) == "") {
				known = false
			}
			return known
		})
		if known && b.collateFails(e, colType) {
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

// operandType returns the lowercased system type of a variable or a
// column of a bound query: the type every declaration the variable may
// mean declares, or the synced type of a column of the block's one
// source; "" when the walk does not know it. An ORDER BY name the select
// list may give is not known.
func (q *binder) operandType(e ast.ExprNode, sc *qscope) string {
	switch v := e.(type) {
	case *ast.VariableRef:
		found := ""
		for _, d := range q.vars {
			if differ(strings.TrimPrefix(d.name, "@"), strings.TrimPrefix(v.Name, "@")) {
				continue
			}
			if d.typ == nil || !isSystemBaseType(d.typ) {
				return ""
			}
			t := strings.ToLower(strings.TrimSpace(d.typ.Name))
			if found != "" && found != t {
				return ""
			}
			found = t
		}
		return found
	case *ast.ColumnRef:
		if sc == nil || sc.opaque || len(sc.sources) != 1 || v.Schema != "" || v.Database != "" || v.Server != "" {
			return ""
		}
		src := sc.sources[0]
		if src.path == "" || v.Table != "" && differ(src.exposed, v.Table) {
			return ""
		}
		if sc.order && v.Table == "" && hasName(sc.aliases, v.Column) {
			return ""
		}
		return q.b.pathColType(src.path, v.Column)
	}
	return ""
}

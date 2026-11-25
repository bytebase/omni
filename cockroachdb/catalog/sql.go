package catalog

import (
	"fmt"
	"strings"

	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	plpgsqlparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// parseOne parses sql, which must hold exactly one statement, with
// CockroachDB's parser.
func parseOne(sql string) (tree.Statement, error) {
	stmts, err := crparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected one statement, got %d", len(stmts))
	}
	return stmts[0].AST, nil
}

// dropDatabase unqualifies every table name in node that names database, in
// place, and reports whether it changed one. CockroachDB stores the relations
// of a view, an index and a SQL routine body under their database, so a
// definition replayed on a database of another name would still read the old
// one. The formatter's hook is the one walk CockroachDB's tree offers over
// every TableName; the hook receives the node itself, so the change sticks.
func dropDatabase(node tree.NodeFormatter, database string) bool {
	if database == "" {
		return false
	}
	changed := false
	ctx := tree.NewFmtCtx(tree.FmtSimple, tree.FmtReformatTableNames(func(ctx *tree.FmtCtx, tn *tree.TableName) {
		if tn.ExplicitCatalog && string(tn.CatalogName) == database {
			tn.ExplicitCatalog = false
			tn.CatalogName = ""
			changed = true
		}
		ctx.WithReformatTableNames(nil, func() { ctx.FormatNode(tn) })
	}))
	ctx.FormatNode(node)
	ctx.Close()
	return changed
}

// relationQuery returns a view's query with the database's own name dropped
// from its relations. A query that names no relation in database, or that
// does not parse, is returned as stored.
func relationQuery(query, database string) string {
	query = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(query), ";"))
	stmt, err := parseOne(query)
	if err != nil {
		return query
	}
	if !dropDatabase(stmt, database) {
		return query
	}
	return tree.AsString(stmt)
}

// routineDefinition returns a CREATE FUNCTION or CREATE PROCEDURE with the
// database's own name dropped from the relations its body names, SQL or
// PL/pgSQL. A definition that names no relation in database, or that does not
// parse, is returned as stored.
func routineDefinition(definition, database string) string {
	definition = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(definition), ";"))
	stmt, err := parseOne(definition)
	if err != nil {
		return definition
	}
	routine, ok := stmt.(*tree.CreateRoutine)
	if !ok {
		return definition
	}
	language := ""
	for _, option := range routine.Options {
		if l, ok := option.(tree.RoutineLanguage); ok {
			language = strings.ToLower(string(l))
		}
	}
	changed := false
	for i, option := range routine.Options {
		body, ok := option.(tree.RoutineBodyStr)
		if !ok {
			continue
		}
		var text string
		switch language {
		case "sql":
			stmts, err := crparser.Parse(string(body))
			if err != nil {
				return definition
			}
			var parts []string
			for _, s := range stmts {
				if dropDatabase(s.AST, database) {
					changed = true
				}
				parts = append(parts, tree.AsString(s.AST)+";")
			}
			text = "\n" + strings.Join(parts, "\n") + "\n"
		case "plpgsql":
			block, err := plpgsqlparser.Parse(string(body))
			if err != nil {
				return definition
			}
			if dropDatabase(block.AST, database) {
				changed = true
			}
			text = "\n" + tree.AsString(block.AST) + "\n"
		default:
			return definition
		}
		routine.Options[i] = tree.RoutineBodyStr(text)
	}
	if !changed {
		return definition
	}
	return tree.AsString(routine)
}

// isProcedure reports whether a routine definition creates a procedure.
func isProcedure(definition string) bool {
	stmt, err := parseOne(definition)
	if err != nil {
		return false
	}
	routine, ok := stmt.(*tree.CreateRoutine)
	return ok && routine.IsProcedure
}

// parseIndex parses a synced index definition, as pg_indexes states it.
func parseIndex(definition string) (*tree.CreateIndex, bool) {
	if definition == "" {
		return nil, false
	}
	stmt, err := parseOne(definition)
	if err != nil {
		return nil, false
	}
	index, ok := stmt.(*tree.CreateIndex)
	return index, ok
}

// indexStatement renders a secondary index on schema.table from its synced
// definition, through CockroachDB's own formatter: the definition carries what
// IndexMetadata has no field for (STORING, a hash-sharded or inverted index, a
// partial index's predicate, visibility). A definition that does not parse is
// rendered from the index's key list.
func indexStatement(schema, table string, index *metadataIndex) string {
	if parsed, ok := parseIndex(index.Definition); ok {
		parsed.Name = tree.Name(index.Name)
		parsed.Table = tree.MakeTableNameWithSchema("", tree.Name(schema), tree.Name(table))
		parsed.Table.ExplicitCatalog = false
		parsed.IfNotExists = false
		parsed.Concurrently = false
		return tree.AsString(parsed)
	}
	var buf strings.Builder
	if index.Unique {
		buf.WriteString("CREATE UNIQUE INDEX ")
	} else {
		buf.WriteString("CREATE INDEX ")
	}
	fmt.Fprintf(&buf, "%s ON %s.%s (%s)", quoteIdent(index.Name), quoteIdent(schema), quoteIdent(table), keyList(index))
	return buf.String()
}

// primaryKeyClause renders a primary key as the constraint CREATE TABLE
// declares, with the key's directions and hash sharding from its synced
// definition.
func primaryKeyClause(index *metadataIndex) string {
	if parsed, ok := parseIndex(index.Definition); ok {
		clause := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", quoteIdent(index.Name), tree.AsString(&parsed.Columns))
		if parsed.Sharded != nil {
			clause += tree.AsString(parsed.Sharded)
		}
		if len(parsed.StorageParams) > 0 {
			clause += " WITH (" + tree.AsString(&parsed.StorageParams) + ")"
		}
		return clause
	}
	return fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", quoteIdent(index.Name), keyList(index))
}

func keyList(index *metadataIndex) string {
	parts := make([]string, len(index.Expressions))
	for i, expression := range index.Expressions {
		parts[i] = expression
		if i < len(index.Descending) && index.Descending[i] {
			parts[i] += " DESC"
		}
	}
	return strings.Join(parts, ", ")
}

// defaultExpression returns a column default as DEFAULT can take it.
// DEFAULT takes a b_expr, so an expression only an a_expr derives (COLLATE,
// a boolean operator, IS, IN, LIKE) is parenthesized: written bare, its
// COLLATE would read as the column's own collation.
func defaultExpression(expr string) string {
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return expr
	}
	switch parsed.(type) {
	case *tree.CollateExpr, *tree.AndExpr, *tree.OrExpr, *tree.NotExpr, *tree.IsNullExpr, *tree.IsNotNullExpr,
		*tree.ComparisonExpr, *tree.RangeCond, *tree.IsOfTypeExpr:
		return "(" + expr + ")"
	}
	return expr
}

// usedSequences returns the sequences, schema-qualified (public when the name
// is bare), whose nextval or currval expr calls.
func usedSequences(expr string) []string {
	if strings.TrimSpace(expr) == "" {
		return nil
	}
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return nil
	}
	v := &sequenceCollector{}
	tree.WalkExprConst(v, parsed)
	return v.names
}

type sequenceCollector struct{ names []string }

func (v *sequenceCollector) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	call, ok := expr.(*tree.FuncExpr)
	if !ok || len(call.Exprs) == 0 {
		return true, expr
	}
	name, ok := call.Func.FunctionReference.(*tree.UnresolvedName)
	if !ok || (name.Parts[0] != "nextval" && name.Parts[0] != "currval") {
		return true, expr
	}
	arg := call.Exprs[0]
	if cast, ok := arg.(*tree.CastExpr); ok {
		arg = cast.Expr
	}
	literal, ok := arg.(*tree.StrVal)
	if !ok {
		return true, expr
	}
	if sequence, err := crparser.ParseQualifiedTableName(literal.RawString()); err == nil {
		v.names = append(v.names, schemaOf(sequence)+"."+string(sequence.ObjectName))
	}
	return true, expr
}

func (*sequenceCollector) VisitPost(expr tree.Expr) tree.Expr { return expr }

// calledFunctions returns the names of the functions expr calls, as written:
// qualified or not.
func calledFunctions(expr string) []string {
	if strings.TrimSpace(expr) == "" {
		return nil
	}
	parsed, err := crparser.ParseExpr(expr)
	if err != nil {
		return nil
	}
	v := &funcCollector{}
	tree.WalkExprConst(v, parsed)
	return v.names
}

type funcCollector struct{ names []string }

func (v *funcCollector) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	if call, ok := expr.(*tree.FuncExpr); ok {
		if name, ok := call.Func.FunctionReference.(*tree.UnresolvedName); ok {
			v.names = append(v.names, name.String())
		}
	}
	return true, expr
}

func (*funcCollector) VisitPost(expr tree.Expr) tree.Expr { return expr }

// columnType returns a synced type name as CockroachDB reads it back. The
// sync states INT4 as "integer", which CockroachDB reads as INT8 under its
// default integer size.
func columnType(t string) string {
	if t == "integer" {
		return "INT4"
	}
	return t
}

// quoteIdent double-quotes a name, as pg's renderer writes every name.
func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

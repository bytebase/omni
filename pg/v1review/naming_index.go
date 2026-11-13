package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// defaultSchema is the schema Bytebase's V1 advisors looked a name up in
// when the statement gave none. The session's search_path is not consulted.
const defaultSchema = "public"

// indexCaseSensitive is whether the pre-change schema's index names are
// compared case sensitively: PostgreSQL's are.
const indexCaseSensitive = true

// NamingIndex judges the names top-level statements give non-unique indexes
// against cfg: the index a CREATE INDEX creates, and the new name an ALTER
// INDEX ... RENAME TO gives an index schema holds. CREATE UNIQUE INDEX, a
// CREATE INDEX without a name, and a statement in a CREATE SCHEMA or in a
// routine body are not judged.
//
// A renamed index is looked up by its old name in schema's "public" schema
// whatever schema the statement names, and judged only when it is found
// there and is neither unique nor primary. Its table and columns are the
// snapshot's: the table or materialized view that holds it and the index's
// expressions.
//
// cfg.Format is expanded per index with its table, without its schema, and
// its column list joined with "_". A created index's columns are its plain
// column elements as the parser stores them; an expression element is left
// out. A format that does not compile once expanded gives one
// PatternInvalid finding and that index is judged no further.
//
// The name is judged as the parser stores it: folded to lower case unless
// quoted and truncated to 63 bytes. Its length is in bytes. Every finding
// anchors the whole statement; Bytebase's V1 advisor reported the
// statement text's first non-blank line. Ported from Bytebase's V1 rule
// NAMING_INDEX_IDX.
func NamingIndex(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	naming := &templateNaming{cfg: cfg}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.IndexStmt:
			if n.Unique || n.Idxname == "" {
				continue
			}
			table := relname(n.Relation)
			if table == "" {
				continue
			}
			naming.check(n.Loc, n.Idxname, table, indexColumns(n))
		case *ast.RenameStmt:
			if n.RenameType != ast.OBJECT_INDEX {
				continue
			}
			old := relname(n.Relation)
			if old == "" {
				continue
			}
			table, index := schema.SchemaIndex(defaultSchema, old, indexCaseSensitive)
			if index != nil && !index.GetUnique() && !index.GetPrimary() {
				naming.check(n.Loc, n.Newname, table, index.GetExpressions())
			}
		}
	}
	return naming.findings
}

// templateNaming applies the config of an index or key naming rule.
type templateNaming struct {
	cfg      v1review.TemplateConfig
	findings []v1review.Finding
}

// check judges name, an index or key on table over columns.
func (c *templateNaming) check(loc ast.Loc, name, table string, columns []string) {
	start, end := span(loc)
	pattern, err := c.cfg.Expand(map[string]string{
		v1review.TableToken:      table,
		v1review.ColumnListToken: strings.Join(columns, "_"),
	})
	if err != nil {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.PatternInvalid, Name: name, Table: table, Detail: err.Error(), Start: start, End: end})
		return
	}
	if !pattern.MatchString(name) {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern.String(), Start: start, End: end})
	}
	if c.cfg.MaxLength > 0 && len(name) > c.cfg.MaxLength {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Table: table, Start: start, End: end})
	}
}

// indexColumns returns the plain columns a CREATE INDEX lists.
func indexColumns(n *ast.IndexStmt) []string {
	if n.IndexParams == nil {
		return nil
	}
	var columns []string
	for _, item := range n.IndexParams.Items {
		if elem, ok := item.(*ast.IndexElem); ok && elem.Name != "" {
			columns = append(columns, elem.Name)
		}
	}
	return columns
}

// schemaOrDefault returns the schema a relation is looked up in.
func schemaOrDefault(rv *ast.RangeVar) string {
	if rv == nil || rv.Schemaname == "" {
		return defaultSchema
	}
	return rv.Schemaname
}

// columnConstraints returns the constraints a column definition carries.
func columnConstraints(col *ast.ColumnDef) []*ast.Constraint {
	if col.Constraints == nil {
		return nil
	}
	var constraints []*ast.Constraint
	for _, item := range col.Constraints.Items {
		if c, ok := item.(*ast.Constraint); ok {
			constraints = append(constraints, c)
		}
	}
	return constraints
}

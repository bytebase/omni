package catalog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// Error is a statement Exec refuses, as Trino would: Name is Trino's
// errorName (TABLE_NOT_FOUND, NOT_SUPPORTED, ...).
type Error struct {
	Name    string
	Message string
}

func (e *Error) Error() string { return e.Name + ": " + e.Message }

func errorf(name, format string, args ...any) *Error {
	return &Error{Name: name, Message: fmt.Sprintf(format, args...)}
}

// ExecOptions controls Exec.
type ExecOptions struct {
	// ContinueOnError runs the statements after one that fails.
	ContinueOnError bool
}

// ExecResult is the outcome of one statement Exec ran.
type ExecResult struct {
	Index int
	SQL   string
	Line  int // 1-based start line in the SQL given to Exec
	Error error
}

// Exec runs the DDL statements of sql in c, in the session c's current
// catalog and schema stand for. It returns a parse error when sql does not
// parse, and otherwise one result per statement, stopping after the first
// that fails unless opts says to continue.
func (c *Catalog) Exec(sql string, opts *ExecOptions) ([]ExecResult, error) {
	file, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	var results []ExecResult
	for i, stmt := range file.Stmts {
		r := ExecResult{Index: i}
		if s, ok := stmt.(interface{ Span() ast.Loc }); ok {
			if span := s.Span(); span.Start >= 0 && span.Start <= span.End && span.End <= len(sql) {
				r.SQL = sql[span.Start:span.End]
				r.Line = 1 + strings.Count(sql[:span.Start], "\n")
			}
		}
		r.Error = c.exec(stmt, sql)
		results = append(results, r)
		if r.Error != nil && (opts == nil || !opts.ContinueOnError) {
			break
		}
	}
	return results, nil
}

var plainName = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

// quoteName writes a stored name as Trino's SHOW CREATE does: bare when it
// is a plain lower-case word that is not reserved, else in double quotes.
func quoteName(name string) string {
	if plainName.MatchString(name) {
		if kind, ok := parser.KeywordToken(name); !ok || !parser.IsReserved(kind) {
			return name
		}
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// ShowCreate writes t as Trino's SHOW CREATE TABLE prints it, in catalog
// catalogName.
func (t *Table) ShowCreate(catalogName string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TABLE %s.%s.%s (\n", quoteName(catalogName), quoteName(t.schema.Name), quoteName(t.Name))
	for i, col := range t.columns {
		if i > 0 {
			b.WriteString(",\n")
		}
		fmt.Fprintf(&b, "   %s %s", quoteName(col.Name), sqlType(col.Type))
		if col.Default != "" {
			b.WriteString(" DEFAULT " + col.Default)
		}
		if !col.Nullable {
			b.WriteString(" NOT NULL")
		}
		if col.Comment != nil {
			b.WriteString(" COMMENT " + quoteString(*col.Comment))
		}
	}
	b.WriteString("\n)")
	if t.Comment != nil {
		b.WriteString("\nCOMMENT " + quoteString(*t.Comment))
	}
	return b.String()
}

// ShowCreateHead writes the head of v as Trino's SHOW CREATE VIEW prints
// it, up to the AS before its query, in catalog catalogName.
func (v *View) ShowCreateHead(catalogName string) string {
	s := fmt.Sprintf("CREATE VIEW %s.%s.%s", quoteName(catalogName), quoteName(v.schema.Name), quoteName(v.Name))
	if v.Comment != nil {
		s += " COMMENT " + quoteString(*v.Comment)
	}
	if v.Invoker {
		return s + " SECURITY INVOKER"
	}
	return s + " SECURITY DEFINER"
}

// quoteString writes s as a Trino string literal.
func quoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

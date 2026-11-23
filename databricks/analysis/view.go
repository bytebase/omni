package analysis

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/databricks/parser"
)

// ViewDependency identifies a table or view referenced by a Databricks view
// definition or query.
type ViewDependency struct {
	Catalog string
	Schema  string
	Table   string
}

// FullName returns the dot-separated name of the dependency, omitting empty
// leading components.
func (d ViewDependency) FullName() string {
	parts := make([]string, 0, 3)
	if d.Catalog != "" {
		parts = append(parts, d.Catalog)
	}
	if d.Schema != "" {
		parts = append(parts, d.Schema)
	}
	if d.Table != "" {
		parts = append(parts, d.Table)
	}
	return strings.Join(parts, ".")
}

// ExtractViewDependencies parses sql using databricks/parser and returns the
// deduplicated list of external tables and views referenced by the query or
// view definition. It accepts full CREATE [OR REPLACE] [MATERIALIZED] VIEW and
// ALTER VIEW ... AS <query> statements as well as bare view definition queries
// (SELECT, WITH, FROM ... |>, TABLE, VALUES). CTE names in scope (including
// recursive CTEs) and dynamic IDENTIFIER(...) references are excluded.
//
// Unqualified 1-part names inherit defaultCatalog and defaultSchema; 2-part
// `schema.table` names inherit defaultCatalog; 3-part `catalog.schema.table`
// names keep all three parts as written.
func ExtractViewDependencies(sql, defaultCatalog, defaultSchema string) ([]ViewDependency, error) {
	if err := firstLexError(sql); err != nil {
		return nil, err
	}
	seg, ok := singleNonEmptySegment(sql)
	if !ok {
		return nil, errors.New("expected a single Databricks SQL view or query statement")
	}
	res, err := parser.ParseDatabricks(seg.Text)
	if err != nil {
		return nil, offsetParseError(err, seg.ByteStart)
	}
	if res == nil || len(res.Statements) != 1 {
		return nil, errors.New("expected a single Databricks SQL view or query statement")
	}
	ps, ok := res.Statements[0].(*parser.ParsedStmt)
	if !ok || ps == nil {
		return nil, nil
	}
	switch ps.Kind {
	case parser.StmtQuery, parser.StmtCreateView, parser.StmtAlterView:
	default:
		return nil, nil
	}

	seen := make(map[string]struct{}, len(ps.TableRefs))
	deps := make([]ViewDependency, 0, len(ps.TableRefs))
	for _, ref := range ps.TableRefs {
		dep, ok := resolveViewDependency(ref.Parts, defaultCatalog, defaultSchema)
		if !ok {
			continue
		}
		key := strings.ToLower(dep.Catalog) + "\x00" + strings.ToLower(dep.Schema) + "\x00" + strings.ToLower(dep.Table)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		deps = append(deps, dep)
	}
	return deps, nil
}

func resolveViewDependency(parts []string, defaultCatalog, defaultSchema string) (ViewDependency, bool) {
	switch len(parts) {
	case 1:
		return ViewDependency{Catalog: defaultCatalog, Schema: defaultSchema, Table: parts[0]}, true
	case 2:
		return ViewDependency{Catalog: defaultCatalog, Schema: parts[0], Table: parts[1]}, true
	case 3:
		return ViewDependency{Catalog: parts[0], Schema: parts[1], Table: parts[2]}, true
	default:
		if len(parts) > 3 {
			n := len(parts)
			return ViewDependency{Catalog: parts[n-3], Schema: parts[n-2], Table: parts[n-1]}, true
		}
		return ViewDependency{}, false
	}
}

// QualifyViewDDL replaces a CREATE VIEW statement's target with qualifiedName
// using databricks/parser's AST and lexer. A CREATE MATERIALIZED VIEW or any
// non-view statement is returned unchanged.
//
// Databricks' SHOW CREATE TABLE states a view as `CREATE VIEW schema.view`,
// omitting the catalog; qualifying the target ensures a dump replays into the
// right catalog.
func QualifyViewDDL(ddl, qualifiedName string) string {
	// Prefer the full AST TargetLoc when the DDL parses cleanly.
	if res, err := parser.ParseDatabricks(ddl); err == nil && res != nil && len(res.Statements) == 1 {
		if ps, ok := res.Statements[0].(*parser.ParsedStmt); ok && ps.Kind == parser.StmtCreateView && !ps.Materialized {
			if ps.TargetLoc.End > ps.TargetLoc.Start && ps.TargetLoc.End <= len(ddl) {
				return ddl[:ps.TargetLoc.Start] + qualifiedName + ddl[ps.TargetLoc.End:]
			}
		}
		return ddl
	}

	// Fallback to databricks/parser's lexer for view DDLs with Databricks-only
	// clauses beyond the Spark 4.0.5 grammar (still rejecting MATERIALIZED VIEW
	// and handling comments and escaped backticks accurately).
	l := parser.NewLexer(ddl)
	tok := l.NextToken()
	if l.Err() != nil || tok.Kind != parser.KwCreate {
		return ddl
	}
	tok = l.NextToken()
	if tok.Kind == parser.KwOr {
		tok = l.NextToken()
		if tok.Kind != parser.KwReplace {
			return ddl
		}
		tok = l.NextToken()
	}
	if tok.Kind == parser.KwGlobal {
		tok = l.NextToken()
	}
	if tok.Kind == parser.KwTemporary {
		tok = l.NextToken()
	}
	if tok.Kind != parser.KwView {
		return ddl
	}
	tok = l.NextToken()
	if tok.Kind == parser.KwIf {
		if l.NextToken().Kind != parser.KwNot || l.NextToken().Kind != parser.KwExists {
			return ddl
		}
		tok = l.NextToken()
	}
	if l.Err() != nil || !tok.Kind.IsIdentifier() {
		return ddl
	}
	start, end := tok.Start, tok.End
	for {
		next := l.NextToken()
		if l.Err() != nil || next.Kind != parser.TokDot {
			break
		}
		part := l.NextToken()
		if l.Err() != nil || !part.Kind.IsIdentifier() {
			return ddl
		}
		end = part.End
	}
	return ddl[:start] + qualifiedName + ddl[end:]
}

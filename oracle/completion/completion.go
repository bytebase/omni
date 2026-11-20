// Package completion provides parser-native auto-complete candidates for Oracle SQL.
package completion

import (
	"cmp"
	"slices"
	"strings"
	"unicode"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle"
	"github.com/bytebase/omni/oracle/analysis"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/catalog"
	"github.com/bytebase/omni/oracle/parser"
)

// CandidateType classifies an Oracle completion candidate.
type CandidateType int

const (
	CandidateKeyword CandidateType = iota
	CandidateSchema
	CandidateTable
	CandidateView
	CandidateMaterializedView
	CandidateColumn
	CandidateFunction
	CandidateProcedure
	CandidatePackage
	CandidateSequence
	CandidateType_
	CandidateIndex
	CandidateTrigger
	CandidateSynonym
	CandidateDatabaseLink
	CandidateConstraint
)

// String returns the lower-case kind name of CandidateType.
func (t CandidateType) String() string {
	switch t {
	case CandidateKeyword:
		return "keyword"
	case CandidateSchema:
		return "schema"
	case CandidateTable:
		return "table"
	case CandidateView:
		return "view"
	case CandidateMaterializedView:
		return "materialized_view"
	case CandidateColumn:
		return "column"
	case CandidateFunction:
		return "function"
	case CandidateProcedure:
		return "procedure"
	case CandidatePackage:
		return "package"
	case CandidateSequence:
		return "sequence"
	case CandidateType_:
		return "type"
	case CandidateIndex:
		return "index"
	case CandidateTrigger:
		return "trigger"
	case CandidateSynonym:
		return "synonym"
	case CandidateDatabaseLink:
		return "database_link"
	case CandidateConstraint:
		return "constraint"
	default:
		return "unknown"
	}
}

// Candidate is a single completion suggestion. Text holds the unquoted
// identifier or keyword; callers that require identifier quoting apply it at
// the editor boundary.
type Candidate struct {
	Text       string
	Type       CandidateType
	Definition string
	Comment    string
}

type options struct {
	defaultSchema string
	schemas       map[string]*metadatapb.DatabaseSchemaMetadata
	schemaNames   []string
	schemaLoader  func(schema string) *metadatapb.DatabaseSchemaMetadata
	catalog       *catalog.Catalog
}

// Option configures Complete.
type Option func(*options)

// WithDefaultSchema sets the default Oracle schema used for unqualified object
// resolution.
func WithDefaultSchema(schema string) Option {
	return func(o *options) {
		o.defaultSchema = schema
	}
}

// WithSchema registers metadata for a named Oracle schema.
func WithSchema(name string, meta *metadatapb.DatabaseSchemaMetadata) Option {
	return func(o *options) {
		if o.schemas == nil {
			o.schemas = make(map[string]*metadatapb.DatabaseSchemaMetadata)
		}
		o.schemas[name] = meta
	}
}

// WithSchemas registers metadata for multiple named Oracle schemas.
func WithSchemas(schemas map[string]*metadatapb.DatabaseSchemaMetadata) Option {
	return func(o *options) {
		if o.schemas == nil {
			o.schemas = make(map[string]*metadatapb.DatabaseSchemaMetadata, len(schemas))
		}
		for k, v := range schemas {
			o.schemas[k] = v
		}
	}
}

// WithSchemaNames registers known schema names for schema completion and
// fallback column enumeration.
func WithSchemaNames(names []string) Option {
	return func(o *options) {
		o.schemaNames = append(o.schemaNames, names...)
	}
}

// WithSchemaLoader sets a callback that lazily loads metadata for a schema
// name on demand.
func WithSchemaLoader(fn func(schema string) *metadatapb.DatabaseSchemaMetadata) Option {
	return func(o *options) {
		o.schemaLoader = fn
	}
}

// WithCatalog uses an oracle/catalog.Catalog as metadata source.
func WithCatalog(cat *catalog.Catalog) Option {
	return func(o *options) {
		o.catalog = cat
	}
}

// Complete returns Oracle SQL completion candidates for sql at byte offset
// cursorOffset. Both meta and opts are optional; when nil/empty, only keyword
// and statement-local candidates (such as CTEs and subquery aliases) are
// returned.
func Complete(sql string, cursorOffset int, meta *metadatapb.DatabaseSchemaMetadata, opts ...Option) []Candidate {
	var cfg options
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	store := newSchemaStore(meta, cfg)

	stmtSQL, localOffset := isolateStatementAtCursor(sql, cursorOffset)
	prefix := extractPrefix(stmtSQL, localOffset)
	collectOffset := localOffset - len(prefix)
	if collectOffset > 0 && collectOffset <= len(stmtSQL) && stmtSQL[collectOffset-1] == '"' && insideDoubleQuotes(stmtSQL, localOffset) {
		collectOffset--
	}

	candidates := standardComplete(stmtSQL, collectOffset, store)
	if len(candidates) == 0 {
		candidates = trickyComplete(stmtSQL, collectOffset, store)
	}

	return filterByPrefix(dedupAndSort(candidates), prefix)
}

// isolateStatementAtCursor slices sql to the statement containing cursorOffset
// so earlier statements or syntax errors do not leak into completion scope.
func isolateStatementAtCursor(sql string, cursorOffset int) (string, int) {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	stmts := analysis.SplitStatementsForCompletion(sql)
	if countNonEmptyStatements(stmts) <= 1 {
		return sql, cursorOffset
	}
	for i, stmt := range stmts {
		if cursorOffset < stmt.Start {
			continue
		}
		// At the exact ';' boundary between two statements, prefer the
		// following statement so completion sees a fresh statement start.
		if cursorOffset == stmt.End && i+1 < len(stmts) && strings.HasSuffix(strings.TrimSpace(stmt.Text), ";") {
			continue
		}
		if cursorOffset <= stmt.End {
			return sql[stmt.Start:], cursorOffset - stmt.Start
		}
	}
	return sql, cursorOffset
}

func countNonEmptyStatements(stmts []analysis.Statement) int {
	n := 0
	for _, s := range stmts {
		if !s.Empty {
			n++
		}
	}
	return n
}

func extractPrefix(sql string, cursorOffset int) string {
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	i := cursorOffset
	for i > 0 && isOracleIdentByte(sql[i-1]) {
		i--
	}
	return sql[i:cursorOffset]
}

func isOracleIdentByte(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '$' || c == '#'
}

func insideDoubleQuotes(sql string, offset int) bool {
	if offset <= 0 {
		return false
	}
	if offset > len(sql) {
		offset = len(sql)
	}
	lex := parser.NewLexer(sql[:offset])
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
	}
	return lex.Err != nil && strings.Contains(lex.Err.Error(), "unterminated")
}

func standardComplete(sql string, collectOffset int, store *schemaStore) []Candidate {
	ctx := parser.CollectCompletion(sql, collectOffset)
	if ctx == nil {
		return nil
	}
	r := newResolver(sql, collectOffset, ctx, store)
	return r.resolve()
}

// trickyComplete retries completion when standard collection yields nothing,
// first stripping any preceding SELECT that lacked a semicolon separator and
// then attempting placeholder-token patching.
func trickyComplete(sql string, collectOffset int, store *schemaStore) []Candidate {
	if collectOffset > len(sql) {
		collectOffset = len(sql)
	}
	if collectOffset < 0 {
		collectOffset = 0
	}

	if selectOffset := latestLineStartSelectOffset(sql, collectOffset); selectOffset > 0 && selectOffset < collectOffset {
		sliced := sql[selectOffset:]
		slicedOffset := collectOffset - selectOffset
		if candidates := standardComplete(sliced, slicedOffset, store); len(candidates) > 0 {
			return candidates
		}
	}

	prefix := sql[:collectOffset]
	suffix := sql[collectOffset:]
	strategies := []string{
		prefix + " __placeholder__" + suffix,
		prefix + " __placeholder__ " + suffix,
		prefix + " 1" + suffix,
	}
	for _, patched := range strategies {
		ctx := parser.CollectCompletion(patched, collectOffset)
		if ctx != nil && ctx.Candidates != nil && (len(ctx.Candidates.Tokens) > 0 || len(ctx.Candidates.Rules) > 0) {
			r := newResolver(sql, collectOffset, ctx, store)
			if candidates := r.resolve(); len(candidates) > 0 {
				return candidates
			}
		}
	}
	return nil
}

func latestLineStartSelectOffset(sql string, limit int) int {
	tokens := parser.Tokenize(sql[:limit])
	latest := -1
	for _, tok := range tokens {
		if tok.Loc >= limit {
			break
		}
		if tok.Type == parser.SELECT && (tok.Loc == 0 || sql[tok.Loc-1] == '\n') {
			latest = tok.Loc
		}
	}
	return latest
}

func filterByPrefix(candidates []Candidate, prefix string) []Candidate {
	if prefix == "" {
		return candidates
	}
	upperPrefix := strings.ToUpper(prefix)
	filtered := make([]Candidate, 0, len(candidates))
	for _, c := range candidates {
		text := strings.TrimSuffix(c.Text, "()")
		if strings.HasPrefix(strings.ToUpper(text), upperPrefix) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func dedupAndSort(candidates []Candidate) []Candidate {
	type key struct {
		typ  CandidateType
		text string
		def  string
	}
	seen := make(map[key]Candidate, len(candidates))
	for _, c := range candidates {
		k := key{typ: c.Type, text: c.Text, def: c.Definition}
		if _, exists := seen[k]; !exists {
			seen[k] = c
		}
	}
	result := make([]Candidate, 0, len(seen))
	for _, c := range seen {
		result = append(result, c)
	}
	slices.SortFunc(result, func(a, b Candidate) int {
		return cmp.Or(
			cmp.Compare(a.Type, b.Type),
			cmp.Compare(a.Text, b.Text),
			cmp.Compare(a.Definition, b.Definition),
			cmp.Compare(a.Comment, b.Comment),
		)
	})
	return result
}

func identifierNeedsQuoting(s string) bool {
	if s == "" {
		return true
	}
	if parser.IsReservedKeyword(s) {
		return true
	}
	if s != strings.ToUpper(s) {
		return true
	}
	for i, r := range s {
		if i == 0 && !unicode.IsLetter(r) {
			return true
		}
		if i > 0 && !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '$' && r != '#' {
			return true
		}
	}
	return false
}

func quoteIdentifierForSQL(identifier string) string {
	if identifierNeedsQuoting(identifier) {
		return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
	}
	return identifier
}

func parsePackageMembers(definition string) (procs []string, funcs []string) {
	if strings.TrimSpace(definition) == "" {
		return nil, nil
	}
	stmts, err := oracle.Parse(definition)
	if err != nil {
		stmts, err = oracle.Parse("CREATE OR REPLACE " + definition)
		if err != nil {
			return nil, nil
		}
	}
	// Prefer package specification declarations over body declarations when both
	// are present so private body helpers are not exposed.
	var specStmts, bodyStmts []*ast.CreatePackageStmt
	for _, s := range stmts {
		pkg, ok := s.AST.(*ast.CreatePackageStmt)
		if !ok || pkg == nil {
			continue
		}
		if pkg.IsBody {
			bodyStmts = append(bodyStmts, pkg)
		} else {
			specStmts = append(specStmts, pkg)
		}
	}
	selected := specStmts
	if len(selected) == 0 {
		selected = bodyStmts
	}
	for _, pkg := range selected {
		if pkg.Body == nil {
			continue
		}
		for _, item := range pkg.Body.Items {
			switch decl := item.(type) {
			case *ast.CreateProcedureStmt:
				if decl.Name != nil && decl.Name.Name != "" {
					procs = append(procs, decl.Name.Name)
				}
			case *ast.CreateFunctionStmt:
				if decl.Name != nil && decl.Name.Name != "" {
					funcs = append(funcs, decl.Name.Name)
				}
			}
		}
	}
	return procs, funcs
}

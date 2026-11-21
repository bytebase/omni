package catalog

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// DefaultSDLDatabase is the synthetic database name LoadSDLInDatabase and
// BuildSDLMigrationPlan use when no explicit database is specified.
const DefaultSDLDatabase = "bbcatalog"

// SDLMigrationOptions configures BuildSDLMigrationPlan.
type SDLMigrationOptions struct {
	// Database is the synthetic database context both source and target SDL
	// are loaded under so unqualified DDL resolves. Empty defaults to
	// DefaultSDLDatabase.
	Database string

	// SourceSessionContext, when non-nil, carries per-object session context
	// (sql_mode, character_set_client, collation_connection, and event
	// time_zone) of the source schema. It is stamped onto the source catalog
	// after load and before Diff so recreate/ALTER statements preserve the
	// source object's original session context.
	SourceSessionContext *SessionContextMap
}

// LoadSDLInDatabase parses a single-database MySQL schema text under dbName
// (defaulting to DefaultSDLDatabase when empty) and returns a catalog
// reflecting it.
//
// The database is pre-created directly on the catalog rather than prepended as
// a synthetic SQL prologue, so parse errors report exact caller-relative line
// and column numbers without offset rewriting.
//
// LoadSDL is the canonical, dependency-aware path. LoadSQL is tried only as a
// fallback for raw dumps containing non-SDL statements (e.g. SET, SELECT,
// LOCK TABLES) that LoadSDL rejects; its session defaults already match MySQL
// 8.0 (including explicit_defaults_for_timestamp). CREATE TABLE ... AS SELECT
// is never accepted via the raw-SQL fallback because createTable cannot
// materialize its schema and would silently omit the table. When both paths
// fail, the LoadSDL error is returned because it is more diagnostic than
// LoadSQL's order-dependent failures.
func LoadSDLInDatabase(sql, dbName string) (*Catalog, error) {
	if dbName == "" {
		dbName = DefaultSDLDatabase
	}
	c, sdlErr := loadSDL(sql, MySQL80, dbName)
	if sdlErr == nil {
		return c, nil
	}
	if errors.Is(sdlErr, errSDLDisallowedCTAS) {
		return nil, sdlErr
	}
	if c, err := loadSQLInDatabase(sql, dbName); err == nil {
		return c, nil
	} else if errors.Is(err, errSDLDisallowedCTAS) {
		return nil, err
	}
	return nil, sdlErr
}

// loadSQLInDatabase pre-creates and selects dbName, then executes sql via
// Catalog.Exec. It rejects CREATE TABLE ... AS SELECT so the raw-SQL fallback
// cannot silently omit tables that loadSDL refused.
func loadSQLInDatabase(sql, dbName string) (*Catalog, error) {
	list, err := mysqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list != nil {
		for _, stmt := range list.Items {
			if ct, ok := stmt.(*nodes.CreateTableStmt); ok && ct.Select != nil {
				return nil, errSDLDisallowedCTAS
			}
		}
	}
	c := New()
	if dbName != "" {
		if err := c.createDatabase(&nodes.CreateDatabaseStmt{Name: dbName, IfNotExists: true}); err != nil {
			return nil, err
		}
		c.SetCurrentDatabase(dbName)
	}
	results, err := c.Exec(sql, nil)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.Error != nil {
			return c, r.Error
		}
	}
	return c, nil
}

// BuildSDLMigrationPlan loads sourceSDL and targetSDL into catalogs under the
// same database context, applies opts.SourceSessionContext to the source
// catalog when present, diffs them, and returns the resulting MigrationPlan.
// Returns (nil, nil) when there are no changes.
func BuildSDLMigrationPlan(sourceSDL, targetSDL string, opts SDLMigrationOptions) (*MigrationPlan, error) {
	dbName := opts.Database
	if dbName == "" {
		dbName = DefaultSDLDatabase
	}
	from, err := LoadSDLInDatabase(sourceSDL, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to load source schema: %w", err)
	}
	if opts.SourceSessionContext != nil {
		from.ApplySessionContext(*opts.SourceSessionContext)
	}
	to, err := LoadSDLInDatabase(targetSDL, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to load target schema: %w", err)
	}
	diff := Diff(from, to)
	if diff.IsEmpty() {
		return nil, nil
	}
	return GenerateMigration(from, to, diff), nil
}

// SessionContextMapFromMetadata extracts per-object session context
// (sql_mode, character_set_client, collation_connection, and event time_zone)
// from a synced database metadata snapshot, keyed by lower-cased object name.
// Returns nil when meta is nil or contains no routines, triggers, or events.
func SessionContextMapFromMetadata(meta *metadata.DatabaseSchemaMetadata) *SessionContextMap {
	if meta == nil {
		return nil
	}
	m := &SessionContextMap{
		Functions:  map[string]SessionContext{},
		Procedures: map[string]SessionContext{},
		Triggers:   map[string]SessionContext{},
		Events:     map[string]SessionContext{},
	}
	for _, sm := range meta.GetSchemas() {
		for _, fn := range sm.GetFunctions() {
			m.Functions[toLower(fn.GetName())] = SessionContext{
				SQLMode:             fn.GetSqlMode(),
				CharacterSetClient:  fn.GetCharacterSetClient(),
				CollationConnection: fn.GetCollationConnection(),
			}
		}
		for _, proc := range sm.GetProcedures() {
			m.Procedures[toLower(proc.GetName())] = SessionContext{
				SQLMode:             proc.GetSqlMode(),
				CharacterSetClient:  proc.GetCharacterSetClient(),
				CollationConnection: proc.GetCollationConnection(),
			}
		}
		for _, event := range sm.GetEvents() {
			m.Events[toLower(event.GetName())] = SessionContext{
				SQLMode:             event.GetSqlMode(),
				CharacterSetClient:  event.GetCharacterSetClient(),
				CollationConnection: event.GetCollationConnection(),
				TimeZone:            event.GetTimeZone(),
			}
		}
		for _, table := range sm.GetTables() {
			for _, trigger := range table.GetTriggers() {
				m.Triggers[toLower(trigger.GetName())] = SessionContext{
					SQLMode:             trigger.GetSqlMode(),
					CharacterSetClient:  trigger.GetCharacterSetClient(),
					CollationConnection: trigger.GetCollationConnection(),
				}
			}
		}
	}
	if len(m.Functions) == 0 && len(m.Procedures) == 0 && len(m.Triggers) == 0 && len(m.Events) == 0 {
		return nil
	}
	return m
}

// StripDatabaseContext removes the synthetic database qualifier (`<dbName>`. )
// from generated DDL in database-qualifier AST positions while preserving
// string literals, comments, and table or alias qualifiers that happen to share
// dbName. Empty dbName defaults to DefaultSDLDatabase.
func StripDatabaseContext(sql, dbName string) string {
	if dbName == "" {
		dbName = DefaultSDLDatabase
	}
	qualifier := quoteIdent(dbName) + "."
	if !strings.Contains(sql, qualifier) {
		return sql
	}

	segs := mysqlparser.Split(sql)
	if len(segs) == 0 {
		return sql
	}

	var starts []int
	seen := make(map[int]bool)
	record := func(pos int) {
		if pos >= 0 && pos+len(qualifier) <= len(sql) && strings.HasPrefix(sql[pos:], qualifier) && !seen[pos] {
			seen[pos] = true
			starts = append(starts, pos)
		}
	}

	for _, seg := range segs {
		if seg.Empty() || !strings.Contains(seg.Text, qualifier) {
			continue
		}
		list, err := mysqlparser.Parse(seg.Text)
		if err != nil || list == nil {
			continue
		}
		base := seg.ByteStart
		for _, stmt := range list.Items {
			collectDatabaseQualifierStarts(stmt, seg.Text, base, dbName, qualifier, record)
		}
	}

	if len(starts) == 0 {
		return sql
	}
	slices.Sort(starts)

	var out strings.Builder
	out.Grow(len(sql) - len(starts)*len(qualifier))
	prev := 0
	for _, start := range starts {
		if start < prev {
			continue
		}
		out.WriteString(sql[prev:start])
		prev = start + len(qualifier)
	}
	out.WriteString(sql[prev:])
	return out.String()
}

// collectDatabaseQualifierStarts walks a parsed statement AST and invokes
// record(absOffset) for each database qualifier matching dbName. Because
// 2-part table/alias column references populate ColumnRef.Table (leaving
// ColumnRef.Schema empty), aliases and tables named dbName are preserved.
func collectDatabaseQualifierStarts(stmt nodes.Node, segText string, base int, dbName, qualifier string, record func(int)) {
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		switch n := n.(type) {
		case *nodes.TableRef:
			if strings.EqualFold(n.Schema, dbName) && n.Loc.Start >= 0 {
				record(base + n.Loc.Start)
			}
		case *nodes.ColumnRef:
			if strings.EqualFold(n.Schema, dbName) && n.Loc.Start >= 0 {
				record(base + n.Loc.Start)
			}
		case *nodes.FuncCallExpr:
			if strings.EqualFold(n.Schema, dbName) && n.Loc.Start >= 0 {
				record(base + n.Loc.Start)
			}
		case *nodes.CreateTriggerStmt:
			if strings.EqualFold(n.Schema, dbName) {
				limit := n.Loc.End
				if n.Table != nil && n.Table.Loc.Start >= 0 && n.Table.Loc.Start < limit {
					limit = n.Table.Loc.Start
				}
				if pos := findQualifiedIdent(segText, n.Loc.Start, limit, qualifier, n.Name); pos >= 0 {
					record(base + pos)
				}
			}
		case *nodes.CreateEventStmt:
			if strings.EqualFold(n.Schema, dbName) {
				if pos := findQualifiedIdent(segText, n.Loc.Start, n.Loc.End, qualifier, n.Name); pos >= 0 {
					record(base + pos)
				}
			}
		case *nodes.AlterEventStmt:
			renameStart := n.Loc.Start
			if strings.EqualFold(n.Schema, dbName) {
				if pos := findQualifiedIdent(segText, n.Loc.Start, n.Loc.End, qualifier, n.Name); pos >= 0 {
					record(base + pos)
					renameStart = pos + len(qualifier)
				}
			}
			if strings.EqualFold(n.RenameSchema, dbName) {
				if pos := findQualifiedIdent(segText, renameStart, n.Loc.End, qualifier, n.RenameTo); pos >= 0 {
					record(base + pos)
				}
			}
		case *nodes.DropEventStmt:
			if strings.EqualFold(n.Schema, dbName) {
				if pos := findQualifiedIdent(segText, n.Loc.Start, n.Loc.End, qualifier, n.Name); pos >= 0 {
					record(base + pos)
				}
			}
		}
		return true
	})
}

// findQualifiedIdent locates qualifier followed by name (quoted or bare) in
// sql[start:end], skipping string literals and non-executable comments. Used
// for trigger/event statements whose Schema/Name fields do not carry a separate
// name Loc.
func findQualifiedIdent(sql string, start, end int, qualifier, name string) int {
	if start < 0 {
		start = 0
	}
	if end <= start || end > len(sql) {
		end = len(sql)
	}
	quotedDB := strings.TrimSuffix(qualifier, ".")
	quotedName := quoteIdent(name)
	for i := start; i < end; {
		if isViewBodySpaceByte(sql[i]) {
			i++
			continue
		}
		if cEnd, isComment := scanMySQLComment(sql, i); isComment {
			i = cEnd
			continue
		}
		tokEnd := scanViewBodyToken(sql, i)
		if tokEnd > end {
			tokEnd = end
		}
		if sql[i:tokEnd] == quotedDB && tokEnd < end && sql[tokEnd] == '.' {
			nextStart := tokEnd + 1
			if nextStart < end {
				nextEnd := scanViewBodyToken(sql, nextStart)
				if nextEnd > end {
					nextEnd = end
				}
				nextTok := sql[nextStart:nextEnd]
				if nextTok == quotedName || strings.EqualFold(nextTok, name) {
					return i
				}
			}
		}
		i = tokEnd
	}
	return -1
}

// scanMySQLComment reports whether a non-executable MySQL comment starts at i
// and returns the offset just past it. Version-gated executable comments
// (/*!...*/) are left to the normal token scan because their interior is live
// SQL.
func scanMySQLComment(sql string, i int) (int, bool) {
	switch sql[i] {
	case '#':
		return mysqlLineCommentEnd(sql, i), true
	case '-':
		if i+1 < len(sql) && sql[i+1] == '-' && (i+2 >= len(sql) || isViewBodySpaceByte(sql[i+2])) {
			return mysqlLineCommentEnd(sql, i), true
		}
	case '/':
		if i+1 < len(sql) && sql[i+1] == '*' && (i+2 >= len(sql) || sql[i+2] != '!') {
			if idx := strings.Index(sql[i+2:], "*/"); idx >= 0 {
				return i + 2 + idx + 2, true
			}
			return len(sql), true
		}
	}
	return 0, false
}

func mysqlLineCommentEnd(sql string, i int) int {
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	return i
}

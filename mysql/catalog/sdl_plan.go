package catalog

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/mysql/ast"
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
// 8.0 (including explicit_defaults_for_timestamp). When both fail, the LoadSDL
// error is returned because it is more diagnostic than LoadSQL's
// order-dependent failures.
func LoadSDLInDatabase(sql, dbName string) (*Catalog, error) {
	if dbName == "" {
		dbName = DefaultSDLDatabase
	}
	c, sdlErr := loadSDL(sql, MySQL80, dbName)
	if sdlErr == nil {
		return c, nil
	}
	if c, err := loadSQLInDatabase(sql, dbName); err == nil {
		return c, nil
	}
	return nil, sdlErr
}

// loadSQLInDatabase pre-creates and selects dbName, then executes sql via
// Catalog.Exec. It is the LoadSQL fallback paired with loadSDL.
func loadSQLInDatabase(sql, dbName string) (*Catalog, error) {
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
// from generated DDL in identifier positions while preserving string literals
// and non-executable comments. Empty dbName defaults to DefaultSDLDatabase.
//
// Executable comments (/*!...*/) are scanned as live SQL; non-executable
// comments (#..., -- ..., /*...*/) are copied opaquely so unbalanced quotes
// inside routine comments cannot swallow subsequent qualifiers.
func StripDatabaseContext(sql, dbName string) string {
	if dbName == "" {
		dbName = DefaultSDLDatabase
	}
	quoted := "`" + strings.ReplaceAll(dbName, "`", "``") + "`"
	if !strings.Contains(sql, quoted+".") {
		return sql
	}
	var out strings.Builder
	out.Grow(len(sql))
	for i := 0; i < len(sql); {
		c := sql[i]
		if isViewBodySpaceByte(c) {
			out.WriteByte(c)
			i++
			continue
		}
		if end, isComment := scanMySQLComment(sql, i); isComment {
			out.WriteString(sql[i:end])
			i = end
			continue
		}
		end := scanViewBodyToken(sql, i)
		tok := sql[i:end]
		if tok == quoted && end < len(sql) && sql[end] == '.' && end+1 < len(sql) && sql[end+1] == '`' {
			identEnd := scanViewBodyToken(sql, end+1)
			out.WriteString(sql[end+1 : identEnd])
			i = identEnd
			continue
		}
		out.WriteString(tok)
		i = end
	}
	return out.String()
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

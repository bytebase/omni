// Package mysqlcomplete implements shared MySQL-family session completion,
// cross-database qualifier analysis, and identifier quoting helpers.
package mysqlcomplete

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
)

// CandidateType classifies a completion candidate. Values 0..13 match the
// CandidateType constants in mysql/completion, mariadb/completion,
// oceanbase/completion, and tidb/completion.
type CandidateType int

const (
	CandidateNone      CandidateType = -1
	CandidateKeyword   CandidateType = 0
	CandidateDatabase  CandidateType = 1
	CandidateTable     CandidateType = 2
	CandidateView      CandidateType = 3
	CandidateColumn    CandidateType = 4
	CandidateFunction  CandidateType = 5
	CandidateProcedure CandidateType = 6
	CandidateIndex     CandidateType = 7
	CandidateTrigger   CandidateType = 8
	CandidateEvent     CandidateType = 9
	CandidateVariable  CandidateType = 10
	CandidateCharset   CandidateType = 11
	CandidateEngine    CandidateType = 12
	CandidateType_     CandidateType = 13
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string
	Type       CandidateType
	Definition string
	Comment    string
}

// SessionConfig configures multi-database, session-aware SQL completion.
type SessionConfig struct {
	// ReadOnly drops statement-initiating write/DDL/admin keywords when the
	// caret sits at the start of a statement (e.g. a read-only query editor).
	ReadOnly bool
	// DefaultDatabase is the active session database, if any.
	DefaultDatabase string
	// LoadDatabase fetches metadata for a database name on demand.
	// resolvedName may canonicalize the name (or be empty to keep name).
	LoadDatabase func(ctx context.Context, name string) (resolvedName string, meta *metadatapb.DatabaseSchemaMetadata, caseSensitive bool, ok bool)
	// ListDatabaseNames enumerates known databases when the cursor position can
	// suggest database names or when resolving a case-folded qualifier lazily.
	ListDatabaseNames func(ctx context.Context) ([]string, error)
}

// CatalogHooks binds an engine's catalog instance for a single completion call.
type CatalogHooks struct {
	RegisterDB   func(name string)
	LoadDB       func(name string, meta *metadatapb.DatabaseSchemaMetadata)
	SetCurrentDB func(name string)
	Complete     func(stmt string, pos int) []Candidate
}

// Dialect supplies engine-specific parser and catalog hooks for session completion.
type Dialect struct {
	CurrentStatement  func(statement string, pos int) (string, int)
	IsReservedKeyword func(name string) bool
	HasDatabaseRule   func(sql string, offset int) bool
	NewCatalog        func(ctx context.Context) CatalogHooks
}

// writeStatementKeywords are statement-initiating keywords for write, DDL,
// transaction-control, and admin statements. They are dropped from completion
// in read-only mode, leaving only read statements (SELECT, WITH, EXPLAIN,
// DESC/DESCRIBE, SHOW, HELP, TABLE, VALUES) and clause keywords.
var writeStatementKeywords = map[string]bool{
	// DML / data loading.
	"INSERT": true, "UPDATE": true, "DELETE": true, "REPLACE": true,
	"LOAD": true, "IMPORT": true, "BATCH": true,
	// DDL.
	"CREATE": true, "ALTER": true, "DROP": true, "TRUNCATE": true, "RENAME": true,
	// Privileges / session.
	"GRANT": true, "REVOKE": true, "SET": true, "USE": true,
	// Transaction control.
	"BEGIN": true, "START": true, "COMMIT": true, "ROLLBACK": true,
	"SAVEPOINT": true, "RELEASE": true, "XA": true,
	// Admin / maintenance / replication.
	"ANALYZE": true, "OPTIMIZE": true, "FLUSH": true, "REPAIR": true,
	"CHECK": true, "CHECKSUM": true, "RESET": true, "KILL": true,
	"LOCK": true, "UNLOCK": true, "CACHE": true,
	"SHUTDOWN": true, "RESTART": true, "CLONE": true,
	"INSTALL": true, "UNINSTALL": true, "CHANGE": true, "STOP": true,
	"PURGE": true, "BINLOG": true, "HANDLER": true,
	"SIGNAL": true, "RESIGNAL": true,
	// TiDB maintenance / backup / resource / region statements.
	"ADMIN": true, "SPLIT": true, "DISTRIBUTE": true, "CALIBRATE": true,
	"QUERY": true, "BACKUP": true, "RESTORE": true, "RECOVER": true,
	"FLASHBACK": true, "PLAN": true,
	// Prepared statements / routines.
	"PREPARE": true, "EXECUTE": true, "DEALLOCATE": true, "CALL": true, "DO": true,
}

// CompleteSession runs session-aware auto-completion for a MySQL-family dialect.
// CaretLine is 1-based and caretOffset is a 0-based UTF-16 code-unit column.
func CompleteSession(
	ctx context.Context,
	cfg SessionConfig,
	dialect Dialect,
	statement string,
	caretLine int,
	caretOffset int,
) []Candidate {
	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	stmt, pos := statement, bytePos
	if dialect.CurrentStatement != nil {
		stmt, pos = dialect.CurrentStatement(statement, bytePos)
	}
	pos = max(0, min(pos, len(stmt)))

	cat := dialect.NewCatalog(ctx)
	PopulateCatalog(ctx, cfg, stmt, pos, dialect.HasDatabaseRule, cat.RegisterDB, cat.LoadDB, cat.SetCurrentDB)

	caretInBacktick := CaretInsideBacktickIdentifier(stmt, pos)
	type candKey struct {
		typ        CandidateType
		text       string
		definition string
	}
	candidateMap := make(map[candKey]Candidate)
	for _, c := range cat.Complete(stmt, pos) {
		t := c.Type
		if t == CandidateNone {
			continue
		}
		text := c.Text
		switch {
		case isObjectIdentifierCandidate(t):
			text = QuoteIdentifierIfNeeded(c.Text, caretInBacktick, dialect.IsReservedKeyword)
		case t == CandidateFunction:
			if !strings.HasSuffix(text, "()") {
				text += "()"
			}
		default:
			// Keywords and value-like candidates pass through unchanged.
		}
		candidate := Candidate{
			Type:       t,
			Text:       text,
			Definition: c.Definition,
			Comment:    c.Comment,
		}
		candidateMap[candKey{typ: t, text: text, definition: c.Definition}] = candidate
	}

	// In read-only mode, drop keywords that initiate writes at statement start
	// while keeping the same words when valid as sub-keywords (e.g. SHOW CREATE
	// TABLE, SELECT ... FOR UPDATE).
	if cfg.ReadOnly && CaretAtStatementStart(stmt, pos) {
		for key, c := range candidateMap {
			if c.Type == CandidateKeyword && writeStatementKeywords[strings.ToUpper(c.Text)] {
				delete(candidateMap, key)
			}
		}
	}

	result := make([]Candidate, 0, len(candidateMap))
	for _, c := range candidateMap {
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

// PopulateCatalog populates a MySQL-family catalog for statement at pos.
// It loads the default database plus any database referenced as a database
// qualifier in the statement, and enumerates unreferenced database names only
// when the cursor position can suggest database_ref candidates or when a
// case-folded qualifier needs canonical resolution.
func PopulateCatalog(
	ctx context.Context,
	cfg SessionConfig,
	statement string,
	pos int,
	hasDatabaseRule func(sql string, offset int) bool,
	registerDB func(name string),
	loadDB func(name string, meta *metadatapb.DatabaseSchemaMetadata),
	setCurrentDB func(name string),
) {
	suggestsDB := hasDatabaseRule != nil && suggestsDatabaseRef(statement, pos, hasDatabaseRule)

	var allNames []string
	enumerated := false
	ensureEnumerated := func() {
		if enumerated || cfg.ListDatabaseNames == nil {
			return
		}
		enumerated = true
		if names, err := cfg.ListDatabaseNames(ctx); err == nil {
			allNames = names
		}
	}
	if suggestsDB {
		ensureEnumerated()
	}

	findExactFolded := func(name string) (string, bool) {
		ensureEnumerated()
		lower := strings.ToLower(name)
		for _, n := range allNames {
			if strings.ToLower(n) == lower {
				return n, n != name
			}
		}
		return "", false
	}

	// Load the default database first so its metadata determines whether this
	// instance treats database identifiers as case-sensitive or case-folding.
	caseSensitive := true
	defaultLoaded := false
	seen := map[string]bool{}
	keyFor := func(name string) string {
		if caseSensitive {
			return name
		}
		return strings.ToLower(name)
	}

	defaultName := cfg.DefaultDatabase
	var defaultProto *metadatapb.DatabaseSchemaMetadata
	if cfg.LoadDatabase != nil && defaultName != "" {
		if resolvedName, proto, cs, ok := cfg.LoadDatabase(ctx, defaultName); ok {
			if resolvedName != "" {
				defaultName = resolvedName
			}
			caseSensitive = cs
			defaultProto = proto
			defaultLoaded = true
		}
		seen[keyFor(defaultName)] = true
	}

	qualifiers, caretDBs := AnalyzeStatementQualifiers(statement, pos, cfg.DefaultDatabase, caseSensitive)

	// Omni's catalog maps database names case-insensitively. On a case-sensitive
	// instance where multiple DBs differ only by case (e.g. `Foo` vs `foo`),
	// prioritize the databases active at the caret (with projected CTE/subquery
	// DBs first) over colliding siblings, and defer defaultDB only when the
	// caret's leading DB is the colliding sibling.
	if caseSensitive && len(caretDBs) > 0 {
		for idx := len(caretDBs) - 1; idx >= 0; idx-- {
			target := caretDBs[idx]
			for i, q := range qualifiers {
				if q == target && i > 0 {
					qualifiers = append([]string{q}, append(qualifiers[:i:i], qualifiers[i+1:]...)...)
					break
				}
			}
		}
	}
	deferDefault := caseSensitive && shouldDeferDefaultDB(defaultName, caretDBs)
	loadDefault := func() {
		if cfg.LoadDatabase == nil || defaultName == "" {
			return
		}
		if defaultProto != nil {
			loadDB(defaultName, defaultProto)
		} else {
			registerDB(defaultName)
		}
	}
	if !deferDefault {
		loadDefault()
	}

	if cfg.LoadDatabase != nil {
		for _, q := range qualifiers {
			loadName := q
			if !caseSensitive && enumerated {
				if exact, diffCase := findExactFolded(q); diffCase {
					loadName = exact
				}
			}
			if seen[keyFor(loadName)] {
				continue
			}
			var proto *metadatapb.DatabaseSchemaMetadata
			if resolvedName, p, _, ok := cfg.LoadDatabase(ctx, loadName); ok {
				if resolvedName != "" {
					loadName = resolvedName
				}
				proto = p
			} else if !defaultLoaded || !caseSensitive {
				// Case-folding instances may store "db_b" while the query uses
				// "DB_B"; resolve canonically on miss when defaultDB is
				// case-folding or unavailable.
				if exact, diffCase := findExactFolded(loadName); diffCase {
					if resolvedName, p, cs, ok := cfg.LoadDatabase(ctx, exact); ok && !cs {
						loadName = exact
						if resolvedName != "" {
							loadName = resolvedName
						}
						proto = p
					}
				}
			}
			seen[keyFor(loadName)] = true
			seen[keyFor(q)] = true
			if proto != nil {
				loadDB(loadName, proto)
			} else if suggestsDB {
				registerDB(loadName)
			}
		}
	}
	if deferDefault {
		loadDefault()
	}

	if suggestsDB {
		for _, name := range allNames {
			if !seen[keyFor(name)] {
				seen[keyFor(name)] = true
				registerDB(name)
			}
		}
	}

	if defaultName != "" {
		setCurrentDB(defaultName)
	}
}

func suggestsDatabaseRef(stmt string, pos int, hasDatabaseRule func(sql string, offset int) bool) bool {
	pos = max(0, min(pos, len(stmt)))
	collectOffset := pos - extractIdentPrefixLen(stmt, pos)
	if hasDatabaseRule(stmt, collectOffset) {
		return true
	}
	prefix := stmt[:collectOffset]
	suffix := stmt[collectOffset:]
	for _, patched := range []string{
		prefix + " __placeholder__" + suffix,
		prefix + " __placeholder__ " + suffix,
		prefix + " 1" + suffix,
	} {
		if hasDatabaseRule(patched, collectOffset) {
			return true
		}
	}
	return false
}

func extractIdentPrefixLen(sql string, cursorOffset int) int {
	end := max(0, min(cursorOffset, len(sql)))
	i := end
	for i > 0 {
		if c := sql[i-1]; c < utf8.RuneSelf {
			if !isIdentByte(c) {
				break
			}
			i--
			continue
		}
		r, size := utf8.DecodeLastRuneInString(sql[:i])
		if r == utf8.RuneError || !isIdentRune(r) {
			break
		}
		i -= size
	}
	return end - i
}

func shouldDeferDefaultDB(defaultName string, prioritizedDBs []string) bool {
	if defaultName == "" {
		return false
	}
	for _, db := range prioritizedDBs {
		if db == defaultName {
			return false
		}
		if strings.EqualFold(db, defaultName) {
			return true
		}
	}
	return false
}

func appendUniqueDB(dst []string, dbs ...string) []string {
	for _, db := range dbs {
		if db == "" || slices.Contains(dst, db) {
			continue
		}
		dst = append(dst, db)
	}
	return dst
}

// LineOffsetToBytePos converts a 1-based line and 0-based UTF-16 code-unit
// column offset into a byte offset into s, matching LSP/Monaco positions.
func LineOffsetToBytePos(s string, line, offset int) int {
	currentLine := 1
	col := 0
	for i := 0; i < len(s); {
		if currentLine == line && col >= offset {
			return i
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\n' {
			if currentLine == line {
				return i
			}
			currentLine++
			col = 0
		} else if r <= 0xFFFF {
			col++
		} else {
			col += 2
		}
		i += size
	}
	return len(s)
}

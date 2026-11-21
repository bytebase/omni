package completion

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/catalog"
	"github.com/bytebase/omni/mariadb/parser"
	metadatapb "github.com/bytebase/omni/metadata"
)

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

// CompleteSession runs session-aware auto-completion over statement at the
// 1-based line and 0-based UTF-16 code-unit caretOffset.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	stmt, pos := currentStatement(statement, bytePos)
	pos = max(0, min(pos, len(stmt)))

	cat := catalog.New()
	populateCatalog(
		ctx,
		cfg,
		stmt,
		pos,
		hasDatabaseRule,
		func(name string) {
			_ = cat.DefineDatabase(&ast.CreateDatabaseStmt{Name: name, IfNotExists: true})
		},
		func(name string, meta *metadatapb.DatabaseSchemaMetadata) {
			_, _ = cat.LoadMetadataInto(ctx, meta, name)
		},
		cat.SetCurrentDatabase,
	)

	return postProcessSessionCandidates(stmt, pos, cfg.ReadOnly, Complete(stmt, pos, cat), hasFunctionInvocationRule)
}

func postProcessSessionCandidates(
	stmt string,
	pos int,
	readOnly bool,
	raw []Candidate,
	isFuncInvocation func(stmt string, pos int) bool,
) []Candidate {
	caretInBacktick := CaretInsideBacktickIdentifier(stmt, pos)
	_, _, dbQualified := identBeforeDot(stmt, max(0, pos-extractIdentPrefixLen(stmt, pos)))
	var funcInvocationChecked bool
	var appendFuncParens bool
	shouldAppendFuncParens := func() bool {
		if !funcInvocationChecked {
			funcInvocationChecked = true
			appendFuncParens = !caretInBacktick && !hasFollowingOpenParen(stmt, pos) && (isFuncInvocation == nil || isFuncInvocation(stmt, pos))
		}
		return appendFuncParens
	}

	type candKey struct {
		typ        CandidateType
		text       string
		definition string
	}
	candidateMap := make(map[candKey]Candidate, len(raw))
	for _, c := range raw {
		text := c.Text
		switch {
		case isObjectIdentifierCandidate(c.Type):
			text = QuoteIdentifierIfNeeded(c.Text, caretInBacktick)
		case c.Type == CandidateFunction:
			// Preserve built-in function syntax unquoted (several built-ins such
			// as IF/LEFT/REPLACE are reserved keywords); quote catalog stored
			// functions using the complete raw identifier before attaching call parens.
			if caretInBacktick || dbQualified || !isBuiltinFunction(text) {
				text = QuoteIdentifierIfNeeded(text, caretInBacktick)
			}
			if shouldAppendFuncParens() {
				text += "()"
			}
		default:
			// Keywords and value-like candidates pass through unchanged.
		}
		candidate := Candidate{
			Type:       c.Type,
			Text:       text,
			Definition: c.Definition,
			Comment:    c.Comment,
		}
		candidateMap[candKey{typ: c.Type, text: text, definition: c.Definition}] = candidate
	}

	// In read-only mode, drop keywords that initiate writes at statement start
	// while keeping the same words when valid as sub-keywords (e.g. SHOW CREATE
	// TABLE, SELECT ... FOR UPDATE).
	if readOnly && caretAtStatementStart(stmt, pos) {
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

func populateCatalog(
	ctx context.Context,
	cfg SessionConfig,
	statement string,
	pos int,
	hasDBRule func(sql string, offset int) bool,
	registerDB func(name string),
	loadDB func(name string, meta *metadatapb.DatabaseSchemaMetadata),
	setCurrentDB func(name string),
) {
	suggestsDB := hasDBRule != nil && suggestsDatabaseRef(statement, pos, hasDBRule)

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
	// Omni's catalog stores databases under lowercased keys and merges
	// subsequent LoadMetadataInto calls targeting the same lowercased key.
	// Track loaded lowercased catalog slots so case-colliding siblings (e.g.
	// `Foo` vs `foo`) never merge their objects into the prioritized database.
	loadedCatalogSlots := map[string]bool{}
	keyFor := func(name string) string {
		if caseSensitive {
			return name
		}
		return strings.ToLower(name)
	}
	loadOnceIntoCatalog := func(name string, meta *metadatapb.DatabaseSchemaMetadata) {
		slot := strings.ToLower(name)
		if loadedCatalogSlots[slot] {
			return
		}
		loadedCatalogSlots[slot] = true
		if meta != nil {
			loadDB(name, meta)
		} else {
			registerDB(name)
		}
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

	qualifiers, caretDBs := analyzeStatementQualifiers(statement, pos, cfg.DefaultDatabase, caseSensitive)

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
		loadOnceIntoCatalog(defaultName, defaultProto)
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
			if caseSensitive && loadedCatalogSlots[strings.ToLower(loadName)] {
				seen[keyFor(loadName)] = true
				seen[keyFor(q)] = true
				continue
			}
			var proto *metadatapb.DatabaseSchemaMetadata
			var found bool
			if resolvedName, p, _, ok := cfg.LoadDatabase(ctx, loadName); ok {
				if resolvedName != "" {
					loadName = resolvedName
				}
				proto = p
				found = true
			} else if !defaultLoaded || !caseSensitive {
				if exact, diffCase := findExactFolded(loadName); diffCase {
					if resolvedName, p, cs, ok := cfg.LoadDatabase(ctx, exact); ok && !cs {
						loadName = exact
						if resolvedName != "" {
							loadName = resolvedName
						}
						proto = p
						found = true
					}
				}
			}
			seen[keyFor(loadName)] = true
			seen[keyFor(q)] = true
			if found && proto != nil {
				loadOnceIntoCatalog(loadName, proto)
			} else if suggestsDB {
				loadOnceIntoCatalog(loadName, nil)
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

var builtinFunctionSet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(builtinFunctions))
	for _, name := range builtinFunctions {
		set[strings.ToUpper(name)] = struct{}{}
	}
	return set
}()

func isBuiltinFunction(name string) bool {
	_, ok := builtinFunctionSet[strings.ToUpper(name)]
	return ok
}

func isReservedKeyword(name string) bool {
	return parser.IsReservedKeyword(name)
}

func hasDatabaseRule(sql string, offset int) bool {
	cs := parser.Collect(sql, offset)
	return cs != nil && cs.HasRule("database_ref")
}

func hasFunctionInvocationRule(stmt string, pos int) bool {
	return suggestsRule(stmt, pos, func(sql string, offset int) bool {
		cs := parser.Collect(sql, offset)
		return cs != nil && cs.HasRule("func_name")
	})
}

func currentStatement(statement string, pos int) (string, int) {
	for _, seg := range parser.Split(statement) {
		if pos >= seg.ByteStart && pos <= seg.ByteEnd {
			return seg.Text, pos - seg.ByteStart
		}
	}
	return "", 0
}

func suggestsDatabaseRef(stmt string, pos int, hasDBRule func(sql string, offset int) bool) bool {
	return suggestsRule(stmt, pos, hasDBRule)
}

func suggestsRule(stmt string, pos int, hasRule func(sql string, offset int) bool) bool {
	pos = max(0, min(pos, len(stmt)))
	collectOffset := pos - extractIdentPrefixLen(stmt, pos)
	if hasRule(stmt, collectOffset) {
		return true
	}
	prefix := stmt[:collectOffset]
	suffix := stmt[collectOffset:]
	for _, patched := range []string{
		prefix + " __placeholder__" + suffix,
		prefix + " __placeholder__ " + suffix,
		prefix + " 1" + suffix,
	} {
		if hasRule(patched, collectOffset) {
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

func hasFollowingOpenParen(sql string, pos int) bool {
	i := max(0, min(pos, len(sql)))
	for i < len(sql) {
		sz, isIdent := nextIdentRune(sql, i)
		if !isIdent {
			break
		}
		i += sz
	}
	for i < len(sql) && isSpaceByte(sql[i]) {
		i++
	}
	return i < len(sql) && sql[i] == '('
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

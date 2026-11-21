package mysqlcomplete

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type qualToken struct {
	text   string
	upper  string // non-empty only for unquoted word tokens
	punct  byte   // non-zero for '.', ',', '(', ')', ';', '*'
	isWord bool
	end    int
}

func scanQualTokens(statement string) []qualToken {
	s := StripStringsAndComments(statement)
	var toks []qualToken
	n := len(s)
	for i := 0; i < n; {
		b := s[i]
		if isSpaceByte(b) {
			i++
			continue
		}
		switch b {
		case '.', ',', '(', ')', ';', '*':
			i++
			toks = append(toks, qualToken{punct: b, end: i})
		case '`':
			var buf strings.Builder
			i++
			closed := false
			for i < n {
				if s[i] == '`' {
					if i+1 < n && s[i+1] == '`' {
						buf.WriteByte('`')
						i += 2
						continue
					}
					closed = true
					i++
					break
				}
				buf.WriteByte(s[i])
				i++
			}
			if closed && buf.Len() > 0 {
				toks = append(toks, qualToken{text: buf.String(), isWord: true, end: i})
			}
		default:
			if size, ok := nextIdentRune(s, i); ok {
				start := i
				i += size
				for i < n {
					sz, isIdent := nextIdentRune(s, i)
					if !isIdent {
						break
					}
					i += sz
				}
				word := s[start:i]
				toks = append(toks, qualToken{
					text:   word,
					upper:  strings.ToUpper(word),
					isWord: true,
					end:    i,
				})
			} else {
				// Emit a neutral separator token for operators/symbols so adjacent
				// dotted expressions (e.g. `ON a. = o.user_id`) never coalesce into
				// a synthetic three-part chain (`a . o .`).
				i += size
				if len(toks) == 0 || toks[len(toks)-1].punct != ' ' {
					toks = append(toks, qualToken{punct: ' ', end: i})
				} else {
					toks[len(toks)-1].end = i
				}
			}
		}
	}
	return toks
}

type projRef struct {
	db    string
	alias string
}

type qualScope struct {
	aliasScopeID          int
	cteName               string
	stmtStart             bool
	stmtVerb              string
	inSelectList          bool
	pastFromOrWhere       bool
	isSubquery            bool
	inTableClause         bool
	inInsertReplaceTarget bool
	inGrantPrivList       bool
	allowCommaTable       bool
	allowFromAsTable      bool
	afterCreate           bool
	afterRename           bool
	afterDDLOnObj         bool
	prevKeyword           string
}

type aliasMap struct {
	exact map[string][]string
	lower map[string][]string
}

// DatabaseQualifiersInStatement returns the non-default database qualifiers
// referenced in statement.
func DatabaseQualifiersInStatement(statement, defaultDB string, caseSensitive bool) []string {
	qualifiers, _ := AnalyzeStatementQualifiers(statement, 0, defaultDB, caseSensitive)
	return qualifiers
}

// AnalyzeStatementQualifiers extracts database qualifiers referenced in
// statement (excluding defaultDB), distinguishing them from ordinary
// table/alias qualifiers such as `SELECT a. FROM users a`, and resolves the
// databases associated with any qualifier or table/alias active at pos:
//   - Three-part chains (`db.tbl.col` or `db.tbl.`) always identify `db`.
//   - Two-part chains (`db.obj` or `db.`) only identify `db` inside table/object
//     positions (FROM, JOIN, UPDATE, INSERT/REPLACE [INTO], DELETE, USING,
//     DDL object headers, REFERENCES, CALL, CREATE ... LIKE, etc.).
func AnalyzeStatementQualifiers(statement string, pos int, defaultDB string, caseSensitive bool) ([]string, []string) {
	if !strings.Contains(statement, ".") {
		return nil, nil
	}
	toks := scanQualTokens(statement)
	if len(toks) == 0 {
		return nil, nil
	}

	keyFor := func(name string) string {
		if caseSensitive {
			return name
		}
		return strings.ToLower(name)
	}
	defaultKey := keyFor(defaultDB)
	seen := make(map[string]bool)
	var out []string
	addDB := func(name string) {
		if name == "" {
			return
		}
		key := keyFor(name)
		if key == defaultKey || seen[key] {
			return
		}
		seen[key] = true
		out = append(out, name)
	}

	nextAliasScopeID := 1
	aliasesByScope := map[int]*aliasMap{}
	parentAliasScope := map[int]int{}
	scopeIsQueryBlock := map[int]bool{0: true}
	scopeUsesDefaultUnqualified := map[int]bool{}
	scopeProjRefs := map[int][]projRef{}
	scopeFromDBs := map[int][]string{}
	cteByScope := map[int]map[string][]string{}
	recordProjRef := func(scopeID int, ref projRef) {
		if ref.db == "" && ref.alias == "" {
			return
		}
		scopeProjRefs[scopeID] = append(scopeProjRefs[scopeID], ref)
	}
	recordScopeFromDB := func(scopeID int, dbs ...string) {
		scopeFromDBs[scopeID] = appendUniqueDB(scopeFromDBs[scopeID], dbs...)
	}
	resolveAliasDBs := func(scopeID int, ident string) ([]string, bool) {
		lower := strings.ToLower(ident)
		for sid, ok := scopeID, true; ok; sid, ok = parentAliasScope[sid] {
			if am := aliasesByScope[sid]; am != nil {
				if dbs, hit := am.exact[ident]; hit {
					return dbs, true
				}
				if dbs, hit := am.lower[lower]; hit {
					return dbs, true
				}
			}
		}
		return nil, false
	}
	finalizeScopeDBs := func(scopeID int) ([]string, []string) {
		var projDBs []string
		for _, ref := range scopeProjRefs[scopeID] {
			if ref.db != "" {
				projDBs = appendUniqueDB(projDBs, ref.db)
				continue
			}
			if dbs, ok := resolveAliasDBs(scopeID, ref.alias); ok {
				projDBs = appendUniqueDB(projDBs, dbs...)
			}
		}
		allDBs := appendUniqueDB(slices.Clone(projDBs), scopeFromDBs[scopeID]...)
		return projDBs, allDBs
	}
	recordCTE := func(scopeID int, name string, dbs []string) {
		if name == "" {
			return
		}
		m := cteByScope[scopeID]
		if m == nil {
			m = make(map[string][]string)
			cteByScope[scopeID] = m
		}
		lower := strings.ToLower(name)
		m[lower] = appendUniqueDB(m[lower], dbs...)
	}
	lookupCTE := func(scopeID int, name string) ([]string, bool) {
		lower := strings.ToLower(name)
		for sid, ok := scopeID, true; ok; sid, ok = parentAliasScope[sid] {
			if m := cteByScope[sid]; m != nil {
				if dbs, hit := m[lower]; hit {
					return dbs, true
				}
			}
		}
		return nil, false
	}
	recordAlias := func(scopeID int, name string, overwrite bool, dbs ...string) {
		if name == "" || len(dbs) == 0 {
			return
		}
		am := aliasesByScope[scopeID]
		if am == nil {
			am = &aliasMap{
				exact: make(map[string][]string),
				lower: make(map[string][]string),
			}
			aliasesByScope[scopeID] = am
		}
		lower := strings.ToLower(name)
		if overwrite {
			am.exact[name] = appendUniqueDB(nil, dbs...)
			am.lower[lower] = appendUniqueDB(nil, dbs...)
			return
		}
		if _, ok := am.exact[name]; !ok {
			am.exact[name] = appendUniqueDB(nil, dbs...)
		}
		if _, ok := am.lower[lower]; !ok {
			am.lower[lower] = appendUniqueDB(nil, dbs...)
		}
	}

	scopes := []qualScope{{aliasScopeID: 0, stmtStart: true, allowFromAsTable: true}}
	cur := func() *qualScope { return &scopes[len(scopes)-1] }
	caretAliasScopeID := 0

	for i := 0; i < len(toks); i++ {
		tok := toks[i]
		sc := cur()
		prevPunct := byte(0)
		if i > 0 {
			prevPunct = toks[i-1].punct
		}

		switch tok.punct {
		case ';':
			aliasID := nextAliasScopeID
			nextAliasScopeID++
			scopeIsQueryBlock[aliasID] = true
			scopes = []qualScope{{aliasScopeID: aliasID, stmtStart: true, allowFromAsTable: true}}
			if tok.end <= pos {
				caretAliasScopeID = aliasID
			}
			continue
		case '(':
			child := qualScope{
				aliasScopeID:    sc.aliasScopeID,
				stmtStart:       true,
				inSelectList:    sc.inSelectList,
				pastFromOrWhere: sc.pastFromOrWhere,
			}
			if sc.stmtVerb == "WITH" {
				child.cteName = cteNameBeforeAS(toks, i)
				if child.cteName != "" {
					recordCTE(sc.aliasScopeID, child.cteName, nil)
				}
			}
			// Parenthesized table-reference / join group, e.g. FROM (db1.t1 JOIN ...),
			// JOIN (db1.t1), UPDATE [LOW_PRIORITY] [IGNORE] (db1.t1 JOIN ...),
			// DELETE ... USING (db1.t1 JOIN ...), or nested ((db1.t1)) / , (db1.t1).
			if sc.inTableClause && isParenthesizedTablePrefix(toks, i, sc) {
				child.inTableClause = true
				child.allowCommaTable = true
				child.allowFromAsTable = true
				child.inSelectList = false
			} else {
				aliasID := nextAliasScopeID
				nextAliasScopeID++
				parentAliasScope[aliasID] = sc.aliasScopeID
				child.aliasScopeID = aliasID
			}
			scopes = append(scopes, child)
			if tok.end <= pos {
				caretAliasScopeID = child.aliasScopeID
			}
			continue
		case ')':
			var closed qualScope
			if len(scopes) > 1 {
				closed = scopes[len(scopes)-1]
				scopes = scopes[:len(scopes)-1]
			}
			parent := cur()
			if closed.aliasScopeID != parent.aliasScopeID && !closed.isSubquery && len(scopeFromDBs[closed.aliasScopeID]) == 0 {
				// Expression parentheses inside SELECT (e.g. COALESCE(y.col, 0))
				// defer alias resolution to the enclosing query block's FROM/JOINs.
				if parent.inSelectList {
					for _, ref := range scopeProjRefs[closed.aliasScopeID] {
						recordProjRef(parent.aliasScopeID, ref)
					}
				}
			} else {
				subProjDBs, subDBs := finalizeScopeDBs(closed.aliasScopeID)
				if closed.cteName != "" {
					recordCTE(parent.aliasScopeID, closed.cteName, subDBs)
					recordAlias(parent.aliasScopeID, closed.cteName, true, subDBs...)
				}
				if len(subDBs) > 0 {
					if closed.aliasScopeID != parent.aliasScopeID {
						if parent.inSelectList {
							for _, db := range subDBs {
								recordProjRef(parent.aliasScopeID, projRef{db: db})
							}
							recordScopeFromDB(parent.aliasScopeID, scopeFromDBs[closed.aliasScopeID]...)
						} else if !parent.inTableClause && !parent.pastFromOrWhere && len(scopeProjRefs[parent.aliasScopeID]) == 0 && len(scopeFromDBs[parent.aliasScopeID]) == 0 {
							for _, db := range subProjDBs {
								recordProjRef(parent.aliasScopeID, projRef{db: db})
							}
							recordScopeFromDB(parent.aliasScopeID, scopeFromDBs[closed.aliasScopeID]...)
							parent.isSubquery = parent.isSubquery || closed.isSubquery
						} else {
							recordScopeFromDB(parent.aliasScopeID, subDBs...)
						}
					}
					if parent.inTableClause {
						if aliasTok, ok := scanTableAliasAfter(toks, i+1); ok {
							recordAlias(parent.aliasScopeID, aliasTok, true, subDBs...)
						}
					}
				}
			}
			parent.prevKeyword = ""
			parent.stmtStart = false
			if tok.end <= pos {
				caretAliasScopeID = parent.aliasScopeID
			}
			continue
		case ',':
			sc.inTableClause = sc.allowCommaTable
			sc.inInsertReplaceTarget = false
			sc.prevKeyword = ""
			sc.stmtStart = false
			if tok.end <= pos {
				caretAliasScopeID = sc.aliasScopeID
			}
			continue
		case '.':
			sc.prevKeyword = ""
			sc.stmtStart = false
			if tok.end <= pos {
				caretAliasScopeID = sc.aliasScopeID
			}
			continue
		default:
			// Neutral separator or zero punct falls through to word handling.
		}

		if tok.end <= pos {
			caretAliasScopeID = sc.aliasScopeID
		}

		if !tok.isWord {
			continue
		}

		// Record database qualifiers only at the head of a dotted identifier
		// chain (not when preceded by '.'). In DELETE target lists before FROM,
		// ignore `alias.*` wildcards while keeping three-part `db.tbl.*` and `db.`.
		precededByDot := i > 0 && toks[i-1].punct == '.'
		if !precededByDot && i+1 < len(toks) && toks[i+1].punct == '.' {
			isThreePart := i+3 < len(toks) && toks[i+2].isWord && !isClauseBoundaryAfterDot(toks[i+2].upper) && toks[i+3].punct == '.'
			isDeleteStar := sc.stmtVerb == "DELETE" && i+2 < len(toks) && toks[i+2].punct == '*'
			if isThreePart || (sc.inTableClause && !isDeleteStar) {
				addDB(tok.text)
				if sc.inTableClause {
					recordScopeFromDB(sc.aliasScopeID, tok.text)
				} else if sc.inSelectList {
					recordProjRef(sc.aliasScopeID, projRef{db: tok.text})
				} else {
					recordScopeFromDB(sc.aliasScopeID, tok.text)
				}
			} else if sc.inSelectList && !isDeleteStar {
				recordProjRef(sc.aliasScopeID, projRef{alias: tok.text})
			}
			if sc.inTableClause && !isDeleteStar && i+2 < len(toks) && toks[i+2].isWord && !isClauseBoundaryAfterDot(toks[i+2].upper) {
				recordAlias(sc.aliasScopeID, toks[i+2].text, true, tok.text)
				if aliasTok, ok := scanTableAliasAfter(toks, i+3); ok {
					recordAlias(sc.aliasScopeID, aliasTok, true, tok.text)
				}
			}
		} else if sc.inTableClause && !precededByDot && (i+1 >= len(toks) || (toks[i+1].punct != '.' && toks[i+1].punct != '(')) &&
			isTableStartPosition(sc.prevKeyword, prevPunct) && !isNonAliasClauseKeyword(tok.upper) {
			targetDBs, fromCTE := lookupCTE(sc.aliasScopeID, tok.text)
			if !fromCTE && defaultDB != "" {
				scopeUsesDefaultUnqualified[sc.aliasScopeID] = true
				targetDBs = []string{defaultDB}
			}
			if len(targetDBs) > 0 {
				recordScopeFromDB(sc.aliasScopeID, targetDBs...)
				recordAlias(sc.aliasScopeID, tok.text, fromCTE, targetDBs...)
				if aliasTok, ok := scanTableAliasAfter(toks, i+1); ok {
					recordAlias(sc.aliasScopeID, aliasTok, fromCTE, targetDBs...)
				}
			}
		}

		// Keyword-driven clause tracking. Backtick-quoted words have tok.upper == "".
		// An unquoted member after '.' (e.g. SELECT u. FROM otherdb.users u) only
		// triggers clause transitions when it is a reserved clause boundary.
		if tok.upper == "" || (precededByDot && !isClauseBoundaryAfterDot(tok.upper)) {
			sc.prevKeyword = ""
			sc.stmtStart = false
			continue
		}

		kw := tok.upper
		wasStmtStart := sc.stmtStart
		if wasStmtStart {
			sc.stmtVerb = kw
		} else if sc.stmtVerb == "EXPLAIN" {
			switch kw {
			case "SELECT", "UPDATE", "DELETE", "INSERT", "REPLACE", "TABLE":
				sc.stmtVerb = kw
			default:
			}
		}
		prevKW := sc.prevKeyword
		sc.stmtStart = false
		sc.prevKeyword = kw

		// While scanning privilege names/clauses before ON in GRANT/REVOKE, keep
		// privilege keywords (e.g. SELECT, INSERT, UPDATE, TRIGGER, REFERENCES)
		// from opening a table clause until ON arrives.
		if sc.inGrantPrivList {
			if kw == "ON" {
				sc.inGrantPrivList = false
				sc.inTableClause = true
				sc.allowCommaTable = false
				sc.afterDDLOnObj = false
			}
			continue
		}

		switch kw {
		case "SELECT", "WITH":
			if len(scopes) > 1 {
				// Derived tables / subqueries isolate their table aliases from outer scopes.
				if _, ok := parentAliasScope[sc.aliasScopeID]; !ok && sc.aliasScopeID == scopes[len(scopes)-2].aliasScopeID {
					aliasID := nextAliasScopeID
					nextAliasScopeID++
					parentAliasScope[aliasID] = sc.aliasScopeID
					sc.aliasScopeID = aliasID
					if tok.end <= pos {
						caretAliasScopeID = aliasID
					}
				}
			}
			scopeIsQueryBlock[sc.aliasScopeID] = true
			sc.inSelectList = kw == "SELECT"
			sc.pastFromOrWhere = false
			if kw == "SELECT" {
				sc.isSubquery = true
			}
			sc.inTableClause = false
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
			sc.allowFromAsTable = true
			sc.afterCreate = false
			sc.afterRename = false
			sc.afterDDLOnObj = false
		case "WHERE", "SET", "VALUES", "VALUE", "GROUP", "HAVING", "ORDER", "LIMIT",
			"UNION", "INTERSECT", "EXCEPT", "FOR", "WINDOW", "RETURNING", "RETURN":
			if sc.isSubquery || len(scopes) == 1 {
				sc.inSelectList = false
				sc.pastFromOrWhere = true
			}
			sc.inTableClause = false
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
		case "ON":
			sc.inTableClause = sc.afterDDLOnObj
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
			sc.afterDDLOnObj = false
		case "FROM":
			if sc.isSubquery || len(scopes) == 1 {
				sc.inSelectList = false
				sc.pastFromOrWhere = true
			}
			if sc.stmtVerb == "REVOKE" {
				sc.inTableClause = false
				sc.allowCommaTable = false
			} else if sc.allowFromAsTable || sc.stmtVerb == "SHOW" || sc.stmtVerb == "DELETE" || prevKW == "COLUMNS" || prevKW == "FIELDS" || prevKW == "INDEX" || prevKW == "INDEXES" || prevKW == "KEYS" {
				sc.inTableClause = true
				sc.allowCommaTable = true
				sc.inInsertReplaceTarget = false
			}
		case "JOIN", "STRAIGHT_JOIN":
			sc.inTableClause = true
			sc.allowCommaTable = true
			sc.inInsertReplaceTarget = false
		case "UPDATE":
			if prevKW != "ON" && prevKW != "KEY" && prevKW != "FOR" && prevKW != "BEFORE" && prevKW != "AFTER" {
				sc.stmtVerb = "UPDATE"
				sc.inTableClause = true
				sc.allowCommaTable = true
				sc.inInsertReplaceTarget = false
			}
		case "DELETE":
			if prevKW != "ON" && prevKW != "BEFORE" && prevKW != "AFTER" {
				sc.stmtVerb = "DELETE"
				sc.inTableClause = true
				sc.allowCommaTable = true
				sc.inInsertReplaceTarget = false
			}
		case "INSERT", "REPLACE":
			sc.inTableClause = true
			sc.inInsertReplaceTarget = true
			sc.allowCommaTable = false
		case "INTO":
			if sc.inInsertReplaceTarget || sc.stmtVerb == "IMPORT" {
				sc.inTableClause = true
				sc.allowCommaTable = false
				sc.inInsertReplaceTarget = false
			} else {
				if sc.isSubquery || len(scopes) == 1 {
					sc.inSelectList = false
					sc.pastFromOrWhere = true
				}
				sc.inTableClause = false
				sc.allowCommaTable = false
			}
		case "USING":
			sc.inInsertReplaceTarget = false
			if sc.stmtVerb != "DELETE" && i+1 < len(toks) && toks[i+1].punct == '(' {
				sc.inTableClause = false
				sc.allowCommaTable = false
			} else {
				sc.inTableClause = true
				sc.allowCommaTable = true
			}
		case "CREATE":
			sc.afterCreate = true
			sc.inTableClause = false
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
		case "RENAME":
			sc.afterRename = true
			sc.inTableClause = false
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
		case "GRANT", "REVOKE":
			sc.inGrantPrivList = true
			sc.inTableClause = false
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
			sc.allowFromAsTable = false
		case "INDEX", "TRIGGER":
			sc.inInsertReplaceTarget = false
			if kw == "INDEX" && sc.stmtVerb == "CACHE" && prevKW == "CACHE" {
				sc.inTableClause = true
				sc.allowCommaTable = true
				sc.afterDDLOnObj = false
			} else {
				sc.afterDDLOnObj = true
				if kw == "TRIGGER" {
					sc.inTableClause = true
					sc.allowCommaTable = false
				}
			}
		case "CACHE":
			if sc.stmtVerb == "LOAD" && prevKW == "INTO" {
				sc.inTableClause = true
				sc.allowCommaTable = true
				sc.inInsertReplaceTarget = false
			}
		case "TABLE", "TABLES":
			sc.inTableClause = true
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = isMultiTableStatementVerb(sc.stmtVerb)
		case "VIEW":
			sc.inTableClause = true
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = sc.stmtVerb == "DROP"
		case "SEQUENCE":
			sc.inTableClause = true
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = sc.stmtVerb == "DROP"
		case "TRUNCATE":
			if (wasStmtStart || prevKW == "BEGIN") && (i+1 >= len(toks) || toks[i+1].punct != '(') {
				sc.inTableClause = true
				sc.inInsertReplaceTarget = false
				sc.allowCommaTable = false
			}
		case "EVENT", "FUNCTION", "PROCEDURE", "CALL", "HANDLER", "REFERENCES":
			sc.inTableClause = true
			sc.inInsertReplaceTarget = false
			sc.allowCommaTable = false
		case "DESC", "DESCRIBE", "EXPLAIN":
			if wasStmtStart {
				sc.inTableClause = true
				sc.allowCommaTable = false
			}
		case "TO":
			if sc.afterRename {
				sc.inTableClause = true
			} else {
				sc.inTableClause = false
				sc.allowCommaTable = false
			}
		case "LIKE":
			if sc.afterCreate {
				sc.inTableClause = true
			}
		case "IN":
			if sc.stmtVerb == "CACHE" {
				sc.inTableClause = false
				sc.allowCommaTable = false
			} else if prevKW == "COLUMNS" || prevKW == "FIELDS" || prevKW == "INDEX" || prevKW == "INDEXES" || prevKW == "KEYS" {
				sc.inTableClause = true
				sc.allowCommaTable = false
			}
		default:
			// Non-clause keywords leave table-clause state unchanged.
		}
	}

	var caretDBs []string
	if caretIdent := databaseQualifierAtCaret(statement, pos); caretIdent != "" {
		if caretIdent == defaultDB {
			caretDBs = []string{defaultDB}
		} else {
			for _, q := range out {
				if q == caretIdent {
					caretDBs = []string{q}
					break
				}
			}
		}
		if len(caretDBs) == 0 {
			if dbs, hit := resolveAliasDBs(caretAliasScopeID, caretIdent); hit {
				caretDBs = dbs
			}
		}
	} else {
		scopeID := caretAliasScopeID
		for scopeID != 0 && !scopeIsQueryBlock[scopeID] {
			parentID, ok := parentAliasScope[scopeID]
			if !ok {
				break
			}
			scopeID = parentID
		}
		if !scopeUsesDefaultUnqualified[scopeID] {
			if _, subDBs := finalizeScopeDBs(scopeID); len(subDBs) > 0 {
				caretDBs = subDBs
			}
		}
	}
	return out, caretDBs
}

// StatementReferencesDatabase reports whether dbName appears as a database
// qualifier (`dbName.`) in statement, at an identifier boundary.
func StatementReferencesDatabase(statement, dbName string) bool {
	if dbName == "" {
		return false
	}
	s := strings.ToLower(StripStringsAndComments(statement))
	name := strings.ToLower(dbName)
	for _, q := range []string{name, "`" + name + "`"} {
		from := 0
		for {
			i := strings.Index(s[from:], q)
			if i < 0 {
				break
			}
			at := from + i
			from = at + 1
			if q[0] != '`' && at > 0 {
				prev, _ := utf8.DecodeLastRuneInString(s[:at])
				if isIdentRune(prev) {
					continue
				}
			}
			j := at + len(q)
			for j < len(s) && isSpaceByte(s[j]) {
				j++
			}
			if j < len(s) && s[j] == '.' {
				return true
			}
		}
	}
	return false
}

func isMultiTableStatementVerb(verb string) bool {
	switch verb {
	case "DROP", "RENAME", "LOCK", "TRUNCATE", "ANALYZE", "CHECK", "CHECKSUM", "OPTIMIZE", "REPAIR", "FLUSH":
		return true
	default:
		return false
	}
}

// databaseQualifierAtCaret returns the leading database identifier of a dotted
// qualifier chain active at pos (e.g. `Foo` for `Foo.t.` or `Foo.`), or "".
func databaseQualifierAtCaret(statement string, pos int) string {
	pos = max(0, min(pos, len(statement)))
	prefix := StripStringsAndComments(statement[:pos])
	toks := scanQualTokens(prefix)
	if len(toks) < 2 {
		return ""
	}
	trimmed := strings.TrimRight(prefix, " \t\r\n")
	if trimmed == "" {
		return ""
	}
	lastIdx := len(toks) - 1
	if trimmed[len(trimmed)-1] == '.' {
		if toks[lastIdx].punct != '.' {
			return ""
		}
		lastIdx--
	} else if isSpaceByte(prefix[len(prefix)-1]) || !toks[lastIdx].isWord || toks[lastIdx-1].punct != '.' {
		return ""
	}
	if lastIdx < 0 || !toks[lastIdx].isWord {
		return ""
	}
	for steps := 0; steps < 2 && lastIdx >= 2 && toks[lastIdx-1].punct == '.' && toks[lastIdx-2].isWord; steps++ {
		lastIdx -= 2
	}
	return toks[lastIdx].text
}

func isParenthesizedTablePrefix(toks []qualToken, parenIdx int, sc *qualScope) bool {
	if parenIdx == 0 {
		return true
	}
	prev := toks[parenIdx-1]
	if prev.punct == '(' || prev.punct == ',' {
		return true
	}
	switch prev.upper {
	case "FROM", "JOIN", "STRAIGHT_JOIN", "UPDATE":
		return true
	case "USING":
		return sc.stmtVerb == "DELETE"
	case "LOW_PRIORITY", "IGNORE":
		return sc.stmtVerb == "UPDATE"
	default:
		return false
	}
}

func isTableStartPosition(prevKW string, prevPunct byte) bool {
	if prevPunct == ',' {
		return true
	}
	switch prevKW {
	case "FROM", "JOIN", "STRAIGHT_JOIN", "UPDATE", "INTO":
		return true
	default:
		return false
	}
}

func scanTableAliasAfter(toks []qualToken, startIdx int) (string, bool) {
	if startIdx >= len(toks) {
		return "", false
	}
	idx := startIdx
	if toks[idx].upper == "AS" {
		idx++
	}
	if idx >= len(toks) || !toks[idx].isWord || isNonAliasClauseKeyword(toks[idx].upper) {
		return "", false
	}
	return toks[idx].text, true
}

func cteNameBeforeAS(toks []qualToken, openParenIdx int) string {
	if openParenIdx < 2 || toks[openParenIdx-1].upper != "AS" {
		return ""
	}
	k := openParenIdx - 2
	if toks[k].punct == ')' {
		depth := 1
		k--
		for k >= 0 && depth > 0 {
			switch toks[k].punct {
			case ')':
				depth++
			case '(':
				depth--
			default:
			}
			k--
		}
	}
	if k >= 0 && toks[k].isWord && toks[k].upper != "WITH" && toks[k].upper != "RECURSIVE" {
		return toks[k].text
	}
	return ""
}

func isNonAliasClauseKeyword(upper string) bool {
	switch upper {
	case "AS", "ON", "USING", "WHERE", "SET", "VALUES", "VALUE",
		"GROUP", "HAVING", "ORDER", "LIMIT", "OFFSET", "FOR", "LOCK",
		"UNION", "INTERSECT", "EXCEPT", "WINDOW", "PARTITION",
		"USE", "IGNORE", "FORCE", "JOIN", "INNER", "LEFT", "RIGHT",
		"FULL", "CROSS", "NATURAL", "STRAIGHT_JOIN", "SELECT", "FROM",
		"INTO", "TO", "RETURNING", "RETURN", "LOW_PRIORITY", "DELAYED",
		"HIGH_PRIORITY", "QUICK", "ALL", "DISTINCT", "WITH", "LATERAL":
		return true
	default:
		return false
	}
}

func isClauseBoundaryAfterDot(upper string) bool {
	switch upper {
	case "SELECT", "WITH", "FROM", "JOIN", "STRAIGHT_JOIN",
		"LEFT", "RIGHT", "INNER", "CROSS", "NATURAL",
		"WHERE", "SET", "VALUES", "GROUP", "HAVING", "ORDER", "LIMIT",
		"UNION", "INTERSECT", "EXCEPT", "ON", "USING", "INTO", "FOR",
		"WINDOW", "RETURNING":
		return true
	default:
		return false
	}
}

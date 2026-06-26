package catalog

import (
	"sort"
	"strings"
)

// MySQL SDL diff — views.
//
// This is the MySQL analog of PG's view differ (pg/catalog/diff_view.go). It compares the view
// set of two catalogs and reports per-view changes via SchemaDiff.Views. It is wired into
// DiffWithNormalizer (diff.go); a view is a database-level object keyed by (database, name), so —
// like diffTables — it takes the whole catalog rather than a single table.
//
// IDEMPOTENCE — the view body is the risk (correctness-protocol.md). MySQL stores a rewritten
// canonical form of the SELECT body in SHOW CREATE VIEW (fully-qualified columns, backtick
// quoting, expanded *, lower-cased function names, parenthesized predicates). The omni loader
// already reproduces most of that form: createView routes the body through deparseViewSelect
// (Resolver → RewriteSelectStmt → DeparseSelect), so a user's `SELECT a,b FROM t` and the
// engine's readback `select ... AS ... from ...` land on the same View.Definition — EXCEPT the
// database-name qualifier. MySQL 8.0 qualifies a same-database column/table reference with the
// database name (`db`.`t`.`a`); the omni deparse preserves whatever the parser produced, so the
// engine readback keeps the `db.` prefix while the user form has none. The two would phantom-diff
// forever. canonicalViewBody folds that prefix away (the view's OWN database only — a genuine
// cross-database reference keeps its qualifier), so the user form and the stored form collapse
// onto one key. This is the one canonicalization this node owns; everything else in the body is
// already canonical out of deparse. (Normalization flag in the PR: the db-qualifier stripping
// ideally belongs in mysql/deparse's resolver so the stored Definition is symmetric; it is done
// here because this node may not edit deparse.)
//
// VERSION DIVERGENCE (5.7 vs 8.0), proven against the live oracle:
//   - Explicit column list: 8.0 stores `CREATE ... VIEW v (c1,c2) AS select e1 AS a1,e2 AS a2`
//     (the list is kept, the SELECT aliases stay the ORIGINAL expression names). 5.7 DROPS the
//     list and instead RENAMES the SELECT aliases to the view column names
//     (`select e1 AS c1,e2 AS c2`). So on 5.7 the explicit column list is folded into the body
//     and View.ExplicitColumns round-trips as false. Comparing the explicit column list is
//     therefore an 8.0-only concern (explicitColumnsComparable); on 5.7 the renamed aliases
//     already live in the body, which the body key compares. The 5.7 alias-rename itself is NOT
//     reproduced by the omni loader (the user's explicit-list form keeps `AS a1`), so a 5.7
//     round-trip of an explicit-column-list view can still differ — a documented, flagged
//     limitation (PR normalization flags), not a silently-wrong empty.
//   - view-on-view: 8.0 qualifies the referenced view with its database, 5.7 does not — both
//     handled by canonicalViewBody stripping the own-database qualifier.
//
// DEFINER is intentionally NOT compared. On the release path the synced `from` view is loaded
// from a metadata SDL dump that carries no DEFINER (createView defaults it to `root`@`%`), while a
// user's target SDL DEFINER is environment-specific; comparing it would phantom-diff every release
// whose live view has a non-root definer. ALGORITHM, SQL SECURITY and WITH CHECK OPTION ARE
// compared — they are stored-form-significant and identical across 5.7/8.0 (oracle-verified).
func diffViews(from, to *Catalog, n *Normalizer) []ViewDiffEntry {
	fromMap := buildViewMap(from)
	toMap := buildViewMap(to)

	var result []ViewDiffEntry

	// Dropped: in from but not in to.
	for key, fromView := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, ViewDiffEntry{
				Action:   DiffDrop,
				Database: key.db,
				Name:     key.name,
				From:     fromView,
			})
		}
	}

	// Added or modified: in to.
	for key, toView := range toMap {
		fromView, ok := fromMap[key]
		if !ok {
			result = append(result, ViewDiffEntry{
				Action:   DiffAdd,
				Database: key.db,
				Name:     key.name,
				To:       toView,
			})
			continue
		}
		if viewsChanged(fromView, toView, n) {
			result = append(result, ViewDiffEntry{
				Action:   DiffModify,
				Database: key.db,
				Name:     key.name,
				From:     fromView,
				To:       toView,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Database != result[j].Database {
			return result[i].Database < result[j].Database
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// buildViewMap indexes every view in a catalog by (database, name) lower-cased, mirroring
// buildTableMap. Tables and views share one namespace in MySQL, but they live in separate
// catalog maps, so a name collision between a table and a view cannot occur within one valid
// catalog; the table differ and the view differ each own their own map.
func buildViewMap(c *Catalog) map[tableKey]*View {
	m := make(map[tableKey]*View)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, v := range db.Views {
			if v == nil {
				continue
			}
			m[tableKey{db: toLower(db.Name), name: name}] = v
		}
	}
	return m
}

// viewsChanged reports whether two same-identity views differ, comparing their canonical keys.
// canonicalView folds every stored-form-significant attribute into one key, so this differ never
// re-implements a per-attribute comparison (mirroring checksChanged / indexesChanged).
func viewsChanged(a, b *View, n *Normalizer) bool {
	return canonicalView(a, n) != canonicalView(b, n)
}

// canonicalView returns a single stable comparison key for a view, folding the attributes that
// survive into MySQL's stored form:
//   - the SELECT body, canonicalized by canonicalViewBody (own-database qualifier stripped);
//   - ALGORITHM (UNDEFINED/MERGE/TEMPTABLE), upper-cased, defaulting to UNDEFINED;
//   - SQL SECURITY (DEFINER/INVOKER), upper-cased, defaulting to DEFINER;
//   - WITH [CASCADED|LOCAL] CHECK OPTION, upper-cased ("" when absent);
//   - the explicit column list, but ONLY on a version that represents it in the stored form
//     (8.0). On 5.7 the list is folded into the body's SELECT aliases (see diffViews), so
//     comparing it there would phantom-diff; explicitColumnsComparable gates it out.
//
// DEFINER is deliberately excluded (ignore-in-diff; see diffViews). The view NAME is the identity
// key (handled by the caller's map), not part of this content key.
func canonicalView(v *View, n *Normalizer) string {
	fields := []string{
		"body", canonicalViewBody(v),
		"algorithm", canonicalViewAlgorithm(v),
		"sqlsecurity", canonicalViewSQLSecurity(v),
		"checkoption", strings.ToUpper(v.CheckOption),
	}
	if explicitColumnsComparable(v, n) {
		fields = append(fields, "columns", strings.Join(lowerAll(v.Columns), ","))
	}
	return encodeKeyFields(fields...)
}

// canonicalViewAlgorithm returns the upper-cased ALGORITHM, defaulting to UNDEFINED (the form
// SHOW CREATE VIEW emits when none was specified). The loader already defaults it on load, so
// this is belt-and-suspenders for a view assembled without going through createView.
func canonicalViewAlgorithm(v *View) string {
	if v.Algorithm == "" {
		return "UNDEFINED"
	}
	return strings.ToUpper(v.Algorithm)
}

// canonicalViewSQLSecurity returns the upper-cased SQL SECURITY, defaulting to DEFINER (the
// SHOW CREATE VIEW default), mirroring canonicalViewAlgorithm.
func canonicalViewSQLSecurity(v *View) string {
	if v.SqlSecurity == "" {
		return "DEFINER"
	}
	return strings.ToUpper(v.SqlSecurity)
}

// explicitColumnsComparable reports whether a view's explicit column list is part of the stored
// form for the target version and is therefore safe to compare directly. It is only on 8.0:
// 8.0 keeps `CREATE ... VIEW v (c1,c2) AS ...` verbatim, so the column list is an independent
// stored attribute. On 5.7 MySQL rewrites the list into the SELECT aliases and drops the clause
// (View.ExplicitColumns becomes false on the readback), so the column identity already lives in
// the body key — comparing the list separately would phantom-diff an explicit-list `to` against
// its alias-folded 5.7 `from`. When neither side has an explicit list, there is nothing to
// compare and this returns false regardless of version.
func explicitColumnsComparable(v *View, n *Normalizer) bool {
	if !v.ExplicitColumns {
		return false
	}
	return n.Version == MySQL80
}

// canonicalViewBody returns the view's SELECT body in a database-qualifier-neutral canonical form:
// the view's OWN database prefix (`db`.) is stripped from every table/column reference, so the
// engine readback's 8.0-style `db`.`t`.`a` and the user form's `t`.`a` (both already produced by
// deparseViewSelect) collapse onto the same string. Only the view's own database is stripped — a
// reference into a different database keeps its qualifier, because that qualifier is semantically
// load-bearing (it names a different object). A view with no Database (defensive; the loader
// always sets it) is returned unchanged.
//
// This is the single canonicalization this node performs on top of deparse. It is a targeted
// fold, not a re-deparse: deparse has already done the structural canonicalization (qualified
// columns, expanded *, function casing, parenthesization); only the redundant own-database prefix
// remains, and stripping it yields a portable, version-symmetric body. (Normalization flag: this
// ideally belongs in mysql/deparse's resolver so View.Definition is symmetric at load time.)
//
// The strip is POSITION-AWARE, not a blind substring replace, to avoid two corruptions:
//   - inside a string literal: a single-quoted literal may contain the text "`db`." (e.g.
//     `SELECT '`vt`.x'`); blind replacement would mutate the literal's value and make two
//     genuinely different views collide. So single-quoted regions are skipped.
//   - a relation literally named the same as the database: `\`vt\`.\`vt\`.\`a\“ (db vt, table vt,
//     column a) must lose only the LEADING `\`vt\`.` (the database slot), yielding
//     `\`vt\`.\`a\“ (table vt, column a) — NOT both, which would destroy the table qualifier. So
//     the prefix is stripped only when its `\`db\“ segment is the HEAD of a dotted chain (it is
//     not itself immediately preceded by a '.', which would make it a table/column segment).
func canonicalViewBody(v *View) string {
	body := v.Definition
	if v.Database == nil || v.Database.Name == "" {
		return body
	}
	prefix := "`" + strings.ReplaceAll(v.Database.Name, "`", "``") + "`."
	return stripOwnDatabaseQualifier(body, prefix)
}

// stripOwnDatabaseQualifier removes every database-qualifier occurrence of prefix ("`db`.") from a
// deparsed SQL body. Two safeguards keep the strip from corrupting the body (see canonicalViewBody):
//   - single-quoted string literals are skipped, so a literal containing the text "`db`." is left
//     intact (different literal values must not collapse);
//   - the prefix is only stripped at the HEAD of a dotted-identifier chain (its leading backtick is
//     not itself immediately preceded by '.'), AND after a strip the immediately-following
//     identifier segment is copied verbatim. The latter is what protects a relation/column literally
//     named the same as the database: in `\`vt\`.\`vt\`.\`a\“ the leading `\`vt\`.` (database) is
//     dropped, then the next `\`vt\“ (the table) is copied untouched, yielding `\`vt\`.\`a\“ — not
//     `\`a\“ (which would destroy the table qualifier).
func stripOwnDatabaseQualifier(body, prefix string) string {
	var b strings.Builder
	b.Grow(len(body))
	inString := false
	prev := byte(0) // last byte emitted (for the "chain head = not preceded by '.'" test)
	i := 0
	for i < len(body) {
		c := body[i]
		if inString {
			// Inside a single-quoted literal: copy verbatim. A backslash escapes the next byte; a
			// lone quote ends the literal (a doubled '' stays inside via the escaped-next-byte path
			// not applying — MySQL's deparse uses backslash escaping, but handle '' defensively).
			b.WriteByte(c)
			if c == '\\' && i+1 < len(body) {
				b.WriteByte(body[i+1])
				i += 2
				continue
			}
			if c == '\'' {
				if i+1 < len(body) && body[i+1] == '\'' { // doubled quote inside the literal
					b.WriteByte(body[i+1])
					i += 2
					continue
				}
				inString = false
			}
			prev = c
			i++
			continue
		}
		if c == '\'' {
			inString = true
			b.WriteByte(c)
			prev = c
			i++
			continue
		}
		// Outside a string: strip the prefix only at a chain head (not preceded by '.').
		if prev != '.' && strings.HasPrefix(body[i:], prefix) {
			i += len(prefix)
			// Copy the following identifier segment (the table) verbatim so a table named like the
			// database is not re-stripped as another chain head.
			i = copyBacktickIdent(&b, body, i)
			prev = '`' // the segment ended with a closing backtick (or nothing was copied)
			continue
		}
		b.WriteByte(c)
		prev = c
		i++
	}
	return b.String()
}

// copyBacktickIdent copies a single backtick-quoted identifier starting at i (if one is there) from
// src into b, returning the index just past it. A backtick inside the identifier is escaped by
// doubling (“); the scan consumes both. If src[i] is not an opening backtick, nothing is copied and
// i is returned unchanged.
func copyBacktickIdent(b *strings.Builder, src string, i int) int {
	if i >= len(src) || src[i] != '`' {
		return i
	}
	b.WriteByte('`')
	i++
	for i < len(src) {
		if src[i] == '`' {
			if i+1 < len(src) && src[i+1] == '`' { // doubled backtick = escaped, stays in the ident
				b.WriteString("``")
				i += 2
				continue
			}
			b.WriteByte('`')
			i++
			break
		}
		b.WriteByte(src[i])
		i++
	}
	return i
}

// lowerAll lower-cases every string in a slice, returning a new slice. Used to make the explicit
// column-list comparison case-insensitive (MySQL view column names are case-insensitive for
// identity, like all MySQL identifiers in the stored form).
func lowerAll(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = toLower(s)
	}
	return out
}

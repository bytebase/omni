package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// MySQL SDL generate — views.
//
// This is the MySQL analog of PG's view generator (pg/catalog/migration_view.go). It turns the
// view diff (SchemaDiff.Views, populated by diff_view.go) into CREATE OR REPLACE VIEW / DROP VIEW
// DDL. It is wired into GenerateMigrationWithNormalizer (migration.go) against the
// OpCreateView/OpDropView op-types at priorityView.
//
// Op mapping:
//   - DiffAdd / DiffModify → CREATE OR REPLACE VIEW. MySQL has no ALTER-for-redefine that covers
//     every attribute (ALTER VIEW exists but CREATE OR REPLACE is the dump/restore form and the
//     one a declarative migration wants — it both creates a new view and atomically redefines an
//     existing one). So an added view and a changed view emit the SAME statement; "view change =
//     CREATE OR REPLACE VIEW".
//   - DiffDrop → DROP VIEW IF EXISTS.
//
// PHASE / ORDERING (apply-correctness, oracle-verified):
//   - CREATE runs in PhaseMain at priorityView (=50), AFTER table creation (priorityTable=10),
//     column ALTERs (priorityColumn=20), and routine CREATEs (priorityRoutineCreate=45). A view
//     over a freshly created/altered table therefore applies against the final table — creating a
//     view whose referenced table does not yet exist fails (live-verified), and creating a view
//     whose body calls a not-yet-created stored function fails with Error 1305 (live-verified on
//     both versions; see migration_routine.go for the routine-side ordering rationale) — so this
//     ordering is required.
//   - view-on-view: a view that references another view being created in the SAME plan must be
//     created AFTER it (the dependency must exist first — live-verified to fail otherwise). The
//     ops all share priorityView, so the intra-priority tie-break (sortName, then stable original
//     order) must encode the dependency order. viewCreateOrder computes a topological depth per
//     view from the in-batch reference graph and prefixes the sortName with a zero-padded depth,
//     so the global stable sort (Phase → Priority → sortName) emits a referenced view before its
//     referrer. CREATE OR REPLACE of an already-existing dependency while a dependent exists does
//     NOT require ordering (live-verified to succeed), so this only matters for brand-new
//     view-on-view chains — but the depth ordering is harmless and uniform, so it is always
//     applied.
//   - DROP runs in PhasePre (before all creates), so a view dropped and re-created in one plan
//     never collides. MySQL permits dropping a view another view still references (the dependent
//     merely becomes invalid — live-verified), so drops need no inter-view ordering; a plain
//     name sort is enough.
//
// RENDERING routes through formatCreateOrReplaceView, which emits the stored SHOW CREATE VIEW form
// (ALGORITHM / SQL SECURITY preamble, optional explicit column list, AS <body>, optional WITH CHECK
// OPTION) minus DEFINER — so the readback of the emitted DDL canonicalizes equal to `to` and
// apply-correctness holds. DEFINER is omitted on purpose (ignore-in-diff + a least-privilege apply
// hazard; see formatCreateOrReplaceView). The body is emitted as stored (v.Definition); MySQL
// re-qualifies a database-unqualified reference against the view's database on apply (live-verified),
// and the diff's canonical comparison — not the emitted text — is what absorbs the db-qualifier
// difference on the readback.
func generateViewDDL(_, _ *Catalog, diff *SchemaDiff, _ *Normalizer) []MigrationOp {
	if diff == nil || len(diff.Views) == 0 {
		return nil
	}

	// Depth per (created/modified) view, encoding view-on-view dependency order so a referenced
	// view's CREATE sorts before its referrer's. Keyed by lower-cased (database, name).
	depth := viewCreateOrder(diff.Views)

	var ops []MigrationOp
	for i := range diff.Views {
		entry := &diff.Views[i]
		switch entry.Action {
		case DiffAdd, DiffModify:
			if entry.To == nil {
				continue
			}
			ops = append(ops, createViewOp(entry, depth))
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, dropViewOp(entry))
		}
	}

	return ops
}

// createViewOp builds a CREATE OR REPLACE VIEW op (PhaseMain, priorityView). The sortName is
// prefixed with the view's topological depth so the global stable sort creates a referenced view
// before its referrer (see generateViewDDL ordering notes).
func createViewOp(entry *ViewDiffEntry, depth map[tableKey]int) MigrationOp {
	v := entry.To
	d := depth[tableKey{db: toLower(entry.Database), name: toLower(entry.Name)}]
	return MigrationOp{
		Type:       OpCreateView,
		Database:   entry.Database,
		ObjectName: entry.Name,
		SQL:        formatCreateOrReplaceView(v),
		Phase:      PhaseMain,
		Priority:   priorityView,
		sortName:   viewCreateSortName(d, entry.Database, entry.Name),
	}
}

// dropViewOp builds a DROP VIEW IF EXISTS op (PhasePre, priorityView). IF EXISTS keeps the drop
// idempotent against a view already absent on the target (defensive; the diff only emits a drop
// for a view present in `from`). Drops run in PhasePre so a drop-then-recreate of the same name
// never collides, and need no inter-view ordering (MySQL allows dropping a referenced view).
func dropViewOp(entry *ViewDiffEntry) MigrationOp {
	return MigrationOp{
		Type:       OpDropView,
		Database:   entry.Database,
		ObjectName: entry.Name,
		SQL:        fmt.Sprintf("DROP VIEW IF EXISTS %s", viewIdent(entry.From, entry.Database, entry.Name)),
		Phase:      PhasePre,
		Priority:   priorityView,
		sortName:   viewSortName(entry.Database, entry.Name),
	}
}

// formatCreateOrReplaceView renders a view's full CREATE OR REPLACE statement in MySQL's canonical
// stored form, mirroring show.go's showCreateView (ALGORITHM / SQL SECURITY preamble, optional
// explicit column list, AS body, optional WITH CHECK OPTION) but with CREATE OR REPLACE and a
// database-qualified view name. Because the rendered statement is the stored form for the target
// version, the readback of the applied DDL canonicalizes equal to `to`.
//
// DEFINER is intentionally NOT emitted. It is ignore-in-diff (see diff_view.go), so emitting it
// would add no declarative meaning — and `CREATE OR REPLACE DEFINER=<other> VIEW` requires the
// applying session to hold SET_USER_ID / SUPER (else MySQL errno 1227), which would fail the whole
// migration on a least-privilege deployment. Omitting the clause makes MySQL default the definer to
// the current user, which is always permitted and re-reads as a definer the diff ignores anyway, so
// apply-correctness still holds. (The loader defaults Definer to `root`@`%` on load; forcing that on
// apply is exactly the privilege hazard we avoid.)
func formatCreateOrReplaceView(v *View) string {
	var b strings.Builder

	b.WriteString("CREATE OR REPLACE")

	algorithm := v.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	fmt.Fprintf(&b, " ALGORITHM=%s", strings.ToUpper(algorithm))

	sqlSecurity := v.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}
	fmt.Fprintf(&b, " SQL SECURITY %s", strings.ToUpper(sqlSecurity))

	fmt.Fprintf(&b, " VIEW %s", viewIdent(v, v.databaseName(), v.Name))

	// Explicit column list (only when the user specified one). On 5.7 the engine folds this into
	// the SELECT aliases and drops the clause; emitting it on a 5.7 target is still accepted
	// (MySQL applies the rename), and the diff's 5.7 handling does not compare the list, so this
	// rendering is safe on both versions.
	if v.ExplicitColumns && len(v.Columns) > 0 {
		cols := make([]string, len(v.Columns))
		for i, c := range v.Columns {
			cols[i] = migrationQuoteIdent(c)
		}
		fmt.Fprintf(&b, " (%s)", strings.Join(cols, ","))
	}

	b.WriteString(" AS ")
	b.WriteString(v.Definition)

	if v.CheckOption != "" {
		fmt.Fprintf(&b, " WITH %s CHECK OPTION", strings.ToUpper(v.CheckOption))
	}

	return b.String()
}

// viewIdent returns the backtick-quoted database.view identifier for migration DDL, using the
// view's ORIGINAL-CASE name and database (mirroring tableIdent for tables). The database
// qualifier is taken from the diff entry's Database when the View carries no Database back-pointer
// (defensive — the loader always sets it). A view with no database resolves to a bare name (the
// synced single-database release path may load with no qualifying database).
func viewIdent(v *View, database, name string) string {
	// Prefer the view's own original-case names; fall back to the diff entry's lower-cased keys.
	dbName := database
	viewName := name
	if v != nil {
		if v.Name != "" {
			viewName = v.Name
		}
		if dn := v.databaseName(); dn != "" {
			dbName = dn
		}
	}
	if dbName == "" {
		return migrationQuoteIdent(viewName)
	}
	return migrationQuoteIdent(dbName) + "." + migrationQuoteIdent(viewName)
}

// databaseName returns the view's database name, or "" if it has no database back-pointer.
func (v *View) databaseName() string {
	if v == nil || v.Database == nil {
		return ""
	}
	return v.Database.Name
}

// viewCreateOrder computes a topological depth for every created/modified view, so that a view
// referencing another view in the same batch sorts after it. Depth 0 = references no other batch
// view (only tables, or views not in this plan); depth d = one more than the deepest batch view it
// references. The depth becomes a zero-padded sortName prefix, so the global stable sort emits
// shallower (referenced) views before deeper (referring) ones.
//
// Dependencies are read from the deparsed body by extracting the relations it references at FROM /
// JOIN positions (referencedRelations), then keeping only edges to OTHER views in this same batch.
// Position-anchored extraction is what makes this precise: a relation appears only after `from `,
// `join `, or an opening `(` (the deparse parenthesizes join arms), whereas a COLUMN is rendered
// `\`tbl\`.\`col\“ and a column ALIAS as `AS \`name\“ — neither sits at a relation anchor — so a
// column or alias that merely shares a batch view's name does NOT produce a false edge (the bug an
// earlier substring scan had). Matching is on the full (database, name) key, so a same-named view
// in a different database is not conflated either.
//
// Self-references and cycles cannot arise in a valid MySQL view set (a view cannot reference itself
// or form a cycle), so the longest-path depth via memoized DFS terminates; a defensive visiting-
// guard breaks any cycle at depth 0 on malformed input rather than recursing forever.
func viewCreateOrder(entries []ViewDiffEntry) map[tableKey]int {
	// The views being created/modified, keyed by identity, with the relation set each references.
	type viewInfo struct {
		key  tableKey
		refs map[tableKey]bool // relations referenced at FROM/JOIN positions
	}
	creates := make([]viewInfo, 0, len(entries))
	idxByKey := make(map[tableKey]int) // identity -> index into creates
	for i := range entries {
		e := &entries[i]
		if e.Action == DiffDrop || e.To == nil {
			continue
		}
		k := tableKey{db: toLower(e.Database), name: toLower(e.Name)}
		idxByKey[k] = len(creates)
		creates = append(creates, viewInfo{key: k, refs: referencedRelations(e.To.Definition, k.db)})
	}

	// Edge i -> j when view i references batch view j at a FROM/JOIN position.
	refIdx := make([][]int, len(creates))
	for i, ci := range creates {
		for j, cj := range creates {
			if i == j {
				continue
			}
			if ci.refs[cj.key] {
				refIdx[i] = append(refIdx[i], j)
			}
		}
	}

	depth := make(map[tableKey]int, len(creates))
	memo := make([]int, len(creates))
	for i := range memo {
		memo[i] = -1
	}
	visiting := make([]bool, len(creates))
	var compute func(i int) int
	compute = func(i int) int {
		if memo[i] >= 0 {
			return memo[i]
		}
		if visiting[i] {
			// Defensive: a cycle (impossible for valid views) — break it at 0.
			return 0
		}
		visiting[i] = true
		best := 0
		for _, j := range refIdx[i] {
			if d := compute(j) + 1; d > best {
				best = d
			}
		}
		visiting[i] = false
		memo[i] = best
		return best
	}
	for i, ci := range creates {
		depth[ci.key] = compute(i)
	}
	return depth
}

// referencedRelations returns the set of (database, name) keys a deparsed view body references as
// real relations (base tables / views), lower-cased, with an unqualified reference resolved against
// ownDB (the view's own database) — matching how the diff keys views, so a same-database reference
// resolves to the batch view's key.
//
// It works on the PARSED AST, not the body text. That is what makes it robust where a text scan is
// not: a CTE name, a column or a column alias that matches a relation name, a relation name inside a
// string literal, and a column reference inside an ON-clause are all correctly excluded by the
// structural walk. Two top-level body shapes are handled:
//   - a SELECT (including UNION/INTERSECT/EXCEPT and parenthesized selects, all *SelectStmt) — refs
//     come from selectTableRefs, the same FROM/JOIN extractor the SDL loader uses, which drops names
//     bound by a top-level CTE;
//   - a `TABLE t` query primary (MySQL 8.0.19+; deparsed as `table \`t\“ and re-parsed to a
//     *TableStmt) — a TABLE primary is `SELECT * FROM t`, so its single Table ref is the relation.
//
// If the body fails to re-parse, or is some other statement shape, the empty set is returned, so the
// view simply carries no view-on-view ordering constraint; that fallback is conservative (a missing
// edge could only matter for a brand-new view-on-view chain).
func referencedRelations(body, ownDB string) map[tableKey]bool {
	out := make(map[tableKey]bool)
	list, err := mysqlparser.Parse(body)
	if err != nil || list == nil || len(list.Items) == 0 {
		return out
	}
	addRef := func(ref *nodes.TableRef) {
		if ref == nil || ref.Name == "" {
			return
		}
		db := ref.Schema
		if db == "" {
			db = ownDB
		}
		out[tableKey{db: toLower(db), name: toLower(ref.Name)}] = true
	}
	switch stmt := list.Items[0].(type) {
	case *nodes.SelectStmt:
		for _, ref := range selectTableRefs(stmt) {
			addRef(ref)
		}
	case *nodes.TableStmt:
		// `TABLE t` is sugar for `SELECT * FROM t`; its single relation is the dependency.
		addRef(stmt.Table)
	}
	return out
}

// viewSortName is the stable secondary sort key for a view-level op: lower-cased database.view.
func viewSortName(database, name string) string {
	return strings.ToLower(database) + "." + strings.ToLower(name)
}

// viewCreateSortName prefixes the view sort name with a zero-padded topological depth so that,
// among the priorityView CREATE ops, a referenced view (smaller depth) sorts before its referrer
// (larger depth). The width is 5 only to keep realistic depths fixed-width; %d never truncates, so
// even an (impossible) depth ≥ 100000 still encodes its true magnitude — it would merely lose the
// fixed-width alignment, and a view graph that deep cannot exist anyway.
func viewCreateSortName(depth int, database, name string) string {
	return fmt.Sprintf("%05d.%s", depth, viewSortName(database, name))
}

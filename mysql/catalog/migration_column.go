package catalog

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// generateColumnDDL produces ALTER TABLE ... ADD/DROP/MODIFY COLUMN operations for every
// column change within modified tables. It is the MySQL analog of PG's generateColumnDDL
// (pg/catalog/migration_column.go:12), but where PG decomposes a column change into a chain
// of incremental ALTER COLUMN sub-statements (DROP DEFAULT → ALTER TYPE → SET DEFAULT, ...),
// MySQL re-declares the whole column with a single MODIFY COLUMN. So one diff column entry
// maps to exactly one op:
//
//   - DiffAdd    → ALTER TABLE ... ADD COLUMN <full canonical def>
//   - DiffDrop   → ALTER TABLE ... DROP COLUMN `name`
//   - DiffModify → ALTER TABLE ... MODIFY COLUMN <full canonical def>
//
// MODIFY (not CHANGE) is used because the column name is unchanged — diff-core keys columns
// by name, so a rename is seen as DROP+ADD, never a MODIFY. The full canonical definition
// routes through normalize-core (the same formatColumnDefinition CREATE TABLE uses), so the
// readback of the ALTER canonicalizes equal to `to`.
//
// Determinism: drops are emitted before adds/modifies (so a drop-then-add of the same name
// never collides). Drops and modifies sort by lower-cased column name, adds by the column's
// position in `to` (columnAddSortName); sortMigrationOps re-imposes the global order. An ADD
// COLUMN carries no placement here: positionColumnAdds gives it FIRST / AFTER once the final
// statement order is known.
func generateColumnDDL(_, to *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		// Only MODIFY tables carry column changes (ADD tables ship columns inline in CREATE
		// TABLE; DROP tables vanish whole). diff-core only populates Columns on a DiffModify
		// entry, but guard the action explicitly so a future ADD-with-columns shape cannot
		// leak a duplicate ALTER.
		if entry.Action != DiffModify || len(entry.Columns) == 0 {
			continue
		}
		// Use the table's original-case identifier (a MODIFY table always has To). Lower-casing
		// the diff identity key would target the wrong object on a case-sensitive server.
		table := tableIdent(entry.To)

		for _, col := range entry.Columns {
			switch col.Action {
			case DiffAdd:
				if col.To == nil {
					continue
				}
				ops = append(ops, addColumnOp(entry, table, col.To, n))

			case DiffDrop:
				if col.From == nil {
					continue
				}
				ops = append(ops, dropColumnOp(entry, table, col.Name, col.From))

			case DiffModify:
				if col.To == nil {
					continue
				}
				// MySQL rejects an in-place MODIFY that flips a generated column's storage mode
				// (VIRTUAL↔STORED) or that converts a plain column to/from generated (error 3106).
				// Express such a change as DROP + ADD instead of MODIFY.
				if generatedShapeChanged(col.From, col.To) {
					ops = append(ops, dropColumnOp(entry, table, col.Name, col.From))
					ops = append(ops, addColumnOp(entry, table, col.To, n))
					continue
				}
				clause := fmt.Sprintf("MODIFY COLUMN %s", formatColumnDefinition(entry.To, col.To, n))
				ops = append(ops, MigrationOp{
					Type:         OpModifyColumn,
					Database:     entry.Database,
					ObjectName:   entry.Name,
					ParentObject: entry.Name,
					SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
					Phase:        PhaseMain,
					Priority:     priorityColumn,
					sortName:     columnSortName(entry.Database, entry.Name, col.Name),
					alterClause:  clause,
					column:       toLower(col.Name),
				})
			}
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			// Drops (PhasePre) ahead of adds/modifies — also enforced by phase, but keep the
			// local order stable for readability when grouped by table.
			return columnOpRank(ops[i].Type) < columnOpRank(ops[j].Type)
		}
		// Within the same op type, lower Priority first (generated-column drops precede plain
		// drops so a generated→plain dependency never breaks), then name for determinism.
		if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addColumnOp builds an ADD COLUMN op. A generated column is added AFTER plain columns
// (priorityColumn+1) so that any base column it references — whether freshly added or being
// MODIFYed in the same PhaseMain batch — is already in place; this is the symmetric partner of
// dropColumnOp's generated-first rule (generated columns drop before, and add after, the plain
// columns they depend on). The global sortMigrationOps honors Priority within PhaseMain.
func addColumnOp(entry *TableDiffEntry, table string, to *Column, n *Normalizer) MigrationOp {
	clause := fmt.Sprintf("ADD COLUMN %s", formatColumnDefinition(entry.To, to, n))
	return MigrationOp{
		Type:         OpAddColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhaseMain,
		Priority:     columnAddPriority(to),
		sortName:     columnAddSortName(entry, to),
		alterClause:  clause,
		column:       toLower(to.Name),
	}
}

// columnAddPriority returns priorityColumn+1 for a generated column (add it after the plain
// columns its expression may reference) and priorityColumn for a plain column.
func columnAddPriority(c *Column) int {
	if c != nil && c.Generated != nil {
		return priorityColumn + 1
	}
	return priorityColumn
}

// dropColumnOp builds a DROP COLUMN op, giving a generated column a lower priority so it is
// dropped before the plain columns its expression may reference (MySQL rejects dropping a
// column a generated column depends on — error 3108). The global sortMigrationOps also honors
// Priority within PhasePre, so the ordering holds across the whole plan.
func dropColumnOp(entry *TableDiffEntry, table, colName string, from *Column) MigrationOp {
	clause := fmt.Sprintf("DROP COLUMN %s", migrationQuoteIdent(colName))
	return MigrationOp{
		Type:         OpDropColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhasePre,
		Priority:     columnDropPriority(from),
		sortName:     columnSortName(entry.Database, entry.Name, colName),
		alterClause:  clause,
		column:       toLower(colName),
	}
}

// columnDropPriority returns priorityColumn-1 for a generated column (drop it first) and
// priorityColumn for a plain column, so a generated column that references a plain column is
// always dropped before that plain column.
func columnDropPriority(c *Column) int {
	if c != nil && c.Generated != nil {
		return priorityColumn - 1
	}
	return priorityColumn
}

// generatedShapeChanged reports whether a column's GENERATED shape changed in a way MySQL
// cannot express as an in-place MODIFY: gaining or losing the generated attribute, or flipping
// VIRTUAL↔STORED. Such a change must be rendered as DROP + ADD. A change to only the generation
// EXPRESSION (storage mode unchanged) is a normal MODIFY and is not flagged here.
func generatedShapeChanged(from, to *Column) bool {
	if from == nil || to == nil {
		return false
	}
	fromGen, toGen := from.Generated != nil, to.Generated != nil
	if fromGen != toGen {
		return true
	}
	if fromGen && toGen && from.Generated.Stored != to.Generated.Stored {
		return true
	}
	return false
}

// columnOpRank orders column op types so drops sort ahead of adds/modifies locally.
func columnOpRank(t MigrationOpType) int {
	if t == OpDropColumn {
		return 0
	}
	return 1
}

// columnSortName is the stable secondary sort key for a column op: lower-cased
// database.table.column.
func columnSortName(database, table, column string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(column)
}

// columnAddSortName is the sort key of an ADD COLUMN: the table's key, a NUL, and the column's
// zero-padded position in `to` (MySQL allows at most 4096 columns). An identifier cannot contain
// U+0000, so a table's plain adds sort before its modifies, which may read them (a generated
// column's new expression).
func columnAddSortName(entry *TableDiffEntry, c *Column) string {
	return columnSortName(entry.Database, entry.Name, fmt.Sprintf("\x00%04d", columnPosition(entry.To, c.Name)))
}

// positionColumnAdds appends FIRST or AFTER to every ADD COLUMN clause of the sorted plan so the
// column lands where `to` places it; MySQL otherwise appends it. It runs after
// mergeAutoIncrementKeyOps, which can move an ADD COLUMN into another statement of its table, and
// walks the clauses in execution order, so an AFTER only names a column that exists by then: one
// that survives from `from`, or one an earlier clause added. A column that no existing column
// follows in `to` gets no clause.
func positionColumnAdds(ops []MigrationOp, diff *SchemaDiff) []MigrationOp {
	parts := make([][]alterPart, len(ops))
	added := make(map[string]map[string]bool)
	for i, op := range ops {
		parts[i] = op.alterParts
		if parts[i] == nil && op.Type == OpAddColumn {
			parts[i] = []alterPart{{clause: op.alterClause, addsColumn: op.column}}
		}
		for _, part := range parts[i] {
			if part.addsColumn == "" {
				continue
			}
			key := aiTableKey(op.Database, op.ObjectName)
			if added[key] == nil {
				added[key] = make(map[string]bool)
			}
			added[key][part.addsColumn] = true
		}
	}
	layouts := make(map[string]*columnLayout)
	for i := range diff.Tables {
		entry := &diff.Tables[i]
		key := aiTableKey(entry.Database, entry.Name)
		if added[key] != nil && entry.Action == DiffModify && entry.From != nil && entry.To != nil {
			layouts[key] = newColumnLayout(entry, added[key])
		}
	}
	for i := range ops {
		op := &ops[i]
		l := layouts[aiTableKey(op.Database, op.ObjectName)]
		joined := joinAlterParts(parts[i])
		if l == nil || parts[i] == nil || !strings.HasSuffix(op.SQL, joined) {
			continue
		}
		placed := make([]alterPart, len(parts[i]))
		for j, part := range parts[i] {
			placed[j] = part
			if part.addsColumn != "" {
				placed[j].clause += l.place(part.addsColumn)
			}
		}
		op.SQL = strings.TrimSuffix(op.SQL, joined) + joinAlterParts(placed)
	}
	return ops
}

// columnLayout tracks, for one modified table, which of the target's columns exist as the plan's
// ADD COLUMN clauses run.
type columnLayout struct {
	columns []*Column       // the target's columns in order, hidden ones excluded
	exists  map[string]bool // lower-cased name → the column exists at this point of the plan
}

// newColumnLayout returns the layout of a modified table whose plan adds the named columns. A
// column the plan drops and adds again does not exist until its ADD runs.
func newColumnLayout(entry *TableDiffEntry, added map[string]bool) *columnLayout {
	l := &columnLayout{exists: make(map[string]bool)}
	for _, c := range entry.To.Columns {
		if c == nil || c.Hidden != ColumnHiddenNone {
			continue
		}
		l.columns = append(l.columns, c)
		name := toLower(c.Name)
		l.exists[name] = !added[name] && entry.From.GetColumn(c.Name) != nil
	}
	return l
}

// place returns the placement clause for adding the named column and records that it exists.
func (l *columnLayout) place(name string) string {
	defer func() { l.exists[name] = true }()
	pos := slices.IndexFunc(l.columns, func(c *Column) bool { return toLower(c.Name) == name })
	if pos < 0 || !slices.ContainsFunc(l.columns[pos+1:], func(c *Column) bool { return l.exists[toLower(c.Name)] }) {
		return ""
	}
	for k := pos - 1; k >= 0; k-- {
		if l.exists[toLower(l.columns[k].Name)] {
			return " AFTER " + migrationQuoteIdent(l.columns[k].Name)
		}
	}
	return " FIRST"
}

// formatColumnDefinition renders a single column in MySQL's canonical stored form for the
// target version: `name` <type> [CHARACTER SET .. COLLATE ..] [GENERATED ALWAYS AS (..)
// VIRTUAL|STORED] [NULL|NOT NULL] [DEFAULT ..] [AUTO_INCREMENT] [ON UPDATE ..] [COMMENT ..]
// [INVISIBLE] [/*!80003 SRID n */]. It is shared by CREATE TABLE, ADD COLUMN, and MODIFY
// COLUMN, so all three express a column identically and round-trip through the same
// normalize-core canonical form diff-core compares — no ad-hoc surface rebuilding that would
// re-introduce diff noise.
//
// The clause ORDER follows MySQL's own SHOW CREATE rendering (show.go), which is also the
// order MySQL accepts on input, so a rendered definition is both valid DDL and canonical.
func formatColumnDefinition(tbl *Table, c *Column, n *Normalizer) string {
	var b strings.Builder
	b.WriteString(migrationQuoteIdent(c.Name))
	b.WriteString(" ")
	b.WriteString(n.CanonicalColumnType(c))

	// CHARACTER SET / COLLATE — emit the resolved pair when it differs from the table's
	// effective pair, so the column's stored charset/collation matches `to` after readback. A
	// bare column (same effective pair as the table) emits nothing and inherits, exactly as
	// MySQL stores it; an explicit divergent charset/collation is rendered in full. Resolving
	// through ResolveColumnCharsetCollation (the same key diff-core compares) keeps the
	// version-flagged echo asymmetry out of the rendering.
	if cs, coll := n.ResolveColumnCharsetCollation(tbl, c); cs != "" {
		tcs := foldCharset(tbl.Charset)
		tcoll := n.effectiveTableCollation(tbl)
		if cs != tcs {
			b.WriteString(" CHARACTER SET ")
			b.WriteString(cs)
			// The column declares a DIFFERENT charset than the table, so it cannot inherit the
			// table COLLATE; a bare `CHARACTER SET cs` falls back to cs's SERVER default. Emit
			// COLLATE only when the target collation is NOT that charset's default — so a
			// column whose collation IS the default renders no COLLATE, while a non-default one
			// is explicit.
			if coll != "" && !n.isCharsetDefaultCollation(cs, coll) {
				b.WriteString(" COLLATE ")
				b.WriteString(coll)
			}
		} else if coll != "" && coll != tcoll {
			// Same charset as the table but a collation that differs from what the column would
			// INHERIT from the table (tcoll) → emit COLLATE. The comparison is against the
			// table's effective collation, NOT the charset default: when the table itself
			// carries a non-default collation, a column at the charset default still differs
			// from the table and MUST be rendered explicitly, or it would wrongly inherit the
			// table's collation on apply (non-convergent). tcoll is already the version-correct
			// effective table collation, so this stays version-robust.
			b.WriteString(" COLLATE ")
			b.WriteString(coll)
		}
	}

	// Generated columns: GENERATED ALWAYS AS (expr) VIRTUAL|STORED. A generated column cannot
	// carry a DEFAULT/AUTO_INCREMENT/ON UPDATE, so this branch renders nullability + comment +
	// invisibility and returns. The expression is rendered verbatim (its stored form); the
	// diff compares it via CanonicalGeneratedExpr, so a round-trip matches.
	if c.Generated != nil {
		mode := "VIRTUAL"
		if c.Generated.Stored {
			mode = "STORED"
		}
		b.WriteString(" GENERATED ALWAYS AS (")
		b.WriteString(c.Generated.Expr)
		b.WriteString(") ")
		b.WriteString(mode)
		if n.CanonicalNotNull(tbl, c) {
			b.WriteString(" NOT NULL")
		}
		if c.NotSecondary {
			b.WriteString(" NOT SECONDARY")
		}
		writeColumnComment(&b, c)
		if c.Invisible {
			b.WriteString(" /*!80023 INVISIBLE */")
		}
		return b.String()
	}

	// Nullability. Render NOT NULL when the canonical (post-rewrite) state is NOT NULL. For a
	// nullable column we render NULL only for TIMESTAMP (MySQL's SHOW CREATE shows an explicit
	// NULL for a nullable TIMESTAMP, and rendering it keeps the EDFT-off magic from forcing
	// NOT NULL on apply); other nullable columns omit the keyword (the default).
	if n.CanonicalNotNull(tbl, c) {
		b.WriteString(" NOT NULL")
	} else if isTimestampType(c.DataType) {
		b.WriteString(" NULL")
	}
	if c.NotSecondary {
		b.WriteString(" NOT SECONDARY")
	}

	// SRID (spatial) — MySQL places it after NOT NULL, before DEFAULT. Emit on PRESENCE so
	// an explicit `SRID 0` round-trips (SHOW CREATE echoes `/*!80003 SRID 0 */`); a no-SRID
	// column leaves HasSRID false and emits nothing.
	if c.HasSRID {
		fmt.Fprintf(&b, " /*!80003 SRID %d */", c.SRID)
	}

	// DEFAULT. Render the column's declared default in canonical surface form (reusing
	// show.go's formatDefault, which normalizes CURRENT_TIMESTAMP and quoting). A bare first
	// TIMESTAMP under EDFT=OFF has no declared default but the engine injects
	// DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP on apply; CanonicalTimestampDefaults
	// predicts the same for `to`, so omitting it here still round-trips. We render only what is
	// declared to avoid emitting an injected default that a re-apply would also inject.
	if c.Default != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(formatDefault(*c.Default, c))
	}

	// AUTO_INCREMENT.
	if c.AutoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	}

	// ON UPDATE.
	if c.OnUpdate != "" {
		b.WriteString(" ON UPDATE ")
		b.WriteString(formatOnUpdate(c.OnUpdate, c.OnUpdateKind))
	}

	writeColumnComment(&b, c)

	if c.Invisible {
		b.WriteString(" /*!80023 INVISIBLE */")
	}

	return b.String()
}

// writeColumnComment appends a COMMENT clause when the column has one, escaped the way
// MySQL's stored form does.
func writeColumnComment(b *strings.Builder, c *Column) {
	if c.Comment != "" {
		b.WriteString(" COMMENT ")
		b.WriteString(quoteStringLiteral(c.Comment))
	}
}

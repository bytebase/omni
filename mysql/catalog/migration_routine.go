package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — stored routines (functions + procedures).
//
// This is the MySQL analog of PG's generateFunctionDDL (pg/catalog/migration_function.go). A
// single generate hook covers both routine kinds: wired into GenerateMigrationWithNormalizer
// (migration.go), it consumes BOTH SchemaDiff.Functions and SchemaDiff.Procedures and emits the
// CREATE / DROP / ALTER ops (OpCreateFunction / OpDropFunction / OpCreateProcedure /
// OpDropProcedure, priorityRoutine).
//
// MySQL has NO `CREATE OR REPLACE` and NO `ALTER` for a routine's body/params/RETURNS. So a
// change to any of those is a DROP followed by a CREATE. `ALTER FUNCTION/PROCEDURE` can change
// ONLY the characteristics SQL SECURITY, COMMENT, and DATA ACCESS (verified on both live
// engines — DETERMINISTIC is NOT alterable). diff_routine.go's routineAlterSuffices decides
// which path a MODIFY takes; this generator renders it:
//   - DiffAdd                       → CREATE.
//   - DiffDrop                      → DROP.
//   - DiffModify, ALTER suffices    → ALTER (one statement, the changed characteristics).
//   - DiffModify, ALTER insufficient → DROP + CREATE.
//
// Ordering (mirrors the index node's drop-before-add discipline): DROPs run in PhasePre, CREATEs
// and ALTERs in PhaseMain, both at priorityRoutine. So a MODIFY's DROP (PhasePre) always precedes
// its CREATE (PhaseMain) — the name is free for re-creation — and a routine dropped in one place
// and created in another never reverses. sortMigrationOps re-imposes the global phase/priority/
// name order.
//
// Rendering routes through the same form SHOW CREATE uses (showCreateRoutineBody), MINUS the
// DEFINER clause: DEFINER is ignored in the diff (it is environment, not schema-as-code), so the
// emitted CREATE must NOT pin a definer — a DEFINER-less CREATE adopts CURRENT_USER on apply,
// which is the correct release behavior and avoids requiring SET_USER_ID/SUPER for an account the
// user never named. Because the body/params/characteristics are rendered in the engine's stored
// form, the readback of the emitted CREATE canonicalizes equal to `to` and apply-correctness holds.
func generateRoutineDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	if diff == nil {
		return nil
	}
	var ops []MigrationOp
	for i := range diff.Functions {
		ops = append(ops, routineOps(&diff.Functions[i], n)...)
	}
	for i := range diff.Procedures {
		ops = append(ops, routineOps(&diff.Procedures[i], n)...)
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Phase != ops[j].Phase {
			return ops[i].Phase < ops[j].Phase
		}
		return ops[i].sortName < ops[j].sortName
	})
	return ops
}

// routineOps renders the ops for a single routine diff entry, choosing CREATE / DROP / ALTER /
// DROP+CREATE per the entry's action and the ALTER-suffices test.
func routineOps(e *RoutineDiffEntry, n *Normalizer) []MigrationOp {
	switch e.Action {
	case DiffAdd:
		if e.To == nil {
			return nil
		}
		return []MigrationOp{createRoutineOp(e, e.To, n)}
	case DiffDrop:
		if e.From == nil {
			return nil
		}
		return []MigrationOp{dropRoutineOp(e, e.From)}
	case DiffModify:
		if e.From == nil || e.To == nil {
			return nil
		}
		if routineAlterSuffices(e.From, e.To, n) {
			return []MigrationOp{alterRoutineOp(e, e.To)}
		}
		// Body/params/RETURNS/DETERMINISTIC changed → DROP (PhasePre) then CREATE (PhaseMain).
		return []MigrationOp{dropRoutineOp(e, e.From), createRoutineOp(e, e.To, n)}
	}
	return nil
}

// createRoutineOp builds a CREATE FUNCTION/PROCEDURE op (PhaseMain). The SQL is the DEFINER-less
// stored form (see the file doc).
func createRoutineOp(e *RoutineDiffEntry, r *Routine, n *Normalizer) MigrationOp {
	return MigrationOp{
		Type:       createRoutineOpType(e.IsProcedure),
		Database:   e.Database,
		ObjectName: e.Name,
		SQL:        renderCreateRoutine(e.Database, r, n),
		Phase:      PhaseMain,
		Priority:   priorityRoutine,
		sortName:   routineSortName(e.Database, e.Name),
	}
}

// dropRoutineOp builds a DROP FUNCTION/PROCEDURE op (PhasePre, so it precedes any re-CREATE in
// the same plan). IF EXISTS is NOT used: the diff already established the routine is present in
// `from`, and an unconditional DROP keeps the emitted DDL minimal and explicit (matching the
// table node's DROP TABLE, which is likewise unconditional).
func dropRoutineOp(e *RoutineDiffEntry, r *Routine) MigrationOp {
	kw := "FUNCTION"
	if e.IsProcedure {
		kw = "PROCEDURE"
	}
	return MigrationOp{
		Type:       dropRoutineOpType(e.IsProcedure),
		Database:   e.Database,
		ObjectName: e.Name,
		SQL:        fmt.Sprintf("DROP %s %s", kw, routineIdent(e.Database, r.Name)),
		Phase:      PhasePre,
		Priority:   priorityRoutine,
		sortName:   routineSortName(e.Database, e.Name),
	}
}

// alterRoutineOp builds an ALTER FUNCTION/PROCEDURE op carrying the routine's full ALTER-able
// characteristic set (DATA ACCESS, SQL SECURITY, COMMENT) in their canonical default-folded
// form. It re-states ALL three (not just the changed ones): re-applying a characteristic to its
// existing value is a harmless no-op, and emitting the whole set keeps the op self-contained and
// order-independent of what the synced side held. The op runs in PhaseMain at priorityRoutine.
//
// ALTER is chosen only when routineAlterSuffices(from, to) held — i.e. params, RETURNS, body, and
// DETERMINISTIC are unchanged — so this never needs to touch a DROP+CREATE-only attribute.
func alterRoutineOp(e *RoutineDiffEntry, to *Routine) MigrationOp {
	kw := "FUNCTION"
	if e.IsProcedure {
		kw = "PROCEDURE"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "ALTER %s %s", kw, routineIdent(e.Database, to.Name))
	for _, clause := range routineAlterableCharacteristicClauses(to) {
		b.WriteString(" ")
		b.WriteString(clause)
	}
	return MigrationOp{
		Type:       createRoutineOpType(e.IsProcedure), // ALTER reuses the CREATE op-type key (no OpAlter* in the set)
		Database:   e.Database,
		ObjectName: e.Name,
		SQL:        b.String(),
		Phase:      PhaseMain,
		Priority:   priorityRoutine,
		sortName:   routineSortName(e.Database, e.Name),
	}
}

// renderCreateRoutine renders a DEFINER-less CREATE FUNCTION/PROCEDURE in MySQL's canonical
// stored form. It mirrors showCreateRoutine (routinecmds.go) — same parameter rendering, same
// RETURNS, same characteristic lines, same body placement — but omits the DEFINER clause (see the
// file doc) and emits characteristics in default-folded canonical form so the readback matches.
//
// Characteristic emission policy: only NON-DEFAULT characteristics are emitted (SHOW CREATE drops
// defaults, so emitting them would still round-trip, but matching SHOW CREATE's policy keeps the
// emitted DDL byte-identical to a subsequent readback's shape):
//   - DETERMINISTIC only when YES (NOT DETERMINISTIC is the default and is omitted).
//   - DATA ACCESS only when not CONTAINS SQL (the default).
//   - SQL SECURITY only when INVOKER (DEFINER is the default).
//   - COMMENT only when non-empty.
func renderCreateRoutine(database string, r *Routine, _ *Normalizer) string {
	var b strings.Builder
	kw := "FUNCTION"
	if r.IsProcedure {
		kw = "PROCEDURE"
	}
	fmt.Fprintf(&b, "CREATE %s %s(", kw, routineIdent(database, r.Name))

	// Parameters: "DIR `name` type" for procedures (direction shown), "`name` type" for functions.
	// The parameter NAME is backtick-quoted: a name that is a reserved word (e.g. `select`) or
	// needs quoting would otherwise produce invalid DDL. Always-quoting is safe for round-trip —
	// the loader strips backticks from a parameter name (the canonical key uses the decoded name),
	// so a quoted emission and an unquoted user form canonicalize equal. MySQL stores the
	// parameter list verbatim, so the readback echoes the quoted form, which reloads to the same
	// decoded name.
	for i, p := range r.Params {
		if i > 0 {
			b.WriteString(", ")
		}
		if r.IsProcedure && strings.TrimSpace(p.Direction) != "" {
			fmt.Fprintf(&b, "%s %s %s", strings.ToUpper(p.Direction), migrationQuoteIdent(p.Name), p.TypeName)
		} else {
			fmt.Fprintf(&b, "%s %s", migrationQuoteIdent(p.Name), p.TypeName)
		}
	}
	b.WriteString(")")

	// RETURNS (functions only).
	if !r.IsProcedure && strings.TrimSpace(r.Returns) != "" {
		fmt.Fprintf(&b, " RETURNS %s", r.Returns)
	}

	// Characteristics, default-folded (non-defaults only), each on its own indented line, in the
	// order SHOW CREATE uses for a freshly created routine (DETERMINISTIC, DATA ACCESS, SQL
	// SECURITY, COMMENT).
	if strings.EqualFold(r.Characteristics["DETERMINISTIC"], "YES") {
		b.WriteString("\n    DETERMINISTIC")
	}
	if da := strings.ToUpper(strings.TrimSpace(r.Characteristics["DATA ACCESS"])); da != "" && da != "CONTAINS SQL" {
		fmt.Fprintf(&b, "\n    %s", da)
	}
	if sec := strings.ToUpper(strings.TrimSpace(r.Characteristics["SQL SECURITY"])); sec == "INVOKER" {
		b.WriteString("\n    SQL SECURITY INVOKER")
	}
	if cmt := r.Characteristics["COMMENT"]; cmt != "" {
		fmt.Fprintf(&b, "\n    COMMENT '%s'", escapeComment(cmt))
	}

	// Body on its own line (matches SHOW CREATE).
	if r.Body != "" {
		b.WriteString("\n")
		b.WriteString(r.Body)
	}
	return b.String()
}

// routineAlterableCharacteristicClauses returns the ALTER-able characteristic clauses for a
// routine in canonical default-folded form, always all three (DATA ACCESS, SQL SECURITY,
// COMMENT) so the ALTER fully establishes the target characteristic state regardless of the
// synced side. DETERMINISTIC is excluded (not ALTER-able). The DATA ACCESS clause is the keyword
// itself (CONTAINS SQL / NO SQL / READS SQL DATA / MODIFIES SQL DATA).
func routineAlterableCharacteristicClauses(r *Routine) []string {
	clauses := make([]string, 0, 3)

	da := strings.ToUpper(strings.TrimSpace(r.Characteristics["DATA ACCESS"]))
	if da == "" {
		da = "CONTAINS SQL"
	}
	clauses = append(clauses, da)

	sec := strings.ToUpper(strings.TrimSpace(r.Characteristics["SQL SECURITY"]))
	if sec == "" {
		sec = "DEFINER"
	}
	clauses = append(clauses, "SQL SECURITY "+sec)

	clauses = append(clauses, fmt.Sprintf("COMMENT '%s'", escapeComment(r.Characteristics["COMMENT"])))
	return clauses
}

// createRoutineOpType returns the CREATE op-type for the routine kind.
func createRoutineOpType(isProcedure bool) MigrationOpType {
	if isProcedure {
		return OpCreateProcedure
	}
	return OpCreateFunction
}

// dropRoutineOpType returns the DROP op-type for the routine kind.
func dropRoutineOpType(isProcedure bool) MigrationOpType {
	if isProcedure {
		return OpDropProcedure
	}
	return OpDropFunction
}

// routineIdent returns the backtick-quoted database.routine identifier for migration DDL, using
// ORIGINAL-CASE names (matching tableIdent's case discipline — a case-sensitive server stores the
// declared casing). The database qualifier is omitted when empty (the synced single-database
// release path may load with no qualifying database).
func routineIdent(database, name string) string {
	if database == "" {
		return migrationQuoteIdent(name)
	}
	return migrationQuoteIdent(database) + "." + migrationQuoteIdent(name)
}

// routineSortName is the stable secondary sort key for a routine op: lower-cased database.name.
func routineSortName(database, name string) string {
	return strings.ToLower(database) + "." + strings.ToLower(name)
}

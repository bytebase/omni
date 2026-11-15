package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateRoutineDDL turns SchemaDiff.Functions and SchemaDiff.Procedures into CREATE, DROP, and
// ALTER ops. ALTER FUNCTION/PROCEDURE changes only SQL SECURITY, COMMENT, and DATA ACCESS, so a
// change to the body, parameters, RETURNS, or DETERMINISTIC is a DROP followed by a CREATE;
// routineAlterSuffices decides which path a modify takes.
//
// DROPs run in PhasePre at priorityRoutineDrop and CREATEs and ALTERs in PhaseMain at
// priorityRoutine, before view creates, so a view body can call a routine created in the same
// plan.
//
// The rendered CREATE omits DEFINER, which the diff ignores; the server then uses the current
// user, and the applying account needs no privilege to set another definer.
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
		return []MigrationOp{createRoutineOp(e, e.To, nil, n)}
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
		return []MigrationOp{dropRoutineOp(e, e.From), createRoutineOp(e, e.To, e.From, n)}
	}
	return nil
}

// createRoutineOp builds a CREATE FUNCTION/PROCEDURE op, wrapped in a save/restore of the
// session context it is created under (renderWithSessionContext): r's own where it is known,
// else that of `from`, the routine it recreates. The sql_mode a routine is created under is
// part of what its body computes, PIPES_AS_CONCAT for one.
func createRoutineOp(e *RoutineDiffEntry, r, from *Routine, _ *Normalizer) MigrationOp {
	// The context the statement executes under — r's own, else that of the routine
	// it recreates — governs both the framing and how its comment literal is
	// escaped, so the two cannot disagree (a backslash under NO_BACKSLASH_ESCAPES).
	ctx := r
	if !r.HasSessionContext && from != nil {
		ctx = from
	}
	sql := renderCreateRoutine(e.Database, r, ctx.HasSessionContext && noBackslashEscapes(ctx.SQLMode))
	if ctx.HasSessionContext {
		sql = renderWithSessionContext(sql, routineSessionContext(ctx), false)
	}
	return MigrationOp{
		Type:       createRoutineOpType(e.IsProcedure),
		Database:   e.Database,
		ObjectName: e.Name,
		SQL:        sql,
		Phase:      PhaseMain,
		Priority:   priorityRoutine,
		sortName:   routineSortName(e.Database, e.Name),
	}
}

// routineSessionContext bundles a routine's captured session context for the recreate framing.
// Routines have no time_zone.
func routineSessionContext(r *Routine) SessionContext {
	return SessionContext{
		SQLMode:             r.SQLMode,
		CharacterSetClient:  r.CharacterSetClient,
		CollationConnection: r.CollationConnection,
	}
}

// priorityRoutineDrop orders a routine DROP within PhasePre after view drops, so a dropped
// dependent view goes before the routine it calls.
const priorityRoutineDrop = priorityView + 10

// dropRoutineOp builds a DROP FUNCTION/PROCEDURE op in PhasePre. It has no IF EXISTS because the
// diff established the routine exists in `from`.
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
		Priority:   priorityRoutineDrop,
		sortName:   routineSortName(e.Database, e.Name),
	}
}

// alterRoutineOp builds an ALTER FUNCTION/PROCEDURE op that states all three ALTER-able
// characteristics (DATA ACCESS, SQL SECURITY, COMMENT), so it does not depend on what the synced
// side held. It is used only when routineAlterSuffices held.
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
	sql := b.String()
	// The comment is escaped for the routine's captured sql_mode, so the
	// ALTER runs under that mode, as the CREATE does.
	if to.HasSessionContext {
		sql = renderWithSessionContext(sql, routineSessionContext(to), false)
	}
	return MigrationOp{
		Type:       createRoutineOpType(e.IsProcedure), // ALTER reuses the CREATE op type
		Database:   e.Database,
		ObjectName: e.Name,
		SQL:        sql,
		Phase:      PhaseMain,
		Priority:   priorityRoutine,
		sortName:   routineSortName(e.Database, e.Name),
	}
}

// renderCreateRoutine renders a DEFINER-less CREATE FUNCTION/PROCEDURE like showCreateRoutine
// (routinecmds.go). Only non-default characteristics are emitted:
//   - DETERMINISTIC only when YES.
//   - DATA ACCESS only when not CONTAINS SQL.
//   - SQL SECURITY only when INVOKER.
//   - COMMENT only when non-empty.
func renderCreateRoutine(database string, r *Routine, noBackslash bool) string {
	var b strings.Builder
	kw := "FUNCTION"
	if r.IsProcedure {
		kw = "PROCEDURE"
	}
	fmt.Fprintf(&b, "CREATE %s %s(", kw, routineIdent(database, r.Name))

	// Parameter names are always quoted so a reserved word stays valid; the loader strips the
	// backticks, so the quoted and unquoted forms compare equal.
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

	if !r.IsProcedure && strings.TrimSpace(r.Returns) != "" {
		fmt.Fprintf(&b, " RETURNS %s", r.Returns)
	}

	// Characteristics in the order SHOW CREATE writes them.
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
		fmt.Fprintf(&b, "\n    COMMENT '%s'", escapeCommentMode(cmt, noBackslash))
	}

	if r.Body != "" {
		b.WriteString("\n")
		b.WriteString(r.Body)
	}
	return b.String()
}

// routineAlterableCharacteristicClauses returns the DATA ACCESS, SQL SECURITY, and COMMENT
// clauses, always all three. DETERMINISTIC is not ALTER-able.
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

	clauses = append(clauses, fmt.Sprintf("COMMENT '%s'", escapeCommentMode(r.Characteristics["COMMENT"], r.HasSessionContext && noBackslashEscapes(r.SQLMode))))
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

// routineIdent returns the backtick-quoted database.routine identifier in original case,
// omitting the qualifier when database is empty.
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

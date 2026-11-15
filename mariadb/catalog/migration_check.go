package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateCheckDDL turns the check diff (TableDiffEntry.Checks) into ALTER TABLE ADD CONSTRAINT
// ... CHECK and DROP CONSTRAINT ops. MariaDB has no DROP CHECK.
//
// A modified table follows the per-check diff. A new table gets every check added after its
// CREATE, because formatCreateTable does not inline checks. A modified check (changed
// expression or enforcement) is a DROP in PhasePre followed by an ADD in PhaseMain, so the old
// definition has released the name.
//
// Check ADDs run at priorityConstraint, after the table CREATE and column ALTERs, so a check
// that references an added or modified column applies against the final column.
func generateCheckDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, addCheckOpsForNewTable(entry)...)
		case DiffModify:
			ops = append(ops, checkOpsForModifiedTable(entry)...)
		case DiffDrop:
			// DROP TABLE removes the checks; nothing to emit.
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			return checkOpRank(ops[i].Type) < checkOpRank(ops[j].Type)
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addCheckOpsForNewTable emits ADD ops for every table-level CHECK constraint on a new table,
// ordered by lower-cased name.
func addCheckOpsForNewTable(entry *TableDiffEntry) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, con := range orderedChecks(entry.To) {
		ops = append(ops, addCheckOp(entry, table, con))
	}
	return ops
}

// checkOpsForModifiedTable emits ADD, DROP, or DROP+ADD ops from the per-check diff of a
// modified table. A check cannot be altered in place.
func checkOpsForModifiedTable(entry *TableDiffEntry) []MigrationOp {
	if entry.To == nil || len(entry.Checks) == 0 {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, ce := range entry.Checks {
		switch ce.Action {
		case DiffAdd:
			if ce.To == nil {
				continue
			}
			ops = append(ops, addCheckOp(entry, table, ce.To))
		case DiffDrop:
			if ce.From == nil {
				continue
			}
			ops = append(ops, dropCheckOp(entry, table, ce.From))
		case DiffModify:
			if ce.From == nil || ce.To == nil {
				continue
			}
			ops = append(ops, dropCheckOp(entry, table, ce.From))
			ops = append(ops, addCheckOp(entry, table, ce.To))
		}
	}
	return ops
}

// addCheckOp builds an ALTER TABLE ... ADD CONSTRAINT `name` CHECK (...) op.
func addCheckOp(entry *TableDiffEntry, table string, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpAddCheck,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s ADD %s", table, formatCheckClause(con)),
		Phase:        PhaseMain,
		Priority:     priorityConstraint,
		sortName:     checkSortName(entry.Database, entry.Name, con.Name),
	}
}

// priorityCheckDrop orders a CHECK DROP within PhasePre after table drops and before any column
// drop, so the DROP CONSTRAINT runs while the constraint still exists and a modified check's
// drop precedes its PhaseMain re-add.
const priorityCheckDrop = priorityColumn - 5

// dropCheckOp drops a table check by name with DROP CONSTRAINT.
func dropCheckOp(entry *TableDiffEntry, table string, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpDropCheck,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", table, migrationQuoteIdent(con.Name)),
		Phase:        PhasePre,
		Priority:     priorityCheckDrop,
		sortName:     checkSortName(entry.Database, entry.Name, con.Name),
	}
}

// formatCheckClause is a table check as SHOW CREATE TABLE writes it.
func formatCheckClause(con *Constraint) string {
	return fmt.Sprintf("CONSTRAINT %s CHECK (%s)", migrationQuoteIdent(con.Name), con.CheckExpr)
}

// orderedChecks returns a table's table-level CHECK constraints ordered by lower-cased name.
func orderedChecks(t *Table) []*Constraint {
	var out []*Constraint
	if t == nil {
		return out
	}
	for _, con := range t.Constraints {
		// A column's check is declared with the column.
		if con == nil || con.Type != ConCheck || con.Name == "" || isColumnCheck(con) {
			continue
		}
		out = append(out, con)
	}
	sort.Slice(out, func(i, j int) bool {
		return toLower(out[i].Name) < toLower(out[j].Name)
	})
	return out
}

// checkOpRank orders check op types so drops sort ahead of adds.
func checkOpRank(t MigrationOpType) int {
	if t == OpDropCheck {
		return 0
	}
	return 1
}

// checkSortName is the stable secondary sort key for a check op: lower-cased
// database.table.constraint.
func checkSortName(database, table, check string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(check)
}

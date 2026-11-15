package catalog

import (
	"fmt"
	"strings"
)

// An application-time period is dropped after the keys WITHOUT OVERLAPS over it
// and ahead of the columns it covers, and added after the columns, ahead of the
// keys. Measured on mariadb:11.8.8: DROP PERIOD while such a key holds it is
// 4156. System versioning is added once the table's columns are in place, and
// dropped after the keys and checks the plan drops, which may cover the ROW
// START and ROW END columns it takes with it, and ahead of the column changes.
const (
	priorityPeriodDrop     = priorityColumn - 2
	priorityPeriodAdd      = priorityIndex - 1
	priorityVersioningDrop = priorityIndexDrop + 1
	priorityVersioningAdd  = priorityIndex - 2
)

// addedWithVersioning reports whether col is a ROW START or ROW END column the
// plan adds as it makes the table system-versioned. MariaDB takes such a column
// only in the ALTER that adds SYSTEM VERSIONING (4124 otherwise), so that ALTER
// adds it. Measured on mariadb:11.8.8.
func addedWithVersioning(e *TableDiffEntry, col *Column) bool {
	return col.Generated != nil && col.Generated.RowBound != "" && e.From != nil && e.To != nil && e.To.SystemVersioned && (!e.From.SystemVersioned || systemPeriodReplaced(e))
}

// systemPeriodReplaced reports whether a table that stays system-versioned
// changes the columns of its PERIOD FOR SYSTEM_TIME, which MariaDB changes only
// with the versioning: the plan drops and adds it again.
func systemPeriodReplaced(e *TableDiffEntry) bool {
	return e.From.SystemVersioned && e.To.SystemVersioned &&
		(!eqFoldStr(e.From.PeriodStartCol, e.To.PeriodStartCol) || !eqFoldStr(e.From.PeriodEndCol, e.To.PeriodEndCol))
}

// droppedWithVersioning is addedWithVersioning for the ALTER that drops SYSTEM
// VERSIONING, which has to drop the table's ROW START and ROW END columns with
// it (4125 otherwise).
func droppedWithVersioning(e *TableDiffEntry, col *Column) bool {
	return col.Generated != nil && col.Generated.RowBound != "" && e.From != nil && e.To != nil && e.From.SystemVersioned && (!e.To.SystemVersioned || systemPeriodReplaced(e))
}

func generateVersioningDDL(diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify || e.From == nil || e.To == nil {
			continue
		}
		from, to := e.From, e.To
		table := tableIdent(to)
		op := func(sql string, phase MigrationPhase, priority int) {
			ops = append(ops, MigrationOp{
				Type: OpAlterVersioning, Database: e.Database, ObjectName: e.Name, ParentObject: e.Name,
				SQL: sql, Phase: phase, Priority: priority, sortName: tableSortName(e.Database, e.Name),
			})
		}
		if appPeriodChanged(from, to) {
			if from.AppPeriodName != "" {
				op(fmt.Sprintf("ALTER TABLE %s DROP PERIOD FOR %s", table, backtickIDQ(from.AppPeriodName)), PhasePre, priorityPeriodDrop)
			}
			if to.AppPeriodName != "" {
				op(fmt.Sprintf("ALTER TABLE %s ADD PERIOD FOR %s (%s, %s)", table, backtickIDQ(to.AppPeriodName), backtickIDQ(to.AppPeriodStartCol), backtickIDQ(to.AppPeriodEndCol)), PhaseMain, priorityPeriodAdd)
			}
		}
		replaced := systemPeriodReplaced(e)
		if from.SystemVersioned && (!to.SystemVersioned || replaced) {
			clauses := []string{"DROP SYSTEM VERSIONING"}
			// A ROW START or ROW END column kept as a plain one is dropped
			// here and added again by the column generator.
			for _, ce := range e.Columns {
				if (ce.Action == DiffDrop || ce.Action == DiffModify && generatedShapeChanged(ce.From, ce.To)) && ce.From != nil && droppedWithVersioning(e, ce.From) {
					clauses = append(clauses, "DROP COLUMN "+backtickIDQ(ce.From.Name))
				}
			}
			op(fmt.Sprintf("ALTER TABLE %s %s", table, strings.Join(clauses, ", ")), PhasePre, priorityVersioningDrop)
			ops[len(ops)-1].dropsVersionedColumns = len(clauses) > 1
			if replaced {
				ops[len(ops)-1].Warning = fmt.Sprintf("changing the system-time period of %s drops its system versioning and adds it again, which deletes the history rows", table)
			}
		}
		if to.SystemVersioned && (!from.SystemVersioned || replaced) {
			// A column that becomes ROW START or ROW END is added here too,
			// after its plain self is dropped.
			added := make(map[string]bool)
			for _, ce := range e.Columns {
				if ce.Action == DiffAdd || ce.Action == DiffModify && generatedShapeChanged(ce.From, ce.To) {
					added[toLower(ce.Name)] = true
				}
			}
			var clauses []string
			for _, col := range to.Columns {
				if added[toLower(col.Name)] && addedWithVersioning(e, col) {
					clauses = append(clauses, "ADD COLUMN "+formatColumnDefinition(to, col, n))
				}
			}
			if len(clauses) > 0 && to.PeriodStartCol != "" {
				clauses = append(clauses, fmt.Sprintf("ADD PERIOD FOR SYSTEM_TIME (%s, %s)", backtickIDQ(to.PeriodStartCol), backtickIDQ(to.PeriodEndCol)))
			}
			clauses = append(clauses, "ADD SYSTEM VERSIONING")
			op(fmt.Sprintf("ALTER TABLE %s %s", table, strings.Join(clauses, ", ")), PhaseMain, priorityVersioningAdd)
		}
	}
	return ops
}

// keepVersionedHistory puts SET system_versioning_alter_history = KEEP ahead of
// a plan that changes a system-versioned table's columns, keys or constraints,
// which MariaDB otherwise refuses (4119). KEEP leaves the history rows as they
// are. Measured on mariadb:11.8.8.
func keepVersionedHistory(ops []MigrationOp, diff *SchemaDiff) []MigrationOp {
	versioned := make(map[string]bool)
	for _, e := range diff.Tables {
		if e.Action == DiffModify && e.From != nil && e.From.SystemVersioned {
			versioned[e.Database+"."+toLower(e.Name)] = true
		}
	}
	for _, op := range ops {
		switch {
		case op.Type == OpAlterVersioning && !op.dropsVersionedColumns, op.Type == OpDropTable, op.Type == OpCreateTable:
			continue
		}
		if op.ParentObject != "" && versioned[op.Database+"."+toLower(op.ParentObject)] {
			keep := MigrationOp{Type: OpAlterVersioning, SQL: "SET @@system_versioning_alter_history = KEEP", Phase: PhasePre}
			return append([]MigrationOp{keep}, ops...)
		}
	}
	return ops
}

func appPeriodChanged(from, to *Table) bool {
	return !eqFoldStr(from.AppPeriodName, to.AppPeriodName) || !eqFoldStr(from.AppPeriodStartCol, to.AppPeriodStartCol) || !eqFoldStr(from.AppPeriodEndCol, to.AppPeriodEndCol)
}

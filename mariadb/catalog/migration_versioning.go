package catalog

import "fmt"

// An application-time period is dropped after the keys WITHOUT OVERLAPS over it
// and ahead of the columns it covers, and added after the columns, ahead of the
// keys. Measured on mariadb:11.8.8: DROP PERIOD while such a key holds it is
// 4156. System
// versioning is added once the table's columns are in place and dropped ahead
// of column changes.
const (
	priorityPeriodDrop     = priorityColumn - 2
	priorityPeriodAdd      = priorityIndex - 1
	priorityVersioningDrop = priorityColumn - 7
	priorityVersioningAdd  = priorityIndex - 2
)

func generateVersioningDDL(diff *SchemaDiff) []MigrationOp {
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
		if !eqFoldStr(from.AppPeriodName, to.AppPeriodName) || !eqFoldStr(from.AppPeriodStartCol, to.AppPeriodStartCol) || !eqFoldStr(from.AppPeriodEndCol, to.AppPeriodEndCol) {
			if from.AppPeriodName != "" {
				op(fmt.Sprintf("ALTER TABLE %s DROP PERIOD FOR %s", table, backtickIDQ(from.AppPeriodName)), PhasePre, priorityPeriodDrop)
			}
			if to.AppPeriodName != "" {
				op(fmt.Sprintf("ALTER TABLE %s ADD PERIOD FOR %s (%s, %s)", table, backtickIDQ(to.AppPeriodName), backtickIDQ(to.AppPeriodStartCol), backtickIDQ(to.AppPeriodEndCol)), PhaseMain, priorityPeriodAdd)
			}
		}
		switch {
		case from.SystemVersioned && !to.SystemVersioned:
			op(fmt.Sprintf("ALTER TABLE %s DROP SYSTEM VERSIONING", table), PhasePre, priorityVersioningDrop)
		case !from.SystemVersioned && to.SystemVersioned:
			op(fmt.Sprintf("ALTER TABLE %s ADD SYSTEM VERSIONING", table), PhaseMain, priorityVersioningAdd)
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
		switch op.Type {
		case OpAlterVersioning, OpDropTable, OpCreateTable:
			continue
		}
		if op.ParentObject != "" && versioned[op.Database+"."+toLower(op.ParentObject)] {
			keep := MigrationOp{Type: OpAlterVersioning, SQL: "SET @@system_versioning_alter_history = KEEP", Phase: PhasePre}
			return append([]MigrationOp{keep}, ops...)
		}
	}
	return ops
}

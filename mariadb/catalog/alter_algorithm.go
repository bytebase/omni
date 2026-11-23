package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The ALGORITHM clause of ALTER TABLE and CREATE INDEX on InnoDB. MariaDB
// orders the algorithms COPY < INPLACE < NOCOPY < INSTANT and runs a statement
// under the one it names only when each of its operations can run so: one
// that cannot is 1846 with the reason InnoDB gives, or 1845 where it gives
// none, naming the best algorithm the statement can have (the least of its
// operations'). It judges this after everything else about the statement. Each
// classification below was read on mariadb:11.8.8 over CREATE TABLE t (b INT
// PRIMARY KEY, c INT, d INT, KEY idx_d (d)):
//
//   - ADD INDEX and CREATE INDEX, NOCOPY ("ADD INDEX");
//   - DROP INDEX, and DROP COLUMN of an indexed column, NOCOPY ("DROP INDEX");
//   - MODIFY c BIGINT, COPY ("Cannot change column type");
//   - MODIFY c INT NOT NULL, INPLACE;
//   - DROP PRIMARY KEY, COPY ("Dropping a primary key is not allowed without
//     also adding a new primary key"), and with ADD PRIMARY KEY INPLACE;
//   - ENGINE = InnoDB, INPLACE;
//   - ADD CONSTRAINT ... CHECK, COPY.
//
// ADD COLUMN (FIRST too), DROP COLUMN of a column no key reads, RENAME COLUMN,
// SET DEFAULT, RENAME INDEX and ALTER INDEX ... IGNORED run INSTANT. An
// operation the catalog does not classify is taken to run under any
// algorithm, so the catalog refuses only what MariaDB certainly refuses. LOCK
// is not judged: MariaDB copies a table online, and took LOCK = NONE beside
// every operation above.

// Algorithm levels, from the one any operation runs under up.
const (
	algoCopy = iota
	algoInplace
	algoNocopy
	algoInstant
)

var algoNames = [...]string{"COPY", "INPLACE", "NOCOPY", "INSTANT"}

// algoLimit is the best algorithm one operation runs under, and the reason
// InnoDB gives for refusing a better one ("" for none).
type algoLimit struct {
	best   int
	reason string
}

func algoLevel(name string) (int, bool) {
	for i, n := range algoNames {
		if strings.EqualFold(strings.TrimSpace(name), n) {
			return i, true
		}
	}
	return 0, false
}

// alterAlgorithmLimits classifies the operations of an ALTER TABLE over old,
// the table before it.
func alterAlgorithmLimits(old *Table, commands []*nodes.AlterTableCmd) []algoLimit {
	addsPrimary := false
	for _, cmd := range commands {
		if cmd != nil && (cmd.Type == nodes.ATAddConstraint || cmd.Type == nodes.ATAddIndex) && cmd.Constraint != nil && cmd.Constraint.Type == nodes.ConstrPrimaryKey {
			addsPrimary = true
		}
	}
	// An IF [NOT] EXISTS clause that finds what it looks for is a note, and
	// runs no operation: over a table holding key k, ADD INDEX IF NOT EXISTS k
	// (c), ADD COLUMN d INT, ALGORITHM = INSTANT runs (mariadb:11.8.8).
	_, skippedKeys := computeSkippedAlterAdditions(old, commands)
	var limits []algoLimit
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATAddConstraint, nodes.ATAddIndex:
			if cmd.Constraint == nil || skippedKeys[cmd.Constraint] {
				continue
			}
			switch cmd.Constraint.Type {
			case nodes.ConstrIndex:
				limits = append(limits, algoLimit{algoNocopy, "ADD INDEX"})
			case nodes.ConstrCheck:
				limits = append(limits, algoLimit{algoCopy, ""})
			}
		case nodes.ATDropIndex, nodes.ATDropConstraint:
			// DROP PRIMARY KEY is a dropped constraint or index named PRIMARY
			// (alterDropsPrimaryKey).
			switch {
			case strings.EqualFold(cmd.Name, "PRIMARY") && addsPrimary:
				limits = append(limits, algoLimit{algoInplace, ""})
			case strings.EqualFold(cmd.Name, "PRIMARY"):
				limits = append(limits, algoLimit{algoCopy, "Dropping a primary key is not allowed without also adding a new primary key"})
			case cmd.Type == nodes.ATDropIndex && !(cmd.IfExists && !indexNameExists(old, cmd.Name)):
				limits = append(limits, algoLimit{algoNocopy, "DROP INDEX"})
			}
		case nodes.ATDropColumn:
			if columnIndexed(old, cmd.Name) {
				limits = append(limits, algoLimit{algoNocopy, "DROP INDEX"})
			}
		case nodes.ATModifyColumn, nodes.ATChangeColumn:
			// A definition with more than its type and nullability (a default,
			// a key, a check) is left unclassified.
			if cmd.Column == nil || !onlyNullability(cmd.Column) {
				continue
			}
			was := old.GetColumn(cmd.Name)
			if was == nil || cmd.Column.TypeName == nil {
				continue
			}
			is := normalizedColumnDataType(cmd.Column.TypeName)
			switch {
			case isIntegerTypeName(was.DataType) && isIntegerTypeName(is) && !strings.EqualFold(was.DataType, is):
				limits = append(limits, algoLimit{algoCopy, "Cannot change column type"})
			case strings.EqualFold(was.ColumnType, formatColumnType(cmd.Column.TypeName)) && was.Nullable && columnDefNotNull(cmd.Column):
				limits = append(limits, algoLimit{algoInplace, ""})
			}
		case nodes.ATTableOption:
			if cmd.Option != nil && (eqFoldStr(cmd.Option.Name, "ENGINE") || eqFoldStr(cmd.Option.Name, "STORAGE ENGINE")) && isInnoDB(old.Engine) && isInnoDB(canonicalEngineName(cmd.Option.Value)) {
				limits = append(limits, algoLimit{algoInplace, ""})
			}
		}
	}
	return limits
}

func isIntegerTypeName(dt string) bool {
	switch strings.ToLower(dt) {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		return true
	}
	return false
}

func onlyNullability(def *nodes.ColumnDef) bool {
	for _, cc := range def.Constraints {
		if cc != nil && cc.Type != nodes.ColConstrNotNull && cc.Type != nodes.ColConstrNull {
			return false
		}
	}
	return def.Generated == nil && !def.AutoIncrement
}

func columnDefNotNull(def *nodes.ColumnDef) bool {
	for _, cc := range def.Constraints {
		if cc != nil && cc.Type == nodes.ColConstrNotNull {
			return true
		}
	}
	return false
}

func columnIndexed(tbl *Table, name string) bool {
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if strings.EqualFold(ic.Name, name) {
				return true
			}
		}
	}
	return false
}

// checkAlgorithm holds a statement's operations to the ALGORITHM it names.
func checkAlgorithm(algorithm string, limits []algoLimit) error {
	want, ok := algoLevel(algorithm)
	if !ok || len(limits) == 0 {
		return nil
	}
	worst := limits[0]
	for _, l := range limits[1:] {
		if l.best < worst.best {
			worst = l
		}
	}
	if want <= worst.best {
		return nil
	}
	if worst.reason == "" {
		return &Error{Code: 1845, SQLState: "0A000", Message: fmt.Sprintf(
			"ALGORITHM=%s is not supported for this operation. Try ALGORITHM=%s", algoNames[want], algoNames[worst.best])}
	}
	return &Error{Code: 1846, SQLState: "0A000", Message: fmt.Sprintf(
		"ALGORITHM=%s is not supported. Reason: %s. Try ALGORITHM=%s", algoNames[want], worst.reason, algoNames[worst.best])}
}

// alterAlgorithm returns the ALGORITHM an ALTER TABLE names, the last one.
func alterAlgorithm(commands []*nodes.AlterTableCmd) string {
	algorithm := ""
	for _, cmd := range commands {
		if cmd != nil && cmd.Type == nodes.ATAlgorithm {
			algorithm = cmd.Name
		}
	}
	return algorithm
}

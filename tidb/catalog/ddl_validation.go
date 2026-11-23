package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// mysqlValidTableEngineNames is TiDB v8.5.5's storage engine allowlist from
// pingcap/tidb pkg/planner/core/preprocess.go (checkTableEngine).
var mysqlValidTableEngineNames = map[string]bool{
	"archive":    true,
	"blackhole":  true,
	"csv":        true,
	"example":    true,
	"federated":  true,
	"innodb":     true,
	"memory":     true,
	"merge":      true,
	"mgr_myisam": true,
	"myisam":     true,
	"ndb":        true,
	"heap":       true,
}

// checkUnsupportedTableOptions mirrors TiDB v8.5.5's preprocess-time table
// option checks (UNION -> 8232, INSERT_METHOD -> 8233, unknown ENGINE -> 1286),
// which run before database/table lookup.
func checkUnsupportedTableOptions(options []*nodes.TableOption) error {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		if err := checkUnsupportedTableOption(opt); err != nil {
			return err
		}
	}
	return nil
}

func checkUnsupportedTableOption(opt *nodes.TableOption) error {
	switch toLower(opt.Name) {
	case "union":
		return &Error{
			Code:     ErrTableOptionUnionUnsupported,
			SQLState: sqlState(ErrTableOptionUnionUnsupported),
			Message:  "CREATE/ALTER table with union option is not supported",
		}
	case "insert_method":
		return &Error{
			Code:     ErrTableOptionInsertMethodUnsupported,
			SQLState: sqlState(ErrTableOptionInsertMethodUnsupported),
			Message:  "CREATE/ALTER table with insert method option is not supported",
		}
	case "engine", "storage engine":
		if !mysqlValidTableEngineNames[strings.ToLower(opt.Value)] {
			return &Error{
				Code:     ErrUnknownStorageEngine,
				SQLState: sqlState(ErrUnknownStorageEngine),
				Message:  fmt.Sprintf("Unknown storage engine '%s'", opt.Value),
			}
		}
	}
	return nil
}

// checkAlterUnsupportedTableOptions checks every ATTableOption in an ALTER
// TABLE statement during preprocessing, before database/table lookup.
func checkAlterUnsupportedTableOptions(cmds []*nodes.AlterTableCmd) error {
	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATTableOption || cmd.Option == nil {
			continue
		}
		if err := checkUnsupportedTableOption(cmd.Option); err != nil {
			return err
		}
	}
	return nil
}

// checkUnsupportedAddColumnConstraint mirrors TiDB v8.5.5's
// pkg/ddl/add_column.go checkUnsupportedColumnConstraint, which refuses inline
// AUTO_INCREMENT (including SERIAL / SERIAL DEFAULT VALUE), PRIMARY KEY,
// UNIQUE KEY (8200), and AUTO_RANDOM (8216) on ALTER TABLE ADD COLUMN ahead of
// duplicate-column (1060) and multiple-PK (1068) checks.
func checkUnsupportedAddColumnConstraint(tbl *Table, colDef *nodes.ColumnDef) error {
	dbName := ""
	if tbl.Database != nil {
		dbName = tbl.Database.Name
	}
	if colDef.AutoRandom {
		return errAutoRandomAddColumn(colDef.Name, dbName, tbl.Name)
	}
	if colDef.TypeName != nil && strings.EqualFold(colDef.TypeName.Name, "serial") {
		return &Error{
			Code:     ErrUnsupportedDDL,
			SQLState: sqlState(ErrUnsupportedDDL),
			Message:  fmt.Sprintf("unsupported add column '%s' constraint AUTO_INCREMENT when altering '%s.%s'", colDef.Name, dbName, tbl.Name),
		}
	}
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrAutoIncrement:
			return &Error{
				Code:     ErrUnsupportedDDL,
				SQLState: sqlState(ErrUnsupportedDDL),
				Message:  fmt.Sprintf("unsupported add column '%s' constraint AUTO_INCREMENT when altering '%s.%s'", colDef.Name, dbName, tbl.Name),
			}
		case nodes.ColConstrPrimaryKey:
			return &Error{
				Code:     ErrUnsupportedDDL,
				SQLState: sqlState(ErrUnsupportedDDL),
				Message:  fmt.Sprintf("unsupported add column '%s' constraint PRIMARY KEY when altering '%s.%s'", colDef.Name, dbName, tbl.Name),
			}
		case nodes.ColConstrUnique:
			return &Error{
				Code:     ErrUnsupportedDDL,
				SQLState: sqlState(ErrUnsupportedDDL),
				Message:  fmt.Sprintf("unsupported add column '%s' constraint UNIQUE KEY when altering '%s.%s'", colDef.Name, dbName, tbl.Name),
			}
		}
	}
	if colDef.AutoIncrement {
		return &Error{
			Code:     ErrUnsupportedDDL,
			SQLState: sqlState(ErrUnsupportedDDL),
			Message:  fmt.Sprintf("unsupported add column '%s' constraint AUTO_INCREMENT when altering '%s.%s'", colDef.Name, dbName, tbl.Name),
		}
	}
	return nil
}

// checkAlterTableCharset enforces TiDB v8.5.5's checkModifyCharsetAndCollation
// (pkg/ddl/executor.go): only utf8->utf8mb4, latin1->utf8mb4, utf8->utf8,
// utf8mb4->utf8mb4, or keeping the exact same charset and collation are
// supported on ALTER TABLE.
func checkAlterTableCharset(origCharset, origCollate, toCharset, toCollate string) error {
	origCS := normalizeCharsetName(origCharset)
	toCS := normalizeCharsetName(toCharset)
	origColl := normalizeCollationName(origCollate)
	toColl := normalizeCollationName(toCollate)
	if origCS == "" {
		return nil
	}
	if (origCS == "utf8" && toCS == "utf8mb4") ||
		(origCS == "utf8" && toCS == "utf8") ||
		(origCS == "utf8mb4" && toCS == "utf8mb4") ||
		(origCS == "latin1" && toCS == "utf8mb4") {
		return nil
	}
	if toCS != origCS {
		return &Error{
			Code:     ErrUnsupportedDDL,
			SQLState: sqlState(ErrUnsupportedDDL),
			Message:  fmt.Sprintf("Unsupported modify charset from %s to %s", origCS, toCS),
		}
	}
	if toColl != origColl {
		return &Error{
			Code:     ErrUnsupportedDDL,
			SQLState: sqlState(ErrUnsupportedDDL),
			Message:  fmt.Sprintf("Unsupported modify change collate from %s to %s", origColl, toColl),
		}
	}
	return nil
}

// checkAlterAttributes validates TiDB's ATTRIBUTES = '...' specification
// (pkg/ddl/label/attributes.go): comma-separated non-empty key=value pairs with
// no conflicting values for the same key (8237).
func checkAlterAttributes(spec *nodes.AttributesSpec) error {
	if spec == nil || spec.Default || strings.TrimSpace(spec.Value) == "" {
		return nil
	}
	seen := make(map[string]string)
	for _, raw := range strings.Split(spec.Value, ",") {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}
		kv := strings.Split(part, "=")
		if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" || strings.TrimSpace(kv[1]) == "" {
			return &Error{
				Code:     ErrInvalidAttributes,
				SQLState: sqlState(ErrInvalidAttributes),
				Message:  fmt.Sprintf("Invalid attributes: attributes should be in format 'key=value': %s", part),
			}
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if prev, ok := seen[k]; ok && prev != v {
			return &Error{
				Code:     ErrInvalidAttributes,
				SQLState: sqlState(ErrInvalidAttributes),
				Message:  fmt.Sprintf("Invalid attributes: '%s=%s' and '%s=%s' are conflicted", k, v, k, prev),
			}
		}
		seen[k] = v
	}
	return nil
}

// checkAlterAffinity validates TiDB v8.5.5's AFFINITY table option value
// (pkg/parser/ast/ddl.go NormalizeTableAffinityLevel): only "", "none",
// "table", and "partition" (case-insensitive) are accepted; anything else is 8266.
func checkAlterAffinity(level string) error {
	switch strings.ToLower(level) {
	case "", "none", "table", "partition":
		return nil
	default:
		return &Error{
			Code:     ErrInvalidTableAffinity,
			SQLState: sqlState(ErrInvalidTableAffinity),
			Message:  fmt.Sprintf("Invalid AFFINITY '%s'", level),
		}
	}
}

// errCompactNonPartitionedTable returns TiDB v8.5.5's 1105 error when COMPACT
// PARTITION specifies partition names on a non-partitioned table or view.
func errCompactNonPartitionedTable(tableName string, partNames []string) error {
	return &Error{
		Code:     ErrUnknownError,
		SQLState: sqlState(ErrUnknownError),
		Message:  fmt.Sprintf("table:%s is not a partition table, but user specify partition name list:[%s]", tableName, strings.Join(partNames, " ")),
	}
}

// checkModifyPrimaryKeyColumnType mirrors TiDB v8.5.5's pkg/ddl/modify_column.go
// checkModifyTypes: modifying the data type of a column that belongs to the
// primary key when the type change requires reorganization (e.g. narrowing an
// integer type or changing signedness/type family) is refused with 8200.
func checkModifyPrimaryKeyColumnType(tbl *Table, oldCol, newCol *Column) error {
	if oldCol == nil || newCol == nil || !columnInPrimaryKey(tbl, oldCol.Name) {
		return nil
	}
	if modifyColumnTypeNeedsReorg(oldCol, newCol) {
		return &Error{
			Code:     ErrUnsupportedDDL,
			SQLState: sqlState(ErrUnsupportedDDL),
			Message:  "Unsupported modify column: this column has primary key flag",
		}
	}
	return nil
}

func intTypeRank(dt string) int {
	switch strings.ToLower(dt) {
	case "tinyint":
		return 1
	case "smallint":
		return 2
	case "mediumint":
		return 3
	case "int", "integer":
		return 4
	case "bigint":
		return 5
	}
	return 0
}

func normalizeDataTypeAlias(dt string) string {
	switch strings.ToLower(strings.TrimSpace(dt)) {
	case "bool", "boolean", "int1":
		return "tinyint"
	case "int2":
		return "smallint"
	case "int3":
		return "mediumint"
	case "integer", "int4":
		return "int"
	case "int8":
		return "bigint"
	case "numeric", "dec", "fixed":
		return "decimal"
	case "real", "double precision":
		return "double"
	case "character", "nchar", "national char", "national character":
		return "char"
	case "character varying", "nvarchar", "national varchar", "national character varying":
		return "varchar"
	case "serial":
		return "bigint"
	default:
		return strings.ToLower(strings.TrimSpace(dt))
	}
}

func modifyColumnTypeNeedsReorg(oldCol, newCol *Column) bool {
	oldDT := normalizeDataTypeAlias(oldCol.DataType)
	newDT := normalizeDataTypeAlias(newCol.DataType)
	oldRank := intTypeRank(oldDT)
	newRank := intTypeRank(newDT)
	if oldRank > 0 && newRank > 0 {
		oldUnsigned := strings.Contains(strings.ToLower(oldCol.ColumnType), "unsigned")
		newUnsigned := strings.Contains(strings.ToLower(newCol.ColumnType), "unsigned")
		if oldUnsigned != newUnsigned {
			return true
		}
		return newRank < oldRank
	}
	if oldDT != newDT {
		return true
	}
	return false
}

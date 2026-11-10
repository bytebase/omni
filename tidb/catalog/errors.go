package catalog

import (
	"fmt"
	"strings"
)

type Error struct {
	Code     int
	SQLState string
	Message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ERROR %d (%s): %s", e.Code, e.SQLState, e.Message)
}

const (
	ErrDupDatabase                       = 1007
	ErrUnknownDatabase                   = 1049
	ErrDupTable                          = 1050
	ErrUnknownTable                      = 1051
	ErrDupColumn                         = 1060
	ErrWrongDBName                       = 1102
	ErrWrongTableName                    = 1103
	ErrDupKeyName                        = 1061
	ErrDupEntry                          = 1062
	ErrMultiplePriKey                    = 1068
	ErrNoSuchTable                       = 1146
	ErrNoSuchColumn                      = 1054
	ErrNoDatabaseSelected                = 1046
	ErrDupIndex                          = 1831
	ErrFKNoRefTable                      = 1824
	ErrCantDropKey                       = 1091
	ErrCheckConstraintViolated           = 3819
	ErrFKCannotDropParent                = 3730
	ErrFKMissingIndex                    = 1822
	ErrFKIncompatibleColumns             = 3780
	ErrNoSuchFunction                    = 1305
	ErrNoSuchProcedure                   = 1305
	ErrDupFunction                       = 1304
	ErrDupProcedure                      = 1304
	ErrNoSuchTrigger                     = 1360
	ErrDupTrigger                        = 1359
	ErrNoSuchEvent                       = 1539
	ErrDupEvent                          = 1537
	ErrUnsupportedGeneratedStorageChange = 3106
	ErrDependentByGenCol                 = 3108
	ErrWrongArguments                    = 1210
	ErrUnknownCharacterSet               = 1115
	ErrBlobCantHaveDefault               = 1101
	// ErrBlobKeyWithoutLength is 1170, which TiDB answers for a key part over a
	// BLOB or TEXT column that names no prefix length (read on 2026-11-09). A
	// FULLTEXT index is exempt; a JSON column gives 3152, not this.
	ErrBlobKeyWithoutLength = 1170
	// ErrUnsupportedDDL is TiDB's 8200, for DDL it parses and does not run.
	ErrUnsupportedDDL = 8200
	// ErrStoredRoutineUnsupported is TiDB's 8108. TiDB compiles CREATE and
	// DROP PROCEDURE for MySQL compatibility but will not run a stored
	// routine, rejecting it at execution. (CREATE/DROP/ALTER FUNCTION and
	// ALTER PROCEDURE never reach here: the parser refuses them, 1064.)
	ErrStoredRoutineUnsupported = 8108
	ErrUnknownCollation         = 1273
	ErrCollationCharsetMismatch = 1253
	// ErrCantAggregate2Collations, 3 and N are "Illegal mix of collations",
	// which TiDB answers when the arguments of a function or comparison do
	// not aggregate (see view_collation.go).
	ErrCantAggregate2Collations = 1267
	ErrCantAggregate3Collations = 1270
	ErrCantAggregateNCollations = 1271
	// ErrWrongNameForIndex is "Incorrect index name '<name>'", which TiDB
	// v8.5.5 answers for an UPDATE INDEXES entry naming an index the table
	// does not have (read on 2026-11-04).
	ErrWrongNameForIndex = 1280
	// ErrUnknownError is the code TiDB gives an error with no code of its own,
	// such as "range name is not supported" from ALTER RANGE.
	ErrUnknownError = 1105
	// TiDB error codes. 8239 is the closest match in TiDB's error
	// catalog for "Placement policy <name> already exists" (err.go:261
	// upstream). 8241 covers "unknown placement policy" (err.go:266).
	// 8240 covers "policy still in use" — surfaced on DROP PLACEMENT
	// POLICY when a database/table still references the policy.
	ErrPlacementPolicyExists  = 8239
	ErrPlacementPolicyInUse   = 8240
	ErrPlacementPolicyMissing = 8241
	// Column type / default validation (measured on TiDB v8.5.5, 2026-11-09).
	// 1171 rejects a PRIMARY KEY part declared explicitly NULL; 1425/1426/1427
	// bound a DECIMAL's scale/precision; 3770 rejects a DEFAULT expression whose
	// function is not on TiDB's allow list.
	ErrPrimaryCantHaveNull    = 1171
	ErrTooBigScale            = 1425
	ErrTooBigPrecision        = 1426
	ErrMBiggerThanD           = 1427
	ErrDefValFunctionNotAllow = 3770
	// FK / generated-column / drop-column over-acceptance validation (measured
	// on TiDB v8.5.5, 2026-11-09). 1830 rejects a NOT NULL referencing column
	// under ON DELETE/UPDATE SET NULL; 3733 rejects a foreign key on a VIRTUAL
	// generated column; 3102 rejects a generated-column expression that contains
	// a disallowed (non-deterministic) function; 3106 (shared with the stored-
	// status change above) rejects a VIRTUAL generated column as a primary key;
	// 1090 rejects dropping a table's only remaining column.
	ErrFKColumnCannotBeNull     = 1830
	ErrFKCannotUseVirtualColumn = 3733
	ErrGeneratedColumnFunction  = 3102
	ErrCantRemoveAllFields      = 1090
	// Temporal fsp validation (measured on TiDB v8.5.5, 2026-11-09). 1426 bounds
	// a DATETIME/TIMESTAMP/TIME fractional-seconds precision at 6; 1067 rejects a
	// temporal column whose NOW()-family DEFAULT names a precision other than the
	// column's own.
	ErrInvalidDefault = 1067
	// ErrInvalidOnUpdate is 1294 for an ON UPDATE CURRENT_TIMESTAMP clause that
	// is not valid for the column: on a non-TIMESTAMP/DATETIME column, or whose
	// fractional-seconds precision does not match the column's (measured on TiDB
	// v8.5.5, 2026-11-09).
	ErrInvalidOnUpdate = 1294
	// ErrWrongSubKey is 1089 for an index prefix length that exceeds the column's
	// own length (measured on TiDB v8.5.5, 2026-11-09).
	ErrWrongSubKey = 1089
	// Expression (functional) index restrictions TiDB v8.5.5 enforces (measured
	// 2026-11-09): 3762 rejects an expression index over a bare column, 3753 one
	// whose expression returns JSON or GEOMETRY, 3757 one that returns BLOB or
	// TEXT, and 3758 one that calls a non-deterministic / disallowed function.
	ErrFunctionalIndexOnField              = 3762
	ErrFunctionalIndexOnJSONOrGeo          = 3753
	ErrFunctionalIndexOnBlob               = 3757
	ErrFunctionalIndexFunctionIsNotAllowed = 3758
)

var sqlStateMap = map[int]string{
	ErrDupDatabase:                         "HY000",
	ErrUnknownDatabase:                     "42000",
	ErrDupTable:                            "42S01",
	ErrUnknownTable:                        "42S02",
	ErrDupColumn:                           "42S21",
	ErrWrongDBName:                         "42000",
	ErrWrongTableName:                      "42000",
	ErrWrongNameForIndex:                   "42000",
	ErrDupKeyName:                          "42000",
	ErrDupEntry:                            "23000",
	ErrMultiplePriKey:                      "42000",
	ErrNoSuchTable:                         "42S02",
	ErrNoSuchColumn:                        "42S22",
	ErrNoDatabaseSelected:                  "3D000",
	ErrDupIndex:                            "42000",
	ErrFKNoRefTable:                        "HY000",
	ErrCantDropKey:                         "42000",
	ErrCheckConstraintViolated:             "HY000",
	ErrFKCannotDropParent:                  "HY000",
	ErrFKMissingIndex:                      "HY000",
	ErrFKIncompatibleColumns:               "HY000",
	ErrNoSuchFunction:                      "42000",
	ErrDupFunction:                         "HY000",
	ErrNoSuchEvent:                         "HY000",
	ErrDupEvent:                            "HY000",
	ErrUnsupportedGeneratedStorageChange:   "HY000",
	ErrStoredRoutineUnsupported:            "HY000",
	ErrDependentByGenCol:                   "HY000",
	ErrWrongArguments:                      "HY000",
	ErrUnknownCharacterSet:                 "42000",
	ErrBlobCantHaveDefault:                 "42000",
	ErrBlobKeyWithoutLength:                "42000",
	ErrUnknownCollation:                    "HY000",
	ErrCollationCharsetMismatch:            "42000",
	ErrPlacementPolicyExists:               "HY000",
	ErrPlacementPolicyInUse:                "HY000",
	ErrUnknownError:                        "HY000",
	ErrPlacementPolicyMissing:              "HY000",
	ErrPrimaryCantHaveNull:                 "42000",
	ErrTooBigScale:                         "42000",
	ErrTooBigPrecision:                     "42000",
	ErrMBiggerThanD:                        "42000",
	ErrDefValFunctionNotAllow:              "HY000",
	ErrFKColumnCannotBeNull:                "HY000",
	ErrFKCannotUseVirtualColumn:            "HY000",
	ErrGeneratedColumnFunction:             "HY000",
	ErrCantRemoveAllFields:                 "42000",
	ErrInvalidDefault:                      "42000",
	ErrInvalidOnUpdate:                     "HY000",
	ErrWrongSubKey:                         "HY000",
	ErrFunctionalIndexOnField:              "HY000",
	ErrFunctionalIndexOnJSONOrGeo:          "HY000",
	ErrFunctionalIndexOnBlob:               "HY000",
	ErrFunctionalIndexFunctionIsNotAllowed: "HY000",
}

func sqlState(code int) string {
	if s, ok := sqlStateMap[code]; ok {
		return s
	}
	return "HY000"
}

func errDupDatabase(name string) error {
	return &Error{Code: ErrDupDatabase, SQLState: sqlState(ErrDupDatabase),
		Message: fmt.Sprintf("Can't create database '%s'; database exists", name)}
}

func errUnknownDatabase(name string) error {
	return &Error{Code: ErrUnknownDatabase, SQLState: sqlState(ErrUnknownDatabase),
		Message: fmt.Sprintf("Unknown database '%s'", name)}
}

func errDupPlacementPolicy(name string) error {
	return &Error{Code: ErrPlacementPolicyExists, SQLState: sqlState(ErrPlacementPolicyExists),
		Message: fmt.Sprintf("Placement policy '%s' already exists", name)}
}

func errUnknownPlacementPolicy(name string) error {
	return &Error{Code: ErrPlacementPolicyMissing, SQLState: sqlState(ErrPlacementPolicyMissing),
		Message: fmt.Sprintf("Unknown placement policy '%s'", name)}
}

func errRangeNameNotSupported() error {
	return &Error{Code: ErrUnknownError, SQLState: sqlState(ErrUnknownError),
		Message: "range name is not supported"}
}

func errPlacementPolicyInUse(name, ref string) error {
	return &Error{Code: ErrPlacementPolicyInUse, SQLState: sqlState(ErrPlacementPolicyInUse),
		Message: fmt.Sprintf("Placement policy '%s' is still in use by %s", name, ref)}
}

func errNoDatabaseSelected() error {
	return &Error{Code: ErrNoDatabaseSelected, SQLState: sqlState(ErrNoDatabaseSelected),
		Message: "No database selected"}
}

func errDupTable(name string) error {
	return &Error{Code: ErrDupTable, SQLState: sqlState(ErrDupTable),
		Message: fmt.Sprintf("Table '%s' already exists", name)}
}

// errDupTableQualified reports a CREATE collision the way TiDB v8.5.5 does,
// qualifying the table with its schema: "Table 'db.tbl' already exists".
func errDupTableQualified(db, name string) error {
	return errDupTable(db + "." + name)
}

func errNoSuchTable(db, name string) error {
	return &Error{Code: ErrNoSuchTable, SQLState: sqlState(ErrNoSuchTable),
		Message: fmt.Sprintf("Table '%s.%s' doesn't exist", db, name)}
}

func errWrongDBName(name string) error {
	return &Error{Code: ErrWrongDBName, SQLState: sqlState(ErrWrongDBName),
		Message: fmt.Sprintf("Incorrect database name '%s'", name)}
}

func errWrongTableName(name string) error {
	return &Error{Code: ErrWrongTableName, SQLState: sqlState(ErrWrongTableName),
		Message: fmt.Sprintf("Incorrect table name '%s'", name)}
}

// checkColumnQualifier holds what a column's qualifier must name. CREATE TABLE, MODIFY and CHANGE
// check it, and ADD, ALTER and DROP ignore it (verified on v8.5.5). The database answers first.
func checkColumnQualifier(schema, table, dbName, tableName string) error {
	if schema != "" && !strings.EqualFold(schema, dbName) {
		return errWrongDBName(schema)
	}
	if table != "" && !strings.EqualFold(table, tableName) {
		return errWrongTableName(table)
	}
	return nil
}

func errDupColumn(name string) error {
	return &Error{Code: ErrDupColumn, SQLState: sqlState(ErrDupColumn),
		Message: fmt.Sprintf("Duplicate column name '%s'", name)}
}

func errDupKeyName(name string) error {
	return &Error{Code: ErrDupKeyName, SQLState: sqlState(ErrDupKeyName),
		Message: fmt.Sprintf("Duplicate key name '%s'", name)}
}

func errMultiplePriKey() error {
	return &Error{Code: ErrMultiplePriKey, SQLState: sqlState(ErrMultiplePriKey),
		Message: "Multiple primary key defined"}
}

func errNoSuchColumn(name, context string) error {
	return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
		Message: fmt.Sprintf("Unknown column '%s' in '%s'", name, context)}
}

func errUnknownTable(db, name string) error {
	return &Error{Code: ErrUnknownTable, SQLState: sqlState(ErrUnknownTable),
		Message: fmt.Sprintf("Unknown table '%s.%s'", db, name)}
}

func errFKCannotDropParent(table, fkName, refTable string) error {
	return &Error{Code: ErrFKCannotDropParent, SQLState: sqlState(ErrFKCannotDropParent),
		Message: fmt.Sprintf("Cannot drop table '%s' referenced by a foreign key constraint '%s' on table '%s'", table, fkName, refTable)}
}

func errWrongNameForIndex(name string) error {
	return &Error{Code: ErrWrongNameForIndex, SQLState: sqlState(ErrWrongNameForIndex),
		Message: fmt.Sprintf("Incorrect index name '%s'", name)}
}

func errCantDropKey(name string) error {
	return &Error{Code: ErrCantDropKey, SQLState: sqlState(ErrCantDropKey),
		Message: fmt.Sprintf("Can't DROP '%s'; check that column/key exists", name)}
}

// errIndexDoesntExist is TiDB v8.5.5's 1091 for DROP INDEX / ALTER DROP
// INDEX|KEY of an index the table does not have: "index X doesn't exist". DROP
// COLUMN of a missing column keeps the "Can't DROP ..." wording (errCantDropKey).
// errWrongSubKey is TiDB v8.5.5's 1089 for an index prefix length that exceeds
// the column's own length.
func errWrongSubKey(_ string) error {
	return &Error{Code: ErrWrongSubKey, SQLState: sqlState(ErrWrongSubKey),
		Message: "Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys"}
}

func errIndexDoesntExist(name string) error {
	return &Error{Code: ErrCantDropKey, SQLState: sqlState(ErrCantDropKey),
		Message: fmt.Sprintf("index %s doesn't exist", name)}
}

// errDropColumnCovered is TiDB v8.5.5's 8200 for DROP COLUMN of a column that
// participates in a composite (multi-column) index or in the primary key (any
// primary key except the single integer row handle, which has its own message).
func errDropColumnCovered(col string) error {
	return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL),
		Message: fmt.Sprintf("can't drop column %s with composite index covered or Primary Key covered now", col)}
}

// errDropIntegerPrimaryKey is TiDB v8.5.5's 8200 for DROP COLUMN of the single
// integer, clustered primary key column (pkIsHandle).
func errDropIntegerPrimaryKey() error {
	return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL),
		Message: "Unsupported drop integer primary key"}
}

// errDropClusteredPrimaryKey is TiDB v8.5.5's 8200 for DROP PRIMARY KEY (or
// DROP INDEX `PRIMARY`) when the table clusters by that primary key.
func errDropClusteredPrimaryKey() error {
	return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL),
		Message: "Unsupported drop primary key when the table is using clustered index"}
}

func errFKNoRefTable(table string) error {
	return &Error{Code: ErrFKNoRefTable, SQLState: sqlState(ErrFKNoRefTable),
		Message: fmt.Sprintf("Failed to open the referenced table '%s'", table)}
}

func errFKMissingIndex(constraint, refTable string) error {
	return &Error{Code: ErrFKMissingIndex, SQLState: sqlState(ErrFKMissingIndex),
		Message: fmt.Sprintf("Failed to add the foreign key constraint. Missing index for constraint '%s' in the referenced table '%s'", constraint, refTable)}
}

func errFKIncompatibleColumns(col, refCol, constraint string) error {
	return &Error{Code: ErrFKIncompatibleColumns, SQLState: sqlState(ErrFKIncompatibleColumns),
		Message: fmt.Sprintf("Referencing column '%s' and referenced column '%s' in foreign key constraint '%s' are incompatible.", col, refCol, constraint)}
}

func errDupFunction(name string) error {
	return &Error{Code: ErrDupFunction, SQLState: sqlState(ErrDupFunction),
		Message: fmt.Sprintf("FUNCTION %s already exists", name)}
}

func errDupProcedure(name string) error {
	return &Error{Code: ErrDupProcedure, SQLState: sqlState(ErrDupProcedure),
		Message: fmt.Sprintf("PROCEDURE %s already exists", name)}
}

func errDupTrigger(name string) error {
	return &Error{Code: ErrDupTrigger, SQLState: sqlState(ErrDupTrigger),
		Message: fmt.Sprintf("Trigger already exists")}
}

func errNoSuchTrigger(db, name string) error {
	return &Error{Code: ErrNoSuchTrigger, SQLState: sqlState(ErrNoSuchTrigger),
		Message: fmt.Sprintf("Trigger does not exist")}
}

func errDupEvent(name string) error {
	return &Error{Code: ErrDupEvent, SQLState: sqlState(ErrDupEvent),
		Message: fmt.Sprintf("Event '%s' already exists", name)}
}

func errNoSuchEvent(db, name string) error {
	return &Error{Code: ErrNoSuchEvent, SQLState: sqlState(ErrNoSuchEvent),
		Message: fmt.Sprintf("Unknown event '%s'", name)}
}

func errUnsupportedGeneratedStorageChange(col, table string) error {
	return &Error{Code: ErrUnsupportedGeneratedStorageChange, SQLState: sqlState(ErrUnsupportedGeneratedStorageChange),
		Message: fmt.Sprintf("'Changing the STORED status' is not supported for generated columns.")}
}

// errStoredRoutineUnsupported is TiDB v8.5.5's rejection of a stored-routine
// statement it parses but will not run. CREATE PROCEDURE and DROP PROCEDURE
// (including DROP PROCEDURE IF EXISTS of a routine that does not exist) all
// answer this at execution; the message is TiDB's own (read 2026-11-09).
// MySQL runs stored routines, so this rejection is TiDB-only.
func errStoredRoutineUnsupported() error {
	return &Error{Code: ErrStoredRoutineUnsupported, SQLState: sqlState(ErrStoredRoutineUnsupported),
		Message: "Unsupported type *resolve.NodeW"}
}

func errDependentByGeneratedColumn(column, genColumn, table string) error {
	return &Error{Code: ErrDependentByGenCol, SQLState: sqlState(ErrDependentByGenCol),
		Message: fmt.Sprintf("Column '%s' has a generated column dependency and cannot be dropped or renamed. A generated column '%s' refers to this column in table '%s'.", column, genColumn, table)}
}

func errWrongArguments(fn string) error {
	return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments),
		Message: fmt.Sprintf("Incorrect arguments to %s", fn)}
}

// checkLiteralDefaultAllowed refuses a literal default on a BLOB, TEXT or JSON
// column, which TiDB answers 1101 for in CREATE TABLE, ADD, MODIFY and CHANGE
// COLUMN and ALTER COLUMN SET DEFAULT; DEFAULT NULL and an expression default
// in parentheses are allowed (v8.5.5, read on 2026-11-06).
func checkLiteralDefaultAllowed(col *Column) error {
	if col.Default == nil {
		return nil
	}
	v := strings.TrimSpace(*col.Default)
	if strings.EqualFold(v, "NULL") || strings.HasPrefix(v, "(") {
		return nil
	}
	switch strings.ToLower(col.DataType) {
	case "tinyblob", "blob", "mediumblob", "longblob",
		"tinytext", "text", "mediumtext", "longtext", "json":
		return &Error{Code: ErrBlobCantHaveDefault, SQLState: sqlState(ErrBlobCantHaveDefault),
			Message: fmt.Sprintf("BLOB/TEXT/JSON column '%s' can't have a default value", col.Name)}
	}
	return nil
}

func errBlobKeyWithoutLength(column string) error {
	return &Error{Code: ErrBlobKeyWithoutLength, SQLState: sqlState(ErrBlobKeyWithoutLength),
		Message: fmt.Sprintf("BLOB/TEXT column '%s' used in key specification without a key length", column)}
}

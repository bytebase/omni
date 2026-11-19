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
	ErrPartitionMgmtOnNonPartitioned     = 1505
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
	// ErrSequenceInvalidData is 4136, for sequence options that conflict, and
	// ErrUnknownSequence 4139, for DROP SEQUENCE of a missing one.
	ErrSequenceInvalidData = 4136
	ErrUnknownSequence     = 4139
	// TiDB's placement-policy codes, read on v8.5.5: a second CREATE of a
	// policy is 8238, a policy a statement names and the catalog does not
	// have is 8239 ("Unknown placement policy '<name>'"), and DROP of one a
	// database or table still references is 8241 ("Placement policy '<name>'
	// is still in use by <ref>"). 8240 is TiDB's placement-rule check, which
	// no catalog path answers.
	ErrPlacementPolicyExists  = 8238
	ErrPlacementPolicyMissing = 8239
	ErrPlacementPolicyInUse   = 8241
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
	// ErrKeyColumnDoesNotExist is 1072 for an index key part naming a column the
	// table does not have (measured on TiDB v8.5.5, 2026-11-10) — including a
	// column added later in the same multi-command ALTER, which TiDB resolves
	// against the pre-ALTER schema.
	ErrKeyColumnDoesNotExist = 1072
	// ErrIncorrectTableDefinition is 1075 for a table that declares more than
	// one AUTO_INCREMENT column (measured on TiDB v8.5.5, 2026-11-10).
	ErrIncorrectTableDefinition = 1075

	// ErrTableMustHaveColumns is 1113, which TiDB answers for a CREATE TABLE
	// that lists no column at all (measured on v8.5.5, 2026-11-19).
	ErrTableMustHaveColumns = 1113
	// ErrWrongObject is 1347, "'<db>.<name>' is not VIEW" or "is not SEQUENCE":
	// an object of the name exists and is of another kind, which DROP VIEW and
	// a sequence reference refuse even under IF EXISTS (measured on v8.5.5,
	// 2026-11-19).
	ErrWrongObject = 1347
	// ErrKeyDoesNotExist is 1176, "Key '<name>' doesn't exist in table
	// '<table>'", which ALTER TABLE RENAME KEY/INDEX and ALTER INDEX ...
	// VISIBLE/INVISIBLE answer for a name the table has no index under
	// (measured on v8.5.5, 2026-11-19).
	ErrKeyDoesNotExist = 1176
	// ErrDupEnumValue is 1291, a duplicated member in an ENUM or SET type
	// (measured on v8.5.5, 2026-11-19).
	ErrDupEnumValue = 1291
	// ErrTTLColumnType is 8148, a TTL config whose column is not a DATETIME,
	// DATE or TIMESTAMP one (measured on v8.5.5, 2026-11-19).
	ErrTTLColumnType = 8148
	// ErrAutoRandom is 8216, TiDB's "Invalid auto random" for an ADD COLUMN
	// carrying AUTO_RANDOM and for AUTO_RANDOM_BASE on a table that has no
	// auto-random column (measured on v8.5.5, 2026-11-19).
	ErrAutoRandom = 8216
	// ErrFKKeyRefMismatch is 1239 for a table-level FOREIGN KEY that names no
	// referenced column, which TiDB answers whether or not the referenced
	// table exists (measured on v8.5.5, 2026-11-19).
	ErrFKKeyRefMismatch = 1239
	// ErrPartitionFieldNotFound is 1488, a partition clause naming a column the
	// table does not have (measured on v8.5.5, 2026-11-19).
	ErrPartitionFieldNotFound = 1488
	// ErrCacheTableUnsupported is 8242, a cached table an operation cannot
	// touch: every DDL statement but NOCACHE, and CACHE itself on a
	// partitioned table (measured on v8.5.5, 2026-11-19).
	ErrCacheTableUnsupported = 8242
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
	ErrPartitionMgmtOnNonPartitioned:       "HY000",
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
	ErrKeyColumnDoesNotExist:               "42000",
	ErrIncorrectTableDefinition:            "42000",
	ErrTableMustHaveColumns:                "42000",
	ErrWrongObject:                         "HY000",
	ErrKeyDoesNotExist:                     "42000",
	ErrDupEnumValue:                        "HY000",
	ErrTTLColumnType:                       "HY000",
	ErrAutoRandom:                          "HY000",
	ErrFKKeyRefMismatch:                    "42000",
	ErrPartitionFieldNotFound:              "HY000",
	ErrCacheTableUnsupported:               "HY000",
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

// errWrongObjectKind is 1347 for a name that exists as another kind of object.
// The message names the object by its schema, as TiDB's does: `DROP VIEW a`
// over the table a answers "'d.a' is not VIEW".
func errWrongObjectKind(db, name, kind string) error {
	return &Error{Code: ErrWrongObject, SQLState: sqlState(ErrWrongObject),
		Message: fmt.Sprintf("'%s.%s' is not %s", db, name, kind)}
}

// errKeyDoesNotExist is 1176 for an index name the table has no index under.
func errKeyDoesNotExist(name, table string) error {
	return &Error{Code: ErrKeyDoesNotExist, SQLState: sqlState(ErrKeyDoesNotExist),
		Message: fmt.Sprintf("Key '%s' doesn't exist in table '%s'", name, table)}
}

// errDupEnumValue is 1291 for a duplicated ENUM or SET member.
func errDupEnumValue(col, value, kind string) error {
	return &Error{Code: ErrDupEnumValue, SQLState: sqlState(ErrDupEnumValue),
		Message: fmt.Sprintf("Column '%s' has duplicated value '%s' in %s", col, value, kind)}
}

// errTTLColumnType is 8148 for a TTL column whose type the config cannot use.
func errTTLColumnType(col string) error {
	return &Error{Code: ErrTTLColumnType, SQLState: sqlState(ErrTTLColumnType),
		Message: fmt.Sprintf("Field '%s' is of a not supported type for TTL config, expect DATETIME, DATE or TIMESTAMP", col)}
}

// errAutoRandomAddColumn is 8216 for an ADD COLUMN that carries AUTO_RANDOM,
// which TiDB refuses whatever the column's type: the row handle is fixed once
// the table exists.
func errAutoRandomAddColumn(col, db, table string) error {
	return &Error{Code: ErrAutoRandom, SQLState: sqlState(ErrAutoRandom),
		Message: fmt.Sprintf("Invalid auto random: unsupported add column '%s' constraint AUTO_RANDOM when altering '%s.%s'", col, db, table)}
}

// errAutoRandomBaseNoColumn is 8216 for AUTO_RANDOM_BASE on a table with no
// auto-random column.
func errAutoRandomBaseNoColumn() error {
	return &Error{Code: ErrAutoRandom, SQLState: sqlState(ErrAutoRandom),
		Message: "Invalid auto random: alter auto_random_base of a non auto_random table"}
}

// errFKKeyRefMismatch is 1239 for a table-level FOREIGN KEY that names no
// referenced column: TiDB requires the referenced column list there, and
// answers this whether or not the referenced table exists.
func errFKKeyRefMismatch(name string) error {
	return &Error{Code: ErrFKKeyRefMismatch, SQLState: sqlState(ErrFKKeyRefMismatch),
		Message: fmt.Sprintf("Incorrect foreign key definition for '%s': Key reference and table reference don't match", name)}
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

// errTableMustHaveColumns is TiDB 1113 for a CREATE TABLE whose element list
// holds no column. TiDB reads the definition while it prepares the statement,
// so it answers this before it looks the name up: `CREATE TABLE a`, on a
// database that already holds a table a, is 1113 rather than 1050 (read on
// v8.5.5).
func errTableMustHaveColumns() error {
	return &Error{Code: ErrTableMustHaveColumns, SQLState: sqlState(ErrTableMustHaveColumns),
		Message: "A table must have at least 1 column"}
}

// errStarNamesNoTable is TiDB 1051 for a view body's star that has no table to
// expand: table is the qualifier the star wrote, empty for a bare `*`. TiDB
// resolves the body before it looks the view's own name up, so
// `CREATE VIEW a AS SELECT *` answers this even when a table a exists.
func errStarNamesNoTable(table string) error {
	return &Error{Code: ErrUnknownTable, SQLState: sqlState(ErrUnknownTable),
		Message: fmt.Sprintf("Unknown table '%s'", table)}
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

// errPartitionMgmtOnNonPartitioned is 1505, a partition action on a table that
// has no partitions: TiDB v8.5.5 answers it with MySQL's message.
func errPartitionMgmtOnNonPartitioned() error {
	return &Error{Code: ErrPartitionMgmtOnNonPartitioned, SQLState: sqlState(ErrPartitionMgmtOnNonPartitioned),
		Message: "Partition management on a not partitioned table is not possible"}
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

// errKeyColumnDoesNotExist is TiDB v8.5.5's 1072 for an index key part over a
// column the table does not have.
func errKeyColumnDoesNotExist(col string) error {
	return &Error{Code: ErrKeyColumnDoesNotExist, SQLState: sqlState(ErrKeyColumnDoesNotExist),
		Message: fmt.Sprintf("column does not exist: %s", col)}
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

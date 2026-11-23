package catalog

import "fmt"

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
	ErrNonUniqTable                      = 1066
	ErrUpdateTableUsed                   = 1093
	ErrNoTablesUsed                      = 1096
	ErrCantReopenTable                   = 1137
	ErrOperandColumns                    = 1241
	ErrWrongObject                       = 1347
	ErrViewSelectClause                  = 1350
	ErrRowIsReferenced                   = 1451
	ErrIsAView                           = 1965
	ErrTableMustHaveColumns              = 1113
	ErrDupColumn                         = 1060
	ErrDupKeyName                        = 1061
	ErrDupEntry                          = 1062
	ErrMultiplePriKey                    = 1068
	ErrInvalidDefault                    = 1067
	ErrIncorrectTableDefinition          = 1075
	ErrNoSuchTable                       = 1146
	ErrUnknownView                       = 4092
	ErrNoSuchColumn                      = 1054
	ErrNoDatabaseSelected                = 1046
	ErrCantRemoveAllFields               = 1090
	ErrDupIndex                          = 1831
	ErrCantCreateTable                   = 1005
	ErrCantDropKey                       = 1091
	ErrCheckConstraintViolated           = 3819
	ErrFKCannotDropParent                = 3730
	ErrNoSuchFunction                    = 1305
	ErrFuncInexistentNameCollision       = 1630
	ErrFunctionNotDefined                = 1128
	ErrNoSuchProcedure                   = 1305
	ErrDupFunction                       = 1304
	ErrDupProcedure                      = 1304
	ErrNoSuchTrigger                     = 1360
	ErrDupTrigger                        = 1359
	ErrNoSuchEvent                       = 1539
	ErrTriggerInWrongSchema              = 1435
	ErrDupEvent                          = 1537
	ErrUnsupportedGeneratedStorageChange = 3106
	ErrDependentByGenCol                 = 3108
	ErrWrongArguments                    = 1210
	ErrWrongNameForIndex                 = 1280
	ErrCannotPartOfFulltext              = 1283
	ErrUnknownStorageEngine              = 1286
	ErrKeyRefMismatch                    = 1239
	ErrInvalidOnUpdate                   = 1294
	ErrTooBigPrecision                   = 1426
	// The window-reference errors, measured on mariadb:11.8.8.
	ErrWindowNoSuchWindow         = 4009
	ErrWindowDuplicateName        = 4010
	ErrWindowRefWithPartition     = 4011
	ErrWindowRefHasOrder          = 4012
	ErrWindowRefHasFrame          = 4013
	ErrWindowFrameNotAllowed      = 4017
	ErrFKDupName                  = 1826
	ErrCheckConstraintNotAllowed  = 3815
	ErrCheckConstraintDupName     = 3822
	ErrWrongSystemVersioning      = 4123
	ErrRowColumnWrongType         = 4110
	ErrNotSystemVersioned         = 4124
	ErrAlterOnSystemVersioned     = 4119
	ErrDuplicateRowColumn         = 4134
	ErrAlreadySystemVersioned     = 4135
	ErrTruncateSystemVersioned    = 4137
	ErrMissingSystemVersioning    = 4125
	ErrPeriodWrongColumns         = 4126
	ErrWrongPartitionType         = 4113
	ErrWrongSystemTimePartitions  = 4128
	ErrValuesOnlyForRange         = 1480
	ErrMultipleAppPeriods         = 4154
	ErrColumnSpecifiedTwice       = 1110
	ErrIncorrectColumnSpecifier   = 1063
	ErrPeriodNotFound             = 4156
	ErrKeyPartZeroLength          = 1391
	ErrKeyColumnDoesNotExist      = 1072
	ErrKeyIncludesPeriodColumn    = 4170
	ErrPeriodFieldsDifferentTypes = 4153
	ErrKeyDoesNotExist            = 1176
	ErrPKIndexCantBeIgnored       = 4174
	ErrWrongDBName                = 1102
)

var sqlStateMap = map[int]string{
	ErrDupDatabase:                       "HY000",
	ErrUnknownDatabase:                   "42000",
	ErrDupTable:                          "42S01",
	ErrUnknownTable:                      "42S02",
	ErrNonUniqTable:                      "42000",
	ErrUpdateTableUsed:                   "HY000",
	ErrNoTablesUsed:                      "HY000",
	ErrCantReopenTable:                   "HY000",
	ErrOperandColumns:                    "21000",
	ErrWrongObject:                       "HY000",
	ErrViewSelectClause:                  "HY000",
	ErrRowIsReferenced:                   "23000",
	ErrIsAView:                           "42S02",
	ErrTableMustHaveColumns:              "42000",
	ErrDupColumn:                         "42S21",
	ErrDupKeyName:                        "42000",
	ErrDupEntry:                          "23000",
	ErrMultiplePriKey:                    "42000",
	ErrInvalidDefault:                    "42000",
	ErrIncorrectTableDefinition:          "42000",
	ErrNoSuchTable:                       "42S02",
	ErrUnknownView:                       "42S02",
	ErrNoSuchColumn:                      "42S22",
	ErrNoDatabaseSelected:                "3D000",
	ErrCantRemoveAllFields:               "42000",
	ErrDupIndex:                          "42000",
	ErrCantCreateTable:                   "HY000",
	ErrCantDropKey:                       "42000",
	ErrCheckConstraintViolated:           "HY000",
	ErrFKCannotDropParent:                "HY000",
	ErrNoSuchFunction:                    "42000",
	ErrFuncInexistentNameCollision:       "42000",
	ErrFunctionNotDefined:                "HY000",
	ErrDupFunction:                       "HY000",
	ErrNoSuchEvent:                       "HY000",
	ErrTriggerInWrongSchema:              "HY000",
	ErrDupEvent:                          "HY000",
	ErrUnsupportedGeneratedStorageChange: "HY000",
	ErrDependentByGenCol:                 "HY000",
	ErrWrongSystemVersioning:             "HY000",
	ErrRowColumnWrongType:                "HY000",
	ErrNotSystemVersioned:                "HY000",
	ErrAlterOnSystemVersioned:            "HY000",
	ErrDuplicateRowColumn:                "HY000",
	ErrAlreadySystemVersioned:            "HY000",
	ErrTruncateSystemVersioned:           "HY000",
	ErrMissingSystemVersioning:           "HY000",
	ErrPeriodWrongColumns:                "HY000",
	ErrWrongPartitionType:                "HY000",
	ErrWrongSystemTimePartitions:         "HY000",
	ErrValuesOnlyForRange:                "HY000",
	ErrMultipleAppPeriods:                "HY000",
	ErrColumnSpecifiedTwice:              "42000",
	ErrIncorrectColumnSpecifier:          "42000",
	ErrPeriodNotFound:                    "HY000",
	ErrKeyPartZeroLength:                 "HY000",
	ErrKeyColumnDoesNotExist:             "42000",
	ErrKeyIncludesPeriodColumn:           "HY000",
	ErrPeriodFieldsDifferentTypes:        "HY000",
	ErrWrongArguments:                    "HY000",
	ErrWrongNameForIndex:                 "42000",
	ErrCannotPartOfFulltext:              "HY000",
	ErrUnknownStorageEngine:              "42000",
	ErrKeyRefMismatch:                    "42000",
	ErrInvalidOnUpdate:                   "HY000",
	ErrTooBigPrecision:                   "42000",
	ErrFKDupName:                         "HY000",
	ErrCheckConstraintNotAllowed:         "HY000",
	ErrCheckConstraintDupName:            "HY000",
	ErrKeyDoesNotExist:                   "42000",
	ErrPKIndexCantBeIgnored:              "HY000",
	ErrWrongDBName:                       "42000",
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

func errNoDatabaseSelected() error {
	return &Error{Code: ErrNoDatabaseSelected, SQLState: sqlState(ErrNoDatabaseSelected),
		Message: "No database selected"}
}

// errWrongDBName is 1102, what MariaDB answers when a 3-part field_ident's
// database qualifier does not match the target table's database (including the
// session's current database for an unqualified table).
func errWrongDBName(name string) error {
	return &Error{Code: ErrWrongDBName, SQLState: sqlState(ErrWrongDBName),
		Message: fmt.Sprintf("Incorrect database name '%s'", name)}
}

func errDupTable(name string) error {
	return &Error{Code: ErrDupTable, SQLState: sqlState(ErrDupTable),
		Message: fmt.Sprintf("Table '%s' already exists", name)}
}

func errNoSuchTable(db, name string) error {
	return &Error{Code: ErrNoSuchTable, SQLState: sqlState(ErrNoSuchTable),
		Message: fmt.Sprintf("Table '%s.%s' doesn't exist", db, name)}
}

// errUnknownView is what DROP VIEW reports for a name the database holds no
// view of, whether it holds a table of that name or nothing at all.
func errUnknownView(db, name string) error {
	return &Error{Code: ErrUnknownView, SQLState: sqlState(ErrUnknownView),
		Message: fmt.Sprintf("Unknown VIEW: '%s.%s'", db, name)}
}

// errFunctionNotDefined is 1128, what MariaDB answers for a function or a
// fulltext parser plugin it does not have.
func errFunctionNotDefined(name string) error {
	return &Error{Code: ErrFunctionNotDefined, SQLState: sqlState(ErrFunctionNotDefined),
		Message: fmt.Sprintf("Function '%s' is not defined", name)}
}

// errPartitionMgmtOnNonpartitioned is 1505, a partition action on a table with no
// partitions.
func errPartitionMgmtOnNonpartitioned() error {
	return &Error{Code: 1505, SQLState: "HY000",
		Message: "Partition management on a not partitioned table is not possible"}
}

// errOnlyOnRangeListPartitioning is 1512, an action RANGE and LIST partitioning
// alone take.
func errOnlyOnRangeListPartitioning(action string) error {
	return &Error{Code: 1512, SQLState: "HY000",
		Message: fmt.Sprintf("%s can only be used on RANGE/LIST partitions", action)}
}

// errDropPartitionNonExistent is 1507, a partition the table does not have.
func errDropPartitionNonExistent() error {
	return &Error{Code: 1507, SQLState: "HY000", Message: "Wrong partition name or partition list"}
}

// errDropLastPartition is 1508, an action that would leave the table no partition.
func errDropLastPartition() error {
	return &Error{Code: 1508, SQLState: "HY000", Message: "Cannot remove all partitions, use DROP TABLE instead"}
}

// errCoalesceOnlyOnHashPartition is 1509, COALESCE PARTITION on a non-HASH/KEY table.
func errCoalesceOnlyOnHashPartition() error {
	return &Error{Code: 1509, SQLState: "HY000", Message: "COALESCE PARTITION can only be used on HASH/KEY partitions"}
}

// errReorgNoParam is 1511, REORGANIZE PARTITION without parameters on a non-HASH/KEY table.
func errReorgNoParam() error {
	return &Error{Code: 1511, SQLState: "HY000", Message: "REORGANIZE PARTITION without parameters can only be used on auto-partitioned tables using HASH PARTITIONs"}
}

// errPartitionsMustBeDefined is 1492, ADD PARTITION without definitions on a RANGE/LIST table.
func errPartitionsMustBeDefined(kind string) error {
	return &Error{Code: 1492, SQLState: "HY000", Message: fmt.Sprintf("For %s partitions each partition must be defined", kind)}
}

// errForeignKeyOnPartitioned is 1506, foreign key operations on a partitioned table.
func errForeignKeyOnPartitioned() error {
	return &Error{Code: 1506, SQLState: "HY000", Message: "Foreign keys are not yet supported in conjunction with partitioning"}
}

// errDropIndexFK is 1553, dropping an index or column whose backing index is needed by an FK.
func errDropIndexFK(name string) error {
	return &Error{Code: 1553, SQLState: "HY000", Message: fmt.Sprintf("Cannot drop index '%s': needed in a foreign key constraint", name)}
}

// errRowIsReferencedOld is 1217, what MariaDB reports when partitioning a table
// that has or is referenced by a foreign key constraint.
func errRowIsReferencedOld() error {
	return &Error{Code: 1217, SQLState: "23000", Message: "Cannot delete or update a parent row: a foreign key constraint fails"}
}

// errChangeSchemaNotAllowed is 1450, renaming a view into another database.
func errChangeSchemaNotAllowed(fromDB, toDB string) error {
	return &Error{Code: 1450, SQLState: "HY000", Message: fmt.Sprintf("Changing schema from '%s' to '%s' is not allowed", fromDB, toDB)}
}

// errSameNamePartition is 1517, a partition name the table already has.
func errSameNamePartition(name string) error {
	return &Error{Code: 1517, SQLState: "HY000", Message: fmt.Sprintf("Duplicate partition name %s", name)}
}

// errTablesDifferentMetadata is 1736, two tables an action needs alike that are not.
func errTablesDifferentMetadata() error {
	return &Error{Code: 1736, SQLState: "HY000", Message: "Tables have different definitions"}
}

func errDupColumn(name string) error {
	return &Error{Code: ErrDupColumn, SQLState: sqlState(ErrDupColumn),
		Message: fmt.Sprintf("Duplicate column name '%s'", name)}
}

func errDupKeyName(name string) error {
	return &Error{Code: ErrDupKeyName, SQLState: sqlState(ErrDupKeyName),
		Message: fmt.Sprintf("Duplicate key name '%s'", name)}
}

func errWrongNameForIndex(name string) error {
	return &Error{Code: ErrWrongNameForIndex, SQLState: sqlState(ErrWrongNameForIndex),
		Message: fmt.Sprintf("Incorrect index name '%s'", name)}
}

func errMultiplePriKey() error {
	return &Error{Code: ErrMultiplePriKey, SQLState: sqlState(ErrMultiplePriKey),
		Message: "Multiple primary key defined"}
}

// errKeyDoesNotExist is MariaDB's answer to an ALTER INDEX that names no index
// of the table (1176, verified on 11.8.8).
func errKeyDoesNotExist(key, table string) error {
	return &Error{Code: ErrKeyDoesNotExist, SQLState: sqlState(ErrKeyDoesNotExist),
		Message: fmt.Sprintf("Key '%s' doesn't exist in table '%s'", key, table)}
}

// errPKIndexCantBeIgnored is MariaDB's answer to an ignored primary key,
// explicit or implicit (4174, verified on 11.8.8).
func errPKIndexCantBeIgnored() error {
	return &Error{Code: ErrPKIndexCantBeIgnored, SQLState: sqlState(ErrPKIndexCantBeIgnored),
		Message: "A primary key cannot be marked as IGNORE"}
}

func errInvalidDefault(name string) error {
	return &Error{Code: ErrInvalidDefault, SQLState: sqlState(ErrInvalidDefault),
		Message: fmt.Sprintf("Invalid default value for '%s'", name)}
}

func errIncorrectAutoIncrementDefinition() error {
	return &Error{Code: ErrIncorrectTableDefinition, SQLState: sqlState(ErrIncorrectTableDefinition),
		Message: "Incorrect table definition; there can be only one auto column and it must be defined as a key"}
}

func errInvalidOnUpdate(name string) error {
	return &Error{Code: ErrInvalidOnUpdate, SQLState: sqlState(ErrInvalidOnUpdate),
		Message: fmt.Sprintf("Invalid ON UPDATE clause for '%s' column", name)}
}

func errTooBigPrecision(precision int, typ string) error {
	return &Error{Code: ErrTooBigPrecision, SQLState: sqlState(ErrTooBigPrecision),
		Message: fmt.Sprintf("Too big precision %d specified for '%s'. Maximum is 6.", precision, typ)}
}

// The window-reference errors. Every message is the server's own, read off
// mariadb:11.8.8.
func errWindowNoSuchWindow(name string) error {
	return &Error{Code: ErrWindowNoSuchWindow, SQLState: sqlState(ErrWindowNoSuchWindow),
		Message: fmt.Sprintf("Window specification with name '%s' is not defined", name)}
}

func errWindowDuplicateName(name string) error {
	return &Error{Code: ErrWindowDuplicateName, SQLState: sqlState(ErrWindowDuplicateName),
		Message: fmt.Sprintf("Multiple window specifications with the same name '%s'", name)}
}

func errWindowRefWithPartition(name string) error {
	return &Error{Code: ErrWindowRefWithPartition, SQLState: sqlState(ErrWindowRefWithPartition),
		Message: fmt.Sprintf("Window specification referencing another one '%s' cannot contain partition list", name)}
}

func errWindowRefHasOrder(name string) error {
	return &Error{Code: ErrWindowRefHasOrder, SQLState: sqlState(ErrWindowRefHasOrder),
		Message: fmt.Sprintf("Referenced window specification '%s' already contains order list", name)}
}

func errWindowRefHasFrame(name string) error {
	return &Error{Code: ErrWindowRefHasFrame, SQLState: sqlState(ErrWindowRefHasFrame),
		Message: fmt.Sprintf("Referenced window specification '%s' cannot contain window frame", name)}
}

func errWindowFrameNotAllowed(name string) error {
	return &Error{Code: ErrWindowFrameNotAllowed, SQLState: sqlState(ErrWindowFrameNotAllowed),
		Message: fmt.Sprintf("Window frame is not allowed with '%s'", name)}
}

func errNoTablesUsed() error {
	return &Error{Code: ErrNoTablesUsed, SQLState: sqlState(ErrNoTablesUsed),
		Message: "No tables used"}
}

func errOperandColumns(n int) error {
	return &Error{Code: ErrOperandColumns, SQLState: sqlState(ErrOperandColumns),
		Message: fmt.Sprintf("Operand should contain %d column(s)", n)}
}

func errViewSelectClause(clause string) error {
	return &Error{Code: ErrViewSelectClause, SQLState: sqlState(ErrViewSelectClause),
		Message: fmt.Sprintf("View's SELECT contains a '%s' clause", clause)}
}

func errNoSuchColumn(name, context string) error {
	return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
		Message: fmt.Sprintf("Unknown column '%s' in '%s'", name, context)}
}

func errUnknownTable(db, name string) error {
	return &Error{Code: ErrUnknownTable, SQLState: sqlState(ErrUnknownTable),
		Message: fmt.Sprintf("Unknown table '%s.%s'", db, name)}
}

// The six errors below are the ones CREATE OR REPLACE TABLE and CREATE
// TABLE ... LIKE raise in MariaDB; messages as mariadb:11.8.8 prints them.

func errNonUniqTable(name string) error {
	return &Error{Code: ErrNonUniqTable, SQLState: sqlState(ErrNonUniqTable),
		Message: fmt.Sprintf("Not unique table/alias: '%s'", name)}
}

func errUpdateTableUsed(name, operation string) error {
	return &Error{Code: ErrUpdateTableUsed, SQLState: sqlState(ErrUpdateTableUsed),
		Message: fmt.Sprintf("Table '%s' is specified twice, both as a target for '%s' and as a separate source for data", name, operation)}
}

func errCantReopenTable(alias string) error {
	return &Error{Code: ErrCantReopenTable, SQLState: sqlState(ErrCantReopenTable),
		Message: fmt.Sprintf("Can't reopen table: '%s'", alias)}
}

func errWrongObject(db, name, kind string) error {
	return &Error{Code: ErrWrongObject, SQLState: sqlState(ErrWrongObject),
		Message: fmt.Sprintf("'%s.%s' is not of type '%s'", db, name, kind)}
}

// errRowIsReferenced is what MariaDB reports when CREATE OR REPLACE TABLE
// would drop a table another table's foreign key references.
func errRowIsReferenced() error {
	return &Error{Code: ErrRowIsReferenced, SQLState: sqlState(ErrRowIsReferenced),
		Message: "Cannot delete or update a parent row: a foreign key constraint fails"}
}

func errIsAView(db, name string) error {
	return &Error{Code: ErrIsAView, SQLState: sqlState(ErrIsAView),
		Message: fmt.Sprintf("'%s.%s' is a view", db, name)}
}

// errTableMustHaveColumns is CREATE TABLE of a definition without columns.
func errTableMustHaveColumns() error {
	return &Error{Code: ErrTableMustHaveColumns, SQLState: sqlState(ErrTableMustHaveColumns),
		Message: "A table must have at least 1 column"}
}

func errFKCannotDropParent(table, fkName, refTable string) error {
	return &Error{Code: ErrFKCannotDropParent, SQLState: sqlState(ErrFKCannotDropParent),
		Message: fmt.Sprintf("Cannot drop table '%s' referenced by a foreign key constraint '%s' on table '%s'", table, fkName, refTable)}
}

// errCantDropKey reports 1091 for a DROP of something that is not there.
// MariaDB names what was being dropped and backticks it:
// "Can't DROP COLUMN `x`; check that it exists", and the same for INDEX,
// FOREIGN KEY and CONSTRAINT. MySQL 8.0 writes one message for all of them,
// "Can't DROP 'x'; check that column/key exists". Measured on mariadb:11.8.8.
func errCantDropKey(kind, name string) error {
	return &Error{Code: ErrCantDropKey, SQLState: sqlState(ErrCantDropKey),
		Message: fmt.Sprintf("Can't DROP %s `%s`; check that it exists", kind, name)}
}

// errCannotPartOfFulltext is 1283, a FULLTEXT key over a column whose type it
// cannot index: every type but the character ones, BLOB and VARBINARY included.
func errCannotPartOfFulltext(name string) error {
	return &Error{Code: ErrCannotPartOfFulltext, SQLState: sqlState(ErrCannotPartOfFulltext),
		Message: fmt.Sprintf("Column '%s' cannot be part of FULLTEXT index", name)}
}

// errUnknownStorageEngine is 1286, an ENGINE option naming an engine the pinned
// server does not ship.
func errUnknownStorageEngine(name string) error {
	return &Error{Code: ErrUnknownStorageEngine, SQLState: sqlState(ErrUnknownStorageEngine),
		Message: fmt.Sprintf("Unknown storage engine '%s'", name)}
}

// errKeyRefMismatch is 1239, a foreign key whose local columns and referenced
// columns do not count the same: `FOREIGN KEY (a) REFERENCES z(a, a)`.
func errKeyRefMismatch(name string) error {
	return &Error{Code: ErrKeyRefMismatch, SQLState: sqlState(ErrKeyRefMismatch),
		Message: fmt.Sprintf("Incorrect foreign key definition for '%s': Key reference and table reference don't match", name)}
}

func errCantRemoveAllFields() error {
	return &Error{Code: ErrCantRemoveAllFields, SQLState: sqlState(ErrCantRemoveAllFields),
		Message: "You can't delete all columns with ALTER TABLE; use DROP TABLE instead"}
}

// errFKIncorrectlyFormed reports the one error MariaDB gives for a foreign key
// it cannot form, whatever the reason: a missing referenced table, a referenced
// table with no index over the referenced columns, or column types that do not
// match. MySQL 8.0 distinguishes them as 1824, 1822 and 3780 with a message
// naming the cause. Measured on mariadb:11.8.8: all three, and the
// ALTER TABLE ADD FOREIGN KEY form too, answer
//
//	ERROR 1005 (HY000): Can't create table `<db>`.`<table>` (errno: 150
//	"Foreign key constraint is incorrectly formed")
func errFKIncorrectlyFormed(db, table string) error {
	return &Error{Code: ErrCantCreateTable, SQLState: sqlState(ErrCantCreateTable),
		Message: fmt.Sprintf("Can't create table `%s`.`%s` (errno: 150 "+
			"\"Foreign key constraint is incorrectly formed\")", db, table)}
}

func errFKDupName(name string) error {
	return &Error{Code: ErrFKDupName, SQLState: sqlState(ErrFKDupName),
		Message: fmt.Sprintf("Duplicate foreign key constraint name '%s'", name)}
}

func errCheckConstraintDupName(name string) error {
	return &Error{Code: ErrFKDupName, SQLState: sqlState(ErrFKDupName),
		Message: fmt.Sprintf("Duplicate CHECK constraint name '%s'", name)}
}

func errCheckConstraintNotAllowed(name string) error {
	return &Error{Code: ErrCheckConstraintNotAllowed, SQLState: sqlState(ErrCheckConstraintNotAllowed),
		Message: fmt.Sprintf("An expression of a check constraint '%s' contains disallowed function.", name)}
}

func errDupFunction(name string) error {
	return &Error{Code: ErrDupFunction, SQLState: sqlState(ErrDupFunction),
		Message: fmt.Sprintf("FUNCTION %s already exists", name)}
}

func errDupProcedure(name string) error {
	return &Error{Code: ErrDupProcedure, SQLState: sqlState(ErrDupProcedure),
		Message: fmt.Sprintf("PROCEDURE %s already exists", name)}
}

func errNoSuchFunction(name string) error {
	return &Error{Code: ErrNoSuchFunction, SQLState: sqlState(ErrNoSuchFunction),
		Message: fmt.Sprintf("FUNCTION %s does not exist", name)}
}

func errFuncInexistentNameCollision(name string) error {
	return &Error{Code: ErrFuncInexistentNameCollision, SQLState: sqlState(ErrFuncInexistentNameCollision),
		Message: fmt.Sprintf("FUNCTION %s does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual", name)}
}

func errNoSuchProcedure(db, name string) error {
	return &Error{Code: ErrNoSuchProcedure, SQLState: sqlState(ErrNoSuchProcedure),
		Message: fmt.Sprintf("PROCEDURE %s.%s does not exist", db, name)}
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

func errTriggerInWrongSchema() error {
	return &Error{Code: ErrTriggerInWrongSchema, SQLState: sqlState(ErrTriggerInWrongSchema), Message: "Trigger in wrong schema"}
}

func errNoSuchEvent(db, name string) error {
	return &Error{Code: ErrNoSuchEvent, SQLState: sqlState(ErrNoSuchEvent),
		Message: fmt.Sprintf("Unknown event '%s'", name)}
}

func errUnsupportedGeneratedStorageChange(col, table string) error {
	return &Error{Code: ErrUnsupportedGeneratedStorageChange, SQLState: sqlState(ErrUnsupportedGeneratedStorageChange),
		Message: fmt.Sprintf("'Changing the STORED status' is not supported for generated columns.")}
}

func errMissingSystemVersioning(table string) error {
	return &Error{Code: ErrMissingSystemVersioning, SQLState: sqlState(ErrMissingSystemVersioning),
		Message: fmt.Sprintf("Wrong parameters for `%s`: missing 'WITH SYSTEM VERSIONING'", table)}
}

func errPeriodWrongColumns(rowStart, rowEnd string) error {
	return &Error{Code: ErrPeriodWrongColumns, SQLState: sqlState(ErrPeriodWrongColumns),
		Message: fmt.Sprintf("PERIOD FOR SYSTEM_TIME must use columns `%s` and `%s`", rowStart, rowEnd)}
}

func errMissingRowStart(table string) error {
	return &Error{Code: ErrMissingSystemVersioning, SQLState: sqlState(ErrMissingSystemVersioning),
		Message: fmt.Sprintf("Wrong parameters for `%s`: missing 'AS ROW START'", table)}
}

func errMissingPeriod(table string) error {
	return &Error{Code: ErrMissingSystemVersioning, SQLState: sqlState(ErrMissingSystemVersioning),
		Message: fmt.Sprintf("Wrong parameters for `%s`: missing 'PERIOD FOR SYSTEM_TIME'", table)}
}

func errWrongPartitionType(byType string) error {
	return &Error{Code: ErrWrongPartitionType, SQLState: sqlState(ErrWrongPartitionType),
		Message: fmt.Sprintf("Wrong partition type `SYSTEM_TIME` for partitioning by `%s`", byType)}
}

func errWrongSystemTimePartitions(table string) error {
	return &Error{Code: ErrWrongSystemTimePartitions, SQLState: sqlState(ErrWrongSystemTimePartitions),
		Message: fmt.Sprintf("Wrong partitions for `%s`: must have at least one HISTORY and exactly one last CURRENT", table)}
}

func errValuesOnlyForRange() error {
	return &Error{Code: ErrValuesOnlyForRange, SQLState: sqlState(ErrValuesOnlyForRange),
		Message: "Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition"}
}

func errMultipleAppPeriods() error {
	return &Error{Code: ErrMultipleAppPeriods, SQLState: sqlState(ErrMultipleAppPeriods),
		Message: "Cannot specify more than one application-time period"}
}

// errKeyPartZeroLength reports a key part declared with an explicit zero-length
// prefix, e.g. `KEY (c(0))`. MariaDB rejects this with 1391 on every key path
// and column type; a zero-length prefix is distinct from no prefix at all.
func errKeyPartZeroLength(name string) error {
	return &Error{Code: ErrKeyPartZeroLength, SQLState: sqlState(ErrKeyPartZeroLength),
		Message: fmt.Sprintf("Key part '%s' length cannot be 0", name)}
}

// errKeyColumnDoesNotExist reports a key referencing a name that is not a column
// (a nonexistent column, or a period name used as an ordinary key part). 1072.
func errKeyColumnDoesNotExist(name string) error {
	return &Error{Code: ErrKeyColumnDoesNotExist, SQLState: sqlState(ErrKeyColumnDoesNotExist),
		Message: fmt.Sprintf("Key column '%s' doesn't exist in table", name)}
}

// errKeyIncludesPeriodColumn reports an application-time period column used as an
// ordinary part of a WITHOUT OVERLAPS key. 4170.
func errKeyIncludesPeriodColumn(keyName, colName string) error {
	return &Error{Code: ErrKeyIncludesPeriodColumn, SQLState: sqlState(ErrKeyIncludesPeriodColumn),
		Message: fmt.Sprintf("Key `%s` cannot explicitly include column `%s`", keyName, colName)}
}

func errColumnSpecifiedTwice(col string) error {
	return &Error{Code: ErrColumnSpecifiedTwice, SQLState: sqlState(ErrColumnSpecifiedTwice),
		Message: fmt.Sprintf("Column '%s' specified twice", col)}
}

// errCantDropPeriod reports DROP PERIOD for a period that does not exist. 1091.
func errCantDropPeriod(name string) error {
	return &Error{Code: ErrCantDropKey, SQLState: sqlState(ErrCantDropKey),
		Message: fmt.Sprintf("Can't DROP PERIOD `%s`; check that it exists", name)}
}

// errPeriodFieldsDifferentTypes reports that a period's start and end columns
// have different types — they must match exactly, including precision. 4153.
func errPeriodFieldsDifferentTypes(name string) error {
	return &Error{Code: ErrPeriodFieldsDifferentTypes, SQLState: sqlState(ErrPeriodFieldsDifferentTypes),
		Message: fmt.Sprintf("Fields of PERIOD FOR `%s` have different types", name)}
}

func errPeriodNotFound(name string) error {
	return &Error{Code: ErrPeriodNotFound, SQLState: sqlState(ErrPeriodNotFound),
		Message: fmt.Sprintf("Period `%s` is not found in table", name)}
}

func errIncorrectColumnSpecifier(col string) error {
	return &Error{Code: ErrIncorrectColumnSpecifier, SQLState: sqlState(ErrIncorrectColumnSpecifier),
		Message: fmt.Sprintf("Incorrect column specifier for column '%s'", col)}
}

func errNotSystemVersioned(table string) error {
	return &Error{Code: ErrNotSystemVersioned, SQLState: sqlState(ErrNotSystemVersioned),
		Message: fmt.Sprintf("Table `%s` is not system-versioned", table)}
}

func errTruncateSystemVersioned(table string) error {
	return &Error{Code: ErrTruncateSystemVersioned, SQLState: sqlState(ErrTruncateSystemVersioned),
		Message: fmt.Sprintf("System-versioned tables do not support TRUNCATE TABLE, table `%s`", table)}
}

func errAlreadySystemVersioned(table string) error {
	return &Error{Code: ErrAlreadySystemVersioned, SQLState: sqlState(ErrAlreadySystemVersioned),
		Message: fmt.Sprintf("Table `%s` is already system-versioned", table)}
}

func errAlterColumnOnSystemVersioned(table string) error {
	return &Error{Code: ErrAlterOnSystemVersioned, SQLState: sqlState(ErrAlterOnSystemVersioned),
		Message: fmt.Sprintf("Not allowed for system-versioned `%s`. Change @@system_versioning_alter_history to proceed with ALTER.", table)}
}

func errNoVersionedColumn(table string) error {
	return &Error{Code: ErrWrongSystemVersioning, SQLState: sqlState(ErrWrongSystemVersioning),
		Message: fmt.Sprintf("Table `%s` must have at least one versioned column", table)}
}

func errDuplicateRowColumn(bound, col string) error {
	return &Error{Code: ErrDuplicateRowColumn, SQLState: sqlState(ErrDuplicateRowColumn),
		Message: fmt.Sprintf("Duplicate %s column `%s`", bound, col)}
}

func errWrongRowColumnType(table, col string) error {
	return &Error{Code: ErrRowColumnWrongType, SQLState: sqlState(ErrRowColumnWrongType),
		Message: fmt.Sprintf("`%s` in `%s` must be of type TIMESTAMP(6) or BIGINT(20) UNSIGNED", col, table)}
}

func errDependentByGeneratedColumn(column, genColumn, table string) error {
	return &Error{Code: ErrDependentByGenCol, SQLState: sqlState(ErrDependentByGenCol),
		Message: fmt.Sprintf("Column '%s' has a generated column dependency and cannot be dropped or renamed. A generated column '%s' refers to this column in table '%s'.", column, genColumn, table)}
}

func errCannotDropColumnFunctionalIndex(column string) error {
	return &Error{Code: ErrDependentByGenCol, SQLState: sqlState(ErrDependentByGenCol),
		Message: fmt.Sprintf("Cannot drop column '%s' because it is used by a functional index.", column)}
}

func errWrongArguments(fn string) error {
	return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments),
		Message: fmt.Sprintf("Incorrect arguments to %s", fn)}
}

// errUnknownOption is 1911: an engine-defined option the table's engine does
// not define. Measured on mariadb:11.8.8.
func errUnknownOption(name string) error {
	return &Error{Code: 1911, SQLState: "HY000", Message: fmt.Sprintf("Unknown option '%s'", name)}
}

// errTablespaceExists is 1813, ALTER TABLE ... IMPORT TABLESPACE over a table
// whose tablespace was not DISCARDed first, which is every table the server
// created itself: measured on mariadb:11.8.8, `ALTER TABLE a IMPORT TABLESPACE`
// answers "Tablespace for table 'db/a' exists. Please DISCARD the tablespace
// before IMPORT".
func errTablespaceExists(db, table string) error {
	return &Error{Code: 1813, SQLState: "HY000",
		Message: fmt.Sprintf("Tablespace for table '%s/%s' exists. Please DISCARD the tablespace before IMPORT", db, table)}
}

// refuseTablespaceEngine refuses DISCARD and IMPORT TABLESPACE on a table that
// has no InnoDB file-per-table tablespace: 1031 for MyISAM, Aria, MEMORY and
// partitioned InnoDB tables, and 3007 for an InnoDB temporary table (after the
// engine check, so a MyISAM temporary table still answers 1031). Measured on
// mariadb:11.8.8.
func refuseTablespaceEngine(tbl *Table) error {
	// ENGINE='DEFAULT' is the server's own engine, InnoDB on the pin.
	engine := tbl.Engine
	if engine == "" || eqFoldStr(engine, "InnoDB") || eqFoldStr(engine, "DEFAULT") {
		if tbl.Partitioning != nil {
			return &Error{Code: 1031, SQLState: "HY000",
				Message: fmt.Sprintf("Storage engine InnoDB of the table `%s`.`%s` doesn't have this option", tbl.Database.Name, tbl.Name)}
		}
		if tbl.Temporary {
			return &Error{Code: 3007, SQLState: "HY000",
				Message: "Cannot DISCARD/IMPORT tablespace associated with temporary table"}
		}
		return nil
	}
	return &Error{Code: 1031, SQLState: "HY000",
		Message: fmt.Sprintf("Storage engine %s of the table `%s`.`%s` doesn't have this option", engine, tbl.Database.Name, tbl.Name)}
}

// errTablespaceDiscarded is 1814, an ALTER that needs a copy of a table whose
// tablespace the statement DISCARDed. Measured on mariadb:11.8.8.
func errTablespaceDiscarded(table string) error {
	return &Error{Code: 1814, SQLState: "HY000",
		Message: fmt.Sprintf("Tablespace has been discarded for table `%s`", table)}
}

// errTablespaceMissing is 1030, IMPORT TABLESPACE over a table whose tablespace
// the statement DISCARDed and which holds no secondary index: there is no file
// to import, which the server reports as storage engine error 194. Measured on
// mariadb:11.8.8.
func errTablespaceMissing() error {
	return &Error{Code: 1030, SQLState: "HY000",
		Message: "Got error 194 \"Tablespace is missing for a table\" from storage engine InnoDB"}
}

// errImportNeedsIndexesDropped is 1815, IMPORT TABLESPACE over a table whose
// tablespace the statement DISCARDed and which still holds a secondary index:
// the import reads the .cfg file DISCARD removed, and asks for the indexes to go
// first. Measured on mariadb:11.8.8.
func errImportNeedsIndexesDropped(db, table string) error {
	return &Error{Code: 1815, SQLState: "HY000",
		Message: fmt.Sprintf("Internal error: Drop all secondary indexes before importing table %s/%s when .cfg file is missing.", db, table)}
}

// errDupColumnInKey is 1060 for a key naming a column twice: measured on
// mariadb:11.8.8, `ALTER TABLE a ADD KEY(b, b)` answers "Duplicate column name
// 'b'", as CREATE TABLE with the same key does.
func errDupColumnInKey(name string) error {
	return &Error{Code: ErrDupColumn, SQLState: sqlState(ErrDupColumn),
		Message: fmt.Sprintf("Duplicate column name '%s'", name)}
}

// errKeyPartNotString is 1089 for a key part with a length over a column that
// takes none: measured on mariadb:11.8.8, `ADD KEY(b(1))` over an INT column b
// answers it, as the same key part in CREATE TABLE and CREATE INDEX does.
func errKeyPartNotString() error {
	return &Error{Code: 1089, SQLState: "HY000",
		Message: "Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys"}
}

// errDupEnumValue is 1291: an ENUM or SET column listing a value twice.
// MariaDB compares the values case-insensitively and names the first spelling.
// Measured on mariadb:11.8.8: `CREATE TABLE t (e ENUM('s','S'))` answers
// "Column 'e' has duplicated value 's' in ENUM".
func errDupEnumValue(column, value, kind string) error {
	return &Error{Code: 1291, SQLState: "HY000",
		Message: fmt.Sprintf("Column '%s' has duplicated value '%s' in %s", column, value, kind)}
}

// errSequenceInvalidTableStructure is 4086: SEQUENCE=1 on a table whose columns
// are not a sequence's. Measured on mariadb:11.8.8.
func errSequenceInvalidTableStructure(db, table string) error {
	return &Error{Code: 4086, SQLState: "HY000",
		Message: fmt.Sprintf("Sequence '%s.%s' table structure is invalid (Wrong number of columns)", db, table)}
}

// errSequenceOutOfRange is 4085: sequence options that leave no valid value,
// such as a start below the minimum. Measured on mariadb:11.8.8.
func errSequenceOutOfRange(db, name string) error {
	return &Error{Code: 4085, SQLState: "HY000",
		Message: fmt.Sprintf("Sequence '%s.%s' has out of range value for options", db, name)}
}

// errNotSequence is DROP SEQUENCE naming a table, 4090, or ALTER SEQUENCE
// doing so, 4089. Measured on mariadb:11.8.8.
func errNotSequence(code int, db, name string) error {
	return &Error{Code: code, SQLState: "42S02",
		Message: fmt.Sprintf("'%s.%s' is not a SEQUENCE", db, name)}
}

// errUnknownSequence is 4091: DROP or ALTER SEQUENCE naming nothing. DROP
// qualifies the name and ALTER does not. Measured on mariadb:11.8.8.
func errUnknownSequence(name string) error {
	return &Error{Code: 4091, SQLState: "42S02",
		Message: fmt.Sprintf("Unknown SEQUENCE: '%s'", name)}
}

// errEventSameName is 1551: ALTER EVENT ... RENAME TO naming the same event.
// Measured on mariadb:11.8.8.
func errEventSameName() error {
	return &Error{Code: 1551, SQLState: "HY000", Message: "Same old and new event name"}
}

// errOutOfRangeColumn is 1264: a numeric value out of range for a column
// (e.g. ALTER SEQUENCE narrowing AS <type> with an out-of-range existing limit).
// Measured on mariadb:11.8.8.
func errOutOfRangeColumn(col string) error {
	return &Error{Code: 1264, SQLState: "22003",
		Message: fmt.Sprintf("Out of range value for column '%s' at row 1", col)}
}

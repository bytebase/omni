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
	ErrDupKeyName                        = 1061
	ErrDupEntry                          = 1062
	ErrNonUniqueTable                    = 1066
	ErrKeyDoesNotExist                   = 1176
	ErrIncorrectUsage                    = 1221
	ErrInvalidEncryption                 = 3184
	ErrStartTransactionInAlter           = 3979
	ErrMultiplePriKey                    = 1068
	ErrInvalidDefault                    = 1067
	ErrIncorrectTableDefinition          = 1075
	ErrTableStorageEngineOption          = 1031
	ErrUnknownStorageEngine              = 1286
	ErrInvalidPerformanceSchema          = 1683
	ErrNoSuchTable                       = 1146
	ErrNoSuchColumn                      = 1054
	ErrKeyColumnDoesNotExist             = 1072
	ErrFunctionNotDefined                = 1128
	ErrWrongObject                       = 1347
	ErrTooBigDisplayWidth                = 1439
	ErrPartitionMgmtOnNonPartitioned     = 1505
	ErrNoDatabaseSelected                = 1046
	ErrCantRemoveAllFields               = 1090
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
	ErrTriggerInWrongSchema              = 1435
	ErrDupEvent                          = 1537
	ErrFKCannotUseVirtualColumn          = 3104
	ErrUnsupportedGeneratedStorageChange = 3106
	ErrDependentByGenCol                 = 3108
	ErrWrongArguments                    = 1210
	ErrCannotPartOfFulltext              = 1283
	ErrWrongNameForIndex                 = 1280
	ErrInvalidOnUpdate                   = 1294
	ErrTooBigPrecision                   = 1426
	ErrInvalidYearColumnLength           = 1818
	ErrFKDupName                         = 1826
	ErrCheckConstraintNotAllowed         = 3815
	ErrCheckNotBoolean                   = 3812
	ErrCheckConstraintDupName            = 3822

	// ErrNoTablesUsed is 1096, which mysql:8.0.46 answers for a SELECT whose
	// star has no table to expand, a CREATE VIEW body included.
	ErrNoTablesUsed = 1096
	// ErrTableMustHaveColumns is 4028, which mysql:8.0.46 answers for a CREATE
	// TABLE that lists no column at all.
	ErrTableMustHaveColumns        = 4028
	ErrPartitionFieldNotFound      = 1488
	ErrPartitionSubpartMix         = 1500
	ErrDupPartitionField           = 1652
	ErrSpatialColumnNotGeometrical = 1687
	ErrSRSMissingAttribute         = 3708

	ErrIncorrectPrefixKey                = 1089
	ErrNotAllowedCommand                 = 1148
	ErrWrongFkDef                        = 1239
	ErrDuplicatedValueInType             = 1291
	ErrTablespaceMissing                 = 1812
	ErrTablespaceExists                  = 1813
	ErrTablespaceDiscarded               = 1814
	ErrCannotDiscardImportTempTable      = 3007
	ErrWrongTablespaceName               = 3119
	ErrTablespaceIsNotEmpty              = 3120
	ErrJSONUsedAsKey                     = 3152
	ErrTablespaceDoesNotExist            = 3510
	ErrReservedTablespaceName            = 3723
	ErrFunctionalIndexPrimaryKey         = 3756
	ErrFunctionalIndexOnField            = 3762
	ErrDependentByDefaultVal             = 3766
	ErrDefaultValExprNamedFuncNotAllowed = 3767
	ErrDefaultValExprRefAutoInc          = 3768
	ErrDefaultValExprContainSubquery     = 3769
	ErrDefaultValExprRefVariable         = 3772
	ErrSecondaryEngineNotDefined         = 3889
	ErrAutoextendSizeNotMultiple         = 4023
	ErrAutoextendSizeOutOfRange          = 4025
)

const (
	// ErrCheckConstraintNotFound is 3821, which mysql:8.0.46 answers when a
	// DROP CHECK or ALTER CHECK names no check constraint of the table:
	// "Check constraint '%s' is not found in the table." The CONSTRAINT
	// spelling of the same name answers ErrConstraintDoesNotExist instead.
	ErrCheckConstraintNotFound = 3821
	// ErrConstraintDoesNotExist is 3940, which mysql:8.0.46 answers when a
	// DROP CONSTRAINT or ALTER CONSTRAINT names no constraint of the table.
	ErrConstraintDoesNotExist = 3940
)

var sqlStateMap = map[int]string{
	ErrDupDatabase:                       "HY000",
	ErrUnknownDatabase:                   "42000",
	ErrDupTable:                          "42S01",
	ErrUnknownTable:                      "42S02",
	ErrNoTablesUsed:                      "HY000",
	ErrDupColumn:                         "42S21",
	ErrDupKeyName:                        "42000",
	ErrDupEntry:                          "23000",
	ErrNonUniqueTable:                    "42000",
	ErrKeyDoesNotExist:                   "42000",
	ErrIncorrectUsage:                    "HY000",
	ErrInvalidEncryption:                 "HY000",
	ErrStartTransactionInAlter:           "HY000",
	ErrMultiplePriKey:                    "42000",
	ErrInvalidDefault:                    "42000",
	ErrIncorrectTableDefinition:          "42000",
	ErrTableStorageEngineOption:          "HY000",
	ErrUnknownStorageEngine:              "42000",
	ErrInvalidPerformanceSchema:          "HY000",
	ErrNoSuchTable:                       "42S02",
	ErrNoSuchColumn:                      "42S22",
	ErrKeyColumnDoesNotExist:             "42000",
	ErrFunctionNotDefined:                "HY000",
	ErrWrongObject:                       "HY000",
	ErrTooBigDisplayWidth:                "42000",
	ErrPartitionMgmtOnNonPartitioned:     "HY000",
	ErrNoDatabaseSelected:                "3D000",
	ErrCantRemoveAllFields:               "42000",
	ErrDupIndex:                          "42000",
	ErrFKNoRefTable:                      "HY000",
	ErrCantDropKey:                       "42000",
	ErrCheckConstraintViolated:           "HY000",
	ErrFKCannotDropParent:                "HY000",
	ErrFKMissingIndex:                    "HY000",
	ErrFKIncompatibleColumns:             "HY000",
	ErrNoSuchFunction:                    "42000",
	ErrDupFunction:                       "HY000",
	ErrNoSuchEvent:                       "HY000",
	ErrTriggerInWrongSchema:              "HY000",
	ErrDupEvent:                          "HY000",
	ErrFKCannotUseVirtualColumn:          "HY000",
	ErrUnsupportedGeneratedStorageChange: "HY000",
	ErrDependentByGenCol:                 "HY000",
	ErrWrongArguments:                    "HY000",
	ErrCannotPartOfFulltext:              "HY000",
	ErrWrongNameForIndex:                 "42000",
	ErrInvalidOnUpdate:                   "HY000",
	ErrTooBigPrecision:                   "42000",
	ErrInvalidYearColumnLength:           "HY000",
	ErrFKDupName:                         "HY000",
	ErrCheckConstraintNotAllowed:         "HY000",
	ErrCheckNotBoolean:                   "HY000",
	ErrCheckConstraintDupName:            "HY000",
	ErrCheckConstraintNotFound:           "HY000",
	ErrConstraintDoesNotExist:            "HY000",
	ErrTableMustHaveColumns:              "HY000",
	ErrSpatialColumnNotGeometrical:       "42000",
	ErrSRSMissingAttribute:               "SR006",
	ErrIncorrectPrefixKey:                "HY000",
	ErrNotAllowedCommand:                 "42000",
	ErrWrongFkDef:                        "42000",
	ErrDuplicatedValueInType:             "HY000",
	ErrTablespaceMissing:                 "HY000",
	ErrTablespaceExists:                  "HY000",
	ErrTablespaceDiscarded:               "HY000",
	ErrCannotDiscardImportTempTable:      "HY000",
	ErrWrongTablespaceName:               "42000",
	ErrTablespaceIsNotEmpty:              "HY000",
	ErrJSONUsedAsKey:                     "42000",
	ErrTablespaceDoesNotExist:            "HY000",
	ErrReservedTablespaceName:            "HY000",
	ErrFunctionalIndexPrimaryKey:         "HY000",
	ErrFunctionalIndexOnField:            "HY000",
	ErrDependentByDefaultVal:             "HY000",
	ErrDefaultValExprNamedFuncNotAllowed: "HY000",
	ErrDefaultValExprRefAutoInc:          "HY000",
	ErrDefaultValExprContainSubquery:     "HY000",
	ErrDefaultValExprRefVariable:         "HY000",
	ErrSecondaryEngineNotDefined:         "HY000",
	ErrAutoextendSizeNotMultiple:         "HY000",
	ErrAutoextendSizeOutOfRange:          "HY000",
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

func errDupTable(name string) error {
	return &Error{Code: ErrDupTable, SQLState: sqlState(ErrDupTable),
		Message: fmt.Sprintf("Table '%s' already exists", name)}
}

func errNoSuchTable(db, name string) error {
	return &Error{Code: ErrNoSuchTable, SQLState: sqlState(ErrNoSuchTable),
		Message: fmt.Sprintf("Table '%s.%s' doesn't exist", db, name)}
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

// errNoTablesUsed is MySQL 1096 for a star with no table to expand. mysql's
// CREATE VIEW resolves the body before the view's own name, so this answers
// ahead of a name collision (1050). A star that qualifies its table is not this
// error: MySQL names the table it cannot find (1051).
func errNoTablesUsed() error {
	return &Error{Code: ErrNoTablesUsed, SQLState: sqlState(ErrNoTablesUsed),
		Message: "No tables used"}
}

// errUnknownTableName is MySQL 1051 as a view body's qualified star words it:
// `CREATE VIEW v AS SELECT t.*`, with no FROM naming t, is "Unknown table 't'"
// (mysql:8.0.46).
func errUnknownTableName(name string) error {
	return &Error{Code: ErrUnknownTable, SQLState: sqlState(ErrUnknownTable),
		Message: fmt.Sprintf("Unknown table '%s'", name)}
}

// errTableMustHaveColumns is MySQL 4028 for a CREATE TABLE whose element list
// holds no column. mysql:8.0.46 answers it after the name collision (1050) and
// after a key's columns are resolved (1072).
func errTableMustHaveColumns() error {
	return &Error{Code: ErrTableMustHaveColumns, SQLState: sqlState(ErrTableMustHaveColumns),
		Message: "A table must have at least one visible column."}
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

func errInvalidYearColumnLength() error {
	return &Error{Code: ErrInvalidYearColumnLength, SQLState: sqlState(ErrInvalidYearColumnLength),
		Message: "Invalid YEAR column length"}
}

func errNoSuchColumn(name, context string) error {
	return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
		Message: fmt.Sprintf("Unknown column '%s' in '%s'", name, context)}
}

// errPartitionFieldNotFound is 1488, a column a KEY partitioning names and the
// table does not hold. MySQL words it without the name.
func errPartitionFieldNotFound() error {
	return &Error{Code: ErrPartitionFieldNotFound, SQLState: sqlState(ErrPartitionFieldNotFound),
		Message: "Field in list of fields for partition function not found in table"}
}

// errDupPartitionField is 1652, the same column named twice in one partition
// function.
func errDupPartitionField(name string) error {
	return &Error{Code: ErrDupPartitionField, SQLState: sqlState(ErrDupPartitionField),
		Message: fmt.Sprintf("Duplicate partition field name '%s'", name)}
}

// errPartitionSubpartMix is 1500, a HASH or KEY partitioning that also writes a
// SUBPARTITION BY: MySQL mixes subpartitioning with RANGE and LIST only.
func errPartitionSubpartMix() error {
	return &Error{Code: ErrPartitionSubpartMix, SQLState: sqlState(ErrPartitionSubpartMix),
		Message: "It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning"}
}

// errPartitionsMustBeDefined is 1492, a RANGE or LIST partitioning that defines
// no partition.
func errPartitionsMustBeDefined() error {
	return &Error{Code: 1492, SQLState: sqlState(1492),
		Message: "Partitions must be defined for RANGE/LIST partitioning"}
}

// errSpatialNotGeometrical is 1687, a SPATIAL key over a column of another type.
func errSpatialNotGeometrical() error {
	return &Error{Code: ErrSpatialColumnNotGeometrical, SQLState: sqlState(ErrSpatialColumnNotGeometrical),
		Message: "A SPATIAL index may only contain a geometrical type column"}
}

// errSRSMandatoryAttribute is 3708, a CREATE SPATIAL REFERENCE SYSTEM that
// leaves out the attribute MySQL requires.
func errSRSMandatoryAttribute(name string) error {
	return &Error{Code: ErrSRSMissingAttribute, SQLState: sqlState(ErrSRSMissingAttribute),
		Message: fmt.Sprintf("Missing mandatory attribute %s.", name)}
}

func errUnknownTable(db, name string) error {
	return &Error{Code: ErrUnknownTable, SQLState: sqlState(ErrUnknownTable),
		Message: fmt.Sprintf("Unknown table '%s.%s'", db, name)}
}

// errWrongObject is what ALTER VIEW and DROP VIEW answer for a name a table
// holds: mysql:8.0.46 says "'db.a' is not VIEW".
func errWrongObject(db, name, kind string) error {
	return &Error{Code: ErrWrongObject, SQLState: sqlState(ErrWrongObject),
		Message: fmt.Sprintf("'%s.%s' is not %s", db, name, kind)}
}

// errNonUniqueTable is 1066, the same table named twice in one statement.
func errNonUniqueTable(name string) error {
	return &Error{Code: ErrNonUniqueTable, SQLState: sqlState(ErrNonUniqueTable),
		Message: fmt.Sprintf("Not unique table/alias: '%s'", name)}
}

// errKeyDoesNotExist is 1176, an ALTER INDEX, RENAME KEY or DROP KEY naming a
// key the table does not hold. mysql:8.0.46 words it "Key '%s' doesn't exist in
// table '%s'"; 1061 is a name collision instead.
func errKeyDoesNotExist(name, table string) error {
	return &Error{Code: ErrKeyDoesNotExist, SQLState: sqlState(ErrKeyDoesNotExist),
		Message: fmt.Sprintf("Key '%s' doesn't exist in table '%s'", name, table)}
}

// errIncorrectSRIDUsage is 1221 for a SRID clause on a column whose type is not
// a geometry, answered while the statement is prepared.
func errIncorrectSRIDUsage() error {
	return &Error{Code: ErrIncorrectUsage, SQLState: sqlState(ErrIncorrectUsage),
		Message: "Incorrect usage of SRID and non-geometry column"}
}

// errIncorrectValidationUsage is 1221 for WITH VALIDATION or WITHOUT VALIDATION
// outside EXCHANGE PARTITION, answered while the statement is prepared.
func errIncorrectValidationUsage() error {
	return &Error{Code: ErrIncorrectUsage, SQLState: sqlState(ErrIncorrectUsage),
		Message: "Incorrect usage of ALTER and WITH VALIDATION"}
}

// errInvalidEncryptionOption is 3184, an ENCRYPTION value other than 'Y' or
// 'N' on ALTER TABLE.
func errInvalidEncryptionOption() error {
	return &Error{Code: ErrInvalidEncryption, SQLState: sqlState(ErrInvalidEncryption),
		Message: "Invalid encryption option."}
}

// errStartTransactionInAlter is 3979, START TRANSACTION as an ALTER TABLE
// option; MySQL answers it before the table it names is looked up.
func errStartTransactionInAlter() error {
	return &Error{Code: ErrStartTransactionInAlter, SQLState: sqlState(ErrStartTransactionInAlter),
		Message: "START TRANSACTION clause cannot be used with ALTER TABLE command."}
}

// errCheckNotBoolean is 3812 for a CHECK expression whose top level is not a
// boolean: mysql:8.0.46 refuses a bare column, literal, arithmetic or CAST and
// names the constraint the statement would create.
func errCheckNotBoolean(name string) error {
	return &Error{Code: ErrCheckNotBoolean, SQLState: sqlState(ErrCheckNotBoolean),
		Message: fmt.Sprintf("An expression of non-boolean type specified to a check constraint '%s'.", name)}
}

// errUnknownStorageEngine is 1286, an ENGINE option naming an engine the pinned
// server does not ship.
func errUnknownStorageEngine(name string) error {
	return &Error{Code: ErrUnknownStorageEngine, SQLState: sqlState(ErrUnknownStorageEngine),
		Message: fmt.Sprintf("Unknown storage engine '%s'", name)}
}

// errTableStorageEngineOption is 1031, what ALTER TABLE answers for an ENGINE it
// cannot put on a table that exists: mysql:8.0.46 says "Table storage engine for
// 'a' doesn't have this option" for PERFORMANCE_SCHEMA.
func errTableStorageEngineOption(table string) error {
	return &Error{Code: ErrTableStorageEngineOption, SQLState: sqlState(ErrTableStorageEngineOption),
		Message: fmt.Sprintf("Table storage engine for '%s' doesn't have this option", table)}
}

// errInvalidPerformanceSchema is 1683, what CREATE TABLE answers for the engine
// the performance schema keeps its own tables in.
func errInvalidPerformanceSchema() error {
	return &Error{Code: ErrInvalidPerformanceSchema, SQLState: sqlState(ErrInvalidPerformanceSchema),
		Message: "Invalid performance_schema usage."}
}

// errCannotPartOfFulltext is 1283, a FULLTEXT key over a column whose type it
// cannot index.
func errCannotPartOfFulltext(name string) error {
	return &Error{Code: ErrCannotPartOfFulltext, SQLState: sqlState(ErrCannotPartOfFulltext),
		Message: fmt.Sprintf("Column '%s' cannot be part of FULLTEXT index", name)}
}

// errKeyColumnDoesNotExist is 1072, a key part naming a column the table does
// not have.
func errKeyColumnDoesNotExist(name string) error {
	return &Error{Code: ErrKeyColumnDoesNotExist, SQLState: sqlState(ErrKeyColumnDoesNotExist),
		Message: fmt.Sprintf("Key column '%s' doesn't exist in table", name)}
}

// errFunctionNotDefined is 1128, what MySQL answers for a function or a
// fulltext parser plugin it does not have.
func errFunctionNotDefined(name string) error {
	return &Error{Code: ErrFunctionNotDefined, SQLState: sqlState(ErrFunctionNotDefined),
		Message: fmt.Sprintf("Function '%s' is not defined", name)}
}

// errTooBigDisplayWidth is 1439, a type's length argument past what MySQL holds
// for it; max is the widest the type takes.
func errTooBigDisplayWidth(name string, max int) error {
	return &Error{Code: ErrTooBigDisplayWidth, SQLState: sqlState(ErrTooBigDisplayWidth),
		Message: fmt.Sprintf("Display width out of range for column '%s' (max = %d)", name, max)}
}

// errPartitionMgmtOnNonPartitioned is 1505, a partition action on a table that
// has no partitions.
func errPartitionMgmtOnNonPartitioned() error {
	return &Error{Code: ErrPartitionMgmtOnNonPartitioned, SQLState: sqlState(ErrPartitionMgmtOnNonPartitioned),
		Message: "Partition management on a not partitioned table is not possible"}
}

func errFKCannotDropParent(table, fkName, refTable string) error {
	return &Error{Code: ErrFKCannotDropParent, SQLState: sqlState(ErrFKCannotDropParent),
		Message: fmt.Sprintf("Cannot drop table '%s' referenced by a foreign key constraint '%s' on table '%s'", table, fkName, refTable)}
}

func errCantDropKey(name string) error {
	return &Error{Code: ErrCantDropKey, SQLState: sqlState(ErrCantDropKey),
		Message: fmt.Sprintf("Can't DROP '%s'; check that column/key exists", name)}
}

// errCheckConstraintNotFound is 3821, the answer a DROP CHECK or ALTER CHECK
// gets when the name is no check constraint of the table, whatever else it
// names.
func errCheckConstraintNotFound(name string) error {
	return &Error{Code: ErrCheckConstraintNotFound, SQLState: sqlState(ErrCheckConstraintNotFound),
		Message: fmt.Sprintf("Check constraint '%s' is not found in the table.", name)}
}

// errConstraintDoesNotExist is 3940, the answer a DROP CONSTRAINT or ALTER
// CONSTRAINT gets when the table holds no constraint of that name.
func errConstraintDoesNotExist(name string) error {
	return &Error{Code: ErrConstraintDoesNotExist, SQLState: sqlState(ErrConstraintDoesNotExist),
		Message: fmt.Sprintf("Constraint '%s' does not exist.", name)}
}

func errCantRemoveAllFields() error {
	return &Error{Code: ErrCantRemoveAllFields, SQLState: sqlState(ErrCantRemoveAllFields),
		Message: "You can't delete all columns with ALTER TABLE; use DROP TABLE instead"}
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

func errFKDupName(name string) error {
	return &Error{Code: ErrFKDupName, SQLState: sqlState(ErrFKDupName),
		Message: fmt.Sprintf("Duplicate foreign key constraint name '%s'", name)}
}

func errCheckConstraintDupName(name string) error {
	return &Error{Code: ErrCheckConstraintDupName, SQLState: sqlState(ErrCheckConstraintDupName),
		Message: fmt.Sprintf("Duplicate check constraint name '%s'.", name)}
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

func errDependentByGeneratedColumn(column, genColumn, table string) error {
	return &Error{Code: ErrDependentByGenCol, SQLState: sqlState(ErrDependentByGenCol),
		Message: fmt.Sprintf("Column '%s' has a generated column dependency and cannot be dropped or renamed. A generated column '%s' refers to this column in table '%s'.", column, genColumn, table)}
}

func errCannotDropColumnFunctionalIndex(column string) error {
	return &Error{Code: ErrDependentByGenCol, SQLState: sqlState(ErrDependentByGenCol),
		Message: fmt.Sprintf("Cannot drop column '%s' because it is used by a functional index.", column)}
}

func errFKCannotUseVirtualColumn(col string) error {
	return &Error{Code: ErrFKCannotUseVirtualColumn, SQLState: sqlState(ErrFKCannotUseVirtualColumn),
		Message: fmt.Sprintf("Cannot define foreign key with clause on a generated column '%s'", col)}
}

func errWrongArguments(fn string) error {
	return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments),
		Message: fmt.Sprintf("Incorrect arguments to %s", fn)}
}

// errWrongParametersToNativeFct is 1583: an argument alias on a native function.
func errWrongParametersToNativeFct(fn string) error {
	return &Error{Code: 1583, SQLState: "42000",
		Message: fmt.Sprintf("Incorrect parameters in the call to native function '%s'", fn)}
}

// errWrongParametersToStoredFct is 1584: an argument alias on a stored function.
func errWrongParametersToStoredFct(fn string) error {
	return &Error{Code: 1584, SQLState: "42000",
		Message: fmt.Sprintf("Incorrect parameters in the call to stored function `%s`", strings.ToLower(fn))}
}

func errIncorrectPrefixKey() error {
	return &Error{Code: ErrIncorrectPrefixKey, SQLState: sqlState(ErrIncorrectPrefixKey),
		Message: "Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys"}
}

func errNotAllowedCommand() error {
	return &Error{Code: ErrNotAllowedCommand, SQLState: sqlState(ErrNotAllowedCommand),
		Message: "The used command is not allowed with this MySQL version"}
}

func errWrongFkDef(forStr, reason string) error {
	return &Error{Code: ErrWrongFkDef, SQLState: sqlState(ErrWrongFkDef),
		Message: fmt.Sprintf("Incorrect foreign key definition for '%s': %s", forStr, reason)}
}

func errDuplicatedValueInType(colName, val, typ string) error {
	return &Error{Code: ErrDuplicatedValueInType, SQLState: sqlState(ErrDuplicatedValueInType),
		Message: fmt.Sprintf("Column '%s' has duplicated value '%s' in %s", colName, val, typ)}
}

func errTablespaceMissing(db, table string) error {
	return &Error{Code: ErrTablespaceMissing, SQLState: sqlState(ErrTablespaceMissing),
		Message: fmt.Sprintf("Tablespace is missing for table `%s`.`%s`.", db, table)}
}

func errTablespaceExists(name string) error {
	return &Error{Code: ErrTablespaceExists, SQLState: sqlState(ErrTablespaceExists),
		Message: fmt.Sprintf("Tablespace '%s' exists.", name)}
}

func errTablespaceDiscarded(db, table string) error {
	return &Error{Code: ErrTablespaceDiscarded, SQLState: sqlState(ErrTablespaceDiscarded),
		Message: fmt.Sprintf("Tablespace has been discarded for table '%s/%s'", db, table)}
}

func errCannotDiscardImportTempTable() error {
	return &Error{Code: ErrCannotDiscardImportTempTable, SQLState: sqlState(ErrCannotDiscardImportTempTable),
		Message: "Cannot DISCARD/IMPORT tablespace associated with temporary table"}
}

func errReservedInnoDBTablespacePrefix() error {
	return &Error{Code: ErrWrongTablespaceName, SQLState: sqlState(ErrWrongTablespaceName),
		Message: "InnoDB: Tablespace names starting with `innodb_` are reserved."}
}

func errTablespaceIsNotEmpty(name string) error {
	return &Error{Code: ErrTablespaceIsNotEmpty, SQLState: sqlState(ErrTablespaceIsNotEmpty),
		Message: fmt.Sprintf("Tablespace '%s' is not empty.", name)}
}

func errJSONUsedAsKey(colName string) error {
	return &Error{Code: ErrJSONUsedAsKey, SQLState: sqlState(ErrJSONUsedAsKey),
		Message: fmt.Sprintf("JSON column '%s' supports indexing only via generated columns on a specified JSON path.", colName)}
}

func errTablespaceDoesNotExist(name string) error {
	return &Error{Code: ErrTablespaceDoesNotExist, SQLState: sqlState(ErrTablespaceDoesNotExist),
		Message: fmt.Sprintf("Tablespace %s doesn't exist.", name)}
}

func errReservedTablespaceName(table, ts string) error {
	return &Error{Code: ErrReservedTablespaceName, SQLState: sqlState(ErrReservedTablespaceName),
		Message: fmt.Sprintf("The table '%s' may not be created in the reserved tablespace '%s'.", table, ts)}
}

func errTempTablespaceNonTempTable() error {
	return &Error{Code: 1478, SQLState: "HY000",
		Message: "InnoDB: Tablespace `innodb_temporary` can only contain TEMPORARY tables."}
}

func errPartitionedTableSharedTablespace() error {
	return &Error{Code: 1478, SQLState: "HY000",
		Message: "InnoDB : A partitioned table is not allowed in a shared tablespace."}
}

func errDependentByDefaultValExpr(colName, tableName string) error {
	return &Error{Code: ErrDependentByDefaultVal, SQLState: sqlState(ErrDependentByDefaultVal),
		Message: fmt.Sprintf("Column '%s' of table '%s' has a default value expression dependency and cannot be dropped or renamed.", colName, tableName)}
}

func errFunctionalIndexPrimaryKey() error {
	return &Error{Code: ErrFunctionalIndexPrimaryKey, SQLState: sqlState(ErrFunctionalIndexPrimaryKey),
		Message: "The primary key cannot be a functional index"}
}

func errFunctionalIndexOnField() error {
	return &Error{Code: ErrFunctionalIndexOnField, SQLState: sqlState(ErrFunctionalIndexOnField),
		Message: "Functional index on a column is not supported. Consider using a regular index instead."}
}

func errDefaultValExprForwardRef(colName string) error {
	return &Error{Code: ErrDefaultValExprNamedFuncNotAllowed, SQLState: sqlState(ErrDefaultValExprNamedFuncNotAllowed),
		Message: fmt.Sprintf("Default value expression of column '%s' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.", colName)}
}

func errDefaultValExprRefAutoInc(colName string) error {
	return &Error{Code: ErrDefaultValExprRefAutoInc, SQLState: sqlState(ErrDefaultValExprRefAutoInc),
		Message: fmt.Sprintf("Default value expression of column '%s' cannot refer to an auto-increment column.", colName)}
}

func errDefaultValExprSubquery(colName string) error {
	return &Error{Code: ErrDefaultValExprContainSubquery, SQLState: sqlState(ErrDefaultValExprContainSubquery),
		Message: fmt.Sprintf("Default value expression of column '%s' cannot contain subqueries.", colName)}
}

func errDefaultValExprVariable(colName string) error {
	return &Error{Code: ErrDefaultValExprRefVariable, SQLState: sqlState(ErrDefaultValExprRefVariable),
		Message: fmt.Sprintf("Default value expression of column '%s' cannot refer to user or system variables.", colName)}
}

func errSecondaryEngineNotDefined() error {
	return &Error{Code: ErrSecondaryEngineNotDefined, SQLState: sqlState(ErrSecondaryEngineNotDefined),
		Message: "Secondary engine operation failed. No secondary engine defined."}
}

func errAutoextendSizeNotMultiple() error {
	return &Error{Code: ErrAutoextendSizeNotMultiple, SQLState: sqlState(ErrAutoextendSizeNotMultiple),
		Message: "AUTOEXTEND_SIZE should be a multiple of 4M"}
}

func errAutoextendSizeOutOfRange() error {
	return &Error{Code: ErrAutoextendSizeOutOfRange, SQLState: sqlState(ErrAutoextendSizeOutOfRange),
		Message: "AUTOEXTEND_SIZE value should be between 4M and 4096M"}
}

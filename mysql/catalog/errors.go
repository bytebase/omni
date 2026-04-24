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
	ErrDupColumn                         = 1060
	ErrDupKeyName                        = 1061
	ErrDupEntry                          = 1062
	ErrMultiplePriKey                    = 1068
	ErrNoSuchTable                       = 1146
	ErrNoSuchColumn                      = 1054
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
	ErrDupEvent                          = 1537
	ErrUnsupportedGeneratedStorageChange = 3106
	ErrDependentByGenCol                 = 3108
	ErrWrongArguments                    = 1210
	ErrWrongNameForIndex                 = 1280
	ErrFKDupName                         = 1826
	ErrCheckConstraintDupName            = 3822
)

var sqlStateMap = map[int]string{
	ErrDupDatabase:                       "HY000",
	ErrUnknownDatabase:                   "42000",
	ErrDupTable:                          "42S01",
	ErrUnknownTable:                      "42S02",
	ErrDupColumn:                         "42S21",
	ErrDupKeyName:                        "42000",
	ErrDupEntry:                          "23000",
	ErrMultiplePriKey:                    "42000",
	ErrNoSuchTable:                       "42S02",
	ErrNoSuchColumn:                      "42S22",
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
	ErrDupEvent:                          "HY000",
	ErrUnsupportedGeneratedStorageChange: "HY000",
	ErrDependentByGenCol:                 "HY000",
	ErrWrongArguments:                    "HY000",
	ErrWrongNameForIndex:                 "42000",
	ErrFKDupName:                         "HY000",
	ErrCheckConstraintDupName:            "HY000",
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

func errCantDropKey(name string) error {
	return &Error{Code: ErrCantDropKey, SQLState: sqlState(ErrCantDropKey),
		Message: fmt.Sprintf("Can't DROP '%s'; check that column/key exists", name)}
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

func errWrongArguments(fn string) error {
	return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments),
		Message: fmt.Sprintf("Incorrect arguments to %s", fn)}
}

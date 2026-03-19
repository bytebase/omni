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
	ErrDupDatabase             = 1007
	ErrUnknownDatabase         = 1049
	ErrDupTable                = 1050
	ErrUnknownTable            = 1051
	ErrDupColumn               = 1060
	ErrDupKeyName              = 1061
	ErrDupEntry                = 1062
	ErrMultiplePriKey          = 1068
	ErrNoSuchTable             = 1146
	ErrNoSuchColumn            = 1054
	ErrNoDatabaseSelected      = 1046
	ErrDupIndex                = 1831
	ErrFKNoRefTable            = 1824
	ErrCantDropKey             = 1091
	ErrCheckConstraintViolated = 3819
	ErrFKCannotDropParent      = 3730
)

var sqlStateMap = map[int]string{
	ErrDupDatabase:             "HY000",
	ErrUnknownDatabase:         "42000",
	ErrDupTable:                "42S01",
	ErrUnknownTable:            "42S02",
	ErrDupColumn:               "42S21",
	ErrDupKeyName:              "42000",
	ErrDupEntry:                "23000",
	ErrMultiplePriKey:          "42000",
	ErrNoSuchTable:             "42S02",
	ErrNoSuchColumn:            "42S22",
	ErrNoDatabaseSelected:      "3D000",
	ErrDupIndex:                "42000",
	ErrFKNoRefTable:            "HY000",
	ErrCantDropKey:             "42000",
	ErrCheckConstraintViolated: "HY000",
	ErrFKCannotDropParent:      "HY000",
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

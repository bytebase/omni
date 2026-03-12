package ddl

import "fmt"

// Warning represents a non-fatal notice emitted during DDL processing,
// analogous to PostgreSQL's NOTICE/WARNING messages.
type Warning struct {
	Code    string // SQLSTATE code
	Message string
}

// SQLSTATE warning codes matching PostgreSQL.
const (
	CodeWarning     = "01000" // generic warning
	CodeWarningSkip = "00000" // IF NOT EXISTS / IF EXISTS skipped
)

// SQLSTATE error codes matching PostgreSQL.
const (
	CodeDuplicateSchema       = "42P06"
	CodeDuplicateTable        = "42P07"
	CodeDuplicateColumn       = "42701"
	CodeDuplicateObject       = "42710"
	CodeUndefinedSchema       = "3F000"
	CodeUndefinedTable        = "42P01"
	CodeUndefinedColumn       = "42703"
	CodeUndefinedObject       = "42704"
	CodeSchemaNotEmpty        = "2BP01"
	CodeDependentObjects      = "2BP01"
	CodeWrongObjectType       = "42809"
	CodeInvalidParameterValue = "22023"
	CodeInvalidFK             = "42830"
	CodeInvalidTableDefinition  = "42P16"
	CodeDuplicatePKey           = "42P16" // same SQLSTATE as InvalidTableDefinition
	CodeDatatypeMismatch        = "42804"
	CodeUndefinedFunction       = "42883"
	CodeAmbiguousColumn         = "42702"
	CodeAmbiguousFunction       = "42725"
	CodeInvalidColumnDefinition = "42611"
	CodeTooManyColumns          = "54011"
	CodeFeatureNotSupported     = "0A000"
	CodeDuplicateFunction       = "42723"
	CodeInvalidObjectDefinition = "42P17"
	CodeSyntaxError             = "42601"
	CodeInvalidFunctionDefinition   = "42P13"
	CodeCheckViolation              = "23514"
	CodeNotNullViolation            = "23502"
	CodeForeignKeyViolation         = "23503"
	CodeUniqueViolation             = "23505"
	CodeIndeterminateCollation      = "42P22"
	CodeObjectNotInPrerequisiteState = "55000"
	CodeInvalidGrantOperation       = "0LP01"
	CodeProgramLimitExceeded        = "54000"
	CodeReservedName                = "42939"
)

// Error represents a PostgreSQL-compatible error with an SQLSTATE code.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ERROR: %s (SQLSTATE %s)", e.Message, e.Code)
}

func errDuplicateSchema(name string) error {
	return &Error{Code: CodeDuplicateSchema, Message: fmt.Sprintf("schema %q already exists", name)}
}

func errDuplicateTable(name string) error {
	return &Error{Code: CodeDuplicateTable, Message: fmt.Sprintf("relation %q already exists", name)}
}

func errDuplicateColumn(name string) error {
	return &Error{Code: CodeDuplicateColumn, Message: fmt.Sprintf("column %q specified more than once", name)}
}

func errUndefinedSchema(name string) error {
	return &Error{Code: CodeUndefinedSchema, Message: fmt.Sprintf("schema %q does not exist", name)}
}

func errUndefinedTable(name string) error {
	return &Error{Code: CodeUndefinedTable, Message: fmt.Sprintf("relation %q does not exist", name)}
}

func errUndefinedType(name string) error {
	return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("type %q does not exist", name)}
}

func errSchemaNotEmpty(name string) error {
	return &Error{Code: CodeSchemaNotEmpty, Message: fmt.Sprintf("cannot drop schema %q because other objects depend on it", name)}
}

func errInvalidParameterValue(msg string) error {
	return &Error{Code: CodeInvalidParameterValue, Message: msg}
}

func errDuplicateObject(kind, name string) error {
	return &Error{Code: CodeDuplicateObject, Message: fmt.Sprintf("%s %q already exists", kind, name)}
}

func errUndefinedColumn(name string) error {
	return &Error{Code: CodeUndefinedColumn, Message: fmt.Sprintf("column %q does not exist", name)}
}

func errInvalidFK(msg string) error {
	return &Error{Code: CodeInvalidFK, Message: msg}
}

func errDuplicatePKey(table string) error {
	return &Error{Code: CodeDuplicatePKey, Message: fmt.Sprintf("multiple primary keys for table %q are not allowed", table)}
}

func errDatatypeMismatch(msg string) error {
	return &Error{Code: CodeDatatypeMismatch, Message: msg}
}

func errDependentObjects(kind, name string) error {
	return &Error{Code: CodeDependentObjects, Message: fmt.Sprintf("cannot drop %s %q because other objects depend on it", kind, name)}
}

func errUndefinedFunction(name string, argTypes []uint32) error {
	return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("function %s(%s) does not exist", name, oidListString(argTypes))}
}

func errAmbiguousColumn(name string) error {
	return &Error{Code: CodeAmbiguousColumn, Message: fmt.Sprintf("column reference %q is ambiguous", name)}
}

func errAmbiguousFunction(name string) error {
	return &Error{Code: CodeAmbiguousFunction, Message: fmt.Sprintf("function %q is not unique", name)}
}

func errDuplicateFunction(name string) error {
	return &Error{Code: CodeDuplicateFunction, Message: fmt.Sprintf("function %q already exists", name)}
}

func errUndefinedSequence(name string) error {
	return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("sequence %q does not exist", name)}
}

func errInvalidObjectDefinition(msg string) error {
	return &Error{Code: CodeInvalidObjectDefinition, Message: msg}
}

func errUndefinedObject(kind, name string) error {
	return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("%s %q does not exist", kind, name)}
}

func errWrongObjectType(name, expected string) error {
	return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not %s", name, expected)}
}

func errInvalidFunctionDefinition(msg string) *Error {
	return &Error{Code: CodeInvalidFunctionDefinition, Message: msg}
}
func errUndefinedTrigger(trigName, tableName string) error {
	return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("trigger %q on relation %q does not exist", trigName, tableName)}
}

func oidListString(oids []uint32) string {
	s := ""
	for i, o := range oids {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d", o)
	}
	return s
}

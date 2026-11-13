package v1review

// The conditions the system rules, the storage engine rule and the drop
// database rule report. An allow-list or a length limit is passed to its
// rule function as a plain value.
const (
	// CharsetNotAllowed: a character set is not in the allowed list. Name
	// holds it as the V1 advisor printed it.
	CharsetNotAllowed Code = "charset-not-allowed"
	// CollationNotAllowed: a collation is not in the allowed list. Name
	// holds it as the V1 advisor printed it.
	CollationNotAllowed Code = "collation-not-allowed"
	// CommentOnTooLong: the comment a COMMENT ON statement sets is longer
	// than the limit.
	CommentOnTooLong Code = "comment-on-too-long"
	// ProcedureCreated, FunctionCreated, EventCreated, ViewCreated: a
	// statement creates an object of a kind the rule forbids.
	ProcedureCreated Code = "procedure-created"
	FunctionCreated  Code = "function-created"
	EventCreated     Code = "event-created"
	ViewCreated      Code = "view-created"
	// FunctionDisallowed: a statement calls a function in the disallowed
	// list. Name holds the function's name as the engine's parser holds it.
	FunctionDisallowed Code = "function-disallowed"
	// EngineNotInnoDB: a statement chooses a storage engine other than
	// InnoDB.
	EngineNotInnoDB Code = "engine-not-innodb"

	// DatabaseNotCurrent: a statement drops a database other than the one
	// the pre-change schema is of. Name holds the dropped database.
	DatabaseNotCurrent Code = "database-not-current"
	// DatabaseNotEmpty: a statement drops the pre-change schema's database
	// while it holds a table. Name holds the dropped database.
	DatabaseNotEmpty Code = "database-not-empty"
)

// The changes SchemaBackwardCompatibility reports as possibly breaking the
// existing data or code. Where an engine's V1 advisor named the object in
// its advice, Name holds that text.
const (
	IncompatibleDropDatabase  Code = "incompatible-drop-database"
	IncompatibleDropSchema    Code = "incompatible-drop-schema"
	IncompatibleDropTable     Code = "incompatible-drop-table"
	IncompatibleRenameTable   Code = "incompatible-rename-table"
	IncompatibleDropColumn    Code = "incompatible-drop-column"
	IncompatibleRenameColumn  Code = "incompatible-rename-column"
	IncompatibleAlterColumn   Code = "incompatible-alter-column"
	IncompatibleAddPrimaryKey Code = "incompatible-add-primary-key"
	IncompatibleAddUniqueKey  Code = "incompatible-add-unique-key"
	IncompatibleAddForeignKey Code = "incompatible-add-foreign-key"
	IncompatibleAddCheck      Code = "incompatible-add-check"
	IncompatibleAlterCheck    Code = "incompatible-alter-check"
	// IncompatibleAddForeignKeyNoCheck, IncompatibleAddCheckNoCheck: a
	// constraint is added without checking the existing rows.
	IncompatibleAddForeignKeyNoCheck Code = "incompatible-add-foreign-key-no-check"
	IncompatibleAddCheckNoCheck      Code = "incompatible-add-check-no-check"
	// IncompatibleDropMaybeNewDatabase, IncompatibleDropMaybeNewSchema,
	// IncompatibleDropMaybeNewTable: a statement drops an object the change
	// itself created earlier, but with IF NOT EXISTS, so the object may
	// predate the change.
	IncompatibleDropMaybeNewDatabase Code = "incompatible-drop-maybe-new-database"
	IncompatibleDropMaybeNewSchema   Code = "incompatible-drop-maybe-new-schema"
	IncompatibleDropMaybeNewTable    Code = "incompatible-drop-maybe-new-table"
)

// DatabaseName returns the name of the database the snapshot is of, "" when
// there is no snapshot.
func (s Schema) DatabaseName() string {
	return s.Metadata.GetName()
}

// HasNoTable reports whether no schema of the snapshot lists a table. Among
// schemas of the same name only the last is read, as the V1 advisors read
// them.
func (s Schema) HasNoTable() bool {
	for _, schema := range lastOfEachName(s.Metadata.GetSchemas(), s.ObjectCaseSensitive, schemaNameOf) {
		if len(schema.GetTables()) > 0 {
			return false
		}
	}
	return true
}

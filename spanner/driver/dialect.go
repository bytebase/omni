package driver

// dialect is one of the two SQL dialects Spanner serves. A database's dialect
// is fixed when the database is created and decides its SQL and its
// INFORMATION_SCHEMA, so it decides the catalog a sync reads and the grammar a
// script is split by.
type dialect int

const (
	// googleSQL is Spanner's own dialect, the admin API's default, and what a
	// database that states none is read as.
	googleSQL dialect = iota
	// postgresSQL is a database's PostgreSQL dialect.
	postgresSQL
)

// dialectOf reads the dialect a database stated, as Database.Dialect spells it.
func dialectOf(name string) dialect {
	if name == DialectPostgresSQL {
		return postgresSQL
	}
	return googleSQL
}

// dialectName is the dialect as the admin API spells it, and empty for the API's
// own, GoogleSQL, which is what a creation asks for when it names nothing.
func dialectName(d dialect) string {
	if d == postgresSQL {
		return DialectPostgresSQL
	}
	return ""
}

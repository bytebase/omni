package analysis

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/metadata"
)

// mysqlMigrationDatabase is the database both snapshots of a MySQL-family
// migration load under, so that snapshots of differently named databases
// compare object by object. The generated DDL then qualifies objects with it,
// and stripMigrationDatabase removes that qualifier: the migration runs in a
// connection already scoped to the target database.
const mysqlMigrationDatabase = "bbcatalog"

// MySQLMetadataMigration is the metadata migration MySQL, MariaDB, OceanBase
// and TiDB share, on the engine's catalog C. Its MetadataMigration method
// implements MetadataMigrator.
type MySQLMetadataMigration[C any] struct {
	// Load installs a snapshot into a new catalog under the database named
	// database. It returns the objects the catalog could not install as the
	// snapshot defines them, keyed as the catalog's LoadMetadataReport is.
	Load func(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, database string) (C, []map[string]error, error)
	// Migration renders the DDL that takes from to to, "" when they agree.
	Migration func(from, to C) string
}

// MetadataMigration loads both snapshots and renders the migration between
// them. It refuses relation names that differ only by case when either side
// is case-sensitive, since the catalogs key relations by lower-cased name.
// An object a load could not install as defined is safe to leave only when it
// is the same on both sides, so a migration that changes one is refused.
func (m MySQLMetadataMigration[C]) MetadataMigration(ctx context.Context, request MetadataMigrationRequest) (string, error) {
	fromMeta, toMeta := migrationSnapshot(request.From), migrationSnapshot(request.To)
	fromCaseSensitive := request.From != nil && request.From.ObjectCaseSensitive
	toCaseSensitive := request.To != nil && request.To.ObjectCaseSensitive
	if err := rejectCaseOnlyIdentities(fromMeta, toMeta, fromCaseSensitive, toCaseSensitive); err != nil {
		return "", err
	}
	from, fromProblems, err := m.Load(ctx, fromMeta, mysqlMigrationDatabase)
	if err != nil {
		return "", fmt.Errorf("failed to load source schema: %w", err)
	}
	to, toProblems, err := m.Load(ctx, toMeta, mysqlMigrationDatabase)
	if err != nil {
		return "", fmt.Errorf("failed to load target schema: %w", err)
	}
	if err := unloadableChanges(fromMeta, toMeta, slices.Concat(fromProblems, toProblems)); err != nil {
		return "", err
	}
	return stripMigrationDatabase(m.Migration(from, to)), nil
}

func migrationSnapshot(database *Database) *metadata.DatabaseSchemaMetadata {
	if database == nil || database.Metadata == nil {
		return &metadata.DatabaseSchemaMetadata{}
	}
	return database.Metadata
}

// rejectCaseOnlyIdentities refuses a snapshot whose table or view names
// differ only by case, and a case-only rename between the snapshots.
func rejectCaseOnlyIdentities(fromMeta, toMeta *metadata.DatabaseSchemaMetadata, fromCaseSensitive, toCaseSensitive bool) error {
	if !fromCaseSensitive && !toCaseSensitive {
		return nil
	}
	for _, meta := range []*metadata.DatabaseSchemaMetadata{fromMeta, toMeta} {
		for _, schema := range meta.GetSchemas() {
			seen := make(map[string]string)
			for _, name := range relationNames(schema) {
				lower := strings.ToLower(name)
				if previous, ok := seen[lower]; ok && previous != name {
					return fmt.Errorf("case-sensitive schema has %q and %q, which the catalog cannot distinguish", previous, name)
				}
				seen[lower] = name
			}
		}
	}
	fromNames := make(map[string]string)
	for _, schema := range fromMeta.GetSchemas() {
		for _, name := range relationNames(schema) {
			fromNames[strings.ToLower(name)] = name
		}
	}
	for _, schema := range toMeta.GetSchemas() {
		for _, name := range relationNames(schema) {
			if previous, ok := fromNames[strings.ToLower(name)]; ok && previous != name {
				return fmt.Errorf("case-only rename %q to %q is not supported by the catalog diff", previous, name)
			}
		}
	}
	return nil
}

func relationNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.GetTables())+len(schema.GetViews()))
	for _, table := range schema.GetTables() {
		names = append(names, table.GetName())
	}
	for _, view := range schema.GetViews() {
		names = append(names, view.GetName())
	}
	return names
}

// unloadableChanges returns an error naming each object a load could not
// install as defined that differs between the snapshots. An object in several
// problem maps takes the error of the first.
func unloadableChanges(fromMeta, toMeta *metadata.DatabaseSchemaMetadata, problems []map[string]error) error {
	reasons := make(map[string]error)
	for _, objects := range problems {
		for key, err := range objects {
			if _, seen := reasons[key]; !seen {
				reasons[key] = err
			}
		}
	}
	var changed []string
	for key, err := range reasons {
		if !proto.Equal(snapshotObject(fromMeta, key), snapshotObject(toMeta, key)) {
			changed = append(changed, key+": "+err.Error())
		}
	}
	if len(changed) == 0 {
		return nil
	}
	slices.Sort(changed)
	return fmt.Errorf("cannot generate a migration for objects whose definitions could not be loaded: %s", strings.Join(changed, "; "))
}

// snapshotObject finds the object a LoadMetadataReport key names, or nil.
func snapshotObject(meta *metadata.DatabaseSchemaMetadata, key string) proto.Message {
	kind, name, _ := strings.Cut(key, ":")
	for _, s := range meta.GetSchemas() {
		switch kind {
		case "table":
			for _, t := range s.GetTables() {
				if t.GetName() == name {
					// Statistics change with the data, not the definition.
					t = proto.CloneOf(t)
					t.RowCount, t.DataSize, t.IndexSize, t.DataFree = 0, 0, 0, 0
					return t
				}
			}
		case "view":
			for _, v := range s.GetViews() {
				if v.GetName() == name {
					return v
				}
			}
		case "sequence":
			for _, sequence := range s.GetSequences() {
				if sequence.GetName() == name {
					clone := proto.CloneOf(sequence)
					clone.LastValue = ""
					return clone
				}
			}
		case "function":
			for _, f := range s.GetFunctions() {
				if f.GetName() == name {
					return f
				}
			}
		case "procedure":
			for _, p := range s.GetProcedures() {
				if p.GetName() == name {
					return p
				}
			}
		case "trigger":
			table, trigger, _ := strings.Cut(name, ".")
			for _, t := range s.GetTables() {
				if t.GetName() != table {
					continue
				}
				for _, trg := range t.GetTriggers() {
					if trg.GetName() == trigger {
						return trg
					}
				}
			}
		case "event":
			for _, e := range s.GetEvents() {
				if e.GetName() == name {
					return e
				}
			}
		default:
		}
	}
	return nil
}

// stripMigrationDatabase removes the mysqlMigrationDatabase qualifier from
// generated DDL. The scan copies string literals and the identifier after a
// removed qualifier whole, so only qualifiers in identifier position go, even
// for an object named like the database. It copies #, "-- " and /*...*/
// comments whole too, since routine bodies are emitted verbatim and a quote
// in a comment would otherwise open a phantom literal; the inside of an
// executable /*!...*/ comment is live SQL and is scanned. A table alias named
// like the database loses its qualifier too.
func stripMigrationDatabase(sql string) string {
	quoted := "`" + mysqlMigrationDatabase + "`"
	if !strings.Contains(sql, quoted+".") {
		return sql
	}
	var out strings.Builder
	out.Grow(len(sql))
	for i := 0; i < len(sql); {
		c := sql[i]
		if isMySQLSpaceByte(c) {
			out.WriteByte(c)
			i++
			continue
		}
		if end, isComment := scanMySQLComment(sql, i); isComment {
			out.WriteString(sql[i:end])
			i = end
			continue
		}
		end := scanMySQLToken(sql, i)
		if sql[i:end] == quoted && end+1 < len(sql) && sql[end] == '.' && sql[end+1] == '`' {
			identEnd := scanMySQLToken(sql, end+1)
			out.WriteString(sql[end+1 : identEnd])
			i = identEnd
			continue
		}
		out.WriteString(sql[i:end])
		i = end
	}
	return out.String()
}

// scanMySQLComment reports whether a comment that is not executable starts at
// i and, if so, the offset just past it. "#" and "-- " (a double dash followed
// by whitespace or the end) run to the end of the line, "/*" to its "*/".
func scanMySQLComment(sql string, i int) (int, bool) {
	switch sql[i] {
	case '#':
		return mysqlLineEnd(sql, i), true
	case '-':
		if i+1 < len(sql) && sql[i+1] == '-' && (i+2 >= len(sql) || isMySQLSpaceByte(sql[i+2])) {
			return mysqlLineEnd(sql, i), true
		}
	case '/':
		if i+1 < len(sql) && sql[i+1] == '*' && (i+2 >= len(sql) || sql[i+2] != '!') {
			if idx := strings.Index(sql[i+2:], "*/"); idx >= 0 {
				return i + 2 + idx + 2, true
			}
			return len(sql), true
		}
	default:
	}
	return 0, false
}

// mysqlLineEnd returns the offset of the newline ending the line at i, or the
// end of sql.
func mysqlLineEnd(sql string, i int) int {
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	return i
}

func isMySQLSpaceByte(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// scanMySQLToken returns the end offset of the token starting at the
// non-space byte i: a whole '...' or "..." literal, a whole `...` identifier,
// a word, or one punctuation byte. Literals honor backslash escapes and
// doubled delimiters, identifiers doubled backticks; an unterminated one runs
// to the end.
func scanMySQLToken(sql string, i int) int {
	switch c := sql[i]; c {
	case '\'', '"':
		for j := i + 1; j < len(sql); {
			switch {
			case sql[j] == '\\' && j+1 < len(sql):
				j += 2
			case sql[j] == c:
				if j+1 < len(sql) && sql[j+1] == c {
					j += 2
					continue
				}
				return j + 1
			default:
				j++
			}
		}
		return len(sql)
	case '`':
		for j := i + 1; j < len(sql); {
			if sql[j] == '`' {
				if j+1 < len(sql) && sql[j+1] == '`' {
					j += 2
					continue
				}
				return j + 1
			}
			j++
		}
		return len(sql)
	default:
		if !isMySQLWordByte(sql[i]) {
			return i + 1
		}
		j := i + 1
		for j < len(sql) && isMySQLWordByte(sql[j]) {
			j++
		}
		return j
	}
}

// isMySQLWordByte reports whether b can be part of a bare word; a byte of a
// multibyte character can.
func isMySQLWordByte(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b == '_' || b == '$' || b >= 0x80
}

// MySQLMetadataExtraction is the metadata extraction MySQL, MariaDB,
// OceanBase and TiDB share. Its ExtractMetadata method implements
// MetadataExtractor.
type MySQLMetadataExtraction struct {
	// Database is the database the script runs in.
	Database   string
	NewCatalog func() MySQLScriptCatalog
	// Export returns the metadata the script defined; a nil Export returns
	// Database's.
	Export func(c MySQLScriptCatalog) *metadata.DatabaseSchemaMetadata
}

// MySQLScriptCatalog is a MySQL-family catalog a schema script runs on.
type MySQLScriptCatalog interface {
	// Run executes a script, continuing past a statement that fails.
	Run(script string) ([]CatalogResult, error)
	CurrentDatabase() string
	DatabaseMetadata(database string) *metadata.DatabaseSchemaMetadata
}

// ExtractMetadata runs the script in Database on a new catalog and exports
// the result. The first statement that fails fails the extraction.
func (x MySQLMetadataExtraction) ExtractMetadata(_ context.Context, script string) (*metadata.DatabaseSchemaMetadata, error) {
	if strings.TrimSpace(script) == "" {
		return &metadata.DatabaseSchemaMetadata{}, nil
	}
	c := x.NewCatalog()
	setup := fmt.Sprintf("SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `%s`;\nUSE `%s`;", x.Database, x.Database)
	if _, err := c.Run(setup); err != nil {
		return nil, fmt.Errorf("failed to initialize catalog: %w", err)
	}
	results, err := c.Run(script)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema DDL: %w", err)
	}
	for _, r := range results {
		if r.Error != nil {
			return nil, fmt.Errorf("failed to execute schema DDL: %w", r.Error)
		}
	}
	if x.Export != nil {
		return x.Export(c), nil
	}
	return c.DatabaseMetadata(x.Database), nil
}

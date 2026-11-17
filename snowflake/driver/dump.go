package driver

import (
	"context"
	"fmt"
	"io"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/snowflake/parser"
)

// Dump writes the server's own definition of the database the driver is open
// on: what GET_DDL('DATABASE', ...) states for it, with four passages taken out
// so that the script creates the database's objects rather than replacing the
// database itself.
//
// It is a server-DDL dump, so it ignores the metadata it is given: an edited or
// narrowed snapshot changes nothing in what is written. Bytebase's Snowflake
// dump was the same call and the same four transforms; VERIFY.md records what
// GET_DDL states for a pinned database before they are trusted.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	query := fmt.Sprintf("SELECT GET_DDL('DATABASE', %s, TRUE)", sqlLiteral(quoteIdent(d.databaseName)))
	var definition string
	if err := d.db.QueryRowContext(ctx, query).Scan(&definition); err != nil {
		return fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if _, err := io.WriteString(out, databaseDefinition(d.databaseName, definition)); err != nil {
		return err
	}
	return nil
}

// databaseDefinition is the script a dump writes for a database, from what
// GET_DDL states for it.
func databaseDefinition(databaseName, definition string) string {
	definition = dropCreateDatabase(databaseName, definition)
	definition = dropPublicSchema(databaseName, definition)
	return keepCreate(databaseName, definition)
}

// dropCreateDatabase removes the CREATE DATABASE statement GET_DDL opens with,
// and the blank line after it: a dump writes the objects of a database, and the
// database itself is the caller's to create.
//
// The statement is only removed when the definition opens with it; a definition
// that states something else first is written as the server stated it.
func dropCreateDatabase(databaseName, definition string) string {
	lines := strings.Split(definition, "\n")
	if !isCreateDatabaseStatement(lines[0]) {
		return definition
	}
	// The statement, and the blank line after it when the server wrote one.
	if len(lines) <= 2 {
		return ""
	}
	return strings.Join(lines[2:], "\n")
}

// isCreateDatabaseStatement reports whether line opens with CREATE [OR REPLACE]
// [TRANSIENT] DATABASE, which is what GET_DDL states for the database it
// describes. TRANSIENT is the modifier a transient database is created with.
func isCreateDatabaseStatement(line string) bool {
	words := strings.Fields(strings.ToLower(line))
	if len(words) == 0 || words[0] != "create" {
		return false
	}
	words = words[1:]
	if len(words) >= 2 && words[0] == "or" && words[1] == "replace" {
		words = words[2:]
	}
	if len(words) >= 1 && words[0] == "transient" {
		words = words[1:]
	}
	return len(words) >= 1 && words[0] == "database"
}

// dropPublicSchema removes the statement that creates the default PUBLIC schema,
// which every database holds whether or not it was created with one.
//
// The statement is recognised by what it creates rather than by its text: the
// header may say OR REPLACE, and the database may be named as the quoted
// identifier its DDL requires.
func dropPublicSchema(databaseName, definition string) string {
	lines := strings.Split(definition, "\n")
	kept := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(strings.ToLower(line), "create ") && strings.EqualFold(keepCreateLine(databaseName, line), "create schema PUBLIC;") {
			// The statement, and the blank line after it when the server wrote
			// one.
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
				i++
			}
			continue
		}
		kept = append(kept, lines[i])
	}
	return strings.Join(kept, "\n")
}

// keepCreate takes the database qualifier off a CREATE statement and turns
// CREATE OR REPLACE into a CREATE: a dump restores into objects that may
// already be there, and replacing one would drop what it holds.
//
// Only the statement's head changes. A literal or a comment that spells
// `create or replace` or the database's name is part of the object's definition,
// so changing it would change what the object answers.
func keepCreate(databaseName, definition string) string {
	lines := strings.Split(definition, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(strings.ToLower(line), "create ") {
			continue
		}
		lines[i] = keepCreateLine(databaseName, line)
	}
	return strings.Join(lines, "\n")
}

// createOrReplace is the clause a dump's CREATE statement loses.
const createOrReplace = "create or replace "

// keepCreateLine rewrites one CREATE statement's head: the leading CREATE OR
// REPLACE becomes CREATE, and the object name loses the database it qualifies.
func keepCreateLine(databaseName, line string) string {
	if len(line) >= len(createOrReplace) && strings.EqualFold(line[:len(createOrReplace)], createOrReplace) {
		line = "create " + line[len(createOrReplace):]
	}
	tokens, _ := parser.Tokenize(line)
	for i := 0; i+1 < len(tokens); i++ {
		name, dot := tokens[i], tokens[i+1]
		if dot.Type != '.' || !isNameText(line, name) {
			continue
		}
		if !isDatabaseName(line[name.Loc.Start:name.Loc.End], databaseName) {
			continue
		}
		return line[:name.Loc.Start] + line[dot.Loc.End:]
	}
	return line
}

// isNameText reports whether the token at that source offset is a name rather
// than a literal: a string keeps its bytes and is not rewritten.
func isNameText(line string, token parser.Token) bool {
	if token.Loc.Start < 0 || token.Loc.Start >= len(line) {
		return false
	}
	switch c := line[token.Loc.Start]; {
	case c == '"', c == '_':
		return true
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z':
		return true
	}
	return false
}

// isDatabaseName reports whether the name the statement states is the
// database's: a quoted name is compared as written, an unquoted one by case,
// since Snowflake folds it.
func isDatabaseName(name, databaseName string) bool {
	if len(name) > 1 && name[0] == '"' && name[len(name)-1] == '"' {
		return strings.ReplaceAll(name[1:len(name)-1], `""`, `"`) == databaseName
	}
	return strings.EqualFold(name, databaseName)
}

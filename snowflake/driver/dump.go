package driver

import (
	"context"
	"fmt"
	"io"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
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
	first := strings.ToLower(strings.TrimSpace(lines[0]))
	if !strings.HasPrefix(first, "create or replace database ") && !strings.HasPrefix(first, "create database ") {
		return definition
	}
	// The statement, and the blank line after it when the server wrote one.
	if len(lines) <= 2 {
		return ""
	}
	return strings.Join(lines[2:], "\n")
}

// dropPublicSchema removes the statement that creates the default PUBLIC schema,
// which every database holds.
func dropPublicSchema(databaseName, definition string) string {
	schemaStatement := fmt.Sprintf("create or replace schema %s.PUBLIC;", databaseName)
	definition = strings.ReplaceAll(definition, schemaStatement+"\n\n", "")
	return strings.ReplaceAll(definition, schemaStatement, "")
}

// keepCreate takes the database qualifier off a CREATE statement and turns every
// CREATE OR REPLACE into a CREATE: a dump restores into objects that may already
// be there, and replacing one would drop what it holds.
func keepCreate(databaseName, definition string) string {
	lines := strings.Split(definition, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(strings.ToLower(line), "create ") {
			continue
		}
		line = strings.ReplaceAll(line, " "+databaseName+".", " ")
		lines[i] = strings.ReplaceAll(line, "create or replace ", "create ")
	}
	return strings.Join(lines, "\n")
}

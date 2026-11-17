package driver

import (
	"context"
	"database/sql"
	"io"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

// Dump writes the database's DDL as CockroachDB states it. The server's own
// SHOW CREATE answers are the only DDL CockroachDB publishes for its objects,
// and the caller's metadata narrows nothing: the parameter is the contract's,
// not a source of text.
func (d *Driver) Dump(ctx context.Context, w io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	sb := &strings.Builder{}
	if err := crdb.ExecuteTx(ctx, d.db, &sql.TxOptions{
		ReadOnly: true,
	}, func(txn *sql.Tx) error {
		createSchemas, err := dumpCreateSchemas(ctx, txn)
		if err != nil {
			return err
		}
		_, _ = sb.WriteString(createSchemas)

		createTables, err := dumpCreateTables(ctx, txn)
		if err != nil {
			return err
		}
		_, _ = sb.WriteString(createTables)

		return nil
	}); err != nil {
		return err
	}

	_, err := w.Write([]byte(sb.String()))
	return err
}

func dumpCreateTables(ctx context.Context, txn *sql.Tx) (string, error) {
	sb := &strings.Builder{}
	rows, err := txn.QueryContext(ctx, "SELECT create_statement FROM [SHOW CREATE ALL TABLES];")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return "", err
		}
		_, _ = sb.WriteString(payload)
		_, _ = sb.WriteString("\n\n")
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func dumpCreateSchemas(ctx context.Context, txn *sql.Tx) (string, error) {
	sb := &strings.Builder{}
	rows, err := txn.QueryContext(ctx, "SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS];")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return "", err
		}
		// Skip public schema
		if strings.HasSuffix(payload, " public;") {
			continue
		}
		_, _ = sb.WriteString(payload)
		_, _ = sb.WriteString("\n\n")
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

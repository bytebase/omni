package driver

import (
	"context"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// TestColumnMetadata covers the mapping from one INFORMATION_SCHEMA.COLUMNS row
// onto a synced column, including the two absences a row states: no collation, and
// no default, which reads as the text "NULL".
func TestColumnMetadata(t *testing.T) {
	column, err := columnMetadata(columnRow{
		TableName:       "customers",
		ColumnName:      "city",
		OrdinalPosition: 2,
		IsNullable:      "YES",
		DataType:        "STRING",
	})
	require.NoError(t, err)
	require.Equal(t, &metadatapb.ColumnMetadata{
		Name:     "city",
		Position: 2,
		Nullable: true,
		Type:     "STRING",
		Default:  "NULL",
	}, column)

	column, err = columnMetadata(columnRow{
		TableName:       "customers",
		ColumnName:      "id",
		OrdinalPosition: 1,
		IsNullable:      "NO",
		DataType:        "INT64",
		CollationName:   bigquery.NullString{StringVal: "und:ci", Valid: true},
		ColumnDefault:   bigquery.NullString{StringVal: "0", Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, &metadatapb.ColumnMetadata{
		Name:      "id",
		Position:  1,
		Nullable:  false,
		Type:      "INT64",
		Collation: "und:ci",
		Default:   "0",
	}, column)

	// A row whose nullability BigQuery does not answer is an error rather than a
	// guess.
	_, err = columnMetadata(columnRow{ColumnName: "flag", IsNullable: "MAYBE", DataType: "BOOL"})
	require.Error(t, err)
	require.Contains(t, err.Error(), `unrecognized is_nullable value "MAYBE"`)
}

// TestConvertYesNo pins the two answers INFORMATION_SCHEMA states.
func TestConvertYesNo(t *testing.T) {
	yes, err := convertYesNo("YES")
	require.NoError(t, err)
	require.True(t, yes)

	no, err := convertYesNo("NO")
	require.NoError(t, err)
	require.False(t, no)

	_, err = convertYesNo("yes")
	require.Error(t, err)
}

// TestSyncDatabaseWithoutADataset pins that a driver open on no dataset reports
// the contract's not-found instead of asking the service about an empty name.
func TestSyncDatabaseWithoutADataset(t *testing.T) {
	_, err := (&Driver{}).SyncDatabase(context.Background())
	require.ErrorIs(t, err, driver.ErrNotFound)
}

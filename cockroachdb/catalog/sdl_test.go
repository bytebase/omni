package catalog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/catalog"
)

func TestLoadSDLRefuses(t *testing.T) {
	for name, sdl := range map[string]string{
		"syntax":           "CREATE TABLE (",
		"data":             "INSERT INTO t VALUES (1)",
		"missing table":    "CREATE INDEX i ON app.t (a)",
		"missing column":   "CREATE TABLE t (a INT8); COMMENT ON COLUMN t.b IS 'x'",
		"table as select":  "CREATE TABLE t AS SELECT 1",
		"unknown sequence": "ALTER SEQUENCE s INCREMENT BY 2",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := catalog.LoadSDL(sdl)
			require.Error(t, err)
		})
	}
}

func TestBuildSDLMigrationPlan(t *testing.T) {
	source := `SET statement_timeout = 0;
SELECT pg_catalog.set_config('search_path', '', false);
CREATE TABLE public.users (id INT8 PRIMARY KEY, team INT8 REFERENCES public.teams);
CREATE TABLE public.teams (id INT8 PRIMARY KEY);`
	plan, err := catalog.BuildSDLMigrationPlan(source, source, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan)

	_, err = catalog.BuildSDLMigrationPlan("CREATE TABLE (", source, catalog.SDLMigrationOptions{})
	require.ErrorContains(t, err, "failed to load source schema")
	_, err = catalog.BuildSDLMigrationPlan(source, "CREATE TABLE (", catalog.SDLMigrationOptions{})
	require.ErrorContains(t, err, "failed to load target schema")

	target := `CREATE TABLE public.users (id INT8 PRIMARY KEY, email STRING);
CREATE TABLE public.teams (id INT8 PRIMARY KEY);`
	plan, err = catalog.BuildSDLMigrationPlan(source, target, catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Equal(t, "ALTER TABLE \"public\".\"users\" DROP CONSTRAINT \"users_team_fkey\";\n"+
		"ALTER TABLE \"public\".\"users\" DROP COLUMN \"team\";\n"+
		"ALTER TABLE \"public\".\"users\" ADD COLUMN \"email\" STRING", plan.SQL())
	var types []catalog.MigrationOpType
	for _, op := range plan.Ops {
		types = append(types, op.Type)
	}
	require.Equal(t, []catalog.MigrationOpType{catalog.OpDropConstraint, catalog.OpDropColumn, catalog.OpAddColumn}, types)
}

// TestLoadSDLReadsDump reads the DDL cockroachdb/driver dumps (SHOW CREATE),
// which Bytebase hands the SDL diff as a raw schema, and requires the SDL
// rendered from it to hold the same schema.
func TestLoadSDLReadsDump(t *testing.T) {
	dump, err := os.ReadFile(filepath.Join("..", "driver", "testdata", "dump", "default.sql"))
	require.NoError(t, err)
	loaded, err := catalog.LoadSDL(string(dump))
	require.NoError(t, err)
	sdl, err := catalog.DatabaseDefinition(loaded, catalog.DefinitionOptions{SDLFormat: true})
	require.NoError(t, err)
	require.NotContains(t, sdl, "omni_testdump")
	plan, err := catalog.BuildSDLMigrationPlan(sdl, string(dump), catalog.SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "SDL:\n%s\ndiffers from the dump by:\n%s", sdl, plan.SQL())
}

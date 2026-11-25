package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// TestSequenceSimulation holds the catalog to testdata/sequences/simulate.tsv:
// every sequence statement OceanBase accepts in the parser's compat.tsv, and a
// curated matrix of the options' values and checks, run on each fixture, with
// what OceanBase did and what Exec does.
//
// OceanBase has no SHOW CREATE SEQUENCE (1064) and lists no sequence in SHOW
// FULL TABLES, so a sequence is read from oceanbase.DBA_SEQUENCES and written
// as one CREATE SEQUENCE text with every option (sequenceForm); tables and
// views are compared as TestSimulation compares them.
//
// OceanBase's side is on the test's tape
// (testdata/tapes/TestSequenceSimulation.tape), so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned OceanBase again; -update-compat
// commits the rows a catalog fix brought to OceanBase's outcome; and
// -simulate-show=<text> prints both schemas of the rows that hold the text.
func TestSequenceSimulation(t *testing.T) {
	rows, err := sequenceSimulationRows(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "sequences", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   testcontainer.Codec{}.MadeUp(),
		Omni:     sequenceSimulationOmni,
		Engine:   sequenceSimulationEngine,
	})
}

// sequenceSimulationFixtures: no object named a (one_table@missing), a table
// and a view named a, whose names are apart from a sequence's, and a sequence
// named a with the default options (objects@sequence) and with every option
// set, descending and cyclic (objects@sequence_options).
var sequenceSimulationFixtures = []string{
	"one_table@missing", "one_table@table", "objects@view",
	"objects@sequence", "objects@sequence_options",
}

func sequenceSimulated(statement string) bool {
	for _, words := range []string{"CREATE SEQUENCE ", "ALTER SEQUENCE ", "DROP SEQUENCE "} {
		if strings.HasPrefix(statement, words) {
			return true
		}
	}
	return false
}

// curatedSequenceStatements cover what compat.tsv's rows, one per grammar
// alternative, leave out: each check OceanBase's option builder makes, the
// order it makes them in, values past 64 bits and with a fraction, every
// ALTER SEQUENCE option on a sequence that exists, and names in another
// database, in another case and too long.
//
// Ref: ob_sequence_resolver.h, ob_sequence_option_builder.cpp and
// ob_sequence_ddl_proxy.cpp (oceanbase v4.3.5_CE_BP6).
var curatedSequenceStatements = []programSpec{
	{"sequence/create_start_increment", "CREATE SEQUENCE a START WITH 5 INCREMENT BY 2"},
	{"sequence/create_descending", "CREATE SEQUENCE a INCREMENT BY -1"},
	{"sequence/create_descending_minvalue", "CREATE SEQUENCE a MINVALUE -10 INCREMENT BY -2"},
	{"sequence/create_descending_start_above_max", "CREATE SEQUENCE a INCREMENT BY -3 START WITH 4"},
	{"sequence/create_every_option", "CREATE SEQUENCE a START WITH 5 INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE ORDER"},
	{"sequence/create_cycle_cache_fits", "CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 10"},
	{"sequence/create_cycle_cache_two", "CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 3"},
	{"sequence/create_cycle_cache_too_large", "CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 200"},
	{"sequence/create_cycle_cache_one_cycle", "CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 2"},
	{"sequence/create_cycle_default_cache", "CREATE SEQUENCE a MAXVALUE 5 CYCLE"},
	{"sequence/create_descending_cycle", "CREATE SEQUENCE a INCREMENT BY -1 CYCLE"},
	{"sequence/create_descending_cycle_minvalue", "CREATE SEQUENCE a INCREMENT BY -1 MINVALUE -100 CYCLE"},
	{"sequence/create_cache_zero", "CREATE SEQUENCE a CACHE 0"},
	{"sequence/create_cache_two", "CREATE SEQUENCE a CACHE 2"},
	{"sequence/create_increment_zero", "CREATE SEQUENCE a INCREMENT BY 0"},
	{"sequence/create_increment_zero_cache_one", "CREATE SEQUENCE a INCREMENT BY 0 CACHE 1"},
	{"sequence/create_cache_one_maxvalue_one", "CREATE SEQUENCE a CACHE 1 MAXVALUE 1"},
	{"sequence/create_start_below_min", "CREATE SEQUENCE a START WITH 0"},
	{"sequence/create_start_negative", "CREATE SEQUENCE a START WITH -5"},
	{"sequence/create_start_above_max", "CREATE SEQUENCE a START WITH 100 MAXVALUE 10"},
	{"sequence/create_min_above_max", "CREATE SEQUENCE a MINVALUE 10 MAXVALUE 5"},
	{"sequence/create_min_at_int64", "CREATE SEQUENCE a MINVALUE 9223372036854775807"},
	{"sequence/create_increment_too_large", "CREATE SEQUENCE a MINVALUE 1 MAXVALUE 3 INCREMENT BY 5"},
	{"sequence/create_maxvalue_past_int64", "CREATE SEQUENCE a MAXVALUE 9223372036854775808"},
	{"sequence/create_maxvalue_past_uint64", "CREATE SEQUENCE a MAXVALUE 99999999999999999999999999999"},
	{"sequence/create_minvalue_past_int64", "CREATE SEQUENCE a MINVALUE -99999999999999999999999999999 INCREMENT BY -1"},
	{"sequence/create_maxvalue_whole_decimal", "CREATE SEQUENCE a MAXVALUE 10.0"},
	{"sequence/create_maxvalue_wide_fraction", "CREATE SEQUENCE a MAXVALUE 99999999999999999999.5"},
	{"sequence/create_maxvalue_fraction_cycle", "CREATE SEQUENCE a MAXVALUE 1.5 CYCLE"},
	{"sequence/create_minvalue_fraction", "CREATE SEQUENCE a MINVALUE 1.5"},
	{"sequence/create_increment_fraction", "CREATE SEQUENCE a INCREMENT BY 1.5"},
	{"sequence/create_start_fraction", "CREATE SEQUENCE a START WITH 1.5"},
	{"sequence/create_cache_fraction", "CREATE SEQUENCE a CACHE 1.5"},
	{"sequence/create_maxvalue_exponent", "CREATE SEQUENCE a MAXVALUE 1e5"},
	{"sequence/create_nomaxvalue_descending", "CREATE SEQUENCE a NOMAXVALUE INCREMENT BY -1"},
	{"sequence/create_dup_increment", "CREATE SEQUENCE a INCREMENT BY 1 INCREMENT BY 2"},
	{"sequence/create_dup_start", "CREATE SEQUENCE a START WITH 1 START WITH 2"},
	{"sequence/create_dup_maxvalue", "CREATE SEQUENCE a MAXVALUE 5 MAXVALUE 6"},
	{"sequence/create_dup_minvalue", "CREATE SEQUENCE a MINVALUE 1 MINVALUE 2"},
	{"sequence/create_dup_nominvalue", "CREATE SEQUENCE a NOMINVALUE NOMINVALUE"},
	{"sequence/create_dup_cache", "CREATE SEQUENCE a CACHE 2 CACHE 3"},
	{"sequence/create_dup_nocache", "CREATE SEQUENCE a NOCACHE NOCACHE"},
	{"sequence/create_dup_cycle", "CREATE SEQUENCE a CYCLE CYCLE"},
	{"sequence/create_dup_order", "CREATE SEQUENCE a ORDER ORDER"},
	{"sequence/create_dup_restart", "CREATE SEQUENCE a RESTART RESTART"},
	{"sequence/create_conflict_maxvalue", "CREATE SEQUENCE a MAXVALUE 5 NOMAXVALUE"},
	{"sequence/create_conflict_minvalue", "CREATE SEQUENCE a MINVALUE 5 NOMINVALUE"},
	{"sequence/create_conflict_cache", "CREATE SEQUENCE a CACHE 5 NOCACHE"},
	{"sequence/create_conflict_order", "CREATE SEQUENCE a ORDER NOORDER"},
	{"sequence/create_conflict_cycle", "CREATE SEQUENCE a NOCYCLE CYCLE"},
	{"sequence/create_conflicts_in_order", "CREATE SEQUENCE a CYCLE NOCYCLE ORDER NOORDER"},
	{"sequence/create_if_not_exists_cycle", "CREATE SEQUENCE IF NOT EXISTS a CYCLE"},
	{"sequence/create_if_not_exists_conflict", "CREATE SEQUENCE IF NOT EXISTS a CYCLE NOCYCLE"},
	{"sequence/create_unknown_db", "CREATE SEQUENCE nope.a"},
	{"sequence/create_unknown_db_cycle", "CREATE SEQUENCE nope.a CYCLE"},
	{"sequence/create_unknown_db_conflict", "CREATE SEQUENCE nope.a CYCLE NOCYCLE"},
	{"sequence/create_other_case", "CREATE SEQUENCE A"},
	{"sequence/create_other_name", "CREATE SEQUENCE b"},
	{"sequence/create_name_64", "CREATE SEQUENCE " + strings.Repeat("b", 64)},
	{"sequence/create_name_65", "CREATE SEQUENCE " + strings.Repeat("b", 65)},

	{"sequence/alter_nomaxvalue", "ALTER SEQUENCE a NOMAXVALUE"},
	{"sequence/alter_nominvalue", "ALTER SEQUENCE a NOMINVALUE"},
	{"sequence/alter_maxvalue_one", "ALTER SEQUENCE a MAXVALUE 1"},
	{"sequence/alter_maxvalue", "ALTER SEQUENCE a MAXVALUE 1000"},
	{"sequence/alter_minvalue", "ALTER SEQUENCE a MINVALUE 1"},
	{"sequence/alter_minvalue_negative", "ALTER SEQUENCE a MINVALUE -50"},
	{"sequence/alter_start", "ALTER SEQUENCE a START WITH 1"},
	{"sequence/alter_restart", "ALTER SEQUENCE a RESTART"},
	{"sequence/alter_restart_start", "ALTER SEQUENCE a RESTART START WITH 5"},
	{"sequence/alter_restart_start_above_max", "ALTER SEQUENCE a RESTART START WITH 1000"},
	{"sequence/alter_increment_descending", "ALTER SEQUENCE a INCREMENT BY -1"},
	{"sequence/alter_increment", "ALTER SEQUENCE a INCREMENT BY 1"},
	{"sequence/alter_increment_zero", "ALTER SEQUENCE a INCREMENT BY 0"},
	{"sequence/alter_increment_descending_restart", "ALTER SEQUENCE a INCREMENT BY -1 RESTART"},
	{"sequence/alter_cycle", "ALTER SEQUENCE a CYCLE"},
	{"sequence/alter_nocycle", "ALTER SEQUENCE a NOCYCLE"},
	{"sequence/alter_cache_one", "ALTER SEQUENCE a CACHE 1"},
	{"sequence/alter_cache", "ALTER SEQUENCE a CACHE 2"},
	{"sequence/alter_cache_large", "ALTER SEQUENCE a CACHE 1000"},
	{"sequence/alter_nocache", "ALTER SEQUENCE a NOCACHE"},
	{"sequence/alter_order", "ALTER SEQUENCE a ORDER"},
	{"sequence/alter_noorder", "ALTER SEQUENCE a NOORDER"},
	{"sequence/alter_maxvalue_fraction", "ALTER SEQUENCE a MAXVALUE 1.5"},
	{"sequence/alter_nomaxvalue_cycle", "ALTER SEQUENCE a NOMAXVALUE CYCLE"},
	{"sequence/alter_conflict_cycle", "ALTER SEQUENCE a CYCLE NOCYCLE"},
	{"sequence/alter_dup_restart", "ALTER SEQUENCE a RESTART RESTART"},
	{"sequence/alter_other_case", "ALTER SEQUENCE A NOCACHE"},
	{"sequence/alter_missing", "ALTER SEQUENCE b NOCACHE"},
	{"sequence/alter_missing_no_option", "ALTER SEQUENCE b"},
	{"sequence/alter_unknown_db", "ALTER SEQUENCE nope.a NOCACHE"},
	{"sequence/alter_name_65", "ALTER SEQUENCE " + strings.Repeat("b", 65)},

	{"sequence/drop_other_case", "DROP SEQUENCE A"},
	{"sequence/drop_missing", "DROP SEQUENCE b"},
	{"sequence/drop_if_exists_missing", "DROP SEQUENCE IF EXISTS b"},
	{"sequence/drop_unknown_db", "DROP SEQUENCE nope.a"},
	{"sequence/drop_if_exists_unknown_db", "DROP SEQUENCE IF EXISTS nope.a"},
}

func sequenceSimulationRows(compatPath string) ([]simulate.Generated, error) {
	fromCompat, err := simulate.FromCompat(compatPath, []string{"_"}, sequenceSimulated)
	if err != nil {
		return nil, err
	}
	var statements []programSpec
	seen := make(map[string]bool)
	for _, r := range fromCompat {
		seen[r.Statement] = true
		statements = append(statements, programSpec{rule: r.Rule, stmt: r.Statement})
	}
	for _, s := range curatedSequenceStatements {
		if !seen[s.stmt] {
			seen[s.stmt] = true
			statements = append(statements, s)
		}
	}
	var rows []simulate.Generated
	for _, fx := range sequenceSimulationFixtures {
		for _, s := range statements {
			rows = append(rows, simulate.Generated{Fixture: fx, Rule: s.rule, Statement: s.stmt})
		}
	}
	return rows, nil
}

// sequenceSimulationEngine asks OceanBase, or the tape: the database dropped
// and built again with the fixture, the statement, and after one OceanBase
// took, the database's tables, views and sequences.
func sequenceSimulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.OpenWith(t, testcontainer.Options{MultiStatements: true})
	ctx := context.Background()
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		build := fmt.Sprintf("DROP DATABASE IF EXISTS %[1]s; CREATE DATABASE %[1]s; USE %[1]s; %[2]s", simulationDatabase, fx.SQL)
		if _, err := db.ExecContext(ctx, build); err != nil {
			return simulate.Answer{}, fmt.Errorf("build fixture %s: %w", fx.Name, err)
		}
		if _, err := db.ExecContext(ctx, statement); err != nil {
			var refusal *gomysql.MySQLError
			if !errors.As(err, &refusal) {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: strconv.Itoa(int(refusal.Number))}, nil
		}
		forms, err := simulationForms(ctx, db)
		if err != nil {
			return simulate.Answer{}, err
		}
		sequences, err := sequenceForms(ctx, db)
		return simulate.Answer{Forms: append(forms, sequences...)}, err
	}
}

// sequenceForms reads the selected database's sequences in name order. The
// view has no START WITH; LAST_NUMBER is the value the sequence hands out
// next, which for a sequence nothing has read is its start.
func sequenceForms(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	rows, err := db.QueryContext(ctx, "SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG "+
		"FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var forms []simulate.Form
	for rows.Next() {
		var name, start, lo, hi, step, cache, order, cycle string
		if err := rows.Scan(&name, &start, &lo, &hi, &step, &cache, &order, &cycle); err != nil {
			return nil, err
		}
		forms = append(forms, simulate.Form{Kind: "SEQUENCE", Name: name,
			Text: sequenceForm(name, start, lo, hi, step, cache, order == "Y", cycle == "Y")})
	}
	return forms, rows.Err()
}

// sequenceForm writes a sequence as one CREATE SEQUENCE with every option.
func sequenceForm(name, start, lo, hi, step, cache string, order, cycle bool) string {
	words := func(on bool, yes, no string) string {
		if on {
			return yes
		}
		return no
	}
	return fmt.Sprintf("CREATE SEQUENCE `%s` START WITH %s MINVALUE %s MAXVALUE %s INCREMENT BY %s CACHE %s %s %s",
		name, start, lo, hi, step, cache, words(order, "ORDER", "NOORDER"), words(cycle, "CYCLE", "NOCYCLE"))
}

// sequenceSimulationOmni runs the row in the catalog: its tables and views are
// held to OceanBase's as TestSimulation holds them, and its sequences to the
// ones OceanBase lists.
func sequenceSimulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	o := simulationOmni(fx, statement, engine)
	if engine.Refused() || o.Code != "" || !o.Same {
		return o
	}
	database := fmt.Sprintf("CREATE DATABASE %[1]s DEFAULT CHARSET=utf8mb4;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQL(database + fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	if _, err := c.Exec(statement, nil); err != nil {
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	// The pinned tenant runs with lower_case_table_names=1 and stores a
	// sequence's name folded, where the catalog keeps every name as written and
	// finds it by its folded key, so the names are held to each other folded,
	// as Diff holds tables'.
	var mine []simulate.Form
	for key, seq := range c.GetDatabase(simulationDatabase).Sequences {
		mine = append(mine, simulate.Form{Kind: "SEQUENCE", Name: key, Text: sequenceForm(key,
			seq.Start.String(), seq.Min.String(), seq.Max.String(), seq.Increment.String(), seq.Cache.String(), seq.Order, seq.Cycle)})
	}
	slices.SortFunc(mine, func(a, b simulate.Form) int { return strings.Compare(a.Name, b.Name) })
	theirs := slices.DeleteFunc(slices.Clone(engine.Forms), func(f simulate.Form) bool { return f.Kind != "SEQUENCE" })
	if slices.Equal(mine, theirs) {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: fmt.Sprintf("OceanBase stores the sequences %q, the catalog %q", formTexts(theirs), formTexts(mine))}
}

func formTexts(forms []simulate.Form) []string {
	var texts []string
	for _, f := range forms {
		texts = append(texts, f.Text)
	}
	return texts
}

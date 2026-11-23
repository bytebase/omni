package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/sijms/go-ora/v2/network"

	"github.com/bytebase/omni/enginetest/simulate"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
)

// TestSimulation holds the catalog to testdata/simulate.tsv: core schema DDL
// statements Oracle accepts in the parser's compat.tsv, run on each fixture,
// with what Oracle did and what Exec does. An outcome is the ORA- error number
// the statement is refused with, or ok: and a hash of the stored schema DDL
// read back from Oracle's data dictionary afterwards.
//
// Oracle's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned Oracle again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	// Keep only grammar-generated positive rules: compat.tsv's hand corpora
	// (suite:*) use manifest-specific table names rather than a, and near-miss
	// mutations (!insert/!swap/!truncate) only appear as accept when Oracle
	// raised a semantic lookup error before parsing the mutated tail.
	rows = slices.DeleteFunc(rows, func(g simulate.Generated) bool {
		return handLabel(g.Rule) || strings.Contains(g.Rule, "!")
	})
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   (testcontainer.Codec{}).MadeUp(),
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

func handLabel(rule string) bool {
	i := strings.IndexAny(rule, ":/")
	return i >= 0 && rule[i] == ':'
}

// simulationFixtures are the starting schemas in testdata/fixtures. A
// compat.tsv statement writes a for every identifier it holds, so what a
// statement finds is determined by what the fixture names a.
var simulationFixtures = []string{
	"one_table@missing", "one_table@table", "one_table@column", "one_table@index",
	"objects@view", "objects@sequence", "objects@synonym", "objects@fk",
}

// simulatedStatement selects core relational schema DDL from compat.tsv:
// tables, indexes, views, sequences, synonyms, table/column comments,
// TRUNCATE TABLE, and RENAME. Physical/storage-only ALTER TABLE actions,
// partitioning, object/XMLType tables & views, clusters, and database/server
// admin statements are left out so the suite stays focused and compact.
var simulatedStatement = regexp.MustCompile(`^(?:` +
	`CREATE TABLE (?:IF NOT EXISTS )?(?:a\.a|a)\s*(?:\(|AS SELECT )` +
	`|ALTER TABLE (?:IF EXISTS )?(?:a\.a|a) (?:ADD(?:\(| CONSTRAINT | PRIMARY | UNIQUE | FOREIGN | CHECK )|MODIFY(?:\(| [a-z]+ (?:VISIBLE|INVISIBLE|DROP IDENTITY))|DROP (?:COLUMN |CONSTRAINT |PRIMARY |UNIQUE|\()|RENAME (?:TO |COLUMN |CONSTRAINT )|SET UNUSED )` +
	`|DROP TABLE ` +
	`|CREATE (?:UNIQUE |BITMAP )?INDEX (?:IF NOT EXISTS )?(?:a\.a|a) ON (?:a\.a|a)(?: a)?\([^)]*\)$` +
	`|DROP INDEX ` +
	`|CREATE (?:OR REPLACE )?(?:(?:NO)?FORCE )?(?:(?:NON)?EDITIONABLE |EDITIONING )?VIEW (?:IF NOT EXISTS )?(?:a\.a|a)(?:\(| | AS )` +
	`|ALTER VIEW ` +
	`|DROP VIEW ` +
	`|CREATE SEQUENCE ` +
	`|ALTER SEQUENCE ` +
	`|DROP SEQUENCE ` +
	`|CREATE (?:OR REPLACE )?(?:(?:NON)?EDITIONABLE )?(?:PUBLIC )?SYNONYM ` +
	`|DROP (?:PUBLIC )?SYNONYM ` +
	`|COMMENT ON (?:TABLE|COLUMN) ` +
	`|TRUNCATE TABLE ` +
	`|RENAME ` +
	`)`)

func simulated(statement string) bool {
	if !simulatedStatement.MatchString(statement) {
		return false
	}
	if strings.Contains(statement, " OF ") || strings.Contains(statement, " XMLTYPE ") {
		return false
	}
	if strings.HasPrefix(statement, "CREATE TABLE a(PARTITION") {
		return false
	}
	return true
}

const simUser = "OMNI_SIM"

const simHelperPkgSpec = `CREATE OR REPLACE PACKAGE system.omni_sim AUTHID CURRENT_USER AS
  PROCEDURE reset_schema(p_fixture IN CLOB);
  FUNCTION dump_forms RETURN CLOB;
END omni_sim;`

const simHelperPkgBody = `CREATE OR REPLACE PACKAGE BODY system.omni_sim AS
  FUNCTION q_id(p_name IN VARCHAR2) RETURN VARCHAR2 IS
  BEGIN
    RETURN '"' || REPLACE(p_name, '"', '""') || '"';
  END;

  FUNCTION q_str(p_val IN VARCHAR2) RETURN VARCHAR2 IS
  BEGIN
    RETURN '''' || REPLACE(p_val, '''', '''''') || '''';
  END;

  FUNCTION fmt_type(
    p_type IN VARCHAR2,
    p_len  IN NUMBER,
    p_prec IN NUMBER,
    p_scal IN NUMBER,
    p_cuse IN VARCHAR2,
    p_clen IN NUMBER
  ) RETURN VARCHAR2 IS
  BEGIN
    IF p_type IN ('VARCHAR2', 'CHAR') THEN
      IF p_cuse = 'C' AND p_clen IS NOT NULL THEN
        RETURN p_type || '(' || p_clen || ' CHAR)';
      END IF;
      RETURN p_type || '(' || p_len || ' BYTE)';
    ELSIF p_type IN ('NVARCHAR2', 'NCHAR') THEN
      RETURN p_type || '(' || NVL(p_clen, p_len) || ')';
    ELSIF p_type IN ('RAW', 'UROWID') THEN
      RETURN p_type || '(' || p_len || ')';
    ELSIF p_type = 'NUMBER' THEN
      IF p_prec IS NULL THEN
        IF p_scal IS NOT NULL THEN
          RETURN 'NUMBER(*,' || p_scal || ')';
        END IF;
        RETURN 'NUMBER';
      ELSIF p_prec > 0 THEN
        IF p_scal IS NOT NULL AND p_scal <> 0 THEN
          RETURN 'NUMBER(' || p_prec || ',' || p_scal || ')';
        END IF;
        RETURN 'NUMBER(' || p_prec || ')';
      END IF;
      RETURN 'NUMBER';
    ELSIF p_type = 'FLOAT' THEN
      IF p_prec IS NOT NULL AND p_prec > 0 THEN
        RETURN 'FLOAT(' || p_prec || ')';
      END IF;
      RETURN 'FLOAT';
    END IF;
    RETURN p_type;
  END;

  PROCEDURE reset_schema(p_fixture IN CLOB) IS
    v_left VARCHAR2(4000);
    v_pos  PLS_INTEGER := 1;
    v_semi PLS_INTEGER;
    v_len  PLS_INTEGER;
    v_stmt VARCHAR2(32767);
  BEGIN
    FOR pass IN 1 .. 3 LOOP
      FOR o IN (
        SELECT object_name, object_type
        FROM user_objects
        WHERE object_type NOT IN (
          'LOB', 'LOB PARTITION', 'LOB SUBPARTITION',
          'TABLE PARTITION', 'TABLE SUBPARTITION',
          'INDEX PARTITION', 'INDEX SUBPARTITION'
        )
        ORDER BY DECODE(object_type,
          'MATERIALIZED VIEW', 1, 'TRIGGER', 2, 'VIEW', 3, 'SYNONYM', 4,
          'SEQUENCE', 5, 'TABLE', 6, 'TYPE', 8, 7), object_name
      ) LOOP
        BEGIN
          IF o.object_type = 'TABLE' THEN
            EXECUTE IMMEDIATE 'ALTER TABLE ' || q_id(o.object_name) || ' ENABLE TABLE LOCK';
          END IF;
        EXCEPTION
          WHEN OTHERS THEN NULL;
        END;
        BEGIN
          EXECUTE IMMEDIATE 'DROP ' || o.object_type || ' ' || q_id(o.object_name) ||
            CASE o.object_type
              WHEN 'TABLE' THEN ' CASCADE CONSTRAINTS PURGE'
              WHEN 'TYPE'  THEN ' FORCE'
              ELSE ''
            END;
        EXCEPTION
          WHEN OTHERS THEN NULL;
        END;
      END LOOP;
    END LOOP;

    FOR ps IN (
      SELECT synonym_name
      FROM all_synonyms
      WHERE owner = 'PUBLIC' AND (table_owner = USER OR synonym_name = 'A')
      ORDER BY synonym_name
    ) LOOP
      BEGIN
        EXECUTE IMMEDIATE 'DROP PUBLIC SYNONYM ' || q_id(ps.synonym_name);
      EXCEPTION
        WHEN OTHERS THEN NULL;
      END;
    END LOOP;

    EXECUTE IMMEDIATE 'PURGE RECYCLEBIN';

    SELECT LISTAGG(object_type || ' ' || object_name, ', ') WITHIN GROUP (ORDER BY object_type, object_name)
      INTO v_left FROM user_objects;
    IF v_left IS NOT NULL THEN
      RAISE_APPLICATION_ERROR(-20998, 'omni_sim: schema not empty: ' || v_left);
    END IF;

    IF p_fixture IS NULL THEN
      RETURN;
    END IF;
    v_len := DBMS_LOB.GETLENGTH(p_fixture);
    WHILE v_pos <= v_len LOOP
      v_semi := DBMS_LOB.INSTR(p_fixture, ';', v_pos);
      IF v_semi = 0 THEN
        v_stmt := TRIM(CHR(10) FROM TRIM(CHR(13) FROM TRIM(DBMS_LOB.SUBSTR(p_fixture, v_len - v_pos + 1, v_pos))));
        v_pos := v_len + 1;
      ELSE
        v_stmt := TRIM(CHR(10) FROM TRIM(CHR(13) FROM TRIM(DBMS_LOB.SUBSTR(p_fixture, v_semi - v_pos, v_pos))));
        v_pos := v_semi + 1;
      END IF;
      IF v_stmt IS NOT NULL THEN
        EXECUTE IMMEDIATE v_stmt;
      END IF;
    END LOOP;
  END;

  FUNCTION dump_forms RETURN CLOB IS
    v_out   CLOB;
    v_buf   CLOB;
    v_first BOOLEAN;
    v_cols  VARCHAR2(32767);
    v_rcols VARCHAR2(32767);
    v_rtab  VARCHAR2(256);
    v_rown  VARCHAR2(256);
    v_def   VARCHAR2(32767);
    v_expr  VARCHAR2(32767);
    v_vtxt  VARCHAR2(32767);
    v_idgen VARCHAR2(64);
    v_idopt VARCHAR2(4000);
    v_icnt  PLS_INTEGER;
    PROCEDURE append_line(p_line IN VARCHAR2) IS
    BEGIN
      IF p_line IS NOT NULL AND LENGTH(p_line) > 0 THEN
        DBMS_LOB.WRITEAPPEND(v_out, LENGTH(p_line), p_line);
      END IF;
    END;
  BEGIN
    DBMS_LOB.CREATETEMPORARY(v_out, TRUE);

    -- 1. SEQUENCES
    FOR s IN (
      SELECT sequence_name, min_value, max_value, increment_by, cycle_flag, order_flag, cache_size
      FROM user_sequences
      WHERE sequence_name NOT LIKE 'ISEQ$$_%'
      ORDER BY sequence_name
    ) LOOP
      append_line('===SEQUENCE' || CHR(9) || s.sequence_name || CHR(10));
      append_line(
        'CREATE SEQUENCE ' || q_id(s.sequence_name) ||
        ' INCREMENT BY ' || s.increment_by ||
        ' MINVALUE ' || s.min_value ||
        ' MAXVALUE ' || s.max_value ||
        CASE WHEN s.cache_size = 0 THEN ' NOCACHE' ELSE ' CACHE ' || s.cache_size END ||
        CASE WHEN s.cycle_flag = 'Y' THEN ' CYCLE' ELSE ' NOCYCLE' END ||
        CASE WHEN s.order_flag = 'Y' THEN ' ORDER' ELSE ' NOORDER' END ||
        ';' || CHR(10)
      );
    END LOOP;

    -- 2. TABLES
    FOR t IN (
      SELECT table_name
      FROM user_tables
      WHERE dropped = 'NO'
      ORDER BY table_name
    ) LOOP
      append_line('===TABLE' || CHR(9) || t.table_name || CHR(10));
      append_line('CREATE TABLE ' || q_id(t.table_name) || ' (');
      v_first := TRUE;
      FOR c IN (
        SELECT column_name, data_type, data_length, data_precision, data_scale,
               char_used, char_length, data_default, nullable, collation,
               default_on_null, virtual_column, identity_column, hidden_column
        FROM user_tab_cols
        WHERE table_name = t.table_name AND user_generated = 'YES'
        ORDER BY column_id NULLS LAST, internal_column_id
      ) LOOP
        IF NOT v_first THEN
          append_line(', ');
        END IF;
        v_first := FALSE;
        append_line(q_id(c.column_name) || ' ' ||
          fmt_type(c.data_type, c.data_length, c.data_precision, c.data_scale, c.char_used, c.char_length));
        IF c.collation IS NOT NULL AND c.collation <> 'USING_NLS_COMP' THEN
          append_line(' COLLATE ' || c.collation);
        END IF;
        IF c.identity_column = 'YES' THEN
          SELECT ic.generation_type,
                 'INCREMENT BY ' || sq.increment_by ||
                 ' MINVALUE ' || sq.min_value ||
                 ' MAXVALUE ' || sq.max_value ||
                 CASE WHEN sq.cache_size = 0 THEN ' NOCACHE' ELSE ' CACHE ' || sq.cache_size END ||
                 CASE WHEN sq.cycle_flag = 'Y' THEN ' CYCLE' ELSE ' NOCYCLE' END ||
                 CASE WHEN sq.order_flag = 'Y' THEN ' ORDER' ELSE ' NOORDER' END
          INTO v_idgen, v_idopt
          FROM user_tab_identity_cols ic
          JOIN user_sequences sq ON sq.sequence_name = ic.sequence_name
          WHERE ic.table_name = t.table_name AND ic.column_name = c.column_name;
          append_line(' GENERATED ' || v_idgen ||
            CASE WHEN c.default_on_null = 'YES' THEN ' ON NULL' ELSE '' END ||
            ' AS IDENTITY (' || v_idopt || ')');
        ELSIF c.virtual_column = 'YES' THEN
          v_def := TRIM(c.data_default);
          append_line(' GENERATED ALWAYS AS (' || v_def || ') VIRTUAL');
        ELSIF c.data_default IS NOT NULL THEN
          v_def := TRIM(c.data_default);
          append_line(' DEFAULT ' || CASE WHEN c.default_on_null = 'YES' THEN 'ON NULL ' ELSE '' END || v_def);
        END IF;
        IF c.hidden_column = 'YES' THEN
          append_line(' INVISIBLE');
        END IF;
        IF c.nullable = 'N' THEN
          append_line(' NOT NULL');
        END IF;
      END LOOP;
      append_line(');' || CHR(10));

      -- Constraints
      FOR k IN (
        SELECT constraint_name, constraint_type, search_condition_vc, generated,
               r_owner, r_constraint_name, delete_rule, status
        FROM user_constraints
        WHERE table_name = t.table_name AND constraint_type IN ('P', 'U', 'R', 'C')
        ORDER BY constraint_name
      ) LOOP
        IF k.constraint_type = 'C' THEN
          IF k.generated = 'GENERATED NAME' AND REGEXP_LIKE(TRIM(k.search_condition_vc), '^"[^"]+" IS NOT NULL$') THEN
            CONTINUE;
          END IF;
          append_line('ALTER TABLE ' || q_id(t.table_name) || ' ADD CONSTRAINT ' || q_id(k.constraint_name) ||
            ' CHECK (' || TRIM(k.search_condition_vc) || ')' ||
            CASE WHEN k.status = 'DISABLED' THEN ' DISABLE' ELSE '' END || ';' || CHR(10));
        ELSIF k.constraint_type IN ('P', 'U') THEN
          SELECT LISTAGG('"' || REPLACE(column_name, '"', '""') || '"', ', ') WITHIN GROUP (ORDER BY position)
            INTO v_cols
            FROM user_cons_columns
            WHERE constraint_name = k.constraint_name;
          append_line('ALTER TABLE ' || q_id(t.table_name) || ' ADD CONSTRAINT ' || q_id(k.constraint_name) ||
            CASE k.constraint_type WHEN 'P' THEN ' PRIMARY KEY (' ELSE ' UNIQUE (' END ||
            v_cols || ')' ||
            CASE WHEN k.status = 'DISABLED' THEN ' DISABLE' ELSE '' END || ';' || CHR(10));
        ELSIF k.constraint_type = 'R' THEN
          SELECT LISTAGG('"' || REPLACE(column_name, '"', '""') || '"', ', ') WITHIN GROUP (ORDER BY position)
            INTO v_cols
            FROM user_cons_columns
            WHERE constraint_name = k.constraint_name;
          SELECT MAX(table_name), LISTAGG('"' || REPLACE(column_name, '"', '""') || '"', ', ') WITHIN GROUP (ORDER BY position)
            INTO v_rtab, v_rcols
            FROM all_cons_columns
            WHERE owner = k.r_owner AND constraint_name = k.r_constraint_name;
          v_rown := CASE WHEN k.r_owner <> USER THEN q_id(k.r_owner) || '.' ELSE '' END;
          append_line('ALTER TABLE ' || q_id(t.table_name) || ' ADD CONSTRAINT ' || q_id(k.constraint_name) ||
            ' FOREIGN KEY (' || v_cols || ') REFERENCES ' || v_rown || q_id(v_rtab) || ' (' || v_rcols || ')' ||
            CASE WHEN k.delete_rule IN ('CASCADE', 'SET NULL') THEN ' ON DELETE ' || k.delete_rule ELSE '' END ||
            CASE WHEN k.status = 'DISABLED' THEN ' DISABLE' ELSE '' END || ';' || CHR(10));
        END IF;
      END LOOP;

      -- Standalone indexes
      FOR ix IN (
        SELECT index_name, index_type, uniqueness, visibility
        FROM user_indexes
        WHERE table_name = t.table_name
          AND index_name NOT IN (
            SELECT index_name FROM user_constraints
            WHERE table_name = t.table_name AND index_name IS NOT NULL
          )
        ORDER BY index_name
      ) LOOP
        v_cols := '';
        v_first := TRUE;
        FOR ic IN (
          SELECT column_name, column_position, descend
          FROM user_ind_columns
          WHERE index_name = ix.index_name
          ORDER BY column_position
        ) LOOP
          IF NOT v_first THEN
            v_cols := v_cols || ', ';
          END IF;
          v_first := FALSE;
          IF ic.column_name LIKE 'SYS_NC%' THEN
            SELECT column_expression INTO v_expr
            FROM user_ind_expressions
            WHERE index_name = ix.index_name AND column_position = ic.column_position;
            v_cols := v_cols || TRIM(v_expr);
          ELSE
            v_cols := v_cols || q_id(ic.column_name);
          END IF;
          IF ic.descend = 'DESC' THEN
            v_cols := v_cols || ' DESC';
          END IF;
        END LOOP;
        append_line('CREATE ' ||
          CASE WHEN ix.uniqueness = 'UNIQUE' THEN 'UNIQUE ' ELSE '' END ||
          CASE WHEN ix.index_type LIKE '%BITMAP%' THEN 'BITMAP ' ELSE '' END ||
          'INDEX ' || q_id(ix.index_name) || ' ON ' || q_id(t.table_name) || ' (' || v_cols || ')' ||
          CASE WHEN ix.visibility = 'INVISIBLE' THEN ' INVISIBLE' ELSE '' END || ';' || CHR(10));
      END LOOP;

      -- Table and column comments
      FOR tc IN (
        SELECT comments FROM user_tab_comments
        WHERE table_name = t.table_name AND table_type = 'TABLE' AND comments IS NOT NULL
      ) LOOP
        append_line('COMMENT ON TABLE ' || q_id(t.table_name) || ' IS ' || q_str(tc.comments) || ';' || CHR(10));
      END LOOP;
      FOR cc IN (
        SELECT column_name, comments FROM user_col_comments
        WHERE table_name = t.table_name AND comments IS NOT NULL
        ORDER BY column_name
      ) LOOP
        append_line('COMMENT ON COLUMN ' || q_id(t.table_name) || '.' || q_id(cc.column_name) ||
          ' IS ' || q_str(cc.comments) || ';' || CHR(10));
      END LOOP;
    END LOOP;

    -- 3. VIEWS
    FOR v IN (
      SELECT view_name, text, read_only
      FROM user_views
      ORDER BY view_name
    ) LOOP
      append_line('===VIEW' || CHR(9) || v.view_name || CHR(10));
      v_cols := '';
      v_first := TRUE;
      FOR vc IN (
        SELECT column_name, hidden_column
        FROM user_tab_cols
        WHERE table_name = v.view_name
        ORDER BY internal_column_id
      ) LOOP
        IF NOT v_first THEN
          v_cols := v_cols || ', ';
        END IF;
        v_first := FALSE;
        v_cols := v_cols || q_id(vc.column_name) ||
          CASE WHEN vc.hidden_column = 'YES' THEN ' INVISIBLE' ELSE '' END;
      END LOOP;
      v_vtxt := RTRIM(TRIM(v.text), ';' || CHR(10) || CHR(13) || ' ');
      IF v.read_only = 'Y' AND INSTR(UPPER(v_vtxt), 'READ ONLY') = 0 THEN
        v_vtxt := v_vtxt || ' WITH READ ONLY';
      END IF;
      append_line('CREATE VIEW ' || q_id(v.view_name) ||
        CASE WHEN v_cols IS NOT NULL THEN ' (' || v_cols || ')' ELSE '' END ||
        ' AS ' || v_vtxt || ';' || CHR(10));
      FOR tc IN (
        SELECT comments FROM user_tab_comments
        WHERE table_name = v.view_name AND table_type = 'VIEW' AND comments IS NOT NULL
      ) LOOP
        append_line('COMMENT ON TABLE ' || q_id(v.view_name) || ' IS ' || q_str(tc.comments) || ';' || CHR(10));
      END LOOP;
    END LOOP;

    -- 4. PRIVATE SYNONYMS
    FOR sy IN (
      SELECT synonym_name, table_owner, table_name, db_link
      FROM user_synonyms
      ORDER BY synonym_name
    ) LOOP
      append_line('===SYNONYM' || CHR(9) || sy.synonym_name || CHR(10));
      append_line('CREATE SYNONYM ' || q_id(sy.synonym_name) || ' FOR ' ||
        CASE WHEN sy.table_owner IS NOT NULL AND (sy.table_owner <> USER OR sy.db_link IS NOT NULL)
             THEN q_id(sy.table_owner) || '.' ELSE '' END ||
        q_id(sy.table_name) ||
        CASE WHEN sy.db_link IS NOT NULL THEN '@' || sy.db_link ELSE '' END || ';' || CHR(10));
    END LOOP;

    -- 5. PUBLIC SYNONYMS CREATED BY THIS USER
    FOR sy IN (
      SELECT synonym_name, table_owner, table_name, db_link
      FROM all_synonyms
      WHERE owner = 'PUBLIC' AND (table_owner = USER OR synonym_name = 'A')
      ORDER BY synonym_name
    ) LOOP
      append_line('===PUBLIC SYNONYM' || CHR(9) || sy.synonym_name || CHR(10));
      append_line('CREATE PUBLIC SYNONYM ' || q_id(sy.synonym_name) || ' FOR ' ||
        CASE WHEN sy.table_owner IS NOT NULL AND (sy.table_owner <> USER OR sy.db_link IS NOT NULL)
             THEN q_id(sy.table_owner) || '.' ELSE '' END ||
        q_id(sy.table_name) ||
        CASE WHEN sy.db_link IS NOT NULL THEN '@' || sy.db_link ELSE '' END || ';' || CHR(10));
    END LOOP;

    RETURN v_out;
  END;
END omni_sim;`

func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	ctx := context.Background()
	sysDB := testcontainer.Open(t, "system")
	t.Cleanup(func() { _ = sysDB.Close() })

	setupStmts := []string{
		fmt.Sprintf(`CREATE USER "%s" IDENTIFIED BY "%s" QUOTA UNLIMITED ON USERS`, simUser, testcontainer.Password),
		fmt.Sprintf(`GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PUBLIC SYNONYM, DROP PUBLIC SYNONYM TO "%s"`, simUser),
		simHelperPkgSpec,
		simHelperPkgBody,
		fmt.Sprintf(`GRANT EXECUTE ON system.omni_sim TO "%s"`, simUser),
	}
	for _, stmt := range setupStmts {
		if _, err := sysDB.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("Oracle simulation setup (%s): %v", stmt[:min(60, len(stmt))], err)
		}
	}

	userDB := testcontainer.Open(t, simUser)
	t.Cleanup(func() { _ = userDB.Close() })

	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		reset := fmt.Sprintf("BEGIN system.omni_sim.reset_schema(q'[%s]'); END;", fx.SQL)
		if _, err := userDB.ExecContext(ctx, reset); err != nil {
			return simulate.Answer{}, fmt.Errorf("build fixture %s: %w", fx.Name, err)
		}
		if _, err := userDB.ExecContext(ctx, statement); err != nil {
			var oe *network.OracleError
			if !errors.As(err, &oe) {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: strconv.Itoa(oe.ErrCode)}, nil
		}
		var raw sql.NullString
		if err := userDB.QueryRowContext(ctx, "SELECT system.omni_sim.dump_forms FROM dual").Scan(&raw); err != nil {
			return simulate.Answer{}, fmt.Errorf("read stored schema: %w", err)
		}
		return simulate.Answer{Forms: parseSimulationForms(raw.String)}, nil
	}
}

func parseSimulationForms(raw string) []simulate.Form {
	var forms []simulate.Form
	var cur *simulate.Form
	var body strings.Builder
	flush := func() {
		if cur != nil {
			cur.Text = body.String()
			forms = append(forms, *cur)
			cur = nil
			body.Reset()
		}
	}
	for _, line := range strings.Split(raw, "\n") {
		if strings.HasPrefix(line, "===") {
			flush()
			kind, name, _ := strings.Cut(strings.TrimPrefix(line, "==="), "\t")
			cur = &simulate.Form{Kind: kind, Name: name}
			continue
		}
		if cur != nil && line != "" {
			body.WriteString(line)
			body.WriteByte('\n')
		}
	}
	flush()
	return forms
}

func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c, err := LoadSQL(fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	// Tag the loaded fixture with the scratch schema owner so cross-schema
	// references (e.g. a.a) resolve against OMNI_SIM rather than adopting a.
	c.meta.Name = simUser
	if len(c.meta.Schemas) > 0 {
		c.meta.Schemas[0].Name = simUser
	}
	results, err := c.Exec(statement, nil)
	if err != nil {
		var refusal *Error
		if errors.As(err, &refusal) {
			if code := oracleErrorCode(refusal.Code); code != "" {
				return simulate.Omni{Code: code, Why: err.Error()}
			}
		}
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	for _, r := range results {
		var refusal *Error
		switch {
		case r.Error == nil:
		case errors.As(r.Error, &refusal):
			if code := oracleErrorCode(refusal.Code); code != "" {
				return simulate.Omni{Code: code, Why: r.Error.Error()}
			}
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		default:
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
	}
	if engine.Refused() {
		return simulate.Omni{}
	}

	want, err := loadSimulationForms(engine.Forms)
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what Oracle stores: %v", err)}
	}
	gotNorm := normalizeForSimulation(c.Metadata())
	wantNorm := normalizeForSimulation(want.Metadata())
	alignGeneratedConstraintNames(gotNorm, wantNorm)

	gotCat := &Catalog{meta: gotNorm}
	wantCat := &Catalog{meta: wantNorm}
	there := Diff(gotCat, wantCat)
	back := Diff(wantCat, gotCat)
	if there.IsEmpty() && back.IsEmpty() {
		return simulate.Omni{Same: true}
	}
	plan := GenerateMigration(gotCat, wantCat, nil).SQL()
	if there.IsEmpty() {
		plan = "nothing, though from there it would run: " + GenerateMigration(wantCat, gotCat, nil).SQL()
	}
	return simulate.Omni{Why: "to store what Oracle stores, the catalog would run: " + strings.Join(strings.Fields(plan), " ")}
}

func oracleErrorCode(code string) string {
	code = strings.TrimSpace(code)
	code = strings.TrimPrefix(strings.ToUpper(code), "ORA-")
	n, err := strconv.Atoi(code)
	if err != nil || n <= 0 {
		return ""
	}
	return strconv.Itoa(n)
}

func loadSimulationForms(forms []simulate.Form) (*Catalog, error) {
	// Statements in Oracle's dump_forms are grouped per object in alphabetical
	// order, so a table's FOREIGN KEY (e.g. A -> P or A_FK before A_PK) can
	// precede the referenced table or key. Partition into dependency phases:
	//   1. Sequences and CREATE TABLE statements
	//   2. Non-FK ALTER TABLE (PK, UNIQUE, CHECK), CREATE INDEX, COMMENT ON
	//   3. Foreign-key ALTER TABLE statements
	//   4. Views (as CREATE FORCE VIEW, since Oracle already stored them) and synonyms
	var phase1, phase2, phase3, phase4, phase5 []string
	for _, f := range forms {
		for _, rawLine := range strings.Split(f.Text, "\n") {
			line := strings.TrimSpace(rawLine)
			if line == "" {
				continue
			}
			// Oracle's default ascending MAXVALUE (10^28 - 1) overflows int64
			// inside oracle/parser's identity clause parser; strip it only from
			// AS IDENTITY (...) clauses while keeping CREATE SEQUENCE intact.
			if strings.Contains(line, "AS IDENTITY (") {
				line = strings.ReplaceAll(line, " MAXVALUE 9999999999999999999999999999", "")
			}
			upper := strings.ToUpper(line)
			switch {
			case strings.HasPrefix(upper, "CREATE SEQUENCE "), strings.HasPrefix(upper, "CREATE TABLE "):
				phase1 = append(phase1, line)
			case strings.HasPrefix(upper, "ALTER TABLE ") && strings.Contains(upper, " FOREIGN KEY "):
				phase3 = append(phase3, line)
			case strings.HasPrefix(upper, "CREATE VIEW "):
				phase4 = append(phase4, "CREATE FORCE VIEW "+line[len("CREATE VIEW "):])
			case strings.HasPrefix(upper, "CREATE SYNONYM "), strings.HasPrefix(upper, "CREATE PUBLIC SYNONYM "):
				phase4 = append(phase4, line)
			case strings.HasPrefix(upper, "COMMENT ON "):
				phase5 = append(phase5, line)
			default:
				phase2 = append(phase2, line)
			}
		}
	}
	var all []string
	all = append(all, phase1...)
	all = append(all, phase2...)
	all = append(all, phase3...)
	all = append(all, phase4...)
	all = append(all, phase5...)
	if len(all) == 0 {
		return New(), nil
	}
	return LoadSQL(strings.Join(all, "\n") + "\n")
}

var oracleMadeUpConstraint = regexp.MustCompile(`^SYS_C\d{6,}$`)

func normalizeForSimulation(meta *metadatapb.DatabaseSchemaMetadata) *metadatapb.DatabaseSchemaMetadata {
	out := cloneMetadata(meta)
	out.Name = ""
	if len(out.Schemas) == 0 {
		out.Schemas = []*metadatapb.SchemaMetadata{{}}
	}
	for _, s := range out.Schemas {
		s.Name = ""
		for _, seq := range s.Sequences {
			normalizeSequenceDefaults(seq)
		}
		for _, t := range s.Tables {
			for _, col := range t.Columns {
				if col.IsIdentity {
					normalizeIdentityOptions(col)
				}
			}
		}
		for _, syn := range s.Synonyms {
			if strings.EqualFold(syn.TargetSchema, simUser) && syn.TargetDbLink == "" {
				syn.TargetSchema = ""
			}
		}
	}
	for _, syn := range out.PublicSynonyms {
		if strings.EqualFold(syn.TargetSchema, simUser) && syn.TargetDbLink == "" {
			syn.TargetSchema = ""
		}
	}
	return out
}

func normalizeSequenceDefaults(seq *metadatapb.SequenceMetadata) {
	if seq.Start == "1" {
		seq.Start = ""
	}
	if seq.Increment == "1" {
		seq.Increment = ""
	}
	if seq.MinValue == "1" {
		seq.MinValue = ""
	}
	if seq.MaxValue == "9999999999999999999999999999" {
		seq.MaxValue = ""
	}
	if seq.CacheSize == "20" {
		seq.CacheSize = ""
	}
}

func normalizeIdentityOptions(col *metadatapb.ColumnMetadata) {
	if col.IdentitySeed == 0 {
		col.IdentitySeed = 1
	}
	if col.IdentityIncrement == 0 {
		col.IdentityIncrement = 1
	}
	var filtered []string
	for _, opt := range col.Options {
		switch opt {
		case "MINVALUE=1", "MAXVALUE=9999999999999999999999999999", "CACHE=20", "CYCLE=false", "ORDER=false":
			continue
		}
		filtered = append(filtered, opt)
	}
	col.Options = filtered
}

// alignGeneratedConstraintNames rewrites system-generated Oracle constraint
// names (SYS_C...) in want to match omni's synthetic names (PK_*, UK_*, FK_*,
// CHK_*) when their definitions match structurally. Explicitly named
// constraints are untouched and compared by exact name.
func alignGeneratedConstraintNames(got, want *metadatapb.DatabaseSchemaMetadata) {
	if len(got.Schemas) == 0 || len(want.Schemas) == 0 {
		return
	}
	gotTables := tableByName(got.Schemas[0])
	for _, wt := range want.Schemas[0].Tables {
		gt := gotTables[wt.Name]
		if gt == nil {
			continue
		}
		// Indexes (PK/UK constraints)
		usedGotIdx := map[int]bool{}
		for _, wi := range wt.Indexes {
			if !wi.IsConstraint || !oracleMadeUpConstraint.MatchString(wi.Name) {
				continue
			}
			for j, gi := range gt.Indexes {
				if usedGotIdx[j] || !gi.IsConstraint || !isOmniSyntheticConstraintName(gi.Name, gt.Name) {
					continue
				}
				if indexEqual(gi, wi) {
					wi.Name = gi.Name
					usedGotIdx[j] = true
					break
				}
			}
		}
		// Foreign keys
		usedGotFK := map[int]bool{}
		for _, wfk := range wt.ForeignKeys {
			if !oracleMadeUpConstraint.MatchString(wfk.Name) {
				continue
			}
			for j, gfk := range gt.ForeignKeys {
				if usedGotFK[j] || !isOmniSyntheticConstraintName(gfk.Name, gt.Name) {
					continue
				}
				if foreignKeyEqual(gfk, wfk) {
					wfk.Name = gfk.Name
					usedGotFK[j] = true
					break
				}
			}
		}
		// Check constraints
		usedGotChk := map[int]bool{}
		for _, wchk := range wt.CheckConstraints {
			if !oracleMadeUpConstraint.MatchString(wchk.Name) {
				continue
			}
			for j, gchk := range gt.CheckConstraints {
				if usedGotChk[j] || !isOmniSyntheticConstraintName(gchk.Name, gt.Name) {
					continue
				}
				if normalizeExpression(gchk.Expression) == normalizeExpression(wchk.Expression) {
					wchk.Name = gchk.Name
					usedGotChk[j] = true
					break
				}
			}
		}
	}
}

func isOmniSyntheticConstraintName(name, table string) bool {
	return name == "PK_"+table ||
		strings.HasPrefix(name, "UK_"+table+"_") ||
		strings.HasPrefix(name, "FK_"+table+"_") ||
		strings.HasPrefix(name, "CHK_"+table+"_")
}

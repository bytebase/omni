-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestSimulation default#1
> exec
| CREATE USER "OMNI_SIM" IDENTIFIED BY "omni_test" QUOTA UNLIMITED ON USERS
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE PUBLIC SYNONYM, DROP PUBLIC SYNONYM TO "OMNI_SIM"
< ok 0 0
> exec
| CREATE OR REPLACE PACKAGE system.omni_sim AUTHID CURRENT_USER AS
|   PROCEDURE reset_schema(p_fixture IN CLOB);
|   FUNCTION dump_forms RETURN CLOB;
| END omni_sim;
< ok 0 0
> exec
| CREATE OR REPLACE PACKAGE BODY system.omni_sim AS
|   FUNCTION q_id(p_name IN VARCHAR2) RETURN VARCHAR2 IS
|   BEGIN
|     RETURN '"' || REPLACE(p_name, '"', '""') || '"';
|   END;
|
|   FUNCTION q_str(p_val IN VARCHAR2) RETURN VARCHAR2 IS
|   BEGIN
|     RETURN '''' || REPLACE(p_val, '''', '''''') || '''';
|   END;
|
|   FUNCTION fmt_type(
|     p_type IN VARCHAR2,
|     p_len  IN NUMBER,
|     p_prec IN NUMBER,
|     p_scal IN NUMBER,
|     p_cuse IN VARCHAR2,
|     p_clen IN NUMBER
|   ) RETURN VARCHAR2 IS
|   BEGIN
|     IF p_type IN ('VARCHAR2', 'CHAR') THEN
|       IF p_cuse = 'C' AND p_clen IS NOT NULL THEN
|         RETURN p_type || '(' || p_clen || ' CHAR)';
|       END IF;
|       RETURN p_type || '(' || p_len || ' BYTE)';
|     ELSIF p_type IN ('NVARCHAR2', 'NCHAR') THEN
|       RETURN p_type || '(' || NVL(p_clen, p_len) || ')';
|     ELSIF p_type IN ('RAW', 'UROWID') THEN
|       RETURN p_type || '(' || p_len || ')';
|     ELSIF p_type = 'NUMBER' THEN
|       IF p_prec IS NULL THEN
|         IF p_scal IS NOT NULL THEN
|           RETURN 'NUMBER(*,' || p_scal || ')';
|         END IF;
|         RETURN 'NUMBER';
|       ELSIF p_prec > 0 THEN
|         IF p_scal IS NOT NULL AND p_scal <> 0 THEN
|           RETURN 'NUMBER(' || p_prec || ',' || p_scal || ')';
|         END IF;
|         RETURN 'NUMBER(' || p_prec || ')';
|       END IF;
|       RETURN 'NUMBER';
|     ELSIF p_type = 'FLOAT' THEN
|       IF p_prec IS NOT NULL AND p_prec > 0 THEN
|         RETURN 'FLOAT(' || p_prec || ')';
|       END IF;
|       RETURN 'FLOAT';
|     END IF;
|     RETURN p_type;
|   END;
|
|   PROCEDURE reset_schema(p_fixture IN CLOB) IS
|     v_left VARCHAR2(4000);
|     v_pos  PLS_INTEGER := 1;
|     v_semi PLS_INTEGER;
|     v_len  PLS_INTEGER;
|     v_stmt VARCHAR2(32767);
|   BEGIN
|     FOR pass IN 1 .. 3 LOOP
|       FOR o IN (
|         SELECT object_name, object_type
|         FROM user_objects
|         WHERE object_type NOT IN (
|           'LOB', 'LOB PARTITION', 'LOB SUBPARTITION',
|           'TABLE PARTITION', 'TABLE SUBPARTITION',
|           'INDEX PARTITION', 'INDEX SUBPARTITION'
|         )
|         ORDER BY DECODE(object_type,
|           'MATERIALIZED VIEW', 1, 'TRIGGER', 2, 'VIEW', 3, 'SYNONYM', 4,
|           'SEQUENCE', 5, 'TABLE', 6, 'TYPE', 8, 7), object_name
|       ) LOOP
|         BEGIN
|           IF o.object_type = 'TABLE' THEN
|             EXECUTE IMMEDIATE 'ALTER TABLE ' || q_id(o.object_name) || ' ENABLE TABLE LOCK';
|           END IF;
|         EXCEPTION
|           WHEN OTHERS THEN NULL;
|         END;
|         BEGIN
|           EXECUTE IMMEDIATE 'DROP ' || o.object_type || ' ' || q_id(o.object_name) ||
|             CASE o.object_type
|               WHEN 'TABLE' THEN ' CASCADE CONSTRAINTS PURGE'
|               WHEN 'TYPE'  THEN ' FORCE'
|               ELSE ''
|             END;
|         EXCEPTION
|           WHEN OTHERS THEN NULL;
|         END;
|       END LOOP;
|     END LOOP;
|
|     FOR ps IN (
|       SELECT synonym_name
|       FROM all_synonyms
|       WHERE owner = 'PUBLIC' AND (table_owner = USER OR synonym_name = 'A')
|       ORDER BY synonym_name
|     ) LOOP
|       BEGIN
|         EXECUTE IMMEDIATE 'DROP PUBLIC SYNONYM ' || q_id(ps.synonym_name);
|       EXCEPTION
|         WHEN OTHERS THEN NULL;
|       END;
|     END LOOP;
|
|     EXECUTE IMMEDIATE 'PURGE RECYCLEBIN';
|
|     SELECT LISTAGG(object_type || ' ' || object_name, ', ') WITHIN GROUP (ORDER BY object_type, object_name)
|       INTO v_left FROM user_objects;
|     IF v_left IS NOT NULL THEN
|       RAISE_APPLICATION_ERROR(-20998, 'omni_sim: schema not empty: ' || v_left);
|     END IF;
|
|     IF p_fixture IS NULL THEN
|       RETURN;
|     END IF;
|     v_len := DBMS_LOB.GETLENGTH(p_fixture);
|     WHILE v_pos <= v_len LOOP
|       v_semi := DBMS_LOB.INSTR(p_fixture, ';', v_pos);
|       IF v_semi = 0 THEN
|         v_stmt := TRIM(CHR(10) FROM TRIM(CHR(13) FROM TRIM(DBMS_LOB.SUBSTR(p_fixture, v_len - v_pos + 1, v_pos))));
|         v_pos := v_len + 1;
|       ELSE
|         v_stmt := TRIM(CHR(10) FROM TRIM(CHR(13) FROM TRIM(DBMS_LOB.SUBSTR(p_fixture, v_semi - v_pos, v_pos))));
|         v_pos := v_semi + 1;
|       END IF;
|       IF v_stmt IS NOT NULL THEN
|         EXECUTE IMMEDIATE v_stmt;
|       END IF;
|     END LOOP;
|   END;
|
|   FUNCTION dump_forms RETURN CLOB IS
|     v_out   CLOB;
|     v_buf   CLOB;
|     v_first BOOLEAN;
|     v_cols  VARCHAR2(32767);
|     v_rcols VARCHAR2(32767);
|     v_rtab  VARCHAR2(256);
|     v_rown  VARCHAR2(256);
|     v_def   VARCHAR2(32767);
|     v_expr  VARCHAR2(32767);
|     v_vtxt  VARCHAR2(32767);
|     v_idgen VARCHAR2(64);
|     v_idopt VARCHAR2(4000);
|     v_icnt  PLS_INTEGER;
|     PROCEDURE append_line(p_line IN VARCHAR2) IS
|     BEGIN
|       IF p_line IS NOT NULL AND LENGTH(p_line) > 0 THEN
|         DBMS_LOB.WRITEAPPEND(v_out, LENGTH(p_line), p_line);
|       END IF;
|     END;
|   BEGIN
|     DBMS_LOB.CREATETEMPORARY(v_out, TRUE);
|
|     -- 1. SEQUENCES
|     FOR s IN (
|       SELECT sequence_name, min_value, max_value, increment_by, cycle_flag, order_flag, cache_size
|       FROM user_sequences
|       WHERE sequence_name NOT LIKE 'ISEQ$$_%'
|       ORDER BY sequence_name
|     ) LOOP
|       append_line('===SEQUENCE' || CHR(9) || s.sequence_name || CHR(10));
|       append_line(
|         'CREATE SEQUENCE ' || q_id(s.sequence_name) ||
|         ' INCREMENT BY ' || s.increment_by ||
|         ' MINVALUE ' || s.min_value ||
|         ' MAXVALUE ' || s.max_value ||
|         CASE WHEN s.cache_size = 0 THEN ' NOCACHE' ELSE ' CACHE ' || s.cache_size END ||
|         CASE WHEN s.cycle_flag = 'Y' THEN ' CYCLE' ELSE ' NOCYCLE' END ||
|         CASE WHEN s.order_flag = 'Y' THEN ' ORDER' ELSE ' NOORDER' END ||
|         ';' || CHR(10)
|       );
|     END LOOP;
|
|     -- 2. TABLES
|     FOR t IN (
|       SELECT table_name
|       FROM user_tables
|       WHERE dropped = 'NO'
|       ORDER BY table_name
|     ) LOOP
|       append_line('===TABLE' || CHR(9) || t.table_name || CHR(10));
|       append_line('CREATE TABLE ' || q_id(t.table_name) || ' (');
|       v_first := TRUE;
|       FOR c IN (
|         SELECT column_name, data_type, data_length, data_precision, data_scale,
|                char_used, char_length, data_default, nullable, collation,
|                default_on_null, virtual_column, identity_column, hidden_column
|         FROM user_tab_cols
|         WHERE table_name = t.table_name AND user_generated = 'YES'
|         ORDER BY column_id NULLS LAST, internal_column_id
|       ) LOOP
|         IF NOT v_first THEN
|           append_line(', ');
|         END IF;
|         v_first := FALSE;
|         append_line(q_id(c.column_name) || ' ' ||
|           fmt_type(c.data_type, c.data_length, c.data_precision, c.data_scale, c.char_used, c.char_length));
|         IF c.collation IS NOT NULL AND c.collation <> 'USING_NLS_COMP' THEN
|           append_line(' COLLATE ' || c.collation);
|         END IF;
|         IF c.identity_column = 'YES' THEN
|           SELECT ic.generation_type,
|                  'INCREMENT BY ' || sq.increment_by ||
|                  ' MINVALUE ' || sq.min_value ||
|                  ' MAXVALUE ' || sq.max_value ||
|                  CASE WHEN sq.cache_size = 0 THEN ' NOCACHE' ELSE ' CACHE ' || sq.cache_size END ||
|                  CASE WHEN sq.cycle_flag = 'Y' THEN ' CYCLE' ELSE ' NOCYCLE' END ||
|                  CASE WHEN sq.order_flag = 'Y' THEN ' ORDER' ELSE ' NOORDER' END
|           INTO v_idgen, v_idopt
|           FROM user_tab_identity_cols ic
|           JOIN user_sequences sq ON sq.sequence_name = ic.sequence_name
|           WHERE ic.table_name = t.table_name AND ic.column_name = c.column_name;
|           append_line(' GENERATED ' || v_idgen ||
|             CASE WHEN c.default_on_null = 'YES' THEN ' ON NULL' ELSE '' END ||
|             ' AS IDENTITY (' || v_idopt || ')');
|         ELSIF c.virtual_column = 'YES' THEN
|           v_def := TRIM(c.data_default);
|           append_line(' GENERATED ALWAYS AS (' || v_def || ') VIRTUAL');
|         ELSIF c.data_default IS NOT NULL THEN
|           v_def := TRIM(c.data_default);
|           append_line(' DEFAULT ' || CASE WHEN c.default_on_null = 'YES' THEN 'ON NULL ' ELSE '' END || v_def);
|         END IF;
|         IF c.hidden_column = 'YES' THEN
|           append_line(' INVISIBLE');
|         END IF;
|         IF c.nullable = 'N' THEN
|           append_line(' NOT NULL');
|         END IF;
|       END LOOP;
|       append_line(');' || CHR(10));
|
|       -- Constraints
|       -- A period's check constraint is named after the period, its hidden
|       -- virtual column, and six hex digits Oracle draws anew on every run
|       -- (A453584, then A482f68), with GENERATED = 'USER NAME'. It is read as
|       -- <period>$PERIOD so the stored schema reads the same on every run.
|       FOR k IN (
|         SELECT constraint_name, constraint_type, search_condition_vc, generated,
|                r_owner, r_constraint_name, delete_rule, status,
|                NVL2(period_name, period_name || '$PERIOD', constraint_name) AS shown_name
|         FROM (
|           SELECT uc.*, (
|               SELECT MAX(p.column_name)
|               FROM user_tab_cols p
|               WHERE p.table_name = uc.table_name
|                 AND p.hidden_column = 'YES' AND p.user_generated = 'NO' AND p.virtual_column = 'YES'
|                 AND uc.constraint_type = 'C'
|                 AND LENGTH(uc.constraint_name) = LENGTH(p.column_name) + 6
|                 AND SUBSTR(uc.constraint_name, 1, LENGTH(p.column_name)) = p.column_name
|                 AND REGEXP_LIKE(SUBSTR(uc.constraint_name, -6), '^[0-9a-f]{6}$', 'c')
|             ) AS period_name
|           FROM user_constraints uc
|           WHERE uc.table_name = t.table_name AND uc.constraint_type IN ('P', 'U', 'R', 'C')
|         )
|         ORDER BY shown_name
|       ) LOOP
|         IF k.constraint_type = 'C' THEN
|           IF k.generated = 'GENERATED NAME' AND REGEXP_LIKE(TRIM(k.search_condition_vc), '^"[^"]+" IS NOT NULL$') THEN
|             CONTINUE;
|           END IF;
|           append_line('ALTER TABLE ' || q_id(t.table_name) || ' ADD CONSTRAINT ' || q_id(k.shown_name) ||
|             ' CHECK (' || TRIM(k.search_condition_vc) || ')' ||
|             CASE WHEN k.status = 'DISABLED' THEN ' DISABLE' ELSE '' END || ';' || CHR(10));
|         ELSIF k.constraint_type IN ('P', 'U') THEN
|           SELECT LISTAGG('"' || REPLACE(column_name, '"', '""') || '"', ', ') WITHIN GROUP (ORDER BY position)
|             INTO v_cols
|             FROM user_cons_columns
|             WHERE constraint_name = k.constraint_name;
|           append_line('ALTER TABLE ' || q_id(t.table_name) || ' ADD CONSTRAINT ' || q_id(k.constraint_name) ||
|             CASE k.constraint_type WHEN 'P' THEN ' PRIMARY KEY (' ELSE ' UNIQUE (' END ||
|             v_cols || ')' ||
|             CASE WHEN k.status = 'DISABLED' THEN ' DISABLE' ELSE '' END || ';' || CHR(10));
|         ELSIF k.constraint_type = 'R' THEN
|           SELECT LISTAGG('"' || REPLACE(column_name, '"', '""') || '"', ', ') WITHIN GROUP (ORDER BY position)
|             INTO v_cols
|             FROM user_cons_columns
|             WHERE constraint_name = k.constraint_name;
|           SELECT MAX(table_name), LISTAGG('"' || REPLACE(column_name, '"', '""') || '"', ', ') WITHIN GROUP (ORDER BY position)
|             INTO v_rtab, v_rcols
|             FROM all_cons_columns
|             WHERE owner = k.r_owner AND constraint_name = k.r_constraint_name;
|           v_rown := CASE WHEN k.r_owner <> USER THEN q_id(k.r_owner) || '.' ELSE '' END;
|           append_line('ALTER TABLE ' || q_id(t.table_name) || ' ADD CONSTRAINT ' || q_id(k.constraint_name) ||
|             ' FOREIGN KEY (' || v_cols || ') REFERENCES ' || v_rown || q_id(v_rtab) || ' (' || v_rcols || ')' ||
|             CASE WHEN k.delete_rule IN ('CASCADE', 'SET NULL') THEN ' ON DELETE ' || k.delete_rule ELSE '' END ||
|             CASE WHEN k.status = 'DISABLED' THEN ' DISABLE' ELSE '' END || ';' || CHR(10));
|         END IF;
|       END LOOP;
|
|       -- Standalone indexes
|       FOR ix IN (
|         SELECT index_name, index_type, uniqueness, visibility
|         FROM user_indexes
|         WHERE table_name = t.table_name
|           AND index_name NOT IN (
|             SELECT index_name FROM user_constraints
|             WHERE table_name = t.table_name AND index_name IS NOT NULL
|           )
|         ORDER BY index_name
|       ) LOOP
|         v_cols := '';
|         v_first := TRUE;
|         FOR ic IN (
|           SELECT column_name, column_position, descend
|           FROM user_ind_columns
|           WHERE index_name = ix.index_name
|           ORDER BY column_position
|         ) LOOP
|           IF NOT v_first THEN
|             v_cols := v_cols || ', ';
|           END IF;
|           v_first := FALSE;
|           IF ic.column_name LIKE 'SYS_NC%' THEN
|             SELECT column_expression INTO v_expr
|             FROM user_ind_expressions
|             WHERE index_name = ix.index_name AND column_position = ic.column_position;
|             v_cols := v_cols || TRIM(v_expr);
|           ELSE
|             v_cols := v_cols || q_id(ic.column_name);
|           END IF;
|           IF ic.descend = 'DESC' THEN
|             v_cols := v_cols || ' DESC';
|           END IF;
|         END LOOP;
|         append_line('CREATE ' ||
|           CASE WHEN ix.uniqueness = 'UNIQUE' THEN 'UNIQUE ' ELSE '' END ||
|           CASE WHEN ix.index_type LIKE '%BITMAP%' THEN 'BITMAP ' ELSE '' END ||
|           'INDEX ' || q_id(ix.index_name) || ' ON ' || q_id(t.table_name) || ' (' || v_cols || ')' ||
|           CASE WHEN ix.visibility = 'INVISIBLE' THEN ' INVISIBLE' ELSE '' END || ';' || CHR(10));
|       END LOOP;
|
|       -- Table and column comments
|       FOR tc IN (
|         SELECT comments FROM user_tab_comments
|         WHERE table_name = t.table_name AND table_type = 'TABLE' AND comments IS NOT NULL
|       ) LOOP
|         append_line('COMMENT ON TABLE ' || q_id(t.table_name) || ' IS ' || q_str(tc.comments) || ';' || CHR(10));
|       END LOOP;
|       FOR cc IN (
|         SELECT column_name, comments FROM user_col_comments
|         WHERE table_name = t.table_name AND comments IS NOT NULL
|         ORDER BY column_name
|       ) LOOP
|         append_line('COMMENT ON COLUMN ' || q_id(t.table_name) || '.' || q_id(cc.column_name) ||
|           ' IS ' || q_str(cc.comments) || ';' || CHR(10));
|       END LOOP;
|     END LOOP;
|
|     -- 3. VIEWS
|     FOR v IN (
|       SELECT view_name, text, read_only
|       FROM user_views
|       ORDER BY view_name
|     ) LOOP
|       append_line('===VIEW' || CHR(9) || v.view_name || CHR(10));
|       v_cols := '';
|       v_first := TRUE;
|       FOR vc IN (
|         SELECT column_name, hidden_column
|         FROM user_tab_cols
|         WHERE table_name = v.view_name
|         ORDER BY internal_column_id
|       ) LOOP
|         IF NOT v_first THEN
|           v_cols := v_cols || ', ';
|         END IF;
|         v_first := FALSE;
|         v_cols := v_cols || q_id(vc.column_name) ||
|           CASE WHEN vc.hidden_column = 'YES' THEN ' INVISIBLE' ELSE '' END;
|       END LOOP;
|       v_vtxt := RTRIM(TRIM(v.text), ';' || CHR(10) || CHR(13) || ' ');
|       IF v.read_only = 'Y' AND INSTR(UPPER(v_vtxt), 'READ ONLY') = 0 THEN
|         v_vtxt := v_vtxt || ' WITH READ ONLY';
|       END IF;
|       append_line('CREATE VIEW ' || q_id(v.view_name) ||
|         CASE WHEN v_cols IS NOT NULL THEN ' (' || v_cols || ')' ELSE '' END ||
|         ' AS ' || v_vtxt || ';' || CHR(10));
|       FOR tc IN (
|         SELECT comments FROM user_tab_comments
|         WHERE table_name = v.view_name AND table_type = 'VIEW' AND comments IS NOT NULL
|       ) LOOP
|         append_line('COMMENT ON TABLE ' || q_id(v.view_name) || ' IS ' || q_str(tc.comments) || ';' || CHR(10));
|       END LOOP;
|     END LOOP;
|
|     -- 4. PRIVATE SYNONYMS
|     FOR sy IN (
|       SELECT synonym_name, table_owner, table_name, db_link
|       FROM user_synonyms
|       ORDER BY synonym_name
|     ) LOOP
|       append_line('===SYNONYM' || CHR(9) || sy.synonym_name || CHR(10));
|       append_line('CREATE SYNONYM ' || q_id(sy.synonym_name) || ' FOR ' ||
|         CASE WHEN sy.table_owner IS NOT NULL AND (sy.table_owner <> USER OR sy.db_link IS NOT NULL)
|              THEN q_id(sy.table_owner) || '.' ELSE '' END ||
|         q_id(sy.table_name) ||
|         CASE WHEN sy.db_link IS NOT NULL THEN '@' || sy.db_link ELSE '' END || ';' || CHR(10));
|     END LOOP;
|
|     -- 5. PUBLIC SYNONYMS CREATED BY THIS USER
|     FOR sy IN (
|       SELECT synonym_name, table_owner, table_name, db_link
|       FROM all_synonyms
|       WHERE owner = 'PUBLIC' AND (table_owner = USER OR synonym_name = 'A')
|       ORDER BY synonym_name
|     ) LOOP
|       append_line('===PUBLIC SYNONYM' || CHR(9) || sy.synonym_name || CHR(10));
|       append_line('CREATE PUBLIC SYNONYM ' || q_id(sy.synonym_name) || ' FOR ' ||
|         CASE WHEN sy.table_owner IS NOT NULL AND (sy.table_owner <> USER OR sy.db_link IS NOT NULL)
|              THEN q_id(sy.table_owner) || '.' ELSE '' END ||
|         q_id(sy.table_name) ||
|         CASE WHEN sy.db_link IS NOT NULL THEN '@' || sy.db_link ELSE '' END || ';' || CHR(10));
|     END LOOP;
|
|     RETURN v_out;
|   END;
| END omni_sim;
< ok 0 0
> exec
| GRANT EXECUTE ON system.omni_sim TO "OMNI_SIM"
< ok 0 0

== TestSimulation default#2
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 4043 ORA-04043: Object A does not exist.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4004 ORA-04004: The MINVALUE integer value is not less than the MAXVALUE integer value.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4015 ORA-04015: MAXVALUE must be specified for ascending sequences that cycle.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 NOCACHE NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4010 ORA-04010: the number of values to CACHE must be greater than 1
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE ORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "A"@A;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ();
| ALTER TABLE "A" ADD CONSTRAINT "A$PERIOD" CHECK ("A_START" < "A_END");
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("A" NUMBER);
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 54039 ORA-54039: table must have at least one column that is not invisible
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "T_D_IDX" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX t_d_idx ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A$PERIOD" CHECK ("A_START" < "A_END");
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 23292 ORA-23292: The constraint does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "C" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A$PERIOD" CHECK (("A" < "A") and ("A" > 0));
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 2443 ORA-02443: Cannot drop constraint  - nonexistent constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "C" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "C" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| CREATE UNIQUE INDEX "A_PK" ON "A" ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "C" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "C" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER INVISIBLE);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER DEFAULT 1);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS IDENTITY (INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER) NOT NULL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER DEFAULT ON NULL 1 NOT NULL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 28365 ORA-28365: Wallet is not open.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 22893 ORA-22893: constraint can be specified only for REF columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL INVISIBLE);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 904 ORA-00904: "PARTITION": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 904 ORA-00904: "PARTITION": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 22893 ORA-22893: constraint can be specified only for REF columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000001" FOREIGN KEY ("A") REFERENCES "A" ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "A" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000002" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE CASCADE;
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000003" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE SET NULL;
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 38802 ORA-38802: edition does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 38802 ORA-38802: edition does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 38802 ORA-38802: edition does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER, "A" NUMBER GENERATED ALWAYS AS (1) VIRTUAL);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| COMMENT ON TABLE "A" IS 's';
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" (1);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 23292 ORA-23292: The constraint does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "A" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" UNIQUE ("A");
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 2260 ORA-02260: table can have only one primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000004" FOREIGN KEY ("A") REFERENCES "A" ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000005" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE CASCADE;
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000006" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE SET NULL;
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "A" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 2443 ORA-02443: Cannot drop constraint  - nonexistent constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 2442 ORA-02442: Cannot drop nonexistent unique key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "A" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "A" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| CREATE UNIQUE INDEX "A_PK" ON "A" ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "A" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "A" NUMBER, "D" NUMBER);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER DEFAULT 1, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER DEFAULT ON NULL 1 NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 28365 ORA-28365: Wallet is not open.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 28333 ORA-28333: The column is not encrypted.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 1451 ORA-01451: column to be modified to NULL cannot be modified to NULL
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER, "A" NUMBER INVISIBLE);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 904 ORA-00904: "PARTITION": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 904 ORA-00904: "PARTITION": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 28365 ORA-28365: Wallet is not open.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 28365 ORA-28365: Wallet is not open.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 28365 ORA-28365: Wallet is not open.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 1451 ORA-01451: column to be modified to NULL cannot be modified to NULL
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000007" UNIQUE ("A");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 2260 ORA-02260: table can have only one primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000008" FOREIGN KEY ("A") REFERENCES "A" ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000009" CHECK (1 = 1);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000010" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE CASCADE;
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000011" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE SET NULL;
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| COMMENT ON COLUMN "A"."A" IS 's';
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| COMMENT ON TABLE "A" IS 's';
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A" DESC);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" (1);
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: null
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A_D_IDX" ON "A" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 4043 ORA-04043: Object A does not exist.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4004 ORA-04004: The MINVALUE integer value is not less than the MAXVALUE integer value.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4015 ORA-04015: MAXVALUE must be specified for ascending sequences that cycle.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 NOCACHE NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4010 ORA-04010: the number of values to CACHE must be greater than 1
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE ORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "A"@A;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ();
| ALTER TABLE "A" ADD CONSTRAINT "A$PERIOD" CHECK ("A_START" < "A_END");
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("A" NUMBER);
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 2268 ORA-02268: referenced table does not have a primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 54039 ORA-54039: table must have at least one column that is not invisible
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
| ===VIEW	A
| CREATE VIEW "A" ("A") AS SELECT 1 FROM dual;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 29862 ORA-29862: cannot specify FORCE option for dropping non-domain index
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" ("D");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a ON t (d);]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP VIEW a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 42310 ORA-42310: Implementation restriction: this operation is only permitted on Editioning views
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 25129 ORA-25129: cannot modify constraint (A) - no such constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 2443 ORA-02443: Cannot drop constraint  - nonexistent constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 42310 ORA-42310: Implementation restriction: this operation is only permitted on Editioning views
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 25129 ORA-25129: cannot modify constraint (A) - no such constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
| COMMENT ON TABLE "A" IS 's';
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
| ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 54039 ORA-54039: table must have at least one column that is not invisible
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===VIEW	A
| CREATE VIEW "A" ("B", "A", "D") AS SELECT b, a, d FROM t;
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE VIEW a AS SELECT b, a, d FROM t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< error
| 1702 ORA-01702: a view is not appropriate here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 4004 ORA-04004: The MINVALUE integer value is not less than the MAXVALUE integer value.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 4007 ORA-04007: MINVALUE cannot exceed the current value
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 4007 ORA-04007: MINVALUE cannot exceed the current value
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 CYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 NOCACHE NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 4010 ORA-04010: the number of values to CACHE must be greater than 1
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE ORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===SEQUENCE	A
| CREATE SEQUENCE "A" INCREMENT BY 1 MINVALUE 1 MAXVALUE 9999999999999999999999999999 CACHE 20 NOCYCLE NOORDER;
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SEQUENCE a START WITH 1 INCREMENT BY 1;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< error
| 2201 ORA-02201: sequence not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| COMMENT ON TABLE "T" IS 's';
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "T" (1);
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	T
| CREATE TABLE "T" ("B" NUMBER NOT NULL, "C" NUMBER, "D" NUMBER);
| ALTER TABLE "T" ADD CONSTRAINT "T_PK" PRIMARY KEY ("B");
| ===SYNONYM	A
| CREATE SYNONYM "A" FOR "T";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER, d NUMBER);
| CREATE SYNONYM a FOR t;]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a COMPILE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SYNONYM a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP TABLE a
< error
| 2449 ORA-02449: unique/primary keys in table referenced by foreign keys
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP VIEW a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| RENAME a TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a KEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOKEEP
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSCALE
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE EXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SCALE NOEXTEND
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a NOSHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a VISIBLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME TO a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME COLUMN a TO a
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO a
< error
| 2264 ORA-02264: name already used by an existing constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 2273 ORA-02273: this unique/primary key is referenced by some foreign keys
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER, a NUMBER)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1))
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1), a AS(1))
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a INVALIDATE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CASCADE CONSTRAINTS
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a SET UNUSED COLUMN a ONLINE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP COLUMN a CHECKPOINT 1
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(PERIOD FOR a(a, a))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a UNIQUE(a)
< error
| 2264 ORA-02264: name already used by an existing constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(a)
< error
| 2260 ORA-02260: table can have only one primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000012" FOREIGN KEY ("A") REFERENCES "A" ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a.a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a(a)
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000013" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE CASCADE;
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON DELETE SET NULL
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000014" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE SET NULL;
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD CONSTRAINT a FOREIGN KEY(a) REFERENCES a
< error
| 2264 ORA-02264: name already used by an existing constraint
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP CONSTRAINT a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a)
< error
| 2442 ORA-02442: Cannot drop nonexistent unique key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP UNIQUE(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY KEEP INDEX
< error
| 2273 ORA-02273: this unique/primary key is referenced by some foreign keys
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY DROP INDEX
< error
| 2273 ORA-02273: this unique/primary key is referenced by some foreign keys
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a DROP PRIMARY KEY ONLINE
< error
| 2297 ORA-02297: cannot disable constraint (OMNI_SIM.A_PK) - dependencies exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SORT)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER VISIBLE)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER INVISIBLE)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT 1)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER GENERATED ALWAYS AS IDENTITY)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER DEFAULT ON NULL 1)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ENCRYPT)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER NULL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER WITH ROWID)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER ANNOTATIONS(ADD a))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a VISIBLE AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a INVISIBLE AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a GENERATED ALWAYS AS(1))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) VIRTUAL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) NULL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a)
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT 1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER DEFAULT 1, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DROP IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DEFAULT ON NULL 1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER DEFAULT ON NULL 1 NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT)
< error
| 28335 ORA-28335: referenced or referencing FK constraint column cannot be encrypted
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a DECRYPT)
< error
| 28333 ORA-28333: The column is not encrypted.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NULL)
< error
| 1451 ORA-01451: column to be modified to NULL cannot be modified to NULL
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ANNOTATIONS(ADD a))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NUMBER AS(1))
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a COLLATE a AS(1))
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS(1))
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) VIRTUAL)
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a AS(1) NULL)
< error
| 54026 ORA-54026: Real column cannot have an expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY a INVISIBLE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "D" NUMBER, "A" NUMBER INVISIBLE);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION)
< error
| 904 ORA-00904: "PARTITION": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(PARTITION a)
< error
| 904 ORA-00904: "PARTITION": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT IDENTIFIED BY a)
< error
| 28335 ORA-28335: referenced or referencing FK constraint column cannot be encrypted
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT SALT)
< error
| 28335 ORA-28335: referenced or referencing FK constraint column cannot be encrypted
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a ENCRYPT NO SALT)
< error
| 28335 ORA-28335: referenced or referencing FK constraint column cannot be encrypted
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CONSTRAINT a NULL)
< error
| 1451 ORA-01451: column to be modified to NULL cannot be modified to NULL
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a UNIQUE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000015" UNIQUE ("A");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a PRIMARY KEY)
< error
| 2260 ORA-02260: table can have only one primary key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000016" FOREIGN KEY ("A") REFERENCES "A" ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a CHECK(1 = 1))
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000017" CHECK (1 = 1);
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a NOT NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a.a)
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a(a))
< error
| 2270 ORA-02270: no matching unique or primary key for this column-list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE CASCADE)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000018" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE CASCADE;
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a REFERENCES a ON DELETE SET NULL)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ALTER TABLE "A" ADD CONSTRAINT "SYS_C9000019" FOREIGN KEY ("A") REFERENCES "A" ("B") ON DELETE SET NULL;
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER SCOPE IS a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a.a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER CONSTRAINT a REFERENCES a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a(a))
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE CASCADE)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a NUMBER REFERENCES a ON DELETE SET NULL)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING EDITION a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) EVALUATE USING NULL EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a ADD(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 1430 ORA-01430: column being added already exists in table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER TABLE a MODIFY(a GENERATED ALWAYS AS IDENTITY(START WITH LIMIT VALUE))
< error
| 30673 ORA-30673: column to be modified is not an identity column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW IF EXISTS a COMPILE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a.a COMPILE
< error
| 1917 ORA-01917: user or role '' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ ONLY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a RELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a DROP CONSTRAINT a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a EDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a READ WRITE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a MODIFY CONSTRAINT a NORELY
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a NONEDITIONABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| ALTER VIEW a ADD UNIQUE(a) DISABLE
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| COMMENT ON COLUMN a.a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| COMMENT ON COLUMN "A"."A" IS 's';
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| COMMENT ON TABLE a IS 's'
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| COMMENT ON TABLE "A" IS 's';
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a.a(a)
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a a(a)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a ASC)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a DESC)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" ("A" DESC);
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE INDEX a ON a(1)
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| CREATE INDEX "A" ON "A" (1);
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a.a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = METADATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = DATA
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARING = NONE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a KEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOKEEP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a NOSCALE
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE EXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SCALE NOEXTEND
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SESSION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a GLOBAL
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SEQUENCE a SHARD
< error
| 2511 ORA-02511: SQL query not allowed; the shard DDL is disabled.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE SYNONYM a FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE PUBLIC SYNONYM a FOR a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
| ===PUBLIC SYNONYM	A
| CREATE PUBLIC SYNONYM "A" FOR "A";
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a.a FOR a
< error
| 1031 ORA-01031: insufficient privileges
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = METADATA FOR a
< error
| 65021 ORA-65021: illegal use of SHARING clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a SHARING = NONE FOR a
< error
| 1471 ORA-01471: cannot create a synonym with same name as object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE SYNONYM a FOR a@a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a, a)
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a)
< error
| 11504 ORA-11504: The domain specified does not exist or the user does not have privileges on the domain for the operation.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DOMAIN a.a)
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SORT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a DEFAULT ON NULL FOR INSERT ONLY 1)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a WITH ROWID)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ANNOTATIONS(ADD a))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NUMBER AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a VISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a INVISIBLE AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS(1))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) VIRTUAL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) NULL)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) ANNOTATIONS(ADD a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(PERIOD FOR a(a, a))
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a UNIQUE)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a PRIMARY KEY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CHECK(1 = 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a NOT NULL)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a.a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a(a))
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE CASCADE)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a REFERENCES a ON DELETE SET NULL)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a SCOPE IS a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a CONSTRAINT a REFERENCES a)
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED BY DEFAULT ON NULL AS IDENTITY)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY())
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT IDENTIFIED BY a)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a ENCRYPT NO SALT)
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) EVALUATE USING EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEFORE EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH CURRENT EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH EDITION a)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a AS(1) UNUSABLE BEGINNING WITH NULL EDITION)
< error
| 54037 ORA-54037: table must have at least 1 column that is not virtual
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(START WITH 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(INCREMENT BY 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMAXVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MAXVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOMINVALUE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(MINVALUE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCYCLE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOCACHE))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(CACHE 1))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(ORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(NOORDER))
< error
| 2263 ORA-02263: need to specify the data type for this column
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE OR REPLACE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE FORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE NOFORCE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE EDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE NONEDITIONABLE VIEW a AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a.a AS SELECT 1 FROM dual
< error
| 1917 ORA-01917: user or role 'A' does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a VISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a INVISIBLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, a) AS SELECT 1 FROM dual
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a ANNOTATIONS(a) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a DEFAULT COLLATION a AS SELECT 1 FROM dual
< error
| 43929 ORA-43929: Collation cannot be specified if parameter MAX_STRING_SIZE=STANDARD is set.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH CURRENT_USER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a BEQUEATH DEFINER AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINER_MAP
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual CONTAINERS_DEFAULT
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a UNIQUE DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a(a, UNIQUE(a) DISABLE) AS SELECT 1 FROM dual
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH CHECK OPTION
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| CREATE VIEW a AS SELECT 1 FROM dual WITH READ ONLY CONSTRAINT a
< error
| 955 ORA-00955: name is already used by an existing object
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX a.a
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX a ONLINE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX a FORCE
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX a DEFERRED INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP INDEX a IMMEDIATE INVALIDATION
< error
| 1418 ORA-01418: specified index does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SEQUENCE a.a
< error
| 2289 ORA-02289: sequence does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP PUBLIC SYNONYM a
< error
| 1432 ORA-01432: public synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SYNONYM IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SYNONYM a.a
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP SYNONYM a FORCE
< error
| 1434 ORA-01434: private synonym to be dropped does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP TABLE IF EXISTS a
< error
| 2449 ORA-02449: unique/primary keys in table referenced by foreign keys
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP TABLE a CASCADE CONSTRAINTS
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP TABLE a PURGE
< error
| 2449 ORA-02449: unique/primary keys in table referenced by foreign keys
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP VIEW a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| DROP VIEW a CASCADE CONSTRAINTS
< error
| 942 ORA-00942: table or view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PRESERVE MATERIALIZED VIEW LOG
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a PURGE MATERIALIZED VIEW LOG
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a DROP ALL STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a REUSE STORAGE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY, d NUMBER);
| CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
| CREATE TABLE ch (b NUMBER CONSTRAINT ch_pk PRIMARY KEY, a_b NUMBER CONSTRAINT ch_fk REFERENCES a (b));]'); END;
< ok 1 0
> exec
| TRUNCATE TABLE a CASCADE
< ok 0 0
> query
| SELECT system.omni_sim.dump_forms FROM dual
< result
< column "DUMP_FORMS" "LongVarChar"
< row
: string
| ===TABLE	A
| CREATE TABLE "A" ("B" NUMBER NOT NULL, "A" NUMBER, "D" NUMBER);
| ALTER TABLE "A" ADD CONSTRAINT "A" FOREIGN KEY ("A") REFERENCES "P" ("B");
| ALTER TABLE "A" ADD CONSTRAINT "A_CHK" CHECK (d > 0);
| ALTER TABLE "A" ADD CONSTRAINT "A_PK" PRIMARY KEY ("B");
| ===TABLE	CH
| CREATE TABLE "CH" ("B" NUMBER NOT NULL, "A_B" NUMBER);
| ALTER TABLE "CH" ADD CONSTRAINT "CH_FK" FOREIGN KEY ("A_B") REFERENCES "A" ("B");
| ALTER TABLE "CH" ADD CONSTRAINT "CH_PK" PRIMARY KEY ("B");
| ===TABLE	P
| CREATE TABLE "P" ("B" NUMBER NOT NULL, "D" NUMBER);
| ALTER TABLE "P" ADD CONSTRAINT "P_PK" PRIMARY KEY ("B");
|

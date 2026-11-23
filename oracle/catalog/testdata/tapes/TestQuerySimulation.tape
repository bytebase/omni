-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestQuerySimulation default#1
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

== TestQuerySimulation default#2
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH FUNCTION a RETURN NUMBER IS BEGIN RETURN 1; END; SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH FUNCTION a RETURN NUMBER IS BEGIN RETURN 1; END; FUNCTION a RETURN NUMBER IS BEGIN RETURN 1; END; SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH a AS(SELECT *) SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH a ANALYTIC VIEW AS(USING a) SELECT *
< error
| 18385 ORA-18385: recursive WITH clause not allowed in a derived analytic view
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * UNION SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * INTERSECT SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * MINUS SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * UNION ALL SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE OF a
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE OF a.a
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE OF a.a.a
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE OF a, a
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE NOWAIT
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE WAIT 1
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FOR UPDATE SKIP LOCKED
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following 'CREATE VIEW v AS SELECT * '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH PROCEDURE a IS BEGIN NULL; END; SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH a(a) AS(SELECT *) SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH a AS(SELECT *), a AS(SELECT *) SELECT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH a ANALYTIC VIEW AS(USING a), a ANALYTIC VIEW AS(USING a) SELECT *
< error
| 18385 ORA-18385: recursive WITH clause not allowed in a derived analytic view
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT /*+ a */ *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DISTINCT *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT UNIQUE *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ALL *
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a zzz
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a zzz
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a zzz
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a zzz
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a CROSS JOIN a
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a CROSS JOIN a zzz
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a, a
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a, a zzz
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * WHERE 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * CONNECT BY 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * WINDOW a AS()
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * QUALIFY 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY a
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY 1 ASC
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY 1 DESC
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY 1 NULLS FIRST
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY 1 NULLS LAST
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * ORDER BY 1, 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * OFFSET 1 ROW
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * OFFSET 1 ROWS
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH FIRST ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH FIRST zzz ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH NEXT ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH NEXT zzz ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH FIRST 1 ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH FIRST 1 PERCENT ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH FIRST PERCENT ROW ONLY
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FETCH FIRST ROW WITH TIES
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a)
< error
| 18307 ORA-18307: analytic view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a.a)
< error
| 18307 ORA-18307: analytic view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES())
< error
| 18370 ORA-18370: The HIERARCHIES clause may only be specified for cubes.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a FILTER FACT(a TO 1 = 1))
< error
| 18307 ORA-18307: analytic view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT a.*
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1 zzz
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT 1 zzz
: string
| N
: string
| ZZZ
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT a.a.*
< error
| 964 ORA-00964: table name not in FROM list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT a
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1 a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT 1 a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1 AS a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT 1 AS a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM(a CROSS JOIN a)
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a a
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a a
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a CROSS JOIN a CROSS JOIN a
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a CROSS JOIN a CROSS JOIN a zzz
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL LEFT JOIN a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT JOIN a
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT JOIN a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT JOIN a
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL LEFT JOIN a zzz
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT JOIN a zzz
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT JOIN a zzz
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT JOIN a zzz
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a CROSS APPLY a
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a CROSS APPLY a zzz
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * START WITH 1 = 1 CONNECT BY 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * CONNECT BY NOCYCLE 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * CONNECT BY 1 = 1 START WITH 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY ROLLUP(1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY(1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY ROOT(1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY GROUPING SETS(1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY SETS(1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY 1, 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY 1 HAVING 1 = 1
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES(zzz))
< error
| 18370 ORA-18370: The HIERARCHIES clause may only be specified for cubes.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES(a))
< error
| 18370 ORA-18370: The HIERARCHIES clause may only be specified for cubes.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES(a.a))
< error
| 18370 ORA-18370: The HIERARCHIES clause may only be specified for cubes.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES(a, a))
< error
| 18370 ORA-18370: The HIERARCHIES clause may only be specified for cubes.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a FILTER FACT(a TO 1 = 1, a TO 1 = 1))
< error
| 18307 ORA-18307: analytic view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM(SELECT *)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM(SELECT *) zzz
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM TABLE(1)
< error
| 22905 ORA-22905: Cannot access rows from unsupported input to TABLE operator.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM TABLE(1) zzz
< error
| 22905 ORA-22905: Cannot access rows from unsupported input to TABLE operator.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM TABLES(1)
< error
| 904 ORA-00904: "TABLES": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a) zzz
< error
| 18307 ORA-18307: analytic view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM zzz a
< error
| 942 ORA-00942: table or view "OMNI_SIM"."ZZZ" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a.a
< error
| 942 ORA-00942: table or view "A"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a.zzz a
< error
| 942 ORA-00942: table or view "A"."ZZZ" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PARTITION(a)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a(a)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PARTITION(a) zzz
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a@a
< error
| 2019 ORA-02019: connection description for remote database not found
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a@a zzz
< error
| 2019 ORA-02019: connection description for remote database not found
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a AS OF SCN 1
< error
| 8181 ORA-08181: specified number is not a valid system change number
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a AS OF SCN 1 zzz
< error
| 8181 ORA-08181: specified number is not a valid system change number
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SAMPLE(1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1)
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1)
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a(1)
< error
| 904 ORA-00904: "A": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SAMPLE(1) zzz
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) zzz
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) zzz
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) zzz
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT((1) FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a() FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN())
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1)) zzz
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR a IN(a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR a IN(a)) zzz
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1) zzz
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a JOIN a ON 1 = 1
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL JOIN a
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL JOIN a zzz
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL JOIN a zzz
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL JOIN a zzz
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL JOIN a zzz
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL INNER JOIN a
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a JOIN a USING(a)
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a JOIN a USING(a, a)
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a OUTER APPLY a
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY CUBE(1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY GROUPING SETS(ROLLUP(1))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY GROUPING SETS((1))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY GROUPING SETS(ROOT(1))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY GROUPING SETS(1, 1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * WINDOW a AS(PARTITION BY 1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * WINDOW a AS(ORDER BY 1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a ADD MEASURES(a AS(1)))
< error
| 18307 ORA-18307: analytic view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PARTITION FOR(1)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PARTITION(1)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PARTITION FOR(1) zzz
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SUBPARTITION(a)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SUBPARTITION(a) zzz
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SUBPARTITION FOR(1)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SUBPARTITION(1)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SUBPARTITION FOR(1) zzz
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PARTITION FOR(1, 1)
< error
| 14501 ORA-14501: object is not partitioned
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM TABLE(1)(+)
< error
| 22905 ORA-22905: Cannot access rows from unsupported input to TABLE operator.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM TABLE(1)(+) zzz
< error
| 22905 ORA-22905: Cannot access rows from unsupported input to TABLE operator.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM ANALYTIC VIEW(USING a) a
< error
| 18307 ORA-18307: analytic view "OMNI_SIM"."A" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a VERSIONS BETWEEN SCN 1 AND 1
< error
| 8181 ORA-08181: specified number is not a valid system change number
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a VERSIONS BETWEEN SCN 1 AND 1 zzz
< error
| 8181 ORA-08181: specified number is not a valid system change number
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a VERSIONS PERIOD FOR a BETWEEN 1 AND 1
< error
| 55603 ORA-55603: invalid flashback archive or valid time period command
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a VERSIONS PERIOD FOR a BETWEEN 1 AND 1 zzz
< error
| 55603 ORA-55603: invalid flashback archive or valid time period command
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a AS OF PERIOD FOR a 1
< error
| 55603 ORA-55603: invalid flashback archive or valid time period command
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a AS OF PERIOD FOR a 1 zzz
< error
| 55603 ORA-55603: invalid flashback archive or valid time period command
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a AS OF TIMESTAMP 1
< error
| 932 ORA-00932: expression is of data type NUMBER, which is incompatible with expected data type TIMESTAMP
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SAMPLE BLOCK(1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE BLOCK(1)
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE BLOCK(1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE BLOCK(1)
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SAMPLE
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SAMPLE(1) SEED(1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) SEED(1)
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) SEED(1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) SEED(1)
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a SAMPLE(1) SEED(1) zzz
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) SEED(1) zzz
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) SEED(1) zzz
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a SAMPLE(1) SEED(1) zzz
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a PIVOT
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a PIVOT
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a PIVOT
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) zzz FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) a FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) AS a FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1), a(1) FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1),(1) FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1), a() FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1), a(1) zzz FOR a IN(1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT INCLUDE NULLS(a FOR a IN(a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a UNPIVOT
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a UNPIVOT
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a UNPIVOT
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT EXCLUDE NULLS(a FOR a IN(a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT((a) FOR a IN(a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PARTITION BY a PATTERN(a) DEFINE a AS 1 = 1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "A" FROM a MATCH_RECOGNIZE(PARTITION BY a PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PARTITION BY a, a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ORDER BY a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ORDER BY a ASC PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ORDER BY a DESC PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ORDER BY a, a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ORDER BY a, a ASC PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(MEASURES 1 AS a PATTERN(a) DEFINE a AS 1 = 1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "A" FROM a MATCH_RECOGNIZE(MEASURES 1 AS a PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(MEASURES 1 a PATTERN(a) DEFINE a AS 1 = 1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "A" FROM a MATCH_RECOGNIZE(MEASURES 1 a PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(MEASURES 1 AS a, 1 AS a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(MEASURES 1 AS a, 1 a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ONE ROW PER MATCH PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH PATTERN(a) DEFINE a AS 1 = 1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH SHOW EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH SHOW EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH SHOW EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH SHOW EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH OMIT EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH OMIT EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH OMIT EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a MATCH_RECOGNIZE(ALL ROWS PER MATCH OMIT EMPTY MATCHES PATTERN(a) DEFINE a AS 1 = 1)
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(AFTER MATCH SKIP TO a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(AFTER MATCH SKIP PAST LAST ROW PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(AFTER MATCH SKIP TO NEXT ROW PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(AFTER MATCH SKIP TO FIRST a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(AFTER MATCH SKIP TO LAST a PATTERN(a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) SUBSET a =(a) DEFINE a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) SUBSET a =(a), a =(a) DEFINE a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) DEFINE a AS 1 = 1, a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * WINDOW a AS(PARTITION BY(1))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * WINDOW a AS(PARTITION BY 1, 1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL RIGHT JOIN a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL RIGHT JOIN a
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL RIGHT JOIN a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL RIGHT JOIN a
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL FULL JOIN a
< error
| 918 ORA-00918: : column ambiguously specified - appears in A and A
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a NATURAL LEFT OUTER JOIN a
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT OUTER JOIN a
: string
| N
: string
| B
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT OUTER JOIN a
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "B","A","D" FROM a NATURAL LEFT OUTER JOIN a
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY ROLLUP(1, 1)
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR(a) IN(a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR(a, a) IN(a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1 zzz))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(ANYSCHEMA))
< error
| 56901 ORA-56901: non-constant expression in pivot or unpivot values clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN((1)))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN((1) zzz))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1 a))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1 AS a))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1, 1))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a PIVOT(a(1) FOR a IN(1, 1 zzz))
< error
| 56902 ORA-56902: non-aggregate expressions inside the PIVOT clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR a IN((a)))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR a IN(a AS 1))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR a IN(a AS(1)))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a UNPIVOT(a FOR a IN(a, a))
< error
| 918 ORA-00918: A: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a zzz) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a) SUBSET a =(a, a) DEFINE a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY ROLLUP((zzz))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY ROLLUP((1))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * GROUP BY ROLLUP((1, 1))
< error
| 54058 ORA-54058: '*' not allowed in SELECT expressions without FROM clause
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a a) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(a a zzz) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(()) DEFINE a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(() zzz) DEFINE a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PERMUTE(a)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PERMUTE(a) zzz) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PFILE(a)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN((zzz)) DEFINE a AS 1 = 1)
< error
| 62503 ORA-62503: illegal variable definition
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN((a)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN((a zzz)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PERMUTE(a zzz)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PERMUTE(a, a)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PERMUTE(a a)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, a NUMBER, d NUMBER);
| CREATE INDEX a_d_idx ON a (d);]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM a MATCH_RECOGNIZE(PATTERN(PERMUTE(a, a zzz)) DEFINE a AS 1 = 1)
< error
| 30732 ORA-30732: table contains no user-visible columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO" FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.* FROM emp e
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO" FROM emp e
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp, dept
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.*, d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e."EMPNO",e."ENAME",e."JOB",e."MGR",e."HIREDATE",e."SAL",e."COMM",e."DEPTNO", d.dname FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT emp.ename, dept.dname FROM emp, dept WHERE emp.deptno = dept.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT emp.ename, dept.dname FROM emp, dept WHERE emp.deptno = dept.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT emp.ename, dept.dname FROM emp, dept WHERE emp.deptno = dept.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT omni_sim.emp.ename FROM omni_sim.emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT omni_sim.emp.ename FROM omni_sim.emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT emp.ename FROM emp e
< error
| 904 ORA-00904: "EMP"."ENAME": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT emp.* FROM emp e
< error
| 904 ORA-00904: "EMP": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename "Name", job "job title" FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename "Name", job "job title" FROM emp
: string
| N
: string
| Name
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename "Name", job "job title" FROM emp
: string
| N
: string
| job title
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT "ename" FROM emp
< error
| 904 ORA-00904: "ename": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename AS name, job AS title FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename AS name, job AS title FROM emp
: string
| N
: string
| NAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename AS name, job AS title FROM emp
: string
| N
: string
| TITLE
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, ename FROM emp
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename AS a, job AS a FROM emp
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, m.ename FROM emp e, emp m WHERE e.mgr = m.empno
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1 FROM dual
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal + 1 FROM emp
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT UPPER(ename) FROM emp
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT nosuch FROM emp
< error
| 904 ORA-00904: "NOSUCH": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM nosuch
< error
| 942 ORA-00942: table or view "OMNI_SIM"."NOSUCH" does not exist
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.nosuch FROM emp e
< error
| 904 ORA-00904: "E"."NOSUCH": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT x.ename FROM emp e
< error
| 904 ORA-00904: "X"."ENAME": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE nosuch = 1
< error
| 904 ORA-00904: "NOSUCH": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp ORDER BY nosuch
< error
| 904 ORA-00904: "NOSUCH": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp GROUP BY nosuch
< error
| 904 ORA-00904: "NOSUCH": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename e1 FROM emp WHERE e1 = 'X'
< error
| 904 ORA-00904: "E1": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp ORDER BY 2
< error
| 1785 ORA-01785: ORDER BY item must be the number of a SELECT-list expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM emp, dept
< error
| 918 ORA-00918: DEPTNO: column ambiguously specified - appears in DEPT and EMP
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal FROM emp e JOIN emp m ON e.mgr = m.empno
< error
| 918 ORA-00918: SAL: column ambiguously specified - appears in  and 
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp_v
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp_v
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp_v
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","DEPTNO" FROM emp_v
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","DEPTNO" FROM emp_v
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","DEPTNO" FROM emp_v
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp_v v JOIN dept d ON v.deptno = d.deptno
< error
| 957 ORA-00957: duplicate column name
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ROWNUM rn, ename FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ROWNUM rn, ename FROM emp
: string
| N
: string
| RN
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ROWNUM rn, ename FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ROWNUM, ename FROM emp
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ROWID rid, ename FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ROWID rid, ename FROM emp
: string
| N
: string
| RID
: string
| ROWID
: string
| 10
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ROWID rid, ename FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DISTINCT job FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DISTINCT job FROM emp
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT UNIQUE job FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UNIQUE job FROM emp
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| ID
: string
| NUMBER
: string
| 22
: null
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| CODE
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| TITLE
: string
| VARCHAR2
: string
| 80
: null
: null
: string
| C
: string
| 20
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| NAME
: string
| NVARCHAR2
: string
| 40
: null
: null
: string
| C
: string
| 20
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| NOTE
: string
| CLOB
: string
| 4000
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| TS
: string
| TIMESTAMP(6)
: string
| 11
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| TZ
: string
| TIMESTAMP(3) WITH TIME ZONE
: string
| 13
: null
: string
| 3
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| BD
: string
| BINARY_DOUBLE
: string
| 8
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| BF
: string
| BINARY_FLOAT
: string
| 4
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| R
: string
| RAW
: string
| 16
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| FL
: string
| FLOAT
: string
| 22
: string
| 126
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| YM
: string
| INTERVAL YEAR(2) TO MONTH
: string
| 5
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| DS
: string
| INTERVAL DAY(2) TO SECOND(6)
: string
| 11
: string
| 2
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ID","CODE","TITLE","NAME","NOTE","TS","TZ","BD","BF","R","FL","YM","DS","N","FLAG" FROM misc
: string
| N
: string
| FLAG
: string
| BOOLEAN
: string
| 1
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| ID
: string
| NUMBER
: string
| 22
: null
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| CODE
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| TITLE
: string
| VARCHAR2
: string
| 80
: null
: null
: string
| C
: string
| 20
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| NAME
: string
| NVARCHAR2
: string
| 40
: null
: null
: string
| C
: string
| 20
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| NOTE
: string
| CLOB
: string
| 4000
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| TS
: string
| TIMESTAMP(6)
: string
| 11
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| TZ
: string
| TIMESTAMP(3) WITH TIME ZONE
: string
| 13
: null
: string
| 3
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| BD
: string
| BINARY_DOUBLE
: string
| 8
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| BF
: string
| BINARY_FLOAT
: string
| 4
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| R
: string
| RAW
: string
| 16
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| FL
: string
| FLOAT
: string
| 22
: string
| 126
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| YM
: string
| INTERVAL YEAR(2) TO MONTH
: string
| 5
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| DS
: string
| INTERVAL DAY(2) TO SECOND(6)
: string
| 11
: string
| 2
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT id, code, title, name, note, ts, tz, bd, bf, r, fl, ym, ds, n, flag FROM misc
: string
| N
: string
| FLAG
: string
| BOOLEAN
: string
| 1
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE deptno = :1
< error
| 1027 ORA-01027: bind variables not allowed for data definition operations
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp SAMPLE (10)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp SAMPLE (10)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WITH READ ONLY
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WITH READ ONLY
: string
| Y
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION CONSTRAINT v_chk
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION 
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp WHERE deptno = 10 WITH CHECK OPTION 
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp FOR UPDATE
< error
| 3048 ORA-03048: SQL reserved word 'FOR' is not syntactically valid following '...v AS SELECT ename FROM emp '
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp ORDER BY sal DESC, ename
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp ORDER BY sal DESC, ename
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp ORDER BY sal DESC, ename
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname FROM emp e JOIN dept d ON e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename FROM emp e INNER JOIN dept d ON (e.deptno = d.deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename FROM emp e INNER JOIN dept d ON (e.deptno = d.deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, dname FROM emp JOIN dept USING (deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, dname FROM emp JOIN dept USING (deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, dname FROM emp JOIN dept USING (deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp JOIN dept USING (deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp JOIN dept USING (deptno)
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.deptno FROM emp e JOIN dept d USING (deptno)
< error
| 25154 ORA-25154: column part of USING clause cannot have qualifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename FROM emp e JOIN dept d USING (nosuch)
< error
| 904 ORA-00904: "D"."NOSUCH": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp NATURAL JOIN dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp NATURAL JOIN dept
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM dept NATURAL JOIN emp_v
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME" FROM dept NATURAL JOIN emp_v
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME" FROM dept NATURAL JOIN emp_v
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME" FROM dept NATURAL JOIN emp_v
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME" FROM dept NATURAL JOIN emp_v
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME" FROM dept NATURAL JOIN emp_v
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno FROM emp NATURAL JOIN emp_v
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT empno FROM emp NATURAL JOIN emp_v
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname FROM emp e LEFT OUTER JOIN dept d ON e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e LEFT OUTER JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e LEFT OUTER JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname FROM emp e RIGHT JOIN dept d ON e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e RIGHT JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e RIGHT JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname FROM emp e FULL OUTER JOIN dept d ON e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e FULL OUTER JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e FULL OUTER JOIN dept d ON e.deptno = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp FULL JOIN dept USING (deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp FULL JOIN dept USING (deptno)
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp LEFT JOIN dept USING (deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp LEFT JOIN dept USING (deptno)
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp RIGHT JOIN dept USING (deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DNAME","LOC" FROM emp RIGHT JOIN dept USING (deptno)
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM dept JOIN emp USING (deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC","EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM" FROM dept JOIN emp USING (deptno)
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT CAST(1 AS NUMBER(3)) k FROM dual) JOIN (SELECT CAST(1 AS NUMBER(5,1)) k FROM dual) USING (k)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "K" FROM (SELECT CAST(1 AS NUMBER(3)) k FROM dual) JOIN (SELECT CAST(1 AS NUMBER(5,1)) k FROM dual) USING (k)
: string
| N
: string
| K
: string
| NUMBER
: string
| 22
: string
| 5
: string
| 1
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT CAST(1 AS NUMBER(3)) k FROM dual) LEFT JOIN (SELECT CAST(1 AS NUMBER(5,1)) k FROM dual) USING (k)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "K" FROM (SELECT CAST(1 AS NUMBER(3)) k FROM dual) LEFT JOIN (SELECT CAST(1 AS NUMBER(5,1)) k FROM dual) USING (k)
: string
| N
: string
| K
: string
| NUMBER
: string
| 22
: string
| 3
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, m.ename AS mgr_name FROM emp e LEFT JOIN emp m ON e.mgr = m.empno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, m.ename AS mgr_name FROM emp e LEFT JOIN emp m ON e.mgr = m.empno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, m.ename AS mgr_name FROM emp e LEFT JOIN emp m ON e.mgr = m.empno
: string
| N
: string
| MGR_NAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.deptno FROM emp e CROSS JOIN dept d
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.deptno FROM emp e CROSS JOIN dept d
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname, m.job FROM emp e JOIN dept d ON e.deptno = d.deptno JOIN emp m ON m.empno = e.mgr
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname, m.job FROM emp e JOIN dept d ON e.deptno = d.deptno JOIN emp m ON m.empno = e.mgr
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname, m.job FROM emp e JOIN dept d ON e.deptno = d.deptno JOIN emp m ON m.empno = e.mgr
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname, m.job FROM emp e JOIN dept d ON e.deptno = d.deptno JOIN emp m ON m.empno = e.mgr
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp e JOIN dept d ON e.deptno = d.deptno AND d.loc = 'DALLAS'
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp e JOIN dept d ON e.deptno = d.deptno AND d.loc = 'DALLAS'
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.loc FROM emp e, dept d WHERE e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.loc FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.loc FROM emp e, dept d WHERE e.deptno = d.deptno
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname FROM emp e, dept d WHERE e.deptno = d.deptno(+)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e, dept d WHERE e.deptno = d.deptno(+)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname FROM emp e, dept d WHERE e.deptno = d.deptno(+)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename, d.dname, d.loc FROM emp e, dept d WHERE e.deptno(+) = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname, d.loc FROM emp e, dept d WHERE e.deptno(+) = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname, d.loc FROM emp e, dept d WHERE e.deptno(+) = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT e.ename, d.dname, d.loc FROM emp e, dept d WHERE e.deptno(+) = d.deptno
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp e, dept d WHERE e.deptno = d.deptno(+) AND d.loc(+) = 'X'
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp e, dept d WHERE e.deptno = d.deptno(+) AND d.loc(+) = 'X'
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename FROM emp e, dept d WHERE e.deptno(+) = d.deptno(+)
< error
| 1468 ORA-01468: a predicate may reference only one outer-joined table
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT e.ename FROM emp e LEFT JOIN dept d ON e.deptno = d.deptno(+)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT e.ename FROM emp e LEFT JOIN dept d ON e.deptno = d.deptno(+)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT d.dname, e.ename FROM dept d CROSS APPLY (SELECT ename FROM emp WHERE emp.deptno = d.deptno) e
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, e.ename FROM dept d CROSS APPLY (SELECT ename FROM emp WHERE emp.deptno = d.deptno) e
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, e.ename FROM dept d CROSS APPLY (SELECT ename FROM emp WHERE emp.deptno = d.deptno) e
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT d.dname, x.cnt FROM dept d OUTER APPLY (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, x.cnt FROM dept d OUTER APPLY (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, x.cnt FROM dept d OUTER APPLY (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x
: string
| N
: string
| CNT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT d.dname, x.cnt FROM dept d, LATERAL (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, x.cnt FROM dept d, LATERAL (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, x.cnt FROM dept d, LATERAL (SELECT COUNT(*) cnt FROM emp e WHERE e.deptno = d.deptno) x
: string
| N
: string
| CNT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","DEPTNO","DNAME" FROM emp e, LATERAL (SELECT dname FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT d.dname, COUNT(e.empno) cnt FROM dept d LEFT JOIN emp e ON e.deptno = d.deptno GROUP BY d.dname
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, COUNT(e.empno) cnt FROM dept d LEFT JOIN emp e ON e.deptno = d.deptno GROUP BY d.dname
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT d.dname, COUNT(e.empno) cnt FROM dept d LEFT JOIN emp e ON e.deptno = d.deptno GROUP BY d.dname
: string
| N
: string
| CNT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT x.ename FROM (SELECT ename, sal FROM emp WHERE sal > 1000) x
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT x.ename FROM (SELECT ename, sal FROM emp WHERE sal > 1000) x
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT ename, sal * 2 FROM emp)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ENAME","SAL*2" FROM (SELECT ename, sal * 2 FROM emp)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "ENAME","SAL*2" FROM (SELECT ename, sal * 2 FROM emp)
: string
| N
: string
| SAL*2
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT ename, sal * 2 dbl FROM emp)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "ENAME","DBL" FROM (SELECT ename, sal * 2 dbl FROM emp)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "ENAME","DBL" FROM (SELECT ename, sal * 2 dbl FROM emp)
: string
| N
: string
| DBL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT * FROM dept)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC" FROM (SELECT * FROM dept)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC" FROM (SELECT * FROM dept)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC" FROM (SELECT * FROM dept)
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT x.* FROM (SELECT d.* FROM dept d) x
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT x."DEPTNO",x."DNAME",x."LOC" FROM (SELECT d.* FROM dept d) x
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT x."DEPTNO",x."DNAME",x."LOC" FROM (SELECT d.* FROM dept d) x
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT x."DEPTNO",x."DNAME",x."LOC" FROM (SELECT d.* FROM dept d) x
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT x FROM (SELECT 1 x FROM dual UNION ALL SELECT 2 FROM dual)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT x FROM (SELECT 1 x FROM dual UNION ALL SELECT 2 FROM dual)
: string
| N
: string
| X
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE deptno IN (SELECT deptno FROM dept WHERE loc = 'DALLAS')
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE deptno IN (SELECT deptno FROM dept WHERE loc = 'DALLAS')
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp e WHERE EXISTS (SELECT 1 FROM dept d WHERE d.deptno = e.deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp e WHERE EXISTS (SELECT 1 FROM dept d WHERE d.deptno = e.deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT d.deptno, d.dname FROM dept d WHERE NOT EXISTS (SELECT NULL FROM emp e WHERE e.deptno = d.deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT d.deptno, d.dname FROM dept d WHERE NOT EXISTS (SELECT NULL FROM emp e WHERE e.deptno = d.deptno)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT d.deptno, d.dname FROM dept d WHERE NOT EXISTS (SELECT NULL FROM emp e WHERE e.deptno = d.deptno)
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, (SELECT dname FROM dept d WHERE d.deptno = e.deptno) dname FROM emp e
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, (SELECT dname FROM dept d WHERE d.deptno = e.deptno) dname FROM emp e
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, (SELECT dname FROM dept d WHERE d.deptno = e.deptno) dname FROM emp e
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, (SELECT dname, loc FROM dept d WHERE d.deptno = e.deptno) x FROM emp e
< error
| 913 ORA-00913: too many values
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE sal > ALL (SELECT sal FROM emp WHERE deptno = 30)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE sal > ALL (SELECT sal FROM emp WHERE deptno = 30)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE deptno = ANY (10, 20)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE deptno = ANY (10, 20)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE (deptno, job) IN (SELECT deptno, job FROM emp WHERE ename = 'SMITH')
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE (deptno, job) IN (SELECT deptno, job FROM emp WHERE ename = 'SMITH')
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp e WHERE sal > (SELECT AVG(sal) FROM emp WHERE deptno = e.deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp e WHERE sal > (SELECT AVG(sal) FROM emp WHERE deptno = e.deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp e WHERE sal > (SELECT AVG(sal) FROM emp WHERE deptno = e.deptno)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno, ename FROM emp WHERE empno = (SELECT MAX(empno) FROM emp)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT empno, ename FROM emp WHERE empno = (SELECT MAX(empno) FROM emp)
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT empno, ename FROM emp WHERE empno = (SELECT MAX(empno) FROM emp)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE empno = (SELECT empno, ename FROM emp)
< error
| 913 ORA-00913: too many values
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE ename LIKE 'A%' AND sal BETWEEN 1000 AND 3000 AND deptno IN (10, 20) AND comm IS NOT NULL
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE ename LIKE 'A%' AND sal BETWEEN 1000 AND 3000 AND deptno IN (10, 20) AND comm IS NOT NULL
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT d.dname, CURSOR(SELECT ename FROM emp e WHERE e.deptno = d.deptno) c FROM dept d
< error
| 22902 ORA-22902: CURSOR expression is not allowed in this context.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH d AS (SELECT deptno, dname FROM dept) SELECT * FROM d
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH d AS (SELECT deptno, dname FROM dept) SELECT "DEPTNO","DNAME" FROM d
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| WITH d AS (SELECT deptno, dname FROM dept) SELECT "DEPTNO","DNAME" FROM d
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH d (no, name) AS (SELECT deptno, dname FROM dept) SELECT * FROM d
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH d (no, name) AS (SELECT deptno, dname FROM dept) SELECT "NO","NAME" FROM d
: string
| N
: string
| NO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| WITH d (no, name) AS (SELECT deptno, dname FROM dept) SELECT "NO","NAME" FROM d
: string
| N
: string
| NAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH d (no) AS (SELECT deptno, dname FROM dept) SELECT * FROM d
< error
| 32038 ORA-32038: number of WITH clause column names does not match number of elements in select list
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH a AS (SELECT empno, mgr FROM emp), b AS (SELECT empno FROM a WHERE mgr IS NULL) SELECT * FROM b
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH a AS (SELECT empno, mgr FROM emp), b AS (SELECT empno FROM a WHERE mgr IS NULL) SELECT "EMPNO" FROM b
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH d AS (SELECT * FROM dept) SELECT d.dname, e.ename FROM d JOIN emp e ON e.deptno = d.deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH d AS (SELECT * FROM dept) SELECT d.dname, e.ename FROM d JOIN emp e ON e.deptno = d.deptno
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| N
: string
| NO
< row
: string
| V
: string
| WITH d AS (SELECT * FROM dept) SELECT d.dname, e.ename FROM d JOIN emp e ON e.deptno = d.deptno
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH r (empno, ename, lvl) AS (SELECT empno, ename, 1 FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename, r.lvl + 1 FROM emp e JOIN r ON e.mgr = r.empno) SELECT * FROM r
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (empno, ename, lvl) AS (SELECT empno, ename, 1 FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename, r.lvl + 1 FROM emp e JOIN r ON e.mgr = r.empno) SELECT "EMPNO","ENAME","LVL" FROM r
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (empno, ename, lvl) AS (SELECT empno, ename, 1 FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename, r.lvl + 1 FROM emp e JOIN r ON e.mgr = r.empno) SELECT "EMPNO","ENAME","LVL" FROM r
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (empno, ename, lvl) AS (SELECT empno, ename, 1 FROM emp WHERE mgr IS NULL UNION ALL SELECT e.empno, e.ename, r.lvl + 1 FROM emp e JOIN r ON e.mgr = r.empno) SELECT "EMPNO","ENAME","LVL" FROM r
: string
| N
: string
| LVL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH r AS (SELECT 1 n FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r
< error
| 32039 ORA-32039: missing column alias list in recursive WITH clause element R
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SEARCH DEPTH FIRST BY n SET ord CYCLE n SET is_cycle TO 'Y' DEFAULT 'N' SELECT * FROM r
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SEARCH DEPTH FIRST BY n SET ord CYCLE n SET is_cycle TO 'Y' DEFAULT 'N' SELECT "N","ORD","IS_CYCLE" FROM r
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SEARCH DEPTH FIRST BY n SET ord CYCLE n SET is_cycle TO 'Y' DEFAULT 'N' SELECT "N","ORD","IS_CYCLE" FROM r
: string
| N
: string
| ORD
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH r (n) AS (SELECT 1 FROM dual UNION ALL SELECT n + 1 FROM r WHERE n < 5) SEARCH DEPTH FIRST BY n SET ord CYCLE n SET is_cycle TO 'Y' DEFAULT 'N' SELECT "N","ORD","IS_CYCLE" FROM r
: string
| N
: string
| IS_CYCLE
: string
| VARCHAR2
: string
| 1
: null
: null
: string
| B
: string
| 1
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH d AS (SELECT deptno FROM dept) SELECT ename FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH d AS (SELECT deptno FROM dept) SELECT ename FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH d AS (SELECT deptno FROM dept) SELECT dname FROM d
< error
| 904 ORA-00904: "DNAME": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS WITH FUNCTION dbl(x NUMBER) RETURN NUMBER IS BEGIN RETURN x * 2; END; SELECT dbl(sal) d FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| WITH FUNCTION dbl(x NUMBER) RETURN NUMBER IS BEGIN RETURN x * 2; END; SELECT dbl(sal) d FROM emp
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno, ename FROM emp UNION SELECT deptno, dname FROM dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT empno, ename FROM emp UNION SELECT deptno, dname FROM dept
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT empno, ename FROM emp UNION SELECT deptno, dname FROM dept
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno FROM emp UNION ALL SELECT deptno FROM dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT empno FROM emp UNION ALL SELECT deptno FROM dept
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM dept MINUS SELECT deptno FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno FROM dept MINUS SELECT deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM dept INTERSECT SELECT deptno FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno FROM dept INTERSECT SELECT deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM dept EXCEPT SELECT deptno FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno FROM dept EXCEPT SELECT deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM dept MINUS ALL SELECT deptno FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno FROM dept MINUS ALL SELECT deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM dept INTERSECT ALL SELECT deptno FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno FROM dept INTERSECT ALL SELECT deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno, ename FROM emp UNION SELECT deptno FROM dept
< error
| 1789 ORA-01789: query block has incorrect number of result columns
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno FROM emp UNION SELECT dname FROM dept
< error
| 1790 ORA-01790: expression must have same datatype as corresponding expression
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp UNION SELECT dname FROM dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp UNION SELECT dname FROM dept
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT job FROM emp UNION SELECT loc FROM dept ORDER BY 1
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT job FROM emp UNION SELECT loc FROM dept ORDER BY 1
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp UNION SELECT dname FROM dept ORDER BY ename
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp UNION SELECT dname FROM dept ORDER BY ename
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp UNION SELECT dname FROM dept ORDER BY dname
< error
| 904 ORA-00904: "DNAME": invalid identifier
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS (SELECT ename FROM emp) UNION ALL (SELECT dname FROM dept)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| (SELECT ename FROM emp) UNION ALL (SELECT dname FROM dept)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal FROM emp UNION ALL SELECT comm FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal FROM emp UNION ALL SELECT comm FROM emp
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal FROM emp UNION SELECT deptno FROM dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal FROM emp UNION SELECT deptno FROM dept
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT code FROM misc UNION SELECT title FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT code FROM misc UNION SELECT title FROM misc
: string
| N
: string
| CODE
: string
| VARCHAR2
: string
| 80
: null
: null
: string
| C
: string
| 20
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NULL x FROM dual UNION SELECT ename FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NULL x FROM dual UNION SELECT ename FROM emp
: string
| N
: string
| X
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT hiredate FROM emp UNION SELECT ts FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT hiredate FROM emp UNION SELECT ts FROM misc
: string
| N
: string
| HIREDATE
: string
| TIMESTAMP(6)
: string
| 11
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename n FROM emp UNION SELECT dname FROM dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename n FROM emp UNION SELECT dname FROM dept
: string
| N
: string
| N
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 'a' c FROM dual UNION SELECT 'bcd' FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'a' c FROM dual UNION SELECT 'bcd' FROM dual
: string
| N
: string
| C
: string
| VARCHAR2
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 'a' c FROM dual UNION ALL SELECT ename FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'a' c FROM dual UNION ALL SELECT ename FROM emp
: string
| N
: string
| C
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1 n FROM dual UNION SELECT 2.5 FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 1 n FROM dual UNION SELECT 2.5 FROM dual
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM dept UNION SELECT * FROM dept
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC" FROM dept UNION SELECT "DEPTNO","DNAME","LOC" FROM dept
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC" FROM dept UNION SELECT "DEPTNO","DNAME","LOC" FROM dept
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 14
: null
: null
: string
| B
: string
| 14
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","DNAME","LOC" FROM dept UNION SELECT "DEPTNO","DNAME","LOC" FROM dept
: string
| N
: string
| LOC
: string
| VARCHAR2
: string
| 13
: null
: null
: string
| B
: string
| 13
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT code FROM misc UNION SELECT code FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT code FROM misc UNION SELECT code FROM misc
: string
| N
: string
| CODE
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT note FROM misc UNION SELECT note FROM misc
< error
| 22848 ORA-22848: cannot use CLOB type as comparison key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT note FROM misc UNION ALL SELECT note FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT note FROM misc UNION ALL SELECT note FROM misc
: string
| N
: string
| NOTE
: string
| CLOB
: string
| 4000
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
: string
| N
: string
| CNT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
: string
| N
: string
| TOTAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
: string
| N
: string
| AVG_SAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
: string
| N
: string
| FIRST_HIRE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt, SUM(sal) total, AVG(sal) avg_sal, MIN(hiredate) first_hire, MAX(ename) last_name FROM emp GROUP BY deptno
: string
| N
: string
| LAST_NAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT MAX(sal) mx, MIN(empno) mn, MAX(code) mc FROM misc, emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT MAX(sal) mx, MIN(empno) mn, MAX(code) mc FROM misc, emp
: string
| N
: string
| MX
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT MAX(sal) mx, MIN(empno) mn, MAX(code) mc FROM misc, emp
: string
| N
: string
| MN
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT MAX(sal) mx, MIN(empno) mn, MAX(code) mc FROM misc, emp
: string
| N
: string
| MC
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, COUNT(*) cnt FROM emp GROUP BY deptno HAVING COUNT(*) > 3
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt FROM emp GROUP BY deptno HAVING COUNT(*) > 3
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, COUNT(*) cnt FROM emp GROUP BY deptno HAVING COUNT(*) > 3
: string
| N
: string
| CNT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, ename, COUNT(*) cnt FROM emp GROUP BY deptno
< error
| 979 ORA-00979: "ENAME": must appear in the GROUP BY clause or be used in an aggregate function
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, COUNT(*) cnt FROM emp
< error
| 937 ORA-00937: not a single-group group function
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno FROM emp WHERE SUM(sal) > 10 GROUP BY deptno
< error
| 934 ORA-00934: group function is not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, job, SUM(sal) s FROM emp GROUP BY ROLLUP (deptno, job)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s FROM emp GROUP BY ROLLUP (deptno, job)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s FROM emp GROUP BY ROLLUP (deptno, job)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s FROM emp GROUP BY ROLLUP (deptno, job)
: string
| N
: string
| S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, job, SUM(sal) s, GROUPING(job) g FROM emp GROUP BY CUBE (deptno, job)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s, GROUPING(job) g FROM emp GROUP BY CUBE (deptno, job)
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s, GROUPING(job) g FROM emp GROUP BY CUBE (deptno, job)
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s, GROUPING(job) g FROM emp GROUP BY CUBE (deptno, job)
: string
| N
: string
| S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s, GROUPING(job) g FROM emp GROUP BY CUBE (deptno, job)
: string
| N
: string
| G
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, job, SUM(sal) s FROM emp GROUP BY GROUPING SETS ((deptno), (job), ())
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s FROM emp GROUP BY GROUPING SETS ((deptno), (job), ())
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s FROM emp GROUP BY GROUPING SETS ((deptno), (job), ())
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, job, SUM(sal) s FROM emp GROUP BY GROUPING SETS ((deptno), (job), ())
: string
| N
: string
| S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT COUNT(DISTINCT job) jobs FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT COUNT(DISTINCT job) jobs FROM emp
: string
| N
: string
| JOBS
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, STDDEV(sal) sd, VARIANCE(sal) v, MEDIAN(sal) med FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, STDDEV(sal) sd, VARIANCE(sal) v, MEDIAN(sal) med FROM emp GROUP BY deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, STDDEV(sal) sd, VARIANCE(sal) v, MEDIAN(sal) med FROM emp GROUP BY deptno
: string
| N
: string
| SD
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, STDDEV(sal) sd, VARIANCE(sal) v, MEDIAN(sal) med FROM emp GROUP BY deptno
: string
| N
: string
| V
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, STDDEV(sal) sd, VARIANCE(sal) v, MEDIAN(sal) med FROM emp GROUP BY deptno
: string
| N
: string
| MED
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY sal) p50 FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY sal) p50 FROM emp GROUP BY deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY sal) p50 FROM emp GROUP BY deptno
: string
| N
: string
| P50
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 1 + COUNT(*) c FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 1 + COUNT(*) c FROM emp
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno d, COUNT(*) c FROM emp GROUP BY d
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno d, COUNT(*) c FROM emp GROUP BY d
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno d, COUNT(*) c FROM emp GROUP BY d
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, COUNT(*) c FROM emp GROUP BY 1
< error
| 3162 ORA-03162: "DEPTNO": must appear in the GROUP BY clause or be used in an aggregate function as 'group_by_position_enabled' is FALSE
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT MAX(AVG(sal)) m FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT MAX(AVG(sal)) m FROM emp GROUP BY deptno
: string
| N
: string
| M
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT MAX(AVG(sal)) m FROM emp
< error
| 978 ORA-00978: nested group function without GROUP BY
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, MAX(AVG(sal)) m FROM emp GROUP BY deptno
< error
| 937 ORA-00937: not a single-group group function
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, LISTAGG(ename, ',') WITHIN GROUP (ORDER BY ename) names FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, LISTAGG(ename, ',') WITHIN GROUP (ORDER BY ename) names FROM emp GROUP BY deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, LISTAGG(ename, ',') WITHIN GROUP (ORDER BY ename) names FROM emp GROUP BY deptno
: string
| N
: string
| NAMES
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, MAX(sal) KEEP (DENSE_RANK LAST ORDER BY hiredate) last_sal FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, MAX(sal) KEEP (DENSE_RANK LAST ORDER BY hiredate) last_sal FROM emp GROUP BY deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, MAX(sal) KEEP (DENSE_RANK LAST ORDER BY hiredate) last_sal FROM emp GROUP BY deptno
: string
| N
: string
| LAST_SAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT SUM(bd) s, AVG(bf) a, MAX(ts) m, MIN(code) mc, COUNT(note) cn FROM misc
< error
| 22849 ORA-22849: Type CLOB is not supported for this function or operator.
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT COUNT(*) c FROM misc GROUP BY note
< error
| 22848 ORA-22848: cannot use CLOB type as comparison key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DISTINCT note FROM misc
< error
| 22848 ORA-22848: cannot use CLOB type as comparison key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT note FROM misc ORDER BY note
< error
| 22848 ORA-22848: cannot use CLOB type as comparison key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, ROW_NUMBER() OVER (ORDER BY sal DESC) rn FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, ROW_NUMBER() OVER (ORDER BY sal DESC) rn FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, ROW_NUMBER() OVER (ORDER BY sal DESC) rn FROM emp
: string
| N
: string
| RN
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, RANK() OVER (PARTITION BY deptno ORDER BY sal) rnk, DENSE_RANK() OVER (PARTITION BY deptno ORDER BY sal) drnk FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, RANK() OVER (PARTITION BY deptno ORDER BY sal) rnk, DENSE_RANK() OVER (PARTITION BY deptno ORDER BY sal) drnk FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, RANK() OVER (PARTITION BY deptno ORDER BY sal) rnk, DENSE_RANK() OVER (PARTITION BY deptno ORDER BY sal) drnk FROM emp
: string
| N
: string
| RNK
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, RANK() OVER (PARTITION BY deptno ORDER BY sal) rnk, DENSE_RANK() OVER (PARTITION BY deptno ORDER BY sal) drnk FROM emp
: string
| N
: string
| DRNK
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, SUM(sal) OVER (PARTITION BY deptno) dept_sal FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, SUM(sal) OVER (PARTITION BY deptno) dept_sal FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, SUM(sal) OVER (PARTITION BY deptno) dept_sal FROM emp
: string
| N
: string
| DEPT_SAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, AVG(sal) OVER (ORDER BY hiredate ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) mov FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, AVG(sal) OVER (ORDER BY hiredate ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) mov FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, AVG(sal) OVER (ORDER BY hiredate ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) mov FROM emp
: string
| N
: string
| MOV
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, LAG(sal) OVER (ORDER BY hiredate) prev, LEAD(sal, 1, 0) OVER (ORDER BY hiredate) nxt FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, LAG(sal) OVER (ORDER BY hiredate) prev, LEAD(sal, 1, 0) OVER (ORDER BY hiredate) nxt FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, LAG(sal) OVER (ORDER BY hiredate) prev, LEAD(sal, 1, 0) OVER (ORDER BY hiredate) nxt FROM emp
: string
| N
: string
| PREV
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, LAG(sal) OVER (ORDER BY hiredate) prev, LEAD(sal, 1, 0) OVER (ORDER BY hiredate) nxt FROM emp
: string
| N
: string
| NXT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, FIRST_VALUE(ename) OVER (PARTITION BY deptno ORDER BY sal) lowest FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, FIRST_VALUE(ename) OVER (PARTITION BY deptno ORDER BY sal) lowest FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, FIRST_VALUE(ename) OVER (PARTITION BY deptno ORDER BY sal) lowest FROM emp
: string
| N
: string
| LOWEST
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, NTILE(4) OVER (ORDER BY sal) q FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, NTILE(4) OVER (ORDER BY sal) q FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, NTILE(4) OVER (ORDER BY sal) q FROM emp
: string
| N
: string
| Q
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, RATIO_TO_REPORT(sal) OVER () pct FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, RATIO_TO_REPORT(sal) OVER () pct FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, RATIO_TO_REPORT(sal) OVER () pct FROM emp
: string
| N
: string
| PCT
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, COUNT(*) OVER () total FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, COUNT(*) OVER () total FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, COUNT(*) OVER () total FROM emp
: string
| N
: string
| TOTAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, PERCENT_RANK() OVER (ORDER BY sal) pr, CUME_DIST() OVER (ORDER BY sal) cd FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, PERCENT_RANK() OVER (ORDER BY sal) pr, CUME_DIST() OVER (ORDER BY sal) cd FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, PERCENT_RANK() OVER (ORDER BY sal) pr, CUME_DIST() OVER (ORDER BY sal) cd FROM emp
: string
| N
: string
| PR
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, PERCENT_RANK() OVER (ORDER BY sal) pr, CUME_DIST() OVER (ORDER BY sal) cd FROM emp
: string
| N
: string
| CD
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp WHERE ROW_NUMBER() OVER (ORDER BY sal) = 1
< error
| 30483 ORA-30483: window  functions are not allowed here
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, SUM(sal) OVER w s FROM emp WINDOW w AS (PARTITION BY deptno)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, SUM(sal) OVER w s FROM emp WINDOW w AS (PARTITION BY deptno)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, SUM(sal) OVER w s FROM emp WINDOW w AS (PARTITION BY deptno)
: string
| N
: string
| S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal, sal - AVG(sal) OVER (PARTITION BY deptno) diff FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal, sal - AVG(sal) OVER (PARTITION BY deptno) diff FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal, sal - AVG(sal) OVER (PARTITION BY deptno) diff FROM emp
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal, sal - AVG(sal) OVER (PARTITION BY deptno) diff FROM emp
: string
| N
: string
| DIFF
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT deptno, SUM(SUM(sal)) OVER () total FROM emp GROUP BY deptno
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, SUM(SUM(sal)) OVER () total FROM emp GROUP BY deptno
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT deptno, SUM(SUM(sal)) OVER () total FROM emp GROUP BY deptno
: string
| N
: string
| TOTAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ROW_NUMBER() OVER (ORDER BY note) r FROM misc
< error
| 22848 ORA-22848: cannot use CLOB type as comparison key
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT LEVEL lvl, ename FROM emp CONNECT BY PRIOR empno = mgr START WITH mgr IS NULL
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LEVEL lvl, ename FROM emp CONNECT BY PRIOR empno = mgr START WITH mgr IS NULL
: string
| N
: string
| LVL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LEVEL lvl, ename FROM emp CONNECT BY PRIOR empno = mgr START WITH mgr IS NULL
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') path, CONNECT_BY_ISLEAF leaf, CONNECT_BY_ROOT ename root FROM emp START WITH mgr IS NULL CONNECT BY NOCYCLE PRIOR empno = mgr
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') path, CONNECT_BY_ISLEAF leaf, CONNECT_BY_ROOT ename root FROM emp START WITH mgr IS NULL CONNECT BY NOCYCLE PRIOR empno = mgr
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') path, CONNECT_BY_ISLEAF leaf, CONNECT_BY_ROOT ename root FROM emp START WITH mgr IS NULL CONNECT BY NOCYCLE PRIOR empno = mgr
: string
| N
: string
| PATH
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') path, CONNECT_BY_ISLEAF leaf, CONNECT_BY_ROOT ename root FROM emp START WITH mgr IS NULL CONNECT BY NOCYCLE PRIOR empno = mgr
: string
| N
: string
| LEAF
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, SYS_CONNECT_BY_PATH(ename, '/') path, CONNECT_BY_ISLEAF leaf, CONNECT_BY_ROOT ename root FROM emp START WITH mgr IS NULL CONNECT BY NOCYCLE PRIOR empno = mgr
: string
| N
: string
| ROOT
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, LEVEL FROM emp CONNECT BY PRIOR empno = mgr ORDER SIBLINGS BY ename
< error
| 998 ORA-00998: must name this expression with a column alias
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT LEVEL l FROM dual CONNECT BY LEVEL <= 10
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LEVEL l FROM dual CONNECT BY LEVEL <= 10
: string
| N
: string
| L
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp START WITH mgr IS NULL CONNECT BY PRIOR empno = mgr
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT PRIOR ename p, ename FROM emp CONNECT BY PRIOR empno = mgr
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT PRIOR ename p, ename FROM emp CONNECT BY PRIOR empno = mgr
: string
| N
: string
| P
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT PRIOR ename p, ename FROM emp CONNECT BY PRIOR empno = mgr
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp ORDER BY sal DESC FETCH FIRST 5 ROWS ONLY
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp ORDER BY sal DESC FETCH FIRST 5 ROWS ONLY
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp ORDER BY sal DESC FETCH FIRST 5 ROWS ONLY
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp ORDER BY sal OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp ORDER BY sal OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp ORDER BY sal FETCH FIRST 10 PERCENT ROWS WITH TIES
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp ORDER BY sal FETCH FIRST 10 PERCENT ROWS WITH TIES
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE ROWNUM <= 5
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE ROWNUM <= 5
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE ROWNUM < 3 ORDER BY ename
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp WHERE ROWNUM < 3 ORDER BY ename
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR deptno IN (10 AS d10, 20 AS d20, 30 AS d30))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","D10","D20","D30" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR deptno IN (10 AS d10, 20 AS d20, 30 AS d30))
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","D10","D20","D30" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR deptno IN (10 AS d10, 20 AS d20, 30 AS d30))
: string
| N
: string
| D10
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","D10","D20","D30" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR deptno IN (10 AS d10, 20 AS d20, 30 AS d30))
: string
| N
: string
| D20
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","D10","D20","D30" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR deptno IN (10 AS d10, 20 AS d20, 30 AS d30))
: string
| N
: string
| D30
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","10_S","10_C","20_S","20_C" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","10_S","10_C","20_S","20_C" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))
: string
| N
: string
| 10_S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","10_S","10_C","20_S","20_C" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))
: string
| N
: string
| 10_C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","10_S","10_C","20_S","20_C" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))
: string
| N
: string
| 20_S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "JOB","10_S","10_C","20_S","20_C" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) s, COUNT(*) c FOR deptno IN (10, 20))
: string
| N
: string
| 20_C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT job, deptno, sal FROM emp) PIVOT (MAX(sal) FOR (job, deptno) IN (('CLERK', 10) AS c10, ('MANAGER', 20) AS m20))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "C10","M20" FROM (SELECT job, deptno, sal FROM emp) PIVOT (MAX(sal) FOR (job, deptno) IN (('CLERK', 10) AS c10, ('MANAGER', 20) AS m20))
: string
| N
: string
| C10
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "C10","M20" FROM (SELECT job, deptno, sal FROM emp) PIVOT (MAX(sal) FOR (job, deptno) IN (('CLERK', 10) AS c10, ('MANAGER', 20) AS m20))
: string
| N
: string
| M20
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR job IN ('CLERK', 'MANAGER'))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","'CLERK'","'MANAGER'" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR job IN ('CLERK', 'MANAGER'))
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","'CLERK'","'MANAGER'" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR job IN ('CLERK', 'MANAGER'))
: string
| N
: string
| 'CLERK'
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "DEPTNO","'CLERK'","'MANAGER'" FROM (SELECT deptno, job, sal FROM emp) PIVOT (SUM(sal) FOR job IN ('CLERK', 'MANAGER'))
: string
| N
: string
| 'MANAGER'
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| JOB
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| MGR
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| HIREDATE
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| COMM
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| D10
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","ENAME","JOB","MGR","HIREDATE","SAL","COMM","D10","D20" FROM emp PIVOT (COUNT(*) FOR deptno IN (10 AS d10, 20 AS d20))
: string
| N
: string
| D20
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT empno, sal, comm FROM emp) UNPIVOT (amount FOR kind IN (sal AS 'SAL', comm AS 'COMM'))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","KIND","AMOUNT" FROM (SELECT empno, sal, comm FROM emp) UNPIVOT (amount FOR kind IN (sal AS 'SAL', comm AS 'COMM'))
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","KIND","AMOUNT" FROM (SELECT empno, sal, comm FROM emp) UNPIVOT (amount FOR kind IN (sal AS 'SAL', comm AS 'COMM'))
: string
| N
: string
| KIND
: string
| VARCHAR2
: string
| 4
: null
: null
: string
| B
: string
| 4
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","KIND","AMOUNT" FROM (SELECT empno, sal, comm FROM emp) UNPIVOT (amount FOR kind IN (sal AS 'SAL', comm AS 'COMM'))
: string
| N
: string
| AMOUNT
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT empno, sal, comm FROM emp) UNPIVOT INCLUDE NULLS (amount FOR kind IN (sal, comm))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","KIND","AMOUNT" FROM (SELECT empno, sal, comm FROM emp) UNPIVOT INCLUDE NULLS (amount FOR kind IN (sal, comm))
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","KIND","AMOUNT" FROM (SELECT empno, sal, comm FROM emp) UNPIVOT INCLUDE NULLS (amount FOR kind IN (sal, comm))
: string
| N
: string
| KIND
: string
| VARCHAR2
: string
| 4
: null
: null
: string
| B
: string
| 4
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "EMPNO","KIND","AMOUNT" FROM (SELECT empno, sal, comm FROM emp) UNPIVOT INCLUDE NULLS (amount FOR kind IN (sal, comm))
: string
| N
: string
| AMOUNT
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CASE WHEN sal > 2000 THEN 'high' ELSE 'low' END band FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CASE WHEN sal > 2000 THEN 'high' ELSE 'low' END band FROM emp
: string
| N
: string
| BAND
: string
| VARCHAR2
: string
| 4
: null
: null
: string
| B
: string
| 4
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CASE deptno WHEN 10 THEN 'ACCOUNTING' WHEN 20 THEN 'RESEARCH' END dname FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CASE deptno WHEN 10 THEN 'ACCOUNTING' WHEN 20 THEN 'RESEARCH' END dname FROM emp
: string
| N
: string
| DNAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CASE WHEN sal > 2000 THEN sal ELSE comm END x FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CASE WHEN sal > 2000 THEN sal ELSE comm END x FROM emp
: string
| N
: string
| X
: string
| NUMBER
: string
| 22
: string
| 7
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CASE WHEN sal > 2000 THEN 'x' ELSE 1 END x FROM emp
< error
| 932 ORA-00932: expression (1) is of data type NUMBER, which is incompatible with expected data type CHAR
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CASE WHEN id > 1 THEN code ELSE title END x FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CASE WHEN id > 1 THEN code ELSE title END x FROM misc
: string
| N
: string
| X
: string
| VARCHAR2
: string
| 80
: null
: null
: string
| B
: string
| 80
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CASE WHEN id > 1 THEN ts ELSE NULL END x FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CASE WHEN id > 1 THEN ts ELSE NULL END x FROM misc
: string
| N
: string
| X
: string
| TIMESTAMP(6)
: string
| 11
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DECODE(deptno, 10, 'A', 20, 'BB', 'CCC') d FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DECODE(deptno, 10, 'A', 20, 'BB', 'CCC') d FROM emp
: string
| N
: string
| D
: string
| VARCHAR2
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DECODE(deptno, 10, sal, comm) d FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DECODE(deptno, 10, sal, comm) d FROM emp
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DECODE(job, 'CLERK', 1, 2) d FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DECODE(job, 'CLERK', 1, 2) d FROM emp
: string
| N
: string
| D
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DECODE(id, 1, NULL, code) x FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DECODE(id, 1, NULL, code) x FROM misc
: string
| N
: string
| X
: string
| VARCHAR2
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NVL(comm, 0) c FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NVL(comm, 0) c FROM emp
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NVL(job, 'none') j FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NVL(job, 'none') j FROM emp
: string
| N
: string
| J
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NVL(NULL, 1) x FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NVL(NULL, 1) x FROM dual
: string
| N
: string
| X
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NVL(title, name) x FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NVL(title, name) x FROM misc
: string
| N
: string
| X
: string
| VARCHAR2
: string
| 320
: null
: null
: string
| C
: string
| 80
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NVL2(comm, 'Y', 'N') has_comm FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NVL2(comm, 'Y', 'N') has_comm FROM emp
: string
| N
: string
| HAS_COMM
: string
| VARCHAR2
: string
| 1
: null
: null
: string
| B
: string
| 1
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT COALESCE(comm, sal, 0) c FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT COALESCE(comm, sal, 0) c FROM emp
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT COALESCE(NULL, code) x FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT COALESCE(NULL, code) x FROM misc
: string
| N
: string
| X
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT NULLIF(sal, comm) n FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT NULLIF(sal, comm) n FROM emp
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT GREATEST(sal, comm) g, LEAST(ename, job) l FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT GREATEST(sal, comm) g, LEAST(ename, job) l FROM emp
: string
| N
: string
| G
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT GREATEST(sal, comm) g, LEAST(ename, job) l FROM emp
: string
| N
: string
| L
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename || ' ' || job label FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename || ' ' || job label FROM emp
: string
| N
: string
| LABEL
: string
| VARCHAR2
: string
| 20
: null
: null
: string
| B
: string
| 20
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CONCAT(ename, job) c FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CONCAT(ename, job) c FROM emp
: string
| N
: string
| C
: string
| VARCHAR2
: string
| 19
: null
: null
: string
| B
: string
| 19
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal || 'x' x FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal || 'x' x FROM emp
: string
| N
: string
| X
: string
| VARCHAR2
: string
| 41
: null
: null
: string
| B
: string
| 41
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal * 12 + NVL(comm, 0) annual FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal * 12 + NVL(comm, 0) annual FROM emp
: string
| N
: string
| ANNUAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal + ename x FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal + ename x FROM emp
: string
| N
: string
| X
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| M
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| R
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| T
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| F
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| P
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT sal / 3 s, MOD(empno, 2) m, ROUND(sal, 1) r, TRUNC(sal) t, ABS(sal) a, CEIL(sal) c, FLOOR(sal) f, POWER(sal, 2) p, SQRT(sal) q FROM emp
: string
| N
: string
| Q
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT SIGN(sal) s, LN(sal) l, EXP(1) e, LOG(10, sal) lg FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SIGN(sal) s, LN(sal) l, EXP(1) e, LOG(10, sal) lg FROM emp
: string
| N
: string
| S
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SIGN(sal) s, LN(sal) l, EXP(1) e, LOG(10, sal) lg FROM emp
: string
| N
: string
| L
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SIGN(sal) s, LN(sal) l, EXP(1) e, LOG(10, sal) lg FROM emp
: string
| N
: string
| E
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SIGN(sal) s, LN(sal) l, EXP(1) e, LOG(10, sal) lg FROM emp
: string
| N
: string
| LG
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
: string
| N
: string
| U
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
: string
| N
: string
| L
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
: string
| N
: string
| I
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
: string
| N
: string
| S
: string
| VARCHAR2
: string
| 12
: null
: null
: string
| B
: string
| 12
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(ename) u, LOWER(ename) l, INITCAP(ename) i, LENGTH(ename) n, SUBSTR(ename, 1, 3) s, INSTR(ename, 'A') p FROM emp
: string
| N
: string
| P
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| L
: string
| VARCHAR2
: string
| 60
: null
: null
: string
| B
: string
| 60
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| R
: string
| VARCHAR2
: string
| 48
: null
: null
: string
| B
: string
| 48
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| T
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| LT
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| RT
: string
| VARCHAR2
: string
| 9
: null
: null
: string
| B
: string
| 9
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| RP
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT LPAD(ename, 15, '*') l, RPAD(job, 12) r, TRIM(ename) t, LTRIM(job) lt, RTRIM(job) rt, REPLACE(ename, 'A', 'a') rp, TRANSLATE(job, 'AB', 'ab') tr FROM emp
: string
| N
: string
| TR
: string
| VARCHAR2
: string
| 36
: null
: null
: string
| B
: string
| 36
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT REGEXP_SUBSTR(ename, '[A-Z]+') rs, REGEXP_REPLACE(job, 'A', 'B') rr, REGEXP_INSTR(ename, 'A') ri, REGEXP_COUNT(ename, 'A') rc FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT REGEXP_SUBSTR(ename, '[A-Z]+') rs, REGEXP_REPLACE(job, 'A', 'B') rr, REGEXP_INSTR(ename, 'A') ri, REGEXP_COUNT(ename, 'A') rc FROM emp
: string
| N
: string
| RS
: string
| VARCHAR2
: string
| 40
: null
: null
: string
| B
: string
| 40
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT REGEXP_SUBSTR(ename, '[A-Z]+') rs, REGEXP_REPLACE(job, 'A', 'B') rr, REGEXP_INSTR(ename, 'A') ri, REGEXP_COUNT(ename, 'A') rc FROM emp
: string
| N
: string
| RR
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT REGEXP_SUBSTR(ename, '[A-Z]+') rs, REGEXP_REPLACE(job, 'A', 'B') rr, REGEXP_INSTR(ename, 'A') ri, REGEXP_COUNT(ename, 'A') rc FROM emp
: string
| N
: string
| RI
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT REGEXP_SUBSTR(ename, '[A-Z]+') rs, REGEXP_REPLACE(job, 'A', 'B') rr, REGEXP_INSTR(ename, 'A') ri, REGEXP_COUNT(ename, 'A') rc FROM emp
: string
| N
: string
| RC
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CHR(65) c, ASCII(ename) a FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CHR(65) c, ASCII(ename) a FROM emp
: string
| N
: string
| C
: string
| VARCHAR2
: string
| 1
: null
: null
: string
| B
: string
| 1
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CHR(65) c, ASCII(ename) a FROM emp
: string
| N
: string
| A
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT DUMP(ename) d, VSIZE(ename) v FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DUMP(ename) d, VSIZE(ename) v FROM emp
: string
| N
: string
| D
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT DUMP(ename) d, VSIZE(ename) v FROM emp
: string
| N
: string
| V
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT TO_CHAR(hiredate, 'YYYY-MM-DD') d, TO_CHAR(sal) s, TO_CHAR(sal, '9999.99') f FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CHAR(hiredate, 'YYYY-MM-DD') d, TO_CHAR(sal) s, TO_CHAR(sal, '9999.99') f FROM emp
: string
| N
: string
| D
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CHAR(hiredate, 'YYYY-MM-DD') d, TO_CHAR(sal) s, TO_CHAR(sal, '9999.99') f FROM emp
: string
| N
: string
| S
: string
| VARCHAR2
: string
| 40
: null
: null
: string
| B
: string
| 40
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CHAR(hiredate, 'YYYY-MM-DD') d, TO_CHAR(sal) s, TO_CHAR(sal, '9999.99') f FROM emp
: string
| N
: string
| F
: string
| VARCHAR2
: string
| 8
: null
: null
: string
| B
: string
| 8
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT TO_NUMBER('42') n, TO_DATE('2024-01-01', 'YYYY-MM-DD') d, TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD') t FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_NUMBER('42') n, TO_DATE('2024-01-01', 'YYYY-MM-DD') d, TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD') t FROM dual
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_NUMBER('42') n, TO_DATE('2024-01-01', 'YYYY-MM-DD') d, TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD') t FROM dual
: string
| N
: string
| D
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_NUMBER('42') n, TO_DATE('2024-01-01', 'YYYY-MM-DD') d, TO_TIMESTAMP('2024-01-01', 'YYYY-MM-DD') t FROM dual
: string
| N
: string
| T
: string
| TIMESTAMP(9)
: string
| 11
: null
: string
| 9
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual
: string
| N
: string
| NOW
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual
: string
| N
: string
| TS
: string
| TIMESTAMP(6) WITH TIME ZONE
: string
| 13
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual
: string
| N
: string
| CD
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual
: string
| N
: string
| CT
: string
| TIMESTAMP(6) WITH TIME ZONE
: string
| 13
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT SYSDATE now, SYSTIMESTAMP ts, CURRENT_DATE cd, CURRENT_TIMESTAMP ct, LOCALTIMESTAMP lt FROM dual
: string
| N
: string
| LT
: string
| TIMESTAMP(6)
: string
| 11
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT USER u, UID i, SYS_GUID() g, USERENV('LANG') l FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT USER u, UID i, SYS_GUID() g, USERENV('LANG') l FROM dual
: string
| N
: string
| U
: string
| VARCHAR2
: string
| 128
: null
: null
: string
| B
: string
| 128
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT USER u, UID i, SYS_GUID() g, USERENV('LANG') l FROM dual
: string
| N
: string
| I
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT USER u, UID i, SYS_GUID() g, USERENV('LANG') l FROM dual
: string
| N
: string
| G
: string
| RAW
: string
| 16
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT USER u, UID i, SYS_GUID() g, USERENV('LANG') l FROM dual
: string
| N
: string
| L
: string
| VARCHAR2
: string
| 52
: null
: null
: string
| B
: string
| 52
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
: string
| N
: string
| A
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
: string
| N
: string
| M
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
: string
| N
: string
| L
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
: string
| N
: string
| N
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
: string
| N
: string
| T
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ADD_MONTHS(hiredate, 1) a, MONTHS_BETWEEN(SYSDATE, hiredate) m, LAST_DAY(hiredate) l, NEXT_DAY(hiredate, 'MONDAY') n, TRUNC(hiredate, 'MM') t, EXTRACT(YEAR FROM hiredate) y FROM emp
: string
| N
: string
| Y
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT hiredate + 1 d, hiredate - hiredate n, SYSDATE - hiredate days FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT hiredate + 1 d, hiredate - hiredate n, SYSDATE - hiredate days FROM emp
: string
| N
: string
| D
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT hiredate + 1 d, hiredate - hiredate n, SYSDATE - hiredate days FROM emp
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT hiredate + 1 d, hiredate - hiredate n, SYSDATE - hiredate days FROM emp
: string
| N
: string
| DAYS
: string
| NUMBER
: string
| 22
: null
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp WHERE hiredate > 5
< error
| 932 ORA-00932: expression (5) is of data type NUMBER, which is incompatible with expected data type DATE
|
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CAST(sal AS VARCHAR2(20)) s, CAST(ename AS CHAR(12)) c, CAST(empno AS NUMBER(10,2)) n, CAST(hiredate AS TIMESTAMP) t FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(sal AS VARCHAR2(20)) s, CAST(ename AS CHAR(12)) c, CAST(empno AS NUMBER(10,2)) n, CAST(hiredate AS TIMESTAMP) t FROM emp
: string
| N
: string
| S
: string
| VARCHAR2
: string
| 20
: null
: null
: string
| B
: string
| 20
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(sal AS VARCHAR2(20)) s, CAST(ename AS CHAR(12)) c, CAST(empno AS NUMBER(10,2)) n, CAST(hiredate AS TIMESTAMP) t FROM emp
: string
| N
: string
| C
: string
| CHAR
: string
| 12
: null
: null
: string
| B
: string
| 12
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(sal AS VARCHAR2(20)) s, CAST(ename AS CHAR(12)) c, CAST(empno AS NUMBER(10,2)) n, CAST(hiredate AS TIMESTAMP) t FROM emp
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: string
| 10
: string
| 2
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(sal AS VARCHAR2(20)) s, CAST(ename AS CHAR(12)) c, CAST(empno AS NUMBER(10,2)) n, CAST(hiredate AS TIMESTAMP) t FROM emp
: string
| N
: string
| T
: string
| TIMESTAMP(6)
: string
| 11
: null
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT CAST(NULL AS DATE) d, CAST('1' AS INTEGER) i, CAST(1 AS BINARY_DOUBLE) b FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(NULL AS DATE) d, CAST('1' AS INTEGER) i, CAST(1 AS BINARY_DOUBLE) b FROM dual
: string
| N
: string
| D
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(NULL AS DATE) d, CAST('1' AS INTEGER) i, CAST(1 AS BINARY_DOUBLE) b FROM dual
: string
| N
: string
| I
: string
| NUMBER
: string
| 22
: null
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT CAST(NULL AS DATE) d, CAST('1' AS INTEGER) i, CAST(1 AS BINARY_DOUBLE) b FROM dual
: string
| N
: string
| B
: string
| BINARY_DOUBLE
: string
| 8
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| S
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| N
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| F
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| NS
: string
| NCHAR
: string
| 2
: null
: null
: string
| C
: string
| 1
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| D
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| T
: string
| TIMESTAMP(9)
: string
| 11
: null
: string
| 9
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT 'abc' s, 1 n, 1.5 f, N'x' ns, DATE '2024-01-01' d, TIMESTAMP '2024-01-01 00:00:00' t, NULL z FROM dual
: string
| N
: string
| Z
: string
| VARCHAR2
: string
| 0
: null
: null
: string
| B
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT TRUE b FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TRUE b FROM dual
: string
| N
: string
| B
: string
| BOOLEAN
: string
| 1
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| A
: string
| CHAR
: string
| 3
: null
: null
: string
| B
: string
| 3
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| B
: string
| VARCHAR2
: string
| 83
: null
: null
: string
| B
: string
| 83
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| D
: string
| CLOB
: string
| 32600
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| E
: string
| VARCHAR2
: string
| 40
: null
: null
: string
| B
: string
| 40
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| F
: string
| BINARY_DOUBLE
: string
| 8
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| G
: string
| BINARY_FLOAT
: string
| 4
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| H
: string
| DATE
: string
| 7
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| I
: string
| INTERVAL DAY(9) TO SECOND(6)
: string
| 24
: string
| 9
: string
| 6
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT UPPER(code) a, code || title b, LENGTH(note) c, SUBSTR(note, 1, 10) d, TO_CHAR(ts) e, bd + 1 f, bf * 2 g, ts + 1 h, ts - ts i, ds * 2 j FROM misc
: string
| N
: string
| J
: string
| INTERVAL DAY(9) TO SECOND(9)
: string
| 24
: string
| 9
: string
| 9
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT TO_CLOB(title) c, TO_NCHAR(title) nc, RAWTOHEX(r) h, HEXTORAW('AB') rr FROM misc
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CLOB(title) c, TO_NCHAR(title) nc, RAWTOHEX(r) h, HEXTORAW('AB') rr FROM misc
: string
| N
: string
| C
: string
| CLOB
: string
| 32600
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CLOB(title) c, TO_NCHAR(title) nc, RAWTOHEX(r) h, HEXTORAW('AB') rr FROM misc
: string
| N
: string
| NC
: string
| NVARCHAR2
: string
| 400
: null
: null
: string
| C
: string
| 200
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CLOB(title) c, TO_NCHAR(title) nc, RAWTOHEX(r) h, HEXTORAW('AB') rr FROM misc
: string
| N
: string
| H
: string
| VARCHAR2
: string
| 32
: null
: null
: string
| B
: string
| 32
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT TO_CLOB(title) c, TO_NCHAR(title) nc, RAWTOHEX(r) h, HEXTORAW('AB') rr FROM misc
: string
| N
: string
| RR
: string
| RAW
: string
| 1
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT JSON_OBJECT('name' VALUE ename) j FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT JSON_OBJECT('name' VALUE ename) j FROM emp
: string
| N
: string
| J
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{"a":1}', '$.a') j FROM dual
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT JSON_VALUE('{"a":1}', '$.a') j FROM dual
: string
| N
: string
| J
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT JSON_ARRAYAGG(ename) j FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT JSON_ARRAYAGG(ename) j FROM emp
: string
| N
: string
| J
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT XMLELEMENT("e", ename) x FROM emp
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT XMLELEMENT("e", ename) x FROM emp
: string
| N
: string
| X
: string
| XMLTYPE
: string
| 64
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT empno, xt.c FROM emp, JSON_TABLE('{"a":[1,2]}', '$.a[*]' COLUMNS (c NUMBER PATH '$')) xt
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT empno, xt.c FROM emp, JSON_TABLE('{"a":[1,2]}', '$.a[*]' COLUMNS (c NUMBER PATH '$')) xt
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| V
: string
| SELECT empno, xt.c FROM emp, JSON_TABLE('{"a":[1,2]}', '$.a[*]' COLUMNS (c NUMBER PATH '$')) xt
: string
| N
: string
| C
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT column_value cv FROM TABLE(sys.odcinumberlist(1, 2, 3))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT column_value cv FROM TABLE(sys.odcinumberlist(1, 2, 3))
: string
| N
: string
| CV
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM TABLE(sys.odcivarchar2list('a', 'b'))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "COLUMN_VALUE" FROM TABLE(sys.odcivarchar2list('a', 'b'))
: string
| N
: string
| COLUMN_VALUE
: string
| VARCHAR2
: string
| 4000
: null
: null
: string
| B
: string
| 4000
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename, sal FROM emp MODEL DIMENSION BY (empno) MEASURES (sal, ename) RULES (sal[7369] = 0)
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp MODEL DIMENSION BY (empno) MEASURES (sal, ename) RULES (sal[7369] = 0)
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename, sal FROM emp MODEL DIMENSION BY (empno) MEASURES (sal, ename) RULES (sal[7369] = 0)
: string
| N
: string
| SAL
: string
| NUMBER
: string
| 22
: null
: null
: null
: string
| 0
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT * FROM emp MATCH_RECOGNIZE (ORDER BY hiredate MEASURES FIRST(ename) AS first_name PATTERN (a b*) DEFINE b AS sal > PREV(sal))
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT "FIRST_NAME" FROM emp MATCH_RECOGNIZE (ORDER BY hiredate MEASURES FIRST(ename) AS first_name PATTERN (a b*) DEFINE b AS sal > PREV(sal))
: string
| N
: string
| FIRST_NAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| Y
: string
| NO
> exec
| BEGIN system.omni_sim.reset_schema(q'[CREATE TABLE dept (deptno NUMBER(2) CONSTRAINT dept_pk PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
| CREATE TABLE emp (empno NUMBER(4) CONSTRAINT emp_pk PRIMARY KEY, ename VARCHAR2(10) NOT NULL, job VARCHAR2(9), mgr NUMBER(4), hiredate DATE, sal NUMBER(7,2), comm NUMBER(7,2), deptno NUMBER(2) CONSTRAINT emp_dept_fk REFERENCES dept (deptno));
| CREATE TABLE misc (id INTEGER CONSTRAINT misc_pk PRIMARY KEY, code CHAR(3), title VARCHAR2(20 CHAR), name NVARCHAR2(20), note CLOB, ts TIMESTAMP, tz TIMESTAMP(3) WITH TIME ZONE, bd BINARY_DOUBLE, bf BINARY_FLOAT, r RAW(16), fl FLOAT, ym INTERVAL YEAR TO MONTH, ds INTERVAL DAY TO SECOND, n NUMBER(*,2), flag BOOLEAN);
| CREATE VIEW emp_v AS SELECT empno, ename, deptno FROM emp;]'); END;
< ok 1 0
> exec
| CREATE VIEW v AS SELECT ename FROM emp AS OF TIMESTAMP SYSDATE - 1
< ok 0 0
> query
| SELECT v.view_name, v.text_vc, v.read_only,
|   c.column_name, c.data_type, c.data_length, c.data_precision, c.data_scale,
|   c.char_used, c.char_length, c.nullable, c.hidden_column
| FROM user_views v JOIN user_tab_cols c ON c.table_name = v.view_name
| ORDER BY v.view_name, c.internal_column_id
< result
< column "VIEW_NAME" "NCHAR" "string"
< column "TEXT_VC" "NCHAR" "string"
< column "READ_ONLY" "NCHAR" "string"
< column "COLUMN_NAME" "NCHAR" "string"
< column "DATA_TYPE" "NCHAR" "string"
< column "DATA_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "DATA_PRECISION" "NUMBER" "float64" precision 38 scale 255
< column "DATA_SCALE" "NUMBER" "float64" precision 38 scale 255
< column "CHAR_USED" "NCHAR" "string"
< column "CHAR_LENGTH" "NUMBER" "float64" precision 38 scale 255
< column "NULLABLE" "NCHAR" "string"
< column "HIDDEN_COLUMN" "NCHAR" "string"
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| EMPNO
: string
| NUMBER
: string
| 22
: string
| 4
: string
| 0
: null
: string
| 0
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO
< row
: string
| EMP_V
: string
| SELECT empno, ename, deptno FROM emp
: string
| N
: string
| DEPTNO
: string
| NUMBER
: string
| 22
: string
| 2
: string
| 0
: null
: string
| 0
: string
| Y
: string
| NO
< row
: string
| V
: string
| SELECT ename FROM emp AS OF TIMESTAMP SYSDATE - 1
: string
| N
: string
| ENAME
: string
| VARCHAR2
: string
| 10
: null
: null
: string
| B
: string
| 10
: string
| N
: string
| NO

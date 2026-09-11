# PL/pgSQL Parser Scenarios

> Goal: Hand-written PL/pgSQL recursive descent parser matching omni pg/parser architecture with full PostgreSQL PL/pgSQL compatibility
> Verification: Go unit tests — each scenario is a PL/pgSQL body string that must parse to the correct AST (or produce the correct error)
> Reference sources: PostgreSQL src/pl/plpgsql/src/pl_gram.y, plpgsql.h, pl_scanner.c, pl_reserved_kwlist.h, pl_unreserved_kwlist.h; PostgreSQL docs Chapter 43

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Foundation

### 1.1 AST Node Types

- [x] PLBlock node — cmd_type, label, body (statement list), declarations, exception block
- [x] PLDeclare node — variable name, type reference, constant flag, not-null flag, collation, default expression
- [x] PLCursorDecl node — cursor name, scroll option, argument list, query text
- [x] PLAliasDecl node — alias name, referenced variable name
- [x] PLIf node — condition text, then-body, elsif list, else-body
- [x] PLCase node — test expression text, case-when list, else-body, has-test flag
- [x] PLLoop node — label, body
- [x] PLWhile node — label, condition text, body
- [x] PLForI node — label, variable, lower/upper/step expression texts, reverse flag, body
- [x] PLForS node — label, variable, query text, body
- [x] PLForC node — label, variable, cursor variable, arg query text, body
- [x] PLForDynS node — label, variable, query text, params list, body
- [x] PLForEachA node — label, variable, slice dimension, array expression text, body
- [x] PLReturn node — expression text (optional)
- [x] PLReturnNext node — expression text
- [x] PLReturnQuery node — query text (static), dynamic query text, params list
- [x] PLAssign node — variable target, expression text
- [x] PLExecSQL node — SQL text, into-target, strict flag
- [x] PLDynExecute node — query text, into-target, strict flag, params list
- [x] PLPerform node — expression text
- [x] PLCall node — SQL text (CALL or DO)
- [x] PLRaise node — level, condition name, sqlstate, message text, params, options list
- [x] PLAssert node — condition text, message text
- [x] PLGetDiag node — is-stacked flag, diagnostic items list
- [x] PLOpen node — cursor variable, query text, dynamic query text, params, scroll option
- [x] PLFetch node — cursor variable, direction, direction count, into-target, is-move flag
- [x] PLClose node — cursor variable
- [x] PLExit node — is-continue flag, label, condition text
- [x] PLCommit node — chain flag
- [x] PLRollback node — chain flag
- [x] PLNull node (null statement)
- [x] All nodes implement omni Node interface with Tag() method
- [x] All nodes carry Loc (start, end byte offsets)

### 1.2 Keywords and Tokens

- [x] Reserved keywords: ALL, BEGIN, BY, CASE, DECLARE, ELSE, END, EXECUTE, FOR, FOREACH, FROM, IF, IN, INTO, LOOP, NOT, NULL, OR, STRICT, THEN, TO, USING, WHEN, WHILE
- [x] Unreserved keywords recognized: ABSOLUTE, ALIAS, AND, ARRAY, ASSERT, BACKWARD, CALL, CHAIN, CLOSE, COLLATE, COLUMN, COLUMN_NAME, COMMIT, CONSTANT, CONSTRAINT, CONSTRAINT_NAME, CONTINUE, CURRENT, CURSOR, DATATYPE, DEBUG, DEFAULT, DETAIL, DIAGNOSTICS, DO, DUMP, ELSIF, ELSEIF, ERRCODE, ERROR, EXCEPTION, EXIT, FETCH, FIRST, FORWARD, GET, HINT, IMPORT, INFO, INSERT, IS, LAST, LOG, MERGE, MESSAGE, MESSAGE_TEXT, MOVE, NEXT, NO, NOTICE, OPEN, OPTION, PERFORM, PG_CONTEXT, PG_DATATYPE_NAME, PG_EXCEPTION_CONTEXT, PG_EXCEPTION_DETAIL, PG_EXCEPTION_HINT, PG_ROUTINE_OID, PRINT_STRICT_PARAMS, PRIOR, QUERY, RAISE, RELATIVE, RETURN, RETURNED_SQLSTATE, REVERSE, ROLLBACK, ROW_COUNT, ROWTYPE, SCHEMA, SCHEMA_NAME, SCROLL, SLICE, SQLSTATE, STACKED, TABLE, TABLE_NAME, TYPE, USE_COLUMN, USE_VARIABLE, VARIABLE_CONFLICT, WARNING
- [x] Unreserved keywords usable as variable names in declarations
- [x] Reserved keywords NOT usable as variable names
- [x] Keyword lookup is case-insensitive

### 1.3 Compiler Options

- [x] `#option dump` directive before block
- [x] `#print_strict_params on` directive
- [x] `#print_strict_params off` directive
- [x] `#variable_conflict error` directive
- [x] `#variable_conflict use_variable` directive
- [x] `#variable_conflict use_column` directive
- [x] Multiple directives before block
- [x] Unknown `#option` produces error

### 1.4 Parser Scaffold and Minimal Block

- [x] Entry point accepts PL/pgSQL function body string, returns AST root
- [x] Empty block: `BEGIN END` parses to PLBlock with empty body
- [x] Empty block with semicolons: `BEGIN END;` parses correctly
- [x] Block with DECLARE and no variables: `DECLARE BEGIN END`
- [x] Block with trailing label: `<<lbl>> BEGIN END lbl`
- [x] Label mismatch at END produces error
- [x] Mismatched END keyword (e.g. `END IF` in block) produces error
- [x] Unexpected token at top level produces error with position
- [x] Missing END produces error
- [x] Multiple statements in block separated by semicolons
- [x] Nested block: `BEGIN BEGIN END; END`
- [x] Labeled nested block: `<<inner>> BEGIN END inner;`
- [x] Label between DECLARE and BEGIN produces error: `DECLARE <<lbl>> BEGIN END;`
- [x] Extra DECLARE keyword within section: `DECLARE DECLARE x int; BEGIN END` parses correctly

### 1.5 DECLARE Section — Variable Declarations

- [x] Scalar: `DECLARE x integer;`
- [x] With default `:=`: `DECLARE x integer := 0;`
- [x] With default `=`: `DECLARE x integer = 0;`
- [x] With default `DEFAULT`: `DECLARE x integer DEFAULT 0;`
- [x] CONSTANT: `DECLARE x CONSTANT integer := 42;`
- [x] NOT NULL: `DECLARE x integer NOT NULL := 1;`
- [x] CONSTANT NOT NULL: `DECLARE x CONSTANT integer NOT NULL := 1;`
- [x] COLLATE: `DECLARE x text COLLATE "en_US";`
- [x] COLLATE with NOT NULL: `DECLARE x text COLLATE "C" NOT NULL := 'a';`
- [x] Record type: `DECLARE r record;`
- [x] Qualified type: `DECLARE x public.my_type;`
- [x] %TYPE reference: `DECLARE x my_table.my_column%TYPE;`
- [x] %ROWTYPE reference: `DECLARE r my_table%ROWTYPE;`
- [x] Multiple variables: `DECLARE x int; y text; z boolean;`
- [x] Cursor declaration: `DECLARE c CURSOR FOR SELECT 1;`
- [x] Cursor with SCROLL: `DECLARE c SCROLL CURSOR FOR SELECT 1;`
- [x] Cursor with NO SCROLL: `DECLARE c NO SCROLL CURSOR FOR SELECT 1;`
- [x] Cursor with parameters: `DECLARE c CURSOR (p1 int, p2 text) FOR SELECT p1, p2;`
- [x] Cursor with IS (synonym for FOR): `DECLARE c CURSOR IS SELECT 1;`
- [x] Alias declaration: `DECLARE a ALIAS FOR $1;`
- [x] Alias for named parameter: `DECLARE a ALIAS FOR param_name;`
- [x] Expression default spans until semicolon: `DECLARE x int := 1 + 2 * 3;`
- [x] Array type: `DECLARE a integer[];`
- [x] Array type with dimension: `DECLARE a integer[10];`
- [x] NOT NULL without default produces error: `DECLARE x integer NOT NULL;`
- [x] CONSTANT without default produces error: `DECLARE x CONSTANT integer;`

---

## Phase 2: Control Flow

### 2.1 IF / ELSIF / ELSE

- [x] Simple IF-THEN-END IF: `IF x > 0 THEN y := 1; END IF;`
- [x] IF-THEN-ELSE: `IF x > 0 THEN y := 1; ELSE y := 0; END IF;`
- [x] IF-THEN-ELSIF-ELSE: `IF x > 0 THEN y := 1; ELSIF x = 0 THEN y := 0; ELSE y := -1; END IF;`
- [x] Multiple ELSIF branches: three or more ELSIF clauses
- [x] ELSEIF synonym: `IF x THEN y := 1; ELSEIF z THEN y := 2; END IF;`
- [x] Nested IF inside IF: `IF a THEN IF b THEN c := 1; END IF; END IF;`
- [x] Condition expression spans until THEN: `IF x > 0 AND y < 10 THEN ...`
- [x] Empty THEN body: `IF x THEN END IF;` — parses with empty body
- [x] Multiple statements in THEN body
- [x] Multiple statements in ELSE body
- [x] Missing THEN produces error
- [x] Missing END IF produces error

### 2.2 CASE Statement

- [x] Searched CASE: `CASE WHEN x > 0 THEN y := 1; WHEN x = 0 THEN y := 0; END CASE;`
- [x] Searched CASE with ELSE: `CASE WHEN x > 0 THEN y := 1; ELSE y := -1; END CASE;`
- [x] Simple CASE: `CASE x WHEN 1 THEN y := 'a'; WHEN 2 THEN y := 'b'; END CASE;`
- [x] Simple CASE with ELSE: `CASE x WHEN 1 THEN y := 'a'; ELSE y := 'z'; END CASE;`
- [x] Multiple WHEN branches
- [x] Multiple statements per WHEN body
- [x] Nested CASE inside CASE
- [x] Missing END CASE produces error
- [x] No WHEN clause produces error

### 2.3 Basic Loops — LOOP, WHILE

- [x] Infinite LOOP: `LOOP x := x + 1; EXIT WHEN x > 10; END LOOP;`
- [x] Labeled LOOP: `<<myloop>> LOOP EXIT myloop; END LOOP myloop;`
- [x] WHILE loop: `WHILE x > 0 LOOP x := x - 1; END LOOP;`
- [x] Labeled WHILE: `<<w>> WHILE x > 0 LOOP x := x - 1; END LOOP w;`
- [x] Condition expression spans until LOOP keyword
- [x] Nested loops: LOOP inside LOOP
- [x] Empty loop body: `LOOP END LOOP;`
- [x] Missing END LOOP produces error
- [x] WHILE missing LOOP keyword produces error

### 2.4 FOR Loops — All Variants

- [x] Integer FOR: `FOR i IN 1..10 LOOP x := x + i; END LOOP;`
- [x] Integer FOR REVERSE: `FOR i IN REVERSE 10..1 LOOP x := x + i; END LOOP;`
- [x] Integer FOR with BY: `FOR i IN 1..100 BY 5 LOOP x := x + i; END LOOP;`
- [x] Integer FOR REVERSE with BY: `FOR i IN REVERSE 100..1 BY 10 LOOP END LOOP;`
- [x] Integer FOR with expression bounds: `FOR i IN a+1..b*2 LOOP END LOOP;`
- [x] Query FOR: `FOR rec IN SELECT * FROM t LOOP x := rec.a; END LOOP;`
- [x] Query FOR with complex query: `FOR rec IN SELECT a, b FROM t WHERE c > 0 ORDER BY a LOOP END LOOP;`
- [x] Cursor FOR: `FOR rec IN cur LOOP END LOOP;` (cur is cursor variable)
- [x] Cursor FOR with arguments: `FOR rec IN cur(1, 'a') LOOP END LOOP;`
- [x] Dynamic FOR: `FOR rec IN EXECUTE 'SELECT * FROM ' || tbl LOOP END LOOP;`
- [x] Dynamic FOR with USING: `FOR rec IN EXECUTE 'SELECT * FROM t WHERE id = $1' USING my_id LOOP END LOOP;`
- [x] FOREACH ARRAY: `FOREACH x IN ARRAY arr LOOP END LOOP;`
- [x] FOREACH ARRAY with SLICE: `FOREACH x SLICE 1 IN ARRAY arr LOOP END LOOP;`
- [x] Labeled FOR: `<<fl>> FOR i IN 1..10 LOOP END LOOP fl;`
- [x] Label mismatch on FOR loop END produces error
- [x] Nested FOR loops

### 2.5 EXIT and CONTINUE

- [x] Bare EXIT: `EXIT;`
- [x] EXIT with label: `EXIT myloop;`
- [x] EXIT with WHEN: `EXIT WHEN x > 10;`
- [x] EXIT with label and WHEN: `EXIT myloop WHEN x > 10;`
- [x] Bare CONTINUE: `CONTINUE;`
- [x] CONTINUE with label: `CONTINUE myloop;`
- [x] CONTINUE with WHEN: `CONTINUE WHEN x = 0;`
- [x] CONTINUE with label and WHEN: `CONTINUE myloop WHEN x = 0;`
- [x] WHEN condition expression spans until semicolon

---

## Phase 3: Statements

### 3.1 Variable Assignment

- [x] Simple assignment: `x := 1;`
- [x] Expression assignment: `x := a + b * c;`
- [x] Function call assignment: `x := my_func(a, b);`
- [x] Subquery assignment: `x := (SELECT max(a) FROM t);`
- [x] Record field assignment: `rec.field := 42;`
- [x] Array element assignment: `arr[1] := 'hello';`
- [x] Array slice assignment: `arr[1:3] := ARRAY[1,2,3];`
- [x] Multi-level field: `rec.nested.field := 1;`
- [x] Assignment with complex RHS spanning until semicolon
- [x] Assignment with `=` operator (PG also accepts this): `x = 1;`

### 3.2 RETURN Variants

- [x] Simple RETURN: `RETURN;`
- [x] RETURN with expression: `RETURN x + 1;`
- [x] RETURN with subquery: `RETURN (SELECT count(*) FROM t);`
- [x] RETURN NEXT: `RETURN NEXT x;`
- [x] RETURN NEXT with record: `RETURN NEXT rec;`
- [x] RETURN NEXT bare (for OUT params): `RETURN NEXT;`
- [x] RETURN QUERY static: `RETURN QUERY SELECT * FROM t;`
- [x] RETURN QUERY with WHERE: `RETURN QUERY SELECT a, b FROM t WHERE c > 0;`
- [x] RETURN QUERY EXECUTE: `RETURN QUERY EXECUTE 'SELECT * FROM t';`
- [x] RETURN QUERY EXECUTE with USING: `RETURN QUERY EXECUTE 'SELECT * FROM t WHERE id = $1' USING my_id;`
- [x] RETURN QUERY EXECUTE with multiple USING: `RETURN QUERY EXECUTE $q$ SELECT $1, $2 $q$ USING a, b;`
- [x] Expression text spans until semicolon for RETURN value

### 3.3 PERFORM and Bare SQL

- [x] PERFORM: `PERFORM my_func(1, 2);`
- [x] PERFORM with query: `PERFORM * FROM t WHERE a > 0;` (converted from SELECT)
- [x] Bare SQL — INSERT: `INSERT INTO t VALUES (1, 2);`
- [x] Bare SQL — UPDATE: `UPDATE t SET a = 1 WHERE b = 2;`
- [x] Bare SQL — DELETE: `DELETE FROM t WHERE a = 1;`
- [x] SQL with INTO: `SELECT a, b INTO x, y FROM t WHERE id = 1;`
- [x] SQL with INTO STRICT: `SELECT a INTO STRICT x FROM t WHERE id = 1;`
- [x] SQL statement text extracted until semicolon
- [x] SQL keywords (INSERT, UPDATE, DELETE, SELECT, MERGE, IMPORT) recognized as statement starters
- [x] Bare SQL — IMPORT: `IMPORT FOREIGN SCHEMA s FROM SERVER srv INTO public;`

### 3.4 EXECUTE Dynamic SQL

- [x] Basic EXECUTE: `EXECUTE 'SELECT 1';`
- [x] EXECUTE with INTO: `EXECUTE 'SELECT a FROM t' INTO x;`
- [x] EXECUTE with INTO STRICT: `EXECUTE 'SELECT a FROM t' INTO STRICT x;`
- [x] EXECUTE with USING: `EXECUTE 'INSERT INTO t VALUES($1)' USING val;`
- [x] EXECUTE with multiple USING: `EXECUTE $q$ SELECT $1, $2 $q$ USING a, b;`
- [x] EXECUTE with INTO and USING: `EXECUTE 'SELECT a FROM t WHERE id=$1' INTO x USING my_id;`
- [x] EXECUTE with USING before INTO: `EXECUTE 'SELECT $1' USING a INTO x;`
- [x] EXECUTE expression concatenation: `EXECUTE 'SELECT * FROM ' || quote_ident(tbl);`
- [x] EXECUTE with format(): `EXECUTE format('SELECT * FROM %I', tbl);`
- [x] Duplicate INTO produces error: `EXECUTE 'SELECT 1' INTO x INTO y;`
- [x] Duplicate USING produces error: `EXECUTE 'SELECT $1' USING a USING b;`

---

## Phase 4: Cursor Operations and Diagnostics

### 4.1 OPEN Cursor

- [x] Open bound cursor: `OPEN cur;`
- [x] Open bound cursor with arguments: `OPEN cur(1, 'a');`
- [x] Open unbound cursor for query: `OPEN cur FOR SELECT * FROM t;`
- [x] Open unbound cursor with SCROLL: `OPEN cur SCROLL FOR SELECT * FROM t;`
- [x] Open unbound cursor with NO SCROLL: `OPEN cur NO SCROLL FOR SELECT * FROM t;`
- [x] Open for EXECUTE: `OPEN cur FOR EXECUTE 'SELECT * FROM t';`
- [x] Open for EXECUTE with USING: `OPEN cur FOR EXECUTE 'SELECT * FROM t WHERE id=$1' USING my_id;`

### 4.2 FETCH and MOVE

- [x] FETCH default (NEXT): `FETCH cur INTO x;`
- [x] FETCH NEXT: `FETCH NEXT FROM cur INTO x;`
- [x] FETCH PRIOR: `FETCH PRIOR FROM cur INTO x;`
- [x] FETCH FIRST: `FETCH FIRST FROM cur INTO x;`
- [x] FETCH LAST: `FETCH LAST FROM cur INTO x;`
- [x] FETCH ABSOLUTE n: `FETCH ABSOLUTE 5 FROM cur INTO x;`
- [x] FETCH RELATIVE n: `FETCH RELATIVE -1 FROM cur INTO x;`
- [x] FETCH FORWARD: `FETCH FORWARD FROM cur INTO x;`
- [x] FETCH FORWARD n: `FETCH FORWARD 5 FROM cur INTO x;`
- [x] FETCH FORWARD ALL: `FETCH FORWARD ALL FROM cur INTO x;`
- [x] FETCH BACKWARD: `FETCH BACKWARD FROM cur INTO x;`
- [x] FETCH BACKWARD n: `FETCH BACKWARD 3 FROM cur INTO x;`
- [x] FETCH BACKWARD ALL: `FETCH BACKWARD ALL FROM cur INTO x;`
- [x] FETCH with IN instead of FROM: `FETCH NEXT IN cur INTO x;`
- [x] MOVE (no INTO): `MOVE NEXT FROM cur;`
- [x] MOVE with direction: `MOVE FORWARD ALL FROM cur;`
- [x] MOVE ABSOLUTE: `MOVE ABSOLUTE 0 FROM cur;`

### 4.3 CLOSE and GET DIAGNOSTICS

- [x] CLOSE cursor: `CLOSE cur;`
- [x] GET DIAGNOSTICS (current): `GET DIAGNOSTICS rowcount = ROW_COUNT;`
- [x] GET CURRENT DIAGNOSTICS: `GET CURRENT DIAGNOSTICS rowcount = ROW_COUNT;`
- [x] GET DIAGNOSTICS multiple items: `GET DIAGNOSTICS rc = ROW_COUNT, ctx = PG_CONTEXT;`
- [x] GET STACKED DIAGNOSTICS: `GET STACKED DIAGNOSTICS msg = MESSAGE_TEXT;`
- [x] GET STACKED DIAGNOSTICS multiple: `GET STACKED DIAGNOSTICS st = RETURNED_SQLSTATE, msg = MESSAGE_TEXT, det = PG_EXCEPTION_DETAIL;`
- [x] Assignment operator `:=` in diagnostics: `GET DIAGNOSTICS rc := ROW_COUNT;`
- [x] All current items: ROW_COUNT, PG_CONTEXT, PG_ROUTINE_OID
- [x] All stacked items: RETURNED_SQLSTATE, MESSAGE_TEXT, PG_EXCEPTION_DETAIL, PG_EXCEPTION_HINT, PG_EXCEPTION_CONTEXT, COLUMN_NAME, CONSTRAINT_NAME, PG_DATATYPE_NAME, TABLE_NAME, SCHEMA_NAME

---

## Phase 5: Advanced Statements

### 5.1 RAISE

- [x] Bare RAISE (re-raise): `RAISE;`
- [x] RAISE with level and message: `RAISE NOTICE 'hello %', name;`
- [x] RAISE DEBUG: `RAISE DEBUG 'debug info';`
- [x] RAISE LOG: `RAISE LOG 'log entry';`
- [x] RAISE INFO: `RAISE INFO 'info message';`
- [x] RAISE WARNING: `RAISE WARNING 'warning %', detail;`
- [x] RAISE EXCEPTION (default level): `RAISE EXCEPTION 'error occurred';`
- [x] RAISE without level (defaults to EXCEPTION): `RAISE 'error occurred';`
- [x] RAISE with multiple params: `RAISE NOTICE '% and % and %', a, b, c;`
- [x] RAISE with USING: `RAISE EXCEPTION 'fail' USING ERRCODE = 'P0001';`
- [x] RAISE with multiple USING options: `RAISE EXCEPTION 'fail' USING ERRCODE = 'P0001', DETAIL = 'more info', HINT = 'try this';`
- [x] RAISE with condition name: `RAISE division_by_zero;`
- [x] RAISE with SQLSTATE: `RAISE SQLSTATE '22012';`
- [x] RAISE USING only: `RAISE USING MESSAGE = 'dynamic error', ERRCODE = 'P0001';`
- [x] All USING options: MESSAGE, DETAIL, HINT, ERRCODE, COLUMN, CONSTRAINT, DATATYPE, TABLE, SCHEMA

### 5.2 ASSERT, CALL, Transaction Control, NULL

- [x] ASSERT: `ASSERT x > 0;`
- [x] ASSERT with message: `ASSERT x > 0, 'x must be positive';`
- [x] CALL procedure: `CALL my_proc(1, 2);`
- [x] CALL with schema: `CALL public.my_proc();`
- [x] DO block: `DO $$ BEGIN RAISE NOTICE 'hi'; END $$;`
- [x] COMMIT: `COMMIT;`
- [x] COMMIT AND CHAIN: `COMMIT AND CHAIN;`
- [x] COMMIT AND NO CHAIN: `COMMIT AND NO CHAIN;`
- [x] ROLLBACK: `ROLLBACK;`
- [x] ROLLBACK AND CHAIN: `ROLLBACK AND CHAIN;`
- [x] ROLLBACK AND NO CHAIN: `ROLLBACK AND NO CHAIN;`
- [x] NULL statement: `NULL;`
- [x] NULL statement in empty branches: `IF x THEN NULL; ELSE do_something(); END IF;`

### 5.3 Exception Handling

- [x] Single WHEN clause: `EXCEPTION WHEN division_by_zero THEN x := 0;`
- [x] Multiple WHEN clauses: `EXCEPTION WHEN no_data_found THEN ... WHEN too_many_rows THEN ...`
- [x] OR conditions: `EXCEPTION WHEN division_by_zero OR unique_violation THEN ...`
- [x] SQLSTATE code: `EXCEPTION WHEN SQLSTATE '22012' THEN ...`
- [x] OTHERS catch-all: `EXCEPTION WHEN OTHERS THEN ...`
- [x] Multiple statements in handler body
- [x] GET STACKED DIAGNOSTICS inside handler
- [x] RAISE inside handler (re-raise): `EXCEPTION WHEN OTHERS THEN RAISE;`
- [x] Nested block with own exception handler
- [x] Exception block with no statements before EXCEPTION
- [x] Exception handler with SQLSTATE and named conditions combined via OR: `WHEN SQLSTATE '23505' OR unique_violation THEN`

---

## Phase 6: Integration and Edge Cases

### 6.1 Nested and Complex Structures

- [x] Deeply nested blocks (3 levels): BEGIN BEGIN BEGIN ... END; END; END
- [x] Block inside IF branch
- [x] Block inside FOR loop body
- [x] Block inside exception handler
- [x] Loop inside IF inside block
- [x] FOR loop containing IF containing RETURN QUERY
- [x] Exception handler containing loop
- [x] Multiple DECLARE sections in nested blocks
- [x] Label scoping: inner label shadows outer

### 6.2 SQL Fragment Extraction

- [x] SQL expression in IF condition extracted as text until THEN
- [x] SQL expression in WHILE condition extracted until LOOP
- [x] SQL expression in RETURN extracted until semicolon
- [x] SQL in RETURN QUERY extracted as full SELECT text until semicolon
- [x] SQL in FOR..IN extracted correctly (distinguish integer range `..` from query)
- [x] SQL in EXECUTE extracted until INTO/USING/semicolon
- [x] SQL in PERFORM extracted until semicolon
- [x] SQL in assignment RHS extracted until semicolon
- [x] Parenthesized expressions in SQL fragments don't cause false termination
- [x] Dollar-quoted strings within SQL fragments preserved
- [x] String literals with semicolons in SQL fragments don't cause false termination
- [x] Nested parentheses in SQL fragments tracked correctly

### 6.3 Error Reporting

- [x] Error position points to offending token
- [x] Missing semicolon after statement
- [x] Unexpected keyword in expression context
- [x] Unterminated block (missing END)
- [x] Unterminated IF (missing END IF)
- [x] Unterminated LOOP (missing END LOOP)
- [x] Unterminated CASE (missing END CASE)
- [x] Unknown statement keyword produces meaningful error
- [x] Label without matching block/loop
- [x] END label mismatch: `<<a>> BEGIN END b;`
- [x] Empty input produces error or empty block
- [x] Duplicate DECLARE sections in same block (should be single DECLARE before BEGIN)

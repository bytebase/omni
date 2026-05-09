# Oracle Splitter Production Coverage

## Phase 1: Ordinary SQL Boundaries

- [x] Empty input returns no segments.
- [x] Whitespace-only input returns no segments.
- [x] Comment-only input returns no segments.
- [x] A single semicolon-terminated SQL statement returns one segment without the delimiter.
- [x] A single SQL statement without a semicolon returns one segment.
- [x] Multiple semicolon-terminated SQL statements return separate segments.
- [x] Mixed terminated and unterminated SQL statements return separate segments.
- [x] Consecutive semicolons do not produce empty segments.
- [x] Leading whitespace before SQL is preserved in segment text and range.
- [x] Inter-statement whitespace is attached to the following segment unless it is a skipped SQL*Plus command.
- [x] Trailing whitespace after the final statement is preserved in the final segment.
- [x] CRLF line endings do not corrupt boundaries.
- [x] UTF-8 comments do not corrupt byte ranges.
- [x] Arithmetic division does not act as a slash delimiter.
- [x] Schema-qualified and dblink-qualified names do not affect splitting.

## Phase 2: Literals, Identifiers, Comments, and Hints

- [x] Semicolons inside single-quoted string literals are ignored.
- [x] Slashes inside single-quoted string literals are ignored.
- [x] Escaped single quotes using doubled quotes are handled.
- [x] National character literals `N'...'` are handled.
- [x] Oracle `q'[ ... ]'` quoting is handled.
- [x] Oracle `q'{ ... }'` quoting is handled.
- [x] Oracle `q'( ... )'` quoting is handled.
- [x] Oracle `q'< ... >'` quoting is handled.
- [x] Oracle custom `q'!...!'` quoting is handled.
- [x] Semicolons inside double-quoted identifiers are ignored.
- [x] Slashes inside double-quoted identifiers are ignored.
- [x] Line comments containing semicolons are ignored.
- [x] Line comments containing slash-only visual markers are ignored.
- [x] Block comments containing semicolons are ignored.
- [x] Block comments containing slash-only lines are ignored.
- [x] Nested block comments are handled.
- [x] Hints `/*+ ... */` containing delimiters are ignored for splitting.

## Phase 3: SQL*Plus Slash and Buffer Commands

- [x] A slash alone on a line flushes the current SQL segment.
- [x] A slash line with surrounding horizontal whitespace flushes the current SQL segment.
- [x] A slash at EOF flushes the current SQL segment.
- [x] A slash followed by non-whitespace on the same line is not a delimiter.
- [x] A slash preceded by non-whitespace on the same line is not a delimiter.
- [x] `RUN` alone on a line flushes the current SQL segment.
- [x] `R` alone on a line flushes the current SQL segment.
- [x] Slash commands with leading command lines before SQL do not produce empty segments.
- [x] Multiple slash commands in a row do not produce empty segments.
- [x] Slash commands after PL/SQL units do not appear in segment text.

## Phase 4: Anonymous PL/SQL Blocks

- [x] `BEGIN ... END;` is returned as one segment.
- [x] `DECLARE ... BEGIN ... END;` is returned as one segment.
- [x] Labeled anonymous blocks are returned as one segment.
- [x] Nested anonymous blocks do not split at inner semicolons.
- [x] Exception handlers do not split the outer block.
- [x] `IF ... END IF;` stays inside the block.
- [x] `CASE ... END CASE;` stays inside the block.
- [x] `LOOP ... END LOOP;` stays inside the block.
- [x] `WHILE ... LOOP ... END LOOP;` stays inside the block.
- [x] `FOR ... LOOP ... END LOOP;` stays inside the block.
- [x] `FORALL` statements stay inside the block.
- [x] Dynamic SQL strings containing semicolons stay inside the block.

## Phase 5: Stored PL/SQL Units

- [x] `CREATE PROCEDURE ... END; /` is one segment.
- [x] `CREATE OR REPLACE PROCEDURE ... END; /` is one segment.
- [x] `CREATE EDITIONABLE PROCEDURE ... END; /` is one segment.
- [x] `CREATE NONEDITIONABLE PROCEDURE ... END; /` is one segment.
- [x] `CREATE FUNCTION ... END; /` is one segment.
- [x] Function bodies with arithmetic division do not split at `/`.
- [x] `CREATE PACKAGE ... END; /` is one segment.
- [x] `CREATE PACKAGE BODY ... END; /` is one segment.
- [x] Package body nested procedure bodies do not terminate the package.
- [x] Package body nested function bodies do not terminate the package.
- [x] `CREATE TRIGGER ... BEGIN ... END; /` is one segment.
- [x] Compound triggers are one segment.
- [x] `CREATE TYPE BODY ... END; /` is one segment.
- [x] Type body member function bodies do not terminate the type body.
- [x] Stored units followed by ordinary SQL split correctly.

## Phase 6: SQL*Plus Line Commands

- [x] `SET` command lines are skipped.
- [x] `SHOW` command lines are skipped.
- [x] `PROMPT` command lines are skipped.
- [x] `SPOOL` command lines are skipped.
- [x] `COLUMN` command lines are skipped.
- [x] `BREAK` command lines are skipped.
- [x] `COMPUTE` command lines are skipped.
- [x] `TTITLE` and `BTITLE` command lines are skipped.
- [x] `REPHEADER` and `REPFOOTER` command lines are skipped.
- [x] `DEFINE` and `UNDEFINE` command lines are skipped.
- [x] `ACCEPT` command lines are skipped.
- [x] `VARIABLE` and `PRINT` command lines are skipped.
- [x] `EXECUTE` command lines are skipped as SQL*Plus commands.
- [x] `CONNECT` and `DISCONNECT` command lines are skipped.
- [x] `EXIT` and `QUIT` command lines are skipped.
- [x] `WHENEVER SQLERROR` command lines are skipped.
- [x] `WHENEVER OSERROR` command lines are skipped.
- [x] `@`, `@@`, and `START` script invocation lines are skipped.
- [x] `REM` and `REMARK` command lines are skipped.
- [x] `HOST` and `!` shell command lines are skipped.
- [x] Buffer manipulation commands `LIST`, `APPEND`, `CHANGE`, `DEL`, `INPUT`, `SAVE`, `GET`, and `EDIT` are skipped.
- [x] SQL*Plus command words inside SQL are not skipped.
- [x] SQL*Plus command words inside PL/SQL are not skipped.

## Phase 7: Soft-Fail and Safety

- [x] Unterminated single-quoted strings do not panic.
- [x] Unterminated q-quotes do not panic.
- [x] Unterminated double-quoted identifiers do not panic.
- [x] Unterminated block comments do not panic.
- [x] Missing PL/SQL `END` returns a best-effort trailing segment.
- [x] Truncated `CREATE PROCEDURE` returns a best-effort trailing segment.
- [x] Binary-ish bytes do not panic.
- [x] Every returned segment has a valid non-overlapping byte range.
- [x] Every returned segment text equals `input[ByteStart:ByteEnd]`.
- [x] Reconstructing returned ranges never includes skipped SQL*Plus command lines.

## Phase 8: Facade Integration

- [x] `oracle.Parse` parses splitter output from ordinary SQL scripts.
- [x] `oracle.Parse` parses splitter output from slash-terminated PL/SQL scripts.
- [x] `oracle.Parse` skips SQL*Plus commands before parsing.
- [x] `oracle.Parse` preserves original byte offsets in AST locations.
- [x] `oracle.Parse` reports start line and column based on the original script.
- [x] `oracle.Parse` handles multiple statements around skipped command lines.

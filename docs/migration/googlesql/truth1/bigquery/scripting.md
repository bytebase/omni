# BigQuery GoogleSQL truth1 — Scripting / Procedural Language

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language

---

## SCRIPT-001: DECLARE

```
DECLARE variable_name[, ...] [variable_type] [DEFAULT expression];
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#declare
version: bigquery-current

```sql
DECLARE x INT64;
DECLARE d DATE DEFAULT CURRENT_DATE();
DECLARE x, y, z INT64 DEFAULT 0;
DECLARE item DEFAULT (SELECT item FROM schema1.products LIMIT 1);
```

---

## SCRIPT-002: SET

```
-- Single variable:
SET variable_name = expression;

-- Multiple variables (tuple assignment):
SET (variable_name[, ...]) = (expression[, ...]);
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#set
version: bigquery-current

```sql
SET x = 5;
SET (a, b, c) = (1 + 3, 'foo', false);
SET (corpus_count, word_count) = (
  SELECT AS STRUCT COUNT(DISTINCT corpus), SUM(word_count)
  FROM bigquery-public-data.samples.shakespeare
  WHERE LOWER(word) = target_word
);
```

---

## SCRIPT-003: EXECUTE IMMEDIATE

```
EXECUTE IMMEDIATE sql_expression [ INTO variable[, ...] ] [ USING identifier[, ...] ];

sql_expression:
  { "query_statement" | expression("query_statement") }

identifier:
  { variable | value } [ AS alias ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#execute_immediate
version: bigquery-current

```sql
-- Positional parameters:
DECLARE y INT64;
EXECUTE IMMEDIATE "SELECT ? * (? + 2)" INTO y USING 1, 3;

-- Named parameters:
EXECUTE IMMEDIATE "SELECT @a * (@b + 2)" INTO y USING 1 AS a, 3 AS b;

-- DDL via EXECUTE IMMEDIATE:
EXECUTE IMMEDIATE "CREATE TEMP TABLE Books (title STRING, publish_date INT64)";

-- Save result:
EXECUTE IMMEDIATE "SELECT MIN(publish_date) FROM Books LIMIT 1" INTO first_date;
```

---

## SCRIPT-004: BEGIN...END

```
BEGIN
  sql_statement_list
END;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#beginend
version: bigquery-current

```sql
DECLARE x INT64 DEFAULT 10;
BEGIN
  DECLARE y INT64;
  SET y = x;
  SELECT y;
END;
SELECT x;
```

---

## SCRIPT-005: BEGIN...EXCEPTION...END

```
BEGIN
  sql_statement_list
EXCEPTION WHEN ERROR THEN
  sql_statement_list
END;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#beginexceptionend
version: bigquery-current

```sql
BEGIN
  CALL schema1.proc2();
EXCEPTION WHEN ERROR THEN
  SELECT @@error.message, @@error.stack_trace, @@error.statement_text,
         @@error.formatted_stack_trace;
END;
```

---

## SCRIPT-006: CASE statement (search)

```
CASE
  WHEN boolean_expression THEN sql_statement_list
  [...]
  [ELSE sql_statement_list]
END CASE;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#case
version: bigquery-current

```sql
DECLARE target_product_id INT64 DEFAULT 103;
CASE
  WHEN EXISTS(SELECT 1 FROM schema.products_a WHERE product_id = target_product_id)
    THEN SELECT 'found product in products_a table';
  WHEN EXISTS(SELECT 1 FROM schema.products_b WHERE product_id = target_product_id)
    THEN SELECT 'found product in products_b table';
  ELSE SELECT 'did not find product';
END CASE;
```

---

## SCRIPT-007: CASE statement (simple)

```
CASE search_expression
  WHEN expression THEN sql_statement_list
  [...]
  [ELSE sql_statement_list]
END CASE;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#case_search_expression
version: bigquery-current

```sql
DECLARE product_id INT64 DEFAULT 1;
CASE product_id
  WHEN 1 THEN SELECT CONCAT('Product one');
  WHEN 2 THEN SELECT CONCAT('Product two');
  ELSE SELECT CONCAT('Invalid product');
END CASE;
```

---

## SCRIPT-008: IF / ELSEIF / ELSE

```
IF condition THEN [sql_statement_list]
  [ELSEIF condition THEN sql_statement_list]
  [...]
  [ELSE sql_statement_list]
END IF;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#if
version: bigquery-current

```sql
DECLARE target_product_id INT64 DEFAULT 103;
IF EXISTS(SELECT 1 FROM schema.products WHERE product_id = target_product_id) THEN
  SELECT CONCAT('found product ', CAST(target_product_id AS STRING));
ELSEIF EXISTS(SELECT 1 FROM schema.more_products WHERE product_id = target_product_id) THEN
  SELECT CONCAT('found product from more_products table', CAST(target_product_id AS STRING));
ELSE
  SELECT CONCAT('did not find product ', CAST(target_product_id AS STRING));
END IF;
```

---

## SCRIPT-009: Labels

```
label_name: BEGIN
  block_statement_list
END [label_name];

label_name: LOOP
  loop_statement_list
END LOOP [label_name];

label_name: WHILE condition DO
  loop_statement_list
END WHILE [label_name];

label_name: FOR variable IN query DO
  loop_statement_list
END FOR [label_name];

label_name: REPEAT
  loop_statement_list
  UNTIL boolean_condition
END REPEAT [label_name];

-- break/continue with label:
{ BREAK | LEAVE } label_name;
{ BREAK | LEAVE | CONTINUE | ITERATE } label_name;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#labels
version: bigquery-current

```sql
label_1: BEGIN
  SELECT 1;
  BREAK label_1;
  SELECT 2; -- Unreached
END;

label_1: LOOP
  WHILE x < 1 DO
    IF y < 1 THEN
      CONTINUE label_1;
    ELSE
      BREAK label_1;
    END IF;
  END WHILE;
END LOOP label_1;
```

---

## SCRIPT-010: LOOP

```
LOOP
  sql_statement_list
END LOOP;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#loop
version: bigquery-current

```sql
DECLARE x INT64 DEFAULT 0;
LOOP
  SET x = x + 1;
  IF x >= 10 THEN
    LEAVE;
  END IF;
END LOOP;
SELECT x;
```

---

## SCRIPT-011: REPEAT ... UNTIL

```
REPEAT
  sql_statement_list
  UNTIL boolean_condition
END REPEAT;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#repeat
version: bigquery-current

```sql
DECLARE x INT64 DEFAULT 0;
REPEAT
  SET x = x + 1;
  SELECT x;
  UNTIL x >= 3
END REPEAT;
```

---

## SCRIPT-012: WHILE

```
WHILE boolean_expression DO
  sql_statement_list
END WHILE;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#while
version: bigquery-current

```sql
DECLARE x INT64 DEFAULT 0;
WHILE x < 10 DO
  SET x = x + 1;
END WHILE;
```

---

## SCRIPT-013: BREAK / LEAVE

```
BREAK;
LEAVE;    -- synonym for BREAK
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#break
version: bigquery-current

```sql
LOOP
  IF condition THEN
    BREAK;
  END IF;
END LOOP;
```

---

## SCRIPT-014: CONTINUE / ITERATE

```
CONTINUE;
ITERATE;   -- synonym for CONTINUE
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#continue
version: bigquery-current

```sql
LOOP
  IF condition THEN
    CONTINUE;
  END IF;
  SELECT 'not skipped';
END LOOP;
```

---

## SCRIPT-015: FOR ... IN

```
FOR loop_variable_name IN (table_expression)
DO
  sql_expression_list
END FOR;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#for_in
version: bigquery-current

```sql
FOR record IN
  (SELECT word, word_count
   FROM bigquery-public-data.samples.shakespeare
   LIMIT 5)
DO
  SELECT record.word, record.word_count;
END FOR;
```

---

## SCRIPT-016: BEGIN TRANSACTION / COMMIT TRANSACTION / ROLLBACK TRANSACTION

```
BEGIN [TRANSACTION];
COMMIT [TRANSACTION];
ROLLBACK [TRANSACTION];
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#begin_transaction
version: bigquery-current

```sql
BEGIN TRANSACTION;
  INSERT INTO myschema.NewArrivals VALUES ('top load washer', 100, 'warehouse #1');
COMMIT TRANSACTION;

BEGIN
  BEGIN TRANSACTION;
  SELECT 1/0;  -- triggers error
  COMMIT TRANSACTION;
EXCEPTION WHEN ERROR THEN
  SELECT @@error.message;
  ROLLBACK TRANSACTION;
END;
```

---

## SCRIPT-017: RAISE

```
RAISE [USING MESSAGE = message];
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#raise
version: bigquery-current

```sql
BEGIN
  RAISE USING MESSAGE = 'Something went wrong';
EXCEPTION WHEN ERROR THEN
  SELECT @@error.message;
END;

-- Re-raise caught exception (must be inside EXCEPTION clause):
BEGIN
  SELECT 1/0;
EXCEPTION WHEN ERROR THEN
  RAISE;  -- re-raises with original stack trace
END;
```

---

## SCRIPT-018: RETURN

```
RETURN;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#return
version: bigquery-current

```sql
-- Stops execution of the multi-statement query
RETURN;
```

---

## SCRIPT-019: CALL

```
CALL procedure_name (procedure_argument[, ...])
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/procedural-language#call
version: bigquery-current

```sql
DECLARE retCode INT64;
-- Procedure signature: (IN account_id STRING, OUT retCode INT64)
CALL mySchema.UpdateSomeTables('someAccountId', retCode);
SELECT retCode;
```

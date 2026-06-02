# BigQuery GoogleSQL truth1 — DML Syntax

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/dml-syntax

---

## DML-001: INSERT statement

```
INSERT [INTO] target_name
 [(column_1 [, ..., column_n ] )]
 input

input ::=
 VALUES (expr_1 [, ..., expr_n ] )
        [, ..., (expr_k_1 [, ..., expr_k_n ] ) ]
| SELECT_QUERY

expr ::= value_expression | DEFAULT
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/dml-syntax#insert_statement
version: bigquery-current

```sql
INSERT dataset.Inventory (product, quantity)
VALUES('top load washer', 10),
      ('front load washer', 20),
      ('dryer', 30);

INSERT dataset.Warehouse (warehouse, state)
SELECT * FROM UNNEST([('warehouse #1', 'WA'), ('warehouse #2', 'CA')]);

INSERT dataset.Warehouse VALUES('warehouse #4', 'WA'), ('warehouse #5', 'NY');

-- Using DEFAULT keyword:
INSERT dataset.NewArrivals (product, quantity, warehouse)
VALUES('top load washer', DEFAULT, 'warehouse #1');
```

---

## DML-002: DELETE statement

```
DELETE [FROM] target_name [alias]
WHERE condition
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/dml-syntax#delete_statement
version: bigquery-current

```sql
DELETE dataset.Inventory WHERE quantity = 0;
DELETE FROM dataset.Inventory WHERE product = 'oven';
```

---

## DML-003: TRUNCATE TABLE statement

```
TRUNCATE TABLE [[project_name.]dataset_name.]table_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/dml-syntax#truncate_table_statement
version: bigquery-current

```sql
TRUNCATE TABLE dataset.Inventory;
```

---

## DML-004: UPDATE statement

```
UPDATE target_name [[AS] alias]
SET set_clause
[FROM from_clause]
WHERE condition

set_clause ::= update_item[, ...]

update_item ::= column_name = expression
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/dml-syntax#update_statement
version: bigquery-current

```sql
UPDATE dataset.Inventory
SET quantity = quantity - 10
WHERE product = 'oven';

-- UPDATE with FROM clause (join-based update):
UPDATE dataset.Inventory AS I
SET I.quantity = I.quantity + W.quantity
FROM dataset.NewArrivals AS W
WHERE I.product = W.product;

-- UPDATE using DEFAULT:
UPDATE dataset.mytable SET col = DEFAULT WHERE id = 1;
```

---

## DML-005: MERGE statement

```
MERGE [INTO] target_name [[AS] alias]
USING source_name [[AS] alias]
ON merge_condition
{ when_clause } +

when_clause ::= matched_clause | not_matched_by_target_clause | not_matched_by_source_clause

matched_clause ::= WHEN MATCHED [ AND search_condition ] THEN { merge_update_clause | merge_delete_clause }

not_matched_by_target_clause ::= WHEN NOT MATCHED [BY TARGET] [ AND search_condition ] THEN merge_insert_clause

not_matched_by_source_clause ::= WHEN NOT MATCHED BY SOURCE [ AND search_condition ] THEN { merge_update_clause | merge_delete_clause }

merge_condition ::= bool_expression

search_condition ::= bool_expression

merge_update_clause ::= UPDATE SET update_item [, update_item]*
update_item ::= column_name = expression

merge_delete_clause ::= DELETE

merge_insert_clause ::= INSERT [(column_1 [, ..., column_n ])] input

input ::= VALUES (expr_1 [, ..., expr_n ]) | ROW

expr ::= expression | DEFAULT
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/dml-syntax#merge_statement
version: bigquery-current

```sql
MERGE dataset.DetailedInventory T
USING dataset.Inventory S
ON T.product = S.product
WHEN NOT MATCHED AND quantity < 20 THEN
  INSERT(product, quantity, supply_constrained)
  VALUES(product, quantity, true)
WHEN NOT MATCHED THEN
  INSERT(product, quantity, supply_constrained)
  VALUES(product, quantity, false);

-- WHEN MATCHED UPDATE:
MERGE dataset.Inventory AS I
USING dataset.NewArrivals AS N
ON I.product = N.product
WHEN MATCHED THEN
  UPDATE SET quantity = I.quantity + N.quantity
WHEN NOT MATCHED BY TARGET THEN
  INSERT(product, quantity) VALUES(N.product, N.quantity)
WHEN NOT MATCHED BY SOURCE THEN
  DELETE;
```

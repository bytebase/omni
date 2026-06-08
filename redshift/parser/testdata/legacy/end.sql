-- END command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_END.html

-- Basic END statements
END;
end;

-- END with WORK keyword (optional synonym)
END WORK;
end work;

-- END with TRANSACTION keyword (optional synonym)
END TRANSACTION;
end transaction;

-- Examples from AWS documentation
-- The AWS docs show these three variations:
end;
end work;
end transaction;

-- Mixed case variations
End;
END;
EnD;
end Work;
END WORK;
End Work;
end transaction;
END TRANSACTION;
End Transaction;
END Transaction;

-- Real-world usage examples in transaction blocks
-- Transaction block with basic end
BEGIN;
INSERT INTO sales VALUES (1001, 'Widget A', 25.50);
END;

-- Transaction block with end work
BEGIN WORK;
UPDATE inventory SET quantity = quantity - 1 WHERE product_id = 'Widget A';
END WORK;

-- Transaction block with end transaction
BEGIN TRANSACTION;
INSERT INTO sales_summary SELECT product_name, SUM(amount) FROM sales GROUP BY product_name;
END TRANSACTION;

-- Nested transaction example (Redshift processes as single transaction)
BEGIN;
CREATE TEMP TABLE temp_sales AS SELECT * FROM sales WHERE amount > 100;
INSERT INTO high_value_sales SELECT * FROM temp_sales;
DROP TABLE temp_sales;
END;

-- Complex transaction with multiple operations
BEGIN WORK;
CREATE TABLE test_table (id INTEGER, name VARCHAR(50));
INSERT INTO test_table VALUES (1, 'Test Item 1');
INSERT INTO test_table VALUES (2, 'Test Item 2');
UPDATE test_table SET name = 'Updated Item' WHERE id = 1;
DELETE FROM test_table WHERE id = 2;
END WORK;

-- Transaction with DDL and DML operations
BEGIN TRANSACTION;
ALTER TABLE products ADD COLUMN description TEXT;
UPDATE products SET description = 'No description available' WHERE description IS NULL;
END TRANSACTION;

-- Error handling scenarios (these should parse correctly)
-- END after failed operations (would rollback in real usage)
BEGIN;
-- This would normally fail, but syntax should parse
INSERT INTO non_existent_table VALUES (1, 'test');
END;

-- Multiple END variations in sequence (each would be separate transactions)
BEGIN; INSERT INTO test VALUES (1); END;
BEGIN WORK; UPDATE test SET col1 = 'a' WHERE id = 1; END WORK;
BEGIN TRANSACTION; DELETE FROM test WHERE id = 1; END TRANSACTION;

-- Case insensitive parsing verification
begin; select 1; end;
BEGIN; SELECT 2; END;
Begin; Select 3; End;
begin work; select 4; end work;
BEGIN WORK; SELECT 5; END WORK;
Begin Work; Select 6; End Work;
begin transaction; select 7; end transaction;
BEGIN TRANSACTION; SELECT 8; END TRANSACTION;
Begin Transaction; Select 9; End Transaction;
-- COMMIT command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_COMMIT.html

-- Basic COMMIT statements
COMMIT;
commit;

-- COMMIT with WORK keyword (optional synonym)
COMMIT WORK;
commit work;

-- COMMIT with TRANSACTION keyword (optional synonym)
COMMIT TRANSACTION;
commit transaction;

-- END statements (alternative syntax for COMMIT)
END;
end;

-- END with WORK keyword
END WORK;
end work;

-- END with TRANSACTION keyword
END TRANSACTION;
end transaction;

-- Examples from AWS documentation
commit;
commit work;
commit transaction;

-- Mixed case variations
Commit;
COMMIT;
CoMmIt WoRk;
End;
END;
end WORK;
END transaction;
End Transaction;

-- Real-world usage examples
-- Transaction block with commit
BEGIN;
INSERT INTO test_table VALUES (1, 'test');
COMMIT;

-- Transaction block with commit work
BEGIN WORK;
UPDATE test_table SET column1 = 'updated' WHERE id = 1;
COMMIT WORK;

-- Transaction block with commit transaction
BEGIN TRANSACTION;
DELETE FROM test_table WHERE id = 1;
COMMIT TRANSACTION;

-- Transaction block with end
BEGIN;
CREATE TEMP TABLE temp_test AS SELECT * FROM source_table;
END;

-- Transaction block with end work
BEGIN WORK;
INSERT INTO target_table SELECT * FROM temp_test;
END WORK;

-- Transaction block with end transaction
BEGIN TRANSACTION;
DROP TABLE temp_test;
END TRANSACTION;
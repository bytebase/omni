-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_ReservedKeywordAcceptance default#1
> exec
| PREPARE omni_kw_stmt FROM ?
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ABSENT)
< error
| 1115 42000 Unknown character set: 'ABSENT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ACCESSIBLE)
< error
| 1064 42000 …near 'ACCESSIBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ACCOUNT)
< error
| 1064 42000 …near 'ACCOUNT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ACTION)
< error
| 1064 42000 …near 'ACTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ACTIVE)
< error
| 1064 42000 …near 'ACTIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ADD)
< error
| 1064 42000 …near 'ADD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ADDDATE)
< error
| 1064 42000 …near 'ADDDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ADMIN)
< error
| 1064 42000 …near 'ADMIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AFTER)
< error
| 1064 42000 …near 'AFTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AGAINST)
< error
| 1064 42000 …near 'AGAINST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AGGREGATE)
< error
| 1064 42000 …near 'AGGREGATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ALGORITHM)
< error
| 1064 42000 …near 'ALGORITHM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ALL)
< error
| 1064 42000 …near 'ALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ALLOW_MISSING_FILES)
< error
| 1115 42000 Unknown character set: 'ALLOW_MISSING_FILES'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ALTER)
< error
| 1064 42000 …near 'ALTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ALWAYS)
< error
| 1064 42000 …near 'ALWAYS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ANALYZE)
< error
| 1064 42000 …near 'ANALYZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AND)
< error
| 1064 42000 …near 'AND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ANY)
< error
| 1064 42000 …near 'ANY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AS)
< error
| 1064 42000 …near 'AS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ASC)
< error
| 1064 42000 …near 'ASC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ASCII)
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ASENSITIVE)
< error
| 1064 42000 …near 'ASENSITIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS)
< error
| 1115 42000 Unknown character set: 'ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AT)
< error
| 1064 42000 …near 'AT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ATTRIBUTE)
< error
| 1064 42000 …near 'ATTRIBUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AUTHENTICATION)
< error
| 1115 42000 Unknown character set: 'AUTHENTICATION'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AUTO)
< error
| 1064 42000 …near 'AUTO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AUTOEXTEND_SIZE)
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AUTO_INCREMENT)
< error
| 1064 42000 …near 'AUTO_INCREMENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AUTO_REFRESH)
< error
| 1064 42000 …near 'AUTO_REFRESH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AUTO_REFRESH_SOURCE)
< error
| 1115 42000 Unknown character set: 'AUTO_REFRESH_SOURCE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AVG)
< error
| 1064 42000 …near 'AVG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET AVG_ROW_LENGTH)
< error
| 1064 42000 …near 'AVG_ROW_LENGTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BACKUP)
< error
| 1064 42000 …near 'BACKUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BEFORE)
< error
| 1064 42000 …near 'BEFORE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BEGIN)
< error
| 1064 42000 …near 'BEGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BERNOULLI)
< error
| 1115 42000 Unknown character set: 'BERNOULLI'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BETWEEN)
< error
| 1064 42000 …near 'BETWEEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BINLOG)
< error
| 1064 42000 …near 'BINLOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BIT)
< error
| 1064 42000 …near 'BIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BIT_AND)
< error
| 1064 42000 …near 'BIT_AND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BIT_OR)
< error
| 1064 42000 …near 'BIT_OR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BIT_XOR)
< error
| 1064 42000 …near 'BIT_XOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BLOB)
< error
| 1064 42000 …near 'BLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BLOCK)
< error
| 1064 42000 …near 'BLOCK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BOOL)
< error
| 1064 42000 …near 'BOOL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BOTH)
< error
| 1064 42000 …near 'BOTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BTREE)
< error
| 1064 42000 …near 'BTREE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BUCKETS)
< error
| 1064 42000 …near 'BUCKETS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BULK)
< error
| 1115 42000 Unknown character set: 'BULK'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BY)
< error
| 1064 42000 …near 'BY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET BYTE)
< error
| 1064 42000 …near 'BYTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CACHE)
< error
| 1064 42000 …near 'CACHE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CALL)
< error
| 1064 42000 …near 'CALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CASCADE)
< error
| 1064 42000 …near 'CASCADE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CASCADED)
< error
| 1064 42000 …near 'CASCADED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CASE)
< error
| 1064 42000 …near 'CASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CAST)
< error
| 1064 42000 …near 'CAST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CATALOG_NAME)
< error
| 1064 42000 …near 'CATALOG_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHAIN)
< error
| 1064 42000 …near 'CHAIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHALLENGE_RESPONSE)
< error
| 1115 42000 Unknown character set: 'CHALLENGE_RESPONSE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHANGE)
< error
| 1064 42000 …near 'CHANGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHANGED)
< error
| 1064 42000 …near 'CHANGED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHANNEL)
< error
| 1115 42000 Unknown character set: 'CHANNEL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHAR)
< error
| 1064 42000 …near 'CHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHARACTER)
< error
| 1064 42000 …near 'CHARACTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHARSET)
< error
| 1064 42000 …near 'CHARSET)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHECK)
< error
| 1064 42000 …near 'CHECK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CHECKSUM)
< error
| 1064 42000 …near 'CHECKSUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CIPHER)
< error
| 1064 42000 …near 'CIPHER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CLASS_ORIGIN)
< error
| 1064 42000 …near 'CLASS_ORIGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CLIENT)
< error
| 1064 42000 …near 'CLIENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CLONE)
< error
| 1064 42000 …near 'CLONE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CLOSE)
< error
| 1064 42000 …near 'CLOSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COALESCE)
< error
| 1064 42000 …near 'COALESCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CODE)
< error
| 1064 42000 …near 'CODE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COLLATE)
< error
| 1064 42000 …near 'COLLATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COLLATION)
< error
| 1064 42000 …near 'COLLATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COLUMN)
< error
| 1064 42000 …near 'COLUMN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COLUMNS)
< error
| 1064 42000 …near 'COLUMNS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COLUMN_FORMAT)
< error
| 1064 42000 …near 'COLUMN_FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COLUMN_NAME)
< error
| 1064 42000 …near 'COLUMN_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMMENT)
< error
| 1064 42000 …near 'COMMENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMMIT)
< error
| 1064 42000 …near 'COMMIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMMITTED)
< error
| 1064 42000 …near 'COMMITTED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMPACT)
< error
| 1064 42000 …near 'COMPACT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMPLETION)
< error
| 1064 42000 …near 'COMPLETION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMPONENT)
< error
| 1115 42000 Unknown character set: 'COMPONENT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMPRESSED)
< error
| 1064 42000 …near 'COMPRESSED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COMPRESSION)
< error
| 1064 42000 …near 'COMPRESSION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONCAT)
< error
| 1115 42000 Unknown character set: 'CONCAT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONCURRENT)
< error
| 1064 42000 …near 'CONCURRENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONDITION)
< error
| 1064 42000 …near 'CONDITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONNECTION)
< error
| 1064 42000 …near 'CONNECTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONSISTENT)
< error
| 1064 42000 …near 'CONSISTENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONSTRAINT)
< error
| 1064 42000 …near 'CONSTRAINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONSTRAINT_CATALOG)
< error
| 1064 42000 …near 'CONSTRAINT_CATALOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONSTRAINT_NAME)
< error
| 1064 42000 …near 'CONSTRAINT_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONSTRAINT_SCHEMA)
< error
| 1064 42000 …near 'CONSTRAINT_SCHEMA)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONTAINS)
< error
| 1064 42000 …near 'CONTAINS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONTEXT)
< error
| 1064 42000 …near 'CONTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONTINUE)
< error
| 1064 42000 …near 'CONTINUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CONVERT)
< error
| 1064 42000 …near 'CONVERT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET COUNT)
< error
| 1064 42000 …near 'COUNT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CPU)
< error
| 1064 42000 …near 'CPU)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CREATE)
< error
| 1064 42000 …near 'CREATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CROSS)
< error
| 1064 42000 …near 'CROSS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CSV)
< error
| 1115 42000 Unknown character set: 'CSV'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CUBE)
< error
| 1064 42000 …near 'CUBE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CUME_DIST)
< error
| 1064 42000 …near 'CUME_DIST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURDATE)
< error
| 1064 42000 …near 'CURDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURRENT)
< error
| 1064 42000 …near 'CURRENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURRENT_DATE)
< error
| 1064 42000 …near 'CURRENT_DATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURRENT_TIME)
< error
| 1064 42000 …near 'CURRENT_TIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURRENT_TIMESTAMP)
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURRENT_USER)
< error
| 1064 42000 …near 'CURRENT_USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURSOR)
< error
| 1064 42000 …near 'CURSOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURSOR_NAME)
< error
| 1064 42000 …near 'CURSOR_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET CURTIME)
< error
| 1064 42000 …near 'CURTIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATA)
< error
| 1064 42000 …near 'DATA)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATABASE)
< error
| 1064 42000 …near 'DATABASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATABASES)
< error
| 1064 42000 …near 'DATABASES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATAFILE)
< error
| 1064 42000 …near 'DATAFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATE)
< error
| 1064 42000 …near 'DATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATETIME)
< error
| 1064 42000 …near 'DATETIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATE_ADD)
< error
| 1064 42000 …near 'DATE_ADD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DATE_SUB)
< error
| 1064 42000 …near 'DATE_SUB)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DAY)
< error
| 1064 42000 …near 'DAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DAY_HOUR)
< error
| 1064 42000 …near 'DAY_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DAY_MICROSECOND)
< error
| 1064 42000 …near 'DAY_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DAY_MINUTE)
< error
| 1064 42000 …near 'DAY_MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DAY_SECOND)
< error
| 1064 42000 …near 'DAY_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DEALLOCATE)
< error
| 1064 42000 …near 'DEALLOCATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DEC)
< error
| 1064 42000 …near 'DEC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DECIMAL)
< error
| 1064 42000 …near 'DECIMAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DECLARE)
< error
| 1064 42000 …near 'DECLARE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DEFAULT)
< error
| 1064 42000 …near 'DEFAULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DEFAULT_AUTH)
< error
| 1064 42000 …near 'DEFAULT_AUTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DEFINER)
< error
| 1064 42000 …near 'DEFINER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DEFINITION)
< error
| 1115 42000 Unknown character set: 'DEFINITION'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DELAYED)
< error
| 1064 42000 …near 'DELAYED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DELAY_KEY_WRITE)
< error
| 1064 42000 …near 'DELAY_KEY_WRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DELETE)
< error
| 1064 42000 …near 'DELETE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DENSE_RANK)
< error
| 1064 42000 …near 'DENSE_RANK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DESC)
< error
| 1064 42000 …near 'DESC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DESCRIBE)
< error
| 1064 42000 …near 'DESCRIBE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DESCRIPTION)
< error
| 1064 42000 …near 'DESCRIPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DETERMINISTIC)
< error
| 1064 42000 …near 'DETERMINISTIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DIAGNOSTICS)
< error
| 1064 42000 …near 'DIAGNOSTICS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DIRECTORY)
< error
| 1064 42000 …near 'DIRECTORY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DISABLE)
< error
| 1064 42000 …near 'DISABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DISCARD)
< error
| 1064 42000 …near 'DISCARD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DISK)
< error
| 1064 42000 …near 'DISK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DISTINCT)
< error
| 1064 42000 …near 'DISTINCT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DISTINCTROW)
< error
| 1064 42000 …near 'DISTINCTROW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DIV)
< error
| 1064 42000 …near 'DIV)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DO)
< error
| 1064 42000 …near 'DO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DOUBLE)
< error
| 1064 42000 …near 'DOUBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DROP)
< error
| 1064 42000 …near 'DROP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DUAL)
< error
| 1064 42000 …near 'DUAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DUALITY)
< error
| 1115 42000 Unknown character set: 'DUALITY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DUMPFILE)
< error
| 1064 42000 …near 'DUMPFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DUPLICATE)
< error
| 1064 42000 …near 'DUPLICATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET DYNAMIC)
< error
| 1064 42000 …near 'DYNAMIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EACH)
< error
| 1064 42000 …near 'EACH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ELSE)
< error
| 1064 42000 …near 'ELSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ELSEIF)
< error
| 1064 42000 …near 'ELSEIF)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EMPTY)
< error
| 1064 42000 …near 'EMPTY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENABLE)
< error
| 1064 42000 …near 'ENABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENCLOSED)
< error
| 1064 42000 …near 'ENCLOSED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENCRYPTION)
< error
| 1064 42000 …near 'ENCRYPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET END)
< error
| 1064 42000 …near 'END)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENDS)
< error
| 1064 42000 …near 'ENDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENFORCED)
< error
| 1064 42000 …near 'ENFORCED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENGINE)
< error
| 1064 42000 …near 'ENGINE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENGINES)
< error
| 1064 42000 …near 'ENGINES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENGINE_ATTRIBUTE)
< error
| 1115 42000 Unknown character set: 'ENGINE_ATTRIBUTE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ENUM)
< error
| 1064 42000 …near 'ENUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ERROR)
< error
| 1064 42000 …near 'ERROR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ERRORS)
< error
| 1064 42000 …near 'ERRORS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ESCAPE)
< error
| 1064 42000 …near 'ESCAPE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ESCAPED)
< error
| 1064 42000 …near 'ESCAPED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EVENT)
< error
| 1064 42000 …near 'EVENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EVENTS)
< error
| 1064 42000 …near 'EVENTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EVERY)
< error
| 1064 42000 …near 'EVERY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXCEPT)
< error
| 1064 42000 …near 'EXCEPT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXCHANGE)
< error
| 1064 42000 …near 'EXCHANGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXCLUDE)
< error
| 1115 42000 Unknown character set: 'EXCLUDE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXECUTE)
< error
| 1064 42000 …near 'EXECUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXISTS)
< error
| 1064 42000 …near 'EXISTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXIT)
< error
| 1064 42000 …near 'EXIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXPANSION)
< error
| 1064 42000 …near 'EXPANSION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXPIRE)
< error
| 1064 42000 …near 'EXPIRE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXPLAIN)
< error
| 1064 42000 …near 'EXPLAIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXPORT)
< error
| 1064 42000 …near 'EXPORT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXTENDED)
< error
| 1064 42000 …near 'EXTENDED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXTENT_SIZE)
< error
| 1064 42000 …near 'EXTENT_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXTERNAL)
< error
| 1064 42000 …near 'EXTERNAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXTERNAL_FORMAT)
< error
| 1115 42000 Unknown character set: 'EXTERNAL_FORMAT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET EXTRACT)
< error
| 1064 42000 …near 'EXTRACT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FACTOR)
< error
| 1115 42000 Unknown character set: 'FACTOR'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FAILED_LOGIN_ATTEMPTS)
< error
| 1115 42000 Unknown character set: 'FAILED_LOGIN_ATTEMPTS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FALSE)
< error
| 1064 42000 …near 'FALSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FAST)
< error
| 1064 42000 …near 'FAST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FAULTS)
< error
| 1064 42000 …near 'FAULTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FETCH)
< error
| 1064 42000 …near 'FETCH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FIELDS)
< error
| 1064 42000 …near 'FIELDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILE)
< error
| 1064 42000 …near 'FILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILES)
< error
| 1064 42000 …near 'FILES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILE_BLOCK_SIZE)
< error
| 1064 42000 …near 'FILE_BLOCK_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILE_FORMAT)
< error
| 1115 42000 Unknown character set: 'FILE_FORMAT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILE_NAME)
< error
| 1115 42000 Unknown character set: 'FILE_NAME'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILE_PATTERN)
< error
| 1115 42000 Unknown character set: 'FILE_PATTERN'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILE_PREFIX)
< error
| 1115 42000 Unknown character set: 'FILE_PREFIX'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FILTER)
< error
| 1064 42000 …near 'FILTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FINISH)
< error
| 1115 42000 Unknown character set: 'FINISH'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FIRST)
< error
| 1064 42000 …near 'FIRST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FIRST_VALUE)
< error
| 1064 42000 …near 'FIRST_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FIXED)
< error
| 1064 42000 …near 'FIXED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FLOAT)
< error
| 1064 42000 …near 'FLOAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FLOAT4)
< error
| 1064 42000 …near 'FLOAT4)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FLOAT8)
< error
| 1064 42000 …near 'FLOAT8)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FLUSH)
< error
| 1064 42000 …near 'FLUSH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FOLLOWING)
< error
| 1064 42000 …near 'FOLLOWING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FOLLOWS)
< error
| 1115 42000 Unknown character set: 'FOLLOWS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FOR)
< error
| 1064 42000 …near 'FOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FORCE)
< error
| 1064 42000 …near 'FORCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FOREIGN)
< error
| 1064 42000 …near 'FOREIGN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FORMAT)
< error
| 1064 42000 …near 'FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FOUND)
< error
| 1064 42000 …near 'FOUND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FROM)
< error
| 1064 42000 …near 'FROM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FULL)
< error
| 1064 42000 …near 'FULL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FULLTEXT)
< error
| 1064 42000 …near 'FULLTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET FUNCTION)
< error
| 1064 42000 …near 'FUNCTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GENERAL)
< error
| 1064 42000 …near 'GENERAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GENERATE)
< error
| 1115 42000 Unknown character set: 'GENERATE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GENERATED)
< error
| 1064 42000 …near 'GENERATED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GEOMCOLLECTION)
< error
| 1064 42000 …near 'GEOMCOLLECTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GEOMETRY)
< error
| 1064 42000 …near 'GEOMETRY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GEOMETRYCOLLECTION)
< error
| 1064 42000 …near 'GEOMETRYCOLLECTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GET)
< error
| 1064 42000 …near 'GET)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GET_FORMAT)
< error
| 1064 42000 …near 'GET_FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GET_SOURCE_PUBLIC_KEY)
< error
| 1115 42000 Unknown character set: 'GET_SOURCE_PUBLIC_KEY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GLOBAL)
< error
| 1064 42000 …near 'GLOBAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GRANT)
< error
| 1064 42000 …near 'GRANT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GRANTS)
< error
| 1064 42000 …near 'GRANTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GREATEST)
< error
| 1115 42000 Unknown character set: 'GREATEST'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GROUP)
< error
| 1064 42000 …near 'GROUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GROUPING)
< error
| 1064 42000 …near 'GROUPING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GROUPS)
< error
| 1115 42000 Unknown character set: 'GROUPS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GROUP_CONCAT)
< error
| 1064 42000 …near 'GROUP_CONCAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GROUP_REPLICATION)
< error
| 1115 42000 Unknown character set: 'GROUP_REPLICATION'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GTIDS)
< error
| 1115 42000 Unknown character set: 'GTIDS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GTID_ONLY)
< error
| 1115 42000 Unknown character set: 'GTID_ONLY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET GUIDED)
< error
| 1115 42000 Unknown character set: 'GUIDED'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HANDLER)
< error
| 1064 42000 …near 'HANDLER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HASH)
< error
| 1064 42000 …near 'HASH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HAVING)
< error
| 1064 42000 …near 'HAVING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HEADER)
< error
| 1115 42000 Unknown character set: 'HEADER'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HELP)
< error
| 1064 42000 …near 'HELP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HIGH_PRIORITY)
< error
| 1064 42000 …near 'HIGH_PRIORITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HISTOGRAM)
< error
| 1064 42000 …near 'HISTOGRAM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HISTORY)
< error
| 1115 42000 Unknown character set: 'HISTORY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HOST)
< error
| 1064 42000 …near 'HOST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HOSTS)
< error
| 1064 42000 …near 'HOSTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HOUR)
< error
| 1064 42000 …near 'HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HOUR_MICROSECOND)
< error
| 1064 42000 …near 'HOUR_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HOUR_MINUTE)
< error
| 1064 42000 …near 'HOUR_MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET HOUR_SECOND)
< error
| 1064 42000 …near 'HOUR_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IDENTIFIED)
< error
| 1064 42000 …near 'IDENTIFIED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IF)
< error
| 1064 42000 …near 'IF)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IGNORE)
< error
| 1064 42000 …near 'IGNORE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IGNORE_SERVER_IDS)
< error
| 1064 42000 …near 'IGNORE_SERVER_IDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IMPORT)
< error
| 1064 42000 …near 'IMPORT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IN)
< error
| 1064 42000 …near 'IN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INACTIVE)
< error
| 1115 42000 Unknown character set: 'INACTIVE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INDEX)
< error
| 1064 42000 …near 'INDEX)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INDEXES)
< error
| 1064 42000 …near 'INDEXES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INFILE)
< error
| 1064 42000 …near 'INFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INITIAL)
< error
| 1115 42000 Unknown character set: 'INITIAL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INITIAL_SIZE)
< error
| 1064 42000 …near 'INITIAL_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INITIATE)
< error
| 1115 42000 Unknown character set: 'INITIATE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INNER)
< error
| 1064 42000 …near 'INNER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INOUT)
< error
| 1064 42000 …near 'INOUT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INSENSITIVE)
< error
| 1064 42000 …near 'INSENSITIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INSERT)
< error
| 1064 42000 …near 'INSERT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INSERT_METHOD)
< error
| 1064 42000 …near 'INSERT_METHOD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INSTALL)
< error
| 1064 42000 …near 'INSTALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INSTANCE)
< error
| 1064 42000 …near 'INSTANCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INT1)
< error
| 1064 42000 …near 'INT1)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INT2)
< error
| 1064 42000 …near 'INT2)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INT3)
< error
| 1064 42000 …near 'INT3)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INT4)
< error
| 1064 42000 …near 'INT4)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INT8)
< error
| 1064 42000 …near 'INT8)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INTEGER)
< error
| 1064 42000 …near 'INTEGER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INTERSECT)
< error
| 1064 42000 …near 'INTERSECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INTERVAL)
< error
| 1064 42000 …near 'INTERVAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INTO)
< error
| 1064 42000 …near 'INTO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INVISIBLE)
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET INVOKER)
< error
| 1064 42000 …near 'INVOKER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IO)
< error
| 1064 42000 …near 'IO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IO_AFTER_GTIDS)
< error
| 1064 42000 …near 'IO_AFTER_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IO_BEFORE_GTIDS)
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IO_THREAD)
< error
| 1064 42000 …near 'IO_THREAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IPC)
< error
| 1064 42000 …near 'IPC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET IS)
< error
| 1064 42000 …near 'IS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ISOLATION)
< error
| 1064 42000 …near 'ISOLATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ISSUER)
< error
| 1064 42000 …near 'ISSUER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ITERATE)
< error
| 1064 42000 …near 'ITERATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JOIN)
< error
| 1064 42000 …near 'JOIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JSON)
< error
| 1064 42000 …near 'JSON)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JSON_ARRAYAGG)
< error
| 1064 42000 …near 'JSON_ARRAYAGG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JSON_DUALITY_OBJECT)
< error
| 1115 42000 Unknown character set: 'JSON_DUALITY_OBJECT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JSON_OBJECTAGG)
< error
| 1064 42000 …near 'JSON_OBJECTAGG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JSON_TABLE)
< error
| 1064 42000 …near 'JSON_TABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET JSON_VALUE)
< error
| 1064 42000 …near 'JSON_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET KEY)
< error
| 1064 42000 …near 'KEY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET KEYRING)
< error
| 1115 42000 Unknown character set: 'KEYRING'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET KEYS)
< error
| 1064 42000 …near 'KEYS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET KEY_BLOCK_SIZE)
< error
| 1064 42000 …near 'KEY_BLOCK_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET KILL)
< error
| 1064 42000 …near 'KILL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LAG)
< error
| 1064 42000 …near 'LAG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LANGUAGE)
< error
| 1064 42000 …near 'LANGUAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LAST)
< error
| 1064 42000 …near 'LAST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LAST_VALUE)
< error
| 1064 42000 …near 'LAST_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LATERAL)
< error
| 1064 42000 …near 'LATERAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEAD)
< error
| 1064 42000 …near 'LEAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEADING)
< error
| 1064 42000 …near 'LEADING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEAST)
< error
| 1115 42000 Unknown character set: 'LEAST'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEAVE)
< error
| 1064 42000 …near 'LEAVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEAVES)
< error
| 1064 42000 …near 'LEAVES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEFT)
< error
| 1064 42000 …near 'LEFT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LESS)
< error
| 1064 42000 …near 'LESS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LEVEL)
< error
| 1064 42000 …near 'LEVEL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LIBRARY)
< error
| 1115 42000 Unknown character set: 'LIBRARY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LIKE)
< error
| 1064 42000 …near 'LIKE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LIMIT)
< error
| 1064 42000 …near 'LIMIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LINEAR)
< error
| 1064 42000 …near 'LINEAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LINES)
< error
| 1064 42000 …near 'LINES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LINESTRING)
< error
| 1064 42000 …near 'LINESTRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LIST)
< error
| 1064 42000 …near 'LIST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOAD)
< error
| 1064 42000 …near 'LOAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOCAL)
< error
| 1064 42000 …near 'LOCAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOCALTIME)
< error
| 1064 42000 …near 'LOCALTIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOCALTIMESTAMP)
< error
| 1064 42000 …near 'LOCALTIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOCK)
< error
| 1064 42000 …near 'LOCK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOCKED)
< error
| 1064 42000 …near 'LOCKED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOCKS)
< error
| 1064 42000 …near 'LOCKS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOG)
< error
| 1064 42000 …near 'LOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOGFILE)
< error
| 1064 42000 …near 'LOGFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOGS)
< error
| 1064 42000 …near 'LOGS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LONG)
< error
| 1064 42000 …near 'LONG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LONGBLOB)
< error
| 1064 42000 …near 'LONGBLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LONGTEXT)
< error
| 1064 42000 …near 'LONGTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOOP)
< error
| 1064 42000 …near 'LOOP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET LOW_PRIORITY)
< error
| 1064 42000 …near 'LOW_PRIORITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MANIFEST)
< error
| 1115 42000 Unknown character set: 'MANIFEST'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MANUAL)
< error
| 1064 42000 …near 'MANUAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MASTER)
< error
| 1064 42000 …near 'MASTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MATCH)
< error
| 1064 42000 …near 'MATCH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MATERIALIZED)
< error
| 1064 42000 …near 'MATERIALIZED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX)
< error
| 1064 42000 …near 'MAX)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAXVALUE)
< error
| 1064 42000 …near 'MAXVALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX_CONNECTIONS_PER_HOUR)
< error
| 1064 42000 …near 'MAX_CONNECTIONS_PER_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX_QUERIES_PER_HOUR)
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX_ROWS)
< error
| 1064 42000 …near 'MAX_ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX_SIZE)
< error
| 1064 42000 …near 'MAX_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX_UPDATES_PER_HOUR)
< error
| 1064 42000 …near 'MAX_UPDATES_PER_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MAX_USER_CONNECTIONS)
< error
| 1064 42000 …near 'MAX_USER_CONNECTIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MEDIUM)
< error
| 1064 42000 …near 'MEDIUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MEDIUMBLOB)
< error
| 1064 42000 …near 'MEDIUMBLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MEDIUMINT)
< error
| 1064 42000 …near 'MEDIUMINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MEDIUMTEXT)
< error
| 1064 42000 …near 'MEDIUMTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MEMBER)
< error
| 1064 42000 …near 'MEMBER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MEMORY)
< error
| 1064 42000 …near 'MEMORY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MERGE)
< error
| 1064 42000 …near 'MERGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MESSAGE_TEXT)
< error
| 1064 42000 …near 'MESSAGE_TEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MICROSECOND)
< error
| 1064 42000 …near 'MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MID)
< error
| 1115 42000 Unknown character set: 'MID'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MIDDLEINT)
< error
| 1064 42000 …near 'MIDDLEINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MIGRATE)
< error
| 1064 42000 …near 'MIGRATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MIN)
< error
| 1064 42000 …near 'MIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MINUTE)
< error
| 1064 42000 …near 'MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MINUTE_MICROSECOND)
< error
| 1064 42000 …near 'MINUTE_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MINUTE_SECOND)
< error
| 1064 42000 …near 'MINUTE_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MIN_ROWS)
< error
| 1064 42000 …near 'MIN_ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MOD)
< error
| 1064 42000 …near 'MOD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MODE)
< error
| 1064 42000 …near 'MODE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MODIFIES)
< error
| 1064 42000 …near 'MODIFIES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MODIFY)
< error
| 1064 42000 …near 'MODIFY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MONTH)
< error
| 1064 42000 …near 'MONTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MULTILINESTRING)
< error
| 1064 42000 …near 'MULTILINESTRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MULTIPOINT)
< error
| 1064 42000 …near 'MULTIPOINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MULTIPOLYGON)
< error
| 1064 42000 …near 'MULTIPOLYGON)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MUTEX)
< error
| 1064 42000 …near 'MUTEX)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET MYSQL_ERRNO)
< error
| 1064 42000 …near 'MYSQL_ERRNO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NAME)
< error
| 1064 42000 …near 'NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NAMES)
< error
| 1064 42000 …near 'NAMES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NATIONAL)
< error
| 1064 42000 …near 'NATIONAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NATURAL)
< error
| 1064 42000 …near 'NATURAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NCHAR)
< error
| 1064 42000 …near 'NCHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NDB)
< error
| 1064 42000 …near 'NDB)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NDBCLUSTER)
< error
| 1064 42000 …near 'NDBCLUSTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NESTED)
< error
| 1064 42000 …near 'NESTED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NETWORK_NAMESPACE)
< error
| 1115 42000 Unknown character set: 'NETWORK_NAMESPACE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NEVER)
< error
| 1064 42000 …near 'NEVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NEW)
< error
| 1064 42000 …near 'NEW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NEXT)
< error
| 1064 42000 …near 'NEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NO)
< error
| 1064 42000 …near 'NO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NODEGROUP)
< error
| 1064 42000 …near 'NODEGROUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NONE)
< error
| 1064 42000 …near 'NONE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NOT)
< error
| 1064 42000 …near 'NOT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NOW)
< error
| 1064 42000 …near 'NOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NOWAIT)
< error
| 1064 42000 …near 'NOWAIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NO_WAIT)
< error
| 1064 42000 …near 'NO_WAIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NO_WRITE_TO_BINLOG)
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NTH_VALUE)
< error
| 1064 42000 …near 'NTH_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NTILE)
< error
| 1064 42000 …near 'NTILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NULLIF)
< error
| 1115 42000 Unknown character set: 'NULLIF'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NULLS)
< error
| 1064 42000 …near 'NULLS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NUMBER)
< error
| 1064 42000 …near 'NUMBER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NUMERIC)
< error
| 1064 42000 …near 'NUMERIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET NVARCHAR)
< error
| 1064 42000 …near 'NVARCHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OF)
< error
| 1064 42000 …near 'OF)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OFF)
< error
| 1064 42000 …near 'OFF)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OFFSET)
< error
| 1064 42000 …near 'OFFSET)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OJ)
< error
| 1064 42000 …near 'OJ)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OLD)
< error
| 1064 42000 …near 'OLD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ON)
< error
| 1064 42000 …near 'ON)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ONE)
< error
| 1064 42000 …near 'ONE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ONLY)
< error
| 1064 42000 …near 'ONLY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPEN)
< error
| 1064 42000 …near 'OPEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPTIMIZE)
< error
| 1064 42000 …near 'OPTIMIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPTIMIZER_COSTS)
< error
| 1064 42000 …near 'OPTIMIZER_COSTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPTION)
< error
| 1064 42000 …near 'OPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPTIONAL)
< error
| 1115 42000 Unknown character set: 'OPTIONAL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPTIONALLY)
< error
| 1064 42000 …near 'OPTIONALLY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OPTIONS)
< error
| 1064 42000 …near 'OPTIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OR)
< error
| 1064 42000 …near 'OR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ORDER)
< error
| 1064 42000 …near 'ORDER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ORDINALITY)
< error
| 1064 42000 …near 'ORDINALITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ORGANIZATION)
< error
| 1064 42000 …near 'ORGANIZATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OTHERS)
< error
| 1115 42000 Unknown character set: 'OTHERS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OUT)
< error
| 1064 42000 …near 'OUT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OUTER)
< error
| 1064 42000 …near 'OUTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OUTFILE)
< error
| 1064 42000 …near 'OUTFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OVER)
< error
| 1064 42000 …near 'OVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OVERLAY)
< error
| 1115 42000 Unknown character set: 'OVERLAY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OVERWRITE)
< error
| 1064 42000 …near 'OVERWRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET OWNER)
< error
| 1064 42000 …near 'OWNER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PACK_KEYS)
< error
| 1064 42000 …near 'PACK_KEYS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PAGE)
< error
| 1064 42000 …near 'PAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARALLEL)
< error
| 1064 42000 …near 'PARALLEL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARAMETERS)
< error
| 1064 42000 …near 'PARAMETERS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARSER)
< error
| 1064 42000 …near 'PARSER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARSE_TREE)
< error
| 1115 42000 Unknown character set: 'PARSE_TREE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARTIAL)
< error
| 1064 42000 …near 'PARTIAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARTITION)
< error
| 1064 42000 …near 'PARTITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARTITIONING)
< error
| 1064 42000 …near 'PARTITIONING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PARTITIONS)
< error
| 1064 42000 …near 'PARTITIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PASSWORD)
< error
| 1064 42000 …near 'PASSWORD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PASSWORD_LOCK_TIME)
< error
| 1115 42000 Unknown character set: 'PASSWORD_LOCK_TIME'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PATH)
< error
| 1064 42000 …near 'PATH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PERCENT_RANK)
< error
| 1064 42000 …near 'PERCENT_RANK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PERSIST)
< error
| 1115 42000 Unknown character set: 'PERSIST'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PERSIST_ONLY)
< error
| 1115 42000 Unknown character set: 'PERSIST_ONLY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PHASE)
< error
| 1064 42000 …near 'PHASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PLACING)
< error
| 1115 42000 Unknown character set: 'PLACING'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PLUGIN)
< error
| 1064 42000 …near 'PLUGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PLUGINS)
< error
| 1064 42000 …near 'PLUGINS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PLUGIN_DIR)
< error
| 1064 42000 …near 'PLUGIN_DIR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET POINT)
< error
| 1064 42000 …near 'POINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET POLYGON)
< error
| 1064 42000 …near 'POLYGON)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PORT)
< error
| 1064 42000 …near 'PORT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET POSITION)
< error
| 1064 42000 …near 'POSITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRECEDES)
< error
| 1115 42000 Unknown character set: 'PRECEDES'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRECEDING)
< error
| 1064 42000 …near 'PRECEDING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRECISION)
< error
| 1064 42000 …near 'PRECISION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PREFIX)
< error
| 1115 42000 Unknown character set: 'PREFIX'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PREPARE)
< error
| 1064 42000 …near 'PREPARE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRESERVE)
< error
| 1064 42000 …near 'PRESERVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PREV)
< error
| 1064 42000 …near 'PREV)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRIMARY)
< error
| 1064 42000 …near 'PRIMARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRIVILEGES)
< error
| 1064 42000 …near 'PRIVILEGES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PRIVILEGE_CHECKS_USER)
< error
| 1115 42000 Unknown character set: 'PRIVILEGE_CHECKS_USER'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PROCEDURE)
< error
| 1064 42000 …near 'PROCEDURE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PROCESS)
< error
| 1064 42000 …near 'PROCESS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PROCESSLIST)
< error
| 1064 42000 …near 'PROCESSLIST)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PROFILE)
< error
| 1064 42000 …near 'PROFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PROFILES)
< error
| 1064 42000 …near 'PROFILES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PROXY)
< error
| 1064 42000 …near 'PROXY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET PURGE)
< error
| 1064 42000 …near 'PURGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET QUALIFY)
< error
| 1115 42000 Unknown character set: 'QUALIFY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET QUARTER)
< error
| 1064 42000 …near 'QUARTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET QUERY)
< error
| 1064 42000 …near 'QUERY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET QUICK)
< error
| 1064 42000 …near 'QUICK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RANDOM)
< error
| 1115 42000 Unknown character set: 'RANDOM'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RANGE)
< error
| 1064 42000 …near 'RANGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RANK)
< error
| 1064 42000 …near 'RANK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET READ)
< error
| 1064 42000 …near 'READ)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET READS)
< error
| 1064 42000 …near 'READS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET READ_ONLY)
< error
| 1064 42000 …near 'READ_ONLY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET READ_WRITE)
< error
| 1064 42000 …near 'READ_WRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REAL)
< error
| 1064 42000 …near 'REAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REBUILD)
< error
| 1064 42000 …near 'REBUILD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RECOVER)
< error
| 1064 42000 …near 'RECOVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RECURSIVE)
< error
| 1064 42000 …near 'RECURSIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REDO_BUFFER_SIZE)
< error
| 1064 42000 …near 'REDO_BUFFER_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REDUNDANT)
< error
| 1064 42000 …near 'REDUNDANT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REFERENCE)
< error
| 1115 42000 Unknown character set: 'REFERENCE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REFERENCES)
< error
| 1064 42000 …near 'REFERENCES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REGEXP)
< error
| 1064 42000 …near 'REGEXP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REGISTRATION)
< error
| 1115 42000 Unknown character set: 'REGISTRATION'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELATIONAL)
< error
| 1115 42000 Unknown character set: 'RELATIONAL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELAY)
< error
| 1064 42000 …near 'RELAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELAYLOG)
< error
| 1064 42000 …near 'RELAYLOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELAY_LOG_FILE)
< error
| 1064 42000 …near 'RELAY_LOG_FILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELAY_LOG_POS)
< error
| 1064 42000 …near 'RELAY_LOG_POS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELAY_THREAD)
< error
| 1064 42000 …near 'RELAY_THREAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELEASE)
< error
| 1064 42000 …near 'RELEASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RELOAD)
< error
| 1064 42000 …near 'RELOAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REMOVE)
< error
| 1064 42000 …near 'REMOVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RENAME)
< error
| 1064 42000 …near 'RENAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REORGANIZE)
< error
| 1064 42000 …near 'REORGANIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPAIR)
< error
| 1064 42000 …near 'REPAIR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPEAT)
< error
| 1064 42000 …near 'REPEAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPEATABLE)
< error
| 1064 42000 …near 'REPEATABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLACE)
< error
| 1064 42000 …near 'REPLACE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICA)
< error
| 1064 42000 …near 'REPLICA)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICAS)
< error
| 1115 42000 Unknown character set: 'REPLICAS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_DO_DB)
< error
| 1115 42000 Unknown character set: 'REPLICATE_DO_DB'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_DO_TABLE)
< error
| 1115 42000 Unknown character set: 'REPLICATE_DO_TABLE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_IGNORE_DB)
< error
| 1115 42000 Unknown character set: 'REPLICATE_IGNORE_DB'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_IGNORE_TABLE)
< error
| 1115 42000 Unknown character set: 'REPLICATE_IGNORE_TABLE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_REWRITE_DB)
< error
| 1115 42000 Unknown character set: 'REPLICATE_REWRITE_DB'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_WILD_DO_TABLE)
< error
| 1115 42000 Unknown character set: 'REPLICATE_WILD_DO_TABLE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATE_WILD_IGNORE_TABLE)
< error
| 1115 42000 Unknown character set: 'REPLICATE_WILD_IGNORE_TABLE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REPLICATION)
< error
| 1064 42000 …near 'REPLICATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REQUIRE)
< error
| 1064 42000 …near 'REQUIRE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REQUIRE_ROW_FORMAT)
< error
| 1115 42000 Unknown character set: 'REQUIRE_ROW_FORMAT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REQUIRE_TABLE_PRIMARY_KEY_CHECK)
< error
| 1115 42000 Unknown character set: 'REQUIRE_TABLE_PRIMARY_KEY_CHECK'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESET)
< error
| 1064 42000 …near 'RESET)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESIGNAL)
< error
| 1064 42000 …near 'RESIGNAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESOURCE)
< error
| 1064 42000 …near 'RESOURCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESPECT)
< error
| 1064 42000 …near 'RESPECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESTART)
< error
| 1064 42000 …near 'RESTART)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESTORE)
< error
| 1064 42000 …near 'RESTORE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESTRICT)
< error
| 1064 42000 …near 'RESTRICT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RESUME)
< error
| 1064 42000 …near 'RESUME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RETAIN)
< error
| 1115 42000 Unknown character set: 'RETAIN'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RETURN)
< error
| 1064 42000 …near 'RETURN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RETURNED_SQLSTATE)
< error
| 1064 42000 …near 'RETURNED_SQLSTATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RETURNING)
< error
| 1064 42000 …near 'RETURNING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RETURNS)
< error
| 1064 42000 …near 'RETURNS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REUSE)
< error
| 1115 42000 Unknown character set: 'REUSE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REVERSE)
< error
| 1064 42000 …near 'REVERSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET REVOKE)
< error
| 1064 42000 …near 'REVOKE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RIGHT)
< error
| 1064 42000 …near 'RIGHT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RLIKE)
< error
| 1064 42000 …near 'RLIKE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROLE)
< error
| 1064 42000 …near 'ROLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROLLBACK)
< error
| 1064 42000 …near 'ROLLBACK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROLLUP)
< error
| 1064 42000 …near 'ROLLUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROTATE)
< error
| 1064 42000 …near 'ROTATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROUTINE)
< error
| 1064 42000 …near 'ROUTINE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROW)
< error
| 1064 42000 …near 'ROW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROWS)
< error
| 1064 42000 …near 'ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROW_COUNT)
< error
| 1064 42000 …near 'ROW_COUNT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROW_FORMAT)
< error
| 1064 42000 …near 'ROW_FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ROW_NUMBER)
< error
| 1064 42000 …near 'ROW_NUMBER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET RTREE)
< error
| 1064 42000 …near 'RTREE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET S3)
< error
| 1115 42000 Unknown character set: 'S3'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SAVEPOINT)
< error
| 1064 42000 …near 'SAVEPOINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SCHEDULE)
< error
| 1064 42000 …near 'SCHEDULE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SCHEMA)
< error
| 1064 42000 …near 'SCHEMA)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SCHEMAS)
< error
| 1064 42000 …near 'SCHEMAS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SCHEMA_NAME)
< error
| 1064 42000 …near 'SCHEMA_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECOND)
< error
| 1064 42000 …near 'SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECONDARY)
< error
| 1115 42000 Unknown character set: 'SECONDARY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECONDARY_ENGINE)
< error
| 1115 42000 Unknown character set: 'SECONDARY_ENGINE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECONDARY_ENGINE_ATTRIBUTE)
< error
| 1115 42000 Unknown character set: 'SECONDARY_ENGINE_ATTRIBUTE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECONDARY_LOAD)
< error
| 1115 42000 Unknown character set: 'SECONDARY_LOAD'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECONDARY_UNLOAD)
< error
| 1115 42000 Unknown character set: 'SECONDARY_UNLOAD'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECOND_MICROSECOND)
< error
| 1064 42000 …near 'SECOND_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SECURITY)
< error
| 1064 42000 …near 'SECURITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SELECT)
< error
| 1064 42000 …near 'SELECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SENSITIVE)
< error
| 1064 42000 …near 'SENSITIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SEPARATOR)
< error
| 1064 42000 …near 'SEPARATOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SERIAL)
< error
| 1064 42000 …near 'SERIAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SERIALIZABLE)
< error
| 1064 42000 …near 'SERIALIZABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SERVER)
< error
| 1064 42000 …near 'SERVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SESSION)
< error
| 1064 42000 …near 'SESSION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SESSION_USER)
< error
| 1064 42000 …near 'SESSION_USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SET)
< error
| 1064 42000 …near 'SET)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SETS)
< error
| 1064 42000 …near 'SETS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SHARE)
< error
| 1064 42000 …near 'SHARE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SHOW)
< error
| 1064 42000 …near 'SHOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SHUTDOWN)
< error
| 1064 42000 …near 'SHUTDOWN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SIGNAL)
< error
| 1064 42000 …near 'SIGNAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SIMPLE)
< error
| 1064 42000 …near 'SIMPLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SKIP)
< error
| 1064 42000 …near 'SKIP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SLAVE)
< error
| 1064 42000 …near 'SLAVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SLOW)
< error
| 1064 42000 …near 'SLOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SMALLINT)
< error
| 1064 42000 …near 'SMALLINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SNAPSHOT)
< error
| 1064 42000 …near 'SNAPSHOT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOCKET)
< error
| 1064 42000 …near 'SOCKET)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOME)
< error
| 1064 42000 …near 'SOME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SONAME)
< error
| 1064 42000 …near 'SONAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOUNDS)
< error
| 1064 42000 …near 'SOUNDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE)
< error
| 1064 42000 …near 'SOURCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_AUTO_POSITION)
< error
| 1115 42000 Unknown character set: 'SOURCE_AUTO_POSITION'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_BIND)
< error
| 1115 42000 Unknown character set: 'SOURCE_BIND'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_COMPRESSION_ALGORITHMS)
< error
| 1115 42000 Unknown character set: 'SOURCE_COMPRESSION_ALGORITHMS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_CONNECTION_AUTO_FAILOVER)
< error
| 1115 42000 Unknown character set: 'SOURCE_CONNECTION_AUTO_FAILOVER'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_CONNECT_RETRY)
< error
| 1115 42000 Unknown character set: 'SOURCE_CONNECT_RETRY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_DELAY)
< error
| 1115 42000 Unknown character set: 'SOURCE_DELAY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_HEARTBEAT_PERIOD)
< error
| 1115 42000 Unknown character set: 'SOURCE_HEARTBEAT_PERIOD'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_HOST)
< error
| 1115 42000 Unknown character set: 'SOURCE_HOST'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_LOG_FILE)
< error
| 1115 42000 Unknown character set: 'SOURCE_LOG_FILE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_LOG_POS)
< error
| 1115 42000 Unknown character set: 'SOURCE_LOG_POS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_PASSWORD)
< error
| 1115 42000 Unknown character set: 'SOURCE_PASSWORD'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_PORT)
< error
| 1115 42000 Unknown character set: 'SOURCE_PORT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_PUBLIC_KEY_PATH)
< error
| 1115 42000 Unknown character set: 'SOURCE_PUBLIC_KEY_PATH'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_RETRY_COUNT)
< error
| 1115 42000 Unknown character set: 'SOURCE_RETRY_COUNT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_CA)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_CA'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_CAPATH)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_CAPATH'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_CERT)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_CERT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_CIPHER)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_CIPHER'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_CRL)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_CRL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_CRLPATH)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_CRLPATH'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_KEY)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_KEY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_SSL_VERIFY_SERVER_CERT)
< error
| 1115 42000 Unknown character set: 'SOURCE_SSL_VERIFY_SERVER_CERT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_TLS_CIPHERSUITES)
< error
| 1115 42000 Unknown character set: 'SOURCE_TLS_CIPHERSUITES'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_TLS_VERSION)
< error
| 1115 42000 Unknown character set: 'SOURCE_TLS_VERSION'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_USER)
< error
| 1115 42000 Unknown character set: 'SOURCE_USER'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SOURCE_ZSTD_COMPRESSION_LEVEL)
< error
| 1115 42000 Unknown character set: 'SOURCE_ZSTD_COMPRESSION_LEVEL'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SPATIAL)
< error
| 1064 42000 …near 'SPATIAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SPECIFIC)
< error
| 1064 42000 …near 'SPECIFIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL)
< error
| 1064 42000 …near 'SQL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQLEXCEPTION)
< error
| 1064 42000 …near 'SQLEXCEPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQLSTATE)
< error
| 1064 42000 …near 'SQLSTATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQLWARNING)
< error
| 1064 42000 …near 'SQLWARNING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_AFTER_GTIDS)
< error
| 1064 42000 …near 'SQL_AFTER_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_AFTER_MTS_GAPS)
< error
| 1064 42000 …near 'SQL_AFTER_MTS_GAPS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_BEFORE_GTIDS)
< error
| 1064 42000 …near 'SQL_BEFORE_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_BIG_RESULT)
< error
| 1064 42000 …near 'SQL_BIG_RESULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_BUFFER_RESULT)
< error
| 1064 42000 …near 'SQL_BUFFER_RESULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_CALC_FOUND_ROWS)
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_NO_CACHE)
< error
| 1064 42000 …near 'SQL_NO_CACHE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_SMALL_RESULT)
< error
| 1064 42000 …near 'SQL_SMALL_RESULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_THREAD)
< error
| 1064 42000 …near 'SQL_THREAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_DAY)
< error
| 1064 42000 …near 'SQL_TSI_DAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_HOUR)
< error
| 1064 42000 …near 'SQL_TSI_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_MINUTE)
< error
| 1064 42000 …near 'SQL_TSI_MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_MONTH)
< error
| 1064 42000 …near 'SQL_TSI_MONTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_QUARTER)
< error
| 1064 42000 …near 'SQL_TSI_QUARTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_SECOND)
< error
| 1064 42000 …near 'SQL_TSI_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_WEEK)
< error
| 1064 42000 …near 'SQL_TSI_WEEK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SQL_TSI_YEAR)
< error
| 1064 42000 …near 'SQL_TSI_YEAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SRID)
< error
| 1064 42000 …near 'SRID)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SSE_KMS)
< error
| 1115 42000 Unknown character set: 'SSE_KMS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SSE_S3)
< error
| 1115 42000 Unknown character set: 'SSE_S3'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SSL)
< error
| 1064 42000 …near 'SSL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STACKED)
< error
| 1064 42000 …near 'STACKED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET START)
< error
| 1064 42000 …near 'START)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STARTING)
< error
| 1064 42000 …near 'STARTING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STARTS)
< error
| 1064 42000 …near 'STARTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STATS_AUTO_RECALC)
< error
| 1064 42000 …near 'STATS_AUTO_RECALC)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STATS_PERSISTENT)
< error
| 1064 42000 …near 'STATS_PERSISTENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STATS_SAMPLE_PAGES)
< error
| 1064 42000 …near 'STATS_SAMPLE_PAGES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STATUS)
< error
| 1064 42000 …near 'STATUS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STD)
< error
| 1064 42000 …near 'STD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STDDEV)
< error
| 1064 42000 …near 'STDDEV)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STDDEV_POP)
< error
| 1064 42000 …near 'STDDEV_POP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STDDEV_SAMP)
< error
| 1064 42000 …near 'STDDEV_SAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STOP)
< error
| 1064 42000 …near 'STOP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STORAGE)
< error
| 1064 42000 …near 'STORAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STORED)
< error
| 1064 42000 …near 'STORED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STRAIGHT_JOIN)
< error
| 1064 42000 …near 'STRAIGHT_JOIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STREAM)
< error
| 1115 42000 Unknown character set: 'STREAM'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STRICT_LOAD)
< error
| 1115 42000 Unknown character set: 'STRICT_LOAD'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET STRING)
< error
| 1064 42000 …near 'STRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ST_COLLECT)
< error
| 1115 42000 Unknown character set: 'ST_COLLECT'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBCLASS_ORIGIN)
< error
| 1064 42000 …near 'SUBCLASS_ORIGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBDATE)
< error
| 1064 42000 …near 'SUBDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBJECT)
< error
| 1064 42000 …near 'SUBJECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBPARTITION)
< error
| 1064 42000 …near 'SUBPARTITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBPARTITIONS)
< error
| 1064 42000 …near 'SUBPARTITIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBSTR)
< error
| 1064 42000 …near 'SUBSTR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUBSTRING)
< error
| 1064 42000 …near 'SUBSTRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUM)
< error
| 1064 42000 …near 'SUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUPER)
< error
| 1064 42000 …near 'SUPER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SUSPEND)
< error
| 1064 42000 …near 'SUSPEND)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SWAPS)
< error
| 1064 42000 …near 'SWAPS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SWITCHES)
< error
| 1064 42000 …near 'SWITCHES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SYSDATE)
< error
| 1064 42000 …near 'SYSDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SYSTEM)
< error
| 1064 42000 …near 'SYSTEM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET SYSTEM_USER)
< error
| 1064 42000 …near 'SYSTEM_USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TABLE)
< error
| 1064 42000 …near 'TABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TABLES)
< error
| 1064 42000 …near 'TABLES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TABLESAMPLE)
< error
| 1115 42000 Unknown character set: 'TABLESAMPLE'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TABLESPACE)
< error
| 1064 42000 …near 'TABLESPACE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TABLE_CHECKSUM)
< error
| 1064 42000 …near 'TABLE_CHECKSUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TABLE_NAME)
< error
| 1064 42000 …near 'TABLE_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TEMPORARY)
< error
| 1064 42000 …near 'TEMPORARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TEMPTABLE)
< error
| 1064 42000 …near 'TEMPTABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TERMINATED)
< error
| 1064 42000 …near 'TERMINATED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET THAN)
< error
| 1064 42000 …near 'THAN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET THEN)
< error
| 1064 42000 …near 'THEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET THREAD_PRIORITY)
< error
| 1115 42000 Unknown character set: 'THREAD_PRIORITY'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TIES)
< error
| 1115 42000 Unknown character set: 'TIES'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TIME)
< error
| 1064 42000 …near 'TIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TIMESTAMP)
< error
| 1064 42000 …near 'TIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TIMESTAMPADD)
< error
| 1064 42000 …near 'TIMESTAMPADD)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TIMESTAMPDIFF)
< error
| 1064 42000 …near 'TIMESTAMPDIFF)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TINYBLOB)
< error
| 1064 42000 …near 'TINYBLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TINYINT)
< error
| 1064 42000 …near 'TINYINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TINYTEXT)
< error
| 1064 42000 …near 'TINYTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TLS)
< error
| 1115 42000 Unknown character set: 'TLS'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TO)
< error
| 1064 42000 …near 'TO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRAILING)
< error
| 1064 42000 …near 'TRAILING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRANSACTION)
< error
| 1064 42000 …near 'TRANSACTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRIGGER)
< error
| 1064 42000 …near 'TRIGGER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRIGGERS)
< error
| 1064 42000 …near 'TRIGGERS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRIM)
< error
| 1064 42000 …near 'TRIM)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRUE)
< error
| 1064 42000 …near 'TRUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TRUNCATE)
< error
| 1064 42000 …near 'TRUNCATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TYPE)
< error
| 1064 42000 …near 'TYPE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET TYPES)
< error
| 1064 42000 …near 'TYPES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNBOUNDED)
< error
| 1064 42000 …near 'UNBOUNDED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNCOMMITTED)
< error
| 1064 42000 …near 'UNCOMMITTED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNDEFINED)
< error
| 1064 42000 …near 'UNDEFINED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNDO)
< error
| 1064 42000 …near 'UNDO)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNDOFILE)
< error
| 1064 42000 …near 'UNDOFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNDO_BUFFER_SIZE)
< error
| 1064 42000 …near 'UNDO_BUFFER_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNINSTALL)
< error
| 1064 42000 …near 'UNINSTALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNION)
< error
| 1064 42000 …near 'UNION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNIQUE)
< error
| 1064 42000 …near 'UNIQUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNKNOWN)
< error
| 1064 42000 …near 'UNKNOWN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNLOCK)
< error
| 1064 42000 …near 'UNLOCK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNREGISTER)
< error
| 1115 42000 Unknown character set: 'UNREGISTER'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UNTIL)
< error
| 1064 42000 …near 'UNTIL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UPDATE)
< error
| 1064 42000 …near 'UPDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UPGRADE)
< error
| 1064 42000 …near 'UPGRADE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET URI)
< error
| 1115 42000 Unknown character set: 'URI'
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET USAGE)
< error
| 1064 42000 …near 'USAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET USE)
< error
| 1064 42000 …near 'USE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET USER)
< error
| 1064 42000 …near 'USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET USING)
< error
| 1064 42000 …near 'USING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UTC_DATE)
< error
| 1064 42000 …near 'UTC_DATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UTC_TIME)
< error
| 1064 42000 …near 'UTC_TIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET UTC_TIMESTAMP)
< error
| 1064 42000 …near 'UTC_TIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VALIDATION)
< error
| 1064 42000 …near 'VALIDATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VALUE)
< error
| 1064 42000 …near 'VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VALUES)
< error
| 1064 42000 …near 'VALUES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VARBINARY)
< error
| 1064 42000 …near 'VARBINARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VARCHAR)
< error
| 1064 42000 …near 'VARCHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VARIABLES)
< error
| 1064 42000 …near 'VARIABLES)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VARYING)
< error
| 1064 42000 …near 'VARYING)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VIEW)
< error
| 1064 42000 …near 'VIEW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VIRTUAL)
< error
| 1064 42000 …near 'VIRTUAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET VISIBLE)
< error
| 1064 42000 …near 'VISIBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WAIT)
< error
| 1064 42000 …near 'WAIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WARNINGS)
< error
| 1064 42000 …near 'WARNINGS)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WHEN)
< error
| 1064 42000 …near 'WHEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WHERE)
< error
| 1064 42000 …near 'WHERE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WHILE)
< error
| 1064 42000 …near 'WHILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WINDOW)
< error
| 1064 42000 …near 'WINDOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WITH)
< error
| 1064 42000 …near 'WITH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WITHOUT)
< error
| 1064 42000 …near 'WITHOUT)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WORK)
< error
| 1064 42000 …near 'WORK)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WRAPPER)
< error
| 1064 42000 …near 'WRAPPER)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET WRITE)
< error
| 1064 42000 …near 'WRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET X509)
< error
| 1064 42000 …near 'X509)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET XA)
< error
| 1064 42000 …near 'XA)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET XID)
< error
| 1064 42000 …near 'XID)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET XML)
< error
| 1064 42000 …near 'XML)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET XOR)
< error
| 1064 42000 …near 'XOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET YEAR)
< error
| 1064 42000 …near 'YEAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET YEAR_MONTH)
< error
| 1064 42000 …near 'YEAR_MONTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_cs_test (a VARCHAR(50) CHARACTER SET ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ABSENT)
< error
| 1273 HY000 Unknown collation: 'ABSENT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ACCESSIBLE)
< error
| 1064 42000 …near 'ACCESSIBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ACCOUNT)
< error
| 1064 42000 …near 'ACCOUNT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ACTION)
< error
| 1064 42000 …near 'ACTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ACTIVE)
< error
| 1064 42000 …near 'ACTIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ADD)
< error
| 1064 42000 …near 'ADD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ADDDATE)
< error
| 1064 42000 …near 'ADDDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ADMIN)
< error
| 1064 42000 …near 'ADMIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AFTER)
< error
| 1064 42000 …near 'AFTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AGAINST)
< error
| 1064 42000 …near 'AGAINST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AGGREGATE)
< error
| 1064 42000 …near 'AGGREGATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ALGORITHM)
< error
| 1064 42000 …near 'ALGORITHM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ALL)
< error
| 1064 42000 …near 'ALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ALLOW_MISSING_FILES)
< error
| 1273 HY000 Unknown collation: 'ALLOW_MISSING_FILES'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ALTER)
< error
| 1064 42000 …near 'ALTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ALWAYS)
< error
| 1064 42000 …near 'ALWAYS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ANALYZE)
< error
| 1064 42000 …near 'ANALYZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AND)
< error
| 1064 42000 …near 'AND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ANY)
< error
| 1064 42000 …near 'ANY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AS)
< error
| 1064 42000 …near 'AS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ASC)
< error
| 1064 42000 …near 'ASC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ASENSITIVE)
< error
| 1064 42000 …near 'ASENSITIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS)
< error
| 1273 HY000 Unknown collation: 'ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AT)
< error
| 1064 42000 …near 'AT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ATTRIBUTE)
< error
| 1064 42000 …near 'ATTRIBUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AUTHENTICATION)
< error
| 1273 HY000 Unknown collation: 'AUTHENTICATION'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AUTO)
< error
| 1064 42000 …near 'AUTO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AUTOEXTEND_SIZE)
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AUTO_INCREMENT)
< error
| 1064 42000 …near 'AUTO_INCREMENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AUTO_REFRESH)
< error
| 1064 42000 …near 'AUTO_REFRESH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AUTO_REFRESH_SOURCE)
< error
| 1273 HY000 Unknown collation: 'AUTO_REFRESH_SOURCE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AVG)
< error
| 1064 42000 …near 'AVG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE AVG_ROW_LENGTH)
< error
| 1064 42000 …near 'AVG_ROW_LENGTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BACKUP)
< error
| 1064 42000 …near 'BACKUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BEFORE)
< error
| 1064 42000 …near 'BEFORE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BEGIN)
< error
| 1064 42000 …near 'BEGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BERNOULLI)
< error
| 1273 HY000 Unknown collation: 'BERNOULLI'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BETWEEN)
< error
| 1064 42000 …near 'BETWEEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BIGINT)
< error
| 1064 42000 …near 'BIGINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BINLOG)
< error
| 1064 42000 …near 'BINLOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BIT)
< error
| 1064 42000 …near 'BIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BIT_AND)
< error
| 1064 42000 …near 'BIT_AND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BIT_OR)
< error
| 1064 42000 …near 'BIT_OR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BIT_XOR)
< error
| 1064 42000 …near 'BIT_XOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BLOB)
< error
| 1064 42000 …near 'BLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BLOCK)
< error
| 1064 42000 …near 'BLOCK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BOOL)
< error
| 1064 42000 …near 'BOOL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BOOLEAN)
< error
| 1064 42000 …near 'BOOLEAN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BOTH)
< error
| 1064 42000 …near 'BOTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BTREE)
< error
| 1064 42000 …near 'BTREE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BUCKETS)
< error
| 1064 42000 …near 'BUCKETS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BULK)
< error
| 1273 HY000 Unknown collation: 'BULK'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BY)
< error
| 1064 42000 …near 'BY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE BYTE)
< error
| 1064 42000 …near 'BYTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CACHE)
< error
| 1064 42000 …near 'CACHE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CALL)
< error
| 1064 42000 …near 'CALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CASCADE)
< error
| 1064 42000 …near 'CASCADE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CASCADED)
< error
| 1064 42000 …near 'CASCADED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CASE)
< error
| 1064 42000 …near 'CASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CAST)
< error
| 1064 42000 …near 'CAST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CATALOG_NAME)
< error
| 1064 42000 …near 'CATALOG_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHAIN)
< error
| 1064 42000 …near 'CHAIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHALLENGE_RESPONSE)
< error
| 1273 HY000 Unknown collation: 'CHALLENGE_RESPONSE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHANGE)
< error
| 1064 42000 …near 'CHANGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHANGED)
< error
| 1064 42000 …near 'CHANGED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHANNEL)
< error
| 1273 HY000 Unknown collation: 'CHANNEL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHAR)
< error
| 1064 42000 …near 'CHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHARACTER)
< error
| 1064 42000 …near 'CHARACTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHARSET)
< error
| 1064 42000 …near 'CHARSET)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHECK)
< error
| 1064 42000 …near 'CHECK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CHECKSUM)
< error
| 1064 42000 …near 'CHECKSUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CIPHER)
< error
| 1064 42000 …near 'CIPHER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CLASS_ORIGIN)
< error
| 1064 42000 …near 'CLASS_ORIGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CLIENT)
< error
| 1064 42000 …near 'CLIENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CLONE)
< error
| 1064 42000 …near 'CLONE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CLOSE)
< error
| 1064 42000 …near 'CLOSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COALESCE)
< error
| 1064 42000 …near 'COALESCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CODE)
< error
| 1064 42000 …near 'CODE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COLLATE)
< error
| 1064 42000 …near 'COLLATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COLLATION)
< error
| 1064 42000 …near 'COLLATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COLUMN)
< error
| 1064 42000 …near 'COLUMN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COLUMNS)
< error
| 1064 42000 …near 'COLUMNS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COLUMN_FORMAT)
< error
| 1064 42000 …near 'COLUMN_FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COLUMN_NAME)
< error
| 1064 42000 …near 'COLUMN_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMMENT)
< error
| 1064 42000 …near 'COMMENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMMIT)
< error
| 1064 42000 …near 'COMMIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMMITTED)
< error
| 1064 42000 …near 'COMMITTED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMPACT)
< error
| 1064 42000 …near 'COMPACT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMPLETION)
< error
| 1064 42000 …near 'COMPLETION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMPONENT)
< error
| 1273 HY000 Unknown collation: 'COMPONENT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMPRESSED)
< error
| 1064 42000 …near 'COMPRESSED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COMPRESSION)
< error
| 1064 42000 …near 'COMPRESSION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONCAT)
< error
| 1273 HY000 Unknown collation: 'CONCAT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONCURRENT)
< error
| 1064 42000 …near 'CONCURRENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONDITION)
< error
| 1064 42000 …near 'CONDITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONNECTION)
< error
| 1064 42000 …near 'CONNECTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONSISTENT)
< error
| 1064 42000 …near 'CONSISTENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONSTRAINT)
< error
| 1064 42000 …near 'CONSTRAINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONSTRAINT_CATALOG)
< error
| 1064 42000 …near 'CONSTRAINT_CATALOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONSTRAINT_NAME)
< error
| 1064 42000 …near 'CONSTRAINT_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONSTRAINT_SCHEMA)
< error
| 1064 42000 …near 'CONSTRAINT_SCHEMA)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONTAINS)
< error
| 1064 42000 …near 'CONTAINS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONTEXT)
< error
| 1064 42000 …near 'CONTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONTINUE)
< error
| 1064 42000 …near 'CONTINUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CONVERT)
< error
| 1064 42000 …near 'CONVERT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE COUNT)
< error
| 1064 42000 …near 'COUNT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CPU)
< error
| 1064 42000 …near 'CPU)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CREATE)
< error
| 1064 42000 …near 'CREATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CROSS)
< error
| 1064 42000 …near 'CROSS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CSV)
< error
| 1273 HY000 Unknown collation: 'CSV'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CUBE)
< error
| 1064 42000 …near 'CUBE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CUME_DIST)
< error
| 1064 42000 …near 'CUME_DIST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURDATE)
< error
| 1064 42000 …near 'CURDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURRENT)
< error
| 1064 42000 …near 'CURRENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURRENT_DATE)
< error
| 1064 42000 …near 'CURRENT_DATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURRENT_TIME)
< error
| 1064 42000 …near 'CURRENT_TIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURRENT_TIMESTAMP)
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURRENT_USER)
< error
| 1064 42000 …near 'CURRENT_USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURSOR)
< error
| 1064 42000 …near 'CURSOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURSOR_NAME)
< error
| 1064 42000 …near 'CURSOR_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE CURTIME)
< error
| 1064 42000 …near 'CURTIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATA)
< error
| 1064 42000 …near 'DATA)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATABASE)
< error
| 1064 42000 …near 'DATABASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATABASES)
< error
| 1064 42000 …near 'DATABASES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATAFILE)
< error
| 1064 42000 …near 'DATAFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATE)
< error
| 1064 42000 …near 'DATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATETIME)
< error
| 1064 42000 …near 'DATETIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATE_ADD)
< error
| 1064 42000 …near 'DATE_ADD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DATE_SUB)
< error
| 1064 42000 …near 'DATE_SUB)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DAY)
< error
| 1064 42000 …near 'DAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DAY_HOUR)
< error
| 1064 42000 …near 'DAY_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DAY_MICROSECOND)
< error
| 1064 42000 …near 'DAY_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DAY_MINUTE)
< error
| 1064 42000 …near 'DAY_MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DAY_SECOND)
< error
| 1064 42000 …near 'DAY_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DEALLOCATE)
< error
| 1064 42000 …near 'DEALLOCATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DEC)
< error
| 1064 42000 …near 'DEC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DECIMAL)
< error
| 1064 42000 …near 'DECIMAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DECLARE)
< error
| 1064 42000 …near 'DECLARE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DEFAULT)
< error
| 1064 42000 …near 'DEFAULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DEFAULT_AUTH)
< error
| 1064 42000 …near 'DEFAULT_AUTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DEFINER)
< error
| 1064 42000 …near 'DEFINER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DEFINITION)
< error
| 1273 HY000 Unknown collation: 'DEFINITION'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DELAYED)
< error
| 1064 42000 …near 'DELAYED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DELAY_KEY_WRITE)
< error
| 1064 42000 …near 'DELAY_KEY_WRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DELETE)
< error
| 1064 42000 …near 'DELETE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DENSE_RANK)
< error
| 1064 42000 …near 'DENSE_RANK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DESC)
< error
| 1064 42000 …near 'DESC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DESCRIBE)
< error
| 1064 42000 …near 'DESCRIBE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DESCRIPTION)
< error
| 1064 42000 …near 'DESCRIPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DETERMINISTIC)
< error
| 1064 42000 …near 'DETERMINISTIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DIAGNOSTICS)
< error
| 1064 42000 …near 'DIAGNOSTICS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DIRECTORY)
< error
| 1064 42000 …near 'DIRECTORY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DISABLE)
< error
| 1064 42000 …near 'DISABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DISCARD)
< error
| 1064 42000 …near 'DISCARD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DISK)
< error
| 1064 42000 …near 'DISK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DISTINCT)
< error
| 1064 42000 …near 'DISTINCT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DISTINCTROW)
< error
| 1064 42000 …near 'DISTINCTROW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DIV)
< error
| 1064 42000 …near 'DIV)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DO)
< error
| 1064 42000 …near 'DO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DOUBLE)
< error
| 1064 42000 …near 'DOUBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DROP)
< error
| 1064 42000 …near 'DROP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DUAL)
< error
| 1064 42000 …near 'DUAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DUALITY)
< error
| 1273 HY000 Unknown collation: 'DUALITY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DUMPFILE)
< error
| 1064 42000 …near 'DUMPFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DUPLICATE)
< error
| 1064 42000 …near 'DUPLICATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE DYNAMIC)
< error
| 1064 42000 …near 'DYNAMIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EACH)
< error
| 1064 42000 …near 'EACH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ELSE)
< error
| 1064 42000 …near 'ELSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ELSEIF)
< error
| 1064 42000 …near 'ELSEIF)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EMPTY)
< error
| 1064 42000 …near 'EMPTY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENABLE)
< error
| 1064 42000 …near 'ENABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENCLOSED)
< error
| 1064 42000 …near 'ENCLOSED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENCRYPTION)
< error
| 1064 42000 …near 'ENCRYPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE END)
< error
| 1064 42000 …near 'END)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENDS)
< error
| 1064 42000 …near 'ENDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENFORCED)
< error
| 1064 42000 …near 'ENFORCED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENGINE)
< error
| 1064 42000 …near 'ENGINE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENGINES)
< error
| 1064 42000 …near 'ENGINES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENGINE_ATTRIBUTE)
< error
| 1273 HY000 Unknown collation: 'ENGINE_ATTRIBUTE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ENUM)
< error
| 1064 42000 …near 'ENUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ERROR)
< error
| 1064 42000 …near 'ERROR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ERRORS)
< error
| 1064 42000 …near 'ERRORS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ESCAPE)
< error
| 1064 42000 …near 'ESCAPE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ESCAPED)
< error
| 1064 42000 …near 'ESCAPED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EVENT)
< error
| 1064 42000 …near 'EVENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EVENTS)
< error
| 1064 42000 …near 'EVENTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EVERY)
< error
| 1064 42000 …near 'EVERY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXCEPT)
< error
| 1064 42000 …near 'EXCEPT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXCHANGE)
< error
| 1064 42000 …near 'EXCHANGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXCLUDE)
< error
| 1273 HY000 Unknown collation: 'EXCLUDE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXECUTE)
< error
| 1064 42000 …near 'EXECUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXISTS)
< error
| 1064 42000 …near 'EXISTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXIT)
< error
| 1064 42000 …near 'EXIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXPANSION)
< error
| 1064 42000 …near 'EXPANSION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXPIRE)
< error
| 1064 42000 …near 'EXPIRE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXPLAIN)
< error
| 1064 42000 …near 'EXPLAIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXPORT)
< error
| 1064 42000 …near 'EXPORT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXTENDED)
< error
| 1064 42000 …near 'EXTENDED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXTENT_SIZE)
< error
| 1064 42000 …near 'EXTENT_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXTERNAL)
< error
| 1064 42000 …near 'EXTERNAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXTERNAL_FORMAT)
< error
| 1273 HY000 Unknown collation: 'EXTERNAL_FORMAT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE EXTRACT)
< error
| 1064 42000 …near 'EXTRACT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FACTOR)
< error
| 1273 HY000 Unknown collation: 'FACTOR'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FAILED_LOGIN_ATTEMPTS)
< error
| 1273 HY000 Unknown collation: 'FAILED_LOGIN_ATTEMPTS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FALSE)
< error
| 1064 42000 …near 'FALSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FAST)
< error
| 1064 42000 …near 'FAST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FAULTS)
< error
| 1064 42000 …near 'FAULTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FETCH)
< error
| 1064 42000 …near 'FETCH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FIELDS)
< error
| 1064 42000 …near 'FIELDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILE)
< error
| 1064 42000 …near 'FILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILES)
< error
| 1064 42000 …near 'FILES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILE_BLOCK_SIZE)
< error
| 1064 42000 …near 'FILE_BLOCK_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILE_FORMAT)
< error
| 1273 HY000 Unknown collation: 'FILE_FORMAT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILE_NAME)
< error
| 1273 HY000 Unknown collation: 'FILE_NAME'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILE_PATTERN)
< error
| 1273 HY000 Unknown collation: 'FILE_PATTERN'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILE_PREFIX)
< error
| 1273 HY000 Unknown collation: 'FILE_PREFIX'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FILTER)
< error
| 1064 42000 …near 'FILTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FINISH)
< error
| 1273 HY000 Unknown collation: 'FINISH'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FIRST)
< error
| 1064 42000 …near 'FIRST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FIRST_VALUE)
< error
| 1064 42000 …near 'FIRST_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FIXED)
< error
| 1064 42000 …near 'FIXED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FLOAT)
< error
| 1064 42000 …near 'FLOAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FLOAT4)
< error
| 1064 42000 …near 'FLOAT4)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FLOAT8)
< error
| 1064 42000 …near 'FLOAT8)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FLUSH)
< error
| 1064 42000 …near 'FLUSH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FOLLOWING)
< error
| 1064 42000 …near 'FOLLOWING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FOLLOWS)
< error
| 1273 HY000 Unknown collation: 'FOLLOWS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FOR)
< error
| 1064 42000 …near 'FOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FORCE)
< error
| 1064 42000 …near 'FORCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FOREIGN)
< error
| 1064 42000 …near 'FOREIGN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FORMAT)
< error
| 1064 42000 …near 'FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FOUND)
< error
| 1064 42000 …near 'FOUND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FROM)
< error
| 1064 42000 …near 'FROM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FULL)
< error
| 1064 42000 …near 'FULL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FULLTEXT)
< error
| 1064 42000 …near 'FULLTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE FUNCTION)
< error
| 1064 42000 …near 'FUNCTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GENERAL)
< error
| 1064 42000 …near 'GENERAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GENERATE)
< error
| 1273 HY000 Unknown collation: 'GENERATE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GENERATED)
< error
| 1064 42000 …near 'GENERATED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GEOMCOLLECTION)
< error
| 1064 42000 …near 'GEOMCOLLECTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GEOMETRY)
< error
| 1064 42000 …near 'GEOMETRY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GEOMETRYCOLLECTION)
< error
| 1064 42000 …near 'GEOMETRYCOLLECTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GET)
< error
| 1064 42000 …near 'GET)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GET_FORMAT)
< error
| 1064 42000 …near 'GET_FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GET_SOURCE_PUBLIC_KEY)
< error
| 1273 HY000 Unknown collation: 'GET_SOURCE_PUBLIC_KEY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GLOBAL)
< error
| 1064 42000 …near 'GLOBAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GRANT)
< error
| 1064 42000 …near 'GRANT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GRANTS)
< error
| 1064 42000 …near 'GRANTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GREATEST)
< error
| 1273 HY000 Unknown collation: 'GREATEST'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GROUP)
< error
| 1064 42000 …near 'GROUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GROUPING)
< error
| 1064 42000 …near 'GROUPING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GROUPS)
< error
| 1273 HY000 Unknown collation: 'GROUPS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GROUP_CONCAT)
< error
| 1064 42000 …near 'GROUP_CONCAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GROUP_REPLICATION)
< error
| 1273 HY000 Unknown collation: 'GROUP_REPLICATION'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GTIDS)
< error
| 1273 HY000 Unknown collation: 'GTIDS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GTID_ONLY)
< error
| 1273 HY000 Unknown collation: 'GTID_ONLY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE GUIDED)
< error
| 1273 HY000 Unknown collation: 'GUIDED'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HANDLER)
< error
| 1064 42000 …near 'HANDLER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HASH)
< error
| 1064 42000 …near 'HASH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HAVING)
< error
| 1064 42000 …near 'HAVING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HEADER)
< error
| 1273 HY000 Unknown collation: 'HEADER'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HELP)
< error
| 1064 42000 …near 'HELP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HIGH_PRIORITY)
< error
| 1064 42000 …near 'HIGH_PRIORITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HISTOGRAM)
< error
| 1064 42000 …near 'HISTOGRAM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HISTORY)
< error
| 1273 HY000 Unknown collation: 'HISTORY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HOST)
< error
| 1064 42000 …near 'HOST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HOSTS)
< error
| 1064 42000 …near 'HOSTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HOUR)
< error
| 1064 42000 …near 'HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HOUR_MICROSECOND)
< error
| 1064 42000 …near 'HOUR_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HOUR_MINUTE)
< error
| 1064 42000 …near 'HOUR_MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE HOUR_SECOND)
< error
| 1064 42000 …near 'HOUR_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IDENTIFIED)
< error
| 1064 42000 …near 'IDENTIFIED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IF)
< error
| 1064 42000 …near 'IF)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IGNORE)
< error
| 1064 42000 …near 'IGNORE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IGNORE_SERVER_IDS)
< error
| 1064 42000 …near 'IGNORE_SERVER_IDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IMPORT)
< error
| 1064 42000 …near 'IMPORT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IN)
< error
| 1064 42000 …near 'IN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INACTIVE)
< error
| 1273 HY000 Unknown collation: 'INACTIVE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INDEX)
< error
| 1064 42000 …near 'INDEX)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INDEXES)
< error
| 1064 42000 …near 'INDEXES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INFILE)
< error
| 1064 42000 …near 'INFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INITIAL)
< error
| 1273 HY000 Unknown collation: 'INITIAL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INITIAL_SIZE)
< error
| 1064 42000 …near 'INITIAL_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INITIATE)
< error
| 1273 HY000 Unknown collation: 'INITIATE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INNER)
< error
| 1064 42000 …near 'INNER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INOUT)
< error
| 1064 42000 …near 'INOUT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INSENSITIVE)
< error
| 1064 42000 …near 'INSENSITIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INSERT)
< error
| 1064 42000 …near 'INSERT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INSERT_METHOD)
< error
| 1064 42000 …near 'INSERT_METHOD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INSTALL)
< error
| 1064 42000 …near 'INSTALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INSTANCE)
< error
| 1064 42000 …near 'INSTANCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INT1)
< error
| 1064 42000 …near 'INT1)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INT2)
< error
| 1064 42000 …near 'INT2)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INT3)
< error
| 1064 42000 …near 'INT3)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INT4)
< error
| 1064 42000 …near 'INT4)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INT8)
< error
| 1064 42000 …near 'INT8)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INTEGER)
< error
| 1064 42000 …near 'INTEGER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INTERSECT)
< error
| 1064 42000 …near 'INTERSECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INTERVAL)
< error
| 1064 42000 …near 'INTERVAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INTO)
< error
| 1064 42000 …near 'INTO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INVISIBLE)
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE INVOKER)
< error
| 1064 42000 …near 'INVOKER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IO)
< error
| 1064 42000 …near 'IO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IO_AFTER_GTIDS)
< error
| 1064 42000 …near 'IO_AFTER_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IO_BEFORE_GTIDS)
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IO_THREAD)
< error
| 1064 42000 …near 'IO_THREAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IPC)
< error
| 1064 42000 …near 'IPC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE IS)
< error
| 1064 42000 …near 'IS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ISOLATION)
< error
| 1064 42000 …near 'ISOLATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ISSUER)
< error
| 1064 42000 …near 'ISSUER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ITERATE)
< error
| 1064 42000 …near 'ITERATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JOIN)
< error
| 1064 42000 …near 'JOIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JSON)
< error
| 1064 42000 …near 'JSON)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JSON_ARRAYAGG)
< error
| 1064 42000 …near 'JSON_ARRAYAGG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JSON_DUALITY_OBJECT)
< error
| 1273 HY000 Unknown collation: 'JSON_DUALITY_OBJECT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JSON_OBJECTAGG)
< error
| 1064 42000 …near 'JSON_OBJECTAGG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JSON_TABLE)
< error
| 1064 42000 …near 'JSON_TABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE JSON_VALUE)
< error
| 1064 42000 …near 'JSON_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE KEY)
< error
| 1064 42000 …near 'KEY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE KEYRING)
< error
| 1273 HY000 Unknown collation: 'KEYRING'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE KEYS)
< error
| 1064 42000 …near 'KEYS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE KEY_BLOCK_SIZE)
< error
| 1064 42000 …near 'KEY_BLOCK_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE KILL)
< error
| 1064 42000 …near 'KILL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LAG)
< error
| 1064 42000 …near 'LAG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LANGUAGE)
< error
| 1064 42000 …near 'LANGUAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LAST)
< error
| 1064 42000 …near 'LAST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LAST_VALUE)
< error
| 1064 42000 …near 'LAST_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LATERAL)
< error
| 1064 42000 …near 'LATERAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEAD)
< error
| 1064 42000 …near 'LEAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEADING)
< error
| 1064 42000 …near 'LEADING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEAST)
< error
| 1273 HY000 Unknown collation: 'LEAST'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEAVE)
< error
| 1064 42000 …near 'LEAVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEAVES)
< error
| 1064 42000 …near 'LEAVES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEFT)
< error
| 1064 42000 …near 'LEFT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LESS)
< error
| 1064 42000 …near 'LESS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LEVEL)
< error
| 1064 42000 …near 'LEVEL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LIBRARY)
< error
| 1273 HY000 Unknown collation: 'LIBRARY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LIKE)
< error
| 1064 42000 …near 'LIKE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LIMIT)
< error
| 1064 42000 …near 'LIMIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LINEAR)
< error
| 1064 42000 …near 'LINEAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LINES)
< error
| 1064 42000 …near 'LINES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LINESTRING)
< error
| 1064 42000 …near 'LINESTRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LIST)
< error
| 1064 42000 …near 'LIST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOAD)
< error
| 1064 42000 …near 'LOAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOCAL)
< error
| 1064 42000 …near 'LOCAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOCALTIME)
< error
| 1064 42000 …near 'LOCALTIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOCALTIMESTAMP)
< error
| 1064 42000 …near 'LOCALTIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOCK)
< error
| 1064 42000 …near 'LOCK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOCKED)
< error
| 1064 42000 …near 'LOCKED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOCKS)
< error
| 1064 42000 …near 'LOCKS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOG)
< error
| 1064 42000 …near 'LOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOGFILE)
< error
| 1064 42000 …near 'LOGFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOGS)
< error
| 1064 42000 …near 'LOGS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LONG)
< error
| 1064 42000 …near 'LONG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LONGBLOB)
< error
| 1064 42000 …near 'LONGBLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LONGTEXT)
< error
| 1064 42000 …near 'LONGTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOOP)
< error
| 1064 42000 …near 'LOOP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE LOW_PRIORITY)
< error
| 1064 42000 …near 'LOW_PRIORITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MANIFEST)
< error
| 1273 HY000 Unknown collation: 'MANIFEST'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MANUAL)
< error
| 1064 42000 …near 'MANUAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MASTER)
< error
| 1064 42000 …near 'MASTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MATCH)
< error
| 1064 42000 …near 'MATCH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MATERIALIZED)
< error
| 1064 42000 …near 'MATERIALIZED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX)
< error
| 1064 42000 …near 'MAX)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAXVALUE)
< error
| 1064 42000 …near 'MAXVALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX_CONNECTIONS_PER_HOUR)
< error
| 1064 42000 …near 'MAX_CONNECTIONS_PER_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX_QUERIES_PER_HOUR)
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX_ROWS)
< error
| 1064 42000 …near 'MAX_ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX_SIZE)
< error
| 1064 42000 …near 'MAX_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX_UPDATES_PER_HOUR)
< error
| 1064 42000 …near 'MAX_UPDATES_PER_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MAX_USER_CONNECTIONS)
< error
| 1064 42000 …near 'MAX_USER_CONNECTIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MEDIUM)
< error
| 1064 42000 …near 'MEDIUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MEDIUMBLOB)
< error
| 1064 42000 …near 'MEDIUMBLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MEDIUMINT)
< error
| 1064 42000 …near 'MEDIUMINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MEDIUMTEXT)
< error
| 1064 42000 …near 'MEDIUMTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MEMBER)
< error
| 1064 42000 …near 'MEMBER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MEMORY)
< error
| 1064 42000 …near 'MEMORY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MERGE)
< error
| 1064 42000 …near 'MERGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MESSAGE_TEXT)
< error
| 1064 42000 …near 'MESSAGE_TEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MICROSECOND)
< error
| 1064 42000 …near 'MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MID)
< error
| 1273 HY000 Unknown collation: 'MID'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MIDDLEINT)
< error
| 1064 42000 …near 'MIDDLEINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MIGRATE)
< error
| 1064 42000 …near 'MIGRATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MIN)
< error
| 1064 42000 …near 'MIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MINUTE)
< error
| 1064 42000 …near 'MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MINUTE_MICROSECOND)
< error
| 1064 42000 …near 'MINUTE_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MINUTE_SECOND)
< error
| 1064 42000 …near 'MINUTE_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MIN_ROWS)
< error
| 1064 42000 …near 'MIN_ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MOD)
< error
| 1064 42000 …near 'MOD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MODE)
< error
| 1064 42000 …near 'MODE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MODIFIES)
< error
| 1064 42000 …near 'MODIFIES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MODIFY)
< error
| 1064 42000 …near 'MODIFY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MONTH)
< error
| 1064 42000 …near 'MONTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MULTILINESTRING)
< error
| 1064 42000 …near 'MULTILINESTRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MULTIPOINT)
< error
| 1064 42000 …near 'MULTIPOINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MULTIPOLYGON)
< error
| 1064 42000 …near 'MULTIPOLYGON)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MUTEX)
< error
| 1064 42000 …near 'MUTEX)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE MYSQL_ERRNO)
< error
| 1064 42000 …near 'MYSQL_ERRNO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NAME)
< error
| 1064 42000 …near 'NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NAMES)
< error
| 1064 42000 …near 'NAMES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NATIONAL)
< error
| 1064 42000 …near 'NATIONAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NATURAL)
< error
| 1064 42000 …near 'NATURAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NCHAR)
< error
| 1064 42000 …near 'NCHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NDB)
< error
| 1064 42000 …near 'NDB)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NDBCLUSTER)
< error
| 1064 42000 …near 'NDBCLUSTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NESTED)
< error
| 1064 42000 …near 'NESTED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NETWORK_NAMESPACE)
< error
| 1273 HY000 Unknown collation: 'NETWORK_NAMESPACE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NEVER)
< error
| 1064 42000 …near 'NEVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NEW)
< error
| 1064 42000 …near 'NEW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NEXT)
< error
| 1064 42000 …near 'NEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NO)
< error
| 1064 42000 …near 'NO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NODEGROUP)
< error
| 1064 42000 …near 'NODEGROUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NONE)
< error
| 1064 42000 …near 'NONE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NOT)
< error
| 1064 42000 …near 'NOT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NOW)
< error
| 1064 42000 …near 'NOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NOWAIT)
< error
| 1064 42000 …near 'NOWAIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NO_WAIT)
< error
| 1064 42000 …near 'NO_WAIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NO_WRITE_TO_BINLOG)
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NTH_VALUE)
< error
| 1064 42000 …near 'NTH_VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NTILE)
< error
| 1064 42000 …near 'NTILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NULLIF)
< error
| 1273 HY000 Unknown collation: 'NULLIF'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NULLS)
< error
| 1064 42000 …near 'NULLS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NUMBER)
< error
| 1064 42000 …near 'NUMBER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NUMERIC)
< error
| 1064 42000 …near 'NUMERIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE NVARCHAR)
< error
| 1064 42000 …near 'NVARCHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OF)
< error
| 1064 42000 …near 'OF)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OFF)
< error
| 1064 42000 …near 'OFF)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OFFSET)
< error
| 1064 42000 …near 'OFFSET)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OJ)
< error
| 1064 42000 …near 'OJ)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OLD)
< error
| 1064 42000 …near 'OLD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ON)
< error
| 1064 42000 …near 'ON)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ONE)
< error
| 1064 42000 …near 'ONE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ONLY)
< error
| 1064 42000 …near 'ONLY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPEN)
< error
| 1064 42000 …near 'OPEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPTIMIZE)
< error
| 1064 42000 …near 'OPTIMIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPTIMIZER_COSTS)
< error
| 1064 42000 …near 'OPTIMIZER_COSTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPTION)
< error
| 1064 42000 …near 'OPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPTIONAL)
< error
| 1273 HY000 Unknown collation: 'OPTIONAL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPTIONALLY)
< error
| 1064 42000 …near 'OPTIONALLY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OPTIONS)
< error
| 1064 42000 …near 'OPTIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OR)
< error
| 1064 42000 …near 'OR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ORDER)
< error
| 1064 42000 …near 'ORDER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ORDINALITY)
< error
| 1064 42000 …near 'ORDINALITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ORGANIZATION)
< error
| 1064 42000 …near 'ORGANIZATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OTHERS)
< error
| 1273 HY000 Unknown collation: 'OTHERS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OUT)
< error
| 1064 42000 …near 'OUT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OUTER)
< error
| 1064 42000 …near 'OUTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OUTFILE)
< error
| 1064 42000 …near 'OUTFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OVER)
< error
| 1064 42000 …near 'OVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OVERLAY)
< error
| 1273 HY000 Unknown collation: 'OVERLAY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OVERWRITE)
< error
| 1064 42000 …near 'OVERWRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE OWNER)
< error
| 1064 42000 …near 'OWNER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PACK_KEYS)
< error
| 1064 42000 …near 'PACK_KEYS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PAGE)
< error
| 1064 42000 …near 'PAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARALLEL)
< error
| 1064 42000 …near 'PARALLEL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARAMETERS)
< error
| 1064 42000 …near 'PARAMETERS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARSER)
< error
| 1064 42000 …near 'PARSER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARSE_TREE)
< error
| 1273 HY000 Unknown collation: 'PARSE_TREE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARTIAL)
< error
| 1064 42000 …near 'PARTIAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARTITION)
< error
| 1064 42000 …near 'PARTITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARTITIONING)
< error
| 1064 42000 …near 'PARTITIONING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PARTITIONS)
< error
| 1064 42000 …near 'PARTITIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PASSWORD)
< error
| 1064 42000 …near 'PASSWORD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PASSWORD_LOCK_TIME)
< error
| 1273 HY000 Unknown collation: 'PASSWORD_LOCK_TIME'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PATH)
< error
| 1064 42000 …near 'PATH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PERCENT_RANK)
< error
| 1064 42000 …near 'PERCENT_RANK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PERSIST)
< error
| 1273 HY000 Unknown collation: 'PERSIST'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PERSIST_ONLY)
< error
| 1273 HY000 Unknown collation: 'PERSIST_ONLY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PHASE)
< error
| 1064 42000 …near 'PHASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PLACING)
< error
| 1273 HY000 Unknown collation: 'PLACING'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PLUGIN)
< error
| 1064 42000 …near 'PLUGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PLUGINS)
< error
| 1064 42000 …near 'PLUGINS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PLUGIN_DIR)
< error
| 1064 42000 …near 'PLUGIN_DIR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE POINT)
< error
| 1064 42000 …near 'POINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE POLYGON)
< error
| 1064 42000 …near 'POLYGON)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PORT)
< error
| 1064 42000 …near 'PORT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE POSITION)
< error
| 1064 42000 …near 'POSITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRECEDES)
< error
| 1273 HY000 Unknown collation: 'PRECEDES'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRECEDING)
< error
| 1064 42000 …near 'PRECEDING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRECISION)
< error
| 1064 42000 …near 'PRECISION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PREFIX)
< error
| 1273 HY000 Unknown collation: 'PREFIX'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PREPARE)
< error
| 1064 42000 …near 'PREPARE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRESERVE)
< error
| 1064 42000 …near 'PRESERVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PREV)
< error
| 1064 42000 …near 'PREV)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRIMARY)
< error
| 1064 42000 …near 'PRIMARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRIVILEGES)
< error
| 1064 42000 …near 'PRIVILEGES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PRIVILEGE_CHECKS_USER)
< error
| 1273 HY000 Unknown collation: 'PRIVILEGE_CHECKS_USER'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PROCEDURE)
< error
| 1064 42000 …near 'PROCEDURE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PROCESS)
< error
| 1064 42000 …near 'PROCESS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PROCESSLIST)
< error
| 1064 42000 …near 'PROCESSLIST)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PROFILE)
< error
| 1064 42000 …near 'PROFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PROFILES)
< error
| 1064 42000 …near 'PROFILES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PROXY)
< error
| 1064 42000 …near 'PROXY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE PURGE)
< error
| 1064 42000 …near 'PURGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE QUALIFY)
< error
| 1273 HY000 Unknown collation: 'QUALIFY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE QUARTER)
< error
| 1064 42000 …near 'QUARTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE QUERY)
< error
| 1064 42000 …near 'QUERY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE QUICK)
< error
| 1064 42000 …near 'QUICK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RANDOM)
< error
| 1273 HY000 Unknown collation: 'RANDOM'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RANGE)
< error
| 1064 42000 …near 'RANGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RANK)
< error
| 1064 42000 …near 'RANK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE READ)
< error
| 1064 42000 …near 'READ)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE READS)
< error
| 1064 42000 …near 'READS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE READ_ONLY)
< error
| 1064 42000 …near 'READ_ONLY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE READ_WRITE)
< error
| 1064 42000 …near 'READ_WRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REAL)
< error
| 1064 42000 …near 'REAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REBUILD)
< error
| 1064 42000 …near 'REBUILD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RECOVER)
< error
| 1064 42000 …near 'RECOVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RECURSIVE)
< error
| 1064 42000 …near 'RECURSIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REDO_BUFFER_SIZE)
< error
| 1064 42000 …near 'REDO_BUFFER_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REDUNDANT)
< error
| 1064 42000 …near 'REDUNDANT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REFERENCE)
< error
| 1273 HY000 Unknown collation: 'REFERENCE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REFERENCES)
< error
| 1064 42000 …near 'REFERENCES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REGEXP)
< error
| 1064 42000 …near 'REGEXP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REGISTRATION)
< error
| 1273 HY000 Unknown collation: 'REGISTRATION'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELATIONAL)
< error
| 1273 HY000 Unknown collation: 'RELATIONAL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELAY)
< error
| 1064 42000 …near 'RELAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELAYLOG)
< error
| 1064 42000 …near 'RELAYLOG)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELAY_LOG_FILE)
< error
| 1064 42000 …near 'RELAY_LOG_FILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELAY_LOG_POS)
< error
| 1064 42000 …near 'RELAY_LOG_POS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELAY_THREAD)
< error
| 1064 42000 …near 'RELAY_THREAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELEASE)
< error
| 1064 42000 …near 'RELEASE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RELOAD)
< error
| 1064 42000 …near 'RELOAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REMOVE)
< error
| 1064 42000 …near 'REMOVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RENAME)
< error
| 1064 42000 …near 'RENAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REORGANIZE)
< error
| 1064 42000 …near 'REORGANIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPAIR)
< error
| 1064 42000 …near 'REPAIR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPEAT)
< error
| 1064 42000 …near 'REPEAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPEATABLE)
< error
| 1064 42000 …near 'REPEATABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLACE)
< error
| 1064 42000 …near 'REPLACE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICA)
< error
| 1064 42000 …near 'REPLICA)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICAS)
< error
| 1273 HY000 Unknown collation: 'REPLICAS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_DO_DB)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_DO_DB'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_DO_TABLE)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_DO_TABLE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_IGNORE_DB)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_IGNORE_DB'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_IGNORE_TABLE)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_IGNORE_TABLE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_REWRITE_DB)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_REWRITE_DB'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_WILD_DO_TABLE)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_WILD_DO_TABLE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATE_WILD_IGNORE_TABLE)
< error
| 1273 HY000 Unknown collation: 'REPLICATE_WILD_IGNORE_TABLE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REPLICATION)
< error
| 1064 42000 …near 'REPLICATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REQUIRE)
< error
| 1064 42000 …near 'REQUIRE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REQUIRE_ROW_FORMAT)
< error
| 1273 HY000 Unknown collation: 'REQUIRE_ROW_FORMAT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REQUIRE_TABLE_PRIMARY_KEY_CHECK)
< error
| 1273 HY000 Unknown collation: 'REQUIRE_TABLE_PRIMARY_KEY_CHECK'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESET)
< error
| 1064 42000 …near 'RESET)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESIGNAL)
< error
| 1064 42000 …near 'RESIGNAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESOURCE)
< error
| 1064 42000 …near 'RESOURCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESPECT)
< error
| 1064 42000 …near 'RESPECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESTART)
< error
| 1064 42000 …near 'RESTART)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESTORE)
< error
| 1064 42000 …near 'RESTORE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESTRICT)
< error
| 1064 42000 …near 'RESTRICT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RESUME)
< error
| 1064 42000 …near 'RESUME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RETAIN)
< error
| 1273 HY000 Unknown collation: 'RETAIN'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RETURN)
< error
| 1064 42000 …near 'RETURN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RETURNED_SQLSTATE)
< error
| 1064 42000 …near 'RETURNED_SQLSTATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RETURNING)
< error
| 1064 42000 …near 'RETURNING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RETURNS)
< error
| 1064 42000 …near 'RETURNS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REUSE)
< error
| 1273 HY000 Unknown collation: 'REUSE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REVERSE)
< error
| 1064 42000 …near 'REVERSE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE REVOKE)
< error
| 1064 42000 …near 'REVOKE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RIGHT)
< error
| 1064 42000 …near 'RIGHT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RLIKE)
< error
| 1064 42000 …near 'RLIKE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROLE)
< error
| 1064 42000 …near 'ROLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROLLBACK)
< error
| 1064 42000 …near 'ROLLBACK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROLLUP)
< error
| 1064 42000 …near 'ROLLUP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROTATE)
< error
| 1064 42000 …near 'ROTATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROUTINE)
< error
| 1064 42000 …near 'ROUTINE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROW)
< error
| 1064 42000 …near 'ROW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROWS)
< error
| 1064 42000 …near 'ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROW_COUNT)
< error
| 1064 42000 …near 'ROW_COUNT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROW_FORMAT)
< error
| 1064 42000 …near 'ROW_FORMAT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ROW_NUMBER)
< error
| 1064 42000 …near 'ROW_NUMBER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE RTREE)
< error
| 1064 42000 …near 'RTREE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE S3)
< error
| 1273 HY000 Unknown collation: 'S3'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SAVEPOINT)
< error
| 1064 42000 …near 'SAVEPOINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SCHEDULE)
< error
| 1064 42000 …near 'SCHEDULE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SCHEMA)
< error
| 1064 42000 …near 'SCHEMA)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SCHEMAS)
< error
| 1064 42000 …near 'SCHEMAS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SCHEMA_NAME)
< error
| 1064 42000 …near 'SCHEMA_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECOND)
< error
| 1064 42000 …near 'SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECONDARY)
< error
| 1273 HY000 Unknown collation: 'SECONDARY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECONDARY_ENGINE)
< error
| 1273 HY000 Unknown collation: 'SECONDARY_ENGINE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECONDARY_ENGINE_ATTRIBUTE)
< error
| 1273 HY000 Unknown collation: 'SECONDARY_ENGINE_ATTRIBUTE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECONDARY_LOAD)
< error
| 1273 HY000 Unknown collation: 'SECONDARY_LOAD'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECONDARY_UNLOAD)
< error
| 1273 HY000 Unknown collation: 'SECONDARY_UNLOAD'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECOND_MICROSECOND)
< error
| 1064 42000 …near 'SECOND_MICROSECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SECURITY)
< error
| 1064 42000 …near 'SECURITY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SELECT)
< error
| 1064 42000 …near 'SELECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SENSITIVE)
< error
| 1064 42000 …near 'SENSITIVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SEPARATOR)
< error
| 1064 42000 …near 'SEPARATOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SERIAL)
< error
| 1064 42000 …near 'SERIAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SERIALIZABLE)
< error
| 1064 42000 …near 'SERIALIZABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SERVER)
< error
| 1064 42000 …near 'SERVER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SESSION)
< error
| 1064 42000 …near 'SESSION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SESSION_USER)
< error
| 1064 42000 …near 'SESSION_USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SET)
< error
| 1064 42000 …near 'SET)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SETS)
< error
| 1064 42000 …near 'SETS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SHARE)
< error
| 1064 42000 …near 'SHARE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SHOW)
< error
| 1064 42000 …near 'SHOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SHUTDOWN)
< error
| 1064 42000 …near 'SHUTDOWN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SIGNAL)
< error
| 1064 42000 …near 'SIGNAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SIMPLE)
< error
| 1064 42000 …near 'SIMPLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SKIP)
< error
| 1064 42000 …near 'SKIP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SLAVE)
< error
| 1064 42000 …near 'SLAVE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SLOW)
< error
| 1064 42000 …near 'SLOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SMALLINT)
< error
| 1064 42000 …near 'SMALLINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SNAPSHOT)
< error
| 1064 42000 …near 'SNAPSHOT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOCKET)
< error
| 1064 42000 …near 'SOCKET)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOME)
< error
| 1064 42000 …near 'SOME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SONAME)
< error
| 1064 42000 …near 'SONAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOUNDS)
< error
| 1064 42000 …near 'SOUNDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE)
< error
| 1064 42000 …near 'SOURCE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_AUTO_POSITION)
< error
| 1273 HY000 Unknown collation: 'SOURCE_AUTO_POSITION'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_BIND)
< error
| 1273 HY000 Unknown collation: 'SOURCE_BIND'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_COMPRESSION_ALGORITHMS)
< error
| 1273 HY000 Unknown collation: 'SOURCE_COMPRESSION_ALGORITHMS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_CONNECTION_AUTO_FAILOVER)
< error
| 1273 HY000 Unknown collation: 'SOURCE_CONNECTION_AUTO_FAILOVER'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_CONNECT_RETRY)
< error
| 1273 HY000 Unknown collation: 'SOURCE_CONNECT_RETRY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_DELAY)
< error
| 1273 HY000 Unknown collation: 'SOURCE_DELAY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_HEARTBEAT_PERIOD)
< error
| 1273 HY000 Unknown collation: 'SOURCE_HEARTBEAT_PERIOD'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_HOST)
< error
| 1273 HY000 Unknown collation: 'SOURCE_HOST'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_LOG_FILE)
< error
| 1273 HY000 Unknown collation: 'SOURCE_LOG_FILE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_LOG_POS)
< error
| 1273 HY000 Unknown collation: 'SOURCE_LOG_POS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_PASSWORD)
< error
| 1273 HY000 Unknown collation: 'SOURCE_PASSWORD'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_PORT)
< error
| 1273 HY000 Unknown collation: 'SOURCE_PORT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_PUBLIC_KEY_PATH)
< error
| 1273 HY000 Unknown collation: 'SOURCE_PUBLIC_KEY_PATH'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_RETRY_COUNT)
< error
| 1273 HY000 Unknown collation: 'SOURCE_RETRY_COUNT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_CA)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_CA'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_CAPATH)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_CAPATH'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_CERT)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_CERT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_CIPHER)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_CIPHER'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_CRL)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_CRL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_CRLPATH)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_CRLPATH'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_KEY)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_KEY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_SSL_VERIFY_SERVER_CERT)
< error
| 1273 HY000 Unknown collation: 'SOURCE_SSL_VERIFY_SERVER_CERT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_TLS_CIPHERSUITES)
< error
| 1273 HY000 Unknown collation: 'SOURCE_TLS_CIPHERSUITES'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_TLS_VERSION)
< error
| 1273 HY000 Unknown collation: 'SOURCE_TLS_VERSION'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_USER)
< error
| 1273 HY000 Unknown collation: 'SOURCE_USER'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SOURCE_ZSTD_COMPRESSION_LEVEL)
< error
| 1273 HY000 Unknown collation: 'SOURCE_ZSTD_COMPRESSION_LEVEL'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SPATIAL)
< error
| 1064 42000 …near 'SPATIAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SPECIFIC)
< error
| 1064 42000 …near 'SPECIFIC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL)
< error
| 1064 42000 …near 'SQL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQLEXCEPTION)
< error
| 1064 42000 …near 'SQLEXCEPTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQLSTATE)
< error
| 1064 42000 …near 'SQLSTATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQLWARNING)
< error
| 1064 42000 …near 'SQLWARNING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_AFTER_GTIDS)
< error
| 1064 42000 …near 'SQL_AFTER_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_AFTER_MTS_GAPS)
< error
| 1064 42000 …near 'SQL_AFTER_MTS_GAPS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_BEFORE_GTIDS)
< error
| 1064 42000 …near 'SQL_BEFORE_GTIDS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_BIG_RESULT)
< error
| 1064 42000 …near 'SQL_BIG_RESULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_BUFFER_RESULT)
< error
| 1064 42000 …near 'SQL_BUFFER_RESULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_CALC_FOUND_ROWS)
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_NO_CACHE)
< error
| 1064 42000 …near 'SQL_NO_CACHE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_SMALL_RESULT)
< error
| 1064 42000 …near 'SQL_SMALL_RESULT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_THREAD)
< error
| 1064 42000 …near 'SQL_THREAD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_DAY)
< error
| 1064 42000 …near 'SQL_TSI_DAY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_HOUR)
< error
| 1064 42000 …near 'SQL_TSI_HOUR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_MINUTE)
< error
| 1064 42000 …near 'SQL_TSI_MINUTE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_MONTH)
< error
| 1064 42000 …near 'SQL_TSI_MONTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_QUARTER)
< error
| 1064 42000 …near 'SQL_TSI_QUARTER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_SECOND)
< error
| 1064 42000 …near 'SQL_TSI_SECOND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_WEEK)
< error
| 1064 42000 …near 'SQL_TSI_WEEK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SQL_TSI_YEAR)
< error
| 1064 42000 …near 'SQL_TSI_YEAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SRID)
< error
| 1064 42000 …near 'SRID)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SSE_KMS)
< error
| 1273 HY000 Unknown collation: 'SSE_KMS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SSE_S3)
< error
| 1273 HY000 Unknown collation: 'SSE_S3'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SSL)
< error
| 1064 42000 …near 'SSL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STACKED)
< error
| 1064 42000 …near 'STACKED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE START)
< error
| 1064 42000 …near 'START)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STARTING)
< error
| 1064 42000 …near 'STARTING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STARTS)
< error
| 1064 42000 …near 'STARTS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STATS_AUTO_RECALC)
< error
| 1064 42000 …near 'STATS_AUTO_RECALC)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STATS_PERSISTENT)
< error
| 1064 42000 …near 'STATS_PERSISTENT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STATS_SAMPLE_PAGES)
< error
| 1064 42000 …near 'STATS_SAMPLE_PAGES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STATUS)
< error
| 1064 42000 …near 'STATUS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STD)
< error
| 1064 42000 …near 'STD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STDDEV)
< error
| 1064 42000 …near 'STDDEV)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STDDEV_POP)
< error
| 1064 42000 …near 'STDDEV_POP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STDDEV_SAMP)
< error
| 1064 42000 …near 'STDDEV_SAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STOP)
< error
| 1064 42000 …near 'STOP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STORAGE)
< error
| 1064 42000 …near 'STORAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STORED)
< error
| 1064 42000 …near 'STORED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STRAIGHT_JOIN)
< error
| 1064 42000 …near 'STRAIGHT_JOIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STREAM)
< error
| 1273 HY000 Unknown collation: 'STREAM'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STRICT_LOAD)
< error
| 1273 HY000 Unknown collation: 'STRICT_LOAD'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE STRING)
< error
| 1064 42000 …near 'STRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ST_COLLECT)
< error
| 1273 HY000 Unknown collation: 'ST_COLLECT'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBCLASS_ORIGIN)
< error
| 1064 42000 …near 'SUBCLASS_ORIGIN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBDATE)
< error
| 1064 42000 …near 'SUBDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBJECT)
< error
| 1064 42000 …near 'SUBJECT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBPARTITION)
< error
| 1064 42000 …near 'SUBPARTITION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBPARTITIONS)
< error
| 1064 42000 …near 'SUBPARTITIONS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBSTR)
< error
| 1064 42000 …near 'SUBSTR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUBSTRING)
< error
| 1064 42000 …near 'SUBSTRING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUM)
< error
| 1064 42000 …near 'SUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUPER)
< error
| 1064 42000 …near 'SUPER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SUSPEND)
< error
| 1064 42000 …near 'SUSPEND)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SWAPS)
< error
| 1064 42000 …near 'SWAPS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SWITCHES)
< error
| 1064 42000 …near 'SWITCHES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SYSDATE)
< error
| 1064 42000 …near 'SYSDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SYSTEM)
< error
| 1064 42000 …near 'SYSTEM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE SYSTEM_USER)
< error
| 1064 42000 …near 'SYSTEM_USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TABLE)
< error
| 1064 42000 …near 'TABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TABLES)
< error
| 1064 42000 …near 'TABLES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TABLESAMPLE)
< error
| 1273 HY000 Unknown collation: 'TABLESAMPLE'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TABLESPACE)
< error
| 1064 42000 …near 'TABLESPACE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TABLE_CHECKSUM)
< error
| 1064 42000 …near 'TABLE_CHECKSUM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TABLE_NAME)
< error
| 1064 42000 …near 'TABLE_NAME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TEMPORARY)
< error
| 1064 42000 …near 'TEMPORARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TEMPTABLE)
< error
| 1064 42000 …near 'TEMPTABLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TERMINATED)
< error
| 1064 42000 …near 'TERMINATED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE THAN)
< error
| 1064 42000 …near 'THAN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE THEN)
< error
| 1064 42000 …near 'THEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE THREAD_PRIORITY)
< error
| 1273 HY000 Unknown collation: 'THREAD_PRIORITY'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TIES)
< error
| 1273 HY000 Unknown collation: 'TIES'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TIME)
< error
| 1064 42000 …near 'TIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TIMESTAMP)
< error
| 1064 42000 …near 'TIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TIMESTAMPADD)
< error
| 1064 42000 …near 'TIMESTAMPADD)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TIMESTAMPDIFF)
< error
| 1064 42000 …near 'TIMESTAMPDIFF)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TINYBLOB)
< error
| 1064 42000 …near 'TINYBLOB)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TINYINT)
< error
| 1064 42000 …near 'TINYINT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TINYTEXT)
< error
| 1064 42000 …near 'TINYTEXT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TLS)
< error
| 1273 HY000 Unknown collation: 'TLS'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TO)
< error
| 1064 42000 …near 'TO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRAILING)
< error
| 1064 42000 …near 'TRAILING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRANSACTION)
< error
| 1064 42000 …near 'TRANSACTION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRIGGER)
< error
| 1064 42000 …near 'TRIGGER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRIGGERS)
< error
| 1064 42000 …near 'TRIGGERS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRIM)
< error
| 1064 42000 …near 'TRIM)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRUE)
< error
| 1064 42000 …near 'TRUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TRUNCATE)
< error
| 1064 42000 …near 'TRUNCATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TYPE)
< error
| 1064 42000 …near 'TYPE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE TYPES)
< error
| 1064 42000 …near 'TYPES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNBOUNDED)
< error
| 1064 42000 …near 'UNBOUNDED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNCOMMITTED)
< error
| 1064 42000 …near 'UNCOMMITTED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNDEFINED)
< error
| 1064 42000 …near 'UNDEFINED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNDO)
< error
| 1064 42000 …near 'UNDO)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNDOFILE)
< error
| 1064 42000 …near 'UNDOFILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNDO_BUFFER_SIZE)
< error
| 1064 42000 …near 'UNDO_BUFFER_SIZE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNINSTALL)
< error
| 1064 42000 …near 'UNINSTALL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNION)
< error
| 1064 42000 …near 'UNION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNIQUE)
< error
| 1064 42000 …near 'UNIQUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNKNOWN)
< error
| 1064 42000 …near 'UNKNOWN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNLOCK)
< error
| 1064 42000 …near 'UNLOCK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNREGISTER)
< error
| 1273 HY000 Unknown collation: 'UNREGISTER'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UNTIL)
< error
| 1064 42000 …near 'UNTIL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UPDATE)
< error
| 1064 42000 …near 'UPDATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UPGRADE)
< error
| 1064 42000 …near 'UPGRADE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE URI)
< error
| 1273 HY000 Unknown collation: 'URI'
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE USAGE)
< error
| 1064 42000 …near 'USAGE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE USE)
< error
| 1064 42000 …near 'USE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE USER)
< error
| 1064 42000 …near 'USER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE USING)
< error
| 1064 42000 …near 'USING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UTC_DATE)
< error
| 1064 42000 …near 'UTC_DATE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UTC_TIME)
< error
| 1064 42000 …near 'UTC_TIME)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE UTC_TIMESTAMP)
< error
| 1064 42000 …near 'UTC_TIMESTAMP)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VALIDATION)
< error
| 1064 42000 …near 'VALIDATION)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VALUE)
< error
| 1064 42000 …near 'VALUE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VALUES)
< error
| 1064 42000 …near 'VALUES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VARBINARY)
< error
| 1064 42000 …near 'VARBINARY)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VARCHAR)
< error
| 1064 42000 …near 'VARCHAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VARIABLES)
< error
| 1064 42000 …near 'VARIABLES)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VARYING)
< error
| 1064 42000 …near 'VARYING)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VIEW)
< error
| 1064 42000 …near 'VIEW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VIRTUAL)
< error
| 1064 42000 …near 'VIRTUAL)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE VISIBLE)
< error
| 1064 42000 …near 'VISIBLE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WAIT)
< error
| 1064 42000 …near 'WAIT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WARNINGS)
< error
| 1064 42000 …near 'WARNINGS)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WHEN)
< error
| 1064 42000 …near 'WHEN)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WHERE)
< error
| 1064 42000 …near 'WHERE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WHILE)
< error
| 1064 42000 …near 'WHILE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WINDOW)
< error
| 1064 42000 …near 'WINDOW)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WITH)
< error
| 1064 42000 …near 'WITH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WITHOUT)
< error
| 1064 42000 …near 'WITHOUT)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WORK)
< error
| 1064 42000 …near 'WORK)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WRAPPER)
< error
| 1064 42000 …near 'WRAPPER)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE WRITE)
< error
| 1064 42000 …near 'WRITE)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE X509)
< error
| 1064 42000 …near 'X509)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE XA)
< error
| 1064 42000 …near 'XA)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE XID)
< error
| 1064 42000 …near 'XID)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE XML)
< error
| 1064 42000 …near 'XML)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE XOR)
< error
| 1064 42000 …near 'XOR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE YEAR)
< error
| 1064 42000 …near 'YEAR)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE YEAR_MONTH)
< error
| 1064 42000 …near 'YEAR_MONTH)' at line 1
> exec ^
: string
| CREATE TABLE kw_co_test (a VARCHAR(50) COLLATE ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ABSENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ACCESSIBLE
< error
| 1064 42000 …near 'ACCESSIBLE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ACCOUNT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ACTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ACTIVE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ADD
< error
| 1064 42000 …near 'ADD' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ADDDATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ADMIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AFTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AGAINST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AGGREGATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ALGORITHM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ALL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ALLOW_MISSING_FILES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ALTER
< error
| 1064 42000 …near 'ALTER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ALWAYS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ANALYZE
< error
| 1064 42000 …near 'ANALYZE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AND
< error
| 1064 42000 …near 'AND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ANY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ARRAY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AS
< error
| 1064 42000 …near 'AS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ASC
< error
| 1064 42000 …near 'ASC' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ASCII
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ASENSITIVE
< error
| 1064 42000 …near 'ASENSITIVE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ATTRIBUTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AUTHENTICATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AUTO
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AUTOEXTEND_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AUTO_INCREMENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AUTO_REFRESH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AUTO_REFRESH_SOURCE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AVG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=AVG_ROW_LENGTH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BACKUP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BEFORE
< error
| 1064 42000 …near 'BEFORE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BEGIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BERNOULLI
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BETWEEN
< error
| 1064 42000 …near 'BETWEEN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BIGINT
< error
| 1064 42000 …near 'BIGINT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BINLOG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BIT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BIT_AND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BIT_OR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BIT_XOR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BLOB
< error
| 1064 42000 …near 'BLOB' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BLOCK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BOOL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BOOLEAN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BOTH
< error
| 1064 42000 …near 'BOTH' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BTREE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BUCKETS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BULK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BY
< error
| 1064 42000 …near 'BY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=BYTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CACHE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CALL
< error
| 1064 42000 …near 'CALL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CASCADE
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CASCADED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CASE
< error
| 1064 42000 …near 'CASE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CAST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CATALOG_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHAIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHALLENGE_RESPONSE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHANGE
< error
| 1064 42000 …near 'CHANGE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHANGED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHANNEL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHAR
< error
| 1064 42000 …near 'CHAR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHARACTER
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHARSET
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHECK
< error
| 1064 42000 …near 'CHECK' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CHECKSUM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CIPHER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CLASS_ORIGIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CLIENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CLONE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CLOSE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COALESCE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CODE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COLLATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COLUMN
< error
| 1064 42000 …near 'COLUMN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COLUMNS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COLUMN_FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COLUMN_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMMENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMMIT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMMITTED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMPACT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMPLETION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMPONENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMPRESSED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COMPRESSION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONCAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONCURRENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONDITION
< error
| 1064 42000 …near 'CONDITION' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONNECTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONSISTENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONSTRAINT
< error
| 1064 42000 …near 'CONSTRAINT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONSTRAINT_CATALOG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONSTRAINT_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONSTRAINT_SCHEMA
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONTAINS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONTEXT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONTINUE
< error
| 1064 42000 …near 'CONTINUE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CONVERT
< error
| 1064 42000 …near 'CONVERT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=COUNT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CPU
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CREATE
< error
| 1064 42000 …near 'CREATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CROSS
< error
| 1064 42000 …near 'CROSS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CSV
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CUBE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CUME_DIST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURDATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURRENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURRENT_DATE
< error
| 1064 42000 …near 'CURRENT_DATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURRENT_TIME
< error
| 1064 42000 …near 'CURRENT_TIME' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURRENT_TIMESTAMP
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURRENT_USER
< error
| 1064 42000 …near 'CURRENT_USER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURSOR
< error
| 1064 42000 …near 'CURSOR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURSOR_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=CURTIME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATA
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATABASE
< error
| 1064 42000 …near 'DATABASE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATABASES
< error
| 1064 42000 …near 'DATABASES' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATAFILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATETIME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATE_ADD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DATE_SUB
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DAY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DAY_HOUR
< error
| 1064 42000 …near 'DAY_HOUR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DAY_MICROSECOND
< error
| 1064 42000 …near 'DAY_MICROSECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DAY_MINUTE
< error
| 1064 42000 …near 'DAY_MINUTE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DAY_SECOND
< error
| 1064 42000 …near 'DAY_SECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DEALLOCATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DEC
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DECIMAL
< error
| 1064 42000 …near 'DECIMAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DECLARE
< error
| 1064 42000 …near 'DECLARE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DEFAULT_AUTH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DEFINER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DEFINITION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DELAYED
< error
| 1064 42000 …near 'DELAYED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DELAY_KEY_WRITE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DELETE
< error
| 1064 42000 …near 'DELETE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DENSE_RANK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DESC
< error
| 1064 42000 …near 'DESC' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DESCRIBE
< error
| 1064 42000 …near 'DESCRIBE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DESCRIPTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DETERMINISTIC
< error
| 1064 42000 …near 'DETERMINISTIC' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DIAGNOSTICS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DIRECTORY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DISABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DISCARD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DISK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DISTINCT
< error
| 1064 42000 …near 'DISTINCT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DISTINCTROW
< error
| 1064 42000 …near 'DISTINCTROW' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DIV
< error
| 1064 42000 …near 'DIV' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DO
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DOUBLE
< error
| 1064 42000 …near 'DOUBLE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DROP
< error
| 1064 42000 …near 'DROP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DUAL
< error
| 1064 42000 …near 'DUAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DUALITY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DUMPFILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DUPLICATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=DYNAMIC
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EACH
< error
| 1064 42000 …near 'EACH' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ELSE
< error
| 1064 42000 …near 'ELSE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ELSEIF
< error
| 1064 42000 …near 'ELSEIF' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EMPTY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENCLOSED
< error
| 1064 42000 …near 'ENCLOSED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENCRYPTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=END
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENFORCED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENGINE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENGINES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENGINE_ATTRIBUTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ENUM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ERROR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ERRORS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ESCAPE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ESCAPED
< error
| 1064 42000 …near 'ESCAPED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EVENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EVENTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EVERY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXCEPT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXCHANGE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXCLUDE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXECUTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXISTS
< error
| 1064 42000 …near 'EXISTS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXIT
< error
| 1064 42000 …near 'EXIT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXPANSION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXPIRE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXPLAIN
< error
| 1064 42000 …near 'EXPLAIN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXPORT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXTENDED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXTENT_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXTERNAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXTERNAL_FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=EXTRACT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FACTOR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FAILED_LOGIN_ATTEMPTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FALSE
< error
| 1064 42000 …near 'FALSE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FAST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FAULTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FETCH
< error
| 1064 42000 …near 'FETCH' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FIELDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILE_BLOCK_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILE_FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILE_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILE_PATTERN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILE_PREFIX
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FILTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FINISH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FIRST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FIRST_VALUE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FIXED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FLOAT
< error
| 1064 42000 …near 'FLOAT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FLOAT4
< error
| 1064 42000 …near 'FLOAT4' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FLOAT8
< error
| 1064 42000 …near 'FLOAT8' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FLUSH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FOLLOWING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FOLLOWS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FOR
< error
| 1064 42000 …near 'FOR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FORCE
< error
| 1064 42000 …near 'FORCE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FOREIGN
< error
| 1064 42000 …near 'FOREIGN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FOUND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FROM
< error
| 1064 42000 …near 'FROM' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FULL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FULLTEXT
< error
| 1064 42000 …near 'FULLTEXT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=FUNCTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GENERAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GENERATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GENERATED
< error
| 1064 42000 …near 'GENERATED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GEOMCOLLECTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GEOMETRY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GEOMETRYCOLLECTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GET
< error
| 1064 42000 …near 'GET' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GET_FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GET_SOURCE_PUBLIC_KEY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GLOBAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GRANT
< error
| 1064 42000 …near 'GRANT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GRANTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GREATEST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GROUP
< error
| 1064 42000 …near 'GROUP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GROUPING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GROUPS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GROUP_CONCAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GROUP_REPLICATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GTIDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GTID_ONLY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=GUIDED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HANDLER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HASH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HAVING
< error
| 1064 42000 …near 'HAVING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HEADER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HELP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HIGH_PRIORITY
< error
| 1064 42000 …near 'HIGH_PRIORITY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HISTOGRAM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HISTORY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HOST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HOSTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HOUR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HOUR_MICROSECOND
< error
| 1064 42000 …near 'HOUR_MICROSECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HOUR_MINUTE
< error
| 1064 42000 …near 'HOUR_MINUTE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=HOUR_SECOND
< error
| 1064 42000 …near 'HOUR_SECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IDENTIFIED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IF
< error
| 1064 42000 …near 'IF' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IGNORE
< error
| 1064 42000 …near 'IGNORE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IGNORE_SERVER_IDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IMPORT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IN
< error
| 1064 42000 …near 'IN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INACTIVE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INDEXES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INFILE
< error
| 1064 42000 …near 'INFILE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INITIAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INITIAL_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INITIATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INNER
< error
| 1064 42000 …near 'INNER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INOUT
< error
| 1064 42000 …near 'INOUT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INSENSITIVE
< error
| 1064 42000 …near 'INSENSITIVE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INSERT
< error
| 1064 42000 …near 'INSERT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INSERT_METHOD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INSTALL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INSTANCE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INT
< error
| 1064 42000 …near 'INT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INT1
< error
| 1064 42000 …near 'INT1' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INT2
< error
| 1064 42000 …near 'INT2' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INT3
< error
| 1064 42000 …near 'INT3' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INT4
< error
| 1064 42000 …near 'INT4' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INT8
< error
| 1064 42000 …near 'INT8' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INTEGER
< error
| 1064 42000 …near 'INTEGER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INTERSECT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INTERVAL
< error
| 1064 42000 …near 'INTERVAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INTO
< error
| 1064 42000 …near 'INTO' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INVISIBLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=INVOKER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IO
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IO_AFTER_GTIDS
< error
| 1064 42000 …near 'IO_AFTER_GTIDS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IO_BEFORE_GTIDS
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IO_THREAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IPC
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=IS
< error
| 1064 42000 …near 'IS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ISOLATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ISSUER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ITERATE
< error
| 1064 42000 …near 'ITERATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JOIN
< error
| 1064 42000 …near 'JOIN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JSON
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JSON_ARRAYAGG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JSON_DUALITY_OBJECT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JSON_OBJECTAGG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JSON_TABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=JSON_VALUE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=KEY
< error
| 1064 42000 …near 'KEY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=KEYRING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=KEYS
< error
| 1064 42000 …near 'KEYS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=KEY_BLOCK_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=KILL
< error
| 1064 42000 …near 'KILL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LAG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LANGUAGE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LAST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LAST_VALUE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LATERAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEADING
< error
| 1064 42000 …near 'LEADING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEAST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEAVE
< error
| 1064 42000 …near 'LEAVE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEAVES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEFT
< error
| 1064 42000 …near 'LEFT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LESS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LEVEL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LIBRARY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LIKE
< error
| 1064 42000 …near 'LIKE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LIMIT
< error
| 1064 42000 …near 'LIMIT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LINEAR
< error
| 1064 42000 …near 'LINEAR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LINES
< error
| 1064 42000 …near 'LINES' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LINESTRING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LIST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOAD
< error
| 1064 42000 …near 'LOAD' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOCAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOCALTIME
< error
| 1064 42000 …near 'LOCALTIME' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOCALTIMESTAMP
< error
| 1064 42000 …near 'LOCALTIMESTAMP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOCK
< error
| 1064 42000 …near 'LOCK' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOCKED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOCKS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOGFILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOGS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LONG
< error
| 1064 42000 …near 'LONG' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LONGBLOB
< error
| 1064 42000 …near 'LONGBLOB' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LONGTEXT
< error
| 1064 42000 …near 'LONGTEXT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOOP
< error
| 1064 42000 …near 'LOOP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=LOW_PRIORITY
< error
| 1064 42000 …near 'LOW_PRIORITY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MANIFEST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MANUAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MASTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MATCH
< error
| 1064 42000 …near 'MATCH' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MATERIALIZED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAXVALUE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX_CONNECTIONS_PER_HOUR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX_QUERIES_PER_HOUR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX_ROWS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX_UPDATES_PER_HOUR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MAX_USER_CONNECTIONS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MEDIUM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MEDIUMBLOB
< error
| 1064 42000 …near 'MEDIUMBLOB' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MEDIUMINT
< error
| 1064 42000 …near 'MEDIUMINT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MEDIUMTEXT
< error
| 1064 42000 …near 'MEDIUMTEXT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MEMBER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MEMORY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MERGE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MESSAGE_TEXT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MICROSECOND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MID
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MIDDLEINT
< error
| 1064 42000 …near 'MIDDLEINT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MIGRATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MINUTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MINUTE_MICROSECOND
< error
| 1064 42000 …near 'MINUTE_MICROSECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MINUTE_SECOND
< error
| 1064 42000 …near 'MINUTE_SECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MIN_ROWS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MOD
< error
| 1064 42000 …near 'MOD' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MODE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MODIFIES
< error
| 1064 42000 …near 'MODIFIES' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MODIFY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MONTH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MULTILINESTRING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MULTIPOINT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MULTIPOLYGON
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MUTEX
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=MYSQL_ERRNO
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NAMES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NATIONAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NATURAL
< error
| 1064 42000 …near 'NATURAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NCHAR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NDB
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NDBCLUSTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NESTED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NETWORK_NAMESPACE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NEVER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NEW
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NEXT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NO
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NODEGROUP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NONE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NOW
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NOWAIT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NO_WAIT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NO_WRITE_TO_BINLOG
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NTH_VALUE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NTILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NULLIF
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NULLS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NUMBER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NUMERIC
< error
| 1064 42000 …near 'NUMERIC' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=NVARCHAR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OF
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OFF
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OFFSET
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OJ
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OLD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ON
< error
| 1064 42000 …near 'ON' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ONE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ONLY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPEN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPTIMIZE
< error
| 1064 42000 …near 'OPTIMIZE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPTIMIZER_COSTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPTION
< error
| 1064 42000 …near 'OPTION' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPTIONAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPTIONALLY
< error
| 1064 42000 …near 'OPTIONALLY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OPTIONS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OR
< error
| 1064 42000 …near 'OR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ORDER
< error
| 1064 42000 …near 'ORDER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ORDINALITY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ORGANIZATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OTHERS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OUT
< error
| 1064 42000 …near 'OUT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OUTER
< error
| 1064 42000 …near 'OUTER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OUTFILE
< error
| 1064 42000 …near 'OUTFILE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OVER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OVERLAY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OVERWRITE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=OWNER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PACK_KEYS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PAGE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARALLEL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARAMETERS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARSER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARSE_TREE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARTIAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARTITION
< error
| 1064 42000 …near 'PARTITION' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARTITIONING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PARTITIONS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PASSWORD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PASSWORD_LOCK_TIME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PATH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PERCENT_RANK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PERSIST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PERSIST_ONLY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PHASE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PLACING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PLUGIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PLUGINS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PLUGIN_DIR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=POINT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=POLYGON
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PORT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=POSITION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRECEDES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRECEDING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRECISION
< error
| 1064 42000 …near 'PRECISION' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PREFIX
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PREPARE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRESERVE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PREV
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRIMARY
< error
| 1064 42000 …near 'PRIMARY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRIVILEGES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PRIVILEGE_CHECKS_USER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PROCEDURE
< error
| 1064 42000 …near 'PROCEDURE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PROCESS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PROCESSLIST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PROFILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PROFILES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PROXY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=PURGE
< error
| 1064 42000 …near 'PURGE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=QUALIFY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=QUARTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=QUERY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=QUICK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RANDOM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RANGE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RANK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=READ
< error
| 1064 42000 …near 'READ' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=READS
< error
| 1064 42000 …near 'READS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=READ_ONLY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=READ_WRITE
< error
| 1064 42000 …near 'READ_WRITE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REAL
< error
| 1064 42000 …near 'REAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REBUILD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RECOVER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RECURSIVE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REDO_BUFFER_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REDUNDANT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REFERENCE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REFERENCES
< error
| 1064 42000 …near 'REFERENCES' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REGEXP
< error
| 1064 42000 …near 'REGEXP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REGISTRATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELATIONAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELAY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELAYLOG
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELAY_LOG_FILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELAY_LOG_POS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELAY_THREAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELEASE
< error
| 1064 42000 …near 'RELEASE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RELOAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REMOVE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RENAME
< error
| 1064 42000 …near 'RENAME' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REORGANIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPAIR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPEAT
< error
| 1064 42000 …near 'REPEAT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPEATABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLACE
< error
| 1064 42000 …near 'REPLACE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICA
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICAS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_DO_DB
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_DO_TABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_IGNORE_DB
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_IGNORE_TABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_REWRITE_DB
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_WILD_DO_TABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATE_WILD_IGNORE_TABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REPLICATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REQUIRE
< error
| 1064 42000 …near 'REQUIRE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REQUIRE_ROW_FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REQUIRE_TABLE_PRIMARY_KEY_CHECK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESET
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESIGNAL
< error
| 1064 42000 …near 'RESIGNAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESOURCE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESPECT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESTART
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESTORE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESTRICT
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RESUME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RETAIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RETURN
< error
| 1064 42000 …near 'RETURN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RETURNED_SQLSTATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RETURNING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RETURNS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REUSE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REVERSE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=REVOKE
< error
| 1064 42000 …near 'REVOKE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RIGHT
< error
| 1064 42000 …near 'RIGHT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RLIKE
< error
| 1064 42000 …near 'RLIKE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROLLBACK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROLLUP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROTATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROUTINE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROW
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROWS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROW_COUNT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROW_FORMAT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ROW_NUMBER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=RTREE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=S3
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SAVEPOINT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SCHEDULE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SCHEMA
< error
| 1064 42000 …near 'SCHEMA' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SCHEMAS
< error
| 1064 42000 …near 'SCHEMAS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SCHEMA_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECOND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECONDARY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECONDARY_ENGINE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECONDARY_ENGINE_ATTRIBUTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECONDARY_LOAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECONDARY_UNLOAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECOND_MICROSECOND
< error
| 1064 42000 …near 'SECOND_MICROSECOND' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SECURITY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SELECT
< error
| 1064 42000 …near 'SELECT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SENSITIVE
< error
| 1064 42000 …near 'SENSITIVE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SEPARATOR
< error
| 1064 42000 …near 'SEPARATOR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SERIAL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SERIALIZABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SERVER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SESSION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SESSION_USER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SET
< error
| 1064 42000 …near 'SET' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SETS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SHARE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SHOW
< error
| 1064 42000 …near 'SHOW' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SHUTDOWN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SIGNAL
< error
| 1064 42000 …near 'SIGNAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SIGNED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SIMPLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SKIP
< error
| 1064 42000 …near 'SKIP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SLAVE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SLOW
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SMALLINT
< error
| 1064 42000 …near 'SMALLINT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SNAPSHOT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOCKET
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SONAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOUNDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_AUTO_POSITION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_BIND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_COMPRESSION_ALGORITHMS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_CONNECTION_AUTO_FAILOVER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_CONNECT_RETRY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_DELAY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_HEARTBEAT_PERIOD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_HOST
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_LOG_FILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_LOG_POS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_PASSWORD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_PORT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_PUBLIC_KEY_PATH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_RETRY_COUNT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_CA
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_CAPATH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_CERT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_CIPHER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_CRL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_CRLPATH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_KEY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_SSL_VERIFY_SERVER_CERT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_TLS_CIPHERSUITES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_TLS_VERSION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_USER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SOURCE_ZSTD_COMPRESSION_LEVEL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SPATIAL
< error
| 1064 42000 …near 'SPATIAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SPECIFIC
< error
| 1064 42000 …near 'SPECIFIC' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL
< error
| 1064 42000 …near 'SQL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQLEXCEPTION
< error
| 1064 42000 …near 'SQLEXCEPTION' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQLSTATE
< error
| 1064 42000 …near 'SQLSTATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQLWARNING
< error
| 1064 42000 …near 'SQLWARNING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_AFTER_GTIDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_AFTER_MTS_GAPS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_BEFORE_GTIDS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_BIG_RESULT
< error
| 1064 42000 …near 'SQL_BIG_RESULT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_BUFFER_RESULT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_CALC_FOUND_ROWS
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_NO_CACHE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_SMALL_RESULT
< error
| 1064 42000 …near 'SQL_SMALL_RESULT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_THREAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_DAY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_HOUR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_MINUTE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_MONTH
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_QUARTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_SECOND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_WEEK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SQL_TSI_YEAR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SRID
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SSE_KMS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SSE_S3
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SSL
< error
| 1064 42000 …near 'SSL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STACKED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=START
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STARTING
< error
| 1064 42000 …near 'STARTING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STARTS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STATS_AUTO_RECALC
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STATS_PERSISTENT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STATS_SAMPLE_PAGES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STATUS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STDDEV
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STDDEV_POP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STDDEV_SAMP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STOP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STORAGE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STORED
< error
| 1064 42000 …near 'STORED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STRAIGHT_JOIN
< error
| 1064 42000 …near 'STRAIGHT_JOIN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STREAM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STRICT_LOAD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=STRING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ST_COLLECT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBCLASS_ORIGIN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBDATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBJECT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBPARTITION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBPARTITIONS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBSTR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUBSTRING
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUPER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SUSPEND
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SWAPS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SWITCHES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SYSDATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SYSTEM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=SYSTEM_USER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TABLE
< error
| 1064 42000 …near 'TABLE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TABLES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TABLESAMPLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TABLESPACE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TABLE_CHECKSUM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TABLE_NAME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TEMPORARY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TEMPTABLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TERMINATED
< error
| 1064 42000 …near 'TERMINATED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TEXT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=THAN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=THEN
< error
| 1064 42000 …near 'THEN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=THREAD_PRIORITY
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TIES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TIME
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TIMESTAMP
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TIMESTAMPADD
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TIMESTAMPDIFF
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TINYBLOB
< error
| 1064 42000 …near 'TINYBLOB' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TINYINT
< error
| 1064 42000 …near 'TINYINT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TINYTEXT
< error
| 1064 42000 …near 'TINYTEXT' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TLS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TO
< error
| 1064 42000 …near 'TO' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRAILING
< error
| 1064 42000 …near 'TRAILING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRANSACTION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRIGGER
< error
| 1064 42000 …near 'TRIGGER' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRIGGERS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRIM
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRUE
< error
| 1064 42000 …near 'TRUE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TRUNCATE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TYPE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=TYPES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNBOUNDED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNCOMMITTED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNDEFINED
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNDO
< error
| 1064 42000 …near 'UNDO' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNDOFILE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNDO_BUFFER_SIZE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNICODE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNINSTALL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNION
< error
| 1064 42000 …near 'UNION' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNIQUE
< error
| 1064 42000 …near 'UNIQUE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNKNOWN
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNLOCK
< error
| 1064 42000 …near 'UNLOCK' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNREGISTER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNSIGNED
< error
| 1064 42000 …near 'UNSIGNED' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UNTIL
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UPDATE
< error
| 1064 42000 …near 'UPDATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UPGRADE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=URI
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=USAGE
< error
| 1064 42000 …near 'USAGE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=USE
< error
| 1064 42000 …near 'USE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=USER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=USING
< error
| 1064 42000 …near 'USING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UTC_DATE
< error
| 1064 42000 …near 'UTC_DATE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UTC_TIME
< error
| 1064 42000 …near 'UTC_TIME' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=UTC_TIMESTAMP
< error
| 1064 42000 …near 'UTC_TIMESTAMP' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VALIDATION
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VALUE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VALUES
< error
| 1064 42000 …near 'VALUES' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VARBINARY
< error
| 1064 42000 …near 'VARBINARY' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VARCHAR
< error
| 1064 42000 …near 'VARCHAR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VARIABLES
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VARYING
< error
| 1064 42000 …near 'VARYING' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VIEW
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VIRTUAL
< error
| 1064 42000 …near 'VIRTUAL' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=VISIBLE
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WAIT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WARNINGS
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WHEN
< error
| 1064 42000 …near 'WHEN' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WHERE
< error
| 1064 42000 …near 'WHERE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WHILE
< error
| 1064 42000 …near 'WHILE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WINDOW
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WITH
< error
| 1064 42000 …near 'WITH' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WITHOUT
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WORK
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WRAPPER
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=WRITE
< error
| 1064 42000 …near 'WRITE' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=X509
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=XA
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=XID
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=XML
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=XOR
< error
| 1064 42000 …near 'XOR' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=YEAR
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=YEAR_MONTH
< error
| 1064 42000 …near 'YEAR_MONTH' at line 1
> exec ^
: string
| CREATE TABLE kw_eng_test (id INT) ENGINE=ZEROFILL
< error
| 1064 42000 …near 'ZEROFILL' at line 1
> exec
| CREATE TABLE kw_idx_base (id INT, val INT)
< ok 0 0
> exec
| PREPARE omni_kw_stmt FROM ?
: string
| CREATE INDEX ABSENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ACCESSIBLE ON kw_idx_base (val)
< error
| 1064 42000 …near 'ACCESSIBLE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ACCOUNT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ACTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ACTIVE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ADD ON kw_idx_base (val)
< error
| 1064 42000 …near 'ADD ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ADDDATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ADMIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AFTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AGAINST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AGGREGATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ALGORITHM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ALL ON kw_idx_base (val)
< error
| 1064 42000 …near 'ALL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ALLOW_MISSING_FILES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ALTER ON kw_idx_base (val)
< error
| 1064 42000 …near 'ALTER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ALWAYS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ANALYZE ON kw_idx_base (val)
< error
| 1064 42000 …near 'ANALYZE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX AND ON kw_idx_base (val)
< error
| 1064 42000 …near 'AND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ANY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ARRAY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AS ON kw_idx_base (val)
< error
| 1064 42000 …near 'AS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ASC ON kw_idx_base (val)
< error
| 1064 42000 …near 'ASC ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ASCII ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ASENSITIVE ON kw_idx_base (val)
< error
| 1064 42000 …near 'ASENSITIVE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ATTRIBUTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AUTHENTICATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AUTO ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AUTOEXTEND_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AUTO_INCREMENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AUTO_REFRESH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AUTO_REFRESH_SOURCE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AVG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX AVG_ROW_LENGTH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BACKUP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BEFORE ON kw_idx_base (val)
< error
| 1064 42000 …near 'BEFORE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BEGIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BERNOULLI ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BETWEEN ON kw_idx_base (val)
< error
| 1064 42000 …near 'BETWEEN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BIGINT ON kw_idx_base (val)
< error
| 1064 42000 …near 'BIGINT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BINARY ON kw_idx_base (val)
< error
| 1064 42000 …near 'BINARY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BINLOG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BIT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BIT_AND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BIT_OR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BIT_XOR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BLOB ON kw_idx_base (val)
< error
| 1064 42000 …near 'BLOB ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BLOCK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BOOL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BOOLEAN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BOTH ON kw_idx_base (val)
< error
| 1064 42000 …near 'BOTH ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BTREE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BUCKETS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BULK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX BY ON kw_idx_base (val)
< error
| 1064 42000 …near 'BY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX BYTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CACHE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CALL ON kw_idx_base (val)
< error
| 1064 42000 …near 'CALL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CASCADE ON kw_idx_base (val)
< error
| 1064 42000 …near 'CASCADE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CASCADED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CASE ON kw_idx_base (val)
< error
| 1064 42000 …near 'CASE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CAST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CATALOG_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CHAIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CHALLENGE_RESPONSE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CHANGE ON kw_idx_base (val)
< error
| 1064 42000 …near 'CHANGE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CHANGED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CHANNEL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CHAR ON kw_idx_base (val)
< error
| 1064 42000 …near 'CHAR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CHARACTER ON kw_idx_base (val)
< error
| 1064 42000 …near 'CHARACTER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CHARSET ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CHECK ON kw_idx_base (val)
< error
| 1064 42000 …near 'CHECK ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CHECKSUM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CIPHER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CLASS_ORIGIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CLIENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CLONE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CLOSE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COALESCE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CODE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COLLATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'COLLATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX COLLATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COLUMN ON kw_idx_base (val)
< error
| 1064 42000 …near 'COLUMN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX COLUMNS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COLUMN_FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COLUMN_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMMENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMMIT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMMITTED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMPACT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMPLETION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMPONENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMPRESSED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX COMPRESSION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONCAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONCURRENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONDITION ON kw_idx_base (val)
< error
| 1064 42000 …near 'CONDITION ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CONNECTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONSISTENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONSTRAINT ON kw_idx_base (val)
< error
| 1064 42000 …near 'CONSTRAINT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CONSTRAINT_CATALOG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONSTRAINT_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONSTRAINT_SCHEMA ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONTAINS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONTEXT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CONTINUE ON kw_idx_base (val)
< error
| 1064 42000 …near 'CONTINUE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CONVERT ON kw_idx_base (val)
< error
| 1064 42000 …near 'CONVERT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX COUNT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CPU ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CREATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'CREATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CROSS ON kw_idx_base (val)
< error
| 1064 42000 …near 'CROSS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CSV ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CUBE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CUME_DIST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CURDATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CURRENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CURRENT_DATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'CURRENT_DATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CURRENT_TIME ON kw_idx_base (val)
< error
| 1064 42000 …near 'CURRENT_TIME ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CURRENT_TIMESTAMP ON kw_idx_base (val)
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CURRENT_USER ON kw_idx_base (val)
< error
| 1064 42000 …near 'CURRENT_USER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CURSOR ON kw_idx_base (val)
< error
| 1064 42000 …near 'CURSOR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX CURSOR_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX CURTIME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DATA ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DATABASE ON kw_idx_base (val)
< error
| 1064 42000 …near 'DATABASE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DATABASES ON kw_idx_base (val)
< error
| 1064 42000 …near 'DATABASES ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DATAFILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DATETIME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DATE_ADD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DATE_SUB ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DAY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DAY_HOUR ON kw_idx_base (val)
< error
| 1064 42000 …near 'DAY_HOUR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DAY_MICROSECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'DAY_MICROSECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DAY_MINUTE ON kw_idx_base (val)
< error
| 1064 42000 …near 'DAY_MINUTE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DAY_SECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'DAY_SECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DEALLOCATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DEC ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DECIMAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'DECIMAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DECLARE ON kw_idx_base (val)
< error
| 1064 42000 …near 'DECLARE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DEFAULT ON kw_idx_base (val)
< error
| 1064 42000 …near 'DEFAULT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DEFAULT_AUTH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DEFINER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DEFINITION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DELAYED ON kw_idx_base (val)
< error
| 1064 42000 …near 'DELAYED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DELAY_KEY_WRITE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DELETE ON kw_idx_base (val)
< error
| 1064 42000 …near 'DELETE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DENSE_RANK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DESC ON kw_idx_base (val)
< error
| 1064 42000 …near 'DESC ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DESCRIBE ON kw_idx_base (val)
< error
| 1064 42000 …near 'DESCRIBE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DESCRIPTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DETERMINISTIC ON kw_idx_base (val)
< error
| 1064 42000 …near 'DETERMINISTIC ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DIAGNOSTICS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DIRECTORY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DISABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DISCARD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DISK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DISTINCT ON kw_idx_base (val)
< error
| 1064 42000 …near 'DISTINCT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DISTINCTROW ON kw_idx_base (val)
< error
| 1064 42000 …near 'DISTINCTROW ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DIV ON kw_idx_base (val)
< error
| 1064 42000 …near 'DIV ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DO ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DOUBLE ON kw_idx_base (val)
< error
| 1064 42000 …near 'DOUBLE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DROP ON kw_idx_base (val)
< error
| 1064 42000 …near 'DROP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DUAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'DUAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX DUALITY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DUMPFILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DUPLICATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX DYNAMIC ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EACH ON kw_idx_base (val)
< error
| 1064 42000 …near 'EACH ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ELSE ON kw_idx_base (val)
< error
| 1064 42000 …near 'ELSE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ELSEIF ON kw_idx_base (val)
< error
| 1064 42000 …near 'ELSEIF ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX EMPTY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENCLOSED ON kw_idx_base (val)
< error
| 1064 42000 …near 'ENCLOSED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ENCRYPTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX END ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENFORCED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENGINE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENGINES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENGINE_ATTRIBUTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ENUM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ERROR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ERRORS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ESCAPE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ESCAPED ON kw_idx_base (val)
< error
| 1064 42000 …near 'ESCAPED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX EVENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EVENTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EVERY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXCEPT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXCHANGE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXCLUDE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXECUTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXISTS ON kw_idx_base (val)
< error
| 1064 42000 …near 'EXISTS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX EXIT ON kw_idx_base (val)
< error
| 1064 42000 …near 'EXIT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX EXPANSION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXPIRE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXPLAIN ON kw_idx_base (val)
< error
| 1064 42000 …near 'EXPLAIN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX EXPORT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXTENDED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXTENT_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXTERNAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXTERNAL_FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX EXTRACT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FACTOR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FAILED_LOGIN_ATTEMPTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FALSE ON kw_idx_base (val)
< error
| 1064 42000 …near 'FALSE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FAST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FAULTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FETCH ON kw_idx_base (val)
< error
| 1064 42000 …near 'FETCH ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FIELDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILE_BLOCK_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILE_FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILE_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILE_PATTERN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILE_PREFIX ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FILTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FINISH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FIRST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FIRST_VALUE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FIXED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FLOAT ON kw_idx_base (val)
< error
| 1064 42000 …near 'FLOAT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FLOAT4 ON kw_idx_base (val)
< error
| 1064 42000 …near 'FLOAT4 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FLOAT8 ON kw_idx_base (val)
< error
| 1064 42000 …near 'FLOAT8 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FLUSH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FOLLOWING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FOLLOWS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FOR ON kw_idx_base (val)
< error
| 1064 42000 …near 'FOR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FORCE ON kw_idx_base (val)
< error
| 1064 42000 …near 'FORCE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FOREIGN ON kw_idx_base (val)
< error
| 1064 42000 …near 'FOREIGN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FOUND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FROM ON kw_idx_base (val)
< error
| 1064 42000 …near 'FROM ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FULL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX FULLTEXT ON kw_idx_base (val)
< error
| 1064 42000 …near 'FULLTEXT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX FUNCTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GENERAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GENERATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GENERATED ON kw_idx_base (val)
< error
| 1064 42000 …near 'GENERATED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX GEOMCOLLECTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GEOMETRY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GEOMETRYCOLLECTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GET ON kw_idx_base (val)
< error
| 1064 42000 …near 'GET ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX GET_FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GET_SOURCE_PUBLIC_KEY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GLOBAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GRANT ON kw_idx_base (val)
< error
| 1064 42000 …near 'GRANT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX GRANTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GREATEST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GROUP ON kw_idx_base (val)
< error
| 1064 42000 …near 'GROUP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX GROUPING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GROUPS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GROUP_CONCAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GROUP_REPLICATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GTIDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GTID_ONLY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX GUIDED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HANDLER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HASH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HAVING ON kw_idx_base (val)
< error
| 1064 42000 …near 'HAVING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX HEADER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HELP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HIGH_PRIORITY ON kw_idx_base (val)
< error
| 1064 42000 …near 'HIGH_PRIORITY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX HISTOGRAM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HISTORY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HOST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HOSTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HOUR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX HOUR_MICROSECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'HOUR_MICROSECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX HOUR_MINUTE ON kw_idx_base (val)
< error
| 1064 42000 …near 'HOUR_MINUTE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX HOUR_SECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'HOUR_SECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX IDENTIFIED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IF ON kw_idx_base (val)
< error
| 1064 42000 …near 'ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX IGNORE ON kw_idx_base (val)
< error
| 1064 42000 …near 'IGNORE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX IGNORE_SERVER_IDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IMPORT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IN ON kw_idx_base (val)
< error
| 1064 42000 …near 'IN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INACTIVE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INDEX ON kw_idx_base (val)
< error
| 1064 42000 …near 'INDEX ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INDEXES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INFILE ON kw_idx_base (val)
< error
| 1064 42000 …near 'INFILE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INITIAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INITIAL_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INITIATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INNER ON kw_idx_base (val)
< error
| 1064 42000 …near 'INNER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INOUT ON kw_idx_base (val)
< error
| 1064 42000 …near 'INOUT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INSENSITIVE ON kw_idx_base (val)
< error
| 1064 42000 …near 'INSENSITIVE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INSERT ON kw_idx_base (val)
< error
| 1064 42000 …near 'INSERT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INSERT_METHOD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INSTALL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INSTANCE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INT ON kw_idx_base (val)
< error
| 1064 42000 …near 'INT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INT1 ON kw_idx_base (val)
< error
| 1064 42000 …near 'INT1 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INT2 ON kw_idx_base (val)
< error
| 1064 42000 …near 'INT2 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INT3 ON kw_idx_base (val)
< error
| 1064 42000 …near 'INT3 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INT4 ON kw_idx_base (val)
< error
| 1064 42000 …near 'INT4 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INT8 ON kw_idx_base (val)
< error
| 1064 42000 …near 'INT8 ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INTEGER ON kw_idx_base (val)
< error
| 1064 42000 …near 'INTEGER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INTERSECT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INTERVAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'INTERVAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INTO ON kw_idx_base (val)
< error
| 1064 42000 …near 'INTO ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX INVISIBLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX INVOKER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IO ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IO_AFTER_GTIDS ON kw_idx_base (val)
< error
| 1064 42000 …near 'IO_AFTER_GTIDS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX IO_BEFORE_GTIDS ON kw_idx_base (val)
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX IO_THREAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IPC ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX IS ON kw_idx_base (val)
< error
| 1064 42000 …near 'IS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ISOLATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ISSUER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ITERATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'ITERATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX JOIN ON kw_idx_base (val)
< error
| 1064 42000 …near 'JOIN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX JSON ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX JSON_ARRAYAGG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX JSON_DUALITY_OBJECT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX JSON_OBJECTAGG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX JSON_TABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX JSON_VALUE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX KEY ON kw_idx_base (val)
< error
| 1064 42000 …near 'KEY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX KEYRING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX KEYS ON kw_idx_base (val)
< error
| 1064 42000 …near 'KEYS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX KEY_BLOCK_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX KILL ON kw_idx_base (val)
< error
| 1064 42000 …near 'KILL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LAG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LANGUAGE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LAST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LAST_VALUE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LATERAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LEAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LEADING ON kw_idx_base (val)
< error
| 1064 42000 …near 'LEADING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LEAST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LEAVE ON kw_idx_base (val)
< error
| 1064 42000 …near 'LEAVE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LEAVES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LEFT ON kw_idx_base (val)
< error
| 1064 42000 …near 'LEFT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LESS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LEVEL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LIBRARY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LIKE ON kw_idx_base (val)
< error
| 1064 42000 …near 'LIKE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LIMIT ON kw_idx_base (val)
< error
| 1064 42000 …near 'LIMIT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LINEAR ON kw_idx_base (val)
< error
| 1064 42000 …near 'LINEAR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LINES ON kw_idx_base (val)
< error
| 1064 42000 …near 'LINES ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LINESTRING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LIST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LOAD ON kw_idx_base (val)
< error
| 1064 42000 …near 'LOAD ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LOCAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LOCALTIME ON kw_idx_base (val)
< error
| 1064 42000 …near 'LOCALTIME ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LOCALTIMESTAMP ON kw_idx_base (val)
< error
| 1064 42000 …near 'LOCALTIMESTAMP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LOCK ON kw_idx_base (val)
< error
| 1064 42000 …near 'LOCK ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LOCKED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LOCKS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LOG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LOGFILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LOGS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX LONG ON kw_idx_base (val)
< error
| 1064 42000 …near 'LONG ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LONGBLOB ON kw_idx_base (val)
< error
| 1064 42000 …near 'LONGBLOB ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LONGTEXT ON kw_idx_base (val)
< error
| 1064 42000 …near 'LONGTEXT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LOOP ON kw_idx_base (val)
< error
| 1064 42000 …near 'LOOP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX LOW_PRIORITY ON kw_idx_base (val)
< error
| 1064 42000 …near 'LOW_PRIORITY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MANIFEST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MANUAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MASTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MATCH ON kw_idx_base (val)
< error
| 1064 42000 …near 'MATCH ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MATERIALIZED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAXVALUE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX_CONNECTIONS_PER_HOUR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX_QUERIES_PER_HOUR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX_ROWS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX_UPDATES_PER_HOUR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MAX_USER_CONNECTIONS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MEDIUM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MEDIUMBLOB ON kw_idx_base (val)
< error
| 1064 42000 …near 'MEDIUMBLOB ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MEDIUMINT ON kw_idx_base (val)
< error
| 1064 42000 …near 'MEDIUMINT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MEDIUMTEXT ON kw_idx_base (val)
< error
| 1064 42000 …near 'MEDIUMTEXT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MEMBER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MEMORY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MERGE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MESSAGE_TEXT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MICROSECOND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MID ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MIDDLEINT ON kw_idx_base (val)
< error
| 1064 42000 …near 'MIDDLEINT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MIGRATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MINUTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MINUTE_MICROSECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'MINUTE_MICROSECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MINUTE_SECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'MINUTE_SECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MIN_ROWS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MOD ON kw_idx_base (val)
< error
| 1064 42000 …near 'MOD ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MODE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MODIFIES ON kw_idx_base (val)
< error
| 1064 42000 …near 'MODIFIES ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX MODIFY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MONTH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MULTILINESTRING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MULTIPOINT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MULTIPOLYGON ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MUTEX ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX MYSQL_ERRNO ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NAMES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NATIONAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NATURAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'NATURAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX NCHAR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NDB ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NDBCLUSTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NESTED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NETWORK_NAMESPACE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NEVER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NEW ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NEXT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NO ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NODEGROUP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NONE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NOT ON kw_idx_base (val)
< error
| 1064 42000 …near 'NOT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX NOW ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NOWAIT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NO_WAIT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NO_WRITE_TO_BINLOG ON kw_idx_base (val)
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX NTH_VALUE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NTILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NULL ON kw_idx_base (val)
< error
| 1064 42000 …near 'NULL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX NULLIF ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NULLS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NUMBER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX NUMERIC ON kw_idx_base (val)
< error
| 1064 42000 …near 'NUMERIC ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX NVARCHAR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OF ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OFF ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OFFSET ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OJ ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OLD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ON ON kw_idx_base (val)
< error
| 1064 42000 …near 'ON ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ONE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ONLY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OPEN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OPTIMIZE ON kw_idx_base (val)
< error
| 1064 42000 …near 'OPTIMIZE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX OPTIMIZER_COSTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OPTION ON kw_idx_base (val)
< error
| 1064 42000 …near 'OPTION ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX OPTIONAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OPTIONALLY ON kw_idx_base (val)
< error
| 1064 42000 …near 'OPTIONALLY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX OPTIONS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OR ON kw_idx_base (val)
< error
| 1064 42000 …near 'OR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ORDER ON kw_idx_base (val)
< error
| 1064 42000 …near 'ORDER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ORDINALITY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ORGANIZATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OTHERS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OUT ON kw_idx_base (val)
< error
| 1064 42000 …near 'OUT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX OUTER ON kw_idx_base (val)
< error
| 1064 42000 …near 'OUTER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX OUTFILE ON kw_idx_base (val)
< error
| 1064 42000 …near 'OUTFILE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX OVER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OVERLAY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OVERWRITE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX OWNER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PACK_KEYS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PAGE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARALLEL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARAMETERS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARSER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARSE_TREE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARTIAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARTITION ON kw_idx_base (val)
< error
| 1064 42000 …near 'PARTITION ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX PARTITIONING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PARTITIONS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PASSWORD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PASSWORD_LOCK_TIME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PATH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PERCENT_RANK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PERSIST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PERSIST_ONLY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PHASE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PLACING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PLUGIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PLUGINS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PLUGIN_DIR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX POINT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX POLYGON ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PORT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX POSITION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PRECEDES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PRECEDING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PRECISION ON kw_idx_base (val)
< error
| 1064 42000 …near 'PRECISION ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX PREFIX ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PREPARE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PRESERVE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PREV ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PRIMARY ON kw_idx_base (val)
< error
| 1064 42000 …near 'PRIMARY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX PRIVILEGES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PRIVILEGE_CHECKS_USER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PROCEDURE ON kw_idx_base (val)
< error
| 1064 42000 …near 'PROCEDURE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX PROCESS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PROCESSLIST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PROFILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PROFILES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PROXY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX PURGE ON kw_idx_base (val)
< error
| 1064 42000 …near 'PURGE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX QUALIFY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX QUARTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX QUERY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX QUICK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RANDOM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RANGE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RANK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX READ ON kw_idx_base (val)
< error
| 1064 42000 …near 'READ ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX READS ON kw_idx_base (val)
< error
| 1064 42000 …near 'READS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX READ_ONLY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX READ_WRITE ON kw_idx_base (val)
< error
| 1064 42000 …near 'READ_WRITE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'REAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REBUILD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RECOVER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RECURSIVE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REDO_BUFFER_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REDUNDANT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REFERENCE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REFERENCES ON kw_idx_base (val)
< error
| 1064 42000 …near 'REFERENCES ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REGEXP ON kw_idx_base (val)
< error
| 1064 42000 …near 'REGEXP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REGISTRATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELATIONAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELAY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELAYLOG ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELAY_LOG_FILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELAY_LOG_POS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELAY_THREAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RELEASE ON kw_idx_base (val)
< error
| 1064 42000 …near 'RELEASE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX RELOAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REMOVE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RENAME ON kw_idx_base (val)
< error
| 1064 42000 …near 'RENAME ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REORGANIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPAIR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPEAT ON kw_idx_base (val)
< error
| 1064 42000 …near 'REPEAT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REPEATABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLACE ON kw_idx_base (val)
< error
| 1064 42000 …near 'REPLACE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REPLICA ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICAS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_DO_DB ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_DO_TABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_IGNORE_DB ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_IGNORE_TABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_REWRITE_DB ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_WILD_DO_TABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATE_WILD_IGNORE_TABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REPLICATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REQUIRE ON kw_idx_base (val)
< error
| 1064 42000 …near 'REQUIRE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX REQUIRE_ROW_FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REQUIRE_TABLE_PRIMARY_KEY_CHECK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RESET ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RESIGNAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'RESIGNAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX RESOURCE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RESPECT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RESTART ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RESTORE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RESTRICT ON kw_idx_base (val)
< error
| 1064 42000 …near 'RESTRICT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX RESUME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RETAIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RETURN ON kw_idx_base (val)
< error
| 1064 42000 …near 'RETURN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX RETURNED_SQLSTATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RETURNING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RETURNS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REUSE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REVERSE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX REVOKE ON kw_idx_base (val)
< error
| 1064 42000 …near 'REVOKE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX RIGHT ON kw_idx_base (val)
< error
| 1064 42000 …near 'RIGHT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX RLIKE ON kw_idx_base (val)
< error
| 1064 42000 …near 'RLIKE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ROLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROLLBACK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROLLUP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROTATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROUTINE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROW ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROWS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROW_COUNT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROW_FORMAT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ROW_NUMBER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX RTREE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX S3 ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SAVEPOINT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SCHEDULE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SCHEMA ON kw_idx_base (val)
< error
| 1064 42000 …near 'SCHEMA ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SCHEMAS ON kw_idx_base (val)
< error
| 1064 42000 …near 'SCHEMAS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SCHEMA_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECOND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECONDARY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECONDARY_ENGINE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECONDARY_ENGINE_ATTRIBUTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECONDARY_LOAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECONDARY_UNLOAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SECOND_MICROSECOND ON kw_idx_base (val)
< error
| 1064 42000 …near 'SECOND_MICROSECOND ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SECURITY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SELECT ON kw_idx_base (val)
< error
| 1064 42000 …near 'SELECT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SENSITIVE ON kw_idx_base (val)
< error
| 1064 42000 …near 'SENSITIVE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SEPARATOR ON kw_idx_base (val)
< error
| 1064 42000 …near 'SEPARATOR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SERIAL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SERIALIZABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SERVER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SESSION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SESSION_USER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SET ON kw_idx_base (val)
< error
| 1064 42000 …near 'SET ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SETS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SHARE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SHOW ON kw_idx_base (val)
< error
| 1064 42000 …near 'SHOW ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SHUTDOWN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SIGNAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'SIGNAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SIGNED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SIMPLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SKIP ON kw_idx_base (val)
< error
| 1064 42000 …near 'SKIP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SLAVE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SLOW ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SMALLINT ON kw_idx_base (val)
< error
| 1064 42000 …near 'SMALLINT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SNAPSHOT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOCKET ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SONAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOUNDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_AUTO_POSITION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_BIND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_COMPRESSION_ALGORITHMS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_CONNECTION_AUTO_FAILOVER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_CONNECT_RETRY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_DELAY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_HEARTBEAT_PERIOD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_HOST ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_LOG_FILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_LOG_POS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_PASSWORD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_PORT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_PUBLIC_KEY_PATH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_RETRY_COUNT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_CA ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_CAPATH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_CERT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_CIPHER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_CRL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_CRLPATH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_KEY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_SSL_VERIFY_SERVER_CERT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_TLS_CIPHERSUITES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_TLS_VERSION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_USER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SOURCE_ZSTD_COMPRESSION_LEVEL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SPATIAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'SPATIAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SPECIFIC ON kw_idx_base (val)
< error
| 1064 42000 …near 'SPECIFIC ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQL ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQLEXCEPTION ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQLEXCEPTION ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQLSTATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQLSTATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQLWARNING ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQLWARNING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQL_AFTER_GTIDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_AFTER_MTS_GAPS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_BEFORE_GTIDS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_BIG_RESULT ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQL_BIG_RESULT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQL_BUFFER_RESULT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_CALC_FOUND_ROWS ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQL_NO_CACHE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_SMALL_RESULT ON kw_idx_base (val)
< error
| 1064 42000 …near 'SQL_SMALL_RESULT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX SQL_THREAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_DAY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_HOUR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_MINUTE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_MONTH ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_QUARTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_SECOND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_WEEK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SQL_TSI_YEAR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SRID ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SSE_KMS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SSE_S3 ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SSL ON kw_idx_base (val)
< error
| 1064 42000 …near 'SSL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX STACKED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX START ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STARTING ON kw_idx_base (val)
< error
| 1064 42000 …near 'STARTING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX STARTS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STATS_AUTO_RECALC ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STATS_PERSISTENT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STATS_SAMPLE_PAGES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STATUS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STDDEV ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STDDEV_POP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STDDEV_SAMP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STOP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STORAGE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STORED ON kw_idx_base (val)
< error
| 1064 42000 …near 'STORED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX STRAIGHT_JOIN ON kw_idx_base (val)
< error
| 1064 42000 …near 'STRAIGHT_JOIN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX STREAM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STRICT_LOAD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX STRING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX ST_COLLECT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBCLASS_ORIGIN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBDATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBJECT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBPARTITION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBPARTITIONS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBSTR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUBSTRING ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUPER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SUSPEND ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SWAPS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SWITCHES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SYSDATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SYSTEM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX SYSTEM_USER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TABLE ON kw_idx_base (val)
< error
| 1064 42000 …near 'TABLE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TABLES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TABLESAMPLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TABLESPACE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TABLE_CHECKSUM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TABLE_NAME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TEMPORARY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TEMPTABLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TERMINATED ON kw_idx_base (val)
< error
| 1064 42000 …near 'TERMINATED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TEXT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX THAN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX THEN ON kw_idx_base (val)
< error
| 1064 42000 …near 'THEN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX THREAD_PRIORITY ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TIES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TIME ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TIMESTAMP ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TIMESTAMPADD ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TIMESTAMPDIFF ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TINYBLOB ON kw_idx_base (val)
< error
| 1064 42000 …near 'TINYBLOB ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TINYINT ON kw_idx_base (val)
< error
| 1064 42000 …near 'TINYINT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TINYTEXT ON kw_idx_base (val)
< error
| 1064 42000 …near 'TINYTEXT ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TLS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TO ON kw_idx_base (val)
< error
| 1064 42000 …near 'TO ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TRAILING ON kw_idx_base (val)
< error
| 1064 42000 …near 'TRAILING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TRANSACTION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TRIGGER ON kw_idx_base (val)
< error
| 1064 42000 …near 'TRIGGER ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TRIGGERS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TRIM ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TRUE ON kw_idx_base (val)
< error
| 1064 42000 …near 'TRUE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX TRUNCATE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TYPE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX TYPES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNBOUNDED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNCOMMITTED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNDEFINED ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNDO ON kw_idx_base (val)
< error
| 1064 42000 …near 'UNDO ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UNDOFILE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNDO_BUFFER_SIZE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNICODE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNINSTALL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNION ON kw_idx_base (val)
< error
| 1064 42000 …near 'UNION ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UNIQUE ON kw_idx_base (val)
< error
| 1064 42000 …near 'UNIQUE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UNKNOWN ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNLOCK ON kw_idx_base (val)
< error
| 1064 42000 …near 'UNLOCK ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UNREGISTER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UNSIGNED ON kw_idx_base (val)
< error
| 1064 42000 …near 'UNSIGNED ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UNTIL ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX UPDATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'UPDATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UPGRADE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX URI ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX USAGE ON kw_idx_base (val)
< error
| 1064 42000 …near 'USAGE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX USE ON kw_idx_base (val)
< error
| 1064 42000 …near 'USE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX USER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX USING ON kw_idx_base (val)
< error
| 1064 42000 …near 'USING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UTC_DATE ON kw_idx_base (val)
< error
| 1064 42000 …near 'UTC_DATE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UTC_TIME ON kw_idx_base (val)
< error
| 1064 42000 …near 'UTC_TIME ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX UTC_TIMESTAMP ON kw_idx_base (val)
< error
| 1064 42000 …near 'UTC_TIMESTAMP ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX VALIDATION ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX VALUE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX VALUES ON kw_idx_base (val)
< error
| 1064 42000 …near 'VALUES ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX VARBINARY ON kw_idx_base (val)
< error
| 1064 42000 …near 'VARBINARY ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX VARCHAR ON kw_idx_base (val)
< error
| 1064 42000 …near 'VARCHAR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX VARIABLES ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX VARYING ON kw_idx_base (val)
< error
| 1064 42000 …near 'VARYING ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX VIEW ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX VIRTUAL ON kw_idx_base (val)
< error
| 1064 42000 …near 'VIRTUAL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX VISIBLE ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WAIT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WARNINGS ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WHEN ON kw_idx_base (val)
< error
| 1064 42000 …near 'WHEN ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX WHERE ON kw_idx_base (val)
< error
| 1064 42000 …near 'WHERE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX WHILE ON kw_idx_base (val)
< error
| 1064 42000 …near 'WHILE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX WINDOW ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WITH ON kw_idx_base (val)
< error
| 1064 42000 …near 'WITH ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX WITHOUT ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WORK ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WRAPPER ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX WRITE ON kw_idx_base (val)
< error
| 1064 42000 …near 'WRITE ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX X509 ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX XA ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX XID ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX XML ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX XOR ON kw_idx_base (val)
< error
| 1064 42000 …near 'XOR ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX YEAR ON kw_idx_base (val)
< ok 0 0
> exec ^
: string
| CREATE INDEX YEAR_MONTH ON kw_idx_base (val)
< error
| 1064 42000 …near 'YEAR_MONTH ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE INDEX ZEROFILL ON kw_idx_base (val)
< error
| 1064 42000 …near 'ZEROFILL ON kw_idx_base (val)' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ABSENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ACCESSIBLE UNIQUE (id))
< error
| 1064 42000 …near 'ACCESSIBLE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ACCOUNT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ACTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ACTIVE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ADD UNIQUE (id))
< error
| 1064 42000 …near 'ADD UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ADDDATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ADMIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AFTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AGAINST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AGGREGATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ALGORITHM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ALL UNIQUE (id))
< error
| 1064 42000 …near 'ALL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ALLOW_MISSING_FILES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ALTER UNIQUE (id))
< error
| 1064 42000 …near 'ALTER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ALWAYS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ANALYZE UNIQUE (id))
< error
| 1064 42000 …near 'ANALYZE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AND UNIQUE (id))
< error
| 1064 42000 …near 'AND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ANY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ARRAY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AS UNIQUE (id))
< error
| 1064 42000 …near 'AS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ASC UNIQUE (id))
< error
| 1064 42000 …near 'ASC UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ASCII UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ASENSITIVE UNIQUE (id))
< error
| 1064 42000 …near 'ASENSITIVE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ATTRIBUTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AUTHENTICATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AUTO UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AUTOEXTEND_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AUTO_INCREMENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AUTO_REFRESH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AUTO_REFRESH_SOURCE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AVG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT AVG_ROW_LENGTH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BACKUP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BEFORE UNIQUE (id))
< error
| 1064 42000 …near 'BEFORE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BEGIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BERNOULLI UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BETWEEN UNIQUE (id))
< error
| 1064 42000 …near 'BETWEEN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BIGINT UNIQUE (id))
< error
| 1064 42000 …near 'BIGINT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BINARY UNIQUE (id))
< error
| 1064 42000 …near 'BINARY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BINLOG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BIT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BIT_AND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BIT_OR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BIT_XOR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BLOB UNIQUE (id))
< error
| 1064 42000 …near 'BLOB UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BLOCK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BOOL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BOOLEAN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BOTH UNIQUE (id))
< error
| 1064 42000 …near 'BOTH UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BTREE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BUCKETS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BULK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BY UNIQUE (id))
< error
| 1064 42000 …near 'BY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT BYTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CACHE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CALL UNIQUE (id))
< error
| 1064 42000 …near 'CALL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CASCADE UNIQUE (id))
< error
| 1064 42000 …near 'CASCADE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CASCADED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CASE UNIQUE (id))
< error
| 1064 42000 …near 'CASE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CAST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CATALOG_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHAIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHALLENGE_RESPONSE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHANGE UNIQUE (id))
< error
| 1064 42000 …near 'CHANGE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHANGED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHANNEL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHAR UNIQUE (id))
< error
| 1064 42000 …near 'CHAR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHARACTER UNIQUE (id))
< error
| 1064 42000 …near 'CHARACTER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHARSET UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHECK UNIQUE (id))
< error
| 1064 42000 …near 'UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CHECKSUM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CIPHER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CLASS_ORIGIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CLIENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CLONE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CLOSE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COALESCE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CODE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COLLATE UNIQUE (id))
< error
| 1064 42000 …near 'COLLATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COLLATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COLUMN UNIQUE (id))
< error
| 1064 42000 …near 'COLUMN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COLUMNS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COLUMN_FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COLUMN_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMMENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMMIT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMMITTED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMPACT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMPLETION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMPONENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMPRESSED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COMPRESSION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONCAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONCURRENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONDITION UNIQUE (id))
< error
| 1064 42000 …near 'CONDITION UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONNECTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONSISTENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONSTRAINT UNIQUE (id))
< error
| 1064 42000 …near 'CONSTRAINT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONSTRAINT_CATALOG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONSTRAINT_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONSTRAINT_SCHEMA UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONTAINS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONTEXT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONTINUE UNIQUE (id))
< error
| 1064 42000 …near 'CONTINUE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CONVERT UNIQUE (id))
< error
| 1064 42000 …near 'CONVERT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT COUNT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CPU UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CREATE UNIQUE (id))
< error
| 1064 42000 …near 'CREATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CROSS UNIQUE (id))
< error
| 1064 42000 …near 'CROSS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CSV UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CUBE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CUME_DIST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURDATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURRENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURRENT_DATE UNIQUE (id))
< error
| 1064 42000 …near 'CURRENT_DATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURRENT_TIME UNIQUE (id))
< error
| 1064 42000 …near 'CURRENT_TIME UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURRENT_TIMESTAMP UNIQUE (id))
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURRENT_USER UNIQUE (id))
< error
| 1064 42000 …near 'CURRENT_USER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURSOR UNIQUE (id))
< error
| 1064 42000 …near 'CURSOR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURSOR_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT CURTIME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATA UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATABASE UNIQUE (id))
< error
| 1064 42000 …near 'DATABASE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATABASES UNIQUE (id))
< error
| 1064 42000 …near 'DATABASES UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATAFILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATETIME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATE_ADD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DATE_SUB UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DAY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DAY_HOUR UNIQUE (id))
< error
| 1064 42000 …near 'DAY_HOUR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DAY_MICROSECOND UNIQUE (id))
< error
| 1064 42000 …near 'DAY_MICROSECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DAY_MINUTE UNIQUE (id))
< error
| 1064 42000 …near 'DAY_MINUTE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DAY_SECOND UNIQUE (id))
< error
| 1064 42000 …near 'DAY_SECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DEALLOCATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DEC UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DECIMAL UNIQUE (id))
< error
| 1064 42000 …near 'DECIMAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DECLARE UNIQUE (id))
< error
| 1064 42000 …near 'DECLARE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DEFAULT UNIQUE (id))
< error
| 1064 42000 …near 'DEFAULT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DEFAULT_AUTH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DEFINER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DEFINITION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DELAYED UNIQUE (id))
< error
| 1064 42000 …near 'DELAYED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DELAY_KEY_WRITE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DELETE UNIQUE (id))
< error
| 1064 42000 …near 'DELETE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DENSE_RANK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DESC UNIQUE (id))
< error
| 1064 42000 …near 'DESC UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DESCRIBE UNIQUE (id))
< error
| 1064 42000 …near 'DESCRIBE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DESCRIPTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DETERMINISTIC UNIQUE (id))
< error
| 1064 42000 …near 'DETERMINISTIC UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DIAGNOSTICS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DIRECTORY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DISABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DISCARD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DISK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DISTINCT UNIQUE (id))
< error
| 1064 42000 …near 'DISTINCT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DISTINCTROW UNIQUE (id))
< error
| 1064 42000 …near 'DISTINCTROW UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DIV UNIQUE (id))
< error
| 1064 42000 …near 'DIV UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DO UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DOUBLE UNIQUE (id))
< error
| 1064 42000 …near 'DOUBLE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DROP UNIQUE (id))
< error
| 1064 42000 …near 'DROP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DUAL UNIQUE (id))
< error
| 1064 42000 …near 'DUAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DUALITY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DUMPFILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DUPLICATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT DYNAMIC UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EACH UNIQUE (id))
< error
| 1064 42000 …near 'EACH UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ELSE UNIQUE (id))
< error
| 1064 42000 …near 'ELSE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ELSEIF UNIQUE (id))
< error
| 1064 42000 …near 'ELSEIF UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EMPTY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENCLOSED UNIQUE (id))
< error
| 1064 42000 …near 'ENCLOSED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENCRYPTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT END UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENFORCED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENGINE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENGINES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENGINE_ATTRIBUTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ENUM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ERROR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ERRORS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ESCAPE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ESCAPED UNIQUE (id))
< error
| 1064 42000 …near 'ESCAPED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EVENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EVENTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EVERY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXCEPT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXCHANGE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXCLUDE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXECUTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXISTS UNIQUE (id))
< error
| 1064 42000 …near 'EXISTS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXIT UNIQUE (id))
< error
| 1064 42000 …near 'EXIT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXPANSION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXPIRE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXPLAIN UNIQUE (id))
< error
| 1064 42000 …near 'EXPLAIN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXPORT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXTENDED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXTENT_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXTERNAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXTERNAL_FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT EXTRACT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FACTOR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FAILED_LOGIN_ATTEMPTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FALSE UNIQUE (id))
< error
| 1064 42000 …near 'FALSE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FAST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FAULTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FETCH UNIQUE (id))
< error
| 1064 42000 …near 'FETCH UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FIELDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILE_BLOCK_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILE_FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILE_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILE_PATTERN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILE_PREFIX UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FILTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FINISH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FIRST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FIRST_VALUE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FIXED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FLOAT UNIQUE (id))
< error
| 1064 42000 …near 'FLOAT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FLOAT4 UNIQUE (id))
< error
| 1064 42000 …near 'FLOAT4 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FLOAT8 UNIQUE (id))
< error
| 1064 42000 …near 'FLOAT8 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FLUSH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FOLLOWING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FOLLOWS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FOR UNIQUE (id))
< error
| 1064 42000 …near 'FOR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FORCE UNIQUE (id))
< error
| 1064 42000 …near 'FORCE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FOREIGN UNIQUE (id))
< error
| 1064 42000 …near 'UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FOUND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FROM UNIQUE (id))
< error
| 1064 42000 …near 'FROM UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FULL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FULLTEXT UNIQUE (id))
< error
| 1064 42000 …near 'FULLTEXT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT FUNCTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GENERAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GENERATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GENERATED UNIQUE (id))
< error
| 1064 42000 …near 'GENERATED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GEOMCOLLECTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GEOMETRY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GEOMETRYCOLLECTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GET UNIQUE (id))
< error
| 1064 42000 …near 'GET UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GET_FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GET_SOURCE_PUBLIC_KEY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GLOBAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GRANT UNIQUE (id))
< error
| 1064 42000 …near 'GRANT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GRANTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GREATEST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GROUP UNIQUE (id))
< error
| 1064 42000 …near 'GROUP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GROUPING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GROUPS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GROUP_CONCAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GROUP_REPLICATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GTIDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GTID_ONLY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT GUIDED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HANDLER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HASH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HAVING UNIQUE (id))
< error
| 1064 42000 …near 'HAVING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HEADER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HELP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HIGH_PRIORITY UNIQUE (id))
< error
| 1064 42000 …near 'HIGH_PRIORITY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HISTOGRAM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HISTORY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HOST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HOSTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HOUR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HOUR_MICROSECOND UNIQUE (id))
< error
| 1064 42000 …near 'HOUR_MICROSECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HOUR_MINUTE UNIQUE (id))
< error
| 1064 42000 …near 'HOUR_MINUTE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT HOUR_SECOND UNIQUE (id))
< error
| 1064 42000 …near 'HOUR_SECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IDENTIFIED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IF UNIQUE (id))
< error
| 1064 42000 …near 'IF UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IGNORE UNIQUE (id))
< error
| 1064 42000 …near 'IGNORE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IGNORE_SERVER_IDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IMPORT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IN UNIQUE (id))
< error
| 1064 42000 …near 'IN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INACTIVE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INDEX UNIQUE (id))
< error
| 1064 42000 …near 'INDEX UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INDEXES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INFILE UNIQUE (id))
< error
| 1064 42000 …near 'INFILE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INITIAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INITIAL_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INITIATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INNER UNIQUE (id))
< error
| 1064 42000 …near 'INNER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INOUT UNIQUE (id))
< error
| 1064 42000 …near 'INOUT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INSENSITIVE UNIQUE (id))
< error
| 1064 42000 …near 'INSENSITIVE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INSERT UNIQUE (id))
< error
| 1064 42000 …near 'INSERT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INSERT_METHOD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INSTALL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INSTANCE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INT UNIQUE (id))
< error
| 1064 42000 …near 'INT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INT1 UNIQUE (id))
< error
| 1064 42000 …near 'INT1 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INT2 UNIQUE (id))
< error
| 1064 42000 …near 'INT2 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INT3 UNIQUE (id))
< error
| 1064 42000 …near 'INT3 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INT4 UNIQUE (id))
< error
| 1064 42000 …near 'INT4 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INT8 UNIQUE (id))
< error
| 1064 42000 …near 'INT8 UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INTEGER UNIQUE (id))
< error
| 1064 42000 …near 'INTEGER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INTERSECT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INTERVAL UNIQUE (id))
< error
| 1064 42000 …near 'INTERVAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INTO UNIQUE (id))
< error
| 1064 42000 …near 'INTO UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INVISIBLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT INVOKER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IO UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IO_AFTER_GTIDS UNIQUE (id))
< error
| 1064 42000 …near 'IO_AFTER_GTIDS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IO_BEFORE_GTIDS UNIQUE (id))
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IO_THREAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IPC UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT IS UNIQUE (id))
< error
| 1064 42000 …near 'IS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ISOLATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ISSUER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ITERATE UNIQUE (id))
< error
| 1064 42000 …near 'ITERATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JOIN UNIQUE (id))
< error
| 1064 42000 …near 'JOIN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JSON UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JSON_ARRAYAGG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JSON_DUALITY_OBJECT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JSON_OBJECTAGG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JSON_TABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT JSON_VALUE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT KEY UNIQUE (id))
< error
| 1064 42000 …near 'KEY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT KEYRING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT KEYS UNIQUE (id))
< error
| 1064 42000 …near 'KEYS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT KEY_BLOCK_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT KILL UNIQUE (id))
< error
| 1064 42000 …near 'KILL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LAG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LANGUAGE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LAST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LAST_VALUE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LATERAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEADING UNIQUE (id))
< error
| 1064 42000 …near 'LEADING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEAST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEAVE UNIQUE (id))
< error
| 1064 42000 …near 'LEAVE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEAVES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEFT UNIQUE (id))
< error
| 1064 42000 …near 'LEFT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LESS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LEVEL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LIBRARY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LIKE UNIQUE (id))
< error
| 1064 42000 …near 'LIKE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LIMIT UNIQUE (id))
< error
| 1064 42000 …near 'LIMIT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LINEAR UNIQUE (id))
< error
| 1064 42000 …near 'LINEAR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LINES UNIQUE (id))
< error
| 1064 42000 …near 'LINES UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LINESTRING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LIST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOAD UNIQUE (id))
< error
| 1064 42000 …near 'LOAD UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOCAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOCALTIME UNIQUE (id))
< error
| 1064 42000 …near 'LOCALTIME UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOCALTIMESTAMP UNIQUE (id))
< error
| 1064 42000 …near 'LOCALTIMESTAMP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOCK UNIQUE (id))
< error
| 1064 42000 …near 'LOCK UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOCKED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOCKS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOGFILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOGS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LONG UNIQUE (id))
< error
| 1064 42000 …near 'LONG UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LONGBLOB UNIQUE (id))
< error
| 1064 42000 …near 'LONGBLOB UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LONGTEXT UNIQUE (id))
< error
| 1064 42000 …near 'LONGTEXT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOOP UNIQUE (id))
< error
| 1064 42000 …near 'LOOP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT LOW_PRIORITY UNIQUE (id))
< error
| 1064 42000 …near 'LOW_PRIORITY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MANIFEST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MANUAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MASTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MATCH UNIQUE (id))
< error
| 1064 42000 …near 'MATCH UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MATERIALIZED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAXVALUE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX_CONNECTIONS_PER_HOUR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX_QUERIES_PER_HOUR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX_ROWS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX_UPDATES_PER_HOUR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MAX_USER_CONNECTIONS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MEDIUM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MEDIUMBLOB UNIQUE (id))
< error
| 1064 42000 …near 'MEDIUMBLOB UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MEDIUMINT UNIQUE (id))
< error
| 1064 42000 …near 'MEDIUMINT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MEDIUMTEXT UNIQUE (id))
< error
| 1064 42000 …near 'MEDIUMTEXT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MEMBER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MEMORY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MERGE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MESSAGE_TEXT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MICROSECOND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MID UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MIDDLEINT UNIQUE (id))
< error
| 1064 42000 …near 'MIDDLEINT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MIGRATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MINUTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MINUTE_MICROSECOND UNIQUE (id))
< error
| 1064 42000 …near 'MINUTE_MICROSECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MINUTE_SECOND UNIQUE (id))
< error
| 1064 42000 …near 'MINUTE_SECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MIN_ROWS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MOD UNIQUE (id))
< error
| 1064 42000 …near 'MOD UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MODE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MODIFIES UNIQUE (id))
< error
| 1064 42000 …near 'MODIFIES UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MODIFY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MONTH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MULTILINESTRING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MULTIPOINT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MULTIPOLYGON UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MUTEX UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT MYSQL_ERRNO UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NAMES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NATIONAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NATURAL UNIQUE (id))
< error
| 1064 42000 …near 'NATURAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NCHAR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NDB UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NDBCLUSTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NESTED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NETWORK_NAMESPACE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NEVER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NEW UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NEXT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NO UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NODEGROUP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NONE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NOT UNIQUE (id))
< error
| 1064 42000 …near 'NOT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NOW UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NOWAIT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NO_WAIT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NO_WRITE_TO_BINLOG UNIQUE (id))
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NTH_VALUE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NTILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NULL UNIQUE (id))
< error
| 1064 42000 …near 'NULL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NULLIF UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NULLS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NUMBER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NUMERIC UNIQUE (id))
< error
| 1064 42000 …near 'NUMERIC UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT NVARCHAR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OF UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OFF UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OFFSET UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OJ UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OLD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ON UNIQUE (id))
< error
| 1064 42000 …near 'ON UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ONE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ONLY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPEN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPTIMIZE UNIQUE (id))
< error
| 1064 42000 …near 'OPTIMIZE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPTIMIZER_COSTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPTION UNIQUE (id))
< error
| 1064 42000 …near 'OPTION UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPTIONAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPTIONALLY UNIQUE (id))
< error
| 1064 42000 …near 'OPTIONALLY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OPTIONS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OR UNIQUE (id))
< error
| 1064 42000 …near 'OR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ORDER UNIQUE (id))
< error
| 1064 42000 …near 'ORDER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ORDINALITY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ORGANIZATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OTHERS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OUT UNIQUE (id))
< error
| 1064 42000 …near 'OUT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OUTER UNIQUE (id))
< error
| 1064 42000 …near 'OUTER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OUTFILE UNIQUE (id))
< error
| 1064 42000 …near 'OUTFILE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OVER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OVERLAY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OVERWRITE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT OWNER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PACK_KEYS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PAGE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARALLEL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARAMETERS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARSER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARSE_TREE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARTIAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARTITION UNIQUE (id))
< error
| 1064 42000 …near 'PARTITION UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARTITIONING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PARTITIONS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PASSWORD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PASSWORD_LOCK_TIME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PATH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PERCENT_RANK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PERSIST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PERSIST_ONLY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PHASE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PLACING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PLUGIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PLUGINS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PLUGIN_DIR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT POINT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT POLYGON UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PORT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT POSITION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRECEDES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRECEDING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRECISION UNIQUE (id))
< error
| 1064 42000 …near 'PRECISION UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PREFIX UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PREPARE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRESERVE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PREV UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRIMARY UNIQUE (id))
< error
| 1064 42000 …near 'UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRIVILEGES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PRIVILEGE_CHECKS_USER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PROCEDURE UNIQUE (id))
< error
| 1064 42000 …near 'PROCEDURE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PROCESS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PROCESSLIST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PROFILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PROFILES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PROXY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT PURGE UNIQUE (id))
< error
| 1064 42000 …near 'PURGE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT QUALIFY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT QUARTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT QUERY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT QUICK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RANDOM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RANGE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RANK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT READ UNIQUE (id))
< error
| 1064 42000 …near 'READ UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT READS UNIQUE (id))
< error
| 1064 42000 …near 'READS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT READ_ONLY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT READ_WRITE UNIQUE (id))
< error
| 1064 42000 …near 'READ_WRITE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REAL UNIQUE (id))
< error
| 1064 42000 …near 'REAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REBUILD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RECOVER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RECURSIVE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REDO_BUFFER_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REDUNDANT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REFERENCE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REFERENCES UNIQUE (id))
< error
| 1064 42000 …near 'REFERENCES UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REGEXP UNIQUE (id))
< error
| 1064 42000 …near 'REGEXP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REGISTRATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELATIONAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELAY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELAYLOG UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELAY_LOG_FILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELAY_LOG_POS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELAY_THREAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELEASE UNIQUE (id))
< error
| 1064 42000 …near 'RELEASE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RELOAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REMOVE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RENAME UNIQUE (id))
< error
| 1064 42000 …near 'RENAME UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REORGANIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPAIR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPEAT UNIQUE (id))
< error
| 1064 42000 …near 'REPEAT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPEATABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLACE UNIQUE (id))
< error
| 1064 42000 …near 'REPLACE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICA UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICAS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_DO_DB UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_DO_TABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_IGNORE_DB UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_IGNORE_TABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_REWRITE_DB UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_WILD_DO_TABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATE_WILD_IGNORE_TABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REPLICATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REQUIRE UNIQUE (id))
< error
| 1064 42000 …near 'REQUIRE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REQUIRE_ROW_FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REQUIRE_TABLE_PRIMARY_KEY_CHECK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESET UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESIGNAL UNIQUE (id))
< error
| 1064 42000 …near 'RESIGNAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESOURCE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESPECT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESTART UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESTORE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESTRICT UNIQUE (id))
< error
| 1064 42000 …near 'RESTRICT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RESUME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RETAIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RETURN UNIQUE (id))
< error
| 1064 42000 …near 'RETURN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RETURNED_SQLSTATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RETURNING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RETURNS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REUSE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REVERSE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT REVOKE UNIQUE (id))
< error
| 1064 42000 …near 'REVOKE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RIGHT UNIQUE (id))
< error
| 1064 42000 …near 'RIGHT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RLIKE UNIQUE (id))
< error
| 1064 42000 …near 'RLIKE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROLLBACK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROLLUP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROTATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROUTINE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROW UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROWS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROW_COUNT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROW_FORMAT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ROW_NUMBER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT RTREE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT S3 UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SAVEPOINT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SCHEDULE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SCHEMA UNIQUE (id))
< error
| 1064 42000 …near 'SCHEMA UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SCHEMAS UNIQUE (id))
< error
| 1064 42000 …near 'SCHEMAS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SCHEMA_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECOND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECONDARY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECONDARY_ENGINE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECONDARY_ENGINE_ATTRIBUTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECONDARY_LOAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECONDARY_UNLOAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECOND_MICROSECOND UNIQUE (id))
< error
| 1064 42000 …near 'SECOND_MICROSECOND UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SECURITY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SELECT UNIQUE (id))
< error
| 1064 42000 …near 'SELECT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SENSITIVE UNIQUE (id))
< error
| 1064 42000 …near 'SENSITIVE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SEPARATOR UNIQUE (id))
< error
| 1064 42000 …near 'SEPARATOR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SERIAL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SERIALIZABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SERVER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SESSION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SESSION_USER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SET UNIQUE (id))
< error
| 1064 42000 …near 'SET UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SETS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SHARE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SHOW UNIQUE (id))
< error
| 1064 42000 …near 'SHOW UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SHUTDOWN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SIGNAL UNIQUE (id))
< error
| 1064 42000 …near 'SIGNAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SIGNED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SIMPLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SKIP UNIQUE (id))
< error
| 1064 42000 …near 'SKIP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SLAVE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SLOW UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SMALLINT UNIQUE (id))
< error
| 1064 42000 …near 'SMALLINT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SNAPSHOT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOCKET UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SONAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOUNDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_AUTO_POSITION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_BIND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_COMPRESSION_ALGORITHMS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_CONNECTION_AUTO_FAILOVER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_CONNECT_RETRY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_DELAY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_HEARTBEAT_PERIOD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_HOST UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_LOG_FILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_LOG_POS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_PASSWORD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_PORT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_PUBLIC_KEY_PATH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_RETRY_COUNT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_CA UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_CAPATH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_CERT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_CIPHER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_CRL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_CRLPATH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_KEY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_SSL_VERIFY_SERVER_CERT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_TLS_CIPHERSUITES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_TLS_VERSION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_USER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SOURCE_ZSTD_COMPRESSION_LEVEL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SPATIAL UNIQUE (id))
< error
| 1064 42000 …near 'SPATIAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SPECIFIC UNIQUE (id))
< error
| 1064 42000 …near 'SPECIFIC UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL UNIQUE (id))
< error
| 1064 42000 …near 'SQL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQLEXCEPTION UNIQUE (id))
< error
| 1064 42000 …near 'SQLEXCEPTION UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQLSTATE UNIQUE (id))
< error
| 1064 42000 …near 'SQLSTATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQLWARNING UNIQUE (id))
< error
| 1064 42000 …near 'SQLWARNING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_AFTER_GTIDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_AFTER_MTS_GAPS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_BEFORE_GTIDS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_BIG_RESULT UNIQUE (id))
< error
| 1064 42000 …near 'SQL_BIG_RESULT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_BUFFER_RESULT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_CALC_FOUND_ROWS UNIQUE (id))
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_NO_CACHE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_SMALL_RESULT UNIQUE (id))
< error
| 1064 42000 …near 'SQL_SMALL_RESULT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_THREAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_DAY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_HOUR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_MINUTE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_MONTH UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_QUARTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_SECOND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_WEEK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SQL_TSI_YEAR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SRID UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SSE_KMS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SSE_S3 UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SSL UNIQUE (id))
< error
| 1064 42000 …near 'SSL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STACKED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT START UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STARTING UNIQUE (id))
< error
| 1064 42000 …near 'STARTING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STARTS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STATS_AUTO_RECALC UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STATS_PERSISTENT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STATS_SAMPLE_PAGES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STATUS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STDDEV UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STDDEV_POP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STDDEV_SAMP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STOP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STORAGE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STORED UNIQUE (id))
< error
| 1064 42000 …near 'STORED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STRAIGHT_JOIN UNIQUE (id))
< error
| 1064 42000 …near 'STRAIGHT_JOIN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STREAM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STRICT_LOAD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT STRING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ST_COLLECT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBCLASS_ORIGIN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBDATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBJECT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBPARTITION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBPARTITIONS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBSTR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUBSTRING UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUPER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SUSPEND UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SWAPS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SWITCHES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SYSDATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SYSTEM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT SYSTEM_USER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TABLE UNIQUE (id))
< error
| 1064 42000 …near 'TABLE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TABLES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TABLESAMPLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TABLESPACE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TABLE_CHECKSUM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TABLE_NAME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TEMPORARY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TEMPTABLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TERMINATED UNIQUE (id))
< error
| 1064 42000 …near 'TERMINATED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TEXT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT THAN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT THEN UNIQUE (id))
< error
| 1064 42000 …near 'THEN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT THREAD_PRIORITY UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TIES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TIME UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TIMESTAMP UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TIMESTAMPADD UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TIMESTAMPDIFF UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TINYBLOB UNIQUE (id))
< error
| 1064 42000 …near 'TINYBLOB UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TINYINT UNIQUE (id))
< error
| 1064 42000 …near 'TINYINT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TINYTEXT UNIQUE (id))
< error
| 1064 42000 …near 'TINYTEXT UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TLS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TO UNIQUE (id))
< error
| 1064 42000 …near 'TO UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRAILING UNIQUE (id))
< error
| 1064 42000 …near 'TRAILING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRANSACTION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRIGGER UNIQUE (id))
< error
| 1064 42000 …near 'TRIGGER UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRIGGERS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRIM UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRUE UNIQUE (id))
< error
| 1064 42000 …near 'TRUE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TRUNCATE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TYPE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT TYPES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNBOUNDED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNCOMMITTED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNDEFINED UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNDO UNIQUE (id))
< error
| 1064 42000 …near 'UNDO UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNDOFILE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNDO_BUFFER_SIZE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNICODE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNINSTALL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNION UNIQUE (id))
< error
| 1064 42000 …near 'UNION UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNIQUE UNIQUE (id))
< error
| 1064 42000 …near 'UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNKNOWN UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNLOCK UNIQUE (id))
< error
| 1064 42000 …near 'UNLOCK UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNREGISTER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNSIGNED UNIQUE (id))
< error
| 1064 42000 …near 'UNSIGNED UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UNTIL UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UPDATE UNIQUE (id))
< error
| 1064 42000 …near 'UPDATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UPGRADE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT URI UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT USAGE UNIQUE (id))
< error
| 1064 42000 …near 'USAGE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT USE UNIQUE (id))
< error
| 1064 42000 …near 'USE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT USER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT USING UNIQUE (id))
< error
| 1064 42000 …near 'USING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UTC_DATE UNIQUE (id))
< error
| 1064 42000 …near 'UTC_DATE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UTC_TIME UNIQUE (id))
< error
| 1064 42000 …near 'UTC_TIME UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT UTC_TIMESTAMP UNIQUE (id))
< error
| 1064 42000 …near 'UTC_TIMESTAMP UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VALIDATION UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VALUE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VALUES UNIQUE (id))
< error
| 1064 42000 …near 'VALUES UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VARBINARY UNIQUE (id))
< error
| 1064 42000 …near 'VARBINARY UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VARCHAR UNIQUE (id))
< error
| 1064 42000 …near 'VARCHAR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VARIABLES UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VARYING UNIQUE (id))
< error
| 1064 42000 …near 'VARYING UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VIEW UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VIRTUAL UNIQUE (id))
< error
| 1064 42000 …near 'VIRTUAL UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT VISIBLE UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WAIT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WARNINGS UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WHEN UNIQUE (id))
< error
| 1064 42000 …near 'WHEN UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WHERE UNIQUE (id))
< error
| 1064 42000 …near 'WHERE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WHILE UNIQUE (id))
< error
| 1064 42000 …near 'WHILE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WINDOW UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WITH UNIQUE (id))
< error
| 1064 42000 …near 'WITH UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WITHOUT UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WORK UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WRAPPER UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT WRITE UNIQUE (id))
< error
| 1064 42000 …near 'WRITE UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT X509 UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT XA UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT XID UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT XML UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT XOR UNIQUE (id))
< error
| 1064 42000 …near 'XOR UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT YEAR UNIQUE (id))
< ok 0 0
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT YEAR_MONTH UNIQUE (id))
< error
| 1064 42000 …near 'YEAR_MONTH UNIQUE (id))' at line 1
> exec ^
: string
| CREATE TABLE kw_con_test (id INT, CONSTRAINT ZEROFILL UNIQUE (id))
< error
| 1064 42000 …near 'ZEROFILL UNIQUE (id))' at line 1
> exec ^
: string
| SELECT 1 AS ABSENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS ACCESSIBLE
< error
| 1064 42000 …near 'ACCESSIBLE' at line 1
> exec ^
: string
| SELECT 1 AS ACCOUNT
< ok 0 0
> exec ^
: string
| SELECT 1 AS ACTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS ACTIVE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ADD
< error
| 1064 42000 …near 'ADD' at line 1
> exec ^
: string
| SELECT 1 AS ADDDATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ADMIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS AFTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS AGAINST
< ok 0 0
> exec ^
: string
| SELECT 1 AS AGGREGATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ALGORITHM
< ok 0 0
> exec ^
: string
| SELECT 1 AS ALL
< error
| 1064 42000 …near 'ALL' at line 1
> exec ^
: string
| SELECT 1 AS ALLOW_MISSING_FILES
< ok 0 0
> exec ^
: string
| SELECT 1 AS ALTER
< error
| 1064 42000 …near 'ALTER' at line 1
> exec ^
: string
| SELECT 1 AS ALWAYS
< ok 0 0
> exec ^
: string
| SELECT 1 AS ANALYZE
< error
| 1064 42000 …near 'ANALYZE' at line 1
> exec ^
: string
| SELECT 1 AS AND
< error
| 1064 42000 …near 'AND' at line 1
> exec ^
: string
| SELECT 1 AS ANY
< ok 0 0
> exec ^
: string
| SELECT 1 AS ARRAY
< ok 0 0
> exec ^
: string
| SELECT 1 AS AS
< error
| 1064 42000 …near 'AS' at line 1
> exec ^
: string
| SELECT 1 AS ASC
< error
| 1064 42000 …near 'ASC' at line 1
> exec ^
: string
| SELECT 1 AS ASCII
< ok 0 0
> exec ^
: string
| SELECT 1 AS ASENSITIVE
< error
| 1064 42000 …near 'ASENSITIVE' at line 1
> exec ^
: string
| SELECT 1 AS ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS
< ok 0 0
> exec ^
: string
| SELECT 1 AS AT
< ok 0 0
> exec ^
: string
| SELECT 1 AS ATTRIBUTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS AUTHENTICATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS AUTO
< ok 0 0
> exec ^
: string
| SELECT 1 AS AUTOEXTEND_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS AUTO_INCREMENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS AUTO_REFRESH
< ok 0 0
> exec ^
: string
| SELECT 1 AS AUTO_REFRESH_SOURCE
< ok 0 0
> exec ^
: string
| SELECT 1 AS AVG
< ok 0 0
> exec ^
: string
| SELECT 1 AS AVG_ROW_LENGTH
< ok 0 0
> exec ^
: string
| SELECT 1 AS BACKUP
< ok 0 0
> exec ^
: string
| SELECT 1 AS BEFORE
< error
| 1064 42000 …near 'BEFORE' at line 1
> exec ^
: string
| SELECT 1 AS BEGIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS BERNOULLI
< ok 0 0
> exec ^
: string
| SELECT 1 AS BETWEEN
< error
| 1064 42000 …near 'BETWEEN' at line 1
> exec ^
: string
| SELECT 1 AS BIGINT
< error
| 1064 42000 …near 'BIGINT' at line 1
> exec ^
: string
| SELECT 1 AS BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| SELECT 1 AS BINLOG
< ok 0 0
> exec ^
: string
| SELECT 1 AS BIT
< ok 0 0
> exec ^
: string
| SELECT 1 AS BIT_AND
< ok 0 0
> exec ^
: string
| SELECT 1 AS BIT_OR
< ok 0 0
> exec ^
: string
| SELECT 1 AS BIT_XOR
< ok 0 0
> exec ^
: string
| SELECT 1 AS BLOB
< error
| 1064 42000 …near 'BLOB' at line 1
> exec ^
: string
| SELECT 1 AS BLOCK
< ok 0 0
> exec ^
: string
| SELECT 1 AS BOOL
< ok 0 0
> exec ^
: string
| SELECT 1 AS BOOLEAN
< ok 0 0
> exec ^
: string
| SELECT 1 AS BOTH
< error
| 1064 42000 …near 'BOTH' at line 1
> exec ^
: string
| SELECT 1 AS BTREE
< ok 0 0
> exec ^
: string
| SELECT 1 AS BUCKETS
< ok 0 0
> exec ^
: string
| SELECT 1 AS BULK
< ok 0 0
> exec ^
: string
| SELECT 1 AS BY
< error
| 1064 42000 …near 'BY' at line 1
> exec ^
: string
| SELECT 1 AS BYTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CACHE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CALL
< error
| 1064 42000 …near 'CALL' at line 1
> exec ^
: string
| SELECT 1 AS CASCADE
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec ^
: string
| SELECT 1 AS CASCADED
< ok 0 0
> exec ^
: string
| SELECT 1 AS CASE
< error
| 1064 42000 …near 'CASE' at line 1
> exec ^
: string
| SELECT 1 AS CAST
< ok 0 0
> exec ^
: string
| SELECT 1 AS CATALOG_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS CHAIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS CHALLENGE_RESPONSE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CHANGE
< error
| 1064 42000 …near 'CHANGE' at line 1
> exec ^
: string
| SELECT 1 AS CHANGED
< ok 0 0
> exec ^
: string
| SELECT 1 AS CHANNEL
< ok 0 0
> exec ^
: string
| SELECT 1 AS CHAR
< error
| 1064 42000 …near 'CHAR' at line 1
> exec ^
: string
| SELECT 1 AS CHARACTER
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec ^
: string
| SELECT 1 AS CHARSET
< ok 0 0
> exec ^
: string
| SELECT 1 AS CHECK
< error
| 1064 42000 …near 'CHECK' at line 1
> exec ^
: string
| SELECT 1 AS CHECKSUM
< ok 0 0
> exec ^
: string
| SELECT 1 AS CIPHER
< ok 0 0
> exec ^
: string
| SELECT 1 AS CLASS_ORIGIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS CLIENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CLONE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CLOSE
< ok 0 0
> exec ^
: string
| SELECT 1 AS COALESCE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CODE
< ok 0 0
> exec ^
: string
| SELECT 1 AS COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec ^
: string
| SELECT 1 AS COLLATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS COLUMN
< error
| 1064 42000 …near 'COLUMN' at line 1
> exec ^
: string
| SELECT 1 AS COLUMNS
< ok 0 0
> exec ^
: string
| SELECT 1 AS COLUMN_FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS COLUMN_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMMENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMMIT
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMMITTED
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMPACT
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMPLETION
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMPONENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMPRESSED
< ok 0 0
> exec ^
: string
| SELECT 1 AS COMPRESSION
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONCAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONCURRENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONDITION
< error
| 1064 42000 …near 'CONDITION' at line 1
> exec ^
: string
| SELECT 1 AS CONNECTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONSISTENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONSTRAINT
< error
| 1064 42000 …near 'CONSTRAINT' at line 1
> exec ^
: string
| SELECT 1 AS CONSTRAINT_CATALOG
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONSTRAINT_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONSTRAINT_SCHEMA
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONTAINS
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONTEXT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CONTINUE
< error
| 1064 42000 …near 'CONTINUE' at line 1
> exec ^
: string
| SELECT 1 AS CONVERT
< error
| 1064 42000 …near 'CONVERT' at line 1
> exec ^
: string
| SELECT 1 AS COUNT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CPU
< ok 0 0
> exec ^
: string
| SELECT 1 AS CREATE
< error
| 1064 42000 …near 'CREATE' at line 1
> exec ^
: string
| SELECT 1 AS CROSS
< error
| 1064 42000 …near 'CROSS' at line 1
> exec ^
: string
| SELECT 1 AS CSV
< ok 0 0
> exec ^
: string
| SELECT 1 AS CUBE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CUME_DIST
< ok 0 0
> exec ^
: string
| SELECT 1 AS CURDATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS CURRENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS CURRENT_DATE
< error
| 1064 42000 …near 'CURRENT_DATE' at line 1
> exec ^
: string
| SELECT 1 AS CURRENT_TIME
< error
| 1064 42000 …near 'CURRENT_TIME' at line 1
> exec ^
: string
| SELECT 1 AS CURRENT_TIMESTAMP
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP' at line 1
> exec ^
: string
| SELECT 1 AS CURRENT_USER
< error
| 1064 42000 …near 'CURRENT_USER' at line 1
> exec ^
: string
| SELECT 1 AS CURSOR
< error
| 1064 42000 …near 'CURSOR' at line 1
> exec ^
: string
| SELECT 1 AS CURSOR_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS CURTIME
< ok 0 0
> exec ^
: string
| SELECT 1 AS DATA
< ok 0 0
> exec ^
: string
| SELECT 1 AS DATABASE
< error
| 1064 42000 …near 'DATABASE' at line 1
> exec ^
: string
| SELECT 1 AS DATABASES
< error
| 1064 42000 …near 'DATABASES' at line 1
> exec ^
: string
| SELECT 1 AS DATAFILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DATETIME
< ok 0 0
> exec ^
: string
| SELECT 1 AS DATE_ADD
< ok 0 0
> exec ^
: string
| SELECT 1 AS DATE_SUB
< ok 0 0
> exec ^
: string
| SELECT 1 AS DAY
< ok 0 0
> exec ^
: string
| SELECT 1 AS DAY_HOUR
< error
| 1064 42000 …near 'DAY_HOUR' at line 1
> exec ^
: string
| SELECT 1 AS DAY_MICROSECOND
< error
| 1064 42000 …near 'DAY_MICROSECOND' at line 1
> exec ^
: string
| SELECT 1 AS DAY_MINUTE
< error
| 1064 42000 …near 'DAY_MINUTE' at line 1
> exec ^
: string
| SELECT 1 AS DAY_SECOND
< error
| 1064 42000 …near 'DAY_SECOND' at line 1
> exec ^
: string
| SELECT 1 AS DEALLOCATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DEC
< ok 0 0
> exec ^
: string
| SELECT 1 AS DECIMAL
< error
| 1064 42000 …near 'DECIMAL' at line 1
> exec ^
: string
| SELECT 1 AS DECLARE
< error
| 1064 42000 …near 'DECLARE' at line 1
> exec ^
: string
| SELECT 1 AS DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| SELECT 1 AS DEFAULT_AUTH
< ok 0 0
> exec ^
: string
| SELECT 1 AS DEFINER
< ok 0 0
> exec ^
: string
| SELECT 1 AS DEFINITION
< ok 0 0
> exec ^
: string
| SELECT 1 AS DELAYED
< error
| 1064 42000 …near 'DELAYED' at line 1
> exec ^
: string
| SELECT 1 AS DELAY_KEY_WRITE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DELETE
< error
| 1064 42000 …near 'DELETE' at line 1
> exec ^
: string
| SELECT 1 AS DENSE_RANK
< ok 0 0
> exec ^
: string
| SELECT 1 AS DESC
< error
| 1064 42000 …near 'DESC' at line 1
> exec ^
: string
| SELECT 1 AS DESCRIBE
< error
| 1064 42000 …near 'DESCRIBE' at line 1
> exec ^
: string
| SELECT 1 AS DESCRIPTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS DETERMINISTIC
< error
| 1064 42000 …near 'DETERMINISTIC' at line 1
> exec ^
: string
| SELECT 1 AS DIAGNOSTICS
< ok 0 0
> exec ^
: string
| SELECT 1 AS DIRECTORY
< ok 0 0
> exec ^
: string
| SELECT 1 AS DISABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DISCARD
< ok 0 0
> exec ^
: string
| SELECT 1 AS DISK
< ok 0 0
> exec ^
: string
| SELECT 1 AS DISTINCT
< error
| 1064 42000 …near 'DISTINCT' at line 1
> exec ^
: string
| SELECT 1 AS DISTINCTROW
< error
| 1064 42000 …near 'DISTINCTROW' at line 1
> exec ^
: string
| SELECT 1 AS DIV
< error
| 1064 42000 …near 'DIV' at line 1
> exec ^
: string
| SELECT 1 AS DO
< ok 0 0
> exec ^
: string
| SELECT 1 AS DOUBLE
< error
| 1064 42000 …near 'DOUBLE' at line 1
> exec ^
: string
| SELECT 1 AS DROP
< error
| 1064 42000 …near 'DROP' at line 1
> exec ^
: string
| SELECT 1 AS DUAL
< error
| 1064 42000 …near 'DUAL' at line 1
> exec ^
: string
| SELECT 1 AS DUALITY
< ok 0 0
> exec ^
: string
| SELECT 1 AS DUMPFILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DUPLICATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS DYNAMIC
< ok 0 0
> exec ^
: string
| SELECT 1 AS EACH
< error
| 1064 42000 …near 'EACH' at line 1
> exec ^
: string
| SELECT 1 AS ELSE
< error
| 1064 42000 …near 'ELSE' at line 1
> exec ^
: string
| SELECT 1 AS ELSEIF
< error
| 1064 42000 …near 'ELSEIF' at line 1
> exec ^
: string
| SELECT 1 AS EMPTY
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENCLOSED
< error
| 1064 42000 …near 'ENCLOSED' at line 1
> exec ^
: string
| SELECT 1 AS ENCRYPTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS END
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENFORCED
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENGINE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENGINES
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENGINE_ATTRIBUTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ENUM
< ok 0 0
> exec ^
: string
| SELECT 1 AS ERROR
< ok 0 0
> exec ^
: string
| SELECT 1 AS ERRORS
< ok 0 0
> exec ^
: string
| SELECT 1 AS ESCAPE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ESCAPED
< error
| 1064 42000 …near 'ESCAPED' at line 1
> exec ^
: string
| SELECT 1 AS EVENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS EVENTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS EVERY
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXCEPT
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXCHANGE
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXCLUDE
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXECUTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXISTS
< error
| 1064 42000 …near 'EXISTS' at line 1
> exec ^
: string
| SELECT 1 AS EXIT
< error
| 1064 42000 …near 'EXIT' at line 1
> exec ^
: string
| SELECT 1 AS EXPANSION
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXPIRE
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXPLAIN
< error
| 1064 42000 …near 'EXPLAIN' at line 1
> exec ^
: string
| SELECT 1 AS EXPORT
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXTENDED
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXTENT_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXTERNAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXTERNAL_FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS EXTRACT
< ok 0 0
> exec ^
: string
| SELECT 1 AS FACTOR
< ok 0 0
> exec ^
: string
| SELECT 1 AS FAILED_LOGIN_ATTEMPTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS FALSE
< error
| 1064 42000 …near 'FALSE' at line 1
> exec ^
: string
| SELECT 1 AS FAST
< ok 0 0
> exec ^
: string
| SELECT 1 AS FAULTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS FETCH
< error
| 1064 42000 …near 'FETCH' at line 1
> exec ^
: string
| SELECT 1 AS FIELDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILES
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILE_BLOCK_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILE_FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILE_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILE_PATTERN
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILE_PREFIX
< ok 0 0
> exec ^
: string
| SELECT 1 AS FILTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS FINISH
< ok 0 0
> exec ^
: string
| SELECT 1 AS FIRST
< ok 0 0
> exec ^
: string
| SELECT 1 AS FIRST_VALUE
< ok 0 0
> exec ^
: string
| SELECT 1 AS FIXED
< ok 0 0
> exec ^
: string
| SELECT 1 AS FLOAT
< error
| 1064 42000 …near 'FLOAT' at line 1
> exec ^
: string
| SELECT 1 AS FLOAT4
< error
| 1064 42000 …near 'FLOAT4' at line 1
> exec ^
: string
| SELECT 1 AS FLOAT8
< error
| 1064 42000 …near 'FLOAT8' at line 1
> exec ^
: string
| SELECT 1 AS FLUSH
< ok 0 0
> exec ^
: string
| SELECT 1 AS FOLLOWING
< ok 0 0
> exec ^
: string
| SELECT 1 AS FOLLOWS
< ok 0 0
> exec ^
: string
| SELECT 1 AS FOR
< error
| 1064 42000 …near 'FOR' at line 1
> exec ^
: string
| SELECT 1 AS FORCE
< error
| 1064 42000 …near 'FORCE' at line 1
> exec ^
: string
| SELECT 1 AS FOREIGN
< error
| 1064 42000 …near 'FOREIGN' at line 1
> exec ^
: string
| SELECT 1 AS FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS FOUND
< ok 0 0
> exec ^
: string
| SELECT 1 AS FROM
< error
| 1064 42000 …near 'FROM' at line 1
> exec ^
: string
| SELECT 1 AS FULL
< ok 0 0
> exec ^
: string
| SELECT 1 AS FULLTEXT
< error
| 1064 42000 …near 'FULLTEXT' at line 1
> exec ^
: string
| SELECT 1 AS FUNCTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS GENERAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS GENERATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS GENERATED
< error
| 1064 42000 …near 'GENERATED' at line 1
> exec ^
: string
| SELECT 1 AS GEOMCOLLECTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS GEOMETRY
< ok 0 0
> exec ^
: string
| SELECT 1 AS GEOMETRYCOLLECTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS GET
< error
| 1064 42000 …near 'GET' at line 1
> exec ^
: string
| SELECT 1 AS GET_FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS GET_SOURCE_PUBLIC_KEY
< ok 0 0
> exec ^
: string
| SELECT 1 AS GLOBAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS GRANT
< error
| 1064 42000 …near 'GRANT' at line 1
> exec ^
: string
| SELECT 1 AS GRANTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS GREATEST
< ok 0 0
> exec ^
: string
| SELECT 1 AS GROUP
< error
| 1064 42000 …near 'GROUP' at line 1
> exec ^
: string
| SELECT 1 AS GROUPING
< ok 0 0
> exec ^
: string
| SELECT 1 AS GROUPS
< ok 0 0
> exec ^
: string
| SELECT 1 AS GROUP_CONCAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS GROUP_REPLICATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS GTIDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS GTID_ONLY
< ok 0 0
> exec ^
: string
| SELECT 1 AS GUIDED
< ok 0 0
> exec ^
: string
| SELECT 1 AS HANDLER
< ok 0 0
> exec ^
: string
| SELECT 1 AS HASH
< ok 0 0
> exec ^
: string
| SELECT 1 AS HAVING
< error
| 1064 42000 …near 'HAVING' at line 1
> exec ^
: string
| SELECT 1 AS HEADER
< ok 0 0
> exec ^
: string
| SELECT 1 AS HELP
< ok 0 0
> exec ^
: string
| SELECT 1 AS HIGH_PRIORITY
< error
| 1064 42000 …near 'HIGH_PRIORITY' at line 1
> exec ^
: string
| SELECT 1 AS HISTOGRAM
< ok 0 0
> exec ^
: string
| SELECT 1 AS HISTORY
< ok 0 0
> exec ^
: string
| SELECT 1 AS HOST
< ok 0 0
> exec ^
: string
| SELECT 1 AS HOSTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS HOUR
< ok 0 0
> exec ^
: string
| SELECT 1 AS HOUR_MICROSECOND
< error
| 1064 42000 …near 'HOUR_MICROSECOND' at line 1
> exec ^
: string
| SELECT 1 AS HOUR_MINUTE
< error
| 1064 42000 …near 'HOUR_MINUTE' at line 1
> exec ^
: string
| SELECT 1 AS HOUR_SECOND
< error
| 1064 42000 …near 'HOUR_SECOND' at line 1
> exec ^
: string
| SELECT 1 AS IDENTIFIED
< ok 0 0
> exec ^
: string
| SELECT 1 AS IF
< error
| 1064 42000 …near 'IF' at line 1
> exec ^
: string
| SELECT 1 AS IGNORE
< error
| 1064 42000 …near 'IGNORE' at line 1
> exec ^
: string
| SELECT 1 AS IGNORE_SERVER_IDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS IMPORT
< ok 0 0
> exec ^
: string
| SELECT 1 AS IN
< error
| 1064 42000 …near 'IN' at line 1
> exec ^
: string
| SELECT 1 AS INACTIVE
< ok 0 0
> exec ^
: string
| SELECT 1 AS INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec ^
: string
| SELECT 1 AS INDEXES
< ok 0 0
> exec ^
: string
| SELECT 1 AS INFILE
< error
| 1064 42000 …near 'INFILE' at line 1
> exec ^
: string
| SELECT 1 AS INITIAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS INITIAL_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS INITIATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS INNER
< error
| 1064 42000 …near 'INNER' at line 1
> exec ^
: string
| SELECT 1 AS INOUT
< error
| 1064 42000 …near 'INOUT' at line 1
> exec ^
: string
| SELECT 1 AS INSENSITIVE
< error
| 1064 42000 …near 'INSENSITIVE' at line 1
> exec ^
: string
| SELECT 1 AS INSERT
< error
| 1064 42000 …near 'INSERT' at line 1
> exec ^
: string
| SELECT 1 AS INSERT_METHOD
< ok 0 0
> exec ^
: string
| SELECT 1 AS INSTALL
< ok 0 0
> exec ^
: string
| SELECT 1 AS INSTANCE
< ok 0 0
> exec ^
: string
| SELECT 1 AS INT
< error
| 1064 42000 …near 'INT' at line 1
> exec ^
: string
| SELECT 1 AS INT1
< error
| 1064 42000 …near 'INT1' at line 1
> exec ^
: string
| SELECT 1 AS INT2
< error
| 1064 42000 …near 'INT2' at line 1
> exec ^
: string
| SELECT 1 AS INT3
< error
| 1064 42000 …near 'INT3' at line 1
> exec ^
: string
| SELECT 1 AS INT4
< error
| 1064 42000 …near 'INT4' at line 1
> exec ^
: string
| SELECT 1 AS INT8
< error
| 1064 42000 …near 'INT8' at line 1
> exec ^
: string
| SELECT 1 AS INTEGER
< error
| 1064 42000 …near 'INTEGER' at line 1
> exec ^
: string
| SELECT 1 AS INTERSECT
< ok 0 0
> exec ^
: string
| SELECT 1 AS INTERVAL
< error
| 1064 42000 …near 'INTERVAL' at line 1
> exec ^
: string
| SELECT 1 AS INTO
< error
| 1064 42000 …near 'INTO' at line 1
> exec ^
: string
| SELECT 1 AS INVISIBLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS INVOKER
< ok 0 0
> exec ^
: string
| SELECT 1 AS IO
< ok 0 0
> exec ^
: string
| SELECT 1 AS IO_AFTER_GTIDS
< error
| 1064 42000 …near 'IO_AFTER_GTIDS' at line 1
> exec ^
: string
| SELECT 1 AS IO_BEFORE_GTIDS
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS' at line 1
> exec ^
: string
| SELECT 1 AS IO_THREAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS IPC
< ok 0 0
> exec ^
: string
| SELECT 1 AS IS
< error
| 1064 42000 …near 'IS' at line 1
> exec ^
: string
| SELECT 1 AS ISOLATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS ISSUER
< ok 0 0
> exec ^
: string
| SELECT 1 AS ITERATE
< error
| 1064 42000 …near 'ITERATE' at line 1
> exec ^
: string
| SELECT 1 AS JOIN
< error
| 1064 42000 …near 'JOIN' at line 1
> exec ^
: string
| SELECT 1 AS JSON
< ok 0 0
> exec ^
: string
| SELECT 1 AS JSON_ARRAYAGG
< ok 0 0
> exec ^
: string
| SELECT 1 AS JSON_DUALITY_OBJECT
< ok 0 0
> exec ^
: string
| SELECT 1 AS JSON_OBJECTAGG
< ok 0 0
> exec ^
: string
| SELECT 1 AS JSON_TABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS JSON_VALUE
< ok 0 0
> exec ^
: string
| SELECT 1 AS KEY
< error
| 1064 42000 …near 'KEY' at line 1
> exec ^
: string
| SELECT 1 AS KEYRING
< ok 0 0
> exec ^
: string
| SELECT 1 AS KEYS
< error
| 1064 42000 …near 'KEYS' at line 1
> exec ^
: string
| SELECT 1 AS KEY_BLOCK_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS KILL
< error
| 1064 42000 …near 'KILL' at line 1
> exec ^
: string
| SELECT 1 AS LAG
< ok 0 0
> exec ^
: string
| SELECT 1 AS LANGUAGE
< ok 0 0
> exec ^
: string
| SELECT 1 AS LAST
< ok 0 0
> exec ^
: string
| SELECT 1 AS LAST_VALUE
< ok 0 0
> exec ^
: string
| SELECT 1 AS LATERAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS LEAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS LEADING
< error
| 1064 42000 …near 'LEADING' at line 1
> exec ^
: string
| SELECT 1 AS LEAST
< ok 0 0
> exec ^
: string
| SELECT 1 AS LEAVE
< error
| 1064 42000 …near 'LEAVE' at line 1
> exec ^
: string
| SELECT 1 AS LEAVES
< ok 0 0
> exec ^
: string
| SELECT 1 AS LEFT
< error
| 1064 42000 …near 'LEFT' at line 1
> exec ^
: string
| SELECT 1 AS LESS
< ok 0 0
> exec ^
: string
| SELECT 1 AS LEVEL
< ok 0 0
> exec ^
: string
| SELECT 1 AS LIBRARY
< ok 0 0
> exec ^
: string
| SELECT 1 AS LIKE
< error
| 1064 42000 …near 'LIKE' at line 1
> exec ^
: string
| SELECT 1 AS LIMIT
< error
| 1064 42000 …near 'LIMIT' at line 1
> exec ^
: string
| SELECT 1 AS LINEAR
< error
| 1064 42000 …near 'LINEAR' at line 1
> exec ^
: string
| SELECT 1 AS LINES
< error
| 1064 42000 …near 'LINES' at line 1
> exec ^
: string
| SELECT 1 AS LINESTRING
< ok 0 0
> exec ^
: string
| SELECT 1 AS LIST
< ok 0 0
> exec ^
: string
| SELECT 1 AS LOAD
< error
| 1064 42000 …near 'LOAD' at line 1
> exec ^
: string
| SELECT 1 AS LOCAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS LOCALTIME
< error
| 1064 42000 …near 'LOCALTIME' at line 1
> exec ^
: string
| SELECT 1 AS LOCALTIMESTAMP
< error
| 1064 42000 …near 'LOCALTIMESTAMP' at line 1
> exec ^
: string
| SELECT 1 AS LOCK
< error
| 1064 42000 …near 'LOCK' at line 1
> exec ^
: string
| SELECT 1 AS LOCKED
< ok 0 0
> exec ^
: string
| SELECT 1 AS LOCKS
< ok 0 0
> exec ^
: string
| SELECT 1 AS LOG
< ok 0 0
> exec ^
: string
| SELECT 1 AS LOGFILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS LOGS
< ok 0 0
> exec ^
: string
| SELECT 1 AS LONG
< error
| 1064 42000 …near 'LONG' at line 1
> exec ^
: string
| SELECT 1 AS LONGBLOB
< error
| 1064 42000 …near 'LONGBLOB' at line 1
> exec ^
: string
| SELECT 1 AS LONGTEXT
< error
| 1064 42000 …near 'LONGTEXT' at line 1
> exec ^
: string
| SELECT 1 AS LOOP
< error
| 1064 42000 …near 'LOOP' at line 1
> exec ^
: string
| SELECT 1 AS LOW_PRIORITY
< error
| 1064 42000 …near 'LOW_PRIORITY' at line 1
> exec ^
: string
| SELECT 1 AS MANIFEST
< ok 0 0
> exec ^
: string
| SELECT 1 AS MANUAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS MASTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS MATCH
< error
| 1064 42000 …near 'MATCH' at line 1
> exec ^
: string
| SELECT 1 AS MATERIALIZED
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAXVALUE
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX_CONNECTIONS_PER_HOUR
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX_QUERIES_PER_HOUR
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX_ROWS
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX_UPDATES_PER_HOUR
< ok 0 0
> exec ^
: string
| SELECT 1 AS MAX_USER_CONNECTIONS
< ok 0 0
> exec ^
: string
| SELECT 1 AS MEDIUM
< ok 0 0
> exec ^
: string
| SELECT 1 AS MEDIUMBLOB
< error
| 1064 42000 …near 'MEDIUMBLOB' at line 1
> exec ^
: string
| SELECT 1 AS MEDIUMINT
< error
| 1064 42000 …near 'MEDIUMINT' at line 1
> exec ^
: string
| SELECT 1 AS MEDIUMTEXT
< error
| 1064 42000 …near 'MEDIUMTEXT' at line 1
> exec ^
: string
| SELECT 1 AS MEMBER
< ok 0 0
> exec ^
: string
| SELECT 1 AS MEMORY
< ok 0 0
> exec ^
: string
| SELECT 1 AS MERGE
< ok 0 0
> exec ^
: string
| SELECT 1 AS MESSAGE_TEXT
< ok 0 0
> exec ^
: string
| SELECT 1 AS MICROSECOND
< ok 0 0
> exec ^
: string
| SELECT 1 AS MID
< ok 0 0
> exec ^
: string
| SELECT 1 AS MIDDLEINT
< error
| 1064 42000 …near 'MIDDLEINT' at line 1
> exec ^
: string
| SELECT 1 AS MIGRATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS MIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS MINUTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS MINUTE_MICROSECOND
< error
| 1064 42000 …near 'MINUTE_MICROSECOND' at line 1
> exec ^
: string
| SELECT 1 AS MINUTE_SECOND
< error
| 1064 42000 …near 'MINUTE_SECOND' at line 1
> exec ^
: string
| SELECT 1 AS MIN_ROWS
< ok 0 0
> exec ^
: string
| SELECT 1 AS MOD
< error
| 1064 42000 …near 'MOD' at line 1
> exec ^
: string
| SELECT 1 AS MODE
< ok 0 0
> exec ^
: string
| SELECT 1 AS MODIFIES
< error
| 1064 42000 …near 'MODIFIES' at line 1
> exec ^
: string
| SELECT 1 AS MODIFY
< ok 0 0
> exec ^
: string
| SELECT 1 AS MONTH
< ok 0 0
> exec ^
: string
| SELECT 1 AS MULTILINESTRING
< ok 0 0
> exec ^
: string
| SELECT 1 AS MULTIPOINT
< ok 0 0
> exec ^
: string
| SELECT 1 AS MULTIPOLYGON
< ok 0 0
> exec ^
: string
| SELECT 1 AS MUTEX
< ok 0 0
> exec ^
: string
| SELECT 1 AS MYSQL_ERRNO
< ok 0 0
> exec ^
: string
| SELECT 1 AS NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS NAMES
< ok 0 0
> exec ^
: string
| SELECT 1 AS NATIONAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS NATURAL
< error
| 1064 42000 …near 'NATURAL' at line 1
> exec ^
: string
| SELECT 1 AS NCHAR
< ok 0 0
> exec ^
: string
| SELECT 1 AS NDB
< ok 0 0
> exec ^
: string
| SELECT 1 AS NDBCLUSTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS NESTED
< ok 0 0
> exec ^
: string
| SELECT 1 AS NETWORK_NAMESPACE
< ok 0 0
> exec ^
: string
| SELECT 1 AS NEVER
< ok 0 0
> exec ^
: string
| SELECT 1 AS NEW
< ok 0 0
> exec ^
: string
| SELECT 1 AS NEXT
< ok 0 0
> exec ^
: string
| SELECT 1 AS NO
< ok 0 0
> exec ^
: string
| SELECT 1 AS NODEGROUP
< ok 0 0
> exec ^
: string
| SELECT 1 AS NONE
< ok 0 0
> exec ^
: string
| SELECT 1 AS NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec ^
: string
| SELECT 1 AS NOW
< ok 0 0
> exec ^
: string
| SELECT 1 AS NOWAIT
< ok 0 0
> exec ^
: string
| SELECT 1 AS NO_WAIT
< ok 0 0
> exec ^
: string
| SELECT 1 AS NO_WRITE_TO_BINLOG
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG' at line 1
> exec ^
: string
| SELECT 1 AS NTH_VALUE
< ok 0 0
> exec ^
: string
| SELECT 1 AS NTILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| SELECT 1 AS NULLIF
< ok 0 0
> exec ^
: string
| SELECT 1 AS NULLS
< ok 0 0
> exec ^
: string
| SELECT 1 AS NUMBER
< ok 0 0
> exec ^
: string
| SELECT 1 AS NUMERIC
< error
| 1064 42000 …near 'NUMERIC' at line 1
> exec ^
: string
| SELECT 1 AS NVARCHAR
< ok 0 0
> exec ^
: string
| SELECT 1 AS OF
< ok 0 0
> exec ^
: string
| SELECT 1 AS OFF
< ok 0 0
> exec ^
: string
| SELECT 1 AS OFFSET
< ok 0 0
> exec ^
: string
| SELECT 1 AS OJ
< ok 0 0
> exec ^
: string
| SELECT 1 AS OLD
< ok 0 0
> exec ^
: string
| SELECT 1 AS ON
< error
| 1064 42000 …near 'ON' at line 1
> exec ^
: string
| SELECT 1 AS ONE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ONLY
< ok 0 0
> exec ^
: string
| SELECT 1 AS OPEN
< ok 0 0
> exec ^
: string
| SELECT 1 AS OPTIMIZE
< error
| 1064 42000 …near 'OPTIMIZE' at line 1
> exec ^
: string
| SELECT 1 AS OPTIMIZER_COSTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS OPTION
< error
| 1064 42000 …near 'OPTION' at line 1
> exec ^
: string
| SELECT 1 AS OPTIONAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS OPTIONALLY
< error
| 1064 42000 …near 'OPTIONALLY' at line 1
> exec ^
: string
| SELECT 1 AS OPTIONS
< ok 0 0
> exec ^
: string
| SELECT 1 AS OR
< error
| 1064 42000 …near 'OR' at line 1
> exec ^
: string
| SELECT 1 AS ORDER
< error
| 1064 42000 …near 'ORDER' at line 1
> exec ^
: string
| SELECT 1 AS ORDINALITY
< ok 0 0
> exec ^
: string
| SELECT 1 AS ORGANIZATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS OTHERS
< ok 0 0
> exec ^
: string
| SELECT 1 AS OUT
< error
| 1064 42000 …near 'OUT' at line 1
> exec ^
: string
| SELECT 1 AS OUTER
< error
| 1064 42000 …near 'OUTER' at line 1
> exec ^
: string
| SELECT 1 AS OUTFILE
< error
| 1064 42000 …near 'OUTFILE' at line 1
> exec ^
: string
| SELECT 1 AS OVER
< ok 0 0
> exec ^
: string
| SELECT 1 AS OVERLAY
< ok 0 0
> exec ^
: string
| SELECT 1 AS OVERWRITE
< ok 0 0
> exec ^
: string
| SELECT 1 AS OWNER
< ok 0 0
> exec ^
: string
| SELECT 1 AS PACK_KEYS
< ok 0 0
> exec ^
: string
| SELECT 1 AS PAGE
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARALLEL
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARAMETERS
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARSER
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARSE_TREE
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARTIAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARTITION
< error
| 1064 42000 …near 'PARTITION' at line 1
> exec ^
: string
| SELECT 1 AS PARTITIONING
< ok 0 0
> exec ^
: string
| SELECT 1 AS PARTITIONS
< ok 0 0
> exec ^
: string
| SELECT 1 AS PASSWORD
< ok 0 0
> exec ^
: string
| SELECT 1 AS PASSWORD_LOCK_TIME
< ok 0 0
> exec ^
: string
| SELECT 1 AS PATH
< ok 0 0
> exec ^
: string
| SELECT 1 AS PERCENT_RANK
< ok 0 0
> exec ^
: string
| SELECT 1 AS PERSIST
< ok 0 0
> exec ^
: string
| SELECT 1 AS PERSIST_ONLY
< ok 0 0
> exec ^
: string
| SELECT 1 AS PHASE
< ok 0 0
> exec ^
: string
| SELECT 1 AS PLACING
< ok 0 0
> exec ^
: string
| SELECT 1 AS PLUGIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS PLUGINS
< ok 0 0
> exec ^
: string
| SELECT 1 AS PLUGIN_DIR
< ok 0 0
> exec ^
: string
| SELECT 1 AS POINT
< ok 0 0
> exec ^
: string
| SELECT 1 AS POLYGON
< ok 0 0
> exec ^
: string
| SELECT 1 AS PORT
< ok 0 0
> exec ^
: string
| SELECT 1 AS POSITION
< ok 0 0
> exec ^
: string
| SELECT 1 AS PRECEDES
< ok 0 0
> exec ^
: string
| SELECT 1 AS PRECEDING
< ok 0 0
> exec ^
: string
| SELECT 1 AS PRECISION
< error
| 1064 42000 …near 'PRECISION' at line 1
> exec ^
: string
| SELECT 1 AS PREFIX
< ok 0 0
> exec ^
: string
| SELECT 1 AS PREPARE
< ok 0 0
> exec ^
: string
| SELECT 1 AS PRESERVE
< ok 0 0
> exec ^
: string
| SELECT 1 AS PREV
< ok 0 0
> exec ^
: string
| SELECT 1 AS PRIMARY
< error
| 1064 42000 …near 'PRIMARY' at line 1
> exec ^
: string
| SELECT 1 AS PRIVILEGES
< ok 0 0
> exec ^
: string
| SELECT 1 AS PRIVILEGE_CHECKS_USER
< ok 0 0
> exec ^
: string
| SELECT 1 AS PROCEDURE
< error
| 1064 42000 …near 'PROCEDURE' at line 1
> exec ^
: string
| SELECT 1 AS PROCESS
< ok 0 0
> exec ^
: string
| SELECT 1 AS PROCESSLIST
< ok 0 0
> exec ^
: string
| SELECT 1 AS PROFILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS PROFILES
< ok 0 0
> exec ^
: string
| SELECT 1 AS PROXY
< ok 0 0
> exec ^
: string
| SELECT 1 AS PURGE
< error
| 1064 42000 …near 'PURGE' at line 1
> exec ^
: string
| SELECT 1 AS QUALIFY
< ok 0 0
> exec ^
: string
| SELECT 1 AS QUARTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS QUERY
< ok 0 0
> exec ^
: string
| SELECT 1 AS QUICK
< ok 0 0
> exec ^
: string
| SELECT 1 AS RANDOM
< ok 0 0
> exec ^
: string
| SELECT 1 AS RANGE
< ok 0 0
> exec ^
: string
| SELECT 1 AS RANK
< ok 0 0
> exec ^
: string
| SELECT 1 AS READ
< error
| 1064 42000 …near 'READ' at line 1
> exec ^
: string
| SELECT 1 AS READS
< error
| 1064 42000 …near 'READS' at line 1
> exec ^
: string
| SELECT 1 AS READ_ONLY
< ok 0 0
> exec ^
: string
| SELECT 1 AS READ_WRITE
< error
| 1064 42000 …near 'READ_WRITE' at line 1
> exec ^
: string
| SELECT 1 AS REAL
< error
| 1064 42000 …near 'REAL' at line 1
> exec ^
: string
| SELECT 1 AS REBUILD
< ok 0 0
> exec ^
: string
| SELECT 1 AS RECOVER
< ok 0 0
> exec ^
: string
| SELECT 1 AS RECURSIVE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REDO_BUFFER_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REDUNDANT
< ok 0 0
> exec ^
: string
| SELECT 1 AS REFERENCE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REFERENCES
< error
| 1064 42000 …near 'REFERENCES' at line 1
> exec ^
: string
| SELECT 1 AS REGEXP
< error
| 1064 42000 …near 'REGEXP' at line 1
> exec ^
: string
| SELECT 1 AS REGISTRATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELATIONAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELAY
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELAYLOG
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELAY_LOG_FILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELAY_LOG_POS
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELAY_THREAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS RELEASE
< error
| 1064 42000 …near 'RELEASE' at line 1
> exec ^
: string
| SELECT 1 AS RELOAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS REMOVE
< ok 0 0
> exec ^
: string
| SELECT 1 AS RENAME
< error
| 1064 42000 …near 'RENAME' at line 1
> exec ^
: string
| SELECT 1 AS REORGANIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPAIR
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPEAT
< error
| 1064 42000 …near 'REPEAT' at line 1
> exec ^
: string
| SELECT 1 AS REPEATABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLACE
< error
| 1064 42000 …near 'REPLACE' at line 1
> exec ^
: string
| SELECT 1 AS REPLICA
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICAS
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_DO_DB
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_DO_TABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_IGNORE_DB
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_IGNORE_TABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_REWRITE_DB
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_WILD_DO_TABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATE_WILD_IGNORE_TABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REPLICATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS REQUIRE
< error
| 1064 42000 …near 'REQUIRE' at line 1
> exec ^
: string
| SELECT 1 AS REQUIRE_ROW_FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS REQUIRE_TABLE_PRIMARY_KEY_CHECK
< ok 0 0
> exec ^
: string
| SELECT 1 AS RESET
< ok 0 0
> exec ^
: string
| SELECT 1 AS RESIGNAL
< error
| 1064 42000 …near 'RESIGNAL' at line 1
> exec ^
: string
| SELECT 1 AS RESOURCE
< ok 0 0
> exec ^
: string
| SELECT 1 AS RESPECT
< ok 0 0
> exec ^
: string
| SELECT 1 AS RESTART
< ok 0 0
> exec ^
: string
| SELECT 1 AS RESTORE
< ok 0 0
> exec ^
: string
| SELECT 1 AS RESTRICT
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec ^
: string
| SELECT 1 AS RESUME
< ok 0 0
> exec ^
: string
| SELECT 1 AS RETAIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS RETURN
< error
| 1064 42000 …near 'RETURN' at line 1
> exec ^
: string
| SELECT 1 AS RETURNED_SQLSTATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS RETURNING
< ok 0 0
> exec ^
: string
| SELECT 1 AS RETURNS
< ok 0 0
> exec ^
: string
| SELECT 1 AS REUSE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REVERSE
< ok 0 0
> exec ^
: string
| SELECT 1 AS REVOKE
< error
| 1064 42000 …near 'REVOKE' at line 1
> exec ^
: string
| SELECT 1 AS RIGHT
< error
| 1064 42000 …near 'RIGHT' at line 1
> exec ^
: string
| SELECT 1 AS RLIKE
< error
| 1064 42000 …near 'RLIKE' at line 1
> exec ^
: string
| SELECT 1 AS ROLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROLLBACK
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROLLUP
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROTATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROUTINE
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROW
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROWS
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROW_COUNT
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROW_FORMAT
< ok 0 0
> exec ^
: string
| SELECT 1 AS ROW_NUMBER
< ok 0 0
> exec ^
: string
| SELECT 1 AS RTREE
< ok 0 0
> exec ^
: string
| SELECT 1 AS S3
< ok 0 0
> exec ^
: string
| SELECT 1 AS SAVEPOINT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SCHEDULE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SCHEMA
< error
| 1064 42000 …near 'SCHEMA' at line 1
> exec ^
: string
| SELECT 1 AS SCHEMAS
< error
| 1064 42000 …near 'SCHEMAS' at line 1
> exec ^
: string
| SELECT 1 AS SCHEMA_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECOND
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECONDARY
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECONDARY_ENGINE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECONDARY_ENGINE_ATTRIBUTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECONDARY_LOAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECONDARY_UNLOAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS SECOND_MICROSECOND
< error
| 1064 42000 …near 'SECOND_MICROSECOND' at line 1
> exec ^
: string
| SELECT 1 AS SECURITY
< ok 0 0
> exec ^
: string
| SELECT 1 AS SELECT
< error
| 1064 42000 …near 'SELECT' at line 1
> exec ^
: string
| SELECT 1 AS SENSITIVE
< error
| 1064 42000 …near 'SENSITIVE' at line 1
> exec ^
: string
| SELECT 1 AS SEPARATOR
< error
| 1064 42000 …near 'SEPARATOR' at line 1
> exec ^
: string
| SELECT 1 AS SERIAL
< ok 0 0
> exec ^
: string
| SELECT 1 AS SERIALIZABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SERVER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SESSION
< ok 0 0
> exec ^
: string
| SELECT 1 AS SESSION_USER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SET
< error
| 1064 42000 …near 'SET' at line 1
> exec ^
: string
| SELECT 1 AS SETS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SHARE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SHOW
< error
| 1064 42000 …near 'SHOW' at line 1
> exec ^
: string
| SELECT 1 AS SHUTDOWN
< ok 0 0
> exec ^
: string
| SELECT 1 AS SIGNAL
< error
| 1064 42000 …near 'SIGNAL' at line 1
> exec ^
: string
| SELECT 1 AS SIGNED
< ok 0 0
> exec ^
: string
| SELECT 1 AS SIMPLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SKIP
< error
| 1064 42000 …near 'SKIP' at line 1
> exec ^
: string
| SELECT 1 AS SLAVE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SLOW
< ok 0 0
> exec ^
: string
| SELECT 1 AS SMALLINT
< error
| 1064 42000 …near 'SMALLINT' at line 1
> exec ^
: string
| SELECT 1 AS SNAPSHOT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOCKET
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOME
< ok 0 0
> exec ^
: string
| SELECT 1 AS SONAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOUNDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_AUTO_POSITION
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_BIND
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_COMPRESSION_ALGORITHMS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_CONNECTION_AUTO_FAILOVER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_CONNECT_RETRY
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_DELAY
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_HEARTBEAT_PERIOD
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_HOST
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_LOG_FILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_LOG_POS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_PASSWORD
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_PORT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_PUBLIC_KEY_PATH
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_RETRY_COUNT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_CA
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_CAPATH
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_CERT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_CIPHER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_CRL
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_CRLPATH
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_KEY
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_SSL_VERIFY_SERVER_CERT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_TLS_CIPHERSUITES
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_TLS_VERSION
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_USER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SOURCE_ZSTD_COMPRESSION_LEVEL
< ok 0 0
> exec ^
: string
| SELECT 1 AS SPATIAL
< error
| 1064 42000 …near 'SPATIAL' at line 1
> exec ^
: string
| SELECT 1 AS SPECIFIC
< error
| 1064 42000 …near 'SPECIFIC' at line 1
> exec ^
: string
| SELECT 1 AS SQL
< error
| 1064 42000 …near 'SQL' at line 1
> exec ^
: string
| SELECT 1 AS SQLEXCEPTION
< error
| 1064 42000 …near 'SQLEXCEPTION' at line 1
> exec ^
: string
| SELECT 1 AS SQLSTATE
< error
| 1064 42000 …near 'SQLSTATE' at line 1
> exec ^
: string
| SELECT 1 AS SQLWARNING
< error
| 1064 42000 …near 'SQLWARNING' at line 1
> exec ^
: string
| SELECT 1 AS SQL_AFTER_GTIDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_AFTER_MTS_GAPS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_BEFORE_GTIDS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_BIG_RESULT
< error
| 1064 42000 …near 'SQL_BIG_RESULT' at line 1
> exec ^
: string
| SELECT 1 AS SQL_BUFFER_RESULT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_CALC_FOUND_ROWS
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS' at line 1
> exec ^
: string
| SELECT 1 AS SQL_NO_CACHE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_SMALL_RESULT
< error
| 1064 42000 …near 'SQL_SMALL_RESULT' at line 1
> exec ^
: string
| SELECT 1 AS SQL_THREAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_DAY
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_HOUR
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_MINUTE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_MONTH
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_QUARTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_SECOND
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_WEEK
< ok 0 0
> exec ^
: string
| SELECT 1 AS SQL_TSI_YEAR
< ok 0 0
> exec ^
: string
| SELECT 1 AS SRID
< ok 0 0
> exec ^
: string
| SELECT 1 AS SSE_KMS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SSE_S3
< ok 0 0
> exec ^
: string
| SELECT 1 AS SSL
< error
| 1064 42000 …near 'SSL' at line 1
> exec ^
: string
| SELECT 1 AS STACKED
< ok 0 0
> exec ^
: string
| SELECT 1 AS START
< ok 0 0
> exec ^
: string
| SELECT 1 AS STARTING
< error
| 1064 42000 …near 'STARTING' at line 1
> exec ^
: string
| SELECT 1 AS STARTS
< ok 0 0
> exec ^
: string
| SELECT 1 AS STATS_AUTO_RECALC
< ok 0 0
> exec ^
: string
| SELECT 1 AS STATS_PERSISTENT
< ok 0 0
> exec ^
: string
| SELECT 1 AS STATS_SAMPLE_PAGES
< ok 0 0
> exec ^
: string
| SELECT 1 AS STATUS
< ok 0 0
> exec ^
: string
| SELECT 1 AS STD
< ok 0 0
> exec ^
: string
| SELECT 1 AS STDDEV
< ok 0 0
> exec ^
: string
| SELECT 1 AS STDDEV_POP
< ok 0 0
> exec ^
: string
| SELECT 1 AS STDDEV_SAMP
< ok 0 0
> exec ^
: string
| SELECT 1 AS STOP
< ok 0 0
> exec ^
: string
| SELECT 1 AS STORAGE
< ok 0 0
> exec ^
: string
| SELECT 1 AS STORED
< error
| 1064 42000 …near 'STORED' at line 1
> exec ^
: string
| SELECT 1 AS STRAIGHT_JOIN
< error
| 1064 42000 …near 'STRAIGHT_JOIN' at line 1
> exec ^
: string
| SELECT 1 AS STREAM
< ok 0 0
> exec ^
: string
| SELECT 1 AS STRICT_LOAD
< ok 0 0
> exec ^
: string
| SELECT 1 AS STRING
< ok 0 0
> exec ^
: string
| SELECT 1 AS ST_COLLECT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBCLASS_ORIGIN
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBDATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBJECT
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBPARTITION
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBPARTITIONS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBSTR
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUBSTRING
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUM
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUPER
< ok 0 0
> exec ^
: string
| SELECT 1 AS SUSPEND
< ok 0 0
> exec ^
: string
| SELECT 1 AS SWAPS
< ok 0 0
> exec ^
: string
| SELECT 1 AS SWITCHES
< ok 0 0
> exec ^
: string
| SELECT 1 AS SYSDATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS SYSTEM
< ok 0 0
> exec ^
: string
| SELECT 1 AS SYSTEM_USER
< ok 0 0
> exec ^
: string
| SELECT 1 AS TABLE
< error
| 1064 42000 …near 'TABLE' at line 1
> exec ^
: string
| SELECT 1 AS TABLES
< ok 0 0
> exec ^
: string
| SELECT 1 AS TABLESAMPLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS TABLESPACE
< ok 0 0
> exec ^
: string
| SELECT 1 AS TABLE_CHECKSUM
< ok 0 0
> exec ^
: string
| SELECT 1 AS TABLE_NAME
< ok 0 0
> exec ^
: string
| SELECT 1 AS TEMPORARY
< ok 0 0
> exec ^
: string
| SELECT 1 AS TEMPTABLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS TERMINATED
< error
| 1064 42000 …near 'TERMINATED' at line 1
> exec ^
: string
| SELECT 1 AS TEXT
< ok 0 0
> exec ^
: string
| SELECT 1 AS THAN
< ok 0 0
> exec ^
: string
| SELECT 1 AS THEN
< error
| 1064 42000 …near 'THEN' at line 1
> exec ^
: string
| SELECT 1 AS THREAD_PRIORITY
< ok 0 0
> exec ^
: string
| SELECT 1 AS TIES
< ok 0 0
> exec ^
: string
| SELECT 1 AS TIME
< ok 0 0
> exec ^
: string
| SELECT 1 AS TIMESTAMP
< ok 0 0
> exec ^
: string
| SELECT 1 AS TIMESTAMPADD
< ok 0 0
> exec ^
: string
| SELECT 1 AS TIMESTAMPDIFF
< ok 0 0
> exec ^
: string
| SELECT 1 AS TINYBLOB
< error
| 1064 42000 …near 'TINYBLOB' at line 1
> exec ^
: string
| SELECT 1 AS TINYINT
< error
| 1064 42000 …near 'TINYINT' at line 1
> exec ^
: string
| SELECT 1 AS TINYTEXT
< error
| 1064 42000 …near 'TINYTEXT' at line 1
> exec ^
: string
| SELECT 1 AS TLS
< ok 0 0
> exec ^
: string
| SELECT 1 AS TO
< error
| 1064 42000 …near 'TO' at line 1
> exec ^
: string
| SELECT 1 AS TRAILING
< error
| 1064 42000 …near 'TRAILING' at line 1
> exec ^
: string
| SELECT 1 AS TRANSACTION
< ok 0 0
> exec ^
: string
| SELECT 1 AS TRIGGER
< error
| 1064 42000 …near 'TRIGGER' at line 1
> exec ^
: string
| SELECT 1 AS TRIGGERS
< ok 0 0
> exec ^
: string
| SELECT 1 AS TRIM
< ok 0 0
> exec ^
: string
| SELECT 1 AS TRUE
< error
| 1064 42000 …near 'TRUE' at line 1
> exec ^
: string
| SELECT 1 AS TRUNCATE
< ok 0 0
> exec ^
: string
| SELECT 1 AS TYPE
< ok 0 0
> exec ^
: string
| SELECT 1 AS TYPES
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNBOUNDED
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNCOMMITTED
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNDEFINED
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNDO
< error
| 1064 42000 …near 'UNDO' at line 1
> exec ^
: string
| SELECT 1 AS UNDOFILE
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNDO_BUFFER_SIZE
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNICODE
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNINSTALL
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNION
< error
| 1064 42000 …near 'UNION' at line 1
> exec ^
: string
| SELECT 1 AS UNIQUE
< error
| 1064 42000 …near 'UNIQUE' at line 1
> exec ^
: string
| SELECT 1 AS UNKNOWN
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNLOCK
< error
| 1064 42000 …near 'UNLOCK' at line 1
> exec ^
: string
| SELECT 1 AS UNREGISTER
< ok 0 0
> exec ^
: string
| SELECT 1 AS UNSIGNED
< error
| 1064 42000 …near 'UNSIGNED' at line 1
> exec ^
: string
| SELECT 1 AS UNTIL
< ok 0 0
> exec ^
: string
| SELECT 1 AS UPDATE
< error
| 1064 42000 …near 'UPDATE' at line 1
> exec ^
: string
| SELECT 1 AS UPGRADE
< ok 0 0
> exec ^
: string
| SELECT 1 AS URI
< ok 0 0
> exec ^
: string
| SELECT 1 AS USAGE
< error
| 1064 42000 …near 'USAGE' at line 1
> exec ^
: string
| SELECT 1 AS USE
< error
| 1064 42000 …near 'USE' at line 1
> exec ^
: string
| SELECT 1 AS USER
< ok 0 0
> exec ^
: string
| SELECT 1 AS USING
< error
| 1064 42000 …near 'USING' at line 1
> exec ^
: string
| SELECT 1 AS UTC_DATE
< error
| 1064 42000 …near 'UTC_DATE' at line 1
> exec ^
: string
| SELECT 1 AS UTC_TIME
< error
| 1064 42000 …near 'UTC_TIME' at line 1
> exec ^
: string
| SELECT 1 AS UTC_TIMESTAMP
< error
| 1064 42000 …near 'UTC_TIMESTAMP' at line 1
> exec ^
: string
| SELECT 1 AS VALIDATION
< ok 0 0
> exec ^
: string
| SELECT 1 AS VALUE
< ok 0 0
> exec ^
: string
| SELECT 1 AS VALUES
< error
| 1064 42000 …near 'VALUES' at line 1
> exec ^
: string
| SELECT 1 AS VARBINARY
< error
| 1064 42000 …near 'VARBINARY' at line 1
> exec ^
: string
| SELECT 1 AS VARCHAR
< error
| 1064 42000 …near 'VARCHAR' at line 1
> exec ^
: string
| SELECT 1 AS VARIABLES
< ok 0 0
> exec ^
: string
| SELECT 1 AS VARYING
< error
| 1064 42000 …near 'VARYING' at line 1
> exec ^
: string
| SELECT 1 AS VIEW
< ok 0 0
> exec ^
: string
| SELECT 1 AS VIRTUAL
< error
| 1064 42000 …near 'VIRTUAL' at line 1
> exec ^
: string
| SELECT 1 AS VISIBLE
< ok 0 0
> exec ^
: string
| SELECT 1 AS WAIT
< ok 0 0
> exec ^
: string
| SELECT 1 AS WARNINGS
< ok 0 0
> exec ^
: string
| SELECT 1 AS WHEN
< error
| 1064 42000 …near 'WHEN' at line 1
> exec ^
: string
| SELECT 1 AS WHERE
< error
| 1064 42000 …near 'WHERE' at line 1
> exec ^
: string
| SELECT 1 AS WHILE
< error
| 1064 42000 …near 'WHILE' at line 1
> exec ^
: string
| SELECT 1 AS WINDOW
< ok 0 0
> exec ^
: string
| SELECT 1 AS WITH
< error
| 1064 42000 …near 'WITH' at line 1
> exec ^
: string
| SELECT 1 AS WITHOUT
< ok 0 0
> exec ^
: string
| SELECT 1 AS WORK
< ok 0 0
> exec ^
: string
| SELECT 1 AS WRAPPER
< ok 0 0
> exec ^
: string
| SELECT 1 AS WRITE
< error
| 1064 42000 …near 'WRITE' at line 1
> exec ^
: string
| SELECT 1 AS X509
< ok 0 0
> exec ^
: string
| SELECT 1 AS XA
< ok 0 0
> exec ^
: string
| SELECT 1 AS XID
< ok 0 0
> exec ^
: string
| SELECT 1 AS XML
< ok 0 0
> exec ^
: string
| SELECT 1 AS XOR
< error
| 1064 42000 …near 'XOR' at line 1
> exec ^
: string
| SELECT 1 AS YEAR
< ok 0 0
> exec ^
: string
| SELECT 1 AS YEAR_MONTH
< error
| 1064 42000 …near 'YEAR_MONTH' at line 1
> exec ^
: string
| SELECT 1 AS ZEROFILL
< error
| 1064 42000 …near 'ZEROFILL' at line 1
> exec
| CREATE TABLE kw_alias_base (id INT)
< ok 0 0
> exec
| PREPARE omni_kw_stmt FROM ?
: string
| SELECT * FROM kw_alias_base ABSENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ACCESSIBLE
< error
| 1064 42000 …near 'ACCESSIBLE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ACCOUNT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ACTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ACTIVE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ADD
< error
| 1064 42000 …near 'ADD' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ADDDATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ADMIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AFTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AGAINST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AGGREGATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ALGORITHM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ALL
< error
| 1064 42000 …near 'ALL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ALLOW_MISSING_FILES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ALTER
< error
| 1064 42000 …near 'ALTER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ALWAYS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ANALYZE
< error
| 1064 42000 …near 'ANALYZE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base AND
< error
| 1064 42000 …near 'AND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ANY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ARRAY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AS
< error
| 1064 42000 …near 'AS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ASC
< error
| 1064 42000 …near 'ASC' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ASCII
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ASENSITIVE
< error
| 1064 42000 …near 'ASENSITIVE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ASSIGN_GTIDS_TO_ANONYMOUS_TRANSACTIONS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ATTRIBUTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AUTHENTICATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AUTO
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AUTOEXTEND_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AUTO_INCREMENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AUTO_REFRESH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AUTO_REFRESH_SOURCE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AVG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base AVG_ROW_LENGTH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BACKUP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BEFORE
< error
| 1064 42000 …near 'BEFORE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BEGIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BERNOULLI
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BETWEEN
< error
| 1064 42000 …near 'BETWEEN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BIGINT
< error
| 1064 42000 …near 'BIGINT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BINLOG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BIT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BIT_AND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BIT_OR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BIT_XOR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BLOB
< error
| 1064 42000 …near 'BLOB' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BLOCK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BOOL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BOOLEAN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BOTH
< error
| 1064 42000 …near 'BOTH' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BTREE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BUCKETS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BULK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base BY
< error
| 1064 42000 …near 'BY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base BYTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CACHE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CALL
< error
| 1064 42000 …near 'CALL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CASCADE
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CASCADED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CASE
< error
| 1064 42000 …near 'CASE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CAST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CATALOG_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CHAIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CHALLENGE_RESPONSE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CHANGE
< error
| 1064 42000 …near 'CHANGE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CHANGED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CHANNEL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CHAR
< error
| 1064 42000 …near 'CHAR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CHARACTER
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CHARSET
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CHECK
< error
| 1064 42000 …near 'CHECK' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CHECKSUM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CIPHER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CLASS_ORIGIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CLIENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CLONE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CLOSE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COALESCE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CODE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base COLLATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COLUMN
< error
| 1064 42000 …near 'COLUMN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base COLUMNS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COLUMN_FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COLUMN_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMMENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMMIT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMMITTED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMPACT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMPLETION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMPONENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMPRESSED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base COMPRESSION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONCAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONCURRENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONDITION
< error
| 1064 42000 …near 'CONDITION' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CONNECTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONSISTENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONSTRAINT
< error
| 1064 42000 …near 'CONSTRAINT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CONSTRAINT_CATALOG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONSTRAINT_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONSTRAINT_SCHEMA
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONTAINS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONTEXT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CONTINUE
< error
| 1064 42000 …near 'CONTINUE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CONVERT
< error
| 1064 42000 …near 'CONVERT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base COUNT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CPU
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CREATE
< error
| 1064 42000 …near 'CREATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CROSS
< error
| 1064 42000 …near 'CROSS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CSV
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CUBE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CUME_DIST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CURDATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CURRENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CURRENT_DATE
< error
| 1064 42000 …near 'CURRENT_DATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CURRENT_TIME
< error
| 1064 42000 …near 'CURRENT_TIME' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CURRENT_TIMESTAMP
< error
| 1064 42000 …near 'CURRENT_TIMESTAMP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CURRENT_USER
< error
| 1064 42000 …near 'CURRENT_USER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CURSOR
< error
| 1064 42000 …near 'CURSOR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base CURSOR_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base CURTIME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DATA
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DATABASE
< error
| 1064 42000 …near 'DATABASE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DATABASES
< error
| 1064 42000 …near 'DATABASES' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DATAFILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DATETIME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DATE_ADD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DATE_SUB
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DAY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DAY_HOUR
< error
| 1064 42000 …near 'DAY_HOUR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DAY_MICROSECOND
< error
| 1064 42000 …near 'DAY_MICROSECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DAY_MINUTE
< error
| 1064 42000 …near 'DAY_MINUTE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DAY_SECOND
< error
| 1064 42000 …near 'DAY_SECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DEALLOCATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DEC
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DECIMAL
< error
| 1064 42000 …near 'DECIMAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DECLARE
< error
| 1064 42000 …near 'DECLARE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DEFAULT_AUTH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DEFINER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DEFINITION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DELAYED
< error
| 1064 42000 …near 'DELAYED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DELAY_KEY_WRITE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DELETE
< error
| 1064 42000 …near 'DELETE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DENSE_RANK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DESC
< error
| 1064 42000 …near 'DESC' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DESCRIBE
< error
| 1064 42000 …near 'DESCRIBE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DESCRIPTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DETERMINISTIC
< error
| 1064 42000 …near 'DETERMINISTIC' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DIAGNOSTICS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DIRECTORY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DISABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DISCARD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DISK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DISTINCT
< error
| 1064 42000 …near 'DISTINCT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DISTINCTROW
< error
| 1064 42000 …near 'DISTINCTROW' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DIV
< error
| 1064 42000 …near 'DIV' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DO
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DOUBLE
< error
| 1064 42000 …near 'DOUBLE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DROP
< error
| 1064 42000 …near 'DROP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DUAL
< error
| 1064 42000 …near 'DUAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base DUALITY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DUMPFILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DUPLICATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base DYNAMIC
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EACH
< error
| 1064 42000 …near 'EACH' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ELSE
< error
| 1064 42000 …near 'ELSE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ELSEIF
< error
| 1064 42000 …near 'ELSEIF' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base EMPTY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENCLOSED
< error
| 1064 42000 …near 'ENCLOSED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ENCRYPTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base END
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENFORCED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENGINE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENGINES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENGINE_ATTRIBUTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ENUM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ERROR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ERRORS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ESCAPE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ESCAPED
< error
| 1064 42000 …near 'ESCAPED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base EVENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EVENTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EVERY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXCEPT
< error
| 1064 42000 …near 'EXCEPT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base EXCHANGE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXCLUDE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXECUTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXISTS
< error
| 1064 42000 …near 'EXISTS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base EXIT
< error
| 1064 42000 …near 'EXIT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base EXPANSION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXPIRE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXPLAIN
< error
| 1064 42000 …near 'EXPLAIN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base EXPORT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXTENDED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXTENT_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXTERNAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXTERNAL_FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base EXTRACT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FACTOR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FAILED_LOGIN_ATTEMPTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FALSE
< error
| 1064 42000 …near 'FALSE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FAST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FAULTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FETCH
< error
| 1064 42000 …near 'FETCH' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FIELDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILE_BLOCK_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILE_FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILE_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILE_PATTERN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILE_PREFIX
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FILTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FINISH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FIRST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FIRST_VALUE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FIXED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FLOAT
< error
| 1064 42000 …near 'FLOAT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FLOAT4
< error
| 1064 42000 …near 'FLOAT4' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FLOAT8
< error
| 1064 42000 …near 'FLOAT8' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FLUSH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FOLLOWING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FOLLOWS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FOR
< error
| 1064 42000 …near 'FOR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FORCE
< error
| 1064 42000 …near 'FORCE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FOREIGN
< error
| 1064 42000 …near 'FOREIGN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FOUND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FROM
< error
| 1064 42000 …near 'FROM' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FULL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base FULLTEXT
< error
| 1064 42000 …near 'FULLTEXT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base FUNCTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GENERAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GENERATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GENERATED
< error
| 1064 42000 …near 'GENERATED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base GEOMCOLLECTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GEOMETRY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GEOMETRYCOLLECTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GET
< error
| 1064 42000 …near 'GET' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base GET_FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GET_SOURCE_PUBLIC_KEY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GLOBAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GRANT
< error
| 1064 42000 …near 'GRANT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base GRANTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GREATEST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GROUP
< error
| 1064 42000 …near 'GROUP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base GROUPING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GROUPS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GROUP_CONCAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GROUP_REPLICATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GTIDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GTID_ONLY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base GUIDED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HANDLER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HASH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HAVING
< error
| 1064 42000 …near 'HAVING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base HEADER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HELP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HIGH_PRIORITY
< error
| 1064 42000 …near 'HIGH_PRIORITY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base HISTOGRAM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HISTORY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HOST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HOSTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HOUR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base HOUR_MICROSECOND
< error
| 1064 42000 …near 'HOUR_MICROSECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base HOUR_MINUTE
< error
| 1064 42000 …near 'HOUR_MINUTE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base HOUR_SECOND
< error
| 1064 42000 …near 'HOUR_SECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base IDENTIFIED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IF
< error
| 1064 42000 …near 'IF' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base IGNORE
< error
| 1064 42000 …near 'IGNORE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base IGNORE_SERVER_IDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IMPORT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IN
< error
| 1064 42000 …near 'IN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INACTIVE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INDEXES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INFILE
< error
| 1064 42000 …near 'INFILE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INITIAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INITIAL_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INITIATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INNER
< error
| 1064 42000 …near 'INNER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INOUT
< error
| 1064 42000 …near 'INOUT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INSENSITIVE
< error
| 1064 42000 …near 'INSENSITIVE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INSERT
< error
| 1064 42000 …near 'INSERT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INSERT_METHOD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INSTALL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INSTANCE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INT
< error
| 1064 42000 …near 'INT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INT1
< error
| 1064 42000 …near 'INT1' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INT2
< error
| 1064 42000 …near 'INT2' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INT3
< error
| 1064 42000 …near 'INT3' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INT4
< error
| 1064 42000 …near 'INT4' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INT8
< error
| 1064 42000 …near 'INT8' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INTEGER
< error
| 1064 42000 …near 'INTEGER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INTERSECT
< error
| 1064 42000 …near 'INTERSECT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INTERVAL
< error
| 1064 42000 …near 'INTERVAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INTO
< error
| 1064 42000 …near 'INTO' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base INVISIBLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base INVOKER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IO
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IO_AFTER_GTIDS
< error
| 1064 42000 …near 'IO_AFTER_GTIDS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base IO_BEFORE_GTIDS
< error
| 1064 42000 …near 'IO_BEFORE_GTIDS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base IO_THREAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IPC
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base IS
< error
| 1064 42000 …near 'IS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ISOLATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ISSUER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ITERATE
< error
| 1064 42000 …near 'ITERATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base JOIN
< error
| 1064 42000 …near 'JOIN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base JSON
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base JSON_ARRAYAGG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base JSON_DUALITY_OBJECT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base JSON_OBJECTAGG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base JSON_TABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base JSON_VALUE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base KEY
< error
| 1064 42000 …near 'KEY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base KEYRING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base KEYS
< error
| 1064 42000 …near 'KEYS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base KEY_BLOCK_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base KILL
< error
| 1064 42000 …near 'KILL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LAG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LANGUAGE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LAST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LAST_VALUE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LATERAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LEAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LEADING
< error
| 1064 42000 …near 'LEADING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LEAST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LEAVE
< error
| 1064 42000 …near 'LEAVE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LEAVES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LEFT
< error
| 1064 42000 …near 'LEFT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LESS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LEVEL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LIBRARY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LIKE
< error
| 1064 42000 …near 'LIKE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LIMIT
< error
| 1064 42000 …near 'LIMIT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LINEAR
< error
| 1064 42000 …near 'LINEAR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LINES
< error
| 1064 42000 …near 'LINES' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LINESTRING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LIST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LOAD
< error
| 1064 42000 …near 'LOAD' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LOCAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LOCALTIME
< error
| 1064 42000 …near 'LOCALTIME' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LOCALTIMESTAMP
< error
| 1064 42000 …near 'LOCALTIMESTAMP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LOCK
< error
| 1064 42000 …near 'LOCK' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LOCKED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LOCKS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LOG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LOGFILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LOGS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base LONG
< error
| 1064 42000 …near 'LONG' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LONGBLOB
< error
| 1064 42000 …near 'LONGBLOB' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LONGTEXT
< error
| 1064 42000 …near 'LONGTEXT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LOOP
< error
| 1064 42000 …near 'LOOP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base LOW_PRIORITY
< error
| 1064 42000 …near 'LOW_PRIORITY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MANIFEST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MANUAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MASTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MATCH
< error
| 1064 42000 …near 'MATCH' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MATERIALIZED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAXVALUE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX_CONNECTIONS_PER_HOUR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX_QUERIES_PER_HOUR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX_ROWS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX_UPDATES_PER_HOUR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MAX_USER_CONNECTIONS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MEDIUM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MEDIUMBLOB
< error
| 1064 42000 …near 'MEDIUMBLOB' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MEDIUMINT
< error
| 1064 42000 …near 'MEDIUMINT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MEDIUMTEXT
< error
| 1064 42000 …near 'MEDIUMTEXT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MEMBER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MEMORY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MERGE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MESSAGE_TEXT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MICROSECOND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MID
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MIDDLEINT
< error
| 1064 42000 …near 'MIDDLEINT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MIGRATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MINUTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MINUTE_MICROSECOND
< error
| 1064 42000 …near 'MINUTE_MICROSECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MINUTE_SECOND
< error
| 1064 42000 …near 'MINUTE_SECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MIN_ROWS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MOD
< error
| 1064 42000 …near 'MOD' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MODE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MODIFIES
< error
| 1064 42000 …near 'MODIFIES' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base MODIFY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MONTH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MULTILINESTRING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MULTIPOINT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MULTIPOLYGON
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MUTEX
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base MYSQL_ERRNO
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NAMES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NATIONAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NATURAL
< error
| 1064 42000 …near 'NATURAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base NCHAR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NDB
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NDBCLUSTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NESTED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NETWORK_NAMESPACE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NEVER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NEW
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NEXT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NO
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NODEGROUP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NONE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base NOW
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NOWAIT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NO_WAIT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NO_WRITE_TO_BINLOG
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base NTH_VALUE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NTILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base NULLIF
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NULLS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NUMBER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base NUMERIC
< error
| 1064 42000 …near 'NUMERIC' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base NVARCHAR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OF
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OFF
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OFFSET
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OJ
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OLD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ON
< error
| 1064 42000 …near 'ON' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ONE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ONLY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OPEN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OPTIMIZE
< error
| 1064 42000 …near 'OPTIMIZE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base OPTIMIZER_COSTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OPTION
< error
| 1064 42000 …near 'OPTION' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base OPTIONAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OPTIONALLY
< error
| 1064 42000 …near 'OPTIONALLY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base OPTIONS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OR
< error
| 1064 42000 …near 'OR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ORDER
< error
| 1064 42000 …near 'ORDER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ORDINALITY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ORGANIZATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OTHERS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OUT
< error
| 1064 42000 …near 'OUT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base OUTER
< error
| 1064 42000 …near 'OUTER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base OUTFILE
< error
| 1064 42000 …near 'OUTFILE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base OVER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OVERLAY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OVERWRITE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base OWNER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PACK_KEYS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PAGE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARALLEL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARAMETERS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARSER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARSE_TREE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARTIAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARTITION
< error
| 1064 42000 …near 'PARTITION' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base PARTITIONING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PARTITIONS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PASSWORD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PASSWORD_LOCK_TIME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PATH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PERCENT_RANK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PERSIST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PERSIST_ONLY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PHASE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PLACING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PLUGIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PLUGINS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PLUGIN_DIR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base POINT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base POLYGON
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PORT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base POSITION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PRECEDES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PRECEDING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PRECISION
< error
| 1064 42000 …near 'PRECISION' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base PREFIX
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PREPARE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PRESERVE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PREV
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PRIMARY
< error
| 1064 42000 …near 'PRIMARY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base PRIVILEGES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PRIVILEGE_CHECKS_USER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PROCEDURE
< error
| 1064 42000 …near 'PROCEDURE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base PROCESS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PROCESSLIST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PROFILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PROFILES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PROXY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base PURGE
< error
| 1064 42000 …near 'PURGE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base QUALIFY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base QUARTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base QUERY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base QUICK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RANDOM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RANGE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RANK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base READ
< error
| 1064 42000 …near 'READ' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base READS
< error
| 1064 42000 …near 'READS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base READ_ONLY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base READ_WRITE
< error
| 1064 42000 …near 'READ_WRITE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REAL
< error
| 1064 42000 …near 'REAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REBUILD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RECOVER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RECURSIVE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REDO_BUFFER_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REDUNDANT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REFERENCE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REFERENCES
< error
| 1064 42000 …near 'REFERENCES' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REGEXP
< error
| 1064 42000 …near 'REGEXP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REGISTRATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELATIONAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELAY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELAYLOG
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELAY_LOG_FILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELAY_LOG_POS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELAY_THREAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RELEASE
< error
| 1064 42000 …near 'RELEASE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base RELOAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REMOVE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RENAME
< error
| 1064 42000 …near 'RENAME' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REORGANIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPAIR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPEAT
< error
| 1064 42000 …near 'REPEAT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REPEATABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLACE
< error
| 1064 42000 …near 'REPLACE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICA
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICAS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_DO_DB
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_DO_TABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_IGNORE_DB
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_IGNORE_TABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_REWRITE_DB
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_WILD_DO_TABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATE_WILD_IGNORE_TABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REPLICATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REQUIRE
< error
| 1064 42000 …near 'REQUIRE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base REQUIRE_ROW_FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REQUIRE_TABLE_PRIMARY_KEY_CHECK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RESET
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RESIGNAL
< error
| 1064 42000 …near 'RESIGNAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base RESOURCE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RESPECT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RESTART
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RESTORE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RESTRICT
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base RESUME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RETAIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RETURN
< error
| 1064 42000 …near 'RETURN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base RETURNED_SQLSTATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RETURNING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RETURNS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REUSE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REVERSE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base REVOKE
< error
| 1064 42000 …near 'REVOKE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base RIGHT
< error
| 1064 42000 …near 'RIGHT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base RLIKE
< error
| 1064 42000 …near 'RLIKE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ROLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROLLBACK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROLLUP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROTATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROUTINE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROW
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROWS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROW_COUNT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROW_FORMAT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ROW_NUMBER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base RTREE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base S3
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SAVEPOINT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SCHEDULE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SCHEMA
< error
| 1064 42000 …near 'SCHEMA' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SCHEMAS
< error
| 1064 42000 …near 'SCHEMAS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SCHEMA_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECOND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECONDARY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECONDARY_ENGINE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECONDARY_ENGINE_ATTRIBUTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECONDARY_LOAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECONDARY_UNLOAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SECOND_MICROSECOND
< error
| 1064 42000 …near 'SECOND_MICROSECOND' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SECURITY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SELECT
< error
| 1064 42000 …near 'SELECT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SENSITIVE
< error
| 1064 42000 …near 'SENSITIVE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SEPARATOR
< error
| 1064 42000 …near 'SEPARATOR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SERIAL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SERIALIZABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SERVER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SESSION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SESSION_USER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SET
< error
| 1064 42000 …near 'SET' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SETS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SHARE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SHOW
< error
| 1064 42000 …near 'SHOW' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SHUTDOWN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SIGNAL
< error
| 1064 42000 …near 'SIGNAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SIGNED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SIMPLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SKIP
< error
| 1064 42000 …near 'SKIP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SLAVE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SLOW
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SMALLINT
< error
| 1064 42000 …near 'SMALLINT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SNAPSHOT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOCKET
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SONAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOUNDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_AUTO_POSITION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_BIND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_COMPRESSION_ALGORITHMS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_CONNECTION_AUTO_FAILOVER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_CONNECT_RETRY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_DELAY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_HEARTBEAT_PERIOD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_HOST
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_LOG_FILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_LOG_POS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_PASSWORD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_PORT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_PUBLIC_KEY_PATH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_RETRY_COUNT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_CA
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_CAPATH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_CERT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_CIPHER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_CRL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_CRLPATH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_KEY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_SSL_VERIFY_SERVER_CERT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_TLS_CIPHERSUITES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_TLS_VERSION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_USER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SOURCE_ZSTD_COMPRESSION_LEVEL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SPATIAL
< error
| 1064 42000 …near 'SPATIAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SPECIFIC
< error
| 1064 42000 …near 'SPECIFIC' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQL
< error
| 1064 42000 …near 'SQL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQLEXCEPTION
< error
| 1064 42000 …near 'SQLEXCEPTION' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQLSTATE
< error
| 1064 42000 …near 'SQLSTATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQLWARNING
< error
| 1064 42000 …near 'SQLWARNING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_AFTER_GTIDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_AFTER_MTS_GAPS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_BEFORE_GTIDS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_BIG_RESULT
< error
| 1064 42000 …near 'SQL_BIG_RESULT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_BUFFER_RESULT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_CALC_FOUND_ROWS
< error
| 1064 42000 …near 'SQL_CALC_FOUND_ROWS' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_NO_CACHE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_SMALL_RESULT
< error
| 1064 42000 …near 'SQL_SMALL_RESULT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_THREAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_DAY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_HOUR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_MINUTE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_MONTH
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_QUARTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_SECOND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_WEEK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SQL_TSI_YEAR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SRID
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SSE_KMS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SSE_S3
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SSL
< error
| 1064 42000 …near 'SSL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base STACKED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base START
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STARTING
< error
| 1064 42000 …near 'STARTING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base STARTS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STATS_AUTO_RECALC
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STATS_PERSISTENT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STATS_SAMPLE_PAGES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STATUS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STDDEV
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STDDEV_POP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STDDEV_SAMP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STOP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STORAGE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STORED
< error
| 1064 42000 …near 'STORED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base STRAIGHT_JOIN
< error
| 1064 42000 …near 'STRAIGHT_JOIN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base STREAM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STRICT_LOAD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base STRING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base ST_COLLECT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBCLASS_ORIGIN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBDATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBJECT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBPARTITION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBPARTITIONS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBSTR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUBSTRING
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUPER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SUSPEND
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SWAPS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SWITCHES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SYSDATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SYSTEM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base SYSTEM_USER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TABLE
< error
| 1064 42000 …near 'TABLE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TABLES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TABLESAMPLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TABLESPACE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TABLE_CHECKSUM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TABLE_NAME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TEMPORARY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TEMPTABLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TERMINATED
< error
| 1064 42000 …near 'TERMINATED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TEXT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base THAN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base THEN
< error
| 1064 42000 …near 'THEN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base THREAD_PRIORITY
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TIES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TIME
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TIMESTAMP
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TIMESTAMPADD
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TIMESTAMPDIFF
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TINYBLOB
< error
| 1064 42000 …near 'TINYBLOB' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TINYINT
< error
| 1064 42000 …near 'TINYINT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TINYTEXT
< error
| 1064 42000 …near 'TINYTEXT' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TLS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TO
< error
| 1064 42000 …near 'TO' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TRAILING
< error
| 1064 42000 …near 'TRAILING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TRANSACTION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TRIGGER
< error
| 1064 42000 …near 'TRIGGER' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TRIGGERS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TRIM
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TRUE
< error
| 1064 42000 …near 'TRUE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base TRUNCATE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TYPE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base TYPES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNBOUNDED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNCOMMITTED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNDEFINED
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNDO
< error
| 1064 42000 …near 'UNDO' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UNDOFILE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNDO_BUFFER_SIZE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNICODE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNINSTALL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNION
< error
| 1064 42000 …near 'UNION' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UNIQUE
< error
| 1064 42000 …near 'UNIQUE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UNKNOWN
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNLOCK
< error
| 1064 42000 …near 'UNLOCK' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UNREGISTER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UNSIGNED
< error
| 1064 42000 …near 'UNSIGNED' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UNTIL
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base UPDATE
< error
| 1064 42000 …near 'UPDATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UPGRADE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base URI
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base USAGE
< error
| 1064 42000 …near 'USAGE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base USE
< error
| 1064 42000 …near 'USE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base USER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base USING
< error
| 1064 42000 …near 'USING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UTC_DATE
< error
| 1064 42000 …near 'UTC_DATE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UTC_TIME
< error
| 1064 42000 …near 'UTC_TIME' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base UTC_TIMESTAMP
< error
| 1064 42000 …near 'UTC_TIMESTAMP' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base VALIDATION
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base VALUE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base VALUES
< error
| 1064 42000 …near 'VALUES' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base VARBINARY
< error
| 1064 42000 …near 'VARBINARY' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base VARCHAR
< error
| 1064 42000 …near 'VARCHAR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base VARIABLES
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base VARYING
< error
| 1064 42000 …near 'VARYING' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base VIEW
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base VIRTUAL
< error
| 1064 42000 …near 'VIRTUAL' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base VISIBLE
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base WAIT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base WARNINGS
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base WHEN
< error
| 1064 42000 …near 'WHEN' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base WHERE
< error
| 1064 42000 …near 'WHERE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base WHILE
< error
| 1064 42000 …near 'WHILE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base WINDOW
< error
| 1064 42000 …near 'WINDOW' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base WITH
< error
| 1064 42000 …near 'WITH' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base WITHOUT
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base WORK
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base WRAPPER
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base WRITE
< error
| 1064 42000 …near 'WRITE' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base X509
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base XA
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base XID
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base XML
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base XOR
< error
| 1064 42000 …near 'XOR' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base YEAR
< ok 0 0
> exec ^
: string
| SELECT * FROM kw_alias_base YEAR_MONTH
< error
| 1064 42000 …near 'YEAR_MONTH' at line 1
> exec ^
: string
| SELECT * FROM kw_alias_base ZEROFILL
< error
| 1064 42000 …near 'ZEROFILL' at line 1

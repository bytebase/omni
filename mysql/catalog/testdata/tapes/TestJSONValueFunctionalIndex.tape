-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestJSONValueFunctionalIndex default#1
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' ))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varchar(512)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING SIGNED))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| bigint
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING UNSIGNED))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| bigint unsigned
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(10)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varchar(10)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DECIMAL(5,2)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| decimal(5,2)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DATETIME(3)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| datetime(3)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DATE))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| date
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DOUBLE))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| double
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING JSON))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING BINARY(3)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varbinary(3)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING TIME(2)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| time(2)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DATETIME))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| datetime
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING FLOAT))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| float
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING YEAR))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| year
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING NCHAR(3)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varchar(3)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(3) CHARSET latin1))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varchar(3)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DECIMAL))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| decimal(10,0)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(512)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varchar(512)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(513)))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING BINARY(512)))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varbinary(512)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING BINARY(513)))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING BINARY))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(600) CHARSET latin1))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(3) CHARSET binary))))
< ok 0 0
> query
| SHOW EXTENDED COLUMNS FROM f
< result
< column "Field" "VARCHAR"
< column "Type" "TEXT"
< column "Null" "VARCHAR"
< column "Key" "ENUM"
< column "Default" "TEXT"
< column "Extra" "VARCHAR"
< row
: bytes
| a
: bytes
| json
: bytes
| YES
: bytes
|
: null
: bytes
|
< row
: bytes
| !hidden!k!0!0
: bytes
| varbinary(3)
: bytes
| YES
: bytes
| MUL
: null
: bytes
| VIRTUAL GENERATED
< row
: bytes
| DB_ROW_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_TRX_ID
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
< row
: bytes
| DB_ROLL_PTR
: bytes
|
: bytes
| NO
: bytes
|
: null
: bytes
|
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR CHARSET binary))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.

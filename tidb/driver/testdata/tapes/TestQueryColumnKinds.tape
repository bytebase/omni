-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQueryColumnKinds default#1
> exec
| USE `omni_qe_kinds`
< ok 0 0
> exec
| CREATE TABLE qkinds (
|   id BIGINT PRIMARY KEY,
|   body BLOB,
|   tiny TINYBLOB,
|   medium MEDIUMBLOB,
|   large LONGBLOB,
|   raw VARBINARY(8),
|   fixedb BINARY(4),
|   ubig BIGINT UNSIGNED,
|   uintc INT UNSIGNED,
|   umedium MEDIUMINT UNSIGNED,
|   usmall SMALLINT UNSIGNED,
|   utiny TINYINT UNSIGNED
| );
< ok 0 0
> exec
|
| INSERT INTO qkinds VALUES
|   (1, 0x00ff, 0x01, 0x02, 0x03, 0x10, 0xdeadbeef, 18446744073709551615, 4294967295, 16777215, 65535, 255),
|   (2, '', '', '', '', '', '', 0, 0, 0, 0, 0),
|   (3, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
< ok 3 0
> query
| SELECT * FROM qkinds ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "body" "BLOB" "[]uint8"
< column "tiny" "BLOB" "[]uint8"
< column "medium" "BLOB" "[]uint8"
< column "large" "BLOB" "[]uint8"
< column "raw" "VARBINARY" "[]uint8"
< column "fixedb" "BINARY" "[]uint8"
< column "ubig" "UNSIGNED BIGINT"
< column "uintc" "UNSIGNED INT"
< column "umedium" "UNSIGNED MEDIUMINT"
< column "usmall" "UNSIGNED SMALLINT"
< column "utiny" "UNSIGNED TINYINT"
< row
: int64 1
: bytes
|q "\x00\xff"
: bytes
|q "\x01"
: bytes
|q "\x02"
: bytes
|q "\x03"
: bytes
|q "\x10"
: bytes
|q "ޭ\xbe\xef"
: uint64 18446744073709551615
: int64 4294967295
: int64 16777215
: int64 65535
: int64 255
< row
: int64 2
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|q "\x00\x00\x00\x00"
: uint64 0
: int64 0
: int64 0
: int64 0
: int64 0
< row
: int64 3
: null
: null
: null
: null
: null
: null
: null
: null
: null
: null
: null

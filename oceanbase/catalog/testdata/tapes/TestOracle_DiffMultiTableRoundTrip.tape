-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_DiffMultiTableRoundTrip default#1
> exec
| DROP DATABASE IF EXISTS diff_multi
< ok 0 0
> exec
| CREATE DATABASE diff_multi DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE diff_multi
< ok 0 0
> exec
| CREATE TABLE parent (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50) NOT NULL)
< ok 0 0
> exec
| CREATE TABLE child (id INT NOT NULL PRIMARY KEY, pid INT, amount DECIMAL(10,2) DEFAULT 0, KEY pid_idx (pid), CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent (id))
< ok 0 0
> query
| SHOW CREATE TABLE diff_multi.parent
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| parent
: bytes
| CREATE TABLE `parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE diff_multi.child
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| child
: bytes
| CREATE TABLE `child` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   `amount` decimal(10,2) DEFAULT '0.00',
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk_pid` FOREIGN KEY (`pid`) REFERENCES `diff_multi`.`parent`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `pid_idx` (`pid`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0

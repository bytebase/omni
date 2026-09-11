# MySQL Parser Keyword & Lexer Final Alignment Scenarios

> Goal: Drive all 3 golden tests to zero — TestKeywordCompleteness (253→0), TestKeywordClassification (132→0), TestNoEqFoldForRegisteredKeywords (5→0) — plus fix related lexer gaps
> Verification: `go test -short ./mysql/parser/... -run "TestKeyword|TestNoEqFold" -count=1` + `go test -short ./mysql/parser/... -count=1` (full suite, no regressions)
> Reference sources: mysql-server sql/lex.h, sql/sql_yacc.yy, sql/sql_lex.cc; golden list in keyword_completeness_test.go

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Lexer Fixes

Must complete before registering more reserved keywords.

### 1.1 Lexer Dot-Context: Suppress Keyword Lookup After `.`

MySQL's lexer does NOT do keyword lookup after `.` — any word after a dot is an identifier. omni always does keyword lookup, causing `t.select`, `t.from`, `db.table` to fail when the name is a reserved word.

- [x] `SELECT t.select FROM t` parses — `select` after dot treated as identifier
- [x] `SELECT t.from FROM t` parses — `from` after dot treated as identifier
- [x] `SELECT t.table FROM t` parses — `table` after dot treated as identifier
- [x] `SELECT t.where FROM t` parses — `where` after dot treated as identifier
- [x] `SELECT db.select FROM db.select` parses — both schema and table after dot
- [x] `SELECT t.key FROM t` parses — `key` after dot
- [x] `SELECT t.index FROM t` parses — `index` after dot
- [x] `SELECT t.group FROM t` parses — `group` after dot
- [x] `CREATE TABLE t (a INT); SELECT t.order FROM t` — `order` after dot
- [x] `SELECT a.b.c FROM t a` — three-part qualified name works
- [x] All existing parser tests still pass — no regressions
- [x] Backtick-quoted identifiers after dot still work: `` SELECT t.`select` FROM t ``

### 1.2 Lexer @@Variable: Emit Separate Tokens

MySQL emits `@@global.var_name` as 5 tokens: `@`, `@`, `GLOBAL_SYM`, `.`, `ident`. omni emits 1 tokIDENT. This must change so the parser can use keyword tokens for GLOBAL/SESSION/LOCAL.

- [ ] `@@var_name` lexes as: `@@` prefix token + `var_name` ident token
- [ ] `@@global.var_name` lexes as: `@@` prefix token + `global` keyword token + `.` + `var_name` ident token
- [ ] `@@session.var_name` lexes as: `@@` prefix token + `session` keyword token + `.` + `var_name` ident token
- [ ] `@@local.var_name` lexes as: `@@` prefix token + `local` keyword token + `.` + `var_name` ident token
- [ ] `@@var_name` without scope prefix still works (no dot, no scope keyword)
- [ ] `SELECT @@global.max_connections` parses correctly
- [ ] `SELECT @@session.wait_timeout` parses correctly
- [ ] `SELECT @@version` parses correctly (no scope)
- [ ] `SET @@global.max_connections = 100` parses correctly
- [ ] `SET @@session.wait_timeout = 28800` parses correctly
- [ ] Parser variable reference handling updated to use tokens instead of eqFold string splitting
- [ ] All existing SET/SHOW/SELECT @@variable tests still pass

### 1.3 String Escaping Fixes

- [ ] `\b` (backspace, 0x08) correctly escaped in string literals
- [ ] `\Z` (Ctrl-Z, 0x1A) correctly escaped in string literals
- [ ] `\_` in string literals preserves the backslash (for LIKE patterns)
- [ ] `\%` in string literals preserves the backslash (for LIKE patterns)
- [ ] Existing string escape tests still pass (`\n`, `\t`, `\r`, `\0`, `\\`, `\'`, `\"`)

---

## Phase 2: Keyword Registration

Register all 253 missing MySQL 8.0 keywords. Grouped alphabetically into batches of ~50.

### 2.1 Register Missing Keywords — Batch A-C

- [ ] Register `absent` (unambiguous)
- [ ] Register `adddate` (unambiguous)
- [ ] Register `allow_missing_files` (unambiguous)
- [ ] Register `any` (unambiguous)
- [ ] Register `array` (unambiguous)
- [ ] Register `ascii` (ambiguous_2)
- [ ] Register `assign_gtids_to_anonymous_transactions` (unambiguous)
- [ ] Register `authentication` (unambiguous)
- [ ] Register `auto` (unambiguous)
- [ ] Register `auto_refresh` (unambiguous)
- [ ] Register `auto_refresh_source` (unambiguous)
- [ ] Register `autoextend_size` (unambiguous)
- [ ] Register `avg_row_length` (unambiguous)
- [ ] Register `bernoulli` (unambiguous)
- [ ] Register `bit_and` (reserved)
- [ ] Register `bit_or` (reserved)
- [ ] Register `bit_xor` (reserved)
- [ ] Register `block` (unambiguous)
- [ ] Register `buckets` (unambiguous)
- [ ] Register `bulk` (unambiguous)
- [ ] Register `byte` (ambiguous_2)
- [ ] Register `catalog_name` (unambiguous)
- [ ] Register `challenge_response` (unambiguous)
- [ ] Register `class_origin` (unambiguous)
- [ ] Register `client` (unambiguous)
- [ ] Register `column_name` (unambiguous)
- [ ] Register `constraint_catalog` (unambiguous)
- [ ] Register `constraint_name` (unambiguous)
- [ ] Register `constraint_schema` (unambiguous)
- [ ] Register `context` (unambiguous)
- [ ] Register `cpu` (unambiguous)
- [ ] Register `curdate` (reserved)
- [ ] Register `cursor_name` (unambiguous)
- [ ] Register `curtime` (reserved)
- [ ] All newly registered keywords lex correctly and tests pass

### 2.2 Register Missing Keywords — Batch D-F

- [ ] Register `date_add` (reserved)
- [ ] Register `date_sub` (reserved)
- [ ] Register `default_auth` (unambiguous)
- [ ] Register `delay_key_write` (unambiguous)
- [ ] Register `duality` (unambiguous)
- [ ] Register `engine_attribute` (unambiguous)
- [ ] Register `exclude` (unambiguous)
- [ ] Register `extent_size` (unambiguous)
- [ ] Register `external` (reserved)
- [ ] Register `external_format` (unambiguous)
- [ ] Register `factor` (unambiguous)
- [ ] Register `failed_login_attempts` (unambiguous)
- [ ] Register `faults` (unambiguous)
- [ ] Register `file` (ambiguous_3)
- [ ] Register `file_block_size` (unambiguous)
- [ ] Register `file_format` (unambiguous)
- [ ] Register `file_name` (unambiguous)
- [ ] Register `file_pattern` (unambiguous)
- [ ] Register `file_prefix` (unambiguous)
- [ ] Register `files` (unambiguous)
- [ ] Register `finish` (unambiguous)
- [ ] Register `float4` (reserved)
- [ ] Register `float8` (reserved)
- [ ] All newly registered keywords lex correctly and tests pass

### 2.3 Register Missing Keywords — Batch G-L

- [ ] Register `general` (unambiguous)
- [ ] Register `geomcollection` (unambiguous)
- [ ] Register `get_format` (unambiguous)
- [ ] Register `get_source_public_key` (unambiguous)
- [ ] Register `group_replication` (unambiguous)
- [ ] Register `gtid_only` (unambiguous)
- [ ] Register `gtids` (unambiguous)
- [ ] Register `guided` (unambiguous)
- [ ] Register `header` (unambiguous)
- [ ] Register `histogram` (unambiguous)
- [ ] Register `host` (unambiguous)
- [ ] Register `hour` (unambiguous)
- [ ] Register `ignore_server_ids` (unambiguous)
- [ ] Register `initial` (unambiguous)
- [ ] Register `initial_size` (unambiguous)
- [ ] Register `initiate` (unambiguous)
- [ ] Register `int1` (reserved)
- [ ] Register `int2` (reserved)
- [ ] Register `int3` (reserved)
- [ ] Register `int4` (reserved)
- [ ] Register `int8` (reserved)
- [ ] Register `io` (unambiguous)
- [ ] Register `io_thread` (unambiguous)
- [ ] Register `ipc` (unambiguous)
- [ ] Register `json_arrayagg` (reserved)
- [ ] Register `json_duality_object` (reserved)
- [ ] Register `json_objectagg` (reserved)
- [ ] Register `json_value` (unambiguous)
- [ ] Register `key_block_size` (unambiguous)
- [ ] Register `library` (reserved)
- [ ] Register `locks` (unambiguous)
- [ ] Register `log` (unambiguous)
- [ ] Register `long` (reserved)
- [ ] All newly registered keywords lex correctly and tests pass

### 2.4 Register Missing Keywords — Batch M-P

- [ ] Register `manual` (reserved)
- [ ] Register `materialized` (unambiguous)
- [ ] Register `max_connections_per_hour` (unambiguous)
- [ ] Register `max_queries_per_hour` (unambiguous)
- [ ] Register `max_rows` (unambiguous)
- [ ] Register `max_size` (unambiguous)
- [ ] Register `max_updates_per_hour` (unambiguous)
- [ ] Register `max_user_connections` (unambiguous)
- [ ] Register `message_text` (unambiguous)
- [ ] Register `microsecond` (unambiguous)
- [ ] Register `mid` (reserved)
- [ ] Register `middleint` (reserved)
- [ ] Register `min_rows` (unambiguous)
- [ ] Register `minute` (unambiguous)
- [ ] Register `month` (unambiguous)
- [ ] Register `mysql_errno` (unambiguous)
- [ ] Register `ndb` (unambiguous)
- [ ] Register `ndbcluster` (unambiguous)
- [ ] Register `network_namespace` (unambiguous)
- [ ] Register `new` (unambiguous)
- [ ] Register `no_wait` (unambiguous)
- [ ] Register `nodegroup` (unambiguous)
- [ ] Register `now` (reserved)
- [ ] Register `number` (unambiguous)
- [ ] Register `off` (unambiguous)
- [ ] Register `oj` (unambiguous)
- [ ] Register `others` (unambiguous)
- [ ] Register `owner` (unambiguous)
- [ ] Register `pack_keys` (unambiguous)
- [ ] Register `page` (unambiguous)
- [ ] Register `parallel` (reserved)
- [ ] Register `parameters` (unambiguous)
- [ ] Register `parse_tree` (unambiguous)
- [ ] Register `password_lock_time` (unambiguous)
- [ ] Register `persist_only` (ambiguous_4)
- [ ] Register `plugin_dir` (unambiguous)
- [ ] Register `port` (unambiguous)
- [ ] Register `privilege_checks_user` (unambiguous)
- [ ] All newly registered keywords lex correctly and tests pass

### 2.5 Register Missing Keywords — Batch Q-S (part 1)

- [ ] Register `qualify` (reserved)
- [ ] Register `quarter` (unambiguous)
- [ ] Register `read_only` (unambiguous)
- [ ] Register `read_write` (reserved)
- [ ] Register `redo_buffer_size` (unambiguous)
- [ ] Register `registration` (unambiguous)
- [ ] Register `relational` (unambiguous)
- [ ] Register `relay` (unambiguous)
- [ ] Register `relay_log_file` (unambiguous)
- [ ] Register `relay_log_pos` (unambiguous)
- [ ] Register `relay_thread` (unambiguous)
- [ ] Register `replicate_do_db` (unambiguous)
- [ ] Register `replicate_do_table` (unambiguous)
- [ ] Register `replicate_ignore_db` (unambiguous)
- [ ] Register `replicate_ignore_table` (unambiguous)
- [ ] Register `replicate_rewrite_db` (unambiguous)
- [ ] Register `replicate_wild_do_table` (unambiguous)
- [ ] Register `replicate_wild_ignore_table` (unambiguous)
- [ ] Register `require_row_format` (unambiguous)
- [ ] Register `require_table_primary_key_check` (unambiguous)
- [ ] Register `respect` (unambiguous)
- [ ] Register `restore` (unambiguous)
- [ ] Register `returned_sqlstate` (unambiguous)
- [ ] Register `returning` (unambiguous)
- [ ] Register `reverse` (unambiguous)
- [ ] Register `row_count` (unambiguous)
- [ ] Register `rtree` (unambiguous)
- [ ] Register `s3` (unambiguous)
- [ ] All newly registered keywords lex correctly and tests pass

### 2.6 Register Missing Keywords — Batch S (part 2)

- [ ] Register `schema_name` (unambiguous)
- [ ] Register `second` (unambiguous)
- [ ] Register `secondary` (unambiguous)
- [ ] Register `secondary_engine` (unambiguous)
- [ ] Register `secondary_engine_attribute` (unambiguous)
- [ ] Register `secondary_load` (unambiguous)
- [ ] Register `secondary_unload` (unambiguous)
- [ ] Register `session_user` (unambiguous)
- [ ] Register `sets` (reserved)
- [ ] Register `simple` (unambiguous)
- [ ] Register `slow` (unambiguous)
- [ ] Register `socket` (unambiguous)
- [ ] Register `some` (unambiguous)
- [ ] Register `source_auto_position` (unambiguous)
- [ ] Register `source_bind` (unambiguous)
- [ ] Register `source_compression_algorithms` (unambiguous)
- [ ] Register `source_connect_retry` (unambiguous)
- [ ] Register `source_connection_auto_failover` (unambiguous)
- [ ] Register `source_delay` (unambiguous)
- [ ] Register `source_heartbeat_period` (unambiguous)
- [ ] Register `source_host` (unambiguous)
- [ ] Register `source_log_file` (unambiguous)
- [ ] Register `source_log_pos` (unambiguous)
- [ ] Register `source_password` (unambiguous)
- [ ] Register `source_port` (unambiguous)
- [ ] Register `source_public_key_path` (unambiguous)
- [ ] Register `source_retry_count` (unambiguous)
- [ ] Register `source_ssl` (unambiguous)
- [ ] Register `source_ssl_ca` (unambiguous)
- [ ] Register `source_ssl_capath` (unambiguous)
- [ ] Register `source_ssl_cert` (unambiguous)
- [ ] Register `source_ssl_cipher` (unambiguous)
- [ ] Register `source_ssl_crl` (unambiguous)
- [ ] Register `source_ssl_crlpath` (unambiguous)
- [ ] Register `source_ssl_key` (unambiguous)
- [ ] Register `source_ssl_verify_server_cert` (unambiguous)
- [ ] Register `source_tls_ciphersuites` (unambiguous)
- [ ] Register `source_tls_version` (unambiguous)
- [ ] Register `source_user` (unambiguous)
- [ ] Register `source_zstd_compression_level` (unambiguous)
- [ ] All newly registered keywords lex correctly and tests pass

### 2.7 Register Missing Keywords — Batch S-Z

- [ ] Register `sql_after_gtids` (unambiguous)
- [ ] Register `sql_after_mts_gaps` (unambiguous)
- [ ] Register `sql_before_gtids` (unambiguous)
- [ ] Register `sql_thread` (unambiguous)
- [ ] Register `sql_tsi_day` (unambiguous)
- [ ] Register `sql_tsi_hour` (unambiguous)
- [ ] Register `sql_tsi_minute` (unambiguous)
- [ ] Register `sql_tsi_month` (unambiguous)
- [ ] Register `sql_tsi_quarter` (unambiguous)
- [ ] Register `sql_tsi_second` (unambiguous)
- [ ] Register `sql_tsi_week` (unambiguous)
- [ ] Register `sql_tsi_year` (unambiguous)
- [ ] Register `st_collect` (unambiguous)
- [ ] Register `stats_auto_recalc` (unambiguous)
- [ ] Register `stats_persistent` (unambiguous)
- [ ] Register `stats_sample_pages` (unambiguous)
- [ ] Register `std` (reserved)
- [ ] Register `stddev` (reserved)
- [ ] Register `stddev_pop` (reserved)
- [ ] Register `stddev_samp` (reserved)
- [ ] Register `strict_load` (unambiguous)
- [ ] Register `string` (unambiguous)
- [ ] Register `subclass_origin` (unambiguous)
- [ ] Register `subdate` (unambiguous)
- [ ] Register `super` (ambiguous_3)
- [ ] Register `substr` (reserved)
- [ ] Register `swaps` (unambiguous)
- [ ] Register `switches` (unambiguous)
- [ ] Register `sysdate` (reserved)
- [ ] Register `system_user` (unambiguous)
- [ ] Register `table_checksum` (unambiguous)
- [ ] Register `table_name` (unambiguous)
- [ ] Register `tablesample` (reserved)
- [ ] Register `thread_priority` (unambiguous)
- [ ] Register `ties` (unambiguous)
- [ ] Register `timestampadd` (unambiguous)
- [ ] Register `timestampdiff` (unambiguous)
- [ ] Register `types` (unambiguous)
- [ ] Register `undo_buffer_size` (unambiguous)
- [ ] Register `undofile` (unambiguous)
- [ ] Register `unicode` (ambiguous_2)
- [ ] Register `unregister` (unambiguous)
- [ ] Register `uri` (unambiguous)
- [ ] Register `url` (unambiguous)
- [ ] Register `use_frm` (unambiguous)
- [ ] Register `user_resources` (unambiguous)
- [ ] Register `validate` (unambiguous)
- [ ] Register `var_pop` (reserved)
- [ ] Register `var_samp` (reserved)
- [ ] Register `varcharacter` (reserved)
- [ ] Register `variance` (reserved)
- [ ] Register `vcpu` (unambiguous)
- [ ] Register `vector` (unambiguous)
- [ ] Register `verify_key_constraints` (unambiguous)
- [ ] Register `week` (unambiguous)
- [ ] Register `weight_string` (unambiguous)
- [ ] Register `zone` (unambiguous)
- [ ] All newly registered keywords lex correctly and tests pass
- [ ] TestKeywordCompleteness passes — 0 missing keywords

---

## Phase 3: Classification Fixes

Fix 132 misclassified keywords. After this phase, TestKeywordClassification must pass.

### 3.1 Add Missing Classification — Reserved Keywords

127 keywords have no entry in keywordCategories (defaulting to unambiguous). Add them.

- [ ] Classify `analyze` as reserved
- [ ] Classify `before` as reserved
- [ ] Classify `bigint` as reserved
- [ ] Classify `blob` as reserved
- [ ] Classify `both` as reserved
- [ ] Classify `call` as reserved
- [ ] Classify `cascade` as reserved
- [ ] Classify `char` as reserved
- [ ] Classify `character` as reserved
- [ ] Classify `collate` as reserved
- [ ] Classify `condition` as reserved
- [ ] Classify `continue` as reserved
- [ ] Classify `count` as reserved
- [ ] Classify `cursor` as reserved
- [ ] Classify `databases` as reserved
- [ ] Classify `dec` as reserved
- [ ] Classify `decimal` as reserved
- [ ] Classify `declare` as reserved
- [ ] Classify `delayed` as reserved
- [ ] Classify `describe` as reserved
- [ ] Classify `deterministic` as reserved
- [ ] Classify `distinctrow` as reserved
- [ ] Classify `double` as reserved
- [ ] Classify `each` as reserved
- [ ] Classify `elseif` as reserved
- [ ] Classify `empty` as reserved
- [ ] Classify `enclosed` as reserved
- [ ] Classify `escaped` as reserved
- [ ] Classify `exit` as reserved
- [ ] Classify `extract` as reserved
- [ ] Classify `fetch` as reserved
- [ ] Classify `float` as reserved
- [ ] Classify `generated` as reserved
- [ ] Classify `get` as reserved
- [ ] Classify `group_concat` as reserved
- [ ] Classify `groups` as reserved
- [ ] Classify `high_priority` as reserved
- [ ] Classify `infile` as reserved
- [ ] Classify `inout` as reserved
- [ ] Classify `int` as reserved
- [ ] Classify `integer` as reserved
- [ ] Classify `iterate` as reserved
- [ ] Classify `json_table` as reserved
- [ ] Classify `keys` as reserved
- [ ] Classify `lateral` as reserved
- [ ] Classify `leading` as reserved
- [ ] Classify `leave` as reserved
- [ ] Classify `linear` as reserved
- [ ] Classify `lines` as reserved
- [ ] Classify `localtime` as reserved
- [ ] Classify `localtimestamp` as reserved
- [ ] Classify `lock` as reserved
- [ ] Classify `longblob` as reserved
- [ ] Classify `longtext` as reserved
- [ ] Classify `loop` as reserved
- [ ] Classify `low_priority` as reserved
- [ ] Classify `max` as reserved
- [ ] Classify `mediumblob` as reserved
- [ ] Classify `mediumint` as reserved
- [ ] Classify `mediumtext` as reserved
- [ ] Classify `min` as reserved
- [ ] Classify `modifies` as reserved
- [ ] Classify `numeric` as reserved
- [ ] Classify `optimize` as reserved
- [ ] Classify `optionally` as reserved
- [ ] Classify `out` as reserved
- [ ] Classify `position` as reserved
- [ ] Classify `purge` as reserved
- [ ] Classify `read` as reserved
- [ ] Classify `reads` as reserved
- [ ] Classify `real` as reserved
- [ ] Classify `recursive` as reserved
- [ ] Classify `release` as reserved
- [ ] Classify `repeat` as reserved
- [ ] Classify `require` as reserved
- [ ] Classify `resignal` as reserved
- [ ] Classify `restrict` as reserved
- [ ] Classify `return` as reserved
- [ ] Classify `rlike` as reserved
- [ ] Classify `schema` as reserved
- [ ] Classify `separator` as reserved
- [ ] Classify `show` as reserved
- [ ] Classify `signal` as reserved
- [ ] Classify `smallint` as reserved
- [ ] Classify `sql` as reserved
- [ ] Classify `sql_big_result` as reserved
- [ ] Classify `sql_calc_found_rows` as reserved
- [ ] Classify `sql_small_result` as reserved
- [ ] Classify `ssl` as reserved
- [ ] Classify `starting` as reserved
- [ ] Classify `stored` as reserved
- [ ] Classify `straight_join` as reserved
- [ ] Classify `substring` as reserved
- [ ] Classify `sum` as reserved
- [ ] Classify `terminated` as reserved
- [ ] Classify `tinyblob` as reserved
- [ ] Classify `tinyint` as reserved
- [ ] Classify `tinytext` as reserved
- [ ] Classify `trailing` as reserved
- [ ] Classify `trim` as reserved
- [ ] Classify `undo` as reserved
- [ ] Classify `unlock` as reserved
- [ ] Classify `unsigned` as reserved
- [ ] Classify `varbinary` as reserved
- [ ] Classify `varchar` as reserved
- [ ] Classify `virtual` as reserved
- [ ] Classify `write` as reserved
- [ ] Classify `zerofill` as reserved
- [ ] All existing parser tests still pass

### 3.2 Add Missing Classification — Ambiguous Keywords

- [ ] Classify `binlog` as ambiguous_2
- [ ] Classify `cache` as ambiguous_2
- [ ] Classify `charset` as ambiguous_2
- [ ] Classify `checksum` as ambiguous_2
- [ ] Classify `clone` as ambiguous_2
- [ ] Classify `comment` as ambiguous_2
- [ ] Classify `deallocate` as ambiguous_2
- [ ] Classify `handler` as ambiguous_2
- [ ] Classify `help` as ambiguous_2
- [ ] Classify `import` as ambiguous_2
- [ ] Classify `install` as ambiguous_2
- [ ] Classify `language` as ambiguous_2
- [ ] Classify `no` as ambiguous_2
- [ ] Classify `uninstall` as ambiguous_2
- [ ] Classify `restart` as ambiguous_1
- [ ] Classify `shutdown` as ambiguous_1
- [ ] Classify `none` as ambiguous_3
- [ ] Classify `resource` as ambiguous_3
- [ ] Classify `persist` as ambiguous_4
- [ ] All existing parser tests still pass

### 3.3 Fix Wrong Classifications

- [ ] Reclassify `modify` from reserved → unambiguous
- [ ] Reclassify `option` from unambiguous → reserved
- [ ] Reclassify `precision` from unambiguous → reserved
- [ ] Reclassify `schemas` from unambiguous → reserved
- [ ] Reclassify `system` from unambiguous → reserved
- [ ] TestKeywordClassification passes — 0 misclassified keywords
- [ ] All existing parser tests still pass (may need fixes for newly reserved words)

---

## Phase 4: eqFold Final Cleanup

Migrate the 5 remaining eqFold violations, now unblocked by Phase 1 lexer fixes.

### 4.1 Migrate @@Variable eqFold Patterns

After Phase 1.2 (@@variable tokenization), the parser receives separate tokens for scope keywords.

- [ ] `global` eqFold in name.go → parser uses kwGLOBAL token from lexer
- [ ] `session` eqFold in name.go → parser uses kwSESSION token from lexer
- [ ] `local` eqFold in name.go → parser uses kwLOCAL token from lexer

### 4.2 Migrate Remaining eqFold Patterns

- [ ] `proxy` eqFold in grant.go → restructure privilege parsing to use kwPROXY token
- [ ] `current_user` eqFold in stmt.go → restructure definer parsing to use kwCURRENT_USER token

### 4.3 Final Golden Test Verification

- [ ] TestKeywordCompleteness passes — 0 missing
- [ ] TestKeywordClassification passes — 0 misclassified
- [ ] TestNoEqFoldForRegisteredKeywords passes — 0 violations
- [ ] Full parser test suite passes — 0 regressions
- [ ] `go vet ./mysql/parser/...` passes

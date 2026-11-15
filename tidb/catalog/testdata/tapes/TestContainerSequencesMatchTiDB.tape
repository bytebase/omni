-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerSequencesMatchTiDB default#1
> exec
| DROP DATABASE IF EXISTS omni_sequences
< ok 0 0
> exec
| CREATE DATABASE omni_sequences
< ok 0 0
> exec
| USE omni_sequences
< ok 0 0
> exec
| CREATE SEQUENCE s_default
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
> exec
| CREATE SEQUENCE s_opts START WITH 10 INCREMENT BY 2 MINVALUE 5 MAXVALUE 1000 CACHE 20 CYCLE COMMENT 'it''s'
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_desc INCREMENT -1
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_desc_start INCREMENT -2 START -10
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_start
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_start
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_desc_max INCREMENT -1 MAXVALUE -5
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_min MINVALUE 5
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_nocache NOCACHE
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_no START WITH 3 INCREMENT BY 4 NO MINVALUE NO MAXVALUE NO CYCLE
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_engine ENGINE=MyISAM
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_all_three START 1 MINVALUE -10 MAXVALUE 10 INCREMENT -1
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_conflict MINVALUE 10 MAXVALUE 5
< error
| 4136 HY000 Sequence 'omni_sequences.s_conflict' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_start0 START 0
< error
| 4136 HY000 Sequence 'omni_sequences.s_start0' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_incr0 INCREMENT 0
< error
| 4136 HY000 Sequence 'omni_sequences.s_incr0' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_cache0 CACHE 0
< error
| 4136 HY000 Sequence 'omni_sequences.s_cache0' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_max MAXVALUE 9223372036854775807
< error
| 4136 HY000 Sequence 'omni_sequences.s_max' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_desc_start_above INCREMENT -1 START 5
< error
| 4136 HY000 Sequence 'omni_sequences.s_desc_start_above' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE s_default
< error
| 1050 42S01 Table 'omni_sequences.s_default' already exists
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE IF NOT EXISTS s_default
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE TABLE s_default (a int)
< error
| 1050 42S01 Table 'omni_sequences.s_default' already exists
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE VIEW s_default AS SELECT 1
< error
| 1050 42S01 Table 'omni_sequences.s_default' already exists
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE t
< error
| 1050 42S01 Table 'omni_sequences.t' already exists
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| CREATE SEQUENCE IF NOT EXISTS t
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_desc MINVALUE -100
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -100 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_desc CACHE 5 INCREMENT 1
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -100 maxvalue -1 increment by 1 cache 5 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_desc MAXVALUE 1000
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with -1 minvalue -100 maxvalue 1000 increment by 1 cache 5 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_desc START 5
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 cache 5 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_desc NOCACHE
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_desc CACHE 0
< error
| 4136 HY000 Sequence 'omni_sequences.s_desc' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_min MINVALUE 10
< error
| 4136 HY000 Sequence 'omni_sequences.s_min' values are conflicting
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_opts RESTART
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_opts RESTART WITH 50
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 cycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_opts NOCYCLE
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE s_missing INCREMENT 2
< error
| 1146 42S02 Table 'omni_sequences.s_missing' doesn't exist
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| ALTER SEQUENCE IF EXISTS s_missing INCREMENT 2
< error
| 1146 42S02 Table 'omni_sequences.s_missing' doesn't exist
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| DROP SEQUENCE s_missing
< error
| 4139 HY000 Unknown SEQUENCE: 'omni_sequences.s_missing'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| DROP SEQUENCE IF EXISTS s_missing
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_no`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_no
: bytes
| CREATE SEQUENCE `s_no` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 4 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_no
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| DROP SEQUENCE s_no, s_missing
< error
| 4139 HY000 Unknown SEQUENCE: 'omni_sequences.s_missing'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_engine`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_engine
: bytes
| CREATE SEQUENCE `s_engine` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_engine
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| DROP SEQUENCE s_engine
< ok 0 0
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> exec
| DROP TABLE s_nocache
< error
| 1051 42S02 Unknown table 'omni_sequences.s_nocache'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SELECT SEQUENCE_NAME, START, MIN_VALUE, MAX_VALUE, INCREMENT, CYCLE, CACHE_VALUE, IFNULL(COMMENT, '')
| 		FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< column "START" "BIGINT"
< column "MIN_VALUE" "BIGINT"
< column "MAX_VALUE" "BIGINT"
< column "INCREMENT" "BIGINT"
< column "CYCLE" "TINYINT"
< column "CACHE_VALUE" "BIGINT"
< column "IFNULL(COMMENT, '')" "VARCHAR"
< row
: bytes
| s_all_three
: int64 1
: int64 -10
: int64 10
: int64 -1
: int64 0
: int64 1000
: bytes
|
< row
: bytes
| s_default
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
< row
: bytes
| s_desc
: int64 5
: int64 -100
: int64 1000
: int64 1
: int64 0
: int64 0
: bytes
|
< row
: bytes
| s_desc_max
: int64 -5
: int64 -9223372036854775807
: int64 -5
: int64 -1
: int64 0
: int64 1000
: bytes
|
< row
: bytes
| s_desc_start
: int64 -10
: int64 -9223372036854775807
: int64 -1
: int64 -2
: int64 0
: int64 1000
: bytes
|
< row
: bytes
| s_min
: int64 5
: int64 5
: int64 9223372036854775806
: int64 1
: int64 0
: int64 1000
: bytes
|
< row
: bytes
| s_nocache
: int64 1
: int64 1
: int64 9223372036854775806
: int64 1
: int64 0
: int64 0
: bytes
|
< row
: bytes
| s_opts
: int64 10
: int64 5
: int64 1000
: int64 2
: int64 0
: int64 20
: bytes
| it's
> query
| SELECT SEQUENCE_NAME FROM information_schema.SEQUENCES WHERE SEQUENCE_SCHEMA = ? ORDER BY SEQUENCE_NAME
: string
| omni_sequences
< result
< column "SEQUENCE_NAME" "VARCHAR"
< row
: bytes
| s_all_three
< row
: bytes
| s_default
< row
: bytes
| s_desc
< row
: bytes
| s_desc_max
< row
: bytes
| s_desc_start
< row
: bytes
| s_min
< row
: bytes
| s_nocache
< row
: bytes
| s_opts
> query
| SHOW CREATE SEQUENCE `s_all_three`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_all_three
: bytes
| CREATE SEQUENCE `s_all_three` start with 1 minvalue -10 maxvalue 10 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_default`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_default
: bytes
| CREATE SEQUENCE `s_default` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc
: bytes
| CREATE SEQUENCE `s_desc` start with 5 minvalue -100 maxvalue 1000 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_max`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_max
: bytes
| CREATE SEQUENCE `s_desc_max` start with -5 minvalue -9223372036854775807 maxvalue -5 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_desc_start`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_desc_start
: bytes
| CREATE SEQUENCE `s_desc_start` start with -10 minvalue -9223372036854775807 maxvalue -1 increment by -2 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_min`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_min
: bytes
| CREATE SEQUENCE `s_min` start with 5 minvalue 5 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_nocache`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_nocache
: bytes
| CREATE SEQUENCE `s_nocache` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `s_opts`
< result
< column "Sequence" "VARCHAR"
< column "Create Sequence" "VARCHAR"
< row
: bytes
| s_opts
: bytes
| CREATE SEQUENCE `s_opts` start with 10 minvalue 5 maxvalue 1000 increment by 2 cache 20 nocycle ENGINE=InnoDB COMMENT='it''s'
> exec
| USE test
< ok 0 0

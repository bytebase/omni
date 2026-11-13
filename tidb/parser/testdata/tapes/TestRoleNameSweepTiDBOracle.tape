-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestRoleNameSweepTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| create role absent
< ok 0 0
> exec ^
: string
| create role accessible
< ok 0 0
> exec ^
: string
| create role account
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role action
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role active
< ok 0 0
> exec ^
: string
| create role add
< error
| 1064 42000 …line 1 column 15 near "add" 
> exec ^
: string
| create role adddate
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role admin
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role after
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role against
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role aggregate
< ok 0 0
> exec ^
: string
| create role algorithm
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role all
< error
| 1064 42000 …line 1 column 15 near "all" 
> exec ^
: string
| create role allow_missing_files
< ok 0 0
> exec ^
: string
| create role alter
< error
| 1064 42000 …line 1 column 17 near "alter" 
> exec ^
: string
| create role always
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role analyze
< error
| 1064 42000 …line 1 column 19 near "analyze" 
> exec ^
: string
| create role and
< error
| 1064 42000 …line 1 column 15 near "and" 
> exec ^
: string
| create role any
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role array
< error
| 1064 42000 …line 1 column 17 near "array" 
> exec ^
: string
| create role as
< error
| 1064 42000 …line 1 column 14 near "as" 
> exec ^
: string
| create role asc
< error
| 1064 42000 …line 1 column 15 near "asc" 
> exec ^
: string
| create role ascii
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role asensitive
< ok 0 0
> exec ^
: string
| create role assign_gtids_to_anonymous_transactions
< ok 0 0
> exec ^
: string
| create role at
< ok 0 0
> exec ^
: string
| create role attribute
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role authentication
< ok 0 0
> exec ^
: string
| create role auto
< ok 0 0
> exec ^
: string
| create role auto_id_cache
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role auto_increment
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role auto_random
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role auto_random_base
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role auto_refresh
< ok 0 0
> exec ^
: string
| create role auto_refresh_source
< ok 0 0
> exec ^
: string
| create role autoextend_size
< ok 0 0
> exec ^
: string
| create role avg
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role avg_row_length
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role background
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role backup
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role backups
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role batch
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role before
< ok 0 0
> exec ^
: string
| create role begin
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role bernoulli
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role between
< error
| 1064 42000 …line 1 column 19 near "between" 
> exec ^
: string
| create role bigint
< error
| 1064 42000 …line 1 column 18 near "bigint" 
> exec ^
: string
| create role binary
< error
| 1064 42000 …line 1 column 18 near "binary" 
> exec ^
: string
| create role binding
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role binding_cache
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role bindings
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role binlog
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role bit
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role bit_and
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role bit_or
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role bit_xor
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role blob
< error
| 1064 42000 …line 1 column 16 near "blob" 
> exec ^
: string
| create role block
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role bool
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role boolean
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role both
< error
| 1064 42000 …line 1 column 16 near "both" 
> exec ^
: string
| create role btree
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role buckets
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role builtins
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role bulk
< ok 0 0
> exec ^
: string
| create role burstable
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role by
< error
| 1064 42000 …line 1 column 14 near "by" 
> exec ^
: string
| create role byte
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role cache
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role calibrate
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role call
< error
| 1064 42000 …line 1 column 16 near "call" 
> exec ^
: string
| create role cascade
< error
| 1064 42000 …line 1 column 19 near "cascade" 
> exec ^
: string
| create role cascaded
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role case
< error
| 1064 42000 …line 1 column 16 near "case" 
> exec ^
: string
| create role cast
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role catalog_name
< ok 0 0
> exec ^
: string
| create role chain
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role challenge_response
< ok 0 0
> exec ^
: string
| create role change
< error
| 1064 42000 …line 1 column 18 near "change" 
> exec ^
: string
| create role changed
< ok 0 0
> exec ^
: string
| create role channel
< ok 0 0
> exec ^
: string
| create role char
< error
| 1064 42000 …line 1 column 16 near "char" 
> exec ^
: string
| create role character
< error
| 1064 42000 …line 1 column 21 near "character" 
> exec ^
: string
| create role charset
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role check
< error
| 1064 42000 …line 1 column 17 near "check" 
> exec ^
: string
| create role checksum
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role cipher
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role class_origin
< ok 0 0
> exec ^
: string
| create role client
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role clone
< ok 0 0
> exec ^
: string
| create role close
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role clustered
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role coalesce
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role code
< ok 0 0
> exec ^
: string
| create role collate
< error
| 1064 42000 …line 1 column 19 near "collate" 
> exec ^
: string
| create role collation
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role column
< error
| 1064 42000 …line 1 column 18 near "column" 
> exec ^
: string
| create role column_format
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role column_name
< ok 0 0
> exec ^
: string
| create role column_stats_usage
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| create role columns
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role comment
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role commit
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role committed
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role compact
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role completion
< ok 0 0
> exec ^
: string
| create role component
< ok 0 0
> exec ^
: string
| create role compressed
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role compression
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role concat
< ok 0 0
> exec ^
: string
| create role concurrent
< ok 0 0
> exec ^
: string
| create role condition
< ok 0 0
> exec ^
: string
| create role config
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role connection
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role consistent
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role constraint
< error
| 1064 42000 …line 1 column 22 near "constraint" 
> exec ^
: string
| create role constraint_catalog
< ok 0 0
> exec ^
: string
| create role constraint_name
< ok 0 0
> exec ^
: string
| create role constraint_schema
< ok 0 0
> exec ^
: string
| create role constraints
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role contains
< ok 0 0
> exec ^
: string
| create role context
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role continue
< error
| 1064 42000 …line 1 column 20 near "continue" 
> exec ^
: string
| create role convert
< error
| 1064 42000 …line 1 column 19 near "convert" 
> exec ^
: string
| create role cooldown
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role count
< ok 0 0
> exec ^
: string
| create role cpu
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role create
< error
| 1064 42000 …line 1 column 18 near "create" 
> exec ^
: string
| create role cross
< error
| 1064 42000 …line 1 column 17 near "cross" 
> exec ^
: string
| create role cube
< ok 0 0
> exec ^
: string
| create role cume_dist
< error
| 1064 42000 …line 1 column 21 near "cume_dist" 
> exec ^
: string
| create role curdate
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role current
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role current_date
< error
| 1064 42000 …line 1 column 24 near "current_date" 
> exec ^
: string
| create role current_time
< error
| 1064 42000 …line 1 column 24 near "current_time" 
> exec ^
: string
| create role current_timestamp
< error
| 1064 42000 …line 1 column 29 near "current_timestamp" 
> exec ^
: string
| create role current_user
< error
| 1064 42000 …line 1 column 24 near "current_user" 
> exec ^
: string
| create role cursor
< error
| 1064 42000 …line 1 column 18 near "cursor" 
> exec ^
: string
| create role cursor_name
< ok 0 0
> exec ^
: string
| create role curtime
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role cycle
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role data
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role database
< error
| 1064 42000 …line 1 column 20 near "database" 
> exec ^
: string
| create role databases
< error
| 1064 42000 …line 1 column 21 near "databases" 
> exec ^
: string
| create role datafile
< ok 0 0
> exec ^
: string
| create role date
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role date_add
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role date_sub
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role datetime
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role day
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role day_hour
< error
| 1064 42000 …line 1 column 20 near "day_hour" 
> exec ^
: string
| create role day_microsecond
< error
| 1064 42000 …line 1 column 27 near "day_microsecond" 
> exec ^
: string
| create role day_minute
< error
| 1064 42000 …line 1 column 22 near "day_minute" 
> exec ^
: string
| create role day_second
< error
| 1064 42000 …line 1 column 22 near "day_second" 
> exec ^
: string
| create role deallocate
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role dec
< error
| 1064 42000 …line 1 column 15 near "dec" 
> exec ^
: string
| create role decimal
< error
| 1064 42000 …line 1 column 19 near "decimal" 
> exec ^
: string
| create role declare
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role default
< error
| 1064 42000 …line 1 column 19 near "default" 
> exec ^
: string
| create role default_auth
< ok 0 0
> exec ^
: string
| create role defined
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role definer
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role definition
< ok 0 0
> exec ^
: string
| create role delay_key_write
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| create role delayed
< error
| 1064 42000 …line 1 column 19 near "delayed" 
> exec ^
: string
| create role delete
< error
| 1064 42000 …line 1 column 18 near "delete" 
> exec ^
: string
| create role dense_rank
< error
| 1064 42000 …line 1 column 22 near "dense_rank" 
> exec ^
: string
| create role desc
< error
| 1064 42000 …line 1 column 16 near "desc" 
> exec ^
: string
| create role describe
< error
| 1064 42000 …line 1 column 20 near "describe" 
> exec ^
: string
| create role description
< ok 0 0
> exec ^
: string
| create role deterministic
< ok 0 0
> exec ^
: string
| create role diagnostics
< ok 0 0
> exec ^
: string
| create role digest
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role directory
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role disable
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role disabled
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role discard
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role disk
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role distinct
< error
| 1064 42000 …line 1 column 20 near "distinct" 
> exec ^
: string
| create role distinctrow
< error
| 1064 42000 …line 1 column 23 near "distinctrow" 
> exec ^
: string
| create role distribute
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role distribution
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role distributions
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role div
< error
| 1064 42000 …line 1 column 15 near "div" 
> exec ^
: string
| create role do
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role double
< error
| 1064 42000 …line 1 column 18 near "double" 
> exec ^
: string
| create role drop
< error
| 1064 42000 …line 1 column 16 near "drop" 
> exec ^
: string
| create role dry
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role dryrun
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role dual
< error
| 1064 42000 …line 1 column 16 near "dual" 
> exec ^
: string
| create role duality
< ok 0 0
> exec ^
: string
| create role dumpfile
< ok 0 0
> exec ^
: string
| create role duplicate
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role duration
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role dynamic
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role each
< ok 0 0
> exec ^
: string
| create role else
< error
| 1064 42000 …line 1 column 16 near "else" 
> exec ^
: string
| create role elseif
< error
| 1064 42000 …line 1 column 18 near "elseif" 
> exec ^
: string
| create role empty
< ok 0 0
> exec ^
: string
| create role enable
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role enabled
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role enclosed
< error
| 1064 42000 …line 1 column 20 near "enclosed" 
> exec ^
: string
| create role encryption
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role end
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role end_time
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role ends
< ok 0 0
> exec ^
: string
| create role enforced
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role engine
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role engine_attribute
< ok 0 0
> exec ^
: string
| create role engines
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role enum
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role error
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role errors
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role escape
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role escaped
< error
| 1064 42000 …line 1 column 19 near "escaped" 
> exec ^
: string
| create role event
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role events
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role every
< ok 0 0
> exec ^
: string
| create role exact
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role except
< error
| 1064 42000 …line 1 column 18 near "except" 
> exec ^
: string
| create role exchange
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role exclude
< ok 0 0
> exec ^
: string
| create role exec_elapsed
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role execute
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role exists
< error
| 1064 42000 …line 1 column 18 near "exists" 
> exec ^
: string
| create role exit
< error
| 1064 42000 …line 1 column 16 near "exit" 
> exec ^
: string
| create role expansion
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role expire
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role explain
< error
| 1064 42000 …line 1 column 19 near "explain" 
> exec ^
: string
| create role export
< ok 0 0
> exec ^
: string
| create role extended
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role extent_size
< ok 0 0
> exec ^
: string
| create role external
< ok 0 0
> exec ^
: string
| create role external_format
< ok 0 0
> exec ^
: string
| create role extract
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role factor
< ok 0 0
> exec ^
: string
| create role failed_login_attempts
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| create role false
< error
| 1064 42000 …line 1 column 17 near "false" 
> exec ^
: string
| create role fast
< ok 0 0
> exec ^
: string
| create role faults
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role fetch
< error
| 1064 42000 …line 1 column 17 near "fetch" 
> exec ^
: string
| create role fields
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role file
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role file_block_size
< ok 0 0
> exec ^
: string
| create role file_format
< ok 0 0
> exec ^
: string
| create role file_name
< ok 0 0
> exec ^
: string
| create role file_pattern
< ok 0 0
> exec ^
: string
| create role file_prefix
< ok 0 0
> exec ^
: string
| create role files
< ok 0 0
> exec ^
: string
| create role filter
< ok 0 0
> exec ^
: string
| create role finish
< ok 0 0
> exec ^
: string
| create role first
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role first_value
< error
| 1064 42000 …line 1 column 23 near "first_value" 
> exec ^
: string
| create role fixed
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role flashback
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role float
< error
| 1064 42000 …line 1 column 17 near "float" 
> exec ^
: string
| create role float4
< error
| 1064 42000 …line 1 column 18 near "float4" 
> exec ^
: string
| create role float8
< error
| 1064 42000 …line 1 column 18 near "float8" 
> exec ^
: string
| create role flush
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role follower_constraints
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| create role followers
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role following
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role follows
< ok 0 0
> exec ^
: string
| create role for
< error
| 1064 42000 …line 1 column 15 near "for" 
> exec ^
: string
| create role force
< error
| 1064 42000 …line 1 column 17 near "force" 
> exec ^
: string
| create role foreign
< error
| 1064 42000 …line 1 column 19 near "foreign" 
> exec ^
: string
| create role format
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role found
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role from
< error
| 1064 42000 …line 1 column 16 near "from" 
> exec ^
: string
| create role full
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role fulltext
< error
| 1064 42000 …line 1 column 20 near "fulltext" 
> exec ^
: string
| create role function
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role general
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role generate
< ok 0 0
> exec ^
: string
| create role generated
< error
| 1064 42000 …line 1 column 21 near "generated" 
> exec ^
: string
| create role geomcollection
< ok 0 0
> exec ^
: string
| create role geometry
< ok 0 0
> exec ^
: string
| create role geometrycollection
< ok 0 0
> exec ^
: string
| create role get
< ok 0 0
> exec ^
: string
| create role get_format
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role get_source_public_key
< ok 0 0
> exec ^
: string
| create role global
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role grant
< error
| 1064 42000 …line 1 column 17 near "grant" 
> exec ^
: string
| create role grants
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role greatest
< ok 0 0
> exec ^
: string
| create role group
< error
| 1064 42000 …line 1 column 17 near "group" 
> exec ^
: string
| create role group_concat
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role group_replication
< ok 0 0
> exec ^
: string
| create role grouping
< ok 0 0
> exec ^
: string
| create role groups
< error
| 1064 42000 …line 1 column 18 near "groups" 
> exec ^
: string
| create role gtid_only
< ok 0 0
> exec ^
: string
| create role gtids
< ok 0 0
> exec ^
: string
| create role guided
< ok 0 0
> exec ^
: string
| create role handler
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role hash
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role having
< error
| 1064 42000 …line 1 column 18 near "having" 
> exec ^
: string
| create role header
< ok 0 0
> exec ^
: string
| create role help
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role high
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role high_priority
< error
| 1064 42000 …line 1 column 25 near "high_priority" 
> exec ^
: string
| create role histogram
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role histograms_in_flight
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| create role history
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role host
< ok 0 0
> exec ^
: string
| create role hosts
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role hour
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role hour_microsecond
< error
| 1064 42000 …line 1 column 28 near "hour_microsecond" 
> exec ^
: string
| create role hour_minute
< error
| 1064 42000 …line 1 column 23 near "hour_minute" 
> exec ^
: string
| create role hour_second
< error
| 1064 42000 …line 1 column 23 near "hour_second" 
> exec ^
: string
| create role identified
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role if
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role ignore
< error
| 1064 42000 …line 1 column 18 near "ignore" 
> exec ^
: string
| create role ignore_server_ids
< ok 0 0
> exec ^
: string
| create role ilike
< error
| 1064 42000 …line 1 column 17 near "ilike" 
> exec ^
: string
| create role import
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role in
< error
| 1064 42000 …line 1 column 14 near "in" 
> exec ^
: string
| create role inactive
< ok 0 0
> exec ^
: string
| create role increment
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role index
< error
| 1064 42000 …line 1 column 17 near "index" 
> exec ^
: string
| create role indexes
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role infile
< error
| 1064 42000 …line 1 column 18 near "infile" 
> exec ^
: string
| create role initial
< ok 0 0
> exec ^
: string
| create role initial_size
< ok 0 0
> exec ^
: string
| create role initiate
< ok 0 0
> exec ^
: string
| create role inner
< error
| 1064 42000 …line 1 column 17 near "inner" 
> exec ^
: string
| create role inout
< error
| 1064 42000 …line 1 column 17 near "inout" 
> exec ^
: string
| create role insensitive
< ok 0 0
> exec ^
: string
| create role insert
< error
| 1064 42000 …line 1 column 18 near "insert" 
> exec ^
: string
| create role insert_method
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role install
< ok 0 0
> exec ^
: string
| create role instance
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role int
< error
| 1064 42000 …line 1 column 15 near "int" 
> exec ^
: string
| create role int1
< error
| 1064 42000 …line 1 column 16 near "int1" 
> exec ^
: string
| create role int2
< error
| 1064 42000 …line 1 column 16 near "int2" 
> exec ^
: string
| create role int3
< error
| 1064 42000 …line 1 column 16 near "int3" 
> exec ^
: string
| create role int4
< error
| 1064 42000 …line 1 column 16 near "int4" 
> exec ^
: string
| create role int8
< error
| 1064 42000 …line 1 column 16 near "int8" 
> exec ^
: string
| create role integer
< error
| 1064 42000 …line 1 column 19 near "integer" 
> exec ^
: string
| create role intersect
< error
| 1064 42000 …line 1 column 21 near "intersect" 
> exec ^
: string
| create role interval
< error
| 1064 42000 …line 1 column 20 near "interval" 
> exec ^
: string
| create role into
< error
| 1064 42000 …line 1 column 16 near "into" 
> exec ^
: string
| create role invisible
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role invoker
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role io
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role io_after_gtids
< ok 0 0
> exec ^
: string
| create role io_before_gtids
< ok 0 0
> exec ^
: string
| create role io_thread
< ok 0 0
> exec ^
: string
| create role ipc
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role is
< error
| 1064 42000 …line 1 column 14 near "is" 
> exec ^
: string
| create role isolation
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role issuer
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role iterate
< error
| 1064 42000 …line 1 column 19 near "iterate" 
> exec ^
: string
| create role job
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role jobs
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role join
< error
| 1064 42000 …line 1 column 16 near "join" 
> exec ^
: string
| create role json
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role json_arrayagg
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role json_duality_object
< ok 0 0
> exec ^
: string
| create role json_objectagg
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role json_table
< ok 0 0
> exec ^
: string
| create role json_value
< ok 0 0
> exec ^
: string
| create role key
< error
| 1064 42000 …line 1 column 15 near "key" 
> exec ^
: string
| create role key_block_size
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role keyring
< ok 0 0
> exec ^
: string
| create role keys
< error
| 1064 42000 …line 1 column 16 near "keys" 
> exec ^
: string
| create role kill
< error
| 1064 42000 …line 1 column 16 near "kill" 
> exec ^
: string
| create role labels
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role lag
< error
| 1064 42000 …line 1 column 15 near "lag" 
> exec ^
: string
| create role language
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role last
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role last_value
< error
| 1064 42000 …line 1 column 22 near "last_value" 
> exec ^
: string
| create role lateral
< ok 0 0
> exec ^
: string
| create role lead
< error
| 1064 42000 …line 1 column 16 near "lead" 
> exec ^
: string
| create role leader_constraints
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| create role leading
< error
| 1064 42000 …line 1 column 19 near "leading" 
> exec ^
: string
| create role learner_constraints
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| create role learners
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role least
< ok 0 0
> exec ^
: string
| create role leave
< error
| 1064 42000 …line 1 column 17 near "leave" 
> exec ^
: string
| create role leaves
< ok 0 0
> exec ^
: string
| create role left
< error
| 1064 42000 …line 1 column 16 near "left" 
> exec ^
: string
| create role less
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role level
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role library
< ok 0 0
> exec ^
: string
| create role like
< error
| 1064 42000 …line 1 column 16 near "like" 
> exec ^
: string
| create role limit
< error
| 1064 42000 …line 1 column 17 near "limit" 
> exec ^
: string
| create role linear
< error
| 1064 42000 …line 1 column 18 near "linear" 
> exec ^
: string
| create role lines
< error
| 1064 42000 …line 1 column 17 near "lines" 
> exec ^
: string
| create role linestring
< ok 0 0
> exec ^
: string
| create role list
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role load
< error
| 1064 42000 …line 1 column 16 near "load" 
> exec ^
: string
| create role local
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role localtime
< error
| 1064 42000 …line 1 column 21 near "localtime" 
> exec ^
: string
| create role localtimestamp
< error
| 1064 42000 …line 1 column 26 near "localtimestamp" 
> exec ^
: string
| create role location
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role lock
< error
| 1064 42000 …line 1 column 16 near "lock" 
> exec ^
: string
| create role locked
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role locks
< ok 0 0
> exec ^
: string
| create role log
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role logfile
< ok 0 0
> exec ^
: string
| create role logs
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role long
< error
| 1064 42000 …line 1 column 16 near "long" 
> exec ^
: string
| create role longblob
< error
| 1064 42000 …line 1 column 20 near "longblob" 
> exec ^
: string
| create role longtext
< error
| 1064 42000 …line 1 column 20 near "longtext" 
> exec ^
: string
| create role loop
< ok 0 0
> exec ^
: string
| create role low
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role low_priority
< error
| 1064 42000 …line 1 column 24 near "low_priority" 
> exec ^
: string
| create role manual
< ok 0 0
> exec ^
: string
| create role master
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role master_log_file
< ok 0 0
> exec ^
: string
| create role match
< error
| 1064 42000 …line 1 column 17 near "match" 
> exec ^
: string
| create role materialized
< ok 0 0
> exec ^
: string
| create role max
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role max_connections_per_hour
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| create role max_queries_per_hour
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| create role max_rows
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role max_size
< ok 0 0
> exec ^
: string
| create role max_updates_per_hour
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| create role max_user_connections
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| create role maxvalue
< error
| 1064 42000 …line 1 column 20 near "maxvalue" 
> exec ^
: string
| create role medium
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role mediumblob
< error
| 1064 42000 …line 1 column 22 near "mediumblob" 
> exec ^
: string
| create role mediumint
< error
| 1064 42000 …line 1 column 21 near "mediumint" 
> exec ^
: string
| create role mediumtext
< error
| 1064 42000 …line 1 column 22 near "mediumtext" 
> exec ^
: string
| create role member
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role memory
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role merge
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role message_text
< ok 0 0
> exec ^
: string
| create role microsecond
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role mid
< ok 0 0
> exec ^
: string
| create role middleint
< error
| 1064 42000 …line 1 column 21 near "middleint" 
> exec ^
: string
| create role migrate
< ok 0 0
> exec ^
: string
| create role min
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role min_rows
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role minute
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role minute_microsecond
< error
| 1064 42000 …line 1 column 30 near "minute_microsecond" 
> exec ^
: string
| create role minute_second
< error
| 1064 42000 …line 1 column 25 near "minute_second" 
> exec ^
: string
| create role minvalue
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role mod
< error
| 1064 42000 …line 1 column 15 near "mod" 
> exec ^
: string
| create role mode
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role modifies
< ok 0 0
> exec ^
: string
| create role modify
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role month
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role multilinestring
< ok 0 0
> exec ^
: string
| create role multipoint
< ok 0 0
> exec ^
: string
| create role multipolygon
< ok 0 0
> exec ^
: string
| create role mutex
< ok 0 0
> exec ^
: string
| create role mysql_errno
< ok 0 0
> exec ^
: string
| create role name
< ok 0 0
> exec ^
: string
| create role names
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role national
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role natural
< error
| 1064 42000 …line 1 column 19 near "natural" 
> exec ^
: string
| create role nchar
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role ndb
< ok 0 0
> exec ^
: string
| create role ndbcluster
< ok 0 0
> exec ^
: string
| create role nested
< ok 0 0
> exec ^
: string
| create role network_namespace
< ok 0 0
> exec ^
: string
| create role never
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role new
< ok 0 0
> exec ^
: string
| create role next
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role next_row_id
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role no
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role no_wait
< ok 0 0
> exec ^
: string
| create role no_write_to_binlog
< error
| 1064 42000 …line 1 column 30 near "no_write_to_binlog" 
> exec ^
: string
| create role nocache
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role nocycle
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role nodegroup
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role nomaxvalue
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role nominvalue
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role nonclustered
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role none
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role not
< error
| 1064 42000 …line 1 column 15 near "not" 
> exec ^
: string
| create role now
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role nowait
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role nth_value
< error
| 1064 42000 …line 1 column 21 near "nth_value" 
> exec ^
: string
| create role ntile
< error
| 1064 42000 …line 1 column 17 near "ntile" 
> exec ^
: string
| create role null
< error
| 1064 42000 …line 1 column 16 near "null" 
> exec ^
: string
| create role nullif
< ok 0 0
> exec ^
: string
| create role nulls
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role number
< ok 0 0
> exec ^
: string
| create role numeric
< error
| 1064 42000 …line 1 column 19 near "numeric" 
> exec ^
: string
| create role nvarchar
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role of
< error
| 1064 42000 …line 1 column 14 near "of" 
> exec ^
: string
| create role off
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role offset
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role oj
< ok 0 0
> exec ^
: string
| create role old
< ok 0 0
> exec ^
: string
| create role oltp_read_only
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role oltp_read_write
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| create role oltp_write_only
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| create role on
< error
| 1064 42000 …line 1 column 14 near "on" 
> exec ^
: string
| create role one
< ok 0 0
> exec ^
: string
| create role only
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role open
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role optimize
< error
| 1064 42000 …line 1 column 20 near "optimize" 
> exec ^
: string
| create role optimizer_costs
< ok 0 0
> exec ^
: string
| create role option
< error
| 1064 42000 …line 1 column 18 near "option" 
> exec ^
: string
| create role optional
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role optionally
< error
| 1064 42000 …line 1 column 22 near "optionally" 
> exec ^
: string
| create role options
< ok 0 0
> exec ^
: string
| create role or
< error
| 1064 42000 …line 1 column 14 near "or" 
> exec ^
: string
| create role order
< error
| 1064 42000 …line 1 column 17 near "order" 
> exec ^
: string
| create role ordinality
< ok 0 0
> exec ^
: string
| create role organization
< ok 0 0
> exec ^
: string
| create role others
< ok 0 0
> exec ^
: string
| create role out
< error
| 1064 42000 …line 1 column 15 near "out" 
> exec ^
: string
| create role outer
< error
| 1064 42000 …line 1 column 17 near "outer" 
> exec ^
: string
| create role outfile
< error
| 1064 42000 …line 1 column 19 near "outfile" 
> exec ^
: string
| create role over
< error
| 1064 42000 …line 1 column 16 near "over" 
> exec ^
: string
| create role overlay
< ok 0 0
> exec ^
: string
| create role owner
< ok 0 0
> exec ^
: string
| create role pack_keys
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role page
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role parallel
< ok 0 0
> exec ^
: string
| create role parameters
< ok 0 0
> exec ^
: string
| create role parse_tree
< ok 0 0
> exec ^
: string
| create role parser
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role partial
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role partition
< error
| 1064 42000 …line 1 column 21 near "partition" 
> exec ^
: string
| create role partitioning
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role partitions
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role password
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role password_lock_time
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| create role path
< ok 0 0
> exec ^
: string
| create role percent_rank
< error
| 1064 42000 …line 1 column 24 near "percent_rank" 
> exec ^
: string
| create role persist
< ok 0 0
> exec ^
: string
| create role persist_only
< ok 0 0
> exec ^
: string
| create role phase
< ok 0 0
> exec ^
: string
| create role placement
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role placing
< ok 0 0
> exec ^
: string
| create role plan
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role plugin
< ok 0 0
> exec ^
: string
| create role plugin_dir
< ok 0 0
> exec ^
: string
| create role plugins
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role point
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role policy
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role polygon
< ok 0 0
> exec ^
: string
| create role port
< ok 0 0
> exec ^
: string
| create role position
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role pre_split_regions
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role precedes
< ok 0 0
> exec ^
: string
| create role preceding
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role precision
< error
| 1064 42000 …line 1 column 21 near "precision" 
> exec ^
: string
| create role prepare
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role preserve
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role prev
< ok 0 0
> exec ^
: string
| create role primary
< error
| 1064 42000 …line 1 column 19 near "primary" 
> exec ^
: string
| create role primary_region
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role priority
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role privilege_checks_user
< ok 0 0
> exec ^
: string
| create role privileges
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role procedure
< error
| 1064 42000 …line 1 column 21 near "procedure" 
> exec ^
: string
| create role process
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role processed_keys
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role processlist
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role profile
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role profiles
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role proxy
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role purge
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role qualify
< ok 0 0
> exec ^
: string
| create role quarter
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role query
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role query_limit
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role quick
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role random
< ok 0 0
> exec ^
: string
| create role range
< error
| 1064 42000 …line 1 column 17 near "range" 
> exec ^
: string
| create role rank
< error
| 1064 42000 …line 1 column 16 near "rank" 
> exec ^
: string
| create role read
< error
| 1064 42000 …line 1 column 16 near "read" 
> exec ^
: string
| create role read_only
< ok 0 0
> exec ^
: string
| create role read_write
< ok 0 0
> exec ^
: string
| create role reads
< ok 0 0
> exec ^
: string
| create role real
< error
| 1064 42000 …line 1 column 16 near "real" 
> exec ^
: string
| create role rebuild
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role recover
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role recursive
< error
| 1064 42000 …line 1 column 21 near "recursive" 
> exec ^
: string
| create role redo_buffer_size
< ok 0 0
> exec ^
: string
| create role redundant
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role reference
< ok 0 0
> exec ^
: string
| create role references
< error
| 1064 42000 …line 1 column 22 near "references" 
> exec ^
: string
| create role regexp
< error
| 1064 42000 …line 1 column 18 near "regexp" 
> exec ^
: string
| create role region
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role regions
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role registration
< ok 0 0
> exec ^
: string
| create role relational
< ok 0 0
> exec ^
: string
| create role relay
< ok 0 0
> exec ^
: string
| create role relay_log_file
< ok 0 0
> exec ^
: string
| create role relay_log_pos
< ok 0 0
> exec ^
: string
| create role relay_thread
< ok 0 0
> exec ^
: string
| create role relaylog
< ok 0 0
> exec ^
: string
| create role release
< error
| 1064 42000 …line 1 column 19 near "release" 
> exec ^
: string
| create role reload
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role remove
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role rename
< error
| 1064 42000 …line 1 column 18 near "rename" 
> exec ^
: string
| create role reorganize
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role repair
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role repeat
< error
| 1064 42000 …line 1 column 18 near "repeat" 
> exec ^
: string
| create role repeatable
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role replace
< error
| 1064 42000 …line 1 column 19 near "replace" 
> exec ^
: string
| create role replica
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role replicas
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role replicate_do_db
< ok 0 0
> exec ^
: string
| create role replicate_do_table
< ok 0 0
> exec ^
: string
| create role replicate_ignore_db
< ok 0 0
> exec ^
: string
| create role replicate_ignore_table
< ok 0 0
> exec ^
: string
| create role replicate_rewrite_db
< ok 0 0
> exec ^
: string
| create role replicate_wild_do_table
< ok 0 0
> exec ^
: string
| create role replicate_wild_ignore_table
< ok 0 0
> exec ^
: string
| create role replication
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role require
< error
| 1064 42000 …line 1 column 19 near "require" 
> exec ^
: string
| create role require_row_format
< ok 0 0
> exec ^
: string
| create role require_table_primary_key_check
< ok 0 0
> exec ^
: string
| create role reset
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role resignal
< ok 0 0
> exec ^
: string
| create role resource
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role respect
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role restart
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role restore
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role restores
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role restrict
< error
| 1064 42000 …line 1 column 20 near "restrict" 
> exec ^
: string
| create role resume
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role retain
< ok 0 0
> exec ^
: string
| create role return
< ok 0 0
> exec ^
: string
| create role returned_sqlstate
< ok 0 0
> exec ^
: string
| create role returning
< ok 0 0
> exec ^
: string
| create role returns
< ok 0 0
> exec ^
: string
| create role reuse
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role reverse
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role revoke
< error
| 1064 42000 …line 1 column 18 near "revoke" 
> exec ^
: string
| create role right
< error
| 1064 42000 …line 1 column 17 near "right" 
> exec ^
: string
| create role rlike
< error
| 1064 42000 …line 1 column 17 near "rlike" 
> exec ^
: string
| create role role
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role rollback
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role rollup
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role rotate
< ok 0 0
> exec ^
: string
| create role routine
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role row
< error
| 1064 42000 …line 1 column 15 near "row" 
> exec ^
: string
| create role row_count
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role row_format
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role row_number
< error
| 1064 42000 …line 1 column 22 near "row_number" 
> exec ^
: string
| create role rows
< error
| 1064 42000 …line 1 column 16 near "rows" 
> exec ^
: string
| create role rtree
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role ru
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role ru_per_sec
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role rule
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role run
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role s3
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| create role savepoint
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role schedule
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role schema
< error
| 1064 42000 …line 1 column 18 near "schema" 
> exec ^
: string
| create role schema_name
< ok 0 0
> exec ^
: string
| create role schemas
< error
| 1064 42000 …line 1 column 19 near "schemas" 
> exec ^
: string
| create role second
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role second_microsecond
< error
| 1064 42000 …line 1 column 30 near "second_microsecond" 
> exec ^
: string
| create role secondary
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role secondary_engine
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role secondary_engine_attribute
< ok 0 0
> exec ^
: string
| create role secondary_load
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role secondary_unload
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role security
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role select
< error
| 1064 42000 …line 1 column 18 near "select" 
> exec ^
: string
| create role sensitive
< ok 0 0
> exec ^
: string
| create role separator
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role sequence
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role serial
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role serializable
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role server
< ok 0 0
> exec ^
: string
| create role session
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role session_states
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role session_user
< ok 0 0
> exec ^
: string
| create role set
< error
| 1064 42000 …line 1 column 15 near "set" 
> exec ^
: string
| create role sets
< ok 0 0
> exec ^
: string
| create role shard_row_id_bits
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role share
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role show
< error
| 1064 42000 …line 1 column 16 near "show" 
> exec ^
: string
| create role shutdown
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role signal
< ok 0 0
> exec ^
: string
| create role signed
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role similar
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role simple
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role skip
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role slave
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role slow
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role smallint
< error
| 1064 42000 …line 1 column 20 near "smallint" 
> exec ^
: string
| create role snapshot
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role socket
< ok 0 0
> exec ^
: string
| create role some
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role soname
< ok 0 0
> exec ^
: string
| create role sounds
< ok 0 0
> exec ^
: string
| create role source
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role source_auto_position
< ok 0 0
> exec ^
: string
| create role source_bind
< ok 0 0
> exec ^
: string
| create role source_compression_algorithms
< ok 0 0
> exec ^
: string
| create role source_connect_retry
< ok 0 0
> exec ^
: string
| create role source_connection_auto_failover
< ok 0 0
> exec ^
: string
| create role source_delay
< ok 0 0
> exec ^
: string
| create role source_heartbeat_period
< ok 0 0
> exec ^
: string
| create role source_host
< ok 0 0
> exec ^
: string
| create role source_log_file
< ok 0 0
> exec ^
: string
| create role source_log_pos
< ok 0 0
> exec ^
: string
| create role source_password
< ok 0 0
> exec ^
: string
| create role source_port
< ok 0 0
> exec ^
: string
| create role source_public_key_path
< ok 0 0
> exec ^
: string
| create role source_retry_count
< ok 0 0
> exec ^
: string
| create role source_ssl
< ok 0 0
> exec ^
: string
| create role source_ssl_ca
< ok 0 0
> exec ^
: string
| create role source_ssl_capath
< ok 0 0
> exec ^
: string
| create role source_ssl_cert
< ok 0 0
> exec ^
: string
| create role source_ssl_cipher
< ok 0 0
> exec ^
: string
| create role source_ssl_crl
< ok 0 0
> exec ^
: string
| create role source_ssl_crlpath
< ok 0 0
> exec ^
: string
| create role source_ssl_key
< ok 0 0
> exec ^
: string
| create role source_ssl_verify_server_cert
< ok 0 0
> exec ^
: string
| create role source_tls_ciphersuites
< ok 0 0
> exec ^
: string
| create role source_tls_version
< ok 0 0
> exec ^
: string
| create role source_user
< ok 0 0
> exec ^
: string
| create role source_zstd_compression_level
< ok 0 0
> exec ^
: string
| create role spatial
< error
| 1064 42000 …line 1 column 19 near "spatial" 
> exec ^
: string
| create role specific
< ok 0 0
> exec ^
: string
| create role split
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role sql
< error
| 1064 42000 …line 1 column 15 near "sql" 
> exec ^
: string
| create role sql_after_gtids
< ok 0 0
> exec ^
: string
| create role sql_after_mts_gaps
< ok 0 0
> exec ^
: string
| create role sql_before_gtids
< ok 0 0
> exec ^
: string
| create role sql_big_result
< error
| 1064 42000 …line 1 column 26 near "sql_big_result" 
> exec ^
: string
| create role sql_buffer_result
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role sql_calc_found_rows
< error
| 1064 42000 …line 1 column 31 near "sql_calc_found_rows" 
> exec ^
: string
| create role sql_no_cache
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role sql_small_result
< error
| 1064 42000 …line 1 column 28 near "sql_small_result" 
> exec ^
: string
| create role sql_thread
< ok 0 0
> exec ^
: string
| create role sql_tsi_day
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role sql_tsi_hour
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role sql_tsi_minute
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role sql_tsi_month
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role sql_tsi_quarter
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| create role sql_tsi_second
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role sql_tsi_week
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role sql_tsi_year
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role sqlexception
< error
| 1064 42000 …line 1 column 24 near "sqlexception" 
> exec ^
: string
| create role sqlstate
< error
| 1064 42000 …line 1 column 20 near "sqlstate" 
> exec ^
: string
| create role sqlwarning
< error
| 1064 42000 …line 1 column 22 near "sqlwarning" 
> exec ^
: string
| create role srid
< ok 0 0
> exec ^
: string
| create role ssl
< error
| 1064 42000 …line 1 column 15 near "ssl" 
> exec ^
: string
| create role st_collect
< ok 0 0
> exec ^
: string
| create role stacked
< ok 0 0
> exec ^
: string
| create role start
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role start_time
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role starting
< error
| 1064 42000 …line 1 column 20 near "starting" 
> exec ^
: string
| create role starts
< ok 0 0
> exec ^
: string
| create role stats_auto_recalc
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role stats_buckets
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role stats_col_choice
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role stats_col_list
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role stats_extended
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role stats_healthy
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role stats_histograms
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role stats_locked
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role stats_meta
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role stats_persistent
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role stats_sample_pages
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| create role stats_sample_rate
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role stats_topn
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role status
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role std
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role stddev
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role stddev_pop
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role stddev_samp
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role stop
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role storage
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role stored
< error
| 1064 42000 …line 1 column 18 near "stored" 
> exec ^
: string
| create role straight_join
< error
| 1064 42000 …line 1 column 25 near "straight_join" 
> exec ^
: string
| create role stream
< ok 0 0
> exec ^
: string
| create role strict_load
< ok 0 0
> exec ^
: string
| create role string
< ok 0 0
> exec ^
: string
| create role subclass_origin
< ok 0 0
> exec ^
: string
| create role subdate
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role subject
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role subpartition
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role subpartitions
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role substr
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role substring
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role sum
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role super
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role survival_preferences
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| create role suspend
< ok 0 0
> exec ^
: string
| create role swaps
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role switch_group
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role switches
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role sysdate
< ok 0 0
> exec ^
: string
| create role system
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role system_user
< ok 0 0
> exec ^
: string
| create role table
< error
| 1064 42000 …line 1 column 17 near "table" 
> exec ^
: string
| create role table_checksum
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| create role table_name
< ok 0 0
> exec ^
: string
| create role tables
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role tablesample
< error
| 1064 42000 …line 1 column 23 near "tablesample" 
> exec ^
: string
| create role tablespace
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role task_types
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role temporary
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role temptable
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role terminated
< error
| 1064 42000 …line 1 column 22 near "terminated" 
> exec ^
: string
| create role text
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role than
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role then
< error
| 1064 42000 …line 1 column 16 near "then" 
> exec ^
: string
| create role thread_priority
< ok 0 0
> exec ^
: string
| create role ties
< ok 0 0
> exec ^
: string
| create role tiflash
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role time
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role timeout
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role timestamp
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role timestampadd
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create role timestampdiff
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role tinyblob
< error
| 1064 42000 …line 1 column 20 near "tinyblob" 
> exec ^
: string
| create role tinyint
< error
| 1064 42000 …line 1 column 19 near "tinyint" 
> exec ^
: string
| create role tinytext
< error
| 1064 42000 …line 1 column 20 near "tinytext" 
> exec ^
: string
| create role tls
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role to
< error
| 1064 42000 …line 1 column 14 near "to" 
> exec ^
: string
| create role tpcc
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role tpch_10
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role trailing
< error
| 1064 42000 …line 1 column 20 near "trailing" 
> exec ^
: string
| create role transaction
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role trigger
< error
| 1064 42000 …line 1 column 19 near "trigger" 
> exec ^
: string
| create role triggers
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role trim
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role true
< error
| 1064 42000 …line 1 column 16 near "true" 
> exec ^
: string
| create role truncate
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role ttl
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| create role ttl_enable
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role ttl_job_interval
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| create role type
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role types
< ok 0 0
> exec ^
: string
| create role unbounded
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role uncommitted
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| create role undefined
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role undo
< ok 0 0
> exec ^
: string
| create role undo_buffer_size
< ok 0 0
> exec ^
: string
| create role undofile
< ok 0 0
> exec ^
: string
| create role unicode
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role uninstall
< ok 0 0
> exec ^
: string
| create role union
< error
| 1064 42000 …line 1 column 17 near "union" 
> exec ^
: string
| create role unique
< error
| 1064 42000 …line 1 column 18 near "unique" 
> exec ^
: string
| create role unknown
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role unlimited
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role unlock
< error
| 1064 42000 …line 1 column 18 near "unlock" 
> exec ^
: string
| create role unregister
< ok 0 0
> exec ^
: string
| create role unsigned
< error
| 1064 42000 …line 1 column 20 near "unsigned" 
> exec ^
: string
| create role until
< error
| 1064 42000 …line 1 column 17 near "until" 
> exec ^
: string
| create role update
< error
| 1064 42000 …line 1 column 18 near "update" 
> exec ^
: string
| create role upgrade
< ok 0 0
> exec ^
: string
| create role uri
< ok 0 0
> exec ^
: string
| create role url
< ok 0 0
> exec ^
: string
| create role usage
< error
| 1064 42000 …line 1 column 17 near "usage" 
> exec ^
: string
| create role use
< error
| 1064 42000 …line 1 column 15 near "use" 
> exec ^
: string
| create role use_frm
< ok 0 0
> exec ^
: string
| create role user
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role user_resources
< ok 0 0
> exec ^
: string
| create role using
< error
| 1064 42000 …line 1 column 17 near "using" 
> exec ^
: string
| create role utc_date
< error
| 1064 42000 …line 1 column 20 near "utc_date" 
> exec ^
: string
| create role utc_time
< error
| 1064 42000 …line 1 column 20 near "utc_time" 
> exec ^
: string
| create role utc_timestamp
< error
| 1064 42000 …line 1 column 25 near "utc_timestamp" 
> exec ^
: string
| create role utilization_limit
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role validate
< ok 0 0
> exec ^
: string
| create role validation
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create role value
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role values
< error
| 1064 42000 …line 1 column 18 near "values" 
> exec ^
: string
| create role var_pop
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role var_samp
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role varbinary
< error
| 1064 42000 …line 1 column 21 near "varbinary" 
> exec ^
: string
| create role varchar
< error
| 1064 42000 …line 1 column 19 near "varchar" 
> exec ^
: string
| create role varcharacter
< error
| 1064 42000 …line 1 column 24 near "varcharacter" 
> exec ^
: string
| create role variables
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| create role variance
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role varying
< error
| 1064 42000 …line 1 column 19 near "varying" 
> exec ^
: string
| create role vcpu
< ok 0 0
> exec ^
: string
| create role vector
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role verify_key_constraints
< ok 0 0
> exec ^
: string
| create role view
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role virtual
< error
| 1064 42000 …line 1 column 19 near "virtual" 
> exec ^
: string
| create role visible
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role voter_constraints
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| create role voters
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role wait
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role warnings
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role watch
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| create role week
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role weight_string
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| create role when
< error
| 1064 42000 …line 1 column 16 near "when" 
> exec ^
: string
| create role where
< error
| 1064 42000 …line 1 column 17 near "where" 
> exec ^
: string
| create role while
< error
| 1064 42000 …line 1 column 17 near "while" 
> exec ^
: string
| create role window
< error
| 1064 42000 …line 1 column 18 near "window" 
> exec ^
: string
| create role with
< error
| 1064 42000 …line 1 column 16 near "with" 
> exec ^
: string
| create role without
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| create role work
< ok 0 0
> exec ^
: string
| create role workload
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role wrapper
< ok 0 0
> exec ^
: string
| create role write
< error
| 1064 42000 …line 1 column 17 near "write" 
> exec ^
: string
| create role x509
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role xa
< ok 0 0
> exec ^
: string
| create role xid
< ok 0 0
> exec ^
: string
| create role xml
< ok 0 0
> exec ^
: string
| create role xor
< error
| 1064 42000 …line 1 column 15 near "xor" 
> exec ^
: string
| create role year
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role year_month
< error
| 1064 42000 …line 1 column 22 near "year_month" 
> exec ^
: string
| create role zerofill
< error
| 1064 42000 …line 1 column 20 near "zerofill" 
> exec ^
: string
| create role zone
< ok 0 0
> exec ^
: string
| create role absent@h
< ok 0 0
> exec ^
: string
| create role accessible@h
< ok 0 0
> exec ^
: string
| create role account@h
< ok 0 0
> exec ^
: string
| create role action@h
< ok 0 0
> exec ^
: string
| create role active@h
< ok 0 0
> exec ^
: string
| create role add@h
< error
| 1064 42000 …line 1 column 15 near "add@h" 
> exec ^
: string
| create role adddate@h
< ok 0 0
> exec ^
: string
| create role admin@h
< ok 0 0
> exec ^
: string
| create role after@h
< ok 0 0
> exec ^
: string
| create role against@h
< ok 0 0
> exec ^
: string
| create role aggregate@h
< ok 0 0
> exec ^
: string
| create role algorithm@h
< ok 0 0
> exec ^
: string
| create role all@h
< error
| 1064 42000 …line 1 column 15 near "all@h" 
> exec ^
: string
| create role allow_missing_files@h
< ok 0 0
> exec ^
: string
| create role alter@h
< error
| 1064 42000 …line 1 column 17 near "alter@h" 
> exec ^
: string
| create role always@h
< ok 0 0
> exec ^
: string
| create role analyze@h
< error
| 1064 42000 …line 1 column 19 near "analyze@h" 
> exec ^
: string
| create role and@h
< error
| 1064 42000 …line 1 column 15 near "and@h" 
> exec ^
: string
| create role any@h
< ok 0 0
> exec ^
: string
| create role array@h
< error
| 1064 42000 …line 1 column 17 near "array@h" 
> exec ^
: string
| create role as@h
< error
| 1064 42000 …line 1 column 14 near "as@h" 
> exec ^
: string
| create role asc@h
< error
| 1064 42000 …line 1 column 15 near "asc@h" 
> exec ^
: string
| create role ascii@h
< ok 0 0
> exec ^
: string
| create role asensitive@h
< ok 0 0
> exec ^
: string
| create role assign_gtids_to_anonymous_transactions@h
< ok 0 0
> exec ^
: string
| create role at@h
< ok 0 0
> exec ^
: string
| create role attribute@h
< ok 0 0
> exec ^
: string
| create role authentication@h
< ok 0 0
> exec ^
: string
| create role auto@h
< ok 0 0
> exec ^
: string
| create role auto_id_cache@h
< ok 0 0
> exec ^
: string
| create role auto_increment@h
< ok 0 0
> exec ^
: string
| create role auto_random@h
< ok 0 0
> exec ^
: string
| create role auto_random_base@h
< ok 0 0
> exec ^
: string
| create role auto_refresh@h
< ok 0 0
> exec ^
: string
| create role auto_refresh_source@h
< ok 0 0
> exec ^
: string
| create role autoextend_size@h
< ok 0 0
> exec ^
: string
| create role avg@h
< ok 0 0
> exec ^
: string
| create role avg_row_length@h
< ok 0 0
> exec ^
: string
| create role background@h
< ok 0 0
> exec ^
: string
| create role backup@h
< ok 0 0
> exec ^
: string
| create role backups@h
< ok 0 0
> exec ^
: string
| create role batch@h
< ok 0 0
> exec ^
: string
| create role before@h
< ok 0 0
> exec ^
: string
| create role begin@h
< ok 0 0
> exec ^
: string
| create role bernoulli@h
< ok 0 0
> exec ^
: string
| create role between@h
< error
| 1064 42000 …line 1 column 19 near "between@h" 
> exec ^
: string
| create role bigint@h
< error
| 1064 42000 …line 1 column 18 near "bigint@h" 
> exec ^
: string
| create role binary@h
< error
| 1064 42000 …line 1 column 18 near "binary@h" 
> exec ^
: string
| create role binding@h
< ok 0 0
> exec ^
: string
| create role binding_cache@h
< ok 0 0
> exec ^
: string
| create role bindings@h
< ok 0 0
> exec ^
: string
| create role binlog@h
< ok 0 0
> exec ^
: string
| create role bit@h
< ok 0 0
> exec ^
: string
| create role bit_and@h
< ok 0 0
> exec ^
: string
| create role bit_or@h
< ok 0 0
> exec ^
: string
| create role bit_xor@h
< ok 0 0
> exec ^
: string
| create role blob@h
< error
| 1064 42000 …line 1 column 16 near "blob@h" 
> exec ^
: string
| create role block@h
< ok 0 0
> exec ^
: string
| create role bool@h
< ok 0 0
> exec ^
: string
| create role boolean@h
< ok 0 0
> exec ^
: string
| create role both@h
< error
| 1064 42000 …line 1 column 16 near "both@h" 
> exec ^
: string
| create role btree@h
< ok 0 0
> exec ^
: string
| create role buckets@h
< ok 0 0
> exec ^
: string
| create role builtins@h
< ok 0 0
> exec ^
: string
| create role bulk@h
< ok 0 0
> exec ^
: string
| create role burstable@h
< ok 0 0
> exec ^
: string
| create role by@h
< error
| 1064 42000 …line 1 column 14 near "by@h" 
> exec ^
: string
| create role byte@h
< ok 0 0
> exec ^
: string
| create role cache@h
< ok 0 0
> exec ^
: string
| create role calibrate@h
< ok 0 0
> exec ^
: string
| create role call@h
< error
| 1064 42000 …line 1 column 16 near "call@h" 
> exec ^
: string
| create role cascade@h
< error
| 1064 42000 …line 1 column 19 near "cascade@h" 
> exec ^
: string
| create role cascaded@h
< ok 0 0
> exec ^
: string
| create role case@h
< error
| 1064 42000 …line 1 column 16 near "case@h" 
> exec ^
: string
| create role cast@h
< ok 0 0
> exec ^
: string
| create role catalog_name@h
< ok 0 0
> exec ^
: string
| create role chain@h
< ok 0 0
> exec ^
: string
| create role challenge_response@h
< ok 0 0
> exec ^
: string
| create role change@h
< error
| 1064 42000 …line 1 column 18 near "change@h" 
> exec ^
: string
| create role changed@h
< ok 0 0
> exec ^
: string
| create role channel@h
< ok 0 0
> exec ^
: string
| create role char@h
< error
| 1064 42000 …line 1 column 16 near "char@h" 
> exec ^
: string
| create role character@h
< error
| 1064 42000 …line 1 column 21 near "character@h" 
> exec ^
: string
| create role charset@h
< ok 0 0
> exec ^
: string
| create role check@h
< error
| 1064 42000 …line 1 column 17 near "check@h" 
> exec ^
: string
| create role checksum@h
< ok 0 0
> exec ^
: string
| create role cipher@h
< ok 0 0
> exec ^
: string
| create role class_origin@h
< ok 0 0
> exec ^
: string
| create role client@h
< ok 0 0
> exec ^
: string
| create role clone@h
< ok 0 0
> exec ^
: string
| create role close@h
< ok 0 0
> exec ^
: string
| create role clustered@h
< ok 0 0
> exec ^
: string
| create role coalesce@h
< ok 0 0
> exec ^
: string
| create role code@h
< ok 0 0
> exec ^
: string
| create role collate@h
< error
| 1064 42000 …line 1 column 19 near "collate@h" 
> exec ^
: string
| create role collation@h
< ok 0 0
> exec ^
: string
| create role column@h
< error
| 1064 42000 …line 1 column 18 near "column@h" 
> exec ^
: string
| create role column_format@h
< ok 0 0
> exec ^
: string
| create role column_name@h
< ok 0 0
> exec ^
: string
| create role column_stats_usage@h
< ok 0 0
> exec ^
: string
| create role columns@h
< ok 0 0
> exec ^
: string
| create role comment@h
< ok 0 0
> exec ^
: string
| create role commit@h
< ok 0 0
> exec ^
: string
| create role committed@h
< ok 0 0
> exec ^
: string
| create role compact@h
< ok 0 0
> exec ^
: string
| create role completion@h
< ok 0 0
> exec ^
: string
| create role component@h
< ok 0 0
> exec ^
: string
| create role compressed@h
< ok 0 0
> exec ^
: string
| create role compression@h
< ok 0 0
> exec ^
: string
| create role concat@h
< ok 0 0
> exec ^
: string
| create role concurrent@h
< ok 0 0
> exec ^
: string
| create role condition@h
< ok 0 0
> exec ^
: string
| create role config@h
< ok 0 0
> exec ^
: string
| create role connection@h
< ok 0 0
> exec ^
: string
| create role consistent@h
< ok 0 0
> exec ^
: string
| create role constraint@h
< error
| 1064 42000 …line 1 column 22 near "constraint@h" 
> exec ^
: string
| create role constraint_catalog@h
< ok 0 0
> exec ^
: string
| create role constraint_name@h
< ok 0 0
> exec ^
: string
| create role constraint_schema@h
< ok 0 0
> exec ^
: string
| create role constraints@h
< ok 0 0
> exec ^
: string
| create role contains@h
< ok 0 0
> exec ^
: string
| create role context@h
< ok 0 0
> exec ^
: string
| create role continue@h
< error
| 1064 42000 …line 1 column 20 near "continue@h" 
> exec ^
: string
| create role convert@h
< error
| 1064 42000 …line 1 column 19 near "convert@h" 
> exec ^
: string
| create role cooldown@h
< ok 0 0
> exec ^
: string
| create role count@h
< ok 0 0
> exec ^
: string
| create role cpu@h
< ok 0 0
> exec ^
: string
| create role create@h
< error
| 1064 42000 …line 1 column 18 near "create@h" 
> exec ^
: string
| create role cross@h
< error
| 1064 42000 …line 1 column 17 near "cross@h" 
> exec ^
: string
| create role cube@h
< ok 0 0
> exec ^
: string
| create role cume_dist@h
< error
| 1064 42000 …line 1 column 21 near "cume_dist@h" 
> exec ^
: string
| create role curdate@h
< ok 0 0
> exec ^
: string
| create role current@h
< ok 0 0
> exec ^
: string
| create role current_date@h
< error
| 1064 42000 …line 1 column 24 near "current_date@h" 
> exec ^
: string
| create role current_time@h
< error
| 1064 42000 …line 1 column 24 near "current_time@h" 
> exec ^
: string
| create role current_timestamp@h
< error
| 1064 42000 …line 1 column 29 near "current_timestamp@h" 
> exec ^
: string
| create role current_user@h
< error
| 1064 42000 …line 1 column 24 near "current_user@h" 
> exec ^
: string
| create role cursor@h
< error
| 1064 42000 …line 1 column 18 near "cursor@h" 
> exec ^
: string
| create role cursor_name@h
< ok 0 0
> exec ^
: string
| create role curtime@h
< ok 0 0
> exec ^
: string
| create role cycle@h
< ok 0 0
> exec ^
: string
| create role data@h
< ok 0 0
> exec ^
: string
| create role database@h
< error
| 1064 42000 …line 1 column 20 near "database@h" 
> exec ^
: string
| create role databases@h
< error
| 1064 42000 …line 1 column 21 near "databases@h" 
> exec ^
: string
| create role datafile@h
< ok 0 0
> exec ^
: string
| create role date@h
< ok 0 0
> exec ^
: string
| create role date_add@h
< ok 0 0
> exec ^
: string
| create role date_sub@h
< ok 0 0
> exec ^
: string
| create role datetime@h
< ok 0 0
> exec ^
: string
| create role day@h
< ok 0 0
> exec ^
: string
| create role day_hour@h
< error
| 1064 42000 …line 1 column 20 near "day_hour@h" 
> exec ^
: string
| create role day_microsecond@h
< error
| 1064 42000 …line 1 column 27 near "day_microsecond@h" 
> exec ^
: string
| create role day_minute@h
< error
| 1064 42000 …line 1 column 22 near "day_minute@h" 
> exec ^
: string
| create role day_second@h
< error
| 1064 42000 …line 1 column 22 near "day_second@h" 
> exec ^
: string
| create role deallocate@h
< ok 0 0
> exec ^
: string
| create role dec@h
< error
| 1064 42000 …line 1 column 15 near "dec@h" 
> exec ^
: string
| create role decimal@h
< error
| 1064 42000 …line 1 column 19 near "decimal@h" 
> exec ^
: string
| create role declare@h
< ok 0 0
> exec ^
: string
| create role default@h
< error
| 1064 42000 …line 1 column 19 near "default@h" 
> exec ^
: string
| create role default_auth@h
< ok 0 0
> exec ^
: string
| create role defined@h
< ok 0 0
> exec ^
: string
| create role definer@h
< ok 0 0
> exec ^
: string
| create role definition@h
< ok 0 0
> exec ^
: string
| create role delay_key_write@h
< ok 0 0
> exec ^
: string
| create role delayed@h
< error
| 1064 42000 …line 1 column 19 near "delayed@h" 
> exec ^
: string
| create role delete@h
< error
| 1064 42000 …line 1 column 18 near "delete@h" 
> exec ^
: string
| create role dense_rank@h
< error
| 1064 42000 …line 1 column 22 near "dense_rank@h" 
> exec ^
: string
| create role desc@h
< error
| 1064 42000 …line 1 column 16 near "desc@h" 
> exec ^
: string
| create role describe@h
< error
| 1064 42000 …line 1 column 20 near "describe@h" 
> exec ^
: string
| create role description@h
< ok 0 0
> exec ^
: string
| create role deterministic@h
< ok 0 0
> exec ^
: string
| create role diagnostics@h
< ok 0 0
> exec ^
: string
| create role digest@h
< ok 0 0
> exec ^
: string
| create role directory@h
< ok 0 0
> exec ^
: string
| create role disable@h
< ok 0 0
> exec ^
: string
| create role disabled@h
< ok 0 0
> exec ^
: string
| create role discard@h
< ok 0 0
> exec ^
: string
| create role disk@h
< ok 0 0
> exec ^
: string
| create role distinct@h
< error
| 1064 42000 …line 1 column 20 near "distinct@h" 
> exec ^
: string
| create role distinctrow@h
< error
| 1064 42000 …line 1 column 23 near "distinctrow@h" 
> exec ^
: string
| create role distribute@h
< ok 0 0
> exec ^
: string
| create role distribution@h
< ok 0 0
> exec ^
: string
| create role distributions@h
< ok 0 0
> exec ^
: string
| create role div@h
< error
| 1064 42000 …line 1 column 15 near "div@h" 
> exec ^
: string
| create role do@h
< ok 0 0
> exec ^
: string
| create role double@h
< error
| 1064 42000 …line 1 column 18 near "double@h" 
> exec ^
: string
| create role drop@h
< error
| 1064 42000 …line 1 column 16 near "drop@h" 
> exec ^
: string
| create role dry@h
< ok 0 0
> exec ^
: string
| create role dryrun@h
< ok 0 0
> exec ^
: string
| create role dual@h
< error
| 1064 42000 …line 1 column 16 near "dual@h" 
> exec ^
: string
| create role duality@h
< ok 0 0
> exec ^
: string
| create role dumpfile@h
< ok 0 0
> exec ^
: string
| create role duplicate@h
< ok 0 0
> exec ^
: string
| create role duration@h
< ok 0 0
> exec ^
: string
| create role dynamic@h
< ok 0 0
> exec ^
: string
| create role each@h
< ok 0 0
> exec ^
: string
| create role else@h
< error
| 1064 42000 …line 1 column 16 near "else@h" 
> exec ^
: string
| create role elseif@h
< error
| 1064 42000 …line 1 column 18 near "elseif@h" 
> exec ^
: string
| create role empty@h
< ok 0 0
> exec ^
: string
| create role enable@h
< ok 0 0
> exec ^
: string
| create role enabled@h
< ok 0 0
> exec ^
: string
| create role enclosed@h
< error
| 1064 42000 …line 1 column 20 near "enclosed@h" 
> exec ^
: string
| create role encryption@h
< ok 0 0
> exec ^
: string
| create role end@h
< ok 0 0
> exec ^
: string
| create role end_time@h
< ok 0 0
> exec ^
: string
| create role ends@h
< ok 0 0
> exec ^
: string
| create role enforced@h
< ok 0 0
> exec ^
: string
| create role engine@h
< ok 0 0
> exec ^
: string
| create role engine_attribute@h
< ok 0 0
> exec ^
: string
| create role engines@h
< ok 0 0
> exec ^
: string
| create role enum@h
< ok 0 0
> exec ^
: string
| create role error@h
< ok 0 0
> exec ^
: string
| create role errors@h
< ok 0 0
> exec ^
: string
| create role escape@h
< ok 0 0
> exec ^
: string
| create role escaped@h
< error
| 1064 42000 …line 1 column 19 near "escaped@h" 
> exec ^
: string
| create role event@h
< ok 0 0
> exec ^
: string
| create role events@h
< ok 0 0
> exec ^
: string
| create role every@h
< ok 0 0
> exec ^
: string
| create role exact@h
< ok 0 0
> exec ^
: string
| create role except@h
< error
| 1064 42000 …line 1 column 18 near "except@h" 
> exec ^
: string
| create role exchange@h
< ok 0 0
> exec ^
: string
| create role exclude@h
< ok 0 0
> exec ^
: string
| create role exec_elapsed@h
< ok 0 0
> exec ^
: string
| create role execute@h
< ok 0 0
> exec ^
: string
| create role exists@h
< error
| 1064 42000 …line 1 column 18 near "exists@h" 
> exec ^
: string
| create role exit@h
< error
| 1064 42000 …line 1 column 16 near "exit@h" 
> exec ^
: string
| create role expansion@h
< ok 0 0
> exec ^
: string
| create role expire@h
< ok 0 0
> exec ^
: string
| create role explain@h
< error
| 1064 42000 …line 1 column 19 near "explain@h" 
> exec ^
: string
| create role export@h
< ok 0 0
> exec ^
: string
| create role extended@h
< ok 0 0
> exec ^
: string
| create role extent_size@h
< ok 0 0
> exec ^
: string
| create role external@h
< ok 0 0
> exec ^
: string
| create role external_format@h
< ok 0 0
> exec ^
: string
| create role extract@h
< ok 0 0
> exec ^
: string
| create role factor@h
< ok 0 0
> exec ^
: string
| create role failed_login_attempts@h
< ok 0 0
> exec ^
: string
| create role false@h
< error
| 1064 42000 …line 1 column 17 near "false@h" 
> exec ^
: string
| create role fast@h
< ok 0 0
> exec ^
: string
| create role faults@h
< ok 0 0
> exec ^
: string
| create role fetch@h
< error
| 1064 42000 …line 1 column 17 near "fetch@h" 
> exec ^
: string
| create role fields@h
< ok 0 0
> exec ^
: string
| create role file@h
< ok 0 0
> exec ^
: string
| create role file_block_size@h
< ok 0 0
> exec ^
: string
| create role file_format@h
< ok 0 0
> exec ^
: string
| create role file_name@h
< ok 0 0
> exec ^
: string
| create role file_pattern@h
< ok 0 0
> exec ^
: string
| create role file_prefix@h
< ok 0 0
> exec ^
: string
| create role files@h
< ok 0 0
> exec ^
: string
| create role filter@h
< ok 0 0
> exec ^
: string
| create role finish@h
< ok 0 0
> exec ^
: string
| create role first@h
< ok 0 0
> exec ^
: string
| create role first_value@h
< error
| 1064 42000 …line 1 column 23 near "first_value@h" 
> exec ^
: string
| create role fixed@h
< ok 0 0
> exec ^
: string
| create role flashback@h
< ok 0 0
> exec ^
: string
| create role float@h
< error
| 1064 42000 …line 1 column 17 near "float@h" 
> exec ^
: string
| create role float4@h
< error
| 1064 42000 …line 1 column 18 near "float4@h" 
> exec ^
: string
| create role float8@h
< error
| 1064 42000 …line 1 column 18 near "float8@h" 
> exec ^
: string
| create role flush@h
< ok 0 0
> exec ^
: string
| create role follower_constraints@h
< ok 0 0
> exec ^
: string
| create role followers@h
< ok 0 0
> exec ^
: string
| create role following@h
< ok 0 0
> exec ^
: string
| create role follows@h
< ok 0 0
> exec ^
: string
| create role for@h
< error
| 1064 42000 …line 1 column 15 near "for@h" 
> exec ^
: string
| create role force@h
< error
| 1064 42000 …line 1 column 17 near "force@h" 
> exec ^
: string
| create role foreign@h
< error
| 1064 42000 …line 1 column 19 near "foreign@h" 
> exec ^
: string
| create role format@h
< ok 0 0
> exec ^
: string
| create role found@h
< ok 0 0
> exec ^
: string
| create role from@h
< error
| 1064 42000 …line 1 column 16 near "from@h" 
> exec ^
: string
| create role full@h
< ok 0 0
> exec ^
: string
| create role fulltext@h
< error
| 1064 42000 …line 1 column 20 near "fulltext@h" 
> exec ^
: string
| create role function@h
< ok 0 0
> exec ^
: string
| create role general@h
< ok 0 0
> exec ^
: string
| create role generate@h
< ok 0 0
> exec ^
: string
| create role generated@h
< error
| 1064 42000 …line 1 column 21 near "generated@h" 
> exec ^
: string
| create role geomcollection@h
< ok 0 0
> exec ^
: string
| create role geometry@h
< ok 0 0
> exec ^
: string
| create role geometrycollection@h
< ok 0 0
> exec ^
: string
| create role get@h
< ok 0 0
> exec ^
: string
| create role get_format@h
< ok 0 0
> exec ^
: string
| create role get_source_public_key@h
< ok 0 0
> exec ^
: string
| create role global@h
< ok 0 0
> exec ^
: string
| create role grant@h
< error
| 1064 42000 …line 1 column 17 near "grant@h" 
> exec ^
: string
| create role grants@h
< ok 0 0
> exec ^
: string
| create role greatest@h
< ok 0 0
> exec ^
: string
| create role group@h
< error
| 1064 42000 …line 1 column 17 near "group@h" 
> exec ^
: string
| create role group_concat@h
< ok 0 0
> exec ^
: string
| create role group_replication@h
< ok 0 0
> exec ^
: string
| create role grouping@h
< ok 0 0
> exec ^
: string
| create role groups@h
< error
| 1064 42000 …line 1 column 18 near "groups@h" 
> exec ^
: string
| create role gtid_only@h
< ok 0 0
> exec ^
: string
| create role gtids@h
< ok 0 0
> exec ^
: string
| create role guided@h
< ok 0 0
> exec ^
: string
| create role handler@h
< ok 0 0
> exec ^
: string
| create role hash@h
< ok 0 0
> exec ^
: string
| create role having@h
< error
| 1064 42000 …line 1 column 18 near "having@h" 
> exec ^
: string
| create role header@h
< ok 0 0
> exec ^
: string
| create role help@h
< ok 0 0
> exec ^
: string
| create role high@h
< ok 0 0
> exec ^
: string
| create role high_priority@h
< error
| 1064 42000 …line 1 column 25 near "high_priority@h" 
> exec ^
: string
| create role histogram@h
< ok 0 0
> exec ^
: string
| create role histograms_in_flight@h
< ok 0 0
> exec ^
: string
| create role history@h
< ok 0 0
> exec ^
: string
| create role host@h
< ok 0 0
> exec ^
: string
| create role hosts@h
< ok 0 0
> exec ^
: string
| create role hour@h
< ok 0 0
> exec ^
: string
| create role hour_microsecond@h
< error
| 1064 42000 …line 1 column 28 near "hour_microsecond@h" 
> exec ^
: string
| create role hour_minute@h
< error
| 1064 42000 …line 1 column 23 near "hour_minute@h" 
> exec ^
: string
| create role hour_second@h
< error
| 1064 42000 …line 1 column 23 near "hour_second@h" 
> exec ^
: string
| create role identified@h
< ok 0 0
> exec ^
: string
| create role if@h
< error
| 1064 42000 …line 1 column 16 near "@h" 
> exec ^
: string
| create role ignore@h
< error
| 1064 42000 …line 1 column 18 near "ignore@h" 
> exec ^
: string
| create role ignore_server_ids@h
< ok 0 0
> exec ^
: string
| create role ilike@h
< error
| 1064 42000 …line 1 column 17 near "ilike@h" 
> exec ^
: string
| create role import@h
< ok 0 0
> exec ^
: string
| create role in@h
< error
| 1064 42000 …line 1 column 14 near "in@h" 
> exec ^
: string
| create role inactive@h
< ok 0 0
> exec ^
: string
| create role increment@h
< ok 0 0
> exec ^
: string
| create role index@h
< error
| 1064 42000 …line 1 column 17 near "index@h" 
> exec ^
: string
| create role indexes@h
< ok 0 0
> exec ^
: string
| create role infile@h
< error
| 1064 42000 …line 1 column 18 near "infile@h" 
> exec ^
: string
| create role initial@h
< ok 0 0
> exec ^
: string
| create role initial_size@h
< ok 0 0
> exec ^
: string
| create role initiate@h
< ok 0 0
> exec ^
: string
| create role inner@h
< error
| 1064 42000 …line 1 column 17 near "inner@h" 
> exec ^
: string
| create role inout@h
< error
| 1064 42000 …line 1 column 17 near "inout@h" 
> exec ^
: string
| create role insensitive@h
< ok 0 0
> exec ^
: string
| create role insert@h
< error
| 1064 42000 …line 1 column 18 near "insert@h" 
> exec ^
: string
| create role insert_method@h
< ok 0 0
> exec ^
: string
| create role install@h
< ok 0 0
> exec ^
: string
| create role instance@h
< ok 0 0
> exec ^
: string
| create role int@h
< error
| 1064 42000 …line 1 column 15 near "int@h" 
> exec ^
: string
| create role int1@h
< error
| 1064 42000 …line 1 column 16 near "int1@h" 
> exec ^
: string
| create role int2@h
< error
| 1064 42000 …line 1 column 16 near "int2@h" 
> exec ^
: string
| create role int3@h
< error
| 1064 42000 …line 1 column 16 near "int3@h" 
> exec ^
: string
| create role int4@h
< error
| 1064 42000 …line 1 column 16 near "int4@h" 
> exec ^
: string
| create role int8@h
< error
| 1064 42000 …line 1 column 16 near "int8@h" 
> exec ^
: string
| create role integer@h
< error
| 1064 42000 …line 1 column 19 near "integer@h" 
> exec ^
: string
| create role intersect@h
< error
| 1064 42000 …line 1 column 21 near "intersect@h" 
> exec ^
: string
| create role interval@h
< error
| 1064 42000 …line 1 column 20 near "interval@h" 
> exec ^
: string
| create role into@h
< error
| 1064 42000 …line 1 column 16 near "into@h" 
> exec ^
: string
| create role invisible@h
< ok 0 0
> exec ^
: string
| create role invoker@h
< ok 0 0
> exec ^
: string
| create role io@h
< ok 0 0
> exec ^
: string
| create role io_after_gtids@h
< ok 0 0
> exec ^
: string
| create role io_before_gtids@h
< ok 0 0
> exec ^
: string
| create role io_thread@h
< ok 0 0
> exec ^
: string
| create role ipc@h
< ok 0 0
> exec ^
: string
| create role is@h
< error
| 1064 42000 …line 1 column 14 near "is@h" 
> exec ^
: string
| create role isolation@h
< ok 0 0
> exec ^
: string
| create role issuer@h
< ok 0 0
> exec ^
: string
| create role iterate@h
< error
| 1064 42000 …line 1 column 19 near "iterate@h" 
> exec ^
: string
| create role job@h
< ok 0 0
> exec ^
: string
| create role jobs@h
< ok 0 0
> exec ^
: string
| create role join@h
< error
| 1064 42000 …line 1 column 16 near "join@h" 
> exec ^
: string
| create role json@h
< ok 0 0
> exec ^
: string
| create role json_arrayagg@h
< ok 0 0
> exec ^
: string
| create role json_duality_object@h
< ok 0 0
> exec ^
: string
| create role json_objectagg@h
< ok 0 0
> exec ^
: string
| create role json_table@h
< ok 0 0
> exec ^
: string
| create role json_value@h
< ok 0 0
> exec ^
: string
| create role key@h
< error
| 1064 42000 …line 1 column 15 near "key@h" 
> exec ^
: string
| create role key_block_size@h
< ok 0 0
> exec ^
: string
| create role keyring@h
< ok 0 0
> exec ^
: string
| create role keys@h
< error
| 1064 42000 …line 1 column 16 near "keys@h" 
> exec ^
: string
| create role kill@h
< error
| 1064 42000 …line 1 column 16 near "kill@h" 
> exec ^
: string
| create role labels@h
< ok 0 0
> exec ^
: string
| create role lag@h
< error
| 1064 42000 …line 1 column 15 near "lag@h" 
> exec ^
: string
| create role language@h
< ok 0 0
> exec ^
: string
| create role last@h
< ok 0 0
> exec ^
: string
| create role last_value@h
< error
| 1064 42000 …line 1 column 22 near "last_value@h" 
> exec ^
: string
| create role lateral@h
< ok 0 0
> exec ^
: string
| create role lead@h
< error
| 1064 42000 …line 1 column 16 near "lead@h" 
> exec ^
: string
| create role leader_constraints@h
< ok 0 0
> exec ^
: string
| create role leading@h
< error
| 1064 42000 …line 1 column 19 near "leading@h" 
> exec ^
: string
| create role learner_constraints@h
< ok 0 0
> exec ^
: string
| create role learners@h
< ok 0 0
> exec ^
: string
| create role least@h
< ok 0 0
> exec ^
: string
| create role leave@h
< error
| 1064 42000 …line 1 column 17 near "leave@h" 
> exec ^
: string
| create role leaves@h
< ok 0 0
> exec ^
: string
| create role left@h
< error
| 1064 42000 …line 1 column 16 near "left@h" 
> exec ^
: string
| create role less@h
< ok 0 0
> exec ^
: string
| create role level@h
< ok 0 0
> exec ^
: string
| create role library@h
< ok 0 0
> exec ^
: string
| create role like@h
< error
| 1064 42000 …line 1 column 16 near "like@h" 
> exec ^
: string
| create role limit@h
< error
| 1064 42000 …line 1 column 17 near "limit@h" 
> exec ^
: string
| create role linear@h
< error
| 1064 42000 …line 1 column 18 near "linear@h" 
> exec ^
: string
| create role lines@h
< error
| 1064 42000 …line 1 column 17 near "lines@h" 
> exec ^
: string
| create role linestring@h
< ok 0 0
> exec ^
: string
| create role list@h
< ok 0 0
> exec ^
: string
| create role load@h
< error
| 1064 42000 …line 1 column 16 near "load@h" 
> exec ^
: string
| create role local@h
< ok 0 0
> exec ^
: string
| create role localtime@h
< error
| 1064 42000 …line 1 column 21 near "localtime@h" 
> exec ^
: string
| create role localtimestamp@h
< error
| 1064 42000 …line 1 column 26 near "localtimestamp@h" 
> exec ^
: string
| create role location@h
< ok 0 0
> exec ^
: string
| create role lock@h
< error
| 1064 42000 …line 1 column 16 near "lock@h" 
> exec ^
: string
| create role locked@h
< ok 0 0
> exec ^
: string
| create role locks@h
< ok 0 0
> exec ^
: string
| create role log@h
< ok 0 0
> exec ^
: string
| create role logfile@h
< ok 0 0
> exec ^
: string
| create role logs@h
< ok 0 0
> exec ^
: string
| create role long@h
< error
| 1064 42000 …line 1 column 16 near "long@h" 
> exec ^
: string
| create role longblob@h
< error
| 1064 42000 …line 1 column 20 near "longblob@h" 
> exec ^
: string
| create role longtext@h
< error
| 1064 42000 …line 1 column 20 near "longtext@h" 
> exec ^
: string
| create role loop@h
< ok 0 0
> exec ^
: string
| create role low@h
< ok 0 0
> exec ^
: string
| create role low_priority@h
< error
| 1064 42000 …line 1 column 24 near "low_priority@h" 
> exec ^
: string
| create role manual@h
< ok 0 0
> exec ^
: string
| create role master@h
< ok 0 0
> exec ^
: string
| create role master_log_file@h
< ok 0 0
> exec ^
: string
| create role match@h
< error
| 1064 42000 …line 1 column 17 near "match@h" 
> exec ^
: string
| create role materialized@h
< ok 0 0
> exec ^
: string
| create role max@h
< ok 0 0
> exec ^
: string
| create role max_connections_per_hour@h
< ok 0 0
> exec ^
: string
| create role max_queries_per_hour@h
< ok 0 0
> exec ^
: string
| create role max_rows@h
< ok 0 0
> exec ^
: string
| create role max_size@h
< ok 0 0
> exec ^
: string
| create role max_updates_per_hour@h
< ok 0 0
> exec ^
: string
| create role max_user_connections@h
< ok 0 0
> exec ^
: string
| create role maxvalue@h
< error
| 1064 42000 …line 1 column 20 near "maxvalue@h" 
> exec ^
: string
| create role medium@h
< ok 0 0
> exec ^
: string
| create role mediumblob@h
< error
| 1064 42000 …line 1 column 22 near "mediumblob@h" 
> exec ^
: string
| create role mediumint@h
< error
| 1064 42000 …line 1 column 21 near "mediumint@h" 
> exec ^
: string
| create role mediumtext@h
< error
| 1064 42000 …line 1 column 22 near "mediumtext@h" 
> exec ^
: string
| create role member@h
< ok 0 0
> exec ^
: string
| create role memory@h
< ok 0 0
> exec ^
: string
| create role merge@h
< ok 0 0
> exec ^
: string
| create role message_text@h
< ok 0 0
> exec ^
: string
| create role microsecond@h
< ok 0 0
> exec ^
: string
| create role mid@h
< ok 0 0
> exec ^
: string
| create role middleint@h
< error
| 1064 42000 …line 1 column 21 near "middleint@h" 
> exec ^
: string
| create role migrate@h
< ok 0 0
> exec ^
: string
| create role min@h
< ok 0 0
> exec ^
: string
| create role min_rows@h
< ok 0 0
> exec ^
: string
| create role minute@h
< ok 0 0
> exec ^
: string
| create role minute_microsecond@h
< error
| 1064 42000 …line 1 column 30 near "minute_microsecond@h" 
> exec ^
: string
| create role minute_second@h
< error
| 1064 42000 …line 1 column 25 near "minute_second@h" 
> exec ^
: string
| create role minvalue@h
< ok 0 0
> exec ^
: string
| create role mod@h
< error
| 1064 42000 …line 1 column 15 near "mod@h" 
> exec ^
: string
| create role mode@h
< ok 0 0
> exec ^
: string
| create role modifies@h
< ok 0 0
> exec ^
: string
| create role modify@h
< ok 0 0
> exec ^
: string
| create role month@h
< ok 0 0
> exec ^
: string
| create role multilinestring@h
< ok 0 0
> exec ^
: string
| create role multipoint@h
< ok 0 0
> exec ^
: string
| create role multipolygon@h
< ok 0 0
> exec ^
: string
| create role mutex@h
< ok 0 0
> exec ^
: string
| create role mysql_errno@h
< ok 0 0
> exec ^
: string
| create role name@h
< ok 0 0
> exec ^
: string
| create role names@h
< ok 0 0
> exec ^
: string
| create role national@h
< ok 0 0
> exec ^
: string
| create role natural@h
< error
| 1064 42000 …line 1 column 19 near "natural@h" 
> exec ^
: string
| create role nchar@h
< ok 0 0
> exec ^
: string
| create role ndb@h
< ok 0 0
> exec ^
: string
| create role ndbcluster@h
< ok 0 0
> exec ^
: string
| create role nested@h
< ok 0 0
> exec ^
: string
| create role network_namespace@h
< ok 0 0
> exec ^
: string
| create role never@h
< ok 0 0
> exec ^
: string
| create role new@h
< ok 0 0
> exec ^
: string
| create role next@h
< ok 0 0
> exec ^
: string
| create role next_row_id@h
< ok 0 0
> exec ^
: string
| create role no@h
< ok 0 0
> exec ^
: string
| create role no_wait@h
< ok 0 0
> exec ^
: string
| create role no_write_to_binlog@h
< error
| 1064 42000 …line 1 column 30 near "no_write_to_binlog@h" 
> exec ^
: string
| create role nocache@h
< ok 0 0
> exec ^
: string
| create role nocycle@h
< ok 0 0
> exec ^
: string
| create role nodegroup@h
< ok 0 0
> exec ^
: string
| create role nomaxvalue@h
< ok 0 0
> exec ^
: string
| create role nominvalue@h
< ok 0 0
> exec ^
: string
| create role nonclustered@h
< ok 0 0
> exec ^
: string
| create role none@h
< ok 0 0
> exec ^
: string
| create role not@h
< error
| 1064 42000 …line 1 column 15 near "not@h" 
> exec ^
: string
| create role now@h
< ok 0 0
> exec ^
: string
| create role nowait@h
< ok 0 0
> exec ^
: string
| create role nth_value@h
< error
| 1064 42000 …line 1 column 21 near "nth_value@h" 
> exec ^
: string
| create role ntile@h
< error
| 1064 42000 …line 1 column 17 near "ntile@h" 
> exec ^
: string
| create role null@h
< error
| 1064 42000 …line 1 column 16 near "null@h" 
> exec ^
: string
| create role nullif@h
< ok 0 0
> exec ^
: string
| create role nulls@h
< ok 0 0
> exec ^
: string
| create role number@h
< ok 0 0
> exec ^
: string
| create role numeric@h
< error
| 1064 42000 …line 1 column 19 near "numeric@h" 
> exec ^
: string
| create role nvarchar@h
< ok 0 0
> exec ^
: string
| create role of@h
< error
| 1064 42000 …line 1 column 14 near "of@h" 
> exec ^
: string
| create role off@h
< ok 0 0
> exec ^
: string
| create role offset@h
< ok 0 0
> exec ^
: string
| create role oj@h
< ok 0 0
> exec ^
: string
| create role old@h
< ok 0 0
> exec ^
: string
| create role oltp_read_only@h
< ok 0 0
> exec ^
: string
| create role oltp_read_write@h
< ok 0 0
> exec ^
: string
| create role oltp_write_only@h
< ok 0 0
> exec ^
: string
| create role on@h
< error
| 1064 42000 …line 1 column 14 near "on@h" 
> exec ^
: string
| create role one@h
< ok 0 0
> exec ^
: string
| create role only@h
< ok 0 0
> exec ^
: string
| create role open@h
< ok 0 0
> exec ^
: string
| create role optimize@h
< error
| 1064 42000 …line 1 column 20 near "optimize@h" 
> exec ^
: string
| create role optimizer_costs@h
< ok 0 0
> exec ^
: string
| create role option@h
< error
| 1064 42000 …line 1 column 18 near "option@h" 
> exec ^
: string
| create role optional@h
< ok 0 0
> exec ^
: string
| create role optionally@h
< error
| 1064 42000 …line 1 column 22 near "optionally@h" 
> exec ^
: string
| create role options@h
< ok 0 0
> exec ^
: string
| create role or@h
< error
| 1064 42000 …line 1 column 14 near "or@h" 
> exec ^
: string
| create role order@h
< error
| 1064 42000 …line 1 column 17 near "order@h" 
> exec ^
: string
| create role ordinality@h
< ok 0 0
> exec ^
: string
| create role organization@h
< ok 0 0
> exec ^
: string
| create role others@h
< ok 0 0
> exec ^
: string
| create role out@h
< error
| 1064 42000 …line 1 column 15 near "out@h" 
> exec ^
: string
| create role outer@h
< error
| 1064 42000 …line 1 column 17 near "outer@h" 
> exec ^
: string
| create role outfile@h
< error
| 1064 42000 …line 1 column 19 near "outfile@h" 
> exec ^
: string
| create role over@h
< error
| 1064 42000 …line 1 column 16 near "over@h" 
> exec ^
: string
| create role overlay@h
< ok 0 0
> exec ^
: string
| create role owner@h
< ok 0 0
> exec ^
: string
| create role pack_keys@h
< ok 0 0
> exec ^
: string
| create role page@h
< ok 0 0
> exec ^
: string
| create role parallel@h
< ok 0 0
> exec ^
: string
| create role parameters@h
< ok 0 0
> exec ^
: string
| create role parse_tree@h
< ok 0 0
> exec ^
: string
| create role parser@h
< ok 0 0
> exec ^
: string
| create role partial@h
< ok 0 0
> exec ^
: string
| create role partition@h
< error
| 1064 42000 …line 1 column 21 near "partition@h" 
> exec ^
: string
| create role partitioning@h
< ok 0 0
> exec ^
: string
| create role partitions@h
< ok 0 0
> exec ^
: string
| create role password@h
< ok 0 0
> exec ^
: string
| create role password_lock_time@h
< ok 0 0
> exec ^
: string
| create role path@h
< ok 0 0
> exec ^
: string
| create role percent_rank@h
< error
| 1064 42000 …line 1 column 24 near "percent_rank@h" 
> exec ^
: string
| create role persist@h
< ok 0 0
> exec ^
: string
| create role persist_only@h
< ok 0 0
> exec ^
: string
| create role phase@h
< ok 0 0
> exec ^
: string
| create role placement@h
< ok 0 0
> exec ^
: string
| create role placing@h
< ok 0 0
> exec ^
: string
| create role plan@h
< ok 0 0
> exec ^
: string
| create role plugin@h
< ok 0 0
> exec ^
: string
| create role plugin_dir@h
< ok 0 0
> exec ^
: string
| create role plugins@h
< ok 0 0
> exec ^
: string
| create role point@h
< ok 0 0
> exec ^
: string
| create role policy@h
< ok 0 0
> exec ^
: string
| create role polygon@h
< ok 0 0
> exec ^
: string
| create role port@h
< ok 0 0
> exec ^
: string
| create role position@h
< ok 0 0
> exec ^
: string
| create role pre_split_regions@h
< ok 0 0
> exec ^
: string
| create role precedes@h
< ok 0 0
> exec ^
: string
| create role preceding@h
< ok 0 0
> exec ^
: string
| create role precision@h
< error
| 1064 42000 …line 1 column 21 near "precision@h" 
> exec ^
: string
| create role prepare@h
< ok 0 0
> exec ^
: string
| create role preserve@h
< ok 0 0
> exec ^
: string
| create role prev@h
< ok 0 0
> exec ^
: string
| create role primary@h
< error
| 1064 42000 …line 1 column 19 near "primary@h" 
> exec ^
: string
| create role primary_region@h
< ok 0 0
> exec ^
: string
| create role priority@h
< ok 0 0
> exec ^
: string
| create role privilege_checks_user@h
< ok 0 0
> exec ^
: string
| create role privileges@h
< ok 0 0
> exec ^
: string
| create role procedure@h
< error
| 1064 42000 …line 1 column 21 near "procedure@h" 
> exec ^
: string
| create role process@h
< ok 0 0
> exec ^
: string
| create role processed_keys@h
< ok 0 0
> exec ^
: string
| create role processlist@h
< ok 0 0
> exec ^
: string
| create role profile@h
< ok 0 0
> exec ^
: string
| create role profiles@h
< ok 0 0
> exec ^
: string
| create role proxy@h
< ok 0 0
> exec ^
: string
| create role purge@h
< ok 0 0
> exec ^
: string
| create role qualify@h
< ok 0 0
> exec ^
: string
| create role quarter@h
< ok 0 0
> exec ^
: string
| create role query@h
< ok 0 0
> exec ^
: string
| create role query_limit@h
< ok 0 0
> exec ^
: string
| create role quick@h
< ok 0 0
> exec ^
: string
| create role random@h
< ok 0 0
> exec ^
: string
| create role range@h
< error
| 1064 42000 …line 1 column 17 near "range@h" 
> exec ^
: string
| create role rank@h
< error
| 1064 42000 …line 1 column 16 near "rank@h" 
> exec ^
: string
| create role read@h
< error
| 1064 42000 …line 1 column 16 near "read@h" 
> exec ^
: string
| create role read_only@h
< ok 0 0
> exec ^
: string
| create role read_write@h
< ok 0 0
> exec ^
: string
| create role reads@h
< ok 0 0
> exec ^
: string
| create role real@h
< error
| 1064 42000 …line 1 column 16 near "real@h" 
> exec ^
: string
| create role rebuild@h
< ok 0 0
> exec ^
: string
| create role recover@h
< ok 0 0
> exec ^
: string
| create role recursive@h
< error
| 1064 42000 …line 1 column 21 near "recursive@h" 
> exec ^
: string
| create role redo_buffer_size@h
< ok 0 0
> exec ^
: string
| create role redundant@h
< ok 0 0
> exec ^
: string
| create role reference@h
< ok 0 0
> exec ^
: string
| create role references@h
< error
| 1064 42000 …line 1 column 22 near "references@h" 
> exec ^
: string
| create role regexp@h
< error
| 1064 42000 …line 1 column 18 near "regexp@h" 
> exec ^
: string
| create role region@h
< ok 0 0
> exec ^
: string
| create role regions@h
< ok 0 0
> exec ^
: string
| create role registration@h
< ok 0 0
> exec ^
: string
| create role relational@h
< ok 0 0
> exec ^
: string
| create role relay@h
< ok 0 0
> exec ^
: string
| create role relay_log_file@h
< ok 0 0
> exec ^
: string
| create role relay_log_pos@h
< ok 0 0
> exec ^
: string
| create role relay_thread@h
< ok 0 0
> exec ^
: string
| create role relaylog@h
< ok 0 0
> exec ^
: string
| create role release@h
< error
| 1064 42000 …line 1 column 19 near "release@h" 
> exec ^
: string
| create role reload@h
< ok 0 0
> exec ^
: string
| create role remove@h
< ok 0 0
> exec ^
: string
| create role rename@h
< error
| 1064 42000 …line 1 column 18 near "rename@h" 
> exec ^
: string
| create role reorganize@h
< ok 0 0
> exec ^
: string
| create role repair@h
< ok 0 0
> exec ^
: string
| create role repeat@h
< error
| 1064 42000 …line 1 column 18 near "repeat@h" 
> exec ^
: string
| create role repeatable@h
< ok 0 0
> exec ^
: string
| create role replace@h
< error
| 1064 42000 …line 1 column 19 near "replace@h" 
> exec ^
: string
| create role replica@h
< ok 0 0
> exec ^
: string
| create role replicas@h
< ok 0 0
> exec ^
: string
| create role replicate_do_db@h
< ok 0 0
> exec ^
: string
| create role replicate_do_table@h
< ok 0 0
> exec ^
: string
| create role replicate_ignore_db@h
< ok 0 0
> exec ^
: string
| create role replicate_ignore_table@h
< ok 0 0
> exec ^
: string
| create role replicate_rewrite_db@h
< ok 0 0
> exec ^
: string
| create role replicate_wild_do_table@h
< ok 0 0
> exec ^
: string
| create role replicate_wild_ignore_table@h
< ok 0 0
> exec ^
: string
| create role replication@h
< ok 0 0
> exec ^
: string
| create role require@h
< error
| 1064 42000 …line 1 column 19 near "require@h" 
> exec ^
: string
| create role require_row_format@h
< ok 0 0
> exec ^
: string
| create role require_table_primary_key_check@h
< ok 0 0
> exec ^
: string
| create role reset@h
< ok 0 0
> exec ^
: string
| create role resignal@h
< ok 0 0
> exec ^
: string
| create role resource@h
< ok 0 0
> exec ^
: string
| create role respect@h
< ok 0 0
> exec ^
: string
| create role restart@h
< ok 0 0
> exec ^
: string
| create role restore@h
< ok 0 0
> exec ^
: string
| create role restores@h
< ok 0 0
> exec ^
: string
| create role restrict@h
< error
| 1064 42000 …line 1 column 20 near "restrict@h" 
> exec ^
: string
| create role resume@h
< ok 0 0
> exec ^
: string
| create role retain@h
< ok 0 0
> exec ^
: string
| create role return@h
< ok 0 0
> exec ^
: string
| create role returned_sqlstate@h
< ok 0 0
> exec ^
: string
| create role returning@h
< ok 0 0
> exec ^
: string
| create role returns@h
< ok 0 0
> exec ^
: string
| create role reuse@h
< ok 0 0
> exec ^
: string
| create role reverse@h
< ok 0 0
> exec ^
: string
| create role revoke@h
< error
| 1064 42000 …line 1 column 18 near "revoke@h" 
> exec ^
: string
| create role right@h
< error
| 1064 42000 …line 1 column 17 near "right@h" 
> exec ^
: string
| create role rlike@h
< error
| 1064 42000 …line 1 column 17 near "rlike@h" 
> exec ^
: string
| create role role@h
< ok 0 0
> exec ^
: string
| create role rollback@h
< ok 0 0
> exec ^
: string
| create role rollup@h
< ok 0 0
> exec ^
: string
| create role rotate@h
< ok 0 0
> exec ^
: string
| create role routine@h
< ok 0 0
> exec ^
: string
| create role row@h
< error
| 1064 42000 …line 1 column 15 near "row@h" 
> exec ^
: string
| create role row_count@h
< ok 0 0
> exec ^
: string
| create role row_format@h
< ok 0 0
> exec ^
: string
| create role row_number@h
< error
| 1064 42000 …line 1 column 22 near "row_number@h" 
> exec ^
: string
| create role rows@h
< error
| 1064 42000 …line 1 column 16 near "rows@h" 
> exec ^
: string
| create role rtree@h
< ok 0 0
> exec ^
: string
| create role ru@h
< ok 0 0
> exec ^
: string
| create role ru_per_sec@h
< ok 0 0
> exec ^
: string
| create role rule@h
< ok 0 0
> exec ^
: string
| create role run@h
< ok 0 0
> exec ^
: string
| create role s3@h
< ok 0 0
> exec ^
: string
| create role savepoint@h
< ok 0 0
> exec ^
: string
| create role schedule@h
< ok 0 0
> exec ^
: string
| create role schema@h
< error
| 1064 42000 …line 1 column 18 near "schema@h" 
> exec ^
: string
| create role schema_name@h
< ok 0 0
> exec ^
: string
| create role schemas@h
< error
| 1064 42000 …line 1 column 19 near "schemas@h" 
> exec ^
: string
| create role second@h
< ok 0 0
> exec ^
: string
| create role second_microsecond@h
< error
| 1064 42000 …line 1 column 30 near "second_microsecond@h" 
> exec ^
: string
| create role secondary@h
< ok 0 0
> exec ^
: string
| create role secondary_engine@h
< ok 0 0
> exec ^
: string
| create role secondary_engine_attribute@h
< ok 0 0
> exec ^
: string
| create role secondary_load@h
< ok 0 0
> exec ^
: string
| create role secondary_unload@h
< ok 0 0
> exec ^
: string
| create role security@h
< ok 0 0
> exec ^
: string
| create role select@h
< error
| 1064 42000 …line 1 column 18 near "select@h" 
> exec ^
: string
| create role sensitive@h
< ok 0 0
> exec ^
: string
| create role separator@h
< ok 0 0
> exec ^
: string
| create role sequence@h
< ok 0 0
> exec ^
: string
| create role serial@h
< ok 0 0
> exec ^
: string
| create role serializable@h
< ok 0 0
> exec ^
: string
| create role server@h
< ok 0 0
> exec ^
: string
| create role session@h
< ok 0 0
> exec ^
: string
| create role session_states@h
< ok 0 0
> exec ^
: string
| create role session_user@h
< ok 0 0
> exec ^
: string
| create role set@h
< error
| 1064 42000 …line 1 column 15 near "set@h" 
> exec ^
: string
| create role sets@h
< ok 0 0
> exec ^
: string
| create role shard_row_id_bits@h
< ok 0 0
> exec ^
: string
| create role share@h
< ok 0 0
> exec ^
: string
| create role show@h
< error
| 1064 42000 …line 1 column 16 near "show@h" 
> exec ^
: string
| create role shutdown@h
< ok 0 0
> exec ^
: string
| create role signal@h
< ok 0 0
> exec ^
: string
| create role signed@h
< ok 0 0
> exec ^
: string
| create role similar@h
< ok 0 0
> exec ^
: string
| create role simple@h
< ok 0 0
> exec ^
: string
| create role skip@h
< ok 0 0
> exec ^
: string
| create role slave@h
< ok 0 0
> exec ^
: string
| create role slow@h
< ok 0 0
> exec ^
: string
| create role smallint@h
< error
| 1064 42000 …line 1 column 20 near "smallint@h" 
> exec ^
: string
| create role snapshot@h
< ok 0 0
> exec ^
: string
| create role socket@h
< ok 0 0
> exec ^
: string
| create role some@h
< ok 0 0
> exec ^
: string
| create role soname@h
< ok 0 0
> exec ^
: string
| create role sounds@h
< ok 0 0
> exec ^
: string
| create role source@h
< ok 0 0
> exec ^
: string
| create role source_auto_position@h
< ok 0 0
> exec ^
: string
| create role source_bind@h
< ok 0 0
> exec ^
: string
| create role source_compression_algorithms@h
< ok 0 0
> exec ^
: string
| create role source_connect_retry@h
< ok 0 0
> exec ^
: string
| create role source_connection_auto_failover@h
< ok 0 0
> exec ^
: string
| create role source_delay@h
< ok 0 0
> exec ^
: string
| create role source_heartbeat_period@h
< ok 0 0
> exec ^
: string
| create role source_host@h
< ok 0 0
> exec ^
: string
| create role source_log_file@h
< ok 0 0
> exec ^
: string
| create role source_log_pos@h
< ok 0 0
> exec ^
: string
| create role source_password@h
< ok 0 0
> exec ^
: string
| create role source_port@h
< ok 0 0
> exec ^
: string
| create role source_public_key_path@h
< ok 0 0
> exec ^
: string
| create role source_retry_count@h
< ok 0 0
> exec ^
: string
| create role source_ssl@h
< ok 0 0
> exec ^
: string
| create role source_ssl_ca@h
< ok 0 0
> exec ^
: string
| create role source_ssl_capath@h
< ok 0 0
> exec ^
: string
| create role source_ssl_cert@h
< ok 0 0
> exec ^
: string
| create role source_ssl_cipher@h
< ok 0 0
> exec ^
: string
| create role source_ssl_crl@h
< ok 0 0
> exec ^
: string
| create role source_ssl_crlpath@h
< ok 0 0
> exec ^
: string
| create role source_ssl_key@h
< ok 0 0
> exec ^
: string
| create role source_ssl_verify_server_cert@h
< ok 0 0
> exec ^
: string
| create role source_tls_ciphersuites@h
< ok 0 0
> exec ^
: string
| create role source_tls_version@h
< ok 0 0
> exec ^
: string
| create role source_user@h
< ok 0 0
> exec ^
: string
| create role source_zstd_compression_level@h
< ok 0 0
> exec ^
: string
| create role spatial@h
< error
| 1064 42000 …line 1 column 19 near "spatial@h" 
> exec ^
: string
| create role specific@h
< ok 0 0
> exec ^
: string
| create role split@h
< ok 0 0
> exec ^
: string
| create role sql@h
< error
| 1064 42000 …line 1 column 15 near "sql@h" 
> exec ^
: string
| create role sql_after_gtids@h
< ok 0 0
> exec ^
: string
| create role sql_after_mts_gaps@h
< ok 0 0
> exec ^
: string
| create role sql_before_gtids@h
< ok 0 0
> exec ^
: string
| create role sql_big_result@h
< error
| 1064 42000 …line 1 column 26 near "sql_big_result@h" 
> exec ^
: string
| create role sql_buffer_result@h
< ok 0 0
> exec ^
: string
| create role sql_calc_found_rows@h
< error
| 1064 42000 …line 1 column 31 near "sql_calc_found_rows@h" 
> exec ^
: string
| create role sql_no_cache@h
< ok 0 0
> exec ^
: string
| create role sql_small_result@h
< error
| 1064 42000 …line 1 column 28 near "sql_small_result@h" 
> exec ^
: string
| create role sql_thread@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_day@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_hour@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_minute@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_month@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_quarter@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_second@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_week@h
< ok 0 0
> exec ^
: string
| create role sql_tsi_year@h
< ok 0 0
> exec ^
: string
| create role sqlexception@h
< error
| 1064 42000 …line 1 column 24 near "sqlexception@h" 
> exec ^
: string
| create role sqlstate@h
< error
| 1064 42000 …line 1 column 20 near "sqlstate@h" 
> exec ^
: string
| create role sqlwarning@h
< error
| 1064 42000 …line 1 column 22 near "sqlwarning@h" 
> exec ^
: string
| create role srid@h
< ok 0 0
> exec ^
: string
| create role ssl@h
< error
| 1064 42000 …line 1 column 15 near "ssl@h" 
> exec ^
: string
| create role st_collect@h
< ok 0 0
> exec ^
: string
| create role stacked@h
< ok 0 0
> exec ^
: string
| create role start@h
< ok 0 0
> exec ^
: string
| create role start_time@h
< ok 0 0
> exec ^
: string
| create role starting@h
< error
| 1064 42000 …line 1 column 20 near "starting@h" 
> exec ^
: string
| create role starts@h
< ok 0 0
> exec ^
: string
| create role stats_auto_recalc@h
< ok 0 0
> exec ^
: string
| create role stats_buckets@h
< ok 0 0
> exec ^
: string
| create role stats_col_choice@h
< ok 0 0
> exec ^
: string
| create role stats_col_list@h
< ok 0 0
> exec ^
: string
| create role stats_extended@h
< ok 0 0
> exec ^
: string
| create role stats_healthy@h
< ok 0 0
> exec ^
: string
| create role stats_histograms@h
< ok 0 0
> exec ^
: string
| create role stats_locked@h
< ok 0 0
> exec ^
: string
| create role stats_meta@h
< ok 0 0
> exec ^
: string
| create role stats_persistent@h
< ok 0 0
> exec ^
: string
| create role stats_sample_pages@h
< ok 0 0
> exec ^
: string
| create role stats_sample_rate@h
< ok 0 0
> exec ^
: string
| create role stats_topn@h
< ok 0 0
> exec ^
: string
| create role status@h
< ok 0 0
> exec ^
: string
| create role std@h
< ok 0 0
> exec ^
: string
| create role stddev@h
< ok 0 0
> exec ^
: string
| create role stddev_pop@h
< ok 0 0
> exec ^
: string
| create role stddev_samp@h
< ok 0 0
> exec ^
: string
| create role stop@h
< ok 0 0
> exec ^
: string
| create role storage@h
< ok 0 0
> exec ^
: string
| create role stored@h
< error
| 1064 42000 …line 1 column 18 near "stored@h" 
> exec ^
: string
| create role straight_join@h
< error
| 1064 42000 …line 1 column 25 near "straight_join@h" 
> exec ^
: string
| create role stream@h
< ok 0 0
> exec ^
: string
| create role strict_load@h
< ok 0 0
> exec ^
: string
| create role string@h
< ok 0 0
> exec ^
: string
| create role subclass_origin@h
< ok 0 0
> exec ^
: string
| create role subdate@h
< ok 0 0
> exec ^
: string
| create role subject@h
< ok 0 0
> exec ^
: string
| create role subpartition@h
< ok 0 0
> exec ^
: string
| create role subpartitions@h
< ok 0 0
> exec ^
: string
| create role substr@h
< ok 0 0
> exec ^
: string
| create role substring@h
< ok 0 0
> exec ^
: string
| create role sum@h
< ok 0 0
> exec ^
: string
| create role super@h
< ok 0 0
> exec ^
: string
| create role survival_preferences@h
< ok 0 0
> exec ^
: string
| create role suspend@h
< ok 0 0
> exec ^
: string
| create role swaps@h
< ok 0 0
> exec ^
: string
| create role switch_group@h
< ok 0 0
> exec ^
: string
| create role switches@h
< ok 0 0
> exec ^
: string
| create role sysdate@h
< ok 0 0
> exec ^
: string
| create role system@h
< ok 0 0
> exec ^
: string
| create role system_user@h
< ok 0 0
> exec ^
: string
| create role table@h
< error
| 1064 42000 …line 1 column 17 near "table@h" 
> exec ^
: string
| create role table_checksum@h
< ok 0 0
> exec ^
: string
| create role table_name@h
< ok 0 0
> exec ^
: string
| create role tables@h
< ok 0 0
> exec ^
: string
| create role tablesample@h
< error
| 1064 42000 …line 1 column 23 near "tablesample@h" 
> exec ^
: string
| create role tablespace@h
< ok 0 0
> exec ^
: string
| create role task_types@h
< ok 0 0
> exec ^
: string
| create role temporary@h
< ok 0 0
> exec ^
: string
| create role temptable@h
< ok 0 0
> exec ^
: string
| create role terminated@h
< error
| 1064 42000 …line 1 column 22 near "terminated@h" 
> exec ^
: string
| create role text@h
< ok 0 0
> exec ^
: string
| create role than@h
< ok 0 0
> exec ^
: string
| create role then@h
< error
| 1064 42000 …line 1 column 16 near "then@h" 
> exec ^
: string
| create role thread_priority@h
< ok 0 0
> exec ^
: string
| create role ties@h
< ok 0 0
> exec ^
: string
| create role tiflash@h
< ok 0 0
> exec ^
: string
| create role time@h
< ok 0 0
> exec ^
: string
| create role timeout@h
< ok 0 0
> exec ^
: string
| create role timestamp@h
< ok 0 0
> exec ^
: string
| create role timestampadd@h
< ok 0 0
> exec ^
: string
| create role timestampdiff@h
< ok 0 0
> exec ^
: string
| create role tinyblob@h
< error
| 1064 42000 …line 1 column 20 near "tinyblob@h" 
> exec ^
: string
| create role tinyint@h
< error
| 1064 42000 …line 1 column 19 near "tinyint@h" 
> exec ^
: string
| create role tinytext@h
< error
| 1064 42000 …line 1 column 20 near "tinytext@h" 
> exec ^
: string
| create role tls@h
< ok 0 0
> exec ^
: string
| create role to@h
< error
| 1064 42000 …line 1 column 14 near "to@h" 
> exec ^
: string
| create role tpcc@h
< ok 0 0
> exec ^
: string
| create role tpch_10@h
< ok 0 0
> exec ^
: string
| create role trailing@h
< error
| 1064 42000 …line 1 column 20 near "trailing@h" 
> exec ^
: string
| create role transaction@h
< ok 0 0
> exec ^
: string
| create role trigger@h
< error
| 1064 42000 …line 1 column 19 near "trigger@h" 
> exec ^
: string
| create role triggers@h
< ok 0 0
> exec ^
: string
| create role trim@h
< ok 0 0
> exec ^
: string
| create role true@h
< error
| 1064 42000 …line 1 column 16 near "true@h" 
> exec ^
: string
| create role truncate@h
< ok 0 0
> exec ^
: string
| create role ttl@h
< ok 0 0
> exec ^
: string
| create role ttl_enable@h
< ok 0 0
> exec ^
: string
| create role ttl_job_interval@h
< ok 0 0
> exec ^
: string
| create role type@h
< ok 0 0
> exec ^
: string
| create role types@h
< ok 0 0
> exec ^
: string
| create role unbounded@h
< ok 0 0
> exec ^
: string
| create role uncommitted@h
< ok 0 0
> exec ^
: string
| create role undefined@h
< ok 0 0
> exec ^
: string
| create role undo@h
< ok 0 0
> exec ^
: string
| create role undo_buffer_size@h
< ok 0 0
> exec ^
: string
| create role undofile@h
< ok 0 0
> exec ^
: string
| create role unicode@h
< ok 0 0
> exec ^
: string
| create role uninstall@h
< ok 0 0
> exec ^
: string
| create role union@h
< error
| 1064 42000 …line 1 column 17 near "union@h" 
> exec ^
: string
| create role unique@h
< error
| 1064 42000 …line 1 column 18 near "unique@h" 
> exec ^
: string
| create role unknown@h
< ok 0 0
> exec ^
: string
| create role unlimited@h
< ok 0 0
> exec ^
: string
| create role unlock@h
< error
| 1064 42000 …line 1 column 18 near "unlock@h" 
> exec ^
: string
| create role unregister@h
< ok 0 0
> exec ^
: string
| create role unsigned@h
< error
| 1064 42000 …line 1 column 20 near "unsigned@h" 
> exec ^
: string
| create role until@h
< error
| 1064 42000 …line 1 column 17 near "until@h" 
> exec ^
: string
| create role update@h
< error
| 1064 42000 …line 1 column 18 near "update@h" 
> exec ^
: string
| create role upgrade@h
< ok 0 0
> exec ^
: string
| create role uri@h
< ok 0 0
> exec ^
: string
| create role url@h
< ok 0 0
> exec ^
: string
| create role usage@h
< error
| 1064 42000 …line 1 column 17 near "usage@h" 
> exec ^
: string
| create role use@h
< error
| 1064 42000 …line 1 column 15 near "use@h" 
> exec ^
: string
| create role use_frm@h
< ok 0 0
> exec ^
: string
| create role user@h
< ok 0 0
> exec ^
: string
| create role user_resources@h
< ok 0 0
> exec ^
: string
| create role using@h
< error
| 1064 42000 …line 1 column 17 near "using@h" 
> exec ^
: string
| create role utc_date@h
< error
| 1064 42000 …line 1 column 20 near "utc_date@h" 
> exec ^
: string
| create role utc_time@h
< error
| 1064 42000 …line 1 column 20 near "utc_time@h" 
> exec ^
: string
| create role utc_timestamp@h
< error
| 1064 42000 …line 1 column 25 near "utc_timestamp@h" 
> exec ^
: string
| create role utilization_limit@h
< ok 0 0
> exec ^
: string
| create role validate@h
< ok 0 0
> exec ^
: string
| create role validation@h
< ok 0 0
> exec ^
: string
| create role value@h
< ok 0 0
> exec ^
: string
| create role values@h
< error
| 1064 42000 …line 1 column 18 near "values@h" 
> exec ^
: string
| create role var_pop@h
< ok 0 0
> exec ^
: string
| create role var_samp@h
< ok 0 0
> exec ^
: string
| create role varbinary@h
< error
| 1064 42000 …line 1 column 21 near "varbinary@h" 
> exec ^
: string
| create role varchar@h
< error
| 1064 42000 …line 1 column 19 near "varchar@h" 
> exec ^
: string
| create role varcharacter@h
< error
| 1064 42000 …line 1 column 24 near "varcharacter@h" 
> exec ^
: string
| create role variables@h
< ok 0 0
> exec ^
: string
| create role variance@h
< ok 0 0
> exec ^
: string
| create role varying@h
< error
| 1064 42000 …line 1 column 19 near "varying@h" 
> exec ^
: string
| create role vcpu@h
< ok 0 0
> exec ^
: string
| create role vector@h
< ok 0 0
> exec ^
: string
| create role verify_key_constraints@h
< ok 0 0
> exec ^
: string
| create role view@h
< ok 0 0
> exec ^
: string
| create role virtual@h
< error
| 1064 42000 …line 1 column 19 near "virtual@h" 
> exec ^
: string
| create role visible@h
< ok 0 0
> exec ^
: string
| create role voter_constraints@h
< ok 0 0
> exec ^
: string
| create role voters@h
< ok 0 0
> exec ^
: string
| create role wait@h
< ok 0 0
> exec ^
: string
| create role warnings@h
< ok 0 0
> exec ^
: string
| create role watch@h
< ok 0 0
> exec ^
: string
| create role week@h
< ok 0 0
> exec ^
: string
| create role weight_string@h
< ok 0 0
> exec ^
: string
| create role when@h
< error
| 1064 42000 …line 1 column 16 near "when@h" 
> exec ^
: string
| create role where@h
< error
| 1064 42000 …line 1 column 17 near "where@h" 
> exec ^
: string
| create role while@h
< error
| 1064 42000 …line 1 column 17 near "while@h" 
> exec ^
: string
| create role window@h
< error
| 1064 42000 …line 1 column 18 near "window@h" 
> exec ^
: string
| create role with@h
< error
| 1064 42000 …line 1 column 16 near "with@h" 
> exec ^
: string
| create role without@h
< ok 0 0
> exec ^
: string
| create role work@h
< ok 0 0
> exec ^
: string
| create role workload@h
< ok 0 0
> exec ^
: string
| create role wrapper@h
< ok 0 0
> exec ^
: string
| create role write@h
< error
| 1064 42000 …line 1 column 17 near "write@h" 
> exec ^
: string
| create role x509@h
< ok 0 0
> exec ^
: string
| create role xa@h
< ok 0 0
> exec ^
: string
| create role xid@h
< ok 0 0
> exec ^
: string
| create role xml@h
< ok 0 0
> exec ^
: string
| create role xor@h
< error
| 1064 42000 …line 1 column 15 near "xor@h" 
> exec ^
: string
| create role year@h
< ok 0 0
> exec ^
: string
| create role year_month@h
< error
| 1064 42000 …line 1 column 22 near "year_month@h" 
> exec ^
: string
| create role zerofill@h
< error
| 1064 42000 …line 1 column 20 near "zerofill@h" 
> exec ^
: string
| create role zone@h
< ok 0 0
> exec ^
: string
| drop role absent
< ok 0 0
> exec ^
: string
| drop role accessible
< ok 0 0
> exec ^
: string
| drop role account
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role action
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role active
< ok 0 0
> exec ^
: string
| drop role add
< error
| 1064 42000 …line 1 column 13 near "add" 
> exec ^
: string
| drop role adddate
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role admin
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role after
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role against
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role aggregate
< ok 0 0
> exec ^
: string
| drop role algorithm
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role all
< error
| 1064 42000 …line 1 column 13 near "all" 
> exec ^
: string
| drop role allow_missing_files
< ok 0 0
> exec ^
: string
| drop role alter
< error
| 1064 42000 …line 1 column 15 near "alter" 
> exec ^
: string
| drop role always
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role analyze
< error
| 1064 42000 …line 1 column 17 near "analyze" 
> exec ^
: string
| drop role and
< error
| 1064 42000 …line 1 column 13 near "and" 
> exec ^
: string
| drop role any
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role array
< error
| 1064 42000 …line 1 column 15 near "array" 
> exec ^
: string
| drop role as
< error
| 1064 42000 …line 1 column 12 near "as" 
> exec ^
: string
| drop role asc
< error
| 1064 42000 …line 1 column 13 near "asc" 
> exec ^
: string
| drop role ascii
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role asensitive
< ok 0 0
> exec ^
: string
| drop role assign_gtids_to_anonymous_transactions
< ok 0 0
> exec ^
: string
| drop role at
< ok 0 0
> exec ^
: string
| drop role attribute
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role authentication
< ok 0 0
> exec ^
: string
| drop role auto
< ok 0 0
> exec ^
: string
| drop role auto_id_cache
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role auto_increment
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role auto_random
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role auto_random_base
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role auto_refresh
< ok 0 0
> exec ^
: string
| drop role auto_refresh_source
< ok 0 0
> exec ^
: string
| drop role autoextend_size
< ok 0 0
> exec ^
: string
| drop role avg
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role avg_row_length
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role background
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role backup
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role backups
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role batch
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role before
< ok 0 0
> exec ^
: string
| drop role begin
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role bernoulli
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role between
< error
| 1064 42000 …line 1 column 17 near "between" 
> exec ^
: string
| drop role bigint
< error
| 1064 42000 …line 1 column 16 near "bigint" 
> exec ^
: string
| drop role binary
< error
| 1064 42000 …line 1 column 16 near "binary" 
> exec ^
: string
| drop role binding
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role binding_cache
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role bindings
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role binlog
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role bit
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role bit_and
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role bit_or
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role bit_xor
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role blob
< error
| 1064 42000 …line 1 column 14 near "blob" 
> exec ^
: string
| drop role block
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role bool
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role boolean
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role both
< error
| 1064 42000 …line 1 column 14 near "both" 
> exec ^
: string
| drop role btree
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role buckets
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role builtins
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role bulk
< ok 0 0
> exec ^
: string
| drop role burstable
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role by
< error
| 1064 42000 …line 1 column 12 near "by" 
> exec ^
: string
| drop role byte
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role cache
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role calibrate
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role call
< error
| 1064 42000 …line 1 column 14 near "call" 
> exec ^
: string
| drop role cascade
< error
| 1064 42000 …line 1 column 17 near "cascade" 
> exec ^
: string
| drop role cascaded
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role case
< error
| 1064 42000 …line 1 column 14 near "case" 
> exec ^
: string
| drop role cast
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role catalog_name
< ok 0 0
> exec ^
: string
| drop role chain
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role challenge_response
< ok 0 0
> exec ^
: string
| drop role change
< error
| 1064 42000 …line 1 column 16 near "change" 
> exec ^
: string
| drop role changed
< ok 0 0
> exec ^
: string
| drop role channel
< ok 0 0
> exec ^
: string
| drop role char
< error
| 1064 42000 …line 1 column 14 near "char" 
> exec ^
: string
| drop role character
< error
| 1064 42000 …line 1 column 19 near "character" 
> exec ^
: string
| drop role charset
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role check
< error
| 1064 42000 …line 1 column 15 near "check" 
> exec ^
: string
| drop role checksum
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role cipher
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role class_origin
< ok 0 0
> exec ^
: string
| drop role client
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role clone
< ok 0 0
> exec ^
: string
| drop role close
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role clustered
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role coalesce
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role code
< ok 0 0
> exec ^
: string
| drop role collate
< error
| 1064 42000 …line 1 column 17 near "collate" 
> exec ^
: string
| drop role collation
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role column
< error
| 1064 42000 …line 1 column 16 near "column" 
> exec ^
: string
| drop role column_format
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role column_name
< ok 0 0
> exec ^
: string
| drop role column_stats_usage
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| drop role columns
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role comment
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role commit
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role committed
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role compact
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role completion
< ok 0 0
> exec ^
: string
| drop role component
< ok 0 0
> exec ^
: string
| drop role compressed
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role compression
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role concat
< ok 0 0
> exec ^
: string
| drop role concurrent
< ok 0 0
> exec ^
: string
| drop role condition
< ok 0 0
> exec ^
: string
| drop role config
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role connection
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role consistent
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role constraint
< error
| 1064 42000 …line 1 column 20 near "constraint" 
> exec ^
: string
| drop role constraint_catalog
< ok 0 0
> exec ^
: string
| drop role constraint_name
< ok 0 0
> exec ^
: string
| drop role constraint_schema
< ok 0 0
> exec ^
: string
| drop role constraints
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role contains
< ok 0 0
> exec ^
: string
| drop role context
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role continue
< error
| 1064 42000 …line 1 column 18 near "continue" 
> exec ^
: string
| drop role convert
< error
| 1064 42000 …line 1 column 17 near "convert" 
> exec ^
: string
| drop role cooldown
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role count
< ok 0 0
> exec ^
: string
| drop role cpu
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role create
< error
| 1064 42000 …line 1 column 16 near "create" 
> exec ^
: string
| drop role cross
< error
| 1064 42000 …line 1 column 15 near "cross" 
> exec ^
: string
| drop role cube
< ok 0 0
> exec ^
: string
| drop role cume_dist
< error
| 1064 42000 …line 1 column 19 near "cume_dist" 
> exec ^
: string
| drop role curdate
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role current
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role current_date
< error
| 1064 42000 …line 1 column 22 near "current_date" 
> exec ^
: string
| drop role current_time
< error
| 1064 42000 …line 1 column 22 near "current_time" 
> exec ^
: string
| drop role current_timestamp
< error
| 1064 42000 …line 1 column 27 near "current_timestamp" 
> exec ^
: string
| drop role current_user
< error
| 1064 42000 …line 1 column 22 near "current_user" 
> exec ^
: string
| drop role cursor
< error
| 1064 42000 …line 1 column 16 near "cursor" 
> exec ^
: string
| drop role cursor_name
< ok 0 0
> exec ^
: string
| drop role curtime
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role cycle
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role data
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role database
< error
| 1064 42000 …line 1 column 18 near "database" 
> exec ^
: string
| drop role databases
< error
| 1064 42000 …line 1 column 19 near "databases" 
> exec ^
: string
| drop role datafile
< ok 0 0
> exec ^
: string
| drop role date
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role date_add
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role date_sub
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role datetime
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role day
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role day_hour
< error
| 1064 42000 …line 1 column 18 near "day_hour" 
> exec ^
: string
| drop role day_microsecond
< error
| 1064 42000 …line 1 column 25 near "day_microsecond" 
> exec ^
: string
| drop role day_minute
< error
| 1064 42000 …line 1 column 20 near "day_minute" 
> exec ^
: string
| drop role day_second
< error
| 1064 42000 …line 1 column 20 near "day_second" 
> exec ^
: string
| drop role deallocate
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role dec
< error
| 1064 42000 …line 1 column 13 near "dec" 
> exec ^
: string
| drop role decimal
< error
| 1064 42000 …line 1 column 17 near "decimal" 
> exec ^
: string
| drop role declare
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role default
< error
| 1064 42000 …line 1 column 17 near "default" 
> exec ^
: string
| drop role default_auth
< ok 0 0
> exec ^
: string
| drop role defined
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role definer
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role definition
< ok 0 0
> exec ^
: string
| drop role delay_key_write
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| drop role delayed
< error
| 1064 42000 …line 1 column 17 near "delayed" 
> exec ^
: string
| drop role delete
< error
| 1064 42000 …line 1 column 16 near "delete" 
> exec ^
: string
| drop role dense_rank
< error
| 1064 42000 …line 1 column 20 near "dense_rank" 
> exec ^
: string
| drop role desc
< error
| 1064 42000 …line 1 column 14 near "desc" 
> exec ^
: string
| drop role describe
< error
| 1064 42000 …line 1 column 18 near "describe" 
> exec ^
: string
| drop role description
< ok 0 0
> exec ^
: string
| drop role deterministic
< ok 0 0
> exec ^
: string
| drop role diagnostics
< ok 0 0
> exec ^
: string
| drop role digest
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role directory
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role disable
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role disabled
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role discard
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role disk
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role distinct
< error
| 1064 42000 …line 1 column 18 near "distinct" 
> exec ^
: string
| drop role distinctrow
< error
| 1064 42000 …line 1 column 21 near "distinctrow" 
> exec ^
: string
| drop role distribute
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role distribution
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role distributions
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role div
< error
| 1064 42000 …line 1 column 13 near "div" 
> exec ^
: string
| drop role do
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role double
< error
| 1064 42000 …line 1 column 16 near "double" 
> exec ^
: string
| drop role drop
< error
| 1064 42000 …line 1 column 14 near "drop" 
> exec ^
: string
| drop role dry
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role dryrun
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role dual
< error
| 1064 42000 …line 1 column 14 near "dual" 
> exec ^
: string
| drop role duality
< ok 0 0
> exec ^
: string
| drop role dumpfile
< ok 0 0
> exec ^
: string
| drop role duplicate
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role duration
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role dynamic
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role each
< ok 0 0
> exec ^
: string
| drop role else
< error
| 1064 42000 …line 1 column 14 near "else" 
> exec ^
: string
| drop role elseif
< error
| 1064 42000 …line 1 column 16 near "elseif" 
> exec ^
: string
| drop role empty
< ok 0 0
> exec ^
: string
| drop role enable
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role enabled
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role enclosed
< error
| 1064 42000 …line 1 column 18 near "enclosed" 
> exec ^
: string
| drop role encryption
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role end
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role end_time
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role ends
< ok 0 0
> exec ^
: string
| drop role enforced
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role engine
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role engine_attribute
< ok 0 0
> exec ^
: string
| drop role engines
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role enum
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role error
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role errors
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role escape
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role escaped
< error
| 1064 42000 …line 1 column 17 near "escaped" 
> exec ^
: string
| drop role event
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role events
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role every
< ok 0 0
> exec ^
: string
| drop role exact
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role except
< error
| 1064 42000 …line 1 column 16 near "except" 
> exec ^
: string
| drop role exchange
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role exclude
< ok 0 0
> exec ^
: string
| drop role exec_elapsed
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role execute
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role exists
< error
| 1064 42000 …line 1 column 16 near "exists" 
> exec ^
: string
| drop role exit
< error
| 1064 42000 …line 1 column 14 near "exit" 
> exec ^
: string
| drop role expansion
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role expire
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role explain
< error
| 1064 42000 …line 1 column 17 near "explain" 
> exec ^
: string
| drop role export
< ok 0 0
> exec ^
: string
| drop role extended
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role extent_size
< ok 0 0
> exec ^
: string
| drop role external
< ok 0 0
> exec ^
: string
| drop role external_format
< ok 0 0
> exec ^
: string
| drop role extract
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role factor
< ok 0 0
> exec ^
: string
| drop role failed_login_attempts
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| drop role false
< error
| 1064 42000 …line 1 column 15 near "false" 
> exec ^
: string
| drop role fast
< ok 0 0
> exec ^
: string
| drop role faults
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role fetch
< error
| 1064 42000 …line 1 column 15 near "fetch" 
> exec ^
: string
| drop role fields
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role file
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role file_block_size
< ok 0 0
> exec ^
: string
| drop role file_format
< ok 0 0
> exec ^
: string
| drop role file_name
< ok 0 0
> exec ^
: string
| drop role file_pattern
< ok 0 0
> exec ^
: string
| drop role file_prefix
< ok 0 0
> exec ^
: string
| drop role files
< ok 0 0
> exec ^
: string
| drop role filter
< ok 0 0
> exec ^
: string
| drop role finish
< ok 0 0
> exec ^
: string
| drop role first
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role first_value
< error
| 1064 42000 …line 1 column 21 near "first_value" 
> exec ^
: string
| drop role fixed
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role flashback
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role float
< error
| 1064 42000 …line 1 column 15 near "float" 
> exec ^
: string
| drop role float4
< error
| 1064 42000 …line 1 column 16 near "float4" 
> exec ^
: string
| drop role float8
< error
| 1064 42000 …line 1 column 16 near "float8" 
> exec ^
: string
| drop role flush
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role follower_constraints
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| drop role followers
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role following
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role follows
< ok 0 0
> exec ^
: string
| drop role for
< error
| 1064 42000 …line 1 column 13 near "for" 
> exec ^
: string
| drop role force
< error
| 1064 42000 …line 1 column 15 near "force" 
> exec ^
: string
| drop role foreign
< error
| 1064 42000 …line 1 column 17 near "foreign" 
> exec ^
: string
| drop role format
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role found
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role from
< error
| 1064 42000 …line 1 column 14 near "from" 
> exec ^
: string
| drop role full
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role fulltext
< error
| 1064 42000 …line 1 column 18 near "fulltext" 
> exec ^
: string
| drop role function
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role general
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role generate
< ok 0 0
> exec ^
: string
| drop role generated
< error
| 1064 42000 …line 1 column 19 near "generated" 
> exec ^
: string
| drop role geomcollection
< ok 0 0
> exec ^
: string
| drop role geometry
< ok 0 0
> exec ^
: string
| drop role geometrycollection
< ok 0 0
> exec ^
: string
| drop role get
< ok 0 0
> exec ^
: string
| drop role get_format
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role get_source_public_key
< ok 0 0
> exec ^
: string
| drop role global
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role grant
< error
| 1064 42000 …line 1 column 15 near "grant" 
> exec ^
: string
| drop role grants
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role greatest
< ok 0 0
> exec ^
: string
| drop role group
< error
| 1064 42000 …line 1 column 15 near "group" 
> exec ^
: string
| drop role group_concat
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role group_replication
< ok 0 0
> exec ^
: string
| drop role grouping
< ok 0 0
> exec ^
: string
| drop role groups
< error
| 1064 42000 …line 1 column 16 near "groups" 
> exec ^
: string
| drop role gtid_only
< ok 0 0
> exec ^
: string
| drop role gtids
< ok 0 0
> exec ^
: string
| drop role guided
< ok 0 0
> exec ^
: string
| drop role handler
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role hash
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role having
< error
| 1064 42000 …line 1 column 16 near "having" 
> exec ^
: string
| drop role header
< ok 0 0
> exec ^
: string
| drop role help
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role high
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role high_priority
< error
| 1064 42000 …line 1 column 23 near "high_priority" 
> exec ^
: string
| drop role histogram
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role histograms_in_flight
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| drop role history
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role host
< ok 0 0
> exec ^
: string
| drop role hosts
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role hour
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role hour_microsecond
< error
| 1064 42000 …line 1 column 26 near "hour_microsecond" 
> exec ^
: string
| drop role hour_minute
< error
| 1064 42000 …line 1 column 21 near "hour_minute" 
> exec ^
: string
| drop role hour_second
< error
| 1064 42000 …line 1 column 21 near "hour_second" 
> exec ^
: string
| drop role identified
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role if
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role ignore
< error
| 1064 42000 …line 1 column 16 near "ignore" 
> exec ^
: string
| drop role ignore_server_ids
< ok 0 0
> exec ^
: string
| drop role ilike
< error
| 1064 42000 …line 1 column 15 near "ilike" 
> exec ^
: string
| drop role import
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role in
< error
| 1064 42000 …line 1 column 12 near "in" 
> exec ^
: string
| drop role inactive
< ok 0 0
> exec ^
: string
| drop role increment
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role index
< error
| 1064 42000 …line 1 column 15 near "index" 
> exec ^
: string
| drop role indexes
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role infile
< error
| 1064 42000 …line 1 column 16 near "infile" 
> exec ^
: string
| drop role initial
< ok 0 0
> exec ^
: string
| drop role initial_size
< ok 0 0
> exec ^
: string
| drop role initiate
< ok 0 0
> exec ^
: string
| drop role inner
< error
| 1064 42000 …line 1 column 15 near "inner" 
> exec ^
: string
| drop role inout
< error
| 1064 42000 …line 1 column 15 near "inout" 
> exec ^
: string
| drop role insensitive
< ok 0 0
> exec ^
: string
| drop role insert
< error
| 1064 42000 …line 1 column 16 near "insert" 
> exec ^
: string
| drop role insert_method
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role install
< ok 0 0
> exec ^
: string
| drop role instance
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role int
< error
| 1064 42000 …line 1 column 13 near "int" 
> exec ^
: string
| drop role int1
< error
| 1064 42000 …line 1 column 14 near "int1" 
> exec ^
: string
| drop role int2
< error
| 1064 42000 …line 1 column 14 near "int2" 
> exec ^
: string
| drop role int3
< error
| 1064 42000 …line 1 column 14 near "int3" 
> exec ^
: string
| drop role int4
< error
| 1064 42000 …line 1 column 14 near "int4" 
> exec ^
: string
| drop role int8
< error
| 1064 42000 …line 1 column 14 near "int8" 
> exec ^
: string
| drop role integer
< error
| 1064 42000 …line 1 column 17 near "integer" 
> exec ^
: string
| drop role intersect
< error
| 1064 42000 …line 1 column 19 near "intersect" 
> exec ^
: string
| drop role interval
< error
| 1064 42000 …line 1 column 18 near "interval" 
> exec ^
: string
| drop role into
< error
| 1064 42000 …line 1 column 14 near "into" 
> exec ^
: string
| drop role invisible
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role invoker
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role io
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role io_after_gtids
< ok 0 0
> exec ^
: string
| drop role io_before_gtids
< ok 0 0
> exec ^
: string
| drop role io_thread
< ok 0 0
> exec ^
: string
| drop role ipc
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role is
< error
| 1064 42000 …line 1 column 12 near "is" 
> exec ^
: string
| drop role isolation
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role issuer
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role iterate
< error
| 1064 42000 …line 1 column 17 near "iterate" 
> exec ^
: string
| drop role job
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role jobs
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role join
< error
| 1064 42000 …line 1 column 14 near "join" 
> exec ^
: string
| drop role json
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role json_arrayagg
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role json_duality_object
< ok 0 0
> exec ^
: string
| drop role json_objectagg
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role json_table
< ok 0 0
> exec ^
: string
| drop role json_value
< ok 0 0
> exec ^
: string
| drop role key
< error
| 1064 42000 …line 1 column 13 near "key" 
> exec ^
: string
| drop role key_block_size
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role keyring
< ok 0 0
> exec ^
: string
| drop role keys
< error
| 1064 42000 …line 1 column 14 near "keys" 
> exec ^
: string
| drop role kill
< error
| 1064 42000 …line 1 column 14 near "kill" 
> exec ^
: string
| drop role labels
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role lag
< error
| 1064 42000 …line 1 column 13 near "lag" 
> exec ^
: string
| drop role language
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role last
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role last_value
< error
| 1064 42000 …line 1 column 20 near "last_value" 
> exec ^
: string
| drop role lateral
< ok 0 0
> exec ^
: string
| drop role lead
< error
| 1064 42000 …line 1 column 14 near "lead" 
> exec ^
: string
| drop role leader_constraints
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| drop role leading
< error
| 1064 42000 …line 1 column 17 near "leading" 
> exec ^
: string
| drop role learner_constraints
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| drop role learners
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role least
< ok 0 0
> exec ^
: string
| drop role leave
< error
| 1064 42000 …line 1 column 15 near "leave" 
> exec ^
: string
| drop role leaves
< ok 0 0
> exec ^
: string
| drop role left
< error
| 1064 42000 …line 1 column 14 near "left" 
> exec ^
: string
| drop role less
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role level
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role library
< ok 0 0
> exec ^
: string
| drop role like
< error
| 1064 42000 …line 1 column 14 near "like" 
> exec ^
: string
| drop role limit
< error
| 1064 42000 …line 1 column 15 near "limit" 
> exec ^
: string
| drop role linear
< error
| 1064 42000 …line 1 column 16 near "linear" 
> exec ^
: string
| drop role lines
< error
| 1064 42000 …line 1 column 15 near "lines" 
> exec ^
: string
| drop role linestring
< ok 0 0
> exec ^
: string
| drop role list
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role load
< error
| 1064 42000 …line 1 column 14 near "load" 
> exec ^
: string
| drop role local
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role localtime
< error
| 1064 42000 …line 1 column 19 near "localtime" 
> exec ^
: string
| drop role localtimestamp
< error
| 1064 42000 …line 1 column 24 near "localtimestamp" 
> exec ^
: string
| drop role location
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role lock
< error
| 1064 42000 …line 1 column 14 near "lock" 
> exec ^
: string
| drop role locked
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role locks
< ok 0 0
> exec ^
: string
| drop role log
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role logfile
< ok 0 0
> exec ^
: string
| drop role logs
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role long
< error
| 1064 42000 …line 1 column 14 near "long" 
> exec ^
: string
| drop role longblob
< error
| 1064 42000 …line 1 column 18 near "longblob" 
> exec ^
: string
| drop role longtext
< error
| 1064 42000 …line 1 column 18 near "longtext" 
> exec ^
: string
| drop role loop
< ok 0 0
> exec ^
: string
| drop role low
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role low_priority
< error
| 1064 42000 …line 1 column 22 near "low_priority" 
> exec ^
: string
| drop role manual
< ok 0 0
> exec ^
: string
| drop role master
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role master_log_file
< ok 0 0
> exec ^
: string
| drop role match
< error
| 1064 42000 …line 1 column 15 near "match" 
> exec ^
: string
| drop role materialized
< ok 0 0
> exec ^
: string
| drop role max
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role max_connections_per_hour
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| drop role max_queries_per_hour
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| drop role max_rows
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role max_size
< ok 0 0
> exec ^
: string
| drop role max_updates_per_hour
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| drop role max_user_connections
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| drop role maxvalue
< error
| 1064 42000 …line 1 column 18 near "maxvalue" 
> exec ^
: string
| drop role medium
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role mediumblob
< error
| 1064 42000 …line 1 column 20 near "mediumblob" 
> exec ^
: string
| drop role mediumint
< error
| 1064 42000 …line 1 column 19 near "mediumint" 
> exec ^
: string
| drop role mediumtext
< error
| 1064 42000 …line 1 column 20 near "mediumtext" 
> exec ^
: string
| drop role member
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role memory
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role merge
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role message_text
< ok 0 0
> exec ^
: string
| drop role microsecond
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role mid
< ok 0 0
> exec ^
: string
| drop role middleint
< error
| 1064 42000 …line 1 column 19 near "middleint" 
> exec ^
: string
| drop role migrate
< ok 0 0
> exec ^
: string
| drop role min
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role min_rows
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role minute
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role minute_microsecond
< error
| 1064 42000 …line 1 column 28 near "minute_microsecond" 
> exec ^
: string
| drop role minute_second
< error
| 1064 42000 …line 1 column 23 near "minute_second" 
> exec ^
: string
| drop role minvalue
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role mod
< error
| 1064 42000 …line 1 column 13 near "mod" 
> exec ^
: string
| drop role mode
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role modifies
< ok 0 0
> exec ^
: string
| drop role modify
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role month
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role multilinestring
< ok 0 0
> exec ^
: string
| drop role multipoint
< ok 0 0
> exec ^
: string
| drop role multipolygon
< ok 0 0
> exec ^
: string
| drop role mutex
< ok 0 0
> exec ^
: string
| drop role mysql_errno
< ok 0 0
> exec ^
: string
| drop role name
< ok 0 0
> exec ^
: string
| drop role names
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role national
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role natural
< error
| 1064 42000 …line 1 column 17 near "natural" 
> exec ^
: string
| drop role nchar
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role ndb
< ok 0 0
> exec ^
: string
| drop role ndbcluster
< ok 0 0
> exec ^
: string
| drop role nested
< ok 0 0
> exec ^
: string
| drop role network_namespace
< ok 0 0
> exec ^
: string
| drop role never
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role new
< ok 0 0
> exec ^
: string
| drop role next
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role next_row_id
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role no
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role no_wait
< ok 0 0
> exec ^
: string
| drop role no_write_to_binlog
< error
| 1064 42000 …line 1 column 28 near "no_write_to_binlog" 
> exec ^
: string
| drop role nocache
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role nocycle
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role nodegroup
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role nomaxvalue
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role nominvalue
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role nonclustered
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role none
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role not
< error
| 1064 42000 …line 1 column 13 near "not" 
> exec ^
: string
| drop role now
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role nowait
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role nth_value
< error
| 1064 42000 …line 1 column 19 near "nth_value" 
> exec ^
: string
| drop role ntile
< error
| 1064 42000 …line 1 column 15 near "ntile" 
> exec ^
: string
| drop role null
< error
| 1064 42000 …line 1 column 14 near "null" 
> exec ^
: string
| drop role nullif
< ok 0 0
> exec ^
: string
| drop role nulls
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role number
< ok 0 0
> exec ^
: string
| drop role numeric
< error
| 1064 42000 …line 1 column 17 near "numeric" 
> exec ^
: string
| drop role nvarchar
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role of
< error
| 1064 42000 …line 1 column 12 near "of" 
> exec ^
: string
| drop role off
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role offset
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role oj
< ok 0 0
> exec ^
: string
| drop role old
< ok 0 0
> exec ^
: string
| drop role oltp_read_only
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role oltp_read_write
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| drop role oltp_write_only
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| drop role on
< error
| 1064 42000 …line 1 column 12 near "on" 
> exec ^
: string
| drop role one
< ok 0 0
> exec ^
: string
| drop role only
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role open
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role optimize
< error
| 1064 42000 …line 1 column 18 near "optimize" 
> exec ^
: string
| drop role optimizer_costs
< ok 0 0
> exec ^
: string
| drop role option
< error
| 1064 42000 …line 1 column 16 near "option" 
> exec ^
: string
| drop role optional
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role optionally
< error
| 1064 42000 …line 1 column 20 near "optionally" 
> exec ^
: string
| drop role options
< ok 0 0
> exec ^
: string
| drop role or
< error
| 1064 42000 …line 1 column 12 near "or" 
> exec ^
: string
| drop role order
< error
| 1064 42000 …line 1 column 15 near "order" 
> exec ^
: string
| drop role ordinality
< ok 0 0
> exec ^
: string
| drop role organization
< ok 0 0
> exec ^
: string
| drop role others
< ok 0 0
> exec ^
: string
| drop role out
< error
| 1064 42000 …line 1 column 13 near "out" 
> exec ^
: string
| drop role outer
< error
| 1064 42000 …line 1 column 15 near "outer" 
> exec ^
: string
| drop role outfile
< error
| 1064 42000 …line 1 column 17 near "outfile" 
> exec ^
: string
| drop role over
< error
| 1064 42000 …line 1 column 14 near "over" 
> exec ^
: string
| drop role overlay
< ok 0 0
> exec ^
: string
| drop role owner
< ok 0 0
> exec ^
: string
| drop role pack_keys
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role page
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role parallel
< ok 0 0
> exec ^
: string
| drop role parameters
< ok 0 0
> exec ^
: string
| drop role parse_tree
< ok 0 0
> exec ^
: string
| drop role parser
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role partial
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role partition
< error
| 1064 42000 …line 1 column 19 near "partition" 
> exec ^
: string
| drop role partitioning
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role partitions
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role password
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role password_lock_time
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| drop role path
< ok 0 0
> exec ^
: string
| drop role percent_rank
< error
| 1064 42000 …line 1 column 22 near "percent_rank" 
> exec ^
: string
| drop role persist
< ok 0 0
> exec ^
: string
| drop role persist_only
< ok 0 0
> exec ^
: string
| drop role phase
< ok 0 0
> exec ^
: string
| drop role placement
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role placing
< ok 0 0
> exec ^
: string
| drop role plan
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role plugin
< ok 0 0
> exec ^
: string
| drop role plugin_dir
< ok 0 0
> exec ^
: string
| drop role plugins
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role point
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role policy
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role polygon
< ok 0 0
> exec ^
: string
| drop role port
< ok 0 0
> exec ^
: string
| drop role position
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role pre_split_regions
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role precedes
< ok 0 0
> exec ^
: string
| drop role preceding
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role precision
< error
| 1064 42000 …line 1 column 19 near "precision" 
> exec ^
: string
| drop role prepare
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role preserve
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role prev
< ok 0 0
> exec ^
: string
| drop role primary
< error
| 1064 42000 …line 1 column 17 near "primary" 
> exec ^
: string
| drop role primary_region
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role priority
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role privilege_checks_user
< ok 0 0
> exec ^
: string
| drop role privileges
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role procedure
< error
| 1064 42000 …line 1 column 19 near "procedure" 
> exec ^
: string
| drop role process
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role processed_keys
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role processlist
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role profile
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role profiles
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role proxy
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role purge
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role qualify
< ok 0 0
> exec ^
: string
| drop role quarter
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role query
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role query_limit
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role quick
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role random
< ok 0 0
> exec ^
: string
| drop role range
< error
| 1064 42000 …line 1 column 15 near "range" 
> exec ^
: string
| drop role rank
< error
| 1064 42000 …line 1 column 14 near "rank" 
> exec ^
: string
| drop role read
< error
| 1064 42000 …line 1 column 14 near "read" 
> exec ^
: string
| drop role read_only
< ok 0 0
> exec ^
: string
| drop role read_write
< ok 0 0
> exec ^
: string
| drop role reads
< ok 0 0
> exec ^
: string
| drop role real
< error
| 1064 42000 …line 1 column 14 near "real" 
> exec ^
: string
| drop role rebuild
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role recover
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role recursive
< error
| 1064 42000 …line 1 column 19 near "recursive" 
> exec ^
: string
| drop role redo_buffer_size
< ok 0 0
> exec ^
: string
| drop role redundant
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role reference
< ok 0 0
> exec ^
: string
| drop role references
< error
| 1064 42000 …line 1 column 20 near "references" 
> exec ^
: string
| drop role regexp
< error
| 1064 42000 …line 1 column 16 near "regexp" 
> exec ^
: string
| drop role region
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role regions
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role registration
< ok 0 0
> exec ^
: string
| drop role relational
< ok 0 0
> exec ^
: string
| drop role relay
< ok 0 0
> exec ^
: string
| drop role relay_log_file
< ok 0 0
> exec ^
: string
| drop role relay_log_pos
< ok 0 0
> exec ^
: string
| drop role relay_thread
< ok 0 0
> exec ^
: string
| drop role relaylog
< ok 0 0
> exec ^
: string
| drop role release
< error
| 1064 42000 …line 1 column 17 near "release" 
> exec ^
: string
| drop role reload
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role remove
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role rename
< error
| 1064 42000 …line 1 column 16 near "rename" 
> exec ^
: string
| drop role reorganize
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role repair
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role repeat
< error
| 1064 42000 …line 1 column 16 near "repeat" 
> exec ^
: string
| drop role repeatable
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role replace
< error
| 1064 42000 …line 1 column 17 near "replace" 
> exec ^
: string
| drop role replica
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role replicas
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role replicate_do_db
< ok 0 0
> exec ^
: string
| drop role replicate_do_table
< ok 0 0
> exec ^
: string
| drop role replicate_ignore_db
< ok 0 0
> exec ^
: string
| drop role replicate_ignore_table
< ok 0 0
> exec ^
: string
| drop role replicate_rewrite_db
< ok 0 0
> exec ^
: string
| drop role replicate_wild_do_table
< ok 0 0
> exec ^
: string
| drop role replicate_wild_ignore_table
< ok 0 0
> exec ^
: string
| drop role replication
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role require
< error
| 1064 42000 …line 1 column 17 near "require" 
> exec ^
: string
| drop role require_row_format
< ok 0 0
> exec ^
: string
| drop role require_table_primary_key_check
< ok 0 0
> exec ^
: string
| drop role reset
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role resignal
< ok 0 0
> exec ^
: string
| drop role resource
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role respect
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role restart
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role restore
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role restores
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role restrict
< error
| 1064 42000 …line 1 column 18 near "restrict" 
> exec ^
: string
| drop role resume
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role retain
< ok 0 0
> exec ^
: string
| drop role return
< ok 0 0
> exec ^
: string
| drop role returned_sqlstate
< ok 0 0
> exec ^
: string
| drop role returning
< ok 0 0
> exec ^
: string
| drop role returns
< ok 0 0
> exec ^
: string
| drop role reuse
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role reverse
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role revoke
< error
| 1064 42000 …line 1 column 16 near "revoke" 
> exec ^
: string
| drop role right
< error
| 1064 42000 …line 1 column 15 near "right" 
> exec ^
: string
| drop role rlike
< error
| 1064 42000 …line 1 column 15 near "rlike" 
> exec ^
: string
| drop role role
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role rollback
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role rollup
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role rotate
< ok 0 0
> exec ^
: string
| drop role routine
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role row
< error
| 1064 42000 …line 1 column 13 near "row" 
> exec ^
: string
| drop role row_count
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role row_format
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role row_number
< error
| 1064 42000 …line 1 column 20 near "row_number" 
> exec ^
: string
| drop role rows
< error
| 1064 42000 …line 1 column 14 near "rows" 
> exec ^
: string
| drop role rtree
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role ru
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role ru_per_sec
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role rule
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role run
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role s3
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| drop role savepoint
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role schedule
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role schema
< error
| 1064 42000 …line 1 column 16 near "schema" 
> exec ^
: string
| drop role schema_name
< ok 0 0
> exec ^
: string
| drop role schemas
< error
| 1064 42000 …line 1 column 17 near "schemas" 
> exec ^
: string
| drop role second
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role second_microsecond
< error
| 1064 42000 …line 1 column 28 near "second_microsecond" 
> exec ^
: string
| drop role secondary
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role secondary_engine
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role secondary_engine_attribute
< ok 0 0
> exec ^
: string
| drop role secondary_load
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role secondary_unload
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role security
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role select
< error
| 1064 42000 …line 1 column 16 near "select" 
> exec ^
: string
| drop role sensitive
< ok 0 0
> exec ^
: string
| drop role separator
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role sequence
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role serial
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role serializable
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role server
< ok 0 0
> exec ^
: string
| drop role session
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role session_states
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role session_user
< ok 0 0
> exec ^
: string
| drop role set
< error
| 1064 42000 …line 1 column 13 near "set" 
> exec ^
: string
| drop role sets
< ok 0 0
> exec ^
: string
| drop role shard_row_id_bits
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role share
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role show
< error
| 1064 42000 …line 1 column 14 near "show" 
> exec ^
: string
| drop role shutdown
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role signal
< ok 0 0
> exec ^
: string
| drop role signed
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role similar
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role simple
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role skip
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role slave
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role slow
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role smallint
< error
| 1064 42000 …line 1 column 18 near "smallint" 
> exec ^
: string
| drop role snapshot
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role socket
< ok 0 0
> exec ^
: string
| drop role some
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role soname
< ok 0 0
> exec ^
: string
| drop role sounds
< ok 0 0
> exec ^
: string
| drop role source
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role source_auto_position
< ok 0 0
> exec ^
: string
| drop role source_bind
< ok 0 0
> exec ^
: string
| drop role source_compression_algorithms
< ok 0 0
> exec ^
: string
| drop role source_connect_retry
< ok 0 0
> exec ^
: string
| drop role source_connection_auto_failover
< ok 0 0
> exec ^
: string
| drop role source_delay
< ok 0 0
> exec ^
: string
| drop role source_heartbeat_period
< ok 0 0
> exec ^
: string
| drop role source_host
< ok 0 0
> exec ^
: string
| drop role source_log_file
< ok 0 0
> exec ^
: string
| drop role source_log_pos
< ok 0 0
> exec ^
: string
| drop role source_password
< ok 0 0
> exec ^
: string
| drop role source_port
< ok 0 0
> exec ^
: string
| drop role source_public_key_path
< ok 0 0
> exec ^
: string
| drop role source_retry_count
< ok 0 0
> exec ^
: string
| drop role source_ssl
< ok 0 0
> exec ^
: string
| drop role source_ssl_ca
< ok 0 0
> exec ^
: string
| drop role source_ssl_capath
< ok 0 0
> exec ^
: string
| drop role source_ssl_cert
< ok 0 0
> exec ^
: string
| drop role source_ssl_cipher
< ok 0 0
> exec ^
: string
| drop role source_ssl_crl
< ok 0 0
> exec ^
: string
| drop role source_ssl_crlpath
< ok 0 0
> exec ^
: string
| drop role source_ssl_key
< ok 0 0
> exec ^
: string
| drop role source_ssl_verify_server_cert
< ok 0 0
> exec ^
: string
| drop role source_tls_ciphersuites
< ok 0 0
> exec ^
: string
| drop role source_tls_version
< ok 0 0
> exec ^
: string
| drop role source_user
< ok 0 0
> exec ^
: string
| drop role source_zstd_compression_level
< ok 0 0
> exec ^
: string
| drop role spatial
< error
| 1064 42000 …line 1 column 17 near "spatial" 
> exec ^
: string
| drop role specific
< ok 0 0
> exec ^
: string
| drop role split
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role sql
< error
| 1064 42000 …line 1 column 13 near "sql" 
> exec ^
: string
| drop role sql_after_gtids
< ok 0 0
> exec ^
: string
| drop role sql_after_mts_gaps
< ok 0 0
> exec ^
: string
| drop role sql_before_gtids
< ok 0 0
> exec ^
: string
| drop role sql_big_result
< error
| 1064 42000 …line 1 column 24 near "sql_big_result" 
> exec ^
: string
| drop role sql_buffer_result
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role sql_calc_found_rows
< error
| 1064 42000 …line 1 column 29 near "sql_calc_found_rows" 
> exec ^
: string
| drop role sql_no_cache
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role sql_small_result
< error
| 1064 42000 …line 1 column 26 near "sql_small_result" 
> exec ^
: string
| drop role sql_thread
< ok 0 0
> exec ^
: string
| drop role sql_tsi_day
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role sql_tsi_hour
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role sql_tsi_minute
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role sql_tsi_month
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role sql_tsi_quarter
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| drop role sql_tsi_second
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role sql_tsi_week
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role sql_tsi_year
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role sqlexception
< error
| 1064 42000 …line 1 column 22 near "sqlexception" 
> exec ^
: string
| drop role sqlstate
< error
| 1064 42000 …line 1 column 18 near "sqlstate" 
> exec ^
: string
| drop role sqlwarning
< error
| 1064 42000 …line 1 column 20 near "sqlwarning" 
> exec ^
: string
| drop role srid
< ok 0 0
> exec ^
: string
| drop role ssl
< error
| 1064 42000 …line 1 column 13 near "ssl" 
> exec ^
: string
| drop role st_collect
< ok 0 0
> exec ^
: string
| drop role stacked
< ok 0 0
> exec ^
: string
| drop role start
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role start_time
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role starting
< error
| 1064 42000 …line 1 column 18 near "starting" 
> exec ^
: string
| drop role starts
< ok 0 0
> exec ^
: string
| drop role stats_auto_recalc
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role stats_buckets
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role stats_col_choice
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role stats_col_list
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role stats_extended
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role stats_healthy
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role stats_histograms
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role stats_locked
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role stats_meta
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role stats_persistent
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role stats_sample_pages
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| drop role stats_sample_rate
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role stats_topn
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role status
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role std
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role stddev
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role stddev_pop
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role stddev_samp
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role stop
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role storage
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role stored
< error
| 1064 42000 …line 1 column 16 near "stored" 
> exec ^
: string
| drop role straight_join
< error
| 1064 42000 …line 1 column 23 near "straight_join" 
> exec ^
: string
| drop role stream
< ok 0 0
> exec ^
: string
| drop role strict_load
< ok 0 0
> exec ^
: string
| drop role string
< ok 0 0
> exec ^
: string
| drop role subclass_origin
< ok 0 0
> exec ^
: string
| drop role subdate
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role subject
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role subpartition
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role subpartitions
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role substr
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role substring
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role sum
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role super
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role survival_preferences
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| drop role suspend
< ok 0 0
> exec ^
: string
| drop role swaps
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role switch_group
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role switches
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role sysdate
< ok 0 0
> exec ^
: string
| drop role system
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role system_user
< ok 0 0
> exec ^
: string
| drop role table
< error
| 1064 42000 …line 1 column 15 near "table" 
> exec ^
: string
| drop role table_checksum
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role table_name
< ok 0 0
> exec ^
: string
| drop role tables
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role tablesample
< error
| 1064 42000 …line 1 column 21 near "tablesample" 
> exec ^
: string
| drop role tablespace
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role task_types
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role temporary
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role temptable
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role terminated
< error
| 1064 42000 …line 1 column 20 near "terminated" 
> exec ^
: string
| drop role text
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role than
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role then
< error
| 1064 42000 …line 1 column 14 near "then" 
> exec ^
: string
| drop role thread_priority
< ok 0 0
> exec ^
: string
| drop role ties
< ok 0 0
> exec ^
: string
| drop role tiflash
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role time
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role timeout
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role timestamp
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role timestampadd
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop role timestampdiff
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role tinyblob
< error
| 1064 42000 …line 1 column 18 near "tinyblob" 
> exec ^
: string
| drop role tinyint
< error
| 1064 42000 …line 1 column 17 near "tinyint" 
> exec ^
: string
| drop role tinytext
< error
| 1064 42000 …line 1 column 18 near "tinytext" 
> exec ^
: string
| drop role tls
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role to
< error
| 1064 42000 …line 1 column 12 near "to" 
> exec ^
: string
| drop role tpcc
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role tpch_10
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role trailing
< error
| 1064 42000 …line 1 column 18 near "trailing" 
> exec ^
: string
| drop role transaction
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role trigger
< error
| 1064 42000 …line 1 column 17 near "trigger" 
> exec ^
: string
| drop role triggers
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role trim
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role true
< error
| 1064 42000 …line 1 column 14 near "true" 
> exec ^
: string
| drop role truncate
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role ttl
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| drop role ttl_enable
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role ttl_job_interval
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop role type
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role types
< ok 0 0
> exec ^
: string
| drop role unbounded
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role uncommitted
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| drop role undefined
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role undo
< ok 0 0
> exec ^
: string
| drop role undo_buffer_size
< ok 0 0
> exec ^
: string
| drop role undofile
< ok 0 0
> exec ^
: string
| drop role unicode
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role uninstall
< ok 0 0
> exec ^
: string
| drop role union
< error
| 1064 42000 …line 1 column 15 near "union" 
> exec ^
: string
| drop role unique
< error
| 1064 42000 …line 1 column 16 near "unique" 
> exec ^
: string
| drop role unknown
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role unlimited
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role unlock
< error
| 1064 42000 …line 1 column 16 near "unlock" 
> exec ^
: string
| drop role unregister
< ok 0 0
> exec ^
: string
| drop role unsigned
< error
| 1064 42000 …line 1 column 18 near "unsigned" 
> exec ^
: string
| drop role until
< error
| 1064 42000 …line 1 column 15 near "until" 
> exec ^
: string
| drop role update
< error
| 1064 42000 …line 1 column 16 near "update" 
> exec ^
: string
| drop role upgrade
< ok 0 0
> exec ^
: string
| drop role uri
< ok 0 0
> exec ^
: string
| drop role url
< ok 0 0
> exec ^
: string
| drop role usage
< error
| 1064 42000 …line 1 column 15 near "usage" 
> exec ^
: string
| drop role use
< error
| 1064 42000 …line 1 column 13 near "use" 
> exec ^
: string
| drop role use_frm
< ok 0 0
> exec ^
: string
| drop role user
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role user_resources
< ok 0 0
> exec ^
: string
| drop role using
< error
| 1064 42000 …line 1 column 15 near "using" 
> exec ^
: string
| drop role utc_date
< error
| 1064 42000 …line 1 column 18 near "utc_date" 
> exec ^
: string
| drop role utc_time
< error
| 1064 42000 …line 1 column 18 near "utc_time" 
> exec ^
: string
| drop role utc_timestamp
< error
| 1064 42000 …line 1 column 23 near "utc_timestamp" 
> exec ^
: string
| drop role utilization_limit
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role validate
< ok 0 0
> exec ^
: string
| drop role validation
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| drop role value
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role values
< error
| 1064 42000 …line 1 column 16 near "values" 
> exec ^
: string
| drop role var_pop
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role var_samp
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role varbinary
< error
| 1064 42000 …line 1 column 19 near "varbinary" 
> exec ^
: string
| drop role varchar
< error
| 1064 42000 …line 1 column 17 near "varchar" 
> exec ^
: string
| drop role varcharacter
< error
| 1064 42000 …line 1 column 22 near "varcharacter" 
> exec ^
: string
| drop role variables
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| drop role variance
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role varying
< error
| 1064 42000 …line 1 column 17 near "varying" 
> exec ^
: string
| drop role vcpu
< ok 0 0
> exec ^
: string
| drop role vector
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role verify_key_constraints
< ok 0 0
> exec ^
: string
| drop role view
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role virtual
< error
| 1064 42000 …line 1 column 17 near "virtual" 
> exec ^
: string
| drop role visible
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role voter_constraints
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| drop role voters
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| drop role wait
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role warnings
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role watch
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop role week
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role weight_string
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role when
< error
| 1064 42000 …line 1 column 14 near "when" 
> exec ^
: string
| drop role where
< error
| 1064 42000 …line 1 column 15 near "where" 
> exec ^
: string
| drop role while
< error
| 1064 42000 …line 1 column 15 near "while" 
> exec ^
: string
| drop role window
< error
| 1064 42000 …line 1 column 16 near "window" 
> exec ^
: string
| drop role with
< error
| 1064 42000 …line 1 column 14 near "with" 
> exec ^
: string
| drop role without
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| drop role work
< ok 0 0
> exec ^
: string
| drop role workload
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| drop role wrapper
< ok 0 0
> exec ^
: string
| drop role write
< error
| 1064 42000 …line 1 column 15 near "write" 
> exec ^
: string
| drop role x509
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role xa
< ok 0 0
> exec ^
: string
| drop role xid
< ok 0 0
> exec ^
: string
| drop role xml
< ok 0 0
> exec ^
: string
| drop role xor
< error
| 1064 42000 …line 1 column 13 near "xor" 
> exec ^
: string
| drop role year
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| drop role year_month
< error
| 1064 42000 …line 1 column 20 near "year_month" 
> exec ^
: string
| drop role zerofill
< error
| 1064 42000 …line 1 column 18 near "zerofill" 
> exec ^
: string
| drop role zone
< ok 0 0
> exec ^
: string
| set role absent
< ok 0 0
> exec ^
: string
| set role accessible
< ok 0 0
> exec ^
: string
| set role account
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role action
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role active
< ok 0 0
> exec ^
: string
| set role add
< error
| 1064 42000 …line 1 column 12 near "add" 
> exec ^
: string
| set role adddate
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role admin
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role after
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role against
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role aggregate
< ok 0 0
> exec ^
: string
| set role algorithm
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role all
< ok 0 0
> exec ^
: string
| set role allow_missing_files
< ok 0 0
> exec ^
: string
| set role alter
< error
| 1064 42000 …line 1 column 14 near "alter" 
> exec ^
: string
| set role always
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role analyze
< error
| 1064 42000 …line 1 column 16 near "analyze" 
> exec ^
: string
| set role and
< error
| 1064 42000 …line 1 column 12 near "and" 
> exec ^
: string
| set role any
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role array
< error
| 1064 42000 …line 1 column 14 near "array" 
> exec ^
: string
| set role as
< error
| 1064 42000 …line 1 column 11 near "as" 
> exec ^
: string
| set role asc
< error
| 1064 42000 …line 1 column 12 near "asc" 
> exec ^
: string
| set role ascii
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role asensitive
< ok 0 0
> exec ^
: string
| set role assign_gtids_to_anonymous_transactions
< ok 0 0
> exec ^
: string
| set role at
< ok 0 0
> exec ^
: string
| set role attribute
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role authentication
< ok 0 0
> exec ^
: string
| set role auto
< ok 0 0
> exec ^
: string
| set role auto_id_cache
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role auto_increment
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role auto_random
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role auto_random_base
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role auto_refresh
< ok 0 0
> exec ^
: string
| set role auto_refresh_source
< ok 0 0
> exec ^
: string
| set role autoextend_size
< ok 0 0
> exec ^
: string
| set role avg
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role avg_row_length
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role background
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role backup
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role backups
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role batch
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role before
< ok 0 0
> exec ^
: string
| set role begin
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role bernoulli
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role between
< error
| 1064 42000 …line 1 column 16 near "between" 
> exec ^
: string
| set role bigint
< error
| 1064 42000 …line 1 column 15 near "bigint" 
> exec ^
: string
| set role binary
< error
| 1064 42000 …line 1 column 15 near "binary" 
> exec ^
: string
| set role binding
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role binding_cache
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role bindings
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role binlog
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role bit
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role bit_and
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role bit_or
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role bit_xor
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role blob
< error
| 1064 42000 …line 1 column 13 near "blob" 
> exec ^
: string
| set role block
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role bool
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role boolean
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role both
< error
| 1064 42000 …line 1 column 13 near "both" 
> exec ^
: string
| set role btree
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role buckets
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role builtins
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role bulk
< ok 0 0
> exec ^
: string
| set role burstable
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role by
< error
| 1064 42000 …line 1 column 11 near "by" 
> exec ^
: string
| set role byte
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role cache
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role calibrate
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role call
< error
| 1064 42000 …line 1 column 13 near "call" 
> exec ^
: string
| set role cascade
< error
| 1064 42000 …line 1 column 16 near "cascade" 
> exec ^
: string
| set role cascaded
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role case
< error
| 1064 42000 …line 1 column 13 near "case" 
> exec ^
: string
| set role cast
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role catalog_name
< ok 0 0
> exec ^
: string
| set role chain
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role challenge_response
< ok 0 0
> exec ^
: string
| set role change
< error
| 1064 42000 …line 1 column 15 near "change" 
> exec ^
: string
| set role changed
< ok 0 0
> exec ^
: string
| set role channel
< ok 0 0
> exec ^
: string
| set role char
< error
| 1064 42000 …line 1 column 13 near "char" 
> exec ^
: string
| set role character
< error
| 1064 42000 …line 1 column 18 near "character" 
> exec ^
: string
| set role charset
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role check
< error
| 1064 42000 …line 1 column 14 near "check" 
> exec ^
: string
| set role checksum
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role cipher
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role class_origin
< ok 0 0
> exec ^
: string
| set role client
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role clone
< ok 0 0
> exec ^
: string
| set role close
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role clustered
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role coalesce
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role code
< ok 0 0
> exec ^
: string
| set role collate
< error
| 1064 42000 …line 1 column 16 near "collate" 
> exec ^
: string
| set role collation
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role column
< error
| 1064 42000 …line 1 column 15 near "column" 
> exec ^
: string
| set role column_format
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role column_name
< ok 0 0
> exec ^
: string
| set role column_stats_usage
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| set role columns
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role comment
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role commit
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role committed
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role compact
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role completion
< ok 0 0
> exec ^
: string
| set role component
< ok 0 0
> exec ^
: string
| set role compressed
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role compression
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role concat
< ok 0 0
> exec ^
: string
| set role concurrent
< ok 0 0
> exec ^
: string
| set role condition
< ok 0 0
> exec ^
: string
| set role config
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role connection
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role consistent
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role constraint
< error
| 1064 42000 …line 1 column 19 near "constraint" 
> exec ^
: string
| set role constraint_catalog
< ok 0 0
> exec ^
: string
| set role constraint_name
< ok 0 0
> exec ^
: string
| set role constraint_schema
< ok 0 0
> exec ^
: string
| set role constraints
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role contains
< ok 0 0
> exec ^
: string
| set role context
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role continue
< error
| 1064 42000 …line 1 column 17 near "continue" 
> exec ^
: string
| set role convert
< error
| 1064 42000 …line 1 column 16 near "convert" 
> exec ^
: string
| set role cooldown
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role count
< ok 0 0
> exec ^
: string
| set role cpu
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role create
< error
| 1064 42000 …line 1 column 15 near "create" 
> exec ^
: string
| set role cross
< error
| 1064 42000 …line 1 column 14 near "cross" 
> exec ^
: string
| set role cube
< ok 0 0
> exec ^
: string
| set role cume_dist
< error
| 1064 42000 …line 1 column 18 near "cume_dist" 
> exec ^
: string
| set role curdate
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role current
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role current_date
< error
| 1064 42000 …line 1 column 21 near "current_date" 
> exec ^
: string
| set role current_time
< error
| 1064 42000 …line 1 column 21 near "current_time" 
> exec ^
: string
| set role current_timestamp
< error
| 1064 42000 …line 1 column 26 near "current_timestamp" 
> exec ^
: string
| set role current_user
< error
| 1064 42000 …line 1 column 21 near "current_user" 
> exec ^
: string
| set role cursor
< error
| 1064 42000 …line 1 column 15 near "cursor" 
> exec ^
: string
| set role cursor_name
< ok 0 0
> exec ^
: string
| set role curtime
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role cycle
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role data
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role database
< error
| 1064 42000 …line 1 column 17 near "database" 
> exec ^
: string
| set role databases
< error
| 1064 42000 …line 1 column 18 near "databases" 
> exec ^
: string
| set role datafile
< ok 0 0
> exec ^
: string
| set role date
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role date_add
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role date_sub
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role datetime
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role day
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role day_hour
< error
| 1064 42000 …line 1 column 17 near "day_hour" 
> exec ^
: string
| set role day_microsecond
< error
| 1064 42000 …line 1 column 24 near "day_microsecond" 
> exec ^
: string
| set role day_minute
< error
| 1064 42000 …line 1 column 19 near "day_minute" 
> exec ^
: string
| set role day_second
< error
| 1064 42000 …line 1 column 19 near "day_second" 
> exec ^
: string
| set role deallocate
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role dec
< error
| 1064 42000 …line 1 column 12 near "dec" 
> exec ^
: string
| set role decimal
< error
| 1064 42000 …line 1 column 16 near "decimal" 
> exec ^
: string
| set role declare
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role default
< ok 0 0
> exec ^
: string
| set role default_auth
< ok 0 0
> exec ^
: string
| set role defined
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role definer
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role definition
< ok 0 0
> exec ^
: string
| set role delay_key_write
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| set role delayed
< error
| 1064 42000 …line 1 column 16 near "delayed" 
> exec ^
: string
| set role delete
< error
| 1064 42000 …line 1 column 15 near "delete" 
> exec ^
: string
| set role dense_rank
< error
| 1064 42000 …line 1 column 19 near "dense_rank" 
> exec ^
: string
| set role desc
< error
| 1064 42000 …line 1 column 13 near "desc" 
> exec ^
: string
| set role describe
< error
| 1064 42000 …line 1 column 17 near "describe" 
> exec ^
: string
| set role description
< ok 0 0
> exec ^
: string
| set role deterministic
< ok 0 0
> exec ^
: string
| set role diagnostics
< ok 0 0
> exec ^
: string
| set role digest
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role directory
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role disable
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role disabled
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role discard
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role disk
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role distinct
< error
| 1064 42000 …line 1 column 17 near "distinct" 
> exec ^
: string
| set role distinctrow
< error
| 1064 42000 …line 1 column 20 near "distinctrow" 
> exec ^
: string
| set role distribute
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role distribution
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role distributions
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role div
< error
| 1064 42000 …line 1 column 12 near "div" 
> exec ^
: string
| set role do
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| set role double
< error
| 1064 42000 …line 1 column 15 near "double" 
> exec ^
: string
| set role drop
< error
| 1064 42000 …line 1 column 13 near "drop" 
> exec ^
: string
| set role dry
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role dryrun
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role dual
< error
| 1064 42000 …line 1 column 13 near "dual" 
> exec ^
: string
| set role duality
< ok 0 0
> exec ^
: string
| set role dumpfile
< ok 0 0
> exec ^
: string
| set role duplicate
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role duration
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role dynamic
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role each
< ok 0 0
> exec ^
: string
| set role else
< error
| 1064 42000 …line 1 column 13 near "else" 
> exec ^
: string
| set role elseif
< error
| 1064 42000 …line 1 column 15 near "elseif" 
> exec ^
: string
| set role empty
< ok 0 0
> exec ^
: string
| set role enable
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role enabled
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role enclosed
< error
| 1064 42000 …line 1 column 17 near "enclosed" 
> exec ^
: string
| set role encryption
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role end
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role end_time
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role ends
< ok 0 0
> exec ^
: string
| set role enforced
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role engine
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role engine_attribute
< ok 0 0
> exec ^
: string
| set role engines
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role enum
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role error
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role errors
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role escape
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role escaped
< error
| 1064 42000 …line 1 column 16 near "escaped" 
> exec ^
: string
| set role event
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role events
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role every
< ok 0 0
> exec ^
: string
| set role exact
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role except
< error
| 1064 42000 …line 1 column 15 near "except" 
> exec ^
: string
| set role exchange
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role exclude
< ok 0 0
> exec ^
: string
| set role exec_elapsed
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role execute
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role exists
< error
| 1064 42000 …line 1 column 15 near "exists" 
> exec ^
: string
| set role exit
< error
| 1064 42000 …line 1 column 13 near "exit" 
> exec ^
: string
| set role expansion
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role expire
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role explain
< error
| 1064 42000 …line 1 column 16 near "explain" 
> exec ^
: string
| set role export
< ok 0 0
> exec ^
: string
| set role extended
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role extent_size
< ok 0 0
> exec ^
: string
| set role external
< ok 0 0
> exec ^
: string
| set role external_format
< ok 0 0
> exec ^
: string
| set role extract
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role factor
< ok 0 0
> exec ^
: string
| set role failed_login_attempts
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| set role false
< error
| 1064 42000 …line 1 column 14 near "false" 
> exec ^
: string
| set role fast
< ok 0 0
> exec ^
: string
| set role faults
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role fetch
< error
| 1064 42000 …line 1 column 14 near "fetch" 
> exec ^
: string
| set role fields
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role file
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role file_block_size
< ok 0 0
> exec ^
: string
| set role file_format
< ok 0 0
> exec ^
: string
| set role file_name
< ok 0 0
> exec ^
: string
| set role file_pattern
< ok 0 0
> exec ^
: string
| set role file_prefix
< ok 0 0
> exec ^
: string
| set role files
< ok 0 0
> exec ^
: string
| set role filter
< ok 0 0
> exec ^
: string
| set role finish
< ok 0 0
> exec ^
: string
| set role first
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role first_value
< error
| 1064 42000 …line 1 column 20 near "first_value" 
> exec ^
: string
| set role fixed
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role flashback
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role float
< error
| 1064 42000 …line 1 column 14 near "float" 
> exec ^
: string
| set role float4
< error
| 1064 42000 …line 1 column 15 near "float4" 
> exec ^
: string
| set role float8
< error
| 1064 42000 …line 1 column 15 near "float8" 
> exec ^
: string
| set role flush
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role follower_constraints
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| set role followers
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role following
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role follows
< ok 0 0
> exec ^
: string
| set role for
< error
| 1064 42000 …line 1 column 12 near "for" 
> exec ^
: string
| set role force
< error
| 1064 42000 …line 1 column 14 near "force" 
> exec ^
: string
| set role foreign
< error
| 1064 42000 …line 1 column 16 near "foreign" 
> exec ^
: string
| set role format
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role found
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role from
< error
| 1064 42000 …line 1 column 13 near "from" 
> exec ^
: string
| set role full
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role fulltext
< error
| 1064 42000 …line 1 column 17 near "fulltext" 
> exec ^
: string
| set role function
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role general
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role generate
< ok 0 0
> exec ^
: string
| set role generated
< error
| 1064 42000 …line 1 column 18 near "generated" 
> exec ^
: string
| set role geomcollection
< ok 0 0
> exec ^
: string
| set role geometry
< ok 0 0
> exec ^
: string
| set role geometrycollection
< ok 0 0
> exec ^
: string
| set role get
< ok 0 0
> exec ^
: string
| set role get_format
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role get_source_public_key
< ok 0 0
> exec ^
: string
| set role global
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role grant
< error
| 1064 42000 …line 1 column 14 near "grant" 
> exec ^
: string
| set role grants
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role greatest
< ok 0 0
> exec ^
: string
| set role group
< error
| 1064 42000 …line 1 column 14 near "group" 
> exec ^
: string
| set role group_concat
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role group_replication
< ok 0 0
> exec ^
: string
| set role grouping
< ok 0 0
> exec ^
: string
| set role groups
< error
| 1064 42000 …line 1 column 15 near "groups" 
> exec ^
: string
| set role gtid_only
< ok 0 0
> exec ^
: string
| set role gtids
< ok 0 0
> exec ^
: string
| set role guided
< ok 0 0
> exec ^
: string
| set role handler
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role hash
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role having
< error
| 1064 42000 …line 1 column 15 near "having" 
> exec ^
: string
| set role header
< ok 0 0
> exec ^
: string
| set role help
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role high
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role high_priority
< error
| 1064 42000 …line 1 column 22 near "high_priority" 
> exec ^
: string
| set role histogram
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role histograms_in_flight
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| set role history
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role host
< ok 0 0
> exec ^
: string
| set role hosts
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role hour
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role hour_microsecond
< error
| 1064 42000 …line 1 column 25 near "hour_microsecond" 
> exec ^
: string
| set role hour_minute
< error
| 1064 42000 …line 1 column 20 near "hour_minute" 
> exec ^
: string
| set role hour_second
< error
| 1064 42000 …line 1 column 20 near "hour_second" 
> exec ^
: string
| set role identified
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role if
< error
| 1064 42000 …line 1 column 11 near "if" 
> exec ^
: string
| set role ignore
< error
| 1064 42000 …line 1 column 15 near "ignore" 
> exec ^
: string
| set role ignore_server_ids
< ok 0 0
> exec ^
: string
| set role ilike
< error
| 1064 42000 …line 1 column 14 near "ilike" 
> exec ^
: string
| set role import
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role in
< error
| 1064 42000 …line 1 column 11 near "in" 
> exec ^
: string
| set role inactive
< ok 0 0
> exec ^
: string
| set role increment
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role index
< error
| 1064 42000 …line 1 column 14 near "index" 
> exec ^
: string
| set role indexes
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role infile
< error
| 1064 42000 …line 1 column 15 near "infile" 
> exec ^
: string
| set role initial
< ok 0 0
> exec ^
: string
| set role initial_size
< ok 0 0
> exec ^
: string
| set role initiate
< ok 0 0
> exec ^
: string
| set role inner
< error
| 1064 42000 …line 1 column 14 near "inner" 
> exec ^
: string
| set role inout
< error
| 1064 42000 …line 1 column 14 near "inout" 
> exec ^
: string
| set role insensitive
< ok 0 0
> exec ^
: string
| set role insert
< error
| 1064 42000 …line 1 column 15 near "insert" 
> exec ^
: string
| set role insert_method
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role install
< ok 0 0
> exec ^
: string
| set role instance
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role int
< error
| 1064 42000 …line 1 column 12 near "int" 
> exec ^
: string
| set role int1
< error
| 1064 42000 …line 1 column 13 near "int1" 
> exec ^
: string
| set role int2
< error
| 1064 42000 …line 1 column 13 near "int2" 
> exec ^
: string
| set role int3
< error
| 1064 42000 …line 1 column 13 near "int3" 
> exec ^
: string
| set role int4
< error
| 1064 42000 …line 1 column 13 near "int4" 
> exec ^
: string
| set role int8
< error
| 1064 42000 …line 1 column 13 near "int8" 
> exec ^
: string
| set role integer
< error
| 1064 42000 …line 1 column 16 near "integer" 
> exec ^
: string
| set role intersect
< error
| 1064 42000 …line 1 column 18 near "intersect" 
> exec ^
: string
| set role interval
< error
| 1064 42000 …line 1 column 17 near "interval" 
> exec ^
: string
| set role into
< error
| 1064 42000 …line 1 column 13 near "into" 
> exec ^
: string
| set role invisible
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role invoker
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role io
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| set role io_after_gtids
< ok 0 0
> exec ^
: string
| set role io_before_gtids
< ok 0 0
> exec ^
: string
| set role io_thread
< ok 0 0
> exec ^
: string
| set role ipc
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role is
< error
| 1064 42000 …line 1 column 11 near "is" 
> exec ^
: string
| set role isolation
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role issuer
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role iterate
< error
| 1064 42000 …line 1 column 16 near "iterate" 
> exec ^
: string
| set role job
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role jobs
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role join
< error
| 1064 42000 …line 1 column 13 near "join" 
> exec ^
: string
| set role json
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role json_arrayagg
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role json_duality_object
< ok 0 0
> exec ^
: string
| set role json_objectagg
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role json_table
< ok 0 0
> exec ^
: string
| set role json_value
< ok 0 0
> exec ^
: string
| set role key
< error
| 1064 42000 …line 1 column 12 near "key" 
> exec ^
: string
| set role key_block_size
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role keyring
< ok 0 0
> exec ^
: string
| set role keys
< error
| 1064 42000 …line 1 column 13 near "keys" 
> exec ^
: string
| set role kill
< error
| 1064 42000 …line 1 column 13 near "kill" 
> exec ^
: string
| set role labels
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role lag
< error
| 1064 42000 …line 1 column 12 near "lag" 
> exec ^
: string
| set role language
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role last
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role last_value
< error
| 1064 42000 …line 1 column 19 near "last_value" 
> exec ^
: string
| set role lateral
< ok 0 0
> exec ^
: string
| set role lead
< error
| 1064 42000 …line 1 column 13 near "lead" 
> exec ^
: string
| set role leader_constraints
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| set role leading
< error
| 1064 42000 …line 1 column 16 near "leading" 
> exec ^
: string
| set role learner_constraints
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| set role learners
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role least
< ok 0 0
> exec ^
: string
| set role leave
< error
| 1064 42000 …line 1 column 14 near "leave" 
> exec ^
: string
| set role leaves
< ok 0 0
> exec ^
: string
| set role left
< error
| 1064 42000 …line 1 column 13 near "left" 
> exec ^
: string
| set role less
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role level
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role library
< ok 0 0
> exec ^
: string
| set role like
< error
| 1064 42000 …line 1 column 13 near "like" 
> exec ^
: string
| set role limit
< error
| 1064 42000 …line 1 column 14 near "limit" 
> exec ^
: string
| set role linear
< error
| 1064 42000 …line 1 column 15 near "linear" 
> exec ^
: string
| set role lines
< error
| 1064 42000 …line 1 column 14 near "lines" 
> exec ^
: string
| set role linestring
< ok 0 0
> exec ^
: string
| set role list
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role load
< error
| 1064 42000 …line 1 column 13 near "load" 
> exec ^
: string
| set role local
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role localtime
< error
| 1064 42000 …line 1 column 18 near "localtime" 
> exec ^
: string
| set role localtimestamp
< error
| 1064 42000 …line 1 column 23 near "localtimestamp" 
> exec ^
: string
| set role location
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role lock
< error
| 1064 42000 …line 1 column 13 near "lock" 
> exec ^
: string
| set role locked
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role locks
< ok 0 0
> exec ^
: string
| set role log
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role logfile
< ok 0 0
> exec ^
: string
| set role logs
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role long
< error
| 1064 42000 …line 1 column 13 near "long" 
> exec ^
: string
| set role longblob
< error
| 1064 42000 …line 1 column 17 near "longblob" 
> exec ^
: string
| set role longtext
< error
| 1064 42000 …line 1 column 17 near "longtext" 
> exec ^
: string
| set role loop
< ok 0 0
> exec ^
: string
| set role low
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role low_priority
< error
| 1064 42000 …line 1 column 21 near "low_priority" 
> exec ^
: string
| set role manual
< ok 0 0
> exec ^
: string
| set role master
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role master_log_file
< ok 0 0
> exec ^
: string
| set role match
< error
| 1064 42000 …line 1 column 14 near "match" 
> exec ^
: string
| set role materialized
< ok 0 0
> exec ^
: string
| set role max
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role max_connections_per_hour
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| set role max_queries_per_hour
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| set role max_rows
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role max_size
< ok 0 0
> exec ^
: string
| set role max_updates_per_hour
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| set role max_user_connections
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| set role maxvalue
< error
| 1064 42000 …line 1 column 17 near "maxvalue" 
> exec ^
: string
| set role medium
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role mediumblob
< error
| 1064 42000 …line 1 column 19 near "mediumblob" 
> exec ^
: string
| set role mediumint
< error
| 1064 42000 …line 1 column 18 near "mediumint" 
> exec ^
: string
| set role mediumtext
< error
| 1064 42000 …line 1 column 19 near "mediumtext" 
> exec ^
: string
| set role member
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role memory
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role merge
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role message_text
< ok 0 0
> exec ^
: string
| set role microsecond
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role mid
< ok 0 0
> exec ^
: string
| set role middleint
< error
| 1064 42000 …line 1 column 18 near "middleint" 
> exec ^
: string
| set role migrate
< ok 0 0
> exec ^
: string
| set role min
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role min_rows
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role minute
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role minute_microsecond
< error
| 1064 42000 …line 1 column 27 near "minute_microsecond" 
> exec ^
: string
| set role minute_second
< error
| 1064 42000 …line 1 column 22 near "minute_second" 
> exec ^
: string
| set role minvalue
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role mod
< error
| 1064 42000 …line 1 column 12 near "mod" 
> exec ^
: string
| set role mode
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role modifies
< ok 0 0
> exec ^
: string
| set role modify
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role month
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role multilinestring
< ok 0 0
> exec ^
: string
| set role multipoint
< ok 0 0
> exec ^
: string
| set role multipolygon
< ok 0 0
> exec ^
: string
| set role mutex
< ok 0 0
> exec ^
: string
| set role mysql_errno
< ok 0 0
> exec ^
: string
| set role name
< ok 0 0
> exec ^
: string
| set role names
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role national
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role natural
< error
| 1064 42000 …line 1 column 16 near "natural" 
> exec ^
: string
| set role nchar
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role ndb
< ok 0 0
> exec ^
: string
| set role ndbcluster
< ok 0 0
> exec ^
: string
| set role nested
< ok 0 0
> exec ^
: string
| set role network_namespace
< ok 0 0
> exec ^
: string
| set role never
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role new
< ok 0 0
> exec ^
: string
| set role next
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role next_row_id
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role no
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| set role no_wait
< ok 0 0
> exec ^
: string
| set role no_write_to_binlog
< error
| 1064 42000 …line 1 column 27 near "no_write_to_binlog" 
> exec ^
: string
| set role nocache
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role nocycle
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role nodegroup
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role nomaxvalue
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role nominvalue
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role nonclustered
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role none
< ok 0 0
> exec ^
: string
| set role not
< error
| 1064 42000 …line 1 column 12 near "not" 
> exec ^
: string
| set role now
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role nowait
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role nth_value
< error
| 1064 42000 …line 1 column 18 near "nth_value" 
> exec ^
: string
| set role ntile
< error
| 1064 42000 …line 1 column 14 near "ntile" 
> exec ^
: string
| set role null
< error
| 1064 42000 …line 1 column 13 near "null" 
> exec ^
: string
| set role nullif
< ok 0 0
> exec ^
: string
| set role nulls
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role number
< ok 0 0
> exec ^
: string
| set role numeric
< error
| 1064 42000 …line 1 column 16 near "numeric" 
> exec ^
: string
| set role nvarchar
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role of
< error
| 1064 42000 …line 1 column 11 near "of" 
> exec ^
: string
| set role off
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role offset
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role oj
< ok 0 0
> exec ^
: string
| set role old
< ok 0 0
> exec ^
: string
| set role oltp_read_only
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role oltp_read_write
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| set role oltp_write_only
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| set role on
< error
| 1064 42000 …line 1 column 11 near "on" 
> exec ^
: string
| set role one
< ok 0 0
> exec ^
: string
| set role only
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role open
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role optimize
< error
| 1064 42000 …line 1 column 17 near "optimize" 
> exec ^
: string
| set role optimizer_costs
< ok 0 0
> exec ^
: string
| set role option
< error
| 1064 42000 …line 1 column 15 near "option" 
> exec ^
: string
| set role optional
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role optionally
< error
| 1064 42000 …line 1 column 19 near "optionally" 
> exec ^
: string
| set role options
< ok 0 0
> exec ^
: string
| set role or
< error
| 1064 42000 …line 1 column 11 near "or" 
> exec ^
: string
| set role order
< error
| 1064 42000 …line 1 column 14 near "order" 
> exec ^
: string
| set role ordinality
< ok 0 0
> exec ^
: string
| set role organization
< ok 0 0
> exec ^
: string
| set role others
< ok 0 0
> exec ^
: string
| set role out
< error
| 1064 42000 …line 1 column 12 near "out" 
> exec ^
: string
| set role outer
< error
| 1064 42000 …line 1 column 14 near "outer" 
> exec ^
: string
| set role outfile
< error
| 1064 42000 …line 1 column 16 near "outfile" 
> exec ^
: string
| set role over
< error
| 1064 42000 …line 1 column 13 near "over" 
> exec ^
: string
| set role overlay
< ok 0 0
> exec ^
: string
| set role owner
< ok 0 0
> exec ^
: string
| set role pack_keys
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role page
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role parallel
< ok 0 0
> exec ^
: string
| set role parameters
< ok 0 0
> exec ^
: string
| set role parse_tree
< ok 0 0
> exec ^
: string
| set role parser
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role partial
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role partition
< error
| 1064 42000 …line 1 column 18 near "partition" 
> exec ^
: string
| set role partitioning
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role partitions
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role password
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role password_lock_time
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| set role path
< ok 0 0
> exec ^
: string
| set role percent_rank
< error
| 1064 42000 …line 1 column 21 near "percent_rank" 
> exec ^
: string
| set role persist
< ok 0 0
> exec ^
: string
| set role persist_only
< ok 0 0
> exec ^
: string
| set role phase
< ok 0 0
> exec ^
: string
| set role placement
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role placing
< ok 0 0
> exec ^
: string
| set role plan
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role plugin
< ok 0 0
> exec ^
: string
| set role plugin_dir
< ok 0 0
> exec ^
: string
| set role plugins
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role point
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role policy
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role polygon
< ok 0 0
> exec ^
: string
| set role port
< ok 0 0
> exec ^
: string
| set role position
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role pre_split_regions
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role precedes
< ok 0 0
> exec ^
: string
| set role preceding
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role precision
< error
| 1064 42000 …line 1 column 18 near "precision" 
> exec ^
: string
| set role prepare
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role preserve
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role prev
< ok 0 0
> exec ^
: string
| set role primary
< error
| 1064 42000 …line 1 column 16 near "primary" 
> exec ^
: string
| set role primary_region
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role priority
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role privilege_checks_user
< ok 0 0
> exec ^
: string
| set role privileges
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role procedure
< error
| 1064 42000 …line 1 column 18 near "procedure" 
> exec ^
: string
| set role process
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role processed_keys
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role processlist
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role profile
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role profiles
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role proxy
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role purge
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role qualify
< ok 0 0
> exec ^
: string
| set role quarter
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role query
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role query_limit
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role quick
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role random
< ok 0 0
> exec ^
: string
| set role range
< error
| 1064 42000 …line 1 column 14 near "range" 
> exec ^
: string
| set role rank
< error
| 1064 42000 …line 1 column 13 near "rank" 
> exec ^
: string
| set role read
< error
| 1064 42000 …line 1 column 13 near "read" 
> exec ^
: string
| set role read_only
< ok 0 0
> exec ^
: string
| set role read_write
< ok 0 0
> exec ^
: string
| set role reads
< ok 0 0
> exec ^
: string
| set role real
< error
| 1064 42000 …line 1 column 13 near "real" 
> exec ^
: string
| set role rebuild
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role recover
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role recursive
< error
| 1064 42000 …line 1 column 18 near "recursive" 
> exec ^
: string
| set role redo_buffer_size
< ok 0 0
> exec ^
: string
| set role redundant
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role reference
< ok 0 0
> exec ^
: string
| set role references
< error
| 1064 42000 …line 1 column 19 near "references" 
> exec ^
: string
| set role regexp
< error
| 1064 42000 …line 1 column 15 near "regexp" 
> exec ^
: string
| set role region
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role regions
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role registration
< ok 0 0
> exec ^
: string
| set role relational
< ok 0 0
> exec ^
: string
| set role relay
< ok 0 0
> exec ^
: string
| set role relay_log_file
< ok 0 0
> exec ^
: string
| set role relay_log_pos
< ok 0 0
> exec ^
: string
| set role relay_thread
< ok 0 0
> exec ^
: string
| set role relaylog
< ok 0 0
> exec ^
: string
| set role release
< error
| 1064 42000 …line 1 column 16 near "release" 
> exec ^
: string
| set role reload
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role remove
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role rename
< error
| 1064 42000 …line 1 column 15 near "rename" 
> exec ^
: string
| set role reorganize
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role repair
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role repeat
< error
| 1064 42000 …line 1 column 15 near "repeat" 
> exec ^
: string
| set role repeatable
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role replace
< error
| 1064 42000 …line 1 column 16 near "replace" 
> exec ^
: string
| set role replica
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role replicas
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role replicate_do_db
< ok 0 0
> exec ^
: string
| set role replicate_do_table
< ok 0 0
> exec ^
: string
| set role replicate_ignore_db
< ok 0 0
> exec ^
: string
| set role replicate_ignore_table
< ok 0 0
> exec ^
: string
| set role replicate_rewrite_db
< ok 0 0
> exec ^
: string
| set role replicate_wild_do_table
< ok 0 0
> exec ^
: string
| set role replicate_wild_ignore_table
< ok 0 0
> exec ^
: string
| set role replication
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role require
< error
| 1064 42000 …line 1 column 16 near "require" 
> exec ^
: string
| set role require_row_format
< ok 0 0
> exec ^
: string
| set role require_table_primary_key_check
< ok 0 0
> exec ^
: string
| set role reset
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role resignal
< ok 0 0
> exec ^
: string
| set role resource
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role respect
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role restart
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role restore
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role restores
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role restrict
< error
| 1064 42000 …line 1 column 17 near "restrict" 
> exec ^
: string
| set role resume
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role retain
< ok 0 0
> exec ^
: string
| set role return
< ok 0 0
> exec ^
: string
| set role returned_sqlstate
< ok 0 0
> exec ^
: string
| set role returning
< ok 0 0
> exec ^
: string
| set role returns
< ok 0 0
> exec ^
: string
| set role reuse
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role reverse
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role revoke
< error
| 1064 42000 …line 1 column 15 near "revoke" 
> exec ^
: string
| set role right
< error
| 1064 42000 …line 1 column 14 near "right" 
> exec ^
: string
| set role rlike
< error
| 1064 42000 …line 1 column 14 near "rlike" 
> exec ^
: string
| set role role
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role rollback
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role rollup
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role rotate
< ok 0 0
> exec ^
: string
| set role routine
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role row
< error
| 1064 42000 …line 1 column 12 near "row" 
> exec ^
: string
| set role row_count
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role row_format
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role row_number
< error
| 1064 42000 …line 1 column 19 near "row_number" 
> exec ^
: string
| set role rows
< error
| 1064 42000 …line 1 column 13 near "rows" 
> exec ^
: string
| set role rtree
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role ru
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| set role ru_per_sec
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role rule
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role run
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role s3
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| set role savepoint
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role schedule
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role schema
< error
| 1064 42000 …line 1 column 15 near "schema" 
> exec ^
: string
| set role schema_name
< ok 0 0
> exec ^
: string
| set role schemas
< error
| 1064 42000 …line 1 column 16 near "schemas" 
> exec ^
: string
| set role second
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role second_microsecond
< error
| 1064 42000 …line 1 column 27 near "second_microsecond" 
> exec ^
: string
| set role secondary
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role secondary_engine
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role secondary_engine_attribute
< ok 0 0
> exec ^
: string
| set role secondary_load
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role secondary_unload
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role security
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role select
< error
| 1064 42000 …line 1 column 15 near "select" 
> exec ^
: string
| set role sensitive
< ok 0 0
> exec ^
: string
| set role separator
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role sequence
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role serial
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role serializable
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role server
< ok 0 0
> exec ^
: string
| set role session
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role session_states
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role session_user
< ok 0 0
> exec ^
: string
| set role set
< error
| 1064 42000 …line 1 column 12 near "set" 
> exec ^
: string
| set role sets
< ok 0 0
> exec ^
: string
| set role shard_row_id_bits
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role share
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role show
< error
| 1064 42000 …line 1 column 13 near "show" 
> exec ^
: string
| set role shutdown
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role signal
< ok 0 0
> exec ^
: string
| set role signed
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role similar
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role simple
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role skip
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role slave
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role slow
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role smallint
< error
| 1064 42000 …line 1 column 17 near "smallint" 
> exec ^
: string
| set role snapshot
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role socket
< ok 0 0
> exec ^
: string
| set role some
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role soname
< ok 0 0
> exec ^
: string
| set role sounds
< ok 0 0
> exec ^
: string
| set role source
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role source_auto_position
< ok 0 0
> exec ^
: string
| set role source_bind
< ok 0 0
> exec ^
: string
| set role source_compression_algorithms
< ok 0 0
> exec ^
: string
| set role source_connect_retry
< ok 0 0
> exec ^
: string
| set role source_connection_auto_failover
< ok 0 0
> exec ^
: string
| set role source_delay
< ok 0 0
> exec ^
: string
| set role source_heartbeat_period
< ok 0 0
> exec ^
: string
| set role source_host
< ok 0 0
> exec ^
: string
| set role source_log_file
< ok 0 0
> exec ^
: string
| set role source_log_pos
< ok 0 0
> exec ^
: string
| set role source_password
< ok 0 0
> exec ^
: string
| set role source_port
< ok 0 0
> exec ^
: string
| set role source_public_key_path
< ok 0 0
> exec ^
: string
| set role source_retry_count
< ok 0 0
> exec ^
: string
| set role source_ssl
< ok 0 0
> exec ^
: string
| set role source_ssl_ca
< ok 0 0
> exec ^
: string
| set role source_ssl_capath
< ok 0 0
> exec ^
: string
| set role source_ssl_cert
< ok 0 0
> exec ^
: string
| set role source_ssl_cipher
< ok 0 0
> exec ^
: string
| set role source_ssl_crl
< ok 0 0
> exec ^
: string
| set role source_ssl_crlpath
< ok 0 0
> exec ^
: string
| set role source_ssl_key
< ok 0 0
> exec ^
: string
| set role source_ssl_verify_server_cert
< ok 0 0
> exec ^
: string
| set role source_tls_ciphersuites
< ok 0 0
> exec ^
: string
| set role source_tls_version
< ok 0 0
> exec ^
: string
| set role source_user
< ok 0 0
> exec ^
: string
| set role source_zstd_compression_level
< ok 0 0
> exec ^
: string
| set role spatial
< error
| 1064 42000 …line 1 column 16 near "spatial" 
> exec ^
: string
| set role specific
< ok 0 0
> exec ^
: string
| set role split
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role sql
< error
| 1064 42000 …line 1 column 12 near "sql" 
> exec ^
: string
| set role sql_after_gtids
< ok 0 0
> exec ^
: string
| set role sql_after_mts_gaps
< ok 0 0
> exec ^
: string
| set role sql_before_gtids
< ok 0 0
> exec ^
: string
| set role sql_big_result
< error
| 1064 42000 …line 1 column 23 near "sql_big_result" 
> exec ^
: string
| set role sql_buffer_result
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role sql_calc_found_rows
< error
| 1064 42000 …line 1 column 28 near "sql_calc_found_rows" 
> exec ^
: string
| set role sql_no_cache
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role sql_small_result
< error
| 1064 42000 …line 1 column 25 near "sql_small_result" 
> exec ^
: string
| set role sql_thread
< ok 0 0
> exec ^
: string
| set role sql_tsi_day
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role sql_tsi_hour
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role sql_tsi_minute
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role sql_tsi_month
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role sql_tsi_quarter
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| set role sql_tsi_second
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role sql_tsi_week
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role sql_tsi_year
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role sqlexception
< error
| 1064 42000 …line 1 column 21 near "sqlexception" 
> exec ^
: string
| set role sqlstate
< error
| 1064 42000 …line 1 column 17 near "sqlstate" 
> exec ^
: string
| set role sqlwarning
< error
| 1064 42000 …line 1 column 19 near "sqlwarning" 
> exec ^
: string
| set role srid
< ok 0 0
> exec ^
: string
| set role ssl
< error
| 1064 42000 …line 1 column 12 near "ssl" 
> exec ^
: string
| set role st_collect
< ok 0 0
> exec ^
: string
| set role stacked
< ok 0 0
> exec ^
: string
| set role start
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role start_time
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role starting
< error
| 1064 42000 …line 1 column 17 near "starting" 
> exec ^
: string
| set role starts
< ok 0 0
> exec ^
: string
| set role stats_auto_recalc
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role stats_buckets
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role stats_col_choice
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role stats_col_list
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role stats_extended
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role stats_healthy
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role stats_histograms
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role stats_locked
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role stats_meta
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role stats_persistent
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role stats_sample_pages
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| set role stats_sample_rate
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role stats_topn
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role status
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role std
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role stddev
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role stddev_pop
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role stddev_samp
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role stop
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role storage
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role stored
< error
| 1064 42000 …line 1 column 15 near "stored" 
> exec ^
: string
| set role straight_join
< error
| 1064 42000 …line 1 column 22 near "straight_join" 
> exec ^
: string
| set role stream
< ok 0 0
> exec ^
: string
| set role strict_load
< ok 0 0
> exec ^
: string
| set role string
< ok 0 0
> exec ^
: string
| set role subclass_origin
< ok 0 0
> exec ^
: string
| set role subdate
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role subject
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role subpartition
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role subpartitions
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role substr
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role substring
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role sum
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role super
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role survival_preferences
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| set role suspend
< ok 0 0
> exec ^
: string
| set role swaps
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role switch_group
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role switches
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role sysdate
< ok 0 0
> exec ^
: string
| set role system
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role system_user
< ok 0 0
> exec ^
: string
| set role table
< error
| 1064 42000 …line 1 column 14 near "table" 
> exec ^
: string
| set role table_checksum
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set role table_name
< ok 0 0
> exec ^
: string
| set role tables
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role tablesample
< error
| 1064 42000 …line 1 column 20 near "tablesample" 
> exec ^
: string
| set role tablespace
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role task_types
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role temporary
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role temptable
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role terminated
< error
| 1064 42000 …line 1 column 19 near "terminated" 
> exec ^
: string
| set role text
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role than
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role then
< error
| 1064 42000 …line 1 column 13 near "then" 
> exec ^
: string
| set role thread_priority
< ok 0 0
> exec ^
: string
| set role ties
< ok 0 0
> exec ^
: string
| set role tiflash
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role time
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role timeout
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role timestamp
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role timestampadd
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set role timestampdiff
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role tinyblob
< error
| 1064 42000 …line 1 column 17 near "tinyblob" 
> exec ^
: string
| set role tinyint
< error
| 1064 42000 …line 1 column 16 near "tinyint" 
> exec ^
: string
| set role tinytext
< error
| 1064 42000 …line 1 column 17 near "tinytext" 
> exec ^
: string
| set role tls
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role to
< error
| 1064 42000 …line 1 column 11 near "to" 
> exec ^
: string
| set role tpcc
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role tpch_10
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role trailing
< error
| 1064 42000 …line 1 column 17 near "trailing" 
> exec ^
: string
| set role transaction
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role trigger
< error
| 1064 42000 …line 1 column 16 near "trigger" 
> exec ^
: string
| set role triggers
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role trim
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role true
< error
| 1064 42000 …line 1 column 13 near "true" 
> exec ^
: string
| set role truncate
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role ttl
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set role ttl_enable
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role ttl_job_interval
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set role type
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role types
< ok 0 0
> exec ^
: string
| set role unbounded
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role uncommitted
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role undefined
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role undo
< ok 0 0
> exec ^
: string
| set role undo_buffer_size
< ok 0 0
> exec ^
: string
| set role undofile
< ok 0 0
> exec ^
: string
| set role unicode
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role uninstall
< ok 0 0
> exec ^
: string
| set role union
< error
| 1064 42000 …line 1 column 14 near "union" 
> exec ^
: string
| set role unique
< error
| 1064 42000 …line 1 column 15 near "unique" 
> exec ^
: string
| set role unknown
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role unlimited
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role unlock
< error
| 1064 42000 …line 1 column 15 near "unlock" 
> exec ^
: string
| set role unregister
< ok 0 0
> exec ^
: string
| set role unsigned
< error
| 1064 42000 …line 1 column 17 near "unsigned" 
> exec ^
: string
| set role until
< error
| 1064 42000 …line 1 column 14 near "until" 
> exec ^
: string
| set role update
< error
| 1064 42000 …line 1 column 15 near "update" 
> exec ^
: string
| set role upgrade
< ok 0 0
> exec ^
: string
| set role uri
< ok 0 0
> exec ^
: string
| set role url
< ok 0 0
> exec ^
: string
| set role usage
< error
| 1064 42000 …line 1 column 14 near "usage" 
> exec ^
: string
| set role use
< error
| 1064 42000 …line 1 column 12 near "use" 
> exec ^
: string
| set role use_frm
< ok 0 0
> exec ^
: string
| set role user
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role user_resources
< ok 0 0
> exec ^
: string
| set role using
< error
| 1064 42000 …line 1 column 14 near "using" 
> exec ^
: string
| set role utc_date
< error
| 1064 42000 …line 1 column 17 near "utc_date" 
> exec ^
: string
| set role utc_time
< error
| 1064 42000 …line 1 column 17 near "utc_time" 
> exec ^
: string
| set role utc_timestamp
< error
| 1064 42000 …line 1 column 22 near "utc_timestamp" 
> exec ^
: string
| set role utilization_limit
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role validate
< ok 0 0
> exec ^
: string
| set role validation
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set role value
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role values
< error
| 1064 42000 …line 1 column 15 near "values" 
> exec ^
: string
| set role var_pop
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role var_samp
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role varbinary
< error
| 1064 42000 …line 1 column 18 near "varbinary" 
> exec ^
: string
| set role varchar
< error
| 1064 42000 …line 1 column 16 near "varchar" 
> exec ^
: string
| set role varcharacter
< error
| 1064 42000 …line 1 column 21 near "varcharacter" 
> exec ^
: string
| set role variables
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role variance
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role varying
< error
| 1064 42000 …line 1 column 16 near "varying" 
> exec ^
: string
| set role vcpu
< ok 0 0
> exec ^
: string
| set role vector
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role verify_key_constraints
< ok 0 0
> exec ^
: string
| set role view
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role virtual
< error
| 1064 42000 …line 1 column 16 near "virtual" 
> exec ^
: string
| set role visible
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role voter_constraints
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set role voters
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set role wait
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role warnings
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role watch
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| set role week
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role weight_string
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set role when
< error
| 1064 42000 …line 1 column 13 near "when" 
> exec ^
: string
| set role where
< error
| 1064 42000 …line 1 column 14 near "where" 
> exec ^
: string
| set role while
< error
| 1064 42000 …line 1 column 14 near "while" 
> exec ^
: string
| set role window
< error
| 1064 42000 …line 1 column 15 near "window" 
> exec ^
: string
| set role with
< error
| 1064 42000 …line 1 column 13 near "with" 
> exec ^
: string
| set role without
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set role work
< ok 0 0
> exec ^
: string
| set role workload
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role wrapper
< ok 0 0
> exec ^
: string
| set role write
< error
| 1064 42000 …line 1 column 14 near "write" 
> exec ^
: string
| set role x509
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role xa
< ok 0 0
> exec ^
: string
| set role xid
< ok 0 0
> exec ^
: string
| set role xml
< ok 0 0
> exec ^
: string
| set role xor
< error
| 1064 42000 …line 1 column 12 near "xor" 
> exec ^
: string
| set role year
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| set role year_month
< error
| 1064 42000 …line 1 column 19 near "year_month" 
> exec ^
: string
| set role zerofill
< error
| 1064 42000 …line 1 column 17 near "zerofill" 
> exec ^
: string
| set role zone
< ok 0 0
> exec ^
: string
| set default role absent to u
< ok 0 0
> exec ^
: string
| set default role accessible to u
< ok 0 0
> exec ^
: string
| set default role account to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role action to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role active to u
< ok 0 0
> exec ^
: string
| set default role add to u
< error
| 1064 42000 …line 1 column 20 near "add to u" 
> exec ^
: string
| set default role adddate to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role admin to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role after to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role against to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role aggregate to u
< ok 0 0
> exec ^
: string
| set default role algorithm to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role all to u
< ok 0 0
> exec ^
: string
| set default role allow_missing_files to u
< ok 0 0
> exec ^
: string
| set default role alter to u
< error
| 1064 42000 …line 1 column 22 near "alter to u" 
> exec ^
: string
| set default role always to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role analyze to u
< error
| 1064 42000 …line 1 column 24 near "analyze to u" 
> exec ^
: string
| set default role and to u
< error
| 1064 42000 …line 1 column 20 near "and to u" 
> exec ^
: string
| set default role any to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role array to u
< error
| 1064 42000 …line 1 column 22 near "array to u" 
> exec ^
: string
| set default role as to u
< error
| 1064 42000 …line 1 column 19 near "as to u" 
> exec ^
: string
| set default role asc to u
< error
| 1064 42000 …line 1 column 20 near "asc to u" 
> exec ^
: string
| set default role ascii to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role asensitive to u
< ok 0 0
> exec ^
: string
| set default role assign_gtids_to_anonymous_transactions to u
< ok 0 0
> exec ^
: string
| set default role at to u
< ok 0 0
> exec ^
: string
| set default role attribute to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role authentication to u
< ok 0 0
> exec ^
: string
| set default role auto to u
< ok 0 0
> exec ^
: string
| set default role auto_id_cache to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role auto_increment to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role auto_random to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role auto_random_base to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role auto_refresh to u
< ok 0 0
> exec ^
: string
| set default role auto_refresh_source to u
< ok 0 0
> exec ^
: string
| set default role autoextend_size to u
< ok 0 0
> exec ^
: string
| set default role avg to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role avg_row_length to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role background to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role backup to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role backups to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role batch to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role before to u
< ok 0 0
> exec ^
: string
| set default role begin to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role bernoulli to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role between to u
< error
| 1064 42000 …line 1 column 24 near "between to u" 
> exec ^
: string
| set default role bigint to u
< error
| 1064 42000 …line 1 column 23 near "bigint to u" 
> exec ^
: string
| set default role binary to u
< error
| 1064 42000 …line 1 column 23 near "binary to u" 
> exec ^
: string
| set default role binding to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role binding_cache to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role bindings to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role binlog to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role bit to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role bit_and to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role bit_or to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role bit_xor to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role blob to u
< error
| 1064 42000 …line 1 column 21 near "blob to u" 
> exec ^
: string
| set default role block to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role bool to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role boolean to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role both to u
< error
| 1064 42000 …line 1 column 21 near "both to u" 
> exec ^
: string
| set default role btree to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role buckets to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role builtins to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role bulk to u
< ok 0 0
> exec ^
: string
| set default role burstable to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role by to u
< error
| 1064 42000 …line 1 column 19 near "by to u" 
> exec ^
: string
| set default role byte to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role cache to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role calibrate to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role call to u
< error
| 1064 42000 …line 1 column 21 near "call to u" 
> exec ^
: string
| set default role cascade to u
< error
| 1064 42000 …line 1 column 24 near "cascade to u" 
> exec ^
: string
| set default role cascaded to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role case to u
< error
| 1064 42000 …line 1 column 21 near "case to u" 
> exec ^
: string
| set default role cast to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role catalog_name to u
< ok 0 0
> exec ^
: string
| set default role chain to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role challenge_response to u
< ok 0 0
> exec ^
: string
| set default role change to u
< error
| 1064 42000 …line 1 column 23 near "change to u" 
> exec ^
: string
| set default role changed to u
< ok 0 0
> exec ^
: string
| set default role channel to u
< ok 0 0
> exec ^
: string
| set default role char to u
< error
| 1064 42000 …line 1 column 21 near "char to u" 
> exec ^
: string
| set default role character to u
< error
| 1064 42000 …line 1 column 26 near "character to u" 
> exec ^
: string
| set default role charset to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role check to u
< error
| 1064 42000 …line 1 column 22 near "check to u" 
> exec ^
: string
| set default role checksum to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role cipher to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role class_origin to u
< ok 0 0
> exec ^
: string
| set default role client to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role clone to u
< ok 0 0
> exec ^
: string
| set default role close to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role clustered to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role coalesce to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role code to u
< ok 0 0
> exec ^
: string
| set default role collate to u
< error
| 1064 42000 …line 1 column 24 near "collate to u" 
> exec ^
: string
| set default role collation to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role column to u
< error
| 1064 42000 …line 1 column 23 near "column to u" 
> exec ^
: string
| set default role column_format to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role column_name to u
< ok 0 0
> exec ^
: string
| set default role column_stats_usage to u
< error
| 1064 42000 …line 1 column 38 near "to u" 
> exec ^
: string
| set default role columns to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role comment to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role commit to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role committed to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role compact to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role completion to u
< ok 0 0
> exec ^
: string
| set default role component to u
< ok 0 0
> exec ^
: string
| set default role compressed to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role compression to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role concat to u
< ok 0 0
> exec ^
: string
| set default role concurrent to u
< ok 0 0
> exec ^
: string
| set default role condition to u
< ok 0 0
> exec ^
: string
| set default role config to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role connection to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role consistent to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role constraint to u
< error
| 1064 42000 …line 1 column 27 near "constraint to u" 
> exec ^
: string
| set default role constraint_catalog to u
< ok 0 0
> exec ^
: string
| set default role constraint_name to u
< ok 0 0
> exec ^
: string
| set default role constraint_schema to u
< ok 0 0
> exec ^
: string
| set default role constraints to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role contains to u
< ok 0 0
> exec ^
: string
| set default role context to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role continue to u
< error
| 1064 42000 …line 1 column 25 near "continue to u" 
> exec ^
: string
| set default role convert to u
< error
| 1064 42000 …line 1 column 24 near "convert to u" 
> exec ^
: string
| set default role cooldown to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role count to u
< ok 0 0
> exec ^
: string
| set default role cpu to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role create to u
< error
| 1064 42000 …line 1 column 23 near "create to u" 
> exec ^
: string
| set default role cross to u
< error
| 1064 42000 …line 1 column 22 near "cross to u" 
> exec ^
: string
| set default role cube to u
< ok 0 0
> exec ^
: string
| set default role cume_dist to u
< error
| 1064 42000 …line 1 column 26 near "cume_dist to u" 
> exec ^
: string
| set default role curdate to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role current to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role current_date to u
< error
| 1064 42000 …line 1 column 29 near "current_date to u" 
> exec ^
: string
| set default role current_time to u
< error
| 1064 42000 …line 1 column 29 near "current_time to u" 
> exec ^
: string
| set default role current_timestamp to u
< error
| 1064 42000 …line 1 column 34 near "current_timestamp to u" 
> exec ^
: string
| set default role current_user to u
< error
| 1064 42000 …line 1 column 29 near "current_user to u" 
> exec ^
: string
| set default role cursor to u
< error
| 1064 42000 …line 1 column 23 near "cursor to u" 
> exec ^
: string
| set default role cursor_name to u
< ok 0 0
> exec ^
: string
| set default role curtime to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role cycle to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role data to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role database to u
< error
| 1064 42000 …line 1 column 25 near "database to u" 
> exec ^
: string
| set default role databases to u
< error
| 1064 42000 …line 1 column 26 near "databases to u" 
> exec ^
: string
| set default role datafile to u
< ok 0 0
> exec ^
: string
| set default role date to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role date_add to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role date_sub to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role datetime to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role day to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role day_hour to u
< error
| 1064 42000 …line 1 column 25 near "day_hour to u" 
> exec ^
: string
| set default role day_microsecond to u
< error
| 1064 42000 …line 1 column 32 near "day_microsecond to u" 
> exec ^
: string
| set default role day_minute to u
< error
| 1064 42000 …line 1 column 27 near "day_minute to u" 
> exec ^
: string
| set default role day_second to u
< error
| 1064 42000 …line 1 column 27 near "day_second to u" 
> exec ^
: string
| set default role deallocate to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role dec to u
< error
| 1064 42000 …line 1 column 20 near "dec to u" 
> exec ^
: string
| set default role decimal to u
< error
| 1064 42000 …line 1 column 24 near "decimal to u" 
> exec ^
: string
| set default role declare to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role default to u
< error
| 1064 42000 …line 1 column 24 near "default to u" 
> exec ^
: string
| set default role default_auth to u
< ok 0 0
> exec ^
: string
| set default role defined to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role definer to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role definition to u
< ok 0 0
> exec ^
: string
| set default role delay_key_write to u
< error
| 1064 42000 …line 1 column 35 near "to u" 
> exec ^
: string
| set default role delayed to u
< error
| 1064 42000 …line 1 column 24 near "delayed to u" 
> exec ^
: string
| set default role delete to u
< error
| 1064 42000 …line 1 column 23 near "delete to u" 
> exec ^
: string
| set default role dense_rank to u
< error
| 1064 42000 …line 1 column 27 near "dense_rank to u" 
> exec ^
: string
| set default role desc to u
< error
| 1064 42000 …line 1 column 21 near "desc to u" 
> exec ^
: string
| set default role describe to u
< error
| 1064 42000 …line 1 column 25 near "describe to u" 
> exec ^
: string
| set default role description to u
< ok 0 0
> exec ^
: string
| set default role deterministic to u
< ok 0 0
> exec ^
: string
| set default role diagnostics to u
< ok 0 0
> exec ^
: string
| set default role digest to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role directory to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role disable to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role disabled to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role discard to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role disk to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role distinct to u
< error
| 1064 42000 …line 1 column 25 near "distinct to u" 
> exec ^
: string
| set default role distinctrow to u
< error
| 1064 42000 …line 1 column 28 near "distinctrow to u" 
> exec ^
: string
| set default role distribute to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role distribution to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role distributions to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role div to u
< error
| 1064 42000 …line 1 column 20 near "div to u" 
> exec ^
: string
| set default role do to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| set default role double to u
< error
| 1064 42000 …line 1 column 23 near "double to u" 
> exec ^
: string
| set default role drop to u
< error
| 1064 42000 …line 1 column 21 near "drop to u" 
> exec ^
: string
| set default role dry to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role dryrun to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role dual to u
< error
| 1064 42000 …line 1 column 21 near "dual to u" 
> exec ^
: string
| set default role duality to u
< ok 0 0
> exec ^
: string
| set default role dumpfile to u
< ok 0 0
> exec ^
: string
| set default role duplicate to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role duration to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role dynamic to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role each to u
< ok 0 0
> exec ^
: string
| set default role else to u
< error
| 1064 42000 …line 1 column 21 near "else to u" 
> exec ^
: string
| set default role elseif to u
< error
| 1064 42000 …line 1 column 23 near "elseif to u" 
> exec ^
: string
| set default role empty to u
< ok 0 0
> exec ^
: string
| set default role enable to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role enabled to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role enclosed to u
< error
| 1064 42000 …line 1 column 25 near "enclosed to u" 
> exec ^
: string
| set default role encryption to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role end to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role end_time to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role ends to u
< ok 0 0
> exec ^
: string
| set default role enforced to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role engine to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role engine_attribute to u
< ok 0 0
> exec ^
: string
| set default role engines to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role enum to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role error to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role errors to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role escape to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role escaped to u
< error
| 1064 42000 …line 1 column 24 near "escaped to u" 
> exec ^
: string
| set default role event to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role events to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role every to u
< ok 0 0
> exec ^
: string
| set default role exact to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role except to u
< error
| 1064 42000 …line 1 column 23 near "except to u" 
> exec ^
: string
| set default role exchange to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role exclude to u
< ok 0 0
> exec ^
: string
| set default role exec_elapsed to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role execute to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role exists to u
< error
| 1064 42000 …line 1 column 23 near "exists to u" 
> exec ^
: string
| set default role exit to u
< error
| 1064 42000 …line 1 column 21 near "exit to u" 
> exec ^
: string
| set default role expansion to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role expire to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role explain to u
< error
| 1064 42000 …line 1 column 24 near "explain to u" 
> exec ^
: string
| set default role export to u
< ok 0 0
> exec ^
: string
| set default role extended to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role extent_size to u
< ok 0 0
> exec ^
: string
| set default role external to u
< ok 0 0
> exec ^
: string
| set default role external_format to u
< ok 0 0
> exec ^
: string
| set default role extract to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role factor to u
< ok 0 0
> exec ^
: string
| set default role failed_login_attempts to u
< error
| 1064 42000 …line 1 column 41 near "to u" 
> exec ^
: string
| set default role false to u
< error
| 1064 42000 …line 1 column 22 near "false to u" 
> exec ^
: string
| set default role fast to u
< ok 0 0
> exec ^
: string
| set default role faults to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role fetch to u
< error
| 1064 42000 …line 1 column 22 near "fetch to u" 
> exec ^
: string
| set default role fields to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role file to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role file_block_size to u
< ok 0 0
> exec ^
: string
| set default role file_format to u
< ok 0 0
> exec ^
: string
| set default role file_name to u
< ok 0 0
> exec ^
: string
| set default role file_pattern to u
< ok 0 0
> exec ^
: string
| set default role file_prefix to u
< ok 0 0
> exec ^
: string
| set default role files to u
< ok 0 0
> exec ^
: string
| set default role filter to u
< ok 0 0
> exec ^
: string
| set default role finish to u
< ok 0 0
> exec ^
: string
| set default role first to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role first_value to u
< error
| 1064 42000 …line 1 column 28 near "first_value to u" 
> exec ^
: string
| set default role fixed to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role flashback to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role float to u
< error
| 1064 42000 …line 1 column 22 near "float to u" 
> exec ^
: string
| set default role float4 to u
< error
| 1064 42000 …line 1 column 23 near "float4 to u" 
> exec ^
: string
| set default role float8 to u
< error
| 1064 42000 …line 1 column 23 near "float8 to u" 
> exec ^
: string
| set default role flush to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role follower_constraints to u
< error
| 1064 42000 …line 1 column 40 near "to u" 
> exec ^
: string
| set default role followers to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role following to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role follows to u
< ok 0 0
> exec ^
: string
| set default role for to u
< error
| 1064 42000 …line 1 column 20 near "for to u" 
> exec ^
: string
| set default role force to u
< error
| 1064 42000 …line 1 column 22 near "force to u" 
> exec ^
: string
| set default role foreign to u
< error
| 1064 42000 …line 1 column 24 near "foreign to u" 
> exec ^
: string
| set default role format to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role found to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role from to u
< error
| 1064 42000 …line 1 column 21 near "from to u" 
> exec ^
: string
| set default role full to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role fulltext to u
< error
| 1064 42000 …line 1 column 25 near "fulltext to u" 
> exec ^
: string
| set default role function to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role general to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role generate to u
< ok 0 0
> exec ^
: string
| set default role generated to u
< error
| 1064 42000 …line 1 column 26 near "generated to u" 
> exec ^
: string
| set default role geomcollection to u
< ok 0 0
> exec ^
: string
| set default role geometry to u
< ok 0 0
> exec ^
: string
| set default role geometrycollection to u
< ok 0 0
> exec ^
: string
| set default role get to u
< ok 0 0
> exec ^
: string
| set default role get_format to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role get_source_public_key to u
< ok 0 0
> exec ^
: string
| set default role global to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role grant to u
< error
| 1064 42000 …line 1 column 22 near "grant to u" 
> exec ^
: string
| set default role grants to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role greatest to u
< ok 0 0
> exec ^
: string
| set default role group to u
< error
| 1064 42000 …line 1 column 22 near "group to u" 
> exec ^
: string
| set default role group_concat to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role group_replication to u
< ok 0 0
> exec ^
: string
| set default role grouping to u
< ok 0 0
> exec ^
: string
| set default role groups to u
< error
| 1064 42000 …line 1 column 23 near "groups to u" 
> exec ^
: string
| set default role gtid_only to u
< ok 0 0
> exec ^
: string
| set default role gtids to u
< ok 0 0
> exec ^
: string
| set default role guided to u
< ok 0 0
> exec ^
: string
| set default role handler to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role hash to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role having to u
< error
| 1064 42000 …line 1 column 23 near "having to u" 
> exec ^
: string
| set default role header to u
< ok 0 0
> exec ^
: string
| set default role help to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role high to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role high_priority to u
< error
| 1064 42000 …line 1 column 30 near "high_priority to u" 
> exec ^
: string
| set default role histogram to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role histograms_in_flight to u
< error
| 1064 42000 …line 1 column 40 near "to u" 
> exec ^
: string
| set default role history to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role host to u
< ok 0 0
> exec ^
: string
| set default role hosts to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role hour to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role hour_microsecond to u
< error
| 1064 42000 …line 1 column 33 near "hour_microsecond to u" 
> exec ^
: string
| set default role hour_minute to u
< error
| 1064 42000 …line 1 column 28 near "hour_minute to u" 
> exec ^
: string
| set default role hour_second to u
< error
| 1064 42000 …line 1 column 28 near "hour_second to u" 
> exec ^
: string
| set default role identified to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role if to u
< error
| 1064 42000 …line 1 column 19 near "if to u" 
> exec ^
: string
| set default role ignore to u
< error
| 1064 42000 …line 1 column 23 near "ignore to u" 
> exec ^
: string
| set default role ignore_server_ids to u
< ok 0 0
> exec ^
: string
| set default role ilike to u
< error
| 1064 42000 …line 1 column 22 near "ilike to u" 
> exec ^
: string
| set default role import to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role in to u
< error
| 1064 42000 …line 1 column 19 near "in to u" 
> exec ^
: string
| set default role inactive to u
< ok 0 0
> exec ^
: string
| set default role increment to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role index to u
< error
| 1064 42000 …line 1 column 22 near "index to u" 
> exec ^
: string
| set default role indexes to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role infile to u
< error
| 1064 42000 …line 1 column 23 near "infile to u" 
> exec ^
: string
| set default role initial to u
< ok 0 0
> exec ^
: string
| set default role initial_size to u
< ok 0 0
> exec ^
: string
| set default role initiate to u
< ok 0 0
> exec ^
: string
| set default role inner to u
< error
| 1064 42000 …line 1 column 22 near "inner to u" 
> exec ^
: string
| set default role inout to u
< error
| 1064 42000 …line 1 column 22 near "inout to u" 
> exec ^
: string
| set default role insensitive to u
< ok 0 0
> exec ^
: string
| set default role insert to u
< error
| 1064 42000 …line 1 column 23 near "insert to u" 
> exec ^
: string
| set default role insert_method to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role install to u
< ok 0 0
> exec ^
: string
| set default role instance to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role int to u
< error
| 1064 42000 …line 1 column 20 near "int to u" 
> exec ^
: string
| set default role int1 to u
< error
| 1064 42000 …line 1 column 21 near "int1 to u" 
> exec ^
: string
| set default role int2 to u
< error
| 1064 42000 …line 1 column 21 near "int2 to u" 
> exec ^
: string
| set default role int3 to u
< error
| 1064 42000 …line 1 column 21 near "int3 to u" 
> exec ^
: string
| set default role int4 to u
< error
| 1064 42000 …line 1 column 21 near "int4 to u" 
> exec ^
: string
| set default role int8 to u
< error
| 1064 42000 …line 1 column 21 near "int8 to u" 
> exec ^
: string
| set default role integer to u
< error
| 1064 42000 …line 1 column 24 near "integer to u" 
> exec ^
: string
| set default role intersect to u
< error
| 1064 42000 …line 1 column 26 near "intersect to u" 
> exec ^
: string
| set default role interval to u
< error
| 1064 42000 …line 1 column 25 near "interval to u" 
> exec ^
: string
| set default role into to u
< error
| 1064 42000 …line 1 column 21 near "into to u" 
> exec ^
: string
| set default role invisible to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role invoker to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role io to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| set default role io_after_gtids to u
< ok 0 0
> exec ^
: string
| set default role io_before_gtids to u
< ok 0 0
> exec ^
: string
| set default role io_thread to u
< ok 0 0
> exec ^
: string
| set default role ipc to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role is to u
< error
| 1064 42000 …line 1 column 19 near "is to u" 
> exec ^
: string
| set default role isolation to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role issuer to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role iterate to u
< error
| 1064 42000 …line 1 column 24 near "iterate to u" 
> exec ^
: string
| set default role job to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role jobs to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role join to u
< error
| 1064 42000 …line 1 column 21 near "join to u" 
> exec ^
: string
| set default role json to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role json_arrayagg to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role json_duality_object to u
< ok 0 0
> exec ^
: string
| set default role json_objectagg to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role json_table to u
< ok 0 0
> exec ^
: string
| set default role json_value to u
< ok 0 0
> exec ^
: string
| set default role key to u
< error
| 1064 42000 …line 1 column 20 near "key to u" 
> exec ^
: string
| set default role key_block_size to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role keyring to u
< ok 0 0
> exec ^
: string
| set default role keys to u
< error
| 1064 42000 …line 1 column 21 near "keys to u" 
> exec ^
: string
| set default role kill to u
< error
| 1064 42000 …line 1 column 21 near "kill to u" 
> exec ^
: string
| set default role labels to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role lag to u
< error
| 1064 42000 …line 1 column 20 near "lag to u" 
> exec ^
: string
| set default role language to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role last to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role last_value to u
< error
| 1064 42000 …line 1 column 27 near "last_value to u" 
> exec ^
: string
| set default role lateral to u
< ok 0 0
> exec ^
: string
| set default role lead to u
< error
| 1064 42000 …line 1 column 21 near "lead to u" 
> exec ^
: string
| set default role leader_constraints to u
< error
| 1064 42000 …line 1 column 38 near "to u" 
> exec ^
: string
| set default role leading to u
< error
| 1064 42000 …line 1 column 24 near "leading to u" 
> exec ^
: string
| set default role learner_constraints to u
< error
| 1064 42000 …line 1 column 39 near "to u" 
> exec ^
: string
| set default role learners to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role least to u
< ok 0 0
> exec ^
: string
| set default role leave to u
< error
| 1064 42000 …line 1 column 22 near "leave to u" 
> exec ^
: string
| set default role leaves to u
< ok 0 0
> exec ^
: string
| set default role left to u
< error
| 1064 42000 …line 1 column 21 near "left to u" 
> exec ^
: string
| set default role less to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role level to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role library to u
< ok 0 0
> exec ^
: string
| set default role like to u
< error
| 1064 42000 …line 1 column 21 near "like to u" 
> exec ^
: string
| set default role limit to u
< error
| 1064 42000 …line 1 column 22 near "limit to u" 
> exec ^
: string
| set default role linear to u
< error
| 1064 42000 …line 1 column 23 near "linear to u" 
> exec ^
: string
| set default role lines to u
< error
| 1064 42000 …line 1 column 22 near "lines to u" 
> exec ^
: string
| set default role linestring to u
< ok 0 0
> exec ^
: string
| set default role list to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role load to u
< error
| 1064 42000 …line 1 column 21 near "load to u" 
> exec ^
: string
| set default role local to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role localtime to u
< error
| 1064 42000 …line 1 column 26 near "localtime to u" 
> exec ^
: string
| set default role localtimestamp to u
< error
| 1064 42000 …line 1 column 31 near "localtimestamp to u" 
> exec ^
: string
| set default role location to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role lock to u
< error
| 1064 42000 …line 1 column 21 near "lock to u" 
> exec ^
: string
| set default role locked to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role locks to u
< ok 0 0
> exec ^
: string
| set default role log to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role logfile to u
< ok 0 0
> exec ^
: string
| set default role logs to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role long to u
< error
| 1064 42000 …line 1 column 21 near "long to u" 
> exec ^
: string
| set default role longblob to u
< error
| 1064 42000 …line 1 column 25 near "longblob to u" 
> exec ^
: string
| set default role longtext to u
< error
| 1064 42000 …line 1 column 25 near "longtext to u" 
> exec ^
: string
| set default role loop to u
< ok 0 0
> exec ^
: string
| set default role low to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role low_priority to u
< error
| 1064 42000 …line 1 column 29 near "low_priority to u" 
> exec ^
: string
| set default role manual to u
< ok 0 0
> exec ^
: string
| set default role master to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role master_log_file to u
< ok 0 0
> exec ^
: string
| set default role match to u
< error
| 1064 42000 …line 1 column 22 near "match to u" 
> exec ^
: string
| set default role materialized to u
< ok 0 0
> exec ^
: string
| set default role max to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role max_connections_per_hour to u
< error
| 1064 42000 …line 1 column 44 near "to u" 
> exec ^
: string
| set default role max_queries_per_hour to u
< error
| 1064 42000 …line 1 column 40 near "to u" 
> exec ^
: string
| set default role max_rows to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role max_size to u
< ok 0 0
> exec ^
: string
| set default role max_updates_per_hour to u
< error
| 1064 42000 …line 1 column 40 near "to u" 
> exec ^
: string
| set default role max_user_connections to u
< error
| 1064 42000 …line 1 column 40 near "to u" 
> exec ^
: string
| set default role maxvalue to u
< error
| 1064 42000 …line 1 column 25 near "maxvalue to u" 
> exec ^
: string
| set default role medium to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role mediumblob to u
< error
| 1064 42000 …line 1 column 27 near "mediumblob to u" 
> exec ^
: string
| set default role mediumint to u
< error
| 1064 42000 …line 1 column 26 near "mediumint to u" 
> exec ^
: string
| set default role mediumtext to u
< error
| 1064 42000 …line 1 column 27 near "mediumtext to u" 
> exec ^
: string
| set default role member to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role memory to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role merge to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role message_text to u
< ok 0 0
> exec ^
: string
| set default role microsecond to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role mid to u
< ok 0 0
> exec ^
: string
| set default role middleint to u
< error
| 1064 42000 …line 1 column 26 near "middleint to u" 
> exec ^
: string
| set default role migrate to u
< ok 0 0
> exec ^
: string
| set default role min to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role min_rows to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role minute to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role minute_microsecond to u
< error
| 1064 42000 …line 1 column 35 near "minute_microsecond to u" 
> exec ^
: string
| set default role minute_second to u
< error
| 1064 42000 …line 1 column 30 near "minute_second to u" 
> exec ^
: string
| set default role minvalue to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role mod to u
< error
| 1064 42000 …line 1 column 20 near "mod to u" 
> exec ^
: string
| set default role mode to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role modifies to u
< ok 0 0
> exec ^
: string
| set default role modify to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role month to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role multilinestring to u
< ok 0 0
> exec ^
: string
| set default role multipoint to u
< ok 0 0
> exec ^
: string
| set default role multipolygon to u
< ok 0 0
> exec ^
: string
| set default role mutex to u
< ok 0 0
> exec ^
: string
| set default role mysql_errno to u
< ok 0 0
> exec ^
: string
| set default role name to u
< ok 0 0
> exec ^
: string
| set default role names to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role national to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role natural to u
< error
| 1064 42000 …line 1 column 24 near "natural to u" 
> exec ^
: string
| set default role nchar to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role ndb to u
< ok 0 0
> exec ^
: string
| set default role ndbcluster to u
< ok 0 0
> exec ^
: string
| set default role nested to u
< ok 0 0
> exec ^
: string
| set default role network_namespace to u
< ok 0 0
> exec ^
: string
| set default role never to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role new to u
< ok 0 0
> exec ^
: string
| set default role next to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role next_row_id to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role no to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| set default role no_wait to u
< ok 0 0
> exec ^
: string
| set default role no_write_to_binlog to u
< error
| 1064 42000 …line 1 column 35 near "no_write_to_binlog to u" 
> exec ^
: string
| set default role nocache to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role nocycle to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role nodegroup to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role nomaxvalue to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role nominvalue to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role nonclustered to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role none to u
< ok 0 0
> exec ^
: string
| set default role not to u
< error
| 1064 42000 …line 1 column 20 near "not to u" 
> exec ^
: string
| set default role now to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role nowait to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role nth_value to u
< error
| 1064 42000 …line 1 column 26 near "nth_value to u" 
> exec ^
: string
| set default role ntile to u
< error
| 1064 42000 …line 1 column 22 near "ntile to u" 
> exec ^
: string
| set default role null to u
< error
| 1064 42000 …line 1 column 21 near "null to u" 
> exec ^
: string
| set default role nullif to u
< ok 0 0
> exec ^
: string
| set default role nulls to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role number to u
< ok 0 0
> exec ^
: string
| set default role numeric to u
< error
| 1064 42000 …line 1 column 24 near "numeric to u" 
> exec ^
: string
| set default role nvarchar to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role of to u
< error
| 1064 42000 …line 1 column 19 near "of to u" 
> exec ^
: string
| set default role off to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role offset to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role oj to u
< ok 0 0
> exec ^
: string
| set default role old to u
< ok 0 0
> exec ^
: string
| set default role oltp_read_only to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role oltp_read_write to u
< error
| 1064 42000 …line 1 column 35 near "to u" 
> exec ^
: string
| set default role oltp_write_only to u
< error
| 1064 42000 …line 1 column 35 near "to u" 
> exec ^
: string
| set default role on to u
< error
| 1064 42000 …line 1 column 19 near "on to u" 
> exec ^
: string
| set default role one to u
< ok 0 0
> exec ^
: string
| set default role only to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role open to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role optimize to u
< error
| 1064 42000 …line 1 column 25 near "optimize to u" 
> exec ^
: string
| set default role optimizer_costs to u
< ok 0 0
> exec ^
: string
| set default role option to u
< error
| 1064 42000 …line 1 column 23 near "option to u" 
> exec ^
: string
| set default role optional to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role optionally to u
< error
| 1064 42000 …line 1 column 27 near "optionally to u" 
> exec ^
: string
| set default role options to u
< ok 0 0
> exec ^
: string
| set default role or to u
< error
| 1064 42000 …line 1 column 19 near "or to u" 
> exec ^
: string
| set default role order to u
< error
| 1064 42000 …line 1 column 22 near "order to u" 
> exec ^
: string
| set default role ordinality to u
< ok 0 0
> exec ^
: string
| set default role organization to u
< ok 0 0
> exec ^
: string
| set default role others to u
< ok 0 0
> exec ^
: string
| set default role out to u
< error
| 1064 42000 …line 1 column 20 near "out to u" 
> exec ^
: string
| set default role outer to u
< error
| 1064 42000 …line 1 column 22 near "outer to u" 
> exec ^
: string
| set default role outfile to u
< error
| 1064 42000 …line 1 column 24 near "outfile to u" 
> exec ^
: string
| set default role over to u
< error
| 1064 42000 …line 1 column 21 near "over to u" 
> exec ^
: string
| set default role overlay to u
< ok 0 0
> exec ^
: string
| set default role owner to u
< ok 0 0
> exec ^
: string
| set default role pack_keys to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role page to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role parallel to u
< ok 0 0
> exec ^
: string
| set default role parameters to u
< ok 0 0
> exec ^
: string
| set default role parse_tree to u
< ok 0 0
> exec ^
: string
| set default role parser to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role partial to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role partition to u
< error
| 1064 42000 …line 1 column 26 near "partition to u" 
> exec ^
: string
| set default role partitioning to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role partitions to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role password to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role password_lock_time to u
< error
| 1064 42000 …line 1 column 38 near "to u" 
> exec ^
: string
| set default role path to u
< ok 0 0
> exec ^
: string
| set default role percent_rank to u
< error
| 1064 42000 …line 1 column 29 near "percent_rank to u" 
> exec ^
: string
| set default role persist to u
< ok 0 0
> exec ^
: string
| set default role persist_only to u
< ok 0 0
> exec ^
: string
| set default role phase to u
< ok 0 0
> exec ^
: string
| set default role placement to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role placing to u
< ok 0 0
> exec ^
: string
| set default role plan to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role plugin to u
< ok 0 0
> exec ^
: string
| set default role plugin_dir to u
< ok 0 0
> exec ^
: string
| set default role plugins to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role point to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role policy to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role polygon to u
< ok 0 0
> exec ^
: string
| set default role port to u
< ok 0 0
> exec ^
: string
| set default role position to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role pre_split_regions to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role precedes to u
< ok 0 0
> exec ^
: string
| set default role preceding to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role precision to u
< error
| 1064 42000 …line 1 column 26 near "precision to u" 
> exec ^
: string
| set default role prepare to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role preserve to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role prev to u
< ok 0 0
> exec ^
: string
| set default role primary to u
< error
| 1064 42000 …line 1 column 24 near "primary to u" 
> exec ^
: string
| set default role primary_region to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role priority to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role privilege_checks_user to u
< ok 0 0
> exec ^
: string
| set default role privileges to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role procedure to u
< error
| 1064 42000 …line 1 column 26 near "procedure to u" 
> exec ^
: string
| set default role process to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role processed_keys to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role processlist to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role profile to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role profiles to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role proxy to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role purge to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role qualify to u
< ok 0 0
> exec ^
: string
| set default role quarter to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role query to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role query_limit to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role quick to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role random to u
< ok 0 0
> exec ^
: string
| set default role range to u
< error
| 1064 42000 …line 1 column 22 near "range to u" 
> exec ^
: string
| set default role rank to u
< error
| 1064 42000 …line 1 column 21 near "rank to u" 
> exec ^
: string
| set default role read to u
< error
| 1064 42000 …line 1 column 21 near "read to u" 
> exec ^
: string
| set default role read_only to u
< ok 0 0
> exec ^
: string
| set default role read_write to u
< ok 0 0
> exec ^
: string
| set default role reads to u
< ok 0 0
> exec ^
: string
| set default role real to u
< error
| 1064 42000 …line 1 column 21 near "real to u" 
> exec ^
: string
| set default role rebuild to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role recover to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role recursive to u
< error
| 1064 42000 …line 1 column 26 near "recursive to u" 
> exec ^
: string
| set default role redo_buffer_size to u
< ok 0 0
> exec ^
: string
| set default role redundant to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role reference to u
< ok 0 0
> exec ^
: string
| set default role references to u
< error
| 1064 42000 …line 1 column 27 near "references to u" 
> exec ^
: string
| set default role regexp to u
< error
| 1064 42000 …line 1 column 23 near "regexp to u" 
> exec ^
: string
| set default role region to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role regions to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role registration to u
< ok 0 0
> exec ^
: string
| set default role relational to u
< ok 0 0
> exec ^
: string
| set default role relay to u
< ok 0 0
> exec ^
: string
| set default role relay_log_file to u
< ok 0 0
> exec ^
: string
| set default role relay_log_pos to u
< ok 0 0
> exec ^
: string
| set default role relay_thread to u
< ok 0 0
> exec ^
: string
| set default role relaylog to u
< ok 0 0
> exec ^
: string
| set default role release to u
< error
| 1064 42000 …line 1 column 24 near "release to u" 
> exec ^
: string
| set default role reload to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role remove to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role rename to u
< error
| 1064 42000 …line 1 column 23 near "rename to u" 
> exec ^
: string
| set default role reorganize to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role repair to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role repeat to u
< error
| 1064 42000 …line 1 column 23 near "repeat to u" 
> exec ^
: string
| set default role repeatable to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role replace to u
< error
| 1064 42000 …line 1 column 24 near "replace to u" 
> exec ^
: string
| set default role replica to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role replicas to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role replicate_do_db to u
< ok 0 0
> exec ^
: string
| set default role replicate_do_table to u
< ok 0 0
> exec ^
: string
| set default role replicate_ignore_db to u
< ok 0 0
> exec ^
: string
| set default role replicate_ignore_table to u
< ok 0 0
> exec ^
: string
| set default role replicate_rewrite_db to u
< ok 0 0
> exec ^
: string
| set default role replicate_wild_do_table to u
< ok 0 0
> exec ^
: string
| set default role replicate_wild_ignore_table to u
< ok 0 0
> exec ^
: string
| set default role replication to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role require to u
< error
| 1064 42000 …line 1 column 24 near "require to u" 
> exec ^
: string
| set default role require_row_format to u
< ok 0 0
> exec ^
: string
| set default role require_table_primary_key_check to u
< ok 0 0
> exec ^
: string
| set default role reset to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role resignal to u
< ok 0 0
> exec ^
: string
| set default role resource to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role respect to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role restart to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role restore to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role restores to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role restrict to u
< error
| 1064 42000 …line 1 column 25 near "restrict to u" 
> exec ^
: string
| set default role resume to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role retain to u
< ok 0 0
> exec ^
: string
| set default role return to u
< ok 0 0
> exec ^
: string
| set default role returned_sqlstate to u
< ok 0 0
> exec ^
: string
| set default role returning to u
< ok 0 0
> exec ^
: string
| set default role returns to u
< ok 0 0
> exec ^
: string
| set default role reuse to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role reverse to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role revoke to u
< error
| 1064 42000 …line 1 column 23 near "revoke to u" 
> exec ^
: string
| set default role right to u
< error
| 1064 42000 …line 1 column 22 near "right to u" 
> exec ^
: string
| set default role rlike to u
< error
| 1064 42000 …line 1 column 22 near "rlike to u" 
> exec ^
: string
| set default role role to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role rollback to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role rollup to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role rotate to u
< ok 0 0
> exec ^
: string
| set default role routine to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role row to u
< error
| 1064 42000 …line 1 column 20 near "row to u" 
> exec ^
: string
| set default role row_count to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role row_format to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role row_number to u
< error
| 1064 42000 …line 1 column 27 near "row_number to u" 
> exec ^
: string
| set default role rows to u
< error
| 1064 42000 …line 1 column 21 near "rows to u" 
> exec ^
: string
| set default role rtree to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role ru to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| set default role ru_per_sec to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role rule to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role run to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role s3 to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| set default role savepoint to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role schedule to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role schema to u
< error
| 1064 42000 …line 1 column 23 near "schema to u" 
> exec ^
: string
| set default role schema_name to u
< ok 0 0
> exec ^
: string
| set default role schemas to u
< error
| 1064 42000 …line 1 column 24 near "schemas to u" 
> exec ^
: string
| set default role second to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role second_microsecond to u
< error
| 1064 42000 …line 1 column 35 near "second_microsecond to u" 
> exec ^
: string
| set default role secondary to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role secondary_engine to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role secondary_engine_attribute to u
< ok 0 0
> exec ^
: string
| set default role secondary_load to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role secondary_unload to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role security to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role select to u
< error
| 1064 42000 …line 1 column 23 near "select to u" 
> exec ^
: string
| set default role sensitive to u
< ok 0 0
> exec ^
: string
| set default role separator to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role sequence to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role serial to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role serializable to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role server to u
< ok 0 0
> exec ^
: string
| set default role session to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role session_states to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role session_user to u
< ok 0 0
> exec ^
: string
| set default role set to u
< error
| 1064 42000 …line 1 column 20 near "set to u" 
> exec ^
: string
| set default role sets to u
< ok 0 0
> exec ^
: string
| set default role shard_row_id_bits to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role share to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role show to u
< error
| 1064 42000 …line 1 column 21 near "show to u" 
> exec ^
: string
| set default role shutdown to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role signal to u
< ok 0 0
> exec ^
: string
| set default role signed to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role similar to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role simple to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role skip to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role slave to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role slow to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role smallint to u
< error
| 1064 42000 …line 1 column 25 near "smallint to u" 
> exec ^
: string
| set default role snapshot to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role socket to u
< ok 0 0
> exec ^
: string
| set default role some to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role soname to u
< ok 0 0
> exec ^
: string
| set default role sounds to u
< ok 0 0
> exec ^
: string
| set default role source to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role source_auto_position to u
< ok 0 0
> exec ^
: string
| set default role source_bind to u
< ok 0 0
> exec ^
: string
| set default role source_compression_algorithms to u
< ok 0 0
> exec ^
: string
| set default role source_connect_retry to u
< ok 0 0
> exec ^
: string
| set default role source_connection_auto_failover to u
< ok 0 0
> exec ^
: string
| set default role source_delay to u
< ok 0 0
> exec ^
: string
| set default role source_heartbeat_period to u
< ok 0 0
> exec ^
: string
| set default role source_host to u
< ok 0 0
> exec ^
: string
| set default role source_log_file to u
< ok 0 0
> exec ^
: string
| set default role source_log_pos to u
< ok 0 0
> exec ^
: string
| set default role source_password to u
< ok 0 0
> exec ^
: string
| set default role source_port to u
< ok 0 0
> exec ^
: string
| set default role source_public_key_path to u
< ok 0 0
> exec ^
: string
| set default role source_retry_count to u
< ok 0 0
> exec ^
: string
| set default role source_ssl to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_ca to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_capath to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_cert to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_cipher to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_crl to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_crlpath to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_key to u
< ok 0 0
> exec ^
: string
| set default role source_ssl_verify_server_cert to u
< ok 0 0
> exec ^
: string
| set default role source_tls_ciphersuites to u
< ok 0 0
> exec ^
: string
| set default role source_tls_version to u
< ok 0 0
> exec ^
: string
| set default role source_user to u
< ok 0 0
> exec ^
: string
| set default role source_zstd_compression_level to u
< ok 0 0
> exec ^
: string
| set default role spatial to u
< error
| 1064 42000 …line 1 column 24 near "spatial to u" 
> exec ^
: string
| set default role specific to u
< ok 0 0
> exec ^
: string
| set default role split to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role sql to u
< error
| 1064 42000 …line 1 column 20 near "sql to u" 
> exec ^
: string
| set default role sql_after_gtids to u
< ok 0 0
> exec ^
: string
| set default role sql_after_mts_gaps to u
< ok 0 0
> exec ^
: string
| set default role sql_before_gtids to u
< ok 0 0
> exec ^
: string
| set default role sql_big_result to u
< error
| 1064 42000 …line 1 column 31 near "sql_big_result to u" 
> exec ^
: string
| set default role sql_buffer_result to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role sql_calc_found_rows to u
< error
| 1064 42000 …line 1 column 36 near "sql_calc_found_rows to u" 
> exec ^
: string
| set default role sql_no_cache to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role sql_small_result to u
< error
| 1064 42000 …line 1 column 33 near "sql_small_result to u" 
> exec ^
: string
| set default role sql_thread to u
< ok 0 0
> exec ^
: string
| set default role sql_tsi_day to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role sql_tsi_hour to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role sql_tsi_minute to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role sql_tsi_month to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role sql_tsi_quarter to u
< error
| 1064 42000 …line 1 column 35 near "to u" 
> exec ^
: string
| set default role sql_tsi_second to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role sql_tsi_week to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role sql_tsi_year to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role sqlexception to u
< error
| 1064 42000 …line 1 column 29 near "sqlexception to u" 
> exec ^
: string
| set default role sqlstate to u
< error
| 1064 42000 …line 1 column 25 near "sqlstate to u" 
> exec ^
: string
| set default role sqlwarning to u
< error
| 1064 42000 …line 1 column 27 near "sqlwarning to u" 
> exec ^
: string
| set default role srid to u
< ok 0 0
> exec ^
: string
| set default role ssl to u
< error
| 1064 42000 …line 1 column 20 near "ssl to u" 
> exec ^
: string
| set default role st_collect to u
< ok 0 0
> exec ^
: string
| set default role stacked to u
< ok 0 0
> exec ^
: string
| set default role start to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role start_time to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role starting to u
< error
| 1064 42000 …line 1 column 25 near "starting to u" 
> exec ^
: string
| set default role starts to u
< ok 0 0
> exec ^
: string
| set default role stats_auto_recalc to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role stats_buckets to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role stats_col_choice to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role stats_col_list to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role stats_extended to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role stats_healthy to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role stats_histograms to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role stats_locked to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role stats_meta to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role stats_persistent to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role stats_sample_pages to u
< error
| 1064 42000 …line 1 column 38 near "to u" 
> exec ^
: string
| set default role stats_sample_rate to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role stats_topn to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role status to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role std to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role stddev to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role stddev_pop to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role stddev_samp to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role stop to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role storage to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role stored to u
< error
| 1064 42000 …line 1 column 23 near "stored to u" 
> exec ^
: string
| set default role straight_join to u
< error
| 1064 42000 …line 1 column 30 near "straight_join to u" 
> exec ^
: string
| set default role stream to u
< ok 0 0
> exec ^
: string
| set default role strict_load to u
< ok 0 0
> exec ^
: string
| set default role string to u
< ok 0 0
> exec ^
: string
| set default role subclass_origin to u
< ok 0 0
> exec ^
: string
| set default role subdate to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role subject to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role subpartition to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role subpartitions to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role substr to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role substring to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role sum to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role super to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role survival_preferences to u
< error
| 1064 42000 …line 1 column 40 near "to u" 
> exec ^
: string
| set default role suspend to u
< ok 0 0
> exec ^
: string
| set default role swaps to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role switch_group to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role switches to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role sysdate to u
< ok 0 0
> exec ^
: string
| set default role system to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role system_user to u
< ok 0 0
> exec ^
: string
| set default role table to u
< error
| 1064 42000 …line 1 column 22 near "table to u" 
> exec ^
: string
| set default role table_checksum to u
< error
| 1064 42000 …line 1 column 34 near "to u" 
> exec ^
: string
| set default role table_name to u
< ok 0 0
> exec ^
: string
| set default role tables to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role tablesample to u
< error
| 1064 42000 …line 1 column 28 near "tablesample to u" 
> exec ^
: string
| set default role tablespace to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role task_types to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role temporary to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role temptable to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role terminated to u
< error
| 1064 42000 …line 1 column 27 near "terminated to u" 
> exec ^
: string
| set default role text to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role than to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role then to u
< error
| 1064 42000 …line 1 column 21 near "then to u" 
> exec ^
: string
| set default role thread_priority to u
< ok 0 0
> exec ^
: string
| set default role ties to u
< ok 0 0
> exec ^
: string
| set default role tiflash to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role time to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role timeout to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role timestamp to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role timestampadd to u
< error
| 1064 42000 …line 1 column 32 near "to u" 
> exec ^
: string
| set default role timestampdiff to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role tinyblob to u
< error
| 1064 42000 …line 1 column 25 near "tinyblob to u" 
> exec ^
: string
| set default role tinyint to u
< error
| 1064 42000 …line 1 column 24 near "tinyint to u" 
> exec ^
: string
| set default role tinytext to u
< error
| 1064 42000 …line 1 column 25 near "tinytext to u" 
> exec ^
: string
| set default role tls to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role to to u
< error
| 1064 42000 …line 1 column 19 near "to to u" 
> exec ^
: string
| set default role tpcc to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role tpch_10 to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role trailing to u
< error
| 1064 42000 …line 1 column 25 near "trailing to u" 
> exec ^
: string
| set default role transaction to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role trigger to u
< error
| 1064 42000 …line 1 column 24 near "trigger to u" 
> exec ^
: string
| set default role triggers to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role trim to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role true to u
< error
| 1064 42000 …line 1 column 21 near "true to u" 
> exec ^
: string
| set default role truncate to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role ttl to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| set default role ttl_enable to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role ttl_job_interval to u
< error
| 1064 42000 …line 1 column 36 near "to u" 
> exec ^
: string
| set default role type to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role types to u
< ok 0 0
> exec ^
: string
| set default role unbounded to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role uncommitted to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| set default role undefined to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role undo to u
< ok 0 0
> exec ^
: string
| set default role undo_buffer_size to u
< ok 0 0
> exec ^
: string
| set default role undofile to u
< ok 0 0
> exec ^
: string
| set default role unicode to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role uninstall to u
< ok 0 0
> exec ^
: string
| set default role union to u
< error
| 1064 42000 …line 1 column 22 near "union to u" 
> exec ^
: string
| set default role unique to u
< error
| 1064 42000 …line 1 column 23 near "unique to u" 
> exec ^
: string
| set default role unknown to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role unlimited to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role unlock to u
< error
| 1064 42000 …line 1 column 23 near "unlock to u" 
> exec ^
: string
| set default role unregister to u
< ok 0 0
> exec ^
: string
| set default role unsigned to u
< error
| 1064 42000 …line 1 column 25 near "unsigned to u" 
> exec ^
: string
| set default role until to u
< error
| 1064 42000 …line 1 column 22 near "until to u" 
> exec ^
: string
| set default role update to u
< error
| 1064 42000 …line 1 column 23 near "update to u" 
> exec ^
: string
| set default role upgrade to u
< ok 0 0
> exec ^
: string
| set default role uri to u
< ok 0 0
> exec ^
: string
| set default role url to u
< ok 0 0
> exec ^
: string
| set default role usage to u
< error
| 1064 42000 …line 1 column 22 near "usage to u" 
> exec ^
: string
| set default role use to u
< error
| 1064 42000 …line 1 column 20 near "use to u" 
> exec ^
: string
| set default role use_frm to u
< ok 0 0
> exec ^
: string
| set default role user to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role user_resources to u
< ok 0 0
> exec ^
: string
| set default role using to u
< error
| 1064 42000 …line 1 column 22 near "using to u" 
> exec ^
: string
| set default role utc_date to u
< error
| 1064 42000 …line 1 column 25 near "utc_date to u" 
> exec ^
: string
| set default role utc_time to u
< error
| 1064 42000 …line 1 column 25 near "utc_time to u" 
> exec ^
: string
| set default role utc_timestamp to u
< error
| 1064 42000 …line 1 column 30 near "utc_timestamp to u" 
> exec ^
: string
| set default role utilization_limit to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role validate to u
< ok 0 0
> exec ^
: string
| set default role validation to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| set default role value to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role values to u
< error
| 1064 42000 …line 1 column 23 near "values to u" 
> exec ^
: string
| set default role var_pop to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role var_samp to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role varbinary to u
< error
| 1064 42000 …line 1 column 26 near "varbinary to u" 
> exec ^
: string
| set default role varchar to u
< error
| 1064 42000 …line 1 column 24 near "varchar to u" 
> exec ^
: string
| set default role varcharacter to u
< error
| 1064 42000 …line 1 column 29 near "varcharacter to u" 
> exec ^
: string
| set default role variables to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| set default role variance to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role varying to u
< error
| 1064 42000 …line 1 column 24 near "varying to u" 
> exec ^
: string
| set default role vcpu to u
< ok 0 0
> exec ^
: string
| set default role vector to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role verify_key_constraints to u
< ok 0 0
> exec ^
: string
| set default role view to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role virtual to u
< error
| 1064 42000 …line 1 column 24 near "virtual to u" 
> exec ^
: string
| set default role visible to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role voter_constraints to u
< error
| 1064 42000 …line 1 column 37 near "to u" 
> exec ^
: string
| set default role voters to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| set default role wait to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role warnings to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role watch to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| set default role week to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role weight_string to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| set default role when to u
< error
| 1064 42000 …line 1 column 21 near "when to u" 
> exec ^
: string
| set default role where to u
< error
| 1064 42000 …line 1 column 22 near "where to u" 
> exec ^
: string
| set default role while to u
< error
| 1064 42000 …line 1 column 22 near "while to u" 
> exec ^
: string
| set default role window to u
< error
| 1064 42000 …line 1 column 23 near "window to u" 
> exec ^
: string
| set default role with to u
< error
| 1064 42000 …line 1 column 21 near "with to u" 
> exec ^
: string
| set default role without to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| set default role work to u
< ok 0 0
> exec ^
: string
| set default role workload to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| set default role wrapper to u
< ok 0 0
> exec ^
: string
| set default role write to u
< error
| 1064 42000 …line 1 column 22 near "write to u" 
> exec ^
: string
| set default role x509 to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role xa to u
< ok 0 0
> exec ^
: string
| set default role xid to u
< ok 0 0
> exec ^
: string
| set default role xml to u
< ok 0 0
> exec ^
: string
| set default role xor to u
< error
| 1064 42000 …line 1 column 20 near "xor to u" 
> exec ^
: string
| set default role year to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| set default role year_month to u
< error
| 1064 42000 …line 1 column 27 near "year_month to u" 
> exec ^
: string
| set default role zerofill to u
< error
| 1064 42000 …line 1 column 25 near "zerofill to u" 
> exec ^
: string
| set default role zone to u
< ok 0 0
> exec ^
: string
| grant absent to u
< ok 0 0
> exec ^
: string
| grant accessible to u
< ok 0 0
> exec ^
: string
| grant account to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant action to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant active to u
< ok 0 0
> exec ^
: string
| grant add to u
< error
| 1064 42000 …line 1 column 9 near "add to u" 
> exec ^
: string
| grant adddate to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant admin to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant after to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant against to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant aggregate to u
< ok 0 0
> exec ^
: string
| grant algorithm to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant all to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant allow_missing_files to u
< ok 0 0
> exec ^
: string
| grant alter to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant always to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant analyze to u
< error
| 1064 42000 …line 1 column 13 near "analyze to u" 
> exec ^
: string
| grant and to u
< error
| 1064 42000 …line 1 column 9 near "and to u" 
> exec ^
: string
| grant any to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant array to u
< error
| 1064 42000 …line 1 column 11 near "array to u" 
> exec ^
: string
| grant as to u
< error
| 1064 42000 …line 1 column 8 near "as to u" 
> exec ^
: string
| grant asc to u
< error
| 1064 42000 …line 1 column 9 near "asc to u" 
> exec ^
: string
| grant ascii to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant asensitive to u
< ok 0 0
> exec ^
: string
| grant assign_gtids_to_anonymous_transactions to u
< ok 0 0
> exec ^
: string
| grant at to u
< ok 0 0
> exec ^
: string
| grant attribute to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant authentication to u
< ok 0 0
> exec ^
: string
| grant auto to u
< ok 0 0
> exec ^
: string
| grant auto_id_cache to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant auto_increment to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant auto_random to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant auto_random_base to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant auto_refresh to u
< ok 0 0
> exec ^
: string
| grant auto_refresh_source to u
< ok 0 0
> exec ^
: string
| grant autoextend_size to u
< ok 0 0
> exec ^
: string
| grant avg to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant avg_row_length to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant background to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant backup to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant backups to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant batch to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant before to u
< ok 0 0
> exec ^
: string
| grant begin to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant bernoulli to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant between to u
< error
| 1064 42000 …line 1 column 13 near "between to u" 
> exec ^
: string
| grant bigint to u
< error
| 1064 42000 …line 1 column 12 near "bigint to u" 
> exec ^
: string
| grant binary to u
< error
| 1064 42000 …line 1 column 12 near "binary to u" 
> exec ^
: string
| grant binding to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant binding_cache to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant bindings to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant binlog to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant bit to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant bit_and to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant bit_or to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant bit_xor to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant blob to u
< error
| 1064 42000 …line 1 column 10 near "blob to u" 
> exec ^
: string
| grant block to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant bool to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant boolean to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant both to u
< error
| 1064 42000 …line 1 column 10 near "both to u" 
> exec ^
: string
| grant btree to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant buckets to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant builtins to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant bulk to u
< ok 0 0
> exec ^
: string
| grant burstable to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant by to u
< error
| 1064 42000 …line 1 column 8 near "by to u" 
> exec ^
: string
| grant byte to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant cache to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant calibrate to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant call to u
< error
| 1064 42000 …line 1 column 10 near "call to u" 
> exec ^
: string
| grant cascade to u
< error
| 1064 42000 …line 1 column 13 near "cascade to u" 
> exec ^
: string
| grant cascaded to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant case to u
< error
| 1064 42000 …line 1 column 10 near "case to u" 
> exec ^
: string
| grant cast to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant catalog_name to u
< ok 0 0
> exec ^
: string
| grant chain to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant challenge_response to u
< ok 0 0
> exec ^
: string
| grant change to u
< error
| 1064 42000 …line 1 column 12 near "change to u" 
> exec ^
: string
| grant changed to u
< ok 0 0
> exec ^
: string
| grant channel to u
< ok 0 0
> exec ^
: string
| grant char to u
< error
| 1064 42000 …line 1 column 10 near "char to u" 
> exec ^
: string
| grant character to u
< error
| 1064 42000 …line 1 column 15 near "character to u" 
> exec ^
: string
| grant charset to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant check to u
< error
| 1064 42000 …line 1 column 11 near "check to u" 
> exec ^
: string
| grant checksum to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant cipher to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant class_origin to u
< ok 0 0
> exec ^
: string
| grant client to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant clone to u
< ok 0 0
> exec ^
: string
| grant close to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant clustered to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant coalesce to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant code to u
< ok 0 0
> exec ^
: string
| grant collate to u
< error
| 1064 42000 …line 1 column 13 near "collate to u" 
> exec ^
: string
| grant collation to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant column to u
< error
| 1064 42000 …line 1 column 12 near "column to u" 
> exec ^
: string
| grant column_format to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant column_name to u
< ok 0 0
> exec ^
: string
| grant column_stats_usage to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| grant columns to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant comment to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant commit to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant committed to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant compact to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant completion to u
< ok 0 0
> exec ^
: string
| grant component to u
< ok 0 0
> exec ^
: string
| grant compressed to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant compression to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant concat to u
< ok 0 0
> exec ^
: string
| grant concurrent to u
< ok 0 0
> exec ^
: string
| grant condition to u
< ok 0 0
> exec ^
: string
| grant config to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant connection to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant consistent to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant constraint to u
< error
| 1064 42000 …line 1 column 16 near "constraint to u" 
> exec ^
: string
| grant constraint_catalog to u
< ok 0 0
> exec ^
: string
| grant constraint_name to u
< ok 0 0
> exec ^
: string
| grant constraint_schema to u
< ok 0 0
> exec ^
: string
| grant constraints to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant contains to u
< ok 0 0
> exec ^
: string
| grant context to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant continue to u
< error
| 1064 42000 …line 1 column 14 near "continue to u" 
> exec ^
: string
| grant convert to u
< error
| 1064 42000 …line 1 column 13 near "convert to u" 
> exec ^
: string
| grant cooldown to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant count to u
< ok 0 0
> exec ^
: string
| grant cpu to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant create to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant cross to u
< error
| 1064 42000 …line 1 column 11 near "cross to u" 
> exec ^
: string
| grant cube to u
< ok 0 0
> exec ^
: string
| grant cume_dist to u
< error
| 1064 42000 …line 1 column 15 near "cume_dist to u" 
> exec ^
: string
| grant curdate to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant current to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant current_date to u
< error
| 1064 42000 …line 1 column 18 near "current_date to u" 
> exec ^
: string
| grant current_time to u
< error
| 1064 42000 …line 1 column 18 near "current_time to u" 
> exec ^
: string
| grant current_timestamp to u
< error
| 1064 42000 …line 1 column 23 near "current_timestamp to u" 
> exec ^
: string
| grant current_user to u
< error
| 1064 42000 …line 1 column 18 near "current_user to u" 
> exec ^
: string
| grant cursor to u
< error
| 1064 42000 …line 1 column 12 near "cursor to u" 
> exec ^
: string
| grant cursor_name to u
< ok 0 0
> exec ^
: string
| grant curtime to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant cycle to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant data to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant database to u
< error
| 1064 42000 …line 1 column 14 near "database to u" 
> exec ^
: string
| grant databases to u
< error
| 1064 42000 …line 1 column 15 near "databases to u" 
> exec ^
: string
| grant datafile to u
< ok 0 0
> exec ^
: string
| grant date to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant date_add to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant date_sub to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant datetime to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant day to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant day_hour to u
< error
| 1064 42000 …line 1 column 14 near "day_hour to u" 
> exec ^
: string
| grant day_microsecond to u
< error
| 1064 42000 …line 1 column 21 near "day_microsecond to u" 
> exec ^
: string
| grant day_minute to u
< error
| 1064 42000 …line 1 column 16 near "day_minute to u" 
> exec ^
: string
| grant day_second to u
< error
| 1064 42000 …line 1 column 16 near "day_second to u" 
> exec ^
: string
| grant deallocate to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant dec to u
< error
| 1064 42000 …line 1 column 9 near "dec to u" 
> exec ^
: string
| grant decimal to u
< error
| 1064 42000 …line 1 column 13 near "decimal to u" 
> exec ^
: string
| grant declare to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant default to u
< error
| 1064 42000 …line 1 column 13 near "default to u" 
> exec ^
: string
| grant default_auth to u
< ok 0 0
> exec ^
: string
| grant defined to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant definer to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant definition to u
< ok 0 0
> exec ^
: string
| grant delay_key_write to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| grant delayed to u
< error
| 1064 42000 …line 1 column 13 near "delayed to u" 
> exec ^
: string
| grant delete to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant dense_rank to u
< error
| 1064 42000 …line 1 column 16 near "dense_rank to u" 
> exec ^
: string
| grant desc to u
< error
| 1064 42000 …line 1 column 10 near "desc to u" 
> exec ^
: string
| grant describe to u
< error
| 1064 42000 …line 1 column 14 near "describe to u" 
> exec ^
: string
| grant description to u
< ok 0 0
> exec ^
: string
| grant deterministic to u
< ok 0 0
> exec ^
: string
| grant diagnostics to u
< ok 0 0
> exec ^
: string
| grant digest to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant directory to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant disable to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant disabled to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant discard to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant disk to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant distinct to u
< error
| 1064 42000 …line 1 column 14 near "distinct to u" 
> exec ^
: string
| grant distinctrow to u
< error
| 1064 42000 …line 1 column 17 near "distinctrow to u" 
> exec ^
: string
| grant distribute to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant distribution to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant distributions to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant div to u
< error
| 1064 42000 …line 1 column 9 near "div to u" 
> exec ^
: string
| grant do to u
< error
| 1064 42000 …line 1 column 11 near "to u" 
> exec ^
: string
| grant double to u
< error
| 1064 42000 …line 1 column 12 near "double to u" 
> exec ^
: string
| grant drop to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant dry to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant dryrun to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant dual to u
< error
| 1064 42000 …line 1 column 10 near "dual to u" 
> exec ^
: string
| grant duality to u
< ok 0 0
> exec ^
: string
| grant dumpfile to u
< ok 0 0
> exec ^
: string
| grant duplicate to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant duration to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant dynamic to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant each to u
< ok 0 0
> exec ^
: string
| grant else to u
< error
| 1064 42000 …line 1 column 10 near "else to u" 
> exec ^
: string
| grant elseif to u
< error
| 1064 42000 …line 1 column 12 near "elseif to u" 
> exec ^
: string
| grant empty to u
< ok 0 0
> exec ^
: string
| grant enable to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant enabled to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant enclosed to u
< error
| 1064 42000 …line 1 column 14 near "enclosed to u" 
> exec ^
: string
| grant encryption to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant end to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant end_time to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant ends to u
< ok 0 0
> exec ^
: string
| grant enforced to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant engine to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant engine_attribute to u
< ok 0 0
> exec ^
: string
| grant engines to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant enum to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant error to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant errors to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant escape to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant escaped to u
< error
| 1064 42000 …line 1 column 13 near "escaped to u" 
> exec ^
: string
| grant event to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant events to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant every to u
< ok 0 0
> exec ^
: string
| grant exact to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant except to u
< error
| 1064 42000 …line 1 column 12 near "except to u" 
> exec ^
: string
| grant exchange to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant exclude to u
< ok 0 0
> exec ^
: string
| grant exec_elapsed to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant execute to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant exists to u
< error
| 1064 42000 …line 1 column 12 near "exists to u" 
> exec ^
: string
| grant exit to u
< error
| 1064 42000 …line 1 column 10 near "exit to u" 
> exec ^
: string
| grant expansion to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant expire to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant explain to u
< error
| 1064 42000 …line 1 column 13 near "explain to u" 
> exec ^
: string
| grant export to u
< ok 0 0
> exec ^
: string
| grant extended to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant extent_size to u
< ok 0 0
> exec ^
: string
| grant external to u
< ok 0 0
> exec ^
: string
| grant external_format to u
< ok 0 0
> exec ^
: string
| grant extract to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant factor to u
< ok 0 0
> exec ^
: string
| grant failed_login_attempts to u
< error
| 1064 42000 …line 1 column 30 near "to u" 
> exec ^
: string
| grant false to u
< error
| 1064 42000 …line 1 column 11 near "false to u" 
> exec ^
: string
| grant fast to u
< ok 0 0
> exec ^
: string
| grant faults to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant fetch to u
< error
| 1064 42000 …line 1 column 11 near "fetch to u" 
> exec ^
: string
| grant fields to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant file to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant file_block_size to u
< ok 0 0
> exec ^
: string
| grant file_format to u
< ok 0 0
> exec ^
: string
| grant file_name to u
< ok 0 0
> exec ^
: string
| grant file_pattern to u
< ok 0 0
> exec ^
: string
| grant file_prefix to u
< ok 0 0
> exec ^
: string
| grant files to u
< ok 0 0
> exec ^
: string
| grant filter to u
< ok 0 0
> exec ^
: string
| grant finish to u
< ok 0 0
> exec ^
: string
| grant first to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant first_value to u
< error
| 1064 42000 …line 1 column 17 near "first_value to u" 
> exec ^
: string
| grant fixed to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant flashback to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant float to u
< error
| 1064 42000 …line 1 column 11 near "float to u" 
> exec ^
: string
| grant float4 to u
< error
| 1064 42000 …line 1 column 12 near "float4 to u" 
> exec ^
: string
| grant float8 to u
< error
| 1064 42000 …line 1 column 12 near "float8 to u" 
> exec ^
: string
| grant flush to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant follower_constraints to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| grant followers to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant following to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant follows to u
< ok 0 0
> exec ^
: string
| grant for to u
< error
| 1064 42000 …line 1 column 9 near "for to u" 
> exec ^
: string
| grant force to u
< error
| 1064 42000 …line 1 column 11 near "force to u" 
> exec ^
: string
| grant foreign to u
< error
| 1064 42000 …line 1 column 13 near "foreign to u" 
> exec ^
: string
| grant format to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant found to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant from to u
< error
| 1064 42000 …line 1 column 10 near "from to u" 
> exec ^
: string
| grant full to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant fulltext to u
< error
| 1064 42000 …line 1 column 14 near "fulltext to u" 
> exec ^
: string
| grant function to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant general to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant generate to u
< ok 0 0
> exec ^
: string
| grant generated to u
< error
| 1064 42000 …line 1 column 15 near "generated to u" 
> exec ^
: string
| grant geomcollection to u
< ok 0 0
> exec ^
: string
| grant geometry to u
< ok 0 0
> exec ^
: string
| grant geometrycollection to u
< ok 0 0
> exec ^
: string
| grant get to u
< ok 0 0
> exec ^
: string
| grant get_format to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant get_source_public_key to u
< ok 0 0
> exec ^
: string
| grant global to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant grant to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant grants to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant greatest to u
< ok 0 0
> exec ^
: string
| grant group to u
< error
| 1064 42000 …line 1 column 11 near "group to u" 
> exec ^
: string
| grant group_concat to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant group_replication to u
< ok 0 0
> exec ^
: string
| grant grouping to u
< ok 0 0
> exec ^
: string
| grant groups to u
< error
| 1064 42000 …line 1 column 12 near "groups to u" 
> exec ^
: string
| grant gtid_only to u
< ok 0 0
> exec ^
: string
| grant gtids to u
< ok 0 0
> exec ^
: string
| grant guided to u
< ok 0 0
> exec ^
: string
| grant handler to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant hash to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant having to u
< error
| 1064 42000 …line 1 column 12 near "having to u" 
> exec ^
: string
| grant header to u
< ok 0 0
> exec ^
: string
| grant help to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant high to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant high_priority to u
< error
| 1064 42000 …line 1 column 19 near "high_priority to u" 
> exec ^
: string
| grant histogram to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant histograms_in_flight to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| grant history to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant host to u
< ok 0 0
> exec ^
: string
| grant hosts to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant hour to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant hour_microsecond to u
< error
| 1064 42000 …line 1 column 22 near "hour_microsecond to u" 
> exec ^
: string
| grant hour_minute to u
< error
| 1064 42000 …line 1 column 17 near "hour_minute to u" 
> exec ^
: string
| grant hour_second to u
< error
| 1064 42000 …line 1 column 17 near "hour_second to u" 
> exec ^
: string
| grant identified to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant if to u
< error
| 1064 42000 …line 1 column 8 near "if to u" 
> exec ^
: string
| grant ignore to u
< error
| 1064 42000 …line 1 column 12 near "ignore to u" 
> exec ^
: string
| grant ignore_server_ids to u
< ok 0 0
> exec ^
: string
| grant ilike to u
< error
| 1064 42000 …line 1 column 11 near "ilike to u" 
> exec ^
: string
| grant import to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant in to u
< error
| 1064 42000 …line 1 column 8 near "in to u" 
> exec ^
: string
| grant inactive to u
< ok 0 0
> exec ^
: string
| grant increment to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant index to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant indexes to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant infile to u
< error
| 1064 42000 …line 1 column 12 near "infile to u" 
> exec ^
: string
| grant initial to u
< ok 0 0
> exec ^
: string
| grant initial_size to u
< ok 0 0
> exec ^
: string
| grant initiate to u
< ok 0 0
> exec ^
: string
| grant inner to u
< error
| 1064 42000 …line 1 column 11 near "inner to u" 
> exec ^
: string
| grant inout to u
< error
| 1064 42000 …line 1 column 11 near "inout to u" 
> exec ^
: string
| grant insensitive to u
< ok 0 0
> exec ^
: string
| grant insert to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant insert_method to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant install to u
< ok 0 0
> exec ^
: string
| grant instance to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant int to u
< error
| 1064 42000 …line 1 column 9 near "int to u" 
> exec ^
: string
| grant int1 to u
< error
| 1064 42000 …line 1 column 10 near "int1 to u" 
> exec ^
: string
| grant int2 to u
< error
| 1064 42000 …line 1 column 10 near "int2 to u" 
> exec ^
: string
| grant int3 to u
< error
| 1064 42000 …line 1 column 10 near "int3 to u" 
> exec ^
: string
| grant int4 to u
< error
| 1064 42000 …line 1 column 10 near "int4 to u" 
> exec ^
: string
| grant int8 to u
< error
| 1064 42000 …line 1 column 10 near "int8 to u" 
> exec ^
: string
| grant integer to u
< error
| 1064 42000 …line 1 column 13 near "integer to u" 
> exec ^
: string
| grant intersect to u
< error
| 1064 42000 …line 1 column 15 near "intersect to u" 
> exec ^
: string
| grant interval to u
< error
| 1064 42000 …line 1 column 14 near "interval to u" 
> exec ^
: string
| grant into to u
< error
| 1064 42000 …line 1 column 10 near "into to u" 
> exec ^
: string
| grant invisible to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant invoker to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant io to u
< error
| 1064 42000 …line 1 column 11 near "to u" 
> exec ^
: string
| grant io_after_gtids to u
< ok 0 0
> exec ^
: string
| grant io_before_gtids to u
< ok 0 0
> exec ^
: string
| grant io_thread to u
< ok 0 0
> exec ^
: string
| grant ipc to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant is to u
< error
| 1064 42000 …line 1 column 8 near "is to u" 
> exec ^
: string
| grant isolation to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant issuer to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant iterate to u
< error
| 1064 42000 …line 1 column 13 near "iterate to u" 
> exec ^
: string
| grant job to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant jobs to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant join to u
< error
| 1064 42000 …line 1 column 10 near "join to u" 
> exec ^
: string
| grant json to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant json_arrayagg to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant json_duality_object to u
< ok 0 0
> exec ^
: string
| grant json_objectagg to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant json_table to u
< ok 0 0
> exec ^
: string
| grant json_value to u
< ok 0 0
> exec ^
: string
| grant key to u
< error
| 1064 42000 …line 1 column 9 near "key to u" 
> exec ^
: string
| grant key_block_size to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant keyring to u
< ok 0 0
> exec ^
: string
| grant keys to u
< error
| 1064 42000 …line 1 column 10 near "keys to u" 
> exec ^
: string
| grant kill to u
< error
| 1064 42000 …line 1 column 10 near "kill to u" 
> exec ^
: string
| grant labels to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant lag to u
< error
| 1064 42000 …line 1 column 9 near "lag to u" 
> exec ^
: string
| grant language to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant last to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant last_value to u
< error
| 1064 42000 …line 1 column 16 near "last_value to u" 
> exec ^
: string
| grant lateral to u
< ok 0 0
> exec ^
: string
| grant lead to u
< error
| 1064 42000 …line 1 column 10 near "lead to u" 
> exec ^
: string
| grant leader_constraints to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| grant leading to u
< error
| 1064 42000 …line 1 column 13 near "leading to u" 
> exec ^
: string
| grant learner_constraints to u
< error
| 1064 42000 …line 1 column 28 near "to u" 
> exec ^
: string
| grant learners to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant least to u
< ok 0 0
> exec ^
: string
| grant leave to u
< error
| 1064 42000 …line 1 column 11 near "leave to u" 
> exec ^
: string
| grant leaves to u
< ok 0 0
> exec ^
: string
| grant left to u
< error
| 1064 42000 …line 1 column 10 near "left to u" 
> exec ^
: string
| grant less to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant level to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant library to u
< ok 0 0
> exec ^
: string
| grant like to u
< error
| 1064 42000 …line 1 column 10 near "like to u" 
> exec ^
: string
| grant limit to u
< error
| 1064 42000 …line 1 column 11 near "limit to u" 
> exec ^
: string
| grant linear to u
< error
| 1064 42000 …line 1 column 12 near "linear to u" 
> exec ^
: string
| grant lines to u
< error
| 1064 42000 …line 1 column 11 near "lines to u" 
> exec ^
: string
| grant linestring to u
< ok 0 0
> exec ^
: string
| grant list to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant load to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant local to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant localtime to u
< error
| 1064 42000 …line 1 column 15 near "localtime to u" 
> exec ^
: string
| grant localtimestamp to u
< error
| 1064 42000 …line 1 column 20 near "localtimestamp to u" 
> exec ^
: string
| grant location to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant lock to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant locked to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant locks to u
< ok 0 0
> exec ^
: string
| grant log to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant logfile to u
< ok 0 0
> exec ^
: string
| grant logs to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant long to u
< error
| 1064 42000 …line 1 column 10 near "long to u" 
> exec ^
: string
| grant longblob to u
< error
| 1064 42000 …line 1 column 14 near "longblob to u" 
> exec ^
: string
| grant longtext to u
< error
| 1064 42000 …line 1 column 14 near "longtext to u" 
> exec ^
: string
| grant loop to u
< ok 0 0
> exec ^
: string
| grant low to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant low_priority to u
< error
| 1064 42000 …line 1 column 18 near "low_priority to u" 
> exec ^
: string
| grant manual to u
< ok 0 0
> exec ^
: string
| grant master to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant master_log_file to u
< ok 0 0
> exec ^
: string
| grant match to u
< error
| 1064 42000 …line 1 column 11 near "match to u" 
> exec ^
: string
| grant materialized to u
< ok 0 0
> exec ^
: string
| grant max to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant max_connections_per_hour to u
< error
| 1064 42000 …line 1 column 33 near "to u" 
> exec ^
: string
| grant max_queries_per_hour to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| grant max_rows to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant max_size to u
< ok 0 0
> exec ^
: string
| grant max_updates_per_hour to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| grant max_user_connections to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| grant maxvalue to u
< error
| 1064 42000 …line 1 column 14 near "maxvalue to u" 
> exec ^
: string
| grant medium to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant mediumblob to u
< error
| 1064 42000 …line 1 column 16 near "mediumblob to u" 
> exec ^
: string
| grant mediumint to u
< error
| 1064 42000 …line 1 column 15 near "mediumint to u" 
> exec ^
: string
| grant mediumtext to u
< error
| 1064 42000 …line 1 column 16 near "mediumtext to u" 
> exec ^
: string
| grant member to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant memory to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant merge to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant message_text to u
< ok 0 0
> exec ^
: string
| grant microsecond to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant mid to u
< ok 0 0
> exec ^
: string
| grant middleint to u
< error
| 1064 42000 …line 1 column 15 near "middleint to u" 
> exec ^
: string
| grant migrate to u
< ok 0 0
> exec ^
: string
| grant min to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant min_rows to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant minute to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant minute_microsecond to u
< error
| 1064 42000 …line 1 column 24 near "minute_microsecond to u" 
> exec ^
: string
| grant minute_second to u
< error
| 1064 42000 …line 1 column 19 near "minute_second to u" 
> exec ^
: string
| grant minvalue to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant mod to u
< error
| 1064 42000 …line 1 column 9 near "mod to u" 
> exec ^
: string
| grant mode to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant modifies to u
< ok 0 0
> exec ^
: string
| grant modify to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant month to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant multilinestring to u
< ok 0 0
> exec ^
: string
| grant multipoint to u
< ok 0 0
> exec ^
: string
| grant multipolygon to u
< ok 0 0
> exec ^
: string
| grant mutex to u
< ok 0 0
> exec ^
: string
| grant mysql_errno to u
< ok 0 0
> exec ^
: string
| grant name to u
< ok 0 0
> exec ^
: string
| grant names to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant national to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant natural to u
< error
| 1064 42000 …line 1 column 13 near "natural to u" 
> exec ^
: string
| grant nchar to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant ndb to u
< ok 0 0
> exec ^
: string
| grant ndbcluster to u
< ok 0 0
> exec ^
: string
| grant nested to u
< ok 0 0
> exec ^
: string
| grant network_namespace to u
< ok 0 0
> exec ^
: string
| grant never to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant new to u
< ok 0 0
> exec ^
: string
| grant next to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant next_row_id to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant no to u
< error
| 1064 42000 …line 1 column 11 near "to u" 
> exec ^
: string
| grant no_wait to u
< ok 0 0
> exec ^
: string
| grant no_write_to_binlog to u
< error
| 1064 42000 …line 1 column 24 near "no_write_to_binlog to u" 
> exec ^
: string
| grant nocache to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant nocycle to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant nodegroup to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant nomaxvalue to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant nominvalue to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant nonclustered to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant none to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant not to u
< error
| 1064 42000 …line 1 column 9 near "not to u" 
> exec ^
: string
| grant now to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant nowait to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant nth_value to u
< error
| 1064 42000 …line 1 column 15 near "nth_value to u" 
> exec ^
: string
| grant ntile to u
< error
| 1064 42000 …line 1 column 11 near "ntile to u" 
> exec ^
: string
| grant null to u
< error
| 1064 42000 …line 1 column 10 near "null to u" 
> exec ^
: string
| grant nullif to u
< ok 0 0
> exec ^
: string
| grant nulls to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant number to u
< ok 0 0
> exec ^
: string
| grant numeric to u
< error
| 1064 42000 …line 1 column 13 near "numeric to u" 
> exec ^
: string
| grant nvarchar to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant of to u
< error
| 1064 42000 …line 1 column 8 near "of to u" 
> exec ^
: string
| grant off to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant offset to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant oj to u
< ok 0 0
> exec ^
: string
| grant old to u
< ok 0 0
> exec ^
: string
| grant oltp_read_only to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant oltp_read_write to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| grant oltp_write_only to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| grant on to u
< error
| 1064 42000 …line 1 column 8 near "on to u" 
> exec ^
: string
| grant one to u
< ok 0 0
> exec ^
: string
| grant only to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant open to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant optimize to u
< error
| 1064 42000 …line 1 column 14 near "optimize to u" 
> exec ^
: string
| grant optimizer_costs to u
< ok 0 0
> exec ^
: string
| grant option to u
< error
| 1064 42000 …line 1 column 12 near "option to u" 
> exec ^
: string
| grant optional to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant optionally to u
< error
| 1064 42000 …line 1 column 16 near "optionally to u" 
> exec ^
: string
| grant options to u
< ok 0 0
> exec ^
: string
| grant or to u
< error
| 1064 42000 …line 1 column 8 near "or to u" 
> exec ^
: string
| grant order to u
< error
| 1064 42000 …line 1 column 11 near "order to u" 
> exec ^
: string
| grant ordinality to u
< ok 0 0
> exec ^
: string
| grant organization to u
< ok 0 0
> exec ^
: string
| grant others to u
< ok 0 0
> exec ^
: string
| grant out to u
< error
| 1064 42000 …line 1 column 9 near "out to u" 
> exec ^
: string
| grant outer to u
< error
| 1064 42000 …line 1 column 11 near "outer to u" 
> exec ^
: string
| grant outfile to u
< error
| 1064 42000 …line 1 column 13 near "outfile to u" 
> exec ^
: string
| grant over to u
< error
| 1064 42000 …line 1 column 10 near "over to u" 
> exec ^
: string
| grant overlay to u
< ok 0 0
> exec ^
: string
| grant owner to u
< ok 0 0
> exec ^
: string
| grant pack_keys to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant page to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant parallel to u
< ok 0 0
> exec ^
: string
| grant parameters to u
< ok 0 0
> exec ^
: string
| grant parse_tree to u
< ok 0 0
> exec ^
: string
| grant parser to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant partial to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant partition to u
< error
| 1064 42000 …line 1 column 15 near "partition to u" 
> exec ^
: string
| grant partitioning to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant partitions to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant password to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant password_lock_time to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| grant path to u
< ok 0 0
> exec ^
: string
| grant percent_rank to u
< error
| 1064 42000 …line 1 column 18 near "percent_rank to u" 
> exec ^
: string
| grant persist to u
< ok 0 0
> exec ^
: string
| grant persist_only to u
< ok 0 0
> exec ^
: string
| grant phase to u
< ok 0 0
> exec ^
: string
| grant placement to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant placing to u
< ok 0 0
> exec ^
: string
| grant plan to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant plugin to u
< ok 0 0
> exec ^
: string
| grant plugin_dir to u
< ok 0 0
> exec ^
: string
| grant plugins to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant point to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant policy to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant polygon to u
< ok 0 0
> exec ^
: string
| grant port to u
< ok 0 0
> exec ^
: string
| grant position to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant pre_split_regions to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant precedes to u
< ok 0 0
> exec ^
: string
| grant preceding to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant precision to u
< error
| 1064 42000 …line 1 column 15 near "precision to u" 
> exec ^
: string
| grant prepare to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant preserve to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant prev to u
< ok 0 0
> exec ^
: string
| grant primary to u
< error
| 1064 42000 …line 1 column 13 near "primary to u" 
> exec ^
: string
| grant primary_region to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant priority to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant privilege_checks_user to u
< ok 0 0
> exec ^
: string
| grant privileges to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant procedure to u
< error
| 1064 42000 …line 1 column 15 near "procedure to u" 
> exec ^
: string
| grant process to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant processed_keys to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant processlist to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant profile to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant profiles to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant proxy to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant purge to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant qualify to u
< ok 0 0
> exec ^
: string
| grant quarter to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant query to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant query_limit to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant quick to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant random to u
< ok 0 0
> exec ^
: string
| grant range to u
< error
| 1064 42000 …line 1 column 11 near "range to u" 
> exec ^
: string
| grant rank to u
< error
| 1064 42000 …line 1 column 10 near "rank to u" 
> exec ^
: string
| grant read to u
< error
| 1064 42000 …line 1 column 10 near "read to u" 
> exec ^
: string
| grant read_only to u
< ok 0 0
> exec ^
: string
| grant read_write to u
< ok 0 0
> exec ^
: string
| grant reads to u
< ok 0 0
> exec ^
: string
| grant real to u
< error
| 1064 42000 …line 1 column 10 near "real to u" 
> exec ^
: string
| grant rebuild to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant recover to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant recursive to u
< error
| 1064 42000 …line 1 column 15 near "recursive to u" 
> exec ^
: string
| grant redo_buffer_size to u
< ok 0 0
> exec ^
: string
| grant redundant to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant reference to u
< ok 0 0
> exec ^
: string
| grant references to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant regexp to u
< error
| 1064 42000 …line 1 column 12 near "regexp to u" 
> exec ^
: string
| grant region to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant regions to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant registration to u
< ok 0 0
> exec ^
: string
| grant relational to u
< ok 0 0
> exec ^
: string
| grant relay to u
< ok 0 0
> exec ^
: string
| grant relay_log_file to u
< ok 0 0
> exec ^
: string
| grant relay_log_pos to u
< ok 0 0
> exec ^
: string
| grant relay_thread to u
< ok 0 0
> exec ^
: string
| grant relaylog to u
< ok 0 0
> exec ^
: string
| grant release to u
< error
| 1064 42000 …line 1 column 13 near "release to u" 
> exec ^
: string
| grant reload to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant remove to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant rename to u
< error
| 1064 42000 …line 1 column 12 near "rename to u" 
> exec ^
: string
| grant reorganize to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant repair to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant repeat to u
< error
| 1064 42000 …line 1 column 12 near "repeat to u" 
> exec ^
: string
| grant repeatable to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant replace to u
< error
| 1064 42000 …line 1 column 13 near "replace to u" 
> exec ^
: string
| grant replica to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant replicas to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant replicate_do_db to u
< ok 0 0
> exec ^
: string
| grant replicate_do_table to u
< ok 0 0
> exec ^
: string
| grant replicate_ignore_db to u
< ok 0 0
> exec ^
: string
| grant replicate_ignore_table to u
< ok 0 0
> exec ^
: string
| grant replicate_rewrite_db to u
< ok 0 0
> exec ^
: string
| grant replicate_wild_do_table to u
< ok 0 0
> exec ^
: string
| grant replicate_wild_ignore_table to u
< ok 0 0
> exec ^
: string
| grant replication to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant require to u
< error
| 1064 42000 …line 1 column 13 near "require to u" 
> exec ^
: string
| grant require_row_format to u
< ok 0 0
> exec ^
: string
| grant require_table_primary_key_check to u
< ok 0 0
> exec ^
: string
| grant reset to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant resignal to u
< ok 0 0
> exec ^
: string
| grant resource to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant respect to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant restart to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant restore to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant restores to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant restrict to u
< error
| 1064 42000 …line 1 column 14 near "restrict to u" 
> exec ^
: string
| grant resume to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant retain to u
< ok 0 0
> exec ^
: string
| grant return to u
< ok 0 0
> exec ^
: string
| grant returned_sqlstate to u
< ok 0 0
> exec ^
: string
| grant returning to u
< ok 0 0
> exec ^
: string
| grant returns to u
< ok 0 0
> exec ^
: string
| grant reuse to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant reverse to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant revoke to u
< error
| 1064 42000 …line 1 column 12 near "revoke to u" 
> exec ^
: string
| grant right to u
< error
| 1064 42000 …line 1 column 11 near "right to u" 
> exec ^
: string
| grant rlike to u
< error
| 1064 42000 …line 1 column 11 near "rlike to u" 
> exec ^
: string
| grant role to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant rollback to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant rollup to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant rotate to u
< ok 0 0
> exec ^
: string
| grant routine to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant row to u
< error
| 1064 42000 …line 1 column 9 near "row to u" 
> exec ^
: string
| grant row_count to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant row_format to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant row_number to u
< error
| 1064 42000 …line 1 column 16 near "row_number to u" 
> exec ^
: string
| grant rows to u
< error
| 1064 42000 …line 1 column 10 near "rows to u" 
> exec ^
: string
| grant rtree to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant ru to u
< error
| 1064 42000 …line 1 column 11 near "to u" 
> exec ^
: string
| grant ru_per_sec to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant rule to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant run to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant s3 to u
< error
| 1064 42000 …line 1 column 11 near "to u" 
> exec ^
: string
| grant savepoint to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant schedule to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant schema to u
< error
| 1064 42000 …line 1 column 12 near "schema to u" 
> exec ^
: string
| grant schema_name to u
< ok 0 0
> exec ^
: string
| grant schemas to u
< error
| 1064 42000 …line 1 column 13 near "schemas to u" 
> exec ^
: string
| grant second to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant second_microsecond to u
< error
| 1064 42000 …line 1 column 24 near "second_microsecond to u" 
> exec ^
: string
| grant secondary to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant secondary_engine to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant secondary_engine_attribute to u
< ok 0 0
> exec ^
: string
| grant secondary_load to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant secondary_unload to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant security to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant select to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant sensitive to u
< ok 0 0
> exec ^
: string
| grant separator to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant sequence to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant serial to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant serializable to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant server to u
< ok 0 0
> exec ^
: string
| grant session to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant session_states to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant session_user to u
< ok 0 0
> exec ^
: string
| grant set to u
< error
| 1064 42000 …line 1 column 9 near "set to u" 
> exec ^
: string
| grant sets to u
< ok 0 0
> exec ^
: string
| grant shard_row_id_bits to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant share to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant show to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant shutdown to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant signal to u
< ok 0 0
> exec ^
: string
| grant signed to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant similar to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant simple to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant skip to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant slave to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant slow to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant smallint to u
< error
| 1064 42000 …line 1 column 14 near "smallint to u" 
> exec ^
: string
| grant snapshot to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant socket to u
< ok 0 0
> exec ^
: string
| grant some to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant soname to u
< ok 0 0
> exec ^
: string
| grant sounds to u
< ok 0 0
> exec ^
: string
| grant source to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant source_auto_position to u
< ok 0 0
> exec ^
: string
| grant source_bind to u
< ok 0 0
> exec ^
: string
| grant source_compression_algorithms to u
< ok 0 0
> exec ^
: string
| grant source_connect_retry to u
< ok 0 0
> exec ^
: string
| grant source_connection_auto_failover to u
< ok 0 0
> exec ^
: string
| grant source_delay to u
< ok 0 0
> exec ^
: string
| grant source_heartbeat_period to u
< ok 0 0
> exec ^
: string
| grant source_host to u
< ok 0 0
> exec ^
: string
| grant source_log_file to u
< ok 0 0
> exec ^
: string
| grant source_log_pos to u
< ok 0 0
> exec ^
: string
| grant source_password to u
< ok 0 0
> exec ^
: string
| grant source_port to u
< ok 0 0
> exec ^
: string
| grant source_public_key_path to u
< ok 0 0
> exec ^
: string
| grant source_retry_count to u
< ok 0 0
> exec ^
: string
| grant source_ssl to u
< ok 0 0
> exec ^
: string
| grant source_ssl_ca to u
< ok 0 0
> exec ^
: string
| grant source_ssl_capath to u
< ok 0 0
> exec ^
: string
| grant source_ssl_cert to u
< ok 0 0
> exec ^
: string
| grant source_ssl_cipher to u
< ok 0 0
> exec ^
: string
| grant source_ssl_crl to u
< ok 0 0
> exec ^
: string
| grant source_ssl_crlpath to u
< ok 0 0
> exec ^
: string
| grant source_ssl_key to u
< ok 0 0
> exec ^
: string
| grant source_ssl_verify_server_cert to u
< ok 0 0
> exec ^
: string
| grant source_tls_ciphersuites to u
< ok 0 0
> exec ^
: string
| grant source_tls_version to u
< ok 0 0
> exec ^
: string
| grant source_user to u
< ok 0 0
> exec ^
: string
| grant source_zstd_compression_level to u
< ok 0 0
> exec ^
: string
| grant spatial to u
< error
| 1064 42000 …line 1 column 13 near "spatial to u" 
> exec ^
: string
| grant specific to u
< ok 0 0
> exec ^
: string
| grant split to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant sql to u
< error
| 1064 42000 …line 1 column 9 near "sql to u" 
> exec ^
: string
| grant sql_after_gtids to u
< ok 0 0
> exec ^
: string
| grant sql_after_mts_gaps to u
< ok 0 0
> exec ^
: string
| grant sql_before_gtids to u
< ok 0 0
> exec ^
: string
| grant sql_big_result to u
< error
| 1064 42000 …line 1 column 20 near "sql_big_result to u" 
> exec ^
: string
| grant sql_buffer_result to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant sql_calc_found_rows to u
< error
| 1064 42000 …line 1 column 25 near "sql_calc_found_rows to u" 
> exec ^
: string
| grant sql_no_cache to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant sql_small_result to u
< error
| 1064 42000 …line 1 column 22 near "sql_small_result to u" 
> exec ^
: string
| grant sql_thread to u
< ok 0 0
> exec ^
: string
| grant sql_tsi_day to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant sql_tsi_hour to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant sql_tsi_minute to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant sql_tsi_month to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant sql_tsi_quarter to u
< error
| 1064 42000 …line 1 column 24 near "to u" 
> exec ^
: string
| grant sql_tsi_second to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant sql_tsi_week to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant sql_tsi_year to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant sqlexception to u
< error
| 1064 42000 …line 1 column 18 near "sqlexception to u" 
> exec ^
: string
| grant sqlstate to u
< error
| 1064 42000 …line 1 column 14 near "sqlstate to u" 
> exec ^
: string
| grant sqlwarning to u
< error
| 1064 42000 …line 1 column 16 near "sqlwarning to u" 
> exec ^
: string
| grant srid to u
< ok 0 0
> exec ^
: string
| grant ssl to u
< error
| 1064 42000 …line 1 column 9 near "ssl to u" 
> exec ^
: string
| grant st_collect to u
< ok 0 0
> exec ^
: string
| grant stacked to u
< ok 0 0
> exec ^
: string
| grant start to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant start_time to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant starting to u
< error
| 1064 42000 …line 1 column 14 near "starting to u" 
> exec ^
: string
| grant starts to u
< ok 0 0
> exec ^
: string
| grant stats_auto_recalc to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant stats_buckets to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant stats_col_choice to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant stats_col_list to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant stats_extended to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant stats_healthy to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant stats_histograms to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant stats_locked to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant stats_meta to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant stats_persistent to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant stats_sample_pages to u
< error
| 1064 42000 …line 1 column 27 near "to u" 
> exec ^
: string
| grant stats_sample_rate to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant stats_topn to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant status to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant std to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant stddev to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant stddev_pop to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant stddev_samp to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant stop to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant storage to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant stored to u
< error
| 1064 42000 …line 1 column 12 near "stored to u" 
> exec ^
: string
| grant straight_join to u
< error
| 1064 42000 …line 1 column 19 near "straight_join to u" 
> exec ^
: string
| grant stream to u
< ok 0 0
> exec ^
: string
| grant strict_load to u
< ok 0 0
> exec ^
: string
| grant string to u
< ok 0 0
> exec ^
: string
| grant subclass_origin to u
< ok 0 0
> exec ^
: string
| grant subdate to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant subject to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant subpartition to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant subpartitions to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant substr to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant substring to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant sum to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant super to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant survival_preferences to u
< error
| 1064 42000 …line 1 column 29 near "to u" 
> exec ^
: string
| grant suspend to u
< ok 0 0
> exec ^
: string
| grant swaps to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant switch_group to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant switches to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant sysdate to u
< ok 0 0
> exec ^
: string
| grant system to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant system_user to u
< ok 0 0
> exec ^
: string
| grant table to u
< error
| 1064 42000 …line 1 column 11 near "table to u" 
> exec ^
: string
| grant table_checksum to u
< error
| 1064 42000 …line 1 column 23 near "to u" 
> exec ^
: string
| grant table_name to u
< ok 0 0
> exec ^
: string
| grant tables to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant tablesample to u
< error
| 1064 42000 …line 1 column 17 near "tablesample to u" 
> exec ^
: string
| grant tablespace to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant task_types to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant temporary to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant temptable to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant terminated to u
< error
| 1064 42000 …line 1 column 16 near "terminated to u" 
> exec ^
: string
| grant text to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant than to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant then to u
< error
| 1064 42000 …line 1 column 10 near "then to u" 
> exec ^
: string
| grant thread_priority to u
< ok 0 0
> exec ^
: string
| grant ties to u
< ok 0 0
> exec ^
: string
| grant tiflash to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant time to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant timeout to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant timestamp to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant timestampadd to u
< error
| 1064 42000 …line 1 column 21 near "to u" 
> exec ^
: string
| grant timestampdiff to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant tinyblob to u
< error
| 1064 42000 …line 1 column 14 near "tinyblob to u" 
> exec ^
: string
| grant tinyint to u
< error
| 1064 42000 …line 1 column 13 near "tinyint to u" 
> exec ^
: string
| grant tinytext to u
< error
| 1064 42000 …line 1 column 14 near "tinytext to u" 
> exec ^
: string
| grant tls to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant to to u
< error
| 1064 42000 …line 1 column 8 near "to to u" 
> exec ^
: string
| grant tpcc to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant tpch_10 to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant trailing to u
< error
| 1064 42000 …line 1 column 14 near "trailing to u" 
> exec ^
: string
| grant transaction to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant trigger to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant triggers to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant trim to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant true to u
< error
| 1064 42000 …line 1 column 10 near "true to u" 
> exec ^
: string
| grant truncate to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant ttl to u
< error
| 1064 42000 …line 1 column 12 near "to u" 
> exec ^
: string
| grant ttl_enable to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant ttl_job_interval to u
< error
| 1064 42000 …line 1 column 25 near "to u" 
> exec ^
: string
| grant type to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant types to u
< ok 0 0
> exec ^
: string
| grant unbounded to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant uncommitted to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant undefined to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant undo to u
< ok 0 0
> exec ^
: string
| grant undo_buffer_size to u
< ok 0 0
> exec ^
: string
| grant undofile to u
< ok 0 0
> exec ^
: string
| grant unicode to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant uninstall to u
< ok 0 0
> exec ^
: string
| grant union to u
< error
| 1064 42000 …line 1 column 11 near "union to u" 
> exec ^
: string
| grant unique to u
< error
| 1064 42000 …line 1 column 12 near "unique to u" 
> exec ^
: string
| grant unknown to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant unlimited to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant unlock to u
< error
| 1064 42000 …line 1 column 12 near "unlock to u" 
> exec ^
: string
| grant unregister to u
< ok 0 0
> exec ^
: string
| grant unsigned to u
< error
| 1064 42000 …line 1 column 14 near "unsigned to u" 
> exec ^
: string
| grant until to u
< error
| 1064 42000 …line 1 column 11 near "until to u" 
> exec ^
: string
| grant update to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant upgrade to u
< ok 0 0
> exec ^
: string
| grant uri to u
< ok 0 0
> exec ^
: string
| grant url to u
< ok 0 0
> exec ^
: string
| grant usage to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant use to u
< error
| 1064 42000 …line 1 column 9 near "use to u" 
> exec ^
: string
| grant use_frm to u
< ok 0 0
> exec ^
: string
| grant user to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant user_resources to u
< ok 0 0
> exec ^
: string
| grant using to u
< error
| 1064 42000 …line 1 column 11 near "using to u" 
> exec ^
: string
| grant utc_date to u
< error
| 1064 42000 …line 1 column 14 near "utc_date to u" 
> exec ^
: string
| grant utc_time to u
< error
| 1064 42000 …line 1 column 14 near "utc_time to u" 
> exec ^
: string
| grant utc_timestamp to u
< error
| 1064 42000 …line 1 column 19 near "utc_timestamp to u" 
> exec ^
: string
| grant utilization_limit to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant validate to u
< ok 0 0
> exec ^
: string
| grant validation to u
< error
| 1064 42000 …line 1 column 19 near "to u" 
> exec ^
: string
| grant value to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant values to u
< error
| 1064 42000 …line 1 column 12 near "values to u" 
> exec ^
: string
| grant var_pop to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant var_samp to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant varbinary to u
< error
| 1064 42000 …line 1 column 15 near "varbinary to u" 
> exec ^
: string
| grant varchar to u
< error
| 1064 42000 …line 1 column 13 near "varchar to u" 
> exec ^
: string
| grant varcharacter to u
< error
| 1064 42000 …line 1 column 18 near "varcharacter to u" 
> exec ^
: string
| grant variables to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| grant variance to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant varying to u
< error
| 1064 42000 …line 1 column 13 near "varying to u" 
> exec ^
: string
| grant vcpu to u
< ok 0 0
> exec ^
: string
| grant vector to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant verify_key_constraints to u
< ok 0 0
> exec ^
: string
| grant view to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant virtual to u
< error
| 1064 42000 …line 1 column 13 near "virtual to u" 
> exec ^
: string
| grant visible to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant voter_constraints to u
< error
| 1064 42000 …line 1 column 26 near "to u" 
> exec ^
: string
| grant voters to u
< error
| 1064 42000 …line 1 column 15 near "to u" 
> exec ^
: string
| grant wait to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant warnings to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant watch to u
< error
| 1064 42000 …line 1 column 14 near "to u" 
> exec ^
: string
| grant week to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant weight_string to u
< error
| 1064 42000 …line 1 column 22 near "to u" 
> exec ^
: string
| grant when to u
< error
| 1064 42000 …line 1 column 10 near "when to u" 
> exec ^
: string
| grant where to u
< error
| 1064 42000 …line 1 column 11 near "where to u" 
> exec ^
: string
| grant while to u
< error
| 1064 42000 …line 1 column 11 near "while to u" 
> exec ^
: string
| grant window to u
< error
| 1064 42000 …line 1 column 12 near "window to u" 
> exec ^
: string
| grant with to u
< error
| 1064 42000 …line 1 column 10 near "with to u" 
> exec ^
: string
| grant without to u
< error
| 1064 42000 …line 1 column 16 near "to u" 
> exec ^
: string
| grant work to u
< ok 0 0
> exec ^
: string
| grant workload to u
< error
| 1064 42000 …line 1 column 17 near "to u" 
> exec ^
: string
| grant wrapper to u
< ok 0 0
> exec ^
: string
| grant write to u
< error
| 1064 42000 …line 1 column 11 near "write to u" 
> exec ^
: string
| grant x509 to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant xa to u
< ok 0 0
> exec ^
: string
| grant xid to u
< ok 0 0
> exec ^
: string
| grant xml to u
< ok 0 0
> exec ^
: string
| grant xor to u
< error
| 1064 42000 …line 1 column 9 near "xor to u" 
> exec ^
: string
| grant year to u
< error
| 1064 42000 …line 1 column 13 near "to u" 
> exec ^
: string
| grant year_month to u
< error
| 1064 42000 …line 1 column 16 near "year_month to u" 
> exec ^
: string
| grant zerofill to u
< error
| 1064 42000 …line 1 column 14 near "zerofill to u" 
> exec ^
: string
| grant zone to u
< ok 0 0
> exec ^
: string
| revoke absent from u
< ok 0 0
> exec ^
: string
| revoke accessible from u
< ok 0 0
> exec ^
: string
| revoke account from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke action from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke active from u
< ok 0 0
> exec ^
: string
| revoke add from u
< error
| 1064 42000 …line 1 column 10 near "add from u" 
> exec ^
: string
| revoke adddate from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke admin from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke after from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke against from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke aggregate from u
< ok 0 0
> exec ^
: string
| revoke algorithm from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke all from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke allow_missing_files from u
< ok 0 0
> exec ^
: string
| revoke alter from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke always from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke analyze from u
< error
| 1064 42000 …line 1 column 14 near "analyze from u" 
> exec ^
: string
| revoke and from u
< error
| 1064 42000 …line 1 column 10 near "and from u" 
> exec ^
: string
| revoke any from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke array from u
< error
| 1064 42000 …line 1 column 12 near "array from u" 
> exec ^
: string
| revoke as from u
< error
| 1064 42000 …line 1 column 9 near "as from u" 
> exec ^
: string
| revoke asc from u
< error
| 1064 42000 …line 1 column 10 near "asc from u" 
> exec ^
: string
| revoke ascii from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke asensitive from u
< ok 0 0
> exec ^
: string
| revoke assign_gtids_to_anonymous_transactions from u
< ok 0 0
> exec ^
: string
| revoke at from u
< ok 0 0
> exec ^
: string
| revoke attribute from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke authentication from u
< ok 0 0
> exec ^
: string
| revoke auto from u
< ok 0 0
> exec ^
: string
| revoke auto_id_cache from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke auto_increment from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke auto_random from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke auto_random_base from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke auto_refresh from u
< ok 0 0
> exec ^
: string
| revoke auto_refresh_source from u
< ok 0 0
> exec ^
: string
| revoke autoextend_size from u
< ok 0 0
> exec ^
: string
| revoke avg from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke avg_row_length from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke background from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke backup from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke backups from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke batch from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke before from u
< ok 0 0
> exec ^
: string
| revoke begin from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke bernoulli from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke between from u
< error
| 1064 42000 …line 1 column 14 near "between from u" 
> exec ^
: string
| revoke bigint from u
< error
| 1064 42000 …line 1 column 13 near "bigint from u" 
> exec ^
: string
| revoke binary from u
< error
| 1064 42000 …line 1 column 13 near "binary from u" 
> exec ^
: string
| revoke binding from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke binding_cache from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke bindings from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke binlog from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke bit from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke bit_and from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke bit_or from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke bit_xor from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke blob from u
< error
| 1064 42000 …line 1 column 11 near "blob from u" 
> exec ^
: string
| revoke block from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke bool from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke boolean from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke both from u
< error
| 1064 42000 …line 1 column 11 near "both from u" 
> exec ^
: string
| revoke btree from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke buckets from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke builtins from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke bulk from u
< ok 0 0
> exec ^
: string
| revoke burstable from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke by from u
< error
| 1064 42000 …line 1 column 9 near "by from u" 
> exec ^
: string
| revoke byte from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke cache from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke calibrate from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke call from u
< error
| 1064 42000 …line 1 column 11 near "call from u" 
> exec ^
: string
| revoke cascade from u
< error
| 1064 42000 …line 1 column 14 near "cascade from u" 
> exec ^
: string
| revoke cascaded from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke case from u
< error
| 1064 42000 …line 1 column 11 near "case from u" 
> exec ^
: string
| revoke cast from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke catalog_name from u
< ok 0 0
> exec ^
: string
| revoke chain from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke challenge_response from u
< ok 0 0
> exec ^
: string
| revoke change from u
< error
| 1064 42000 …line 1 column 13 near "change from u" 
> exec ^
: string
| revoke changed from u
< ok 0 0
> exec ^
: string
| revoke channel from u
< ok 0 0
> exec ^
: string
| revoke char from u
< error
| 1064 42000 …line 1 column 11 near "char from u" 
> exec ^
: string
| revoke character from u
< error
| 1064 42000 …line 1 column 16 near "character from u" 
> exec ^
: string
| revoke charset from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke check from u
< error
| 1064 42000 …line 1 column 12 near "check from u" 
> exec ^
: string
| revoke checksum from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke cipher from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke class_origin from u
< ok 0 0
> exec ^
: string
| revoke client from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke clone from u
< ok 0 0
> exec ^
: string
| revoke close from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke clustered from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke coalesce from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke code from u
< ok 0 0
> exec ^
: string
| revoke collate from u
< error
| 1064 42000 …line 1 column 14 near "collate from u" 
> exec ^
: string
| revoke collation from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke column from u
< error
| 1064 42000 …line 1 column 13 near "column from u" 
> exec ^
: string
| revoke column_format from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke column_name from u
< ok 0 0
> exec ^
: string
| revoke column_stats_usage from u
< error
| 1064 42000 …line 1 column 30 near "from u" 
> exec ^
: string
| revoke columns from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke comment from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke commit from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke committed from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke compact from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke completion from u
< ok 0 0
> exec ^
: string
| revoke component from u
< ok 0 0
> exec ^
: string
| revoke compressed from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke compression from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke concat from u
< ok 0 0
> exec ^
: string
| revoke concurrent from u
< ok 0 0
> exec ^
: string
| revoke condition from u
< ok 0 0
> exec ^
: string
| revoke config from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke connection from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke consistent from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke constraint from u
< error
| 1064 42000 …line 1 column 17 near "constraint from u" 
> exec ^
: string
| revoke constraint_catalog from u
< ok 0 0
> exec ^
: string
| revoke constraint_name from u
< ok 0 0
> exec ^
: string
| revoke constraint_schema from u
< ok 0 0
> exec ^
: string
| revoke constraints from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke contains from u
< ok 0 0
> exec ^
: string
| revoke context from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke continue from u
< error
| 1064 42000 …line 1 column 15 near "continue from u" 
> exec ^
: string
| revoke convert from u
< error
| 1064 42000 …line 1 column 14 near "convert from u" 
> exec ^
: string
| revoke cooldown from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke count from u
< ok 0 0
> exec ^
: string
| revoke cpu from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke create from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke cross from u
< error
| 1064 42000 …line 1 column 12 near "cross from u" 
> exec ^
: string
| revoke cube from u
< ok 0 0
> exec ^
: string
| revoke cume_dist from u
< error
| 1064 42000 …line 1 column 16 near "cume_dist from u" 
> exec ^
: string
| revoke curdate from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke current from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke current_date from u
< error
| 1064 42000 …line 1 column 19 near "current_date from u" 
> exec ^
: string
| revoke current_time from u
< error
| 1064 42000 …line 1 column 19 near "current_time from u" 
> exec ^
: string
| revoke current_timestamp from u
< error
| 1064 42000 …line 1 column 24 near "current_timestamp from u" 
> exec ^
: string
| revoke current_user from u
< error
| 1064 42000 …line 1 column 19 near "current_user from u" 
> exec ^
: string
| revoke cursor from u
< error
| 1064 42000 …line 1 column 13 near "cursor from u" 
> exec ^
: string
| revoke cursor_name from u
< ok 0 0
> exec ^
: string
| revoke curtime from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke cycle from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke data from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke database from u
< error
| 1064 42000 …line 1 column 15 near "database from u" 
> exec ^
: string
| revoke databases from u
< error
| 1064 42000 …line 1 column 16 near "databases from u" 
> exec ^
: string
| revoke datafile from u
< ok 0 0
> exec ^
: string
| revoke date from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke date_add from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke date_sub from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke datetime from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke day from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke day_hour from u
< error
| 1064 42000 …line 1 column 15 near "day_hour from u" 
> exec ^
: string
| revoke day_microsecond from u
< error
| 1064 42000 …line 1 column 22 near "day_microsecond from u" 
> exec ^
: string
| revoke day_minute from u
< error
| 1064 42000 …line 1 column 17 near "day_minute from u" 
> exec ^
: string
| revoke day_second from u
< error
| 1064 42000 …line 1 column 17 near "day_second from u" 
> exec ^
: string
| revoke deallocate from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke dec from u
< error
| 1064 42000 …line 1 column 10 near "dec from u" 
> exec ^
: string
| revoke decimal from u
< error
| 1064 42000 …line 1 column 14 near "decimal from u" 
> exec ^
: string
| revoke declare from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke default from u
< error
| 1064 42000 …line 1 column 14 near "default from u" 
> exec ^
: string
| revoke default_auth from u
< ok 0 0
> exec ^
: string
| revoke defined from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke definer from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke definition from u
< ok 0 0
> exec ^
: string
| revoke delay_key_write from u
< error
| 1064 42000 …line 1 column 27 near "from u" 
> exec ^
: string
| revoke delayed from u
< error
| 1064 42000 …line 1 column 14 near "delayed from u" 
> exec ^
: string
| revoke delete from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke dense_rank from u
< error
| 1064 42000 …line 1 column 17 near "dense_rank from u" 
> exec ^
: string
| revoke desc from u
< error
| 1064 42000 …line 1 column 11 near "desc from u" 
> exec ^
: string
| revoke describe from u
< error
| 1064 42000 …line 1 column 15 near "describe from u" 
> exec ^
: string
| revoke description from u
< ok 0 0
> exec ^
: string
| revoke deterministic from u
< ok 0 0
> exec ^
: string
| revoke diagnostics from u
< ok 0 0
> exec ^
: string
| revoke digest from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke directory from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke disable from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke disabled from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke discard from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke disk from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke distinct from u
< error
| 1064 42000 …line 1 column 15 near "distinct from u" 
> exec ^
: string
| revoke distinctrow from u
< error
| 1064 42000 …line 1 column 18 near "distinctrow from u" 
> exec ^
: string
| revoke distribute from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke distribution from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke distributions from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke div from u
< error
| 1064 42000 …line 1 column 10 near "div from u" 
> exec ^
: string
| revoke do from u
< error
| 1064 42000 …line 1 column 14 near "from u" 
> exec ^
: string
| revoke double from u
< error
| 1064 42000 …line 1 column 13 near "double from u" 
> exec ^
: string
| revoke drop from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke dry from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke dryrun from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke dual from u
< error
| 1064 42000 …line 1 column 11 near "dual from u" 
> exec ^
: string
| revoke duality from u
< ok 0 0
> exec ^
: string
| revoke dumpfile from u
< ok 0 0
> exec ^
: string
| revoke duplicate from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke duration from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke dynamic from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke each from u
< ok 0 0
> exec ^
: string
| revoke else from u
< error
| 1064 42000 …line 1 column 11 near "else from u" 
> exec ^
: string
| revoke elseif from u
< error
| 1064 42000 …line 1 column 13 near "elseif from u" 
> exec ^
: string
| revoke empty from u
< ok 0 0
> exec ^
: string
| revoke enable from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke enabled from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke enclosed from u
< error
| 1064 42000 …line 1 column 15 near "enclosed from u" 
> exec ^
: string
| revoke encryption from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke end from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke end_time from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke ends from u
< ok 0 0
> exec ^
: string
| revoke enforced from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke engine from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke engine_attribute from u
< ok 0 0
> exec ^
: string
| revoke engines from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke enum from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke error from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke errors from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke escape from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke escaped from u
< error
| 1064 42000 …line 1 column 14 near "escaped from u" 
> exec ^
: string
| revoke event from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke events from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke every from u
< ok 0 0
> exec ^
: string
| revoke exact from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke except from u
< error
| 1064 42000 …line 1 column 13 near "except from u" 
> exec ^
: string
| revoke exchange from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke exclude from u
< ok 0 0
> exec ^
: string
| revoke exec_elapsed from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke execute from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke exists from u
< error
| 1064 42000 …line 1 column 13 near "exists from u" 
> exec ^
: string
| revoke exit from u
< error
| 1064 42000 …line 1 column 11 near "exit from u" 
> exec ^
: string
| revoke expansion from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke expire from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke explain from u
< error
| 1064 42000 …line 1 column 14 near "explain from u" 
> exec ^
: string
| revoke export from u
< ok 0 0
> exec ^
: string
| revoke extended from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke extent_size from u
< ok 0 0
> exec ^
: string
| revoke external from u
< ok 0 0
> exec ^
: string
| revoke external_format from u
< ok 0 0
> exec ^
: string
| revoke extract from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke factor from u
< ok 0 0
> exec ^
: string
| revoke failed_login_attempts from u
< error
| 1064 42000 …line 1 column 33 near "from u" 
> exec ^
: string
| revoke false from u
< error
| 1064 42000 …line 1 column 12 near "false from u" 
> exec ^
: string
| revoke fast from u
< ok 0 0
> exec ^
: string
| revoke faults from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke fetch from u
< error
| 1064 42000 …line 1 column 12 near "fetch from u" 
> exec ^
: string
| revoke fields from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke file from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke file_block_size from u
< ok 0 0
> exec ^
: string
| revoke file_format from u
< ok 0 0
> exec ^
: string
| revoke file_name from u
< ok 0 0
> exec ^
: string
| revoke file_pattern from u
< ok 0 0
> exec ^
: string
| revoke file_prefix from u
< ok 0 0
> exec ^
: string
| revoke files from u
< ok 0 0
> exec ^
: string
| revoke filter from u
< ok 0 0
> exec ^
: string
| revoke finish from u
< ok 0 0
> exec ^
: string
| revoke first from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke first_value from u
< error
| 1064 42000 …line 1 column 18 near "first_value from u" 
> exec ^
: string
| revoke fixed from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke flashback from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke float from u
< error
| 1064 42000 …line 1 column 12 near "float from u" 
> exec ^
: string
| revoke float4 from u
< error
| 1064 42000 …line 1 column 13 near "float4 from u" 
> exec ^
: string
| revoke float8 from u
< error
| 1064 42000 …line 1 column 13 near "float8 from u" 
> exec ^
: string
| revoke flush from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke follower_constraints from u
< error
| 1064 42000 …line 1 column 32 near "from u" 
> exec ^
: string
| revoke followers from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke following from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke follows from u
< ok 0 0
> exec ^
: string
| revoke for from u
< error
| 1064 42000 …line 1 column 10 near "for from u" 
> exec ^
: string
| revoke force from u
< error
| 1064 42000 …line 1 column 12 near "force from u" 
> exec ^
: string
| revoke foreign from u
< error
| 1064 42000 …line 1 column 14 near "foreign from u" 
> exec ^
: string
| revoke format from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke found from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke from from u
< error
| 1064 42000 …line 1 column 11 near "from from u" 
> exec ^
: string
| revoke full from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke fulltext from u
< error
| 1064 42000 …line 1 column 15 near "fulltext from u" 
> exec ^
: string
| revoke function from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke general from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke generate from u
< ok 0 0
> exec ^
: string
| revoke generated from u
< error
| 1064 42000 …line 1 column 16 near "generated from u" 
> exec ^
: string
| revoke geomcollection from u
< ok 0 0
> exec ^
: string
| revoke geometry from u
< ok 0 0
> exec ^
: string
| revoke geometrycollection from u
< ok 0 0
> exec ^
: string
| revoke get from u
< ok 0 0
> exec ^
: string
| revoke get_format from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke get_source_public_key from u
< ok 0 0
> exec ^
: string
| revoke global from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke grant from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke grants from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke greatest from u
< ok 0 0
> exec ^
: string
| revoke group from u
< error
| 1064 42000 …line 1 column 12 near "group from u" 
> exec ^
: string
| revoke group_concat from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke group_replication from u
< ok 0 0
> exec ^
: string
| revoke grouping from u
< ok 0 0
> exec ^
: string
| revoke groups from u
< error
| 1064 42000 …line 1 column 13 near "groups from u" 
> exec ^
: string
| revoke gtid_only from u
< ok 0 0
> exec ^
: string
| revoke gtids from u
< ok 0 0
> exec ^
: string
| revoke guided from u
< ok 0 0
> exec ^
: string
| revoke handler from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke hash from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke having from u
< error
| 1064 42000 …line 1 column 13 near "having from u" 
> exec ^
: string
| revoke header from u
< ok 0 0
> exec ^
: string
| revoke help from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke high from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke high_priority from u
< error
| 1064 42000 …line 1 column 20 near "high_priority from u" 
> exec ^
: string
| revoke histogram from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke histograms_in_flight from u
< error
| 1064 42000 …line 1 column 32 near "from u" 
> exec ^
: string
| revoke history from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke host from u
< ok 0 0
> exec ^
: string
| revoke hosts from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke hour from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke hour_microsecond from u
< error
| 1064 42000 …line 1 column 23 near "hour_microsecond from u" 
> exec ^
: string
| revoke hour_minute from u
< error
| 1064 42000 …line 1 column 18 near "hour_minute from u" 
> exec ^
: string
| revoke hour_second from u
< error
| 1064 42000 …line 1 column 18 near "hour_second from u" 
> exec ^
: string
| revoke identified from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke if from u
< error
| 1064 42000 …line 1 column 9 near "if from u" 
> exec ^
: string
| revoke ignore from u
< error
| 1064 42000 …line 1 column 13 near "ignore from u" 
> exec ^
: string
| revoke ignore_server_ids from u
< ok 0 0
> exec ^
: string
| revoke ilike from u
< error
| 1064 42000 …line 1 column 12 near "ilike from u" 
> exec ^
: string
| revoke import from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke in from u
< error
| 1064 42000 …line 1 column 9 near "in from u" 
> exec ^
: string
| revoke inactive from u
< ok 0 0
> exec ^
: string
| revoke increment from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke index from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke indexes from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke infile from u
< error
| 1064 42000 …line 1 column 13 near "infile from u" 
> exec ^
: string
| revoke initial from u
< ok 0 0
> exec ^
: string
| revoke initial_size from u
< ok 0 0
> exec ^
: string
| revoke initiate from u
< ok 0 0
> exec ^
: string
| revoke inner from u
< error
| 1064 42000 …line 1 column 12 near "inner from u" 
> exec ^
: string
| revoke inout from u
< error
| 1064 42000 …line 1 column 12 near "inout from u" 
> exec ^
: string
| revoke insensitive from u
< ok 0 0
> exec ^
: string
| revoke insert from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke insert_method from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke install from u
< ok 0 0
> exec ^
: string
| revoke instance from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke int from u
< error
| 1064 42000 …line 1 column 10 near "int from u" 
> exec ^
: string
| revoke int1 from u
< error
| 1064 42000 …line 1 column 11 near "int1 from u" 
> exec ^
: string
| revoke int2 from u
< error
| 1064 42000 …line 1 column 11 near "int2 from u" 
> exec ^
: string
| revoke int3 from u
< error
| 1064 42000 …line 1 column 11 near "int3 from u" 
> exec ^
: string
| revoke int4 from u
< error
| 1064 42000 …line 1 column 11 near "int4 from u" 
> exec ^
: string
| revoke int8 from u
< error
| 1064 42000 …line 1 column 11 near "int8 from u" 
> exec ^
: string
| revoke integer from u
< error
| 1064 42000 …line 1 column 14 near "integer from u" 
> exec ^
: string
| revoke intersect from u
< error
| 1064 42000 …line 1 column 16 near "intersect from u" 
> exec ^
: string
| revoke interval from u
< error
| 1064 42000 …line 1 column 15 near "interval from u" 
> exec ^
: string
| revoke into from u
< error
| 1064 42000 …line 1 column 11 near "into from u" 
> exec ^
: string
| revoke invisible from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke invoker from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke io from u
< error
| 1064 42000 …line 1 column 14 near "from u" 
> exec ^
: string
| revoke io_after_gtids from u
< ok 0 0
> exec ^
: string
| revoke io_before_gtids from u
< ok 0 0
> exec ^
: string
| revoke io_thread from u
< ok 0 0
> exec ^
: string
| revoke ipc from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke is from u
< error
| 1064 42000 …line 1 column 9 near "is from u" 
> exec ^
: string
| revoke isolation from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke issuer from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke iterate from u
< error
| 1064 42000 …line 1 column 14 near "iterate from u" 
> exec ^
: string
| revoke job from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke jobs from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke join from u
< error
| 1064 42000 …line 1 column 11 near "join from u" 
> exec ^
: string
| revoke json from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke json_arrayagg from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke json_duality_object from u
< ok 0 0
> exec ^
: string
| revoke json_objectagg from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke json_table from u
< ok 0 0
> exec ^
: string
| revoke json_value from u
< ok 0 0
> exec ^
: string
| revoke key from u
< error
| 1064 42000 …line 1 column 10 near "key from u" 
> exec ^
: string
| revoke key_block_size from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke keyring from u
< ok 0 0
> exec ^
: string
| revoke keys from u
< error
| 1064 42000 …line 1 column 11 near "keys from u" 
> exec ^
: string
| revoke kill from u
< error
| 1064 42000 …line 1 column 11 near "kill from u" 
> exec ^
: string
| revoke labels from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke lag from u
< error
| 1064 42000 …line 1 column 10 near "lag from u" 
> exec ^
: string
| revoke language from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke last from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke last_value from u
< error
| 1064 42000 …line 1 column 17 near "last_value from u" 
> exec ^
: string
| revoke lateral from u
< ok 0 0
> exec ^
: string
| revoke lead from u
< error
| 1064 42000 …line 1 column 11 near "lead from u" 
> exec ^
: string
| revoke leader_constraints from u
< error
| 1064 42000 …line 1 column 30 near "from u" 
> exec ^
: string
| revoke leading from u
< error
| 1064 42000 …line 1 column 14 near "leading from u" 
> exec ^
: string
| revoke learner_constraints from u
< error
| 1064 42000 …line 1 column 31 near "from u" 
> exec ^
: string
| revoke learners from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke least from u
< ok 0 0
> exec ^
: string
| revoke leave from u
< error
| 1064 42000 …line 1 column 12 near "leave from u" 
> exec ^
: string
| revoke leaves from u
< ok 0 0
> exec ^
: string
| revoke left from u
< error
| 1064 42000 …line 1 column 11 near "left from u" 
> exec ^
: string
| revoke less from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke level from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke library from u
< ok 0 0
> exec ^
: string
| revoke like from u
< error
| 1064 42000 …line 1 column 11 near "like from u" 
> exec ^
: string
| revoke limit from u
< error
| 1064 42000 …line 1 column 12 near "limit from u" 
> exec ^
: string
| revoke linear from u
< error
| 1064 42000 …line 1 column 13 near "linear from u" 
> exec ^
: string
| revoke lines from u
< error
| 1064 42000 …line 1 column 12 near "lines from u" 
> exec ^
: string
| revoke linestring from u
< ok 0 0
> exec ^
: string
| revoke list from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke load from u
< error
| 1064 42000 …line 1 column 18 near "u" 
> exec ^
: string
| revoke local from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke localtime from u
< error
| 1064 42000 …line 1 column 16 near "localtime from u" 
> exec ^
: string
| revoke localtimestamp from u
< error
| 1064 42000 …line 1 column 21 near "localtimestamp from u" 
> exec ^
: string
| revoke location from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke lock from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke locked from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke locks from u
< ok 0 0
> exec ^
: string
| revoke log from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke logfile from u
< ok 0 0
> exec ^
: string
| revoke logs from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke long from u
< error
| 1064 42000 …line 1 column 11 near "long from u" 
> exec ^
: string
| revoke longblob from u
< error
| 1064 42000 …line 1 column 15 near "longblob from u" 
> exec ^
: string
| revoke longtext from u
< error
| 1064 42000 …line 1 column 15 near "longtext from u" 
> exec ^
: string
| revoke loop from u
< ok 0 0
> exec ^
: string
| revoke low from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke low_priority from u
< error
| 1064 42000 …line 1 column 19 near "low_priority from u" 
> exec ^
: string
| revoke manual from u
< ok 0 0
> exec ^
: string
| revoke master from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke master_log_file from u
< ok 0 0
> exec ^
: string
| revoke match from u
< error
| 1064 42000 …line 1 column 12 near "match from u" 
> exec ^
: string
| revoke materialized from u
< ok 0 0
> exec ^
: string
| revoke max from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke max_connections_per_hour from u
< error
| 1064 42000 …line 1 column 36 near "from u" 
> exec ^
: string
| revoke max_queries_per_hour from u
< error
| 1064 42000 …line 1 column 32 near "from u" 
> exec ^
: string
| revoke max_rows from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke max_size from u
< ok 0 0
> exec ^
: string
| revoke max_updates_per_hour from u
< error
| 1064 42000 …line 1 column 32 near "from u" 
> exec ^
: string
| revoke max_user_connections from u
< error
| 1064 42000 …line 1 column 32 near "from u" 
> exec ^
: string
| revoke maxvalue from u
< error
| 1064 42000 …line 1 column 15 near "maxvalue from u" 
> exec ^
: string
| revoke medium from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke mediumblob from u
< error
| 1064 42000 …line 1 column 17 near "mediumblob from u" 
> exec ^
: string
| revoke mediumint from u
< error
| 1064 42000 …line 1 column 16 near "mediumint from u" 
> exec ^
: string
| revoke mediumtext from u
< error
| 1064 42000 …line 1 column 17 near "mediumtext from u" 
> exec ^
: string
| revoke member from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke memory from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke merge from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke message_text from u
< ok 0 0
> exec ^
: string
| revoke microsecond from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke mid from u
< ok 0 0
> exec ^
: string
| revoke middleint from u
< error
| 1064 42000 …line 1 column 16 near "middleint from u" 
> exec ^
: string
| revoke migrate from u
< ok 0 0
> exec ^
: string
| revoke min from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke min_rows from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke minute from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke minute_microsecond from u
< error
| 1064 42000 …line 1 column 25 near "minute_microsecond from u" 
> exec ^
: string
| revoke minute_second from u
< error
| 1064 42000 …line 1 column 20 near "minute_second from u" 
> exec ^
: string
| revoke minvalue from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke mod from u
< error
| 1064 42000 …line 1 column 10 near "mod from u" 
> exec ^
: string
| revoke mode from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke modifies from u
< ok 0 0
> exec ^
: string
| revoke modify from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke month from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke multilinestring from u
< ok 0 0
> exec ^
: string
| revoke multipoint from u
< ok 0 0
> exec ^
: string
| revoke multipolygon from u
< ok 0 0
> exec ^
: string
| revoke mutex from u
< ok 0 0
> exec ^
: string
| revoke mysql_errno from u
< ok 0 0
> exec ^
: string
| revoke name from u
< ok 0 0
> exec ^
: string
| revoke names from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke national from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke natural from u
< error
| 1064 42000 …line 1 column 14 near "natural from u" 
> exec ^
: string
| revoke nchar from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke ndb from u
< ok 0 0
> exec ^
: string
| revoke ndbcluster from u
< ok 0 0
> exec ^
: string
| revoke nested from u
< ok 0 0
> exec ^
: string
| revoke network_namespace from u
< ok 0 0
> exec ^
: string
| revoke never from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke new from u
< ok 0 0
> exec ^
: string
| revoke next from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke next_row_id from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke no from u
< error
| 1064 42000 …line 1 column 14 near "from u" 
> exec ^
: string
| revoke no_wait from u
< ok 0 0
> exec ^
: string
| revoke no_write_to_binlog from u
< error
| 1064 42000 …line 1 column 25 near "no_write_to_binlog from u" 
> exec ^
: string
| revoke nocache from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke nocycle from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke nodegroup from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke nomaxvalue from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke nominvalue from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke nonclustered from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke none from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke not from u
< error
| 1064 42000 …line 1 column 10 near "not from u" 
> exec ^
: string
| revoke now from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke nowait from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke nth_value from u
< error
| 1064 42000 …line 1 column 16 near "nth_value from u" 
> exec ^
: string
| revoke ntile from u
< error
| 1064 42000 …line 1 column 12 near "ntile from u" 
> exec ^
: string
| revoke null from u
< error
| 1064 42000 …line 1 column 11 near "null from u" 
> exec ^
: string
| revoke nullif from u
< ok 0 0
> exec ^
: string
| revoke nulls from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke number from u
< ok 0 0
> exec ^
: string
| revoke numeric from u
< error
| 1064 42000 …line 1 column 14 near "numeric from u" 
> exec ^
: string
| revoke nvarchar from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke of from u
< error
| 1064 42000 …line 1 column 9 near "of from u" 
> exec ^
: string
| revoke off from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke offset from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke oj from u
< ok 0 0
> exec ^
: string
| revoke old from u
< ok 0 0
> exec ^
: string
| revoke oltp_read_only from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke oltp_read_write from u
< error
| 1064 42000 …line 1 column 27 near "from u" 
> exec ^
: string
| revoke oltp_write_only from u
< error
| 1064 42000 …line 1 column 27 near "from u" 
> exec ^
: string
| revoke on from u
< error
| 1064 42000 …line 1 column 9 near "on from u" 
> exec ^
: string
| revoke one from u
< ok 0 0
> exec ^
: string
| revoke only from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke open from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke optimize from u
< error
| 1064 42000 …line 1 column 15 near "optimize from u" 
> exec ^
: string
| revoke optimizer_costs from u
< ok 0 0
> exec ^
: string
| revoke option from u
< error
| 1064 42000 …line 1 column 13 near "option from u" 
> exec ^
: string
| revoke optional from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke optionally from u
< error
| 1064 42000 …line 1 column 17 near "optionally from u" 
> exec ^
: string
| revoke options from u
< ok 0 0
> exec ^
: string
| revoke or from u
< error
| 1064 42000 …line 1 column 9 near "or from u" 
> exec ^
: string
| revoke order from u
< error
| 1064 42000 …line 1 column 12 near "order from u" 
> exec ^
: string
| revoke ordinality from u
< ok 0 0
> exec ^
: string
| revoke organization from u
< ok 0 0
> exec ^
: string
| revoke others from u
< ok 0 0
> exec ^
: string
| revoke out from u
< error
| 1064 42000 …line 1 column 10 near "out from u" 
> exec ^
: string
| revoke outer from u
< error
| 1064 42000 …line 1 column 12 near "outer from u" 
> exec ^
: string
| revoke outfile from u
< error
| 1064 42000 …line 1 column 14 near "outfile from u" 
> exec ^
: string
| revoke over from u
< error
| 1064 42000 …line 1 column 11 near "over from u" 
> exec ^
: string
| revoke overlay from u
< ok 0 0
> exec ^
: string
| revoke owner from u
< ok 0 0
> exec ^
: string
| revoke pack_keys from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke page from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke parallel from u
< ok 0 0
> exec ^
: string
| revoke parameters from u
< ok 0 0
> exec ^
: string
| revoke parse_tree from u
< ok 0 0
> exec ^
: string
| revoke parser from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke partial from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke partition from u
< error
| 1064 42000 …line 1 column 16 near "partition from u" 
> exec ^
: string
| revoke partitioning from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke partitions from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke password from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke password_lock_time from u
< error
| 1064 42000 …line 1 column 30 near "from u" 
> exec ^
: string
| revoke path from u
< ok 0 0
> exec ^
: string
| revoke percent_rank from u
< error
| 1064 42000 …line 1 column 19 near "percent_rank from u" 
> exec ^
: string
| revoke persist from u
< ok 0 0
> exec ^
: string
| revoke persist_only from u
< ok 0 0
> exec ^
: string
| revoke phase from u
< ok 0 0
> exec ^
: string
| revoke placement from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke placing from u
< ok 0 0
> exec ^
: string
| revoke plan from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke plugin from u
< ok 0 0
> exec ^
: string
| revoke plugin_dir from u
< ok 0 0
> exec ^
: string
| revoke plugins from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke point from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke policy from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke polygon from u
< ok 0 0
> exec ^
: string
| revoke port from u
< ok 0 0
> exec ^
: string
| revoke position from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke pre_split_regions from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke precedes from u
< ok 0 0
> exec ^
: string
| revoke preceding from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke precision from u
< error
| 1064 42000 …line 1 column 16 near "precision from u" 
> exec ^
: string
| revoke prepare from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke preserve from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke prev from u
< ok 0 0
> exec ^
: string
| revoke primary from u
< error
| 1064 42000 …line 1 column 14 near "primary from u" 
> exec ^
: string
| revoke primary_region from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke priority from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke privilege_checks_user from u
< ok 0 0
> exec ^
: string
| revoke privileges from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke procedure from u
< error
| 1064 42000 …line 1 column 16 near "procedure from u" 
> exec ^
: string
| revoke process from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke processed_keys from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke processlist from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke profile from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke profiles from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke proxy from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke purge from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke qualify from u
< ok 0 0
> exec ^
: string
| revoke quarter from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke query from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke query_limit from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke quick from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke random from u
< ok 0 0
> exec ^
: string
| revoke range from u
< error
| 1064 42000 …line 1 column 12 near "range from u" 
> exec ^
: string
| revoke rank from u
< error
| 1064 42000 …line 1 column 11 near "rank from u" 
> exec ^
: string
| revoke read from u
< error
| 1064 42000 …line 1 column 11 near "read from u" 
> exec ^
: string
| revoke read_only from u
< ok 0 0
> exec ^
: string
| revoke read_write from u
< ok 0 0
> exec ^
: string
| revoke reads from u
< ok 0 0
> exec ^
: string
| revoke real from u
< error
| 1064 42000 …line 1 column 11 near "real from u" 
> exec ^
: string
| revoke rebuild from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke recover from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke recursive from u
< error
| 1064 42000 …line 1 column 16 near "recursive from u" 
> exec ^
: string
| revoke redo_buffer_size from u
< ok 0 0
> exec ^
: string
| revoke redundant from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke reference from u
< ok 0 0
> exec ^
: string
| revoke references from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke regexp from u
< error
| 1064 42000 …line 1 column 13 near "regexp from u" 
> exec ^
: string
| revoke region from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke regions from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke registration from u
< ok 0 0
> exec ^
: string
| revoke relational from u
< ok 0 0
> exec ^
: string
| revoke relay from u
< ok 0 0
> exec ^
: string
| revoke relay_log_file from u
< ok 0 0
> exec ^
: string
| revoke relay_log_pos from u
< ok 0 0
> exec ^
: string
| revoke relay_thread from u
< ok 0 0
> exec ^
: string
| revoke relaylog from u
< ok 0 0
> exec ^
: string
| revoke release from u
< error
| 1064 42000 …line 1 column 14 near "release from u" 
> exec ^
: string
| revoke reload from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke remove from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke rename from u
< error
| 1064 42000 …line 1 column 13 near "rename from u" 
> exec ^
: string
| revoke reorganize from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke repair from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke repeat from u
< error
| 1064 42000 …line 1 column 13 near "repeat from u" 
> exec ^
: string
| revoke repeatable from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke replace from u
< error
| 1064 42000 …line 1 column 14 near "replace from u" 
> exec ^
: string
| revoke replica from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke replicas from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke replicate_do_db from u
< ok 0 0
> exec ^
: string
| revoke replicate_do_table from u
< ok 0 0
> exec ^
: string
| revoke replicate_ignore_db from u
< ok 0 0
> exec ^
: string
| revoke replicate_ignore_table from u
< ok 0 0
> exec ^
: string
| revoke replicate_rewrite_db from u
< ok 0 0
> exec ^
: string
| revoke replicate_wild_do_table from u
< ok 0 0
> exec ^
: string
| revoke replicate_wild_ignore_table from u
< ok 0 0
> exec ^
: string
| revoke replication from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke require from u
< error
| 1064 42000 …line 1 column 14 near "require from u" 
> exec ^
: string
| revoke require_row_format from u
< ok 0 0
> exec ^
: string
| revoke require_table_primary_key_check from u
< ok 0 0
> exec ^
: string
| revoke reset from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke resignal from u
< ok 0 0
> exec ^
: string
| revoke resource from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke respect from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke restart from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke restore from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke restores from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke restrict from u
< error
| 1064 42000 …line 1 column 15 near "restrict from u" 
> exec ^
: string
| revoke resume from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke retain from u
< ok 0 0
> exec ^
: string
| revoke return from u
< ok 0 0
> exec ^
: string
| revoke returned_sqlstate from u
< ok 0 0
> exec ^
: string
| revoke returning from u
< ok 0 0
> exec ^
: string
| revoke returns from u
< ok 0 0
> exec ^
: string
| revoke reuse from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke reverse from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke revoke from u
< error
| 1064 42000 …line 1 column 13 near "revoke from u" 
> exec ^
: string
| revoke right from u
< error
| 1064 42000 …line 1 column 12 near "right from u" 
> exec ^
: string
| revoke rlike from u
< error
| 1064 42000 …line 1 column 12 near "rlike from u" 
> exec ^
: string
| revoke role from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke rollback from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke rollup from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke rotate from u
< ok 0 0
> exec ^
: string
| revoke routine from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke row from u
< error
| 1064 42000 …line 1 column 10 near "row from u" 
> exec ^
: string
| revoke row_count from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke row_format from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke row_number from u
< error
| 1064 42000 …line 1 column 17 near "row_number from u" 
> exec ^
: string
| revoke rows from u
< error
| 1064 42000 …line 1 column 11 near "rows from u" 
> exec ^
: string
| revoke rtree from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke ru from u
< error
| 1064 42000 …line 1 column 14 near "from u" 
> exec ^
: string
| revoke ru_per_sec from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke rule from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke run from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke s3 from u
< error
| 1064 42000 …line 1 column 14 near "from u" 
> exec ^
: string
| revoke savepoint from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke schedule from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke schema from u
< error
| 1064 42000 …line 1 column 13 near "schema from u" 
> exec ^
: string
| revoke schema_name from u
< ok 0 0
> exec ^
: string
| revoke schemas from u
< error
| 1064 42000 …line 1 column 14 near "schemas from u" 
> exec ^
: string
| revoke second from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke second_microsecond from u
< error
| 1064 42000 …line 1 column 25 near "second_microsecond from u" 
> exec ^
: string
| revoke secondary from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke secondary_engine from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke secondary_engine_attribute from u
< ok 0 0
> exec ^
: string
| revoke secondary_load from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke secondary_unload from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke security from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke select from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke sensitive from u
< ok 0 0
> exec ^
: string
| revoke separator from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke sequence from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke serial from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke serializable from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke server from u
< ok 0 0
> exec ^
: string
| revoke session from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke session_states from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke session_user from u
< ok 0 0
> exec ^
: string
| revoke set from u
< error
| 1064 42000 …line 1 column 10 near "set from u" 
> exec ^
: string
| revoke sets from u
< ok 0 0
> exec ^
: string
| revoke shard_row_id_bits from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke share from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke show from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke shutdown from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke signal from u
< ok 0 0
> exec ^
: string
| revoke signed from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke similar from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke simple from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke skip from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke slave from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke slow from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke smallint from u
< error
| 1064 42000 …line 1 column 15 near "smallint from u" 
> exec ^
: string
| revoke snapshot from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke socket from u
< ok 0 0
> exec ^
: string
| revoke some from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke soname from u
< ok 0 0
> exec ^
: string
| revoke sounds from u
< ok 0 0
> exec ^
: string
| revoke source from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke source_auto_position from u
< ok 0 0
> exec ^
: string
| revoke source_bind from u
< ok 0 0
> exec ^
: string
| revoke source_compression_algorithms from u
< ok 0 0
> exec ^
: string
| revoke source_connect_retry from u
< ok 0 0
> exec ^
: string
| revoke source_connection_auto_failover from u
< ok 0 0
> exec ^
: string
| revoke source_delay from u
< ok 0 0
> exec ^
: string
| revoke source_heartbeat_period from u
< ok 0 0
> exec ^
: string
| revoke source_host from u
< ok 0 0
> exec ^
: string
| revoke source_log_file from u
< ok 0 0
> exec ^
: string
| revoke source_log_pos from u
< ok 0 0
> exec ^
: string
| revoke source_password from u
< ok 0 0
> exec ^
: string
| revoke source_port from u
< ok 0 0
> exec ^
: string
| revoke source_public_key_path from u
< ok 0 0
> exec ^
: string
| revoke source_retry_count from u
< ok 0 0
> exec ^
: string
| revoke source_ssl from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_ca from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_capath from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_cert from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_cipher from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_crl from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_crlpath from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_key from u
< ok 0 0
> exec ^
: string
| revoke source_ssl_verify_server_cert from u
< ok 0 0
> exec ^
: string
| revoke source_tls_ciphersuites from u
< ok 0 0
> exec ^
: string
| revoke source_tls_version from u
< ok 0 0
> exec ^
: string
| revoke source_user from u
< ok 0 0
> exec ^
: string
| revoke source_zstd_compression_level from u
< ok 0 0
> exec ^
: string
| revoke spatial from u
< error
| 1064 42000 …line 1 column 14 near "spatial from u" 
> exec ^
: string
| revoke specific from u
< ok 0 0
> exec ^
: string
| revoke split from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke sql from u
< error
| 1064 42000 …line 1 column 10 near "sql from u" 
> exec ^
: string
| revoke sql_after_gtids from u
< ok 0 0
> exec ^
: string
| revoke sql_after_mts_gaps from u
< ok 0 0
> exec ^
: string
| revoke sql_before_gtids from u
< ok 0 0
> exec ^
: string
| revoke sql_big_result from u
< error
| 1064 42000 …line 1 column 21 near "sql_big_result from u" 
> exec ^
: string
| revoke sql_buffer_result from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke sql_calc_found_rows from u
< error
| 1064 42000 …line 1 column 26 near "sql_calc_found_rows from u" 
> exec ^
: string
| revoke sql_no_cache from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke sql_small_result from u
< error
| 1064 42000 …line 1 column 23 near "sql_small_result from u" 
> exec ^
: string
| revoke sql_thread from u
< ok 0 0
> exec ^
: string
| revoke sql_tsi_day from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke sql_tsi_hour from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke sql_tsi_minute from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke sql_tsi_month from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke sql_tsi_quarter from u
< error
| 1064 42000 …line 1 column 27 near "from u" 
> exec ^
: string
| revoke sql_tsi_second from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke sql_tsi_week from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke sql_tsi_year from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke sqlexception from u
< error
| 1064 42000 …line 1 column 19 near "sqlexception from u" 
> exec ^
: string
| revoke sqlstate from u
< error
| 1064 42000 …line 1 column 15 near "sqlstate from u" 
> exec ^
: string
| revoke sqlwarning from u
< error
| 1064 42000 …line 1 column 17 near "sqlwarning from u" 
> exec ^
: string
| revoke srid from u
< ok 0 0
> exec ^
: string
| revoke ssl from u
< error
| 1064 42000 …line 1 column 10 near "ssl from u" 
> exec ^
: string
| revoke st_collect from u
< ok 0 0
> exec ^
: string
| revoke stacked from u
< ok 0 0
> exec ^
: string
| revoke start from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke start_time from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke starting from u
< error
| 1064 42000 …line 1 column 15 near "starting from u" 
> exec ^
: string
| revoke starts from u
< ok 0 0
> exec ^
: string
| revoke stats_auto_recalc from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke stats_buckets from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke stats_col_choice from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke stats_col_list from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke stats_extended from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke stats_healthy from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke stats_histograms from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke stats_locked from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke stats_meta from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke stats_persistent from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke stats_sample_pages from u
< error
| 1064 42000 …line 1 column 30 near "from u" 
> exec ^
: string
| revoke stats_sample_rate from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke stats_topn from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke status from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke std from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke stddev from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke stddev_pop from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke stddev_samp from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke stop from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke storage from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke stored from u
< error
| 1064 42000 …line 1 column 13 near "stored from u" 
> exec ^
: string
| revoke straight_join from u
< error
| 1064 42000 …line 1 column 20 near "straight_join from u" 
> exec ^
: string
| revoke stream from u
< ok 0 0
> exec ^
: string
| revoke strict_load from u
< ok 0 0
> exec ^
: string
| revoke string from u
< ok 0 0
> exec ^
: string
| revoke subclass_origin from u
< ok 0 0
> exec ^
: string
| revoke subdate from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke subject from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke subpartition from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke subpartitions from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke substr from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke substring from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke sum from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke super from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke survival_preferences from u
< error
| 1064 42000 …line 1 column 32 near "from u" 
> exec ^
: string
| revoke suspend from u
< ok 0 0
> exec ^
: string
| revoke swaps from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke switch_group from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke switches from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke sysdate from u
< ok 0 0
> exec ^
: string
| revoke system from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke system_user from u
< ok 0 0
> exec ^
: string
| revoke table from u
< error
| 1064 42000 …line 1 column 12 near "table from u" 
> exec ^
: string
| revoke table_checksum from u
< error
| 1064 42000 …line 1 column 26 near "from u" 
> exec ^
: string
| revoke table_name from u
< ok 0 0
> exec ^
: string
| revoke tables from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke tablesample from u
< error
| 1064 42000 …line 1 column 18 near "tablesample from u" 
> exec ^
: string
| revoke tablespace from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke task_types from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke temporary from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke temptable from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke terminated from u
< error
| 1064 42000 …line 1 column 17 near "terminated from u" 
> exec ^
: string
| revoke text from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke than from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke then from u
< error
| 1064 42000 …line 1 column 11 near "then from u" 
> exec ^
: string
| revoke thread_priority from u
< ok 0 0
> exec ^
: string
| revoke ties from u
< ok 0 0
> exec ^
: string
| revoke tiflash from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke time from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke timeout from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke timestamp from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke timestampadd from u
< error
| 1064 42000 …line 1 column 24 near "from u" 
> exec ^
: string
| revoke timestampdiff from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke tinyblob from u
< error
| 1064 42000 …line 1 column 15 near "tinyblob from u" 
> exec ^
: string
| revoke tinyint from u
< error
| 1064 42000 …line 1 column 14 near "tinyint from u" 
> exec ^
: string
| revoke tinytext from u
< error
| 1064 42000 …line 1 column 15 near "tinytext from u" 
> exec ^
: string
| revoke tls from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke to from u
< error
| 1064 42000 …line 1 column 9 near "to from u" 
> exec ^
: string
| revoke tpcc from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke tpch_10 from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke trailing from u
< error
| 1064 42000 …line 1 column 15 near "trailing from u" 
> exec ^
: string
| revoke transaction from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke trigger from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke triggers from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke trim from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke true from u
< error
| 1064 42000 …line 1 column 11 near "true from u" 
> exec ^
: string
| revoke truncate from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke ttl from u
< error
| 1064 42000 …line 1 column 15 near "from u" 
> exec ^
: string
| revoke ttl_enable from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke ttl_job_interval from u
< error
| 1064 42000 …line 1 column 28 near "from u" 
> exec ^
: string
| revoke type from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke types from u
< ok 0 0
> exec ^
: string
| revoke unbounded from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke uncommitted from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| revoke undefined from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke undo from u
< ok 0 0
> exec ^
: string
| revoke undo_buffer_size from u
< ok 0 0
> exec ^
: string
| revoke undofile from u
< ok 0 0
> exec ^
: string
| revoke unicode from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke uninstall from u
< ok 0 0
> exec ^
: string
| revoke union from u
< error
| 1064 42000 …line 1 column 12 near "union from u" 
> exec ^
: string
| revoke unique from u
< error
| 1064 42000 …line 1 column 13 near "unique from u" 
> exec ^
: string
| revoke unknown from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke unlimited from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke unlock from u
< error
| 1064 42000 …line 1 column 13 near "unlock from u" 
> exec ^
: string
| revoke unregister from u
< ok 0 0
> exec ^
: string
| revoke unsigned from u
< error
| 1064 42000 …line 1 column 15 near "unsigned from u" 
> exec ^
: string
| revoke until from u
< error
| 1064 42000 …line 1 column 12 near "until from u" 
> exec ^
: string
| revoke update from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke upgrade from u
< ok 0 0
> exec ^
: string
| revoke uri from u
< ok 0 0
> exec ^
: string
| revoke url from u
< ok 0 0
> exec ^
: string
| revoke usage from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke use from u
< error
| 1064 42000 …line 1 column 10 near "use from u" 
> exec ^
: string
| revoke use_frm from u
< ok 0 0
> exec ^
: string
| revoke user from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke user_resources from u
< ok 0 0
> exec ^
: string
| revoke using from u
< error
| 1064 42000 …line 1 column 12 near "using from u" 
> exec ^
: string
| revoke utc_date from u
< error
| 1064 42000 …line 1 column 15 near "utc_date from u" 
> exec ^
: string
| revoke utc_time from u
< error
| 1064 42000 …line 1 column 15 near "utc_time from u" 
> exec ^
: string
| revoke utc_timestamp from u
< error
| 1064 42000 …line 1 column 20 near "utc_timestamp from u" 
> exec ^
: string
| revoke utilization_limit from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke validate from u
< ok 0 0
> exec ^
: string
| revoke validation from u
< error
| 1064 42000 …line 1 column 22 near "from u" 
> exec ^
: string
| revoke value from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke values from u
< error
| 1064 42000 …line 1 column 13 near "values from u" 
> exec ^
: string
| revoke var_pop from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke var_samp from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke varbinary from u
< error
| 1064 42000 …line 1 column 16 near "varbinary from u" 
> exec ^
: string
| revoke varchar from u
< error
| 1064 42000 …line 1 column 14 near "varchar from u" 
> exec ^
: string
| revoke varcharacter from u
< error
| 1064 42000 …line 1 column 19 near "varcharacter from u" 
> exec ^
: string
| revoke variables from u
< error
| 1064 42000 …line 1 column 21 near "from u" 
> exec ^
: string
| revoke variance from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke varying from u
< error
| 1064 42000 …line 1 column 14 near "varying from u" 
> exec ^
: string
| revoke vcpu from u
< ok 0 0
> exec ^
: string
| revoke vector from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke verify_key_constraints from u
< ok 0 0
> exec ^
: string
| revoke view from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke virtual from u
< error
| 1064 42000 …line 1 column 14 near "virtual from u" 
> exec ^
: string
| revoke visible from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke voter_constraints from u
< error
| 1064 42000 …line 1 column 29 near "from u" 
> exec ^
: string
| revoke voters from u
< error
| 1064 42000 …line 1 column 18 near "from u" 
> exec ^
: string
| revoke wait from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke warnings from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke watch from u
< error
| 1064 42000 …line 1 column 17 near "from u" 
> exec ^
: string
| revoke week from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke weight_string from u
< error
| 1064 42000 …line 1 column 25 near "from u" 
> exec ^
: string
| revoke when from u
< error
| 1064 42000 …line 1 column 11 near "when from u" 
> exec ^
: string
| revoke where from u
< error
| 1064 42000 …line 1 column 12 near "where from u" 
> exec ^
: string
| revoke while from u
< error
| 1064 42000 …line 1 column 12 near "while from u" 
> exec ^
: string
| revoke window from u
< error
| 1064 42000 …line 1 column 13 near "window from u" 
> exec ^
: string
| revoke with from u
< error
| 1064 42000 …line 1 column 11 near "with from u" 
> exec ^
: string
| revoke without from u
< error
| 1064 42000 …line 1 column 19 near "from u" 
> exec ^
: string
| revoke work from u
< ok 0 0
> exec ^
: string
| revoke workload from u
< error
| 1064 42000 …line 1 column 20 near "from u" 
> exec ^
: string
| revoke wrapper from u
< ok 0 0
> exec ^
: string
| revoke write from u
< error
| 1064 42000 …line 1 column 12 near "write from u" 
> exec ^
: string
| revoke x509 from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke xa from u
< ok 0 0
> exec ^
: string
| revoke xid from u
< ok 0 0
> exec ^
: string
| revoke xml from u
< ok 0 0
> exec ^
: string
| revoke xor from u
< error
| 1064 42000 …line 1 column 10 near "xor from u" 
> exec ^
: string
| revoke year from u
< error
| 1064 42000 …line 1 column 16 near "from u" 
> exec ^
: string
| revoke year_month from u
< error
| 1064 42000 …line 1 column 17 near "year_month from u" 
> exec ^
: string
| revoke zerofill from u
< error
| 1064 42000 …line 1 column 15 near "zerofill from u" 
> exec ^
: string
| revoke zone from u
< ok 0 0
> exec ^
: string
| grant absent on *.* to u
< ok 0 0
> exec ^
: string
| grant accessible on *.* to u
< ok 0 0
> exec ^
: string
| grant account on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant action on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant active on *.* to u
< ok 0 0
> exec ^
: string
| grant add on *.* to u
< error
| 1064 42000 …line 1 column 9 near "add on *.* to u" 
> exec ^
: string
| grant adddate on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant admin on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant after on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant against on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant aggregate on *.* to u
< ok 0 0
> exec ^
: string
| grant algorithm on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant all on *.* to u
< ok 0 0
> exec ^
: string
| grant allow_missing_files on *.* to u
< ok 0 0
> exec ^
: string
| grant alter on *.* to u
< ok 0 0
> exec ^
: string
| grant always on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant analyze on *.* to u
< error
| 1064 42000 …line 1 column 13 near "analyze on *.* to u" 
> exec ^
: string
| grant and on *.* to u
< error
| 1064 42000 …line 1 column 9 near "and on *.* to u" 
> exec ^
: string
| grant any on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant array on *.* to u
< error
| 1064 42000 …line 1 column 11 near "array on *.* to u" 
> exec ^
: string
| grant as on *.* to u
< error
| 1064 42000 …line 1 column 8 near "as on *.* to u" 
> exec ^
: string
| grant asc on *.* to u
< error
| 1064 42000 …line 1 column 9 near "asc on *.* to u" 
> exec ^
: string
| grant ascii on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant asensitive on *.* to u
< ok 0 0
> exec ^
: string
| grant assign_gtids_to_anonymous_transactions on *.* to u
< ok 0 0
> exec ^
: string
| grant at on *.* to u
< ok 0 0
> exec ^
: string
| grant attribute on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant authentication on *.* to u
< ok 0 0
> exec ^
: string
| grant auto on *.* to u
< ok 0 0
> exec ^
: string
| grant auto_id_cache on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant auto_increment on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant auto_random on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant auto_random_base on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant auto_refresh on *.* to u
< ok 0 0
> exec ^
: string
| grant auto_refresh_source on *.* to u
< ok 0 0
> exec ^
: string
| grant autoextend_size on *.* to u
< ok 0 0
> exec ^
: string
| grant avg on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant avg_row_length on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant background on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant backup on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant backups on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant batch on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant before on *.* to u
< ok 0 0
> exec ^
: string
| grant begin on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant bernoulli on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant between on *.* to u
< error
| 1064 42000 …line 1 column 13 near "between on *.* to u" 
> exec ^
: string
| grant bigint on *.* to u
< error
| 1064 42000 …line 1 column 12 near "bigint on *.* to u" 
> exec ^
: string
| grant binary on *.* to u
< error
| 1064 42000 …line 1 column 12 near "binary on *.* to u" 
> exec ^
: string
| grant binding on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant binding_cache on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant bindings on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant binlog on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant bit on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant bit_and on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant bit_or on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant bit_xor on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant blob on *.* to u
< error
| 1064 42000 …line 1 column 10 near "blob on *.* to u" 
> exec ^
: string
| grant block on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant bool on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant boolean on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant both on *.* to u
< error
| 1064 42000 …line 1 column 10 near "both on *.* to u" 
> exec ^
: string
| grant btree on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant buckets on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant builtins on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant bulk on *.* to u
< ok 0 0
> exec ^
: string
| grant burstable on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant by on *.* to u
< error
| 1064 42000 …line 1 column 8 near "by on *.* to u" 
> exec ^
: string
| grant byte on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant cache on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant calibrate on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant call on *.* to u
< error
| 1064 42000 …line 1 column 10 near "call on *.* to u" 
> exec ^
: string
| grant cascade on *.* to u
< error
| 1064 42000 …line 1 column 13 near "cascade on *.* to u" 
> exec ^
: string
| grant cascaded on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant case on *.* to u
< error
| 1064 42000 …line 1 column 10 near "case on *.* to u" 
> exec ^
: string
| grant cast on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant catalog_name on *.* to u
< ok 0 0
> exec ^
: string
| grant chain on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant challenge_response on *.* to u
< ok 0 0
> exec ^
: string
| grant change on *.* to u
< error
| 1064 42000 …line 1 column 12 near "change on *.* to u" 
> exec ^
: string
| grant changed on *.* to u
< ok 0 0
> exec ^
: string
| grant channel on *.* to u
< ok 0 0
> exec ^
: string
| grant char on *.* to u
< error
| 1064 42000 …line 1 column 10 near "char on *.* to u" 
> exec ^
: string
| grant character on *.* to u
< error
| 1064 42000 …line 1 column 15 near "character on *.* to u" 
> exec ^
: string
| grant charset on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant check on *.* to u
< error
| 1064 42000 …line 1 column 11 near "check on *.* to u" 
> exec ^
: string
| grant checksum on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant cipher on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant class_origin on *.* to u
< ok 0 0
> exec ^
: string
| grant client on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant clone on *.* to u
< ok 0 0
> exec ^
: string
| grant close on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant clustered on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant coalesce on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant code on *.* to u
< ok 0 0
> exec ^
: string
| grant collate on *.* to u
< error
| 1064 42000 …line 1 column 13 near "collate on *.* to u" 
> exec ^
: string
| grant collation on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant column on *.* to u
< error
| 1064 42000 …line 1 column 12 near "column on *.* to u" 
> exec ^
: string
| grant column_format on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant column_name on *.* to u
< ok 0 0
> exec ^
: string
| grant column_stats_usage on *.* to u
< error
| 1064 42000 …line 1 column 27 near "on *.* to u" 
> exec ^
: string
| grant columns on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant comment on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant commit on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant committed on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant compact on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant completion on *.* to u
< ok 0 0
> exec ^
: string
| grant component on *.* to u
< ok 0 0
> exec ^
: string
| grant compressed on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant compression on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant concat on *.* to u
< ok 0 0
> exec ^
: string
| grant concurrent on *.* to u
< ok 0 0
> exec ^
: string
| grant condition on *.* to u
< ok 0 0
> exec ^
: string
| grant config on *.* to u
< ok 0 0
> exec ^
: string
| grant connection on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant consistent on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant constraint on *.* to u
< error
| 1064 42000 …line 1 column 16 near "constraint on *.* to u" 
> exec ^
: string
| grant constraint_catalog on *.* to u
< ok 0 0
> exec ^
: string
| grant constraint_name on *.* to u
< ok 0 0
> exec ^
: string
| grant constraint_schema on *.* to u
< ok 0 0
> exec ^
: string
| grant constraints on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant contains on *.* to u
< ok 0 0
> exec ^
: string
| grant context on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant continue on *.* to u
< error
| 1064 42000 …line 1 column 14 near "continue on *.* to u" 
> exec ^
: string
| grant convert on *.* to u
< error
| 1064 42000 …line 1 column 13 near "convert on *.* to u" 
> exec ^
: string
| grant cooldown on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant count on *.* to u
< ok 0 0
> exec ^
: string
| grant cpu on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant create on *.* to u
< ok 0 0
> exec ^
: string
| grant cross on *.* to u
< error
| 1064 42000 …line 1 column 11 near "cross on *.* to u" 
> exec ^
: string
| grant cube on *.* to u
< ok 0 0
> exec ^
: string
| grant cume_dist on *.* to u
< error
| 1064 42000 …line 1 column 15 near "cume_dist on *.* to u" 
> exec ^
: string
| grant curdate on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant current on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant current_date on *.* to u
< error
| 1064 42000 …line 1 column 18 near "current_date on *.* to u" 
> exec ^
: string
| grant current_time on *.* to u
< error
| 1064 42000 …line 1 column 18 near "current_time on *.* to u" 
> exec ^
: string
| grant current_timestamp on *.* to u
< error
| 1064 42000 …line 1 column 23 near "current_timestamp on *.* to u" 
> exec ^
: string
| grant current_user on *.* to u
< error
| 1064 42000 …line 1 column 18 near "current_user on *.* to u" 
> exec ^
: string
| grant cursor on *.* to u
< error
| 1064 42000 …line 1 column 12 near "cursor on *.* to u" 
> exec ^
: string
| grant cursor_name on *.* to u
< ok 0 0
> exec ^
: string
| grant curtime on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant cycle on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant data on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant database on *.* to u
< error
| 1064 42000 …line 1 column 14 near "database on *.* to u" 
> exec ^
: string
| grant databases on *.* to u
< error
| 1064 42000 …line 1 column 15 near "databases on *.* to u" 
> exec ^
: string
| grant datafile on *.* to u
< ok 0 0
> exec ^
: string
| grant date on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant date_add on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant date_sub on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant datetime on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant day on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant day_hour on *.* to u
< error
| 1064 42000 …line 1 column 14 near "day_hour on *.* to u" 
> exec ^
: string
| grant day_microsecond on *.* to u
< error
| 1064 42000 …line 1 column 21 near "day_microsecond on *.* to u" 
> exec ^
: string
| grant day_minute on *.* to u
< error
| 1064 42000 …line 1 column 16 near "day_minute on *.* to u" 
> exec ^
: string
| grant day_second on *.* to u
< error
| 1064 42000 …line 1 column 16 near "day_second on *.* to u" 
> exec ^
: string
| grant deallocate on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant dec on *.* to u
< error
| 1064 42000 …line 1 column 9 near "dec on *.* to u" 
> exec ^
: string
| grant decimal on *.* to u
< error
| 1064 42000 …line 1 column 13 near "decimal on *.* to u" 
> exec ^
: string
| grant declare on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant default on *.* to u
< error
| 1064 42000 …line 1 column 13 near "default on *.* to u" 
> exec ^
: string
| grant default_auth on *.* to u
< ok 0 0
> exec ^
: string
| grant defined on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant definer on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant definition on *.* to u
< ok 0 0
> exec ^
: string
| grant delay_key_write on *.* to u
< error
| 1064 42000 …line 1 column 24 near "on *.* to u" 
> exec ^
: string
| grant delayed on *.* to u
< error
| 1064 42000 …line 1 column 13 near "delayed on *.* to u" 
> exec ^
: string
| grant delete on *.* to u
< ok 0 0
> exec ^
: string
| grant dense_rank on *.* to u
< error
| 1064 42000 …line 1 column 16 near "dense_rank on *.* to u" 
> exec ^
: string
| grant desc on *.* to u
< error
| 1064 42000 …line 1 column 10 near "desc on *.* to u" 
> exec ^
: string
| grant describe on *.* to u
< error
| 1064 42000 …line 1 column 14 near "describe on *.* to u" 
> exec ^
: string
| grant description on *.* to u
< ok 0 0
> exec ^
: string
| grant deterministic on *.* to u
< ok 0 0
> exec ^
: string
| grant diagnostics on *.* to u
< ok 0 0
> exec ^
: string
| grant digest on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant directory on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant disable on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant disabled on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant discard on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant disk on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant distinct on *.* to u
< error
| 1064 42000 …line 1 column 14 near "distinct on *.* to u" 
> exec ^
: string
| grant distinctrow on *.* to u
< error
| 1064 42000 …line 1 column 17 near "distinctrow on *.* to u" 
> exec ^
: string
| grant distribute on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant distribution on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant distributions on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant div on *.* to u
< error
| 1064 42000 …line 1 column 9 near "div on *.* to u" 
> exec ^
: string
| grant do on *.* to u
< error
| 1064 42000 …line 1 column 11 near "on *.* to u" 
> exec ^
: string
| grant double on *.* to u
< error
| 1064 42000 …line 1 column 12 near "double on *.* to u" 
> exec ^
: string
| grant drop on *.* to u
< ok 0 0
> exec ^
: string
| grant dry on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant dryrun on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant dual on *.* to u
< error
| 1064 42000 …line 1 column 10 near "dual on *.* to u" 
> exec ^
: string
| grant duality on *.* to u
< ok 0 0
> exec ^
: string
| grant dumpfile on *.* to u
< ok 0 0
> exec ^
: string
| grant duplicate on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant duration on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant dynamic on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant each on *.* to u
< ok 0 0
> exec ^
: string
| grant else on *.* to u
< error
| 1064 42000 …line 1 column 10 near "else on *.* to u" 
> exec ^
: string
| grant elseif on *.* to u
< error
| 1064 42000 …line 1 column 12 near "elseif on *.* to u" 
> exec ^
: string
| grant empty on *.* to u
< ok 0 0
> exec ^
: string
| grant enable on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant enabled on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant enclosed on *.* to u
< error
| 1064 42000 …line 1 column 14 near "enclosed on *.* to u" 
> exec ^
: string
| grant encryption on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant end on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant end_time on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant ends on *.* to u
< ok 0 0
> exec ^
: string
| grant enforced on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant engine on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant engine_attribute on *.* to u
< ok 0 0
> exec ^
: string
| grant engines on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant enum on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant error on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant errors on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant escape on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant escaped on *.* to u
< error
| 1064 42000 …line 1 column 13 near "escaped on *.* to u" 
> exec ^
: string
| grant event on *.* to u
< ok 0 0
> exec ^
: string
| grant events on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant every on *.* to u
< ok 0 0
> exec ^
: string
| grant exact on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant except on *.* to u
< error
| 1064 42000 …line 1 column 12 near "except on *.* to u" 
> exec ^
: string
| grant exchange on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant exclude on *.* to u
< ok 0 0
> exec ^
: string
| grant exec_elapsed on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant execute on *.* to u
< ok 0 0
> exec ^
: string
| grant exists on *.* to u
< error
| 1064 42000 …line 1 column 12 near "exists on *.* to u" 
> exec ^
: string
| grant exit on *.* to u
< error
| 1064 42000 …line 1 column 10 near "exit on *.* to u" 
> exec ^
: string
| grant expansion on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant expire on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant explain on *.* to u
< error
| 1064 42000 …line 1 column 13 near "explain on *.* to u" 
> exec ^
: string
| grant export on *.* to u
< ok 0 0
> exec ^
: string
| grant extended on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant extent_size on *.* to u
< ok 0 0
> exec ^
: string
| grant external on *.* to u
< ok 0 0
> exec ^
: string
| grant external_format on *.* to u
< ok 0 0
> exec ^
: string
| grant extract on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant factor on *.* to u
< ok 0 0
> exec ^
: string
| grant failed_login_attempts on *.* to u
< error
| 1064 42000 …line 1 column 30 near "on *.* to u" 
> exec ^
: string
| grant false on *.* to u
< error
| 1064 42000 …line 1 column 11 near "false on *.* to u" 
> exec ^
: string
| grant fast on *.* to u
< ok 0 0
> exec ^
: string
| grant faults on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant fetch on *.* to u
< error
| 1064 42000 …line 1 column 11 near "fetch on *.* to u" 
> exec ^
: string
| grant fields on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant file on *.* to u
< ok 0 0
> exec ^
: string
| grant file_block_size on *.* to u
< ok 0 0
> exec ^
: string
| grant file_format on *.* to u
< ok 0 0
> exec ^
: string
| grant file_name on *.* to u
< ok 0 0
> exec ^
: string
| grant file_pattern on *.* to u
< ok 0 0
> exec ^
: string
| grant file_prefix on *.* to u
< ok 0 0
> exec ^
: string
| grant files on *.* to u
< ok 0 0
> exec ^
: string
| grant filter on *.* to u
< ok 0 0
> exec ^
: string
| grant finish on *.* to u
< ok 0 0
> exec ^
: string
| grant first on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant first_value on *.* to u
< error
| 1064 42000 …line 1 column 17 near "first_value on *.* to u" 
> exec ^
: string
| grant fixed on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant flashback on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant float on *.* to u
< error
| 1064 42000 …line 1 column 11 near "float on *.* to u" 
> exec ^
: string
| grant float4 on *.* to u
< error
| 1064 42000 …line 1 column 12 near "float4 on *.* to u" 
> exec ^
: string
| grant float8 on *.* to u
< error
| 1064 42000 …line 1 column 12 near "float8 on *.* to u" 
> exec ^
: string
| grant flush on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant follower_constraints on *.* to u
< error
| 1064 42000 …line 1 column 29 near "on *.* to u" 
> exec ^
: string
| grant followers on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant following on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant follows on *.* to u
< ok 0 0
> exec ^
: string
| grant for on *.* to u
< error
| 1064 42000 …line 1 column 9 near "for on *.* to u" 
> exec ^
: string
| grant force on *.* to u
< error
| 1064 42000 …line 1 column 11 near "force on *.* to u" 
> exec ^
: string
| grant foreign on *.* to u
< error
| 1064 42000 …line 1 column 13 near "foreign on *.* to u" 
> exec ^
: string
| grant format on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant found on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant from on *.* to u
< error
| 1064 42000 …line 1 column 10 near "from on *.* to u" 
> exec ^
: string
| grant full on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant fulltext on *.* to u
< error
| 1064 42000 …line 1 column 14 near "fulltext on *.* to u" 
> exec ^
: string
| grant function on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant general on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant generate on *.* to u
< ok 0 0
> exec ^
: string
| grant generated on *.* to u
< error
| 1064 42000 …line 1 column 15 near "generated on *.* to u" 
> exec ^
: string
| grant geomcollection on *.* to u
< ok 0 0
> exec ^
: string
| grant geometry on *.* to u
< ok 0 0
> exec ^
: string
| grant geometrycollection on *.* to u
< ok 0 0
> exec ^
: string
| grant get on *.* to u
< ok 0 0
> exec ^
: string
| grant get_format on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant get_source_public_key on *.* to u
< ok 0 0
> exec ^
: string
| grant global on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant grant on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant grants on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant greatest on *.* to u
< ok 0 0
> exec ^
: string
| grant group on *.* to u
< error
| 1064 42000 …line 1 column 11 near "group on *.* to u" 
> exec ^
: string
| grant group_concat on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant group_replication on *.* to u
< ok 0 0
> exec ^
: string
| grant grouping on *.* to u
< ok 0 0
> exec ^
: string
| grant groups on *.* to u
< error
| 1064 42000 …line 1 column 12 near "groups on *.* to u" 
> exec ^
: string
| grant gtid_only on *.* to u
< ok 0 0
> exec ^
: string
| grant gtids on *.* to u
< ok 0 0
> exec ^
: string
| grant guided on *.* to u
< ok 0 0
> exec ^
: string
| grant handler on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant hash on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant having on *.* to u
< error
| 1064 42000 …line 1 column 12 near "having on *.* to u" 
> exec ^
: string
| grant header on *.* to u
< ok 0 0
> exec ^
: string
| grant help on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant high on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant high_priority on *.* to u
< error
| 1064 42000 …line 1 column 19 near "high_priority on *.* to u" 
> exec ^
: string
| grant histogram on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant histograms_in_flight on *.* to u
< error
| 1064 42000 …line 1 column 29 near "on *.* to u" 
> exec ^
: string
| grant history on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant host on *.* to u
< ok 0 0
> exec ^
: string
| grant hosts on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant hour on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant hour_microsecond on *.* to u
< error
| 1064 42000 …line 1 column 22 near "hour_microsecond on *.* to u" 
> exec ^
: string
| grant hour_minute on *.* to u
< error
| 1064 42000 …line 1 column 17 near "hour_minute on *.* to u" 
> exec ^
: string
| grant hour_second on *.* to u
< error
| 1064 42000 …line 1 column 17 near "hour_second on *.* to u" 
> exec ^
: string
| grant identified on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant if on *.* to u
< error
| 1064 42000 …line 1 column 8 near "if on *.* to u" 
> exec ^
: string
| grant ignore on *.* to u
< error
| 1064 42000 …line 1 column 12 near "ignore on *.* to u" 
> exec ^
: string
| grant ignore_server_ids on *.* to u
< ok 0 0
> exec ^
: string
| grant ilike on *.* to u
< error
| 1064 42000 …line 1 column 11 near "ilike on *.* to u" 
> exec ^
: string
| grant import on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant in on *.* to u
< error
| 1064 42000 …line 1 column 8 near "in on *.* to u" 
> exec ^
: string
| grant inactive on *.* to u
< ok 0 0
> exec ^
: string
| grant increment on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant index on *.* to u
< ok 0 0
> exec ^
: string
| grant indexes on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant infile on *.* to u
< error
| 1064 42000 …line 1 column 12 near "infile on *.* to u" 
> exec ^
: string
| grant initial on *.* to u
< ok 0 0
> exec ^
: string
| grant initial_size on *.* to u
< ok 0 0
> exec ^
: string
| grant initiate on *.* to u
< ok 0 0
> exec ^
: string
| grant inner on *.* to u
< error
| 1064 42000 …line 1 column 11 near "inner on *.* to u" 
> exec ^
: string
| grant inout on *.* to u
< error
| 1064 42000 …line 1 column 11 near "inout on *.* to u" 
> exec ^
: string
| grant insensitive on *.* to u
< ok 0 0
> exec ^
: string
| grant insert on *.* to u
< ok 0 0
> exec ^
: string
| grant insert_method on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant install on *.* to u
< ok 0 0
> exec ^
: string
| grant instance on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant int on *.* to u
< error
| 1064 42000 …line 1 column 9 near "int on *.* to u" 
> exec ^
: string
| grant int1 on *.* to u
< error
| 1064 42000 …line 1 column 10 near "int1 on *.* to u" 
> exec ^
: string
| grant int2 on *.* to u
< error
| 1064 42000 …line 1 column 10 near "int2 on *.* to u" 
> exec ^
: string
| grant int3 on *.* to u
< error
| 1064 42000 …line 1 column 10 near "int3 on *.* to u" 
> exec ^
: string
| grant int4 on *.* to u
< error
| 1064 42000 …line 1 column 10 near "int4 on *.* to u" 
> exec ^
: string
| grant int8 on *.* to u
< error
| 1064 42000 …line 1 column 10 near "int8 on *.* to u" 
> exec ^
: string
| grant integer on *.* to u
< error
| 1064 42000 …line 1 column 13 near "integer on *.* to u" 
> exec ^
: string
| grant intersect on *.* to u
< error
| 1064 42000 …line 1 column 15 near "intersect on *.* to u" 
> exec ^
: string
| grant interval on *.* to u
< error
| 1064 42000 …line 1 column 14 near "interval on *.* to u" 
> exec ^
: string
| grant into on *.* to u
< error
| 1064 42000 …line 1 column 10 near "into on *.* to u" 
> exec ^
: string
| grant invisible on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant invoker on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant io on *.* to u
< error
| 1064 42000 …line 1 column 11 near "on *.* to u" 
> exec ^
: string
| grant io_after_gtids on *.* to u
< ok 0 0
> exec ^
: string
| grant io_before_gtids on *.* to u
< ok 0 0
> exec ^
: string
| grant io_thread on *.* to u
< ok 0 0
> exec ^
: string
| grant ipc on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant is on *.* to u
< error
| 1064 42000 …line 1 column 8 near "is on *.* to u" 
> exec ^
: string
| grant isolation on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant issuer on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant iterate on *.* to u
< error
| 1064 42000 …line 1 column 13 near "iterate on *.* to u" 
> exec ^
: string
| grant job on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant jobs on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant join on *.* to u
< error
| 1064 42000 …line 1 column 10 near "join on *.* to u" 
> exec ^
: string
| grant json on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant json_arrayagg on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant json_duality_object on *.* to u
< ok 0 0
> exec ^
: string
| grant json_objectagg on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant json_table on *.* to u
< ok 0 0
> exec ^
: string
| grant json_value on *.* to u
< ok 0 0
> exec ^
: string
| grant key on *.* to u
< error
| 1064 42000 …line 1 column 9 near "key on *.* to u" 
> exec ^
: string
| grant key_block_size on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant keyring on *.* to u
< ok 0 0
> exec ^
: string
| grant keys on *.* to u
< error
| 1064 42000 …line 1 column 10 near "keys on *.* to u" 
> exec ^
: string
| grant kill on *.* to u
< error
| 1064 42000 …line 1 column 10 near "kill on *.* to u" 
> exec ^
: string
| grant labels on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant lag on *.* to u
< error
| 1064 42000 …line 1 column 9 near "lag on *.* to u" 
> exec ^
: string
| grant language on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant last on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant last_value on *.* to u
< error
| 1064 42000 …line 1 column 16 near "last_value on *.* to u" 
> exec ^
: string
| grant lateral on *.* to u
< ok 0 0
> exec ^
: string
| grant lead on *.* to u
< error
| 1064 42000 …line 1 column 10 near "lead on *.* to u" 
> exec ^
: string
| grant leader_constraints on *.* to u
< error
| 1064 42000 …line 1 column 27 near "on *.* to u" 
> exec ^
: string
| grant leading on *.* to u
< error
| 1064 42000 …line 1 column 13 near "leading on *.* to u" 
> exec ^
: string
| grant learner_constraints on *.* to u
< error
| 1064 42000 …line 1 column 28 near "on *.* to u" 
> exec ^
: string
| grant learners on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant least on *.* to u
< ok 0 0
> exec ^
: string
| grant leave on *.* to u
< error
| 1064 42000 …line 1 column 11 near "leave on *.* to u" 
> exec ^
: string
| grant leaves on *.* to u
< ok 0 0
> exec ^
: string
| grant left on *.* to u
< error
| 1064 42000 …line 1 column 10 near "left on *.* to u" 
> exec ^
: string
| grant less on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant level on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant library on *.* to u
< ok 0 0
> exec ^
: string
| grant like on *.* to u
< error
| 1064 42000 …line 1 column 10 near "like on *.* to u" 
> exec ^
: string
| grant limit on *.* to u
< error
| 1064 42000 …line 1 column 11 near "limit on *.* to u" 
> exec ^
: string
| grant linear on *.* to u
< error
| 1064 42000 …line 1 column 12 near "linear on *.* to u" 
> exec ^
: string
| grant lines on *.* to u
< error
| 1064 42000 …line 1 column 11 near "lines on *.* to u" 
> exec ^
: string
| grant linestring on *.* to u
< ok 0 0
> exec ^
: string
| grant list on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant load on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant local on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant localtime on *.* to u
< error
| 1064 42000 …line 1 column 15 near "localtime on *.* to u" 
> exec ^
: string
| grant localtimestamp on *.* to u
< error
| 1064 42000 …line 1 column 20 near "localtimestamp on *.* to u" 
> exec ^
: string
| grant location on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant lock on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant locked on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant locks on *.* to u
< ok 0 0
> exec ^
: string
| grant log on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant logfile on *.* to u
< ok 0 0
> exec ^
: string
| grant logs on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant long on *.* to u
< error
| 1064 42000 …line 1 column 10 near "long on *.* to u" 
> exec ^
: string
| grant longblob on *.* to u
< error
| 1064 42000 …line 1 column 14 near "longblob on *.* to u" 
> exec ^
: string
| grant longtext on *.* to u
< error
| 1064 42000 …line 1 column 14 near "longtext on *.* to u" 
> exec ^
: string
| grant loop on *.* to u
< ok 0 0
> exec ^
: string
| grant low on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant low_priority on *.* to u
< error
| 1064 42000 …line 1 column 18 near "low_priority on *.* to u" 
> exec ^
: string
| grant manual on *.* to u
< ok 0 0
> exec ^
: string
| grant master on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant master_log_file on *.* to u
< ok 0 0
> exec ^
: string
| grant match on *.* to u
< error
| 1064 42000 …line 1 column 11 near "match on *.* to u" 
> exec ^
: string
| grant materialized on *.* to u
< ok 0 0
> exec ^
: string
| grant max on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant max_connections_per_hour on *.* to u
< error
| 1064 42000 …line 1 column 33 near "on *.* to u" 
> exec ^
: string
| grant max_queries_per_hour on *.* to u
< error
| 1064 42000 …line 1 column 29 near "on *.* to u" 
> exec ^
: string
| grant max_rows on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant max_size on *.* to u
< ok 0 0
> exec ^
: string
| grant max_updates_per_hour on *.* to u
< error
| 1064 42000 …line 1 column 29 near "on *.* to u" 
> exec ^
: string
| grant max_user_connections on *.* to u
< error
| 1064 42000 …line 1 column 29 near "on *.* to u" 
> exec ^
: string
| grant maxvalue on *.* to u
< error
| 1064 42000 …line 1 column 14 near "maxvalue on *.* to u" 
> exec ^
: string
| grant medium on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant mediumblob on *.* to u
< error
| 1064 42000 …line 1 column 16 near "mediumblob on *.* to u" 
> exec ^
: string
| grant mediumint on *.* to u
< error
| 1064 42000 …line 1 column 15 near "mediumint on *.* to u" 
> exec ^
: string
| grant mediumtext on *.* to u
< error
| 1064 42000 …line 1 column 16 near "mediumtext on *.* to u" 
> exec ^
: string
| grant member on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant memory on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant merge on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant message_text on *.* to u
< ok 0 0
> exec ^
: string
| grant microsecond on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant mid on *.* to u
< ok 0 0
> exec ^
: string
| grant middleint on *.* to u
< error
| 1064 42000 …line 1 column 15 near "middleint on *.* to u" 
> exec ^
: string
| grant migrate on *.* to u
< ok 0 0
> exec ^
: string
| grant min on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant min_rows on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant minute on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant minute_microsecond on *.* to u
< error
| 1064 42000 …line 1 column 24 near "minute_microsecond on *.* to u" 
> exec ^
: string
| grant minute_second on *.* to u
< error
| 1064 42000 …line 1 column 19 near "minute_second on *.* to u" 
> exec ^
: string
| grant minvalue on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant mod on *.* to u
< error
| 1064 42000 …line 1 column 9 near "mod on *.* to u" 
> exec ^
: string
| grant mode on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant modifies on *.* to u
< ok 0 0
> exec ^
: string
| grant modify on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant month on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant multilinestring on *.* to u
< ok 0 0
> exec ^
: string
| grant multipoint on *.* to u
< ok 0 0
> exec ^
: string
| grant multipolygon on *.* to u
< ok 0 0
> exec ^
: string
| grant mutex on *.* to u
< ok 0 0
> exec ^
: string
| grant mysql_errno on *.* to u
< ok 0 0
> exec ^
: string
| grant name on *.* to u
< ok 0 0
> exec ^
: string
| grant names on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant national on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant natural on *.* to u
< error
| 1064 42000 …line 1 column 13 near "natural on *.* to u" 
> exec ^
: string
| grant nchar on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant ndb on *.* to u
< ok 0 0
> exec ^
: string
| grant ndbcluster on *.* to u
< ok 0 0
> exec ^
: string
| grant nested on *.* to u
< ok 0 0
> exec ^
: string
| grant network_namespace on *.* to u
< ok 0 0
> exec ^
: string
| grant never on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant new on *.* to u
< ok 0 0
> exec ^
: string
| grant next on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant next_row_id on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant no on *.* to u
< error
| 1064 42000 …line 1 column 11 near "on *.* to u" 
> exec ^
: string
| grant no_wait on *.* to u
< ok 0 0
> exec ^
: string
| grant no_write_to_binlog on *.* to u
< error
| 1064 42000 …line 1 column 24 near "no_write_to_binlog on *.* to u" 
> exec ^
: string
| grant nocache on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant nocycle on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant nodegroup on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant nomaxvalue on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant nominvalue on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant nonclustered on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant none on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant not on *.* to u
< error
| 1064 42000 …line 1 column 9 near "not on *.* to u" 
> exec ^
: string
| grant now on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant nowait on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant nth_value on *.* to u
< error
| 1064 42000 …line 1 column 15 near "nth_value on *.* to u" 
> exec ^
: string
| grant ntile on *.* to u
< error
| 1064 42000 …line 1 column 11 near "ntile on *.* to u" 
> exec ^
: string
| grant null on *.* to u
< error
| 1064 42000 …line 1 column 10 near "null on *.* to u" 
> exec ^
: string
| grant nullif on *.* to u
< ok 0 0
> exec ^
: string
| grant nulls on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant number on *.* to u
< ok 0 0
> exec ^
: string
| grant numeric on *.* to u
< error
| 1064 42000 …line 1 column 13 near "numeric on *.* to u" 
> exec ^
: string
| grant nvarchar on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant of on *.* to u
< error
| 1064 42000 …line 1 column 8 near "of on *.* to u" 
> exec ^
: string
| grant off on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant offset on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant oj on *.* to u
< ok 0 0
> exec ^
: string
| grant old on *.* to u
< ok 0 0
> exec ^
: string
| grant oltp_read_only on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant oltp_read_write on *.* to u
< error
| 1064 42000 …line 1 column 24 near "on *.* to u" 
> exec ^
: string
| grant oltp_write_only on *.* to u
< error
| 1064 42000 …line 1 column 24 near "on *.* to u" 
> exec ^
: string
| grant on on *.* to u
< error
| 1064 42000 …line 1 column 8 near "on on *.* to u" 
> exec ^
: string
| grant one on *.* to u
< ok 0 0
> exec ^
: string
| grant only on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant open on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant optimize on *.* to u
< error
| 1064 42000 …line 1 column 14 near "optimize on *.* to u" 
> exec ^
: string
| grant optimizer_costs on *.* to u
< ok 0 0
> exec ^
: string
| grant option on *.* to u
< error
| 1064 42000 …line 1 column 12 near "option on *.* to u" 
> exec ^
: string
| grant optional on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant optionally on *.* to u
< error
| 1064 42000 …line 1 column 16 near "optionally on *.* to u" 
> exec ^
: string
| grant options on *.* to u
< ok 0 0
> exec ^
: string
| grant or on *.* to u
< error
| 1064 42000 …line 1 column 8 near "or on *.* to u" 
> exec ^
: string
| grant order on *.* to u
< error
| 1064 42000 …line 1 column 11 near "order on *.* to u" 
> exec ^
: string
| grant ordinality on *.* to u
< ok 0 0
> exec ^
: string
| grant organization on *.* to u
< ok 0 0
> exec ^
: string
| grant others on *.* to u
< ok 0 0
> exec ^
: string
| grant out on *.* to u
< error
| 1064 42000 …line 1 column 9 near "out on *.* to u" 
> exec ^
: string
| grant outer on *.* to u
< error
| 1064 42000 …line 1 column 11 near "outer on *.* to u" 
> exec ^
: string
| grant outfile on *.* to u
< error
| 1064 42000 …line 1 column 13 near "outfile on *.* to u" 
> exec ^
: string
| grant over on *.* to u
< error
| 1064 42000 …line 1 column 10 near "over on *.* to u" 
> exec ^
: string
| grant overlay on *.* to u
< ok 0 0
> exec ^
: string
| grant owner on *.* to u
< ok 0 0
> exec ^
: string
| grant pack_keys on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant page on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant parallel on *.* to u
< ok 0 0
> exec ^
: string
| grant parameters on *.* to u
< ok 0 0
> exec ^
: string
| grant parse_tree on *.* to u
< ok 0 0
> exec ^
: string
| grant parser on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant partial on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant partition on *.* to u
< error
| 1064 42000 …line 1 column 15 near "partition on *.* to u" 
> exec ^
: string
| grant partitioning on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant partitions on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant password on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant password_lock_time on *.* to u
< error
| 1064 42000 …line 1 column 27 near "on *.* to u" 
> exec ^
: string
| grant path on *.* to u
< ok 0 0
> exec ^
: string
| grant percent_rank on *.* to u
< error
| 1064 42000 …line 1 column 18 near "percent_rank on *.* to u" 
> exec ^
: string
| grant persist on *.* to u
< ok 0 0
> exec ^
: string
| grant persist_only on *.* to u
< ok 0 0
> exec ^
: string
| grant phase on *.* to u
< ok 0 0
> exec ^
: string
| grant placement on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant placing on *.* to u
< ok 0 0
> exec ^
: string
| grant plan on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant plugin on *.* to u
< ok 0 0
> exec ^
: string
| grant plugin_dir on *.* to u
< ok 0 0
> exec ^
: string
| grant plugins on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant point on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant policy on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant polygon on *.* to u
< ok 0 0
> exec ^
: string
| grant port on *.* to u
< ok 0 0
> exec ^
: string
| grant position on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant pre_split_regions on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant precedes on *.* to u
< ok 0 0
> exec ^
: string
| grant preceding on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant precision on *.* to u
< error
| 1064 42000 …line 1 column 15 near "precision on *.* to u" 
> exec ^
: string
| grant prepare on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant preserve on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant prev on *.* to u
< ok 0 0
> exec ^
: string
| grant primary on *.* to u
< error
| 1064 42000 …line 1 column 13 near "primary on *.* to u" 
> exec ^
: string
| grant primary_region on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant priority on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant privilege_checks_user on *.* to u
< ok 0 0
> exec ^
: string
| grant privileges on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant procedure on *.* to u
< error
| 1064 42000 …line 1 column 15 near "procedure on *.* to u" 
> exec ^
: string
| grant process on *.* to u
< ok 0 0
> exec ^
: string
| grant processed_keys on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant processlist on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant profile on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant profiles on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant proxy on *.* to u
< error
| 1064 42000 …line 1 column 16 near "*.* to u" 
> exec ^
: string
| grant purge on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant qualify on *.* to u
< ok 0 0
> exec ^
: string
| grant quarter on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant query on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant query_limit on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant quick on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant random on *.* to u
< ok 0 0
> exec ^
: string
| grant range on *.* to u
< error
| 1064 42000 …line 1 column 11 near "range on *.* to u" 
> exec ^
: string
| grant rank on *.* to u
< error
| 1064 42000 …line 1 column 10 near "rank on *.* to u" 
> exec ^
: string
| grant read on *.* to u
< error
| 1064 42000 …line 1 column 10 near "read on *.* to u" 
> exec ^
: string
| grant read_only on *.* to u
< ok 0 0
> exec ^
: string
| grant read_write on *.* to u
< ok 0 0
> exec ^
: string
| grant reads on *.* to u
< ok 0 0
> exec ^
: string
| grant real on *.* to u
< error
| 1064 42000 …line 1 column 10 near "real on *.* to u" 
> exec ^
: string
| grant rebuild on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant recover on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant recursive on *.* to u
< error
| 1064 42000 …line 1 column 15 near "recursive on *.* to u" 
> exec ^
: string
| grant redo_buffer_size on *.* to u
< ok 0 0
> exec ^
: string
| grant redundant on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant reference on *.* to u
< ok 0 0
> exec ^
: string
| grant references on *.* to u
< ok 0 0
> exec ^
: string
| grant regexp on *.* to u
< error
| 1064 42000 …line 1 column 12 near "regexp on *.* to u" 
> exec ^
: string
| grant region on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant regions on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant registration on *.* to u
< ok 0 0
> exec ^
: string
| grant relational on *.* to u
< ok 0 0
> exec ^
: string
| grant relay on *.* to u
< ok 0 0
> exec ^
: string
| grant relay_log_file on *.* to u
< ok 0 0
> exec ^
: string
| grant relay_log_pos on *.* to u
< ok 0 0
> exec ^
: string
| grant relay_thread on *.* to u
< ok 0 0
> exec ^
: string
| grant relaylog on *.* to u
< ok 0 0
> exec ^
: string
| grant release on *.* to u
< error
| 1064 42000 …line 1 column 13 near "release on *.* to u" 
> exec ^
: string
| grant reload on *.* to u
< ok 0 0
> exec ^
: string
| grant remove on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant rename on *.* to u
< error
| 1064 42000 …line 1 column 12 near "rename on *.* to u" 
> exec ^
: string
| grant reorganize on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant repair on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant repeat on *.* to u
< error
| 1064 42000 …line 1 column 12 near "repeat on *.* to u" 
> exec ^
: string
| grant repeatable on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant replace on *.* to u
< error
| 1064 42000 …line 1 column 13 near "replace on *.* to u" 
> exec ^
: string
| grant replica on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant replicas on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant replicate_do_db on *.* to u
< ok 0 0
> exec ^
: string
| grant replicate_do_table on *.* to u
< ok 0 0
> exec ^
: string
| grant replicate_ignore_db on *.* to u
< ok 0 0
> exec ^
: string
| grant replicate_ignore_table on *.* to u
< ok 0 0
> exec ^
: string
| grant replicate_rewrite_db on *.* to u
< ok 0 0
> exec ^
: string
| grant replicate_wild_do_table on *.* to u
< ok 0 0
> exec ^
: string
| grant replicate_wild_ignore_table on *.* to u
< ok 0 0
> exec ^
: string
| grant replication on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant require on *.* to u
< error
| 1064 42000 …line 1 column 13 near "require on *.* to u" 
> exec ^
: string
| grant require_row_format on *.* to u
< ok 0 0
> exec ^
: string
| grant require_table_primary_key_check on *.* to u
< ok 0 0
> exec ^
: string
| grant reset on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant resignal on *.* to u
< ok 0 0
> exec ^
: string
| grant resource on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant respect on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant restart on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant restore on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant restores on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant restrict on *.* to u
< error
| 1064 42000 …line 1 column 14 near "restrict on *.* to u" 
> exec ^
: string
| grant resume on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant retain on *.* to u
< ok 0 0
> exec ^
: string
| grant return on *.* to u
< ok 0 0
> exec ^
: string
| grant returned_sqlstate on *.* to u
< ok 0 0
> exec ^
: string
| grant returning on *.* to u
< ok 0 0
> exec ^
: string
| grant returns on *.* to u
< ok 0 0
> exec ^
: string
| grant reuse on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant reverse on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant revoke on *.* to u
< error
| 1064 42000 …line 1 column 12 near "revoke on *.* to u" 
> exec ^
: string
| grant right on *.* to u
< error
| 1064 42000 …line 1 column 11 near "right on *.* to u" 
> exec ^
: string
| grant rlike on *.* to u
< error
| 1064 42000 …line 1 column 11 near "rlike on *.* to u" 
> exec ^
: string
| grant role on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant rollback on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant rollup on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant rotate on *.* to u
< ok 0 0
> exec ^
: string
| grant routine on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant row on *.* to u
< error
| 1064 42000 …line 1 column 9 near "row on *.* to u" 
> exec ^
: string
| grant row_count on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant row_format on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant row_number on *.* to u
< error
| 1064 42000 …line 1 column 16 near "row_number on *.* to u" 
> exec ^
: string
| grant rows on *.* to u
< error
| 1064 42000 …line 1 column 10 near "rows on *.* to u" 
> exec ^
: string
| grant rtree on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant ru on *.* to u
< error
| 1064 42000 …line 1 column 11 near "on *.* to u" 
> exec ^
: string
| grant ru_per_sec on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant rule on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant run on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant s3 on *.* to u
< error
| 1064 42000 …line 1 column 11 near "on *.* to u" 
> exec ^
: string
| grant savepoint on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant schedule on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant schema on *.* to u
< error
| 1064 42000 …line 1 column 12 near "schema on *.* to u" 
> exec ^
: string
| grant schema_name on *.* to u
< ok 0 0
> exec ^
: string
| grant schemas on *.* to u
< error
| 1064 42000 …line 1 column 13 near "schemas on *.* to u" 
> exec ^
: string
| grant second on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant second_microsecond on *.* to u
< error
| 1064 42000 …line 1 column 24 near "second_microsecond on *.* to u" 
> exec ^
: string
| grant secondary on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant secondary_engine on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant secondary_engine_attribute on *.* to u
< ok 0 0
> exec ^
: string
| grant secondary_load on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant secondary_unload on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant security on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant select on *.* to u
< ok 0 0
> exec ^
: string
| grant sensitive on *.* to u
< ok 0 0
> exec ^
: string
| grant separator on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant sequence on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant serial on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant serializable on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant server on *.* to u
< ok 0 0
> exec ^
: string
| grant session on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant session_states on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant session_user on *.* to u
< ok 0 0
> exec ^
: string
| grant set on *.* to u
< error
| 1064 42000 …line 1 column 9 near "set on *.* to u" 
> exec ^
: string
| grant sets on *.* to u
< ok 0 0
> exec ^
: string
| grant shard_row_id_bits on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant share on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant show on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant shutdown on *.* to u
< ok 0 0
> exec ^
: string
| grant signal on *.* to u
< ok 0 0
> exec ^
: string
| grant signed on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant similar on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant simple on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant skip on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant slave on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant slow on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant smallint on *.* to u
< error
| 1064 42000 …line 1 column 14 near "smallint on *.* to u" 
> exec ^
: string
| grant snapshot on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant socket on *.* to u
< ok 0 0
> exec ^
: string
| grant some on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant soname on *.* to u
< ok 0 0
> exec ^
: string
| grant sounds on *.* to u
< ok 0 0
> exec ^
: string
| grant source on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant source_auto_position on *.* to u
< ok 0 0
> exec ^
: string
| grant source_bind on *.* to u
< ok 0 0
> exec ^
: string
| grant source_compression_algorithms on *.* to u
< ok 0 0
> exec ^
: string
| grant source_connect_retry on *.* to u
< ok 0 0
> exec ^
: string
| grant source_connection_auto_failover on *.* to u
< ok 0 0
> exec ^
: string
| grant source_delay on *.* to u
< ok 0 0
> exec ^
: string
| grant source_heartbeat_period on *.* to u
< ok 0 0
> exec ^
: string
| grant source_host on *.* to u
< ok 0 0
> exec ^
: string
| grant source_log_file on *.* to u
< ok 0 0
> exec ^
: string
| grant source_log_pos on *.* to u
< ok 0 0
> exec ^
: string
| grant source_password on *.* to u
< ok 0 0
> exec ^
: string
| grant source_port on *.* to u
< ok 0 0
> exec ^
: string
| grant source_public_key_path on *.* to u
< ok 0 0
> exec ^
: string
| grant source_retry_count on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_ca on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_capath on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_cert on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_cipher on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_crl on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_crlpath on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_key on *.* to u
< ok 0 0
> exec ^
: string
| grant source_ssl_verify_server_cert on *.* to u
< ok 0 0
> exec ^
: string
| grant source_tls_ciphersuites on *.* to u
< ok 0 0
> exec ^
: string
| grant source_tls_version on *.* to u
< ok 0 0
> exec ^
: string
| grant source_user on *.* to u
< ok 0 0
> exec ^
: string
| grant source_zstd_compression_level on *.* to u
< ok 0 0
> exec ^
: string
| grant spatial on *.* to u
< error
| 1064 42000 …line 1 column 13 near "spatial on *.* to u" 
> exec ^
: string
| grant specific on *.* to u
< ok 0 0
> exec ^
: string
| grant split on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant sql on *.* to u
< error
| 1064 42000 …line 1 column 9 near "sql on *.* to u" 
> exec ^
: string
| grant sql_after_gtids on *.* to u
< ok 0 0
> exec ^
: string
| grant sql_after_mts_gaps on *.* to u
< ok 0 0
> exec ^
: string
| grant sql_before_gtids on *.* to u
< ok 0 0
> exec ^
: string
| grant sql_big_result on *.* to u
< error
| 1064 42000 …line 1 column 20 near "sql_big_result on *.* to u" 
> exec ^
: string
| grant sql_buffer_result on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant sql_calc_found_rows on *.* to u
< error
| 1064 42000 …line 1 column 25 near "sql_calc_found_rows on *.* to u" 
> exec ^
: string
| grant sql_no_cache on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant sql_small_result on *.* to u
< error
| 1064 42000 …line 1 column 22 near "sql_small_result on *.* to u" 
> exec ^
: string
| grant sql_thread on *.* to u
< ok 0 0
> exec ^
: string
| grant sql_tsi_day on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_hour on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_minute on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_month on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_quarter on *.* to u
< error
| 1064 42000 …line 1 column 24 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_second on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_week on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant sql_tsi_year on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant sqlexception on *.* to u
< error
| 1064 42000 …line 1 column 18 near "sqlexception on *.* to u" 
> exec ^
: string
| grant sqlstate on *.* to u
< error
| 1064 42000 …line 1 column 14 near "sqlstate on *.* to u" 
> exec ^
: string
| grant sqlwarning on *.* to u
< error
| 1064 42000 …line 1 column 16 near "sqlwarning on *.* to u" 
> exec ^
: string
| grant srid on *.* to u
< ok 0 0
> exec ^
: string
| grant ssl on *.* to u
< error
| 1064 42000 …line 1 column 9 near "ssl on *.* to u" 
> exec ^
: string
| grant st_collect on *.* to u
< ok 0 0
> exec ^
: string
| grant stacked on *.* to u
< ok 0 0
> exec ^
: string
| grant start on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant start_time on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant starting on *.* to u
< error
| 1064 42000 …line 1 column 14 near "starting on *.* to u" 
> exec ^
: string
| grant starts on *.* to u
< ok 0 0
> exec ^
: string
| grant stats_auto_recalc on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant stats_buckets on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant stats_col_choice on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant stats_col_list on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant stats_extended on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant stats_healthy on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant stats_histograms on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant stats_locked on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant stats_meta on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant stats_persistent on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant stats_sample_pages on *.* to u
< error
| 1064 42000 …line 1 column 27 near "on *.* to u" 
> exec ^
: string
| grant stats_sample_rate on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant stats_topn on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant status on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant std on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant stddev on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant stddev_pop on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant stddev_samp on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant stop on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant storage on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant stored on *.* to u
< error
| 1064 42000 …line 1 column 12 near "stored on *.* to u" 
> exec ^
: string
| grant straight_join on *.* to u
< error
| 1064 42000 …line 1 column 19 near "straight_join on *.* to u" 
> exec ^
: string
| grant stream on *.* to u
< ok 0 0
> exec ^
: string
| grant strict_load on *.* to u
< ok 0 0
> exec ^
: string
| grant string on *.* to u
< ok 0 0
> exec ^
: string
| grant subclass_origin on *.* to u
< ok 0 0
> exec ^
: string
| grant subdate on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant subject on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant subpartition on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant subpartitions on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant substr on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant substring on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant sum on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant super on *.* to u
< ok 0 0
> exec ^
: string
| grant survival_preferences on *.* to u
< error
| 1064 42000 …line 1 column 29 near "on *.* to u" 
> exec ^
: string
| grant suspend on *.* to u
< ok 0 0
> exec ^
: string
| grant swaps on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant switch_group on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant switches on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant sysdate on *.* to u
< ok 0 0
> exec ^
: string
| grant system on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant system_user on *.* to u
< ok 0 0
> exec ^
: string
| grant table on *.* to u
< error
| 1064 42000 …line 1 column 11 near "table on *.* to u" 
> exec ^
: string
| grant table_checksum on *.* to u
< error
| 1064 42000 …line 1 column 23 near "on *.* to u" 
> exec ^
: string
| grant table_name on *.* to u
< ok 0 0
> exec ^
: string
| grant tables on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant tablesample on *.* to u
< error
| 1064 42000 …line 1 column 17 near "tablesample on *.* to u" 
> exec ^
: string
| grant tablespace on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant task_types on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant temporary on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant temptable on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant terminated on *.* to u
< error
| 1064 42000 …line 1 column 16 near "terminated on *.* to u" 
> exec ^
: string
| grant text on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant than on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant then on *.* to u
< error
| 1064 42000 …line 1 column 10 near "then on *.* to u" 
> exec ^
: string
| grant thread_priority on *.* to u
< ok 0 0
> exec ^
: string
| grant ties on *.* to u
< ok 0 0
> exec ^
: string
| grant tiflash on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant time on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant timeout on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant timestamp on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant timestampadd on *.* to u
< error
| 1064 42000 …line 1 column 21 near "on *.* to u" 
> exec ^
: string
| grant timestampdiff on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant tinyblob on *.* to u
< error
| 1064 42000 …line 1 column 14 near "tinyblob on *.* to u" 
> exec ^
: string
| grant tinyint on *.* to u
< error
| 1064 42000 …line 1 column 13 near "tinyint on *.* to u" 
> exec ^
: string
| grant tinytext on *.* to u
< error
| 1064 42000 …line 1 column 14 near "tinytext on *.* to u" 
> exec ^
: string
| grant tls on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant to on *.* to u
< error
| 1064 42000 …line 1 column 8 near "to on *.* to u" 
> exec ^
: string
| grant tpcc on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant tpch_10 on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant trailing on *.* to u
< error
| 1064 42000 …line 1 column 14 near "trailing on *.* to u" 
> exec ^
: string
| grant transaction on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant trigger on *.* to u
< ok 0 0
> exec ^
: string
| grant triggers on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant trim on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant true on *.* to u
< error
| 1064 42000 …line 1 column 10 near "true on *.* to u" 
> exec ^
: string
| grant truncate on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant ttl on *.* to u
< error
| 1064 42000 …line 1 column 12 near "on *.* to u" 
> exec ^
: string
| grant ttl_enable on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant ttl_job_interval on *.* to u
< error
| 1064 42000 …line 1 column 25 near "on *.* to u" 
> exec ^
: string
| grant type on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant types on *.* to u
< ok 0 0
> exec ^
: string
| grant unbounded on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant uncommitted on *.* to u
< error
| 1064 42000 …line 1 column 20 near "on *.* to u" 
> exec ^
: string
| grant undefined on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant undo on *.* to u
< ok 0 0
> exec ^
: string
| grant undo_buffer_size on *.* to u
< ok 0 0
> exec ^
: string
| grant undofile on *.* to u
< ok 0 0
> exec ^
: string
| grant unicode on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant uninstall on *.* to u
< ok 0 0
> exec ^
: string
| grant union on *.* to u
< error
| 1064 42000 …line 1 column 11 near "union on *.* to u" 
> exec ^
: string
| grant unique on *.* to u
< error
| 1064 42000 …line 1 column 12 near "unique on *.* to u" 
> exec ^
: string
| grant unknown on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant unlimited on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant unlock on *.* to u
< error
| 1064 42000 …line 1 column 12 near "unlock on *.* to u" 
> exec ^
: string
| grant unregister on *.* to u
< ok 0 0
> exec ^
: string
| grant unsigned on *.* to u
< error
| 1064 42000 …line 1 column 14 near "unsigned on *.* to u" 
> exec ^
: string
| grant until on *.* to u
< error
| 1064 42000 …line 1 column 11 near "until on *.* to u" 
> exec ^
: string
| grant update on *.* to u
< ok 0 0
> exec ^
: string
| grant upgrade on *.* to u
< ok 0 0
> exec ^
: string
| grant uri on *.* to u
< ok 0 0
> exec ^
: string
| grant url on *.* to u
< ok 0 0
> exec ^
: string
| grant usage on *.* to u
< ok 0 0
> exec ^
: string
| grant use on *.* to u
< error
| 1064 42000 …line 1 column 9 near "use on *.* to u" 
> exec ^
: string
| grant use_frm on *.* to u
< ok 0 0
> exec ^
: string
| grant user on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant user_resources on *.* to u
< ok 0 0
> exec ^
: string
| grant using on *.* to u
< error
| 1064 42000 …line 1 column 11 near "using on *.* to u" 
> exec ^
: string
| grant utc_date on *.* to u
< error
| 1064 42000 …line 1 column 14 near "utc_date on *.* to u" 
> exec ^
: string
| grant utc_time on *.* to u
< error
| 1064 42000 …line 1 column 14 near "utc_time on *.* to u" 
> exec ^
: string
| grant utc_timestamp on *.* to u
< error
| 1064 42000 …line 1 column 19 near "utc_timestamp on *.* to u" 
> exec ^
: string
| grant utilization_limit on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant validate on *.* to u
< ok 0 0
> exec ^
: string
| grant validation on *.* to u
< error
| 1064 42000 …line 1 column 19 near "on *.* to u" 
> exec ^
: string
| grant value on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant values on *.* to u
< error
| 1064 42000 …line 1 column 12 near "values on *.* to u" 
> exec ^
: string
| grant var_pop on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant var_samp on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant varbinary on *.* to u
< error
| 1064 42000 …line 1 column 15 near "varbinary on *.* to u" 
> exec ^
: string
| grant varchar on *.* to u
< error
| 1064 42000 …line 1 column 13 near "varchar on *.* to u" 
> exec ^
: string
| grant varcharacter on *.* to u
< error
| 1064 42000 …line 1 column 18 near "varcharacter on *.* to u" 
> exec ^
: string
| grant variables on *.* to u
< error
| 1064 42000 …line 1 column 18 near "on *.* to u" 
> exec ^
: string
| grant variance on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant varying on *.* to u
< error
| 1064 42000 …line 1 column 13 near "varying on *.* to u" 
> exec ^
: string
| grant vcpu on *.* to u
< ok 0 0
> exec ^
: string
| grant vector on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant verify_key_constraints on *.* to u
< ok 0 0
> exec ^
: string
| grant view on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant virtual on *.* to u
< error
| 1064 42000 …line 1 column 13 near "virtual on *.* to u" 
> exec ^
: string
| grant visible on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant voter_constraints on *.* to u
< error
| 1064 42000 …line 1 column 26 near "on *.* to u" 
> exec ^
: string
| grant voters on *.* to u
< error
| 1064 42000 …line 1 column 15 near "on *.* to u" 
> exec ^
: string
| grant wait on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant warnings on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant watch on *.* to u
< error
| 1064 42000 …line 1 column 14 near "on *.* to u" 
> exec ^
: string
| grant week on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant weight_string on *.* to u
< error
| 1064 42000 …line 1 column 22 near "on *.* to u" 
> exec ^
: string
| grant when on *.* to u
< error
| 1064 42000 …line 1 column 10 near "when on *.* to u" 
> exec ^
: string
| grant where on *.* to u
< error
| 1064 42000 …line 1 column 11 near "where on *.* to u" 
> exec ^
: string
| grant while on *.* to u
< error
| 1064 42000 …line 1 column 11 near "while on *.* to u" 
> exec ^
: string
| grant window on *.* to u
< error
| 1064 42000 …line 1 column 12 near "window on *.* to u" 
> exec ^
: string
| grant with on *.* to u
< error
| 1064 42000 …line 1 column 10 near "with on *.* to u" 
> exec ^
: string
| grant without on *.* to u
< error
| 1064 42000 …line 1 column 16 near "on *.* to u" 
> exec ^
: string
| grant work on *.* to u
< ok 0 0
> exec ^
: string
| grant workload on *.* to u
< error
| 1064 42000 …line 1 column 17 near "on *.* to u" 
> exec ^
: string
| grant wrapper on *.* to u
< ok 0 0
> exec ^
: string
| grant write on *.* to u
< error
| 1064 42000 …line 1 column 11 near "write on *.* to u" 
> exec ^
: string
| grant x509 on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant xa on *.* to u
< ok 0 0
> exec ^
: string
| grant xid on *.* to u
< ok 0 0
> exec ^
: string
| grant xml on *.* to u
< ok 0 0
> exec ^
: string
| grant xor on *.* to u
< error
| 1064 42000 …line 1 column 9 near "xor on *.* to u" 
> exec ^
: string
| grant year on *.* to u
< error
| 1064 42000 …line 1 column 13 near "on *.* to u" 
> exec ^
: string
| grant year_month on *.* to u
< error
| 1064 42000 …line 1 column 16 near "year_month on *.* to u" 
> exec ^
: string
| grant zerofill on *.* to u
< error
| 1064 42000 …line 1 column 14 near "zerofill on *.* to u" 
> exec ^
: string
| grant zone on *.* to u
< ok 0 0
> exec ^
: string
| show grants for u using absent
< ok 0 0
> exec ^
: string
| show grants for u using accessible
< ok 0 0
> exec ^
: string
| show grants for u using account
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using action
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using active
< ok 0 0
> exec ^
: string
| show grants for u using add
< error
| 1064 42000 …line 1 column 27 near "add" 
> exec ^
: string
| show grants for u using adddate
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using admin
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using after
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using against
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using aggregate
< ok 0 0
> exec ^
: string
| show grants for u using algorithm
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using all
< error
| 1064 42000 …line 1 column 27 near "all" 
> exec ^
: string
| show grants for u using allow_missing_files
< ok 0 0
> exec ^
: string
| show grants for u using alter
< error
| 1064 42000 …line 1 column 29 near "alter" 
> exec ^
: string
| show grants for u using always
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using analyze
< error
| 1064 42000 …line 1 column 31 near "analyze" 
> exec ^
: string
| show grants for u using and
< error
| 1064 42000 …line 1 column 27 near "and" 
> exec ^
: string
| show grants for u using any
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using array
< error
| 1064 42000 …line 1 column 29 near "array" 
> exec ^
: string
| show grants for u using as
< error
| 1064 42000 …line 1 column 26 near "as" 
> exec ^
: string
| show grants for u using asc
< error
| 1064 42000 …line 1 column 27 near "asc" 
> exec ^
: string
| show grants for u using ascii
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using asensitive
< ok 0 0
> exec ^
: string
| show grants for u using assign_gtids_to_anonymous_transactions
< ok 0 0
> exec ^
: string
| show grants for u using at
< ok 0 0
> exec ^
: string
| show grants for u using attribute
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using authentication
< ok 0 0
> exec ^
: string
| show grants for u using auto
< ok 0 0
> exec ^
: string
| show grants for u using auto_id_cache
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using auto_increment
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using auto_random
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using auto_random_base
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using auto_refresh
< ok 0 0
> exec ^
: string
| show grants for u using auto_refresh_source
< ok 0 0
> exec ^
: string
| show grants for u using autoextend_size
< ok 0 0
> exec ^
: string
| show grants for u using avg
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using avg_row_length
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using background
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using backup
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using backups
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using batch
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using before
< ok 0 0
> exec ^
: string
| show grants for u using begin
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using bernoulli
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using between
< error
| 1064 42000 …line 1 column 31 near "between" 
> exec ^
: string
| show grants for u using bigint
< error
| 1064 42000 …line 1 column 30 near "bigint" 
> exec ^
: string
| show grants for u using binary
< error
| 1064 42000 …line 1 column 30 near "binary" 
> exec ^
: string
| show grants for u using binding
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using binding_cache
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using bindings
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using binlog
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using bit
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using bit_and
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using bit_or
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using bit_xor
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using blob
< error
| 1064 42000 …line 1 column 28 near "blob" 
> exec ^
: string
| show grants for u using block
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using bool
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using boolean
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using both
< error
| 1064 42000 …line 1 column 28 near "both" 
> exec ^
: string
| show grants for u using btree
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using buckets
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using builtins
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using bulk
< ok 0 0
> exec ^
: string
| show grants for u using burstable
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using by
< error
| 1064 42000 …line 1 column 26 near "by" 
> exec ^
: string
| show grants for u using byte
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using cache
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using calibrate
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using call
< error
| 1064 42000 …line 1 column 28 near "call" 
> exec ^
: string
| show grants for u using cascade
< error
| 1064 42000 …line 1 column 31 near "cascade" 
> exec ^
: string
| show grants for u using cascaded
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using case
< error
| 1064 42000 …line 1 column 28 near "case" 
> exec ^
: string
| show grants for u using cast
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using catalog_name
< ok 0 0
> exec ^
: string
| show grants for u using chain
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using challenge_response
< ok 0 0
> exec ^
: string
| show grants for u using change
< error
| 1064 42000 …line 1 column 30 near "change" 
> exec ^
: string
| show grants for u using changed
< ok 0 0
> exec ^
: string
| show grants for u using channel
< ok 0 0
> exec ^
: string
| show grants for u using char
< error
| 1064 42000 …line 1 column 28 near "char" 
> exec ^
: string
| show grants for u using character
< error
| 1064 42000 …line 1 column 33 near "character" 
> exec ^
: string
| show grants for u using charset
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using check
< error
| 1064 42000 …line 1 column 29 near "check" 
> exec ^
: string
| show grants for u using checksum
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using cipher
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using class_origin
< ok 0 0
> exec ^
: string
| show grants for u using client
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using clone
< ok 0 0
> exec ^
: string
| show grants for u using close
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using clustered
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using coalesce
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using code
< ok 0 0
> exec ^
: string
| show grants for u using collate
< error
| 1064 42000 …line 1 column 31 near "collate" 
> exec ^
: string
| show grants for u using collation
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using column
< error
| 1064 42000 …line 1 column 30 near "column" 
> exec ^
: string
| show grants for u using column_format
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using column_name
< ok 0 0
> exec ^
: string
| show grants for u using column_stats_usage
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| show grants for u using columns
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using comment
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using commit
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using committed
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using compact
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using completion
< ok 0 0
> exec ^
: string
| show grants for u using component
< ok 0 0
> exec ^
: string
| show grants for u using compressed
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using compression
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using concat
< ok 0 0
> exec ^
: string
| show grants for u using concurrent
< ok 0 0
> exec ^
: string
| show grants for u using condition
< ok 0 0
> exec ^
: string
| show grants for u using config
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using connection
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using consistent
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using constraint
< error
| 1064 42000 …line 1 column 34 near "constraint" 
> exec ^
: string
| show grants for u using constraint_catalog
< ok 0 0
> exec ^
: string
| show grants for u using constraint_name
< ok 0 0
> exec ^
: string
| show grants for u using constraint_schema
< ok 0 0
> exec ^
: string
| show grants for u using constraints
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using contains
< ok 0 0
> exec ^
: string
| show grants for u using context
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using continue
< error
| 1064 42000 …line 1 column 32 near "continue" 
> exec ^
: string
| show grants for u using convert
< error
| 1064 42000 …line 1 column 31 near "convert" 
> exec ^
: string
| show grants for u using cooldown
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using count
< ok 0 0
> exec ^
: string
| show grants for u using cpu
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using create
< error
| 1064 42000 …line 1 column 30 near "create" 
> exec ^
: string
| show grants for u using cross
< error
| 1064 42000 …line 1 column 29 near "cross" 
> exec ^
: string
| show grants for u using cube
< ok 0 0
> exec ^
: string
| show grants for u using cume_dist
< error
| 1064 42000 …line 1 column 33 near "cume_dist" 
> exec ^
: string
| show grants for u using curdate
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using current
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using current_date
< error
| 1064 42000 …line 1 column 36 near "current_date" 
> exec ^
: string
| show grants for u using current_time
< error
| 1064 42000 …line 1 column 36 near "current_time" 
> exec ^
: string
| show grants for u using current_timestamp
< error
| 1064 42000 …line 1 column 41 near "current_timestamp" 
> exec ^
: string
| show grants for u using current_user
< error
| 1064 42000 …line 1 column 36 near "current_user" 
> exec ^
: string
| show grants for u using cursor
< error
| 1064 42000 …line 1 column 30 near "cursor" 
> exec ^
: string
| show grants for u using cursor_name
< ok 0 0
> exec ^
: string
| show grants for u using curtime
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using cycle
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using data
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using database
< error
| 1064 42000 …line 1 column 32 near "database" 
> exec ^
: string
| show grants for u using databases
< error
| 1064 42000 …line 1 column 33 near "databases" 
> exec ^
: string
| show grants for u using datafile
< ok 0 0
> exec ^
: string
| show grants for u using date
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using date_add
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using date_sub
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using datetime
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using day
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using day_hour
< error
| 1064 42000 …line 1 column 32 near "day_hour" 
> exec ^
: string
| show grants for u using day_microsecond
< error
| 1064 42000 …line 1 column 39 near "day_microsecond" 
> exec ^
: string
| show grants for u using day_minute
< error
| 1064 42000 …line 1 column 34 near "day_minute" 
> exec ^
: string
| show grants for u using day_second
< error
| 1064 42000 …line 1 column 34 near "day_second" 
> exec ^
: string
| show grants for u using deallocate
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using dec
< error
| 1064 42000 …line 1 column 27 near "dec" 
> exec ^
: string
| show grants for u using decimal
< error
| 1064 42000 …line 1 column 31 near "decimal" 
> exec ^
: string
| show grants for u using declare
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using default
< error
| 1064 42000 …line 1 column 31 near "default" 
> exec ^
: string
| show grants for u using default_auth
< ok 0 0
> exec ^
: string
| show grants for u using defined
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using definer
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using definition
< ok 0 0
> exec ^
: string
| show grants for u using delay_key_write
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| show grants for u using delayed
< error
| 1064 42000 …line 1 column 31 near "delayed" 
> exec ^
: string
| show grants for u using delete
< error
| 1064 42000 …line 1 column 30 near "delete" 
> exec ^
: string
| show grants for u using dense_rank
< error
| 1064 42000 …line 1 column 34 near "dense_rank" 
> exec ^
: string
| show grants for u using desc
< error
| 1064 42000 …line 1 column 28 near "desc" 
> exec ^
: string
| show grants for u using describe
< error
| 1064 42000 …line 1 column 32 near "describe" 
> exec ^
: string
| show grants for u using description
< ok 0 0
> exec ^
: string
| show grants for u using deterministic
< ok 0 0
> exec ^
: string
| show grants for u using diagnostics
< ok 0 0
> exec ^
: string
| show grants for u using digest
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using directory
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using disable
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using disabled
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using discard
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using disk
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using distinct
< error
| 1064 42000 …line 1 column 32 near "distinct" 
> exec ^
: string
| show grants for u using distinctrow
< error
| 1064 42000 …line 1 column 35 near "distinctrow" 
> exec ^
: string
| show grants for u using distribute
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using distribution
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using distributions
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using div
< error
| 1064 42000 …line 1 column 27 near "div" 
> exec ^
: string
| show grants for u using do
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| show grants for u using double
< error
| 1064 42000 …line 1 column 30 near "double" 
> exec ^
: string
| show grants for u using drop
< error
| 1064 42000 …line 1 column 28 near "drop" 
> exec ^
: string
| show grants for u using dry
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using dryrun
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using dual
< error
| 1064 42000 …line 1 column 28 near "dual" 
> exec ^
: string
| show grants for u using duality
< ok 0 0
> exec ^
: string
| show grants for u using dumpfile
< ok 0 0
> exec ^
: string
| show grants for u using duplicate
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using duration
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using dynamic
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using each
< ok 0 0
> exec ^
: string
| show grants for u using else
< error
| 1064 42000 …line 1 column 28 near "else" 
> exec ^
: string
| show grants for u using elseif
< error
| 1064 42000 …line 1 column 30 near "elseif" 
> exec ^
: string
| show grants for u using empty
< ok 0 0
> exec ^
: string
| show grants for u using enable
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using enabled
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using enclosed
< error
| 1064 42000 …line 1 column 32 near "enclosed" 
> exec ^
: string
| show grants for u using encryption
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using end
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using end_time
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using ends
< ok 0 0
> exec ^
: string
| show grants for u using enforced
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using engine
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using engine_attribute
< ok 0 0
> exec ^
: string
| show grants for u using engines
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using enum
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using error
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using errors
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using escape
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using escaped
< error
| 1064 42000 …line 1 column 31 near "escaped" 
> exec ^
: string
| show grants for u using event
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using events
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using every
< ok 0 0
> exec ^
: string
| show grants for u using exact
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using except
< error
| 1064 42000 …line 1 column 30 near "except" 
> exec ^
: string
| show grants for u using exchange
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using exclude
< ok 0 0
> exec ^
: string
| show grants for u using exec_elapsed
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using execute
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using exists
< error
| 1064 42000 …line 1 column 30 near "exists" 
> exec ^
: string
| show grants for u using exit
< error
| 1064 42000 …line 1 column 28 near "exit" 
> exec ^
: string
| show grants for u using expansion
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using expire
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using explain
< error
| 1064 42000 …line 1 column 31 near "explain" 
> exec ^
: string
| show grants for u using export
< ok 0 0
> exec ^
: string
| show grants for u using extended
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using extent_size
< ok 0 0
> exec ^
: string
| show grants for u using external
< ok 0 0
> exec ^
: string
| show grants for u using external_format
< ok 0 0
> exec ^
: string
| show grants for u using extract
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using factor
< ok 0 0
> exec ^
: string
| show grants for u using failed_login_attempts
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| show grants for u using false
< error
| 1064 42000 …line 1 column 29 near "false" 
> exec ^
: string
| show grants for u using fast
< ok 0 0
> exec ^
: string
| show grants for u using faults
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using fetch
< error
| 1064 42000 …line 1 column 29 near "fetch" 
> exec ^
: string
| show grants for u using fields
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using file
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using file_block_size
< ok 0 0
> exec ^
: string
| show grants for u using file_format
< ok 0 0
> exec ^
: string
| show grants for u using file_name
< ok 0 0
> exec ^
: string
| show grants for u using file_pattern
< ok 0 0
> exec ^
: string
| show grants for u using file_prefix
< ok 0 0
> exec ^
: string
| show grants for u using files
< ok 0 0
> exec ^
: string
| show grants for u using filter
< ok 0 0
> exec ^
: string
| show grants for u using finish
< ok 0 0
> exec ^
: string
| show grants for u using first
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using first_value
< error
| 1064 42000 …line 1 column 35 near "first_value" 
> exec ^
: string
| show grants for u using fixed
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using flashback
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using float
< error
| 1064 42000 …line 1 column 29 near "float" 
> exec ^
: string
| show grants for u using float4
< error
| 1064 42000 …line 1 column 30 near "float4" 
> exec ^
: string
| show grants for u using float8
< error
| 1064 42000 …line 1 column 30 near "float8" 
> exec ^
: string
| show grants for u using flush
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using follower_constraints
< error
| 1064 42000 …line 1 column 44 near "" 
> exec ^
: string
| show grants for u using followers
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using following
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using follows
< ok 0 0
> exec ^
: string
| show grants for u using for
< error
| 1064 42000 …line 1 column 27 near "for" 
> exec ^
: string
| show grants for u using force
< error
| 1064 42000 …line 1 column 29 near "force" 
> exec ^
: string
| show grants for u using foreign
< error
| 1064 42000 …line 1 column 31 near "foreign" 
> exec ^
: string
| show grants for u using format
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using found
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using from
< error
| 1064 42000 …line 1 column 28 near "from" 
> exec ^
: string
| show grants for u using full
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using fulltext
< error
| 1064 42000 …line 1 column 32 near "fulltext" 
> exec ^
: string
| show grants for u using function
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using general
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using generate
< ok 0 0
> exec ^
: string
| show grants for u using generated
< error
| 1064 42000 …line 1 column 33 near "generated" 
> exec ^
: string
| show grants for u using geomcollection
< ok 0 0
> exec ^
: string
| show grants for u using geometry
< ok 0 0
> exec ^
: string
| show grants for u using geometrycollection
< ok 0 0
> exec ^
: string
| show grants for u using get
< ok 0 0
> exec ^
: string
| show grants for u using get_format
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using get_source_public_key
< ok 0 0
> exec ^
: string
| show grants for u using global
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using grant
< error
| 1064 42000 …line 1 column 29 near "grant" 
> exec ^
: string
| show grants for u using grants
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using greatest
< ok 0 0
> exec ^
: string
| show grants for u using group
< error
| 1064 42000 …line 1 column 29 near "group" 
> exec ^
: string
| show grants for u using group_concat
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using group_replication
< ok 0 0
> exec ^
: string
| show grants for u using grouping
< ok 0 0
> exec ^
: string
| show grants for u using groups
< error
| 1064 42000 …line 1 column 30 near "groups" 
> exec ^
: string
| show grants for u using gtid_only
< ok 0 0
> exec ^
: string
| show grants for u using gtids
< ok 0 0
> exec ^
: string
| show grants for u using guided
< ok 0 0
> exec ^
: string
| show grants for u using handler
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using hash
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using having
< error
| 1064 42000 …line 1 column 30 near "having" 
> exec ^
: string
| show grants for u using header
< ok 0 0
> exec ^
: string
| show grants for u using help
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using high
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using high_priority
< error
| 1064 42000 …line 1 column 37 near "high_priority" 
> exec ^
: string
| show grants for u using histogram
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using histograms_in_flight
< error
| 1064 42000 …line 1 column 44 near "" 
> exec ^
: string
| show grants for u using history
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using host
< ok 0 0
> exec ^
: string
| show grants for u using hosts
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using hour
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using hour_microsecond
< error
| 1064 42000 …line 1 column 40 near "hour_microsecond" 
> exec ^
: string
| show grants for u using hour_minute
< error
| 1064 42000 …line 1 column 35 near "hour_minute" 
> exec ^
: string
| show grants for u using hour_second
< error
| 1064 42000 …line 1 column 35 near "hour_second" 
> exec ^
: string
| show grants for u using identified
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using if
< error
| 1064 42000 …line 1 column 26 near "if" 
> exec ^
: string
| show grants for u using ignore
< error
| 1064 42000 …line 1 column 30 near "ignore" 
> exec ^
: string
| show grants for u using ignore_server_ids
< ok 0 0
> exec ^
: string
| show grants for u using ilike
< error
| 1064 42000 …line 1 column 29 near "ilike" 
> exec ^
: string
| show grants for u using import
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using in
< error
| 1064 42000 …line 1 column 26 near "in" 
> exec ^
: string
| show grants for u using inactive
< ok 0 0
> exec ^
: string
| show grants for u using increment
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using index
< error
| 1064 42000 …line 1 column 29 near "index" 
> exec ^
: string
| show grants for u using indexes
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using infile
< error
| 1064 42000 …line 1 column 30 near "infile" 
> exec ^
: string
| show grants for u using initial
< ok 0 0
> exec ^
: string
| show grants for u using initial_size
< ok 0 0
> exec ^
: string
| show grants for u using initiate
< ok 0 0
> exec ^
: string
| show grants for u using inner
< error
| 1064 42000 …line 1 column 29 near "inner" 
> exec ^
: string
| show grants for u using inout
< error
| 1064 42000 …line 1 column 29 near "inout" 
> exec ^
: string
| show grants for u using insensitive
< ok 0 0
> exec ^
: string
| show grants for u using insert
< error
| 1064 42000 …line 1 column 30 near "insert" 
> exec ^
: string
| show grants for u using insert_method
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using install
< ok 0 0
> exec ^
: string
| show grants for u using instance
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using int
< error
| 1064 42000 …line 1 column 27 near "int" 
> exec ^
: string
| show grants for u using int1
< error
| 1064 42000 …line 1 column 28 near "int1" 
> exec ^
: string
| show grants for u using int2
< error
| 1064 42000 …line 1 column 28 near "int2" 
> exec ^
: string
| show grants for u using int3
< error
| 1064 42000 …line 1 column 28 near "int3" 
> exec ^
: string
| show grants for u using int4
< error
| 1064 42000 …line 1 column 28 near "int4" 
> exec ^
: string
| show grants for u using int8
< error
| 1064 42000 …line 1 column 28 near "int8" 
> exec ^
: string
| show grants for u using integer
< error
| 1064 42000 …line 1 column 31 near "integer" 
> exec ^
: string
| show grants for u using intersect
< error
| 1064 42000 …line 1 column 33 near "intersect" 
> exec ^
: string
| show grants for u using interval
< error
| 1064 42000 …line 1 column 32 near "interval" 
> exec ^
: string
| show grants for u using into
< error
| 1064 42000 …line 1 column 28 near "into" 
> exec ^
: string
| show grants for u using invisible
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using invoker
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using io
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| show grants for u using io_after_gtids
< ok 0 0
> exec ^
: string
| show grants for u using io_before_gtids
< ok 0 0
> exec ^
: string
| show grants for u using io_thread
< ok 0 0
> exec ^
: string
| show grants for u using ipc
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using is
< error
| 1064 42000 …line 1 column 26 near "is" 
> exec ^
: string
| show grants for u using isolation
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using issuer
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using iterate
< error
| 1064 42000 …line 1 column 31 near "iterate" 
> exec ^
: string
| show grants for u using job
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using jobs
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using join
< error
| 1064 42000 …line 1 column 28 near "join" 
> exec ^
: string
| show grants for u using json
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using json_arrayagg
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using json_duality_object
< ok 0 0
> exec ^
: string
| show grants for u using json_objectagg
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using json_table
< ok 0 0
> exec ^
: string
| show grants for u using json_value
< ok 0 0
> exec ^
: string
| show grants for u using key
< error
| 1064 42000 …line 1 column 27 near "key" 
> exec ^
: string
| show grants for u using key_block_size
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using keyring
< ok 0 0
> exec ^
: string
| show grants for u using keys
< error
| 1064 42000 …line 1 column 28 near "keys" 
> exec ^
: string
| show grants for u using kill
< error
| 1064 42000 …line 1 column 28 near "kill" 
> exec ^
: string
| show grants for u using labels
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using lag
< error
| 1064 42000 …line 1 column 27 near "lag" 
> exec ^
: string
| show grants for u using language
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using last
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using last_value
< error
| 1064 42000 …line 1 column 34 near "last_value" 
> exec ^
: string
| show grants for u using lateral
< ok 0 0
> exec ^
: string
| show grants for u using lead
< error
| 1064 42000 …line 1 column 28 near "lead" 
> exec ^
: string
| show grants for u using leader_constraints
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| show grants for u using leading
< error
| 1064 42000 …line 1 column 31 near "leading" 
> exec ^
: string
| show grants for u using learner_constraints
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| show grants for u using learners
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using least
< ok 0 0
> exec ^
: string
| show grants for u using leave
< error
| 1064 42000 …line 1 column 29 near "leave" 
> exec ^
: string
| show grants for u using leaves
< ok 0 0
> exec ^
: string
| show grants for u using left
< error
| 1064 42000 …line 1 column 28 near "left" 
> exec ^
: string
| show grants for u using less
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using level
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using library
< ok 0 0
> exec ^
: string
| show grants for u using like
< error
| 1064 42000 …line 1 column 28 near "like" 
> exec ^
: string
| show grants for u using limit
< error
| 1064 42000 …line 1 column 29 near "limit" 
> exec ^
: string
| show grants for u using linear
< error
| 1064 42000 …line 1 column 30 near "linear" 
> exec ^
: string
| show grants for u using lines
< error
| 1064 42000 …line 1 column 29 near "lines" 
> exec ^
: string
| show grants for u using linestring
< ok 0 0
> exec ^
: string
| show grants for u using list
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using load
< error
| 1064 42000 …line 1 column 28 near "load" 
> exec ^
: string
| show grants for u using local
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using localtime
< error
| 1064 42000 …line 1 column 33 near "localtime" 
> exec ^
: string
| show grants for u using localtimestamp
< error
| 1064 42000 …line 1 column 38 near "localtimestamp" 
> exec ^
: string
| show grants for u using location
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using lock
< error
| 1064 42000 …line 1 column 28 near "lock" 
> exec ^
: string
| show grants for u using locked
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using locks
< ok 0 0
> exec ^
: string
| show grants for u using log
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using logfile
< ok 0 0
> exec ^
: string
| show grants for u using logs
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using long
< error
| 1064 42000 …line 1 column 28 near "long" 
> exec ^
: string
| show grants for u using longblob
< error
| 1064 42000 …line 1 column 32 near "longblob" 
> exec ^
: string
| show grants for u using longtext
< error
| 1064 42000 …line 1 column 32 near "longtext" 
> exec ^
: string
| show grants for u using loop
< ok 0 0
> exec ^
: string
| show grants for u using low
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using low_priority
< error
| 1064 42000 …line 1 column 36 near "low_priority" 
> exec ^
: string
| show grants for u using manual
< ok 0 0
> exec ^
: string
| show grants for u using master
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using master_log_file
< ok 0 0
> exec ^
: string
| show grants for u using match
< error
| 1064 42000 …line 1 column 29 near "match" 
> exec ^
: string
| show grants for u using materialized
< ok 0 0
> exec ^
: string
| show grants for u using max
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using max_connections_per_hour
< error
| 1064 42000 …line 1 column 48 near "" 
> exec ^
: string
| show grants for u using max_queries_per_hour
< error
| 1064 42000 …line 1 column 44 near "" 
> exec ^
: string
| show grants for u using max_rows
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using max_size
< ok 0 0
> exec ^
: string
| show grants for u using max_updates_per_hour
< error
| 1064 42000 …line 1 column 44 near "" 
> exec ^
: string
| show grants for u using max_user_connections
< error
| 1064 42000 …line 1 column 44 near "" 
> exec ^
: string
| show grants for u using maxvalue
< error
| 1064 42000 …line 1 column 32 near "maxvalue" 
> exec ^
: string
| show grants for u using medium
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using mediumblob
< error
| 1064 42000 …line 1 column 34 near "mediumblob" 
> exec ^
: string
| show grants for u using mediumint
< error
| 1064 42000 …line 1 column 33 near "mediumint" 
> exec ^
: string
| show grants for u using mediumtext
< error
| 1064 42000 …line 1 column 34 near "mediumtext" 
> exec ^
: string
| show grants for u using member
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using memory
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using merge
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using message_text
< ok 0 0
> exec ^
: string
| show grants for u using microsecond
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using mid
< ok 0 0
> exec ^
: string
| show grants for u using middleint
< error
| 1064 42000 …line 1 column 33 near "middleint" 
> exec ^
: string
| show grants for u using migrate
< ok 0 0
> exec ^
: string
| show grants for u using min
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using min_rows
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using minute
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using minute_microsecond
< error
| 1064 42000 …line 1 column 42 near "minute_microsecond" 
> exec ^
: string
| show grants for u using minute_second
< error
| 1064 42000 …line 1 column 37 near "minute_second" 
> exec ^
: string
| show grants for u using minvalue
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using mod
< error
| 1064 42000 …line 1 column 27 near "mod" 
> exec ^
: string
| show grants for u using mode
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using modifies
< ok 0 0
> exec ^
: string
| show grants for u using modify
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using month
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using multilinestring
< ok 0 0
> exec ^
: string
| show grants for u using multipoint
< ok 0 0
> exec ^
: string
| show grants for u using multipolygon
< ok 0 0
> exec ^
: string
| show grants for u using mutex
< ok 0 0
> exec ^
: string
| show grants for u using mysql_errno
< ok 0 0
> exec ^
: string
| show grants for u using name
< ok 0 0
> exec ^
: string
| show grants for u using names
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using national
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using natural
< error
| 1064 42000 …line 1 column 31 near "natural" 
> exec ^
: string
| show grants for u using nchar
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using ndb
< ok 0 0
> exec ^
: string
| show grants for u using ndbcluster
< ok 0 0
> exec ^
: string
| show grants for u using nested
< ok 0 0
> exec ^
: string
| show grants for u using network_namespace
< ok 0 0
> exec ^
: string
| show grants for u using never
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using new
< ok 0 0
> exec ^
: string
| show grants for u using next
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using next_row_id
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using no
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| show grants for u using no_wait
< ok 0 0
> exec ^
: string
| show grants for u using no_write_to_binlog
< error
| 1064 42000 …line 1 column 42 near "no_write_to_binlog" 
> exec ^
: string
| show grants for u using nocache
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using nocycle
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using nodegroup
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using nomaxvalue
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using nominvalue
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using nonclustered
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using none
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using not
< error
| 1064 42000 …line 1 column 27 near "not" 
> exec ^
: string
| show grants for u using now
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using nowait
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using nth_value
< error
| 1064 42000 …line 1 column 33 near "nth_value" 
> exec ^
: string
| show grants for u using ntile
< error
| 1064 42000 …line 1 column 29 near "ntile" 
> exec ^
: string
| show grants for u using null
< error
| 1064 42000 …line 1 column 28 near "null" 
> exec ^
: string
| show grants for u using nullif
< ok 0 0
> exec ^
: string
| show grants for u using nulls
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using number
< ok 0 0
> exec ^
: string
| show grants for u using numeric
< error
| 1064 42000 …line 1 column 31 near "numeric" 
> exec ^
: string
| show grants for u using nvarchar
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using of
< error
| 1064 42000 …line 1 column 26 near "of" 
> exec ^
: string
| show grants for u using off
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using offset
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using oj
< ok 0 0
> exec ^
: string
| show grants for u using old
< ok 0 0
> exec ^
: string
| show grants for u using oltp_read_only
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using oltp_read_write
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| show grants for u using oltp_write_only
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| show grants for u using on
< error
| 1064 42000 …line 1 column 26 near "on" 
> exec ^
: string
| show grants for u using one
< ok 0 0
> exec ^
: string
| show grants for u using only
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using open
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using optimize
< error
| 1064 42000 …line 1 column 32 near "optimize" 
> exec ^
: string
| show grants for u using optimizer_costs
< ok 0 0
> exec ^
: string
| show grants for u using option
< error
| 1064 42000 …line 1 column 30 near "option" 
> exec ^
: string
| show grants for u using optional
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using optionally
< error
| 1064 42000 …line 1 column 34 near "optionally" 
> exec ^
: string
| show grants for u using options
< ok 0 0
> exec ^
: string
| show grants for u using or
< error
| 1064 42000 …line 1 column 26 near "or" 
> exec ^
: string
| show grants for u using order
< error
| 1064 42000 …line 1 column 29 near "order" 
> exec ^
: string
| show grants for u using ordinality
< ok 0 0
> exec ^
: string
| show grants for u using organization
< ok 0 0
> exec ^
: string
| show grants for u using others
< ok 0 0
> exec ^
: string
| show grants for u using out
< error
| 1064 42000 …line 1 column 27 near "out" 
> exec ^
: string
| show grants for u using outer
< error
| 1064 42000 …line 1 column 29 near "outer" 
> exec ^
: string
| show grants for u using outfile
< error
| 1064 42000 …line 1 column 31 near "outfile" 
> exec ^
: string
| show grants for u using over
< error
| 1064 42000 …line 1 column 28 near "over" 
> exec ^
: string
| show grants for u using overlay
< ok 0 0
> exec ^
: string
| show grants for u using owner
< ok 0 0
> exec ^
: string
| show grants for u using pack_keys
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using page
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using parallel
< ok 0 0
> exec ^
: string
| show grants for u using parameters
< ok 0 0
> exec ^
: string
| show grants for u using parse_tree
< ok 0 0
> exec ^
: string
| show grants for u using parser
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using partial
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using partition
< error
| 1064 42000 …line 1 column 33 near "partition" 
> exec ^
: string
| show grants for u using partitioning
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using partitions
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using password
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using password_lock_time
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| show grants for u using path
< ok 0 0
> exec ^
: string
| show grants for u using percent_rank
< error
| 1064 42000 …line 1 column 36 near "percent_rank" 
> exec ^
: string
| show grants for u using persist
< ok 0 0
> exec ^
: string
| show grants for u using persist_only
< ok 0 0
> exec ^
: string
| show grants for u using phase
< ok 0 0
> exec ^
: string
| show grants for u using placement
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using placing
< ok 0 0
> exec ^
: string
| show grants for u using plan
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using plugin
< ok 0 0
> exec ^
: string
| show grants for u using plugin_dir
< ok 0 0
> exec ^
: string
| show grants for u using plugins
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using point
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using policy
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using polygon
< ok 0 0
> exec ^
: string
| show grants for u using port
< ok 0 0
> exec ^
: string
| show grants for u using position
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using pre_split_regions
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using precedes
< ok 0 0
> exec ^
: string
| show grants for u using preceding
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using precision
< error
| 1064 42000 …line 1 column 33 near "precision" 
> exec ^
: string
| show grants for u using prepare
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using preserve
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using prev
< ok 0 0
> exec ^
: string
| show grants for u using primary
< error
| 1064 42000 …line 1 column 31 near "primary" 
> exec ^
: string
| show grants for u using primary_region
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using priority
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using privilege_checks_user
< ok 0 0
> exec ^
: string
| show grants for u using privileges
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using procedure
< error
| 1064 42000 …line 1 column 33 near "procedure" 
> exec ^
: string
| show grants for u using process
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using processed_keys
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using processlist
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using profile
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using profiles
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using proxy
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using purge
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using qualify
< ok 0 0
> exec ^
: string
| show grants for u using quarter
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using query
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using query_limit
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using quick
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using random
< ok 0 0
> exec ^
: string
| show grants for u using range
< error
| 1064 42000 …line 1 column 29 near "range" 
> exec ^
: string
| show grants for u using rank
< error
| 1064 42000 …line 1 column 28 near "rank" 
> exec ^
: string
| show grants for u using read
< error
| 1064 42000 …line 1 column 28 near "read" 
> exec ^
: string
| show grants for u using read_only
< ok 0 0
> exec ^
: string
| show grants for u using read_write
< ok 0 0
> exec ^
: string
| show grants for u using reads
< ok 0 0
> exec ^
: string
| show grants for u using real
< error
| 1064 42000 …line 1 column 28 near "real" 
> exec ^
: string
| show grants for u using rebuild
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using recover
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using recursive
< error
| 1064 42000 …line 1 column 33 near "recursive" 
> exec ^
: string
| show grants for u using redo_buffer_size
< ok 0 0
> exec ^
: string
| show grants for u using redundant
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using reference
< ok 0 0
> exec ^
: string
| show grants for u using references
< error
| 1064 42000 …line 1 column 34 near "references" 
> exec ^
: string
| show grants for u using regexp
< error
| 1064 42000 …line 1 column 30 near "regexp" 
> exec ^
: string
| show grants for u using region
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using regions
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using registration
< ok 0 0
> exec ^
: string
| show grants for u using relational
< ok 0 0
> exec ^
: string
| show grants for u using relay
< ok 0 0
> exec ^
: string
| show grants for u using relay_log_file
< ok 0 0
> exec ^
: string
| show grants for u using relay_log_pos
< ok 0 0
> exec ^
: string
| show grants for u using relay_thread
< ok 0 0
> exec ^
: string
| show grants for u using relaylog
< ok 0 0
> exec ^
: string
| show grants for u using release
< error
| 1064 42000 …line 1 column 31 near "release" 
> exec ^
: string
| show grants for u using reload
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using remove
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using rename
< error
| 1064 42000 …line 1 column 30 near "rename" 
> exec ^
: string
| show grants for u using reorganize
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using repair
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using repeat
< error
| 1064 42000 …line 1 column 30 near "repeat" 
> exec ^
: string
| show grants for u using repeatable
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using replace
< error
| 1064 42000 …line 1 column 31 near "replace" 
> exec ^
: string
| show grants for u using replica
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using replicas
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using replicate_do_db
< ok 0 0
> exec ^
: string
| show grants for u using replicate_do_table
< ok 0 0
> exec ^
: string
| show grants for u using replicate_ignore_db
< ok 0 0
> exec ^
: string
| show grants for u using replicate_ignore_table
< ok 0 0
> exec ^
: string
| show grants for u using replicate_rewrite_db
< ok 0 0
> exec ^
: string
| show grants for u using replicate_wild_do_table
< ok 0 0
> exec ^
: string
| show grants for u using replicate_wild_ignore_table
< ok 0 0
> exec ^
: string
| show grants for u using replication
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using require
< error
| 1064 42000 …line 1 column 31 near "require" 
> exec ^
: string
| show grants for u using require_row_format
< ok 0 0
> exec ^
: string
| show grants for u using require_table_primary_key_check
< ok 0 0
> exec ^
: string
| show grants for u using reset
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using resignal
< ok 0 0
> exec ^
: string
| show grants for u using resource
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using respect
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using restart
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using restore
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using restores
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using restrict
< error
| 1064 42000 …line 1 column 32 near "restrict" 
> exec ^
: string
| show grants for u using resume
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using retain
< ok 0 0
> exec ^
: string
| show grants for u using return
< ok 0 0
> exec ^
: string
| show grants for u using returned_sqlstate
< ok 0 0
> exec ^
: string
| show grants for u using returning
< ok 0 0
> exec ^
: string
| show grants for u using returns
< ok 0 0
> exec ^
: string
| show grants for u using reuse
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using reverse
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using revoke
< error
| 1064 42000 …line 1 column 30 near "revoke" 
> exec ^
: string
| show grants for u using right
< error
| 1064 42000 …line 1 column 29 near "right" 
> exec ^
: string
| show grants for u using rlike
< error
| 1064 42000 …line 1 column 29 near "rlike" 
> exec ^
: string
| show grants for u using role
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using rollback
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using rollup
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using rotate
< ok 0 0
> exec ^
: string
| show grants for u using routine
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using row
< error
| 1064 42000 …line 1 column 27 near "row" 
> exec ^
: string
| show grants for u using row_count
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using row_format
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using row_number
< error
| 1064 42000 …line 1 column 34 near "row_number" 
> exec ^
: string
| show grants for u using rows
< error
| 1064 42000 …line 1 column 28 near "rows" 
> exec ^
: string
| show grants for u using rtree
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using ru
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| show grants for u using ru_per_sec
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using rule
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using run
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using s3
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| show grants for u using savepoint
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using schedule
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using schema
< error
| 1064 42000 …line 1 column 30 near "schema" 
> exec ^
: string
| show grants for u using schema_name
< ok 0 0
> exec ^
: string
| show grants for u using schemas
< error
| 1064 42000 …line 1 column 31 near "schemas" 
> exec ^
: string
| show grants for u using second
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using second_microsecond
< error
| 1064 42000 …line 1 column 42 near "second_microsecond" 
> exec ^
: string
| show grants for u using secondary
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using secondary_engine
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using secondary_engine_attribute
< ok 0 0
> exec ^
: string
| show grants for u using secondary_load
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using secondary_unload
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using security
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using select
< error
| 1064 42000 …line 1 column 30 near "select" 
> exec ^
: string
| show grants for u using sensitive
< ok 0 0
> exec ^
: string
| show grants for u using separator
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using sequence
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using serial
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using serializable
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using server
< ok 0 0
> exec ^
: string
| show grants for u using session
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using session_states
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using session_user
< ok 0 0
> exec ^
: string
| show grants for u using set
< error
| 1064 42000 …line 1 column 27 near "set" 
> exec ^
: string
| show grants for u using sets
< ok 0 0
> exec ^
: string
| show grants for u using shard_row_id_bits
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using share
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using show
< error
| 1064 42000 …line 1 column 28 near "show" 
> exec ^
: string
| show grants for u using shutdown
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using signal
< ok 0 0
> exec ^
: string
| show grants for u using signed
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using similar
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using simple
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using skip
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using slave
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using slow
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using smallint
< error
| 1064 42000 …line 1 column 32 near "smallint" 
> exec ^
: string
| show grants for u using snapshot
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using socket
< ok 0 0
> exec ^
: string
| show grants for u using some
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using soname
< ok 0 0
> exec ^
: string
| show grants for u using sounds
< ok 0 0
> exec ^
: string
| show grants for u using source
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using source_auto_position
< ok 0 0
> exec ^
: string
| show grants for u using source_bind
< ok 0 0
> exec ^
: string
| show grants for u using source_compression_algorithms
< ok 0 0
> exec ^
: string
| show grants for u using source_connect_retry
< ok 0 0
> exec ^
: string
| show grants for u using source_connection_auto_failover
< ok 0 0
> exec ^
: string
| show grants for u using source_delay
< ok 0 0
> exec ^
: string
| show grants for u using source_heartbeat_period
< ok 0 0
> exec ^
: string
| show grants for u using source_host
< ok 0 0
> exec ^
: string
| show grants for u using source_log_file
< ok 0 0
> exec ^
: string
| show grants for u using source_log_pos
< ok 0 0
> exec ^
: string
| show grants for u using source_password
< ok 0 0
> exec ^
: string
| show grants for u using source_port
< ok 0 0
> exec ^
: string
| show grants for u using source_public_key_path
< ok 0 0
> exec ^
: string
| show grants for u using source_retry_count
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_ca
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_capath
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_cert
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_cipher
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_crl
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_crlpath
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_key
< ok 0 0
> exec ^
: string
| show grants for u using source_ssl_verify_server_cert
< ok 0 0
> exec ^
: string
| show grants for u using source_tls_ciphersuites
< ok 0 0
> exec ^
: string
| show grants for u using source_tls_version
< ok 0 0
> exec ^
: string
| show grants for u using source_user
< ok 0 0
> exec ^
: string
| show grants for u using source_zstd_compression_level
< ok 0 0
> exec ^
: string
| show grants for u using spatial
< error
| 1064 42000 …line 1 column 31 near "spatial" 
> exec ^
: string
| show grants for u using specific
< ok 0 0
> exec ^
: string
| show grants for u using split
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using sql
< error
| 1064 42000 …line 1 column 27 near "sql" 
> exec ^
: string
| show grants for u using sql_after_gtids
< ok 0 0
> exec ^
: string
| show grants for u using sql_after_mts_gaps
< ok 0 0
> exec ^
: string
| show grants for u using sql_before_gtids
< ok 0 0
> exec ^
: string
| show grants for u using sql_big_result
< error
| 1064 42000 …line 1 column 38 near "sql_big_result" 
> exec ^
: string
| show grants for u using sql_buffer_result
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using sql_calc_found_rows
< error
| 1064 42000 …line 1 column 43 near "sql_calc_found_rows" 
> exec ^
: string
| show grants for u using sql_no_cache
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using sql_small_result
< error
| 1064 42000 …line 1 column 40 near "sql_small_result" 
> exec ^
: string
| show grants for u using sql_thread
< ok 0 0
> exec ^
: string
| show grants for u using sql_tsi_day
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using sql_tsi_hour
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using sql_tsi_minute
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using sql_tsi_month
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using sql_tsi_quarter
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| show grants for u using sql_tsi_second
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using sql_tsi_week
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using sql_tsi_year
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using sqlexception
< error
| 1064 42000 …line 1 column 36 near "sqlexception" 
> exec ^
: string
| show grants for u using sqlstate
< error
| 1064 42000 …line 1 column 32 near "sqlstate" 
> exec ^
: string
| show grants for u using sqlwarning
< error
| 1064 42000 …line 1 column 34 near "sqlwarning" 
> exec ^
: string
| show grants for u using srid
< ok 0 0
> exec ^
: string
| show grants for u using ssl
< error
| 1064 42000 …line 1 column 27 near "ssl" 
> exec ^
: string
| show grants for u using st_collect
< ok 0 0
> exec ^
: string
| show grants for u using stacked
< ok 0 0
> exec ^
: string
| show grants for u using start
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using start_time
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using starting
< error
| 1064 42000 …line 1 column 32 near "starting" 
> exec ^
: string
| show grants for u using starts
< ok 0 0
> exec ^
: string
| show grants for u using stats_auto_recalc
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using stats_buckets
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using stats_col_choice
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using stats_col_list
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using stats_extended
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using stats_healthy
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using stats_histograms
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using stats_locked
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using stats_meta
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using stats_persistent
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using stats_sample_pages
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| show grants for u using stats_sample_rate
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using stats_topn
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using status
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using std
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using stddev
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using stddev_pop
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using stddev_samp
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using stop
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using storage
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using stored
< error
| 1064 42000 …line 1 column 30 near "stored" 
> exec ^
: string
| show grants for u using straight_join
< error
| 1064 42000 …line 1 column 37 near "straight_join" 
> exec ^
: string
| show grants for u using stream
< ok 0 0
> exec ^
: string
| show grants for u using strict_load
< ok 0 0
> exec ^
: string
| show grants for u using string
< ok 0 0
> exec ^
: string
| show grants for u using subclass_origin
< ok 0 0
> exec ^
: string
| show grants for u using subdate
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using subject
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using subpartition
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using subpartitions
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using substr
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using substring
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using sum
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using super
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using survival_preferences
< error
| 1064 42000 …line 1 column 44 near "" 
> exec ^
: string
| show grants for u using suspend
< ok 0 0
> exec ^
: string
| show grants for u using swaps
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using switch_group
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using switches
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using sysdate
< ok 0 0
> exec ^
: string
| show grants for u using system
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using system_user
< ok 0 0
> exec ^
: string
| show grants for u using table
< error
| 1064 42000 …line 1 column 29 near "table" 
> exec ^
: string
| show grants for u using table_checksum
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| show grants for u using table_name
< ok 0 0
> exec ^
: string
| show grants for u using tables
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using tablesample
< error
| 1064 42000 …line 1 column 35 near "tablesample" 
> exec ^
: string
| show grants for u using tablespace
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using task_types
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using temporary
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using temptable
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using terminated
< error
| 1064 42000 …line 1 column 34 near "terminated" 
> exec ^
: string
| show grants for u using text
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using than
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using then
< error
| 1064 42000 …line 1 column 28 near "then" 
> exec ^
: string
| show grants for u using thread_priority
< ok 0 0
> exec ^
: string
| show grants for u using ties
< ok 0 0
> exec ^
: string
| show grants for u using tiflash
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using time
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using timeout
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using timestamp
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using timestampadd
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| show grants for u using timestampdiff
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using tinyblob
< error
| 1064 42000 …line 1 column 32 near "tinyblob" 
> exec ^
: string
| show grants for u using tinyint
< error
| 1064 42000 …line 1 column 31 near "tinyint" 
> exec ^
: string
| show grants for u using tinytext
< error
| 1064 42000 …line 1 column 32 near "tinytext" 
> exec ^
: string
| show grants for u using tls
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using to
< error
| 1064 42000 …line 1 column 26 near "to" 
> exec ^
: string
| show grants for u using tpcc
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using tpch_10
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using trailing
< error
| 1064 42000 …line 1 column 32 near "trailing" 
> exec ^
: string
| show grants for u using transaction
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using trigger
< error
| 1064 42000 …line 1 column 31 near "trigger" 
> exec ^
: string
| show grants for u using triggers
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using trim
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using true
< error
| 1064 42000 …line 1 column 28 near "true" 
> exec ^
: string
| show grants for u using truncate
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using ttl
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| show grants for u using ttl_enable
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using ttl_job_interval
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| show grants for u using type
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using types
< ok 0 0
> exec ^
: string
| show grants for u using unbounded
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using uncommitted
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| show grants for u using undefined
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using undo
< ok 0 0
> exec ^
: string
| show grants for u using undo_buffer_size
< ok 0 0
> exec ^
: string
| show grants for u using undofile
< ok 0 0
> exec ^
: string
| show grants for u using unicode
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using uninstall
< ok 0 0
> exec ^
: string
| show grants for u using union
< error
| 1064 42000 …line 1 column 29 near "union" 
> exec ^
: string
| show grants for u using unique
< error
| 1064 42000 …line 1 column 30 near "unique" 
> exec ^
: string
| show grants for u using unknown
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using unlimited
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using unlock
< error
| 1064 42000 …line 1 column 30 near "unlock" 
> exec ^
: string
| show grants for u using unregister
< ok 0 0
> exec ^
: string
| show grants for u using unsigned
< error
| 1064 42000 …line 1 column 32 near "unsigned" 
> exec ^
: string
| show grants for u using until
< error
| 1064 42000 …line 1 column 29 near "until" 
> exec ^
: string
| show grants for u using update
< error
| 1064 42000 …line 1 column 30 near "update" 
> exec ^
: string
| show grants for u using upgrade
< ok 0 0
> exec ^
: string
| show grants for u using uri
< ok 0 0
> exec ^
: string
| show grants for u using url
< ok 0 0
> exec ^
: string
| show grants for u using usage
< error
| 1064 42000 …line 1 column 29 near "usage" 
> exec ^
: string
| show grants for u using use
< error
| 1064 42000 …line 1 column 27 near "use" 
> exec ^
: string
| show grants for u using use_frm
< ok 0 0
> exec ^
: string
| show grants for u using user
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using user_resources
< ok 0 0
> exec ^
: string
| show grants for u using using
< error
| 1064 42000 …line 1 column 29 near "using" 
> exec ^
: string
| show grants for u using utc_date
< error
| 1064 42000 …line 1 column 32 near "utc_date" 
> exec ^
: string
| show grants for u using utc_time
< error
| 1064 42000 …line 1 column 32 near "utc_time" 
> exec ^
: string
| show grants for u using utc_timestamp
< error
| 1064 42000 …line 1 column 37 near "utc_timestamp" 
> exec ^
: string
| show grants for u using utilization_limit
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using validate
< ok 0 0
> exec ^
: string
| show grants for u using validation
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| show grants for u using value
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using values
< error
| 1064 42000 …line 1 column 30 near "values" 
> exec ^
: string
| show grants for u using var_pop
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using var_samp
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using varbinary
< error
| 1064 42000 …line 1 column 33 near "varbinary" 
> exec ^
: string
| show grants for u using varchar
< error
| 1064 42000 …line 1 column 31 near "varchar" 
> exec ^
: string
| show grants for u using varcharacter
< error
| 1064 42000 …line 1 column 36 near "varcharacter" 
> exec ^
: string
| show grants for u using variables
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| show grants for u using variance
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using varying
< error
| 1064 42000 …line 1 column 31 near "varying" 
> exec ^
: string
| show grants for u using vcpu
< ok 0 0
> exec ^
: string
| show grants for u using vector
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using verify_key_constraints
< ok 0 0
> exec ^
: string
| show grants for u using view
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using virtual
< error
| 1064 42000 …line 1 column 31 near "virtual" 
> exec ^
: string
| show grants for u using visible
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using voter_constraints
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| show grants for u using voters
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| show grants for u using wait
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using warnings
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using watch
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| show grants for u using week
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using weight_string
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| show grants for u using when
< error
| 1064 42000 …line 1 column 28 near "when" 
> exec ^
: string
| show grants for u using where
< error
| 1064 42000 …line 1 column 29 near "where" 
> exec ^
: string
| show grants for u using while
< error
| 1064 42000 …line 1 column 29 near "while" 
> exec ^
: string
| show grants for u using window
< error
| 1064 42000 …line 1 column 30 near "window" 
> exec ^
: string
| show grants for u using with
< error
| 1064 42000 …line 1 column 28 near "with" 
> exec ^
: string
| show grants for u using without
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| show grants for u using work
< ok 0 0
> exec ^
: string
| show grants for u using workload
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| show grants for u using wrapper
< ok 0 0
> exec ^
: string
| show grants for u using write
< error
| 1064 42000 …line 1 column 29 near "write" 
> exec ^
: string
| show grants for u using x509
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using xa
< ok 0 0
> exec ^
: string
| show grants for u using xid
< ok 0 0
> exec ^
: string
| show grants for u using xml
< ok 0 0
> exec ^
: string
| show grants for u using xor
< error
| 1064 42000 …line 1 column 27 near "xor" 
> exec ^
: string
| show grants for u using year
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| show grants for u using year_month
< error
| 1064 42000 …line 1 column 34 near "year_month" 
> exec ^
: string
| show grants for u using zerofill
< error
| 1064 42000 …line 1 column 32 near "zerofill" 
> exec ^
: string
| show grants for u using zone
< ok 0 0

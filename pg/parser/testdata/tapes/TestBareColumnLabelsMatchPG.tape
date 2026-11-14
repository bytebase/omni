-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestBareColumnLabelsMatchPG default#1
> begin
< ok
> exec
| SELECT 1 abort
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 absent
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 absolute
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 access
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 action
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 add
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 admin
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 after
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 aggregate
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 all
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 also
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 alter
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 always
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 analyse
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 analyze
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 and
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 any
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 array
< error
| 42601 syntax error at or near "array"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 as
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 asc
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 asensitive
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 assertion
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 assignment
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 asymmetric
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 at
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 atomic
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 attach
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 attribute
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 authorization
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 backward
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 before
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 begin
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 between
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 bigint
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 binary
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 bit
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 boolean
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 both
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 breadth
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 by
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cache
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 call
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 called
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cascade
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cascaded
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 case
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cast
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 catalog
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 chain
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 char
< error
| 42601 syntax error at or near "char"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 character
< error
| 42601 syntax error at or near "character"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 characteristics
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 check
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 checkpoint
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 class
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 close
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cluster
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 coalesce
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 collate
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 collation
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 column
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 columns
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 comment
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 comments
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 commit
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 committed
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 compression
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 concurrently
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 conditional
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 configuration
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 conflict
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 connection
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 constraint
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 constraints
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 content
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 continue
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 conversion
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 copy
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cost
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 create
< error
| 42601 syntax error at or near "create"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cross
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 csv
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cube
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_catalog
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_date
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_role
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_schema
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_time
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_timestamp
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 current_user
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cursor
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 cycle
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 data
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 database
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 day
< error
| 42601 syntax error at or near "day"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 deallocate
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 dec
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 decimal
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 declare
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 default
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 defaults
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 deferrable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 deferred
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 definer
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 delete
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 delimiter
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 delimiters
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 depends
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 depth
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 desc
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 detach
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 dictionary
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 disable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 discard
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 distinct
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 do
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 document
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 domain
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 double
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 drop
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 each
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 else
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 empty
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 enable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 encoding
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 encrypted
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 end
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 enum
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 error
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 escape
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 event
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 except
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 exclude
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 excluding
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 exclusive
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 execute
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 exists
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 explain
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 expression
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 extension
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 external
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 extract
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 false
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 family
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 fetch
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 filter
< error
| 42601 syntax error at or near "filter"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 finalize
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 first
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 float
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 following
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 for
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 force
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 foreign
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 format
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 forward
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 freeze
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 from
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 full
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 function
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 functions
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 generated
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 global
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 grant
< error
| 42601 syntax error at or near "grant"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 granted
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 greatest
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 group
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 grouping
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 groups
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 handler
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 having
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 header
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 hold
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 hour
< error
| 42601 syntax error at or near "hour"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 identity
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 if
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ilike
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 immediate
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 immutable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 implicit
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 import
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 in
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 include
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 including
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 increment
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 indent
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 index
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 indexes
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 inherit
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 inherits
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 initially
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 inline
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 inner
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 inout
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 input
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 insensitive
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 insert
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 instead
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 int
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 integer
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 intersect
< error
| 42601 syntax error at end of input
| position: 19
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 interval
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 into
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 invoker
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 is
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 isnull
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 isolation
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 join
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_array
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_arrayagg
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_exists
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_object
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_objectagg
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_query
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_scalar
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_serialize
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_table
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 json_value
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 keep
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 key
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 keys
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 label
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 language
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 large
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 last
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 lateral
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 leading
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 leakproof
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 least
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 left
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 level
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 like
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 limit
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 listen
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 load
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 local
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 localtime
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 localtimestamp
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 location
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 lock
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 locked
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 logged
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 mapping
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 match
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 matched
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 materialized
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 maxvalue
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 merge
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 merge_action
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 method
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 minute
< error
| 42601 syntax error at or near "minute"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 minvalue
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 mode
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 month
< error
| 42601 syntax error at or near "month"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 move
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 name
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 names
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 national
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 natural
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nchar
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nested
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 new
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 next
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nfc
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nfd
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nfkc
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nfkd
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 no
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 none
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 normalize
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 normalized
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 not
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nothing
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 notify
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 notnull
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nowait
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 null
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nullif
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 nulls
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 numeric
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 object
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 of
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 off
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 offset
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 oids
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 old
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 omit
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 on
< error
| 42601 syntax error at or near "on"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 only
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 operator
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 option
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 options
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 or
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 order
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ordinality
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 others
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 out
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 outer
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 over
< error
| 42601 syntax error at or near "over"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 overlaps
< error
| 42601 syntax error at or near "overlaps"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 overlay
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 overriding
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 owned
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 owner
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 parallel
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 parameter
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 parser
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 partial
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 partition
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 passing
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 password
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 path
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 placing
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 plan
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 plans
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 policy
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 position
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 preceding
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 precision
< error
| 42601 syntax error at or near "precision"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 prepare
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 prepared
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 preserve
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 primary
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 prior
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 privileges
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 procedural
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 procedure
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 procedures
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 program
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 publication
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 quote
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 quotes
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 range
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 read
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 real
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 reassign
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 recheck
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 recursive
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ref
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 references
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 referencing
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 refresh
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 reindex
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 relative
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 release
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 rename
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 repeatable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 replace
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 replica
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 reset
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 restart
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 restrict
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 return
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 returning
< error
| 42601 syntax error at or near "returning"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 returns
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 revoke
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 right
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 role
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 rollback
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 rollup
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 routine
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 routines
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 row
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 rows
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 rule
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 savepoint
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 scalar
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 schema
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 schemas
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 scroll
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 search
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 second
< error
| 42601 syntax error at or near "second"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 security
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 select
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 sequence
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 sequences
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 serializable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 server
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 session
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 session_user
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 set
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 setof
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 sets
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 share
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 show
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 similar
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 simple
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 skip
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 smallint
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 snapshot
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 some
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 source
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 sql
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 stable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 standalone
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 start
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 statement
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 statistics
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 stdin
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 stdout
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 storage
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 stored
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 strict
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 string
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 strip
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 subscription
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 substring
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 support
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 symmetric
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 sysid
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 system
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 system_user
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 table
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 tables
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 tablesample
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 tablespace
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 target
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 temp
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 template
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 temporary
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 text
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 then
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 ties
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 time
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 timestamp
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 to
< error
| 42601 syntax error at or near "to"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 trailing
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 transaction
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 transform
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 treat
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 trigger
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 trim
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 true
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 truncate
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 trusted
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 type
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 types
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 uescape
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unbounded
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 uncommitted
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unconditional
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unencrypted
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 union
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unique
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unknown
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unlisten
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 unlogged
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 until
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 update
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 user
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 using
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 vacuum
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 valid
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 validate
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 validator
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 value
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 values
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 varchar
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 variadic
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 varying
< error
| 42601 syntax error at or near "varying"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 verbose
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 version
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 view
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 views
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 volatile
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 when
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 where
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 whitespace
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 window
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 with
< error
| 42601 syntax error at or near "with"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 within
< error
| 42601 syntax error at or near "within"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 without
< error
| 42601 syntax error at or near "without"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 work
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 wrapper
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 write
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xml
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlattributes
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlconcat
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlelement
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlexists
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlforest
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlnamespaces
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlparse
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlpi
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlroot
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmlserialize
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 xmltable
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 year
< error
| 42601 syntax error at or near "year"
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 yes
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 zone
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1 omni_probe_ident
< ok 1 -
> rollback
< ok

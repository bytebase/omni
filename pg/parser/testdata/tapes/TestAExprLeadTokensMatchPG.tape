-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestAExprLeadTokensMatchPG default#1
> begin
< ok
> exec
| SELECT abort
< error
| 42703 column "abort" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT absent
< error
| 42703 column "absent" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT absolute
< error
| 42703 column "absolute" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT access
< error
| 42703 column "access" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT action
< error
| 42703 column "action" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT add
< error
| 42703 column "add" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT admin
< error
| 42703 column "admin" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT after
< error
| 42703 column "after" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT aggregate
< error
| 42703 column "aggregate" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT all
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT also
< error
| 42703 column "also" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT alter
< error
| 42703 column "alter" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT always
< error
| 42703 column "always" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT analyse
< error
| 42601 syntax error at or near "analyse"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT analyze
< error
| 42601 syntax error at or near "analyze"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT and
< error
| 42601 syntax error at or near "and"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT any
< error
| 42601 syntax error at or near "any"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT ARRAY[1]
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT as
< error
| 42601 syntax error at or near "as"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT asc
< error
| 42601 syntax error at or near "asc"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT asensitive
< error
| 42703 column "asensitive" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT assertion
< error
| 42703 column "assertion" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT assignment
< error
| 42703 column "assignment" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT asymmetric
< error
| 42601 syntax error at or near "asymmetric"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT at
< error
| 42703 column "at" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT atomic
< error
| 42703 column "atomic" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT attach
< error
| 42703 column "attach" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT attribute
< error
| 42703 column "attribute" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT authorization
< error
| 42601 syntax error at end of input
| position: 21
> rollback
< ok
> begin
< ok
> exec
| SELECT backward
< error
| 42703 column "backward" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT before
< error
| 42703 column "before" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT begin
< error
| 42703 column "begin" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT between
< error
| 42703 column "between" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT bigint
< error
| 42703 column "bigint" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT binary
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT bit
< error
| 42703 column "bit" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT boolean
< error
| 42703 column "boolean" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT both
< error
| 42601 syntax error at or near "both"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT breadth
< error
| 42703 column "breadth" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT by
< error
| 42703 column "by" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cache
< error
| 42703 column "cache" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT call
< error
| 42703 column "call" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT called
< error
| 42703 column "called" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cascade
< error
| 42703 column "cascade" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cascaded
< error
| 42703 column "cascaded" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT CASE WHEN TRUE THEN 1 END
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CAST(NULL AS int)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT catalog
< error
| 42703 column "catalog" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT chain
< error
| 42703 column "chain" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT char
< error
| 42703 column "char" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT character
< error
| 42703 column "character" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT characteristics
< error
| 42703 column "characteristics" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT check
< error
| 42601 syntax error at or near "check"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT checkpoint
< error
| 42703 column "checkpoint" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT class
< error
| 42703 column "class" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT close
< error
| 42703 column "close" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cluster
< error
| 42703 column "cluster" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT COALESCE(1, 2)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT collate
< error
| 42601 syntax error at or near "collate"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT collation
< error
| 42601 syntax error at end of input
| position: 17
> rollback
< ok
> begin
< ok
> exec
| SELECT column
< error
| 42601 syntax error at or near "column"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT columns
< error
| 42703 column "columns" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT comment
< error
| 42703 column "comment" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT comments
< error
| 42703 column "comments" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT commit
< error
| 42703 column "commit" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT committed
< error
| 42703 column "committed" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT compression
< error
| 42703 column "compression" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT concurrently
< error
| 42601 syntax error at end of input
| position: 20
> rollback
< ok
> begin
< ok
> exec
| SELECT conditional
< error
| 42703 column "conditional" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT configuration
< error
| 42703 column "configuration" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT conflict
< error
| 42703 column "conflict" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT connection
< error
| 42703 column "connection" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT constraint
< error
| 42601 syntax error at or near "constraint"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT constraints
< error
| 42703 column "constraints" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT content
< error
| 42703 column "content" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT continue
< error
| 42703 column "continue" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT conversion
< error
| 42703 column "conversion" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT copy
< error
| 42703 column "copy" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cost
< error
| 42703 column "cost" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT create
< error
| 42601 syntax error at or near "create"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cross
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT csv
< error
| 42703 column "csv" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cube
< error
| 42703 column "cube" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT current
< error
| 42703 column "current" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT current_catalog
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CURRENT_DATE
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT current_role
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT current_schema
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CURRENT_TIME
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CURRENT_TIMESTAMP
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT CURRENT_USER
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT cursor
< error
| 42703 column "cursor" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT cycle
< error
| 42703 column "cycle" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT data
< error
| 42703 column "data" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT database
< error
| 42703 column "database" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT day
< error
| 42703 column "day" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT deallocate
< error
| 42703 column "deallocate" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT dec
< error
| 42703 column "dec" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT decimal
< error
| 42703 column "decimal" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT declare
< error
| 42703 column "declare" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT default
< error
| 42601 DEFAULT is not allowed in this context
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT defaults
< error
| 42703 column "defaults" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT deferrable
< error
| 42601 syntax error at or near "deferrable"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT deferred
< error
| 42703 column "deferred" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT definer
< error
| 42703 column "definer" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT delete
< error
| 42703 column "delete" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT delimiter
< error
| 42703 column "delimiter" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT delimiters
< error
| 42703 column "delimiters" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT depends
< error
| 42703 column "depends" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT depth
< error
| 42703 column "depth" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT desc
< error
| 42601 syntax error at or near "desc"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT detach
< error
| 42703 column "detach" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT dictionary
< error
| 42703 column "dictionary" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT disable
< error
| 42703 column "disable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT discard
< error
| 42703 column "discard" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT distinct
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT do
< error
| 42601 syntax error at or near "do"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT document
< error
| 42703 column "document" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT domain
< error
| 42703 column "domain" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT double
< error
| 42703 column "double" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT drop
< error
| 42703 column "drop" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT each
< error
| 42703 column "each" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT else
< error
| 42601 syntax error at or near "else"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT empty
< error
| 42703 column "empty" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT enable
< error
| 42703 column "enable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT encoding
< error
| 42703 column "encoding" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT encrypted
< error
| 42703 column "encrypted" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT end
< error
| 42601 syntax error at or near "end"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT enum
< error
| 42703 column "enum" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT error
< error
| 42703 column "error" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT escape
< error
| 42703 column "escape" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT event
< error
| 42703 column "event" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT except
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT exclude
< error
| 42703 column "exclude" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT excluding
< error
| 42703 column "excluding" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT exclusive
< error
| 42703 column "exclusive" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT execute
< error
| 42703 column "execute" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT EXISTS (SELECT 1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT explain
< error
| 42703 column "explain" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT expression
< error
| 42703 column "expression" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT extension
< error
| 42703 column "extension" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT external
< error
| 42703 column "external" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT extract
< error
| 42703 column "extract" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT false
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT family
< error
| 42703 column "family" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT fetch
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT filter
< error
| 42703 column "filter" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT finalize
< error
| 42703 column "finalize" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT first
< error
| 42703 column "first" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT float
< error
| 42703 column "float" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT following
< error
| 42703 column "following" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT for
< error
| 42601 syntax error at end of input
| position: 11
> rollback
< ok
> begin
< ok
> exec
| SELECT force
< error
| 42703 column "force" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT foreign
< error
| 42601 syntax error at or near "foreign"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT format
< error
| 42703 column "format" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT forward
< error
| 42703 column "forward" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT freeze
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT from
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT full
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT function
< error
| 42703 column "function" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT functions
< error
| 42703 column "functions" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT generated
< error
| 42703 column "generated" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT global
< error
| 42703 column "global" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT grant
< error
| 42601 syntax error at or near "grant"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT granted
< error
| 42703 column "granted" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT GREATEST(1, 2)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT group
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT grouping
< error
| 42703 column "grouping" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT groups
< error
| 42703 column "groups" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT handler
< error
| 42703 column "handler" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT having
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT header
< error
| 42703 column "header" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT hold
< error
| 42703 column "hold" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT hour
< error
| 42703 column "hour" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT identity
< error
| 42703 column "identity" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT if
< error
| 42703 column "if" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT ilike
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT immediate
< error
| 42703 column "immediate" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT immutable
< error
| 42703 column "immutable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT implicit
< error
| 42703 column "implicit" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT import
< error
| 42703 column "import" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT in
< error
| 42601 syntax error at or near "in"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT include
< error
| 42703 column "include" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT including
< error
| 42703 column "including" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT increment
< error
| 42703 column "increment" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT indent
< error
| 42703 column "indent" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT index
< error
| 42703 column "index" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT indexes
< error
| 42703 column "indexes" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT inherit
< error
| 42703 column "inherit" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT inherits
< error
| 42703 column "inherits" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT initially
< error
| 42601 syntax error at or near "initially"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT inline
< error
| 42703 column "inline" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT inner
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT inout
< error
| 42703 column "inout" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT input
< error
| 42703 column "input" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT insensitive
< error
| 42703 column "insensitive" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT insert
< error
| 42703 column "insert" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT instead
< error
| 42703 column "instead" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT int
< error
| 42703 column "int" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT integer
< error
| 42703 column "integer" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT intersect
< error
| 42601 syntax error at end of input
| position: 17
> rollback
< ok
> begin
< ok
> exec
| SELECT INTERVAL '1 day'
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT into
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT invoker
< error
| 42703 column "invoker" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT is
< error
| 42601 syntax error at end of input
| position: 10
> rollback
< ok
> begin
< ok
> exec
| SELECT isnull
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT isolation
< error
| 42703 column "isolation" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT join
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT json
< error
| 42703 column "json" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_array
< error
| 42703 column "json_array" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_arrayagg
< error
| 42703 column "json_arrayagg" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_exists
< error
| 42703 column "json_exists" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_object
< error
| 42703 column "json_object" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_objectagg
< error
| 42703 column "json_objectagg" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_query
< error
| 42703 column "json_query" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_scalar
< error
| 42703 column "json_scalar" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_serialize
< error
| 42703 column "json_serialize" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_table
< error
| 42703 column "json_table" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT json_value
< error
| 42703 column "json_value" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT keep
< error
| 42703 column "keep" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT key
< error
| 42703 column "key" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT keys
< error
| 42703 column "keys" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT label
< error
| 42703 column "label" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT language
< error
| 42703 column "language" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT large
< error
| 42703 column "large" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT last
< error
| 42703 column "last" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT lateral
< error
| 42601 syntax error at or near "lateral"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT leading
< error
| 42601 syntax error at or near "leading"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT leakproof
< error
| 42703 column "leakproof" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT LEAST(1, 2)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT left
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT level
< error
| 42703 column "level" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT like
< error
| 42601 syntax error at end of input
| position: 12
> rollback
< ok
> begin
< ok
> exec
| SELECT limit
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT listen
< error
| 42703 column "listen" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT load
< error
| 42703 column "load" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT local
< error
| 42703 column "local" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT LOCALTIME
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT LOCALTIMESTAMP
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT location
< error
| 42703 column "location" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT lock
< error
| 42703 column "lock" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT locked
< error
| 42703 column "locked" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT logged
< error
| 42703 column "logged" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT mapping
< error
| 42703 column "mapping" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT match
< error
| 42703 column "match" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT matched
< error
| 42703 column "matched" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT materialized
< error
| 42703 column "materialized" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT maxvalue
< error
| 42703 column "maxvalue" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT merge
< error
| 42703 column "merge" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT merge_action
< error
| 42703 column "merge_action" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT method
< error
| 42703 column "method" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT minute
< error
| 42703 column "minute" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT minvalue
< error
| 42703 column "minvalue" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT mode
< error
| 42703 column "mode" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT month
< error
| 42703 column "month" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT move
< error
| 42703 column "move" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT name
< error
| 42703 column "name" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT names
< error
| 42703 column "names" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT national
< error
| 42703 column "national" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT natural
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT nchar
< error
| 42703 column "nchar" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT nested
< error
| 42703 column "nested" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT new
< error
| 42703 column "new" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT next
< error
| 42703 column "next" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT nfc
< error
| 42703 column "nfc" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT nfd
< error
| 42703 column "nfd" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT nfkc
< error
| 42703 column "nfkc" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT nfkd
< error
| 42703 column "nfkd" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT no
< error
| 42703 column "no" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT none
< error
| 42703 column "none" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT normalize
< error
| 42703 column "normalize" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT normalized
< error
| 42703 column "normalized" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT NOT TRUE
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT nothing
< error
| 42703 column "nothing" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT notify
< error
| 42703 column "notify" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT notnull
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT nowait
< error
| 42703 column "nowait" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT null
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT NULLIF(1, 2)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT nulls
< error
| 42703 column "nulls" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT numeric
< error
| 42703 column "numeric" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT object
< error
| 42703 column "object" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT of
< error
| 42703 column "of" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT off
< error
| 42703 column "off" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT offset
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT oids
< error
| 42703 column "oids" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT old
< error
| 42703 column "old" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT omit
< error
| 42703 column "omit" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT on
< error
| 42601 syntax error at or near "on"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT only
< error
| 42601 syntax error at or near "only"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT operator
< error
| 42703 column "operator" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT option
< error
| 42703 column "option" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT options
< error
| 42703 column "options" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT or
< error
| 42601 syntax error at or near "or"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT order
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT ordinality
< error
| 42703 column "ordinality" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT others
< error
| 42703 column "others" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT out
< error
| 42703 column "out" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT outer
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT over
< error
| 42703 column "over" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT overlaps
< error
| 42601 syntax error at end of input
| position: 16
> rollback
< ok
> begin
< ok
> exec
| SELECT overlay
< error
| 42703 column "overlay" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT overriding
< error
| 42703 column "overriding" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT owned
< error
| 42703 column "owned" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT owner
< error
| 42703 column "owner" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT parallel
< error
| 42703 column "parallel" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT parameter
< error
| 42703 column "parameter" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT parser
< error
| 42703 column "parser" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT partial
< error
| 42703 column "partial" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT partition
< error
| 42703 column "partition" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT passing
< error
| 42703 column "passing" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT password
< error
| 42703 column "password" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT path
< error
| 42703 column "path" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT placing
< error
| 42601 syntax error at or near "placing"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT plan
< error
| 42703 column "plan" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT plans
< error
| 42703 column "plans" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT policy
< error
| 42703 column "policy" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT position
< error
| 42703 column "position" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT preceding
< error
| 42703 column "preceding" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT precision
< error
| 42703 column "precision" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT prepare
< error
| 42703 column "prepare" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT prepared
< error
| 42703 column "prepared" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT preserve
< error
| 42703 column "preserve" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT primary
< error
| 42601 syntax error at or near "primary"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT prior
< error
| 42703 column "prior" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT privileges
< error
| 42703 column "privileges" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT procedural
< error
| 42703 column "procedural" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT procedure
< error
| 42703 column "procedure" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT procedures
< error
| 42703 column "procedures" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT program
< error
| 42703 column "program" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT publication
< error
| 42703 column "publication" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT quote
< error
| 42703 column "quote" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT quotes
< error
| 42703 column "quotes" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT range
< error
| 42703 column "range" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT read
< error
| 42703 column "read" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT real
< error
| 42703 column "real" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT reassign
< error
| 42703 column "reassign" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT recheck
< error
| 42703 column "recheck" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT recursive
< error
| 42703 column "recursive" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT ref
< error
| 42703 column "ref" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT references
< error
| 42601 syntax error at or near "references"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT referencing
< error
| 42703 column "referencing" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT refresh
< error
| 42703 column "refresh" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT reindex
< error
| 42703 column "reindex" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT relative
< error
| 42703 column "relative" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT release
< error
| 42703 column "release" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT rename
< error
| 42703 column "rename" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT repeatable
< error
| 42703 column "repeatable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT replace
< error
| 42703 column "replace" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT replica
< error
| 42703 column "replica" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT reset
< error
| 42703 column "reset" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT restart
< error
| 42703 column "restart" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT restrict
< error
| 42703 column "restrict" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT return
< error
| 42703 column "return" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT returning
< error
| 42601 syntax error at or near "returning"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT returns
< error
| 42703 column "returns" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT revoke
< error
| 42703 column "revoke" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT right
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT role
< error
| 42703 column "role" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT rollback
< error
| 42703 column "rollback" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT rollup
< error
| 42703 column "rollup" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT routine
< error
| 42703 column "routine" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT routines
< error
| 42703 column "routines" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT ROW(1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT rows
< error
| 42703 column "rows" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT rule
< error
| 42703 column "rule" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT savepoint
< error
| 42703 column "savepoint" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT scalar
< error
| 42703 column "scalar" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT schema
< error
| 42703 column "schema" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT schemas
< error
| 42703 column "schemas" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT scroll
< error
| 42703 column "scroll" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT search
< error
| 42703 column "search" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT second
< error
| 42703 column "second" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT security
< error
| 42703 column "security" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT select
< error
| 42601 syntax error at or near "select"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT sequence
< error
| 42703 column "sequence" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT sequences
< error
| 42703 column "sequences" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT serializable
< error
| 42703 column "serializable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT server
< error
| 42703 column "server" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT session
< error
| 42703 column "session" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT SESSION_USER
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT set
< error
| 42703 column "set" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT setof
< error
| 42703 column "setof" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT sets
< error
| 42703 column "sets" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT share
< error
| 42703 column "share" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT show
< error
| 42703 column "show" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT similar
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT simple
< error
| 42703 column "simple" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT skip
< error
| 42703 column "skip" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT smallint
< error
| 42703 column "smallint" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT snapshot
< error
| 42703 column "snapshot" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT some
< error
| 42601 syntax error at or near "some"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT source
< error
| 42703 column "source" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT sql
< error
| 42703 column "sql" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT stable
< error
| 42703 column "stable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT standalone
< error
| 42703 column "standalone" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT start
< error
| 42703 column "start" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT statement
< error
| 42703 column "statement" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT statistics
< error
| 42703 column "statistics" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT stdin
< error
| 42703 column "stdin" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT stdout
< error
| 42703 column "stdout" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT storage
< error
| 42703 column "storage" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT stored
< error
| 42703 column "stored" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT strict
< error
| 42703 column "strict" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT string
< error
| 42703 column "string" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT strip
< error
| 42703 column "strip" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT subscription
< error
| 42703 column "subscription" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT substring
< error
| 42703 column "substring" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT support
< error
| 42703 column "support" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT symmetric
< error
| 42601 syntax error at or near "symmetric"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT sysid
< error
| 42703 column "sysid" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT system
< error
| 42703 column "system" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT system_user
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT table
< error
| 42601 syntax error at or near "table"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT tables
< error
| 42703 column "tables" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT tablesample
< error
| 42601 syntax error at end of input
| position: 19
> rollback
< ok
> begin
< ok
> exec
| SELECT tablespace
< error
| 42703 column "tablespace" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT target
< error
| 42703 column "target" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT temp
< error
| 42703 column "temp" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT template
< error
| 42703 column "template" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT temporary
< error
| 42703 column "temporary" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT text
< error
| 42703 column "text" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT then
< error
| 42601 syntax error at or near "then"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT ties
< error
| 42703 column "ties" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT time
< error
| 42703 column "time" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT timestamp
< error
| 42703 column "timestamp" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT to
< error
| 42601 syntax error at or near "to"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT trailing
< error
| 42601 syntax error at or near "trailing"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT transaction
< error
| 42703 column "transaction" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT transform
< error
| 42703 column "transform" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT treat
< error
| 42703 column "treat" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT trigger
< error
| 42703 column "trigger" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT trim
< error
| 42703 column "trim" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT true
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT truncate
< error
| 42703 column "truncate" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT trusted
< error
| 42703 column "trusted" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT type
< error
| 42703 column "type" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT types
< error
| 42703 column "types" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT uescape
< error
| 42703 column "uescape" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT unbounded
< error
| 42703 column "unbounded" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT uncommitted
< error
| 42703 column "uncommitted" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT unconditional
< error
| 42703 column "unconditional" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT unencrypted
< error
| 42703 column "unencrypted" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT union
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT unique
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT unknown
< error
| 42703 column "unknown" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT unlisten
< error
| 42703 column "unlisten" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT unlogged
< error
| 42703 column "unlogged" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT until
< error
| 42703 column "until" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT update
< error
| 42703 column "update" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT USER
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT using
< error
| 42601 syntax error at or near "using"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT vacuum
< error
| 42703 column "vacuum" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT valid
< error
| 42703 column "valid" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT validate
< error
| 42703 column "validate" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT validator
< error
| 42703 column "validator" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT value
< error
| 42703 column "value" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT values
< error
| 42703 column "values" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT varchar
< error
| 42703 column "varchar" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT variadic
< error
| 42601 syntax error at or near "variadic"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT varying
< error
| 42703 column "varying" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT verbose
< error
| 42601 syntax error at end of input
| position: 15
> rollback
< ok
> begin
< ok
> exec
| SELECT version
< error
| 42703 column "version" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT view
< error
| 42703 column "view" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT views
< error
| 42703 column "views" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT volatile
< error
| 42703 column "volatile" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT when
< error
| 42601 syntax error at or near "when"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT where
< error
| 42601 syntax error at end of input
| position: 13
> rollback
< ok
> begin
< ok
> exec
| SELECT whitespace
< error
| 42703 column "whitespace" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT window
< error
| 42601 syntax error at end of input
| position: 14
> rollback
< ok
> begin
< ok
> exec
| SELECT with
< error
| 42601 syntax error at or near "with"
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT within
< error
| 42703 column "within" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT without
< error
| 42703 column "without" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT work
< error
| 42703 column "work" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT wrapper
< error
| 42703 column "wrapper" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT write
< error
| 42703 column "write" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xml
< error
| 42703 column "xml" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlattributes
< error
| 42703 column "xmlattributes" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlconcat
< error
| 42703 column "xmlconcat" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlelement
< error
| 42703 column "xmlelement" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlexists
< error
| 42703 column "xmlexists" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlforest
< error
| 42703 column "xmlforest" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlnamespaces
< error
| 42703 column "xmlnamespaces" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlparse
< error
| 42703 column "xmlparse" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlpi
< error
| 42703 column "xmlpi" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlroot
< error
| 42703 column "xmlroot" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmlserialize
< error
| 42703 column "xmlserialize" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT xmltable
< error
| 42703 column "xmltable" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT year
< error
| 42703 column "year" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT yes
< error
| 42703 column "yes" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT zone
< error
| 42703 column "zone" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT omni_probe_ident
< error
| 42703 column "omni_probe_ident" does not exist
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT 1
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 1.5
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT 'x'
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT B'1'
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT X'AB'
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT $1
< error
| 42P02 there is no parameter $1
| position: 8
> rollback
< ok
> begin
< ok
> exec
| SELECT (1)
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT +1
< ok 1 -
> rollback
< ok
> begin
< ok
> exec
| SELECT -1
< ok 1 -
> rollback
< ok

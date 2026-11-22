package parser

import "strings"

// TokenType identifies a lexical token from Apache Hive 4.0.1's HiveLexer.
type TokenType int

const (
	TokEOF TokenType = iota
	TokInvalid

	// Literals and identifiers.
	TokIdentifier
	TokStringLiteral
	TokCharSetLiteral
	TokCharSetName
	TokIntegralLiteral
	TokNumberLiteral
	TokByteLengthLiteral
	TokNumber
	TokQueryHint

	// Punctuation and operators.
	TokDot                  // .
	TokColon                // :
	TokComma                // ,
	TokSemicolon            // ;
	TokLParen               // (
	TokRParen               // )
	TokLSquare              // [
	TokRSquare              // ]
	TokLCurly               // {
	TokRCurly               // }
	TokEqual                // = or ==
	TokEqualNS              // <=>
	TokNotEqual             // <> or !=
	TokLessThanOrEqualTo    // <=
	TokLessThan             // <
	TokGreaterThanOrEqualTo // >=
	TokGreaterThan          // >
	TokDivide               // /
	TokPlus                 // +
	TokMinus                // -
	TokStar                 // *
	TokMod                  // %
	TokDiv                  // DIV
	TokAmpersand            // &
	TokTilde                // ~
	TokBitwiseOr            // |
	TokConcatenate          // ||
	TokBitwiseXor           // ^
	TokQuestion             // ?
	TokDollar               // $

	// Keywords (alphabetical / matching HiveLexerParent.g).
	KW_ABORT
	KW_ACTIVATE
	KW_ACTIVE
	KW_ADD
	KW_ADMIN
	KW_AFTER
	KW_ALL
	KW_ALLOC_FRACTION
	KW_ALTER
	KW_ANALYZE
	KW_AND
	KW_ANY
	KW_APPLICATION
	KW_ARCHIVE
	KW_ARRAY
	KW_AS
	KW_ASC
	KW_AST
	KW_AT
	KW_AUTHORIZATION
	KW_AUTOCOMMIT
	KW_BEFORE
	KW_BETWEEN
	KW_BIGINT
	KW_BINARY
	KW_BOOLEAN
	KW_BOTH
	KW_BRANCH
	KW_BUCKET
	KW_BUCKETS
	KW_BY
	KW_CACHE
	KW_CASCADE
	KW_CASE
	KW_CAST
	KW_CBO
	KW_CHANGE
	KW_CHAR
	KW_CHECK
	KW_CHERRY_PICK
	KW_CLUSTER
	KW_CLUSTERED
	KW_CLUSTERSTATUS
	KW_COLLECTION
	KW_COLUMN
	KW_COLUMNS
	KW_COMMENT
	KW_COMMIT
	KW_COMPACT
	KW_COMPACT_ID
	KW_COMPACTIONS
	KW_COMPUTE
	KW_CONCATENATE
	KW_CONF
	KW_CONSTRAINT
	KW_CONTINUE
	KW_CONVERT
	KW_COST
	KW_CREATE
	KW_CRON
	KW_CROSS
	KW_CUBE
	KW_CURRENT
	KW_CURRENT_DATE
	KW_CURRENT_TIMESTAMP
	KW_CURSOR
	KW_DATA
	KW_DATABASE
	KW_DATABASES
	KW_DATACONNECTOR
	KW_DATACONNECTORS
	KW_DATE
	KW_DATETIME
	KW_DAY
	KW_DBPROPERTIES
	KW_DCPROPERTIES
	KW_DDL
	KW_DEBUG
	KW_DECIMAL
	KW_DEFAULT
	KW_DEFERRED
	KW_DEFINED
	KW_DELETE
	KW_DELIMITED
	KW_DEPENDENCY
	KW_DESC
	KW_DESCRIBE
	KW_DETAIL
	KW_DIRECTORIES
	KW_DIRECTORY
	KW_DISABLE
	KW_DISTINCT
	KW_DISTRIBUTE
	KW_DISTRIBUTED
	KW_DO
	KW_DOUBLE
	KW_DOW
	KW_DROP
	KW_DUMP
	KW_ELEM_TYPE
	KW_ELSE
	KW_ENABLE
	KW_END
	KW_ENFORCED
	KW_ESCAPED
	KW_EVERY
	KW_EXCEPT
	KW_EXCHANGE
	KW_EXCLUSIVE
	KW_EXECUTE
	KW_EXECUTED
	KW_EXISTS
	KW_EXPIRE_SNAPSHOTS
	KW_EXPLAIN
	KW_EXPORT
	KW_EXPRESSION
	KW_EXTENDED
	KW_EXTERNAL
	KW_EXTRACT
	KW_FALSE
	KW_FAST_FORWARD
	KW_FETCH
	KW_FIELDS
	KW_FILE
	KW_FILEFORMAT
	KW_FIRST
	KW_FLOAT
	KW_FLOOR
	KW_FOLLOWING
	KW_FOR
	KW_FORCE
	KW_FOREIGN
	KW_FORMAT
	KW_FORMATTED
	KW_FROM
	KW_FULL
	KW_FUNCTION
	KW_FUNCTIONS
	KW_GRANT
	KW_GROUP
	KW_GROUPING
	KW_HAVING
	KW_HOUR
	KW_IDXPROPERTIES
	KW_IF
	KW_IGNORE
	KW_IMPORT
	KW_IN
	KW_INDEX
	KW_INDEXES
	KW_INNER
	KW_INPATH
	KW_INPUTDRIVER
	KW_INPUTFORMAT
	KW_INSERT
	KW_INT
	KW_INTERSECT
	KW_INTERVAL
	KW_INTO
	KW_IS
	KW_ISOLATION
	KW_ITEMS
	KW_JAR
	KW_JOIN
	KW_JOINCOST
	KW_KEY
	KW_KEY_TYPE
	KW_KEYS
	KW_KILL
	KW_LAST
	KW_LATERAL
	KW_LEADING
	KW_LEFT
	KW_LESS
	KW_LEVEL
	KW_LIKE
	KW_LIMIT
	KW_LINES
	KW_LOAD
	KW_LOCAL
	KW_LOCATION
	KW_LOCK
	KW_LOCKS
	KW_LOGICAL
	KW_LONG
	KW_MACRO
	KW_MANAGED
	KW_MANAGEDLOCATION
	KW_MANAGEMENT
	KW_MAP
	KW_MAPJOIN
	KW_MAPPING
	KW_MATCHED
	KW_MATERIALIZED
	KW_MERGE
	KW_METADATA
	KW_MINUS
	KW_MINUTE
	KW_MONTH
	KW_MORE
	KW_MOVE
	KW_MSCK
	KW_NONE
	KW_NORELY
	KW_NOSCAN
	KW_NOT
	KW_NOVALIDATE
	KW_NULL
	KW_NULLS
	KW_OF
	KW_OFFSET
	KW_OLDER
	KW_ON
	KW_ONLY
	KW_OPERATOR
	KW_OPTIMIZE
	KW_OPTION
	KW_OR
	KW_ORDER
	KW_ORPHAN_FILES
	KW_OUT
	KW_OUTER
	KW_OUTPUTDRIVER
	KW_OUTPUTFORMAT
	KW_OVER
	KW_OVERWRITE
	KW_OWNER
	KW_PARTITION
	KW_PARTITIONED
	KW_PARTITIONS
	KW_PATH
	KW_PERCENT
	KW_PKFK_JOIN
	KW_PLAN
	KW_PLANS
	KW_PLUS
	KW_POOL
	KW_PRECEDING
	KW_PRECISION
	KW_PREPARE
	KW_PRESERVE
	KW_PRIMARY
	KW_PRINCIPALS
	KW_PROCEDURE
	KW_PURGE
	KW_QUALIFY
	KW_QUARTER
	KW_QUERY
	KW_QUERY_PARALLELISM
	KW_RANGE
	KW_READ
	KW_READS
	KW_REAL
	KW_REBUILD
	KW_RECORDREADER
	KW_RECORDWRITER
	KW_REDUCE
	KW_REFERENCES
	KW_REGEXP
	KW_RELOAD
	KW_RELY
	KW_REMOTE
	KW_RENAME
	KW_REOPTIMIZATION
	KW_REPAIR
	KW_REPL
	KW_REPLACE
	KW_REPLICATION
	KW_RESOURCE
	KW_RESPECT
	KW_RESTRICT
	KW_RETAIN
	KW_RETENTION
	KW_REVOKE
	KW_REWRITE
	KW_RIGHT
	KW_RLIKE
	KW_ROLE
	KW_ROLES
	KW_ROLLBACK
	KW_ROLLUP
	KW_ROW
	KW_ROWS
	KW_SCHEDULED
	KW_SCHEDULING_POLICY
	KW_SCHEMA
	KW_SCHEMAS
	KW_SECOND
	KW_SELECT
	KW_SEMI
	KW_SERDE
	KW_SERDEPROPERTIES
	KW_SERVER
	KW_SET
	KW_SET_CURRENT_SNAPSHOT
	KW_SETS
	KW_SHARED
	KW_SHOW
	KW_SHOW_DATABASE
	KW_SKEWED
	KW_SMALLINT
	KW_SNAPSHOT
	KW_SNAPSHOTS
	KW_SOME
	KW_SORT
	KW_SORTED
	KW_SPEC
	KW_SSL
	KW_START
	KW_STATISTICS
	KW_STATUS
	KW_STORED
	KW_STREAMTABLE
	KW_STRING
	KW_STRUCT
	KW_SUMMARY
	KW_SYNC
	KW_SYSTEM_TIME
	KW_SYSTEM_VERSION
	KW_TABLE
	KW_TABLES
	KW_TABLESAMPLE
	KW_TAG
	KW_TBLPROPERTIES
	KW_TEMPORARY
	KW_TERMINATED
	KW_THAN
	KW_THEN
	KW_TIME
	KW_TIMESTAMP
	KW_TIMESTAMPLOCALTZ
	KW_TINYINT
	KW_TO
	KW_TOUCH
	KW_TRAILING
	KW_TRANSACTION
	KW_TRANSACTIONAL
	KW_TRANSACTIONS
	KW_TRANSFORM
	KW_TRIGGER
	KW_TRIM
	KW_TRUE
	KW_TRUNCATE
	KW_TYPE
	KW_UNARCHIVE
	KW_UNBOUNDED
	KW_UNDO
	KW_UNION
	KW_UNIONTYPE
	KW_UNIQUE
	KW_UNIQUEJOIN
	KW_UNKNOWN
	KW_UNLOCK
	KW_UNMANAGED
	KW_UNSET
	KW_UNSIGNED
	KW_UPDATE
	KW_URI
	KW_URL
	KW_USE
	KW_USER
	KW_USING
	KW_UTC
	KW_UTCTIMESTAMP
	KW_VALIDATE
	KW_VALUE_TYPE
	KW_VALUES
	KW_VARCHAR
	KW_VECTORIZATION
	KW_VIEW
	KW_VIEWS
	KW_WAIT
	KW_WEEK
	KW_WHEN
	KW_WHERE
	KW_WHILE
	KW_WINDOW
	KW_WITH
	KW_WITHIN
	KW_WORK
	KW_WORKLOAD
	KW_WRITE
	KW_YEAR
	KW_ZONE
)

// Loc is a byte range [Start, End) in the input SQL string.
type Loc struct {
	Start int
	End   int
}

// Token is one lexical token produced by the Hive lexer.
type Token struct {
	Type   TokenType
	Str    string // raw or unescaped text
	Quoted bool   // true for backtick-quoted Identifier
	Loc    Loc
}

var keywords = map[string]TokenType{
	"ABORT":                KW_ABORT,
	"ACTIVATE":             KW_ACTIVATE,
	"ACTIVE":               KW_ACTIVE,
	"ADD":                  KW_ADD,
	"ADMIN":                KW_ADMIN,
	"AFTER":                KW_AFTER,
	"ALL":                  KW_ALL,
	"ALLOC_FRACTION":       KW_ALLOC_FRACTION,
	"ALTER":                KW_ALTER,
	"ANALYZE":              KW_ANALYZE,
	"AND":                  KW_AND,
	"ANY":                  KW_ANY,
	"APPLICATION":          KW_APPLICATION,
	"ARCHIVE":              KW_ARCHIVE,
	"ARRAY":                KW_ARRAY,
	"AS":                   KW_AS,
	"ASC":                  KW_ASC,
	"AST":                  KW_AST,
	"AT":                   KW_AT,
	"AUTHORIZATION":        KW_AUTHORIZATION,
	"AUTOCOMMIT":           KW_AUTOCOMMIT,
	"BEFORE":               KW_BEFORE,
	"BETWEEN":              KW_BETWEEN,
	"BIGINT":               KW_BIGINT,
	"BINARY":               KW_BINARY,
	"BOOLEAN":              KW_BOOLEAN,
	"BOTH":                 KW_BOTH,
	"BRANCH":               KW_BRANCH,
	"BUCKET":               KW_BUCKET,
	"BUCKETS":              KW_BUCKETS,
	"BY":                   KW_BY,
	"CACHE":                KW_CACHE,
	"CASCADE":              KW_CASCADE,
	"CASE":                 KW_CASE,
	"CAST":                 KW_CAST,
	"CBO":                  KW_CBO,
	"CHANGE":               KW_CHANGE,
	"CHAR":                 KW_CHAR,
	"CHECK":                KW_CHECK,
	"CHERRY-PICK":          KW_CHERRY_PICK,
	"CLUSTER":              KW_CLUSTER,
	"CLUSTERED":            KW_CLUSTERED,
	"CLUSTERSTATUS":        KW_CLUSTERSTATUS,
	"COLLECTION":           KW_COLLECTION,
	"COLUMN":               KW_COLUMN,
	"COLUMNS":              KW_COLUMNS,
	"COMMENT":              KW_COMMENT,
	"COMMIT":               KW_COMMIT,
	"COMPACT":              KW_COMPACT,
	"COMPACTIONID":         KW_COMPACT_ID,
	"COMPACTIONS":          KW_COMPACTIONS,
	"COMPUTE":              KW_COMPUTE,
	"CONCATENATE":          KW_CONCATENATE,
	"CONF":                 KW_CONF,
	"CONSTRAINT":           KW_CONSTRAINT,
	"CONTINUE":             KW_CONTINUE,
	"CONVERT":              KW_CONVERT,
	"COST":                 KW_COST,
	"CREATE":               KW_CREATE,
	"CRON":                 KW_CRON,
	"CROSS":                KW_CROSS,
	"CUBE":                 KW_CUBE,
	"CURRENT":              KW_CURRENT,
	"CURRENT_DATE":         KW_CURRENT_DATE,
	"CURRENT_TIMESTAMP":    KW_CURRENT_TIMESTAMP,
	"CURSOR":               KW_CURSOR,
	"DATA":                 KW_DATA,
	"DATABASE":             KW_DATABASE,
	"DATABASES":            KW_DATABASES,
	"CONNECTOR":            KW_DATACONNECTOR,
	"CONNECTORS":           KW_DATACONNECTORS,
	"DATE":                 KW_DATE,
	"DATETIME":             KW_DATETIME,
	"DAY":                  KW_DAY,
	"DAYS":                 KW_DAY,
	"DAYOFWEEK":            KW_DOW,
	"DBPROPERTIES":         KW_DBPROPERTIES,
	"DCPROPERTIES":         KW_DCPROPERTIES,
	"DDL":                  KW_DDL,
	"DEBUG":                KW_DEBUG,
	"DEC":                  KW_DECIMAL,
	"DECIMAL":              KW_DECIMAL,
	"NUMERIC":              KW_DECIMAL,
	"DEFAULT":              KW_DEFAULT,
	"DEFERRED":             KW_DEFERRED,
	"DEFINED":              KW_DEFINED,
	"DELETE":               KW_DELETE,
	"DELIMITED":            KW_DELIMITED,
	"DEPENDENCY":           KW_DEPENDENCY,
	"DESC":                 KW_DESC,
	"DESCRIBE":             KW_DESCRIBE,
	"DETAIL":               KW_DETAIL,
	"DIRECTORIES":          KW_DIRECTORIES,
	"DIRECTORY":            KW_DIRECTORY,
	"DISABLE":              KW_DISABLE,
	"DISABLED":             KW_DISABLE,
	"DISTINCT":             KW_DISTINCT,
	"DISTRIBUTE":           KW_DISTRIBUTE,
	"DISTRIBUTED":          KW_DISTRIBUTED,
	"DIV":                  TokDiv,
	"DO":                   KW_DO,
	"DOUBLE":               KW_DOUBLE,
	"DROP":                 KW_DROP,
	"DUMP":                 KW_DUMP,
	"$ELEM$":               KW_ELEM_TYPE,
	"ELSE":                 KW_ELSE,
	"ENABLE":               KW_ENABLE,
	"ENABLED":              KW_ENABLE,
	"END":                  KW_END,
	"ENFORCED":             KW_ENFORCED,
	"ESCAPED":              KW_ESCAPED,
	"EVERY":                KW_EVERY,
	"EXCEPT":               KW_EXCEPT,
	"EXCHANGE":             KW_EXCHANGE,
	"EXCLUSIVE":            KW_EXCLUSIVE,
	"EXECUTE":              KW_EXECUTE,
	"EXECUTED":             KW_EXECUTED,
	"EXISTS":               KW_EXISTS,
	"EXPIRE_SNAPSHOTS":     KW_EXPIRE_SNAPSHOTS,
	"EXPLAIN":              KW_EXPLAIN,
	"EXPORT":               KW_EXPORT,
	"EXPRESSION":           KW_EXPRESSION,
	"EXTENDED":             KW_EXTENDED,
	"EXTERNAL":             KW_EXTERNAL,
	"EXTRACT":              KW_EXTRACT,
	"FALSE":                KW_FALSE,
	"FAST-FORWARD":         KW_FAST_FORWARD,
	"FETCH":                KW_FETCH,
	"FIELDS":               KW_FIELDS,
	"FILE":                 KW_FILE,
	"FILEFORMAT":           KW_FILEFORMAT,
	"FIRST":                KW_FIRST,
	"FLOAT":                KW_FLOAT,
	"FLOOR":                KW_FLOOR,
	"FOLLOWING":            KW_FOLLOWING,
	"FOR":                  KW_FOR,
	"FORCE":                KW_FORCE,
	"FOREIGN":              KW_FOREIGN,
	"FORMAT":               KW_FORMAT,
	"FORMATTED":            KW_FORMATTED,
	"FROM":                 KW_FROM,
	"FULL":                 KW_FULL,
	"FUNCTION":             KW_FUNCTION,
	"FUNCTIONS":            KW_FUNCTIONS,
	"GRANT":                KW_GRANT,
	"GROUP":                KW_GROUP,
	"GROUPING":             KW_GROUPING,
	"HAVING":               KW_HAVING,
	"HOUR":                 KW_HOUR,
	"HOURS":                KW_HOUR,
	"IDXPROPERTIES":        KW_IDXPROPERTIES,
	"IF":                   KW_IF,
	"IGNORE":               KW_IGNORE,
	"IMPORT":               KW_IMPORT,
	"IN":                   KW_IN,
	"INDEX":                KW_INDEX,
	"INDEXES":              KW_INDEXES,
	"INNER":                KW_INNER,
	"INPATH":               KW_INPATH,
	"INPUTDRIVER":          KW_INPUTDRIVER,
	"INPUTFORMAT":          KW_INPUTFORMAT,
	"INSERT":               KW_INSERT,
	"INT":                  KW_INT,
	"INTEGER":              KW_INT,
	"INTERSECT":            KW_INTERSECT,
	"INTERVAL":             KW_INTERVAL,
	"INTO":                 KW_INTO,
	"IS":                   KW_IS,
	"ISOLATION":            KW_ISOLATION,
	"ITEMS":                KW_ITEMS,
	"JAR":                  KW_JAR,
	"JOIN":                 KW_JOIN,
	"JOINCOST":             KW_JOINCOST,
	"KEY":                  KW_KEY,
	"$KEY$":                KW_KEY_TYPE,
	"KEYS":                 KW_KEYS,
	"KILL":                 KW_KILL,
	"LAST":                 KW_LAST,
	"LATERAL":              KW_LATERAL,
	"LEADING":              KW_LEADING,
	"LEFT":                 KW_LEFT,
	"LESS":                 KW_LESS,
	"LEVEL":                KW_LEVEL,
	"LIKE":                 KW_LIKE,
	"LIMIT":                KW_LIMIT,
	"LINES":                KW_LINES,
	"LOAD":                 KW_LOAD,
	"LOCAL":                KW_LOCAL,
	"LOCATION":             KW_LOCATION,
	"LOCK":                 KW_LOCK,
	"LOCKS":                KW_LOCKS,
	"LOGICAL":              KW_LOGICAL,
	"LONG":                 KW_LONG,
	"MACRO":                KW_MACRO,
	"MANAGED":              KW_MANAGED,
	"MANAGEDLOCATION":      KW_MANAGEDLOCATION,
	"MANAGEMENT":           KW_MANAGEMENT,
	"MAP":                  KW_MAP,
	"MAPJOIN":              KW_MAPJOIN,
	"MAPPING":              KW_MAPPING,
	"MATCHED":              KW_MATCHED,
	"MATERIALIZED":         KW_MATERIALIZED,
	"MERGE":                KW_MERGE,
	"METADATA":             KW_METADATA,
	"MINUS":                KW_MINUS,
	"MINUTE":               KW_MINUTE,
	"MINUTES":              KW_MINUTE,
	"MONTH":                KW_MONTH,
	"MONTHS":               KW_MONTH,
	"MORE":                 KW_MORE,
	"MOVE":                 KW_MOVE,
	"MSCK":                 KW_MSCK,
	"NONE":                 KW_NONE,
	"NORELY":               KW_NORELY,
	"NOSCAN":               KW_NOSCAN,
	"NOT":                  KW_NOT,
	"NOVALIDATE":           KW_NOVALIDATE,
	"NULL":                 KW_NULL,
	"NULLS":                KW_NULLS,
	"OF":                   KW_OF,
	"OFFSET":               KW_OFFSET,
	"OLDER":                KW_OLDER,
	"ON":                   KW_ON,
	"ONLY":                 KW_ONLY,
	"OPERATOR":             KW_OPERATOR,
	"OPTIMIZE":             KW_OPTIMIZE,
	"OPTION":               KW_OPTION,
	"OR":                   KW_OR,
	"ORDER":                KW_ORDER,
	"ORPHAN-FILES":         KW_ORPHAN_FILES,
	"OUT":                  KW_OUT,
	"OUTER":                KW_OUTER,
	"OUTPUTDRIVER":         KW_OUTPUTDRIVER,
	"OUTPUTFORMAT":         KW_OUTPUTFORMAT,
	"OVER":                 KW_OVER,
	"OVERWRITE":            KW_OVERWRITE,
	"OWNER":                KW_OWNER,
	"PARTITION":            KW_PARTITION,
	"PARTITIONED":          KW_PARTITIONED,
	"PARTITIONS":           KW_PARTITIONS,
	"PATH":                 KW_PATH,
	"PERCENT":              KW_PERCENT,
	"PKFK_JOIN":            KW_PKFK_JOIN,
	"PLAN":                 KW_PLAN,
	"PLANS":                KW_PLANS,
	"PLUS":                 KW_PLUS,
	"POOL":                 KW_POOL,
	"PRECEDING":            KW_PRECEDING,
	"PRECISION":            KW_PRECISION,
	"PREPARE":              KW_PREPARE,
	"PRESERVE":             KW_PRESERVE,
	"PRIMARY":              KW_PRIMARY,
	"PRINCIPALS":           KW_PRINCIPALS,
	"PROCEDURE":            KW_PROCEDURE,
	"PURGE":                KW_PURGE,
	"QUALIFY":              KW_QUALIFY,
	"QUARTER":              KW_QUARTER,
	"QUERY":                KW_QUERY,
	"QUERY_PARALLELISM":    KW_QUERY_PARALLELISM,
	"RANGE":                KW_RANGE,
	"READ":                 KW_READ,
	"READS":                KW_READS,
	"REAL":                 KW_REAL,
	"REBUILD":              KW_REBUILD,
	"RECORDREADER":         KW_RECORDREADER,
	"RECORDWRITER":         KW_RECORDWRITER,
	"REDUCE":               KW_REDUCE,
	"REFERENCES":           KW_REFERENCES,
	"REGEXP":               KW_REGEXP,
	"RELOAD":               KW_RELOAD,
	"RELY":                 KW_RELY,
	"REMOTE":               KW_REMOTE,
	"RENAME":               KW_RENAME,
	"REOPTIMIZATION":       KW_REOPTIMIZATION,
	"REPAIR":               KW_REPAIR,
	"REPL":                 KW_REPL,
	"REPLACE":              KW_REPLACE,
	"REPLICATION":          KW_REPLICATION,
	"RESOURCE":             KW_RESOURCE,
	"RESPECT":              KW_RESPECT,
	"RESTRICT":             KW_RESTRICT,
	"RETAIN":               KW_RETAIN,
	"RETENTION":            KW_RETENTION,
	"REVOKE":               KW_REVOKE,
	"REWRITE":              KW_REWRITE,
	"RIGHT":                KW_RIGHT,
	"RLIKE":                KW_RLIKE,
	"ROLE":                 KW_ROLE,
	"ROLES":                KW_ROLES,
	"ROLLBACK":             KW_ROLLBACK,
	"ROLLUP":               KW_ROLLUP,
	"ROW":                  KW_ROW,
	"ROWS":                 KW_ROWS,
	"SCHEDULED":            KW_SCHEDULED,
	"SCHEDULING_POLICY":    KW_SCHEDULING_POLICY,
	"SCHEMA":               KW_SCHEMA,
	"SCHEMAS":              KW_SCHEMAS,
	"SECOND":               KW_SECOND,
	"SECONDS":              KW_SECOND,
	"SELECT":               KW_SELECT,
	"SEMI":                 KW_SEMI,
	"SERDE":                KW_SERDE,
	"SERDEPROPERTIES":      KW_SERDEPROPERTIES,
	"SERVER":               KW_SERVER,
	"SET":                  KW_SET,
	"SET_CURRENT_SNAPSHOT": KW_SET_CURRENT_SNAPSHOT,
	"SETS":                 KW_SETS,
	"SHARED":               KW_SHARED,
	"SHOW":                 KW_SHOW,
	"SHOW_DATABASE":        KW_SHOW_DATABASE,
	"SKEWED":               KW_SKEWED,
	"SMALLINT":             KW_SMALLINT,
	"SNAPSHOT":             KW_SNAPSHOT,
	"SNAPSHOTS":            KW_SNAPSHOTS,
	"SOME":                 KW_SOME,
	"SORT":                 KW_SORT,
	"SORTED":               KW_SORTED,
	"SPEC":                 KW_SPEC,
	"SSL":                  KW_SSL,
	"START":                KW_START,
	"STATISTICS":           KW_STATISTICS,
	"STATUS":               KW_STATUS,
	"STORED":               KW_STORED,
	"STREAMTABLE":          KW_STREAMTABLE,
	"STRING":               KW_STRING,
	"STRUCT":               KW_STRUCT,
	"SUMMARY":              KW_SUMMARY,
	"SYNC":                 KW_SYNC,
	"SYSTEM_TIME":          KW_SYSTEM_TIME,
	"SYSTEM_VERSION":       KW_SYSTEM_VERSION,
	"TABLE":                KW_TABLE,
	"TABLES":               KW_TABLES,
	"TABLESAMPLE":          KW_TABLESAMPLE,
	"TAG":                  KW_TAG,
	"TBLPROPERTIES":        KW_TBLPROPERTIES,
	"TEMPORARY":            KW_TEMPORARY,
	"TERMINATED":           KW_TERMINATED,
	"THAN":                 KW_THAN,
	"THEN":                 KW_THEN,
	"TIME":                 KW_TIME,
	"TIMESTAMP":            KW_TIMESTAMP,
	"TIMESTAMPLOCALTZ":     KW_TIMESTAMPLOCALTZ,
	"TINYINT":              KW_TINYINT,
	"TO":                   KW_TO,
	"TOUCH":                KW_TOUCH,
	"TRAILING":             KW_TRAILING,
	"TRANSACTION":          KW_TRANSACTION,
	"TRANSACTIONAL":        KW_TRANSACTIONAL,
	"TRANSACTIONS":         KW_TRANSACTIONS,
	"TRANSFORM":            KW_TRANSFORM,
	"TRIGGER":              KW_TRIGGER,
	"TRIM":                 KW_TRIM,
	"TRUE":                 KW_TRUE,
	"TRUNCATE":             KW_TRUNCATE,
	"TYPE":                 KW_TYPE,
	"UNARCHIVE":            KW_UNARCHIVE,
	"UNBOUNDED":            KW_UNBOUNDED,
	"UNDO":                 KW_UNDO,
	"UNION":                KW_UNION,
	"UNIONTYPE":            KW_UNIONTYPE,
	"UNIQUE":               KW_UNIQUE,
	"UNIQUEJOIN":           KW_UNIQUEJOIN,
	"UNKNOWN":              KW_UNKNOWN,
	"UNLOCK":               KW_UNLOCK,
	"UNMANAGED":            KW_UNMANAGED,
	"UNSET":                KW_UNSET,
	"UNSIGNED":             KW_UNSIGNED,
	"UPDATE":               KW_UPDATE,
	"URI":                  KW_URI,
	"URL":                  KW_URL,
	"USE":                  KW_USE,
	"USER":                 KW_USER,
	"USING":                KW_USING,
	"UTC":                  KW_UTC,
	"UTC_TMESTAMP":         KW_UTCTIMESTAMP,
	"VALIDATE":             KW_VALIDATE,
	"$VALUE$":              KW_VALUE_TYPE,
	"VALUES":               KW_VALUES,
	"VARCHAR":              KW_VARCHAR,
	"VECTORIZATION":        KW_VECTORIZATION,
	"VIEW":                 KW_VIEW,
	"VIEWS":                KW_VIEWS,
	"WAIT":                 KW_WAIT,
	"WEEK":                 KW_WEEK,
	"WEEKS":                KW_WEEK,
	"WHEN":                 KW_WHEN,
	"WHERE":                KW_WHERE,
	"WHILE":                KW_WHILE,
	"WINDOW":               KW_WINDOW,
	"WITH":                 KW_WITH,
	"WITHIN":               KW_WITHIN,
	"WORK":                 KW_WORK,
	"WORKLOAD":             KW_WORKLOAD,
	"WRITE":                KW_WRITE,
	"YEAR":                 KW_YEAR,
	"YEARS":                KW_YEAR,
	"ZONE":                 KW_ZONE,
}

func lookupKeyword(s string) (TokenType, bool) {
	t, ok := keywords[strings.ToUpper(s)]
	return t, ok
}

// nonReserved holds all keywords listed in IdentifiersParser.g's nonReserved rule.
var nonReserved = map[TokenType]bool{
	KW_ABORT: true, KW_ADD: true, KW_ADMIN: true, KW_AFTER: true, KW_ANALYZE: true,
	KW_ARCHIVE: true, KW_ASC: true, KW_BEFORE: true, KW_BUCKET: true, KW_BUCKETS: true,
	KW_CASCADE: true, KW_CBO: true, KW_CHANGE: true, KW_CHECK: true, KW_CLUSTER: true,
	KW_CLUSTERED: true, KW_CLUSTERSTATUS: true, KW_COLLECTION: true, KW_COLUMNS: true,
	KW_COMMENT: true, KW_COMPACT: true, KW_COMPACTIONS: true, KW_COMPUTE: true,
	KW_CONCATENATE: true, KW_CONTINUE: true, KW_COST: true, KW_DATA: true, KW_DAY: true,
	KW_DATABASES: true, KW_DATETIME: true, KW_DBPROPERTIES: true, KW_DCPROPERTIES: true,
	KW_DEFERRED: true, KW_DEFINED: true, KW_DELIMITED: true, KW_DEPENDENCY: true,
	KW_DESC: true, KW_DIRECTORIES: true, KW_DIRECTORY: true, KW_DISABLE: true,
	KW_DISTRIBUTE: true, KW_DISTRIBUTED: true, KW_DOW: true, KW_ELEM_TYPE: true,
	KW_ENABLE: true, KW_ENFORCED: true, KW_ESCAPED: true, KW_EXCLUSIVE: true,
	KW_EXPLAIN: true, KW_EXPORT: true, KW_FIELDS: true, KW_FILE: true, KW_FILEFORMAT: true,
	KW_FIRST: true, KW_FORMAT: true, KW_FORMATTED: true, KW_FUNCTIONS: true, KW_HOUR: true,
	KW_IDXPROPERTIES: true, KW_RESPECT: true, KW_IGNORE: true, KW_INDEX: true,
	KW_INDEXES: true, KW_INPATH: true, KW_INPUTDRIVER: true, KW_INPUTFORMAT: true,
	KW_ITEMS: true, KW_JAR: true, KW_JOINCOST: true, KW_KILL: true, KW_KEYS: true,
	KW_KEY_TYPE: true, KW_LAST: true, KW_LIMIT: true, KW_OFFSET: true, KW_LINES: true,
	KW_LOAD: true, KW_LOCATION: true, KW_LOCK: true, KW_LOCKS: true, KW_LOGICAL: true,
	KW_LONG: true, KW_MANAGED: true, KW_MANAGEDLOCATION: true, KW_MAPJOIN: true,
	KW_MATERIALIZED: true, KW_METADATA: true, KW_MINUTE: true, KW_MONTH: true, KW_MSCK: true,
	KW_NOSCAN: true, KW_NULLS: true, KW_OPTION: true, KW_OUTPUTDRIVER: true,
	KW_OUTPUTFORMAT: true, KW_OVERWRITE: true, KW_OWNER: true, KW_PARTITIONED: true,
	KW_PARTITIONS: true, KW_PLUS: true, KW_PRINCIPALS: true, KW_PURGE: true, KW_QUERY: true,
	KW_QUARTER: true, KW_READ: true, KW_REBUILD: true, KW_RECORDREADER: true,
	KW_RECORDWRITER: true, KW_RELOAD: true, KW_REMOTE: true, KW_RENAME: true, KW_REPAIR: true,
	KW_REPLACE: true, KW_REPLICATION: true, KW_RESTRICT: true, KW_REWRITE: true,
	KW_ROLE: true, KW_ROLES: true, KW_SCHEMA: true, KW_SCHEMAS: true, KW_SECOND: true,
	KW_SEMI: true, KW_SERDE: true, KW_SERDEPROPERTIES: true, KW_SERVER: true, KW_SETS: true,
	KW_SHARED: true, KW_SHOW: true, KW_SHOW_DATABASE: true, KW_SKEWED: true, KW_SORT: true,
	KW_SORTED: true, KW_SSL: true, KW_STATISTICS: true, KW_STORED: true, KW_AST: true,
	KW_STREAMTABLE: true, KW_STRING: true, KW_STRUCT: true, KW_TABLES: true,
	KW_TBLPROPERTIES: true, KW_TEMPORARY: true, KW_TERMINATED: true, KW_TINYINT: true,
	KW_TOUCH: true, KW_TRANSACTIONAL: true, KW_TRANSACTIONS: true, KW_TYPE: true,
	KW_UNARCHIVE: true, KW_UNDO: true, KW_UNIONTYPE: true, KW_UNLOCK: true, KW_UNSET: true,
	KW_UNSIGNED: true, KW_URI: true, KW_URL: true, KW_USE: true, KW_UTC: true,
	KW_UTCTIMESTAMP: true, KW_VALUE_TYPE: true, KW_VIEW: true, KW_WEEK: true, KW_WHILE: true,
	KW_YEAR: true, KW_WORK: true, KW_TRANSACTION: true, KW_WRITE: true, KW_ISOLATION: true,
	KW_LEVEL: true, KW_SNAPSHOT: true, KW_AUTOCOMMIT: true, KW_RELY: true, KW_NORELY: true,
	KW_VALIDATE: true, KW_NOVALIDATE: true, KW_KEY: true, KW_MATCHED: true, KW_REPL: true,
	KW_DUMP: true, KW_STATUS: true, KW_CACHE: true, KW_VIEWS: true, KW_VECTORIZATION: true,
	KW_SUMMARY: true, KW_OPERATOR: true, KW_EXPRESSION: true, KW_DETAIL: true, KW_DEBUG: true,
	KW_WAIT: true, KW_ZONE: true, KW_DEFAULT: true, KW_REOPTIMIZATION: true,
	KW_EXECUTED: true, KW_SCHEDULED: true, KW_CRON: true, KW_EVERY: true, KW_AT: true,
	KW_EXECUTE: true, KW_RESOURCE: true, KW_PLAN: true, KW_PLANS: true,
	KW_QUERY_PARALLELISM: true, KW_ACTIVATE: true, KW_MOVE: true, KW_DO: true, KW_POOL: true,
	KW_ALLOC_FRACTION: true, KW_SCHEDULING_POLICY: true, KW_PATH: true, KW_MAPPING: true,
	KW_WORKLOAD: true, KW_MANAGEMENT: true, KW_ACTIVE: true, KW_UNMANAGED: true,
	KW_UNKNOWN: true, KW_WITHIN: true, KW_TRIM: true, KW_SPEC: true, KW_SYSTEM_TIME: true,
	KW_SYSTEM_VERSION: true, KW_EXPIRE_SNAPSHOTS: true, KW_SET_CURRENT_SNAPSHOT: true,
	KW_BRANCH: true, KW_SNAPSHOTS: true, KW_RETAIN: true, KW_RETENTION: true, KW_TAG: true,
	KW_FAST_FORWARD: true, KW_OPTIMIZE: true,
}

// sql11ReservedFunctionNames matches IdentifiersParser.g's
// sql11ReservedKeywordsUsedAsFunctionName rule.
var sql11ReservedFunctionNames = map[TokenType]bool{
	KW_IF: true, KW_ARRAY: true, KW_MAP: true, KW_BIGINT: true, KW_BINARY: true,
	KW_BOOLEAN: true, KW_CURRENT_DATE: true, KW_CURRENT_TIMESTAMP: true, KW_DATE: true,
	KW_DOUBLE: true, KW_FLOAT: true, KW_REAL: true, KW_GROUPING: true, KW_INT: true,
	KW_SMALLINT: true, KW_TIMESTAMP: true,
}

// sysFuncNames matches IdentifiersParser.g's sysFuncNames rule.
var sysFuncNames = map[TokenType]bool{
	KW_AND: true, KW_OR: true, KW_NOT: true, KW_LIKE: true, KW_IF: true, KW_CASE: true,
	KW_WHEN: true, KW_FLOOR: true, KW_TINYINT: true, KW_SMALLINT: true, KW_INT: true,
	KW_BIGINT: true, KW_FLOAT: true, KW_REAL: true, KW_DOUBLE: true, KW_BOOLEAN: true,
	KW_STRING: true, KW_BINARY: true, KW_ARRAY: true, KW_MAP: true, KW_STRUCT: true,
	KW_UNIONTYPE: true, TokEqual: true, TokEqualNS: true, TokNotEqual: true,
	TokLessThanOrEqualTo: true, TokLessThan: true, TokGreaterThanOrEqualTo: true,
	TokGreaterThan: true, TokDivide: true, TokPlus: true, TokMinus: true, TokStar: true,
	TokMod: true, TokDiv: true, TokAmpersand: true, TokTilde: true, TokBitwiseOr: true,
	TokBitwiseXor: true, KW_RLIKE: true, KW_REGEXP: true, KW_IN: true, KW_BETWEEN: true,
}

// isIdentifier reports whether tt is Identifier or a nonReserved keyword.
func isIdentifier(tt TokenType) bool {
	return tt == TokIdentifier || nonReserved[tt]
}

// isFunctionNameStart reports whether tt can begin a functionName.
func isFunctionNameStart(tt TokenType) bool {
	return isIdentifier(tt) || sql11ReservedFunctionNames[tt]
}

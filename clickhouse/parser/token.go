package parser

import "strings"

// TokenType identifies the lexical kind of a ClickHouse SQL token.
type TokenType int

const (
	TokEOF TokenType = iota
	TokIllegal

	// Identifiers and literals
	TokIdent       // unquoted identifier
	TokQuotedIdent // `backtick` or "double-quoted" identifier
	TokInt         // decimal, octal, hex, or binary integer literal
	TokFloat       // floating-point literal
	TokString      // 'single-quoted' or $tag$...$tag$ string literal
	TokHexString   // X'...'
	TokBinString   // B'...'

	// Punctuation and operators
	TokLParen      // (
	TokRParen      // )
	TokLBracket    // [
	TokRBracket    // ]
	TokLBrace      // {
	TokRBrace      // }
	TokComma       // ,
	TokDot         // .
	TokSemicolon   // ;
	TokColon       // :
	TokDoubleColon // ::
	TokArrow       // ->
	TokQuestion    // ?
	TokPlus        // +
	TokMinus       // -
	TokStar        // *
	TokSlash       // /
	TokPercent     // %
	TokConcat      // ||
	TokEq          // =
	TokEqDouble    // ==
	TokNotEq       // != or <>
	TokLT          // <
	TokLE          // <=
	TokGT          // >
	TokGE          // >=

	// Keywords (alphabetical; all case-insensitive)
	kwStart
	KwAccess
	KwAdd
	KwAfter
	KwAlias
	KwAll
	KwAlter
	KwAnd
	KwAnti
	KwAny
	KwArray
	KwAs
	KwAsc
	KwAsof
	KwAst
	KwAsync
	KwAttach
	KwBackup
	KwBetween
	KwBoth
	KwBy
	KwCaches
	KwCase
	KwCast
	KwChanged
	KwCheck
	KwClear
	KwCluster
	KwClusters
	KwCodec
	KwCollate
	KwColumn
	KwColumns
	KwComment
	KwCompression
	KwConstraint
	KwCreate
	KwCross
	KwCube
	KwCurrent
	KwCurrentUser
	KwDatabase
	KwDatabases
	KwDate
	KwDay
	KwDeduplicate
	KwDefault
	KwDelay
	KwDelete
	KwDesc
	KwDescending
	KwDescribe
	KwDetach
	KwDictionaries
	KwDictionary
	KwDisk
	KwDistinct
	KwDistributed
	KwDrop
	KwElse
	KwEnabled
	KwEnd
	KwEngine
	KwEngines
	KwEstimate
	KwEvents
	KwExcept
	KwExchange
	KwExists
	KwExplain
	KwExpression
	KwExtended
	KwExtract
	KwFalse
	KwFetches
	KwFields
	KwFilesystem
	KwFill
	KwFinal
	KwFirst
	KwFlush
	KwFollowing
	KwFor
	KwFormat
	KwFreeze
	KwFrom
	KwFull
	KwFunction
	KwFunctions
	KwGlobal
	KwGrant
	KwGrants
	KwGranularity
	KwGroup
	KwGrouping
	KwHaving
	KwHierarchical
	KwHour
	KwId
	KwIdentified
	KwIf
	KwIlike
	KwImplicit
	KwIn
	KwIndex
	KwIndexes
	KwIndices
	KwInf
	KwInjective
	KwInner
	KwInsert
	KwInterpolate
	KwIntersect
	KwInterval
	KwInto
	KwIs
	KwIsObjectId
	KwJoin
	KwKey
	KwKeys
	KwKill
	KwLast
	KwLayout
	KwLeading
	KwLeft
	KwLifetime
	KwLike
	KwLimit
	KwLive
	KwLocal
	KwLogs
	KwMaterialize
	KwMaterialized
	KwMax
	KwMerges
	KwMicrosecond
	KwMillisecond
	KwMin
	KwMinute
	KwModify
	KwMonth
	KwMove
	KwMutation
	KwNan
	KwNanosecond
	KwNo
	KwNot
	KwNull
	KwNulls
	KwOffset
	KwOn
	KwOptimize
	KwOr
	KwOrder
	KwOuter
	KwOutfile
	KwOver
	KwOverride
	KwPartition
	KwPermanently
	KwPipeline
	KwPlan
	KwPolicies
	KwPolicy
	KwPopulate
	KwPreceding
	KwPrewhere
	KwPrimary
	KwPrivileges
	KwProcesslist
	KwProfile
	KwProfiles
	KwProjection
	KwQualify
	KwQuarter
	KwQuery
	KwQuota
	KwQuotas
	KwRange
	KwRecursive
	KwRefresh
	KwReload
	KwRemove
	KwRename
	KwReplace
	KwReplica
	KwReplicated
	KwRestore
	KwRevoke
	KwRight
	KwRole
	KwRoles
	KwRollup
	KwRow
	KwRows
	KwSample
	KwSecond
	KwSelect
	KwSemi
	KwSends
	KwSet
	KwSets
	KwSetting
	KwSettings
	KwShow
	KwSource
	KwStart
	KwStep
	KwStop
	KwSubstring
	KwSync
	KwSyntax
	KwSystem
	KwTable
	KwTables
	KwTemporary
	KwTest
	KwThen
	KwTies
	KwTimeout
	KwTimestamp
	KwTo
	KwTop
	KwTotals
	KwTrailing
	KwTree
	KwTrim
	KwTrue
	KwTruncate
	KwTtl
	KwType
	KwUnbounded
	KwUndrop
	KwUnion
	KwUpdate
	KwUse
	KwUser
	KwUsers
	KwUsing
	KwUuid
	KwValues
	KwView
	KwVolume
	KwWatch
	KwWeek
	KwWhen
	KwWhere
	KwWindow
	KwWith
	KwYear
	kwEnd
)

// Token represents one scanned ClickHouse SQL token.
type Token struct {
	Type TokenType
	Raw  string // exact source slice
	Val  string // decoded text (for strings/quoted idents) or raw text
	Pos  int    // start byte offset
	End  int    // end byte offset (exclusive)
}

// IsKeyword reports whether t is a keyword token.
func (t Token) IsKeyword() bool {
	return t.Type > kwStart && t.Type < kwEnd
}

var keywords = map[string]TokenType{
	"ACCESS":       KwAccess,
	"ADD":          KwAdd,
	"AFTER":        KwAfter,
	"ALIAS":        KwAlias,
	"ALL":          KwAll,
	"ALTER":        KwAlter,
	"AND":          KwAnd,
	"ANTI":         KwAnti,
	"ANY":          KwAny,
	"ARRAY":        KwArray,
	"AS":           KwAs,
	"ASC":          KwAsc,
	"ASCENDING":    KwAsc,
	"ASOF":         KwAsof,
	"AST":          KwAst,
	"ASYNC":        KwAsync,
	"ATTACH":       KwAttach,
	"BACKUP":       KwBackup,
	"BETWEEN":      KwBetween,
	"BOTH":         KwBoth,
	"BY":           KwBy,
	"CACHES":       KwCaches,
	"CASE":         KwCase,
	"CAST":         KwCast,
	"CHANGED":      KwChanged,
	"CHECK":        KwCheck,
	"CLEAR":        KwClear,
	"CLUSTER":      KwCluster,
	"CLUSTERS":     KwClusters,
	"CODEC":        KwCodec,
	"COLLATE":      KwCollate,
	"COLUMN":       KwColumn,
	"COLUMNS":      KwColumns,
	"COMMENT":      KwComment,
	"COMPRESSION":  KwCompression,
	"CONSTRAINT":   KwConstraint,
	"CREATE":       KwCreate,
	"CROSS":        KwCross,
	"CUBE":         KwCube,
	"CURRENT":      KwCurrent,
	"CURRENT_USER": KwCurrentUser,
	"DATABASE":     KwDatabase,
	"DATABASES":    KwDatabases,
	"DATE":         KwDate,
	"DAY":          KwDay,
	"DEDUPLICATE":  KwDeduplicate,
	"DEFAULT":      KwDefault,
	"DELAY":        KwDelay,
	"DELETE":       KwDelete,
	"DESC":         KwDesc,
	"DESCENDING":   KwDescending,
	"DESCRIBE":     KwDescribe,
	"DETACH":       KwDetach,
	"DICTIONARIES": KwDictionaries,
	"DICTIONARY":   KwDictionary,
	"DISK":         KwDisk,
	"DISTINCT":     KwDistinct,
	"DISTRIBUTED":  KwDistributed,
	"DROP":         KwDrop,
	"ELSE":         KwElse,
	"ENABLED":      KwEnabled,
	"END":          KwEnd,
	"ENGINE":       KwEngine,
	"ENGINES":      KwEngines,
	"ESTIMATE":     KwEstimate,
	"EVENTS":       KwEvents,
	"EXCEPT":       KwExcept,
	"EXCHANGE":     KwExchange,
	"EXISTS":       KwExists,
	"EXPLAIN":      KwExplain,
	"EXPRESSION":   KwExpression,
	"EXTENDED":     KwExtended,
	"EXTRACT":      KwExtract,
	"FALSE":        KwFalse,
	"FETCHES":      KwFetches,
	"FIELDS":       KwFields,
	"FILESYSTEM":   KwFilesystem,
	"FILL":         KwFill,
	"FINAL":        KwFinal,
	"FIRST":        KwFirst,
	"FLUSH":        KwFlush,
	"FOLLOWING":    KwFollowing,
	"FOR":          KwFor,
	"FORMAT":       KwFormat,
	"FREEZE":       KwFreeze,
	"FROM":         KwFrom,
	"FULL":         KwFull,
	"FUNCTION":     KwFunction,
	"FUNCTIONS":    KwFunctions,
	"GLOBAL":       KwGlobal,
	"GRANT":        KwGrant,
	"GRANTS":       KwGrants,
	"GRANULARITY":  KwGranularity,
	"GROUP":        KwGroup,
	"GROUPING":     KwGrouping,
	"HAVING":       KwHaving,
	"HIERARCHICAL": KwHierarchical,
	"HOUR":         KwHour,
	"ID":           KwId,
	"IDENTIFIED":   KwIdentified,
	"IF":           KwIf,
	"ILIKE":        KwIlike,
	"IMPLICIT":     KwImplicit,
	"IN":           KwIn,
	"INDEX":        KwIndex,
	"INDEXES":      KwIndexes,
	"INDICES":      KwIndices,
	"INF":          KwInf,
	"INFINITY":     KwInf,
	"INJECTIVE":    KwInjective,
	"INNER":        KwInner,
	"INSERT":       KwInsert,
	"INTERPOLATE":  KwInterpolate,
	"INTERSECT":    KwIntersect,
	"INTERVAL":     KwInterval,
	"INTO":         KwInto,
	"IS":           KwIs,
	"IS_OBJECT_ID": KwIsObjectId,
	"JOIN":         KwJoin,
	"KEY":          KwKey,
	"KEYS":         KwKeys,
	"KILL":         KwKill,
	"LAST":         KwLast,
	"LAYOUT":       KwLayout,
	"LEADING":      KwLeading,
	"LEFT":         KwLeft,
	"LIFETIME":     KwLifetime,
	"LIKE":         KwLike,
	"LIMIT":        KwLimit,
	"LIVE":         KwLive,
	"LOCAL":        KwLocal,
	"LOGS":         KwLogs,
	"MATERIALIZE":  KwMaterialize,
	"MATERIALIZED": KwMaterialized,
	"MAX":          KwMax,
	"MERGES":       KwMerges,
	"MICROSECOND":  KwMicrosecond,
	"MILLISECOND":  KwMillisecond,
	"MIN":          KwMin,
	"MINUTE":       KwMinute,
	"MODIFY":       KwModify,
	"MONTH":        KwMonth,
	"MOVE":         KwMove,
	"MUTATION":     KwMutation,
	"NAN":          KwNan,
	"NANOSECOND":   KwNanosecond,
	"NO":           KwNo,
	"NOT":          KwNot,
	"NULL":         KwNull,
	"NULLS":        KwNulls,
	"OFFSET":       KwOffset,
	"ON":           KwOn,
	"OPTIMIZE":     KwOptimize,
	"OR":           KwOr,
	"ORDER":        KwOrder,
	"OUTER":        KwOuter,
	"OUTFILE":      KwOutfile,
	"OVER":         KwOver,
	"OVERRIDE":     KwOverride,
	"PARTITION":    KwPartition,
	"PERMANENTLY":  KwPermanently,
	"PIPELINE":     KwPipeline,
	"PLAN":         KwPlan,
	"POLICIES":     KwPolicies,
	"POLICY":       KwPolicy,
	"POPULATE":     KwPopulate,
	"PRECEDING":    KwPreceding,
	"PREWHERE":     KwPrewhere,
	"PRIMARY":      KwPrimary,
	"PRIVILEGES":   KwPrivileges,
	"PROCESSLIST":  KwProcesslist,
	"PROFILE":      KwProfile,
	"PROFILES":     KwProfiles,
	"PROJECTION":   KwProjection,
	"QUALIFY":      KwQualify,
	"QUARTER":      KwQuarter,
	"QUERY":        KwQuery,
	"QUOTA":        KwQuota,
	"QUOTAS":       KwQuotas,
	"RANGE":        KwRange,
	"RECURSIVE":    KwRecursive,
	"REFRESH":      KwRefresh,
	"RELOAD":       KwReload,
	"REMOVE":       KwRemove,
	"RENAME":       KwRename,
	"REPLACE":      KwReplace,
	"REPLICA":      KwReplica,
	"REPLICATED":   KwReplicated,
	"RESTORE":      KwRestore,
	"REVOKE":       KwRevoke,
	"RIGHT":        KwRight,
	"ROLE":         KwRole,
	"ROLES":        KwRoles,
	"ROLLUP":       KwRollup,
	"ROW":          KwRow,
	"ROWS":         KwRows,
	"SAMPLE":       KwSample,
	"SECOND":       KwSecond,
	"SELECT":       KwSelect,
	"SEMI":         KwSemi,
	"SENDS":        KwSends,
	"SET":          KwSet,
	"SETS":         KwSets,
	"SETTING":      KwSetting,
	"SETTINGS":     KwSettings,
	"SHOW":         KwShow,
	"SOURCE":       KwSource,
	"START":        KwStart,
	"STEP":         KwStep,
	"STOP":         KwStop,
	"SUBSTRING":    KwSubstring,
	"SYNC":         KwSync,
	"SYNTAX":       KwSyntax,
	"SYSTEM":       KwSystem,
	"TABLE":        KwTable,
	"TABLES":       KwTables,
	"TEMPORARY":    KwTemporary,
	"TEST":         KwTest,
	"THEN":         KwThen,
	"TIES":         KwTies,
	"TIMEOUT":      KwTimeout,
	"TIMESTAMP":    KwTimestamp,
	"TO":           KwTo,
	"TOP":          KwTop,
	"TOTALS":       KwTotals,
	"TRAILING":     KwTrailing,
	"TREE":         KwTree,
	"TRIM":         KwTrim,
	"TRUE":         KwTrue,
	"TRUNCATE":     KwTruncate,
	"TTL":          KwTtl,
	"TYPE":         KwType,
	"UNBOUNDED":    KwUnbounded,
	"UNDROP":       KwUndrop,
	"UNION":        KwUnion,
	"UPDATE":       KwUpdate,
	"USE":          KwUse,
	"USER":         KwUser,
	"USERS":        KwUsers,
	"USING":        KwUsing,
	"UUID":         KwUuid,
	"VALUES":       KwValues,
	"VIEW":         KwView,
	"VOLUME":       KwVolume,
	"WATCH":        KwWatch,
	"WEEK":         KwWeek,
	"WHEN":         KwWhen,
	"WHERE":        KwWhere,
	"WINDOW":       KwWindow,
	"WITH":         KwWith,
	"YEAR":         KwYear,
	"YYYY":         KwYear,
}

// lookupKeyword returns the keyword TokenType for word (case-insensitive), or TokIdent.
func lookupKeyword(word string) TokenType {
	if kw, ok := keywords[strings.ToUpper(word)]; ok {
		return kw
	}
	return TokIdent
}

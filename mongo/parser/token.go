package parser

// Token represents a single lexical token.
type Token struct {
	Type int    // token type constant
	Str  string // string content (identifier name, string value, number text, etc.)
	Loc  int    // inclusive start byte offset
	End  int    // exclusive end byte offset
}

// Token type constants for literals and punctuation.
const (
	tokEOF = 0

	// Literal tokens (600+)
	tokNumber = 600 + iota // integer or float
	tokString              // double or single quoted string
	tokRegex               // /pattern/flags
	tokIdent               // identifier (including $-prefixed)

	// Multi-character operators / punctuation
	tokEqEq       // ==
	tokNotEq      // !=
	tokNotEqEq    // !==
	tokEqEqEq     // ===
	tokLessEq     // <=
	tokGreaterEq  // >=
	tokAnd        // &&
	tokOr         // ||
	tokArrow      // =>
	tokSpread     // ...
	tokOptChain   // ?.
	tokNullCoal   // ??
	tokPlusPlus   // ++
	tokMinusMinus // --
	tokPlusEq     // +=
	tokMinusEq    // -=
	tokStarEq     // *=
	tokSlashEq    // /=
	tokPercentEq  // %=
	tokPower      // **
	tokShiftLeft  // <<
	tokShiftRight // >>
)

// Keyword token constants. Values start at 700.
// mongosh keywords are case-sensitive (unlike SQL).
const (
	// --- mongosh top-level objects ---
	kwDb  = 700 + iota // db
	kwRs               // rs
	kwSh               // sh
	kwSp               // sp

	// --- show command variants ---
	kwShow         // show
	kwDbs          // dbs
	kwDatabases    // databases
	kwCollections  // collections
	kwTables       // tables
	kwProfile      // profile
	kwUsers        // users
	kwRoles        // roles
	kwLog          // log
	kwLogs         // logs
	kwStartupWarnings // startupWarnings

	// --- JavaScript literals/keywords ---
	kwTrue      // true
	kwFalse     // false
	kwNull      // null
	kwUndefined // undefined
	kwNew       // new
	kwVar       // var
	kwLet       // let
	kwConst     // const
	kwFunction  // function
	kwReturn    // return
	kwIf        // if
	kwElse      // else
	kwFor       // for
	kwWhile     // while
	kwDo        // do
	kwBreak     // break
	kwContinue  // continue
	kwSwitch    // switch
	kwCase      // case
	kwDefault   // default
	kwThrow     // throw
	kwTry       // try
	kwCatch     // catch
	kwFinally   // finally
	kwTypeof    // typeof
	kwInstanceof // instanceof
	kwVoid      // void
	kwDelete    // delete
	kwIn        // in
	kwOf        // of
	kwThis      // this
	kwClass     // class
	kwExtends   // extends
	kwSuper     // super
	kwYield     // yield
	kwAsync     // async
	kwAwait     // await
	kwExport    // export
	kwImport    // import
	kwFrom      // from

	// --- BSON helpers ---
	kwObjectId       // ObjectId
	kwNumberLong     // NumberLong
	kwNumberInt      // NumberInt
	kwNumberDecimal  // NumberDecimal
	kwTimestamp      // Timestamp
	kwDate           // Date
	kwISODate        // ISODate
	kwUUID           // UUID
	kwMD5            // MD5
	kwHexData        // HexData
	kwBinData        // BinData
	kwCode           // Code
	kwDBRef          // DBRef
	kwMinKey         // MinKey
	kwMaxKey         // MaxKey
	kwRegExp         // RegExp
	kwSymbol         // Symbol

	// --- Collection methods ---
	kwFind             // find
	kwFindOne          // findOne
	kwFindOneAndDelete // findOneAndDelete
	kwFindOneAndReplace // findOneAndReplace
	kwFindOneAndUpdate  // findOneAndUpdate
	kwInsertOne        // insertOne
	kwInsertMany       // insertMany
	kwUpdateOne        // updateOne
	kwUpdateMany       // updateMany
	kwDeleteOne        // deleteOne
	kwDeleteMany       // deleteMany
	kwReplaceOne       // replaceOne
	kwBulkWrite        // bulkWrite
	kwAggregate        // aggregate
	kwCount            // count
	kwCountDocuments   // countDocuments
	kwEstimatedDocumentCount // estimatedDocumentCount
	kwDistinct         // distinct
	kwMapReduce        // mapReduce
	kwWatch            // watch
	kwCreateIndex      // createIndex
	kwCreateIndexes    // createIndexes
	kwDropIndex        // dropIndex
	kwDropIndexes      // dropIndexes
	kwGetIndexes       // getIndexes
	kwReIndex          // reIndex
	kwDrop             // drop
	kwRenameCollection // renameCollection
	kwStats            // stats
	kwDataSize         // dataSize
	kwStorageSize      // storageSize
	kwTotalSize        // totalSize
	kwTotalIndexSize   // totalIndexSize
	kwValidate         // validate
	kwExplain          // explain
	kwGetShardDistribution // getShardDistribution
	kwLatencyStats     // latencyStats

	// --- Cursor methods ---
	kwSort      // sort
	kwLimit     // limit
	kwSkip      // skip
	kwToArray   // toArray
	kwForEach   // forEach
	kwMap       // map
	kwHasNext   // hasNext
	kwNext      // next
	kwItcount   // itcount
	kwSize      // size
	kwPretty    // pretty
	kwHint      // hint
	kwMin       // min
	kwMax       // max
	kwReadPref  // readPref
	kwComment   // comment
	kwBatchSize // batchSize
	kwClose     // close
	kwCollation // collation
	kwNoCursorTimeout // noCursorTimeout
	kwAllowPartialResults // allowPartialResults
	kwReturnKey // returnKey
	kwShowRecordId // showRecordId
	kwAllowDiskUse // allowDiskUse
	kwMaxTimeMS // maxTimeMS
	kwReadConcern // readConcern
	kwWriteConcern // writeConcern
	kwTailable  // tailable
	kwOplogReplay // oplogReplay
	kwProjection  // projection

	// --- Database methods ---
	kwGetName                 // getName
	kwGetSiblingDB            // getSiblingDB
	kwGetMongo                // getMongo
	kwGetCollectionNames      // getCollectionNames
	kwGetCollectionInfos      // getCollectionInfos
	kwGetCollection           // getCollection
	kwCreateCollection        // createCollection
	kwCreateView              // createView
	kwDropDatabase            // dropDatabase
	kwAdminCommand            // adminCommand
	kwRunCommand              // runCommand
	kwGetProfilingStatus      // getProfilingStatus
	kwSetProfilingLevel       // setProfilingLevel
	kwGetLogComponents        // getLogComponents
	kwSetLogLevel             // setLogLevel
	kwFsyncLock               // fsyncLock
	kwFsyncUnlock             // fsyncUnlock
	kwCurrentOp               // currentOp
	kwKillOp                  // killOp
	kwGetUser                 // getUser
	kwGetUsers                // getUsers
	kwCreateUser              // createUser
	kwUpdateUser              // updateUser
	kwDropUser                // dropUser
	kwDropAllUsers            // dropAllUsers
	kwGrantRolesToUser        // grantRolesToUser
	kwRevokeRolesFromUser     // revokeRolesFromUser
	kwGetRole                 // getRole
	kwGetRoles                // getRoles
	kwCreateRole              // createRole
	kwUpdateRole              // updateRole
	kwDropRole                // dropRole
	kwDropAllRoles            // dropAllRoles
	kwGrantPrivilegesToRole   // grantPrivilegesToRole
	kwRevokePrivilegesFromRole // revokePrivilegesFromRole
	kwGrantRolesToRole        // grantRolesToRole
	kwRevokeRolesFromRole     // revokeRolesFromRole
	kwServerStatus            // serverStatus
	kwIsMaster                // isMaster
	kwHello                   // hello
	kwHostInfo                // hostInfo

	// --- Bulk operations ---
	kwInitializeOrderedBulkOp   // initializeOrderedBulkOp
	kwInitializeUnorderedBulkOp // initializeUnorderedBulkOp
	kwInsert                    // insert
	kwUpdate                    // update
	kwRemove                    // remove
	kwExecute                   // execute
	kwTojson                    // tojson
	kwToJSON                    // toJSON

	// --- Replica set methods ---
	kwStatus       // status
	kwConf         // conf
	kwConfig       // config
	kwInitiate     // initiate
	kwReconfig     // reconfig
	kwAdd          // add
	kwAddArb       // addArb
	kwStepDown     // stepDown
	kwFreeze       // freeze
	kwSlaveOk      // slaveOk
	kwSecondaryOk  // secondaryOk
	kwSyncFrom     // syncFrom
	kwPrintReplicationInfo     // printReplicationInfo
	kwPrintSecondaryReplicationInfo // printSecondaryReplicationInfo

	// --- Sharding methods ---
	kwAddShard           // addShard
	kwAddShardTag        // addShardTag
	kwAddShardToZone     // addShardToZone
	kwAddTagRange        // addTagRange
	kwDisableAutoSplit   // disableAutoSplit
	kwEnableAutoSplit    // enableAutoSplit
	kwEnableSharding     // enableSharding
	kwDisableBalancing   // disableBalancing
	kwEnableBalancing    // enableBalancing
	kwGetBalancerState   // getBalancerState
	kwIsBalancerRunning  // isBalancerRunning
	kwMoveChunk          // moveChunk
	kwRemoveRangeFromZone // removeRangeFromZone
	kwRemoveShard        // removeShard
	kwRemoveShardTag     // removeShardTag
	kwRemoveShardFromZone // removeShardFromZone
	kwSetBalancerState   // setBalancerState
	kwShardCollection    // shardCollection
	kwSplitAt            // splitAt
	kwSplitFind          // splitFind
	kwStartBalancer      // startBalancer
	kwStopBalancer       // stopBalancer
	kwUpdateZoneKeyRange // updateZoneKeyRange

	// --- Connection methods ---
	kwMongo       // Mongo
	kwConnect     // connect

	// --- Encryption ---
	kwGetKeyVault          // getKeyVault
	kwGetClientEncryption  // getClientEncryption
	kwCreateKey            // createKey
	kwGetKey               // getKey
	kwGetKeys              // getKeys
	kwDeleteKey            // deleteKey
	kwAddKeyAlternateName  // addKeyAlternateName
	kwRemoveKeyAlternateName // removeKeyAlternateName
	kwRewrapManyDataKey    // rewrapManyDataKey
	kwEncrypt              // encrypt
	kwDecrypt              // decrypt

	// --- Plan cache ---
	kwGetPlanCache // getPlanCache
	kwClear        // clear
	kwList         // list

	// --- Native mongosh functions ---
	kwSleep            // sleep
	kwLoad             // load
	kwPrint            // print
	kwPrintjson        // printjson
	kwQuit             // quit
	kwExit             // exit
	kwVersion          // version
	kwCls              // cls
	kwHelp             // help
	kwIt               // it
	kwNumberIsNaN      // Number.isNaN
	kwIsNaN            // isNaN
	kwIsFinite         // isFinite
	kwParseInt         // parseInt
	kwParseFloat       // parseFloat
	kwEncodeURI        // encodeURI
	kwEncodeURIComponent // encodeURIComponent
	kwDecodeURI        // decodeURI
	kwDecodeURIComponent // decodeURIComponent
	kwJSON             // JSON

	// --- Misc ---
	kwConvert       // convert
	kwTo            // to
	kwType          // type
	kwNumber        // Number
	kwString        // String
	kwBoolean       // Boolean
	kwArray         // Array
	kwObject        // Object
	kwMath          // Math
	kwInfinity      // Infinity
	kwNaN           // NaN
)

// keywords maps keyword strings to their token types.
// mongosh keywords are case-sensitive.
var keywords = map[string]int{
	// Top-level objects
	"db": kwDb,
	"rs": kwRs,
	"sh": kwSh,
	"sp": kwSp,

	// show variants
	"show":            kwShow,
	"dbs":             kwDbs,
	"databases":       kwDatabases,
	"collections":     kwCollections,
	"tables":          kwTables,
	"profile":         kwProfile,
	"users":           kwUsers,
	"roles":           kwRoles,
	"log":             kwLog,
	"logs":            kwLogs,
	"startupWarnings": kwStartupWarnings,

	// JS literals/keywords
	"true":       kwTrue,
	"false":      kwFalse,
	"null":       kwNull,
	"undefined":  kwUndefined,
	"new":        kwNew,
	"var":        kwVar,
	"let":        kwLet,
	"const":      kwConst,
	"function":   kwFunction,
	"return":     kwReturn,
	"if":         kwIf,
	"else":       kwElse,
	"for":        kwFor,
	"while":      kwWhile,
	"do":         kwDo,
	"break":      kwBreak,
	"continue":   kwContinue,
	"switch":     kwSwitch,
	"case":       kwCase,
	"default":    kwDefault,
	"throw":      kwThrow,
	"try":        kwTry,
	"catch":      kwCatch,
	"finally":    kwFinally,
	"typeof":     kwTypeof,
	"instanceof": kwInstanceof,
	"void":       kwVoid,
	"delete":     kwDelete,
	"in":         kwIn,
	"of":         kwOf,
	"this":       kwThis,
	"class":      kwClass,
	"extends":    kwExtends,
	"super":      kwSuper,
	"yield":      kwYield,
	"async":      kwAsync,
	"await":      kwAwait,
	"export":     kwExport,
	"import":     kwImport,
	"from":       kwFrom,

	// BSON helpers
	"ObjectId":      kwObjectId,
	"NumberLong":    kwNumberLong,
	"NumberInt":     kwNumberInt,
	"NumberDecimal": kwNumberDecimal,
	"Timestamp":     kwTimestamp,
	"Date":          kwDate,
	"ISODate":       kwISODate,
	"UUID":          kwUUID,
	"MD5":           kwMD5,
	"HexData":       kwHexData,
	"BinData":       kwBinData,
	"Code":          kwCode,
	"DBRef":         kwDBRef,
	"MinKey":        kwMinKey,
	"MaxKey":        kwMaxKey,
	"RegExp":        kwRegExp,
	"Symbol":        kwSymbol,

	// Collection methods
	"find":                    kwFind,
	"findOne":                 kwFindOne,
	"findOneAndDelete":        kwFindOneAndDelete,
	"findOneAndReplace":       kwFindOneAndReplace,
	"findOneAndUpdate":        kwFindOneAndUpdate,
	"insertOne":               kwInsertOne,
	"insertMany":              kwInsertMany,
	"updateOne":               kwUpdateOne,
	"updateMany":              kwUpdateMany,
	"deleteOne":               kwDeleteOne,
	"deleteMany":              kwDeleteMany,
	"replaceOne":              kwReplaceOne,
	"bulkWrite":               kwBulkWrite,
	"aggregate":               kwAggregate,
	"count":                   kwCount,
	"countDocuments":          kwCountDocuments,
	"estimatedDocumentCount":  kwEstimatedDocumentCount,
	"distinct":                kwDistinct,
	"mapReduce":               kwMapReduce,
	"watch":                   kwWatch,
	"createIndex":             kwCreateIndex,
	"createIndexes":           kwCreateIndexes,
	"dropIndex":               kwDropIndex,
	"dropIndexes":             kwDropIndexes,
	"getIndexes":              kwGetIndexes,
	"reIndex":                 kwReIndex,
	"drop":                    kwDrop,
	"renameCollection":        kwRenameCollection,
	"stats":                   kwStats,
	"dataSize":                kwDataSize,
	"storageSize":             kwStorageSize,
	"totalSize":               kwTotalSize,
	"totalIndexSize":          kwTotalIndexSize,
	"validate":                kwValidate,
	"explain":                 kwExplain,
	"getShardDistribution":    kwGetShardDistribution,
	"latencyStats":            kwLatencyStats,

	// Cursor methods
	"sort":                kwSort,
	"limit":               kwLimit,
	"skip":                kwSkip,
	"toArray":             kwToArray,
	"forEach":             kwForEach,
	"map":                 kwMap,
	"hasNext":             kwHasNext,
	"next":                kwNext,
	"itcount":             kwItcount,
	"size":                kwSize,
	"pretty":              kwPretty,
	"hint":                kwHint,
	"min":                 kwMin,
	"max":                 kwMax,
	"readPref":            kwReadPref,
	"comment":             kwComment,
	"batchSize":           kwBatchSize,
	"close":               kwClose,
	"collation":           kwCollation,
	"noCursorTimeout":     kwNoCursorTimeout,
	"allowPartialResults": kwAllowPartialResults,
	"returnKey":           kwReturnKey,
	"showRecordId":        kwShowRecordId,
	"allowDiskUse":        kwAllowDiskUse,
	"maxTimeMS":           kwMaxTimeMS,
	"readConcern":         kwReadConcern,
	"writeConcern":        kwWriteConcern,
	"tailable":            kwTailable,
	"oplogReplay":         kwOplogReplay,
	"projection":          kwProjection,

	// Database methods
	"getName":            kwGetName,
	"getSiblingDB":       kwGetSiblingDB,
	"getMongo":           kwGetMongo,
	"getCollectionNames": kwGetCollectionNames,
	"getCollectionInfos": kwGetCollectionInfos,
	"getCollection":      kwGetCollection,
	"createCollection":   kwCreateCollection,
	"createView":         kwCreateView,
	"dropDatabase":       kwDropDatabase,
	"adminCommand":       kwAdminCommand,
	"runCommand":         kwRunCommand,
	"getProfilingStatus": kwGetProfilingStatus,
	"setProfilingLevel":  kwSetProfilingLevel,
	"getLogComponents":   kwGetLogComponents,
	"setLogLevel":        kwSetLogLevel,
	"fsyncLock":          kwFsyncLock,
	"fsyncUnlock":        kwFsyncUnlock,
	"currentOp":          kwCurrentOp,
	"killOp":             kwKillOp,
	"getUser":            kwGetUser,
	"getUsers":           kwGetUsers,
	"createUser":         kwCreateUser,
	"updateUser":         kwUpdateUser,
	"dropUser":           kwDropUser,
	"dropAllUsers":       kwDropAllUsers,
	"grantRolesToUser":   kwGrantRolesToUser,
	"revokeRolesFromUser": kwRevokeRolesFromUser,
	"getRole":            kwGetRole,
	"getRoles":           kwGetRoles,
	"createRole":         kwCreateRole,
	"updateRole":         kwUpdateRole,
	"dropRole":           kwDropRole,
	"dropAllRoles":       kwDropAllRoles,
	"grantPrivilegesToRole":   kwGrantPrivilegesToRole,
	"revokePrivilegesFromRole": kwRevokePrivilegesFromRole,
	"grantRolesToRole":        kwGrantRolesToRole,
	"revokeRolesFromRole":     kwRevokeRolesFromRole,
	"serverStatus":            kwServerStatus,
	"isMaster":                kwIsMaster,
	"hello":                   kwHello,
	"hostInfo":                kwHostInfo,

	// Bulk operations
	"initializeOrderedBulkOp":   kwInitializeOrderedBulkOp,
	"initializeUnorderedBulkOp": kwInitializeUnorderedBulkOp,
	"insert":                    kwInsert,
	"update":                    kwUpdate,
	"remove":                    kwRemove,
	"execute":                   kwExecute,
	"tojson":                    kwTojson,
	"toJSON":                    kwToJSON,

	// Replica set methods
	"status":                          kwStatus,
	"conf":                            kwConf,
	"config":                          kwConfig,
	"initiate":                        kwInitiate,
	"reconfig":                        kwReconfig,
	"add":                             kwAdd,
	"addArb":                          kwAddArb,
	"stepDown":                        kwStepDown,
	"freeze":                          kwFreeze,
	"slaveOk":                         kwSlaveOk,
	"secondaryOk":                     kwSecondaryOk,
	"syncFrom":                        kwSyncFrom,
	"printReplicationInfo":            kwPrintReplicationInfo,
	"printSecondaryReplicationInfo":   kwPrintSecondaryReplicationInfo,

	// Sharding methods
	"addShard":            kwAddShard,
	"addShardTag":         kwAddShardTag,
	"addShardToZone":      kwAddShardToZone,
	"addTagRange":         kwAddTagRange,
	"disableAutoSplit":    kwDisableAutoSplit,
	"enableAutoSplit":     kwEnableAutoSplit,
	"enableSharding":      kwEnableSharding,
	"disableBalancing":    kwDisableBalancing,
	"enableBalancing":     kwEnableBalancing,
	"getBalancerState":    kwGetBalancerState,
	"isBalancerRunning":   kwIsBalancerRunning,
	"moveChunk":           kwMoveChunk,
	"removeRangeFromZone": kwRemoveRangeFromZone,
	"removeShard":         kwRemoveShard,
	"removeShardTag":      kwRemoveShardTag,
	"removeShardFromZone": kwRemoveShardFromZone,
	"setBalancerState":    kwSetBalancerState,
	"shardCollection":     kwShardCollection,
	"splitAt":             kwSplitAt,
	"splitFind":           kwSplitFind,
	"startBalancer":       kwStartBalancer,
	"stopBalancer":        kwStopBalancer,
	"updateZoneKeyRange":  kwUpdateZoneKeyRange,

	// Connection methods
	"Mongo":   kwMongo,
	"connect": kwConnect,

	// Encryption
	"getKeyVault":            kwGetKeyVault,
	"getClientEncryption":    kwGetClientEncryption,
	"createKey":              kwCreateKey,
	"getKey":                 kwGetKey,
	"getKeys":                kwGetKeys,
	"deleteKey":              kwDeleteKey,
	"addKeyAlternateName":    kwAddKeyAlternateName,
	"removeKeyAlternateName": kwRemoveKeyAlternateName,
	"rewrapManyDataKey":      kwRewrapManyDataKey,
	"encrypt":                kwEncrypt,
	"decrypt":                kwDecrypt,

	// Plan cache
	"getPlanCache": kwGetPlanCache,
	"clear":        kwClear,
	"list":         kwList,

	// Native functions
	"sleep":                kwSleep,
	"load":                 kwLoad,
	"print":                kwPrint,
	"printjson":            kwPrintjson,
	"quit":                 kwQuit,
	"exit":                 kwExit,
	"version":              kwVersion,
	"cls":                  kwCls,
	"help":                 kwHelp,
	"it":                   kwIt,
	"isNaN":                kwIsNaN,
	"isFinite":             kwIsFinite,
	"parseInt":             kwParseInt,
	"parseFloat":           kwParseFloat,
	"encodeURI":            kwEncodeURI,
	"encodeURIComponent":   kwEncodeURIComponent,
	"decodeURI":            kwDecodeURI,
	"decodeURIComponent":   kwDecodeURIComponent,
	"JSON":                 kwJSON,

	// Misc
	"convert":  kwConvert,
	"to":       kwTo,
	"type":     kwType,
	"Number":   kwNumber,
	"String":   kwString,
	"Boolean":  kwBoolean,
	"Array":    kwArray,
	"Object":   kwObject,
	"Math":     kwMath,
	"Infinity": kwInfinity,
	"NaN":      kwNaN,
}

// tokenName returns a human-readable name for a token type, used in error messages.
func tokenName(t int) string {
	switch {
	case t == tokEOF:
		return "end of input"
	case t == tokNumber:
		return "number"
	case t == tokString:
		return "string"
	case t == tokRegex:
		return "regex"
	case t == tokIdent:
		return "identifier"
	case t == tokEqEq:
		return "=="
	case t == tokNotEq:
		return "!="
	case t == tokNotEqEq:
		return "!=="
	case t == tokEqEqEq:
		return "==="
	case t == tokLessEq:
		return "<="
	case t == tokGreaterEq:
		return ">="
	case t == tokAnd:
		return "&&"
	case t == tokOr:
		return "||"
	case t == tokArrow:
		return "=>"
	case t == tokSpread:
		return "..."
	case t == tokOptChain:
		return "?."
	case t == tokNullCoal:
		return "??"
	case t == tokPlusPlus:
		return "++"
	case t == tokMinusMinus:
		return "--"
	case t == tokPlusEq:
		return "+="
	case t == tokMinusEq:
		return "-="
	case t == tokStarEq:
		return "*="
	case t == tokSlashEq:
		return "/="
	case t == tokPercentEq:
		return "%="
	case t == tokPower:
		return "**"
	case t == tokShiftLeft:
		return "<<"
	case t == tokShiftRight:
		return ">>"
	case t >= 700:
		// keyword — find by reverse lookup
		for name, v := range keywords {
			if v == t {
				return name
			}
		}
		return "keyword"
	case t > 0 && t < 128:
		return string(rune(t))
	default:
		return "?"
	}
}

package completion

import "github.com/bytebase/omni/mongo/catalog"

// candidatesForContext returns the raw candidate list for a given context,
// optionally enriched by the catalog.
func candidatesForContext(ctx completionContext, cat *catalog.Catalog) []Candidate {
	switch ctx {
	case contextTopLevel:
		return topLevelCandidates()
	case contextAfterDbDot:
		return afterDbDotCandidates(cat)
	case contextAfterCollDot:
		return collectionMethodCandidates()
	case contextAfterBracket:
		return bracketCandidates(cat)
	case contextCursorChain:
		return cursorMethodCandidates()
	case contextShowTarget:
		return showTargetCandidates()
	case contextAfterRsDot:
		return rsMethodCandidates()
	case contextAfterShDot:
		return shMethodCandidates()
	case contextAggStage:
		return aggStageCandidates()
	case contextQueryOperator:
		return queryOperatorCandidates()
	case contextInsideArgs:
		return insideArgsCandidates()
	case contextDocumentKey:
		return documentKeyCandidates()
	default:
		return nil
	}
}

func topLevelCandidates() []Candidate {
	keywords := []string{
		"db", "rs", "sh", "sp", "show",
		"sleep", "load", "print", "printjson",
		"quit", "exit", "help", "it", "cls", "version",
	}
	candidates := make([]Candidate, 0, len(keywords)+len(bsonHelpers))
	for _, kw := range keywords {
		candidates = append(candidates, Candidate{Text: kw, Type: CandidateKeyword})
	}
	for _, h := range bsonHelpers {
		candidates = append(candidates, Candidate{Text: h, Type: CandidateBSONHelper})
	}
	return candidates
}

func afterDbDotCandidates(cat *catalog.Catalog) []Candidate {
	var candidates []Candidate
	if cat != nil {
		for _, name := range cat.Collections() {
			candidates = append(candidates, Candidate{Text: name, Type: CandidateCollection})
		}
	}
	for _, m := range dbMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateDbMethod})
	}
	return candidates
}

func collectionMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(collectionMethods))
	for _, m := range collectionMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateMethod})
	}
	return candidates
}

func bracketCandidates(cat *catalog.Catalog) []Candidate {
	if cat == nil {
		return nil
	}
	var candidates []Candidate
	for _, name := range cat.Collections() {
		candidates = append(candidates, Candidate{Text: name, Type: CandidateCollection})
	}
	return candidates
}

func cursorMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(cursorMethods))
	for _, m := range cursorMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateCursorMethod})
	}
	return candidates
}

func showTargetCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(showTargets))
	for _, t := range showTargets {
		candidates = append(candidates, Candidate{Text: t, Type: CandidateShowTarget})
	}
	return candidates
}

func rsMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(rsMethods))
	for _, m := range rsMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateRsMethod})
	}
	return candidates
}

func shMethodCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(shMethods))
	for _, m := range shMethods {
		candidates = append(candidates, Candidate{Text: m, Type: CandidateShMethod})
	}
	return candidates
}

func aggStageCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(aggStages))
	for _, s := range aggStages {
		candidates = append(candidates, Candidate{Text: s, Type: CandidateAggStage})
	}
	return candidates
}

func queryOperatorCandidates() []Candidate {
	candidates := make([]Candidate, 0, len(queryOperators))
	for _, op := range queryOperators {
		candidates = append(candidates, Candidate{Text: op, Type: CandidateQueryOperator})
	}
	return candidates
}

func insideArgsCandidates() []Candidate {
	literals := []string{"true", "false", "null"}
	candidates := make([]Candidate, 0, len(bsonHelpers)+len(literals))
	for _, h := range bsonHelpers {
		candidates = append(candidates, Candidate{Text: h, Type: CandidateBSONHelper})
	}
	for _, l := range literals {
		candidates = append(candidates, Candidate{Text: l, Type: CandidateKeyword})
	}
	return candidates
}

func documentKeyCandidates() []Candidate {
	candidates := queryOperatorCandidates()
	candidates = append(candidates, insideArgsCandidates()...)
	return candidates
}

// --- Hardcoded candidate lists ---

var bsonHelpers = []string{
	"ObjectId", "NumberLong", "NumberInt", "NumberDecimal",
	"Timestamp", "Date", "ISODate", "UUID",
	"MD5", "HexData", "BinData", "Code",
	"DBRef", "MinKey", "MaxKey", "RegExp", "Symbol",
}

var collectionMethods = []string{
	"find", "findOne", "findOneAndDelete", "findOneAndReplace", "findOneAndUpdate",
	"insertOne", "insertMany",
	"updateOne", "updateMany",
	"deleteOne", "deleteMany",
	"replaceOne", "bulkWrite",
	"aggregate",
	"count", "countDocuments", "estimatedDocumentCount",
	"distinct", "mapReduce", "watch",
	"createIndex", "createIndexes",
	"dropIndex", "dropIndexes", "getIndexes", "reIndex",
	"drop", "renameCollection",
	"stats", "dataSize", "storageSize", "totalSize", "totalIndexSize",
	"validate", "explain",
	"getShardDistribution", "latencyStats",
	"getPlanCache",
	"initializeOrderedBulkOp", "initializeUnorderedBulkOp",
}

var cursorMethods = []string{
	"sort", "limit", "skip",
	"toArray", "forEach", "map",
	"hasNext", "next", "itcount", "size",
	"pretty", "hint", "min", "max",
	"readPref", "comment", "batchSize", "close",
	"collation", "noCursorTimeout", "allowPartialResults",
	"returnKey", "showRecordId", "allowDiskUse",
	"maxTimeMS", "readConcern", "writeConcern",
	"tailable", "oplogReplay", "projection",
}

var dbMethods = []string{
	"getName", "getSiblingDB", "getMongo",
	"getCollectionNames", "getCollectionInfos", "getCollection",
	"createCollection", "createView",
	"dropDatabase",
	"adminCommand", "runCommand",
	"getProfilingStatus", "setProfilingLevel",
	"getLogComponents", "setLogLevel",
	"fsyncLock", "fsyncUnlock",
	"currentOp", "killOp",
	"getUser", "getUsers", "createUser", "updateUser",
	"dropUser", "dropAllUsers",
	"grantRolesToUser", "revokeRolesFromUser",
	"getRole", "getRoles", "createRole", "updateRole",
	"dropRole", "dropAllRoles",
	"grantPrivilegesToRole", "revokePrivilegesFromRole",
	"grantRolesToRole", "revokeRolesFromRole",
	"serverStatus", "isMaster", "hello", "hostInfo",
}

var showTargets = []string{
	"dbs", "databases", "collections", "tables",
	"profile", "users", "roles",
	"log", "logs", "startupWarnings",
}

var rsMethods = []string{
	"status", "conf", "config",
	"initiate", "reconfig",
	"add", "addArb",
	"stepDown", "freeze",
	"slaveOk", "secondaryOk",
	"syncFrom",
	"printReplicationInfo", "printSecondaryReplicationInfo",
}

var shMethods = []string{
	"addShard", "addShardTag", "addShardToZone", "addTagRange",
	"disableAutoSplit", "enableAutoSplit",
	"enableSharding", "disableBalancing", "enableBalancing",
	"getBalancerState", "isBalancerRunning",
	"moveChunk",
	"removeRangeFromZone", "removeShard", "removeShardTag", "removeShardFromZone",
	"setBalancerState", "shardCollection",
	"splitAt", "splitFind",
	"startBalancer", "stopBalancer",
	"updateZoneKeyRange",
	"status",
}

var aggStages = []string{
	"$match", "$group", "$project", "$sort", "$limit", "$skip",
	"$unwind", "$lookup", "$graphLookup",
	"$addFields", "$set", "$unset",
	"$out", "$merge",
	"$bucket", "$bucketAuto", "$facet",
	"$replaceRoot", "$replaceWith",
	"$sample", "$count", "$redact",
	"$geoNear", "$setWindowFields", "$fill", "$densify",
	"$unionWith",
	"$collStats", "$indexStats", "$planCacheStats",
	"$search", "$searchMeta", "$changeStream",
}

var queryOperators = []string{
	// Comparison
	"$eq", "$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin",
	// Logical
	"$and", "$or", "$not", "$nor",
	// Element
	"$exists", "$type",
	// Evaluation
	"$regex", "$expr", "$mod", "$text", "$where", "$jsonSchema",
	// Array
	"$all", "$elemMatch", "$size",
	// Geospatial
	"$geoWithin", "$geoIntersects", "$near", "$nearSphere",
}

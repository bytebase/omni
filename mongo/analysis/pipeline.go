package analysis

import (
	"github.com/bytebase/omni/mongo/ast"
)

// shapePreservingStages lists pipeline stages that do not change document shape.
var shapePreservingStages = map[string]bool{
	"$match":           true,
	"$sort":            true,
	"$limit":           true,
	"$skip":            true,
	"$sample":          true,
	"$addFields":       true,
	"$set":             true,
	"$unset":           true,
	"$geoNear":         true,
	"$setWindowFields": true,
	"$fill":            true,
	"$redact":          true,
	"$unwind":          true,
}

// analyzePipelineInto populates a with pipeline information derived from args.
func analyzePipelineInto(args []ast.Node, a *StatementAnalysis) {
	if len(args) == 0 {
		a.ShapePreserving = true
		return
	}

	arr, ok := args[0].(*ast.Array)
	if !ok {
		a.UnsupportedStage = "unknown"
		return
	}

	a.ShapePreserving = true
	fields := make(map[string]struct{})

	for _, elem := range arr.Elements {
		doc, ok := elem.(*ast.Document)
		if !ok || len(doc.Pairs) == 0 {
			a.UnsupportedStage = "unknown"
			a.ShapePreserving = false
			return
		}

		stageName := doc.Pairs[0].Key
		a.PipelineStages = append(a.PipelineStages, stageName)

		if shapePreservingStages[stageName] {
			if stageName == "$match" {
				stageVal := doc.Pairs[0].Value
				if matchDoc, ok := stageVal.(*ast.Document); ok {
					collectFromDocument(matchDoc, "", fields)
				}
			}
			continue
		}

		if stageName == "$lookup" {
			join, unsupported := extractLookup(doc.Pairs[0].Value)
			if unsupported {
				a.UnsupportedStage = "$lookup"
				a.ShapePreserving = false
				return
			}
			if join != nil {
				a.Joins = append(a.Joins, *join)
			}
			continue
		}

		if stageName == "$graphLookup" {
			join := extractGraphLookup(doc.Pairs[0].Value)
			if join != nil {
				a.Joins = append(a.Joins, *join)
			}
			continue
		}

		// Not shape-preserving and not a join stage
		a.UnsupportedStage = stageName
		a.ShapePreserving = false
		return
	}

	if len(fields) > 0 {
		a.PredicateFields = sortedKeys(fields)
	}
}

// extractLookup extracts join info from a $lookup stage value.
// Returns (nil, true) if a pipeline-form lookup is detected (unsupported).
// Returns (*JoinInfo, false) if from+as are found.
func extractLookup(node ast.Node) (*JoinInfo, bool) {
	doc, ok := node.(*ast.Document)
	if !ok {
		return nil, false
	}

	var from, as string
	for _, kv := range doc.Pairs {
		switch kv.Key {
		case "pipeline":
			return nil, true
		case "from":
			if s, ok := kv.Value.(*ast.StringLiteral); ok {
				from = s.Value
			}
		case "as":
			if s, ok := kv.Value.(*ast.StringLiteral); ok {
				as = s.Value
			}
		}
	}

	if from != "" && as != "" {
		return &JoinInfo{Collection: from, AsField: as}, false
	}
	return nil, false
}

// extractGraphLookup extracts join info from a $graphLookup stage value.
func extractGraphLookup(node ast.Node) *JoinInfo {
	doc, ok := node.(*ast.Document)
	if !ok {
		return nil
	}

	var from, as string
	for _, kv := range doc.Pairs {
		switch kv.Key {
		case "from":
			if s, ok := kv.Value.(*ast.StringLiteral); ok {
				from = s.Value
			}
		case "as":
			if s, ok := kv.Value.(*ast.StringLiteral); ok {
				as = s.Value
			}
		}
	}

	if from != "" && as != "" {
		return &JoinInfo{Collection: from, AsField: as}
	}
	return nil
}

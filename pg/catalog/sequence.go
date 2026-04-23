package catalog

import (
	"fmt"
	"math"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Sequence represents a PostgreSQL sequence.
type Sequence struct {
	OID         uint32
	Name        string
	Schema      *Schema
	TypeOID     uint32 // int2/int4/int8
	Start       int64
	Increment   int64
	MinValue    int64
	MaxValue    int64
	CacheValue  int64
	Cycle       bool
	OwnerRelOID uint32 // 0 if not owned
	OwnerAttNum int16
}

// sequence type bounds
var seqBounds = map[uint32][2]int64{
	INT2OID: {math.MinInt16, math.MaxInt16},
	INT4OID: {math.MinInt32, math.MaxInt32},
	INT8OID: {math.MinInt64, math.MaxInt64},
}

// DefineSequence creates a new sequence from a parsed CREATE SEQUENCE statement.
//
// pg: src/backend/commands/sequence.c — DefineSequence
func (c *Catalog) DefineSequence(stmt *nodes.CreateSeqStmt) error {
	ifNotExists := stmt.IfNotExists
	var schemaName, seqName string
	if stmt.Sequence != nil {
		schemaName = stmt.Sequence.Schemaname
		seqName = stmt.Sequence.Relname
	}

	// Process DefElem options.
	var pIncrement, pMinValue, pMaxValue, pStart, pCacheValue *int64
	var cycle bool
	var ownedBy string
	var typeName string

	// Duplicate DefElem detection.
	// pg: src/backend/commands/sequence.c — init_params (duplicate option checks)
	var sawIncrement, sawMinValue, sawMaxValue, sawStart, sawCache, sawCycle, sawOwnedBy, sawAs bool

	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			d, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch d.Defname {
			case "increment":
				if sawIncrement {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawIncrement = true
				if v, ok := defElemInt(d); ok {
					pIncrement = &v
				}
			case "minvalue":
				if sawMinValue {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawMinValue = true
				if v, ok := defElemInt(d); ok {
					pMinValue = &v
				}
			case "maxvalue":
				if sawMaxValue {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawMaxValue = true
				if v, ok := defElemInt(d); ok {
					pMaxValue = &v
				}
			case "start":
				if sawStart {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawStart = true
				if v, ok := defElemInt(d); ok {
					pStart = &v
				}
			case "cache":
				if sawCache {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawCache = true
				if v, ok := defElemInt(d); ok {
					pCacheValue = &v
				}
			case "cycle":
				if sawCycle {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawCycle = true
				cycle = defElemBool(d)
			case "owned_by":
				if sawOwnedBy {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawOwnedBy = true
				// Arg is a List of String nodes: [table, column] or [schema, table, column].
				if list, ok := d.Arg.(*nodes.List); ok {
					parts := stringListItems(list)
					if len(parts) == 1 && strings.ToUpper(parts[0]) == "NONE" {
						// OWNED BY NONE — not applicable for CREATE, ignore.
					} else if len(parts) >= 2 {
						// Join last two elements as "table.column".
						ownedBy = parts[len(parts)-2] + "." + parts[len(parts)-1]
					}
				}
			case "as":
				if sawAs {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawAs = true
				typeName = defElemString(d)
			case "sequence_name":
				// sequence_name is only valid for identity columns, not CREATE SEQUENCE.
				// pg: src/backend/commands/sequence.c — init_params (sequence_name rejection)
				return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
			}
		}
	}

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Check name conflicts.
	if _, exists := schema.Sequences[seqName]; exists {
		if ifNotExists {
			return nil
		}
		return errDuplicateObject("sequence", seqName)
	}
	if _, exists := schema.Relations[seqName]; exists {
		if ifNotExists {
			return nil
		}
		return errDuplicateTable(seqName)
	}
	if _, exists := schema.Indexes[seqName]; exists {
		if ifNotExists {
			return nil
		}
		return errDuplicateObject("index", seqName)
	}

	// Resolve type.
	typeOID, err := c.resolveSeqType(typeName)
	if err != nil {
		return err
	}
	bounds := seqBounds[typeOID]

	// Apply defaults.
	increment := int64(1)
	if pIncrement != nil {
		increment = *pIncrement
	}
	if increment == 0 {
		return errInvalidParameterValue("INCREMENT must not be zero")
	}

	ascending := increment > 0

	minValue := bounds[0]
	if pMinValue != nil {
		minValue = *pMinValue
	} else if ascending {
		minValue = 1
	}

	maxValue := bounds[1]
	if pMaxValue != nil {
		maxValue = *pMaxValue
	} else if !ascending {
		maxValue = -1
	}

	// Validate MAXVALUE against type bounds.
	// pg: src/backend/commands/sequence.c — init_params (MAXVALUE validation)
	if typeOID != INT8OID {
		if maxValue < bounds[0] || maxValue > bounds[1] {
			return errInvalidParameterValue(fmt.Sprintf("MAXVALUE (%d) is out of range for sequence data type %s",
				maxValue, c.seqTypeName(typeOID)))
		}
	}

	// Validate MINVALUE against type bounds.
	// pg: src/backend/commands/sequence.c — init_params (MINVALUE validation)
	if typeOID != INT8OID {
		if minValue < bounds[0] || minValue > bounds[1] {
			return errInvalidParameterValue(fmt.Sprintf("MINVALUE (%d) is out of range for sequence data type %s",
				minValue, c.seqTypeName(typeOID)))
		}
	}

	if minValue > maxValue {
		return errInvalidParameterValue("MINVALUE must be less than MAXVALUE")
	}

	startValue := minValue
	if !ascending {
		startValue = maxValue
	}
	if pStart != nil {
		startValue = *pStart
	}

	if startValue < minValue || startValue > maxValue {
		return errInvalidParameterValue("START value must be between MINVALUE and MAXVALUE")
	}

	cacheValue := int64(1)
	if pCacheValue != nil {
		cacheValue = *pCacheValue
	}
	if cacheValue <= 0 {
		return errInvalidParameterValue("CACHE must be positive")
	}

	seqOID := c.oidGen.Next()
	seq := &Sequence{
		OID:        seqOID,
		Name:       seqName,
		Schema:     schema,
		TypeOID:    typeOID,
		Start:      startValue,
		Increment:  increment,
		MinValue:   minValue,
		MaxValue:   maxValue,
		CacheValue: cacheValue,
		Cycle:      cycle,
	}

	schema.Sequences[seqName] = seq
	c.sequenceByOID[seqOID] = seq

	// If OwnedBy specified, record dependency.
	if ownedBy != "" {
		if err := c.setSequenceOwner(seq, schema, ownedBy); err != nil {
			delete(schema.Sequences, seqName)
			delete(c.sequenceByOID, seqOID)
			return err
		}
	}

	return nil
}

// AlterSequenceStmt alters an existing sequence from a parsed ALTER SEQUENCE statement.
//
// pg: src/backend/commands/sequence.c — AlterSequence
func (c *Catalog) AlterSequenceStmt(stmt *nodes.AlterSeqStmt) error {
	var schemaName, seqName string
	if stmt.Sequence != nil {
		schemaName = stmt.Sequence.Schemaname
		seqName = stmt.Sequence.Relname
	}

	// If MissingOk is set and the sequence doesn't exist, silently succeed.
	if stmt.MissingOk {
		_, err := c.findSequence(schemaName, seqName)
		if err != nil {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q does not exist, skipping", seqName))
			return nil
		}
	}

	seq, err := c.findSequence(schemaName, seqName)
	if err != nil {
		return err
	}

	// Duplicate DefElem detection.
	// pg: src/backend/commands/sequence.c — init_params (duplicate option checks)
	var sawIncrement, sawMinValue, sawMaxValue, sawStart, sawRestart, sawCache, sawCycle, sawOwnedBy bool

	// Process DefElem options.
	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			d, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch d.Defname {
			case "increment":
				if sawIncrement {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawIncrement = true
				if v, ok := defElemInt(d); ok {
					if v == 0 {
						return errInvalidParameterValue("INCREMENT must not be zero")
					}
					seq.Increment = v
				}
			case "minvalue":
				if sawMinValue {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawMinValue = true
				if v, ok := defElemInt(d); ok {
					seq.MinValue = v
				}
			case "maxvalue":
				if sawMaxValue {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawMaxValue = true
				if v, ok := defElemInt(d); ok {
					seq.MaxValue = v
				}
			case "start":
				if sawStart {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawStart = true
				if v, ok := defElemInt(d); ok {
					seq.Start = v
				}
			case "restart":
				if sawRestart {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawRestart = true
				// RESTART [WITH value]: resets the start position.
				// pg: src/backend/commands/sequence.c — AlterSequence (restart)
				if v, ok := defElemInt(d); ok {
					seq.Start = v
				}
			case "cache":
				if sawCache {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawCache = true
				if v, ok := defElemInt(d); ok {
					if v <= 0 {
						return errInvalidParameterValue("CACHE must be positive")
					}
					seq.CacheValue = v
				}
			case "cycle":
				if sawCycle {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawCycle = true
				seq.Cycle = defElemBool(d)
			case "owned_by":
				if sawOwnedBy {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				sawOwnedBy = true
				if list, ok := d.Arg.(*nodes.List); ok {
					parts := stringListItems(list)
					if len(parts) == 1 && strings.ToUpper(parts[0]) == "NONE" {
						seq.OwnerRelOID = 0
						seq.OwnerAttNum = 0
						c.removeDepsOf('s', seq.OID)
					} else if len(parts) >= 2 {
						ownedBy := parts[len(parts)-2] + "." + parts[len(parts)-1]
						c.removeDepsOf('s', seq.OID)
						if err := c.setSequenceOwner(seq, seq.Schema, ownedBy); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// After processing all options, validate bounds against type.
	// pg: src/backend/commands/sequence.c — init_params (post-processing validation)
	bounds := seqBounds[seq.TypeOID]
	if seq.TypeOID != INT8OID {
		if seq.MaxValue < bounds[0] || seq.MaxValue > bounds[1] {
			return errInvalidParameterValue(fmt.Sprintf("MAXVALUE (%d) is out of range for sequence data type %s",
				seq.MaxValue, c.seqTypeName(seq.TypeOID)))
		}
		if seq.MinValue < bounds[0] || seq.MinValue > bounds[1] {
			return errInvalidParameterValue(fmt.Sprintf("MINVALUE (%d) is out of range for sequence data type %s",
				seq.MinValue, c.seqTypeName(seq.TypeOID)))
		}
	}

	if seq.MinValue >= seq.MaxValue {
		return errInvalidParameterValue("MINVALUE must be less than MAXVALUE")
	}

	// Validate START (or RESTART) value is within [MinValue, MaxValue].
	// pg: src/backend/commands/sequence.c — init_params (START crosscheck)
	if seq.Start < seq.MinValue {
		return errInvalidParameterValue(fmt.Sprintf("START value (%d) cannot be less than MINVALUE (%d)",
			seq.Start, seq.MinValue))
	}
	if seq.Start > seq.MaxValue {
		return errInvalidParameterValue(fmt.Sprintf("START value (%d) cannot be greater than MAXVALUE (%d)",
			seq.Start, seq.MaxValue))
	}

	return nil
}

// findSequence locates a sequence by schema (or search path) and name.
func (c *Catalog) findSequence(schemaName, seqName string) (*Sequence, error) {
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, errUndefinedSchema(schemaName)
		}
		seq := s.Sequences[seqName]
		if seq == nil {
			return nil, errUndefinedSequence(seqName)
		}
		return seq, nil
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		if seq := s.Sequences[seqName]; seq != nil {
			return seq, nil
		}
	}
	return nil, errUndefinedSequence(seqName)
}

func (c *Catalog) findIdentitySequence(rel *Relation, attNum int16) (*Sequence, error) {
	for _, seq := range c.sequenceByOID {
		if seq.OwnerRelOID == rel.OID && seq.OwnerAttNum == attNum {
			return seq, nil
		}
	}
	return nil, &Error{Code: CodeUndefinedObject,
		Message: fmt.Sprintf("identity sequence for column %d of relation %q does not exist", attNum, rel.Name)}
}

// removeSequence removes a sequence from the catalog.
func (c *Catalog) removeSequence(schema *Schema, seq *Sequence) {
	delete(schema.Sequences, seq.Name)
	delete(c.sequenceByOID, seq.OID)
	c.removeComments('s', seq.OID)
	c.removeDepsOf('s', seq.OID)
	c.removeDepsOn('s', seq.OID)
}

// seqTypeName returns the display name for a sequence type OID.
// (pgddl helper — PG uses format_type_be)
func (c *Catalog) seqTypeName(typeOID uint32) string {
	switch typeOID {
	case INT2OID:
		return "smallint"
	case INT4OID:
		return "integer"
	default:
		return "bigint"
	}
}

// resolveSeqType returns the OID for a sequence type name.
//
// pg: src/backend/commands/sequence.c — init_params (AS type resolution)
func (c *Catalog) resolveSeqType(typeName string) (uint32, error) {
	switch strings.ToLower(typeName) {
	case "smallint", "int2":
		return INT2OID, nil
	case "integer", "int", "int4":
		return INT4OID, nil
	case "bigint", "int8", "":
		return INT8OID, nil
	default:
		return 0, errUndefinedType(typeName)
	}
}

// setSequenceOwner sets the owner of a sequence to a table.column reference.
//
// pg: src/backend/commands/sequence.c — process_owned_by
func (c *Catalog) setSequenceOwner(seq *Sequence, schema *Schema, ownedBy string) error {
	parts := strings.SplitN(ownedBy, ".", 2)
	if len(parts) != 2 {
		return errInvalidParameterValue(fmt.Sprintf("invalid OWNED BY specification: %s", ownedBy))
	}
	tableName, colName := parts[0], parts[1]
	rel := schema.Relations[tableName]
	if rel == nil {
		return errUndefinedTable(tableName)
	}

	// Validate that the target relation is a table, view, or foreign table.
	// pg: src/backend/commands/sequence.c — process_owned_by (relkind check)
	switch rel.RelKind {
	case 'r', 'f', 'p', 'v':
		// ordinary table, foreign table, partitioned table, view — OK
	default:
		return errInvalidParameterValue(fmt.Sprintf(
			"sequence cannot be owned by relation \"%s\" of type %s",
			tableName, relKindDescription(rel.RelKind)))
	}

	// Validate sequence and table are in the same schema.
	// pg: src/backend/commands/sequence.c — process_owned_by (line 1640-1648)
	if seq.Schema != nil && rel.Schema != nil && seq.Schema.OID != rel.Schema.OID {
		return &Error{Code: CodeFeatureNotSupported,
			Message: fmt.Sprintf("sequence must be in same schema as table it is linked to")}
	}

	idx, ok := rel.colByName[colName]
	if !ok {
		return errUndefinedColumn(colName)
	}

	seq.OwnerRelOID = rel.OID
	seq.OwnerAttNum = rel.Columns[idx].AttNum
	c.recordDependency('s', seq.OID, 0, 'r', rel.OID, int32(rel.Columns[idx].AttNum), DepAuto)
	return nil
}

type identitySequenceSpec struct {
	schema       *Schema
	name         string
	options      *nodes.List
	explicitName bool
}

func (c *Catalog) createIdentitySequence(tableSchema *Schema, relName, colName string, colTypeOID uint32, seqOptions *nodes.List) (*Sequence, *nodes.RangeVar, bool, error) {
	seqTypeOID, err := identitySequenceTypeOID(colTypeOID)
	if err != nil {
		return nil, nil, false, err
	}

	spec, err := c.identitySequenceSpec(tableSchema, relName, colName, seqOptions)
	if err != nil {
		return nil, nil, false, err
	}
	if spec.schema.OID != tableSchema.OID {
		return nil, nil, false, &Error{Code: CodeFeatureNotSupported,
			Message: "sequence must be in same schema as table it is linked to"}
	}

	createOptions := prependSequenceTypeOption(c.seqTypeName(seqTypeOID), spec.options)
	stmt := &nodes.CreateSeqStmt{
		Sequence:    &nodes.RangeVar{Schemaname: spec.schema.Name, Relname: spec.name},
		Options:     createOptions,
		ForIdentity: true,
	}
	if err := c.DefineSequence(stmt); err != nil {
		return nil, nil, false, err
	}

	seq := spec.schema.Sequences[spec.name]
	return seq, stmt.Sequence, spec.explicitName, nil
}

func (c *Catalog) identitySequenceSpec(tableSchema *Schema, relName, colName string, seqOptions *nodes.List) (identitySequenceSpec, error) {
	spec := identitySequenceSpec{
		schema: tableSchema,
		name:   fmt.Sprintf("%s_%s_seq", relName, colName),
	}
	if seqOptions == nil {
		return spec, nil
	}

	filtered := make([]nodes.Node, 0, len(seqOptions.Items))
	for _, item := range seqOptions.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok {
			filtered = append(filtered, item)
			continue
		}
		switch d.Defname {
		case "sequence_name":
			if spec.explicitName {
				return identitySequenceSpec{}, &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
			}
			parts, ok := d.Arg.(*nodes.List)
			if !ok {
				return identitySequenceSpec{}, &Error{Code: CodeSyntaxError, Message: "invalid sequence name option"}
			}
			names := stringListItems(parts)
			if len(names) == 0 {
				return identitySequenceSpec{}, &Error{Code: CodeSyntaxError, Message: "invalid sequence name option"}
			}
			spec.explicitName = true
			spec.name = names[len(names)-1]
			if len(names) >= 2 {
				seqSchema, err := c.resolveTargetSchema(names[len(names)-2])
				if err != nil {
					return identitySequenceSpec{}, err
				}
				spec.schema = seqSchema
			}
		case "logged":
			// LOGGED/UNLOGGED affects persistence in PG. The catalog model does not
			// track sequence persistence, so remove it like PG does before CREATE SEQUENCE.
			continue
		default:
			filtered = append(filtered, item)
		}
	}
	if len(filtered) > 0 {
		spec.options = &nodes.List{Items: filtered}
	}
	return spec, nil
}

func identitySequenceTypeOID(colTypeOID uint32) (uint32, error) {
	switch colTypeOID {
	case INT2OID, INT4OID, INT8OID:
		return colTypeOID, nil
	default:
		return 0, errInvalidParameterValue("identity column type must be smallint, integer, or bigint")
	}
}

func prependSequenceTypeOption(typeName string, options *nodes.List) *nodes.List {
	items := []nodes.Node{&nodes.DefElem{
		Defname: "as",
		Arg: &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: typeName},
		}}},
		Loc: nodes.NoLoc(),
	}}
	if options != nil {
		items = append(items, options.Items...)
	}
	return &nodes.List{Items: items}
}

func identitySequenceDefault(tableSchema, seqSchema *Schema, seqName string, explicitName bool) string {
	target := seqName
	if explicitName || seqSchema.Name != tableSchema.Name {
		target = fmt.Sprintf("%s.%s", seqSchema.Name, seqName)
	}
	return fmt.Sprintf("nextval('%s'::regclass)", target)
}

// relKindDescription returns a human-readable description of a relation kind.
// (pgddl helper)
func relKindDescription(relkind byte) string {
	switch relkind {
	case 'r':
		return "table"
	case 'v':
		return "view"
	case 'm':
		return "materialized view"
	case 'c':
		return "composite type"
	case 'S':
		return "sequence"
	case 'f':
		return "foreign table"
	case 'p':
		return "partitioned table"
	case 'i':
		return "index"
	case 'I':
		return "partitioned index"
	default:
		return fmt.Sprintf("unknown (%c)", relkind)
	}
}

// createSequenceInternal creates an implicit sequence for SERIAL columns.
func (c *Catalog) createSequenceInternal(schema *Schema, name string, typeOID uint32) *Sequence {
	bounds := seqBounds[typeOID]
	seqOID := c.oidGen.Next()
	seq := &Sequence{
		OID:        seqOID,
		Name:       name,
		Schema:     schema,
		TypeOID:    typeOID,
		Start:      1,
		Increment:  1,
		MinValue:   1,
		MaxValue:   bounds[1],
		CacheValue: 1,
	}
	schema.Sequences[name] = seq
	c.sequenceByOID[seqOID] = seq
	return seq
}

// SequencesOf returns all sequences in the given schema.
func (c *Catalog) SequencesOf(schemaName string) []*Sequence {
	s := c.schemaByName[schemaName]
	if s == nil {
		return nil
	}
	result := make([]*Sequence, 0, len(s.Sequences))
	for _, seq := range s.Sequences {
		result = append(result, seq)
	}
	return result
}

package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// commentKey identifies a comment target.
type commentKey struct {
	ObjType byte
	ObjOID  uint32
	SubID   int16 // attnum for columns, 0 otherwise
}

// CommentObject sets or removes a comment on a catalog object from a pgparser AST.
//
// pg: src/backend/commands/comment.c — CommentObject
func (c *Catalog) CommentObject(stmt *nodes.CommentStmt) error {
	var key commentKey

	switch stmt.Objtype {
	case nodes.OBJECT_TABLE, nodes.OBJECT_VIEW:
		schemaName, name := extractCommentObjectName(stmt.Object)
		_, rel, err := c.findRelation(schemaName, name)
		if err != nil {
			return err
		}
		key = commentKey{ObjType: 'r', ObjOID: rel.OID}

	case nodes.OBJECT_MATVIEW:
		schemaName, name := extractCommentObjectName(stmt.Object)
		_, rel, err := c.findRelation(schemaName, name)
		if err != nil {
			return err
		}
		if rel.RelKind != 'm' {
			return errWrongObjectType(name, "a materialized view")
		}
		key = commentKey{ObjType: 'r', ObjOID: rel.OID}

	case nodes.OBJECT_FOREIGN_TABLE:
		schemaName, name := extractCommentObjectName(stmt.Object)
		_, rel, err := c.findRelation(schemaName, name)
		if err != nil {
			return err
		}
		if rel.RelKind != 'f' {
			return errWrongObjectType(name, "a foreign table")
		}
		key = commentKey{ObjType: 'r', ObjOID: rel.OID}

	case nodes.OBJECT_COLUMN:
		var schemaName, tableName, colName string
		if list, ok := stmt.Object.(*nodes.List); ok {
			items := list.Items
			switch len(items) {
			case 2:
				tableName = stringVal(items[0])
				colName = stringVal(items[1])
			case 3:
				schemaName = stringVal(items[0])
				tableName = stringVal(items[1])
				colName = stringVal(items[2])
			}
		}
		_, rel, err := c.findRelation(schemaName, tableName)
		if err != nil {
			return err
		}
		// Validate relkind allows column comments.
		// pg: src/backend/catalog/objectaddress.c — get_relation_by_qualified_name
		if rel.RelKind != 'r' && rel.RelKind != 'v' && rel.RelKind != 'm' &&
			rel.RelKind != 'c' && rel.RelKind != 'f' && rel.RelKind != 'p' {
			return errWrongObjectType(tableName, "a table, view, or composite type")
		}
		idx, ok := rel.colByName[colName]
		if !ok {
			return errUndefinedColumn(colName)
		}
		key = commentKey{ObjType: 'r', ObjOID: rel.OID, SubID: rel.Columns[idx].AttNum}

	case nodes.OBJECT_INDEX:
		schemaName, name := extractCommentObjectName(stmt.Object)
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
		idx, ok := schema.Indexes[name]
		if !ok {
			return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("index %q does not exist", name)}
		}
		key = commentKey{ObjType: 'i', ObjOID: idx.OID}

	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE, nodes.OBJECT_AGGREGATE:
		// pg: src/backend/commands/comment.c — CommentObject (OBJECT_FUNCTION etc.)
		// PROCEDURE, ROUTINE, and AGGREGATE all resolve via the same function lookup.
		if owa, ok := stmt.Object.(*nodes.ObjectWithArgs); ok {
			schemaName, name := qualifiedName(owa.Objname)
			schema, err := c.resolveTargetSchema(schemaName)
			if err != nil {
				return err
			}
			var argOIDs []uint32
			if owa.Objargs != nil {
				for _, item := range owa.Objargs.Items {
					if tn, ok := item.(*nodes.TypeName); ok {
						at := convertTypeNameToInternal(tn)
						oid, _, err := c.ResolveType(at)
						if err != nil {
							return err
						}
						argOIDs = append(argOIDs, oid)
					}
				}
			}
			bp := c.findExactProc(schema, name, argOIDs)
			if bp == nil {
				return errUndefinedFunction(name, argOIDs)
			}
			key = commentKey{ObjType: 'f', ObjOID: bp.OID}
		}

	case nodes.OBJECT_SCHEMA:
		name := extractSimpleObjectName(stmt.Object)
		s := c.schemaByName[name]
		if s == nil {
			return errUndefinedSchema(name)
		}
		key = commentKey{ObjType: 'n', ObjOID: s.OID}

	case nodes.OBJECT_TYPE:
		schemaName, name := extractCommentObjectName(stmt.Object)
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
		bt := c.typeByName[typeKey{ns: schema.OID, name: name}]
		if bt == nil {
			return errUndefinedType(name)
		}
		key = commentKey{ObjType: 't', ObjOID: bt.OID}

	case nodes.OBJECT_SEQUENCE:
		schemaName, name := extractCommentObjectName(stmt.Object)
		seq, err := c.findSequence(schemaName, name)
		if err != nil {
			return err
		}
		key = commentKey{ObjType: 's', ObjOID: seq.OID}

	case nodes.OBJECT_TABCONSTRAINT:
		var schemaName, tableName, conName string
		if list, ok := stmt.Object.(*nodes.List); ok {
			items := list.Items
			switch len(items) {
			case 2:
				tableName = stringVal(items[0])
				conName = stringVal(items[1])
			case 3:
				schemaName = stringVal(items[0])
				tableName = stringVal(items[1])
				conName = stringVal(items[2])
			}
		}
		_, rel, err := c.findRelation(schemaName, tableName)
		if err != nil {
			return err
		}
		var con *Constraint
		for _, cc := range c.consByRel[rel.OID] {
			if cc.Name == conName {
				con = cc
				break
			}
		}
		if con == nil {
			return errUndefinedObject("constraint", conName)
		}
		key = commentKey{ObjType: 'c', ObjOID: con.OID}

	case nodes.OBJECT_POLICY:
		var schemaName, tableName, policyName string
		if list, ok := stmt.Object.(*nodes.List); ok {
			items := list.Items
			switch len(items) {
			case 2:
				tableName = stringVal(items[0])
				policyName = stringVal(items[1])
			case 3:
				schemaName = stringVal(items[0])
				tableName = stringVal(items[1])
				policyName = stringVal(items[2])
			}
		}
		_, rel, err := c.findRelation(schemaName, tableName)
		if err != nil {
			return err
		}
		var policy *Policy
		for _, p := range c.policiesByRel[rel.OID] {
			if p.Name == policyName {
				policy = p
				break
			}
		}
		if policy == nil {
			return errUndefinedObject("policy", policyName)
		}
		key = commentKey{ObjType: 'p', ObjOID: policy.OID}

	case nodes.OBJECT_DOMCONSTRAINT:
		// COMMENT ON CONSTRAINT name ON DOMAIN type
		// pg: src/backend/commands/comment.c — CommentObject (OBJECT_DOMCONSTRAINT)
		var schemaName, domainName, conName string
		if list, ok := stmt.Object.(*nodes.List); ok {
			items := list.Items
			switch len(items) {
			case 2:
				domainName = stringVal(items[0])
				conName = stringVal(items[1])
			case 3:
				schemaName = stringVal(items[0])
				domainName = stringVal(items[1])
				conName = stringVal(items[2])
			}
		}
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
		bt := c.typeByName[typeKey{ns: schema.OID, name: domainName}]
		if bt == nil {
			return errUndefinedType(domainName)
		}
		if bt.Type != 'd' {
			return errWrongObjectType(domainName, "a domain")
		}
		dt := c.domainTypes[bt.OID]
		if dt == nil {
			return errUndefinedType(domainName)
		}
		var foundDC *DomainConstraint
		for _, dc := range dt.Constraints {
			if dc.Name == conName {
				foundDC = dc
				break
			}
		}
		if foundDC == nil {
			return errUndefinedObject("constraint", conName)
		}
		key = commentKey{ObjType: 'd', ObjOID: foundDC.OID}

	case nodes.OBJECT_TRIGGER:
		var schemaName, tableName, trigName string
		if list, ok := stmt.Object.(*nodes.List); ok {
			items := list.Items
			switch len(items) {
			case 2:
				tableName = stringVal(items[0])
				trigName = stringVal(items[1])
			case 3:
				schemaName = stringVal(items[0])
				tableName = stringVal(items[1])
				trigName = stringVal(items[2])
			}
		}
		_, rel, err := c.findRelation(schemaName, tableName)
		if err != nil {
			return err
		}
		var found bool
		for _, trig := range c.triggersByRel[rel.OID] {
			if trig.Name == trigName {
				key = commentKey{ObjType: 'g', ObjOID: trig.OID}
				found = true
				break
			}
		}
		if !found {
			return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("trigger %q does not exist", trigName)}
		}

	case nodes.OBJECT_EVENT_TRIGGER:
		name := extractSimpleObjectName(stmt.Object)
		evt := c.eventTriggerByName[name]
		if evt == nil {
			return errUndefinedObject("event trigger", name)
		}
		key = commentKey{ObjType: 'E', ObjOID: evt.OID}

	case nodes.OBJECT_COLLATION,
		nodes.OBJECT_CONVERSION,
		nodes.OBJECT_OPERATOR,
		nodes.OBJECT_OPCLASS,
		nodes.OBJECT_OPFAMILY,
		nodes.OBJECT_LANGUAGE,
		nodes.OBJECT_FDW,
		nodes.OBJECT_FOREIGN_SERVER,
		nodes.OBJECT_RULE,
		nodes.OBJECT_CAST,
		nodes.OBJECT_TRANSFORM,
		nodes.OBJECT_STATISTIC_EXT,
		nodes.OBJECT_TSPARSER,
		nodes.OBJECT_TSDICTIONARY,
		nodes.OBJECT_TSTEMPLATE,
		nodes.OBJECT_TSCONFIGURATION,
		nodes.OBJECT_EXTENSION,
		nodes.OBJECT_ACCESS_METHOD,
		nodes.OBJECT_PUBLICATION,
		nodes.OBJECT_DATABASE,
		nodes.OBJECT_TABLESPACE,
		nodes.OBJECT_ROLE,
		nodes.OBJECT_LARGEOBJECT:
		// No-op: these object types are not tracked by pgddl.
		return nil

	default:
		return fmt.Errorf("unsupported comment object type %d", stmt.Objtype)
	}

	if stmt.Comment == "" {
		delete(c.comments, key)
	} else {
		c.comments[key] = stmt.Comment
	}
	return nil
}

// extractCommentObjectName extracts schema and name from a comment Object node.
// The Object can be a *nodes.List of String nodes or a *nodes.String.
func extractCommentObjectName(obj nodes.Node) (schema, name string) {
	switch n := obj.(type) {
	case *nodes.List:
		return qualifiedName(n)
	case *nodes.String:
		return "", n.Str
	default:
		return "", ""
	}
}

// extractSimpleObjectName extracts a simple name from a comment Object node.
func extractSimpleObjectName(obj nodes.Node) string {
	switch n := obj.(type) {
	case *nodes.String:
		return n.Str
	case *nodes.List:
		if len(n.Items) > 0 {
			return stringVal(n.Items[len(n.Items)-1])
		}
	}
	return ""
}

// GetComment returns the comment for the given object, if any.
func (c *Catalog) GetComment(objType byte, objOID uint32, subID int16) (string, bool) {
	s, ok := c.comments[commentKey{ObjType: objType, ObjOID: objOID, SubID: subID}]
	return s, ok
}

// removeComments deletes all comments for the given object (including sub-objects like columns).
func (c *Catalog) removeComments(objType byte, objOID uint32) {
	for k := range c.comments {
		if k.ObjType == objType && k.ObjOID == objOID {
			delete(c.comments, k)
		}
	}
}

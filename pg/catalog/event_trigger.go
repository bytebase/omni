package catalog

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// EventTrigger represents a PostgreSQL event trigger in pg_event_trigger.
type EventTrigger struct {
	OID       uint32
	Name      string
	EventName string
	FuncOID   uint32
	Enabled   byte
	Tags      []string
}

// CreateEventTriggerStmt creates a database-level event trigger.
//
// pg: src/backend/commands/event_trigger.c — CreateEventTrigger
func (c *Catalog) CreateEventTriggerStmt(stmt *nodes.CreateEventTrigStmt) error {
	if !validEventTriggerEvent(stmt.Eventname) {
		return &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("unrecognized event name %q", stmt.Eventname)}
	}
	tags, err := eventTriggerTags(stmt.Eventname, stmt.Whenclause)
	if err != nil {
		return err
	}
	if c.eventTriggerByName[stmt.Trigname] != nil {
		return errDuplicateObject("event trigger", stmt.Trigname)
	}

	funcSchema, funcName := qualifiedName(stmt.Funcname)
	funcBP, wrongRet := c.findEventTriggerFunc(funcSchema, funcName)
	if funcBP == nil {
		if wrongRet {
			return errInvalidObjectDefinition(fmt.Sprintf("function %s must return type event_trigger", funcName))
		}
		return errUndefinedFunction(funcName, nil)
	}

	oid := c.oidGen.Next()
	evt := &EventTrigger{
		OID:       oid,
		Name:      stmt.Trigname,
		EventName: stmt.Eventname,
		FuncOID:   funcBP.OID,
		Enabled:   'O',
		Tags:      tags,
	}
	c.eventTriggers[oid] = evt
	c.eventTriggerByName[evt.Name] = evt

	// pg: event_trigger.c — event trigger depends on its function.
	c.recordDependency('E', oid, 0, 'f', funcBP.OID, 0, DepNormal)
	return nil
}

// AlterEventTrigger updates an event trigger's firing configuration.
//
// pg: src/backend/commands/event_trigger.c — AlterEventTrigger
func (c *Catalog) AlterEventTrigger(stmt *nodes.AlterEventTrigStmt) error {
	evt := c.eventTriggerByName[stmt.Trigname]
	if evt == nil {
		return errUndefinedObject("event trigger", stmt.Trigname)
	}
	evt.Enabled = stmt.Tgenabled
	return nil
}

func (c *Catalog) renameEventTrigger(stmt *nodes.RenameStmt) error {
	oldName := extractSimpleObjectName(stmt.Object)
	if oldName == "" {
		oldName = stmt.Subname
	}
	evt := c.eventTriggerByName[oldName]
	if evt == nil {
		return errUndefinedObject("event trigger", oldName)
	}
	if c.eventTriggerByName[stmt.Newname] != nil {
		return errDuplicateObject("event trigger", stmt.Newname)
	}
	delete(c.eventTriggerByName, oldName)
	evt.Name = stmt.Newname
	c.eventTriggerByName[evt.Name] = evt
	return nil
}

func (c *Catalog) removeEventTriggerObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		_, name := extractDropObjectName(obj)
		evt := c.eventTriggerByName[name]
		if evt == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("event trigger %q does not exist, skipping", name))
				continue
			}
			return errUndefinedObject("event trigger", name)
		}
		c.removeEventTrigger(evt)
	}
	return nil
}

func (c *Catalog) removeEventTrigger(evt *EventTrigger) {
	delete(c.eventTriggers, evt.OID)
	delete(c.eventTriggerByName, evt.Name)
	c.removeComments('E', evt.OID)
	c.removeDepsOf('E', evt.OID)
	c.removeDepsOn('E', evt.OID)
}

func validEventTriggerEvent(event string) bool {
	switch event {
	case "ddl_command_start", "ddl_command_end", "sql_drop", "login", "table_rewrite":
		return true
	default:
		return false
	}
}

func eventTriggerTags(event string, when *nodes.List) ([]string, error) {
	var tags []string
	var items []nodes.Node
	if when != nil {
		items = when.Items
	}
	for _, item := range items {
		def, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		if def.Defname != "tag" {
			return nil, &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("unrecognized filter variable %q", def.Defname)}
		}
		if tags != nil {
			return nil, &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("filter variable %q specified more than once", def.Defname)}
		}
		l, ok := def.Arg.(*nodes.List)
		if !ok {
			continue
		}
		for _, tagItem := range l.Items {
			tag := strings.ToUpper(stringVal(tagItem))
			if tag == "" {
				continue
			}
			tags = append(tags, tag)
		}
	}

	if len(tags) == 0 {
		return nil, nil
	}
	switch event {
	case "login":
		return nil, &Error{Code: CodeFeatureNotSupported, Message: "tag filtering is not supported for login event triggers"}
	case "table_rewrite":
		for _, tag := range tags {
			if !eventTriggerTableRewriteTagOK(tag) {
				return nil, &Error{Code: CodeFeatureNotSupported, Message: fmt.Sprintf("event triggers are not supported for %s", tag)}
			}
		}
	default:
		for _, tag := range tags {
			if !eventTriggerKnownDDLTags[tag] && !eventTriggerUnsupportedDDLTags[tag] {
				return nil, &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("filter value %q not recognized for filter variable %q", tag, "tag")}
			}
			if !eventTriggerDDLTagOK(tag) {
				return nil, &Error{Code: CodeFeatureNotSupported, Message: fmt.Sprintf("event triggers are not supported for %s", tag)}
			}
		}
	}
	return tags, nil
}

func eventTriggerDDLTagOK(tag string) bool {
	if tag == "" {
		return false
	}
	if eventTriggerUnsupportedDDLTags[tag] {
		return false
	}
	return eventTriggerKnownDDLTags[tag]
}

func eventTriggerTableRewriteTagOK(tag string) bool {
	return eventTriggerTableRewriteTags[tag]
}

var eventTriggerUnsupportedDDLTags = map[string]bool{
	"ALTER DATABASE":           true,
	"ALTER EVENT TRIGGER":      true,
	"ALTER ROLE":               true,
	"ALTER SYSTEM":             true,
	"ALTER TABLESPACE":         true,
	"ANALYZE":                  true,
	"BEGIN":                    true,
	"CALL":                     true,
	"CHECKPOINT":               true,
	"CLOSE":                    true,
	"CLOSE CURSOR":             true,
	"CLOSE CURSOR ALL":         true,
	"CLUSTER":                  true,
	"COMMIT":                   true,
	"COMMIT PREPARED":          true,
	"COPY":                     true,
	"COPY FROM":                true,
	"CREATE DATABASE":          true,
	"CREATE EVENT TRIGGER":     true,
	"CREATE ROLE":              true,
	"CREATE TABLESPACE":        true,
	"DEALLOCATE":               true,
	"DEALLOCATE ALL":           true,
	"DECLARE CURSOR":           true,
	"DELETE":                   true,
	"DISCARD":                  true,
	"DISCARD ALL":              true,
	"DISCARD PLANS":            true,
	"DISCARD SEQUENCES":        true,
	"DISCARD TEMP":             true,
	"DO":                       true,
	"DROP DATABASE":            true,
	"DROP EVENT TRIGGER":       true,
	"DROP ROLE":                true,
	"DROP TABLESPACE":          true,
	"EXECUTE":                  true,
	"EXPLAIN":                  true,
	"FETCH":                    true,
	"GRANT ROLE":               true,
	"INSERT":                   true,
	"LISTEN":                   true,
	"LOAD":                     true,
	"LOCK TABLE":               true,
	"MERGE":                    true,
	"MOVE":                     true,
	"NOTIFY":                   true,
	"PREPARE":                  true,
	"PREPARE TRANSACTION":      true,
	"REASSIGN OWNED":           true,
	"RELEASE":                  true,
	"RESET":                    true,
	"REVOKE ROLE":              true,
	"ROLLBACK":                 true,
	"ROLLBACK PREPARED":        true,
	"SAVEPOINT":                true,
	"SELECT":                   true,
	"SELECT FOR KEY SHARE":     true,
	"SELECT FOR NO KEY UPDATE": true,
	"SELECT FOR SHARE":         true,
	"SELECT FOR UPDATE":        true,
	"SET":                      true,
	"SET CONSTRAINTS":          true,
	"SHOW":                     true,
	"START TRANSACTION":        true,
	"TRUNCATE TABLE":           true,
	"UNLISTEN":                 true,
	"UPDATE":                   true,
	"VACUUM":                   true,
	"VALUES":                   true,
}

var eventTriggerKnownDDLTags = map[string]bool{
	"ALTER ACCESS METHOD":              true,
	"ALTER AGGREGATE":                  true,
	"ALTER CAST":                       true,
	"ALTER COLLATION":                  true,
	"ALTER CONSTRAINT":                 true,
	"ALTER CONVERSION":                 true,
	"ALTER DEFAULT PRIVILEGES":         true,
	"ALTER DOMAIN":                     true,
	"ALTER EXTENSION":                  true,
	"ALTER FOREIGN DATA WRAPPER":       true,
	"ALTER FOREIGN TABLE":              true,
	"ALTER FUNCTION":                   true,
	"ALTER INDEX":                      true,
	"ALTER LANGUAGE":                   true,
	"ALTER LARGE OBJECT":               true,
	"ALTER MATERIALIZED VIEW":          true,
	"ALTER OPERATOR":                   true,
	"ALTER OPERATOR CLASS":             true,
	"ALTER OPERATOR FAMILY":            true,
	"ALTER POLICY":                     true,
	"ALTER PROCEDURE":                  true,
	"ALTER PUBLICATION":                true,
	"ALTER ROUTINE":                    true,
	"ALTER RULE":                       true,
	"ALTER SCHEMA":                     true,
	"ALTER SEQUENCE":                   true,
	"ALTER SERVER":                     true,
	"ALTER STATISTICS":                 true,
	"ALTER SUBSCRIPTION":               true,
	"ALTER TABLE":                      true,
	"ALTER TEXT SEARCH CONFIGURATION":  true,
	"ALTER TEXT SEARCH DICTIONARY":     true,
	"ALTER TEXT SEARCH PARSER":         true,
	"ALTER TEXT SEARCH TEMPLATE":       true,
	"ALTER TRANSFORM":                  true,
	"ALTER TRIGGER":                    true,
	"ALTER TYPE":                       true,
	"ALTER USER MAPPING":               true,
	"ALTER VIEW":                       true,
	"COMMENT":                          true,
	"CREATE ACCESS METHOD":             true,
	"CREATE AGGREGATE":                 true,
	"CREATE CAST":                      true,
	"CREATE COLLATION":                 true,
	"CREATE CONSTRAINT":                true,
	"CREATE CONVERSION":                true,
	"CREATE DOMAIN":                    true,
	"CREATE EXTENSION":                 true,
	"CREATE FOREIGN DATA WRAPPER":      true,
	"CREATE FOREIGN TABLE":             true,
	"CREATE FUNCTION":                  true,
	"CREATE INDEX":                     true,
	"CREATE LANGUAGE":                  true,
	"CREATE MATERIALIZED VIEW":         true,
	"CREATE OPERATOR":                  true,
	"CREATE OPERATOR CLASS":            true,
	"CREATE OPERATOR FAMILY":           true,
	"CREATE POLICY":                    true,
	"CREATE PROCEDURE":                 true,
	"CREATE PUBLICATION":               true,
	"CREATE ROUTINE":                   true,
	"CREATE RULE":                      true,
	"CREATE SCHEMA":                    true,
	"CREATE SEQUENCE":                  true,
	"CREATE SERVER":                    true,
	"CREATE STATISTICS":                true,
	"CREATE SUBSCRIPTION":              true,
	"CREATE TABLE":                     true,
	"CREATE TABLE AS":                  true,
	"CREATE TEXT SEARCH CONFIGURATION": true,
	"CREATE TEXT SEARCH DICTIONARY":    true,
	"CREATE TEXT SEARCH PARSER":        true,
	"CREATE TEXT SEARCH TEMPLATE":      true,
	"CREATE TRANSFORM":                 true,
	"CREATE TRIGGER":                   true,
	"CREATE TYPE":                      true,
	"CREATE USER MAPPING":              true,
	"CREATE VIEW":                      true,
	"DROP ACCESS METHOD":               true,
	"DROP AGGREGATE":                   true,
	"DROP CAST":                        true,
	"DROP COLLATION":                   true,
	"DROP CONSTRAINT":                  true,
	"DROP CONVERSION":                  true,
	"DROP DOMAIN":                      true,
	"DROP EXTENSION":                   true,
	"DROP FOREIGN DATA WRAPPER":        true,
	"DROP FOREIGN TABLE":               true,
	"DROP FUNCTION":                    true,
	"DROP INDEX":                       true,
	"DROP LANGUAGE":                    true,
	"DROP MATERIALIZED VIEW":           true,
	"DROP OPERATOR":                    true,
	"DROP OPERATOR CLASS":              true,
	"DROP OPERATOR FAMILY":             true,
	"DROP OWNED":                       true,
	"DROP POLICY":                      true,
	"DROP PROCEDURE":                   true,
	"DROP PUBLICATION":                 true,
	"DROP ROUTINE":                     true,
	"DROP RULE":                        true,
	"DROP SCHEMA":                      true,
	"DROP SEQUENCE":                    true,
	"DROP SERVER":                      true,
	"DROP STATISTICS":                  true,
	"DROP SUBSCRIPTION":                true,
	"DROP TABLE":                       true,
	"DROP TEXT SEARCH CONFIGURATION":   true,
	"DROP TEXT SEARCH DICTIONARY":      true,
	"DROP TEXT SEARCH PARSER":          true,
	"DROP TEXT SEARCH TEMPLATE":        true,
	"DROP TRANSFORM":                   true,
	"DROP TRIGGER":                     true,
	"DROP TYPE":                        true,
	"DROP USER MAPPING":                true,
	"DROP VIEW":                        true,
	"GRANT":                            true,
	"IMPORT FOREIGN SCHEMA":            true,
	"LOGIN":                            true,
	"REFRESH MATERIALIZED VIEW":        true,
	"REINDEX":                          true,
	"REVOKE":                           true,
	"SECURITY LABEL":                   true,
	"SELECT INTO":                      true,
}

var eventTriggerTableRewriteTags = map[string]bool{
	"ALTER MATERIALIZED VIEW": true,
	"ALTER TABLE":             true,
	"ALTER TYPE":              true,
}

func (c *Catalog) findEventTriggerFunc(schemaName, funcName string) (*BuiltinProc, bool) {
	foundAny := false
	for _, p := range c.procByName[funcName] {
		if p.NArgs != 0 {
			continue
		}
		if schemaName != "" {
			up := c.userProcs[p.OID]
			if up == nil || c.schemaByName[schemaName] == nil || up.Schema.OID != c.schemaByName[schemaName].OID {
				continue
			}
		}
		foundAny = true
		if p.RetType == EVENTTRIGGEROID {
			return p, false
		}
	}
	return nil, foundAny
}

func diffEventTriggers(from, to *Catalog) []EventTriggerDiffEntry {
	fromMap := make(map[string]*EventTrigger)
	for _, evt := range from.eventTriggers {
		fromMap[evt.Name] = evt
	}
	toMap := make(map[string]*EventTrigger)
	for _, evt := range to.eventTriggers {
		toMap[evt.Name] = evt
	}

	var result []EventTriggerDiffEntry
	for name, fromEvt := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, EventTriggerDiffEntry{Action: DiffDrop, Name: name, From: fromEvt})
		}
	}
	for name, toEvt := range toMap {
		fromEvt, ok := fromMap[name]
		if !ok {
			result = append(result, EventTriggerDiffEntry{Action: DiffAdd, Name: name, To: toEvt})
			continue
		}
		if eventTriggersChanged(from, to, fromEvt, toEvt) {
			result = append(result, EventTriggerDiffEntry{Action: DiffModify, Name: name, From: fromEvt, To: toEvt})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})
	return result
}

func eventTriggersChanged(from, to *Catalog, a, b *EventTrigger) bool {
	if a.EventName != b.EventName || a.Enabled != b.Enabled {
		return true
	}
	if resolveEventTriggerFuncIdentity(from, a.FuncOID) != resolveEventTriggerFuncIdentity(to, b.FuncOID) {
		return true
	}
	return !stringSliceEqual(a.Tags, b.Tags)
}

func resolveEventTriggerFuncIdentity(c *Catalog, funcOID uint32) string {
	up := c.userProcs[funcOID]
	if up != nil {
		return funcIdentity(c, up)
	}
	p := c.procByOID[funcOID]
	if p == nil {
		return ""
	}
	return p.Name + "()"
}

func resolveEventTriggerFuncSQLName(c *Catalog, funcOID uint32) string {
	up := c.userProcs[funcOID]
	if up != nil && up.Schema != nil {
		return migrationQualifiedName(up.Schema.Name, up.Name)
	}
	p := c.procByOID[funcOID]
	if p == nil {
		return quoteIdentAlways("unknown_function")
	}
	return quoteIdentAlways(p.Name)
}

func structuralEventTriggerChange(from, to *Catalog, a, b *EventTrigger) bool {
	if a == nil || b == nil {
		return true
	}
	if a.EventName != b.EventName {
		return true
	}
	if resolveEventTriggerFuncIdentity(from, a.FuncOID) != resolveEventTriggerFuncIdentity(to, b.FuncOID) {
		return true
	}
	return !stringSliceEqual(a.Tags, b.Tags)
}

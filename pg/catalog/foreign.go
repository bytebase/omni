package catalog

import (
	"maps"
	"slices"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ForeignDataWrapper represents a PostgreSQL foreign-data wrapper (pg_foreign_data_wrapper).
type ForeignDataWrapper struct {
	OID          uint32
	Name         string
	OwnerName    string
	Handler      string
	HandlerOID   uint32
	Validator    string
	ValidatorOID uint32
	Options      []string
}

// ForeignServer represents a PostgreSQL foreign server (pg_foreign_server).
type ForeignServer struct {
	OID       uint32
	Name      string
	OwnerName string
	FdwName   string
	Type      *string
	Version   *string
	Options   []string
}

// ForeignDataWrapperDiffEntry describes a foreign-data wrapper change.
type ForeignDataWrapperDiffEntry struct {
	Action DiffAction
	Name   string
	From   *ForeignDataWrapper
	To     *ForeignDataWrapper
}

// ForeignServerDiffEntry describes a foreign server change.
type ForeignServerDiffEntry struct {
	Action DiffAction
	Name   string
	From   *ForeignServer
	To     *ForeignServer
}

func cloneStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}

func stringPtrEqual(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (c *Catalog) resolveFdwFuncOptions(funcOpts *nodes.List) (handler, validator string, handlerOID, validatorOID uint32, hasHandler, hasValidator bool, err error) {
	if funcOpts == nil {
		return "", "", 0, 0, false, false, nil
	}
	for _, item := range funcOpts.Items {
		de, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		switch strings.ToLower(de.Defname) {
		case "handler":
			if hasHandler {
				return "", "", 0, 0, false, false, errInvalidParameterValue("conflicting or redundant options")
			}
			handler, handlerOID, err = c.resolveRoutineRef(de, nil, FDWHANDLEROID, "fdw_handler")
			if err != nil {
				return "", "", 0, 0, false, false, err
			}
			hasHandler = true
		case "validator":
			if hasValidator {
				return "", "", 0, 0, false, false, errInvalidParameterValue("conflicting or redundant options")
			}
			validator, validatorOID, err = c.resolveRoutineRef(de, []uint32{TEXTARRAYOID, OIDOID}, VOIDOID, "void")
			if err != nil {
				return "", "", 0, 0, false, false, err
			}
			hasValidator = true
		}
	}
	return handler, validator, handlerOID, validatorOID, hasHandler, hasValidator, nil
}

func splitDefElemQualifiedName(de *nodes.DefElem) (schemaName, objName string) {
	if de == nil || de.Arg == nil {
		return "", ""
	}
	switch v := de.Arg.(type) {
	case *nodes.List:
		return qualifiedName(v)
	case *nodes.TypeName:
		if v.Names != nil {
			return qualifiedName(v.Names)
		}
	case *nodes.String:
		return "", v.Str
	}
	return "", ""
}

func (c *Catalog) resolveRoutineRef(de *nodes.DefElem, expectedArgs []uint32, expectedRet uint32, retTypeName string) (string, uint32, error) {
	formatted := formatCatalogOptionName(de)
	if formatted == "" {
		return "", 0, nil
	}
	schemaName, funcName := splitDefElemQualifiedName(de)
	if funcName == "" {
		return "", 0, errUndefinedFunction(formatted, expectedArgs)
	}
	oid, err := c.resolveRoutineRefByName(schemaName, funcName, formatted, expectedArgs, expectedRet, retTypeName)
	if err != nil {
		return "", 0, err
	}
	return formatted, oid, nil
}

func (c *Catalog) resolveRoutineRefByName(schemaName, funcName, displayName string, expectedArgs []uint32, expectedRet uint32, retTypeName string) (uint32, error) {
	foundWrongRet := false
	if schemaName != "" {
		ns := c.schemaByName[schemaName]
		if ns != nil || schemaName == "pg_catalog" {
			for _, p := range c.procByName[funcName] {
				up := c.userProcs[p.OID]
				if schemaName == "pg_catalog" {
					if up != nil {
						continue
					}
				} else if up == nil || up.Schema == nil || up.Schema.OID != ns.OID {
					continue
				}
				if !slices.Equal(p.ArgTypes, expectedArgs) {
					continue
				}
				if p.Kind == 'p' || p.Kind == 'a' || p.Kind == 'w' {
					continue
				}
				if expectedRet == 0 || p.RetType == expectedRet {
					return p.OID, nil
				}
				foundWrongRet = true
			}
		}
	} else {
		for _, nsOID := range c.routineSearchPathWithCatalog() {
			matchedInSchema := false
			for _, p := range c.procByName[funcName] {
				up := c.userProcs[p.OID]
				procNS := uint32(PGCatalogNamespace)
				if up != nil && up.Schema != nil {
					procNS = up.Schema.OID
				}
				if procNS != nsOID || !slices.Equal(p.ArgTypes, expectedArgs) {
					continue
				}
				matchedInSchema = true
				if p.Kind == 'p' || p.Kind == 'a' || p.Kind == 'w' {
					continue
				}
				if expectedRet == 0 || p.RetType == expectedRet {
					return p.OID, nil
				}
				foundWrongRet = true
			}
			if matchedInSchema {
				break
			}
		}
	}
	if foundWrongRet {
		return 0, errInvalidObjectDefinition("function " + displayName + " must return type " + retTypeName)
	}
	return 0, errUndefinedFunction(funcName, expectedArgs)
}

func parseRoutineRefName(ref string) (schemaName, funcName string) {
	parts := splitQualifiedIdentifiers(ref)
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return "", parts[0]
	default:
		return parts[len(parts)-2], parts[len(parts)-1]
	}
}

func splitQualifiedIdentifiers(raw string) []string {
	var parts []string
	var cur strings.Builder
	inQuote := false
	wasQuoted := false
	flush := func() {
		s := strings.TrimSpace(cur.String())
		cur.Reset()
		if s != "" || wasQuoted {
			parts = append(parts, s)
		}
		wasQuoted = false
	}
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if inQuote {
			if ch == '"' {
				if i+1 < len(raw) && raw[i+1] == '"' {
					cur.WriteByte('"')
					i++
					continue
				}
				inQuote = false
				continue
			}
			cur.WriteByte(ch)
			continue
		}
		if ch == '"' {
			inQuote = true
			wasQuoted = true
			continue
		}
		if ch == '.' {
			flush()
			continue
		}
		cur.WriteByte(ch)
	}
	flush()
	return parts
}

func (c *Catalog) syncFdwRoutineDeps(fdw *ForeignDataWrapper) {
	if fdw == nil {
		return
	}
	c.removeDepsOf('F', fdw.OID)
	if fdw.HandlerOID >= FirstNormalObjectId {
		c.recordDependency('F', fdw.OID, 0, 'f', fdw.HandlerOID, 0, DepNormal)
	}
	if fdw.ValidatorOID >= FirstNormalObjectId {
		c.recordDependency('F', fdw.OID, 0, 'f', fdw.ValidatorOID, 0, DepNormal)
	}
}

func (c *Catalog) linkForeignWrapperFuncDeps() {
	for _, fdw := range c.fdws {
		if fdw == nil {
			continue
		}
		if fdw.Handler != "" && fdw.HandlerOID == 0 {
			s, fn := parseRoutineRefName(fdw.Handler)
			if oid, err := c.resolveRoutineRefByName(s, fn, fdw.Handler, nil, FDWHANDLEROID, "fdw_handler"); err == nil {
				fdw.HandlerOID = oid
			}
		}
		if fdw.Validator != "" && fdw.ValidatorOID == 0 {
			s, fn := parseRoutineRefName(fdw.Validator)
			if oid, err := c.resolveRoutineRefByName(s, fn, fdw.Validator, []uint32{TEXTARRAYOID, OIDOID}, VOIDOID, "void"); err == nil {
				fdw.ValidatorOID = oid
			}
		}
		c.syncFdwRoutineDeps(fdw)
	}
}

func formatUpdatedRoutineRef(existing, schemaName, newName string, schemaChanged bool) string {
	if existing == "" {
		return ""
	}
	if s, _ := parseRoutineRefName(existing); s != "" || schemaChanged {
		if s != "" && !schemaChanged {
			schemaName = s
		}
		if schemaName != "" {
			return quoteIdentifier(schemaName) + "." + quoteIdentifier(newName)
		}
	}
	return quoteIdentifier(newName)
}

func (c *Catalog) updateFdwRoutineRefs(procOID uint32, schemaName, newName string, schemaChanged bool) {
	if procOID == 0 || newName == "" {
		return
	}
	for _, fdw := range c.fdws {
		if fdw == nil {
			continue
		}
		if fdw.HandlerOID == procOID {
			fdw.Handler = formatUpdatedRoutineRef(fdw.Handler, schemaName, newName, schemaChanged)
		}
		if fdw.ValidatorOID == procOID {
			fdw.Validator = formatUpdatedRoutineRef(fdw.Validator, schemaName, newName, schemaChanged)
		}
	}
}

func (c *Catalog) updateSupportRoutineRefs(procOID uint32, schemaName, newName string, schemaChanged bool) {
	if procOID == 0 || newName == "" {
		return
	}
	for _, up := range c.userProcs {
		if up == nil || up.SupportOID != procOID {
			continue
		}
		up.Support = formatUpdatedRoutineRef(up.Support, schemaName, newName, schemaChanged)
	}
	for _, rt := range c.rangeTypes {
		if rt == nil {
			continue
		}
		if rt.CanonicalOID == procOID {
			rt.Canonical = formatUpdatedRoutineRef(rt.Canonical, schemaName, newName, schemaChanged)
		}
		if rt.SubtypeDiffOID == procOID {
			rt.SubtypeDiff = formatUpdatedRoutineRef(rt.SubtypeDiff, schemaName, newName, schemaChanged)
		}
	}
}

// formatCatalogOptionName formats a qualified/unqualified object reference
// from a DefElem argument, returning "" when NONE or nil.
func formatCatalogOptionName(de *nodes.DefElem) string {
	if de == nil || de.Arg == nil {
		return ""
	}
	var parts []string
	switch v := de.Arg.(type) {
	case *nodes.List:
		for _, item := range v.Items {
			if s := stringVal(item); s != "" {
				parts = append(parts, s)
			}
		}
	case *nodes.TypeName:
		if v.Names != nil {
			for _, item := range v.Names.Items {
				if s := stringVal(item); s != "" {
					parts = append(parts, s)
				}
			}
		}
	case *nodes.String:
		if strings.EqualFold(v.Str, "none") {
			return ""
		}
		return v.Str
	default:
		s := defElemString(de)
		if strings.EqualFold(s, "none") {
			return ""
		}
		return s
	}
	if len(parts) == 0 {
		return ""
	}
	formatted := make([]string, len(parts))
	for i, part := range parts {
		formatted[i] = quoteIdentifier(part)
	}
	return strings.Join(formatted, ".")
}

// CreateFdwStmt registers a foreign-data wrapper in the catalog.
func (c *Catalog) CreateFdwStmt(stmt *nodes.CreateFdwStmt) error {
	if stmt == nil || stmt.Fdwname == "" {
		return nil
	}
	if c.fdws[stmt.Fdwname] != nil {
		return errDuplicateObject("foreign-data wrapper", stmt.Fdwname)
	}
	handler, validator, handlerOID, validatorOID, _, _, err := c.resolveFdwFuncOptions(stmt.FuncOptions)
	if err != nil {
		return err
	}
	var opts []string
	if stmt.Options != nil {
		opts, err = applyForeignOptionChanges(nil, stmt.Options)
		if err != nil {
			return err
		}
	}
	fdw := &ForeignDataWrapper{
		OID:          c.oidGen.Next(),
		Name:         stmt.Fdwname,
		OwnerName:    c.currentOwner(),
		Handler:      handler,
		HandlerOID:   handlerOID,
		Validator:    validator,
		ValidatorOID: validatorOID,
		Options:      opts,
	}
	c.fdws[stmt.Fdwname] = fdw
	c.syncFdwRoutineDeps(fdw)
	return nil
}

// AlterFdwStmt updates a foreign-data wrapper's handler, validator, or options.
func (c *Catalog) AlterFdwStmt(stmt *nodes.AlterFdwStmt) error {
	if stmt == nil || stmt.Fdwname == "" {
		return nil
	}
	fdw := c.fdws[stmt.Fdwname]
	if fdw == nil {
		return errUndefinedObject("foreign-data wrapper", stmt.Fdwname)
	}
	handler, validator, handlerOID, validatorOID, hasHandler, hasValidator, err := c.resolveFdwFuncOptions(stmt.FuncOptions)
	if err != nil {
		return err
	}
	opts := fdw.Options
	if stmt.Options != nil {
		opts, err = applyForeignOptionChanges(fdw.Options, stmt.Options)
		if err != nil {
			return err
		}
	}
	if hasHandler {
		fdw.Handler = handler
		fdw.HandlerOID = handlerOID
	}
	if hasValidator {
		fdw.Validator = validator
		fdw.ValidatorOID = validatorOID
	}
	if stmt.Options != nil {
		fdw.Options = opts
	}
	if hasHandler || hasValidator {
		c.syncFdwRoutineDeps(fdw)
	}
	return nil
}

// CreateForeignServerStmt registers a foreign server in the catalog.
func (c *Catalog) CreateForeignServerStmt(stmt *nodes.CreateForeignServerStmt) error {
	if stmt == nil || stmt.Servername == "" {
		return nil
	}
	if c.foreignServers[stmt.Servername] != nil {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, "server \""+stmt.Servername+"\" already exists, skipping")
			return nil
		}
		return errDuplicateObject("server", stmt.Servername)
	}
	if stmt.Fdwname == "" || c.fdws[stmt.Fdwname] == nil {
		return errUndefinedObject("foreign-data wrapper", stmt.Fdwname)
	}
	var opts []string
	if stmt.Options != nil {
		var err error
		opts, err = applyForeignOptionChanges(nil, stmt.Options)
		if err != nil {
			return err
		}
	}
	srv := &ForeignServer{
		OID:       c.oidGen.Next(),
		Name:      stmt.Servername,
		OwnerName: c.currentOwner(),
		FdwName:   stmt.Fdwname,
		Options:   opts,
	}
	if stmt.HasServertype || stmt.Servertype != "" {
		t := stmt.Servertype
		srv.Type = &t
	}
	if stmt.HasVersionNull {
		srv.Version = nil
	} else if stmt.HasVersion || stmt.Version != "" {
		v := stmt.Version
		srv.Version = &v
	}
	c.foreignServers[stmt.Servername] = srv
	return nil
}

// AlterForeignServerStmt updates a foreign server's version or options.
func (c *Catalog) AlterForeignServerStmt(stmt *nodes.AlterForeignServerStmt) error {
	if stmt == nil || stmt.Servername == "" {
		return nil
	}
	srv := c.foreignServers[stmt.Servername]
	if srv == nil {
		return errUndefinedObject("server", stmt.Servername)
	}
	opts := srv.Options
	if stmt.Options != nil {
		var err error
		opts, err = applyForeignOptionChanges(srv.Options, stmt.Options)
		if err != nil {
			return err
		}
	}
	if stmt.HasVersion {
		if stmt.HasVersionNull {
			srv.Version = nil
		} else {
			v := stmt.Version
			srv.Version = &v
		}
	}
	if stmt.Options != nil {
		srv.Options = opts
	}
	return nil
}

func (c *Catalog) removeFdwObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	var targets []*ForeignDataWrapper
	seen := make(map[string]bool)
	for _, obj := range stmt.Objects.Items {
		name := extractSimpleObjectName(obj)
		fdw := c.fdws[name]
		if fdw == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, "foreign-data wrapper \""+name+"\" does not exist, skipping")
				continue
			}
			return errUndefinedObject("foreign-data wrapper", name)
		}
		if c.isExtensionMember(fdw.OID) {
			return errDependentObjects("extension", c.extensionMembers[fdw.OID])
		}
		if !seen[name] {
			seen[name] = true
			targets = append(targets, fdw)
		}
	}
	depsByFdw := make(map[string][]*ForeignServer, len(targets))
	for _, fdw := range targets {
		for _, srv := range c.foreignServers {
			if srv.FdwName == fdw.Name {
				depsByFdw[fdw.Name] = append(depsByFdw[fdw.Name], srv)
			}
		}
		if len(depsByFdw[fdw.Name]) > 0 && !cascade {
			return errDependentObjects("foreign-data wrapper", fdw.Name)
		}
	}
	for _, fdw := range targets {
		if err := c.dropFdwWithServers(fdw, depsByFdw[fdw.Name], cascade); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) dropFdw(fdw *ForeignDataWrapper, cascade bool) error {
	var depServers []*ForeignServer
	for _, srv := range c.foreignServers {
		if srv.FdwName == fdw.Name {
			depServers = append(depServers, srv)
		}
	}
	return c.dropFdwWithServers(fdw, depServers, cascade)
}

func (c *Catalog) dropFdwWithServers(fdw *ForeignDataWrapper, depServers []*ForeignServer, cascade bool) error {
	for _, srv := range depServers {
		if err := c.dropForeignServer(srv, cascade); err != nil {
			return err
		}
	}
	delete(c.fdws, fdw.Name)
	delete(c.extensionMembers, fdw.OID)
	c.removeDepsOf('F', fdw.OID)
	return nil
}

func (c *Catalog) removeForeignServerObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	var targets []*ForeignServer
	seen := make(map[string]bool)
	for _, obj := range stmt.Objects.Items {
		name := extractSimpleObjectName(obj)
		srv := c.foreignServers[name]
		if srv == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, "server \""+name+"\" does not exist, skipping")
				continue
			}
			return errUndefinedObject("server", name)
		}
		if c.isExtensionMember(srv.OID) {
			return errDependentObjects("extension", c.extensionMembers[srv.OID])
		}
		if !seen[name] {
			seen[name] = true
			targets = append(targets, srv)
		}
	}
	if !cascade {
		for _, srv := range targets {
			if len(c.foreignTablesOfServer(srv.Name)) > 0 {
				return errDependentObjects("server", srv.Name)
			}
		}
	}
	for _, srv := range targets {
		if err := c.dropForeignServer(srv, cascade); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) foreignTablesOfServer(serverName string) []*Relation {
	var depTables []*Relation
	for _, rel := range c.relationByOID {
		if rel != nil && rel.RelKind == 'f' && rel.ForeignServer == serverName {
			depTables = append(depTables, rel)
		}
	}
	return depTables
}

func (c *Catalog) dropForeignServer(srv *ForeignServer, cascade bool) error {
	depTables := c.foreignTablesOfServer(srv.Name)
	if len(depTables) > 0 && !cascade {
		return errDependentObjects("server", srv.Name)
	}
	for _, rel := range depTables {
		if c.relationByOID[rel.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropRelation(rel)
			continue
		}
		if rel.Schema != nil {
			c.removeRelation(rel.Schema, rel.Name, rel)
		}
	}
	delete(c.foreignServers, srv.Name)
	delete(c.extensionMembers, srv.OID)
	return nil
}

func (c *Catalog) renameFdw(stmt *nodes.RenameStmt) error {
	oldName := extractSimpleObjectName(stmt.Object)
	fdw := c.fdws[oldName]
	if fdw == nil {
		return errUndefinedObject("foreign-data wrapper", oldName)
	}
	if c.fdws[stmt.Newname] != nil {
		return errDuplicateObject("foreign-data wrapper", stmt.Newname)
	}
	delete(c.fdws, oldName)
	fdw.Name = stmt.Newname
	c.fdws[fdw.Name] = fdw
	for _, srv := range c.foreignServers {
		if srv.FdwName == oldName {
			srv.FdwName = stmt.Newname
		}
	}
	return nil
}

func (c *Catalog) renameForeignServer(stmt *nodes.RenameStmt) error {
	oldName := extractSimpleObjectName(stmt.Object)
	srv := c.foreignServers[oldName]
	if srv == nil {
		return errUndefinedObject("server", oldName)
	}
	if c.foreignServers[stmt.Newname] != nil {
		return errDuplicateObject("server", stmt.Newname)
	}
	delete(c.foreignServers, oldName)
	srv.Name = stmt.Newname
	c.foreignServers[srv.Name] = srv
	for _, rel := range c.relationByOID {
		if rel != nil && rel.RelKind == 'f' && rel.ForeignServer == oldName {
			rel.ForeignServer = stmt.Newname
		}
	}
	return nil
}

func (c *Catalog) alterFdwOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	fdw := c.fdws[name]
	if fdw == nil {
		return errUndefinedObject("foreign-data wrapper", name)
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	fdw.OwnerName = owner
	return nil
}

func (c *Catalog) alterForeignServerOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	srv := c.foreignServers[name]
	if srv == nil {
		return errUndefinedObject("server", name)
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	srv.OwnerName = owner
	return nil
}

func foreignOptionsEqual(a, b []string) bool {
	if slices.Equal(a, b) {
		return true
	}
	aMap := make(map[string]string, len(a))
	for _, opt := range a {
		k, v := splitForeignOption(opt)
		aMap[k] = v
	}
	bMap := make(map[string]string, len(b))
	for _, opt := range b {
		k, v := splitForeignOption(opt)
		bMap[k] = v
	}
	return maps.Equal(aMap, bMap)
}

func diffForeignDataWrappers(from, to *Catalog) []ForeignDataWrapperDiffEntry {
	fromMap := make(map[string]*ForeignDataWrapper)
	if from != nil {
		for name, fdw := range from.fdws {
			if fdw != nil && !from.isExtensionMember(fdw.OID) {
				fromMap[name] = fdw
			}
		}
	}
	toMap := make(map[string]*ForeignDataWrapper)
	if to != nil {
		for name, fdw := range to.fdws {
			if fdw != nil && !to.isExtensionMember(fdw.OID) {
				toMap[name] = fdw
			}
		}
	}
	var out []ForeignDataWrapperDiffEntry
	matchedFrom := make(map[string]bool)
	matchedTo := make(map[string]bool)
	unmatchedToByOID := make(map[uint32]*ForeignDataWrapper)
	for name, toFdw := range toMap {
		if _, ok := fromMap[name]; !ok && toFdw.OID != 0 {
			unmatchedToByOID[toFdw.OID] = toFdw
		}
	}
	for name, fromFdw := range fromMap {
		if _, ok := toMap[name]; !ok && fromFdw.OID != 0 {
			if toFdw, ok := unmatchedToByOID[fromFdw.OID]; ok {
				matchedFrom[name] = true
				matchedTo[toFdw.Name] = true
				out = append(out, ForeignDataWrapperDiffEntry{Action: DiffModify, Name: toFdw.Name, From: fromFdw, To: toFdw})
				delete(unmatchedToByOID, fromFdw.OID)
			}
		}
	}
	for name, fromFdw := range fromMap {
		if matchedFrom[name] {
			continue
		}
		if _, ok := toMap[name]; !ok {
			out = append(out, ForeignDataWrapperDiffEntry{Action: DiffDrop, Name: name, From: fromFdw})
		}
	}
	for name, toFdw := range toMap {
		if matchedTo[name] {
			continue
		}
		fromFdw, ok := fromMap[name]
		if !ok {
			out = append(out, ForeignDataWrapperDiffEntry{Action: DiffAdd, Name: name, To: toFdw})
			continue
		}
		if fromFdw.Handler != toFdw.Handler ||
			fromFdw.Validator != toFdw.Validator ||
			effectiveForeignOwner(fromFdw.OwnerName) != effectiveForeignOwner(toFdw.OwnerName) ||
			!foreignOptionsEqual(fromFdw.Options, toFdw.Options) {
			out = append(out, ForeignDataWrapperDiffEntry{Action: DiffModify, Name: name, From: fromFdw, To: toFdw})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Action < out[j].Action
	})
	return out
}

func diffForeignServers(from, to *Catalog) []ForeignServerDiffEntry {
	fromMap := make(map[string]*ForeignServer)
	if from != nil {
		for name, srv := range from.foreignServers {
			if srv != nil && !from.isExtensionMember(srv.OID) {
				fromMap[name] = srv
			}
		}
	}
	toMap := make(map[string]*ForeignServer)
	if to != nil {
		for name, srv := range to.foreignServers {
			if srv != nil && !to.isExtensionMember(srv.OID) {
				toMap[name] = srv
			}
		}
	}
	var out []ForeignServerDiffEntry
	matchedFrom := make(map[string]bool)
	matchedTo := make(map[string]bool)
	unmatchedToByOID := make(map[uint32]*ForeignServer)
	for name, toSrv := range toMap {
		if _, ok := fromMap[name]; !ok && toSrv.OID != 0 {
			unmatchedToByOID[toSrv.OID] = toSrv
		}
	}
	for name, fromSrv := range fromMap {
		if _, ok := toMap[name]; !ok && fromSrv.OID != 0 {
			if toSrv, ok := unmatchedToByOID[fromSrv.OID]; ok {
				matchedFrom[name] = true
				matchedTo[toSrv.Name] = true
				out = append(out, ForeignServerDiffEntry{Action: DiffModify, Name: toSrv.Name, From: fromSrv, To: toSrv})
				delete(unmatchedToByOID, fromSrv.OID)
			}
		}
	}
	for name, fromSrv := range fromMap {
		if matchedFrom[name] {
			continue
		}
		if _, ok := toMap[name]; !ok {
			out = append(out, ForeignServerDiffEntry{Action: DiffDrop, Name: name, From: fromSrv})
		}
	}
	for name, toSrv := range toMap {
		if matchedTo[name] {
			continue
		}
		fromSrv, ok := fromMap[name]
		if !ok {
			out = append(out, ForeignServerDiffEntry{Action: DiffAdd, Name: name, To: toSrv})
			continue
		}
		if fromSrv.FdwName != toSrv.FdwName ||
			!stringPtrEqual(fromSrv.Type, toSrv.Type) ||
			!stringPtrEqual(fromSrv.Version, toSrv.Version) ||
			effectiveForeignOwner(fromSrv.OwnerName) != effectiveForeignOwner(toSrv.OwnerName) ||
			!foreignOptionsEqual(fromSrv.Options, toSrv.Options) {
			out = append(out, ForeignServerDiffEntry{Action: DiffModify, Name: name, From: fromSrv, To: toSrv})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Action < out[j].Action
	})
	return out
}

func effectiveForeignOwner(owner string) string {
	if owner == "" {
		return "postgres"
	}
	return owner
}

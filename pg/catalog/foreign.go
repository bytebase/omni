package catalog

import (
	"slices"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ForeignDataWrapper represents a PostgreSQL foreign-data wrapper (pg_foreign_data_wrapper).
type ForeignDataWrapper struct {
	OID       uint32
	Name      string
	OwnerName string
	Handler   string
	Validator string
	Options   []string
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

func extractFdwFuncOptions(funcOpts *nodes.List) (handler, validator string, hasHandler, hasValidator bool) {
	if funcOpts == nil {
		return "", "", false, false
	}
	for _, item := range funcOpts.Items {
		de, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		val := formatCatalogOptionName(de)
		switch strings.ToLower(de.Defname) {
		case "handler":
			handler = val
			hasHandler = true
		case "validator":
			validator = val
			hasValidator = true
		}
	}
	return handler, validator, hasHandler, hasValidator
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
	if len(parts) == 1 && strings.EqualFold(parts[0], "none") {
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
	handler, validator, _, _ := extractFdwFuncOptions(stmt.FuncOptions)
	var opts []string
	if stmt.Options != nil {
		var err error
		opts, err = applyForeignOptionChanges(nil, stmt.Options, c.snapshot)
		if err != nil {
			return err
		}
	}
	c.fdws[stmt.Fdwname] = &ForeignDataWrapper{
		OID:       c.oidGen.Next(),
		Name:      stmt.Fdwname,
		OwnerName: c.currentOwner(),
		Handler:   handler,
		Validator: validator,
		Options:   opts,
	}
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
	handler, validator, hasHandler, hasValidator := extractFdwFuncOptions(stmt.FuncOptions)
	if hasHandler {
		fdw.Handler = handler
	}
	if hasValidator {
		fdw.Validator = validator
	}
	if stmt.Options != nil {
		opts, err := applyForeignOptionChanges(fdw.Options, stmt.Options, c.snapshot)
		if err != nil {
			return err
		}
		fdw.Options = opts
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
	var opts []string
	if stmt.Options != nil {
		var err error
		opts, err = applyForeignOptionChanges(nil, stmt.Options, c.snapshot)
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
	if stmt.HasVersion || stmt.Version != "" {
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
	if stmt.HasVersion {
		v := stmt.Version
		srv.Version = &v
	}
	if stmt.Options != nil {
		opts, err := applyForeignOptionChanges(srv.Options, stmt.Options, c.snapshot)
		if err != nil {
			return err
		}
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
	for _, fdw := range targets {
		var depServers []*ForeignServer
		for _, srv := range c.foreignServers {
			if srv.FdwName == fdw.Name {
				depServers = append(depServers, srv)
			}
		}
		if len(depServers) > 0 && !cascade {
			return errDependentObjects("foreign-data wrapper", fdw.Name)
		}
		for _, srv := range depServers {
			if err := c.dropForeignServer(srv, cascade); err != nil {
				return err
			}
		}
		delete(c.fdws, fdw.Name)
		delete(c.extensionMembers, fdw.OID)
	}
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
	for _, srv := range targets {
		if err := c.dropForeignServer(srv, cascade); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) dropForeignServer(srv *ForeignServer, cascade bool) error {
	var depTables []*Relation
	for _, rel := range c.relationByOID {
		if rel != nil && rel.RelKind == 'f' && rel.ForeignServer == srv.Name {
			depTables = append(depTables, rel)
		}
	}
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
	for name, fromFdw := range fromMap {
		if _, ok := toMap[name]; !ok {
			out = append(out, ForeignDataWrapperDiffEntry{Action: DiffDrop, Name: name, From: fromFdw})
		}
	}
	for name, toFdw := range toMap {
		fromFdw, ok := fromMap[name]
		if !ok {
			out = append(out, ForeignDataWrapperDiffEntry{Action: DiffAdd, Name: name, To: toFdw})
			continue
		}
		if fromFdw.Handler != toFdw.Handler ||
			fromFdw.Validator != toFdw.Validator ||
			!slices.Equal(fromFdw.Options, toFdw.Options) {
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
	for name, fromSrv := range fromMap {
		if _, ok := toMap[name]; !ok {
			out = append(out, ForeignServerDiffEntry{Action: DiffDrop, Name: name, From: fromSrv})
		}
	}
	for name, toSrv := range toMap {
		fromSrv, ok := fromMap[name]
		if !ok {
			out = append(out, ForeignServerDiffEntry{Action: DiffAdd, Name: name, To: toSrv})
			continue
		}
		if fromSrv.FdwName != toSrv.FdwName ||
			!stringPtrEqual(fromSrv.Type, toSrv.Type) ||
			!stringPtrEqual(fromSrv.Version, toSrv.Version) ||
			!slices.Equal(fromSrv.Options, toSrv.Options) {
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

package ddl

// Clone creates a deep copy of the Catalog.
//
// Builtin immutable data (casts, operators, builtin procs, builtin types) is shared
// by reference. All user-created objects are deep-copied so modifications to the
// clone do not affect the original, and vice versa.
func (c *Catalog) Clone() *Catalog {
	clone := &Catalog{
		oidGen: &OIDGenerator{next: c.oidGen.next},
	}

	// --- Shared builtin indexes (immutable, pointer-safe) ---
	clone.castIndex = c.castIndex
	clone.operByOID = c.operByOID
	clone.operByKey = c.operByKey

	// Procs: start from builtin, user procs added below.
	clone.procByOID = make(map[uint32]*BuiltinProc, len(c.procByOID))
	clone.procByName = make(map[string][]*BuiltinProc, len(c.procByName))
	for oid, p := range c.procByOID {
		if oid < FirstNormalObjectId {
			clone.procByOID[oid] = p // builtin: share pointer
		}
		// user procs handled after userProcs deep copy
	}
	for name, list := range c.procByName {
		var cloned []*BuiltinProc
		for _, p := range list {
			if p.OID < FirstNormalObjectId {
				cloned = append(cloned, p)
			}
			// user procs handled below
		}
		if len(cloned) > 0 {
			clone.procByName[name] = cloned
		}
	}

	// --- Types: builtin shared, user deep-copied ---
	clone.typeByOID = make(map[uint32]*BuiltinType, len(c.typeByOID))
	clone.typeByName = make(map[typeKey]*BuiltinType, len(c.typeByName))
	for oid, t := range c.typeByOID {
		if oid < FirstNormalObjectId {
			clone.typeByOID[oid] = t // builtin: share pointer
			clone.typeByName[typeKey{ns: t.Namespace, name: t.TypeName}] = t
		} else {
			ct := *t // deep copy struct
			clone.typeByOID[oid] = &ct
			clone.typeByName[typeKey{ns: ct.Namespace, name: ct.TypeName}] = &ct
		}
	}

	// --- Schemas: deep copy ---
	clone.schemas = make(map[uint32]*Schema, len(c.schemas))
	clone.schemaByName = make(map[string]*Schema, len(c.schemaByName))
	for oid, s := range c.schemas {
		cs := &Schema{
			OID:       s.OID,
			Name:      s.Name,
			Owner:     s.Owner,
			Relations: make(map[string]*Relation, len(s.Relations)),
			Indexes:   make(map[string]*Index, len(s.Indexes)),
			Sequences: make(map[string]*Sequence, len(s.Sequences)),
		}
		clone.schemas[oid] = cs
		clone.schemaByName[s.Name] = cs
	}

	// --- Relations: deep copy + fix Schema pointer ---
	clone.relationByOID = make(map[uint32]*Relation, len(c.relationByOID))
	for oid, r := range c.relationByOID {
		cr := *r // shallow copy struct
		// Deep copy slices.
		if r.Columns != nil {
			cr.Columns = make([]*Column, len(r.Columns))
			for i, col := range r.Columns {
				cc := *col
				cr.Columns[i] = &cc
			}
		}
		if r.colByName != nil {
			cr.colByName = make(map[string]int, len(r.colByName))
			for k, v := range r.colByName {
				cr.colByName[k] = v
			}
		}
		if r.InhParents != nil {
			cr.InhParents = make([]uint32, len(r.InhParents))
			copy(cr.InhParents, r.InhParents)
		}
		// PartitionInfo: deep copy if non-nil.
		if r.PartitionInfo != nil {
			pi := *r.PartitionInfo
			if r.PartitionInfo.KeyAttNums != nil {
				pi.KeyAttNums = make([]int16, len(r.PartitionInfo.KeyAttNums))
				copy(pi.KeyAttNums, r.PartitionInfo.KeyAttNums)
			}
			cr.PartitionInfo = &pi
		}
		// PartitionBound: deep copy if non-nil.
		if r.PartitionBound != nil {
			pb := *r.PartitionBound
			if r.PartitionBound.ListValues != nil {
				pb.ListValues = make([]string, len(r.PartitionBound.ListValues))
				copy(pb.ListValues, r.PartitionBound.ListValues)
			}
			if r.PartitionBound.LowerBound != nil {
				pb.LowerBound = make([]string, len(r.PartitionBound.LowerBound))
				copy(pb.LowerBound, r.PartitionBound.LowerBound)
			}
			if r.PartitionBound.UpperBound != nil {
				pb.UpperBound = make([]string, len(r.PartitionBound.UpperBound))
				copy(pb.UpperBound, r.PartitionBound.UpperBound)
			}
			cr.PartitionBound = &pb
		}
		// AnalyzedQuery: shallow copy (shared, not mutated by dry-run).
		// Fix Schema pointer.
		cr.Schema = clone.schemas[r.Schema.OID]
		clone.relationByOID[oid] = &cr
		// Register in schema.
		cr.Schema.Relations[cr.Name] = &cr
	}

	// --- Constraints: deep copy + rebuild consByRel ---
	clone.constraints = make(map[uint32]*Constraint, len(c.constraints))
	clone.consByRel = make(map[uint32][]*Constraint, len(c.consByRel))
	for oid, con := range c.constraints {
		cc := *con
		if con.Columns != nil {
			cc.Columns = make([]int16, len(con.Columns))
			copy(cc.Columns, con.Columns)
		}
		if con.FColumns != nil {
			cc.FColumns = make([]int16, len(con.FColumns))
			copy(cc.FColumns, con.FColumns)
		}
		if con.ExclOps != nil {
			cc.ExclOps = make([]string, len(con.ExclOps))
			copy(cc.ExclOps, con.ExclOps)
		}
		if con.PFEqOp != nil {
			cc.PFEqOp = make([]uint32, len(con.PFEqOp))
			copy(cc.PFEqOp, con.PFEqOp)
		}
		if con.PPEqOp != nil {
			cc.PPEqOp = make([]uint32, len(con.PPEqOp))
			copy(cc.PPEqOp, con.PPEqOp)
		}
		if con.FFEqOp != nil {
			cc.FFEqOp = make([]uint32, len(con.FFEqOp))
			copy(cc.FFEqOp, con.FFEqOp)
		}
		if con.FKDelSetCols != nil {
			cc.FKDelSetCols = make([]int16, len(con.FKDelSetCols))
			copy(cc.FKDelSetCols, con.FKDelSetCols)
		}
		clone.constraints[oid] = &cc
		clone.consByRel[cc.RelOID] = append(clone.consByRel[cc.RelOID], &cc)
	}

	// --- Indexes: deep copy + fix Schema pointer + rebuild indexesByRel ---
	clone.indexes = make(map[uint32]*Index, len(c.indexes))
	clone.indexesByRel = make(map[uint32][]*Index, len(c.indexesByRel))
	for oid, idx := range c.indexes {
		ci := *idx
		if idx.Columns != nil {
			ci.Columns = make([]int16, len(idx.Columns))
			copy(ci.Columns, idx.Columns)
		}
		if idx.IndOption != nil {
			ci.IndOption = make([]int16, len(idx.IndOption))
			copy(ci.IndOption, idx.IndOption)
		}
		ci.Schema = clone.schemas[idx.Schema.OID]
		clone.indexes[oid] = &ci
		clone.indexesByRel[ci.RelOID] = append(clone.indexesByRel[ci.RelOID], &ci)
		// Register in schema.
		ci.Schema.Indexes[ci.Name] = &ci
	}

	// --- Sequences: deep copy + fix Schema pointer ---
	clone.sequenceByOID = make(map[uint32]*Sequence, len(c.sequenceByOID))
	for oid, seq := range c.sequenceByOID {
		cs := *seq
		cs.Schema = clone.schemas[seq.Schema.OID]
		clone.sequenceByOID[oid] = &cs
		// Register in schema.
		cs.Schema.Sequences[cs.Name] = &cs
	}

	// --- Enum types: deep copy ---
	clone.enumTypes = make(map[uint32]*EnumType, len(c.enumTypes))
	for oid, et := range c.enumTypes {
		ce := EnumType{
			TypeOID: et.TypeOID,
		}
		if et.Values != nil {
			ce.Values = make([]*EnumValue, len(et.Values))
			ce.labelMap = make(map[string]*EnumValue, len(et.Values))
			for i, v := range et.Values {
				cv := *v
				ce.Values[i] = &cv
				ce.labelMap[cv.Label] = &cv
			}
		}
		clone.enumTypes[oid] = &ce
	}

	// --- Domain types: deep copy ---
	clone.domainTypes = make(map[uint32]*DomainType, len(c.domainTypes))
	for oid, dt := range c.domainTypes {
		cd := DomainType{
			TypeOID:     dt.TypeOID,
			BaseTypeOID: dt.BaseTypeOID,
			BaseTypMod:  dt.BaseTypMod,
			NotNull:     dt.NotNull,
			Default:     dt.Default,
		}
		if dt.Constraints != nil {
			cd.Constraints = make([]*DomainConstraint, len(dt.Constraints))
			for i, dc := range dt.Constraints {
				cc := *dc
				cd.Constraints[i] = &cc
			}
		}
		clone.domainTypes[oid] = &cd
	}

	// --- Range types: deep copy (all value types) ---
	clone.rangeTypes = make(map[uint32]*RangeType, len(c.rangeTypes))
	for oid, rt := range c.rangeTypes {
		cr := *rt
		clone.rangeTypes[oid] = &cr
	}

	// --- User procs: deep copy + fix Schema pointer + register in procByOID/procByName ---
	clone.userProcs = make(map[uint32]*UserProc, len(c.userProcs))
	for oid, up := range c.userProcs {
		cup := *up
		if up.ArgTypes != nil {
			cup.ArgTypes = make([]uint32, len(up.ArgTypes))
			copy(cup.ArgTypes, up.ArgTypes)
		}
		cup.Schema = clone.schemas[up.Schema.OID]
		clone.userProcs[oid] = &cup

		// Also deep-copy the corresponding BuiltinProc entry.
		origBP := c.procByOID[oid]
		if origBP != nil {
			cbp := *origBP
			if origBP.ArgTypes != nil {
				cbp.ArgTypes = make([]uint32, len(origBP.ArgTypes))
				copy(cbp.ArgTypes, origBP.ArgTypes)
			}
			if origBP.AllArgTypes != nil {
				cbp.AllArgTypes = make([]uint32, len(origBP.AllArgTypes))
				copy(cbp.AllArgTypes, origBP.AllArgTypes)
			}
			clone.procByOID[oid] = &cbp
			clone.procByName[cbp.Name] = append(clone.procByName[cbp.Name], &cbp)
		}
	}

	// --- Triggers: deep copy + rebuild triggersByRel ---
	clone.triggers = make(map[uint32]*Trigger, len(c.triggers))
	clone.triggersByRel = make(map[uint32][]*Trigger, len(c.triggersByRel))
	for oid, trig := range c.triggers {
		ct := *trig
		if trig.Columns != nil {
			ct.Columns = make([]int16, len(trig.Columns))
			copy(ct.Columns, trig.Columns)
		}
		if trig.Args != nil {
			ct.Args = make([]string, len(trig.Args))
			copy(ct.Args, trig.Args)
		}
		clone.triggers[oid] = &ct
		clone.triggersByRel[ct.RelOID] = append(clone.triggersByRel[ct.RelOID], &ct)
	}

	// --- Comments: copy map ---
	clone.comments = make(map[commentKey]string, len(c.comments))
	for k, v := range c.comments {
		clone.comments[k] = v
	}

	// --- Grants: copy slice, deep copy Columns ---
	if c.grants != nil {
		clone.grants = make([]Grant, len(c.grants))
		for i, g := range c.grants {
			clone.grants[i] = g
			if g.Columns != nil {
				clone.grants[i].Columns = make([]string, len(g.Columns))
				copy(clone.grants[i].Columns, g.Columns)
			}
		}
	}

	// --- Policies: deep copy + rebuild policiesByRel ---
	clone.policies = make(map[uint32]*Policy, len(c.policies))
	clone.policiesByRel = make(map[uint32][]*Policy, len(c.policiesByRel))
	for oid, p := range c.policies {
		cp := *p
		if p.Roles != nil {
			cp.Roles = make([]string, len(p.Roles))
			copy(cp.Roles, p.Roles)
		}
		clone.policies[oid] = &cp
		clone.policiesByRel[cp.RelOID] = append(clone.policiesByRel[cp.RelOID], &cp)
	}

	// --- Inheritance entries: copy slice (value types) ---
	if c.inhEntries != nil {
		clone.inhEntries = make([]InhEntry, len(c.inhEntries))
		copy(clone.inhEntries, c.inhEntries)
	}

	// --- Dependencies: copy slice (value types) ---
	if c.deps != nil {
		clone.deps = make([]DepEntry, len(c.deps))
		copy(clone.deps, c.deps)
	}

	// --- Extensions: deep copy ---
	clone.extensions = make(map[uint32]*Extension, len(c.extensions))
	clone.extByName = make(map[string]*Extension, len(c.extByName))
	for oid, ext := range c.extensions {
		ce := *ext
		clone.extensions[oid] = &ce
		clone.extByName[ce.Name] = &ce
	}

	// --- Access methods: share built-in, deep copy user-created ---
	clone.accessMethods = make(map[uint32]*AccessMethod, len(c.accessMethods))
	clone.accessMethodByName = make(map[string]*AccessMethod, len(c.accessMethodByName))
	for oid, am := range c.accessMethods {
		if oid < FirstNormalObjectId {
			clone.accessMethods[oid] = am // built-in: share
			clone.accessMethodByName[am.Name] = am
		} else {
			ca := *am
			clone.accessMethods[oid] = &ca
			clone.accessMethodByName[ca.Name] = &ca
		}
	}

	// --- OpFamilies: deep copy ---
	clone.opFamilies = make(map[uint32]*OpFamily, len(c.opFamilies))
	for oid, fam := range c.opFamilies {
		cf := *fam
		clone.opFamilies[oid] = &cf
	}

	// --- OpClasses: deep copy + rebuild opClassByKey ---
	clone.opClasses = make(map[uint32]*OpClass, len(c.opClasses))
	clone.opClassByKey = make(map[opClassKey]*OpClass, len(c.opClassByKey))
	for oid, opc := range c.opClasses {
		cc := *opc
		clone.opClasses[oid] = &cc
		if cc.IsDefault {
			clone.opClassByKey[opClassKey{amOID: cc.AMOID, typeOID: cc.TypeOID}] = &cc
		}
	}

	// --- Search path: copy slice ---
	if c.searchPath != nil {
		clone.searchPath = make([]string, len(c.searchPath))
		copy(clone.searchPath, c.searchPath)
	}

	// Warnings: start clean (not copied from source).
	// clone.warnings is nil by default.

	return clone
}

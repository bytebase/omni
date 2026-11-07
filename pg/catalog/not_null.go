package catalog

import "fmt"

// ensureNotNull records PostgreSQL 18's first-class NOT NULL constraint.
// A primary key, identity, serial, and explicit NOT NULL all share one
// constraint per column. LIKE and inheritance preserve the source constraint name.
// https://github.com/postgres/postgres/blob/REL_18_STABLE/src/backend/catalog/heap.c
func (c *Catalog) ensureNotNull(rel *Relation, col *Column, name string) error {
	for _, con := range c.consByRel[rel.OID] {
		if con.Type == ConstraintNotNull && len(con.Columns) == 1 && con.Columns[0] == col.AttNum {
			col.NotNull = true
			col.NotNullName = con.Name
			return nil
		}
	}
	if name == "" {
		name = generateConstraintName(rel.Name, []string{col.Name}, ConstraintNotNull)
	}
	base := name
	for i := 1; ; i++ {
		collision := false
		for _, con := range c.consByRel[rel.OID] {
			if con.Name == name {
				collision = true
				break
			}
		}
		if !collision {
			break
		}
		name = fmt.Sprintf("%s%d", base, i)
	}
	col.NotNull = true
	col.NotNullName = name
	con := &Constraint{OID: c.oidGen.Next(), Name: name, Type: ConstraintNotNull, RelOID: rel.OID, Namespace: rel.Schema.OID, Columns: []int16{col.AttNum}, Validated: true, ConIsLocal: col.IsLocal, ConInhCount: int16(col.InhCount)}
	c.registerConstraint(rel.OID, con)
	c.recordDependency('c', con.OID, 0, 'r', rel.OID, int32(col.AttNum), DepAuto)
	return nil
}

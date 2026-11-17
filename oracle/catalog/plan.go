package catalog

import "strings"

// MigrationOp is one statement of a migration plan.
type MigrationOp struct {
	Type       string
	SchemaName string
	ObjectName string
	SQL        string
	Warning    string
}

// MigrationPlan is an ordered set of statements.
type MigrationPlan struct {
	ops []MigrationOp
}

// Ops returns the plan's operations.
func (p *MigrationPlan) Ops() []MigrationOp {
	if p == nil {
		return nil
	}
	return p.ops
}

// SQL renders the plan.
func (p *MigrationPlan) SQL() string {
	if p == nil {
		return ""
	}
	statements := make([]string, 0, len(p.ops))
	for _, op := range p.ops {
		text := strings.TrimSpace(op.SQL)
		if text != "" {
			statements = append(statements, text)
		}
	}
	return strings.Join(statements, "\n")
}

// IsEmpty reports whether the plan has no statements.
func (p *MigrationPlan) IsEmpty() bool {
	return p == nil || strings.TrimSpace(p.SQL()) == ""
}

package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateFunctionDDL produces CREATE/DROP FUNCTION/PROCEDURE operations from the diff.
func generateFunctionDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Functions {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			ops = append(ops, buildCreateFunctionOp(to, entry))

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, buildDropFunctionOp(from, entry))

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			if signatureChanged(from, entry.From, to, entry.To) {
				// Signature changed → must DROP + CREATE.
				ops = append(ops, buildDropFunctionOp(from, entry))
				ops = append(ops, buildCreateFunctionOp(to, entry))
			} else if entry.To.Kind == 'p' {
				// Procedure body-only change → CREATE OR REPLACE PROCEDURE.
				ops = append(ops, buildReplaceProcedureOp(to, entry))
			} else {
				// Function body-only change → CREATE OR REPLACE FUNCTION.
				ops = append(ops, buildReplaceeFunctionOp(to, entry))
			}
		}
	}

	// Deterministic ordering: drops first (sorted by identity), then creates (sorted).
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			if ops[i].Type == OpDropFunction {
				return true
			}
			if ops[j].Type == OpDropFunction {
				return false
			}
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// buildCreateFunctionOp generates a CREATE FUNCTION or CREATE PROCEDURE op.
func buildCreateFunctionOp(c *Catalog, entry FunctionDiffEntry) MigrationOp {
	proc := entry.To
	var b strings.Builder

	if proc.Kind == 'p' {
		b.WriteString("CREATE PROCEDURE ")
	} else {
		b.WriteString("CREATE FUNCTION ")
	}

	b.WriteString(migrationQualifiedName(entry.SchemaName, proc.Name))
	b.WriteString("(")
	b.WriteString(formatFuncParams(c, proc))
	b.WriteString(")")

	// Return type (functions only).
	if proc.Kind != 'p' {
		b.WriteString("\nRETURNS ")
		if proc.RetSet {
			b.WriteString("SETOF ")
		}
		b.WriteString(c.FormatType(proc.RetType, -1))
	}

	// Language.
	b.WriteString("\nLANGUAGE ")
	b.WriteString(proc.Language)

	// Attributes (functions only).
	if proc.Kind != 'p' {
		b.WriteString(formatFuncAttributes(proc))
	}

	// Body.
	b.WriteString("\nAS ")
	b.WriteString(dollarQuote(proc.Body))

	return MigrationOp{
		Type:          OpCreateFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
	}
}

// buildReplaceeFunctionOp generates a CREATE OR REPLACE FUNCTION op.
func buildReplaceeFunctionOp(c *Catalog, entry FunctionDiffEntry) MigrationOp {
	proc := entry.To
	var b strings.Builder

	b.WriteString("CREATE OR REPLACE FUNCTION ")
	b.WriteString(migrationQualifiedName(entry.SchemaName, proc.Name))
	b.WriteString("(")
	b.WriteString(formatFuncParams(c, proc))
	b.WriteString(")")

	b.WriteString("\nRETURNS ")
	if proc.RetSet {
		b.WriteString("SETOF ")
	}
	b.WriteString(c.FormatType(proc.RetType, -1))

	b.WriteString("\nLANGUAGE ")
	b.WriteString(proc.Language)

	b.WriteString(formatFuncAttributes(proc))

	b.WriteString("\nAS ")
	b.WriteString(dollarQuote(proc.Body))

	return MigrationOp{
		Type:          OpAlterFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
	}
}

// buildDropFunctionOp generates a DROP FUNCTION or DROP PROCEDURE op.
func buildDropFunctionOp(c *Catalog, entry FunctionDiffEntry) MigrationOp {
	proc := entry.From
	var b strings.Builder

	if proc.Kind == 'p' {
		b.WriteString("DROP PROCEDURE ")
	} else {
		b.WriteString("DROP FUNCTION ")
	}

	b.WriteString(migrationQualifiedName(entry.SchemaName, proc.Name))
	b.WriteString("(")
	b.WriteString(formatDropSignature(c, proc))
	b.WriteString(")")

	return MigrationOp{
		Type:          OpDropFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
	}
}

// formatFuncParams formats the parameter list for CREATE FUNCTION/PROCEDURE.
func formatFuncParams(c *Catalog, proc *UserProc) string {
	if len(proc.ArgTypes) == 0 && len(proc.AllArgTypes) == 0 {
		return ""
	}

	// Use AllArgTypes if available (includes OUT params), otherwise ArgTypes.
	argTypes := proc.AllArgTypes
	if argTypes == nil {
		argTypes = proc.ArgTypes
	}

	var parts []string
	for i, toid := range argTypes {
		var param strings.Builder

		// Mode prefix.
		if proc.ArgModes != nil && i < len(proc.ArgModes) {
			switch proc.ArgModes[i] {
			case 'o':
				param.WriteString("OUT ")
			case 'b':
				param.WriteString("INOUT ")
			case 'v':
				param.WriteString("VARIADIC ")
			// 'i' (IN) is implicit, 't' (TABLE) handled separately
			}
		}

		// Parameter name.
		if proc.ArgNames != nil && i < len(proc.ArgNames) && proc.ArgNames[i] != "" {
			param.WriteString(quoteIdentAlways(proc.ArgNames[i]))
			param.WriteString(" ")
		}

		// Type.
		param.WriteString(c.FormatType(toid, -1))

		parts = append(parts, param.String())
	}

	return strings.Join(parts, ", ")
}

// formatDropSignature formats the argument type list for DROP FUNCTION/PROCEDURE.
// Only IN, INOUT, and VARIADIC parameters are included (not OUT).
func formatDropSignature(c *Catalog, proc *UserProc) string {
	if len(proc.ArgTypes) == 0 && len(proc.AllArgTypes) == 0 {
		return ""
	}

	// If no AllArgTypes, all params are IN — use ArgTypes directly.
	if proc.AllArgTypes == nil {
		var parts []string
		for _, toid := range proc.ArgTypes {
			parts = append(parts, c.FormatType(toid, -1))
		}
		return strings.Join(parts, ", ")
	}

	// Filter: only include IN ('i'), INOUT ('b'), VARIADIC ('v').
	var parts []string
	for i, toid := range proc.AllArgTypes {
		if proc.ArgModes != nil && i < len(proc.ArgModes) {
			mode := proc.ArgModes[i]
			if mode == 'o' || mode == 't' {
				continue
			}
		}
		parts = append(parts, c.FormatType(toid, -1))
	}

	return strings.Join(parts, ", ")
}

// formatFuncAttributes formats function attributes (volatility, strictness, etc.).
func formatFuncAttributes(proc *UserProc) string {
	var parts []string

	// Volatility.
	switch proc.Volatile {
	case 'i':
		parts = append(parts, "IMMUTABLE")
	case 's':
		parts = append(parts, "STABLE")
	// 'v' (VOLATILE) is default, omit.
	}

	// Strictness.
	if proc.IsStrict {
		parts = append(parts, "STRICT")
	}

	// Security.
	if proc.SecDef {
		parts = append(parts, "SECURITY DEFINER")
	}

	// Leakproof.
	if proc.LeakProof {
		parts = append(parts, "LEAKPROOF")
	}

	// Parallel safety.
	switch proc.Parallel {
	case 's':
		parts = append(parts, "PARALLEL SAFE")
	case 'r':
		parts = append(parts, "PARALLEL RESTRICTED")
	// 'u' (UNSAFE) is default, omit.
	}

	if len(parts) == 0 {
		return ""
	}
	return "\n" + strings.Join(parts, "\n")
}

// buildReplaceProcedureOp generates a CREATE OR REPLACE PROCEDURE op.
func buildReplaceProcedureOp(c *Catalog, entry FunctionDiffEntry) MigrationOp {
	proc := entry.To
	var b strings.Builder

	b.WriteString("CREATE OR REPLACE PROCEDURE ")
	b.WriteString(migrationQualifiedName(entry.SchemaName, proc.Name))
	b.WriteString("(")
	b.WriteString(formatFuncParams(c, proc))
	b.WriteString(")")

	b.WriteString("\nLANGUAGE ")
	b.WriteString(proc.Language)

	b.WriteString("\nAS ")
	b.WriteString(dollarQuote(proc.Body))

	return MigrationOp{
		Type:          OpAlterFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
	}
}

// signatureChanged returns true if the function/procedure signature differs
// between two versions (kind, argument types, return type, or setof).
func signatureChanged(fromCat *Catalog, from *UserProc, toCat *Catalog, to *UserProc) bool {
	if from.Kind != to.Kind {
		return true
	}
	if len(from.ArgTypes) != len(to.ArgTypes) {
		return true
	}
	for i := range from.ArgTypes {
		if fromCat.FormatType(from.ArgTypes[i], -1) != toCat.FormatType(to.ArgTypes[i], -1) {
			return true
		}
	}
	if from.Kind != 'p' {
		// Only functions have return types.
		if fromCat.FormatType(from.RetType, -1) != toCat.FormatType(to.RetType, -1) {
			return true
		}
		if from.RetSet != to.RetSet {
			return true
		}
	}
	return false
}

// dollarQuote wraps a function body in dollar-quoting.
// If the body contains $$, uses $fn$...$fn$ instead.
func dollarQuote(body string) string {
	if strings.Contains(body, "$$") {
		return fmt.Sprintf("$fn$%s$fn$", body)
	}
	return fmt.Sprintf("$$%s$$", body)
}

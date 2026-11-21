package catalog

import (
	"fmt"
	"slices"
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
		if hasTableReturn(proc) {
			b.WriteString("\nRETURNS ")
			b.WriteString(formatReturnsTable(c, proc))
		} else {
			b.WriteString("\nRETURNS ")
			if proc.RetSet {
				b.WriteString("SETOF ")
			}
			b.WriteString(c.FormatType(proc.RetType, -1))
		}
	}

	// Language.
	b.WriteString("\nLANGUAGE ")
	b.WriteString(proc.Language)

	// Attributes (functions only), or SECURITY DEFINER and SET options for procedures.
	if proc.Kind != 'p' {
		b.WriteString(formatFuncAttributes(proc))
	} else {
		if proc.SecDef {
			b.WriteString("\nSECURITY DEFINER")
		}
		for _, sc := range proc.SetConfig {
			b.WriteByte('\n')
			b.WriteString(sc)
		}
	}

	// Body.
	b.WriteString(formatFunctionBodyClause(proc))

	return MigrationOp{
		Type:          OpCreateFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'f',
		ObjOID:        entry.To.OID,
		Priority:      PriorityFunction,
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

	if hasTableReturn(proc) {
		b.WriteString("\nRETURNS ")
		b.WriteString(formatReturnsTable(c, proc))
	} else {
		b.WriteString("\nRETURNS ")
		if proc.RetSet {
			b.WriteString("SETOF ")
		}
		b.WriteString(c.FormatType(proc.RetType, -1))
	}

	b.WriteString("\nLANGUAGE ")
	b.WriteString(proc.Language)

	b.WriteString(formatFuncAttributes(proc))

	b.WriteString(formatFunctionBodyClause(proc))

	return MigrationOp{
		Type:          OpAlterFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'f',
		ObjOID:        entry.To.OID,
		Priority:      PriorityFunction,
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
		Phase:         PhasePre,
		ObjType:       'f',
		ObjOID:        entry.From.OID,
		Priority:      PriorityFunction,
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
	inputIdx := 0
	for i, toid := range argTypes {
		// Skip TABLE mode params — they go in RETURNS TABLE(...).
		var mode byte
		if proc.ArgModes != nil && i < len(proc.ArgModes) {
			mode = proc.ArgModes[i]
		}
		if mode == 't' {
			continue
		}

		var param strings.Builder

		// Mode prefix.
		switch mode {
		case 'o':
			param.WriteString("OUT ")
		case 'b':
			param.WriteString("INOUT ")
		case 'v':
			param.WriteString("VARIADIC ")
			// 'i' (IN) is implicit, 't' (TABLE) handled separately
		}

		// Parameter name.
		if proc.ArgNames != nil && i < len(proc.ArgNames) && proc.ArgNames[i] != "" {
			param.WriteString(quoteIdentAlways(proc.ArgNames[i]))
			param.WriteString(" ")
		}

		// Type.
		param.WriteString(c.FormatType(toid, -1))

		// Default expression for input parameters.
		if mode != 'o' {
			if proc.ArgDefaults != nil && inputIdx < len(proc.ArgDefaults) && proc.ArgDefaults[inputIdx] != "" {
				param.WriteString(" DEFAULT ")
				param.WriteString(proc.ArgDefaults[inputIdx])
			}
			inputIdx++
		}

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

	if proc.IsWindow {
		parts = append(parts, "WINDOW")
	}
	if proc.Cost != "" {
		parts = append(parts, "COST "+proc.Cost)
	}
	if proc.Rows != "" {
		parts = append(parts, "ROWS "+proc.Rows)
	}
	parts = append(parts, proc.SetConfig...)
	if proc.Support != "" {
		parts = append(parts, "SUPPORT "+proc.Support)
	}

	if len(parts) == 0 {
		return ""
	}
	return "\n" + strings.Join(parts, "\n")
}

func formatFunctionBodyClause(proc *UserProc) string {
	if proc.HasSqlBody {
		return "\n" + proc.Body
	}
	if proc.LinkSymbol != "" {
		return "\nAS " + quoteLiteral(proc.Body) + ", " + quoteLiteral(proc.LinkSymbol)
	}
	return "\nAS " + dollarQuote(proc.Body)
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

	if proc.SecDef {
		b.WriteString("\nSECURITY DEFINER")
	}
	for _, sc := range proc.SetConfig {
		b.WriteByte('\n')
		b.WriteString(sc)
	}

	b.WriteString(formatFunctionBodyClause(proc))

	return MigrationOp{
		Type:          OpAlterFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'f',
		ObjOID:        entry.To.OID,
		Priority:      PriorityFunction,
	}
}

// signatureChanged returns true if the function/procedure signature differs
// between two versions (kind, argument types, return type, or setof).
func signatureChanged(fromCat *Catalog, from *UserProc, toCat *Catalog, to *UserProc) bool {
	if from.Kind != to.Kind {
		return true
	}
	if from.IsWindow != to.IsWindow {
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
	// Check arg names and modes — PG doesn't allow changing them via CREATE OR REPLACE.
	if !slices.Equal(from.ArgNames, to.ArgNames) {
		return true
	}
	if !slices.Equal(from.ArgModes, to.ArgModes) {
		return true
	}
	if len(from.AllArgTypes) != len(to.AllArgTypes) {
		return true
	}
	for i := range from.AllArgTypes {
		if fromCat.FormatType(from.AllArgTypes[i], -1) != toCat.FormatType(to.AllArgTypes[i], -1) {
			return true
		}
	}
	// PG doesn't allow removing input parameter defaults via CREATE OR REPLACE.
	if from.NArgDefaults > to.NArgDefaults {
		return true
	}
	return false
}

// hasTableReturn returns true if the function uses RETURNS TABLE mode.
func hasTableReturn(proc *UserProc) bool {
	for _, m := range proc.ArgModes {
		if m == 't' {
			return true
		}
	}
	return false
}

// formatReturnsTable formats a RETURNS TABLE(...) clause.
func formatReturnsTable(c *Catalog, proc *UserProc) string {
	var cols []string
	for i, mode := range proc.ArgModes {
		if mode == 't' {
			name := ""
			if i < len(proc.ArgNames) {
				name = proc.ArgNames[i]
			}
			typeName := c.FormatType(proc.AllArgTypes[i], -1)
			cols = append(cols, fmt.Sprintf("%s %s", quoteIdentAlways(name), typeName))
		}
	}
	return "TABLE(" + strings.Join(cols, ", ") + ")"
}

// dollarQuote wraps a function body in dollar-quoting using $$ when absent
// from body, or another safe dollar tag.
func dollarQuote(body string) string {
	return dollarQuoteWithDefault(body, "$$")
}

func dollarQuoteWithDefault(body, preferred string) string {
	for _, tag := range []string{preferred, "$fn$", "$body$"} {
		if tag != "" && !strings.Contains(body, tag) {
			return tag + body + tag
		}
	}
	for i := 1; ; i++ {
		tag := fmt.Sprintf("$fn%d$", i)
		if !strings.Contains(body, tag) {
			return tag + body + tag
		}
	}
}

package catalog

import (
	"regexp"
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

const (
	// Match bracketed, double-quoted, or unquoted T-SQL identifiers, including
	// non-ASCII Unicode characters accepted by the lexer (isIdentStart: byte >= 128).
	identPartExpr = `(?:\[(?:[^\]]|\]\])+\]|"(?:[^"]|"")+"|[A-Za-z_@#\x{0080}-\x{10FFFF}][A-Za-z0-9_@#$\x{0080}-\x{10FFFF}]*)`
	multipartExpr = identPartExpr + `(?:\s*\.\s*` + identPartExpr + `)*`
)

var (
	schemaBindingPattern = regexp.MustCompile(`(?i)\bSCHEMABINDING\b`)
	functionHeaderPrefix = regexp.MustCompile(`(?is)^\s*(?:CREATE\s+OR\s+ALTER|CREATE|ALTER)\s+FUNCTION\s+(` + multipartExpr + `)`)
	functionCallPattern  = regexp.MustCompile(`(?i)(` + multipartExpr + `)\s*\(`)
	multipartIdent       = regexp.MustCompile(multipartExpr)
	identPartPattern     = regexp.MustCompile(identPartExpr)
)

// isSchemaBound reports whether a view or function definition carries
// WITH SCHEMABINDING outside string literals and comments.
func isSchemaBound(definition string) bool {
	if definition == "" {
		return false
	}
	return schemaBindingPattern.MatchString(stripLiteralsAndComments(definition))
}

// expressionReferencesFunction reports whether a SQL expression or query body
// calls targetSchema.targetName outside string literals and comments.
func expressionReferencesFunction(expression, ownerSchema, targetSchema, targetName string) bool {
	if expression == "" || targetName == "" {
		return false
	}
	code := stripLiteralsAndComments(expression)
	return codeReferencesFunction(code, ownerSchema, targetSchema, targetName)
}

// definitionReferencesFunction reports whether a CREATE/ALTER FUNCTION body
// calls targetSchema.targetName, ignoring the function's own header name.
func definitionReferencesFunction(definition, ownerSchema, targetSchema, targetName string) bool {
	if definition == "" || targetName == "" {
		return false
	}
	code := stripLiteralsAndComments(definition)
	if loc := functionHeaderPrefix.FindStringIndex(code); loc != nil {
		code = code[loc[1]:]
	}
	return codeReferencesFunction(code, ownerSchema, targetSchema, targetName)
}

func codeReferencesFunction(code, ownerSchema, targetSchema, targetName string) bool {
	matches := functionCallPattern.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		parts := splitIdentifierParts(match[1])
		if len(parts) == 0 || !strings.EqualFold(parts[len(parts)-1], targetName) {
			continue
		}
		if len(parts) == 1 {
			if strings.EqualFold(targetSchema, ownerSchema) || strings.EqualFold(targetSchema, defaultSchema) {
				return true
			}
			continue
		}
		if strings.EqualFold(parts[len(parts)-2], targetSchema) {
			return true
		}
	}
	return false
}

// definitionReferencesTable reports whether a function or view definition
// references targetSchema.targetTable.
func definitionReferencesTable(definition, defaultSchemaName, targetSchema, targetTable string) bool {
	if definition == "" || targetTable == "" {
		return false
	}
	if defaultSchemaName == "" {
		defaultSchemaName = defaultSchema
	}
	if targetSchema == "" {
		targetSchema = defaultSchema
	}
	if refs, ok := parseDefinitionTableRefs(definition, defaultSchemaName); ok {
		for _, ref := range refs {
			if strings.EqualFold(ref.schema, targetSchema) && strings.EqualFold(ref.name, targetTable) {
				return true
			}
		}
		return false
	}
	// Fallback when parsing fails: inspect stripped code outside literals and
	// comments, ignoring the module header and function-call identifiers.
	code := stripLiteralsAndComments(definition)
	if loc := functionHeaderPrefix.FindStringIndex(code); loc != nil {
		code = code[loc[1]:]
	}
	for _, loc := range multipartIdent.FindAllStringIndex(code, -1) {
		after := strings.TrimLeft(code[loc[1]:], " \t\r\n")
		if strings.HasPrefix(after, "(") {
			continue
		}
		parts := splitIdentifierParts(code[loc[0]:loc[1]])
		if len(parts) == 0 || !strings.EqualFold(parts[len(parts)-1], targetTable) {
			continue
		}
		if len(parts) == 1 {
			if strings.EqualFold(targetSchema, defaultSchemaName) {
				return true
			}
			continue
		}
		if strings.EqualFold(parts[len(parts)-2], targetSchema) {
			return true
		}
	}
	return false
}

type tableRefID struct {
	schema string
	name   string
}

func parseDefinitionTableRefs(definition, defaultSchemaName string) ([]tableRefID, bool) {
	stmts, err := mssql.Parse(creationDefinition(definition))
	if err != nil || len(stmts) == 0 {
		return nil, false
	}
	var root ast.Node
	for i := range stmts {
		switch n := stmts[i].AST.(type) {
		case *ast.CreateFunctionStmt:
			if n.Body != nil {
				root = n.Body
			}
		case *ast.CreateViewStmt:
			if n.Query != nil {
				root = n.Query
			}
		}
		if root != nil {
			break
		}
	}
	if root == nil {
		return nil, false
	}
	var refs []tableRefID
	seen := make(map[string]bool)
	var collect func(node ast.Node)
	collect = func(node ast.Node) {
		ast.Inspect(node, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			switch x := n.(type) {
			case *ast.FuncCallExpr:
				// Skip x.Name (it's a function, not a table) but inspect arguments and window clauses.
				if x.Args != nil {
					collect(x.Args)
				}
				if x.Over != nil {
					collect(x.Over)
				}
				if x.Within != nil {
					collect(x.Within)
				}
				return false
			case *ast.TableRef:
				if x.Object == "" {
					return false
				}
				schemaName := x.Schema
				if schemaName == "" {
					schemaName = defaultSchemaName
				}
				key := strings.ToLower(objectID(schemaName, x.Object))
				if !seen[key] {
					seen[key] = true
					refs = append(refs, tableRefID{schema: schemaName, name: x.Object})
				}
				return false
			}
			return true
		})
	}
	collect(root)
	return refs, true
}

// viewTableDependencies returns every table or view objectID referenced by a
// view definition, combining lineage analysis with direct AST TableRef
// extraction so views like SELECT 1 FROM [dbo].[v1] are captured too.
func viewTableDependencies(viewDef, schemaName string) []string {
	seen := make(map[string]bool)
	var deps []string
	add := func(s, name string) {
		if s == "" {
			s = schemaName
		}
		if s == "" {
			s = defaultSchema
		}
		id := objectID(s, name)
		key := strings.ToLower(id)
		if !seen[key] {
			seen[key] = true
			deps = append(deps, id)
		}
	}
	if lineageDeps, err := getViewDependencies(viewDef, schemaName); err == nil {
		for _, dep := range lineageDeps {
			key := strings.ToLower(dep)
			if !seen[key] {
				seen[key] = true
				deps = append(deps, dep)
			}
		}
	}
	if refs, ok := parseDefinitionTableRefs(viewDef, schemaName); ok {
		for _, ref := range refs {
			add(ref.schema, ref.name)
		}
	}
	slices.Sort(deps)
	return deps
}

func splitIdentifierParts(raw string) []string {
	parts := identPartPattern.FindAllString(raw, -1)
	for i, p := range parts {
		parts[i] = unquoteIdentifier(p)
	}
	return parts
}

func unquoteIdentifier(part string) string {
	if len(part) >= 2 && part[0] == '[' && part[len(part)-1] == ']' {
		return strings.ReplaceAll(part[1:len(part)-1], "]]", "]")
	}
	if len(part) >= 2 && part[0] == '"' && part[len(part)-1] == '"' {
		return strings.ReplaceAll(part[1:len(part)-1], `""`, `"`)
	}
	return part
}

func stripLiteralsAndComments(expression string) string {
	var out strings.Builder
	inSingle, inDouble, inBracket, inLine := false, false, false, false
	blockDepth := 0
	for i := 0; i < len(expression); i++ {
		ch := expression[i]
		switch {
		case inLine:
			if ch == '\n' {
				inLine = false
				out.WriteByte(' ')
			}
		case blockDepth > 0:
			if ch == '/' && i+1 < len(expression) && expression[i+1] == '*' {
				blockDepth++
				i++
			} else if ch == '*' && i+1 < len(expression) && expression[i+1] == '/' {
				blockDepth--
				i++
				if blockDepth == 0 {
					out.WriteByte(' ')
				}
			}
		case inSingle:
			if ch == '\'' {
				if i+1 < len(expression) && expression[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
				out.WriteByte(' ')
			}
		case inDouble:
			out.WriteByte(ch)
			if ch == '"' {
				if i+1 < len(expression) && expression[i+1] == '"' {
					out.WriteByte('"')
					i++
					continue
				}
				inDouble = false
			}
		case inBracket:
			out.WriteByte(ch)
			if ch == ']' {
				if i+1 < len(expression) && expression[i+1] == ']' {
					out.WriteByte(']')
					i++
					continue
				}
				inBracket = false
			}
		default:
			switch {
			case ch == '-' && i+1 < len(expression) && expression[i+1] == '-':
				inLine = true
				i++
			case ch == '/' && i+1 < len(expression) && expression[i+1] == '*':
				blockDepth = 1
				i++
			case ch == '\'':
				inSingle = true
			case ch == '"':
				inDouble = true
				out.WriteByte(ch)
			case ch == '[':
				inBracket = true
				out.WriteByte(ch)
			default:
				out.WriteByte(ch)
			}
		}
	}
	return out.String()
}

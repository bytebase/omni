package catalog

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	omnianalysis "github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// Engine is SQL Server's metadata migrator and extractor. It lives here rather
// than in mssql/analysis because this package imports that one.
var Engine engine

type engine struct{}

var (
	_ omnianalysis.MetadataMigrator  = Engine
	_ omnianalysis.MetadataExtractor = Engine
)

// ExtractMetadata returns GetDatabaseMetadata's metadata for the script.
func (engine) ExtractMetadata(_ context.Context, script string) (*metadata.DatabaseSchemaMetadata, error) {
	return GetDatabaseMetadata(script)
}

// MetadataMigration diffs the two snapshots after reconcileSkipDump. Diff and
// GenerateMigration match names ignoring case, so when either side is
// case-sensitive it refuses names that differ only by case, within a snapshot
// or as a rename between them, rather than merge them into one object.
func (engine) MetadataMigration(ctx context.Context, request omnianalysis.MetadataMigrationRequest) (string, error) {
	fromMeta := omnianalysis.WithoutSchemas(databaseMetadata(request.From), request.SkipSchemas)
	toMeta := omnianalysis.WithoutSchemas(databaseMetadata(request.To), request.SkipSchemas)
	caseSensitive := objectCaseSensitive(request.From) || objectCaseSensitive(request.To)
	if caseSensitive {
		if err := rejectCaseOnlyIdentities(fromMeta, toMeta); err != nil {
			return "", err
		}
	}
	oldMeta, newMeta := reconcileSkipDump(fromMeta, toMeta, caseSensitive)
	from := New()
	if err := from.LoadMetadata(ctx, oldMeta); err != nil {
		return "", fmt.Errorf("failed to load source schema: %w", err)
	}
	to := New()
	if err := to.LoadMetadata(ctx, newMeta); err != nil {
		return "", fmt.Errorf("failed to load target schema: %w", err)
	}
	diff := Diff(from, to)
	if diff.IsEmpty() {
		return "", nil
	}
	return GenerateMigration(from, to, diff).SQL(), nil
}

// rejectCaseOnlyIdentities refuses two schema, object, column, index or check
// constraint names that differ only by case, in one snapshot or as a rename
// between the snapshots.
func rejectCaseOnlyIdentities(oldMeta, newMeta *metadata.DatabaseSchemaMetadata) error {
	folded := make([]map[string]string, 2)
	for i, meta := range []*metadata.DatabaseSchemaMetadata{oldMeta, newMeta} {
		folded[i] = make(map[string]string)
		for _, id := range identities(meta) {
			fold := strings.ToLower(id.key)
			if previous, ok := folded[i][fold]; ok && previous != id.name {
				return fmt.Errorf("case-sensitive database has %q and %q, which the catalog cannot distinguish", previous, id.name)
			}
			folded[i][fold] = id.name
		}
	}
	for _, id := range identities(newMeta) {
		if previous, ok := folded[0][strings.ToLower(id.key)]; ok && previous != id.name {
			return fmt.Errorf("case-only rename %q to %q is not supported by the catalog diff", previous, id.name)
		}
	}
	return nil
}

type identity struct {
	// key is name's kind and path. Two identities whose keys fold alike are
	// one object to the catalog.
	key, name string
}

// identities lists the names Diff and GenerateMigration match ignoring case.
// Tables, views, functions, procedures and sequences share a schema's object
// namespace.
func identities(meta *metadata.DatabaseSchemaMetadata) []identity {
	var ids []identity
	add := func(kind string, path ...string) {
		ids = append(ids, identity{key: kind + "\x00" + strings.Join(path, "\x00"), name: strings.Join(path, ".")})
	}
	for _, schema := range meta.GetSchemas() {
		s := schema.GetName()
		add("schema", s)
		for _, table := range schema.GetTables() {
			add("object", s, table.GetName())
			for _, column := range table.GetColumns() {
				add("column", s, table.GetName(), column.GetName())
			}
			for _, index := range table.GetIndexes() {
				add("index", s, table.GetName(), index.GetName())
			}
			for _, check := range table.GetCheckConstraints() {
				add("check", s, table.GetName(), check.GetName())
			}
		}
		for _, view := range schema.GetViews() {
			add("object", s, view.GetName())
		}
		for _, function := range schema.GetFunctions() {
			add("object", s, function.GetName())
		}
		for _, procedure := range schema.GetProcedures() {
			add("object", s, procedure.GetName())
		}
		for _, sequence := range schema.GetSequences() {
			add("object", s, sequence.GetName())
		}
	}
	return ids
}

func databaseMetadata(database *omnianalysis.Database) *metadata.DatabaseSchemaMetadata {
	if database == nil {
		return nil
	}
	return database.Metadata
}

func objectCaseSensitive(database *omnianalysis.Database) bool {
	return database != nil && database.ObjectCaseSensitive
}

// reconcileSkipDump makes an object skipped on either side take no part in the
// migration. Both snapshots carry the source object unchanged in that case, or
// neither carries it when the source lacks it, so the diff emits neither a drop
// nor a create or alter. It returns copies and leaves its arguments unchanged.
func reconcileSkipDump(oldMeta, newMeta *metadata.DatabaseSchemaMetadata, caseSensitive bool) (*metadata.DatabaseSchemaMetadata, *metadata.DatabaseSchemaMetadata) {
	key := func(name string) string {
		if caseSensitive {
			return name
		}
		return strings.ToLower(name)
	}
	if oldMeta == nil {
		oldMeta = &metadata.DatabaseSchemaMetadata{}
	}
	if newMeta == nil {
		newMeta = &metadata.DatabaseSchemaMetadata{}
	}
	oldOut, newOut := proto.CloneOf(oldMeta), proto.CloneOf(newMeta)
	for _, schema := range oldOut.GetSchemas() {
		schema.SkipDump = false
	}
	for _, schema := range newOut.GetSchemas() {
		schema.SkipDump = false
	}
	for _, oldSchema := range oldOut.GetSchemas() {
		if findByName(newOut.GetSchemas(), oldSchema.GetName(), key) == nil {
			clearSkipDump(oldSchema)
		}
	}
	for _, oldSchema := range oldOut.GetSchemas() {
		newSchema := findByName(newOut.GetSchemas(), oldSchema.GetName(), key)
		if newSchema == nil {
			continue
		}
		newSchema.Tables = reconcileObjects(oldSchema.GetTables(), newSchema.GetTables(),
			func(t *metadata.TableMetadata) string { return key(t.GetName()) },
			func(t *metadata.TableMetadata) bool { return t.GetSkipDump() },
			func(t *metadata.TableMetadata) { t.SkipDump = false })
		newSchema.Views = reconcileObjects(oldSchema.GetViews(), newSchema.GetViews(),
			func(v *metadata.ViewMetadata) string { return key(v.GetName()) },
			func(v *metadata.ViewMetadata) bool { return v.GetSkipDump() },
			func(v *metadata.ViewMetadata) { v.SkipDump = false })
		newSchema.Functions = reconcileObjects(oldSchema.GetFunctions(), newSchema.GetFunctions(),
			func(f *metadata.FunctionMetadata) string { return key(f.GetName()) },
			func(f *metadata.FunctionMetadata) bool { return f.GetSkipDump() },
			func(f *metadata.FunctionMetadata) { f.SkipDump = false })
		newSchema.Procedures = reconcileObjects(oldSchema.GetProcedures(), newSchema.GetProcedures(),
			func(p *metadata.ProcedureMetadata) string { return key(p.GetName()) },
			func(p *metadata.ProcedureMetadata) bool { return p.GetSkipDump() },
			func(p *metadata.ProcedureMetadata) { p.SkipDump = false })
		newSchema.Sequences = reconcileObjects(oldSchema.GetSequences(), newSchema.GetSequences(),
			func(s *metadata.SequenceMetadata) string { return key(s.GetName()) },
			func(s *metadata.SequenceMetadata) bool { return s.GetSkipDump() },
			func(s *metadata.SequenceMetadata) { s.SkipDump = false })
	}
	for _, newSchema := range newOut.GetSchemas() {
		oldSchema := findByName(oldOut.GetSchemas(), newSchema.GetName(), key)
		if oldSchema == nil {
			continue
		}
		for _, table := range newSchema.GetTables() {
			oldTable := findByName(oldSchema.GetTables(), table.GetName(), key)
			if oldTable == nil {
				continue
			}
			table.Triggers = reconcileObjects(oldTable.GetTriggers(), table.GetTriggers(),
				func(t *metadata.TriggerMetadata) string { return key(t.GetName()) },
				func(t *metadata.TriggerMetadata) bool { return t.GetSkipDump() },
				func(t *metadata.TriggerMetadata) { t.SkipDump = false })
		}
		for _, view := range newSchema.GetViews() {
			oldView := findByName(oldSchema.GetViews(), view.GetName(), key)
			if oldView == nil {
				continue
			}
			view.Triggers = reconcileObjects(oldView.GetTriggers(), view.GetTriggers(),
				func(t *metadata.TriggerMetadata) string { return key(t.GetName()) },
				func(t *metadata.TriggerMetadata) bool { return t.GetSkipDump() },
				func(t *metadata.TriggerMetadata) { t.SkipDump = false })
		}
	}
	return oldOut, newOut
}

func reconcileObjects[T any](oldObjs, newObjs []T, name func(T) string, skipped func(T) bool, unskip func(T)) []T {
	oldByName := make(map[string]T, len(oldObjs))
	for _, obj := range oldObjs {
		oldByName[name(obj)] = obj
	}
	seen := make(map[string]bool, len(newObjs))
	result := make([]T, 0, len(newObjs)+len(oldObjs))
	for _, newObj := range newObjs {
		objectName := name(newObj)
		oldObj, hasOld := oldByName[objectName]
		switch {
		case hasOld && skipped(oldObj):
			unskip(oldObj)
			result = append(result, oldObj)
		case skipped(newObj):
			if hasOld {
				unskip(oldObj)
				result = append(result, oldObj)
			}
		default:
			result = append(result, newObj)
		}
		seen[objectName] = true
	}
	for _, oldObj := range oldObjs {
		if skipped(oldObj) && !seen[name(oldObj)] {
			unskip(oldObj)
			result = append(result, oldObj)
		}
	}
	return result
}

// clearSkipDump exposes every object of a schema the target removes: the schema
// cannot be dropped while a skipped table or view still occupies it.
func clearSkipDump(schema *metadata.SchemaMetadata) {
	for _, table := range schema.GetTables() {
		table.SkipDump = false
		for _, trigger := range table.GetTriggers() {
			trigger.SkipDump = false
		}
	}
	for _, view := range schema.GetViews() {
		view.SkipDump = false
		for _, trigger := range view.GetTriggers() {
			trigger.SkipDump = false
		}
	}
	for _, function := range schema.GetFunctions() {
		function.SkipDump = false
	}
	for _, procedure := range schema.GetProcedures() {
		procedure.SkipDump = false
	}
	for _, sequence := range schema.GetSequences() {
		sequence.SkipDump = false
	}
}

func findByName[T interface{ GetName() string }](items []T, name string, key func(string) string) T {
	want := key(name)
	for _, item := range items {
		if key(item.GetName()) == want {
			return item
		}
	}
	var zero T
	return zero
}

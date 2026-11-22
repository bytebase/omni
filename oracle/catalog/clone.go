package catalog

import (
	metadatapb "github.com/bytebase/omni/metadata"
	"reflect"
)

// The catalog keeps its own copies of snapshots and objects. The metadata
// package's generated messages have no dependency-free clone of their own, so a
// small reflective copy keeps production catalog code on the standard library.

func clone[T any](value T) T {
	var out T
	copyValue(reflect.ValueOf(&out).Elem(), reflect.ValueOf(value))
	return out
}

func copyValue(dst, src reflect.Value) {
	if !src.IsValid() || !dst.CanSet() {
		return
	}
	switch src.Kind() {
	case reflect.Pointer:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.New(src.Type().Elem()))
		copyValue(dst.Elem(), src.Elem())
	case reflect.Interface:
		if src.IsNil() {
			return
		}
		inner := reflect.New(src.Elem().Type()).Elem()
		copyValue(inner, src.Elem())
		dst.Set(inner)
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Len()))
		for i := 0; i < src.Len(); i++ {
			copyValue(dst.Index(i), src.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
		iter := src.MapRange()
		for iter.Next() {
			key := reflect.New(src.Type().Key()).Elem()
			copyValue(key, iter.Key())
			value := reflect.New(src.Type().Elem()).Elem()
			copyValue(value, iter.Value())
			dst.SetMapIndex(key, value)
		}
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			if dst.Field(i).CanSet() {
				copyValue(dst.Field(i), src.Field(i))
			}
		}
	default:
		dst.Set(src)
	}
}

func cloneMetadata(meta *metadatapb.DatabaseSchemaMetadata) *metadatapb.DatabaseSchemaMetadata {
	return clone(meta)
}

func cloneSchema(schema *metadatapb.SchemaMetadata) *metadatapb.SchemaMetadata {
	return clone(schema)
}

func cloneTable(table *metadatapb.TableMetadata) *metadatapb.TableMetadata {
	return clone(table)
}

func cloneView(view *metadatapb.ViewMetadata) *metadatapb.ViewMetadata {
	return clone(view)
}

func cloneMaterializedView(view *metadatapb.MaterializedViewMetadata) *metadatapb.MaterializedViewMetadata {
	return clone(view)
}

func cloneFunction(function *metadatapb.FunctionMetadata) *metadatapb.FunctionMetadata {
	return clone(function)
}

func cloneProcedure(procedure *metadatapb.ProcedureMetadata) *metadatapb.ProcedureMetadata {
	return clone(procedure)
}

func cloneSequence(sequence *metadatapb.SequenceMetadata) *metadatapb.SequenceMetadata {
	return clone(sequence)
}

func clonePackage(pkg *metadatapb.PackageMetadata) *metadatapb.PackageMetadata {
	return clone(pkg)
}

func cloneSynonym(synonym *metadatapb.SynonymMetadata) *metadatapb.SynonymMetadata {
	return clone(synonym)
}

func cloneColumn(column *metadatapb.ColumnMetadata) *metadatapb.ColumnMetadata {
	return clone(column)
}

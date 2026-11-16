package driver

import (
	"fmt"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// collectionOptions returns the options a listCollections answer reports for a
// collection, as canonical Extended JSON, or "" when it has none.
func collectionOptions(info bson.M) (string, error) {
	options, ok := info["options"]
	if !ok || options == nil {
		return "", nil
	}
	return canonicalDocument(options)
}

// indexDefinition returns the options a listIndexes answer reports for an
// index beyond its key, name and uniqueness, as canonical Extended JSON. The
// key is an expression, and the name and uniqueness have fields of their own.
func indexDefinition(info bson.M) (string, error) {
	keys := make([]string, 0, len(info))
	for key := range info {
		switch key {
		case "key", "name", "unique":
			continue
		}
		keys = append(keys, key)
	}
	slices.Sort(keys)
	document := make(bson.D, 0, len(keys))
	for _, key := range keys {
		document = append(document, bson.E{Key: key, Value: info[key]})
	}
	return canonicalDocument(document)
}

// canonicalDocument encodes value as canonical Extended JSON, which keeps
// every BSON kind. A bson.M has no order, so its keys are sorted; a bson.D
// keeps the server's order.
func canonicalDocument(value any) (string, error) {
	document := orderedDocument(value)
	if document == nil {
		return "", fmt.Errorf("the options are a %T, not a document", value)
	}
	if len(document) == 0 {
		return "", nil
	}
	text, err := bson.MarshalExtJSON(document, true, false)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func orderedDocument(value any) bson.D {
	switch document := value.(type) {
	case bson.D:
		ordered := make(bson.D, len(document))
		for i, element := range document {
			ordered[i] = bson.E{Key: element.Key, Value: orderedValue(element.Value)}
		}
		return ordered
	case bson.M:
		keys := make([]string, 0, len(document))
		for key := range document {
			keys = append(keys, key)
		}
		slices.Sort(keys)
		ordered := make(bson.D, len(keys))
		for i, key := range keys {
			ordered[i] = bson.E{Key: key, Value: orderedValue(document[key])}
		}
		return ordered
	}
	return nil
}

func orderedValue(value any) any {
	switch value := value.(type) {
	case bson.D, bson.M:
		return orderedDocument(value)
	case bson.A:
		ordered := make(bson.A, len(value))
		for i, element := range value {
			ordered[i] = orderedValue(element)
		}
		return ordered
	case []any:
		ordered := make(bson.A, len(value))
		for i, element := range value {
			ordered[i] = orderedValue(element)
		}
		return ordered
	}
	return value
}

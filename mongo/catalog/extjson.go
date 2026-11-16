package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// renderExtendedJSON renders a canonical Extended JSON document as the
// mongosh object literal that evaluates to the same BSON. The dump writes the
// literal into a command, so an option or pipeline keeps the BSON types a
// server reported rather than the JSON types an ordinary decoder would give
// them.
//
// Only the BSON kinds gomongo can build from mongosh are rendered: numbers,
// strings, booleans, null, documents, arrays, ObjectID, date, timestamp,
// regular expression and Decimal128. A binary, min/max key, symbol, JavaScript
// value or DBPointer has no mongosh spelling in gomongo, so an option carrying
// one is refused rather than silently written as something else.
func renderExtendedJSON(text string) (string, error) {
	dec := json.NewDecoder(strings.NewReader(text))
	dec.UseNumber()
	value, err := decodeJSONValue(dec)
	if err != nil {
		return "", err
	}
	if _, err := dec.Token(); err != io.EOF {
		return "", fmt.Errorf("trailing data after the document")
	}
	return value.literal()
}

// jsonValue is one value of an Extended JSON document, in source order.
type jsonValue struct {
	kind     jsonKind
	boolean  bool
	text     string
	number   json.Number
	fields   []jsonField
	elements []*jsonValue
}

type jsonField struct {
	key   string
	value *jsonValue
}

type jsonKind uint8

const (
	jsonNull jsonKind = iota
	jsonBool
	jsonString
	jsonNumber
	jsonArray
	jsonObject
)

// decodeJSONValue reads one JSON value, preserving the order of an object's
// fields: a MongoDB document is ordered, and a pipeline or an index key can
// depend on that order.
func decodeJSONValue(dec *json.Decoder) (*jsonValue, error) {
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch token := token.(type) {
	case json.Delim:
		switch token {
		case '{':
			value := &jsonValue{kind: jsonObject}
			for dec.More() {
				keyToken, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("object key is %T, not a string", keyToken)
				}
				field, err := decodeJSONValue(dec)
				if err != nil {
					return nil, err
				}
				value.fields = append(value.fields, jsonField{key: key, value: field})
			}
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			return value, nil
		case '[':
			value := &jsonValue{kind: jsonArray}
			for dec.More() {
				element, err := decodeJSONValue(dec)
				if err != nil {
					return nil, err
				}
				value.elements = append(value.elements, element)
			}
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			return value, nil
		}
		return nil, fmt.Errorf("unexpected delimiter %q", token)
	case string:
		return &jsonValue{kind: jsonString, text: token}, nil
	case json.Number:
		return &jsonValue{kind: jsonNumber, number: token}, nil
	case bool:
		return &jsonValue{kind: jsonBool, boolean: token}, nil
	case nil:
		return &jsonValue{kind: jsonNull}, nil
	default:
		return nil, fmt.Errorf("unexpected JSON token %T", token)
	}
}

func (v *jsonValue) literal() (string, error) {
	switch v.kind {
	case jsonNull:
		return "null", nil
	case jsonBool:
		return strconv.FormatBool(v.boolean), nil
	case jsonString:
		return quoteString(v.text), nil
	case jsonNumber:
		return v.number.String(), nil
	case jsonArray:
		parts := make([]string, len(v.elements))
		for i, element := range v.elements {
			literal, err := element.literal()
			if err != nil {
				return "", err
			}
			parts[i] = literal
		}
		return "[" + strings.Join(parts, ", ") + "]", nil
	case jsonObject:
		return v.objectLiteral()
	}
	return "", fmt.Errorf("unknown JSON value kind %d", v.kind)
}

func (v *jsonValue) objectLiteral() (string, error) {
	if len(v.fields) == 1 {
		if literal, handled, err := extendedJSONLiteral(v.fields[0]); handled || err != nil {
			return literal, err
		}
	}
	parts := make([]string, len(v.fields))
	for i, field := range v.fields {
		literal, err := field.value.literal()
		if err != nil {
			return "", err
		}
		parts[i] = quoteString(field.key) + ": " + literal
	}
	return "{" + strings.Join(parts, ", ") + "}", nil
}

// extendedJSONLiteral renders a field that is one of Extended JSON's type
// wrappers. handled is false for a field that is an ordinary document field,
// including an aggregation operator such as $match or $gte.
func extendedJSONLiteral(field jsonField) (literal string, handled bool, err error) {
	switch field.key {
	case "$numberInt":
		text, ok := field.value.stringValue()
		if !ok {
			return "", true, fmt.Errorf("$numberInt is not a string")
		}
		number, err := strconv.ParseInt(text, 10, 32)
		if err != nil {
			return "", true, fmt.Errorf("$numberInt %q: %w", text, err)
		}
		return strconv.FormatInt(number, 10), true, nil
	case "$numberLong":
		text, ok := field.value.stringValue()
		if !ok {
			return "", true, fmt.Errorf("$numberLong is not a string")
		}
		number, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return "", true, fmt.Errorf("$numberLong %q: %w", text, err)
		}
		return "NumberLong(" + quoteString(strconv.FormatInt(number, 10)) + ")", true, nil
	case "$numberDouble":
		text, ok := field.value.stringValue()
		if !ok {
			return "", true, fmt.Errorf("$numberDouble is not a string")
		}
		literal, err := doubleLiteral(text)
		return literal, true, err
	case "$numberDecimal":
		text, ok := field.value.stringValue()
		if !ok {
			return "", true, fmt.Errorf("$numberDecimal is not a string")
		}
		return "NumberDecimal(" + quoteString(text) + ")", true, nil
	case "$oid":
		text, ok := field.value.stringValue()
		if !ok {
			return "", true, fmt.Errorf("$oid is not a string")
		}
		return "ObjectId(" + quoteString(text) + ")", true, nil
	case "$date":
		return dateLiteral(field.value)
	case "$timestamp":
		return timestampLiteral(field.value)
	case "$regularExpression":
		return regexLiteral(field.value)
	case "$binary", "$minKey", "$maxKey", "$symbol", "$code", "$scope", "$dbPointer", "$undefined", "$uuid":
		return "", true, fmt.Errorf("%s has no mongosh spelling gomongo can execute", field.key)
	}
	return "", false, nil
}

func (v *jsonValue) stringValue() (string, bool) {
	if v.kind != jsonString {
		return "", false
	}
	return v.text, true
}

func doubleLiteral(text string) (string, error) {
	switch text {
	case "Infinity", "-Infinity", "NaN":
		return "", fmt.Errorf("the double %s has no mongosh literal gomongo can execute", text)
	}
	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return "", fmt.Errorf("double %q: %w", text, err)
	}
	literal := strconv.FormatFloat(number, 'g', -1, 64)
	if !strings.ContainsAny(literal, ".eE") {
		literal += ".0"
	}
	return literal, nil
}

func dateLiteral(value *jsonValue) (string, bool, error) {
	if text, ok := value.stringValue(); ok {
		return "ISODate(" + quoteString(text) + ")", true, nil
	}
	if value.kind != jsonObject || len(value.fields) != 1 || value.fields[0].key != "$numberLong" {
		return "", false, nil
	}
	text, ok := value.fields[0].value.stringValue()
	if !ok {
		return "", true, fmt.Errorf("$date.$numberLong is not a string")
	}
	milliseconds, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return "", true, fmt.Errorf("$date.$numberLong %q: %w", text, err)
	}
	iso := time.UnixMilli(milliseconds).UTC().Format("2006-01-02T15:04:05.000Z")
	return "ISODate(" + quoteString(iso) + ")", true, nil
}

func timestampLiteral(value *jsonValue) (string, bool, error) {
	if value.kind != jsonObject {
		return "", true, fmt.Errorf("$timestamp is not a document")
	}
	seconds, hasSeconds, err := uint32Field(value, "t")
	if err != nil {
		return "", true, err
	}
	increment, hasIncrement, err := uint32Field(value, "i")
	if err != nil {
		return "", true, err
	}
	if !hasSeconds || !hasIncrement {
		return "", true, fmt.Errorf("$timestamp needs t and i")
	}
	return "Timestamp(" + strconv.FormatUint(uint64(seconds), 10) + ", " + strconv.FormatUint(uint64(increment), 10) + ")", true, nil
}

func uint32Field(doc *jsonValue, key string) (uint32, bool, error) {
	for _, field := range doc.fields {
		if field.key != key {
			continue
		}
		var text string
		switch {
		case field.value.kind == jsonNumber:
			text = field.value.number.String()
		case field.value.kind == jsonObject && len(field.value.fields) == 1:
			wrapper := field.value.fields[0]
			if wrapper.key != "$numberInt" && wrapper.key != "$numberLong" {
				return 0, false, fmt.Errorf("$timestamp.%s is not an integer", key)
			}
			var ok bool
			text, ok = wrapper.value.stringValue()
			if !ok {
				return 0, false, fmt.Errorf("$timestamp.%s is not an integer", key)
			}
		default:
			return 0, false, fmt.Errorf("$timestamp.%s is not an integer", key)
		}
		number, err := strconv.ParseUint(text, 10, 32)
		if err != nil {
			return 0, false, fmt.Errorf("$timestamp.%s %q: %w", key, text, err)
		}
		return uint32(number), true, nil
	}
	return 0, false, nil
}

func regexLiteral(value *jsonValue) (string, bool, error) {
	if value.kind != jsonObject {
		return "", true, fmt.Errorf("$regularExpression is not a document")
	}
	var pattern, options string
	for _, field := range value.fields {
		text, ok := field.value.stringValue()
		if !ok {
			return "", true, fmt.Errorf("$regularExpression.%s is not a string", field.key)
		}
		switch field.key {
		case "pattern":
			pattern = text
		case "options":
			options = text
		}
	}
	// A regex literal is the one spelling gomongo turns into a BSON regex.
	escaped := strings.ReplaceAll(pattern, `\/`, `/`)
	escaped = strings.ReplaceAll(escaped, `/`, `\/`)
	return "/" + escaped + "/" + options, true, nil
}

// quoteString renders s as a JavaScript string literal. JSON's string grammar
// is JavaScript's, so a document's strings survive the round trip.
func quoteString(s string) string {
	var b strings.Builder
	encoder := json.NewEncoder(&b)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(s); err != nil {
		return strconv.Quote(s)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

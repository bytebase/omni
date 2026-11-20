package catalog

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

type relOptType int

const (
	relOptBool relOptType = iota
	relOptInt
	relOptReal
	relOptEnum
)

type relOptSpec struct {
	typ      relOptType
	minInt   int64
	maxInt   int64
	minReal  float64
	maxReal  float64
	enumVals []string
}

var heapRelOpts = map[string]relOptSpec{
	"autovacuum_enabled":                    {typ: relOptBool},
	"user_catalog_table":                    {typ: relOptBool},
	"vacuum_truncate":                       {typ: relOptBool},
	"fillfactor":                            {typ: relOptInt, minInt: 10, maxInt: 100},
	"toast_tuple_target":                    {typ: relOptInt, minInt: 128, maxInt: 8160},
	"parallel_workers":                      {typ: relOptInt, minInt: 0, maxInt: 1024},
	"autovacuum_vacuum_threshold":           {typ: relOptInt, minInt: -1, maxInt: math.MaxInt32},
	"autovacuum_vacuum_insert_threshold":    {typ: relOptInt, minInt: -2, maxInt: math.MaxInt32},
	"autovacuum_analyze_threshold":          {typ: relOptInt, minInt: -1, maxInt: math.MaxInt32},
	"autovacuum_vacuum_cost_limit":          {typ: relOptInt, minInt: -1, maxInt: 10000},
	"autovacuum_freeze_min_age":             {typ: relOptInt, minInt: 0, maxInt: 1000000000},
	"autovacuum_freeze_max_age":             {typ: relOptInt, minInt: 100000, maxInt: 2000000000},
	"autovacuum_freeze_table_age":           {typ: relOptInt, minInt: 0, maxInt: 2000000000},
	"autovacuum_multixact_freeze_min_age":   {typ: relOptInt, minInt: 0, maxInt: 1000000000},
	"autovacuum_multixact_freeze_max_age":   {typ: relOptInt, minInt: 10000, maxInt: 2000000000},
	"autovacuum_multixact_freeze_table_age": {typ: relOptInt, minInt: 0, maxInt: 2000000000},
	"log_autovacuum_min_duration":           {typ: relOptInt, minInt: -1, maxInt: math.MaxInt32},
	"autovacuum_vacuum_cost_delay":          {typ: relOptReal, minReal: -1, maxReal: 100},
	"autovacuum_vacuum_scale_factor":        {typ: relOptReal, minReal: 0, maxReal: 100},
	"autovacuum_vacuum_insert_scale_factor": {typ: relOptReal, minReal: 0, maxReal: 100},
	"autovacuum_analyze_scale_factor":       {typ: relOptReal, minReal: 0, maxReal: 100},
	"vacuum_index_cleanup":                  {typ: relOptEnum, enumVals: []string{"auto", "on", "off", "true", "false", "yes", "no", "1", "0"}},
}

var toastRelOpts = map[string]relOptSpec{
	"autovacuum_enabled":                    {typ: relOptBool},
	"vacuum_truncate":                       {typ: relOptBool},
	"autovacuum_vacuum_threshold":           {typ: relOptInt, minInt: -1, maxInt: math.MaxInt32},
	"autovacuum_vacuum_insert_threshold":    {typ: relOptInt, minInt: -2, maxInt: math.MaxInt32},
	"autovacuum_vacuum_cost_limit":          {typ: relOptInt, minInt: -1, maxInt: 10000},
	"autovacuum_freeze_min_age":             {typ: relOptInt, minInt: 0, maxInt: 1000000000},
	"autovacuum_freeze_max_age":             {typ: relOptInt, minInt: 100000, maxInt: 2000000000},
	"autovacuum_freeze_table_age":           {typ: relOptInt, minInt: 0, maxInt: 2000000000},
	"autovacuum_multixact_freeze_min_age":   {typ: relOptInt, minInt: 0, maxInt: 1000000000},
	"autovacuum_multixact_freeze_max_age":   {typ: relOptInt, minInt: 10000, maxInt: 2000000000},
	"autovacuum_multixact_freeze_table_age": {typ: relOptInt, minInt: 0, maxInt: 2000000000},
	"log_autovacuum_min_duration":           {typ: relOptInt, minInt: -1, maxInt: math.MaxInt32},
	"autovacuum_vacuum_cost_delay":          {typ: relOptReal, minReal: -1, maxReal: 100},
	"autovacuum_vacuum_scale_factor":        {typ: relOptReal, minReal: 0, maxReal: 100},
	"autovacuum_vacuum_insert_scale_factor": {typ: relOptReal, minReal: 0, maxReal: 100},
	"vacuum_index_cleanup":                  {typ: relOptEnum, enumVals: []string{"auto", "on", "off", "true", "false", "yes", "no", "1", "0"}},
}

var viewRelOpts = map[string]relOptSpec{
	"security_barrier": {typ: relOptBool},
	"security_invoker": {typ: relOptBool},
	"check_option":     {typ: relOptEnum, enumVals: []string{"local", "cascaded"}},
}

var indexRelOptsByAM = map[string]map[string]relOptSpec{
	"btree": {
		"fillfactor":                        {typ: relOptInt, minInt: 10, maxInt: 100},
		"deduplicate_items":                 {typ: relOptBool},
		"vacuum_cleanup_index_scale_factor": {typ: relOptReal, minReal: 0, maxReal: 1e10},
	},
	"hash": {
		"fillfactor": {typ: relOptInt, minInt: 10, maxInt: 100},
	},
	"gist": {
		"fillfactor": {typ: relOptInt, minInt: 10, maxInt: 100},
		"buffering":  {typ: relOptEnum, enumVals: []string{"auto", "on", "off"}},
	},
	"spgist": {
		"fillfactor": {typ: relOptInt, minInt: 10, maxInt: 100},
	},
	"gin": {
		"fastupdate":             {typ: relOptBool},
		"gin_pending_list_limit": {typ: relOptInt, minInt: 64, maxInt: math.MaxInt32},
	},
	"brin": {
		"pages_per_range": {typ: relOptInt, minInt: 1, maxInt: 131072},
		"autosummarize":   {typ: relOptBool},
	},
}

var attrRelOpts = map[string]relOptSpec{
	"n_distinct":           {typ: relOptReal, minReal: -1, maxReal: math.MaxFloat64},
	"n_distinct_inherited": {typ: relOptReal, minReal: -1, maxReal: math.MaxFloat64},
}

const (
	// toastTupleThreshold is PostgreSQL's TOAST_TUPLE_THRESHOLD on 8KB pages:
	// MaximumBytesPerTuple(TOAST_TUPPLES_PER_PAGE=4) = 2032.
	toastTupleThreshold   = 2032
	sizeofHeapTupleHeader = 23
	utf8MaxMBLen          = 4
)

// columnsNeedToastTable reports whether a heap relation with cols creates a
// TOAST table in PostgreSQL: at least one column has non-plain storage ('x',
// 'e', or 'm') and the maximum possible tuple width (accounting for bounded
// varlena type modifiers) exceeds TOAST_TUPLE_THRESHOLD.
//
// pg: src/backend/access/heap/heapam_handler.c — heapam_relation_needs_toast_table
func columnsNeedToastTable(cols []*Column) bool {
	var dataLen int64
	hasToastable := false
	unbounded := false
	hasNull := false

	for _, col := range cols {
		if !col.NotNull {
			hasNull = true
		}
		dataLen = alignNominal(dataLen, col.Align)
		if col.Len > 0 {
			dataLen += int64(col.Len)
		} else if col.Len == -1 {
			if col.Storage == 'x' || col.Storage == 'e' || col.Storage == 'm' {
				hasToastable = true
			}
			maxWidth := typeMaximumSize(col.TypeOID, col.TypeMod)
			if maxWidth < 0 {
				unbounded = true
			} else {
				dataLen += maxWidth
			}
		} else {
			// Null-terminated cstring (Len == -2): unbounded.
			unbounded = true
		}
	}
	if !hasToastable {
		return false
	}
	if unbounded {
		return true
	}
	headerLen := int64(sizeofHeapTupleHeader)
	if hasNull {
		headerLen += int64((len(cols) + 7) / 8)
	}
	headerLen = (headerLen + 7) &^ 7
	tupleLen := headerLen + ((dataLen + 7) &^ 7)
	return tupleLen > toastTupleThreshold
}

func alignNominal(offset int64, align byte) int64 {
	switch align {
	case 'd':
		return (offset + 7) &^ 7
	case 'i':
		return (offset + 3) &^ 3
	case 's':
		return (offset + 1) &^ 1
	default:
		return offset
	}
}

// typeMaximumSize mirrors PostgreSQL's type_maximum_size (UTF-8 maxmblen=4),
// returning the maximum on-disk varlena size in bytes or -1 if unbounded.
// Note that PostgreSQL does not unwrap domain types here.
//
// pg: src/backend/utils/adt/format_type.c — type_maximum_size
func typeMaximumSize(typeOID uint32, typemod int32) int64 {
	if typemod < 0 {
		return -1
	}
	switch typeOID {
	case BPCHAROID, VARCHAROID:
		if typemod < 4 {
			return -1
		}
		return int64(typemod-4)*utf8MaxMBLen + 4
	case BITOID, VARBITOID:
		if typemod <= 0 {
			return -1
		}
		return (int64(typemod)+7)/8 + 8
	case NUMERICOID:
		if typemod < 4 {
			return -1
		}
		precision := int64(((typemod - 4) >> 16) & 0xffff)
		return (((precision + 1 + 3) / 4) * 2) + 12
	default:
		return -1
	}
}

// validateRelOptions validates storage parameters (WITH (...) / SET (...)) for
// tables, views, materialized views, and indexes.
//
// pg: src/backend/access/common/reloptions.c — transformRelOptions / parseRelOptions
func (c *Catalog) validateRelOptions(opts *nodes.List, relkind byte, accessMethod string) error {
	return c.validateRelOptionsForCols(opts, relkind, accessMethod, true)
}

func (c *Catalog) validateRelOptionsForCols(opts *nodes.List, relkind byte, accessMethod string, hasToast bool) error {
	if opts == nil || len(opts.Items) == 0 {
		return nil
	}
	switch relkind {
	case 'p':
		return &Error{Code: CodeWrongObjectType, Message: "cannot specify storage parameters for a partitioned table"}
	case 'c', 'f', 'S':
		return &Error{Code: CodeWrongObjectType, Message: "cannot specify storage parameters for this relation"}
	}

	seen := make(map[string]bool, len(opts.Items))
	for _, item := range opts.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		ns := strings.ToLower(d.Defnamespace)
		name := strings.ToLower(d.Defname)
		key := ns + "\x00" + name
		if seen[key] {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q specified more than once", d.Defname))
		}
		seen[key] = true

		var table map[string]relOptSpec
		switch relkind {
		case 'r', 'm':
			if ns == "" {
				table = heapRelOpts
			} else if ns == "toast" {
				if !hasToast {
					// Relations without a TOAST table ignore toast.* parameters.
					continue
				}
				table = toastRelOpts
			} else {
				return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter namespace %q", d.Defnamespace))
			}
		case 'v':
			if ns != "" {
				return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter namespace %q", d.Defnamespace))
			}
			table = viewRelOpts
		case 'i', 'I':
			if ns != "" {
				return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter namespace %q", d.Defnamespace))
			}
			amName := accessMethod
			if amName == "" {
				amName = "btree"
			}
			if am := c.accessMethodByName[amName]; (am != nil && am.OID >= FirstNormalObjectId) || (am == nil && c.snapshot) {
				continue
			}
			table = indexRelOptsByAM[amName]
		default:
			return &Error{Code: CodeWrongObjectType, Message: "cannot specify storage parameters for this relation"}
		}

		spec, exists := table[name]
		if !exists {
			return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter %q", d.Defname))
		}
		if err := validateRelOptValue(d, spec); err != nil {
			return err
		}
	}
	return nil
}

// validateResetRelOptions validates ALTER ... RESET (...) parameters.
//
// pg: src/backend/access/common/reloptions.c — transformRelOptions (isReset=true)
// Note: In PostgreSQL 17, transformRelOptions with isReset=true only removes matching
// names from existing Datum reloptions and never appends reset entries, so
// parseRelOptionsInternal does not reject unknown option names on RESET
// (verified against PostgreSQL 17 in testdata/simulate.tsv:366,637,908 and
// TestSimulationCausesOracle.tape for ALTER_TABLE_RESET_unknown_reloption_on_index).
func (c *Catalog) validateResetRelOptions(opts *nodes.List, relkind byte) error {
	switch relkind {
	case 'c', 'f', 'S':
		return &Error{Code: CodeWrongObjectType, Message: "cannot reset storage parameters for this relation"}
	}
	if opts == nil {
		return nil
	}
	seen := make(map[string]bool, len(opts.Items))
	for _, item := range opts.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		if d.Arg != nil {
			return &Error{Code: CodeSyntaxError, Message: "RESET must not include values for parameters"}
		}
		ns := strings.ToLower(d.Defnamespace)
		name := strings.ToLower(d.Defname)
		key := ns + "\x00" + name
		if seen[key] {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q specified more than once", d.Defname))
		}
		seen[key] = true
		if ns != "" && !(ns == "toast" && (relkind == 'r' || relkind == 'm' || relkind == 'p')) {
			return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter namespace %q", d.Defnamespace))
		}
	}
	return nil
}

// validateAttributeRelOptions validates ALTER TABLE ... ALTER COLUMN ... SET/RESET (...).
//
// pg: src/backend/commands/tablecmds.c — ATExecSetOptions
func (c *Catalog) validateAttributeRelOptions(opts *nodes.List, isReset bool) error {
	if opts == nil {
		return nil
	}
	seen := make(map[string]bool, len(opts.Items))
	for _, item := range opts.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		if isReset && d.Arg != nil {
			return &Error{Code: CodeSyntaxError, Message: "RESET must not include values for parameters"}
		}
		ns := strings.ToLower(d.Defnamespace)
		name := strings.ToLower(d.Defname)
		key := ns + "\x00" + name
		if seen[key] {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q specified more than once", d.Defname))
		}
		seen[key] = true
		if ns != "" {
			return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter namespace %q", d.Defnamespace))
		}
		if isReset {
			continue
		}
		spec, exists := attrRelOpts[name]
		if !exists {
			return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter %q", d.Defname))
		}
		if err := validateRelOptValue(d, spec); err != nil {
			return err
		}
	}
	return nil
}

func validateRelOptValue(d *nodes.DefElem, spec relOptSpec) error {
	switch spec.typ {
	case relOptBool:
		if d.Arg == nil {
			return nil
		}
		if _, ok := parseDefElemBoolStrict(d); !ok {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q requires a Boolean value", d.Defname))
		}
		return nil
	case relOptInt:
		if d.Arg == nil {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q requires an integer value", d.Defname))
		}
		v, ok := parseDefElemIntStrict(d)
		if !ok {
			return errInvalidParameterValue(fmt.Sprintf("invalid value for integer option %q", d.Defname))
		}
		if v < spec.minInt || v > spec.maxInt {
			return errInvalidParameterValue(fmt.Sprintf("value %d out of bounds for option %q", v, d.Defname))
		}
		return nil
	case relOptReal:
		if d.Arg == nil {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q requires a numeric value", d.Defname))
		}
		v, ok := parseDefElemFloatStrict(d)
		if !ok {
			return errInvalidParameterValue(fmt.Sprintf("invalid value for floating point option %q", d.Defname))
		}
		if v < spec.minReal || v > spec.maxReal {
			return errInvalidParameterValue(fmt.Sprintf("value %g out of bounds for option %q", v, d.Defname))
		}
		return nil
	case relOptEnum:
		if d.Arg == nil {
			return errInvalidParameterValue(fmt.Sprintf("parameter %q requires a string value", d.Defname))
		}
		s, ok := parseDefElemStringStrict(d)
		if !ok {
			return errInvalidParameterValue(fmt.Sprintf("invalid value for enum option %q", d.Defname))
		}
		sl := strings.ToLower(s)
		for _, allowed := range spec.enumVals {
			if sl == allowed {
				return nil
			}
		}
		return errInvalidParameterValue(fmt.Sprintf("unrecognized %s value: %q", d.Defname, s))
	}
	return nil
}

func parseDefElemStringStrict(d *nodes.DefElem) (string, bool) {
	if d == nil || d.Arg == nil {
		return "", false
	}
	switch v := d.Arg.(type) {
	case *nodes.String:
		return v.Str, true
	case *nodes.Integer:
		return strconv.FormatInt(v.Ival, 10), true
	case *nodes.Float:
		return v.Fval, true
	case *nodes.Boolean:
		if v.Boolval {
			return "true", true
		}
		return "false", true
	case *nodes.TypeName:
		_, name := typeNameParts(v)
		if name != "" {
			return name, true
		}
	case *nodes.A_Const:
		if v.Isnull {
			return "", false
		}
		return deparseAConst(v), true
	}
	return "", false
}

func parseDefElemBoolStrict(d *nodes.DefElem) (bool, bool) {
	if b, ok := d.Arg.(*nodes.Boolean); ok {
		return b.Boolval, true
	}
	if ac, ok := d.Arg.(*nodes.A_Const); ok && !ac.Isnull {
		if b, ok := ac.Val.(*nodes.Boolean); ok {
			return b.Boolval, true
		}
	}
	s, ok := parseDefElemStringStrict(d)
	if !ok {
		return false, false
	}
	return parsePGBool(s)
}

// parsePGBool reads a Boolean as PostgreSQL's parse_bool does (bool.c): any
// prefix of true, false, yes or no, a prefix of on or off at least two
// characters long, or exactly 1 or 0, case-insensitively and untrimmed.
func parsePGBool(s string) (bool, bool) {
	v := strings.ToLower(s)
	if v == "" {
		return false, false
	}
	switch v[0] {
	case 't':
		return true, strings.HasPrefix("true", v)
	case 'f':
		return false, strings.HasPrefix("false", v)
	case 'y':
		return true, strings.HasPrefix("yes", v)
	case 'n':
		return false, strings.HasPrefix("no", v)
	case 'o':
		if len(v) < 2 {
			return false, false
		}
		if strings.HasPrefix("on", v) {
			return true, true
		}
		return false, strings.HasPrefix("off", v)
	case '1':
		return true, v == "1"
	case '0':
		return false, v == "0"
	}
	return false, false
}

func parseDefElemIntStrict(d *nodes.DefElem) (int64, bool) {
	if i, ok := defElemInt(d); ok {
		return i, true
	}
	s, ok := parseDefElemStringStrict(d)
	if !ok {
		return 0, false
	}
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func parseDefElemFloatStrict(d *nodes.DefElem) (float64, bool) {
	if i, ok := defElemInt(d); ok {
		return float64(i), true
	}
	s, ok := parseDefElemStringStrict(d)
	if !ok {
		return 0, false
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

// resolveColumnStorage resolves a column STORAGE directive against typ.
//
// pg: src/backend/commands/tablecmds.c — GetAttributeStorage
func (c *Catalog) resolveColumnStorage(typ *BuiltinType, storageName string) (byte, error) {
	defStorage := byte('p')
	typeName := "unknown"
	if typ != nil {
		defStorage = typ.Storage
		typeName = typ.TypeName
	}
	if storageName == "" {
		return defStorage, nil
	}
	var stor byte
	switch strings.ToLower(storageName) {
	case "default":
		stor = defStorage
	case "plain":
		stor = 'p'
	case "external":
		stor = 'e'
	case "extended":
		stor = 'x'
	case "main":
		stor = 'm'
	default:
		return 0, errInvalidParameterValue(fmt.Sprintf("invalid storage type %q", storageName))
	}
	if defStorage == 'p' && stor != 'p' {
		return 0, &Error{
			Code:    CodeFeatureNotSupported,
			Message: fmt.Sprintf("column data type %s can only have storage PLAIN", typeName),
		}
	}
	return stor, nil
}

// resolveColumnCompression resolves a column COMPRESSION directive against typ.
// PostgreSQL checks type compressibility (0A000) before validating the method
// name (22023).
//
// pg: src/backend/commands/tablecmds.c — GetAttributeCompression
func (c *Catalog) resolveColumnCompression(typ *BuiltinType, compName string) (byte, error) {
	if compName == "" {
		return 0, nil
	}
	stor := byte('p')
	typeName := "unknown"
	if typ != nil {
		stor = typ.Storage
		typeName = typ.TypeName
	}
	if stor != 'x' && stor != 'm' {
		return 0, &Error{
			Code:    CodeFeatureNotSupported,
			Message: fmt.Sprintf("column data type %s does not support compression", typeName),
		}
	}
	switch strings.ToLower(compName) {
	case "default":
		return 0, nil
	case "pglz":
		return 'p', nil
	case "lz4":
		return 'l', nil
	default:
		return 0, errInvalidParameterValue(fmt.Sprintf("invalid compression method %q", compName))
	}
}

// checkCollationList validates a qualified or unqualified collation name list.
//
// pg: src/backend/parser/parse_collate.c — LookupCollation
func (c *Catalog) checkCollationList(collname *nodes.List) (string, error) {
	name, _, err := c.resolveCollationList(collname)
	return name, err
}

func (c *Catalog) resolveCollationList(collname *nodes.List) (string, uint32, error) {
	if collname == nil || len(collname.Items) == 0 {
		return "", 0, nil
	}
	schemaName, name := qualifiedName(collname)
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return "", 0, errUndefinedSchema(schemaName)
		}
		if c.snapshot {
			var uid uint32
			if uc := s.Collations[name]; uc != nil {
				uid = uc.OID
			}
			return name, uid, nil
		}
		if s.OID == PGCatalogNamespace && isBuiltinCollationName(name) {
			return name, 0, nil
		}
		if uc := s.Collations[name]; uc != nil {
			return name, uc.OID, nil
		}
		return "", 0, errUndefinedObject("collation", name)
	}
	if !c.snapshot {
		for _, nsOID := range c.searchPathWithCatalog() {
			s := c.schemas[nsOID]
			if s == nil {
				continue
			}
			if nsOID == PGCatalogNamespace && isBuiltinCollationName(name) {
				return name, 0, nil
			}
			if uc := s.Collations[name]; uc != nil {
				return name, uc.OID, nil
			}
		}
		if isBuiltinCollationName(name) {
			return name, 0, nil
		}
		return "", 0, errUndefinedObject("collation", name)
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		if s := c.schemas[nsOID]; s != nil {
			if uc := s.Collations[name]; uc != nil {
				return name, uc.OID, nil
			}
		}
	}
	return name, 0, nil
}

func isBuiltinCollationName(name string) bool {
	switch name {
	case "default", "C", "POSIX", "ucs_basic", "pg_c_utf8", "unicode", "C.UTF-8", "C.utf8", "en_US", "en_US.utf8", "en_US.UTF-8":
		return true
	}
	return strings.HasSuffix(strings.ToLower(name), "-x-icu")
}

// checkColumnCollation validates an explicit COLLATE clause on a column or
// index key element: collation existence (42704) is checked before checking
// whether the target type supports collations (42804).
//
// pg: src/backend/commands/tablecmds.c — GetColumnDefCollation
func (c *Catalog) checkColumnCollation(collname *nodes.List, typeOID uint32) (string, error) {
	name, _, err := c.resolveColumnCollation(collname, typeOID)
	return name, err
}

func (c *Catalog) resolveColumnCollation(collname *nodes.List, typeOID uint32) (string, uint32, error) {
	if collname == nil || len(collname.Items) == 0 {
		return "", 0, nil
	}
	name, userOID, err := c.resolveCollationList(collname)
	if err != nil {
		return "", 0, err
	}
	if typeOID != 0 && c.typeCollation(typeOID) == 0 {
		tname := "unknown"
		if bt := c.typeByOID[typeOID]; bt != nil {
			tname = bt.TypeName
		}
		return "", 0, errDatatypeMismatch(fmt.Sprintf("collations are not supported by type %s", tname))
	}
	return name, userOID, nil
}

func (c *Catalog) resolvedCollationSchema(collname *nodes.List, leafName string, userOID uint32) string {
	if userOID != 0 {
		for _, s := range c.schemas {
			if s != nil {
				if uc := s.Collations[leafName]; uc != nil && uc.OID == userOID {
					if s.Name != "public" || isBuiltinCollationName(leafName) {
						return s.Name
					}
					if collname != nil {
						if sc, _ := qualifiedName(collname); sc == "public" {
							return "public"
						}
					}
					return ""
				}
			}
		}
	}
	if collname != nil {
		sc, _ := qualifiedName(collname)
		if sc != "" {
			return sc
		}
	}
	return ""
}

func (c *Catalog) effectiveColumnCollation(typeOID uint32, collName string, userOID uint32) uint32 {
	typeColl := c.typeCollation(typeOID)
	if userOID != 0 {
		return userOID
	}
	if collName != "" && typeColl != 0 {
		if oid := builtinCollationOID(collName); oid != 0 {
			return oid
		}
	}
	return typeColl
}

// checkCompositeSelfReference refuses making a composite row type a direct or
// transitive member of itself (42P16).
//
// pg: src/backend/commands/tablecmds.c — find_composite_type_dependencies
func (c *Catalog) checkCompositeSelfReference(targetRowTypeOID, colTypeOID uint32, relName string) error {
	if targetRowTypeOID == 0 || colTypeOID == 0 {
		return nil
	}
	visited := map[uint32]bool{}
	var contains func(oid uint32) bool
	contains = func(oid uint32) bool {
		if oid == 0 || visited[oid] {
			return false
		}
		if oid == targetRowTypeOID {
			return true
		}
		visited[oid] = true
		bt := c.typeByOID[oid]
		if bt == nil {
			return false
		}
		if bt.Elem != 0 {
			if contains(bt.Elem) {
				return true
			}
		}
		if bt.Type == 'd' {
			if dt := c.domainTypes[oid]; dt != nil && contains(dt.BaseTypeOID) {
				return true
			}
		}
		if bt.Type == 'r' {
			if rt := c.rangeTypes[oid]; rt != nil && contains(rt.SubTypeOID) {
				return true
			}
		}
		if bt.Type == 'c' && bt.RelID != 0 {
			if rel := c.relationByOID[bt.RelID]; rel != nil {
				for _, col := range rel.Columns {
					if contains(col.TypeOID) {
						return true
					}
				}
			}
		}
		return false
	}
	if contains(colTypeOID) {
		return &Error{
			Code:    CodeInvalidTableDefinition,
			Message: fmt.Sprintf("composite type %s cannot be made a member of itself", relName),
		}
	}
	return nil
}

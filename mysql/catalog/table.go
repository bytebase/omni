package catalog

type Table struct {
	Name          string
	Database      *Database
	Columns       []*Column
	colByName     map[string]int // lowered name -> index
	Indexes       []*Index
	Constraints   []*Constraint
	Engine        string
	Charset       string
	Collation     string
	Comment       string
	AutoIncrement int64
	Temporary     bool
	RowFormat     string
	KeyBlockSize  int
	Partitioning  *PartitionInfo
}

// PartitionInfo holds partition metadata for a table.
type PartitionInfo struct {
	Type       string // RANGE, LIST, HASH, KEY
	Linear     bool   // LINEAR HASH or LINEAR KEY
	Expr       string // partition expression (for RANGE/LIST/HASH)
	Columns    []string // partition columns (for RANGE COLUMNS/LIST COLUMNS/KEY)
	Algorithm  int    // ALGORITHM={1|2} for KEY partitioning
	NumParts   int    // PARTITIONS num
	Partitions []*PartitionDefInfo
	SubType    string // subpartition type (HASH or KEY, "" if none)
	SubLinear  bool   // LINEAR for subpartition
	SubExpr    string // subpartition expression
	SubColumns []string // subpartition columns
	SubAlgo    int    // subpartition ALGORITHM
	NumSubParts int   // SUBPARTITIONS num
}

// PartitionDefInfo holds a single partition definition.
type PartitionDefInfo struct {
	Name          string
	ValueExpr     string // "LESS THAN (...)" or "IN (...)" or ""
	Engine        string // ENGINE option for this partition
	Comment       string // COMMENT option for this partition
	SubPartitions []*SubPartitionDefInfo
}

// SubPartitionDefInfo holds a single subpartition definition.
type SubPartitionDefInfo struct {
	Name    string
	Engine  string
	Comment string
}

type Column struct {
	Position       int
	Name           string
	DataType       string // normalized (int, varchar, etc.)
	ColumnType     string // full type string (varchar(100), int unsigned)
	Nullable       bool
	Default        *string
	DefaultDropped bool // true when ALTER COLUMN DROP DEFAULT was used
	AutoIncrement  bool
	Charset        string
	Collation      string
	Comment        string
	OnUpdate       string
	Generated      *GeneratedColumnInfo
	Invisible      bool
	SRID           int // Spatial Reference ID (0 = not set)
}

type GeneratedColumnInfo struct {
	Expr   string
	Stored bool
}

type View struct {
	Name            string
	Database        *Database
	Definition      string
	Algorithm       string
	Definer         string
	SqlSecurity     string
	CheckOption     string
	Columns         []string // All column names (explicit or derived from SELECT)
	ExplicitColumns bool     // true if the user specified a column list in CREATE VIEW
}

// Routine represents a stored function or procedure in the catalog.
type Routine struct {
	Name            string
	Database        *Database
	IsProcedure     bool
	Definer         string
	Params          []*RoutineParam
	Returns         string // return type string for functions (empty for procedures)
	Body            string
	Characteristics map[string]string // name -> value (DETERMINISTIC, COMMENT, etc.)
}

// RoutineParam represents a parameter of a stored routine.
type RoutineParam struct {
	Direction string // IN, OUT, INOUT (empty for functions)
	Name      string
	TypeName  string // full type string
}

// Trigger represents a trigger in the catalog.
type Trigger struct {
	Name     string
	Database *Database
	Table    string // table name the trigger is on
	Timing   string // BEFORE, AFTER
	Event    string // INSERT, UPDATE, DELETE
	Definer  string
	Body     string
	Order    *TriggerOrderInfo
}

// TriggerOrderInfo represents FOLLOWS/PRECEDES ordering.
type TriggerOrderInfo struct {
	Follows     bool
	TriggerName string
}

// Event represents a scheduled event in the catalog.
type Event struct {
	Name         string
	Database     *Database
	Definer      string
	Schedule     string // raw schedule text (e.g. "EVERY 1 HOUR", "AT '2024-01-01 00:00:00'")
	OnCompletion string // PRESERVE, NOT PRESERVE, or "" (default NOT PRESERVE)
	Enable       string // ENABLE, DISABLE, DISABLE ON SLAVE, or "" (default ENABLE)
	Comment      string
	Body         string
}

// cloneTable returns a deep copy of the table's mutable state.
// The returned Table shares the same Name, Database pointer, and scalar fields,
// but has independent slices and maps so that mutations do not affect the original.
func cloneTable(src *Table) Table {
	dst := *src // shallow copy of all scalar fields

	// Deep copy columns.
	dst.Columns = make([]*Column, len(src.Columns))
	for i, sc := range src.Columns {
		col := *sc
		if sc.Default != nil {
			def := *sc.Default
			col.Default = &def
		}
		if sc.Generated != nil {
			gen := *sc.Generated
			col.Generated = &gen
		}
		dst.Columns[i] = &col
	}

	// Deep copy colByName.
	dst.colByName = make(map[string]int, len(src.colByName))
	for k, v := range src.colByName {
		dst.colByName[k] = v
	}

	// Deep copy indexes.
	dst.Indexes = make([]*Index, len(src.Indexes))
	for i, si := range src.Indexes {
		idx := *si
		idx.Table = src // keep pointing to the original table pointer
		cols := make([]*IndexColumn, len(si.Columns))
		for j, sc := range si.Columns {
			ic := *sc
			cols[j] = &ic
		}
		idx.Columns = cols
		dst.Indexes[i] = &idx
	}

	// Deep copy constraints.
	dst.Constraints = make([]*Constraint, len(src.Constraints))
	for i, sc := range src.Constraints {
		con := *sc
		con.Table = src
		con.Columns = append([]string{}, sc.Columns...)
		con.RefColumns = append([]string{}, sc.RefColumns...)
		dst.Constraints[i] = &con
	}

	// Deep copy partitioning.
	if src.Partitioning != nil {
		pi := *src.Partitioning
		pi.Columns = append([]string{}, src.Partitioning.Columns...)
		pi.SubColumns = append([]string{}, src.Partitioning.SubColumns...)
		pi.Partitions = make([]*PartitionDefInfo, len(src.Partitioning.Partitions))
		for i, sp := range src.Partitioning.Partitions {
			pd := *sp
			pd.SubPartitions = make([]*SubPartitionDefInfo, len(sp.SubPartitions))
			for j, ss := range sp.SubPartitions {
				sd := *ss
				pd.SubPartitions[j] = &sd
			}
			pi.Partitions[i] = &pd
		}
		dst.Partitioning = &pi
	}

	return dst
}

func (t *Table) GetColumn(name string) *Column {
	idx, ok := t.colByName[toLower(name)]
	if !ok {
		return nil
	}
	return t.Columns[idx]
}

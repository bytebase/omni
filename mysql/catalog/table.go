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
}

type Column struct {
	Position      int
	Name          string
	DataType      string // normalized (int, varchar, etc.)
	ColumnType    string // full type string (varchar(100), int unsigned)
	Nullable      bool
	Default       *string
	AutoIncrement bool
	Charset       string
	Collation     string
	Comment       string
	OnUpdate      string
	Generated     *GeneratedColumnInfo
	Invisible     bool
}

type GeneratedColumnInfo struct {
	Expr   string
	Stored bool
}

type View struct {
	Name        string
	Database    *Database
	Definition  string
	Algorithm   string
	Definer     string
	SqlSecurity string
	CheckOption string
	Columns     []string
}

func (t *Table) GetColumn(name string) *Column {
	idx, ok := t.colByName[toLower(name)]
	if !ok {
		return nil
	}
	return t.Columns[idx]
}

package catalog

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

const (
	minAutoextendSizeBytes = 4 * 1024 * 1024    // 4M
	maxAutoextendSizeBytes = 4096 * 1024 * 1024 // 4096M (4G)
)

// validateAutoextendSize validates an AUTOEXTEND_SIZE option value as InnoDB
// enforces it on mysql:8.0.46: 0 disables custom autoextend, and any non-zero
// value must lie in [4M, 4096M] (4025) and be an exact multiple of 4M (4023).
func validateAutoextendSize(raw string) error {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil
	}
	mult := uint64(1)
	switch s[len(s)-1] {
	case 'k', 'K':
		mult = 1024
		s = s[:len(s)-1]
	case 'm', 'M':
		mult = 1024 * 1024
		s = s[:len(s)-1]
	case 'g', 'G':
		mult = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return errAutoextendSizeOutOfRange()
	}
	if val == 0 {
		return nil
	}
	if val > math.MaxUint64/mult {
		return errAutoextendSizeOutOfRange()
	}
	bytes := val * mult
	if bytes < minAutoextendSizeBytes || bytes > maxAutoextendSizeBytes {
		return errAutoextendSizeOutOfRange()
	}
	if bytes%minAutoextendSizeBytes != 0 {
		return errAutoextendSizeNotMultiple()
	}
	return nil
}

// validateInnoDBTablespace validates a TABLESPACE option on an InnoDB table or
// partition (measured on mysql:8.0.46). Non-InnoDB engines ignore TABLESPACE.
// General tablespace names are instance-wide and case-sensitive.
func (c *Catalog) validateInnoDBTablespace(tableName string, temporary bool, rawTS string) error {
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if ts == "" {
		return nil
	}
	switch {
	case ts == "innodb_file_per_table" || ts == "innodb_system":
		return nil
	case ts == "innodb_temporary":
		if temporary {
			return nil
		}
		return errTempTablespaceNonTempTable()
	case ts == "mysql":
		return errReservedTablespaceName(tableName, "mysql")
	case strings.HasPrefix(strings.ToLower(ts), "innodb_"):
		return errReservedInnoDBTablespacePrefix()
	case c.tablespaces[ts]:
		return nil
	case c.loadingSDL:
		// SDL captures schema-scoped DDL without instance-wide CREATE TABLESPACE
		// statements; accept and register general tablespaces referenced by SDL tables.
		c.tablespaces[ts] = true
		return nil
	default:
		return errTablespaceDoesNotExist(ts)
	}
}

func validateInnoDBPartitionedSharedTablespace(rawTS string) error {
	ts := unquoteTableOption(strings.TrimSpace(rawTS))
	if ts != "" && ts != "innodb_file_per_table" {
		return errPartitionedTableSharedTablespace()
	}
	return nil
}

func effectivePartEngine(tblEngine, partEngine, subEngine string) string {
	if subEngine != "" {
		return subEngine
	}
	if partEngine != "" {
		return partEngine
	}
	return tblEngine
}

// tableUsesTablespace reports whether tbl occupies general tablespace tsName in
// InnoDB. On mysql:8.0.46, non-InnoDB engines such as MyISAM retain a syntactic
// TABLESPACE clause in SHOW CREATE TABLE without occupying the InnoDB tablespace,
// so DROP TABLESPACE succeeds while only a non-InnoDB table names it.
func tableUsesTablespace(tbl *Table, tsName string) bool {
	if tbl == nil {
		return false
	}
	if isInnoDBEngine(tbl.Engine) && unquoteTableOption(tbl.Tablespace) == tsName {
		return true
	}
	if tbl.Partitioning != nil {
		for _, pd := range tbl.Partitioning.Partitions {
			if pd == nil {
				continue
			}
			if isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, "")) && unquoteTableOption(pd.Tablespace) == tsName {
				return true
			}
			for _, spd := range pd.SubPartitions {
				if spd != nil && isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, spd.Engine)) && unquoteTableOption(spd.Tablespace) == tsName {
					return true
				}
			}
		}
	}
	return false
}

func renameTableTablespaceRefs(tbl *Table, oldName, newName string) {
	if tbl == nil {
		return
	}
	if isInnoDBEngine(tbl.Engine) && unquoteTableOption(tbl.Tablespace) == oldName {
		tbl.Tablespace = newName
	}
	if tbl.Partitioning != nil {
		for _, pd := range tbl.Partitioning.Partitions {
			if pd == nil {
				continue
			}
			if isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, "")) && unquoteTableOption(pd.Tablespace) == oldName {
				pd.Tablespace = newName
			}
			for _, spd := range pd.SubPartitions {
				if spd != nil && isInnoDBEngine(effectivePartEngine(tbl.Engine, pd.Engine, spd.Engine)) && unquoteTableOption(spd.Tablespace) == oldName {
					spd.Tablespace = newName
				}
			}
		}
	}
}

func validateTablespaceCreateName(name string) error {
	if strings.HasPrefix(strings.ToLower(name), "innodb_") {
		return errReservedInnoDBTablespacePrefix()
	}
	if name == "mysql" {
		return errTablespaceExists(name)
	}
	return nil
}

func (c *Catalog) validateTablespaceCreateEngine(engine, op string) error {
	if engine == "" {
		return nil
	}
	if err := validateStorageEngine(engine, c.noEngineSubstitution()); err != nil {
		return err
	}
	canon := storedEngineName(engine)
	if !isInnoDBEngine(canon) {
		if strings.EqualFold(canon, "MRG_MyISAM") {
			canon = "MRG_MYISAM"
		}
		return errTablespaceUnsupportedEngine(canon, op)
	}
	return nil
}

func (c *Catalog) isGeneralTablespaceName(name string) bool {
	if c.tablespaces[name] {
		return true
	}
	lower := strings.ToLower(name)
	return lower == "innodb_system" || lower == "innodb_temporary" || lower == "mysql"
}

func (c *Catalog) createTablespace(stmt *nodes.CreateTablespaceStmt) error {
	op := "CREATE TABLESPACE"
	if stmt.Undo {
		op = "CREATE UNDO TABLESPACE"
	}
	if err := c.validateTablespaceCreateEngine(stmt.Engine, op); err != nil {
		return err
	}
	if stmt.AutoextendSize != "" {
		if err := validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
	}
	if stmt.Encryption != "" && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
		return errInvalidEncryptionOption()
	}
	if err := validateTablespaceCreateName(stmt.Name); err != nil {
		return err
	}
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo || c.tablespaces[stmt.Name] {
		return errTablespaceExists(stmt.Name)
	}
	if stmt.Undo {
		if stmt.DataFile != "" && !strings.HasSuffix(strings.ToLower(stmt.DataFile), ".ibu") {
			return errUndoTablespaceDatafileExt()
		}
		c.undoTablespaces[stmt.Name] = true
		return nil
	}
	c.tablespaces[stmt.Name] = true
	return nil
}

func (c *Catalog) alterTablespace(stmt *nodes.AlterTablespaceStmt) error {
	if stmt.Undo {
		if err := c.validateTablespaceCreateEngine(stmt.Engine, "ALTER UNDO TABLESPACE"); err != nil {
			return err
		}
		if stmt.AutoextendSize != "" {
			if err := validateAutoextendSize(stmt.AutoextendSize); err != nil {
				return err
			}
		}
		if stmt.Encryption != "" && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
			return errInvalidEncryptionOption()
		}
		if c.isGeneralTablespaceName(stmt.Name) {
			return errWrongTablespaceKind("ALTER UNDO TABLESPACE", stmt.Name, "a general", "ALTER TABLESPACE")
		}
		if _, exists := c.undoTablespaces[stmt.Name]; !exists {
			return errTablespaceDoesNotExist(stmt.Name)
		}
		if stmt.SetInactive {
			c.undoTablespaces[stmt.Name] = false
		} else if stmt.SetActive {
			c.undoTablespaces[stmt.Name] = true
		}
		return nil
	}
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo {
		return errWrongTablespaceKind("ALTER TABLESPACE", stmt.Name, "an undo", "ALTER UNDO TABLESPACE")
	}
	lower := strings.ToLower(stmt.Name)
	exists := c.tablespaces[stmt.Name] || lower == "innodb_system" || lower == "innodb_temporary"
	if !exists && !strings.HasPrefix(lower, "innodb_") && lower != "mysql" {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if exists && stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if stmt.AutoextendSize != "" {
		if err := validateAutoextendSize(stmt.AutoextendSize); err != nil {
			return err
		}
	}
	if stmt.Encryption != "" && !strings.EqualFold(stmt.Encryption, "Y") && !strings.EqualFold(stmt.Encryption, "N") {
		return errInvalidEncryptionOption()
	}
	if err := validateTablespaceCreateName(stmt.Name); err != nil {
		return err
	}
	if !c.tablespaces[stmt.Name] {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if stmt.RenameTo != "" {
		if err := validateTablespaceCreateName(stmt.RenameTo); err != nil {
			return err
		}
		if _, isUndo := c.undoTablespaces[stmt.RenameTo]; stmt.RenameTo != stmt.Name && (c.tablespaces[stmt.RenameTo] || isUndo) {
			return errTablespaceExists(stmt.RenameTo)
		}
		delete(c.tablespaces, stmt.Name)
		c.tablespaces[stmt.RenameTo] = true
		for _, db := range c.databases {
			for _, tbl := range db.Tables {
				renameTableTablespaceRefs(tbl, stmt.Name, stmt.RenameTo)
			}
		}
	}
	return nil
}

func (c *Catalog) dropTablespace(stmt *nodes.DropTablespaceStmt) error {
	if stmt.Undo {
		if err := c.validateTablespaceCreateEngine(stmt.Engine, "DROP UNDO TABLESPACE"); err != nil {
			return err
		}
		if c.isGeneralTablespaceName(stmt.Name) {
			return errWrongTablespaceKind("DROP UNDO TABLESPACE", stmt.Name, "a general", "DROP TABLESPACE")
		}
		if strings.HasPrefix(strings.ToLower(stmt.Name), "innodb_") {
			return errReservedInnoDBTablespacePrefix()
		}
		active, exists := c.undoTablespaces[stmt.Name]
		if !exists {
			return errTablespaceDoesNotExist(stmt.Name)
		}
		if active {
			return errDropUndoTablespaceFailed(stmt.Name)
		}
		delete(c.undoTablespaces, stmt.Name)
		return nil
	}
	if _, isUndo := c.undoTablespaces[stmt.Name]; isUndo && !strings.HasPrefix(strings.ToLower(stmt.Name), "innodb_") {
		return errWrongTablespaceKind("DROP TABLESPACE", stmt.Name, "an undo", "DROP UNDO TABLESPACE")
	}
	lower := strings.ToLower(stmt.Name)
	exists := c.tablespaces[stmt.Name] || lower == "innodb_system" || lower == "innodb_temporary"
	if !exists && !strings.HasPrefix(lower, "innodb_") {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	if exists && stmt.Engine != "" && !strings.EqualFold(stmt.Engine, "InnoDB") {
		return errTablespaceEngineMismatch(stmt.Engine, stmt.Name)
	}
	if strings.HasPrefix(lower, "innodb_") {
		return errReservedInnoDBTablespacePrefix()
	}
	if !c.tablespaces[stmt.Name] {
		return errTablespaceDoesNotExist(stmt.Name)
	}
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			if tableUsesTablespace(tbl, stmt.Name) {
				return errTablespaceIsNotEmpty(stmt.Name)
			}
		}
	}
	delete(c.tablespaces, stmt.Name)
	return nil
}

// alterDiscardImportTablespace handles ALTER TABLE ... DISCARD/IMPORT TABLESPACE
// as InnoDB enforces it on mysql:8.0.46.
func (c *Catalog) alterDiscardImportTablespace(tbl *Table, isImport bool) error {
	if !isInnoDBEngine(tbl.Engine) {
		return errTableStorageEngineOption(tbl.Name)
	}
	if tbl.Temporary {
		return errCannotDiscardImportTempTable()
	}
	if tbl.Tablespace != "" && !strings.EqualFold(unquoteTableOption(tbl.Tablespace), "innodb_file_per_table") {
		return errNotAllowedCommand()
	}
	if !isImport {
		tbl.TablespaceDiscarded = true
		return nil
	}
	dbName := ""
	if tbl.Database != nil {
		dbName = tbl.Database.Name
	}
	if !tbl.TablespaceDiscarded {
		if tbl.Partitioning != nil && len(tbl.Partitioning.Partitions) > 0 {
			return errTablespaceExists(fmt.Sprintf("%s/%s#p#%s", dbName, tbl.Name, tbl.Partitioning.Partitions[0].Name))
		}
		return errTablespaceExists(fmt.Sprintf("%s/%s", dbName, tbl.Name))
	}
	return errTablespaceMissing(dbName, tbl.Name)
}

// validatePartitionTablespaces checks TABLESPACE options on partition and
// subpartition definitions when their effective storage engine is InnoDB.
func (c *Catalog) validatePartitionTablespaces(tableName string, temporary bool, tableEngine string, defs []*nodes.PartitionDef) error {
	for _, pd := range defs {
		if pd == nil {
			continue
		}
		partEngine := tableEngine
		partTS := ""
		for _, opt := range pd.Options {
			if opt == nil {
				continue
			}
			if engineOptionName(opt.Name) {
				partEngine = storedEngineName(opt.Value)
			} else if strings.EqualFold(opt.Name, "tablespace") {
				partTS = opt.Value
			}
		}
		if partTS != "" && isInnoDBEngine(partEngine) {
			if err := c.validateInnoDBTablespace(tableName, temporary, partTS); err != nil {
				return err
			}
			if !c.loadingSDL {
				if err := validateInnoDBPartitionedSharedTablespace(partTS); err != nil {
					return err
				}
			}
		}
		for _, sub := range pd.SubPartitions {
			if sub == nil {
				continue
			}
			subEngine := partEngine
			subTS := ""
			for _, opt := range sub.Options {
				if opt == nil {
					continue
				}
				if engineOptionName(opt.Name) {
					subEngine = storedEngineName(opt.Value)
				} else if strings.EqualFold(opt.Name, "tablespace") {
					subTS = opt.Value
				}
			}
			if subTS != "" && isInnoDBEngine(subEngine) {
				if err := c.validateInnoDBTablespace(tableName, temporary, subTS); err != nil {
					return err
				}
				if !c.loadingSDL {
					if err := validateInnoDBPartitionedSharedTablespace(subTS); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

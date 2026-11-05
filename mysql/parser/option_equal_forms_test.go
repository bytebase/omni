package parser

import (
	"regexp"
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// := wherever MySQL 8.0 takes it in place of =: sql_yacc.yy's opt_equal and equal, which are = or SET_VAR. Everywhere
// the grammar writes a literal = it is no spelling of it, and the last groups hold
// those to it.
//
// The verdicts were measured against mysql:8.0 through PREPARE on 2026-11-05, and
// the recorded answer is what that server gave. TestOptionEqualFormsMatchEngine holds
// them to the server. Assignment lists are in insert_forms_test.go.

// optionEqualCase is a statement and whether the engine's parser accepts it.
type optionEqualCase struct {
	sql    string
	accept bool
}

var optionEqualCases = []optionEqualCase{
	// Table options: create_table_option.
	{`create table t (a int) engine := innodb`, true},
	{`create table t (a int) secondary_engine := null`, true},
	{`create table t (a int) secondary_engine := x`, true},
	{`create table t (a int) max_rows := 1`, true},
	{`create table t (a int) min_rows := 1`, true},
	{`create table t (a int) avg_row_length := 1`, true},
	{`create table t (a int) password := 'x'`, true},
	{`create table t (a int) comment := 'x'`, true},
	{`create table t (a int) compression := 'zlib'`, true},
	{`create table t (a int) encryption := 'y'`, true},
	{`create table t (a int) auto_increment := 1`, true},
	{`create table t (a int) pack_keys := 1`, true},
	{`create table t (a int) pack_keys := default`, true},
	{`create table t (a int) stats_auto_recalc := 1`, true},
	{`create table t (a int) stats_auto_recalc := default`, true},
	{`create table t (a int) stats_persistent := 1`, true},
	{`create table t (a int) stats_persistent := default`, true},
	{`create table t (a int) stats_sample_pages := 1`, true},
	{`create table t (a int) stats_sample_pages := default`, true},
	{`create table t (a int) checksum := 1`, true},
	{`create table t (a int) table_checksum := 1`, true},
	{`create table t (a int) page_checksum := 1`, false},
	{`create table t (a int) delay_key_write := 1`, true},
	{`create table t (a int) row_format := dynamic`, true},
	{`create table t (a int) union := (u)`, true},
	{`create table t (a int) insert_method := no`, true},
	{`create table t (a int) data directory := 'x'`, true},
	{`create table t (a int) index directory := 'x'`, true},
	{`create table t (a int) tablespace := ts`, true},
	{`create table t (a int) connection := 'x'`, true},
	{`create table t (a int) key_block_size := 1`, true},
	{`create table t (a int) engine_attribute := '{}'`, true},
	{`create table t (a int) secondary_engine_attribute := '{}'`, true},
	{`create table t (a int) transactional := 1`, false},
	{`create table t (a int) sequence := 1`, false},
	{`create table t (a int) character set := utf8mb4`, true},
	{`create table t (a int) default character set := utf8mb4`, true},
	{`create table t (a int) charset := utf8mb4`, true},
	{`create table t (a int) default charset := utf8mb4`, true},
	{`create table t (a int) collate := utf8mb4_bin`, true},
	{`create table t (a int) default collate := utf8mb4_bin`, true},
	{`create table t (a int) engine := innodb, comment := 'x', max_rows = 1`, true},
	{`create table t (a int) engine := innodb comment 'x'`, true},
	{`alter table t engine := innodb`, true},
	{`alter table t comment := 'x'`, true},
	{`alter table t auto_increment := 1`, true},
	{`alter table t row_format := dynamic`, true},
	{`alter table t character set := utf8mb4`, true},
	{`alter table t default charset := utf8mb4`, true},
	{`alter table t collate := utf8mb4_bin`, true},
	{`alter table t key_block_size := 1`, true},
	{`alter table t union := (u)`, true},
	{`create table t (a int) foo := 'x'`, false},
	{`create table t (a int) foo := bar`, false},
	{`create table t (a int) foo := 1`, false},
	{`create table t (a int) foo := default`, false},

	// Algorithm and lock in ALTER, CREATE INDEX and DROP INDEX.
	{`alter table t algorithm := inplace, add column b int`, true},
	{`alter table t algorithm := default, add column b int`, true},
	{`alter table t lock := none, add column b int`, true},
	{`alter table t lock := default, add column b int`, true},
	{`alter table t add column b int, algorithm := inplace, lock := none`, true},
	{`create index i on t (a) algorithm := inplace`, true},
	{`create index i on t (a) lock := none`, true},
	{`create index i on t (a) algorithm := inplace lock := none`, true},
	{`drop index i on t algorithm := inplace`, true},
	{`drop index i on t lock := none`, true},
	{`drop index i on t algorithm := inplace lock := none`, true},

	// Index options: common_index_attribute and its neighbours.
	{`create index i on t (a) key_block_size := 1`, true},
	{`create index i on t (a) engine_attribute := '{}'`, true},
	{`create index i on t (a) secondary_engine_attribute := '{}'`, true},
	{`create index i on t (a) foo := 1`, false},
	{`create table t (a int, key k (a) key_block_size := 1)`, true},
	{`create table t (a int, key k (a) engine_attribute := '{}')`, true},
	{`create table t (a int, unique key k (a) key_block_size := 1)`, true},
	{`create table t (a int, primary key (a) key_block_size := 1)`, true},
	{`alter table t add index i (a) key_block_size := 1`, true},
	{`alter table t add index i (a) engine_attribute := '{}'`, true},
	{`create fulltext index i on t (a) key_block_size := 1`, true},
	{`create spatial index i on t (a) key_block_size := 1`, true},

	// Column attributes.
	{`create table t (a int engine_attribute := '{}')`, true},
	{`create table t (a int secondary_engine_attribute := '{}')`, true},
	{`create table t (a int foo := 1)`, false},
	{`create table t (a blob compressed := zlib)`, false},

	// Database options.
	{`create database d character set := utf8mb4`, true},
	{`create database d default character set := utf8mb4`, true},
	{`create database d charset := utf8mb4`, true},
	{`create database d collate := utf8mb4_bin`, true},
	{`create database d default collate := utf8mb4_bin`, true},
	{`create database d encryption := 'n'`, true},
	{`create database d default encryption := 'n'`, true},
	{`create database d comment := 'x'`, false},
	{`create schema d character set := utf8mb4`, true},
	{`alter database d character set := utf8mb4`, true},
	{`alter database d collate := utf8mb4_bin`, true},
	{`alter database d encryption := 'n'`, true},
	{`alter database d read only := 1`, true},
	{`alter database d read only := default`, true},
	{`alter database d comment := 'x'`, false},
	{`alter schema d default charset := utf8mb4`, true},

	// Partition options.
	{`create table t (a int) partition by hash(a) (partition p0 engine := innodb)`, true},
	{`create table t (a int) partition by hash(a) (partition p0 storage engine := innodb)`, true},
	{`create table t (a int) partition by hash(a) (partition p0 tablespace := ts)`, true},
	{`create table t (a int) partition by hash(a) (partition p0 nodegroup := 1)`, true},
	{`create table t (a int) partition by hash(a) (partition p0 max_rows := 1)`, true},
	{`create table t (a int) partition by hash(a) (partition p0 min_rows := 1)`, true},
	{`create table t (a int) partition by hash(a) (partition p0 data directory := 'x')`, true},
	{`create table t (a int) partition by hash(a) (partition p0 index directory := 'x')`, true},
	{`create table t (a int) partition by hash(a) (partition p0 comment := 'x')`, true},
	{`create table t (a int) partition by hash(a) (partition p0 connection := 'x')`, false},
	{`create table t (a int) partition by range (a) subpartition by hash (a) (partition p0 values less than (1) (subpartition s0 engine := innodb))`, true},
	{`alter table t add partition (partition p1 engine := innodb)`, true},
	{`alter table t reorganize partition p0 into (partition p1 values less than (1) comment := 'x')`, true},

	// Tablespace and logfile group options.
	{`create tablespace ts add datafile 'a.ibd' initial_size := 1024`, true},
	{`create tablespace ts add datafile 'a.ibd' autoextend_size := 4M`, true},
	{`create tablespace ts add datafile 'a.ibd' max_size := 1024`, true},
	{`create tablespace ts add datafile 'a.ibd' extent_size := 1024`, true},
	{`create tablespace ts add datafile 'a.ibd' nodegroup := 1`, true},
	{`create tablespace ts add datafile 'a.ibd' comment := 'x'`, true},
	{`create tablespace ts add datafile 'a.ibd' engine := innodb`, true},
	{`create tablespace ts add datafile 'a.ibd' file_block_size := 1024`, true},
	{`create tablespace ts add datafile 'a.ibd' encryption := 'y'`, true},
	{`create tablespace ts add datafile 'a.ibd' engine_attribute := '{}'`, true},
	{`alter tablespace ts autoextend_size := 4M`, true},
	{`alter tablespace ts engine := innodb`, true},
	{`alter tablespace ts encryption := 'y'`, true},
	{`alter tablespace ts initial_size := 1024`, true},
	{`alter tablespace ts engine_attribute := '{}'`, true},
	{`alter tablespace ts add datafile 'b' initial_size := 1024`, true},
	{`drop tablespace ts engine := innodb`, true},
	{`create logfile group lg add undofile 'u' initial_size := 1024`, true},
	{`create logfile group lg add undofile 'u' undo_buffer_size := 1024`, true},
	{`create logfile group lg add undofile 'u' redo_buffer_size := 1024`, true},
	{`create logfile group lg add undofile 'u' nodegroup := 1`, true},
	{`create logfile group lg add undofile 'u' comment := 'x'`, true},
	{`create logfile group lg add undofile 'u' engine := ndb`, true},
	{`alter logfile group lg add undofile 'u' initial_size := 1024`, true},
	{`alter logfile group lg add undofile 'u' engine := ndb`, true},
	{`drop logfile group lg engine := ndb`, true},

	// Resource groups and CLONE.
	{`create resource group g type := user`, true},
	{`create resource group g type := user vcpu := 1 thread_priority := 1`, true},
	{`alter resource group g vcpu := 1`, true},
	{`alter resource group g thread_priority := 1`, true},
	{`clone local data directory := 'x'`, true},

	// MariaDB sequences.
	{`create sequence s start := 1`, false},
	{`create sequence s increment := 1`, false},
	{`create sequence s minvalue := 1`, false},
	{`create sequence s maxvalue := 1`, false},
	{`create sequence s cache := 1`, false},
	{`alter sequence s restart := 1`, false},
	{`alter sequence s minvalue := 1`, false},

	// SET statements.
	{`set autocommit := 1`, true},
	{`set @@autocommit := 1`, true},
	{`set session autocommit := 1`, true},
	{`set global autocommit := 1`, true},
	{`set @@global.autocommit := 1`, true},
	{`set persist autocommit := 1`, true},
	{`set @a := 1`, true},
	{`set @a := 1, autocommit := 1`, true},
	{`set names := 'utf8mb4'`, false},
	{`set password := 'x'`, true},
	{`set password for u := 'x'`, true},
	{`set password for u@h := 'x'`, true},
	{`set statement max_statement_time := 1 for select 1`, false},
	{`create definer := u view v as select 1`, false},

	// Literal equals: := is no spelling of these.
	{`create algorithm := merge view v as select 1`, false},
	{`alter algorithm := merge view v as select 1`, false},
	{`create definer := u procedure p() select 1`, false},
	{`alter definer := u view v as select 1`, false},
	{`explain format := json select 1`, false},
	{`alter table t rename := u`, false},
	{`create table t (a int) partition by key algorithm := 2 (a)`, false},
	{`signal sqlstate '45000' set message_text := 'x'`, false},
	{`resignal set message_text := 'x'`, false},
	{`get diagnostics @a := number`, false},
	{`get diagnostics condition 1 @a := mysql_errno`, false},
	{`change replication source to source_host := 'x'`, false},
	{`change master to master_host := 'x'`, false},
	{`change master to master_port := 1`, false},
	{`change master to master_log_file := 'x'`, false},
	{`change master to ignore_server_ids := (1)`, false},
	{`change replication filter replicate_do_db := (d)`, false},
	{`start replica until source_log_file := 'x', source_log_pos := 1`, false},
	{`start slave until master_log_file := 'x', master_log_pos := 1`, false},
	{`start replica user := 'x'`, false},
	{`start replica password := 'x'`, false},
	{`start replica default_auth := 'x'`, false},
	{`start replica plugin_dir := 'x'`, false},
	{`start replica until sql_before_gtids := 'x'`, false},
	{`create table t (a int) engine : = innodb`, false},
	{`create table t (a int) engine :=`, false},
	{`create table t (a int) engine :=  = innodb`, false},
	{`create table t (a int) engine = := innodb`, false},

	// An option's value, and a required equals sign.
	{`create table t (a int) charset = binary`, true},
	{`create table t (a int) charset binary`, true},
	{`create table t (a int) charset := binary`, true},
	{`create table t (a int) default charset = binary`, true},
	{`create table t (a int) character set = binary`, true},
	{`create table t (a int) character set binary`, true},
	{`create table t (a int) default character set binary`, true},
	{`create table t (a int) collate = binary`, true},
	{`create table t (a int) collate binary`, true},
	{`create table t (a int) collate := binary`, true},
	{`create table t (a int) default collate = binary`, true},
	{`create table t (a int) charset = default`, false},
	{`create table t (a int) character set = default`, false},
	{`create table t (a int) collate = default`, false},
	{`create table t (a int) charset = 'utf8mb4'`, true},
	{`create table t (a int) collate = 'utf8mb4_bin'`, true},
	{`create table t (a int) charset = select`, false},
	{`create table t (a int) collate = select`, false},
	{`create table t (a int) charset = null`, false},
	{`alter table t charset = binary`, true},
	{`alter table t character set = binary`, true},
	{`alter table t collate = binary`, true},
	{`alter table t default charset binary`, true},
	{`alter table t convert to character set binary`, true},
	{`alter table t convert to charset binary`, true},
	{`alter table t convert to character set binary collate binary`, true},
	{`alter table t convert to character set utf8mb4 collate binary`, true},
	{`alter table t convert to character set default`, true},
	{`alter table t convert to character set default collate utf8mb4_bin`, true},
	{`alter table t convert to character set select`, false},
	{`alter table t convert to character set utf8mb4 collate select`, false},
	{`create database d charset binary`, true},
	{`create database d collate binary`, true},
	{`create table t (a int) engine = binary`, false},
	{`create table t (a int) engine = select`, false},
	{`create table t (a int) row_format = default`, true},
	{`create table t (a int) row_format = binary`, false},
	{`create table t (a int) row_format = select`, false},
	{`create table t (a int) compression = binary`, false},
	{`create table t (a int) insert_method = first`, true},
	{`create table t (a int) insert_method = last`, true},
	{`create table t (a int) insert_method = no`, true},
	{`create table t (a int) tablespace = binary`, false},
	{`create table t (a int) tablespace binary`, false},
	{`create table t (a int) storage binary`, false},
	{`create table t (a int) secondary_engine = binary`, false},
	{`create table t (a int) pack_keys = binary`, false},
	{`create table t (a int) max_rows = binary`, false},
	{`create table t (a int) foo = binary`, false},
	{`create table t (a int) foo = default`, false},
	{`create table t (a int) partition by hash(a) (partition p engine = binary)`, false},
	{`create table t (a int) partition by hash(a) (partition p tablespace = binary)`, false},
	{`alter table t algorithm = default, add column b int`, true},
	{`alter table t lock = default, add column b int`, true},
	{`alter table t algorithm = binary, add column b int`, false},
	{`alter table t lock = binary, add column b int`, false},
	{`alter table t algorithm = select, add column b int`, false},
	{`set password 'x'`, false},
	{`set password for u 'x'`, false},
	{`create table t (a int) foo := null`, false},
	{`create table t (a int) foo := select`, false},
	{`create table t (a int) pack_keys := null`, false},
	{`create table t (a int) max_rows := null`, false},
	{`create table t (a int) comment := null`, false},
	{`create table t (a int) secondary_engine := 'x'`, true},
	{`create table t (a int) auto_increment := 1 engine := innodb`, true},
	{`create tablespace ts add datafile 'a.ibd' engine := ndb nodegroup := 1 comment := 'x'`, true},
	{`clone instance from 'u'@'h':1 identified by 'x' data directory := 'x'`, true},
	{`create index i on t (a) key_block_size := 1 comment 'x'`, true},
	{`drop index i on t algorithm := default`, true},
	{`create table t (a int) storage disk`, true},
	{`create table t (a int) storage memory`, true},
	{`create table t (a int) storage := disk`, false},
	{`create table t (a int) storage = disk`, false},
	{`create table t (a int) storage := memory`, false},
	{`create table t (a int) storage engine := innodb`, false},
}

// TestOptionEqualForms holds the parser to the verdicts above. An accepted statement
// must read as one: Parse takes statements with no separator between them, so a tail
// the grammar should refuse could otherwise pass as a second one.
func TestOptionEqualForms(t *testing.T) {
	for _, c := range optionEqualCases {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if got := err == nil; got != c.accept {
				t.Fatalf("Parse(%q) accepts = %v, the engine answered %v (err: %v)", c.sql, got, c.accept, err)
			}
			if err == nil && len(list.Items) != 1 {
				t.Errorf("Parse(%q) read %d statements, want 1", c.sql, len(list.Items))
			}
		})
	}
}

// optionLocFree drops the offsets, which := and = cannot share, from a printed tree.
var optionLocFree = regexp.MustCompile(` :loc -?[0-9]+`)

// TestOptionColonEqualsReadsAsEquals: := is another spelling of =, and the tree does
// not show which was written. Every accepted row with a := must read the same with =.
func TestOptionColonEqualsReadsAsEquals(t *testing.T) {
	compared := 0
	for _, c := range optionEqualCases {
		if !c.accept || !strings.Contains(c.sql, ":=") {
			continue
		}
		equal := strings.ReplaceAll(c.sql, ":=", "=")
		t.Run(c.sql, func(t *testing.T) {
			colonList, err := Parse(c.sql)
			if err != nil {
				t.Fatal(err)
			}
			equalList, err := Parse(equal)
			if err != nil {
				t.Fatalf("Parse(%q): %v", equal, err)
			}
			got := optionLocFree.ReplaceAllString(ast.NodeToString(colonList.Items[0]), "")
			want := optionLocFree.ReplaceAllString(ast.NodeToString(equalList.Items[0]), "")
			if got != want {
				t.Errorf("trees differ:\n:= %s\n=  %s", got, want)
			}
		})
		compared++
	}
	if compared == 0 {
		t.Fatal("no row to compare")
	}
}

// TestOptionColonEqualsKeepsItsValue: the printed trees above agree only if they hold
// the option at all, so a value is checked where a reader would look for it.
func TestOptionColonEqualsKeepsItsValue(t *testing.T) {
	list, err := Parse("CREATE TABLE t (a INT) ENGINE := InnoDB, COMMENT := 'x'")
	if err != nil {
		t.Fatal(err)
	}
	opts := list.Items[0].(*ast.CreateTableStmt).Options
	if len(opts) != 2 || opts[0].Value != "InnoDB" || opts[1].Value != "x" {
		t.Errorf("Options = %+v", opts)
	}
}

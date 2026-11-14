package catalog

import (
	"strings"
	"testing"
)

const fkDepartments = "CREATE TABLE departments (dept_no char(4) NOT NULL, PRIMARY KEY (dept_no));\n"

func TestMigration_ForeignKeyBackingIndexName(t *testing.T) {
	cases := []struct {
		name     string
		from, to string
		want     []string
	}{
		{
			name: "new table with a backing key named after the column",
			from: fkDepartments,
			to: fkDepartments + "CREATE TABLE dept_emp (emp_no int NOT NULL, dept_no char(4) NOT NULL, PRIMARY KEY (emp_no, dept_no), " +
				"KEY dept_no (dept_no), CONSTRAINT dept_emp_ibfk_2 FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE);",
			want: []string{
				"ALTER TABLE `app`.`dept_emp` ADD KEY `dept_no` (`dept_no`)",
				"ALTER TABLE `app`.`dept_emp` ADD CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE",
			},
		},
		{
			name: "new table with an unnamed foreign key",
			from: fkDepartments,
			to:   fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), FOREIGN KEY (dept_no) REFERENCES departments (dept_no));",
			want: []string{
				"ALTER TABLE `app`.`c` ADD KEY `dept_no` (`dept_no`)",
				"ALTER TABLE `app`.`c` ADD CONSTRAINT `c_ibfk_1` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)",
			},
		},
		{
			name: "foreign key added to an existing table",
			from: fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4));",
			to: fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), KEY dept_no (dept_no), " +
				"CONSTRAINT c_fk FOREIGN KEY (dept_no) REFERENCES departments (dept_no));",
			want: []string{
				"ALTER TABLE `app`.`c` ADD KEY `dept_no` (`dept_no`)",
				"ALTER TABLE `app`.`c` ADD CONSTRAINT `c_fk` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)",
			},
		},
		{
			name: "backing index named after the constraint is left to MySQL",
			from: fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4));",
			to:   fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), CONSTRAINT c_fk FOREIGN KEY (dept_no) REFERENCES departments (dept_no));",
			want: []string{
				"ALTER TABLE `app`.`c` ADD CONSTRAINT `c_fk` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)",
			},
		},
		{
			name: "backing index that already exists is reused",
			from: fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), KEY dept_no (dept_no));",
			to: fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), KEY dept_no (dept_no), " +
				"CONSTRAINT c_fk FOREIGN KEY (dept_no) REFERENCES departments (dept_no));",
			want: []string{
				"ALTER TABLE `app`.`c` ADD CONSTRAINT `c_fk` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := planFor(t, dbDDL+tc.from, dbDDL+tc.to, NormalizerFor(MySQL80))
			var got []string
			for _, op := range plan.Ops {
				if op.Type != OpCreateTable {
					got = append(got, op.SQL)
				}
			}
			if strings.Join(got, ";\n") != strings.Join(tc.want, ";\n") {
				t.Errorf("plan:\n%s\nwant:\n%s", strings.Join(got, ";\n"), strings.Join(tc.want, ";\n"))
			}
		})
	}
}

// A form without the KEY and its readback with the KEY still diff equal.
func TestMigration_ForeignKeyBackingIndexNameNoPhantom(t *testing.T) {
	user := dbDDL + fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), CONSTRAINT c_fk FOREIGN KEY (dept_no) REFERENCES departments (dept_no));"
	readback := dbDDL + fkDepartments + "CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), KEY dept_no (dept_no), " +
		"CONSTRAINT c_fk FOREIGN KEY (dept_no) REFERENCES departments (dept_no));"
	n := NormalizerFor(MySQL80)
	for _, pair := range [][2]string{{user, readback}, {readback, user}} {
		if plan := planFor(t, pair[0], pair[1], n); plan.SQL() != "" {
			t.Errorf("want an empty plan, got:\n%s", plan.SQL())
		}
	}
}

func TestMigration_ForeignKeyRestrictIsEmitted(t *testing.T) {
	to := dbDDL + "CREATE TABLE city (city_id int NOT NULL, PRIMARY KEY (city_id));\n" +
		"CREATE TABLE address (address_id int NOT NULL, city_id int NOT NULL, PRIMARY KEY (address_id), KEY idx_fk_city_id (city_id), " +
		"CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE RESTRICT ON UPDATE CASCADE, " +
		"CONSTRAINT fk_no_action FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE NO ACTION ON UPDATE RESTRICT);"
	plan := planFor(t, dbDDL, to, NormalizerFor(MySQL80))
	var fks []string
	for _, op := range plan.Ops {
		if strings.Contains(op.SQL, "FOREIGN KEY") {
			fks = append(fks, op.SQL)
		}
	}
	want := []string{
		"ALTER TABLE `app`.`address` ADD CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON DELETE RESTRICT ON UPDATE CASCADE",
		"ALTER TABLE `app`.`address` ADD CONSTRAINT `fk_no_action` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON UPDATE RESTRICT",
	}
	if strings.Join(fks, "\n") != strings.Join(want, "\n") {
		t.Errorf("foreign keys:\n%s\nwant:\n%s", strings.Join(fks, "\n"), strings.Join(want, "\n"))
	}
}

package catalog

import (
	"strings"
	"testing"
)

func TestMigration_ForeignKeyNoActionIsEmitted(t *testing.T) {
	to := dbDDL + "CREATE TABLE city (city_id int NOT NULL, PRIMARY KEY (city_id));\n" +
		"CREATE TABLE address (address_id int NOT NULL, city_id int NOT NULL, PRIMARY KEY (address_id), KEY idx_fk_city_id (city_id), " +
		"CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE NO ACTION ON UPDATE CASCADE, " +
		"CONSTRAINT fk_restrict FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE RESTRICT ON UPDATE NO ACTION);"
	plan := planFor(t, dbDDL, to, NormalizerFor(MySQL80))
	var fks []string
	for _, op := range plan.Ops {
		if strings.Contains(op.SQL, "FOREIGN KEY") {
			fks = append(fks, op.SQL)
		}
	}
	want := []string{
		"ALTER TABLE `app`.`address` ADD CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON DELETE NO ACTION ON UPDATE CASCADE",
		"ALTER TABLE `app`.`address` ADD CONSTRAINT `fk_restrict` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON UPDATE NO ACTION",
	}
	if strings.Join(fks, "\n") != strings.Join(want, "\n") {
		t.Errorf("foreign keys:\n%s\nwant:\n%s", strings.Join(fks, "\n"), strings.Join(want, "\n"))
	}
}

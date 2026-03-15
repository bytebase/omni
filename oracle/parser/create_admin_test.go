package parser

import (
	"testing"
)

// TestParseCreateAdminStmts tests parsing of administrative DDL statements
// handled by create_admin.go: CREATE/ALTER/DROP for TABLESPACE, DIRECTORY,
// CLUSTER, DIMENSION, JAVA, LIBRARY, CONTEXT, DOMAIN, INDEXTYPE, OPERATOR,
// MATERIALIZED ZONEMAP, INMEMORY JOIN GROUP, PROPERTY GRAPH, VECTOR INDEX,
// MLE ENV/MODULE, FLASHBACK ARCHIVE, ROLLBACK SEGMENT, EDITION, RESTORE POINT,
// SCHEMA, LOCKDOWN PROFILE, OUTLINE, and more.
func TestParseCreateAdminStmts(t *testing.T) {
	tests := []string{
		// ---- CREATE TABLESPACE ----
		// Permanent tablespace (minimal)
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M`,
		// IF NOT EXISTS
		`CREATE TABLESPACE IF NOT EXISTS users
			DATAFILE '/u01/data/users01.dbf' SIZE 100M`,
		// BIGFILE tablespace
		`CREATE BIGFILE TABLESPACE bigts
			DATAFILE '/u01/data/bigts01.dbf' SIZE 10G`,
		// SMALLFILE tablespace
		`CREATE SMALLFILE TABLESPACE smallts
			DATAFILE '/u01/data/smallts01.dbf' SIZE 500M`,
		// With AUTOEXTEND
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			AUTOEXTEND ON NEXT 100M MAXSIZE 5G`,
		// AUTOEXTEND OFF
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			AUTOEXTEND OFF`,
		// AUTOEXTEND UNLIMITED
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			AUTOEXTEND ON MAXSIZE UNLIMITED`,
		// LOGGING/NOLOGGING
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			LOGGING`,
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			NOLOGGING`,
		// ONLINE/OFFLINE
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			ONLINE`,
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			OFFLINE`,
		// EXTENT MANAGEMENT LOCAL
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			EXTENT MANAGEMENT LOCAL AUTOALLOCATE`,
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			EXTENT MANAGEMENT LOCAL UNIFORM SIZE 1M`,
		// SEGMENT SPACE MANAGEMENT
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			SEGMENT SPACE MANAGEMENT AUTO`,
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			SEGMENT SPACE MANAGEMENT MANUAL`,
		// BLOCKSIZE
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			BLOCKSIZE 8K`,
		// FORCE LOGGING
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			FORCE LOGGING`,
		// FLASHBACK ON/OFF
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			FLASHBACK ON`,
		// ENCRYPTION
		`CREATE TABLESPACE secure_ts
			DATAFILE '/u01/data/secure01.dbf' SIZE 500M
			ENCRYPTION ENCRYPT`,
		`CREATE TABLESPACE secure_ts
			DATAFILE '/u01/data/secure01.dbf' SIZE 500M
			ENCRYPTION USING 'AES256' ENCRYPT`,
		// DEFAULT compression
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			DEFAULT TABLE COMPRESS FOR OLTP`,
		`CREATE TABLESPACE users
			DATAFILE '/u01/data/users01.dbf' SIZE 500M
			DEFAULT TABLE NOCOMPRESS`,
		// TEMPORARY TABLESPACE
		`CREATE TEMPORARY TABLESPACE temp
			TEMPFILE '/u01/data/temp01.dbf' SIZE 500M`,
		// LOCAL TEMPORARY TABLESPACE
		`CREATE LOCAL TEMPORARY TABLESPACE ltemp
			TEMPFILE '/u01/data/ltemp01.dbf' SIZE 200M`,
		// UNDO TABLESPACE
		`CREATE UNDO TABLESPACE undots
			DATAFILE '/u01/data/undo01.dbf' SIZE 500M`,
		// UNDO with retention
		`CREATE UNDO TABLESPACE undots
			DATAFILE '/u01/data/undo01.dbf' SIZE 500M
			RETENTION GUARANTEE`,
		// LOST WRITE PROTECTION
		`CREATE TABLESPACE lwp_ts
			DATAFILE '/u01/data/lwp01.dbf' SIZE 500M
			ENABLE LOST WRITE PROTECTION`,
		// Comprehensive
		`CREATE TABLESPACE prod_data
			DATAFILE '/u01/data/prod01.dbf' SIZE 1G AUTOEXTEND ON NEXT 500M MAXSIZE 10G,
			         '/u01/data/prod02.dbf' SIZE 1G
			BLOCKSIZE 8K
			LOGGING
			FORCE LOGGING
			EXTENT MANAGEMENT LOCAL AUTOALLOCATE
			SEGMENT SPACE MANAGEMENT AUTO
			FLASHBACK ON
			ONLINE`,

		// ---- ALTER TABLESPACE ----
		`ALTER TABLESPACE users ONLINE`,
		`ALTER TABLESPACE users OFFLINE`,
		`ALTER TABLESPACE users READ ONLY`,
		`ALTER TABLESPACE users READ WRITE`,
		`ALTER TABLESPACE users RENAME TO new_users`,
		`ALTER TABLESPACE users ADD DATAFILE '/u01/data/users02.dbf' SIZE 500M`,
		`ALTER TABLESPACE users AUTOEXTEND ON NEXT 100M`,
		`ALTER TABLESPACE users COALESCE`,
		`ALTER TABLESPACE users LOGGING`,
		`ALTER TABLESPACE users NOLOGGING`,
		`ALTER TABLESPACE users FORCE LOGGING`,
		`ALTER TABLESPACE users FLASHBACK ON`,
		`ALTER TABLESPACE users FLASHBACK OFF`,
		`ALTER TABLESPACE temp ADD TEMPFILE '/u01/data/temp02.dbf' SIZE 500M`,
		`ALTER TABLESPACE users BEGIN BACKUP`,
		`ALTER TABLESPACE users END BACKUP`,
		`ALTER TABLESPACE users RETENTION GUARANTEE`,
		`ALTER TABLESPACE users RETENTION NOGUARANTEE`,
		`ALTER TABLESPACE users ENABLE LOST WRITE PROTECTION`,
		`ALTER TABLESPACE users DISABLE LOST WRITE PROTECTION`,

		// ---- DROP TABLESPACE ----
		`DROP TABLESPACE users`,
		`DROP TABLESPACE users INCLUDING CONTENTS`,
		`DROP TABLESPACE users INCLUDING CONTENTS AND DATAFILES`,
		`DROP TABLESPACE users INCLUDING CONTENTS CASCADE CONSTRAINTS`,

		// ---- CREATE DIRECTORY ----
		`CREATE DIRECTORY data_dir AS '/u01/data'`,
		`CREATE OR REPLACE DIRECTORY data_dir AS '/u01/data'`,
		// ---- DROP DIRECTORY ----
		`DROP DIRECTORY data_dir`,

		// ---- CREATE CONTEXT ----
		`CREATE CONTEXT app_ctx USING myschema.ctx_pkg`,
		`CREATE OR REPLACE CONTEXT app_ctx USING ctx_pkg`,
		`CREATE CONTEXT app_ctx USING ctx_pkg INITIALIZED EXTERNALLY`,
		`CREATE CONTEXT app_ctx USING ctx_pkg INITIALIZED GLOBALLY`,
		`CREATE CONTEXT app_ctx USING ctx_pkg ACCESSED GLOBALLY`,

		// ---- CREATE CLUSTER ----
		`CREATE CLUSTER hr.emp_dept_cluster (dept_id NUMBER(4))
			SIZE 512
			TABLESPACE users`,
		`CREATE CLUSTER IF NOT EXISTS hr.dept_cluster (dept_id NUMBER(4))`,
		// Hash cluster
		`CREATE CLUSTER hash_cluster (id NUMBER)
			HASHKEYS 100
			SIZE 256`,
		// Hash with HASH IS expression
		`CREATE CLUSTER hash_cluster (id NUMBER)
			HASHKEYS 100 HASH IS id
			SIZE 256`,
		// INDEX cluster
		`CREATE CLUSTER idx_cluster (dept_id NUMBER(4))
			INDEX`,
		// SINGLE TABLE
		`CREATE CLUSTER st_cluster (id NUMBER)
			HASHKEYS 50
			SINGLE TABLE`,
		// Physical attributes
		`CREATE CLUSTER phys_cluster (id NUMBER)
			PCTFREE 20
			PCTUSED 60
			INITRANS 2
			SIZE 1K
			TABLESPACE users`,
		// PARALLEL
		`CREATE CLUSTER par_cluster (id NUMBER)
			PARALLEL 4`,
		`CREATE CLUSTER par_cluster (id NUMBER)
			NOPARALLEL`,
		// CACHE/NOCACHE
		`CREATE CLUSTER cache_cluster (id NUMBER)
			CACHE`,
		`CREATE CLUSTER cache_cluster (id NUMBER)
			NOCACHE`,
		// ROWDEPENDENCIES
		`CREATE CLUSTER rd_cluster (id NUMBER)
			ROWDEPENDENCIES`,
		`CREATE CLUSTER rd_cluster (id NUMBER)
			NOROWDEPENDENCIES`,
		// Multiple columns
		`CREATE CLUSTER multi_cluster (dept_id NUMBER(4), loc_id NUMBER(4))
			SIZE 512`,
		// SORT
		`CREATE CLUSTER sort_cluster (id NUMBER SORT)
			HASHKEYS 100`,

		// ---- ALTER CLUSTER ----
		`ALTER CLUSTER hr.emp_dept_cluster SIZE 1024`,
		`ALTER CLUSTER hr.emp_dept_cluster PCTFREE 30`,
		`ALTER CLUSTER hr.emp_dept_cluster CACHE`,
		`ALTER CLUSTER hr.emp_dept_cluster NOCACHE`,
		`ALTER CLUSTER hr.emp_dept_cluster PARALLEL 4`,
		`ALTER CLUSTER hr.emp_dept_cluster ALLOCATE EXTENT`,
		`ALTER CLUSTER hr.emp_dept_cluster DEALLOCATE UNUSED`,

		// ---- DROP CLUSTER ----
		`DROP CLUSTER hr.emp_dept_cluster`,
		`DROP CLUSTER hr.emp_dept_cluster INCLUDING TABLES`,
		`DROP CLUSTER hr.emp_dept_cluster INCLUDING TABLES CASCADE CONSTRAINTS`,

		// ---- CREATE DIMENSION ----
		`CREATE DIMENSION time_dim
			LEVEL day IS (time_table.day_id)
			LEVEL month IS (time_table.month_id)
			LEVEL year IS (time_table.year_id)
			HIERARCHY time_hier (
				day CHILD OF month
				CHILD OF year
			)`,
		// With ATTRIBUTE
		`CREATE DIMENSION time_dim
			LEVEL day IS (time_table.day_id)
			LEVEL month IS (time_table.month_id)
			HIERARCHY time_hier (
				day CHILD OF month
			)
			ATTRIBUTE day DETERMINES (time_table.day_name)`,
		// With JOIN KEY
		`CREATE DIMENSION geo_dim
			LEVEL city IS (cities.city_id)
			LEVEL state IS (states.state_id)
			HIERARCHY geo_hier (
				city CHILD OF state
				JOIN KEY (cities.state_id) REFERENCES state
			)`,
		// SKIP WHEN NULL
		`CREATE DIMENSION time_dim
			LEVEL day IS (t.day_id) SKIP WHEN NULL
			LEVEL month IS (t.month_id)
			HIERARCHY h (day CHILD OF month)`,

		// ---- ALTER DIMENSION ----
		`ALTER DIMENSION time_dim
			ADD LEVEL quarter IS (time_table.quarter_id)`,
		`ALTER DIMENSION time_dim DROP LEVEL quarter`,
		`ALTER DIMENSION time_dim COMPILE`,

		// ---- CREATE JAVA ----
		`CREATE JAVA SOURCE NAMED myschema.MyClass AS 'public class MyClass { }'`,
		`CREATE OR REPLACE JAVA SOURCE NAMED MyClass AS 'public class MyClass { }'`,
		`CREATE AND RESOLVE JAVA SOURCE NAMED MyClass AS 'public class MyClass { }'`,
		`CREATE JAVA CLASS NAMED myschema.MyClass
			USING BFILE (java_dir, 'MyClass.class')`,
		`CREATE JAVA RESOURCE NAMED myschema.props
			USING BFILE (java_dir, 'app.properties')`,

		// ---- ALTER JAVA ----
		`ALTER JAVA SOURCE myschema.MyClass COMPILE`,
		`ALTER JAVA SOURCE myschema.MyClass RESOLVE`,
		`ALTER JAVA CLASS myschema.MyClass COMPILE`,

		// ---- DROP JAVA ----
		`DROP JAVA SOURCE myschema.MyClass`,
		`DROP JAVA CLASS myschema.MyClass`,
		`DROP JAVA RESOURCE myschema.props`,

		// ---- CREATE LIBRARY ----
		`CREATE LIBRARY mylib AS '/usr/lib/mylib.so'`,
		`CREATE OR REPLACE LIBRARY mylib AS '/usr/lib/mylib.so'`,
		`CREATE LIBRARY mylib IS '/usr/lib/mylib.so'
			AGENT agent_link`,
		`CREATE LIBRARY mylib AS '/usr/lib/mylib.so'
			CREDENTIAL my_cred`,
		// Note: EDITIONABLE/NONEDITIONABLE before LIBRARY dispatch is handled
		// at CREATE PROCEDURE/FUNCTION/PACKAGE/TYPE level, not LIBRARY directly.

		// ---- ALTER LIBRARY ----
		`ALTER LIBRARY mylib COMPILE`,
		`ALTER LIBRARY mylib EDITIONABLE`,
		`ALTER LIBRARY mylib NONEDITIONABLE`,

		// ---- CREATE SCHEMA ----
		`CREATE SCHEMA AUTHORIZATION hr
			CREATE TABLE emp (id NUMBER)
			CREATE VIEW emp_v AS SELECT id FROM emp
			GRANT SELECT ON emp TO public`,

		// ---- CREATE DOMAIN ----
		// Single column
		`CREATE DOMAIN email_dom AS VARCHAR2(100)
			CONSTRAINT email_chk CHECK (email_dom LIKE '%@%')`,
		`CREATE OR REPLACE DOMAIN phone_dom AS VARCHAR2(20)`,
		`CREATE DOMAIN IF NOT EXISTS age_dom AS NUMBER(3)
			CONSTRAINT age_chk CHECK (age_dom >= 0 AND age_dom <= 150)`,
		// With DEFAULT
		`CREATE DOMAIN status_dom AS VARCHAR2(10)
			DEFAULT 'ACTIVE'`,
		// With DISPLAY and ORDER
		`CREATE DOMAIN currency_dom AS NUMBER(10, 2)
			DISPLAY '$' || TO_CHAR(currency_dom)
			ORDER currency_dom`,
		// Multi-column domain
		`CREATE DOMAIN address_dom AS (
			street AS VARCHAR2(100),
			city AS VARCHAR2(50),
			zip AS VARCHAR2(10)
		)`,
		// STRICT
		`CREATE DOMAIN strict_dom AS NUMBER STRICT`,

		// ---- ALTER DOMAIN ----
		`ALTER DOMAIN email_dom RENAME TO email_address_dom`,

		// ---- CREATE INDEXTYPE ----
		`CREATE INDEXTYPE myschema.spatial_idx
			FOR myschema.contains_op (VARCHAR2)
			USING myschema.spatial_impl`,
		`CREATE OR REPLACE INDEXTYPE my_idx
			FOR my_op (NUMBER)
			USING my_impl_type`,
		`CREATE INDEXTYPE my_idx
			FOR my_op (NUMBER, VARCHAR2)
			USING my_impl
			WITH LOCAL PARTITION`,
		`CREATE INDEXTYPE my_idx
			FOR my_op (NUMBER)
			USING my_impl
			WITH SYSTEM MANAGED STORAGE TABLES`,
		`CREATE INDEXTYPE my_idx
			FOR my_op (NUMBER)
			USING my_impl
			WITH USER MANAGED STORAGE TABLES`,

		// ---- CREATE OPERATOR ----
		`CREATE OPERATOR my_op
			BINDING (NUMBER, NUMBER) RETURN NUMBER
			USING my_func`,
		`CREATE OR REPLACE OPERATOR my_op
			BINDING (VARCHAR2) RETURN NUMBER
			USING my_pkg.my_func`,

		// ---- CREATE MATERIALIZED ZONEMAP ----
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			ON sales (region_id, product_id)`,
		`CREATE MATERIALIZED ZONEMAP IF NOT EXISTS sales_zmap
			ON hr.sales (region_id)`,
		// With refresh
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			REFRESH FAST ON COMMIT
			ON sales (region_id)`,
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			REFRESH COMPLETE ON DEMAND
			ON sales (region_id)`,
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			REFRESH ON LOAD
			ON sales (region_id)`,
		// With attributes
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			TABLESPACE users
			SCALE 10
			PCTFREE 5
			CACHE
			ON sales (region_id)`,
		// ENABLE/DISABLE PRUNING
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			ENABLE PRUNING
			ON sales (region_id)`,
		`CREATE MATERIALIZED ZONEMAP sales_zmap
			DISABLE PRUNING
			ON sales (region_id)`,
		// AS subquery
		`CREATE MATERIALIZED ZONEMAP join_zmap
			AS SELECT s.region_id, p.category_id
			   FROM sales s, products p
			   WHERE s.product_id = p.product_id`,

		// ---- ALTER MATERIALIZED ZONEMAP ----
		`ALTER MATERIALIZED ZONEMAP sales_zmap ENABLE PRUNING`,
		`ALTER MATERIALIZED ZONEMAP sales_zmap DISABLE PRUNING`,
		`ALTER MATERIALIZED ZONEMAP sales_zmap REBUILD`,
		`ALTER MATERIALIZED ZONEMAP sales_zmap COMPILE`,
		`ALTER MATERIALIZED ZONEMAP sales_zmap REFRESH`,

		// ---- CREATE INMEMORY JOIN GROUP ----
		`CREATE INMEMORY JOIN GROUP dept_jg (
			hr.employees (department_id),
			hr.departments (department_id)
		)`,
		`CREATE INMEMORY JOIN GROUP IF NOT EXISTS dept_jg (
			hr.employees (department_id),
			hr.departments (department_id)
		)`,

		// ---- ALTER INMEMORY JOIN GROUP ----
		`ALTER INMEMORY JOIN GROUP dept_jg ADD (hr.locations (location_id))`,
		`ALTER INMEMORY JOIN GROUP dept_jg REMOVE (hr.locations (location_id))`,

		// ---- CREATE PROPERTY GRAPH ----
		`CREATE PROPERTY GRAPH my_graph
			VERTEX TABLES (persons)`,
		`CREATE OR REPLACE PROPERTY GRAPH my_graph
			VERTEX TABLES (persons KEY (person_id))`,
		`CREATE PROPERTY GRAPH IF NOT EXISTS my_graph
			VERTEX TABLES (persons)
			EDGE TABLES (friendships
				SOURCE KEY (person1_id) REFERENCES persons (person_id)
				DESTINATION KEY (person2_id) REFERENCES persons (person_id)
			)`,
		// With labels and properties
		`CREATE PROPERTY GRAPH my_graph
			VERTEX TABLES (
				persons KEY (id)
					LABEL person
					PROPERTIES (name, age)
			)
			EDGE TABLES (
				knows
					SOURCE KEY (src_id) REFERENCES persons (id)
					DESTINATION KEY (dst_id) REFERENCES persons (id)
					LABEL friend
					PROPERTIES (since)
			)`,
		// PROPERTIES ALL COLUMNS
		`CREATE PROPERTY GRAPH my_graph
			VERTEX TABLES (
				persons KEY (id) PROPERTIES ARE ALL COLUMNS
			)`,
		// NO PROPERTIES
		`CREATE PROPERTY GRAPH my_graph
			VERTEX TABLES (
				persons KEY (id) NO PROPERTIES
			)`,

		// ---- CREATE VECTOR INDEX ----
		`CREATE VECTOR INDEX vec_idx
			ON products (description_embedding)`,
		`CREATE VECTOR INDEX IF NOT EXISTS vec_idx
			ON products (embedding)`,
		// With DISTANCE
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			DISTANCE COSINE`,
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			DISTANCE EUCLIDEAN`,
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			DISTANCE DOT`,
		// With ORGANIZATION
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			ORGANIZATION INMEMORY NEIGHBOR GRAPH`,
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			ORGANIZATION NEIGHBOR PARTITIONS`,
		// WITH TARGET ACCURACY
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			DISTANCE COSINE
			WITH TARGET ACCURACY 95`,
		// With PARAMETERS
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			DISTANCE COSINE
			WITH TARGET ACCURACY 95 PARAMETERS (TYPE HNSW, NEIGHBORS 32, EFCONSTRUCTION 200)`,
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			WITH TARGET ACCURACY 90 PARAMETERS (TYPE IVF, NEIGHBOR PARTITIONS 100)`,
		// PARALLEL
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			DISTANCE COSINE
			PARALLEL 4`,
		// ONLINE
		`CREATE VECTOR INDEX vec_idx
			ON products (embedding)
			ONLINE`,

		// ---- CREATE MLE ENVIRONMENT ----
		`CREATE MLE ENV myenv`,
		`CREATE OR REPLACE MLE ENV myenv`,

		// ---- CREATE MLE MODULE ----
		`CREATE MLE MODULE mymod`,
		`CREATE OR REPLACE MLE MODULE mymod`,

		// ---- CREATE FLASHBACK ARCHIVE ----
		`CREATE FLASHBACK ARCHIVE fba1
			TABLESPACE fba_ts QUOTA 10G
			RETENTION 1 YEAR`,
		`CREATE FLASHBACK ARCHIVE DEFAULT fba_default
			TABLESPACE fba_ts QUOTA 50G
			RETENTION 3 YEAR`,

		// ---- ALTER FLASHBACK ARCHIVE ----
		`ALTER FLASHBACK ARCHIVE fba1 SET DEFAULT`,
		`ALTER FLASHBACK ARCHIVE fba1 ADD TABLESPACE fba_ts2 QUOTA 20G`,
		`ALTER FLASHBACK ARCHIVE fba1 MODIFY TABLESPACE fba_ts QUOTA 30G`,
		`ALTER FLASHBACK ARCHIVE fba1 MODIFY RETENTION 2 YEAR`,
		`ALTER FLASHBACK ARCHIVE fba1 PURGE ALL`,
		`ALTER FLASHBACK ARCHIVE fba1 PURGE BEFORE TIMESTAMP SYSTIMESTAMP - 30`,

		// ---- CREATE ROLLBACK SEGMENT ----
		`CREATE ROLLBACK SEGMENT rbs1`,
		`CREATE ROLLBACK SEGMENT rbs1 TABLESPACE rbs_ts`,
		`CREATE PUBLIC ROLLBACK SEGMENT rbs1`,

		// ---- ALTER ROLLBACK SEGMENT ----
		`ALTER ROLLBACK SEGMENT rbs1 ONLINE`,
		`ALTER ROLLBACK SEGMENT rbs1 OFFLINE`,
		`ALTER ROLLBACK SEGMENT rbs1 SHRINK`,
		`ALTER ROLLBACK SEGMENT rbs1 SHRINK TO 10M`,

		// ---- CREATE EDITION ----
		`CREATE EDITION e2`,
		`CREATE EDITION e2 AS CHILD OF e1`,

		// ---- DROP EDITION ----
		`DROP EDITION e2`,
		`DROP EDITION e2 CASCADE`,

		// ---- CREATE RESTORE POINT ----
		`CREATE RESTORE POINT before_upgrade`,
		`CREATE RESTORE POINT before_upgrade GUARANTEE FLASHBACK DATABASE`,
		`CREATE RESTORE POINT before_upgrade AS OF SCN 12345`,

		// ---- CREATE LOCKDOWN PROFILE ----
		`CREATE LOCKDOWN PROFILE lp1`,

		// ---- ALTER LOCKDOWN PROFILE ----
		`ALTER LOCKDOWN PROFILE lp1 DISABLE STATEMENT = ('ALTER SYSTEM')`,
		`ALTER LOCKDOWN PROFILE lp1 ENABLE STATEMENT = ('ALTER SESSION')`,
		`ALTER LOCKDOWN PROFILE lp1 DISABLE FEATURE = ('COMMON_SCHEMA_ACCESS')`,
		`ALTER LOCKDOWN PROFILE lp1 ENABLE FEATURE = ('COMMON_SCHEMA_ACCESS')`,
		`ALTER LOCKDOWN PROFILE lp1 DISABLE OPTION = ('PARALLEL_QUERY')`,

		// ---- CREATE OUTLINE ----
		`CREATE OUTLINE my_outline FOR CATEGORY special
			ON SELECT * FROM employees WHERE department_id = 10`,
		`CREATE OR REPLACE OUTLINE my_outline
			ON SELECT * FROM employees`,

		// ---- ALTER OUTLINE ----
		`ALTER OUTLINE my_outline REBUILD`,
		`ALTER OUTLINE my_outline RENAME TO new_outline`,
		`ALTER OUTLINE my_outline CHANGE CATEGORY TO default_cat`,
		`ALTER OUTLINE my_outline ENABLE`,
		`ALTER OUTLINE my_outline DISABLE`,

		// ---- CREATE PFILE / SPFILE ----
		`CREATE PFILE FROM SPFILE`,
		`CREATE PFILE = '/u01/admin/pfile.ora' FROM SPFILE`,
		`CREATE SPFILE FROM PFILE`,
		`CREATE SPFILE = '/u01/admin/spfile.ora' FROM PFILE`,
		`CREATE SPFILE FROM MEMORY`,

		// ---- CREATE TABLESPACE SET ----
		`CREATE TABLESPACE SET ts_set`,
		`CREATE TABLESPACE SET ts_set USING TEMPLATE
			(DATAFILE SIZE 100M AUTOEXTEND ON NEXT 50M)`,

		// ---- DROP TABLESPACE SET ----
		`DROP TABLESPACE SET ts_set`,
		`DROP TABLESPACE SET ts_set INCLUDING CONTENTS`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() < 1 {
				t.Fatalf("expected at least 1 statement, got %d", result.Len())
			}
		})
	}
}

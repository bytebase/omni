-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestWalkThroughInventory default#1
> query
| SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg\_%' AND nspname <> 'information_schema' ORDER BY 1
< result
< column "nspname" "NAME"
< row
: string
| public
> query
| SELECT evtname FROM pg_event_trigger ORDER BY 1
< result
< column "evtname" "NAME"
> exec
| RESET ALL
< ok 0 -
> exec
| RESET ROLE
< ok 0 -
> exec
| DISCARD TEMP
< ok 0 -
> exec
| DROP SCHEMA "public" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA public
< ok 0 -
> exec
|
| CREATE SCHEMA s;
| CREATE TABLE p (id integer PRIMARY KEY);
| CREATE TABLE t (id serial PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE, p_id integer CONSTRAINT t_p_fk REFERENCES p (id), n integer CONSTRAINT t_n_check CHECK (n > 0));
| CREATE INDEX t_n_idx ON t (n);
| CREATE TABLE s.t (id integer, a integer);
| CREATE VIEW tv AS SELECT code FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT id FROM p;
| CREATE INDEX mv_id_idx ON mv (id);
| CREATE SEQUENCE q;
| CREATE TYPE mood AS ENUM ('a', 'b');
| CREATE TYPE pair AS (x integer, y integer);
| CREATE TABLE parted (k integer, v text) PARTITION BY RANGE (k);
| CREATE TABLE s.parted_1 PARTITION OF parted FOR VALUES FROM (0) TO (10);
| CREATE INDEX parted_k_idx ON parted (k);
| CREATE TABLE ip (a integer CONSTRAINT ip_a_check CHECK (a > 0));
| CREATE TABLE ic (b integer) INHERITS (ip);
| CREATE TABLE nokey (a integer);
| CREATE DOMAIN dom AS int;
| CREATE TEMP TABLE tmp (a int);
< ok 0 -
> exec
| PREPARE walk_inventory AS WITH ns AS (
| 	SELECT n.oid, CASE WHEN n.oid = pg_catalog.pg_my_temp_schema() THEN 'pg_temp' ELSE n.nspname END AS name
| 	FROM pg_catalog.pg_namespace n
| 	WHERE (n.nspname NOT LIKE 'pg\_%' AND n.nspname <> 'information_schema') OR n.oid = pg_catalog.pg_my_temp_schema()
| ), rel AS (
| 	SELECT c.oid, ns.name AS schema, c.relname, c.relkind::text AS kind, coalesce(
| 		(SELECT t.relname FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class t ON t.oid = x.indrelid WHERE x.indexrelid = c.oid),
| 		(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 			WHERE d.classid = 'pg_catalog.pg_class'::pg_catalog.regclass AND d.objid = c.oid AND d.refclassid = d.classid AND d.deptype IN ('a', 'i') ORDER BY 1 LIMIT 1),
| 		c.relname) AS owner
| 	FROM pg_catalog.pg_class c JOIN ns ON ns.oid = c.relnamespace
| )
| SELECT coalesce(string_agg(line, E'\n' ORDER BY line COLLATE "C"), '') FROM (
| 	SELECT concat_ws(E'\t', 'schema', '', '', n.nspname, '', '') AS line FROM pg_catalog.pg_namespace n
| 		WHERE n.nspname NOT LIKE 'pg\_temp\_%' AND n.nspname NOT LIKE 'pg\_toast\_temp\_%'
| 	UNION ALL SELECT concat_ws(E'\t', 'relation', schema, '', relname, kind, owner) FROM rel
| 	UNION ALL SELECT concat_ws(E'\t', 'type', ns.name, '', t.typname, CASE WHEN t.typtype = 'd' THEN 'd' ELSE '' END,
| 			coalesce((SELECT r.owner FROM rel r WHERE r.oid = t.typrelid),
| 				(SELECT r.owner FROM pg_catalog.pg_type e JOIN rel r ON r.oid = e.typrelid WHERE e.oid = t.typelem), ''))
| 		FROM pg_catalog.pg_type t JOIN ns ON ns.oid = t.typnamespace
| 	UNION ALL SELECT concat_ws(E'\t', 'column', rel.schema, rel.relname, a.attname, '', rel.owner)
| 		FROM pg_catalog.pg_attribute a JOIN rel ON rel.oid = a.attrelid WHERE a.attnum > 0 AND NOT a.attisdropped
| 	UNION ALL SELECT concat_ws(E'\t', 'constraint', rel.schema, rel.relname, k.conname, k.contype::text, rel.owner)
| 		FROM pg_catalog.pg_constraint k JOIN rel ON rel.oid = k.conrelid
| ) o
< ok 0 -
> query
| EXECUTE walk_inventory
< result
< column "coalesce" "TEXT"
< row
: string
| column	pg_temp	tmp	a		tmp
| column	public	ic	a		ic
| column	public	ic	b		ic
| column	public	ip	a		ip
| column	public	mv	id		mv
| column	public	mv_id_idx	id		mv
| column	public	nokey	a		nokey
| column	public	p	id		p
| column	public	p_pkey	id		p
| column	public	pair	x		pair
| column	public	pair	y		pair
| column	public	parted	k		parted
| column	public	parted	v		parted
| column	public	parted_k_idx	k		parted
| column	public	q	is_called		q
| column	public	q	last_value		q
| column	public	q	log_cnt		q
| column	public	t	code		t
| column	public	t	id		t
| column	public	t	n		t
| column	public	t	p_id		t
| column	public	t_code_key	code		t
| column	public	t_id_seq	is_called		t
| column	public	t_id_seq	last_value		t
| column	public	t_id_seq	log_cnt		t
| column	public	t_n_idx	n		t
| column	public	t_pkey	id		t
| column	public	tv	code		tv
| column	s	parted_1	k		parted
| column	s	parted_1	v		parted
| column	s	parted_1_k_idx	k		parted_1
| column	s	t	a		t
| column	s	t	id		t
| constraint	public	ic	ip_a_check	c	ic
| constraint	public	ip	ip_a_check	c	ip
| constraint	public	p	p_pkey	p	p
| constraint	public	t	t_code_key	u	t
| constraint	public	t	t_n_check	c	t
| constraint	public	t	t_p_fk	f	t
| constraint	public	t	t_pkey	p	t
| relation	pg_temp		tmp	r	tmp
| relation	public		ic	r	ic
| relation	public		ip	r	ip
| relation	public		mv	m	mv
| relation	public		mv_id_idx	i	mv
| relation	public		nokey	r	nokey
| relation	public		p	r	p
| relation	public		p_pkey	i	p
| relation	public		pair	c	pair
| relation	public		parted	p	parted
| relation	public		parted_k_idx	I	parted
| relation	public		q	S	q
| relation	public		t	r	t
| relation	public		t_code_key	i	t
| relation	public		t_id_seq	S	t
| relation	public		t_n_idx	i	t
| relation	public		t_pkey	i	t
| relation	public		tv	v	tv
| relation	s		parted_1	r	parted
| relation	s		parted_1_k_idx	i	parted_1
| relation	s		t	r	t
| schema			information_schema		
| schema			pg_catalog		
| schema			pg_toast		
| schema			public		
| schema			s		
| type	pg_temp		_tmp		tmp
| type	pg_temp		tmp		tmp
| type	public		_dom		
| type	public		_ic		ic
| type	public		_ip		ip
| type	public		_mood		
| type	public		_mv		mv
| type	public		_nokey		nokey
| type	public		_p		p
| type	public		_pair		pair
| type	public		_parted		parted
| type	public		_t		t
| type	public		_tv		tv
| type	public		dom	d	
| type	public		ic		ic
| type	public		ip		ip
| type	public		mood		
| type	public		mv		mv
| type	public		nokey		nokey
| type	public		p		p
| type	public		pair		pair
| type	public		parted		parted
| type	public		t		t
| type	public		tv		tv
| type	s		_parted_1		parted
| type	s		_t		t
| type	s		parted_1		parted
| type	s		t		t
> exec
| DEALLOCATE walk_inventory
< ok 0 -

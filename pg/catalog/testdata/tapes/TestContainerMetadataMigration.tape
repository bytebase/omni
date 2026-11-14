-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestContainerMetadataMigration/column_ahead_of_key default#1
> exec
| CREATE TABLE "public"."issue" (
|     "id" bigint NOT NULL,
|     "project" text NOT NULL,
|     "plan_id" bigint
| );
|
| CREATE TABLE "public"."plan" (
|     "id" bigint NOT NULL,
|     "creator" text NOT NULL,
|     "project" text NOT NULL,
|     "config" jsonb NOT NULL
| );
|
| ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_pkey" PRIMARY KEY (project, id);
|
| CREATE INDEX "idx_plan_project" ON "public"."plan" (project);
|
| CREATE INDEX "idx_plan_config_has_rollout" ON "public"."plan" ((config ->> 'hasRollout'::text));
|
| CREATE INDEX "idx_plan_config" ON "public"."plan" USING GIN (config jsonb_path_ops);
|
| ALTER TABLE "public"."issue"
|     ADD CONSTRAINT "issue_project_plan_id_fkey" FOREIGN KEY ("project", "plan_id")
|     REFERENCES "public"."plan" ("project", "id");
|
|
< ok 0 -
> exec
| ALTER TABLE public.plan ADD COLUMN last_plan_editor text
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.issue.id bigint collation= notnull=t default= comment=
< row
: string
| column public.issue.plan_id bigint collation= notnull=f default= comment=
< row
: string
| column public.issue.project text collation= notnull=t default= comment=
< row
: string
| column public.plan.config jsonb collation= notnull=t default= comment=
< row
: string
| column public.plan.creator text collation= notnull=t default= comment=
< row
: string
| column public.plan.id bigint collation= notnull=t default= comment=
< row
: string
| column public.plan.last_plan_editor text collation= notnull=f default= comment=
< row
: string
| column public.plan.project text collation= notnull=t default= comment=
< row
: string
| constraint public.issue.issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES plan(project, id) comment=
< row
: string
| constraint public.plan.plan_pkey PRIMARY KEY (project, id) comment=
< row
: string
| relation public.idx_plan_config kind=i key= bound= index=CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops) comment=
< row
: string
| relation public.idx_plan_config_has_rollout kind=i key= bound= index=CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text))) comment=
< row
: string
| relation public.idx_plan_project kind=i key= bound= index=CREATE INDEX idx_plan_project ON public.plan USING btree (project) comment=
< row
: string
| relation public.issue kind=r key= bound= index= comment=
< row
: string
| relation public.plan kind=r key= bound= index= comment=
< row
: string
| relation public.plan_pkey kind=i key= bound= index=CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/column_ahead_of_key default#2
> exec
| CREATE TABLE "public"."issue" (
|     "id" bigint NOT NULL,
|     "project" text NOT NULL,
|     "plan_id" bigint
| );
|
| CREATE TABLE "public"."plan" (
|     "id" bigint NOT NULL,
|     "creator" text NOT NULL,
|     "last_plan_editor" text,
|     "project" text NOT NULL,
|     "config" jsonb NOT NULL
| );
|
| ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_pkey" PRIMARY KEY (project, id);
|
| CREATE INDEX "idx_plan_project" ON "public"."plan" (project);
|
| CREATE INDEX "idx_plan_config_has_rollout" ON "public"."plan" ((config ->> 'hasRollout'::text));
|
| CREATE INDEX "idx_plan_config" ON "public"."plan" USING GIN (config jsonb_path_ops);
|
| ALTER TABLE "public"."issue"
|     ADD CONSTRAINT "issue_project_plan_id_fkey" FOREIGN KEY ("project", "plan_id")
|     REFERENCES "public"."plan" ("project", "id");
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.issue.id bigint collation= notnull=t default= comment=
< row
: string
| column public.issue.plan_id bigint collation= notnull=f default= comment=
< row
: string
| column public.issue.project text collation= notnull=t default= comment=
< row
: string
| column public.plan.config jsonb collation= notnull=t default= comment=
< row
: string
| column public.plan.creator text collation= notnull=t default= comment=
< row
: string
| column public.plan.id bigint collation= notnull=t default= comment=
< row
: string
| column public.plan.last_plan_editor text collation= notnull=f default= comment=
< row
: string
| column public.plan.project text collation= notnull=t default= comment=
< row
: string
| constraint public.issue.issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES plan(project, id) comment=
< row
: string
| constraint public.plan.plan_pkey PRIMARY KEY (project, id) comment=
< row
: string
| relation public.idx_plan_config kind=i key= bound= index=CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops) comment=
< row
: string
| relation public.idx_plan_config_has_rollout kind=i key= bound= index=CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text))) comment=
< row
: string
| relation public.idx_plan_project kind=i key= bound= index=CREATE INDEX idx_plan_project ON public.plan USING btree (project) comment=
< row
: string
| relation public.issue kind=r key= bound= index= comment=
< row
: string
| relation public.plan kind=r key= bound= index= comment=
< row
: string
| relation public.plan_pkey kind=i key= bound= index=CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/column_ahead_of_key_dropped default#1
> exec
| CREATE TABLE "public"."issue" (
|     "id" bigint NOT NULL,
|     "project" text NOT NULL,
|     "plan_id" bigint
| );
|
| CREATE TABLE "public"."plan" (
|     "id" bigint NOT NULL,
|     "creator" text NOT NULL,
|     "last_plan_editor" text,
|     "project" text NOT NULL,
|     "config" jsonb NOT NULL
| );
|
| ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_pkey" PRIMARY KEY (project, id);
|
| CREATE INDEX "idx_plan_project" ON "public"."plan" (project);
|
| CREATE INDEX "idx_plan_config_has_rollout" ON "public"."plan" ((config ->> 'hasRollout'::text));
|
| CREATE INDEX "idx_plan_config" ON "public"."plan" USING GIN (config jsonb_path_ops);
|
| ALTER TABLE "public"."issue"
|     ADD CONSTRAINT "issue_project_plan_id_fkey" FOREIGN KEY ("project", "plan_id")
|     REFERENCES "public"."plan" ("project", "id");
|
|
< ok 0 -
> exec
| ALTER TABLE public.plan DROP COLUMN last_plan_editor
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.issue.id bigint collation= notnull=t default= comment=
< row
: string
| column public.issue.plan_id bigint collation= notnull=f default= comment=
< row
: string
| column public.issue.project text collation= notnull=t default= comment=
< row
: string
| column public.plan.config jsonb collation= notnull=t default= comment=
< row
: string
| column public.plan.creator text collation= notnull=t default= comment=
< row
: string
| column public.plan.id bigint collation= notnull=t default= comment=
< row
: string
| column public.plan.project text collation= notnull=t default= comment=
< row
: string
| constraint public.issue.issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES plan(project, id) comment=
< row
: string
| constraint public.plan.plan_pkey PRIMARY KEY (project, id) comment=
< row
: string
| relation public.idx_plan_config kind=i key= bound= index=CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops) comment=
< row
: string
| relation public.idx_plan_config_has_rollout kind=i key= bound= index=CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text))) comment=
< row
: string
| relation public.idx_plan_project kind=i key= bound= index=CREATE INDEX idx_plan_project ON public.plan USING btree (project) comment=
< row
: string
| relation public.issue kind=r key= bound= index= comment=
< row
: string
| relation public.plan kind=r key= bound= index= comment=
< row
: string
| relation public.plan_pkey kind=i key= bound= index=CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/column_ahead_of_key_dropped default#2
> exec
| CREATE TABLE "public"."issue" (
|     "id" bigint NOT NULL,
|     "project" text NOT NULL,
|     "plan_id" bigint
| );
|
| CREATE TABLE "public"."plan" (
|     "id" bigint NOT NULL,
|     "creator" text NOT NULL,
|     "project" text NOT NULL,
|     "config" jsonb NOT NULL
| );
|
| ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_pkey" PRIMARY KEY (project, id);
|
| CREATE INDEX "idx_plan_project" ON "public"."plan" (project);
|
| CREATE INDEX "idx_plan_config_has_rollout" ON "public"."plan" ((config ->> 'hasRollout'::text));
|
| CREATE INDEX "idx_plan_config" ON "public"."plan" USING GIN (config jsonb_path_ops);
|
| ALTER TABLE "public"."issue"
|     ADD CONSTRAINT "issue_project_plan_id_fkey" FOREIGN KEY ("project", "plan_id")
|     REFERENCES "public"."plan" ("project", "id");
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.issue.id bigint collation= notnull=t default= comment=
< row
: string
| column public.issue.plan_id bigint collation= notnull=f default= comment=
< row
: string
| column public.issue.project text collation= notnull=t default= comment=
< row
: string
| column public.plan.config jsonb collation= notnull=t default= comment=
< row
: string
| column public.plan.creator text collation= notnull=t default= comment=
< row
: string
| column public.plan.id bigint collation= notnull=t default= comment=
< row
: string
| column public.plan.project text collation= notnull=t default= comment=
< row
: string
| constraint public.issue.issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES plan(project, id) comment=
< row
: string
| constraint public.plan.plan_pkey PRIMARY KEY (project, id) comment=
< row
: string
| relation public.idx_plan_config kind=i key= bound= index=CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops) comment=
< row
: string
| relation public.idx_plan_config_has_rollout kind=i key= bound= index=CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text))) comment=
< row
: string
| relation public.idx_plan_project kind=i key= bound= index=CREATE INDEX idx_plan_project ON public.plan USING btree (project) comment=
< row
: string
| relation public.issue kind=r key= bound= index= comment=
< row
: string
| relation public.plan kind=r key= bound= index= comment=
< row
: string
| relation public.plan_pkey kind=i key= bound= index=CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/composite_collation default#1
> exec
| CREATE TYPE "public"."addr" AS (
|     "street" text
| );
|
|
< ok 0 -
> exec
| ALTER TYPE "public"."addr" ADD ATTRIBUTE "zip" text COLLATE "C";
| ALTER TYPE "public"."addr" ALTER ATTRIBUTE "street" TYPE text COLLATE "C"
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.addr.street text collation=C notnull=f default= comment=
< row
: string
| column public.addr.zip text collation=C notnull=f default= comment=
< row
: string
| relation public.addr kind=c key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| type public.addr c comment=

== TestContainerMetadataMigration/composite_collation default#2
> exec
| CREATE TYPE "public"."addr" AS (
|     "street" text COLLATE "C",
|     "zip" text COLLATE "C"
| );
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.addr.street text collation=C notnull=f default= comment=
< row
: string
| column public.addr.zip text collation=C notnull=f default= comment=
< row
: string
| relation public.addr kind=c key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| type public.addr c comment=

== TestContainerMetadataMigration/composite_comments default#1
> exec
| CREATE TYPE "public"."addr" AS (
|     "street" text
| );
|
| COMMENT ON TYPE "public"."addr" IS 'old type comment';
|
| COMMENT ON COLUMN "public"."addr"."street" IS 'old attr comment';
|
|
< ok 0 -
> exec
| COMMENT ON COLUMN "public"."addr"."street" IS 'new attr comment';
| COMMENT ON TYPE "public"."addr" IS 'new type comment'
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.addr.street text collation= notnull=f default= comment=new attr comment
< row
: string
| relation public.addr kind=c key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| type public.addr c comment=new type comment

== TestContainerMetadataMigration/composite_comments default#2
> exec
| CREATE TYPE "public"."addr" AS (
|     "street" text
| );
|
| COMMENT ON TYPE "public"."addr" IS 'new type comment';
|
| COMMENT ON COLUMN "public"."addr"."street" IS 'new attr comment';
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.addr.street text collation= notnull=f default= comment=new attr comment
< row
: string
| relation public.addr kind=c key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| type public.addr c comment=new type comment

== TestContainerMetadataMigration/create_extension default#1
> exec
|
< ok 0 -
> exec
| CREATE EXTENSION "pg_trgm" WITH SCHEMA "public" VERSION '1.6'
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| extension pg_trgm 1.6 in public
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/create_extension default#2
> exec
| CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public" VERSION '1.6';
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| extension pg_trgm 1.6 in public
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/create_with_expression_index default#1
> exec
|
< ok 0 -
> exec
| CREATE TABLE "public"."issue" (
|     "id" bigint NOT NULL,
|     "project" text NOT NULL,
|     "plan_id" bigint
| );
| CREATE TABLE "public"."plan" (
|     "id" bigint NOT NULL,
|     "creator" text NOT NULL,
|     "project" text NOT NULL,
|     "config" jsonb NOT NULL,
|     CONSTRAINT "plan_pkey" PRIMARY KEY ("project", "id")
| );
| CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops);
| CREATE INDEX idx_plan_config_has_rollout ON public.plan ((config ->> 'hasRollout'::text));
| CREATE INDEX idx_plan_project ON public.plan (project);
| ALTER TABLE public.issue ADD CONSTRAINT issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES public.plan (project, id)
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.issue.id bigint collation= notnull=t default= comment=
< row
: string
| column public.issue.plan_id bigint collation= notnull=f default= comment=
< row
: string
| column public.issue.project text collation= notnull=t default= comment=
< row
: string
| column public.plan.config jsonb collation= notnull=t default= comment=
< row
: string
| column public.plan.creator text collation= notnull=t default= comment=
< row
: string
| column public.plan.id bigint collation= notnull=t default= comment=
< row
: string
| column public.plan.project text collation= notnull=t default= comment=
< row
: string
| constraint public.issue.issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES plan(project, id) comment=
< row
: string
| constraint public.plan.plan_pkey PRIMARY KEY (project, id) comment=
< row
: string
| relation public.idx_plan_config kind=i key= bound= index=CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops) comment=
< row
: string
| relation public.idx_plan_config_has_rollout kind=i key= bound= index=CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text))) comment=
< row
: string
| relation public.idx_plan_project kind=i key= bound= index=CREATE INDEX idx_plan_project ON public.plan USING btree (project) comment=
< row
: string
| relation public.issue kind=r key= bound= index= comment=
< row
: string
| relation public.plan kind=r key= bound= index= comment=
< row
: string
| relation public.plan_pkey kind=i key= bound= index=CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/create_with_expression_index default#2
> exec
| CREATE TABLE "public"."issue" (
|     "id" bigint NOT NULL,
|     "project" text NOT NULL,
|     "plan_id" bigint
| );
|
| CREATE TABLE "public"."plan" (
|     "id" bigint NOT NULL,
|     "creator" text NOT NULL,
|     "project" text NOT NULL,
|     "config" jsonb NOT NULL
| );
|
| ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_pkey" PRIMARY KEY (project, id);
|
| CREATE INDEX "idx_plan_project" ON "public"."plan" (project);
|
| CREATE INDEX "idx_plan_config_has_rollout" ON "public"."plan" ((config ->> 'hasRollout'::text));
|
| CREATE INDEX "idx_plan_config" ON "public"."plan" USING GIN (config jsonb_path_ops);
|
| ALTER TABLE "public"."issue"
|     ADD CONSTRAINT "issue_project_plan_id_fkey" FOREIGN KEY ("project", "plan_id")
|     REFERENCES "public"."plan" ("project", "id");
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.issue.id bigint collation= notnull=t default= comment=
< row
: string
| column public.issue.plan_id bigint collation= notnull=f default= comment=
< row
: string
| column public.issue.project text collation= notnull=t default= comment=
< row
: string
| column public.plan.config jsonb collation= notnull=t default= comment=
< row
: string
| column public.plan.creator text collation= notnull=t default= comment=
< row
: string
| column public.plan.id bigint collation= notnull=t default= comment=
< row
: string
| column public.plan.project text collation= notnull=t default= comment=
< row
: string
| constraint public.issue.issue_project_plan_id_fkey FOREIGN KEY (project, plan_id) REFERENCES plan(project, id) comment=
< row
: string
| constraint public.plan.plan_pkey PRIMARY KEY (project, id) comment=
< row
: string
| relation public.idx_plan_config kind=i key= bound= index=CREATE INDEX idx_plan_config ON public.plan USING gin (config jsonb_path_ops) comment=
< row
: string
| relation public.idx_plan_config_has_rollout kind=i key= bound= index=CREATE INDEX idx_plan_config_has_rollout ON public.plan USING btree (((config ->> 'hasRollout'::text))) comment=
< row
: string
| relation public.idx_plan_project kind=i key= bound= index=CREATE INDEX idx_plan_project ON public.plan USING btree (project) comment=
< row
: string
| relation public.issue kind=r key= bound= index= comment=
< row
: string
| relation public.plan kind=r key= bound= index= comment=
< row
: string
| relation public.plan_pkey kind=i key= bound= index=CREATE UNIQUE INDEX plan_pkey ON public.plan USING btree (project, id) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/drop_extension default#1
> exec
| CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public" VERSION '1.6';
|
|
< ok 0 -
> exec
| DROP EXTENSION "pg_trgm"
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/drop_extension default#2
> exec
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/extension_member default#1
> exec
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/extension_member default#2
> exec
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/owned_sequence default#1
> exec
| CREATE SEQUENCE "public"."order_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	NO CYCLE;
|
|
< ok 0 -
> exec
| CREATE TABLE "public"."orders" (
|     "id" bigint NOT NULL
| );
| ALTER SEQUENCE "public"."order_seq" OWNED BY "public"."orders"."id"
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.orders.id bigint collation= notnull=t default= comment=
< row
: string
| relation public.order_seq kind=S key= bound= index= comment=
< row
: string
| relation public.orders kind=r key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| sequence public.order_seq integer start=1 increment=1 min=1 max=2147483647 cache=1 cycle=f owner=orders.id

== TestContainerMetadataMigration/owned_sequence default#2
> exec
| CREATE SEQUENCE "public"."order_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	NO CYCLE;
|
| CREATE TABLE "public"."orders" (
|     "id" bigint NOT NULL
| );
|
| ALTER SEQUENCE "public"."order_seq" OWNED BY "public"."orders"."id";
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.orders.id bigint collation= notnull=t default= comment=
< row
: string
| relation public.order_seq kind=S key= bound= index= comment=
< row
: string
| relation public.orders kind=r key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| sequence public.order_seq integer start=1 increment=1 min=1 max=2147483647 cache=1 cycle=f owner=orders.id

== TestContainerMetadataMigration/partition default#1
> exec
| CREATE TABLE "public"."orders" (
|     "id" integer NOT NULL,
|     "created_at" date NOT NULL
| ) PARTITION BY RANGE (created_at);
|
| CREATE TABLE "public"."orders_2023" (
|     "id" integer NOT NULL,
|     "created_at" date NOT NULL
| );
|
| ALTER TABLE ONLY "public"."orders" ATTACH PARTITION "public"."orders_2023" FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');
|
|
< ok 0 -
> exec
| DROP TABLE "public"."orders_2023" CASCADE;
| CREATE TABLE "public"."orders_2024" PARTITION OF "public"."orders" FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
| CREATE TABLE "public"."orders_other" PARTITION OF "public"."orders" DEFAULT
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.orders.created_at date collation= notnull=t default= comment=
< row
: string
| column public.orders.id integer collation= notnull=t default= comment=
< row
: string
| column public.orders_2024.created_at date collation= notnull=t default= comment=
< row
: string
| column public.orders_2024.id integer collation= notnull=t default= comment=
< row
: string
| column public.orders_other.created_at date collation= notnull=t default= comment=
< row
: string
| column public.orders_other.id integer collation= notnull=t default= comment=
< row
: string
| relation public.orders kind=p key=RANGE (created_at) bound= index= comment=
< row
: string
| relation public.orders_2024 kind=r key= bound=FOR VALUES FROM ('2024-01-01') TO ('2025-01-01') index= comment=
< row
: string
| relation public.orders_other kind=r key= bound=DEFAULT index= comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/partition default#2
> exec
| CREATE TABLE "public"."orders" (
|     "id" integer NOT NULL,
|     "created_at" date NOT NULL
| ) PARTITION BY RANGE (created_at);
|
| CREATE TABLE "public"."orders_2024" (
|     "id" integer NOT NULL,
|     "created_at" date NOT NULL
| );
|
| CREATE TABLE "public"."orders_other" (
|     "id" integer NOT NULL,
|     "created_at" date NOT NULL
| );
|
| ALTER TABLE ONLY "public"."orders" ATTACH PARTITION "public"."orders_2024" FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
|
| ALTER TABLE ONLY "public"."orders" ATTACH PARTITION "public"."orders_other" DEFAULT;
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.orders.created_at date collation= notnull=t default= comment=
< row
: string
| column public.orders.id integer collation= notnull=t default= comment=
< row
: string
| column public.orders_2024.created_at date collation= notnull=t default= comment=
< row
: string
| column public.orders_2024.id integer collation= notnull=t default= comment=
< row
: string
| column public.orders_other.created_at date collation= notnull=t default= comment=
< row
: string
| column public.orders_other.id integer collation= notnull=t default= comment=
< row
: string
| relation public.orders kind=p key=RANGE (created_at) bound= index= comment=
< row
: string
| relation public.orders_2024 kind=r key= bound=FOR VALUES FROM ('2024-01-01') TO ('2025-01-01') index= comment=
< row
: string
| relation public.orders_other kind=r key= bound=DEFAULT index= comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/partition_with_indexes default#1
> exec
| CREATE TABLE "public"."sales" (
|     "id" integer NOT NULL,
|     "region" text NOT NULL,
|     "d" date NOT NULL
| ) PARTITION BY RANGE (d);
|
| CREATE TABLE "public"."sales_2024" (
|     "id" integer NOT NULL,
|     "region" text NOT NULL,
|     "d" date NOT NULL
| );
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2024" FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
|
| ALTER TABLE ONLY "public"."sales" ADD CONSTRAINT "sales_pkey" PRIMARY KEY (id, d);
|
| ALTER TABLE ONLY "public"."sales_2024" ADD CONSTRAINT "sales_2024_pkey" PRIMARY KEY (id, d);
|
| CREATE INDEX "sales_region_idx" ON ONLY "public"."sales" (region);
|
| CREATE INDEX "sales_lower_idx" ON ONLY "public"."sales" (lower(region));
|
| CREATE INDEX "sales_2024_region_idx" ON "public"."sales_2024" (region);
|
| CREATE INDEX "sales_2024_lower_idx" ON "public"."sales_2024" (lower(region));
|
| ALTER INDEX "public"."sales_pkey" ATTACH PARTITION "public"."sales_2024_pkey";
|
| ALTER INDEX "public"."sales_region_idx" ATTACH PARTITION "public"."sales_2024_region_idx";
|
| ALTER INDEX "public"."sales_lower_idx" ATTACH PARTITION "public"."sales_2024_lower_idx";
|
|
< ok 0 -
> exec
| CREATE TABLE "public"."sales_2025" PARTITION OF "public"."sales" FOR VALUES FROM ('2025-01-01') TO ('2026-01-01')
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.sales.d date collation= notnull=t default= comment=
< row
: string
| column public.sales.id integer collation= notnull=t default= comment=
< row
: string
| column public.sales.region text collation= notnull=t default= comment=
< row
: string
| column public.sales_2024.d date collation= notnull=t default= comment=
< row
: string
| column public.sales_2024.id integer collation= notnull=t default= comment=
< row
: string
| column public.sales_2024.region text collation= notnull=t default= comment=
< row
: string
| column public.sales_2025.d date collation= notnull=t default= comment=
< row
: string
| column public.sales_2025.id integer collation= notnull=t default= comment=
< row
: string
| column public.sales_2025.region text collation= notnull=t default= comment=
< row
: string
| constraint public.sales.sales_pkey PRIMARY KEY (id, d) comment=
< row
: string
| constraint public.sales_2024.sales_2024_pkey PRIMARY KEY (id, d) comment=
< row
: string
| constraint public.sales_2025.sales_2025_pkey PRIMARY KEY (id, d) comment=
< row
: string
| relation public.sales kind=p key=RANGE (d) bound= index= comment=
< row
: string
| relation public.sales_2024 kind=r key= bound=FOR VALUES FROM ('2024-01-01') TO ('2025-01-01') index= comment=
< row
: string
| relation public.sales_2024_lower_idx kind=i key= bound= index=CREATE INDEX sales_2024_lower_idx ON public.sales_2024 USING btree (lower(region)) comment=
< row
: string
| relation public.sales_2024_pkey kind=i key= bound= index=CREATE UNIQUE INDEX sales_2024_pkey ON public.sales_2024 USING btree (id, d) comment=
< row
: string
| relation public.sales_2024_region_idx kind=i key= bound= index=CREATE INDEX sales_2024_region_idx ON public.sales_2024 USING btree (region) comment=
< row
: string
| relation public.sales_2025 kind=r key= bound=FOR VALUES FROM ('2025-01-01') TO ('2026-01-01') index= comment=
< row
: string
| relation public.sales_2025_lower_idx kind=i key= bound= index=CREATE INDEX sales_2025_lower_idx ON public.sales_2025 USING btree (lower(region)) comment=
< row
: string
| relation public.sales_2025_pkey kind=i key= bound= index=CREATE UNIQUE INDEX sales_2025_pkey ON public.sales_2025 USING btree (id, d) comment=
< row
: string
| relation public.sales_2025_region_idx kind=i key= bound= index=CREATE INDEX sales_2025_region_idx ON public.sales_2025 USING btree (region) comment=
< row
: string
| relation public.sales_lower_idx kind=I key= bound= index=CREATE INDEX sales_lower_idx ON ONLY public.sales USING btree (lower(region)) comment=
< row
: string
| relation public.sales_pkey kind=I key= bound= index=CREATE UNIQUE INDEX sales_pkey ON ONLY public.sales USING btree (id, d) comment=
< row
: string
| relation public.sales_region_idx kind=I key= bound= index=CREATE INDEX sales_region_idx ON ONLY public.sales USING btree (region) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/partition_with_indexes default#2
> exec
| CREATE TABLE "public"."sales" (
|     "id" integer NOT NULL,
|     "region" text NOT NULL,
|     "d" date NOT NULL
| ) PARTITION BY RANGE (d);
|
| CREATE TABLE "public"."sales_2024" (
|     "id" integer NOT NULL,
|     "region" text NOT NULL,
|     "d" date NOT NULL
| );
|
| CREATE TABLE "public"."sales_2025" (
|     "id" integer NOT NULL,
|     "region" text NOT NULL,
|     "d" date NOT NULL
| );
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2024" FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2025" FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
|
| ALTER TABLE ONLY "public"."sales" ADD CONSTRAINT "sales_pkey" PRIMARY KEY (id, d);
|
| ALTER TABLE ONLY "public"."sales_2024" ADD CONSTRAINT "sales_2024_pkey" PRIMARY KEY (id, d);
|
| ALTER TABLE ONLY "public"."sales_2025" ADD CONSTRAINT "sales_2025_pkey" PRIMARY KEY (id, d);
|
| CREATE INDEX "sales_region_idx" ON ONLY "public"."sales" (region);
|
| CREATE INDEX "sales_lower_idx" ON ONLY "public"."sales" (lower(region));
|
| CREATE INDEX "sales_2024_region_idx" ON "public"."sales_2024" (region);
|
| CREATE INDEX "sales_2024_lower_idx" ON "public"."sales_2024" (lower(region));
|
| CREATE INDEX "sales_2025_region_idx" ON "public"."sales_2025" (region);
|
| CREATE INDEX "sales_2025_lower_idx" ON "public"."sales_2025" (lower(region));
|
| ALTER INDEX "public"."sales_pkey" ATTACH PARTITION "public"."sales_2024_pkey";
|
| ALTER INDEX "public"."sales_region_idx" ATTACH PARTITION "public"."sales_2024_region_idx";
|
| ALTER INDEX "public"."sales_lower_idx" ATTACH PARTITION "public"."sales_2024_lower_idx";
|
| ALTER INDEX "public"."sales_pkey" ATTACH PARTITION "public"."sales_2025_pkey";
|
| ALTER INDEX "public"."sales_region_idx" ATTACH PARTITION "public"."sales_2025_region_idx";
|
| ALTER INDEX "public"."sales_lower_idx" ATTACH PARTITION "public"."sales_2025_lower_idx";
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| column public.sales.d date collation= notnull=t default= comment=
< row
: string
| column public.sales.id integer collation= notnull=t default= comment=
< row
: string
| column public.sales.region text collation= notnull=t default= comment=
< row
: string
| column public.sales_2024.d date collation= notnull=t default= comment=
< row
: string
| column public.sales_2024.id integer collation= notnull=t default= comment=
< row
: string
| column public.sales_2024.region text collation= notnull=t default= comment=
< row
: string
| column public.sales_2025.d date collation= notnull=t default= comment=
< row
: string
| column public.sales_2025.id integer collation= notnull=t default= comment=
< row
: string
| column public.sales_2025.region text collation= notnull=t default= comment=
< row
: string
| constraint public.sales.sales_pkey PRIMARY KEY (id, d) comment=
< row
: string
| constraint public.sales_2024.sales_2024_pkey PRIMARY KEY (id, d) comment=
< row
: string
| constraint public.sales_2025.sales_2025_pkey PRIMARY KEY (id, d) comment=
< row
: string
| relation public.sales kind=p key=RANGE (d) bound= index= comment=
< row
: string
| relation public.sales_2024 kind=r key= bound=FOR VALUES FROM ('2024-01-01') TO ('2025-01-01') index= comment=
< row
: string
| relation public.sales_2024_lower_idx kind=i key= bound= index=CREATE INDEX sales_2024_lower_idx ON public.sales_2024 USING btree (lower(region)) comment=
< row
: string
| relation public.sales_2024_pkey kind=i key= bound= index=CREATE UNIQUE INDEX sales_2024_pkey ON public.sales_2024 USING btree (id, d) comment=
< row
: string
| relation public.sales_2024_region_idx kind=i key= bound= index=CREATE INDEX sales_2024_region_idx ON public.sales_2024 USING btree (region) comment=
< row
: string
| relation public.sales_2025 kind=r key= bound=FOR VALUES FROM ('2025-01-01') TO ('2026-01-01') index= comment=
< row
: string
| relation public.sales_2025_lower_idx kind=i key= bound= index=CREATE INDEX sales_2025_lower_idx ON public.sales_2025 USING btree (lower(region)) comment=
< row
: string
| relation public.sales_2025_pkey kind=i key= bound= index=CREATE UNIQUE INDEX sales_2025_pkey ON public.sales_2025 USING btree (id, d) comment=
< row
: string
| relation public.sales_2025_region_idx kind=i key= bound= index=CREATE INDEX sales_2025_region_idx ON public.sales_2025 USING btree (region) comment=
< row
: string
| relation public.sales_lower_idx kind=I key= bound= index=CREATE INDEX sales_lower_idx ON ONLY public.sales USING btree (lower(region)) comment=
< row
: string
| relation public.sales_pkey kind=I key= bound= index=CREATE UNIQUE INDEX sales_pkey ON ONLY public.sales USING btree (id, d) comment=
< row
: string
| relation public.sales_region_idx kind=I key= bound= index=CREATE INDEX sales_region_idx ON ONLY public.sales USING btree (region) comment=
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/released_sequence default#1
> exec
| CREATE SEQUENCE "public"."order_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	NO CYCLE;
|
| CREATE TABLE "public"."orders" (
|     "id" bigint NOT NULL
| );
|
| ALTER SEQUENCE "public"."order_seq" OWNED BY "public"."orders"."id";
|
|
< ok 0 -
> exec
| ALTER SEQUENCE "public"."order_seq" OWNED BY NONE;
| DROP TABLE "public"."orders" CASCADE
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| relation public.order_seq kind=S key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| sequence public.order_seq integer start=1 increment=1 min=1 max=2147483647 cache=1 cycle=f owner=

== TestContainerMetadataMigration/released_sequence default#2
> exec
| CREATE SEQUENCE "public"."order_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	NO CYCLE;
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| relation public.order_seq kind=S key= bound= index= comment=
< row
: string
| schema public comment=standard public schema
< row
: string
| sequence public.order_seq integer start=1 increment=1 min=1 max=2147483647 cache=1 cycle=f owner=

== TestContainerMetadataMigration/update_extension default#1
> exec
| CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public" VERSION '1.5';
|
|
< ok 0 -
> exec
| ALTER EXTENSION "pg_trgm" UPDATE TO '1.6'
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| extension pg_trgm 1.6 in public
< row
: string
| schema public comment=standard public schema

== TestContainerMetadataMigration/update_extension default#2
> exec
| CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public" VERSION '1.6';
|
|
< ok 0 -
> query
|
| WITH ns AS (
|   SELECT oid, nspname FROM pg_namespace
|   WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg_toast%' AND nspname NOT LIKE 'pg_temp%'
| ), member AS (
|   SELECT objid FROM pg_depend WHERE deptype = 'e'
| )
| SELECT format('schema %s comment=%s', nspname, obj_description(oid, 'pg_namespace')) FROM ns
| UNION ALL
| SELECT format('extension %s %s in %s', e.extname, e.extversion, n.nspname)
|   FROM pg_extension e JOIN pg_namespace n ON n.oid = e.extnamespace WHERE e.extname <> 'plpgsql'
| UNION ALL
| SELECT format('relation %s.%s kind=%s key=%s bound=%s index=%s comment=%s', ns.nspname, c.relname, c.relkind,
|     pg_get_partkeydef(c.oid), pg_get_expr(c.relpartbound, c.oid),
|     CASE WHEN c.relkind IN ('i', 'I') THEN pg_get_indexdef(c.oid) END, obj_description(c.oid, 'pg_class'))
|   FROM pg_class c JOIN ns ON ns.oid = c.relnamespace
|   WHERE c.relkind IN ('r', 'p', 'v', 'm', 'S', 'c', 'i', 'I') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('column %s.%s.%s %s collation=%s notnull=%s default=%s comment=%s', ns.nspname, c.relname, a.attname,
|     format_type(a.atttypid, a.atttypmod), co.collname, a.attnotnull, pg_get_expr(d.adbin, d.adrelid), col_description(c.oid, a.attnum))
|   FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN ns ON ns.oid = c.relnamespace
|   JOIN pg_type ty ON ty.oid = a.atttypid
|   LEFT JOIN pg_collation co ON co.oid = a.attcollation AND a.attcollation <> ty.typcollation
|   LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
|   WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind IN ('r', 'p', 'v', 'm', 'c') AND c.oid NOT IN (SELECT objid FROM member)
| UNION ALL
| SELECT format('constraint %s.%s.%s %s comment=%s', ns.nspname, c.relname, k.conname, pg_get_constraintdef(k.oid), obj_description(k.oid, 'pg_constraint'))
|   FROM pg_constraint k JOIN pg_class c ON c.oid = k.conrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('sequence %s.%s %s start=%s increment=%s min=%s max=%s cache=%s cycle=%s owner=%s', ns.nspname, c.relname,
|     format_type(s.seqtypid, NULL), s.seqstart, s.seqincrement, s.seqmin, s.seqmax, s.seqcache, s.seqcycle,
|     (SELECT format('%s.%s', o.relname, a.attname) FROM pg_depend d
|        JOIN pg_class o ON o.oid = d.refobjid JOIN pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
|        WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i')))
|   FROM pg_sequence s JOIN pg_class c ON c.oid = s.seqrelid JOIN ns ON ns.oid = c.relnamespace
| UNION ALL
| SELECT format('type %s.%s %s comment=%s', ns.nspname, ty.typname, ty.typtype, obj_description(ty.oid, 'pg_type'))
|   FROM pg_type ty JOIN ns ON ns.oid = ty.typnamespace
|   WHERE (ty.typtype IN ('e', 'd') OR ty.typtype = 'c' AND (SELECT relkind FROM pg_class WHERE oid = ty.typrelid) = 'c')
|     AND ty.oid NOT IN (SELECT objid FROM member)
| ORDER BY 1
< result
< column "format" "TEXT"
< row
: string
| extension pg_trgm 1.6 in public
< row
: string
| schema public comment=standard public schema

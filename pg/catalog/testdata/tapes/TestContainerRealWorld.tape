-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestContainerRealWorld default#1
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_1"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_2"
< ok 0 -
> exec
| SET search_path TO "test_1", public;
| CREATE TYPE "priority_level" AS ENUM ('low', 'medium', 'high', 'critical');
| CREATE TYPE "user_status" AS ENUM ('active', 'inactive', 'banned');
| CREATE SEQUENCE "order_id_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 1000
|     CACHE 1
|     NO CYCLE;
| CREATE SEQUENCE "user_id_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 1
|     CACHE 1
|     NO CYCLE;
| CREATE FUNCTION "log_change"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN RETURN NEW; END;$$;
| CREATE TABLE "audit_log" (
|     "id" integer NOT NULL,
|     "table_name" text NOT NULL,
|     "action" text NOT NULL,
|     "old_data" jsonb,
|     "new_data" jsonb,
|     "performed_at" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "audit_log_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "categories" (
|     "id" integer NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "parent_id" integer,
|     CONSTRAINT "categories_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "orders" (
|     "id" integer NOT NULL DEFAULT nextval('order_id_seq'::regclass),
|     "user_id" integer NOT NULL,
|     "product_id" integer NOT NULL,
|     "quantity" integer NOT NULL,
|     "total" numeric(12,2),
|     "created_at" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "orders_pkey" PRIMARY KEY ("id")
| );
| CREATE FUNCTION "get_user_orders"("uid" integer)
| RETURNS bigint
| LANGUAGE sql
| STABLE
| AS $$SELECT count(*) FROM orders WHERE user_id = uid$$;
| CREATE TABLE "products" (
|     "id" integer NOT NULL,
|     "name" character varying(200) NOT NULL,
|     "category_id" integer NOT NULL,
|     "price" numeric(10,2) NOT NULL DEFAULT 0,
|     "priority" priority_level DEFAULT 'medium'::priority_level,
|     "in_stock" boolean DEFAULT true,
|     CONSTRAINT "products_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "users" (
|     "id" integer NOT NULL DEFAULT nextval('user_id_seq'::regclass),
|     "username" character varying(100) NOT NULL,
|     "email" text,
|     "status" user_status DEFAULT 'active'::user_status,
|     "active" boolean DEFAULT true,
|     "created_at" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
|     CONSTRAINT "users_username_unique" UNIQUE ("username")
| );
| ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
| CREATE VIEW "active_users" AS  SELECT id,
|     username,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
| CREATE VIEW "order_summary" AS  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((orders o
|      JOIN users u ON ((u.id = o.user_id)))
|      JOIN products p ON ((p.id = o.product_id)));
| CREATE INDEX idx_audit_action ON audit_log (action);
| CREATE INDEX idx_audit_performed ON audit_log (performed_at);
| CREATE INDEX idx_orders_created ON orders (created_at);
| CREATE INDEX idx_orders_user ON orders (user_id);
| CREATE INDEX idx_products_category ON products (category_id);
| CREATE INDEX idx_products_price ON products (price);
| CREATE INDEX idx_users_email ON users (email);
| CREATE INDEX idx_users_status ON users (status);
| CREATE TRIGGER "orders_audit" AFTER INSERT ON "orders" FOR EACH ROW EXECUTE FUNCTION "log_change"();
| CREATE TRIGGER "users_audit" AFTER UPDATE ON "users" FOR EACH ROW EXECUTE FUNCTION "log_change"();
| CREATE POLICY "users_see_own" ON "users" FOR SELECT TO PUBLIC USING ((active = true));
| COMMENT ON VIEW "active_users" IS 'Only active users';
| COMMENT ON VIEW "order_summary" IS 'Joined order details';
| COMMENT ON TABLE "orders" IS 'Customer orders';
| COMMENT ON TABLE "products" IS 'Product catalog';
| COMMENT ON TABLE "users" IS 'Core user accounts';
| ALTER TABLE categories ADD CONSTRAINT categories_parent_fk FOREIGN KEY (parent_id) REFERENCES categories (id);
| ALTER TABLE orders ADD CONSTRAINT orders_product_fk FOREIGN KEY (product_id) REFERENCES products (id);
| ALTER TABLE orders ADD CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users (id);
| ALTER TABLE products ADD CONSTRAINT products_category_fk FOREIGN KEY (category_id) REFERENCES categories (id)
< ok 0 -
> exec
| SET search_path TO "test_2", public;
|
| CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned');
| CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');
|
| CREATE SEQUENCE user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1 NO CYCLE;
| CREATE SEQUENCE order_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1000 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_id_seq'),
|     username varchar(100) NOT NULL,
|     email text,
|     status user_status DEFAULT 'active',
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT users_username_unique UNIQUE (username)
| );
|
| CREATE TABLE categories (
|     id integer NOT NULL,
|     name varchar(100) NOT NULL,
|     parent_id integer,
|     CONSTRAINT categories_pkey PRIMARY KEY (id),
|     CONSTRAINT categories_parent_fk FOREIGN KEY (parent_id) REFERENCES categories(id)
| );
|
| CREATE TABLE products (
|     id integer NOT NULL,
|     name varchar(200) NOT NULL,
|     category_id integer NOT NULL,
|     price numeric(10,2) NOT NULL DEFAULT 0,
|     priority priority_level DEFAULT 'medium',
|     in_stock boolean DEFAULT true,
|     CONSTRAINT products_pkey PRIMARY KEY (id),
|     CONSTRAINT products_category_fk FOREIGN KEY (category_id) REFERENCES categories(id)
| );
|
| CREATE TABLE orders (
|     id integer NOT NULL DEFAULT nextval('order_id_seq'),
|     user_id integer NOT NULL,
|     product_id integer NOT NULL,
|     quantity integer NOT NULL,
|     total numeric(12,2),
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT orders_pkey PRIMARY KEY (id),
|     CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users(id),
|     CONSTRAINT orders_product_fk FOREIGN KEY (product_id) REFERENCES products(id)
| );
|
| CREATE TABLE audit_log (
|     id integer NOT NULL,
|     table_name text NOT NULL,
|     action text NOT NULL,
|     old_data jsonb,
|     new_data jsonb,
|     performed_at timestamptz DEFAULT now(),
|     CONSTRAINT audit_log_pkey PRIMARY KEY (id)
| );
|
| CREATE INDEX idx_users_status ON users(status);
| CREATE INDEX idx_users_email ON users(email);
| CREATE INDEX idx_products_category ON products(category_id);
| CREATE INDEX idx_products_price ON products(price);
| CREATE INDEX idx_orders_user ON orders(user_id);
| CREATE INDEX idx_orders_created ON orders(created_at);
| CREATE INDEX idx_audit_action ON audit_log(action);
| CREATE INDEX idx_audit_performed ON audit_log(performed_at);
|
| CREATE FUNCTION get_user_orders(uid integer) RETURNS bigint
|     LANGUAGE sql STABLE AS 'SELECT count(*) FROM orders WHERE user_id = uid';
|
| CREATE FUNCTION log_change() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
|
| CREATE VIEW active_users AS
|     SELECT id, username, email, status
|     FROM users
|     WHERE active = true;
|
| CREATE VIEW order_summary AS
|     SELECT o.id, u.username, p.name AS product, o.quantity, o.total
|     FROM orders o
|     JOIN users u ON u.id = o.user_id
|     JOIN products p ON p.id = o.product_id;
|
| CREATE TRIGGER orders_audit AFTER INSERT ON orders
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| CREATE TRIGGER users_audit AFTER UPDATE ON users
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON TABLE products IS 'Product catalog';
| COMMENT ON TABLE orders IS 'Customer orders';
| COMMENT ON VIEW active_users IS 'Only active users';
| COMMENT ON VIEW order_summary IS 'Joined order details';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_see_own ON users FOR SELECT USING (active = true);
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_1
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| audit_log
: string
| r
< row
: string
| categories
: string
| r
< row
: string
| order_summary
: string
| v
< row
: string
| orders
: string
| r
< row
: string
| products
: string
| r
< row
: string
| users
: string
| r
> query ^
: string
| test_2
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| audit_log
: string
| r
< row
: string
| categories
: string
| r
< row
: string
| order_summary
: string
| v
< row
: string
| orders
: string
| r
< row
: string
| products
: string
| r
< row
: string
| users
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_1
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| audit_log
: string
| table_name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| audit_log
: string
| action
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| audit_log
: string
| old_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| new_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| audit_log
: string
| performed_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| categories
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| categories
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| categories
: string
| parent_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| order_summary
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| order_summary
: string
| product
: string
| character varying
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| quantity
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| order_summary
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_1.order_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| orders
: string
| user_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| orders
: string
| product_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| orders
: string
| quantity
: string
| integer
: string
| NO
: null
: null
: null
: int64 4
< row
: string
| orders
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| products
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| products
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| products
: string
| category_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| products
: string
| price
: string
| numeric
: string
| NO
: string
| 0
: null
: null
: int64 4
< row
: string
| products
: string
| priority
: string
| USER-DEFINED
: string
| YES
: string
| 'medium'::test_1.priority_level
: null
: null
: int64 5
< row
: string
| products
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_1.user_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| username
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_1.user_status
: null
: null
: int64 4
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 5
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
> query ^
: string
| test_2
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| audit_log
: string
| table_name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| audit_log
: string
| action
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| audit_log
: string
| old_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| new_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| audit_log
: string
| performed_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| categories
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| categories
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| categories
: string
| parent_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| order_summary
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| order_summary
: string
| product
: string
| character varying
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| quantity
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| order_summary
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| id
: string
| integer
: string
| NO
: string
| nextval('order_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| orders
: string
| user_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| orders
: string
| product_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| orders
: string
| quantity
: string
| integer
: string
| NO
: null
: null
: null
: int64 4
< row
: string
| orders
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| products
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| products
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| products
: string
| category_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| products
: string
| price
: string
| numeric
: string
| NO
: string
| 0
: null
: null
: int64 4
< row
: string
| products
: string
| priority
: string
| USER-DEFINED
: string
| YES
: string
| 'medium'::priority_level
: null
: null
: int64 5
< row
: string
| products
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('user_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| username
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::user_status
: null
: null
: int64 4
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 5
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_1
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| CREATE UNIQUE INDEX audit_log_pkey ON test_1.audit_log USING btree (id)
< row
: string
| categories_pkey
: string
| CREATE UNIQUE INDEX categories_pkey ON test_1.categories USING btree (id)
< row
: string
| idx_audit_action
: string
| CREATE INDEX idx_audit_action ON test_1.audit_log USING btree (action)
< row
: string
| idx_audit_performed
: string
| CREATE INDEX idx_audit_performed ON test_1.audit_log USING btree (performed_at)
< row
: string
| idx_orders_created
: string
| CREATE INDEX idx_orders_created ON test_1.orders USING btree (created_at)
< row
: string
| idx_orders_user
: string
| CREATE INDEX idx_orders_user ON test_1.orders USING btree (user_id)
< row
: string
| idx_products_category
: string
| CREATE INDEX idx_products_category ON test_1.products USING btree (category_id)
< row
: string
| idx_products_price
: string
| CREATE INDEX idx_products_price ON test_1.products USING btree (price)
< row
: string
| idx_users_email
: string
| CREATE INDEX idx_users_email ON test_1.users USING btree (email)
< row
: string
| idx_users_status
: string
| CREATE INDEX idx_users_status ON test_1.users USING btree (status)
< row
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON test_1.orders USING btree (id)
< row
: string
| products_pkey
: string
| CREATE UNIQUE INDEX products_pkey ON test_1.products USING btree (id)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_1.users USING btree (id)
< row
: string
| users_username_unique
: string
| CREATE UNIQUE INDEX users_username_unique ON test_1.users USING btree (username)
> query ^
: string
| test_2
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| CREATE UNIQUE INDEX audit_log_pkey ON test_2.audit_log USING btree (id)
< row
: string
| categories_pkey
: string
| CREATE UNIQUE INDEX categories_pkey ON test_2.categories USING btree (id)
< row
: string
| idx_audit_action
: string
| CREATE INDEX idx_audit_action ON test_2.audit_log USING btree (action)
< row
: string
| idx_audit_performed
: string
| CREATE INDEX idx_audit_performed ON test_2.audit_log USING btree (performed_at)
< row
: string
| idx_orders_created
: string
| CREATE INDEX idx_orders_created ON test_2.orders USING btree (created_at)
< row
: string
| idx_orders_user
: string
| CREATE INDEX idx_orders_user ON test_2.orders USING btree (user_id)
< row
: string
| idx_products_category
: string
| CREATE INDEX idx_products_category ON test_2.products USING btree (category_id)
< row
: string
| idx_products_price
: string
| CREATE INDEX idx_products_price ON test_2.products USING btree (price)
< row
: string
| idx_users_email
: string
| CREATE INDEX idx_users_email ON test_2.users USING btree (email)
< row
: string
| idx_users_status
: string
| CREATE INDEX idx_users_status ON test_2.users USING btree (status)
< row
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON test_2.orders USING btree (id)
< row
: string
| products_pkey
: string
| CREATE UNIQUE INDEX products_pkey ON test_2.products USING btree (id)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_2.users USING btree (id)
< row
: string
| users_username_unique
: string
| CREATE UNIQUE INDEX users_username_unique ON test_2.users USING btree (username)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_1
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| p
: string
| audit_log
: string
| PRIMARY KEY (id)
< row
: string
| categories_parent_fk
: string
| f
: string
| categories
: string
| FOREIGN KEY (parent_id) REFERENCES test_1.categories(id)
< row
: string
| categories_pkey
: string
| p
: string
| categories
: string
| PRIMARY KEY (id)
< row
: string
| orders_pkey
: string
| p
: string
| orders
: string
| PRIMARY KEY (id)
< row
: string
| orders_product_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (product_id) REFERENCES test_1.products(id)
< row
: string
| orders_user_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (user_id) REFERENCES test_1.users(id)
< row
: string
| products_category_fk
: string
| f
: string
| products
: string
| FOREIGN KEY (category_id) REFERENCES test_1.categories(id)
< row
: string
| products_pkey
: string
| p
: string
| products
: string
| PRIMARY KEY (id)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
< row
: string
| users_username_unique
: string
| u
: string
| users
: string
| UNIQUE (username)
> query ^
: string
| test_2
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| p
: string
| audit_log
: string
| PRIMARY KEY (id)
< row
: string
| categories_parent_fk
: string
| f
: string
| categories
: string
| FOREIGN KEY (parent_id) REFERENCES categories(id)
< row
: string
| categories_pkey
: string
| p
: string
| categories
: string
| PRIMARY KEY (id)
< row
: string
| orders_pkey
: string
| p
: string
| orders
: string
| PRIMARY KEY (id)
< row
: string
| orders_product_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (product_id) REFERENCES products(id)
< row
: string
| orders_user_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (user_id) REFERENCES users(id)
< row
: string
| products_category_fk
: string
| f
: string
| products
: string
| FOREIGN KEY (category_id) REFERENCES categories(id)
< row
: string
| products_pkey
: string
| p
: string
| products
: string
| PRIMARY KEY (id)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
< row
: string
| users_username_unique
: string
| u
: string
| users
: string
| UNIQUE (username)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_1
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_orders
: string
| bigint
: string
| uid integer
: string
| SELECT count(*) FROM orders WHERE user_id = uid
: string
| sql
: string
| s
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| log_change
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_2
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_orders
: string
| bigint
: string
| uid integer
: string
| SELECT count(*) FROM orders WHERE user_id = uid
: string
| sql
: string
| s
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| log_change
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_1
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| orders_audit
: string
| orders
: string
| CREATE TRIGGER orders_audit AFTER INSERT ON test_1.orders FOR EACH ROW EXECUTE FUNCTION test_1.log_change()
< row
: string
| users_audit
: string
| users
: string
| CREATE TRIGGER users_audit AFTER UPDATE ON test_1.users FOR EACH ROW EXECUTE FUNCTION test_1.log_change()
> query ^
: string
| test_2
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| orders_audit
: string
| orders
: string
| CREATE TRIGGER orders_audit AFTER INSERT ON test_2.orders FOR EACH ROW EXECUTE FUNCTION log_change()
< row
: string
| users_audit
: string
| users
: string
| CREATE TRIGGER users_audit AFTER UPDATE ON test_2.users FOR EACH ROW EXECUTE FUNCTION log_change()
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_1
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| order_id_seq
: string
| bigint
: int64 1000
: int64 1
: int64 1
: int64 999999
: bool false
< row
: string
| user_id_seq
: string
| bigint
: int64 1
: int64 1
: int64 1
: int64 999999
: bool false
> query ^
: string
| test_2
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| order_id_seq
: string
| bigint
: int64 1000
: int64 1
: int64 1
: int64 999999
: bool false
< row
: string
| user_id_seq
: string
| bigint
: int64 1
: int64 1
: int64 1
: int64 999999
: bool false
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_1
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority_level
: string
| low,medium,high,critical
< row
: string
| user_status
: string
| active,inactive,banned
> query ^
: string
| test_2
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority_level
: string
| low,medium,high,critical
< row
: string
| user_status
: string
| active,inactive,banned
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_1
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     username,
|     email,
|     status
|    FROM test_1.users
|   WHERE (active = true);
: null
< row
: string
| order_summary
: string
|  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((test_1.orders o
|      JOIN test_1.users u ON ((u.id = o.user_id)))
|      JOIN test_1.products p ON ((p.id = o.product_id)));
: null
> query ^
: string
| test_2
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     username,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
: null
< row
: string
| order_summary
: string
|  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((orders o
|      JOIN users u ON ((u.id = o.user_id)))
|      JOIN products p ON ((p.id = o.product_id)));
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_1
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| orders
: string
| Customer orders
< row
: string
| TABLE
: string
| products
: string
| Product catalog
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
< row
: string
| VIEW
: string
| order_summary
: string
| Joined order details
> query ^
: string
| test_2
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| orders
: string
| Customer orders
< row
: string
| TABLE
: string
| products
: string
| Product catalog
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
< row
: string
| VIEW
: string
| order_summary
: string
| Joined order details
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_1
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_2
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_1
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_see_own
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> query ^
: string
| test_2
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_see_own
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_3" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_3"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_4" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_4"
< ok 0 -
> exec
| SET search_path TO "test_3", public;
|
| CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned');
| CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');
|
| CREATE SEQUENCE user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1 NO CYCLE;
| CREATE SEQUENCE order_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1000 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_id_seq'),
|     username varchar(100) NOT NULL,
|     email text,
|     status user_status DEFAULT 'active',
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT users_username_unique UNIQUE (username)
| );
|
| CREATE TABLE categories (
|     id integer NOT NULL,
|     name varchar(100) NOT NULL,
|     parent_id integer,
|     CONSTRAINT categories_pkey PRIMARY KEY (id),
|     CONSTRAINT categories_parent_fk FOREIGN KEY (parent_id) REFERENCES categories(id)
| );
|
| CREATE TABLE products (
|     id integer NOT NULL,
|     name varchar(200) NOT NULL,
|     category_id integer NOT NULL,
|     price numeric(10,2) NOT NULL DEFAULT 0,
|     priority priority_level DEFAULT 'medium',
|     in_stock boolean DEFAULT true,
|     CONSTRAINT products_pkey PRIMARY KEY (id),
|     CONSTRAINT products_category_fk FOREIGN KEY (category_id) REFERENCES categories(id)
| );
|
| CREATE TABLE orders (
|     id integer NOT NULL DEFAULT nextval('order_id_seq'),
|     user_id integer NOT NULL,
|     product_id integer NOT NULL,
|     quantity integer NOT NULL,
|     total numeric(12,2),
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT orders_pkey PRIMARY KEY (id),
|     CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users(id),
|     CONSTRAINT orders_product_fk FOREIGN KEY (product_id) REFERENCES products(id)
| );
|
| CREATE TABLE audit_log (
|     id integer NOT NULL,
|     table_name text NOT NULL,
|     action text NOT NULL,
|     old_data jsonb,
|     new_data jsonb,
|     performed_at timestamptz DEFAULT now(),
|     CONSTRAINT audit_log_pkey PRIMARY KEY (id)
| );
|
| CREATE INDEX idx_users_status ON users(status);
| CREATE INDEX idx_users_email ON users(email);
| CREATE INDEX idx_products_category ON products(category_id);
| CREATE INDEX idx_products_price ON products(price);
| CREATE INDEX idx_orders_user ON orders(user_id);
| CREATE INDEX idx_orders_created ON orders(created_at);
| CREATE INDEX idx_audit_action ON audit_log(action);
| CREATE INDEX idx_audit_performed ON audit_log(performed_at);
|
| CREATE FUNCTION get_user_orders(uid integer) RETURNS bigint
|     LANGUAGE sql STABLE AS 'SELECT count(*) FROM orders WHERE user_id = uid';
|
| CREATE FUNCTION log_change() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
|
| CREATE VIEW active_users AS
|     SELECT id, username, email, status
|     FROM users
|     WHERE active = true;
|
| CREATE VIEW order_summary AS
|     SELECT o.id, u.username, p.name AS product, o.quantity, o.total
|     FROM orders o
|     JOIN users u ON u.id = o.user_id
|     JOIN products p ON p.id = o.product_id;
|
| CREATE TRIGGER orders_audit AFTER INSERT ON orders
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| CREATE TRIGGER users_audit AFTER UPDATE ON users
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON TABLE products IS 'Product catalog';
| COMMENT ON TABLE orders IS 'Customer orders';
| COMMENT ON VIEW active_users IS 'Only active users';
| COMMENT ON VIEW order_summary IS 'Joined order details';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_see_own ON users FOR SELECT USING (active = true);
|
< ok 0 -
> exec
| SET search_path TO "test_3", public;
| ALTER TYPE "user_status" ADD VALUE 'suspended' AFTER 'banned';
| CREATE OR REPLACE FUNCTION "get_user_orders"("uid" integer)
| RETURNS bigint
| LANGUAGE sql
| STABLE
| AS $$SELECT count(*) FROM orders WHERE user_id = uid AND quantity > 0$$;
| ALTER TABLE products ADD COLUMN description text;
| ALTER TABLE users ADD COLUMN bio text;
| CREATE INDEX idx_products_name ON products (name);
| COMMENT ON VIEW "active_users" IS 'Only active users (v2)';
| COMMENT ON TABLE "users" IS 'Core user accounts (updated)'
< ok 0 -
> exec
| SET search_path TO "test_4", public;
|
| CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned', 'suspended');
| CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');
|
| CREATE SEQUENCE user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1 NO CYCLE;
| CREATE SEQUENCE order_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1000 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_id_seq'),
|     username varchar(100) NOT NULL,
|     email text,
|     status user_status DEFAULT 'active',
|     active boolean DEFAULT true,
|     bio text,
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT users_username_unique UNIQUE (username)
| );
|
| CREATE TABLE categories (
|     id integer NOT NULL,
|     name varchar(100) NOT NULL,
|     parent_id integer,
|     CONSTRAINT categories_pkey PRIMARY KEY (id),
|     CONSTRAINT categories_parent_fk FOREIGN KEY (parent_id) REFERENCES categories(id)
| );
|
| CREATE TABLE products (
|     id integer NOT NULL,
|     name varchar(200) NOT NULL,
|     category_id integer NOT NULL,
|     price numeric(10,2) NOT NULL DEFAULT 0,
|     priority priority_level DEFAULT 'medium',
|     in_stock boolean DEFAULT true,
|     description text,
|     CONSTRAINT products_pkey PRIMARY KEY (id),
|     CONSTRAINT products_category_fk FOREIGN KEY (category_id) REFERENCES categories(id)
| );
|
| CREATE TABLE orders (
|     id integer NOT NULL DEFAULT nextval('order_id_seq'),
|     user_id integer NOT NULL,
|     product_id integer NOT NULL,
|     quantity integer NOT NULL,
|     total numeric(12,2),
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT orders_pkey PRIMARY KEY (id),
|     CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users(id),
|     CONSTRAINT orders_product_fk FOREIGN KEY (product_id) REFERENCES products(id)
| );
|
| CREATE TABLE audit_log (
|     id integer NOT NULL,
|     table_name text NOT NULL,
|     action text NOT NULL,
|     old_data jsonb,
|     new_data jsonb,
|     performed_at timestamptz DEFAULT now(),
|     CONSTRAINT audit_log_pkey PRIMARY KEY (id)
| );
|
| CREATE INDEX idx_users_status ON users(status);
| CREATE INDEX idx_users_email ON users(email);
| CREATE INDEX idx_products_category ON products(category_id);
| CREATE INDEX idx_products_price ON products(price);
| CREATE INDEX idx_products_name ON products(name);
| CREATE INDEX idx_orders_user ON orders(user_id);
| CREATE INDEX idx_orders_created ON orders(created_at);
| CREATE INDEX idx_audit_action ON audit_log(action);
| CREATE INDEX idx_audit_performed ON audit_log(performed_at);
|
| CREATE FUNCTION get_user_orders(uid integer) RETURNS bigint
|     LANGUAGE sql STABLE AS 'SELECT count(*) FROM orders WHERE user_id = uid AND quantity > 0';
|
| CREATE FUNCTION log_change() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
|
| CREATE VIEW active_users AS
|     SELECT id, username, email, status
|     FROM users
|     WHERE active = true;
|
| CREATE VIEW order_summary AS
|     SELECT o.id, u.username, p.name AS product, o.quantity, o.total
|     FROM orders o
|     JOIN users u ON u.id = o.user_id
|     JOIN products p ON p.id = o.product_id;
|
| CREATE TRIGGER orders_audit AFTER INSERT ON orders
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| CREATE TRIGGER users_audit AFTER UPDATE ON users
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| COMMENT ON TABLE users IS 'Core user accounts (updated)';
| COMMENT ON TABLE products IS 'Product catalog';
| COMMENT ON TABLE orders IS 'Customer orders';
| COMMENT ON VIEW active_users IS 'Only active users (v2)';
| COMMENT ON VIEW order_summary IS 'Joined order details';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_see_own ON users FOR SELECT USING (active = true);
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_3
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| audit_log
: string
| r
< row
: string
| categories
: string
| r
< row
: string
| order_summary
: string
| v
< row
: string
| orders
: string
| r
< row
: string
| products
: string
| r
< row
: string
| users
: string
| r
> query ^
: string
| test_4
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| audit_log
: string
| r
< row
: string
| categories
: string
| r
< row
: string
| order_summary
: string
| v
< row
: string
| orders
: string
| r
< row
: string
| products
: string
| r
< row
: string
| users
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_3
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| audit_log
: string
| table_name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| audit_log
: string
| action
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| audit_log
: string
| old_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| new_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| audit_log
: string
| performed_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| categories
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| categories
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| categories
: string
| parent_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| order_summary
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| order_summary
: string
| product
: string
| character varying
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| quantity
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| order_summary
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_3.order_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| orders
: string
| user_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| orders
: string
| product_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| orders
: string
| quantity
: string
| integer
: string
| NO
: null
: null
: null
: int64 4
< row
: string
| orders
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| products
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| products
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| products
: string
| category_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| products
: string
| price
: string
| numeric
: string
| NO
: string
| 0
: null
: null
: int64 4
< row
: string
| products
: string
| priority
: string
| USER-DEFINED
: string
| YES
: string
| 'medium'::test_3.priority_level
: null
: null
: int64 5
< row
: string
| products
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| products
: string
| description
: string
| text
: string
| YES
: null
: null
: null
: int64 7
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_3.user_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| username
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_3.user_status
: null
: null
: int64 4
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 5
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| users
: string
| bio
: string
| text
: string
| YES
: null
: null
: null
: int64 7
> query ^
: string
| test_4
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| audit_log
: string
| table_name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| audit_log
: string
| action
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| audit_log
: string
| old_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| new_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| audit_log
: string
| performed_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| categories
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| categories
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| categories
: string
| parent_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| order_summary
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| order_summary
: string
| product
: string
| character varying
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| quantity
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| order_summary
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| id
: string
| integer
: string
| NO
: string
| nextval('order_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| orders
: string
| user_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| orders
: string
| product_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| orders
: string
| quantity
: string
| integer
: string
| NO
: null
: null
: null
: int64 4
< row
: string
| orders
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| products
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| products
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| products
: string
| category_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| products
: string
| price
: string
| numeric
: string
| NO
: string
| 0
: null
: null
: int64 4
< row
: string
| products
: string
| priority
: string
| USER-DEFINED
: string
| YES
: string
| 'medium'::priority_level
: null
: null
: int64 5
< row
: string
| products
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| products
: string
| description
: string
| text
: string
| YES
: null
: null
: null
: int64 7
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('user_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| username
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::user_status
: null
: null
: int64 4
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 5
< row
: string
| users
: string
| bio
: string
| text
: string
| YES
: null
: null
: null
: int64 6
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 7
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_3
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| CREATE UNIQUE INDEX audit_log_pkey ON test_3.audit_log USING btree (id)
< row
: string
| categories_pkey
: string
| CREATE UNIQUE INDEX categories_pkey ON test_3.categories USING btree (id)
< row
: string
| idx_audit_action
: string
| CREATE INDEX idx_audit_action ON test_3.audit_log USING btree (action)
< row
: string
| idx_audit_performed
: string
| CREATE INDEX idx_audit_performed ON test_3.audit_log USING btree (performed_at)
< row
: string
| idx_orders_created
: string
| CREATE INDEX idx_orders_created ON test_3.orders USING btree (created_at)
< row
: string
| idx_orders_user
: string
| CREATE INDEX idx_orders_user ON test_3.orders USING btree (user_id)
< row
: string
| idx_products_category
: string
| CREATE INDEX idx_products_category ON test_3.products USING btree (category_id)
< row
: string
| idx_products_name
: string
| CREATE INDEX idx_products_name ON test_3.products USING btree (name)
< row
: string
| idx_products_price
: string
| CREATE INDEX idx_products_price ON test_3.products USING btree (price)
< row
: string
| idx_users_email
: string
| CREATE INDEX idx_users_email ON test_3.users USING btree (email)
< row
: string
| idx_users_status
: string
| CREATE INDEX idx_users_status ON test_3.users USING btree (status)
< row
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON test_3.orders USING btree (id)
< row
: string
| products_pkey
: string
| CREATE UNIQUE INDEX products_pkey ON test_3.products USING btree (id)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_3.users USING btree (id)
< row
: string
| users_username_unique
: string
| CREATE UNIQUE INDEX users_username_unique ON test_3.users USING btree (username)
> query ^
: string
| test_4
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| CREATE UNIQUE INDEX audit_log_pkey ON test_4.audit_log USING btree (id)
< row
: string
| categories_pkey
: string
| CREATE UNIQUE INDEX categories_pkey ON test_4.categories USING btree (id)
< row
: string
| idx_audit_action
: string
| CREATE INDEX idx_audit_action ON test_4.audit_log USING btree (action)
< row
: string
| idx_audit_performed
: string
| CREATE INDEX idx_audit_performed ON test_4.audit_log USING btree (performed_at)
< row
: string
| idx_orders_created
: string
| CREATE INDEX idx_orders_created ON test_4.orders USING btree (created_at)
< row
: string
| idx_orders_user
: string
| CREATE INDEX idx_orders_user ON test_4.orders USING btree (user_id)
< row
: string
| idx_products_category
: string
| CREATE INDEX idx_products_category ON test_4.products USING btree (category_id)
< row
: string
| idx_products_name
: string
| CREATE INDEX idx_products_name ON test_4.products USING btree (name)
< row
: string
| idx_products_price
: string
| CREATE INDEX idx_products_price ON test_4.products USING btree (price)
< row
: string
| idx_users_email
: string
| CREATE INDEX idx_users_email ON test_4.users USING btree (email)
< row
: string
| idx_users_status
: string
| CREATE INDEX idx_users_status ON test_4.users USING btree (status)
< row
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON test_4.orders USING btree (id)
< row
: string
| products_pkey
: string
| CREATE UNIQUE INDEX products_pkey ON test_4.products USING btree (id)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_4.users USING btree (id)
< row
: string
| users_username_unique
: string
| CREATE UNIQUE INDEX users_username_unique ON test_4.users USING btree (username)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_3
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| p
: string
| audit_log
: string
| PRIMARY KEY (id)
< row
: string
| categories_parent_fk
: string
| f
: string
| categories
: string
| FOREIGN KEY (parent_id) REFERENCES test_3.categories(id)
< row
: string
| categories_pkey
: string
| p
: string
| categories
: string
| PRIMARY KEY (id)
< row
: string
| orders_pkey
: string
| p
: string
| orders
: string
| PRIMARY KEY (id)
< row
: string
| orders_product_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (product_id) REFERENCES test_3.products(id)
< row
: string
| orders_user_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (user_id) REFERENCES test_3.users(id)
< row
: string
| products_category_fk
: string
| f
: string
| products
: string
| FOREIGN KEY (category_id) REFERENCES test_3.categories(id)
< row
: string
| products_pkey
: string
| p
: string
| products
: string
| PRIMARY KEY (id)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
< row
: string
| users_username_unique
: string
| u
: string
| users
: string
| UNIQUE (username)
> query ^
: string
| test_4
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| p
: string
| audit_log
: string
| PRIMARY KEY (id)
< row
: string
| categories_parent_fk
: string
| f
: string
| categories
: string
| FOREIGN KEY (parent_id) REFERENCES categories(id)
< row
: string
| categories_pkey
: string
| p
: string
| categories
: string
| PRIMARY KEY (id)
< row
: string
| orders_pkey
: string
| p
: string
| orders
: string
| PRIMARY KEY (id)
< row
: string
| orders_product_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (product_id) REFERENCES products(id)
< row
: string
| orders_user_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (user_id) REFERENCES users(id)
< row
: string
| products_category_fk
: string
| f
: string
| products
: string
| FOREIGN KEY (category_id) REFERENCES categories(id)
< row
: string
| products_pkey
: string
| p
: string
| products
: string
| PRIMARY KEY (id)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
< row
: string
| users_username_unique
: string
| u
: string
| users
: string
| UNIQUE (username)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_3
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_orders
: string
| bigint
: string
| uid integer
: string
| SELECT count(*) FROM orders WHERE user_id = uid AND quantity > 0
: string
| sql
: string
| s
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| log_change
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_4
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_orders
: string
| bigint
: string
| uid integer
: string
| SELECT count(*) FROM orders WHERE user_id = uid AND quantity > 0
: string
| sql
: string
| s
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| log_change
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_3
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| orders_audit
: string
| orders
: string
| CREATE TRIGGER orders_audit AFTER INSERT ON test_3.orders FOR EACH ROW EXECUTE FUNCTION test_3.log_change()
< row
: string
| users_audit
: string
| users
: string
| CREATE TRIGGER users_audit AFTER UPDATE ON test_3.users FOR EACH ROW EXECUTE FUNCTION test_3.log_change()
> query ^
: string
| test_4
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| orders_audit
: string
| orders
: string
| CREATE TRIGGER orders_audit AFTER INSERT ON test_4.orders FOR EACH ROW EXECUTE FUNCTION log_change()
< row
: string
| users_audit
: string
| users
: string
| CREATE TRIGGER users_audit AFTER UPDATE ON test_4.users FOR EACH ROW EXECUTE FUNCTION log_change()
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_3
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| order_id_seq
: string
| bigint
: int64 1000
: int64 1
: int64 1
: int64 999999
: bool false
< row
: string
| user_id_seq
: string
| bigint
: int64 1
: int64 1
: int64 1
: int64 999999
: bool false
> query ^
: string
| test_4
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| order_id_seq
: string
| bigint
: int64 1000
: int64 1
: int64 1
: int64 999999
: bool false
< row
: string
| user_id_seq
: string
| bigint
: int64 1
: int64 1
: int64 1
: int64 999999
: bool false
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_3
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority_level
: string
| low,medium,high,critical
< row
: string
| user_status
: string
| active,inactive,banned,suspended
> query ^
: string
| test_4
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority_level
: string
| low,medium,high,critical
< row
: string
| user_status
: string
| active,inactive,banned,suspended
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_3
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     username,
|     email,
|     status
|    FROM test_3.users
|   WHERE (active = true);
: null
< row
: string
| order_summary
: string
|  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((test_3.orders o
|      JOIN test_3.users u ON ((u.id = o.user_id)))
|      JOIN test_3.products p ON ((p.id = o.product_id)));
: null
> query ^
: string
| test_4
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     username,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
: null
< row
: string
| order_summary
: string
|  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((orders o
|      JOIN users u ON ((u.id = o.user_id)))
|      JOIN products p ON ((p.id = o.product_id)));
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_3
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| orders
: string
| Customer orders
< row
: string
| TABLE
: string
| products
: string
| Product catalog
< row
: string
| TABLE
: string
| users
: string
| Core user accounts (updated)
< row
: string
| VIEW
: string
| active_users
: string
| Only active users (v2)
< row
: string
| VIEW
: string
| order_summary
: string
| Joined order details
> query ^
: string
| test_4
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| orders
: string
| Customer orders
< row
: string
| TABLE
: string
| products
: string
| Product catalog
< row
: string
| TABLE
: string
| users
: string
| Core user accounts (updated)
< row
: string
| VIEW
: string
| active_users
: string
| Only active users (v2)
< row
: string
| VIEW
: string
| order_summary
: string
| Joined order details
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_3
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_4
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_3
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_see_own
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> query ^
: string
| test_4
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_see_own
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> exec
| DROP SCHEMA IF EXISTS "test_4" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_3" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_5" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_5"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_6" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_6"
< ok 0 -
> exec
| SET search_path TO "test_5", public;
|
| CREATE TABLE users (
|     id integer NOT NULL,
|     username varchar(100) NOT NULL,
|     active boolean DEFAULT true,
|     CONSTRAINT users_pkey PRIMARY KEY (id)
| );
| CREATE TABLE products (
|     id integer NOT NULL,
|     name varchar(200) NOT NULL,
|     CONSTRAINT products_pkey PRIMARY KEY (id)
| );
| CREATE INDEX idx_users_active ON users(active);
| CREATE INDEX idx_products_name ON products(name);
| CREATE FUNCTION log_change() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
| CREATE TRIGGER products_audit AFTER INSERT ON products
|     FOR EACH ROW EXECUTE FUNCTION log_change();
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_see_own ON users FOR SELECT USING (active = true);
|
< ok 0 -
> exec
| SET search_path TO "test_5", public;
| DROP TABLE "products" CASCADE;
| DROP TABLE "users" CASCADE;
| DROP FUNCTION "log_change"()
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_5
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> query ^
: string
| test_6
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_5
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
> query ^
: string
| test_6
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_5
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query ^
: string
| test_6
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_5
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query ^
: string
| test_6
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_5
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query ^
: string
| test_6
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_5
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_6
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_5
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query ^
: string
| test_6
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_5
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_6
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_5
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_6
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_5
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_6
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_5
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_6
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_5
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_6
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_6" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_5" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_7" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_7"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_8" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_8"
< ok 0 -
> exec
| SET search_path TO "test_7", public;
| CREATE TYPE "priority_level" AS ENUM ('low', 'medium', 'high', 'critical');
| CREATE TYPE "user_status" AS ENUM ('active', 'inactive', 'banned');
| CREATE SEQUENCE "order_id_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 1000
|     CACHE 1
|     NO CYCLE;
| CREATE SEQUENCE "user_id_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 1
|     CACHE 1
|     NO CYCLE;
| CREATE FUNCTION "log_change"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN RETURN NEW; END;$$;
| CREATE TABLE "audit_log" (
|     "id" integer NOT NULL,
|     "table_name" text NOT NULL,
|     "action" text NOT NULL,
|     "old_data" jsonb,
|     "new_data" jsonb,
|     "performed_at" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "audit_log_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "categories" (
|     "id" integer NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "parent_id" integer,
|     CONSTRAINT "categories_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "orders" (
|     "id" integer NOT NULL DEFAULT nextval('order_id_seq'::regclass),
|     "user_id" integer NOT NULL,
|     "product_id" integer NOT NULL,
|     "quantity" integer NOT NULL,
|     "total" numeric(12,2),
|     "created_at" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "orders_pkey" PRIMARY KEY ("id")
| );
| CREATE FUNCTION "get_user_orders"("uid" integer)
| RETURNS bigint
| LANGUAGE sql
| STABLE
| AS $$SELECT count(*) FROM orders WHERE user_id = uid$$;
| CREATE TABLE "products" (
|     "id" integer NOT NULL,
|     "name" character varying(200) NOT NULL,
|     "category_id" integer NOT NULL,
|     "price" numeric(10,2) NOT NULL DEFAULT 0,
|     "priority" priority_level DEFAULT 'medium'::priority_level,
|     "in_stock" boolean DEFAULT true,
|     CONSTRAINT "products_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "users" (
|     "id" integer NOT NULL DEFAULT nextval('user_id_seq'::regclass),
|     "username" character varying(100) NOT NULL,
|     "email" text,
|     "status" user_status DEFAULT 'active'::user_status,
|     "active" boolean DEFAULT true,
|     "created_at" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
|     CONSTRAINT "users_username_unique" UNIQUE ("username")
| );
| ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
| CREATE VIEW "active_users" AS  SELECT id,
|     username,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
| CREATE VIEW "order_summary" AS  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((orders o
|      JOIN users u ON ((u.id = o.user_id)))
|      JOIN products p ON ((p.id = o.product_id)));
| CREATE INDEX idx_audit_action ON audit_log (action);
| CREATE INDEX idx_audit_performed ON audit_log (performed_at);
| CREATE INDEX idx_orders_created ON orders (created_at);
| CREATE INDEX idx_orders_user ON orders (user_id);
| CREATE INDEX idx_products_category ON products (category_id);
| CREATE INDEX idx_products_price ON products (price);
| CREATE INDEX idx_users_email ON users (email);
| CREATE INDEX idx_users_status ON users (status);
| CREATE TRIGGER "orders_audit" AFTER INSERT ON "orders" FOR EACH ROW EXECUTE FUNCTION "log_change"();
| CREATE TRIGGER "users_audit" AFTER UPDATE ON "users" FOR EACH ROW EXECUTE FUNCTION "log_change"();
| CREATE POLICY "users_see_own" ON "users" FOR SELECT TO PUBLIC USING ((active = true));
| COMMENT ON VIEW "active_users" IS 'Only active users';
| COMMENT ON VIEW "order_summary" IS 'Joined order details';
| COMMENT ON TABLE "orders" IS 'Customer orders';
| COMMENT ON TABLE "products" IS 'Product catalog';
| COMMENT ON TABLE "users" IS 'Core user accounts';
| ALTER TABLE categories ADD CONSTRAINT categories_parent_fk FOREIGN KEY (parent_id) REFERENCES categories (id);
| ALTER TABLE orders ADD CONSTRAINT orders_product_fk FOREIGN KEY (product_id) REFERENCES products (id);
| ALTER TABLE orders ADD CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users (id);
| ALTER TABLE products ADD CONSTRAINT products_category_fk FOREIGN KEY (category_id) REFERENCES categories (id)
< ok 0 -
> exec
| SET search_path TO "test_8", public;
|
| CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned');
| CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');
|
| CREATE SEQUENCE user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1 NO CYCLE;
| CREATE SEQUENCE order_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1000 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_id_seq'),
|     username varchar(100) NOT NULL,
|     email text,
|     status user_status DEFAULT 'active',
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT users_username_unique UNIQUE (username)
| );
|
| CREATE TABLE categories (
|     id integer NOT NULL,
|     name varchar(100) NOT NULL,
|     parent_id integer,
|     CONSTRAINT categories_pkey PRIMARY KEY (id),
|     CONSTRAINT categories_parent_fk FOREIGN KEY (parent_id) REFERENCES categories(id)
| );
|
| CREATE TABLE products (
|     id integer NOT NULL,
|     name varchar(200) NOT NULL,
|     category_id integer NOT NULL,
|     price numeric(10,2) NOT NULL DEFAULT 0,
|     priority priority_level DEFAULT 'medium',
|     in_stock boolean DEFAULT true,
|     CONSTRAINT products_pkey PRIMARY KEY (id),
|     CONSTRAINT products_category_fk FOREIGN KEY (category_id) REFERENCES categories(id)
| );
|
| CREATE TABLE orders (
|     id integer NOT NULL DEFAULT nextval('order_id_seq'),
|     user_id integer NOT NULL,
|     product_id integer NOT NULL,
|     quantity integer NOT NULL,
|     total numeric(12,2),
|     created_at timestamptz DEFAULT now(),
|     CONSTRAINT orders_pkey PRIMARY KEY (id),
|     CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users(id),
|     CONSTRAINT orders_product_fk FOREIGN KEY (product_id) REFERENCES products(id)
| );
|
| CREATE TABLE audit_log (
|     id integer NOT NULL,
|     table_name text NOT NULL,
|     action text NOT NULL,
|     old_data jsonb,
|     new_data jsonb,
|     performed_at timestamptz DEFAULT now(),
|     CONSTRAINT audit_log_pkey PRIMARY KEY (id)
| );
|
| CREATE INDEX idx_users_status ON users(status);
| CREATE INDEX idx_users_email ON users(email);
| CREATE INDEX idx_products_category ON products(category_id);
| CREATE INDEX idx_products_price ON products(price);
| CREATE INDEX idx_orders_user ON orders(user_id);
| CREATE INDEX idx_orders_created ON orders(created_at);
| CREATE INDEX idx_audit_action ON audit_log(action);
| CREATE INDEX idx_audit_performed ON audit_log(performed_at);
|
| CREATE FUNCTION get_user_orders(uid integer) RETURNS bigint
|     LANGUAGE sql STABLE AS 'SELECT count(*) FROM orders WHERE user_id = uid';
|
| CREATE FUNCTION log_change() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
|
| CREATE VIEW active_users AS
|     SELECT id, username, email, status
|     FROM users
|     WHERE active = true;
|
| CREATE VIEW order_summary AS
|     SELECT o.id, u.username, p.name AS product, o.quantity, o.total
|     FROM orders o
|     JOIN users u ON u.id = o.user_id
|     JOIN products p ON p.id = o.product_id;
|
| CREATE TRIGGER orders_audit AFTER INSERT ON orders
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| CREATE TRIGGER users_audit AFTER UPDATE ON users
|     FOR EACH ROW EXECUTE FUNCTION log_change();
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON TABLE products IS 'Product catalog';
| COMMENT ON TABLE orders IS 'Customer orders';
| COMMENT ON VIEW active_users IS 'Only active users';
| COMMENT ON VIEW order_summary IS 'Joined order details';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_see_own ON users FOR SELECT USING (active = true);
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_7
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| audit_log
: string
| r
< row
: string
| categories
: string
| r
< row
: string
| order_summary
: string
| v
< row
: string
| orders
: string
| r
< row
: string
| products
: string
| r
< row
: string
| users
: string
| r
> query ^
: string
| test_8
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| audit_log
: string
| r
< row
: string
| categories
: string
| r
< row
: string
| order_summary
: string
| v
< row
: string
| orders
: string
| r
< row
: string
| products
: string
| r
< row
: string
| users
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_7
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| audit_log
: string
| table_name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| audit_log
: string
| action
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| audit_log
: string
| old_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| new_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| audit_log
: string
| performed_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| categories
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| categories
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| categories
: string
| parent_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| order_summary
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| order_summary
: string
| product
: string
| character varying
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| quantity
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| order_summary
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_7.order_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| orders
: string
| user_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| orders
: string
| product_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| orders
: string
| quantity
: string
| integer
: string
| NO
: null
: null
: null
: int64 4
< row
: string
| orders
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| products
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| products
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| products
: string
| category_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| products
: string
| price
: string
| numeric
: string
| NO
: string
| 0
: null
: null
: int64 4
< row
: string
| products
: string
| priority
: string
| USER-DEFINED
: string
| YES
: string
| 'medium'::test_7.priority_level
: null
: null
: int64 5
< row
: string
| products
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_7.user_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| username
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_7.user_status
: null
: null
: int64 4
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 5
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
> query ^
: string
| test_8
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| audit_log
: string
| table_name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| audit_log
: string
| action
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| audit_log
: string
| old_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| audit_log
: string
| new_data
: string
| jsonb
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| audit_log
: string
| performed_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| categories
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| categories
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| categories
: string
| parent_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| order_summary
: string
| username
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| order_summary
: string
| product
: string
| character varying
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| order_summary
: string
| quantity
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| order_summary
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| id
: string
| integer
: string
| NO
: string
| nextval('order_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| orders
: string
| user_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| orders
: string
| product_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| orders
: string
| quantity
: string
| integer
: string
| NO
: null
: null
: null
: int64 4
< row
: string
| orders
: string
| total
: string
| numeric
: string
| YES
: null
: null
: null
: int64 5
< row
: string
| orders
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
< row
: string
| products
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| products
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| products
: string
| category_id
: string
| integer
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| products
: string
| price
: string
| numeric
: string
| NO
: string
| 0
: null
: null
: int64 4
< row
: string
| products
: string
| priority
: string
| USER-DEFINED
: string
| YES
: string
| 'medium'::priority_level
: null
: null
: int64 5
< row
: string
| products
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('user_id_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| username
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::user_status
: null
: null
: int64 4
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 5
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 6
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_7
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| CREATE UNIQUE INDEX audit_log_pkey ON test_7.audit_log USING btree (id)
< row
: string
| categories_pkey
: string
| CREATE UNIQUE INDEX categories_pkey ON test_7.categories USING btree (id)
< row
: string
| idx_audit_action
: string
| CREATE INDEX idx_audit_action ON test_7.audit_log USING btree (action)
< row
: string
| idx_audit_performed
: string
| CREATE INDEX idx_audit_performed ON test_7.audit_log USING btree (performed_at)
< row
: string
| idx_orders_created
: string
| CREATE INDEX idx_orders_created ON test_7.orders USING btree (created_at)
< row
: string
| idx_orders_user
: string
| CREATE INDEX idx_orders_user ON test_7.orders USING btree (user_id)
< row
: string
| idx_products_category
: string
| CREATE INDEX idx_products_category ON test_7.products USING btree (category_id)
< row
: string
| idx_products_price
: string
| CREATE INDEX idx_products_price ON test_7.products USING btree (price)
< row
: string
| idx_users_email
: string
| CREATE INDEX idx_users_email ON test_7.users USING btree (email)
< row
: string
| idx_users_status
: string
| CREATE INDEX idx_users_status ON test_7.users USING btree (status)
< row
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON test_7.orders USING btree (id)
< row
: string
| products_pkey
: string
| CREATE UNIQUE INDEX products_pkey ON test_7.products USING btree (id)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_7.users USING btree (id)
< row
: string
| users_username_unique
: string
| CREATE UNIQUE INDEX users_username_unique ON test_7.users USING btree (username)
> query ^
: string
| test_8
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| CREATE UNIQUE INDEX audit_log_pkey ON test_8.audit_log USING btree (id)
< row
: string
| categories_pkey
: string
| CREATE UNIQUE INDEX categories_pkey ON test_8.categories USING btree (id)
< row
: string
| idx_audit_action
: string
| CREATE INDEX idx_audit_action ON test_8.audit_log USING btree (action)
< row
: string
| idx_audit_performed
: string
| CREATE INDEX idx_audit_performed ON test_8.audit_log USING btree (performed_at)
< row
: string
| idx_orders_created
: string
| CREATE INDEX idx_orders_created ON test_8.orders USING btree (created_at)
< row
: string
| idx_orders_user
: string
| CREATE INDEX idx_orders_user ON test_8.orders USING btree (user_id)
< row
: string
| idx_products_category
: string
| CREATE INDEX idx_products_category ON test_8.products USING btree (category_id)
< row
: string
| idx_products_price
: string
| CREATE INDEX idx_products_price ON test_8.products USING btree (price)
< row
: string
| idx_users_email
: string
| CREATE INDEX idx_users_email ON test_8.users USING btree (email)
< row
: string
| idx_users_status
: string
| CREATE INDEX idx_users_status ON test_8.users USING btree (status)
< row
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON test_8.orders USING btree (id)
< row
: string
| products_pkey
: string
| CREATE UNIQUE INDEX products_pkey ON test_8.products USING btree (id)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_8.users USING btree (id)
< row
: string
| users_username_unique
: string
| CREATE UNIQUE INDEX users_username_unique ON test_8.users USING btree (username)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_7
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| p
: string
| audit_log
: string
| PRIMARY KEY (id)
< row
: string
| categories_parent_fk
: string
| f
: string
| categories
: string
| FOREIGN KEY (parent_id) REFERENCES test_7.categories(id)
< row
: string
| categories_pkey
: string
| p
: string
| categories
: string
| PRIMARY KEY (id)
< row
: string
| orders_pkey
: string
| p
: string
| orders
: string
| PRIMARY KEY (id)
< row
: string
| orders_product_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (product_id) REFERENCES test_7.products(id)
< row
: string
| orders_user_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (user_id) REFERENCES test_7.users(id)
< row
: string
| products_category_fk
: string
| f
: string
| products
: string
| FOREIGN KEY (category_id) REFERENCES test_7.categories(id)
< row
: string
| products_pkey
: string
| p
: string
| products
: string
| PRIMARY KEY (id)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
< row
: string
| users_username_unique
: string
| u
: string
| users
: string
| UNIQUE (username)
> query ^
: string
| test_8
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| audit_log_pkey
: string
| p
: string
| audit_log
: string
| PRIMARY KEY (id)
< row
: string
| categories_parent_fk
: string
| f
: string
| categories
: string
| FOREIGN KEY (parent_id) REFERENCES categories(id)
< row
: string
| categories_pkey
: string
| p
: string
| categories
: string
| PRIMARY KEY (id)
< row
: string
| orders_pkey
: string
| p
: string
| orders
: string
| PRIMARY KEY (id)
< row
: string
| orders_product_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (product_id) REFERENCES products(id)
< row
: string
| orders_user_fk
: string
| f
: string
| orders
: string
| FOREIGN KEY (user_id) REFERENCES users(id)
< row
: string
| products_category_fk
: string
| f
: string
| products
: string
| FOREIGN KEY (category_id) REFERENCES categories(id)
< row
: string
| products_pkey
: string
| p
: string
| products
: string
| PRIMARY KEY (id)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
< row
: string
| users_username_unique
: string
| u
: string
| users
: string
| UNIQUE (username)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_7
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_orders
: string
| bigint
: string
| uid integer
: string
| SELECT count(*) FROM orders WHERE user_id = uid
: string
| sql
: string
| s
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| log_change
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_8
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_orders
: string
| bigint
: string
| uid integer
: string
| SELECT count(*) FROM orders WHERE user_id = uid
: string
| sql
: string
| s
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| log_change
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_7
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| orders_audit
: string
| orders
: string
| CREATE TRIGGER orders_audit AFTER INSERT ON test_7.orders FOR EACH ROW EXECUTE FUNCTION test_7.log_change()
< row
: string
| users_audit
: string
| users
: string
| CREATE TRIGGER users_audit AFTER UPDATE ON test_7.users FOR EACH ROW EXECUTE FUNCTION test_7.log_change()
> query ^
: string
| test_8
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| orders_audit
: string
| orders
: string
| CREATE TRIGGER orders_audit AFTER INSERT ON test_8.orders FOR EACH ROW EXECUTE FUNCTION log_change()
< row
: string
| users_audit
: string
| users
: string
| CREATE TRIGGER users_audit AFTER UPDATE ON test_8.users FOR EACH ROW EXECUTE FUNCTION log_change()
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_7
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| order_id_seq
: string
| bigint
: int64 1000
: int64 1
: int64 1
: int64 999999
: bool false
< row
: string
| user_id_seq
: string
| bigint
: int64 1
: int64 1
: int64 1
: int64 999999
: bool false
> query ^
: string
| test_8
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| order_id_seq
: string
| bigint
: int64 1000
: int64 1
: int64 1
: int64 999999
: bool false
< row
: string
| user_id_seq
: string
| bigint
: int64 1
: int64 1
: int64 1
: int64 999999
: bool false
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_7
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority_level
: string
| low,medium,high,critical
< row
: string
| user_status
: string
| active,inactive,banned
> query ^
: string
| test_8
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority_level
: string
| low,medium,high,critical
< row
: string
| user_status
: string
| active,inactive,banned
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_7
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     username,
|     email,
|     status
|    FROM test_7.users
|   WHERE (active = true);
: null
< row
: string
| order_summary
: string
|  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((test_7.orders o
|      JOIN test_7.users u ON ((u.id = o.user_id)))
|      JOIN test_7.products p ON ((p.id = o.product_id)));
: null
> query ^
: string
| test_8
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     username,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
: null
< row
: string
| order_summary
: string
|  SELECT o.id,
|     u.username,
|     p.name AS product,
|     o.quantity,
|     o.total
|    FROM ((orders o
|      JOIN users u ON ((u.id = o.user_id)))
|      JOIN products p ON ((p.id = o.product_id)));
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_7
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| orders
: string
| Customer orders
< row
: string
| TABLE
: string
| products
: string
| Product catalog
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
< row
: string
| VIEW
: string
| order_summary
: string
| Joined order details
> query ^
: string
| test_8
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| orders
: string
| Customer orders
< row
: string
| TABLE
: string
| products
: string
| Product catalog
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
< row
: string
| VIEW
: string
| order_summary
: string
| Joined order details
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_7
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_8
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_7
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_see_own
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> query ^
: string
| test_8
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_see_own
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> exec
| DROP SCHEMA IF EXISTS "test_8" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_7" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_9"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_10"
< ok 0 -
> exec
| SET search_path TO "test_9", public;
| CREATE FUNCTION "stamp_fn"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN RETURN NEW; END;$$;
| CREATE TABLE "items" (
|     "id" integer NOT NULL,
|     "name" text NOT NULL,
|     "in_stock" boolean DEFAULT true,
|     CONSTRAINT "items_pkey" PRIMARY KEY ("id")
| );
| CREATE VIEW "active_items" AS  SELECT id,
|     name
|    FROM items
|   WHERE (in_stock = true);
| CREATE INDEX idx_items_name ON items (name);
| CREATE TRIGGER "items_stamp" AFTER INSERT ON "items" FOR EACH ROW EXECUTE FUNCTION "stamp_fn"()
< ok 0 -
> exec
| SET search_path TO "test_10", public;
| CREATE FUNCTION "stamp_fn"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN RETURN NEW; END;$$;
| CREATE TABLE "items" (
|     "id" integer NOT NULL,
|     "name" text NOT NULL,
|     "in_stock" boolean DEFAULT true,
|     CONSTRAINT "items_pkey" PRIMARY KEY ("id")
| );
| CREATE VIEW "active_items" AS  SELECT id,
|     name
|    FROM items
|   WHERE (in_stock = true);
| CREATE INDEX idx_items_name ON items (name);
| CREATE TRIGGER "items_stamp" AFTER INSERT ON "items" FOR EACH ROW EXECUTE FUNCTION "stamp_fn"()
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_9
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_items
: string
| v
< row
: string
| items
: string
| r
> query ^
: string
| test_10
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_items
: string
| v
< row
: string
| items
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_9
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_items
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_items
: string
| name
: string
| text
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| items
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| items
: string
| name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| items
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 3
> query ^
: string
| test_10
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_items
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_items
: string
| name
: string
| text
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| items
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| items
: string
| name
: string
| text
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| items
: string
| in_stock
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 3
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_9
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_items_name
: string
| CREATE INDEX idx_items_name ON test_9.items USING btree (name)
< row
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON test_9.items USING btree (id)
> query ^
: string
| test_10
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_items_name
: string
| CREATE INDEX idx_items_name ON test_10.items USING btree (name)
< row
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON test_10.items USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_9
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| items_pkey
: string
| p
: string
| items
: string
| PRIMARY KEY (id)
> query ^
: string
| test_10
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| items_pkey
: string
| p
: string
| items
: string
| PRIMARY KEY (id)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_9
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| stamp_fn
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_10
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| stamp_fn
: string
| trigger
: string
|
: string
| BEGIN RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_9
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| items_stamp
: string
| items
: string
| CREATE TRIGGER items_stamp AFTER INSERT ON test_9.items FOR EACH ROW EXECUTE FUNCTION test_9.stamp_fn()
> query ^
: string
| test_10
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| items_stamp
: string
| items
: string
| CREATE TRIGGER items_stamp AFTER INSERT ON test_10.items FOR EACH ROW EXECUTE FUNCTION stamp_fn()
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_9
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query ^
: string
| test_10
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_9
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_10
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_9
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_items
: string
|  SELECT id,
|     name
|    FROM test_9.items
|   WHERE (in_stock = true);
: null
> query ^
: string
| test_10
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_items
: string
|  SELECT id,
|     name
|    FROM items
|   WHERE (in_stock = true);
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_9
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_10
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_9
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_10
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_9
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_10
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_11"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_12"
< ok 0 -
> exec
| SET search_path TO "test_11", public;
| CREATE TABLE "member" (
|     "id" integer NOT NULL,
|     "team_id" integer,
|     CONSTRAINT "member_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "team" (
|     "id" integer NOT NULL,
|     "lead_member_id" integer,
|     CONSTRAINT "team_pkey" PRIMARY KEY ("id")
| );
| ALTER TABLE member ADD CONSTRAINT member_team_fk FOREIGN KEY (team_id) REFERENCES team (id);
| ALTER TABLE team ADD CONSTRAINT team_lead_fk FOREIGN KEY (lead_member_id) REFERENCES member (id)
< ok 0 -
> exec
| SET search_path TO "test_12", public;
| CREATE TABLE "member" (
|     "id" integer NOT NULL,
|     "team_id" integer,
|     CONSTRAINT "member_pkey" PRIMARY KEY ("id")
| );
| CREATE TABLE "team" (
|     "id" integer NOT NULL,
|     "lead_member_id" integer,
|     CONSTRAINT "team_pkey" PRIMARY KEY ("id")
| );
| ALTER TABLE member ADD CONSTRAINT member_team_fk FOREIGN KEY (team_id) REFERENCES team (id);
| ALTER TABLE team ADD CONSTRAINT team_lead_fk FOREIGN KEY (lead_member_id) REFERENCES member (id)
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_11
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| member
: string
| r
< row
: string
| team
: string
| r
> query ^
: string
| test_12
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| member
: string
| r
< row
: string
| team
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_11
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| member
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| member
: string
| team_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| team
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| team
: string
| lead_member_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
> query ^
: string
| test_12
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| member
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| member
: string
| team_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| team
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| team
: string
| lead_member_id
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_11
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| member_pkey
: string
| CREATE UNIQUE INDEX member_pkey ON test_11.member USING btree (id)
< row
: string
| team_pkey
: string
| CREATE UNIQUE INDEX team_pkey ON test_11.team USING btree (id)
> query ^
: string
| test_12
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| member_pkey
: string
| CREATE UNIQUE INDEX member_pkey ON test_12.member USING btree (id)
< row
: string
| team_pkey
: string
| CREATE UNIQUE INDEX team_pkey ON test_12.team USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_11
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| member_pkey
: string
| p
: string
| member
: string
| PRIMARY KEY (id)
< row
: string
| member_team_fk
: string
| f
: string
| member
: string
| FOREIGN KEY (team_id) REFERENCES test_11.team(id)
< row
: string
| team_lead_fk
: string
| f
: string
| team
: string
| FOREIGN KEY (lead_member_id) REFERENCES test_11.member(id)
< row
: string
| team_pkey
: string
| p
: string
| team
: string
| PRIMARY KEY (id)
> query ^
: string
| test_12
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| member_pkey
: string
| p
: string
| member
: string
| PRIMARY KEY (id)
< row
: string
| member_team_fk
: string
| f
: string
| member
: string
| FOREIGN KEY (team_id) REFERENCES team(id)
< row
: string
| team_lead_fk
: string
| f
: string
| team
: string
| FOREIGN KEY (lead_member_id) REFERENCES member(id)
< row
: string
| team_pkey
: string
| p
: string
| team
: string
| PRIMARY KEY (id)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_11
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query ^
: string
| test_12
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_11
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_12
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_11
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query ^
: string
| test_12
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_11
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_12
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_11
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_12
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_11
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_12
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_11
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_12
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_11
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_12
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_13"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_14"
< ok 0 -
> exec
| SET search_path TO "test_13", public;
| CREATE FUNCTION "is_valid"("val" integer)
| RETURNS boolean
| LANGUAGE sql
| IMMUTABLE
| AS $$SELECT val > 0$$;
| CREATE FUNCTION "on_insert"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN
|         IF NOT is_valid(NEW.qty) THEN RAISE EXCEPTION 'invalid qty'; END IF;
|         RETURN NEW;
|     END;$$;
| CREATE TABLE "items" (
|     "id" integer NOT NULL,
|     "qty" integer,
|     "name" text,
|     CONSTRAINT "items_pkey" PRIMARY KEY ("id"),
|     CONSTRAINT "items_qty_check" CHECK (is_valid(qty))
| );
| CREATE VIEW "valid_items" AS  SELECT id,
|     name
|    FROM items
|   WHERE is_valid(qty);
| CREATE TRIGGER "items_validate" BEFORE INSERT ON "items" FOR EACH ROW EXECUTE FUNCTION "on_insert"()
< ok 0 -
> exec
| SET search_path TO "test_14", public;
| CREATE FUNCTION "is_valid"("val" integer)
| RETURNS boolean
| LANGUAGE sql
| IMMUTABLE
| AS $$SELECT val > 0$$;
| CREATE FUNCTION "on_insert"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN
|         IF NOT is_valid(NEW.qty) THEN RAISE EXCEPTION 'invalid qty'; END IF;
|         RETURN NEW;
|     END;$$;
| CREATE TABLE "items" (
|     "id" integer NOT NULL,
|     "qty" integer,
|     "name" text,
|     CONSTRAINT "items_pkey" PRIMARY KEY ("id"),
|     CONSTRAINT "items_qty_check" CHECK (is_valid(qty))
| );
| CREATE VIEW "valid_items" AS  SELECT id,
|     name
|    FROM items
|   WHERE is_valid(qty);
| CREATE TRIGGER "items_validate" BEFORE INSERT ON "items" FOR EACH ROW EXECUTE FUNCTION "on_insert"()
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_13
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| items
: string
| r
< row
: string
| valid_items
: string
| v
> query ^
: string
| test_14
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| items
: string
| r
< row
: string
| valid_items
: string
| v
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_13
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| items
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| items
: string
| qty
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| items
: string
| name
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| valid_items
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| valid_items
: string
| name
: string
| text
: string
| YES
: null
: null
: null
: int64 2
> query ^
: string
| test_14
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| items
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| items
: string
| qty
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| items
: string
| name
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| valid_items
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| valid_items
: string
| name
: string
| text
: string
| YES
: null
: null
: null
: int64 2
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_13
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON test_13.items USING btree (id)
> query ^
: string
| test_14
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON test_14.items USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_13
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| items_pkey
: string
| p
: string
| items
: string
| PRIMARY KEY (id)
< row
: string
| items_qty_check
: string
| c
: string
| items
: string
| CHECK (test_13.is_valid(qty))
> query ^
: string
| test_14
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| items_pkey
: string
| p
: string
| items
: string
| PRIMARY KEY (id)
< row
: string
| items_qty_check
: string
| c
: string
| items
: string
| CHECK (is_valid(qty))
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_13
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| is_valid
: string
| boolean
: string
| val integer
: string
| SELECT val > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| on_insert
: string
| trigger
: string
|
: string
| BEGIN
|         IF NOT is_valid(NEW.qty) THEN RAISE EXCEPTION 'invalid qty'; END IF;
|         RETURN NEW;
|     END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_14
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| is_valid
: string
| boolean
: string
| val integer
: string
| SELECT val > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
< row
: string
| on_insert
: string
| trigger
: string
|
: string
| BEGIN
|         IF NOT is_valid(NEW.qty) THEN RAISE EXCEPTION 'invalid qty'; END IF;
|         RETURN NEW;
|     END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_13
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| items_validate
: string
| items
: string
| CREATE TRIGGER items_validate BEFORE INSERT ON test_13.items FOR EACH ROW EXECUTE FUNCTION test_13.on_insert()
> query ^
: string
| test_14
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| items_validate
: string
| items
: string
| CREATE TRIGGER items_validate BEFORE INSERT ON test_14.items FOR EACH ROW EXECUTE FUNCTION on_insert()
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_13
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query ^
: string
| test_14
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_13
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_14
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_13
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| valid_items
: string
|  SELECT id,
|     name
|    FROM test_13.items
|   WHERE test_13.is_valid(qty);
: null
> query ^
: string
| test_14
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| valid_items
: string
|  SELECT id,
|     name
|    FROM items
|   WHERE is_valid(qty);
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_13
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_14
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_13
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_14
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_13
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_14
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -

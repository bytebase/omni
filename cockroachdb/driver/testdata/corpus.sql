-- The schema TestSyncDatabase syncs and the parity harness compares with
-- Bytebase's driver. It covers each object kind CockroachDB's catalogs report,
-- in the forms that change how the sync reads them.

CREATE SCHEMA app;
COMMENT ON SCHEMA app IS 'the application';
CREATE SCHEMA "Mixed Case";
CREATE SCHEMA audit;

CREATE TYPE app.mood AS ENUM ('happy', 'sad', 'it''s complicated');
COMMENT ON TYPE app.mood IS 'how it feels';
CREATE TYPE app.address AS (street STRING, city VARCHAR(50), zip INT8);

CREATE SEQUENCE app.counter AS INT8 INCREMENT BY 5 MINVALUE -100 MAXVALUE 1000 START WITH 10 CACHE 3;
CREATE SEQUENCE audit.big_seq AS INT8;

CREATE FUNCTION app.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b $$;

CREATE TABLE app.customer (
    id INT8 PRIMARY KEY DEFAULT unique_rowid(),
    alt_id INT8 UNIQUE,
    email STRING NOT NULL,
    name VARCHAR(50) COLLATE "en_US" NOT NULL DEFAULT 'anon',
    mood app.mood DEFAULT 'happy',
    home app.address,
    tags STRING[] DEFAULT ARRAY[]::STRING[],
    balance DECIMAL(12, 2) DEFAULT 0 CHECK (balance >= 0),
    ratio FLOAT8,
    small INT2,
    created TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT customer_email_check CHECK (email != '')
);
COMMENT ON TABLE app.customer IS 'a customer';
COMMENT ON COLUMN app.customer.email IS 'login';
CREATE INDEX customer_name_idx ON app.customer (name DESC) STORING (balance);
COMMENT ON INDEX app.customer@customer_name_idx IS 'names';

CREATE TABLE app.orders (
    id INT8 PRIMARY KEY,
    customer_id INT8 NOT NULL REFERENCES app.customer(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    amount DECIMAL(10, 2) NOT NULL,
    placed_at TIMESTAMP,
    CONSTRAINT orders_amount_positive CHECK (amount > 0)
);
COMMENT ON TABLE app.orders IS 'an order';
CREATE INDEX orders_customer_idx ON app.orders (customer_id);

-- No primary key: CockroachDB adds a hidden rowid, which the sync marks
-- invisible.
CREATE TABLE audit.events (at TIMESTAMPTZ, note STRING);

CREATE VIEW app.customer_orders AS
    SELECT c.id, c.name, o.id AS order_id, o.amount
    FROM app.customer AS c JOIN app.orders AS o ON o.customer_id = c.id;

-- pg_get_functiondef answers a procedure with its body alone.
CREATE PROCEDURE app.touch_orders() LANGUAGE SQL AS $$ UPDATE app.orders SET amount = amount $$;

CREATE MATERIALIZED VIEW app.order_totals AS
    SELECT customer_id, sum(amount) AS total FROM app.orders GROUP BY customer_id;

INSERT INTO app.customer (id, alt_id, email, name, mood, balance, ratio, small, created) VALUES
    (1, 11, 'one@example.com', 'one', 'happy', 10.50, 1.5, 1, '2024-01-01 00:00:00+00'),
    (2, 12, 'two@example.com', 'two', 'sad', 20.25, 2.5, 2, '2024-01-02 00:00:00+00');
INSERT INTO app.orders (id, customer_id, amount, placed_at) VALUES
    (1, 1, 10.50, '2024-02-01 12:00:00'),
    (2, 1, 4.50, '2024-02-02 12:00:00'),
    (3, 2, 20.25, '2024-02-03 12:00:00');

ANALYZE app.customer;
ANALYZE app.orders;

-- The corpus the tapes in testdata/tapes and the golden in testdata/sync were
-- recorded from. Load it once into an empty database of a scratch cluster
-- before recording with OMNI_RECORD=1; the tests replay the answers and need
-- nothing loaded.
--
-- Redshift refuses ON DELETE CASCADE on a foreign key, which is why the
-- constraint states no referential action. A primary key is the only index a
-- Redshift table reports; there is no CREATE INDEX, so the corpus writes none.
CREATE SCHEMA app;

CREATE TABLE public.lookup (id int PRIMARY KEY, label varchar(20) NOT NULL DEFAULT 'x');
COMMENT ON TABLE public.lookup IS 'a synced comment';
COMMENT ON COLUMN public.lookup.label IS 'a column comment';

CREATE TABLE app.facts (
  id bigint PRIMARY KEY,
  lookup_id int REFERENCES public.lookup (id),
  amount numeric(12,2),
  flag boolean,
  seen timestamp,
  notes varchar(80)
);

CREATE VIEW app.recent AS
  SELECT id, amount FROM app.facts WHERE seen > '2026-01-01';

-- 228 characters: shorter than the varchar(256) PG_VIEWS.definition used to
-- state, so it pinned nothing about a long view's text.
CREATE VIEW app.long_view AS
  SELECT f.id AS id, f.amount AS amount, f.notes AS notes, f.flag AS flag,
         f.seen AS seen, f.lookup_id AS lookup_id, f.id + 1 AS next_id,
         f.amount * 2 AS doubled
  FROM app.facts AS f
  WHERE f.seen > '2025-01-01' AND f.flag = true;

-- 546 characters. Redshift 1.0.539507 states it whole in PG_VIEWS.definition
-- as well, so this pins that a long definition comes back whole, not that
-- SHOW VIEW is needed to get it.
CREATE VIEW app.very_long_view AS
  SELECT f.id AS identifier_column_one,
         f.amount AS amount_column_two,
         f.notes AS notes_column_three,
         f.flag AS flag_column_four,
         f.seen AS seen_column_five,
         f.lookup_id AS lookup_column_six,
         f.id + 1 AS next_identifier_seven,
         f.amount * 2 AS doubled_amount_eight,
         f.id * 3 AS tripled_identifier_nine,
         f.notes || 'x' AS notes_with_suffix_ten
  FROM app.facts AS f
  WHERE f.seen > '2025-01-01' AND f.flag = true AND f.amount > 0 AND f.notes IS NOT NULL;

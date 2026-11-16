-- The tables and rows every cassandra/driver test records against. A recording
-- starts from an empty keyspace, so the tape holds these statements' answers
-- and nothing a previous test left behind.

CREATE TABLE omni_types (
    id int PRIMARY KEY,
    txt text,
    ascii_col ascii,
    vc varchar,
    flag boolean,
    bi bigint,
    ti tinyint,
    si smallint,
    dbl double,
    flt float,
    dec decimal,
    vi varint,
    ts timestamp,
    dt date,
    tm time,
    uid uuid,
    tuid timeuuid,
    bl blob,
    dur duration,
    lst list<int>,
    st set<text>,
    mp map<text,int>,
    addr inet,
    tup tuple<int, text>
);

INSERT INTO omni_types (
    id, txt, ascii_col, vc, flag, bi, ti, si, dbl, flt, dec, vi, ts, dt, tm, uid, tuid, bl, dur, lst, st, mp, addr, tup
) VALUES (
    1, 'text value', 'ascii', 'varchar', true, 9223372036854775807, 7, 300, 3.5, 1.25, 1.5,
    123456789012345678901234567890, '2026-01-02 03:04:05.678', '2026-01-02', '08:12:54.123456789',
    6afaf7c2-c0c2-0ddc-7f22-71b8f0a0b0c0, 50554d6e-29bb-11e5-b345-feff819cdc9f,
    0xdeadbeef, 1mo2d3h, [1, 2, 3], {'a', 'b'}, {'a': 1, 'b': 2}, '10.0.0.1', (7, 'seven')
);

INSERT INTO omni_types (id, txt) VALUES (2, NULL);

CREATE TABLE omni_rows (
    id int PRIMARY KEY,
    label text
);

INSERT INTO omni_rows (id, label) VALUES (1, 'row-1');
INSERT INTO omni_rows (id, label) VALUES (2, 'row-2');
INSERT INTO omni_rows (id, label) VALUES (3, 'row-3');
INSERT INTO omni_rows (id, label) VALUES (4, 'row-4');
INSERT INTO omni_rows (id, label) VALUES (5, 'row-5');
INSERT INTO omni_rows (id, label) VALUES (6, 'row-6');
INSERT INTO omni_rows (id, label) VALUES (7, 'row-7');
INSERT INTO omni_rows (id, label) VALUES (8, 'row-8');
INSERT INTO omni_rows (id, label) VALUES (9, 'row-9');
INSERT INTO omni_rows (id, label) VALUES (10, 'row-10');

CREATE TABLE omni_events (
    group_id int,
    bucket int,
    event_id int,
    payload text,
    note text static,
    PRIMARY KEY ((group_id), bucket, event_id)
) WITH CLUSTERING ORDER BY (bucket DESC, event_id ASC);

INSERT INTO omni_events (group_id, bucket, event_id, payload, note) VALUES (1, 2, 3, 'p1', 'first group');
INSERT INTO omni_events (group_id, bucket, event_id, payload) VALUES (1, 1, 4, 'p2');
INSERT INTO omni_events (group_id, bucket, event_id, payload, note) VALUES (2, 9, 1, 'p3', 'second group');

CREATE INDEX omni_events_payload_idx ON omni_events (payload);

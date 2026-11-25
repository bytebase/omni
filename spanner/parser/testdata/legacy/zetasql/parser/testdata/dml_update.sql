UPDATE T
SET
  x = y;
UPDATE T WITH OFFSET AS offset
SET
  x = y;
UPDATE T.(a.b).c
SET
  x = y;
UPDATE T.a[0].b
SET
  x = y;
UPDATE T
SET
  id.(path.`to`.extension) = 5;
UPDATE T
SET
  id.(path.`to`.extension) = DEFAULT;
UPDATE T
SET
  id[0] = DEFAULT;
UPDATE T
SET
  id1.id2.(path.`to`.extension) = 5;
UPDATE T
SET
  id1.(id2).(id3) = 5;
UPDATE T
SET
  id1.(a.b.c).(d.e.f).id2.(g.h.i).id3.id4 = 5;
UPDATE T
SET
  id1[0] = 5;
UPDATE T
SET
  id1[0].(a.b.c).id1.(d.e.f)[1].id3 = 5;
UPDATE T
SET
  id1[0][1] = 5;
UPDATE c.T
SET
  x = y,
  a.b.c = 1 + (
    SELECT
      count(*)
    FROM
      t2
  )
WHERE
  zzz + yyy = 55
ASSERT_ROWS_MODIFIED @row_count;
UPDATE T
SET
  x = null,
  y = DEFAULT,
  z = z;
UPDATE T AS a
SET
  a.x = 1;
UPDATE T
SET
  x = y,
  (
    DELETE x
  ),
  (
    INSERT INTO y.z
    VALUES
      (5)
  ),
  (
    UPDATE a
    SET
      b = c
  ),
  z = DEFAULT;
UPDATE T
SET
  (
    DELETE a.(b.c)
  ),
  (
    INSERT INTO d.(e.f)
    VALUES
      (5)
  ),
  (
    UPDATE a.(b.c)
    SET
      d = 10
  );
UPDATE T
SET
  (
    DELETE a.(b.c).x
  ),
  (
    INSERT INTO d.(e.f).y
    VALUES
      (5)
  ),
  (
    UPDATE a.(b.c).x
    SET
      d = 10
  );
UPDATE T
SET
  (
    DELETE a.(b.c)[0].d
  ),
  (
    INSERT INTO d.(e.f)[1]
    VALUES
      (5)
  ),
  (
    UPDATE a[2].(b.c)
    SET
      d = 10
  );
UPDATE T AS a
SET
  a.x = y,
  (
    DELETE a.x AS na
    WHERE
      na.x = 1
  ),
  (
    UPDATE a.x AS na
    SET
      na.x = 1
  ),
  (
    INSERT INTO a.x
    VALUES
      (5)
  );
UPDATE T
SET
  (
    UPDATE c1.c2
    SET
      (
        UPDATE c3
        SET
          (
            DELETE c4
            WHERE
              false
            ASSERT_ROWS_MODIFIED 5
          )
        WHERE
          true
        ASSERT_ROWS_MODIFIED 4
      )
    WHERE
      false
    ASSERT_ROWS_MODIFIED 3
  )
WHERE
  true
ASSERT_ROWS_MODIFIED 2;
UPDATE T
SET
  x = T1.y
FROM
  T1
WHERE
  T.a = T1.b;
UPDATE T
SET
  x = T2.c
FROM
  T1
  JOIN
  T2
  ON T1.x = T2.y
WHERE
  T.a < T1.b;
UPDATE T
SET
  x = T2.c
FROM
  T1,
  T2
WHERE
  T.a < T1.b;
UPDATE T
SET
  x = 1
FROM
  (
    T1
    JOIN
    T2
  )
  JOIN
  T3
WHERE
  true;

SELECT *
  FROM db1.public.htt1
    AT(TIMESTAMP => '2024-06-05 17:50:00'::TIMESTAMP_LTZ) h
    JOIN db1.public.tt1
    AT(TIMESTAMP => '2024-06-05 17:50:00'::TIMESTAMP_LTZ) t
    ON h.c1=t.c1;

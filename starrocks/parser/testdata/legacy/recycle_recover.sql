RECOVER DATABASE example_db;
RECOVER TABLE example_db.example_tbl;
RECOVER PARTITION p1 FROM example_tbl;
RECOVER DATABASE example_db 12345;
RECOVER TABLE example_db.example_tbl 12346;
RECOVER PARTITION p1 12347 FROM example_tbl;
RECOVER DATABASE example_db 12345 AS new_example_db;
RECOVER TABLE example_db.example_tbl AS new_example_tbl;
RECOVER PARTITION p1 12347 AS new_p1 FROM example_tbl

COPY INTO t1 (c1) FROM (SELECT d.$1 FROM @mystage/file1.csv.gz d);

COPY INTO @%orderstiny/result/data_
  FROM orderstiny FILE_FORMAT = (FORMAT_NAME ='myformat' COMPRESSION='GZIP');

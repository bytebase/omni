COPY INTO @my_stage/result/data_ FROM (SELECT * FROM orderstiny)
   file_format=(format_name='myformat' compression='gzip');

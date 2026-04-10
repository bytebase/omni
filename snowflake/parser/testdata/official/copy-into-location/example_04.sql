COPY INTO 's3://mybucket/unload/'
  FROM mytable
  CREDENTIALS = (AWS_KEY_ID='xxxx' AWS_SECRET_KEY='xxxxx' AWS_TOKEN='xxxxxx')
  FILE_FORMAT = (FORMAT_NAME = my_csv_format);

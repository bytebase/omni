COPY INTO 'azure://myaccount.blob.core.windows.net/unload/'
  FROM mytable
  STORAGE_INTEGRATION = myint
  FILE_FORMAT = (FORMAT_NAME = my_csv_format);

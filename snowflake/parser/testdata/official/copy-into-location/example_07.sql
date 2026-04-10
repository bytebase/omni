COPY INTO 'azure://myaccount.blob.core.windows.net/mycontainer/unload/'
  FROM mytable
  CREDENTIALS=(AZURE_SAS_TOKEN='xxxx')
  FILE_FORMAT = (FORMAT_NAME = my_csv_format);

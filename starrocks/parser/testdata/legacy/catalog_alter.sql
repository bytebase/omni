ALTER CATALOG ctlg_hive RENAME hive;
ALTER CATALOG hive SET PROPERTIES ('hive.metastore.uris'='thrift://172.21.0.1:9083');
ALTER CATALOG hive MODIFY COMMENT "new catalog comment"

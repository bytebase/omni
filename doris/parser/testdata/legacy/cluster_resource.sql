CREATE EXTERNAL RESOURCE "spark0"
PROPERTIES(
  "type" = "spark",
  "spark.master" = "yarn",
  "spark.submit.deployMode" = "cluster",
  "spark.jars" = "xxx.jar,yyy.jar",
  "spark.files" = "/tmp/aaa,/tmp/bbb",
  "spark.executor.memory" = "1g",
  "spark.yarn.queue" = "queue0",
  "spark.hadoop.yarn.resourcemanager.address" = "127.0.0.1:9999",
  "spark.hadoop.fs.defaultFS" = "hdfs://127.0.0.1:10000",
  "working_dir" = "hdfs://127.0.0.1:10000/tmp/doris",
  "broker" = "broker0",
  "broker.username" = "user0",
  "broker.password" = "password0"
);
CREATE EXTERNAL RESOURCE `oracle_odbc`
PROPERTIES (
  "type" = "odbc_catalog",
  "host" = "192.168.0.1",
  "port" = "8086",
  "user" = "test",
  "password" = "test",
  "database" = "test",
  "odbc_type" = "oracle",
  "driver" = "Oracle 19 ODBC driver"
);
CREATE RESOURCE "remote_s3"
PROPERTIES(
  "type" = "s3",
  "s3.endpoint" = "bj.s3.com",
  "s3.region" = "bj",
  "s3.access_key" = "bbb",
  "s3.secret_key" = "aaaa",
  "s3.connection.maximum" = "50",
  "s3.connection.request.timeout" = "3000",
  "s3.connection.timeout" = "1000"
);
CREATE RESOURCE mysql_resource PROPERTIES (
  "type"="jdbc",
  "user"="root",
  "password"="123456",
  "jdbc_url" = "jdbc:mysql://127.0.0.1:3316/doris_test?useSSL=false",
  "driver_url" = "https://doris-community-test-1308700295.cos.ap-hongkong.myqcloud.com/jdbc_driver/mysql-connector-java-8.0.25.jar",
  "driver_class" = "com.mysql.cj.jdbc.Driver"
);
CREATE RESOURCE hdfs_resource PROPERTIES (
  "type"="hdfs",
  "hadoop.username"="user",
  "dfs.nameservices" = "my_ha",
  "dfs.ha.namenodes.my_ha" = "my_namenode1, my_namenode2",
  "dfs.namenode.rpc-address.my_ha.my_namenode1" = "nn1_host:rpc_port",
  "dfs.namenode.rpc-address.my_ha.my_namenode2" = "nn2_host:rpc_port",
  "dfs.client.failover.proxy.provider.my_ha" = "org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider"
);
CREATE RESOURCE es_resource PROPERTIES (
  "type"="es",
  "hosts"="http://127.0.0.1:29200",
  "nodes_discovery"="false",
  "enable_keyword_sniff"="true"
);
DROP RESOURCE 'spark0';
SHOW RESOURCES;
SHOW RESOURCES WHERE NAME LIKE "2014_01_02" LIMIT 10;
SHOW RESOURCES WHERE NAME = "20140102" ORDER BY `KEY` DESC;
SHOW RESOURCES LIKE "jdbc%"

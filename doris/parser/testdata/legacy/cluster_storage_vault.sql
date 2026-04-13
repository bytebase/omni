CREATE STORAGE VAULT IF NOT EXISTS hdfs_vault_demo PROPERTIES (
    "type" = "hdfs",
    "fs.defaultFS" = "hdfs://127.0.0.1:8020",
    "path_prefix" = "big/data",
    "hadoop.username" = "user",
    "hadoop.security.authentication" = "kerberos",
    "hadoop.kerberos.principal" = "hadoop/127.0.0.1@XXX",
    "hadoop.kerberos.keytab" = "/etc/emr.keytab"
);
CREATE STORAGE VAULT IF NOT EXISTS oss_demo_vault PROPERTIES (
    "type" = "S3",
    "s3.endpoint" = "oss-cn-beijing.aliyuncs.com",
    "s3.access_key" = "xxxxxx",
    "s3.secret_key" = "xxxxxx",
    "s3.region" = "cn-beijing",
    "s3.root.path" = "oss_demo_vault_prefix",
    "s3.bucket" = "xxxxxx",
    "provider" = "OSS",
    "use_path_style" = "false"
);
CREATE STORAGE VAULT IF NOT EXISTS cos_demo_vault PROPERTIES (
    "type" = "S3",
    "s3.endpoint" = "cos.ap-guangzhou.myqcloud.com",
    "s3.access_key" = "xxxxxx",
    "s3.secret_key" = "xxxxxx",
    "s3.region" = "ap-guangzhou",
    "s3.root.path" = "cos_demo_vault_prefix",
    "s3.bucket" = "xxxxxx",
    "provider" = "COS",
    "use_path_style" = "false"
);
CREATE STORAGE VAULT IF NOT EXISTS s3_demo_vault PROPERTIES (
    "type" = "S3",
    "s3.endpoint" = "s3.us-east-1.amazonaws.com",
    "s3.access_key" = "xxxxxx",
    "s3.secret_key" = "xxxxxx",
    "s3.region" = "us-east-1",
    "s3.root.path" = "s3_demo_vault_prefix",
    "s3.bucket" = "xxxxxx",
    "provider" = "S3",
    "use_path_style" = "false"
);
CREATE STORAGE VAULT IF NOT EXISTS minio_demo_vault PROPERTIES (
    "type" = "S3",
    "s3.endpoint" = "127.0.0.1:9000",
    "s3.access_key" = "xxxxxx",
    "s3.secret_key" = "xxxxxx",
    "s3.region" = "us-east-1",
    "s3.root.path" = "minio_demo_vault_prefix",
    "s3.bucket" = "xxxxxx",
    "provider" = "S3",
    "use_path_style" = "true"
);
CREATE STORAGE VAULT IF NOT EXISTS azure_demo_vault PROPERTIES (
    "type" = "S3",
    "s3.endpoint" = "blob.core.windows.net",
    "s3.access_key" = "xxxxxx",
    "s3.secret_key" = "xxxxxx",
    "s3.region" = "us-east-1",
    "s3.root.path" = "azure_demo_vault_prefix",
    "s3.bucket" = "xxxxxx",
    "provider" = "AZURE"
);
CREATE STORAGE VAULT IF NOT EXISTS gcp_demo_vault PROPERTIES (
    "type" = "S3",
    "s3.endpoint" = "storage.googleapis.com",
    "s3.access_key" = "xxxxxx",
    "s3.secret_key" = "xxxxxx",
    "s3.region" = "us-east-1",
    "s3.root.path" = "gcp_demo_vault_prefix",
    "s3.bucket" = "xxxxxx",
    "provider" = "GCP"
);
ALTER STORAGE VAULT old_vault_name PROPERTIES (
    "type"="S3",
    "VAULT_NAME" = "new_vault_name",
    "s3.access_key" = "new_ak"
);
SHOW STORAGE VAULTS;
SHOW COMPUTE GROUPS

-- pin: driver gcp:bigquery:bytebase-dev:US sha256:a38c0d3800b3dfd8d1c6e3f3644af8fd02fc0f8ed110efc924159d3415533cd7 2026-11-21

== TestSyncDatabaseMissing driver#1
> raw
| {"op":"datasetMetadata","dataset":"omni_driver_missing_dataset"}
< raw
| {"error":{"kind":"googleapi","code":404,"message":"Not found: Dataset bytebase-dev:omni_driver_missing_dataset","errors":[{"reason":"notFound","message":"Not found: Dataset bytebase-dev:omni_driver_missing_dataset"}]}}

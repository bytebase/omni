-- pin: driver gcp:bigquery:bytebase-dev:US sha256:92a85df45321648c0feb59ee0c222e8156cb9332d322ddb9f15eaf37c8c12533 2026-11-26

== TestSyncDatabaseMissing driver#1
> raw
| {"op":"datasetMetadata","dataset":"omni_driver_missing_dataset"}
< raw
| {"error":{"kind":"googleapi","code":404,"message":"Not found: Dataset bytebase-dev:omni_driver_missing_dataset","errors":[{"reason":"notFound","message":"Not found: Dataset bytebase-dev:omni_driver_missing_dataset"}]}}

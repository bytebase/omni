-- pin: driver gcp:bigquery:bytebase-dev:US sha256:a38c0d3800b3dfd8d1c6e3f3644af8fd02fc0f8ed110efc924159d3415533cd7 2026-11-21

== TestSyncInstance driver#1
> raw
| {"op":"ping"}
< raw
| {}
> raw
| {"op":"listDatasets"}
< raw
| {"datasets":["ds1","ds2","dwh","hello","omni_record","verizon_demo","zp"]}

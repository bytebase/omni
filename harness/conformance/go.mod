module github.com/bytebase/omni/harness/conformance

go 1.26

replace github.com/bytebase/omni => ../..

require (
	github.com/bytebase/omni v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.9.3
)

require filippo.io/edwards25519 v1.1.0 // indirect

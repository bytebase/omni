// Package testcontainer pins the MariaDB server the tests of this engine run
// against. Each package starts its own container on first use, so the pin is
// what keeps them on one version: a test that records a verdict records it from
// Image, and an upgrade here means rerunning every oracle test and reading the
// drift.
package testcontainer

// Image is the pinned MariaDB. mariadb/parser, mariadb/catalog, mariadb/review
// and mariadb/validate all start it.
const Image = "mariadb:11.8.8"

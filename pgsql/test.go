package pgsql

import "testing"

type DBTest struct {
	t *testing.T
	*DB
}

func NewDBTest(t *testing.T, db *DB) *DBTest {
	return &DBTest{t, db}
}

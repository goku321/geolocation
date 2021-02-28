package e2e_test

import (
	"testing"

	geo "github.com/goku321/geolocation/geolocation"
	"github.com/goku321/geolocation/store"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var createTableStmt = `
CREATE TABLE IF NOT EXISTS geolocations (
ip character varying NOT NULL,
country_code character varying,
country character varying,
city character varying,
latitude double precision,
longitude double precision,
mystery_value bigint,
CONSTRAINT geolocations_pkey PRIMARY KEY (ip)
);
`

func createTableHelper(db *sqlx.DB, t *testing.T) {
	t.Helper()
	_, err := db.Exec(createTableStmt)
	require.NoError(t, err)
}

func TestImportService(t *testing.T) {
	connStr := "postgres://postgres:password@127.0.0.1:6432/postgres?sslmode=disable"
	file := "sample.csv"
	pg, err := sqlx.Open("postgres", connStr)
	require.NoError(t, err)

	// create table
	createTableHelper(pg, t)

	db := store.New(pg)
	defer db.Close()
	importer := geo.NewCSVImporter(db)
	x, err := importer.Parse(file)
	require.NoError(t, err)

	s, err := importer.Import(x)
	require.NoError(t, err)

	assert.Equal(t, 4, s.Inserted)
	assert.Equal(t, 1, s.Skipped)
}

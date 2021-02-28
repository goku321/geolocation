package store

import (
	"testing"

	geo "github.com/goku321/geolocation/geolocation"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var postgresTestDB *Store
var db *sqlx.DB
var createTableStmt = `
	CREATE TABLE geolocations (
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

var dropTableStmt = "DROP TABLE geolocations"

func setup(t *testing.T) {
	t.Helper()
	var err error
	connStr := "host=127.0.0.1 user=postgres dbname=postgres password=postgres port=6432 sslmode=disable"
	db, err = sqlx.Open("postgres", connStr)
	require.NoError(t, err)
	postgresTestDB = New(db)
	_, err = db.Exec(createTableStmt)
	require.NoError(t, err)
}

func teardown(t *testing.T) {
	t.Helper()
	_, err := db.Exec(dropTableStmt)
	require.NoError(t, err)
}

// todo - update this test to check both get and saveall method.
func TestSaveAll(t *testing.T) {
	setup(t)
	defer teardown(t)

	geolocations := map[string]*geo.GeoData{
		"fake":   {IP: "ip-1"},
		"fake-2": {IP: "ip-2"},
	}

	err := postgresTestDB.SaveAll(geolocations)
	require.NoError(t, err)

	t.Run("db should only have 2 records", func(t *testing.T) {
		query, args, err := postgresTestDB.stmtBuilder.Select("count(*)").From("geolocations").ToSql()
		require.NoError(t, err)

		res := db.QueryRow(query, args...)
		// require.NoError(t, err)

		var count int
		err = res.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})

	// more test case:
	// 1. primary key constaint check.
	// 2. if the records are saved properly.
}

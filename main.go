package main

import (
	"log"
	"os"

	geo "github.com/goku321/geolocation/geolocation"
	"github.com/goku321/geolocation/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func createTable(db *sqlx.DB) error {
	_, err := db.Exec(createTableStmt)
	return err
}

func main() {
	connStr := os.Getenv("DB_CONN_STR")
	// pass file name through flag
	file := os.Getenv("CSV_FILE")
	pg, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = createTable(pg); err != nil {
		log.Fatalf("failed to create table: %s", err)
	}

	db := store.New(pg)
	defer db.Close()
	importer := geo.NewCSVImporter(db)
	var s geo.Stats
	var x map[string]*geo.GeoData
	if x, err = importer.Parse(file); err != nil {
		log.Fatalf("failed to parse csv: %s", err)
	}
	s, err = importer.Import(x)
	if err != nil {
		log.Fatalf("failed to import csv: %s", err)
	}
	log.Println("csv imported successfully")
	// print stats
	s.Print()
}

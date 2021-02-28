package store

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	geo "github.com/goku321/geolocation/geolocation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Store represents a on-disk store.
type Store struct {
	stmtBuilder sq.StatementBuilderType
	db          *sqlx.DB
}

// New creates a new DB instance to interact with on disk storage.
func New(db *sqlx.DB) *Store {
	return &Store{
		stmtBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		db:          db,
	}
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Get returns a geolocation based on given ip.
func (s *Store) Get(ip string) (*geo.GeoData, error) {
	query, args, err := s.stmtBuilder.
		Select("*").
		From("geolocations").
		Where(sq.Eq{"ip": ip}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %s", err)
	}

	row := s.db.QueryRowxContext(context.Background(), query, args...)
	location := &geo.GeoData{}
	if err = row.StructScan(location); err != nil {
		return location, fmt.Errorf("error scanning geolocation: %s", err)
	}

	return location, nil
}

// SaveAll persists bulk geolocation records on disk.
func (s *Store) SaveAll(x map[string]*geo.GeoData) error {
	txn, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(pq.CopyIn("geolocations", "ip", "country", "country_code", "city", "latitude", "longitude", "mystery_value"))
	if err != nil {
		return err
	}

	for _, v := range x {
		_, err := stmt.Exec(v.IP, v.Country, v.CountryCode, v.City, v.Latitude, v.Longitude, v.MysteryValue)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	if err = txn.Commit(); err != nil {
		return err
	}

	return nil
}

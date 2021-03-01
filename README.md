# Geolocation service

## Goals
- to provide an interface to access geolocation data.
- to provide a default csv importer that parses and imports csv data to a database.

## Design
This library can be summarized using two interfaces:

### Importer
```
// Importer provides API to parse and import geolocation data.
type Importer interface {
	Parse(file string) (map[string]*GeoData, error)
	Import(map[string]*GeoData) (Stats, error)
}
```

CSV importer implements this interface. This can easily be extended to other importers like, a json importer.

### GeoDataProvider
```
// GeoDataProvider gives access to geolocation data.
type GeoDataProvider interface {
	Get(ip string) (*GeoData, error)
	SaveAll(data map[string]*GeoData) error
}
```

It allows to access or store geolocation data. `store` layer implements this interface using postgres but it can be easily be replaced with other user-defined implementation of store.

## Trade offs
- I've used postgres for store layer but a key-value store would have been a better choice.
- Case when csv headers are in a different order is not handled.
- Table name is hardcoded, could be passed from the environment.

## How to run the import script

1. Spin up a postgres container using `make start-postgres-container`

2. Set env `DB_CONN_STR` to the url of the database.

	`export DB_CONN_STR="postgres://postgres:password@127.0.0.1:6432/postgres?sslmode=disable"`

3. Set env  `CSV_FILE` to the file name of csv.

	`export CSV_FILE=filename`

4. `go run main.go`

5. Stop the postgres container with `make stop-postgres-container`


## Tests

### To run tests for the entire repo

`make test`

### Run e2e test to check the import functionality
`make e2e-test`

### Run tests for store/db layer
`make test-store`
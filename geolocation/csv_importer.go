package geolocation

import (
	"encoding/csv"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type csvImporter struct {
	Store GeoDataProvider
	stats *Stats
}

// NewCSVImporter creates a csv importer.
func NewCSVImporter(store GeoDataProvider) Importer {
	return &csvImporter{
		Store: store,
		stats: &Stats{},
	}
}

// Parse implements Importer.Parse method.
func (c *csvImporter) Parse(file string) (map[string]*GeoData, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	// Skip header row
	_, err = csvReader.Read()
	if err != nil {
		return nil, err
	}

	geoLocations := map[string]*GeoData{}
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return geoLocations, err
		}

		g, err := parse(row)
		if err != nil {
			c.stats.Skipped++
			continue
		}
		geoLocations[g.IP] = g
	}
	return geoLocations, nil
}

// Import implements Importer.Import method.
func (c *csvImporter) Import(x map[string]*GeoData) (Stats, error) {
	start := time.Now()
	if err := c.Store.SaveAll(x); err != nil {
		c.stats.Skipped = len(x)
		c.stats.timeElapsed = time.Now().Sub(start)
		return *c.stats, err
	}
	c.stats.Inserted = len(x)
	c.stats.timeElapsed = time.Now().Sub(start)
	return *c.stats, nil
}

// Parses a single row given in a csv format.
func parse(x []string) (*GeoData, error) {
	if len(x) < 7 {
		return nil, errors.New("not enough columns in a row")
	}

	x = sanitize(x)
	// Check for blank and invalid ip address.
	if isBlank(x[0]) || net.ParseIP(x[0]) == nil {
		return nil, errors.New("ip address is blank or invalid")
	}

	g := &GeoData{}
	var err error
	if g.Latitude, err = strconv.ParseFloat(x[4], 64); err != nil {
		return nil, err
	}
	if g.Longitude, err = strconv.ParseFloat(x[5], 64); err != nil {
		return nil, err
	}
	if g.MysteryValue, err = strconv.ParseInt(x[6], 10, 64); err != nil {
		return nil, err
	}
	g.IP = x[0]
	g.CountryCode = x[1]
	g.Country = x[2]
	g.City = x[3]
	return g, nil
}

func sanitize(row []string) []string {
	for i, col := range row {
		row[i] = strings.TrimSpace(col)
	}
	return row
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

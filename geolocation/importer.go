package geolocation

import (
	"fmt"
	"time"
)

// Importer provides API to parse and import geolocation data.
type Importer interface {
	Parse(file string) (map[string]*GeoData, error)
	Import(map[string]*GeoData) (Stats, error)
}

// Stats contains import statistics such as:
// time elapsed, records skipped/inserted etc.
type Stats struct {
	Skipped     int
	Inserted    int
	timeElapsed time.Duration
}

// Print prints the stats.
func (s *Stats) Print() {
	fmt.Printf("time elapsed: %d\ninserted: %d\nskipped: %d\n", s.timeElapsed, s.Inserted, s.Skipped)
}

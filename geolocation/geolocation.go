package geolocation

// GeoDataProvider gives access to geolocation data.
type GeoDataProvider interface {
	Get(ip string) (*GeoData, error)
	SaveAll(data map[string]*GeoData) error
}

// GeoData represent geolocation data.
type GeoData struct {
	IP           string  `db:"ip"`
	CountryCode  string  `db:"country_code"`
	Country      string  `db:"country"`
	City         string  `db:"city"`
	Latitude     float64 `db:"latitude"`
	Longitude    float64 `db:"longitude"`
	MysteryValue int64   `db:"mystery_value"`
}

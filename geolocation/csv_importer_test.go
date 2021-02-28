package geolocation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStore struct{}

func (m *mockStore) Get(ip string) (*GeoData, error) {
	return &GeoData{}, nil
}

func (m *mockStore) SaveAll(data map[string]*GeoData) error {
	return nil
}

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		file    string
		want    map[string]*GeoData
		wantErr error
	}{
		{
			name:    "file not exists",
			file:    "not-exists",
			want:    nil,
			wantErr: errors.New("no such file or directory"),
		},
		{
			name: "file exists with valid 5 rows",
			file: "sample.csv",
			want: map[string]*GeoData{
				"200.106.141.15": {
					"200.106.141.15",
					"SI",
					"Nepal",
					"DuBuquemouth",
					-84.87503094689836,
					7.206435933364332,
					7823011346,
				},
				"160.103.7.140": {
					"160.103.7.140",
					"CZ",
					"Nicaragua",
					"New Neva",
					-68.31023296602508,
					-37.62435199624531,
					7301823115,
				},
				"70.95.73.73": {
					"70.95.73.73",
					"TL",
					"Saudi Arabia",
					"Gradymouth",
					-49.16675918861615,
					-86.05920084416894,
					2559997162,
				},
				"125.159.20.54": {
					"125.159.20.54",
					"LI",
					"Guyana",
					"Port Karson",
					-78.2274228596799,
					-163.26218895343357,
					1337885276,
				},
			},
			wantErr: nil,
		},
		// Handle more edge cases.
	}

	mockDB := &mockStore{}
	csvImporter := NewCSVImporter(mockDB)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			x, err := csvImporter.Parse(c.file)
			if c.wantErr != nil {
				require.Contains(t, err.Error(), c.wantErr.Error())
			}
			require.Len(t, x, len(c.want))
			require.Equal(t, c.want, x)
		})
	}
}

func TestImport(t *testing.T) {

}

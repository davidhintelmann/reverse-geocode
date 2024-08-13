package node_test

import (
	_ "embed"
	"testing"

	. "github.com/davidhintelmann/reverse-geocode/node"
)

// Test Haversine with some known distances
func TestHaversine(t *testing.T) {
	tests := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
		want       float64
	}{
		{"New York to Los Angeles", 40.7128, -74.0060, 34.0522, -118.2437, 3935.9},
		{"London to Paris", 51.5074, -0.1278, 48.8566, 2.3522, 346.3},
		// {"Tokyo to New York", 35.6895, 139.6919, 40.7128, -74.0060, 9231.3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2); !MathEqualWithinAbsRel(got, tt.want, 5.0) {
				t.Errorf("Haversine(...) = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

func TestHaversineSpecialCases(t *testing.T) {
	// Test with identical points
	if got := Haversine(40.7128, -74.0060, 40.7128, -74.0060); got != 0 {
		t.Errorf("Haversine(...) = %.2f, want 0", got)
	}

	// // Test with antipodal points
	// if got := Haversine(40.7128, -74.0060, 59.9455, 124.6866); got != 20000 {
	// 	t.Errorf("Haversine(...) = %.2f, want 20000", got)
	// }
}

func TestDistance(t *testing.T) {
	c1 := City{Latitude: 40.7128, Longitude: -74.0060}
	c2 := City{Latitude: 34.0522, Longitude: -118.2437}

	dist := Distance(c1, c2)
	if dist < Haversine(c1.Latitude, c1.Longitude, c2.Latitude, c2.Longitude) {
		t.Errorf("Expected Euclidean distance to be at least Haversine distance")
	}
}

func TestMedian(t *testing.T) {
	cities := []City{
		{Latitude: 40.7128, Longitude: -74.0060},
		{Latitude: 34.0522, Longitude: -118.2437},
		{Latitude: 37.7749, Longitude: -122.4194},
	}

	median := Median(cities, 0)
	if median.Latitude != 37.7749 {
		t.Errorf("Expected latitude to be 37.7749")
	}
}

func TestBuildKDTree(t *testing.T) {
	cities := []City{
		{Latitude: 40.7128, Longitude: -74.0060},
		{Latitude: 34.0522, Longitude: -118.2437},
		{Latitude: 37.7749, Longitude: -122.4194},
	}

	tree := NewKDTree(cities)
	if tree.Root == nil {
		t.Errorf("Expected tree to be non-nil")
	}
}

// build kd tree and create some test cases
func TestKDTree(t *testing.T) {
	tests := []struct {
		Latitude  float64
		Longitude float64
		CityName  string
		Country   string
	}{
		{Latitude: 40.71, Longitude: -74.00, CityName: "New York City", Country: "US"},
		{Latitude: 37.77, Longitude: -122.42, CityName: "San Francisco", Country: "US"},
		{Latitude: 51.5072, Longitude: -0.1275, CityName: "London", Country: "GB"},
		{Latitude: 35.6828, Longitude: 139.7594, CityName: "Tokyo", Country: "JP"},
		{Latitude: -33.8678, Longitude: 151.2100, CityName: "Sydney", Country: "AU"},
		{Latitude: 30.04, Longitude: 31.23, CityName: "Cairo", Country: "EG"},
		{Latitude: -22.9111, Longitude: -43.2056, CityName: "Rio de Janeiro", Country: "BR"},
		{Latitude: 48.85, Longitude: 2.35, CityName: "Paris", Country: "FR"},
		{Latitude: -33.9253, Longitude: 18.4239, CityName: "Cape Town", Country: "ZA"},
		{Latitude: 39.90, Longitude: 116.40, CityName: "Beijing", Country: "CH"},
		{Latitude: 19.43, Longitude: -99.13, CityName: "Mexico City", Country: "MX"},
		{Latitude: 43.7417, Longitude: -79.3733, CityName: "Bridle Path-Sunnybrook-York Mills", Country: "CA"},
		{Latitude: 43.65, Longitude: -79.38, CityName: "Moss Park", Country: "CA"},
		{Latitude: 44.30, Longitude: -78.31, CityName: "Peterborough", Country: "CA"},
		{Latitude: 45.50, Longitude: -73.56, CityName: "Montréal", Country: "CA"},
		{Latitude: 52.23, Longitude: 21.01, CityName: "Warsaw", Country: "PL"},
	}

	err := ParseEmbeddedCSV()
	if err != nil {
		t.Fatalf("Encountered a problem parsing the embedded csv.\nError: %v\n", err)
	}

	kdTree := NewKDTree(DataPoints)

	for _, test := range tests {
		t.Run(test.CityName, func(t *testing.T) {
			actual := kdTree.FindNearestNeighbor(test)
			if actual.City.CityName != test.CityName {
				t.Errorf(
					"FindNearestNeighbor(%v) = %v, want %v",
					test,
					actual.City.CityName,
					test.CityName)
			}
		})
	}
}

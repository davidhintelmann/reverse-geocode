//go:build generate

package main

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	node "github.com/davidhintelmann/reverse-geocode/node"
)

// Geoname ID
// Name
// ASCII Name
// Alternate Names
// Feature Class
// Feature Code
// Country Code
// Country name EN
// Country Code 2
// Admin1 Code
// Admin2 Code
// Admin3 Code
// Admin4 Code
// Population
// Elevation
// DIgital Elevation Model
// Timezone
// Modification date
// LABEL EN
// Coordinates

//go:embed geonames.csv
var csvdata string

// var csvData string

func main() {
	// Create some sample points
	// points := []node.Point{
	// 	{40.71, -74.00, "New York City, USA"}, // New York City, USA
	// 	{51.50, -0.12, "London, GBR"},         // London, United Kingdom
	// 	{35.6828, 139.7594, "Tokyo, JPN"},     // Tokyo, Japan
	// 	{33.86, 151.20, "Sydney, AUS"},        // Sydney, Australia
	// 	{30.04, 31.23, "Cairo, EGY"},          // Cairo, Egypt
	// 	{22.90, 43.17, "Rio de Janeiro, BRA"}, // Rio de Janeiro, Brazil
	// 	{48.85, 2.35, "Paris, FRA"},           // Paris, France
	// 	{33.92, 18.42, "Cape Town, ZAF"},      // Cape Town, South Africa
	// 	{39.90, 116.40, "Beijing, CHA"},       // Beijing, China
	// 	{19.43, -99.13, "Mexico City, MEX"},   // Mexico City, Mexico
	// 	{43.65, -79.38, "Toronto, CAN"},       // Toronto, Canada
	// 	{44.30, -78.31, "Peterborough, CAN"},  // Peterborough, Canada
	// 	{45.50, -73.56, "Montreal, CAN"},      // Montreal, Canada
	// 	{52.23, 21.01, "Warsaw, POL"},         // Warsaw, Poland
	// }

	// queryPoint := node.City{Latitude: 37.7749, Longitude: -122.4194, CityName: "San Francisco", Country: "USA"}
	queryPoint := node.City{Latitude: 44.03, Longitude: -79.30, CityName: "Toronto", Country: "Canada"}
	startParseEmbeddedCSV := time.Now()
	err := ParseEmbeddedCSV()
	if err != nil {
		log.Fatalf("Encountered a problem parsing the embedded csv.\nError: %v\n", err)
	}
	fmt.Printf("Parsing embedded csv line by line took %v seconds.\n", time.Since(startParseEmbeddedCSV))

	// Build a KD-tree from the sample points
	// kdTree := node.NewNode(points, 0)
	startNewNode := time.Now()
	kdTree := node.NewKDTree(dataPoints)
	endNewNode := time.Now()
	durationNewNode := endNewNode.Sub(startNewNode)
	fmt.Printf("Building k-d tree took %.02f seconds.\n", durationNewNode.Seconds())
	// Query a point to find its nearest neighbor
	// queryPoint := Point{6, 5}
	// queryPoint := node.Point{35.91, 127.77, "Korea"} // 52.52, 13.41, "Berlin"   35.91, 127.77, "Korea"
	nearestNeighbor := kdTree.FindNearestNeighbor(queryPoint)

	fmt.Printf("Want Latitude: %v, Longitude: %v, CityName: %v, Country: %v\n", queryPoint.Latitude, queryPoint.Longitude, queryPoint.CityName, queryPoint.Country)
	fmt.Printf("Got  Latitude: %v, Longitude: %v, CityName: %v, Country: %v\n", nearestNeighbor.City.Latitude, nearestNeighbor.City.Longitude, nearestNeighbor.City.CityName, nearestNeighbor.City.Country)
}

// points from embedded csv
var dataPoints []node.City

func ParseEmbeddedCSV() error {
	reader := csv.NewReader(strings.NewReader(csvdata))
	// reader := csv.NewReader(strings.NewReader(csvData))
	reader.Comma = ';'
	// skip header
	_, err := reader.Read()
	if err != nil {
		return err
	}
	// fmt.Println(header)

	for data, err := reader.Read(); err != io.EOF; data, err = reader.Read() {
		singleRow := data[19]
		singleCoordinate := strings.Split(singleRow, ", ")
		singleLatFloat, err := strconv.ParseFloat(singleCoordinate[0], 64)
		if err != nil {
			return err
		}
		singleLongFloat, err := strconv.ParseFloat(singleCoordinate[1], 64)
		if err != nil {
			return err
		}

		// node.Point{Latitude, Longitude, City Name, Country Name}
		dataPoint := node.City{singleLatFloat, singleLongFloat, data[1], data[6]}
		dataPoints = append(dataPoints, dataPoint)
	}
	return nil
}

package main

import (
	_ "embed"
	"fmt"
	"log"
	"time"

	node "github.com/davidhintelmann/reverse-geocode/node"
)

// //go:embed geonames.csv
// var csvdata string

// // points from embedded csv
// var dataPoints []node.City

func main() {
	// queryPoint := node.City{Latitude: 37.7749, Longitude: -122.4194, CityName: "San Francisco", Country: "USA"}
	queryPoint := node.City{Latitude: 44.03, Longitude: -79.30, CityName: "Toronto", Country: "Canada"}
	// queryPoint := node.City{Latitude: 44.03, Longitude: -79.30}
	startParseEmbeddedCSV := time.Now()
	err := node.ParseEmbeddedCSV()
	if err != nil {
		log.Fatalf("Encountered a problem parsing the embedded csv.\nError: %v\n", err)
	}
	fmt.Printf("Parsing embedded csv line by line took %v\n", time.Since(startParseEmbeddedCSV))

	// Build a KD-tree from the sample points
	startNewNode := time.Now()
	kdTree := node.NewKDTree(node.DataPoints)
	fmt.Printf("Building k-d tree took %v\n", time.Since(startNewNode))
	nearestNeighbor := kdTree.FindNearestNeighbor(queryPoint)

	fmt.Printf("Want Latitude: %v, Longitude: %v, CityName: %v, Country: %v\n", queryPoint.Latitude, queryPoint.Longitude, queryPoint.CityName, queryPoint.Country)
	fmt.Printf("Got  Latitude: %v, Longitude: %v, CityName: %v, Country: %v\n", nearestNeighbor.City.Latitude, nearestNeighbor.City.Longitude, nearestNeighbor.City.CityName, nearestNeighbor.City.Country)
}

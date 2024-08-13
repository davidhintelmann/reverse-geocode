package node

import (
	"encoding/csv"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

// var CSVData string

// points from embedded csv
var DataPoints []City

type City struct {
	Latitude, Longitude float64
	CityName, Country   string
}

type KDTreeNode struct {
	City        City
	Left, Right *KDTreeNode
	Depth       int
}

type KDTree struct {
	Root *KDTreeNode
}

// Haversine distance between two cities (approximation of the great-circle distance)
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of Earth in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func MathEqualWithinAbsRel(a, b float64, absTol float64) bool {
	return math.Abs(a-b) <= absTol
}

func Distance(city1, city2 City) float64 {
	return Haversine(city1.Latitude, city1.Longitude, city2.Latitude, city2.Longitude)
}

// Get the median city based on the specified dimension (latitude or longitude)
func Median(cities []City, dim int) City {
	sort.Slice(cities, func(i, j int) bool {
		if dim == 0 {
			return cities[i].Latitude < cities[j].Latitude
		}
		return cities[i].Longitude < cities[j].Longitude
	})
	return cities[len(cities)/2]
}

// KD-Tree construction function
func BuildKDTree(cities []City, depth int) *KDTreeNode {
	if len(cities) == 0 {
		return nil
	}

	axis := depth % 2 // We only have 2 dimensions: Latitude and Longitude

	medianCity := Median(cities, axis)
	medianIndex := len(cities) / 2

	return &KDTreeNode{
		City:  medianCity,
		Depth: axis,
		Left:  BuildKDTree(cities[:medianIndex], depth+1),
		Right: BuildKDTree(cities[medianIndex+1:], depth+1),
	}
}

// NewKDTree initializes a new KD-Tree
func NewKDTree(cities []City) *KDTree {
	return &KDTree{
		Root: BuildKDTree(cities, 0),
	}
}

// Nearest neighbor search in the KD-Tree
func (tree *KDTree) NearestNeighbor(root *KDTreeNode, target City, depth int, best *KDTreeNode, bestDist *float64) *KDTreeNode {
	if root == nil {
		return best
	}

	d := Distance(root.City, target)
	if d < *bestDist {
		*bestDist = d
		best = root
	}

	axis := depth % 2

	var nextBranch, otherBranch *KDTreeNode
	if (axis == 0 && target.Latitude < root.City.Latitude) || (axis == 1 && target.Longitude < root.City.Longitude) {
		nextBranch = root.Left
		otherBranch = root.Right
	} else {
		nextBranch = root.Right
		otherBranch = root.Left
	}

	best = tree.NearestNeighbor(nextBranch, target, depth+1, best, bestDist)

	var planeDist float64
	if axis == 0 {
		planeDist = target.Latitude - root.City.Latitude
	} else {
		planeDist = target.Longitude - root.City.Longitude
	}

	if math.Abs(planeDist) < *bestDist {
		best = tree.NearestNeighbor(otherBranch, target, depth+1, best, bestDist)
	}

	return best
}

// FindNearestNeighbor is a method of KDTree to find the nearest neighbor of a target city
func (tree *KDTree) FindNearestNeighbor(target City) *KDTreeNode {
	bestDist := math.Inf(1)
	return tree.NearestNeighbor(tree.Root, target, 0, nil, &bestDist)
}

// This function will parse the embedded csv file into the go binary.
// If the csv file is lost this go program will still work
func ParseEmbeddedCSV() error {
	reader := csv.NewReader(strings.NewReader(CSVData))
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
		dataPoint := City{singleLatFloat, singleLongFloat, data[1], data[6]}
		DataPoints = append(DataPoints, dataPoint)
	}
	return nil
}

package node

import (
	"math"
	"sort"
)

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
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of Earth in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func distance(city1, city2 City) float64 {
	return haversine(city1.Latitude, city1.Longitude, city2.Latitude, city2.Longitude)
}

// Get the median city based on the specified dimension (latitude or longitude)
func median(cities []City, dim int) City {
	sort.Slice(cities, func(i, j int) bool {
		if dim == 0 {
			return cities[i].Latitude < cities[j].Latitude
		}
		return cities[i].Longitude < cities[j].Longitude
	})
	return cities[len(cities)/2]
}

// KD-Tree construction function
func buildKDTree(cities []City, depth int) *KDTreeNode {
	if len(cities) == 0 {
		return nil
	}

	axis := depth % 2 // We only have 2 dimensions: Latitude and Longitude

	medianCity := median(cities, axis)
	medianIndex := len(cities) / 2

	return &KDTreeNode{
		City:  medianCity,
		Depth: axis,
		Left:  buildKDTree(cities[:medianIndex], depth+1),
		Right: buildKDTree(cities[medianIndex+1:], depth+1),
	}
}

// NewKDTree initializes a new KD-Tree
func NewKDTree(cities []City) *KDTree {
	return &KDTree{
		Root: buildKDTree(cities, 0),
	}
}

// Nearest neighbor search in the KD-Tree
func (tree *KDTree) nearestNeighbor(root *KDTreeNode, target City, depth int, best *KDTreeNode, bestDist *float64) *KDTreeNode {
	if root == nil {
		return best
	}

	d := distance(root.City, target)
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

	best = tree.nearestNeighbor(nextBranch, target, depth+1, best, bestDist)

	var planeDist float64
	if axis == 0 {
		planeDist = target.Latitude - root.City.Latitude
	} else {
		planeDist = target.Longitude - root.City.Longitude
	}

	if math.Abs(planeDist) < *bestDist {
		best = tree.nearestNeighbor(otherBranch, target, depth+1, best, bestDist)
	}

	return best
}

// FindNearestNeighbor is a method of KDTree to find the nearest neighbor of a target city
func (tree *KDTree) FindNearestNeighbor(target City) *KDTreeNode {
	bestDist := math.Inf(1)
	return tree.nearestNeighbor(tree.Root, target, 0, nil, &bestDist)
}

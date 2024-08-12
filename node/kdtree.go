package node

import (
	"math"
	"sort"
)

type Point struct {
	Latitude, Longitude float64
	City, Country       string
}

type KDTree struct {
	Point       Point
	Left, Right *KDTree
	Depth       int
}

func NewKDTree(points []Point, depth int) *KDTree {
	if len(points) == 0 {
		return nil
	}

	k := 2 // We're dealing with 2D points (latitude and longitude)
	axis := depth % k

	// Sort points by the current axis (latitude or longitude)
	sort.Slice(points, func(i, j int) bool {
		if axis == 0 {
			return points[i].Latitude < points[j].Latitude
		}
		return points[i].Longitude < points[j].Longitude
	})

	// Select the median point
	median := len(points) / 2

	// Recursively build the KD-tree
	return &KDTree{
		Point: points[median],
		Left:  NewKDTree(points[:median], depth+1),
		Right: NewKDTree(points[median+1:], depth+1),
		Depth: depth,
	}
}

func (t *KDTree) FindNearestNeighbor(target Point) Point {
	return t.nearestNeighborHelper(target, t, math.Inf(1)).Point
}

func (t *KDTree) nearestNeighborHelper(target Point, best *KDTree, bestDist float64) *KDTree {
	if t == nil {
		return best
	}

	// Calculate current distance
	distance := haversine(t.Point.Latitude, t.Point.Longitude, target.Latitude, target.Longitude)

	if distance < bestDist {
		best = t
		bestDist = distance
	}

	axis := t.Depth % 2

	var nextBranch, otherBranch *KDTree
	if (axis == 0 && target.Latitude < t.Point.Latitude) || (axis == 1 && target.Longitude < t.Point.Longitude) {
		nextBranch = t.Left
		otherBranch = t.Right
	} else {
		nextBranch = t.Right
		otherBranch = t.Left
	}

	best = nextBranch.nearestNeighborHelper(target, best, bestDist)

	// Check if we need to search the other branch
	var diff float64
	if axis == 0 {
		diff = math.Abs(target.Latitude - t.Point.Latitude)
	} else {
		diff = math.Abs(target.Longitude - t.Point.Longitude)
	}

	if diff < bestDist {
		best = otherBranch.nearestNeighborHelper(target, best, bestDist)
	}

	return best
}

// Haversine formula to calculate the great-circle distance between two points on the Earth's surface
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1 = lat1 * (math.Pi / 180.0)
	lat2 = lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

package node

import "math"

type Point struct {
	X, Y    float64
	City    string
	Country string
}

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

type Node struct {
	Point                 Point
	Left, Right           *Node
	SplitAxis, SplitValue int
}

func NewNode(points []Point, depth int) *Node {
	if len(points) == 0 {
		return nil
	}

	axis := depth % 2 // Alternating between X and Y axes in 2D space
	median := len(points) / 2

	// Sort points along the current axis
	if axis == 0 {
		sortByX(points)
	} else {
		sortByY(points)
	}

	node := &Node{
		Point:      points[median],
		SplitAxis:  axis,
		SplitValue: median,
		Left:       NewNode(points[:median], depth+1),
		Right:      NewNode(points[median+1:], depth+1),
	}

	return node
}

func sortByX(points []Point) {
	// Implement a simple insertion sort for sorting points by X coordinate
	for i := 1; i < len(points); i++ {
		j := i
		for j > 0 && points[j-1].X > points[j].X {
			points[j], points[j-1] = points[j-1], points[j]
			j--
		}
	}
}

func sortByY(points []Point) {
	// Implement a simple insertion sort for sorting points by Y coordinate
	for i := 1; i < len(points); i++ {
		j := i
		for j > 0 && points[j-1].Y > points[j].Y {
			points[j], points[j-1] = points[j-1], points[j]
			j--
		}
	}
}

func distanceSquared(p1, p2 Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return dx*dx + dy*dy
}

func (n *Node) FindNearestNeighbor(target Point, bestNeighbor *Node) *Node {
	if n == nil {
		return bestNeighbor
	}

	currentBest := bestNeighbor
	// variable below is not being used
	// currentAxis := n.SplitAxis

	// Choose the next branch to search
	var nearChild, farChild *Node
	if target.X < n.Point.X || (target.X == n.Point.X && target.Y < n.Point.Y) {
		nearChild, farChild = n.Left, n.Right
	} else {
		nearChild, farChild = n.Right, n.Left
	}

	// Recursively search the nearest subtree
	currentBest = nearChild.FindNearestNeighbor(target, currentBest)

	// Check if the current node is closer than the best found so far
	if currentBest == nil || distanceSquared(target, n.Point) < distanceSquared(target, currentBest.Point) {
		currentBest = n
	}

	// Check the other side of the splitting plane if necessary
	if currentBest == nil || math.Pow(target.X-n.Point.X, 2)+math.Pow(target.Y-n.Point.Y, 2) <= distanceSquared(target, currentBest.Point) {
		currentBest = farChild.FindNearestNeighbor(target, currentBest)
	}

	return currentBest
}

func (n *Node) FindKNearestNeighbors(target Point, k int) ([]*Node, error) {
	if k < 1 {
		return nil, &KGreaterThanZeroError{}
	} else if k == 1 {
		return nil, &KGreaterThanOneError{}
	} else {
		bestNeighbors := make([]*Node, 0, k)
		n.findKNearestNeighbors(target, k, &bestNeighbors)
		return bestNeighbors, nil
	}
}

func (n *Node) findKNearestNeighbors(target Point, k int, bestNeighbors *[]*Node) {
	if n == nil {
		return
	}

	// variable below is not being used
	// currentAxis := n.SplitAxis

	// Choose the next branch to search
	var nearChild, farChild *Node
	if target.X < n.Point.X || (target.X == n.Point.X && target.Y < n.Point.Y) {
		nearChild, farChild = n.Left, n.Right
	} else {
		nearChild, farChild = n.Right, n.Left
	}

	// Recursively search the nearest subtree
	nearChild.findKNearestNeighbors(target, k, bestNeighbors)

	// Check if the current node is closer than the farthest neighbor in the list
	if len(*bestNeighbors) < k {
		*bestNeighbors = append(*bestNeighbors, n)
		// check if current node is nearest neighbor in bestNeighbors
		// if it is update bestNeighbors
		// else continue
		for i := 0; i < len(*bestNeighbors); i++ {
			if distanceSquared(target, n.Point) < distanceSquared(target, (*bestNeighbors)[i].Point) {
				temp := (*bestNeighbors)[i] // previous node
				(*bestNeighbors)[i] = n     // new node
				(*bestNeighbors)[len(*bestNeighbors)-1] = temp
			}
		}
	}

	// Check the other side of the splitting plane if necessary
	if len(*bestNeighbors) < k || math.Pow(target.X-n.Point.X, 2)+math.Pow(target.Y-n.Point.Y, 2) <= distanceSquared(target, (*bestNeighbors)[0].Point) {
		farChild.findKNearestNeighbors(target, k, bestNeighbors)
	}
}

type NodeList []*Node

func (nl NodeList) Len() int {
	return len(nl)
}

func (nl NodeList) Max(target Point) *Node {
	if len(nl) == 0 {
		return nil
	}

	maxNode := nl[0]
	maxDistance := distanceSquared(nl[0].Point, target)

	for _, node := range nl {
		dist := distanceSquared(node.Point, target)
		if dist > maxDistance {
			maxNode = node
			maxDistance = dist
		}
	}

	return maxNode
}

func (nl *NodeList) Add(node *Node) {
	*nl = append(*nl, node)
}

func (nl *NodeList) RemoveMax(target Point) {
	maxIndex := 0
	maxDistance := distanceSquared((*nl)[0].Point, target)

	for i, node := range *nl {
		dist := distanceSquared(node.Point, target)
		if dist > maxDistance {
			maxIndex = i
			maxDistance = dist
		}
	}

	*nl = append((*nl)[:maxIndex], (*nl)[maxIndex+1:]...)
}

// custom error code, enforce k > 0
// for FindKNearestNeighbors function.
type KGreaterThanZeroError struct{}

func (m *KGreaterThanZeroError) Error() string {
	return "k parameter must be greater than zero for FindKNearestNeighbors function."
}

// custom error code, enforce k > 1
// for FindKNearestNeighbors function.
type KGreaterThanOneError struct{}

func (m *KGreaterThanOneError) Error() string {
	return "k parameter must be greater than '1' (one) for FindKNearestNeighbors function.\nConsider using FindNearestNeighbor function for when k=1 instead."
}

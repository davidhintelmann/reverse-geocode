# Reverse Geocode Project

This project utilizes a KD Tree data structure to efficiently find the nearest city given a pair of longitude and latitude coordinates. The code is written in Go.

Additionally, this repo aims to reverse geocode embedded CSV data into a format suitable for further processing. The goal is to extract latitude and longitude coordinates from the CSV data and store them in a structured format.

[Opendatasoft](https://public.opendatasoft.com/explore/dataset/geonames-all-cities-with-a-population-1000/table/?disjunctive.cou_name_en&sort=name) has as a dataset from [GeoNames](https://www.geonames.org/about.html) which contains all cities with a population greater than a thousand people. This project has been inspired by Richard Penman's [reverse-geocode](https://pypi.org/project/reverse-geocode/) project. Like Richard Penman's approach, I am also using a [k-d tree](https://en.wikipedia.org/wiki/K-d_tree) which is implemented using the [SciPy](https://docs.scipy.org/doc/scipy/reference/generated/scipy.spatial.KDTree.html) library.

This package will take a latitude and longitude coordinate which will return the country, city name, alternate name(s), population, timezone, and exact coordinates for that city. 

## Prerequisites
Go 1.21 or later (as specified in go.mod)
- Familiarity with Go programming language
- Embedded CSV data available (not shown in this code snippet, but stored in csvdata variable)

## Code Walkthrough and Key Features

- The NewKDTree function (node/kdtree.go) creates a new KD Tree by recursively splitting the input list of cities along their median values. This ensures that each node in the tree has roughly half the number of cities as its parent.
- The FindNearestNeighbor function (node/kdtree.go) uses the constructed KD Tree to find the nearest city to a given target city. It traverses the tree by recursively selecting the subtree with cities closer to the target and calculating their distances.
- Efficient search algorithm with an average time complexity of O(log n)
- Supports searching for cities in 2D space (longitude and latitude)

## Usage and How it works

To use this project, you will need to:
1. Create a list of City structs (node/kdtree.go), each representing a city with longitude and latitude coordinates.
2. Call the NewKDTree function to construct a KD Tree from your list of cities.
3. Pass the constructed KD Tree along with a target city to the FindNearestNeighbor function to find the nearest city.

To simply run the program, as is, run the PowerShell script `build.ps1`. The program will parse the embedded CSV data and print the header row. It then reads each row of the CSV file and extracts the latitude and longitude coordinates from the 20th column (using zero-based indexing). These coordinates are stored in a struct (node.Point) along with the city name and country name.

The NewKDTree function creates a new KD tree by recursively partitioning the input list of cities based on their median value. The FindNearestNeighbor function then uses this KD tree to search for the nearest city to a given target location.

### Example Usage

```go
cities := []City{
    {Latitude: 40.7128, Longitude: -74.0060},
    {Latitude: 59.9455, Longitude: 124.6866},
    {Latitude: 51.5074, Longitude: -0.1278},
}
```

We can create a KD Tree from this list and find the nearest city to a target city with longitude and latitude (37.7749, -122.4197) as follows:

```go
root := NewKDTree(cities, 0)
target := City{Latitude: 37.7749, Longitude: -122.4197}
nearestCity := FindNearestNeighbor(root, target).Point
```

This will output the nearest city to the target city.

## API Documentation

### Functions

```go
NewKDTree(cities []City, depth int) *KDTree
```

Creates a new KD tree from the input list of cities. The depth parameter determines the number of recursive partitions to perform.

```go
FindNearestNeighbor(root *KDTree, target City) *KDTree
```

Uses the given KD tree to search for the nearest city to the target location. Returns the root node of the KD tree containing the nearest city.

### Types

```go
City
```

Represents a city with a latitude and longitude coordinate.

```go
KDTree
```

Represents a KD tree data structure, which is used to efficiently search for cities in 2D space.

## Contributing
This project is open to contributions! If you'd like to help improve the parsing logic or add new features, please submit a pull request.

## Note

This README assumes that the CSV data is properly formatted and available. You may want to add more information about how to obtain or prepare the CSV data if it's not self-explanatory.
# Reverse Geocode Project
This project aims to reverse geocode embedded CSV data into a format suitable for further processing. The goal is to extract latitude and longitude coordinates from the CSV data and store them in a structured format.

## Prerequisites
Go 1.21 or later (as specified in go.mod)
- Familiarity with Go programming language
- Embedded CSV data available (not shown in this code snippet, but stored in csvdata variable)

## Usage
To run the program, simply execute main.go. The program will parse the embedded CSV data and print the header row. It then reads each row of the CSV file and extracts the latitude and longitude coordinates from the 20th column (using zero-based indexing). These coordinates are stored in a struct (node.Point) along with the city name and country name.

## Contributing
This project is open to contributions! If you'd like to help improve the parsing logic or add new features, please submit a pull request.

## Note

This README assumes that the CSV data is properly formatted and available. You may want to add more information about how to obtain or prepare the CSV data if it's not self-explanatory.
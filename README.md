# Traveling Salesman

## Overview

This is a Go implementation of the
[traveling salesman problem](https://en.wikipedia.org/wiki/Travelling_salesman_problem) which will
find the shortest distance of a round-trip path between a list of cartesian points.

## Input format

The app reads from a text file. The first line of the file must be the number of points.
Subsequent lines must contain the _x_ and _y_ coordinates of the points, separated by whitespace.
See `test_set.txt` for an example.

## Usage

```
go run main.go -i <input_data.txt>
```

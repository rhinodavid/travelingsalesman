package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
	subset "github.com/rhinodavid/travellingsalesman/set"
)

var (
	input            = flag.String("i", "test_set.txt", "input filename")
	maxFloat float32 = math.MaxFloat32
)

type coords struct {
	x, y float32
}

func cartesianDist(a, b *coords) float32 {
	return float32(math.Sqrt(float64((a.y-b.y)*(a.y-b.y) + (a.x-b.x)*(a.x-b.x))))
}

func main() {
	flag.Parse()
	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	params := scanner.Text()
	p := strings.Fields(params)
	n, _ := strconv.Atoi(p[0])

	a := make([]*coords, n)
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		x, err := strconv.ParseFloat(fields[0], 32)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseFloat(fields[1], 32)
		if err != nil {
			panic(err)
		}
		a[i] = &coords{x: float32(x), y: float32(y)}
		i++
	}
	distances := make([][]float32, n)
	for i := 0; i < n; i++ {
		distances[i] = make([]float32, n)
		for j := 0; j < n; j++ {
			distances[i][j] = cartesianDist(a[i], a[j])
		}
	}
	log.Printf("Finsihed generating distances\n")
	subsets := subset.GenerateSubsets(n)
	log.Printf("Finished generating subsets\n")
	oldCache := subset.GenerateCache(subset.FilterByCardinality(subsets, 0), n)
	// initialize empty set
	for i := 1; i < n; i++ {
		dist := distances[0][i]
		oldCache[subset.Hash(subsets[0].(mapset.Set))][i] = dist
	}
	for m := 1; m < n-1; m++ {
		fmt.Printf("\n\nInitiating subsets of size %d\n", m)
		cardinalSets := subset.FilterByCardinality(subsets, m)
		newCache := subset.GenerateCache(cardinalSets, n)
		lenSets := len(cardinalSets)
		for i, ss := range cardinalSets {
			fmt.Printf("\rProcessing subset %d of %d", i+1, lenSets)
			// iterate over each item in the subset
			for k := range ss.(mapset.Set).Iter() {
				for j := 1; j < n; j++ {
					if k.(int) == j {
						continue
					}
					kj := distances[k.(int)][j]
					sPrime := ss.(mapset.Set).Clone()
					sPrime.Remove(k)
					cv, ok := oldCache[subset.Hash(sPrime)]
					if !ok {
						log.Fatalf("Error looking up set %s in cache", subset.Hash(sPrime))
					}
					sk := cv[j]
					if sk+kj < newCache[subset.Hash(ss.(mapset.Set))][j] {
						newCache[subset.Hash(ss.(mapset.Set))][j] = sk + kj
					}
				}
			}
		}
		oldCache = newCache
	}

	result := maxFloat
	for k := range subsets[len(subsets)-1].(mapset.Set).Iter() {
		kj := distances[k.(int)][0]
		sPrime := subsets[len(subsets)-1].(mapset.Set).Clone()
		sPrime.Remove(k)
		cv, ok := oldCache[subset.Hash(sPrime)]
		if !ok {
			log.Fatalf("Error looking up set %s in cache", sPrime)
		}
		sk := cv[k.(int)]
		if sk+kj < result {
			result = sk + kj
		}
	}
	fmt.Printf("\nResult: %f\n", result)
}

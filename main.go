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

	"github.com/rhinodavid/bitset"
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

func generateCache(subsets []bitset.Bitset, n int) map[bitset.Bitset]map[int]float32 {
	r := make(map[bitset.Bitset]map[int]float32)
	for _, sS := range subsets {
		r[sS] = make(map[int]float32)

	}
	return r
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

	// compute distance matrix for reuse
	distances := make([][]float32, n)
	for i := 0; i < n; i++ {
		distances[i] = make([]float32, n)
		for j := 0; j < n; j++ {
			distances[i][j] = cartesianDist(a[i], a[j])
		}
	}
	log.Printf("Finsihed generating distances\n")

	// generate all subset combinations
	indexes := []int{}
	for i := 1; i < n; i++ {
		indexes = append(indexes, i)
	}
	fullSet := bitset.NewFromSlice(indexes)
	subsets := fullSet.PowerSet()
	log.Printf("Finished generating subsets\n")

	// build and initialize first cache
	oldCache := generateCache(subsets[0], n)
	for i := 1; i < n; i++ {
		dist := distances[0][i]
		for key := range oldCache {
			// only 1 key
			oldCache[key][i] = dist
		}
	}

	// iterate over lengths of subsets m
	for m := 1; m < n-1; m++ {
		newCache := generateCache(subsets[m], n)
		fmt.Printf("\rInitiating subsets of size %d", m)
		i := 1
		for ss := range newCache {
			i++
			// iterate over each item in the subset
			for _, k := range ss.ToSlice() {
				for j := 1; j < n; j++ {
					if ss.Contains(j) {
						continue
					} else {
						if _, ok := newCache[ss][j]; !ok {
							newCache[ss][j] = maxFloat
						}
					}
					kj := distances[k][j]
					sPrime := ss.RemoveMember(k)
					cv, ok := oldCache[sPrime]
					if !ok {
						log.Fatalf("Error looking up set %v in cache", sPrime)
					}
					sk := cv[k]
					if sk+kj < newCache[ss][j] {
						newCache[ss][j] = sk + kj
					}
				}
			}
		}
		oldCache = newCache
	}

	// compute final subset
	fmt.Println()
	log.Printf("Computing final result from last cache\n")
	intSlice := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		intSlice[i] = i + 1
	}
	finalSS := bitset.NewFromSlice(intSlice)

	// compute result from last cache
	result := maxFloat
	for _, k := range finalSS.ToSlice() {
		kj := distances[k][0]
		sPrime := finalSS.RemoveMember(k)
		cv, ok := oldCache[sPrime]
		if !ok {
			log.Fatalf("Error looking up set %v in cache", sPrime)
		}
		sk := cv[k]
		if sk+kj < result {
			result = sk + kj
		}
	}
	fmt.Printf("\nResult: %f\n", result)
}

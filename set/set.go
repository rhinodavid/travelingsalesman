package subset

import (
	"fmt"
	"sort"

	mapset "github.com/deckarep/golang-set"
	"github.com/rhinodavid/bitset"
)

func Hash(s mapset.Set) string {
	a := s.ToSlice()
	sort.Slice(a, func(x, y int) bool { return a[x].(int) < a[y].(int) })
	return fmt.Sprintf("%v", a)
}

func GenerateSubsets(n int) []interface{} {
	s := mapset.NewSet()
	for i := 1; i < n; i++ {
		s.Add(i)
	}
	ps := s.PowerSet().ToSlice()
	sort.Slice(ps, func(i, j int) bool { return ps[i].(mapset.Set).Cardinality() < ps[j].(mapset.Set).Cardinality() })
	return ps
}

func FilterByCardinality(sets []interface{}, n int) []interface{} {
	r := make([]interface{}, 0)
	for _, v := range sets {
		if v.(mapset.Set).Cardinality() == n {
			r = append(r, v)
		}
	}
	return r
}

func GenerateCache(subsets []interface{}, n int) map[bitset.Bitset]map[int]float32 {
	r := make(map[bitset.Bitset]map[int]float32)
	for _, sS := range subsets {
		interfaceSlice := sS.(mapset.Set).ToSlice()
		intSlice := make([]int, len(interfaceSlice))
		for i, v := range interfaceSlice {
			intSlice[i] = v.(int)
		}
		r[bitset.NewFromSlice(intSlice)] = make(map[int]float32)

	}
	return r
}

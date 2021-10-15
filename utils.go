package main

import (
	"math"
	"strconv"
)

func toInt(ss []string) []int {
	ret := make([]int, 0, len(ss))
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		ret = append(ret, i)
	}
	return ret
}

func toFloat(ss []string) []float64 {
	ret := make([]float64, 0, len(ss))
	for _, s := range ss {
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		ret = append(ret, i)
	}
	return ret
}

// integer absolute value
func iabs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// value of a with the sign of b
func isign(a, b int) int {
	if b < 0 {
		return -a
	}
	return a
}

func toDeg(a float64) float64 {
	return a * 180.0 / math.Pi
}

func toRad(a float64) float64 {
	return a * math.Pi / 180.0
}

// set s[i] to v, extending the slice if necessary
func extend(s []float64, i int, v float64) []float64 {
	for i >= len(s) {
		s = append(s, 0)
	}
	s[i] = v
	return s
}

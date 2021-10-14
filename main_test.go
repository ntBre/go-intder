package main

import (
	"reflect"
	"testing"
)

func TestReadInput(t *testing.T) {
	gconf, gsiics, gsyics, gcart := ReadInput("intder.in")
	wconf := Config{3, 3, 3, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 14}
	wsiics := []Siic{
		{"STRE", []int{1, 2}},
		{"STRE", []int{2, 3}},
		{"BEND", []int{1, 2, 3}},
	}
	wysics := [][]int{
		{1, 2},
		{3},
		{1, -2},
	}
	wcart := []float64{
		0.000000000, 1.431390207, 0.986041184,
		0.000000000, 0.000000000, -0.124238453,
		0.000000000, -1.431390207, 0.986041184,
	}
	if !reflect.DeepEqual(gconf, wconf) {
		t.Errorf("got\n%v, wanted\n%v\n", gconf, wconf)
	}
	if !reflect.DeepEqual(gsiics, wsiics) {
		t.Errorf("got\n%v, wanted\n%v\n", gsiics, wsiics)
	}
	if !reflect.DeepEqual(gsyics, wysics) {
		t.Errorf("got\n%v, wanted\n%v\n", gsyics, wysics)
	}
	if !reflect.DeepEqual(gcart, wcart) {
		t.Errorf("got\n%v, wanted\n%v\n", gcart, wcart)
	}
}

func TestSiICVals(t *testing.T) {
	_, siics, _, carts := ReadInput("intder.in")
	got := SiICVals(siics, carts)
	want := []float64{
		0.9586143064, 0.9586143064, 104.4010205969,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

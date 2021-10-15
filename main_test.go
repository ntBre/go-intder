package main

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestReadInput(t *testing.T) {
	gconf, gsiics, gsyics, gcart, gdisp := ReadInput("intder.in")
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
	wdisp := [][]float64{
		{-0.005, -0.005, -0.01}, {-0.005, -0.005}, {-0.005, -0.005, 0.01},
		{-0.005, -0.01}, {-0.005, -0.015}, {-0.005, 0, -0.01}, {-0.005},
		{-0.005, 0, 0.01}, {-0.005, 0.005, -0.01}, {-0.005, 0.005},
		{-0.005, 0.005, 0.01}, {-0.005, 0.01}, {-0.005, 0.015},
		{-0.01, -0.005}, {-0.01, -0.01}, {-0.01, 0, -0.01}, {-0.01},
		{-0.01, 0, 0.01}, {-0.01, 0.005}, {-0.01, 0.01}, {-0.015, -0.005},
		{-0.015}, {-0.015, 0.005}, {-0.02}, {0, -0.005, -0.01},
		{0, -0.005}, {0, -0.005, 0.01}, {0, -0.01, -0.01},
		{0, -0.01}, {0, -0.01, 0.01}, {0, -0.015}, {0, -0.02},
		{0, 0, -0.01}, {0, 0, -0.02}, {}, {0, 0, 0.01}, {0, 0, 0.02},
		{0, 0.005, -0.01}, {0, 0.005}, {0, 0.005, 0.01}, {0, 0.01, -0.01},
		{0, 0.01}, {0, 0.01, 0.01}, {0, 0.015}, {0, 0.02},
		{0.005, -0.005, -0.01}, {0.005, -0.005}, {0.005, -0.005, 0.01},
		{0.005, -0.01}, {0.005, -0.015}, {0.005, 0, -0.01}, {0.005},
		{0.005, 0, 0.01}, {0.005, 0.005, -0.01}, {0.005, 0.005},
		{0.005, 0.005, 0.01}, {0.005, 0.01}, {0.005, 0.015},
		{0.01, -0.005}, {0.01, -0.01}, {0.01, 0, -0.01}, {0.01},
		{0.01, 0, 0.01}, {0.01, 0.005}, {0.01, 0.01}, {0.015, -0.005},
		{0.015}, {0.015, 0.005}, {0.02},
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
	if !reflect.DeepEqual(gdisp, wdisp) {
		for _, i := range gdisp {
			fmt.Print("{")
			for _, v := range i {
				fmt.Printf("%v, ", v)
			}
			fmt.Print("},\n")
		}
		t.Errorf("got\n%v, wanted\n%v\n", gdisp, wdisp)
	}
}

func nearby(a, b []float64, eps float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > eps {
			return false
		}
	}
	return true
}

func TestSiICVals(t *testing.T) {
	_, siics, _, carts, _ := ReadInput("intder.in")
	got := SiICVals(siics, carts)
	want := []float64{
		0.9586143064, 0.9586143064, toRad(104.4010205969),
	}
	if !nearby(got, want, 1e-10) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestSyICVals(t *testing.T) {
	_, siics, syics, carts, _ := ReadInput("intder.in")
	sics := SiICVals(siics, carts)
	got := SyICVals(syics, sics)
	want := []float64{
		1.3556853532, 1.8221415519, 0.0000000000,
	}
	if !nearby(got, want, 1e-10) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

func TestDisp(t *testing.T) {
	// conf, siics, syics, carts, disps := ReadInput("intder.in")
	_, siics, syics, carts, disps := ReadInput("intder.in")
	sics := SyICVals(syics, SiICVals(siics, carts))
	got := Disp(siics, sics, carts, disps[0:1])[0]
	want := []float64{
		0.0000000000, 1.4186597974, 0.9822041564,
		0.0000000000, 0.0094006500, -0.1238566934,
		0.0000000000, -1.4280604475, 0.9894964520,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v\n", got, want)
	}
}

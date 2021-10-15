package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Input options, see the manual for additional details
type IOPT uint8

const (
	NA     IOPT = iota // number of atoms
	NS                 // number of simple internal coordinates
	NSYM               // number of symmetry internal coordinates
	NDER               // derivative level
	NEQ                // stationary point
	NPRT               // print option
	NINV               // coordinate transformation
	NDUM               // number of dummy atoms
	NTEST              // numerical derivative testing
	NGEOM              // source of Cartesian geometry
	NFREQ              // coordinate system for frequency analysis
	IRINT              // compute IR intensities
	NVEC               // dimension of property of derivative transform
	NSTOP              // stop after forming some matrices
	NDISP              // coordinate system for displacements
	NMODE              // assign normal modes
	THRESH             // threshold for displacement convergence (10^-THRESH)
)

// Physical constants
const (
	BOHR  = 0.529177249e0
	DEBYE = 2.54176548e0
	HART  = 4.3597482e0
	WAVE0 = 1302.7910e0
	CINT  = 42.25472e0
)

const (
	// LIN includes LIN1, LINX, LINY
	COORD_TYPES = "STRE|BEND|LIN|TORS|OUT|SPF|RCOM"
)

type Config []int

type Siic struct {
	Type  string
	Atoms []int
}

func (s Siic) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "%5s", s.Type)
	for _, d := range s.Atoms {
		fmt.Fprintf(&str, "%5d", d)
	}
	return str.String()
}

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

// ReadInput reads an intder input file. TODO handle freqs input
func ReadInput(filename string) (conf Config, siics []Siic, syics [][]int,
	carts []float64) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var (
		handler func([]string)
		line    string
		fields  []string
	)
	cartHandler := func(s []string) {
		if strings.Contains(line, "DISP") {
			handler = nil
			return
		}
		carts = append(carts, toFloat(s)...)
	}
	syicHandler := func(s []string) {
		if strings.Contains(line, " 0") {
			handler = cartHandler
			return
		}
		tmp := make([]int, 0)
		for i := 1; i < len(s)-1; i += 2 {
			f, err := strconv.ParseFloat(s[i+1], 64)
			if err != nil {
				panic(err)
			}
			d, err := strconv.Atoi(s[i])
			if err != nil {
				panic(err)
			}
			tmp = append(tmp, int(f)*d)
		}
		syics = append(syics, tmp)
	}
	siicHandler := func(s []string) {
		if !strings.Contains(COORD_TYPES, s[0]) {
			handler = syicHandler
			handler(s)
			return
		}
		siics = append(siics, Siic{s[0], toInt(s[1:])})
	}
	inpHandler := func(s []string) {
		conf = toInt(s)
		handler = siicHandler
	}
	for scanner.Scan() {
		line = scanner.Text()
		fields = strings.Fields(line)
		switch {
		case strings.Contains(line, "INTDER"):
			handler = inpHandler
		case handler != nil:
			handler(fields)
		}
	}
	return
}

// At accesses matr as if it were a nx3 matrix
func At(matr []float64, i, j int) float64 {
	return matr[3*i+j]
}

// Stre computes the distance between a and b
func Stre(a, b []float64) float64 {
	return Mag(Sub(a, b))
}

// Sub returns the vector distance between a and b
func Sub(a, b []float64) (ret []float64) {
	for i := range a {
		ret = append(ret, a[i]-b[i])
	}
	return ret
}

// Dot returns the dot product of a and b
func Dot(a, b []float64) (dot float64) {
	for i := range a {
		dot += a[i] * b[i]
	}
	return
}

// Mag returns the magnitude of a
func Mag(a []float64) (mag float64) {
	for _, v := range a {
		v *= v
		mag += v
	}
	return math.Sqrt(mag)
}

func toDeg(a float64) float64 {
	return a * 180.0 / math.Pi
}

func toRad(a float64) float64 {
	return a * math.Pi / 180.0
}

// SiICVals computes the values of the simple internal coordinates at
// a given cartesian geometry
func SiICVals(siics []Siic, carts []float64) (ret []float64) {
	for _, s := range siics {
		switch s.Type {
		case "STRE":
			a, b := 3*(s.Atoms[0]-1), 3*(s.Atoms[1]-1)
			d := Stre(carts[a:a+3], carts[b:b+3])
			ret = append(ret, d*BOHR)
		case "BEND":
			// b is the central atom
			a, b, c := 3*(s.Atoms[0]-1), 3*(s.Atoms[1]-1), 3*(s.Atoms[2]-1)
			// cosΘ = ba·bc / ||ba||×||bc||
			ba := Sub(carts[b:b+3], carts[a:a+3])
			bc := Sub(carts[b:b+3], carts[c:c+3])
			ret = append(ret, math.Acos(Dot(ba, bc)/(Mag(ba)*Mag(bc))))
			// TODO the rest of the coordinate types
		default:
			panic("unrecognized internal coordinate type")
		}
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

func SyICVals(syics [][]int, siics []float64) (ret []float64) {
	var f, sum float64
	for _, s := range syics {
		sum = 0
		switch len(s) {
		case 1:
			f = 1
		case 2:
			f = 1 / math.Sqrt2
		default:
			panic("unrecognized internal coordinate type")
		}
		for _, t := range s {
			sum += float64(isign(1, t)) * f * siics[iabs(t)-1]
		}
		ret = append(ret, sum)
	}
	return
}

func main() {
	conf, siics, syics, carts := ReadInput("intder.in")
	fmt.Println(conf, siics, syics, carts)
}

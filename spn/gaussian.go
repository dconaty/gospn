package spn

/*
#cgo LDFLAGS: -lgsl -lgslcblas
#include <gsl/gsl_randist.h>
*/
import "C"

import (
	//"fmt"
	"math"
)

const (
	// GaussMax is the maximum value of a standard gaussian, namely 1/sqrt(2*pi).
	GaussMax = 0.398942280 // 1/sqrt(2*pi) := max value of a standard Gaussian
)

// Gaussian represents a gaussian distribution.
type Gaussian struct {
	Node
	// Variable ID
	varid int
	// Mean
	mean float64
	// Standard deviation
	sd float64
}

// NewGaussian constructs a new Gaussian from a counting slice.
func NewGaussian(varid int, counts []int) *Gaussian {
	var mean, sd float64
	var N int
	n := len(counts)

	// Standardizing gaussian from N(mean, sd) to N(0, 1).

	for i := 0; i < n; i++ {
		mean += float64(counts[i] * i)
		N += counts[i]
	}

	mean /= float64(N)
	//fmt.Printf("Mean: %.5f, ", mean)

	for i := 0; i < n; i++ {
		d := float64(i) - mean
		sd += float64(counts[i]) * d * d
	}
	sd = math.Sqrt(sd)
	if sd == 0 {
		sd = 1
	}

	//fmt.Printf("StdDev: %.5f\n", sd)

	return &Gaussian{Node{sc: []int{varid}}, varid, mean, sd}
}

// Type returns the type of this node.
func (g *Gaussian) Type() string { return "leaf" }

// Value returns the probability of a certain valuation. That is Pr(X=val[varid]), where
// Pr is a probability function over a gaussian distribution.
func (g *Gaussian) Value(val VarSet) float64 {
	v, ok := val[g.varid]
	if ok {
		//fmt.Println("Yelloooo")
		return math.Log(float64(C.gsl_ran_ugaussian_pdf(C.double((float64(v) - g.mean) / g.sd))))
	}
	return 0.0 // ln(1.0) = 0.0
}

// Max returns the MAP state given a valuation.
func (g *Gaussian) Max(val VarSet) float64 {
	v, ok := val[g.varid]
	if ok {
		//fmt.Printf("Preparing Gaussian with parameters z = (x-mu)/sigma = (%.3f - %.3f) / %.3f\n",
		//float64(v), g.mean, g.sd)
		z := math.Log(float64(C.gsl_ran_ugaussian_pdf(C.double((float64(v) - g.mean) / g.sd))))
		//fmt.Printf("Max: Gaussian[%d] = %.5f\n", g.varid, z)
		return z
	}
	return math.Log(GaussMax)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (g *Gaussian) ArgMax(val VarSet) (VarSet, float64) {
	retval := make(VarSet)
	v, ok := val[g.varid]

	if ok {
		retval[g.varid] = v
		//fmt.Printf("Preparing Gaussian with parameters z = (x-mu)/sigma = (%.3f - %.3f) / %.3f\n",
		//float64(v), g.mean, g.sd)
		z := math.Log(float64(C.gsl_ran_ugaussian_pdf(C.double((float64(v) - g.mean) / g.sd))))
		//fmt.Printf("ArgMax: Gaussian[%d] = %.5f\n", g.varid, z)
		return retval, z
	}

	retval[g.varid] = int(g.mean)
	return retval, math.Log(GaussMax)
}

package spn

import (
	"fmt"

	utils "github.com/RenatoGeh/gospn/src/utils"
)

// Mode of a univariate distribution.
type Mode struct {
	// Value of variable when it is the highest.
	index int
	// Highest value of variable.
	val float64
}

// A univariate distribution is a probability distribution with unary scope.
// UnivDist actually represents a multinomial distribution.
type UnivDist struct {
	// Parent node.
	pa Node
	// Variable ID
	varid int
	// Discrete probability distribution
	pr []float64
	// Mode of pr. We pre-compute this to save time.
	mode Mode
	// Scope of this univariate distribution. We store this to avoid new creations of slices.
	sc []int
}

// Constructs a new UnivDist.
func NewUnivDist(varid int, dist []float64) *UnivDist {
	n := len(dist)
	var m float64 = 0
	var mi int = 0

	for i := 0; i < n; i++ {
		if dist[i] > m {
			m = dist[i]
			mi = i
		}
	}

	return &UnivDist{nil, varid, dist, Mode{mi, m}, []int{varid}}
}

func NewCountingUnivDist(varid int, counts []int) *UnivDist {
	n := len(counts)

	pr := make([]float64, n)
	s := 0.0
	for i := 0; i < n; i++ {
		s += 1.0 + float64(counts[i])
		pr[i] = float64(1 + counts[i])
	}

	for i := 0; i < n; i++ {
		pr[i] /= float64(s)
	}

	var m float64 = 0
	var mi int = 0

	for i := 0; i < n; i++ {
		if pr[i] > m {
			m = pr[i]
			mi = i
		}
	}

	return &UnivDist{nil, varid, pr, Mode{mi, m}, []int{varid}}
}

// Constructs a new empty UnivDist for learning. We initialize pr to a uniform distribution.
// Argument m is the cardinality of varid.
func NewEmptyUnivDist(varid, m int) *UnivDist {
	pr := make([]float64, m)

	for i := 0; i < m; i++ {
		pr[i] = 1.0 / float64(m)
	}

	return &UnivDist{nil, varid, pr, Mode{0, pr[0]}, []int{varid}}
}

// Ch returns the set of childre nodes. Since a node is a UnivDist iff it is a leaf, Ch=\emptyset.
func (ud *UnivDist) Ch() []Node { return nil }

// Pa returns the parent node.
func (ud *UnivDist) Pa() Node { return ud.pa }

// Type return this node's type: 'leaf'.
func (ud *UnivDist) Type() string { return "leaf" }

// Sc returns this node's scope.
func (ud *UnivDist) Sc() []int {
	// A univariate distribution has unary scope by definition.
	return ud.sc
}

// Sets the parent node.
func (ud *UnivDist) SetParent(pa Node) { ud.pa = pa }

// Returns nil. Leaves have no weights.
func (ud *UnivDist) Weights() []float64 { return nil }

// Adds a child, but actually doesn't since it's a leaf.
func (ud *UnivDist) AddChild(c Node) {}

// Returns the probability of a certain valuation. That is Pr(X=valuation[varid]), where
// Pr=UnivDist.
func (ud *UnivDist) Value(valuation VarSet) float64 {
	val, ok := valuation[ud.varid]
	if ok {
		fmt.Printf("Value of leaf node: %f\n", ud.pr[val])
		return utils.Log(ud.pr[val])
	}
	fmt.Printf("Value of leaf node: 1.00\n")

	//	return 1.0
	return 0.0 // ln(1.0) = 0.0
}

// Max returns the MAP state given a valuation.
func (ud *UnivDist) Max(valuation VarSet) float64 {
	val, ok := valuation[ud.varid]
	if ok {
		return utils.Log(ud.pr[val])
	}
	return utils.Log(ud.mode.val)
}

// ArgMax returns both the arguments and the value of the MAP state given a certain valuation.
func (ud *UnivDist) ArgMax(valuation VarSet) (VarSet, float64) {
	retval := make(VarSet)
	val, ok := valuation[ud.varid]

	if ok {
		retval[ud.varid] = val
		return retval, utils.Log(ud.pr[val])
	}

	retval[ud.varid] = ud.mode.index
	return retval, utils.Log(ud.mode.val)
}

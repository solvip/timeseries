// Package timeseries provides utilities to manipulate and analyze timeseries data.
// For compatability with Gonum, a Timeseries is simply a pair of float64 slices,
// the X and the Y axis.
// You can manipulate them as you wish, but ensure that you call Sort() if you're
// not appending data points in sorted order.
package timeseries

import (
	"sort"
)

type Timeseries struct {
	Xs []float64
	Ys []float64
}

// Equal - Return true if t and other represent the same time series
func (t Timeseries) Equal(other Timeseries) bool {
	if len(t.Xs) != len(t.Ys) || len(other.Xs) != len(other.Ys) {
		panic("timeseries: Xs and Ys length mismatch")
	}

	if t.Len() != other.Len() {
		return false
	}

	for i := 0; i < t.Len(); i++ {
		if t.Xs[i] != other.Xs[i] || t.Ys[i] != other.Ys[i] {
			return false
		}
	}

	return true
}

// After - Return a shallow copy of the items in the time series having Xs >= x
// The series must be sorted.
func (t Timeseries) After(x float64) Timeseries {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys length mismatch")
	}

	if i := t.findPivot(x); i == t.Len() {
		// After is older than all the items in the series
		return Timeseries{}
	} else {
		return Timeseries{
			Xs: t.Xs[i:],
			Ys: t.Ys[i:],
		}
	}
}

// Before - Return a shallow copy of the items in the time series having Xs < x.
// The series must be sorted.
func (t Timeseries) Before(x float64) Timeseries {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys length mismatch")
	}

	if j := t.findPivot(x); j == t.Len() {
		return t
	} else {
		return Timeseries{
			Xs: t.Xs[:j],
			Ys: t.Ys[:j],
		}
	}
}

// Between - Return a shallow copy of the items in the time series between [x1, x2)
func (t Timeseries) Between(x1, x2 float64) Timeseries {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys length mismatch")
	}

	return t.After(x1).Before(x2)
}

// findPivot - Binary search for the location of x in t and return its index,
// where the index will put i at before <= x < after
func (t Timeseries) findPivot(x float64) int {
	findAfter := func(i int) bool {
		return t.Xs[i] >= x
	}

	return sort.Search(t.Len(), findAfter)
}

// Append - Append value @ time to the timeseries
// Note that you might need a sort if you're inserting points out-of-order
func (t *Timeseries) Append(x float64, y float64) {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys slice length mismatch")
	}

	t.Xs = append(t.Xs, x)
	t.Ys = append(t.Ys, y)
}

// Difference the timeseries N, returning a new series of length len(N)-1
func (t Timeseries) Difference() (ret Timeseries) {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys slice length mismatch")
	}

	if t.Len() < 2 {
		// We must have at least two elements to difference
		return ret
	}

	ret = makeTimeseries(t.Len() - 1)
	for i := 0; i < ret.Len(); i++ {
		ret.Ys[i] = t.Ys[i+1] - t.Ys[i]
		ret.Xs[i] = t.Xs[i+1]
	}

	return ret
}

func makeTimeseries(length int) Timeseries {
	return Timeseries{
		Xs: make([]float64, length),
		Ys: make([]float64, length),
	}
}

// Slice slices the Timeseries equivalently to t[start:end]
func (t Timeseries) Slice(start, end int) Timeseries {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys slice length mismatch")
	}

	return Timeseries{
		Xs: t.Xs[start:end],
		Ys: t.Ys[start:end],
	}
}

func (t Timeseries) Sort() {
	if len(t.Xs) != len(t.Ys) {
		panic("timeseries: Xs and Ys slice length mismatch")
	}

	sort.Sort(t)
}

func (t Timeseries) Len() int {
	if n := len(t.Xs); n != len(t.Ys) {
		panic("timeseries: Xs and Ys slice length mismatch")
	} else {
		return n
	}
}

func (t Timeseries) Swap(i, j int) {
	t.Xs[i], t.Xs[j] = t.Xs[j], t.Xs[i]
	t.Ys[i], t.Ys[j] = t.Ys[j], t.Ys[i]
}
func (t Timeseries) Less(i, j int) bool {
	return t.Xs[i] < t.Xs[j]
}

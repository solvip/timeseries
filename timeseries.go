// Package timeseries provides utilities to manipulate and analyze timeseries data
// For ease of use by clients, the Timeseries is simply a slice of points.
// This, however, means that you might need to call Sort() if you're not appending
// points in sorted order.
package timeseries

import (
	"sort"
	"time"
)

type Timeseries []Point

type Point struct {
	Time  time.Time
	Value float64
}

// Equal - Return true if t and other represent the same time series
func (t Timeseries) Equal(other Timeseries) bool {
	if len(t) != len(other) {
		return false
	}

	for i := range t {
		if t[i] != other[i] {
			return false
		}
	}

	return true
}

// After - Return a shallow copy of the items in the time series having Time >= instant
// The series must be sorted.
func (t Timeseries) After(instant time.Time) Timeseries {
	if i := t.findPivot(instant); i == len(t) {
		// After is older than all the items in the series
		return Timeseries{}
	} else {
		return t[i:]
	}
}

// Before - Return a shallow copy of the items in the time series having Time < instant.
// The series must be sorted.
func (t Timeseries) Before(instant time.Time) Timeseries {
	j := t.findPivot(instant)
	return t[:j]
}

// Between - Return a shallow copy of the items in the time series between [from, to)
func (t Timeseries) Between(from, to time.Time) Timeseries {
	return t.After(from).Before(to)
}

// findPivot - Binary search for the location of instant in t and return its index,
// where the index is before <= instant < after
func (t Timeseries) findPivot(instant time.Time) int {
	findAfter := func(i int) bool {
		return t[i].Time.Equal(instant) || t[i].Time.After(instant)
	}

	return sort.Search(len(t), findAfter)
}

// Append - Append the Point p to the timeseries
// Note that you might need a sort if you're inserting points out-of-order
func (t *Timeseries) AppendPoint(p Point) {
	*t = append(*t, p)
}

// Append - Append value @ time to the timeseries
// Note that you might need a sort if you're inserting points out-of-order
func (t *Timeseries) Append(time time.Time, value float64) {
	t.AppendPoint(Point{Time: time, Value: value})
}

// Difference the timeseries N, returning a new series of length len(N)-1
func (t Timeseries) Difference() (ret Timeseries) {
	if t.Len() < 2 {
		// We must have at least two elements to difference
		return ret
	}

	ret = make(Timeseries, len(t)-1)
	for i := range ret {
		ret[i] = Point{
			Time:  t[i+1].Time,
			Value: t[i+1].Value - t[i].Value,
		}
	}

	return ret
}

func (t Timeseries) Sort() {
	sort.Sort(t)
}

func (t Timeseries) Len() int           { return len(t) }
func (t Timeseries) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Timeseries) Less(i, j int) bool { return t[i].Time.Before(t[j].Time) }

package timeseries

import (
	"math"
	"testing"
)

var emptyTimeseries = Timeseries{}
var minY = math.Inf(-1)
var maxY = math.Inf(1)

func TestAppend(t *testing.T) {
	expected := Timeseries{
		Xs: []float64{123.4, 456.7},
		Ys: []float64{1, 2},
	}

	var actual Timeseries
	actual.Append(expected.Xs[0], expected.Ys[0])
	actual.Append(expected.Xs[1], expected.Ys[1])

	if !actual.Equal(expected) {
		t.Fatalf("expected %v after appending; instead got %v", expected, actual)
	}
}

func TestAfter(t *testing.T) {
	ts := Timeseries{
		Xs: []float64{100.0, 50.0, 100.0},
		Ys: []float64{1, 2, 3},
	}

	if x := emptyTimeseries.After(minY); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected emptyTimeseries.After(minY) to return empty series; instead got %v", x)
	}

	if x := ts.After(maxY); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected no items in time series after maximum y, instead got: %v", x)
	}

	if x := ts.After(0); !ts.Equal(ts) {
		t.Fatalf("expected ts.After(0) to return ts; instead got: %v", x)
	}

	// After is inclusive
	expected := ts.Slice(1, ts.Len())
	if x := ts.After(2); !x.Equal(expected) {
		t.Fatalf("expected ts.After(2) to return %v; instead got %v", expected, x)
	}
}

func TestBefore(t *testing.T) {
	ts := Timeseries{
		Xs: []float64{100.0, 50.0, 100.0},
		Ys: []float64{1, 2, 3},
	}

	if x := emptyTimeseries.Before(math.Inf(-1)); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected nil to be returned on emptyTimeseries.Before(maxTime); instead got %v", x)
	}

	if x := ts.Before(minY); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected no items in time series before zero time, instead got %v", x)
	}

	if x := ts.Before(maxY); !x.Equal(ts) {
		t.Fatalf("expected ts.Before(maxTime) to return ts; instead got: %v", x)
	}

	// Before is exclusive
	expected := ts.Slice(0, 2)
	if x := ts.Before(ts.Ys[2]); !x.Equal(expected) {
		t.Fatalf("expected ts.Before(%v) to return %v; instead got %v", ts.Ys[2], expected, x)
	}
}

func TestBetween(t *testing.T) {
	ts := Timeseries{
		Xs: []float64{100.0, 50.0, 100.0},
		Ys: []float64{1, 2, 3},
	}

	if x := ts.Between(minY, minY); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected ts.Between(minY, minY) to return empty; instead got %v", x)
	}

	// We return in the range [from, to).  We should only get the first two elements.
	firstTwo := ts.Slice(0, 2)
	if x := ts.Between(minY, ts.Ys[2]); !x.Equal(firstTwo) {
		t.Fatalf("expected ts.Between(minY, %v) to return %v; instead got %v", ts.Ys[2], firstTwo, x)
	}

	// we return in the range [from, to).  Ensure that we only get the middle
	// element
	y1 := ts.Ys[1]
	y2 := ts.Ys[2]
	middle := ts.Slice(1, 2)
	if x := ts.Between(y1, y2); !x.Equal(middle) {
		t.Fatalf("expected ts.Between(%v, %v) to return %v; instead got %v", y1, y2, middle, x)
	}
}

func TestDifference(t *testing.T) {
	if x := emptyTimeseries.Difference(); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected difference of empty series to return empty series; instead got %v", x)
	}

	ts := Timeseries{
		Xs: []float64{100.0, 50.0, 100.0},
		Ys: []float64{1, 2, 3},
	}

	if x := ts.Slice(1, 2).Difference(); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected difference of 1-length series to return empty series; instead got %v", x)
	}

	actual := ts.Difference()
	expected := Timeseries{
		Xs: []float64{-50.0, 50.0},
		Ys: []float64{2, 3},
	}

	if !actual.Equal(expected) {
		t.Fatalf("expected ts.Difference() to return %v; instead got %v", expected, actual)
	}
}

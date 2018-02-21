package timeseries

import (
	"math"
	"testing"
)

var emptyTimeseries = Timeseries{}
var minX = math.Inf(-1)
var maxX = math.Inf(1)

func TestAppend(t *testing.T) {
	expected := Timeseries{
		Xs: []float64{1, 2},
		Ys: []float64{123.4, 456.7},
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
		Xs: []float64{1, 2, 3},
		Ys: []float64{100.0, 50.0, 100.0},
	}

	if actual := emptyTimeseries.After(minX); !actual.Equal(emptyTimeseries) {
		t.Fatalf("expected emptyTimeseries.After(minX) to return empty series; instead got %v", actual)
	}

	if actual := ts.After(maxX); !actual.Equal(emptyTimeseries) {
		t.Fatalf("expected no items in time series after maximum y, instead got: %v", actual)
	}

	if actual := ts.After(0); !ts.Equal(ts) {
		t.Fatalf("expected ts.After(0) to return ts; instead got: %v", actual)
	}

	// After is inclusive
	expected := ts.Slice(1, ts.Len())
	if actual := ts.After(2); !actual.Equal(expected) {
		t.Fatalf("expected ts.After(2) to return %v; instead got %v", expected, actual)
	}
}

func TestBefore(t *testing.T) {
	ts := Timeseries{
		Xs: []float64{1, 2, 3},
		Ys: []float64{100.0, 50.0, 100.0},
	}

	if x := emptyTimeseries.Before(math.Inf(-1)); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected nil to be returned on emptyTimeseries.Before(maxTime); instead got %v", x)
	}

	if x := ts.Before(minX); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected no items in time series before zero time, instead got %v", x)
	}

	if x := ts.Before(maxX); !x.Equal(ts) {
		t.Fatalf("expected ts.Before(maxTime) to return ts; instead got: %v", x)
	}

	// Before is exclusive
	expected := ts.Slice(0, 2)
	if actual := ts.Before(ts.Xs[2]); !actual.Equal(expected) {
		t.Fatalf("expected ts.Before(%v) to return %v; instead got %v", ts.Xs[2], expected, actual)
	}
}

func TestBetween(t *testing.T) {
	ts := Timeseries{
		Xs: []float64{1, 2, 3},
		Ys: []float64{100.0, 50.0, 100.0},
	}

	if actual := ts.Between(minX, minX); !actual.Equal(emptyTimeseries) {
		t.Fatalf("expected ts.Between(minX, minX) to return empty; instead got %v", actual)
	}

	// We return in the range [from, to).  We should only get the first two elements.
	firstTwo := ts.Slice(0, 2)
	if x := ts.Between(minX, ts.Xs[2]); !x.Equal(firstTwo) {
		t.Fatalf("expected ts.Between(minX, %v) to return %v; instead got %v", ts.Xs[2], firstTwo, x)
	}

	// we return in the range [from, to).  Ensure that we only get the middle
	// element
	x1, x2 := ts.Xs[1], ts.Xs[2]
	middle := ts.Slice(1, 2)
	if actual := ts.Between(x1, x2); !actual.Equal(middle) {
		t.Fatalf("expected ts.Between(%v, %v) to return %v; instead got %v", x1, x2, middle, actual)
	}
}

func TestDifference(t *testing.T) {
	if x := emptyTimeseries.Difference(); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected difference of empty series to return empty series; instead got %v", x)
	}

	ts := Timeseries{
		Xs: []float64{1, 2, 3},
		Ys: []float64{100.0, 50.0, 100.0},
	}

	if x := ts.Slice(1, 2).Difference(); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected difference of 1-length series to return empty series; instead got %v", x)
	}

	actual := ts.Difference()
	expected := Timeseries{
		Xs: []float64{2, 3},
		Ys: []float64{-50.0, 50.0},
	}

	if !actual.Equal(expected) {
		t.Fatalf("expected ts.Difference() to return %v; instead got %v", expected, actual)
	}
}

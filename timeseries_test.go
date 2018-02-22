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

func TestLinearRegression(t *testing.T) {
	// Ensure that a diagonal results in a perfect fit
	ts1 := Timeseries{
		Xs: []float64{0, 1, 2, 3, 4, 5},
		Ys: []float64{0, 10, 20, 30, 40, 50},
	}

	if alpha, beta, rmse := ts1.LinearRegression(); alpha != 0 || beta != 10 || rmse != 0 {
		t.Fatalf("Expected alpha=0, beta=1, rmse=0, instead got alpha=%v, beta=%v, rmse=%v",
			alpha, beta, rmse)
	}

	// Ensure that a line results in a perfect fit
	ts2 := Timeseries{
		Xs: []float64{0, 1, 2, 3, 4},
		Ys: []float64{5, 5, 5, 5, 5},
	}

	if alpha, beta, rmse := ts2.LinearRegression(); alpha != 5 || beta != 0 || rmse != 0 {
		t.Fatalf("Expected alpha=5, beta=0, rmse=0, instead got alpha=%v, beta=%v, rmse=%v",
			alpha, beta, rmse)
	}
}

func TestFirstLast(t *testing.T) {
	assertPanic(t, "timeseries: empty timeseries", func() { emptyTimeseries.First() })
	assertPanic(t, "timeseries: empty timeseries", func() { emptyTimeseries.Last() })

	ts1 := Timeseries{
		Xs: []float64{0, 1, 2, 3, 4, 5},
		Ys: []float64{0, 10, 20, 30, 40, 50},
	}

	expectedFirstX, expectedFirstY := 0.0, 0.0
	expectedLastX, expectedLastY := 5.0, 50.0

	if x, y := ts1.First(); x != expectedFirstX || y != expectedFirstY {
		t.Fatalf("expected First() = %v, %v, instead got %v, %v", expectedFirstX, expectedFirstY, x, y)
	}

	if x, y := ts1.Last(); x != expectedLastX || y != expectedLastY {
		t.Fatalf("expected Last() = %v, %v, instead got %v, %v", expectedLastX, expectedLastY, x, y)
	}

}

func TestLen(t *testing.T) {
	assertPanic(t, "timeseries: Xs and Ys slice length mismatch", func() {
		ts := Timeseries{Xs: []float64{1}, Ys: []float64{1, 2}}
		ts.Len()
	})

	if n := emptyTimeseries.Len(); n != 0 {
		t.Fatalf("expected Len() = 0, instead got %v", n)
	}

	ts := Timeseries{
		Xs: []float64{1},
		Ys: []float64{2},
	}
	if n := ts.Len(); n != 1 {
		t.Fatalf("expected Len() = 1, instead got %v", n)
	}
}

// assertPanic - Assert that f panics with expectedPanicMsg
func assertPanic(t *testing.T, expectedPanicMsg string, f func()) {
	t.Helper()

	recoverHandler := func() {
		if r := recover(); r == nil {
			t.Errorf("f did not panic as expected")
		} else if r != expectedPanicMsg {
			t.Errorf("expected panic(%s); instead got panic(%s)", expectedPanicMsg, r)
		}
	}

	defer recoverHandler()
	f()
}

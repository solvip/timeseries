package timeseries

import (
	"math"
	"testing"
)

var emptyTimeseries = Timeseries{}
var mismatchedTimeseries = Timeseries{
	Xs: []float64{1, 2, 3, 4},
	Ys: []float64{5, 6},
}
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

func TestEqual(t *testing.T) {
	assertPanic(t, "timeseries: Xs and Ys length mismatch", func() {
		mismatchedTimeseries.Equal(emptyTimeseries)
	})

	assertPanic(t, "timeseries: Xs and Ys length mismatch", func() {
		emptyTimeseries.Equal(mismatchedTimeseries)
	})

	if !emptyTimeseries.Equal(emptyTimeseries) {
		t.Fatalf("expected emptyTimeseries to be equal to emptyTimeseries")
	}

	var ts1, ts2 Timeseries
	ts1.Append(1, 2)
	ts2.Append(1, 2)

	if !ts1.Equal(ts2) {
		t.Fatalf("expected ts1 to be equal to ts2")
	}

	if !ts2.Equal(ts1) {
		t.Fatalf("expected ts2 to be equal to ts1")
	}

	// Test the case where ts1 and ts2 are of non-equal length
	ts1.Append(3, 4)
	if ts1.Equal(ts2) {
		t.Fatalf("ts1 should not be equal to ts2")
	}

	// Test the case where ts1 and ts2 are of equal lenght; but different
	ts2.Append(4, 3)
	if ts1.Equal(ts2) {
		t.Fatalf("ts1 should not be equal to ts2")
	}
}

func TestAfter(t *testing.T) {
	assertPanic(t, "timeseries: Xs and Ys slice length mismatch", func() {
		mismatchedTimeseries.After(0)
	})

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
	assertPanic(t, "timeseries: Xs and Ys slice length mismatch", func() {
		mismatchedTimeseries.Before(0)
	})

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
	assertPanic(t, "timeseries: Xs and Ys slice length mismatch", func() {
		mismatchedTimeseries.Between(0, 1)
	})

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

func TestAt(t *testing.T) {
	ts := Timeseries{
		Xs: []float64{0, 1, 2, 3, 4, 5},
		Ys: []float64{0, 10, 20, 30, 40, 50},
	}

	assertPanic(t, "timeseries: empty timeseries", func() { emptyTimeseries.At(0) })
	assertPanic(t, "timeseries: out of bounds", func() { ts.At(1337) })

	if x, y := ts.At(3); x != 3 || y != 30 {
		t.Fatalf("expected ts.At(3) = 3, 30; instead got %v, %v", x, y)
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

func TestMovingAverage(t *testing.T) {
	assertPanic(t, "timeseries: Xs and Ys slice length mismatch", func() {
		mismatchedTimeseries.MovingAverage(10)
	})

	// A moving average with a window size 1 should be the identity
	ts1 := Timeseries{
		Xs: []float64{1, 2, 3, 4, 5, 6},
		Ys: []float64{1, 2, 4, 8, 16, 32},
	}

	if actual := ts1.MovingAverage(1); !actual.Equal(ts1) {
		t.Fatalf("expected MovingAverage(1) to be the identity of ts1; instead got %v", actual)
	}

	expectedForMA2 := Timeseries{
		Xs: []float64{2, 3, 4, 5, 6},
		Ys: []float64{
			(1.0 + 2.0) / 2,
			(2.0 + 4.0) / 2,
			(4.0 + 8.0) / 2,
			(8.0 + 16.0) / 2,
			(16.0 + 32.0) / 2,
		},
	}
	if actual := ts1.MovingAverage(2); !actual.Equal(expectedForMA2) {
		t.Fatalf("expected MovingAverage(2) to return %v; instead got %v", expectedForMA2, actual)
	}

	expectedForMA4 := Timeseries{
		Xs: []float64{4, 5, 6},
		Ys: []float64{
			(1.0 + 2.0 + 4.0 + 8.0) / 4,
			(2.0 + 4.0 + 8.0 + 16.0) / 4,
			(4.0 + 8.0 + 16.0 + 32.0) / 4,
		},
	}
	if actual := ts1.MovingAverage(4); !actual.Equal(expectedForMA4) {
		t.Fatalf("expected MovingAverage(4) to return %v; instead got %v", expectedForMA4, actual)
	}

}

func TestLen(t *testing.T) {
	assertPanic(t, "timeseries: Xs and Ys slice length mismatch", func() {
		mismatchedTimeseries.Len()
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

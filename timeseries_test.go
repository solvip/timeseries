package timeseries

import (
	"testing"
	"time"
)

var emptyTimeseries = Timeseries{}
var zeroTime = time.Time{}
var maxTime = time.Unix(1<<63-62135596801, 999999999)

func TestAppend(t *testing.T) {
	expected := Timeseries{
		Point{Time: time.Now(), Value: 123.4},
		Point{Time: time.Now(), Value: 456.7},
	}

	var actual Timeseries
	actual.Append(expected[0].Time, expected[0].Value)
	actual.AppendPoint(expected[1])

	if !actual.Equal(expected) {
		t.Fatalf("expected %v after appending; instead got %v", expected, actual)
	}
}

func TestAfter(t *testing.T) {
	ts := Timeseries{
		{time.Date(2018, 2, 15, 22, 41, 31, 0, time.UTC), 100.0},
		{time.Date(2018, 2, 15, 22, 41, 41, 0, time.UTC), 50.0},
		{time.Date(2018, 2, 15, 22, 41, 51, 0, time.UTC), 100.0},
	}

	if x := emptyTimeseries.After(zeroTime); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected empty series to be returned on emptyTimeseries.After(zeroTime); instead got %v", x)
	}

	if x := ts.After(maxTime); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected no items in time series after maximum time, instead got: %v", x)
	}

	if x := ts.After(zeroTime); !ts.Equal(ts) {
		t.Fatalf("expected ts.After(zeroTime) to return ts; instead got: %v", x)
	}

	// After is inclusive
	if x := ts.After(ts[1].Time); !x.Equal(ts[1:]) {
		t.Fatalf("expected ts.After(ts[1].Time) to return %v; instead got %v", ts[1:], x)
	}
}

func TestBefore(t *testing.T) {
	ts := Timeseries{
		{time.Date(2018, 2, 15, 22, 41, 31, 0, time.UTC), 100.0},
		{time.Date(2018, 2, 15, 22, 41, 41, 0, time.UTC), 50.0},
		{time.Date(2018, 2, 15, 22, 41, 51, 0, time.UTC), 100.0},
	}

	if x := emptyTimeseries.Before(maxTime); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected nil to be returned on emptyTimeseries.Before(maxTime); instead got %v", x)
	}

	if x := ts.Before(zeroTime); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected no items in time series before zero time, instead got %v", x)
	}

	if x := ts.Before(maxTime); !x.Equal(ts) {
		t.Fatalf("expected ts.Before(maxTime) to return ts; instead got: %v", x)
	}

	// Before is exclusive
	if x := ts.Before(ts[2].Time); !x.Equal(ts[:2]) {
		t.Fatalf("expected ts.Before(%s) to return %v; instead got %v", ts[2].Time, ts[:2], x)
	}
}

func TestBetween(t *testing.T) {
	ts := Timeseries{
		{time.Date(2018, 2, 15, 22, 41, 31, 0, time.UTC), 100.0},
		{time.Date(2018, 2, 15, 22, 41, 41, 0, time.UTC), 50.0},
		{time.Date(2018, 2, 15, 22, 41, 51, 0, time.UTC), 100.0},
	}

	if x := ts.Between(zeroTime, zeroTime); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected ts.Between(zeroTime, zeroTime) to return empty; instead got %v", x)
	}

	// We return in the range [from, to).  We should only get the first two elements.
	if x := ts.Between(zeroTime, ts[2].Time); !x.Equal(ts[0:2]) {
		t.Fatalf("expected ts.Between(zeroTime, %s) to return %v; instead got %v", ts[2].Time, ts[0:2], x)
	}

	// we return in the range [from, to).  Ensure that we only get the middle
	// element
	a := ts[1].Time
	b := ts[2].Time
	if x := ts.Between(a, b); !x.Equal(ts[1:2]) {
		t.Fatalf("expected ts.Between(%v, %v) to return %v; instead got %v", a, b, ts[1:2], x)
	}
}

func TestDifference(t *testing.T) {
	if x := emptyTimeseries.Difference(); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected difference of empty series to return empty series; instead got %v", x)
	}

	ts := Timeseries{
		Point{time.Date(2018, 2, 15, 22, 41, 31, 0, time.UTC), 100.0},
		Point{time.Date(2018, 2, 15, 22, 41, 41, 0, time.UTC), 50.0},
		Point{time.Date(2018, 2, 15, 22, 41, 51, 0, time.UTC), 100.0},
	}

	if x := ts[0:1].Difference(); !x.Equal(emptyTimeseries) {
		t.Fatalf("expected difference of 1-length series to return empty series; instead got %v", x)
	}

	actual := ts.Difference()
	expected := Timeseries{
		Point{ts[1].Time, -50.0},
		Point{ts[2].Time, 50.0},
	}

	if !actual.Equal(expected) {
		t.Fatalf("expected ts.Difference() to return %v; instead got %v", expected, actual)
	}
}

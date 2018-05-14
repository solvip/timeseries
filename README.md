timeseries
==========

[![Build Status](https://travis-ci.org/solvip/timeseries.svg?branch=master)](https://travis-ci.org/solvip/timeseries)

Package timeseries provides utilities to manipulate and analyze timeseries data.
For compatability with Gonum, a Timeseries is simply a pair of float64 slices,
representing the X and the Y axis.
You can manipulate them as you wish, but ensure two things:

- Many of the methods in this library assume that the data is sorted.  If you
  do not insert in sorted order, ensure that you call Sort()
  
- Ensure that Timeseries.Xs and Timeseries.Ys is always of equal length
  if you manipulate them without the accessors provided.  Violating this constraint will result in pancis.

At the time of this writing(May, 2018), please do not assume API stability.

# License

timeseries is licensed under MIT.  See [LICENSE](LICENSE)


package matrix

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

// Format returns matrix formatter for printing matrices
func Format(m mat.Matrix) fmt.Formatter {
	return mat.Formatted(m, mat.Prefix(""), mat.Squeeze())
}

// NewDenseRand creates a new matrix with provided number of rows and columns
// which is initialized to random numbers uniformly distributed in interval [min, max].
// NewDenseRand fails if non-positive matrix dimensions are requested.
func NewDenseRand(rows, cols int, min, max float64) (*mat.Dense, error) {
	return withValidDims(rows, cols, func() (*mat.Dense, error) {
		// set random seed
		rand.Seed(55)
		// allocate data slice
		randVals := make([]float64, rows*cols)
		for i := range randVals {
			// we need value between 0 and 1.0
			randVals[i] = rand.Float64()*(max-min) + min
		}
		return mat.NewDense(rows, cols, randVals), nil
	})
}

// NewDenseVal returns a matrix with rows x cols whose each element is set to val.
// NewDenseVal fails if invalid matrix dimensions are requested.
func NewDenseVal(rows, cols int, val float64) (*mat.Dense, error) {
	return withValidDims(rows, cols, func() (*mat.Dense, error) {
		// allocate zero matrix and set every element to val
		constMx := mat.NewDense(rows, cols, nil)
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				constMx.Set(i, j, val)
			}
		}
		return constMx, nil
	})
}

// NewDenseValIdentity returns a matrix with size n x n whose diagonal elements are set to val.
// NewDenseValIdentity fails if invalid matrix dimensions are requested.
func NewDenseValIdentity(n int, val float64) (*mat.Dense, error) {
	return withValidDims(n, n, func() (*mat.Dense, error) {
		data := make([]float64, n)
		for i := range data {
			data[i] = val
		}
		diag := mat.NewDiagDense(n, data)

		m := &mat.Dense{}
		m.Clone(diag)

		return m, nil
	})
}

// AddVal adds a constant value to every element of matrix
// It modifies the matrix m passed in as a paramter.
// AddConstant fails with error if empty matrix is supplied
func AddVal(m *mat.Dense, val float64) (*mat.Dense, error) {
	if m == nil {
		return nil, fmt.Errorf("invalid matrix supplied: %v", m)
	}
	rows, cols := m.Dims()
	return withValidDims(rows, cols, func() (*mat.Dense, error) {
		// allocate zero matrix and set every element to val
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				m.Set(i, j, m.At(i, j)+val)
			}
		}
		return m, nil
	})
}

// ColsSum returns a slice of sum values of first cols number of matrix columns
// It returns error if passed in matrix is nil, has zero size or requested number
// of columns exceeds the number of columns in the matrix passed in as parameter.
func ColsSum(cols int, m *mat.Dense) ([]float64, error) {
	return withValidDim("cols", cols, m, mat.Sum)
}

// ColsMax returns a slice of max values of first cols number of matrix columns
// It returns error if passed in matrix is nil, has zero size or requested number
// of columns exceeds the number of columns in the matrix passed in as parameter.
func ColsMax(cols int, m *mat.Dense) ([]float64, error) {
	return withValidDim("cols", cols, m, mat.Max)
}

// ColsMin returns a slice of min values of first cols number of matrix columns
// It returns error if passed in matrix is nil, has zero size or requested number
// of columns exceeds the number of columns in the matrix passed in as parameter.
func ColsMin(cols int, m *mat.Dense) ([]float64, error) {
	return withValidDim("cols", cols, m, mat.Min)
}

// ColsMean returns a slice of mean values of first cols matrix columns
// It returns error if passed in matrix is nil or has zero size or requested number
// of columns exceeds the number of columns in matrix m.
func ColsMean(cols int, m *mat.Dense) ([]float64, error) {
	return withValidDim("cols", cols, m, mean)
}

// ColsStdev returns a slice of standard deviations of first cols matrix columns
// It returns error if passed in matrix is nil or has zero size or requested number
// of columns exceeds the number of columns in matrix m.
func ColsStdev(cols int, m *mat.Dense) ([]float64, error) {
	return withValidDim("cols", cols, m, stdev)
}

// RowsMax returns a slice of max values of first rows matrix rows.
// It returns error if passed in matrix is nil or has zero size or requested number
// of rows exceeds the number of rows in matrix m.
func RowsMax(rows int, m *mat.Dense) ([]float64, error) {
	return withValidDim("rows", rows, m, mat.Max)
}

// RowsSum returns a slice of sum values of first rows number of matrix columns
// It returns error if passed in matrix is nil, has zero size or requested number
// of columns exceeds the number of columns in the matrix passed in as parameter.
func RowsSum(rows int, m *mat.Dense) ([]float64, error) {
	return withValidDim("rows", rows, m, mat.Sum)
}

// RowsMin returns a slice of min values of first rows matrix rows.
// It returns error if passed in matrix is nil or has zero size or requested number
// of rows exceeds the number of rows in matrix m.
func RowsMin(rows int, m *mat.Dense) ([]float64, error) {
	return withValidDim("rows", rows, m, mat.Min)
}

// RowsMean returns a slice of mean values of first rows matrix rows
// It returns error if passed in matrix is nil or has zero size or requested number
// of columns exceeds the number of columns in matrix m.
func RowsMean(rows int, m *mat.Dense) ([]float64, error) {
	return withValidDim("rows", rows, m, mean)
}

// viewFunc defines matrix dimension view function
type viewFunc func(int) mat.Vector

// dimFn applies function fn to first count matrix rows or columns.
// dim can be either set to rows or cols.
// dimFn collects the results into a slice and returns it
func dimFn(dim string, count int, m *mat.Dense, fn func(mat.Matrix) float64) []float64 {
	res := make([]float64, count)
	var viewFn viewFunc
	switch dim {
	case "rows":
		viewFn = m.RowView
	case "cols":
		viewFn = m.ColView
	}
	for i := 0; i < count; i++ {
		res[i] = fn(viewFn(i))
	}
	return res
}

// withValidDim executes function fn on first count of matrix columns or rows.
// It collects the results of each calculation and returns it in a slice.
// It returns error if either matrix m is nil, has zero size or requested number of
// particular dimension is larger than the matrix m dimensions.
func withValidDim(dim string, count int, m *mat.Dense,
	fn func(mat.Matrix) float64) ([]float64, error) {
	// matrix can't be nil
	if m == nil {
		return nil, fmt.Errorf("invalid matrix supplied: %v", m)
	}
	rows, cols := m.Dims()
	switch dim {
	case "rows":
		if count > rows {
			return nil, fmt.Errorf("row count exceeds matrix rows: %d", count)
		}
	case "cols":
		if count > cols {
			return nil, fmt.Errorf("column count exceeds matrix columns: %d", count)
		}
	}
	return dimFn(dim, count, m, fn), nil
}

// withValidDims validates if the rows and cols are valid matrix dimensions
// It returns error if either rows or cols are invalid i.e. non-positive integers
func withValidDims(rows, cols int, fn func() (*mat.Dense, error)) (*mat.Dense, error) {
	// can not create matrix with negative dimensions
	if rows <= 0 {
		return nil, fmt.Errorf("invalid number of rows: %d", rows)
	}
	if cols <= 0 {
		return nil, fmt.Errorf("invalid number of columns: %d", cols)
	}
	return fn()
}

// returns a mean valur for a given matrix
func mean(m mat.Matrix) float64 {
	r, c := m.Dims()
	return mat.Sum(m) / (float64(r) * float64(c))
}

// returns a mean valur for a given matrix
func stdev(m mat.Matrix) float64 {
	r, _ := m.Dims()
	col := make([]float64, r)
	mat.Col(col, 0, m)
	return stat.StdDev(col, nil)
}

// Cov calculates a covariance matrix with data stored in m along dim dimension.
// It returns error if the covariance could not be calculated.
func Cov(m *mat.Dense, dim string) (*mat.SymDense, error) {
	// 1. We will calculate zero mean matrix x of the data
	// 2. 1/(n-1)(x * x^T) will give us covariance of the data
	rows, cols := m.Dims()

	// calculate mean data vector across dimension dim
	var mean []float64
	var count float64
	if strings.EqualFold(dim, "rows") {
		mean, _ = RowsMean(rows, m)
		count = float64(rows)
	} else {
		mean, _ = ColsMean(cols, m)
		count = float64(cols)
	}

	// x is zero-mean matrix with data stored in dimension dim
	x := mat.NewDense(rows, cols, nil)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if strings.EqualFold(dim, "rows") {
				x.Set(r, c, m.At(r, c)-mean[c])
			} else {
				x.Set(r, c, m.At(r, c)-mean[r])
			}
		}
	}

	cov := new(mat.Dense)
	cov.Mul(x, x.T())
	cov.Scale(1/(count-1.0), cov)

	return ToSymDense(cov)
}

// ToSymDense converts m to SymDense (symmetric Dense matrix) if possible.
// It returns error if the provided Dense matrix is not symmetric.
func ToSymDense(m *mat.Dense) (*mat.SymDense, error) {
	r, c := m.Dims()
	if r != c {
		return nil, errors.New("Matrix must be square")
	}

	mT := m.T()
	vals := make([]float64, r*c)
	idx := 0
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if i != j && !floats.EqualWithinAbsOrRel(mT.At(i, j), m.At(i, j), 1e-6, 1e-2) {
				return nil, fmt.Errorf("Matrix not symmetric (%d, %d): %.40f != %.40f\n%v",
					i, j, mT.At(i, j), m.At(i, j), Format(m))
			}
			vals[idx] = m.At(i, j)
			idx++
		}
	}

	return mat.NewSymDense(r, vals), nil
}

// BlockDiag accepts a slice of matrices, turns them into a block diagonal matrix and returns it.
// It skips zero sized matrices when assembling the block diagonal matrix.
func BlockDiag(mx []mat.Matrix) *mat.Dense {
	m := &mat.Dense{}

	for i := range mx {
		r, c := mx[i].Dims()
		if r == 0 || c == 0 {
			continue
		}
		dR, dC := m.Dims()
		m = m.Grow(r, c).(*mat.Dense)
		m.Slice(dR, dR+r, dC, dC+c).(*mat.Dense).Copy(mx[i])
	}

	return m
}

// BlockSymDiag turns a slice of symmetric matrices into a symmetric block diagonal matrix and returns it.
// It skips zero sized matrices when assembling the symmetric block diagonal matrix.
func BlockSymDiag(mx []mat.Symmetric) *mat.SymDense {
	m := &mat.SymDense{}

	for i := range mx {
		n := mx[i].Symmetric()
		if n == 0 {
			continue
		}
		r := m.Symmetric()
		m = m.GrowSquare(n).(*mat.SymDense)
		m.SliceSquare(r, r+n).(*mat.SymDense).CopySym(mx[i])
	}

	return m
}

// Unroll unrolls all elements of matrix into *mat.VecDense and returns it
// Matrix elements can be unrolled either by row or by a column.
func Unroll(m *mat.Dense, byRow bool) *mat.VecDense {
	if byRow {
		return toVecByRow(m)
	}
	return toVecByCol(m)
}

// toVecByRow rolls matrix into a slice by rows
func toVecByRow(m *mat.Dense) *mat.VecDense {
	rows, cols := m.Dims()
	vec := make([]float64, rows*cols)
	for i := 0; i < rows; i++ {
		view := m.RowView(i)
		for j := 0; j < view.Len(); j++ {
			vec[i*cols+j] = view.At(j, 0)
		}
	}

	return mat.NewVecDense(rows*cols, vec)
}

// toVecByCol rolls matrix into a slice by columns
func toVecByCol(m *mat.Dense) *mat.VecDense {
	rows, cols := m.Dims()
	vec := make([]float64, rows*cols)
	for i := 0; i < cols; i++ {
		view := m.ColView(i)
		for j := 0; j < view.Len(); j++ {
			vec[i*rows+j] = view.At(j, 0)
		}
	}

	return mat.NewVecDense(rows*cols, vec)
}

// SetVals sets all elements of a matrix to values stored in vals
// passed in as a parameter. It fails with error if number of elements
// of the matrix is bigger than number of elements of the slice.
func SetVals(m *mat.Dense, vals []float64, byRow bool) (err error) {
	r, c := m.Dims()
	if r*c != len(vals) {
		err = fmt.Errorf("elements count mismatch: Vec: %d, Matrix: %d", len(vals), r*c)
		return
	}
	if byRow {
		setByRow(m, vals)
		return
	}
	setByCol(m, vals)
	return
}

// setByRow sets elements of m from vec by rows
func setByRow(m *mat.Dense, vec []float64) {
	rows, cols := m.Dims()
	acc := 0
	for i := 0; i < rows; i++ {
		m.SetRow(i, vec[acc:(acc+cols)])
		acc += cols
	}
}

// setByCol sets elements of m from vec by columns
func setByCol(m *mat.Dense, vec []float64) {
	rows, cols := m.Dims()
	acc := 0
	for i := 0; i < cols; i++ {
		m.SetCol(i, vec[acc:(acc+rows)])
		acc += rows
	}
}

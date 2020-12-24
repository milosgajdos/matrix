package matrix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

var (
	errInvMx   = "invalid matrix supplied: %v"
	errExcCols = "column count exceeds matrix columns: %d"
	errExcrows = "row count exceeds matrix rows: %d"
)

func TestFormat(t *testing.T) {
	assert := assert.New(t)

	out := `⎡1.2  3.4⎤
⎣4.5  6.7⎦`
	data := []float64{1.2, 3.4, 4.5, 6.7}
	m := mat.NewDense(2, 2, data)
	assert.NotNil(m)

	format := Format(m)
	tstOut := fmt.Sprintf("%v", format)
	assert.Equal(out, tstOut)
}

func TestNewRandDense(t *testing.T) {
	assert := assert.New(t)

	// create new matrix
	rows, cols := 2, 3
	min, max := 1.0, 2.0
	randMx, err := NewDenseRand(rows, cols, min, max)
	assert.NotNil(randMx)
	assert.NoError(err)
	r, c := randMx.Dims()
	assert.Equal(r, rows)
	assert.Equal(c, cols)
	for i := 0; i < c; i++ {
		col := randMx.ColView(i)
		assert.True(max >= mat.Max(col))
		assert.True(min <= mat.Min(col))
	}

	// Can't create new matrix
	randMx, err = NewDenseRand(rows, -6, min, max)
	assert.Nil(randMx)
	assert.Error(err)

	// Can't create new matrix
	randMx, err = NewDenseRand(-10, cols, min, max)
	assert.Nil(randMx)
	assert.Error(err)
}

func TestNewConstDense(t *testing.T) {
	assert := assert.New(t)

	// all elements must be equal to 1.0
	constVec := []float64{1.0, 1.0, 1.0, 1.0}
	constMx := mat.NewDense(2, 2, constVec)
	mx, err := NewDenseVal(2, 2, 1.0)
	assert.NoError(err)
	assert.NotNil(mx)
	assert.True(mat.Equal(constMx, mx))

	// Can't create new matrix
	constMx, err = NewDenseVal(3, -6, 1.0)
	assert.Nil(constMx)
	assert.Error(err)

	// Can't create new matrix
	constMx, err = NewDenseVal(-3, 10, 1.0)
	assert.Nil(constMx)
	assert.Error(err)
}

func TestNewConstEyeDense(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1.0, 0.0, 0.0, 1.0}
	exp := mat.NewDense(2, 2, data)
	m, err := NewDenseValIdentity(2, 1.0)
	assert.NoError(err)
	assert.NotNil(m)
	assert.True(mat.Equal(m, exp))

	// Can't create new matrix
	m, err = NewDenseValIdentity(-6, 1.0)
	assert.Nil(m)
	assert.Error(err)
}

func TestAddConst(t *testing.T) {
	assert := assert.New(t)

	// all elements must be equal to 1.0
	val := 0.5
	mx := mat.NewDense(2, 2, []float64{1.0, 2.0, 2.5, 2.5})
	mc := mat.NewDense(2, 2, []float64{1.5, 2.5, 3.0, 3.0})

	mx, err := AddVal(mx, val)
	assert.NotNil(mx)
	assert.NoError(err)
	assert.True(mat.EqualApprox(mx, mc, 0.01))

	// incorrect matrix passed in
	mx, err = AddVal(nil, val)
	assert.Nil(mx)
	assert.Error(err)
}

func TestRowsColsMax(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	colsMax := []float64{8.9, 10.0}
	rowsMax := []float64{3.4, 6.7, 10.0}
	mx := mat.NewDense(3, 2, data)
	assert.NotNil(mx)

	rows, cols := mx.Dims()
	// check cols max
	max, err := ColsMax(cols, mx)
	assert.NotNil(max)
	assert.NoError(err)
	assert.EqualValues(colsMax, max)

	// check rows max
	max, err = RowsMax(rows, mx)
	assert.NotNil(max)
	assert.NoError(err)
	assert.EqualValues(rowsMax, max)

	// requested number of cols exceeds matrix dims
	max, err = ColsMax(cols+1, mx)
	assert.Nil(max)
	assert.EqualError(err, fmt.Sprintf(errExcCols, cols+1))

	// requested number of rows exceeds matrix dims
	max, err = RowsMax(rows+1, mx)
	assert.Nil(max)
	assert.EqualError(err, fmt.Sprintf(errExcrows, rows+1))

	// should get nil back
	mx = nil
	max, err = ColsMax(cols, mx)
	assert.Nil(max)
	assert.EqualError(err, fmt.Sprintf(errInvMx, mx))
	max, err = RowsMax(rows, mx)
	assert.Nil(max)
	assert.EqualError(err, fmt.Sprintf(errInvMx, mx))
}

func TestRowsColsMin(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	colsMin := []float64{1.2, 3.4}
	rowsMin := []float64{1.2, 4.5, 8.9}
	mx := mat.NewDense(3, 2, data)
	assert.NotNil(mx)

	rows, cols := mx.Dims()
	// check cols
	min, err := ColsMin(cols, mx)
	assert.NotNil(min)
	assert.NoError(err)
	assert.EqualValues(colsMin, min)

	// check rows
	min, err = RowsMin(rows, mx)
	assert.NotNil(min)
	assert.NoError(err)
	assert.EqualValues(rowsMin, min)

	// requested number of cols exceeds matrix dims
	min, err = ColsMin(cols+1, mx)
	assert.Nil(min)
	assert.EqualError(err, fmt.Sprintf(errExcCols, cols+1))

	// requested number of rows exceeds matrix dims
	min, err = RowsMin(rows+1, mx)
	assert.Nil(min)
	assert.EqualError(err, fmt.Sprintf(errExcrows, rows+1))

	// should get nil back
	mx = nil
	min, err = ColsMin(cols, mx)
	assert.Nil(min)
	assert.EqualError(err, fmt.Sprintf(errInvMx, mx))
	min, err = RowsMin(rows, mx)
	assert.Nil(min)
	assert.EqualError(err, fmt.Sprintf(errInvMx, mx))
}

func TestRowsColsSums(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	rowSums := []float64{4.6, 11.2, 18.9}
	colSums := []float64{14.6, 20.1}
	delta := 0.001

	m := mat.NewDense(3, 2, data)
	assert.NotNil(m)
	r, c := m.Dims()

	// check rows
	resRows, err := RowsSum(r, m)
	assert.NotNil(resRows)
	assert.NoError(err)
	assert.InDeltaSlice(rowSums, resRows, delta)

	// check cols
	resCols, err := ColsSum(c, m)
	assert.NotNil(resCols)
	assert.NoError(err)
	assert.InDeltaSlice(colSums, resCols, delta)
}

func TestRowsColsMean(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	mx := mat.NewDense(3, 2, data)
	assert.NotNil(mx)
	colsMean := []float64{4.8667, 6.7000}
	rowsMean := []float64{2.3, 5.6, 9.45}

	rows, cols := mx.Dims()

	// check cols
	me, err := ColsMean(cols, mx)
	assert.NotNil(me)
	assert.NoError(err)
	assert.True(floats.EqualApprox(colsMean, me, 0.01))

	// check rows
	me, err = RowsMean(rows, mx)
	assert.NotNil(me)
	assert.NoError(err)
	assert.True(floats.EqualApprox(rowsMean, me, 0.01))
}

func TestColsStdev(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	mx := mat.NewDense(3, 2, data)
	assert.NotNil(mx)
	colsStdev := []float64{3.8631, 3.3000}

	// check cols
	_, cols := mx.Dims()
	sd, err := ColsStdev(cols, mx)
	assert.NotNil(sd)
	assert.NoError(err)
	assert.True(floats.EqualApprox(colsStdev, sd, 0.01))
}

func TestCov(t *testing.T) {
	assert := assert.New(t)
	data := []float64{1, 2, 2, 4}
	delta := 0.001

	rowCov := mat.NewDense(2, 2, []float64{1.25, -1.25, -1.25, 1.25})
	colCov := mat.NewDense(2, 2, []float64{0.5, 1.0, 1.0, 2.0})

	m := mat.NewDense(2, 2, data)
	assert.NotNil(m)

	cov, err := Cov(m, "rows")
	assert.NotNil(cov)
	assert.NoError(err)

	rows, cols := cov.Dims()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			assert.InDelta(rowCov.At(r, c), cov.At(r, c), delta)
		}
	}

	cov, err = Cov(m, "cols")
	assert.NotNil(cov)
	assert.NoError(err)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			assert.InDelta(colCov.At(r, c), cov.At(r, c), delta)
		}
	}
}

func TestToSymDense(t *testing.T) {
	assert := assert.New(t)

	badMx := mat.NewDense(2, 1, []float64{0.5, 1.0})
	notSymMx := mat.NewDense(2, 2, []float64{0.5, 1.0, 2.0, 2.0})
	symMx := mat.NewDense(2, 2, []float64{0.5, 1.0, 1.0, 2.0})

	sym, err := ToSymDense(badMx)
	assert.Nil(sym)
	assert.Error(err)

	sym, err = ToSymDense(notSymMx)
	assert.Nil(sym)
	assert.Error(err)

	sym, err = ToSymDense(symMx)
	assert.NotNil(sym)
	assert.NoError(err)
}

func TestBlockDiag(t *testing.T) {
	assert := assert.New(t)

	xVec := mat.NewVecDense(3, []float64{1.0, 2.0, 3.0})
	yMat := mat.NewDense(2, 2, []float64{4.0, 5.0, 6.0, 7.0})
	exp := mat.NewDense(5, 3, []float64{
		1.0, 0.0, 0.0,
		2.0, 0.0, 0.0,
		3.0, 0.0, 0.0,
		0.0, 4.0, 5.0,
		0.0, 6.0, 7.0})

	mx := make([]mat.Matrix, 3)
	mx[0] = xVec
	mx[1] = &mat.Dense{}
	mx[2] = yMat

	blkDiag := BlockDiag(mx)

	r, c := blkDiag.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			assert.InDelta(exp.At(i, j), blkDiag.At(i, j), 0.001)
		}
	}
}

func TestBlockSymDiag(t *testing.T) {
	assert := assert.New(t)

	x := mat.NewSymDense(2, []float64{1.0, 2.0, 2.0, 1.0})
	y := mat.NewSymDense(1, []float64{4.0})
	exp := mat.NewSymDense(3, []float64{
		1.0, 2.0, 0.0,
		2.0, 1.0, 0.0,
		0.0, 0.0, 4.0})

	mx := make([]mat.Symmetric, 3)
	mx[0] = x
	mx[1] = &mat.SymDense{}
	mx[2] = y

	m := BlockSymDiag(mx)

	n := m.Symmetric()
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			assert.InDelta(exp.At(i, j), m.At(i, j), 0.001)
		}
	}
}

func TestUnroll(t *testing.T) {
	assert := assert.New(t)

	// expected outputs
	byRow := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	byCol := []float64{1.2, 4.5, 8.9, 3.4, 6.7, 10.0}
	// NewDense creates new matrix by rows
	tstMx := mat.NewDense(3, 2, byRow)

	rVec := mat.NewVecDense(6, byRow)
	cVec := mat.NewVecDense(6, byCol)

	// Check if the marix is rolled into vector by row
	rowVec := Unroll(tstMx, true)
	assert.NotNil(rowVec)
	assert.True(mat.Equal(rowVec, rVec))

	// Check if the marix is rolled into vector by col
	colVec := Unroll(tstMx, false)
	assert.NotNil(colVec)
	assert.True(mat.Equal(colVec, cVec))
}

func TestSetVals(t *testing.T) {
	assert := assert.New(t)

	// expected results
	data := []float64{1.2, 3.4, 4.5, 6.7, 8.9, 10.0}
	mx := mat.NewDense(3, 2, nil)
	assert.NotNil(mx)

	// Set matrix by row
	err := SetVals(mx, data, true)
	rowMx := mat.NewDense(3, 2, data)
	assert.NoError(err)
	assert.NotNil(rowMx)
	assert.True(mat.Equal(mx, rowMx))

	// Set matrix by col
	err = SetVals(mx, data, false)
	colData := []float64{1.2, 6.7, 3.4, 8.9, 4.5, 10.0}
	colMx := mat.NewDense(3, 2, colData)
	assert.NoError(err)
	assert.NotNil(colMx)
	assert.True(mat.Equal(mx, colMx))

	// Vector is smaller than number of matrix elements
	shortVec := []float64{1.3, 2.4}
	err = SetVals(mx, shortVec, true)
	assert.Error(err)
}

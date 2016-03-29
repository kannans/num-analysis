package cholesky

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/ludecomp"
)

// Cholesky represents the Cholesky decomposition
// of a matrix (that is, L*L^T form).
type Cholesky struct {
	size  int
	lower []float64
}

// Decompose computes the Cholesky decomposition of
// a symmetric positive-definite matrix.
//
// This will not verify that the given matrix is
// symmetric positine-definite.
//
// The lower-triangular portion of the given matrix
// will not be accessed.
func Decompose(matrix ludecomp.Matrix) *Cholesky {
	res := &Cholesky{
		size:  matrix.N,
		lower: make([]float64, matrix.N*(matrix.N+1)/2),
	}

	for lowerColumn := 0; lowerColumn < matrix.N; lowerColumn++ {
		summer := kahan.NewSummer64()
		summer.Add(matrix.Get(lowerColumn, lowerColumn))
		for i := 0; i < lowerColumn; i++ {
			summer.Add(-res.Get(lowerColumn, i))
		}
		diagEntry := math.Sqrt(summer.Sum())
		res.set(lowerColumn, lowerColumn, diagEntry)

		for lowerRow := lowerColumn + 1; lowerRow < matrix.N; lowerRow++ {
			summer = kahan.NewSummer64()
			summer.Add(matrix.Get(lowerColumn, lowerRow))
			for i := 0; i < lowerColumn; i++ {
				summer.Add(res.Get(lowerColumn, i) * res.Get(i, lowerRow))
			}
			sum := summer.Sum()
			res.set(lowerRow, lowerColumn, sum/diagEntry)
		}
	}

	return res
}

// Size returns N for this NxN matrix.
func (c *Cholesky) Size() int {
	return c.size
}

// Get returns the entry at the given row and column
// in the lower triangular matrix of L from L*L^T.
func (c *Cholesky) Get(row, col int) float64 {
	if row < 0 || col < 0 || row > c.size || col > c.size {
		panic("index out of bounds")
	}
	if col > row {
		return 0
	}
	return c.lower[col+(row*(row+1))/2]
}

func (c *Cholesky) set(row, col int, v float64) {
	c.lower[col+(row*(row+1))/2] = v
}

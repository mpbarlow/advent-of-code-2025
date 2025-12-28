package augmentedmatrix

import (
	"advent-of-code-2025/support"
	"fmt"
	"maps"
	"math"
	"strings"
)

func NewRefAugmentedMatrix(coefficients [][]int, totals []int) AugmentedMatrix {
	matrix := AugmentedMatrix{
		rows: make([]AugmentedMatrixRow, len(totals)),
	}

	for i := range len(totals) {
		c := make([]int, len(coefficients))
		for j := range len(coefficients) {
			c[j] = coefficients[j][i]
		}

		matrix.rows[i] = AugmentedMatrixRow{
			coefficients: c,
			total:        totals[i],
		}
	}

	matrix.toRowEchelonForm()

	return matrix
}

type AugmentedMatrixRow struct {
	coefficients []int
	total        int
}

type AugmentedMatrix struct {
	rows []AugmentedMatrixRow
}

func (m *AugmentedMatrix) Print() {
	for _, row := range m.rows {
		fmt.Printf(
			"| %s || %d |\n",
			strings.Join(
				support.Map(row.coefficients, func(i int) string { return fmt.Sprint(i) }),
				" ",
			),
			row.total,
		)
	}

	fmt.Println()
}

// Convert an augmented matrix to row echelon form using the following steps:
// - For each pivot value, find the first row that has a 1 in that column
// - Swap it with the row we're looking at, if it's not already there
// - Use matrix operations to eliminate all cases of non-zeros below the row in the pivot column
//
// So this matrix from the example input:
//
// | 1 1 1 0 || 10 |
// | 1 0 1 1 || 11 |
// | 1 0 1 1 || 11 |
// | 1 1 0 0 ||  5 |
// | 1 1 1 0 || 10 |
// | 0 0 1 0 ||  5 |
//
// Becomes:
//
// |  1  1  1  0 ||  10 |
// |  0 -1  0  1 ||   1 |
// |  0  0 -1  0 ||  -5 |
//
// through a combination of reordering rows and row multiplication/addition/subtraction, and eliminating all-zero rows.
func (m *AugmentedMatrix) toRowEchelonForm() {
	pivotCol := 0
	currRow := 0

	for {
		if currRow >= len(m.rows) || pivotCol >= len(m.rows[0].coefficients) {
			break
		}

		// First try to locate a pivot value
		pivotFound := false
		for i := currRow; i < len(m.rows); i++ {
			if m.rows[i].coefficients[pivotCol] != 0 {
				pivotFound = true

				// Swap the rows if necessary to put the pivot in the right place
				if i != currRow {
					tmp := m.rows[currRow]
					m.rows[currRow] = m.rows[i]
					m.rows[i] = tmp
				}

				break
			}
		}

		// Then eliminate non-zero values below it
		if pivotFound {
			m.eliminateRows(currRow, pivotCol)
			// If we didn't find a pivot for a given column we need to try the next pivot on the same row
			currRow++
		}

		pivotCol++
	}

	// Remove rows reduced to all zeroes, as they do not contribute to our solution
	m.removeEmptyRows()
}

func (m *AugmentedMatrix) eliminateRows(pivotRow, pivotCol int) {
	for i := pivotRow + 1; i < len(m.rows); i++ {
		if m.rows[i].coefficients[pivotCol] != 0 {
			// If our pivot row is P and the target row to eliminate is R, then we can reduce R's value in the
			// pivot column to zero by:
			// - Multiplying (scaling) R by the pivot value p
			// - Subtracting P * [R's value in the pivot column] q
			p := m.rows[pivotRow].coefficients[pivotCol]
			q := m.rows[i].coefficients[pivotCol]

			for j := pivotCol; j < len(m.rows[i].coefficients); j++ {
				m.rows[i].coefficients[j] = (m.rows[i].coefficients[j] * p) - (m.rows[pivotRow].coefficients[j] * q)
			}

			// We then also must scale the total by the same factor
			m.rows[i].total = (m.rows[i].total * p) - (m.rows[pivotRow].total * q)
		}
	}
}

func (m *AugmentedMatrix) removeEmptyRows() {
	i := 0
	for {
		if i >= len(m.rows) {
			break
		}

		if coefficientsAllZero(m.rows[i].coefficients) {
			for j := i; j < len(m.rows)-1; j++ {
				m.rows[j] = m.rows[j+1]
			}

			m.rows = m.rows[:len(m.rows)-1]
		} else {
			i++
		}
	}
}

// Return the smallest possible sum of the set of variables solving the matrix
func (m *AugmentedMatrix) Solve() int {
	// Create a maps to store the pivot column for each row, the expressions for each pivot variable in terms of free
	// variables, and any known fixed values.
	pivotMap := make(map[int]int)
	pivotExpressions := make(map[int]linearExpr)

	// Work down the rows and for the first non-zero value in each, mark that column as a dependent variable
	for i := range m.rows {
		for j := range m.rows[i].coefficients {
			if m.rows[i].coefficients[j] != 0 {
				pivotMap[i] = j
				break
			}
		}
	}

	// Then work from the bottom up, identifying invariants and denoting each dependent variable in terms of frees
	for i := len(m.rows) - 1; i >= 0; i-- {
		// Express the pivot value in terms of the free variables
		pivotExpressions[pivotMap[i]] = newLinearExpr(m.rows[i], pivotMap[i], pivotExpressions)
	}

	smallestSumValues := math.MaxInt

	for _, combination := range enumerateFreeVariableCombinations(m.deriveLimits(pivotExpressions)) {
		values := evaluateExpressions(pivotExpressions, combination)
		if values == nil {
			continue
		}

		sum := support.SumSeq(maps.Values(values))

		if sum < smallestSumValues {
			smallestSumValues = sum
		}
	}

	return smallestSumValues
}

type limits struct {
	min int
	max int
}

// Once we have a set of linear expressions in terms of free variables, we can begin to restrict the possible space
// of values for those.
//
// We know the result of all equations in terms of free variables must be >= 0 as a button cannot be pressed a negative
// number of times. For expressions with only one free variable, this allows us to narrow the possibility space easily:
// e = 3 - f
// 3 - f >= 0
// f <= 3 -> a clear max on f
//
// However, for equations with multiple variables we might need to do multiple passes so we can sub in other min/max
// values we've already determined.
func (m *AugmentedMatrix) deriveLimits(pivotExpressions map[int]linearExpr) map[int]limits {
	lims := make(map[int]limits)

	// Initialise all limits to [0, math.MaxInt]
	for col := range m.rows[0].coefficients {
		// Not a free variable
		if _, ok := pivotExpressions[col]; ok {
			continue
		}

		lims[col] = limits{min: 0, max: math.MaxInt}
	}

	// We will be converging on a set of limits so we likely won't be able to get everything on the first pass, so we
	// continue looping until we have no longer narrowed any bounds
	narrowedBounds := true
	for narrowedBounds {
		narrowedBounds = false

		for _, e := range pivotExpressions {
		boundsLoop:
			for col, coeff := range e.freeVarCoeffs {
				// Rearrange each equation in terms of each free variable
				value := e.constant * -1

				freeVarCoeffs := make(map[int]int)
				for otherCol, otherCoeff := range e.freeVarCoeffs {
					if col == otherCol {
						continue
					}

					freeVarCoeffs[otherCol] = otherCoeff * -1
				}

				// A negative coefficient turns our >= 0 expression into a <= when we multiply everything by -1 to turn
				// -x into x; therefore our minimum bound turns into a maximum bound
				if coeff < 0 {
					// Invert the equation to get a positive expression of col
					coeff *= -1
					value *= -1
					for k := range freeVarCoeffs {
						freeVarCoeffs[k] *= -1

						// If we are deriving a max, we want the max value for any free variables with a positive
						// coefficient, and the min value for any with a negative
						if freeVarCoeffs[k] < 0 {
							value += lims[k].min * freeVarCoeffs[k]
						} else {
							// If we've not been able to derive a maximum value for this variable yet, we can't proceed,
							// we need to try again later
							if lims[k].max == math.MaxInt {
								continue boundsLoop
							}

							value += lims[k].max * freeVarCoeffs[k]
						}
					}

					value = int(math.Ceil(float64(value) / float64(coeff)))

					if value > 0 && value < lims[col].max {
						bounds := lims[col]
						bounds.max = value
						lims[col] = bounds

						narrowedBounds = true
					}
				} else {
					// If we are deriving a min, we want the min value for any free variables with a positive
					// coefficient, and the max value for any with a negative
					for k := range freeVarCoeffs {
						if freeVarCoeffs[k] < 0 {
							if lims[k].max == math.MaxInt {
								continue boundsLoop
							}

							value += lims[k].max * freeVarCoeffs[k]
						} else {
							value += lims[k].min * freeVarCoeffs[k]
						}
					}

					value = int(math.Ceil(float64(value) / float64(coeff)))

					if value > lims[col].min {
						bounds := lims[col]
						bounds.min = value
						lims[col] = bounds

						narrowedBounds = true
					}
				}
			}
		}
	}

	for _, val := range lims {
		if val.max == math.MaxInt {
			panic("Could not determine an upper limit!")
		}
	}

	return lims
}

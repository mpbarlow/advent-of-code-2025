package augmentedmatrix

import "advent-of-code-2025/support"

// | 1 0 1 || 5 |
// becomes
// linearExpr{total: 5, freeVarCoeffs: {1: 0, 2: 1}} which is then reduced down to
// linearExpr{total: 5, freeVarCoeffs: {2: 1}} because a 0 coefficient has no effect
type linearExpr struct {
	coefficient   int         // The coefficient of the variable on the LHS
	constant      int         // The value in the rightmost column of the matrix; the RHS of the equation
	freeVarCoeffs map[int]int // Map of matrix column index onto the coefficient
}

func newLinearExpr(row AugmentedMatrixRow, pivotCol int, pivotExpressions map[int]linearExpr) linearExpr {
	freeVarCoeffs := make(map[int]int)
	total := row.total

	coefficients := make([]int, len(row.coefficients))
	copy(coefficients, row.coefficients)

	scaleFactor := 1

	// To account for the fact that pivot expressions might have coefficients, but we want to keep everything as an
	// integer, we need to scale up every coefficient by the least common multiple of all pivot expression coefficients
	for k, v := range pivotExpressions {
		if k <= pivotCol || coefficients[k] == 0 {
			continue
		}

		scaleFactor = support.Lcm(scaleFactor, v.coefficient)
	}

	for i := pivotCol; i < len(coefficients); i++ {
		coefficients[i] *= scaleFactor
	}

	total *= scaleFactor

	for i := pivotCol + 1; i < len(coefficients); i++ {
		if coefficients[i] == 0 {
			continue
		}

		// We want everything in terms of free variables, so replace any instances of a previous pivot variable with its
		// free variable equivalent. Do this by moving the constant from the expression immediately to the right hand
		// side of the equation, and adding the coefficients of the relevant free variables.
		if expr, ok := pivotExpressions[i]; ok {
			// We scaled the whole row up by the LCM of the subsituted coefficients, so we now need to divide back down
			// to account for the coefficient we actually have. e.g. if we have:
			// 3e = 123
			// f + 2e = 234
			// We scale to 3f + 6e = 702, we are then subsituting "two 3e's" so we need to multiply the values in expr
			// by 6 / 3 = 2
			total -= expr.constant * (coefficients[i] / expr.coefficient)

			for j, e := range expr.freeVarCoeffs {
				coefficients[j] += e * (coefficients[i] / expr.coefficient)
			}

			continue
		}

		// * -1 because we need to subtract it from the RHS
		freeVarCoeffs[i] = coefficients[i] * -1
	}

	// When we come to derive limits for free variables later, life is much easier if we've pre-normalised all the
	// equations to have a positive coefficient on the pivot variable.
	if coefficients[pivotCol] < 0 {
		coefficients[pivotCol] *= -1
		total *= -1

		for k := range freeVarCoeffs {
			freeVarCoeffs[k] *= -1
		}
	}

	return linearExpr{
		coefficient:   coefficients[pivotCol],
		constant:      total,
		freeVarCoeffs: freeVarCoeffs,
	}
}

func evaluateExpressions(expressions map[int]linearExpr, freeVarValues map[int]int) map[int]int {
	results := make(map[int]int)
	for k, v := range freeVarValues {
		results[k] = v
	}

	for col, expr := range expressions {
		result := expr.constant

		for col, coeff := range expr.freeVarCoeffs {
			result += coeff * freeVarValues[col]
		}

		// We know all solutions must be integers
		if result%expr.coefficient != 0 {
			return nil
		}

		result /= expr.coefficient

		// We know these free variable values cannot be valid if any expression evaluates to < 0
		if result < 0 {
			return nil
		}

		results[col] = result
	}

	return results
}

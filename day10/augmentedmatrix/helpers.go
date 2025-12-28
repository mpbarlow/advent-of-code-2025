package augmentedmatrix

import (
	"advent-of-code-2025/support"
	"maps"
	"slices"
)

func enumerateFreeVariableCombinations(lims map[int]limits) []map[int]int {
	// This is pretty fuckin gross but at this point I don't really care.
	// If there's no free variables we can return a slice just containing a single nil map: there is only one solution
	// and when we come to evaluate each expression we will end up ignoring this anyway. This still allows us to
	// "iterate" over the free variable combinations even though there are none.
	if len(lims) == 0 {
		return []map[int]int{nil}
	}

	values := make(map[int][]int)

	for k, v := range lims {
		values[k] = support.Range(v.min, v.max+1)
	}

	combinations := make([]map[int]int, 0)
	combination := make(map[int]int)

	keys := slices.Collect(maps.Keys(lims))
	slices.Sort(keys)

	var loop func(int)
	loop = func(key int) {
		variable := keys[key]

		for _, v := range values[variable] {
			combination[variable] = v

			if key == len(keys)-1 {
				comboCopy := make(map[int]int)

				for k, v := range combination {
					comboCopy[k] = v
				}

				combinations = append(combinations, comboCopy)
			} else {
				loop(key + 1)
			}
		}
	}

	loop(0)

	return combinations
}

func coefficientsAllZero(c []int) bool {
	for _, v := range c {
		if v != 0 {
			return false
		}
	}

	return true
}

package support

func Map[I any, O any](input []I, fn func(I) O) []O {
	output := make([]O, 0, len(input))

	for _, item := range input {
		output = append(output, fn(item))
	}

	return output
}

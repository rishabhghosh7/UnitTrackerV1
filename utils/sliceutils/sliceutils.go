package sliceutils

func Map[A, B any](in []A, mapf func(A) B) []B {
	if len(in) == 0 {
		return nil
	}

	var mapResult []B
	for _, a := range in {
		mapResult = append(mapResult, mapf(a))
	}
	return mapResult
}

func MapWithError[A, B any](in []A, mapf func(A) (B, error)) ([]B, error) {
	if len(in) == 0 {
		return nil, nil
	}

	var mapResult []B
	for _, a := range in {
		b, err := mapf(a)
		if err != nil {
			return nil, err
		}
		mapResult = append(mapResult, b)
	}
	return mapResult, nil
}

func Filter[T any](in []T, predicate func(T) bool) []T {
	if len(in) == 0 {
		return nil
	}

	var filterResult []T
	for _, t := range in {
		if predicate(t) {
			filterResult = append(filterResult, t)
		}
	}
	return filterResult
}

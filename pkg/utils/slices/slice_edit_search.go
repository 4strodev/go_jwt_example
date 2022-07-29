package slices

func Contains[T comparable](slice []T, element T) bool {
	var contains bool

	for _, el := range slice {
		if el == element {
			contains = true
			break
		}
	}

	return contains
}

func Filter[T comparable](slice []T, handler func(element T) bool) []T {
	var newSlcie []T

	for _, element := range slice {
		if handler(element) {
			newSlcie = append(newSlcie, element)
		}
	}

	return newSlcie
}

package common

func SlicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, value := range a {
		if value != b[i] {
			return false
		}
	}
	return true
}

func SlicesNotEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, value := range a {
		if value == b[i] {
			return false
		}
	}
	return true
}

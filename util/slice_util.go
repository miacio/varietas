package util

func SliceDistinct[T comparable](vs ...T) []T {
	result := []T{}
	for i := range vs {
		if SliceContain(result, vs[i]) {
			continue
		}
		result = append(result, vs[i])
	}
	return result
}

func SliceContain[T comparable](vs []T, v T) bool {
	for i := range vs {
		if vs[i] == v {
			return true
		}
	}
	return false
}

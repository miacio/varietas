package util

// SliceDistinct slice remove duplicates
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

// SliceContain check the v in vs slice
func SliceContain[T comparable](vs []T, v T) bool {
	for i := range vs {
		if vs[i] == v {
			return true
		}
	}
	return false
}

// SliceMaxNumber
func SliceMaxNumber[T NumberGenericity](vs []T) T {
	var max T
	for _, v := range vs {
		if max < v {
			max = v
		}
	}
	return max
}

// SliceMinNumber
func SliceMinNumber[T NumberGenericity](vs []T) T {
	var min T
	for i, v := range vs {
		if i == 0 {
			min = v
		}
		if min > v {
			min = v
		}
	}
	return min
}

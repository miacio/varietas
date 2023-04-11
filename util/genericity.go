package util

type IntGenericity interface {
	int | int8 | int16 | int32 | int64
}

type UIntGenericity interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type FloatGenericity interface {
	float32 | float64
}

type NumberGenericity interface {
	IntGenericity | UIntGenericity | FloatGenericity
}

func Min[T NumberGenericity](a, b T) T {
	if a <= b {
		return a
	}
	return b
}

func Max[T NumberGenericity](a, b T) T {
	if a >= b {
		return a
	}
	return b
}

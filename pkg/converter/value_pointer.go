package converter

// Value returns the value of the input address
func Value[T any](input *T) T {
	if input == nil {
		return *new(T)
	}

	return *input
}

// Pointer returns the address of the input value
func Pointer[T any](input T) *T {
	return &input
}

// SliceValue returns a slice of values from the input slice of addresses
func SliceValue[T any](input []*T) []T {
	if len(input) == 0 {
		return nil
	}

	out := make([]T, len(input))
	for i := range input {
		out[i] = Value(input[i])
	}

	return out
}

// SlicePointer returns a slice of addresses from the input slice of values
func SlicePointer[T any](input []T) []*T {
	if len(input) == 0 {
		return nil
	}

	out := make([]*T, len(input))
	for i := range input {
		out[i] = Pointer(input[i])
	}

	return out
}

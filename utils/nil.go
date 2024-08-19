package utils

func SafeCall[T any](data any, fn func(T)) {
	switch d := data.(type) {
	case nil:
		return
	default:
		fn(d.(T))
	}
}

package harness

// Condition is a predicate over an object of type T.
type Condition[T any] func(T) bool

func And[T any](cs ...Condition[T]) Condition[T] {
	return func(o T) bool {
		for _, c := range cs {
			if !c(o) {
				return false
			}
		}
		return true
	}
}

func Or[T any](cs ...Condition[T]) Condition[T] {
	return func(o T) bool {
		for _, c := range cs {
			if c(o) {
				return true
			}
		}
		return false
	}
}

func Not[T any](c Condition[T]) Condition[T] {
	return func(o T) bool {
		return !c(o)
	}
}

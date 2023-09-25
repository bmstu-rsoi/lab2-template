package collections

type Countable[T any] struct {
	Total uint64

	Items []T
}

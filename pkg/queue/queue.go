package queue

type queue[T any] struct {
	Push   func(T)
	Pop    func() T
	Length func() int
	Front  func() T
}

func NewQueue[T any]() queue[T] {
	slice := make([]T, 0)
	return queue[T]{
		Push: func(i T) {
			slice = append(slice, i)
		},
		Pop: func() T {
			res := slice[0]
			slice = slice[1:]
			return res
		},
		Length: func() int {
			return len(slice)
		},
		Front: func() T {
			return slice[0]
		},
	}
}

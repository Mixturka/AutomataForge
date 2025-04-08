package stack

type stack[T any] struct {
	Push   func(T)
	Pop    func() T
	Length func() int
	Top    func() T
}

func NewStack[T any]() stack[T] {
	slice := make([]T, 0)
	return stack[T]{
		Push: func(i T) {
			slice = append(slice, i)
		},
		Pop: func() T {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		Length: func() int {
			return len(slice)
		},
		Top: func() T {
			return slice[len(slice)-1]
		},
	}
}

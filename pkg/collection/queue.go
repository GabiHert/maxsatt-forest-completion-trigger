package collection

type Queue[T any] interface {
	Push(item T)
	Pop() *T
}

type queue[T any] struct {
	items []T
}

func NewQueue[T any]() Queue[T] {
	return &queue[T]{}
}

func (q *queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *queue[T]) Pop() *T {
	if len(q.items) == 0 {
		return nil
	}

	item := q.items[0]
	q.items = q.items[1:]

	return &item
}

package internal

type Page[T any] struct {
	TotalPages int
	TotalItems int
	Items      []T
}

type PageRequest struct {
	Limit  int
	Offset int
}

func DefaultPageRequest() PageRequest {
	return PageRequest{
		Limit:  10,
		Offset: 0,
	}
}

func EmptyPage[T any]() Page[T] {
	return Page[T]{
		TotalPages: 0,
		TotalItems: 0,
		Items:      []T{},
	}
}

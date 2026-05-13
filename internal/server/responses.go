package server

type jsonBody[T any] struct {
	Body T
}

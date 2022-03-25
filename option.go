package serverfx

import (
	"net/http"
	"time"
)

// Option is a functional option to modify the default Server instance.
type Option[T http.Handler] func(server *Server[T])

// WithAddress replaces the default Server address with the provided one.
func WithAddress[T http.Handler](address string) Option[T] {
	return func(s *Server[T]) {
		s.Address = address
	}
}

// WithMaxHeaderBytes replaces the default Server MaxHeaderBytes with the provided one.
func WithMaxHeaderBytes[T http.Handler](bytes int) Option[T] {
	return func(s *Server[T]) {
		s.MaxHeaderBytes = bytes
	}
}

// WithGracefulTimeout replaces the default Server GracefulTimeout with the provided one.
func WithGracefulTimeout[T http.Handler](timeout time.Duration) Option[T] {
	return func(s *Server[T]) {
		s.GracefulTimeout = timeout
	}
}

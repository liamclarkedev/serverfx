package serverfx

import "time"

// Option is a functional option to modify the default Server instance.
type Option func(server *Server)

// WithAddress replaces the default Server address with the provided one.
func WithAddress(address string) Option {
	return func(s *Server) {
		s.Address = address
	}
}

// WithMaxHeaderBytes replaces the default Server MaxHeaderBytes with the provided one.
func WithMaxHeaderBytes(bytes int) Option {
	return func(s *Server) {
		s.MaxHeaderBytes = bytes
	}
}

// WithGracefulTimeout replaces the default Server GracefulTimeout with the provided one.
func WithGracefulTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.GracefulTimeout = timeout
	}
}

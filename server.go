package serverfx

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrUnableToGracefulShutdown is the annotated error when the Server cannot be shutdown gracefully.
	ErrUnableToGracefulShutdown = errors.New("unable to gracefully shutdown")

	// ErrUnableToListenAndServe is the annotated error when the Server cannot listen and serve the address.
	ErrUnableToListenAndServe = errors.New("unable to listen and serve")
)

const (
	// DefaultGracefulTimeout is the context timeout before a Server shuts down forcefully.
	DefaultGracefulTimeout = 10 * time.Second

	// DefaultMaxHeaderBytes is the maximum permitted size of a HTTP header for the Server.
	DefaultMaxHeaderBytes = http.DefaultMaxHeaderBytes
)

// Server is a http server.
type Server struct {
	Handler         http.Handler
	Address         string
	MaxHeaderBytes  int
	GracefulTimeout time.Duration

	server       *http.Server
	shutdownOnce sync.Once
	shutdownErr  chan error
}

// New initialises a new Server with the provided Handler and Option parameters.
func New(handler http.Handler, options ...Option) *Server {
	server := &Server{
		Handler:         handler,
		Address:         "",
		MaxHeaderBytes:  DefaultMaxHeaderBytes,
		GracefulTimeout: DefaultGracefulTimeout,

		shutdownErr: make(chan error, 1),
	}

	for _, opt := range options {
		opt(server)
	}

	return server
}

// Serve listens and serves the Server.
func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:           s.Address,
		Handler:        s.Handler,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}

	s.server = srv

	go func(s *Server) {
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-shutdownChan
		s.Shutdown()
	}(s)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.shutdownErr <- errors.Wrap(ErrUnableToListenAndServe, err.Error())
	}

	return <-s.shutdownErr
}

// Shutdown shuts down the Server gracefully until the GracefulTimeout duration expires.
func (s *Server) Shutdown() {
	s.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.GracefulTimeout)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.shutdownErr <- errors.Wrap(ErrUnableToGracefulShutdown, err.Error())
			return
		}

		s.shutdownErr <- nil
	})
}

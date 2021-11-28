package serverfx_test

import (
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/clarke94/serverfx"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		handler http.Handler
		options []serverfx.Option
		want    *serverfx.Server
	}{
		{
			name:    "expect an initialised server with the defaults",
			handler: mockHandler{},
			options: []serverfx.Option{},
			want: &serverfx.Server{
				Handler:         mockHandler{},
				Address:         "",
				MaxHeaderBytes:  serverfx.DefaultMaxHeaderBytes,
				GracefulTimeout: serverfx.DefaultGracefulTimeout,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serverfx.New(tt.handler, tt.options...)

			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(serverfx.Server{})) {
				t.Errorf(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(serverfx.Server{})))
			}
		})
	}
}

func TestServer_Serve(t *testing.T) {
	tests := []struct {
		name    string
		handler http.Handler
		options []serverfx.Option
		wantErr error
	}{
		{
			name: "expect the Server forcefully shutdown when the request exceeds the graceful timeout",
			handler: mockHandler{
				WithSleepTime: 1 * time.Second,
			},
			options: []serverfx.Option{
				func(s *serverfx.Server) {
					s.Address = "localhost:54932"
				},
				func(s *serverfx.Server) {
					s.GracefulTimeout = 1 * time.Millisecond
				},
			},
			wantErr: serverfx.ErrUnableToGracefulShutdown,
		},
		{
			name: "expect the Server gracefully shutdown within the graceful timeout",
			handler: mockHandler{
				WithSleepTime: 1 * time.Second,
			},
			options: []serverfx.Option{
				func(s *serverfx.Server) {
					s.Address = "localhost:54932"
				},
				func(s *serverfx.Server) {
					s.GracefulTimeout = 5 * time.Millisecond
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serverfx.New(tt.handler, tt.options...)

			done := make(chan error, 1)
			go func() {
				done <- s.Serve()
			}()

			go func() {
				time.Sleep(1 * time.Second)
				if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
					t.Error(err)
				}
			}()

			_, err := http.Get("http://localhost:54932/foo")
			if err != nil {
				t.Error(err)
			}

			if !cmp.Equal(<-done, tt.wantErr, cmpopts.EquateErrors()) {
				t.Errorf(cmp.Diff(<-done, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockHandler struct {
	WithSleepTime time.Duration
}

func (m mockHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	time.Sleep(m.WithSleepTime)
	writer.WriteHeader(http.StatusOK)
}

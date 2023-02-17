package serverfx_test

import (
	"testing"
	"time"

	"github.com/liamclarkedev/serverfx"

	"github.com/google/go-cmp/cmp"
)

func TestWithAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{
			name:    "expect default address to be replaced with provided address",
			address: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serverfx.New[mockHandler](mockHandler{}, serverfx.WithAddress[mockHandler](tt.address))

			if !cmp.Equal(s.Address, tt.address) {
				t.Error(cmp.Diff(s.Address, tt.address))
			}
		})
	}
}

func TestWithMaxHeaderBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int
	}{
		{
			name:  "expect default bytes to be replaced with provided bytes",
			bytes: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serverfx.New[mockHandler](mockHandler{}, serverfx.WithMaxHeaderBytes[mockHandler](tt.bytes))

			if !cmp.Equal(s.MaxHeaderBytes, tt.bytes) {
				t.Error(cmp.Diff(s.MaxHeaderBytes, tt.bytes))
			}
		})
	}
}

func TestWithGracefulTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "expect default timeout to be replaced with provided timeout",
			timeout: 20 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serverfx.New[mockHandler](mockHandler{}, serverfx.WithGracefulTimeout[mockHandler](tt.timeout))

			if !cmp.Equal(s.GracefulTimeout, tt.timeout) {
				t.Error(cmp.Diff(s.GracefulTimeout, tt.timeout))
			}
		})
	}
}

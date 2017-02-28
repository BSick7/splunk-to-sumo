package agent

import (
	"fmt"
	"net/http"
	"time"
)

type HttpConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (c *HttpConfig) CreateServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(c.Address),
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
		Handler:      handler,
	}
}

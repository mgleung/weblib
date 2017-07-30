package context

import (
	gc "context"
	"fmt"
	"net/http"
	"time"
)

/*
This package defines custom context for a cleaner
creation and for added customizability should it be needed later
*/

func NewCustomContext() gc.Context {
	// TODO: Would use the gc.WithValue to add IDs, loggers, and other information here but it is not needed
	return gc.Background()
}

type HandleResponseFunc func(*http.Response, error) error

type ContextClient interface {
	ExecuteHttp(gc.Context, *http.Request, HandleResponseFunc) error
}

type ContextHttpClient struct {
	client    *http.Client
	transport *http.Transport
	timeout   time.Duration
}

func NewContextHttpClient(transport *http.Transport, timeout time.Duration) ContextClient {
	if transport == nil {
		transport = &http.Transport{}
	}

	if timeout == 0 {
		// Default timeout is 2 seconds
		timeout = time.Duration(2 * time.Second)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &ContextHttpClient{
		client:    client,
		transport: transport,
		timeout:   timeout,
	}
}

func (t ContextHttpClient) ExecuteHttp(ctx gc.Context, req *http.Request, handlerFunc HandleResponseFunc) error {
	ch := make(chan error, 1)
	defer close(ch)

	// Run the actual HTTP request in a go routine and pass the response to the handler function
	go func() {
		ch <- handlerFunc(t.client.Do(req))
	}()

	select {
	case <-ctx.Done():
		t.transport.CancelRequest(req)
		// Wait for the return
		<-ch
		return fmt.Errorf("%v: %s", ctx.Err(), req.URL.String())

	case err := <-ch:
		return err
	}
}

package mdd

import (
	"context"
	"time"
)

type ClientTransport interface {
	Send([]byte) ([]byte, error)
	Close() error
}

type ServerTransport interface {
	Listen() error
	Handler(handler func([]byte) ([]byte, error))
	Close() error
}

type Client struct {
	Codec     Codec
	Transport ClientTransport
}

func (c *Client) SendMessage(request *Containers) (*Containers, error) {

	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	// Hard code 10 second for now. Make it configurable later
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resultCh := make(chan []byte)
	errorCh := make(chan error)
	go func() {
		respBody, err := c.Transport.Send(reqBody)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- respBody
	}()

	select {
	case respBody := <-resultCh:
		response, err := c.Codec.Decode(respBody)
		if err != nil {
			return nil, err
		}
		return response, nil
	case err := <-errorCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type Server struct {
	Codec     Codec
	Transport ServerTransport
}

func (s *Server) MessageHandler(handler func(*Containers) (*Containers, error)) {

	h := func(reqBody []byte) ([]byte, error) {
		request, err := s.Codec.Decode(reqBody)
		if err != nil {
			return nil, err
		}

		response, err := handler(request)
		if err != nil {
			return nil, err
		}

		respBody, err := s.Codec.Encode(response)
		if err != nil {
			return nil, err
		}

		return respBody, nil
	}

	s.Transport.Handler(h)
}

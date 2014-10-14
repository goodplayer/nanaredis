package nanaredis

import (
	"io"
)

const (
	STATE_START = 1
)

type ClientProtocol struct {
	//TODO
	state int
}

type ClientReceived struct {
	//TODO
}

type ClientRequest interface {
	write(o io.Writer) error
}

func createClientProtocol() *ClientProtocol {
	return &ClientProtocol{
		state: STATE_START,
	}
}

func (p *ClientProtocol) process(data []byte) (*ClientReceived, error) {
	//TODO
	switch p.state {

	}
	return nil, nil
}

func (p *ClientProtocol) send(req ClientRequest, output io.Writer) error {
	return req.write(output)
}

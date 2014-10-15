package nanaredis

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	STATE_START = 1 + iota
	STATE_STRING
	STATE_ERROR
	STATE_INTEGER
	STATE_BULK_STRING
	STATE_ARRAY
	STATE_DONE
)

const (
	SUBSTATE_START    = 1 + iota
	SUBSTATE_CONTINUE //need more data
)

type ClientProtocol struct {
	queue chan *resultStruct

	//TODO state machine
	state    int
	substate int

	// for temporary store when event come
	tmp_Received *ClientReceived
}

type ClientReceived struct {
	//TODO
	t int

	// for: string, error
	data []byte
}

func createClientProtocol(queueSize int) *ClientProtocol {
	return &ClientProtocol{
		state:    STATE_START,
		substate: SUBSTATE_START,
		queue:    make(chan *resultStruct, queueSize),
	}
}

func (p *ClientProtocol) process(data []byte) error {
	total := len(data)
	i := 0
	//TODO
	for i < total {
		switch p.state {
		case STATE_START:
			first := data[i]
			i++
			p.tmp_Received = &ClientReceived{}
			p.tmp_Received.t = int(first)
			switch first {
			case '+':
				p.state = STATE_STRING
			case '-':
				p.state = STATE_ERROR
			case ':':
				fmt.Println(":")
			case '$':
				fmt.Println("$")
			case '*':
				fmt.Println("*")
			default:
				return errors.New(fmt.Sprintf("processing first byte error. unknown charactor: %#x", first))
			}
		case STATE_STRING:
			fallthrough
		case STATE_ERROR:
			//TODO
			d, _ := bytes.NewBuffer(data[i:]).ReadBytes('\n')
			i = i + len(d)
			b := make([]byte, len(d)-2)
			copy(b, d[:len(d)-2])
			p.tmp_Received.data = b
			p.state = STATE_DONE
		case STATE_DONE:
			r := <-p.queue
			switch p.tmp_Received.t {
			case '-':
				r.errorCh <- errors.New(string(p.tmp_Received.data))
			default:
				r.resultCh <- p.tmp_Received
			}
			p.tmp_Received = nil
			p.state = STATE_START
			p.substate = SUBSTATE_START
		}
	}
	if p.state == STATE_DONE {
		r := <-p.queue
		switch p.tmp_Received.t {
		case '-':
			r.errorCh <- errors.New(string(p.tmp_Received.data))
		default:
			r.resultCh <- p.tmp_Received
		}
		p.tmp_Received = nil
		p.state = STATE_START
		p.substate = SUBSTATE_START
	}
	return nil
}

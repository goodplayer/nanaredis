package nanaredis

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type RedisClient interface {
	Conn() (RedisConnection, error)
}

type RedisConnection interface {
	Ping() error
}

type redisClientImpl struct {
	hostString string
}

func (r *redisClientImpl) Conn() (RedisConnection, error) {
	addr, err := net.ResolveTCPAddr("tcp", r.hostString)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	clientProtocol := createClientProtocol(10000)

	//TODO
	go func() {
		b := make([]byte, 1000)
		for {
			n, _ := conn.Read(b)
			e := clientProtocol.process(b[:n])
			if e != nil {
				fmt.Println("error!!!!", e)
			}
		}
	}()

	return &redisConnectionImpl{
		conn:           conn,
		bufwriter:      bufio.NewWriter(conn),
		clientProtocol: clientProtocol,
	}, nil
}

type redisConnectionImpl struct {
	lock           sync.Mutex
	conn           *net.TCPConn
	bufwriter      *bufio.Writer
	clientProtocol *ClientProtocol
}

func CreateRedisClient(hostString string) RedisClient {
	return &redisClientImpl{
		hostString: hostString,
	}
}

func (c *redisConnectionImpl) Ping() error {
	r := createResultStruct()
	c.lock.Lock()
	err := sendPing(c.bufwriter)
	if err == nil {
		c.clientProtocol.queue <- r
	}
	c.lock.Unlock()
	if err != nil {
		return err
	}
	select {
	case <-r.resultCh:
		return nil
	case e := <-r.errorCh:
		return e
	}
	panic("should not enter this branch!")
}

type resultStruct struct {
	resultCh chan *ClientReceived
	errorCh  chan error
}

func createResultStruct() *resultStruct {
	return &resultStruct{
		resultCh: make(chan *ClientReceived, 1),
		errorCh:  make(chan error, 1),
	}
}

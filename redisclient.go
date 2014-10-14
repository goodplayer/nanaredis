package nanaredis

import (
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
	return &redisConnectionImpl{
		conn: conn,
	}, nil
}

type redisConnectionImpl struct {
	lock sync.Mutex
	conn *net.TCPConn
}

func CreateRedisClient(hostString string) RedisClient {
	return &redisClientImpl{
		hostString: hostString,
	}
}

func (c *redisConnectionImpl) Ping() error {
	//TODO
	c.lock.Lock()
	c.lock.Unlock()
	return nil
}

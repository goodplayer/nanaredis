package nanaredis

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func ExampleBasicUsageOfNanaRedis() {
	//TODO
}

func TestBasicUse(t *testing.T) {
	redisClient := CreateRedisClient("192.168.1.40:6379")
	conn, err := redisClient.Conn()
	fmt.Println(conn, err)
	err = conn.Ping()
	fmt.Println("ping done!", err)
	time.Sleep(1 * time.Second)
}

func BenchmarkBasicUse(b *testing.B) {
	b.ReportAllocs()

	redisClient := CreateRedisClient("192.168.1.40:6379")
	conn, err := redisClient.Conn()
	fmt.Println(conn, err)

	for i := 0; i < b.N; i++ {
		err = conn.Ping()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkBasicUseMulti(b *testing.B) {
	b.ReportAllocs()

	runtime.GOMAXPROCS(runtime.NumCPU())
	redisClient := CreateRedisClient("192.168.1.40:6379")
	conn, err := redisClient.Conn()
	if err != nil {
		b.Error(err)
	}
	NODE := 10000
	CNT := b.N / NODE
	var waitGroup sync.WaitGroup
	waitGroup.Add(NODE)

	b.ResetTimer()

	for i := 0; i < NODE; i++ {
		go func() {
			defer waitGroup.Done()
			for i := 0; i < CNT; i++ {
				conn.Ping()
			}
		}()
	}
	waitGroup.Wait()
}

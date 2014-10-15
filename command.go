package nanaredis

import (
	"bufio"
)

var (
	REQ_PING = []byte("*1\r\n$4\r\nPING\r\n")
)

func sendPing(writer *bufio.Writer) error {
	_, err := writer.Write(REQ_PING)
	if err != nil {
		return err
	}
	return writer.Flush()
}

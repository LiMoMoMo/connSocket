package socket

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	//LEN message prefix length
	LEN = 4
	//INTERVAL read step length
	INTERVAL = 1024
)

// Socket instance
type Socket struct {
	Closef    context.CancelFunc
	Ctx       context.Context
	Conn      net.Conn
	ReadChan  chan []byte
	WriteChan chan []byte
}

// WriteMsg wait msg to write.
func (s *Socket) WriteMsg() {
	for {
		select {
		case data := <-s.WriteChan:
			pref := intToBytes(len(data))
			var buffer bytes.Buffer
			buffer.Write(pref)
			buffer.Write(data)
			_, err := s.Conn.Write(buffer.Bytes())
			if err != nil {
				fmt.Println("Send Error,", err)
			}
		case <-s.Ctx.Done():
			fmt.Println("Quit WriteMsg()")
			return
		}
	}
}

// ReadMsg wait msg to read.
func (s *Socket) ReadMsg() {
	tmpBuffer := make([]byte, 0)
	data := make([]byte, INTERVAL)
	for {
		// get length
		n, err := s.Conn.Read(data)
		if err != nil {
			s.Conn.Close()
			fmt.Println("Conn has been Closed.")
			s.Closef()
			break
		}
		tmpBuffer = s.unpack(append(tmpBuffer, data[:n]...))
	}
}

func (s *Socket) unpack(buffer []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+LEN {
			break
		}
		messageLength := bytesToInt(buffer[i : i+LEN])
		if length < i+LEN+messageLength {
			break
		}
		data := buffer[i+LEN : i+LEN+messageLength]
		i += LEN + messageLength - 1
		s.ReadChan <- data
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

func bytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func intToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

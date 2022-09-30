package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("libcamera-still", "-o-", "-s", "-t0", "-ejpg")

	cmd.Stdout = MakeWriter(time.Millisecond * 250)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	log.Printf("PID: %d \n", cmd.Process.Pid)

	cmd.Wait()
}

func MakeWriter(d time.Duration) *CollectiveWriter {
	c := &CollectiveWriter{CollectionTime: d}

	c.buf = &bytes.Buffer{}
	c.data = make(chan []byte, 100)
	c.ok = make(chan struct{})

	go func(w *CollectiveWriter) {
		for {
			select {
			case data := <-w.data:
				log.Printf("First received byte though channel: %v \n", data[0])
				_, err := c.buf.Write(data)
				w.ok <- struct{}{}
				if err != nil {
					log.Printf("Error writing first []byte: %s", err)
				}
				timeout := time.NewTimer(c.CollectionTime).C

			inner:
				for {
					select {
					case data := <-w.data:
						c.buf.Write(data)
						w.ok <- struct{}{}
					case <-timeout:
						break inner
					}
				}

				// dump
				b := c.buf.Bytes()

				file := fmt.Sprintf("out-%d.jpg", time.Now().Unix())
				log.Printf("Data received has length: %d \n filename: %s \n first byte: %v \n", len(b), file, b[0])
				os.WriteFile(file, b, 0755)
				c.buf.Reset()
			}
		}
	}(c)

	return c
}

type CollectiveWriter struct {
	CollectionTime time.Duration

	data chan []byte
	ok   chan struct{}
	buf  *bytes.Buffer
}

func (w *CollectiveWriter) Write(b []byte) (n int, err error) {
	log.Printf("First bit is %v \n", b[0])
	w.data <- b
	<-w.ok
	log.Printf("First bit is %v \n", b[0])

	return len(b), nil
}

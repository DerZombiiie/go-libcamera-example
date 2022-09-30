package collector

import (
	"bytes"
	"io"
	"log"
	"time"
)

func MakeCollector(d time.Duration, callback func(*bytes.Buffer)) *WriteCollector {
	c := &WriteCollector{timeout: d, cb: callback}

	c.buf = &bytes.Buffer{}
	c.data = make(chan []byte, 100)
	c.ok = make(chan struct{})

	go func(w *WriteCollector) {
		for {
			select {
			// wait for any data
			case data := <-w.data:
				_, err := c.buf.Write(data)
				w.ok <- struct{}{}
				if err != nil {
					log.Printf("Error writing first []byte: %s", err)
				}

				// create timeout channel after which data gets flushed to file
				timeout := time.NewTimer(c.timeout).C

			inner:
				for {
					select {
					case data := <-w.data: // handle more data
						c.buf.Write(data)
						w.ok <- struct{}{}
					case <-timeout: // 250ms are over, start dumping
						break inner
					}
				}

				// dump
				if w.cb != nil {
					buf2 := &bytes.Buffer{}
					io.Copy(buf2, w.buf)

					go w.cb(buf2)
				}

				w.buf.Reset()
			}
		}
	}(c)

	return c
}

type WriteCollector struct {
	timeout time.Duration

	data chan []byte
	ok   chan struct{}
	buf  *bytes.Buffer

	cb func(*bytes.Buffer)
}

func (w *WriteCollector) Write(b []byte) (n int, err error) {
	w.data <- b
	<-w.ok // the ok channel is needed for parts to arrive in correct order, having a bufferd channel isn't enough (even with 100)

	return len(b), nil
}

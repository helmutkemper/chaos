package networkdelay

import (
	"errors"
	"io"
	"time"
)

func (e *Proxy) copyContent(closer chan bool, dst io.Writer, src io.Reader, direction string) {
	var err error
	var buf []byte
	var errInvalidWrite = errors.New("invalid write result")
	var ErrShortWrite = errors.New("short write")

	_ = err

	buf = make([]byte, e.bufferSize)

	for {

		if e.min != 0 && e.max != 0 {
			time.Sleep(time.Duration(e.intRange(e.min, e.max)) * time.Millisecond)
		}

		nr, er := src.Read(buf)

		if e.parser != nil {
			nr, err = e.parser.Parser(buf[0:nr], direction)
		}

		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	closer <- true
}

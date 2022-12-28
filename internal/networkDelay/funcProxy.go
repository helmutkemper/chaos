package networkdelay

import (
	"log"
	"net"
)

func (e *Proxy) Proxy(inStringConn, outStringConn string) (err error) {
	var listener net.Listener
	var listenerConn net.Conn
	var dialConn net.Conn

	if e.bufferSize == 0 {
		e.bufferSize = 512
	}

	listener, err = net.Listen("tcp", inStringConn)
	if err != nil {
		return
	}

	for {
		listenerConn, err = listener.Accept()
		if err != nil {
			log.Printf("error accepting connection. address: %v, error: %v", listenerConn.RemoteAddr(), err)
			continue
		}

		go func() {
			defer listenerConn.Close()
			dialConn, err = net.Dial("tcp", outStringConn)
			if err != nil {
				log.Printf("error dialing remote connection. address: %v, error: %v", dialConn.RemoteAddr(), err)
				return
			}
			defer dialConn.Close()
			var closer = make(chan bool, 2)
			go e.copyContent(closer, dialConn, listenerConn, "out")
			go e.copyContent(closer, listenerConn, dialConn, "in")
			<-closer
		}()
	}
}

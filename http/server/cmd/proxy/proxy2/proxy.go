package proxy2

import "net"

type Server struct {
	listener   net.Listener
	addr       string
	credential string
}

// Start a proxy server
func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		servLogger.Fatal(err)
	}

	if s.credential != "" {
		servLogger.Infof("use %s for auth\n", s.credential)
	}
	servLogger.Infof("proxy listen in %s, waiting for connection...\n", s.addr)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			servLogger.Error(err)
			continue
		}
		go s.newConn(conn).serve()
	}
}

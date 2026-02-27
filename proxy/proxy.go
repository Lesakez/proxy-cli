package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Lesakez/proxy-cli/config"
	"github.com/Lesakez/proxy-cli/filter"
	"github.com/Lesakez/proxy-cli/logger"
)

type Server struct {
	addr    string
	proxies []config.ProxyConfig
}

func NewServer(addr string, proxies []config.ProxyConfig) *Server {
	return &Server{addr: addr, proxies: proxies}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.addr, err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(clientConn net.Conn) {
	defer clientConn.Close()

	clientConn.SetDeadline(time.Now().Add(30 * time.Second)) //nolint:errcheck

	req, err := http.ReadRequest(bufio.NewReader(clientConn))
	if err != nil {
		if err != io.EOF {
			logger.Error("read", err)
		}
		return
	}

	if req.Method != http.MethodConnect {
		writeResponse(clientConn, http.StatusMethodNotAllowed, "only CONNECT supported")
		return
	}

	targetAddr := req.Host
	if !strings.Contains(targetAddr, ":") {
		targetAddr += ":443"
	}

	host := targetAddr
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	clientConn.SetDeadline(time.Time{}) //nolint:errcheck

	upstream := filter.FindProxy(host, s.proxies)

	var targetConn net.Conn

	if upstream != nil {
		logger.Proxy(targetAddr, upstream.Scheme, upstream.Name, upstream.Addr())
		targetConn, err = Connect(targetAddr, upstream)
	} else {
		logger.Direct(targetAddr)
		targetConn, err = net.DialTimeout("tcp", targetAddr, dialTimeout)
	}

	if err != nil {
		logger.Error(targetAddr, err)
		writeResponse(clientConn, http.StatusBadGateway, err.Error())
		return
	}
	defer targetConn.Close()

	writeResponse(clientConn, http.StatusOK, "Connection established")
	Pipe(targetConn, clientConn)
}

func writeResponse(conn net.Conn, code int, msg string) {
	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n\r\n", code, msg) //nolint:errcheck
}

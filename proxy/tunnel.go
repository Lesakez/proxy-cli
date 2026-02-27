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
	"golang.org/x/net/proxy"
)

const dialTimeout = 15 * time.Second

func Connect(targetAddr string, p *config.ProxyConfig) (net.Conn, error) {
	switch strings.ToUpper(p.Scheme) {
	case "SOCKS5":
		return connectViaSocks5(targetAddr, p)
	default:
		return connectViaHTTP(targetAddr, p)
	}
}

func connectViaSocks5(targetAddr string, p *config.ProxyConfig) (net.Conn, error) {
	var auth *proxy.Auth
	if p.HasCredentials() {
		auth = &proxy.Auth{
			User:     p.Auth.Credentials.Username,
			Password: p.Auth.Credentials.Password,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", p.Addr(), auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("socks5 %s: %w", p.Addr(), err)
	}

	conn, err := dialer.Dial("tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("socks5 dial %s: %w", targetAddr, err)
	}

	return conn, nil
}

func connectViaHTTP(targetAddr string, p *config.ProxyConfig) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", p.Addr(), dialTimeout)
	if err != nil {
		return nil, fmt.Errorf("connect %s: %w", p.Addr(), err)
	}

	req, err := http.NewRequest(http.MethodConnect, "http://"+targetAddr, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}
	req.Host = targetAddr

	if p.HasCredentials() {
		req.SetBasicAuth(p.Auth.Credentials.Username, p.Auth.Credentials.Password)
	}
	if p.Auth.Token != "" {
		req.Header.Set("Proxy-Authorization", "Bearer "+p.Auth.Token)
	}

	if err := req.Write(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("connect write: %w", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("connect response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("proxy: %s", resp.Status)
	}

	return conn, nil
}

type closeWriter interface {
	CloseWrite() error
}

func tryCloseWrite(conn net.Conn) {
	if cw, ok := conn.(closeWriter); ok {
		cw.CloseWrite() //nolint:errcheck
	} else {
		conn.Close()
	}
}

func Pipe(dst, src net.Conn) {
	done := make(chan struct{}, 2)

	go func() {
		io.Copy(dst, src) //nolint:errcheck
		tryCloseWrite(dst)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(src, dst) //nolint:errcheck
		tryCloseWrite(src)
		done <- struct{}{}
	}()

	<-done
	<-done
}

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	host := flag.String("host", "127.0.0.1", "HTTP/2 server host")
	port := flag.Int("port", 8081, "HTTP/2 server port")
	path := flag.String("path", "/", "request path")
	tlsEnabled := flag.Bool("tls", false, "use TLS with HTTP/2 ALPN")
	insecure := flag.Bool("insecure", false, "skip TLS certificate verification for a local test server")
	flag.Parse()

	if *port < 1 || *port > 65535 {
		fail("port must be between 1 and 65535")
	}
	if !strings.HasPrefix(*path, "/") {
		fail("path must begin with /")
	}

	address := net.JoinHostPort(*host, strconv.Itoa(*port))
	transport := &http2.Transport{}
	if *tlsEnabled {
		transport.TLSClientConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: *insecure, // Explicit opt-in for local self-signed certificates.
		}
	} else {
		transport.AllowHTTP = true
		dialer := &net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}
		transport.DialTLSContext = func(ctx context.Context, network, _ string, _ *tls.Config) (net.Conn, error) {
			return dialer.DialContext(ctx, network, address)
		}
	}
	defer transport.CloseIdleConnections()

	scheme := "http"
	if *tlsEnabled {
		scheme = "https"
	}
	url := scheme + "://" + address + *path
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	for _, test := range []struct {
		method string
		body   []byte
	}{{http.MethodGet, nil}, {http.MethodPost, []byte("typescript-http2-probe")}} {
		request, err := http.NewRequestWithContext(context.Background(), test.method, url, bytes.NewReader(test.body))
		if err != nil {
			fail("create %s request: %v", test.method, err)
		}
		response, err := client.Do(request)
		if err != nil {
			fail("%s request failed: %v", test.method, err)
		}
		body, readErr := io.ReadAll(io.LimitReader(response.Body, 1<<20))
		closeErr := response.Body.Close()
		if readErr != nil {
			fail("read %s response: %v", test.method, readErr)
		}
		if closeErr != nil {
			fail("close %s response: %v", test.method, closeErr)
		}
		if response.ProtoMajor != 2 {
			fail("%s negotiated %q; expected HTTP/2", test.method, response.Proto)
		}
		if response.StatusCode != http.StatusOK {
			fail("%s returned %s; expected 200", test.method, response.Status)
		}
		if len(body) == 0 {
			fail("%s returned an empty body", test.method)
		}
		fmt.Printf("HTTP/2 %s passed: %s, %d response bytes\n", test.method, response.Proto, len(body))
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

package capture

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Proxy is an HTTP/HTTPS proxy that captures AWS API calls.
type Proxy struct {
	capturer *Capturer
	listener net.Listener
	server   *http.Server
	verbose  bool
	wg       sync.WaitGroup
}

// NewProxy creates a new capture proxy.
func NewProxy(capturer *Capturer, verbose bool) *Proxy {
	return &Proxy{
		capturer: capturer,
		verbose:  verbose,
	}
}

// Start starts the proxy on a random available port.
func (p *Proxy) Start() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("failed to start listener: %w", err)
	}
	p.listener = listener

	p.server = &http.Server{
		Handler: p,
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.server.Serve(listener)
	}()

	return listener.Addr().String(), nil
}

// Stop stops the proxy.
func (p *Proxy) Stop() error {
	if p.server != nil {
		p.server.Close()
	}
	p.wg.Wait()
	return nil
}

// ServeHTTP handles proxy requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "CONNECT" {
		p.handleConnect(w, req)
	} else {
		p.handleHTTP(w, req)
	}
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, req *http.Request) {
	// Read and capture the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Capture AWS call
	if call := ParseAWSRequest(req, bodyBytes); call != nil {
		p.capturer.AddCall(*call)
		if p.verbose {
			fmt.Printf("[capture] %s:%s on %s\n", call.Service, call.Action, call.Resource)
		}
	}

	// Forward the request
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	outReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Copy headers
	for k, vv := range req.Header {
		for _, v := range vv {
			outReq.Header.Add(k, v)
		}
	}

	resp, err := client.Do(outReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *Proxy) handleConnect(w http.ResponseWriter, req *http.Request) {
	// For HTTPS, we tunnel the connection and capture at TCP level
	// This is a simple pass-through tunnel - we capture based on host only

	host := req.Host
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	// Record that we saw a connection to this AWS service
	if strings.Contains(host, "amazonaws.com") {
		service, region := parseAWSHost(strings.Split(host, ":")[0])
		if service != "" {
			// We can't see the actual action in CONNECT tunnels without MITM
			// So we just note the service was accessed
			if p.verbose {
				fmt.Printf("[capture] HTTPS connection to %s (region: %s) - action unknown without MITM\n", service, region)
			}
		}
	}

	// Establish connection to target
	targetConn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer targetConn.Close()

	// Hijack the client connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	// Send 200 Connection Established
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Tunnel bidirectionally
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(targetConn, clientConn)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientConn, targetConn)
	}()

	wg.Wait()
}

// ProxyWithMITM is a proxy that performs MITM to capture HTTPS traffic.
type ProxyWithMITM struct {
	*Proxy
	tlsConfig *tls.Config
}

// Note: Full MITM implementation requires CA certificate generation
// and is complex. For MVP, we'll use a simpler approach:
// Parse AWS CLI debug output or use SDK-specific instrumentation.

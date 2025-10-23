package old

import (
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// TestNATSServer provides a reusable embedded NATS server for testing
type TestNATSServer struct {
	Server *server.Server
	Conn   *nats.Conn
}

// NewTestNATSServer creates and starts a new embedded NATS server
func NewTestNATSServer(t *testing.T) *TestNATSServer {
	t.Helper()

	// Create a new NATS server with minimal configuration based on NATS test patterns
	opts := &server.Options{
		Host:                  "127.0.0.1",
		Port:                  -1, // Use random available port (NATS uses -1 not 0)
		NoLog:                 true,
		NoSigs:                true,
		Debug:                 false,
		Trace:                 false,
		JetStream:             false,
		DisableShortFirstPing: true,
		MaxControlLine:        4096,
	}

	// Use NewServer instead of deprecated New
	ns, err := server.NewServer(opts)
	if err != nil {
		t.Fatalf("Failed to create NATS server: %v", err)
	}
	if ns == nil {
		t.Fatal("NATS server is nil")
	}

	// Start the server in a goroutine
	go ns.Start()

	// Wait for server to be ready with timeout
	if !ns.ReadyForConnections(2 * time.Second) {
		ns.Shutdown()
		t.Fatalf("NATS server failed to start at %s", ns.Addr())
	}

	// Connect to the server using the client URL from the server
	url := ns.ClientURL()
	nc, err := nats.Connect(url)
	if err != nil {
		ns.Shutdown()
		t.Fatalf("Failed to connect to NATS server at %s: %v", url, err)
	}

	// Verify connection is ready
	if !nc.IsConnected() {
		ns.Shutdown()
		t.Fatal("NATS connection is not ready")
	}

	return &TestNATSServer{
		Server: ns,
		Conn:   nc,
	}
}

// Shutdown stops the NATS server and closes the connection
func (ts *TestNATSServer) Shutdown(t *testing.T) {
	t.Helper()

	if ts.Conn != nil {
		ts.Conn.Close()
	}
	if ts.Server != nil {
		ts.Server.Shutdown()
	}
}

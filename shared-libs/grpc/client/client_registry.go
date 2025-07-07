package grpc

import (
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ClientRegistry struct {
	connections map[string]*grpc.ClientConn
	mu          sync.RWMutex
}

func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		connections: make(map[string]*grpc.ClientConn),
	}
}

func (r *ClientRegistry) GetConnection(target string) (*grpc.ClientConn, error) {
	r.mu.RLock()
	if conn, exists := r.connections[target]; exists {
		r.mu.RUnlock()
		return conn, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// double-check after acquiring write lock
	if conn, exists := r.connections[target]; exists {
		return conn, nil
	}

	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // note: should be used only on private network
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	r.connections[target] = conn
	return conn, nil
}

func (r *ClientRegistry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for target, conn := range r.connections {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection to %s: %v", target, err)
		}
	}
	r.connections = make(map[string]*grpc.ClientConn)
}

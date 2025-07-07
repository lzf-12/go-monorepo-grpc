package grpc

import "time"

// default limit params used in grpc server
const (
	DefaultMaxRecvMessageSize   int = 1024 * 1024
	DefaultMaxSendMessageSize   int = 1024 * 1024
	DefaultConcurrentUpstream   int = 100
	DefaultKeepAliveMaxIdleTime     = 3 * time.Minute
	DefaultKeepAliveTime            = 5 * time.Second
	DefaultKeepAliveTimeout         = 1 * time.Second
)

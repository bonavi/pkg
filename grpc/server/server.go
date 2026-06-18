package server

import (
	"google.golang.org/grpc"
)

type ServerOptions struct {
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	MaxRecvMsgSize     int // Максимальный размер входящего сообщения в байтах (0 — по умолчанию 4 МБ)
	MaxSendMsgSize     int
}

func NewGRPCServer(opts *ServerOptions) *grpc.Server {
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(opts.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(opts.StreamInterceptors...),
	}

	if opts.MaxRecvMsgSize > 0 {
		grpcOpts = append(grpcOpts, grpc.MaxRecvMsgSize(opts.MaxRecvMsgSize))
	}
	if opts.MaxSendMsgSize > 0 {
		grpcOpts = append(grpcOpts, grpc.MaxSendMsgSize(opts.MaxSendMsgSize))
	}

	return grpc.NewServer(grpcOpts...)
}

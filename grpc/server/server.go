package server

import (
	"google.golang.org/grpc"
)

type ServerOptions struct {
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	MaxRecvMsgSize     int // Максимальный размер входящего сообщения в байтах (0 — по умолчанию 4 МБ)
}

func NewGRPCServer(opts *ServerOptions) *grpc.Server {
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(opts.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(opts.StreamInterceptors...),
	}

	if opts.MaxRecvMsgSize > 0 {
		grpcOpts = append(grpcOpts, grpc.MaxRecvMsgSize(opts.MaxRecvMsgSize))
	}

	return grpc.NewServer(grpcOpts...)
}

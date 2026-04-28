package server

import (
	"google.golang.org/grpc"
)

type ServerOptions struct {
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

func NewGRPCServer(opts *ServerOptions) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(opts.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(opts.StreamInterceptors...),
	)

	return grpcServer
}

package grpc

import (
	"google.golang.org/grpc"
	"grpc-auth/internal/web/gen"
)

type server struct {
	//auth.UnimplementedAuthServer
}

func RegisterServer(grpcServer *grpc.Server) {
	auth.RegisterAuthServer(grpcServer, &server{})
}

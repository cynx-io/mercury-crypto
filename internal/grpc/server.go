package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "mercury/api/proto/crypto"
	"mercury/internal/module/cryptomodule"
	"mercury/internal/pkg/logger"
	"net"
)

type Server struct {
	pb.UnimplementedMercuryCryptoServiceServer
	cryptoService *cryptomodule.CryptoService
}

func NewServer(cryptoService *cryptomodule.CryptoService) *Server {
	return &Server{
		cryptoService: cryptoService,
	}
}

func (s *Server) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterMercuryCryptoServiceServer(server, s)
	reflection.Register(server)

	logger.Infof("Starting gRPC server on %s", address)
	return server.Serve(lis)
}

func (s *Server) SearchCoin(ctx context.Context, req *pb.SearchCoinRequest) (*pb.SearchCoinResponse, error) {
	resp, err := s.cryptoService.SearchCoin(ctx, req.Query)
	return serviceWrapper(resp, err)
}

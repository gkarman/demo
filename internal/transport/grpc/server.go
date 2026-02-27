package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/gkarman/demo/api/gen/go/car/v1"
	carservice "github.com/gkarman/demo/internal/service/car"
	"github.com/gkarman/demo/internal/service/car/requestdto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Addr string
}

type Server struct {
	grpcServer *grpc.Server
	lis        net.Listener
	log        *slog.Logger
}

func NewServer(log *slog.Logger, cfg Config) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	s := grpc.NewServer()
	car.RegisterCarServer(s, &carServer{})

	return &Server{
		grpcServer: s,
		lis:        lis,
		log:        log,
	}, nil
}

func (s *Server) Start() {
	go func() {
		s.log.Info("gRPC server started", "addr", s.lis.Addr())
		if err := s.grpcServer.Serve(s.lis); err != nil {
			s.log.Error("gRPC server failed", "error", err)
		}
	}()
}

func (s *Server) Stop() {
	s.log.Info("stopping gRPC server")
	s.grpcServer.GracefulStop()
}

type carServer struct {
	car.UnimplementedCarServer
	getService *carservice.GetService
}

func (s *carServer) GetCar(ctx context.Context, req *car.GetCarRequest) (*car.GetCarResponse, error) {
	resp, err := s.getService.Execute(ctx, &requestdto.GetCar{CarId: req.Id})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &car.GetCarResponse{
		Id:   resp.Car.ID,
		Name: resp.Car.Name,
	}, nil
}

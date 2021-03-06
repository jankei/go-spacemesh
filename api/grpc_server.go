package api

import (
	"github.com/spacemeshos/go-spacemesh/api/config"
	"github.com/spacemeshos/go-spacemesh/api/pb"
	"github.com/spacemeshos/go-spacemesh/log"
	"strconv"

	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// A grpc server implementing the SpaceMesh API

// server is used to implement SpaceMeshService.Echo.
type SpaceMeshGrpcService struct {
	Server *grpc.Server
	Port   uint
}

func (s SpaceMeshGrpcService) Echo(ctx context.Context, in *pb.SimpleMessage) (*pb.SimpleMessage, error) {
	return &pb.SimpleMessage{in.Value}, nil
}

func (s SpaceMeshGrpcService) StopService() {
	log.Info("Stopping grpc service...")
	s.Server.Stop()
}

func NewGrpcService() *SpaceMeshGrpcService {
	port := config.ConfigValues.GrpcServerPort
	server := grpc.NewServer()
	return &SpaceMeshGrpcService{Server: server, Port: port}
}

func (s SpaceMeshGrpcService) StartService() {
	go s.startServiceInternal()
}

// This is a blocking method designed to be called using a go routine
func (s SpaceMeshGrpcService) startServiceInternal() {
	port := config.ConfigValues.GrpcServerPort
	addr := ":" + strconv.Itoa(int(port))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("failed to listen: %v", err)
		return
	}

	pb.RegisterSpaceMeshServiceServer(s.Server, s)

	// Register reflection service on gRPC server
	reflection.Register(s.Server)

	log.Info("grpc API listening on port %d", port)

	// start serving - this blocks until err or server is stopped
	if err := s.Server.Serve(lis); err != nil {
		log.Error("failed to serve grpc: %v", err)
	}
}

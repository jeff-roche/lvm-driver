package services

import (
	"net"
	"os"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/openshift/lvm-driver/pkg/lvmdriver/utils"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

type GrpcServer interface {
	// start the service
	Start()
	// Graceful shutdown
	Stop()
	// Forced shutdown
	ForceStop()
}

type GrpcServerConfig struct {
	Endpoint   string
	IdServer   csi.IdentityServer
	NodeServer csi.NodeServer
}

// GrpcServer is the primary server for all k8s related communications
type grpcServer struct {
	wg         sync.WaitGroup
	server     *grpc.Server
	endpoint   string
	idServer   csi.IdentityServer
	nodeServer csi.NodeServer
}

func NewGrpcServer(config GrpcServerConfig) GrpcServer {
	return &grpcServer{
		endpoint:   config.Endpoint,
		idServer:   config.IdServer,
		nodeServer: config.NodeServer,
	}
}

func (s *grpcServer) Start() {
	s.wg.Add(1)
	go s.serve()

	s.wg.Wait()
}

func (s *grpcServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

func (s *grpcServer) ForceStop() {
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *grpcServer) serve() {
	proto, addr, err := utils.ParseEndpoint(s.endpoint)
	if err != nil {
		klog.Fatal(err.Error())
	}

	// Handle grpc connections over unix sockets
	if proto == "unix" {
		addr = "/" + addr
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			klog.Fatalf("Failed to remove %s, error: %s", addr, err.Error())
		}
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(utils.GRPCLogger),
	}
	s.server = grpc.NewServer(opts...)

	if s.idServer != nil {
		csi.RegisterIdentityServer(s.server, s.idServer)
	}

	if s.nodeServer != nil {
		csi.RegisterNodeServer(s.server, s.nodeServer)
	}

	klog.Infof("Listening for connections on address: %#v", listener.Addr())

	err = s.server.Serve(listener)
	if err != nil {
		klog.Fatalf("Failed to serve grpc server: %v", err)
	}
}

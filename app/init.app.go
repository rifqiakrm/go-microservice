package app

import (
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/rifqiakrm/go-microservice/pb/another"
	"github.com/rifqiakrm/go-microservice/pb/sample"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var (
	DB *sql.DB
)

// NewSample returns a new server
func NewSample(tr opentracing.Tracer, an *grpc.ClientConn) *Sample {
	return &Sample{
		tracer:        tr,
		anotherClient: another.NewAnotherServiceClient(an),
	}
}

// Server implements the auth_another_another service
type Sample struct {
	sample.UnimplementedSampleServiceServer
	tracer        opentracing.Tracer
	anotherClient another.AnotherServiceClient
}

func (s *Sample) Run(port int) error {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.tracer),
		),
	)

	sample.RegisterSampleServiceServer(srv, s)

	reflection.Register(srv)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return fmt.Errorf("failed to listen : %v", err)
	}
	log.Println("gRPC started!")
	log.Printf("server is running on %v:%v \n", viper.GetString("app.host"), viper.GetString("app.port"))
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve rpc : %v", err)
	}

	return nil
}

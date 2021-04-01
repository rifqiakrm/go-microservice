package app

import (
	"context"
	"fmt"
	"github.com/rifqiakrm/go-microservice/model"
	"github.com/rifqiakrm/go-microservice/pb/another"
	"github.com/rifqiakrm/go-microservice/pb/sample"
	"github.com/rifqiakrm/go-microservice/resources"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Sample) GetHelloWorld(ctx context.Context, req *sample.GetHelloWorldRequest) (*sample.GetHelloWorldResponse, error) {
	msg := fmt.Sprintf("Hi my name is %s. Hello world.", req.GetName())

	_, err := model.InsertSample(ctx, DB, &resources.Sample{
		Message: msg,
	})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	//get from another rpc
	_, err = s.anotherClient.GetHelloWorld(ctx, &another.GetHelloWorldRequest{Name: "Sample"})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &sample.GetHelloWorldResponse{
		Message: msg,
	}, nil
}

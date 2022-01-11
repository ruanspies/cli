package main

import (
	"context"
	pb "go.protobuf.{{.Organisation}}.alis.exchange/{{.Organisation}}/{{.Product}}/{{.Contract}}/{{.Neuron}}/{{.VersionMajor}}"
)

// Create a Service object which we'll register with the Server
type myService struct {
	pb.UnimplementedServiceServer
}

// TODO: Implement all the methods as per the proto
func (s *myService) MyMethod1(ctx context.Context, req *pb.MyMethod1Request) (*pb.MyMethod1Response, error) {

	return &pb.MyMethod1Response{}, nil
}

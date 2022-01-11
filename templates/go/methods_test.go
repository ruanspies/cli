package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "go.protobuf.{{.Organisation}}.alis.exchange/{{.Organisation}}/{{.Product}}/{{.Contract}}/{{.Neuron}}/{{.VersionMajor}}"
)

// Simulate a client object
var client myService

// This init() function will only run when running Go tests.
func init(){
	// Include a link to the file location of where the log originated from
	log.SetFlags(log.Lshortfile)

	client = myService{}
}

func TestServiceService_MyMethod1(t *testing.T) {

	// Construct a request message
	req := pb.MyMethod1Request{}

	// Run a method
	res, err := client.MyMethod1(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

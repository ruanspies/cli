package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"

	"context"
	"log"
	"net"
	"os"

	alis "go.lib.alis.dev"
	"google.golang.org/grpc"

	pb "go.protobuf.{{.Organisation}}.alis.exchange/{{.Organisation}}/{{.Product}}/{{.Contract}}/{{.Neuron}}/{{.VersionMajor}}"
)

// client is a global client, initialized once per cloud run instance.
var (
	firestoreClient *firestore.Client
	bigqueryClient *bigquery.Client
)

func init() {

	// Pre-declare err to avoid shadowing.
	var err error

	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)

	// Retrieve project id from the environment.
	projectID := os.Getenv("ALIS_OS_PROJECT")
	if projectID == "" {
		log.Fatal("ALIS_OS_PROJECT env not set.")
	}

	// TODO: add/remove required clients.
	// Initialise Bigquery client
	bigqueryClient, err = bigquery.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf( "bigquery.NewClient: %v", err)
	}

	// Initialise Firestore client
	firestoreClient, err = firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
}

func main() {
	log.Println(&Entry{Message: "starting server...", Severity: LogSeverity_NOTICE})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println(&Entry{Message: "Defaulting to port " + port, Severity: LogSeverity_WARNING})
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	pb.RegisterServiceServer(grpcServer, &myService{})

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
)

type grpcTokenSource struct {
	oauth.TokenSource
}

// NewServerConnection creates a new gRPC connection to one of alis_' services.
//
// The host should be the domain where the Service is hosted,:
// Format: v1.{neuron}.{resources|services}.{deployment-project}.{domain}.alis.dev
//
// This method also uses the Google Default Credentials workflow.  To run this locally ensure that you have the
// environmental variable GOOGLE_APPLICATION_CREDENTIALS = ../key.json set.
//
// The location defaults to "europe-west1" if not set, i.e. if empty ""
//
// Best practise is to create a new connection at global level, which could be used to run many methods.  This avoids
// unnecessary api calls to retrieve the required ID tokens each time a single method is called.
func NewServerConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {

	tokenSource, err := IDTokenTokenSource(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"NewTokenSource: %s", err,
		)
	}

	// Establishes a connection
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host+":443"))
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	opts = append(opts, grpc.WithPerRPCCredentials(grpcTokenSource{
		TokenSource: oauth.TokenSource{
			tokenSource,
		},
	}))

	// Enable tracing
	//opts = append(opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	//opts = append(opts, grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	conn, err := grpc.Dial(host+":443", opts...)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"grpc.Dail: %s", err,
		)
	}

	return conn, nil
}

func IDTokenTokenSource(ctx context.Context) (oauth2.TokenSource, error) {

	// Get the token for the authorized_user (not the service_account since this CLI is use by users and machines)
	gts, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, err
	}
	ts := oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: gts})

	return ts, nil
}

// idTokenSource is an oauth2.TokenSource that wraps another
// It takes the id_token from TokenSource and passes that on as a bearer token
type idTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (s *idTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("token did not contain an id_token")
	}

	return &oauth2.Token{
		AccessToken: idToken,
		TokenType:   "Bearer",
		Expiry:      token.Expiry,
	}, nil
}

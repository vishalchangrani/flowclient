package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/onflow/flow/protobuf/go/flow/legacy/access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main2() {

	// determine the Access node IP address and port
	accessNodeAddress := "localhost:9000"



	//tlsCredentials, err := loadTLSCredentials()
	//if err != nil {
	//	log.Fatal("cannot load TLS credentials: ", err)
	//}

	// dial to it using GRPC
	//conn, err := grpc.Dial(accessNodeAddress,  grpc.WithTransportCredentials(tlsCredentials))
	//if err != nil {
	//	panic("failed to connect to access node")
	//}

	config := &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate:VerifyPeerCertificate,
	}
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	ctx := context.Background()

	request := access.GetLatestBlockRequest{
	}

	response, err := grpcClient.GetLatestBlock(ctx, &request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Block.String())
}



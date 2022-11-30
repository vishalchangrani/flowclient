package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"google.golang.org/grpc"
)

func mainsdfsdf() {

	// use the archive API endpoint
	accessNodeAddress := "archive.mainnet.nodes.onflow.org:9000"

	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	ctx := context.Background()

	addr := flow.HexToAddress("2cf3c07fe32957f6")
	// create a GetAccount request
	request := access.GetAccountAtBlockHeightRequest{
		Address:     addr.Bytes(), // set the appropriate account address
		BlockHeight: 40171634,     // start height of the current spork as per https://developers.flow.com/nodes/node-operation/past-sporks#mainnet-20
	}

	// call the GetAccount gRPC method
	// https://github.com/dapperlabs/flow/blob/master/docs/access-api-spec.md#getAccount
	response, err := grpcClient.GetAccountAtBlockHeight(ctx, &request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.GetAccount().GetBalance())
}

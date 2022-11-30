package main

import (
	"context"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"google.golang.org/grpc"
)

func mainn(){

	// determine the Access node IP address and port
	accessNodeAddress := "localhost:8187"

	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	ctx := context.Background()

	transaction := &entities.Transaction{
		Script:               nil,
		ReferenceBlockId:     nil,
		//PayerAccount:         nil,
		//ScriptAccounts:       nil,
		//Signatures:           nil,
		//Status:               0,
	}

	// create a Transaction request
	request := access.SendTransactionRequest{
		Transaction: transaction,
	}

	// call the SendTransaction GRPC method
	// https://github.com/dapperlabs/flow/blob/master/docs/access-api-spec.md#sendtransaction
	_, err = grpcClient.SendTransaction(ctx, &request)
	if err != nil {
		panic("failed to send transaction")
	}


}




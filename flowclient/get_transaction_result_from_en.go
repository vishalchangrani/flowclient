package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/execution"
	"google.golang.org/grpc"
)

//var missingTransactions = []struct {
//	blockHeight uint64
//	txID        string
//}{
//	{64920177, "ddacb628e80552cf7db2d75ca3318ac56622f5d1a66b3517f3a1d85cc8570dbb"},
//}

func mainsdfxc() {

	accessNodeAddress := "access.mainnet.nodes.onflow.org:9000"

	MaxGRPCMessageSize := 1024 * 1024 * 20 // 20MB
	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic("failed to connect to access node")
	}

	ctx := context.Background()

	// create an AccessAPIClient
	grpcAccessClient := access.NewAccessAPIClient(conn)

	//getTransactionFromAN(ctx, "aa890ff09415b12005a0c233d09282abec26aaa8c42990259f3b4fc30d50f0d2", grpcAccessClient)

	executionNodeAddress := "execution-001.mainnet20.nodes.onflow.org:9000"

	// dial to it using GRPC
	conn2, err := grpc.Dial(executionNodeAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxGRPCMessageSize)))
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an ExecutionAPIClient
	grpcExecutionClient := execution.NewExecutionAPIClient(conn2)

	err = getTransactionFromEN(ctx, 41881384, "aa890ff09415b12005a0c233d09282abec26aaa8c42990259f3b4fc30d50f0d2", grpcAccessClient, grpcExecutionClient)
	if err == nil {
		panic(err)
	}

}

func getTransactionFromEN(ctx context.Context, height uint64, txIDHex string, grpcAccessClient access.AccessAPIClient, grpcExecutionClient execution.ExecutionAPIClient) error {
	//req := &access.GetBlockByHeightRequest{
	//	Height: height,
	//}
	//resp, err := grpcAccessClient.GetBlockByHeight(ctx, req)
	//if err != nil {
	//	panic(fmt.Sprintf("failed to get block: %w", err))
	//}

	//fmt.Println(block.GetId())
	blockID, err := hex.DecodeString("b1c37c373cc1d7610395ec9f27a748f192a350928eb79282512d1cbc7262f9d9")
	if err != nil {
		panic(err)
	}
	txID, err := hex.DecodeString(txIDHex)
	if err != nil {
		panic(err)
	}

	requestTransactionResult := &execution.GetTransactionResultRequest{
		BlockId:       blockID,
		TransactionId: txID,
	}

	trResp, err := grpcExecutionClient.GetTransactionResult(ctx, requestTransactionResult)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("-------------------")
	fmt.Printf(" Transaction %s", txIDHex)
	fmt.Printf(" Status: %d\n", trResp.GetStatusCode())
	if trResp.GetErrorMessage() != "" {
		fmt.Printf(" Error Message: %s\n", trResp.GetErrorMessage())
	}
	events := trResp.GetEvents()
	fmt.Printf("got events %d\n", len(events))
	for _, e := range events {
		fmt.Printf(" type: %s\n", e.GetType())
		fmt.Printf(" index: %d\n", e.GetEventIndex())
		fmt.Printf(" payload: %x\n", e.GetPayload())
		fmt.Printf(" txid: %x\n", e.GetTransactionId())
		fmt.Printf(" transaction index: %d\n", e.GetTransactionIndex())
	}
	fmt.Println("-------------------")
	return nil
}

func getBlockHeaderFromAN(ctx context.Context, blockID []byte, grpcAccessClient access.AccessAPIClient) {

	req := &access.GetBlockHeaderByIDRequest{
		Id: blockID,
	}
	resp, err := grpcAccessClient.GetBlockHeaderByID(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("failed to get block: %s", err))
	}

	fmt.Printf(" block: %d\n", resp.GetBlock().GetHeight())
}

func getTransactionFromAN(ctx context.Context, txIDHex string, grpcAccessClient access.AccessAPIClient) {

	txID, err := hex.DecodeString(txIDHex)
	if err != nil {
		panic(err)
	}

	req := &access.GetTransactionRequest{
		Id: txID,
	}
	resp, err := grpcAccessClient.GetTransactionResult(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("failed to get transaction: %s", err))
	}

	fmt.Printf(" txid: %x\n", resp.GetTransactionId())
	fmt.Printf(" block: %x\n", resp.GetBlockId())
	fmt.Printf(" status: %d\n", resp.GetStatus())

	getBlockHeaderFromAN(ctx, resp.GetBlockId(), grpcAccessClient)
}

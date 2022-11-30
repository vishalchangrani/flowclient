package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/onflow/flow/protobuf/go/flow/execution"
	"google.golang.org/grpc"
)

func main() {

	MaxGRPCMessageSize := 1024 * 1024 * 20 // 20MB

	ctx := context.Background()

	executionNodeAddress := "execution-001.mainnet20.nodes.onflow.org:9000"

	// dial to it using GRPC
	conn2, err := grpc.Dial(executionNodeAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxGRPCMessageSize)))
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an ExecutionAPIClient
	grpcExecutionClient := execution.NewExecutionAPIClient(conn2)

	err = getTransactionFromExecutionNode(ctx, "b1c37c373cc1d7610395ec9f27a748f192a350928eb79282512d1cbc7262f9d9", "aa890ff09415b12005a0c233d09282abec26aaa8c42990259f3b4fc30d50f0d2", grpcExecutionClient)
	if err == nil {
		panic(err)
	}

}

func getTransactionFromExecutionNode(ctx context.Context, blockIDStr string, txIDHex string, grpcExecutionClient execution.ExecutionAPIClient) error {

	blockID, err := hex.DecodeString(blockIDStr)
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

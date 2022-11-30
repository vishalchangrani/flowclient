package main

//
//import (
//	"context"
//	"encoding/hex"
//	"fmt"
//
//	"github.com/onflow/flow/protobuf/go/flow/access"
//	"github.com/onflow/flow/protobuf/go/flow/execution"
//	"google.golang.org/grpc"
//)
//
//func main() {
//
//	accessNodeAddress := "access-001.mainnet18.nodes.onflow.org:9000"
//
//	// dial to it using GRPC
//	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
//	if err != nil {
//		panic("failed to connect to access node")
//	}
//
//	ctx := context.Background()
//
//	// create an AccessAPIClient
//	grpcAccessClient := access.NewAccessAPIClient(conn)
//
//	executionNodeAddress := "execution-001.mainnet6.nodes.onflow.org:9000"
//
//	// dial to it using GRPC
//	conn2, err := grpc.Dial(executionNodeAddress, grpc.WithInsecure())
//	if err != nil {
//		panic("failed to connect to access node")
//	}
//
//	// create an ExecutionAPIClient
//	grpcExecutionClient := execution.NewExecutionAPIClient(conn2)
//
//	for _, missing := range missingTransactions {
//		getTransactionFromEN(ctx, missing.blockHeight, missing.txID, grpcAccessClient, grpcExecutionClient)
//	}
//}
//
//func getEpochCounter(ctx context.Context, grpcExecutionClient execution.ExecutionAPIClient) {
//	req := &access.GetBlockByHeightRequest{
//		Height: height,
//	}
//	resp, err := grpcAccessClient.GetBlockByHeight(ctx, req)
//	if err != nil {
//		panic(fmt.Sprintf("failed to get block: %w", err))
//	}
//
//	block := resp.Block
//
//	txID, err := hex.DecodeString(txIDHex)
//	if err != nil {
//		panic(err)
//	}
//
//	requestTransactionResult := &execution.ExecuteScriptAtBlockIDRequest{
//		BlockId:       block.Id,
//		TransactionId: txID,
//	}
//
//	trResp, err := grpcExecutionClient.GetTransactionResult(ctx, requestTransactionResult)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("-------------------")
//	fmt.Printf(" Transaction %s", txIDHex)
//	fmt.Printf(" Status: %d\n", trResp.GetStatusCode())
//	if trResp.GetErrorMessage() != "" {
//		fmt.Printf(" Error Message: %s\n", trResp.GetErrorMessage())
//	}
//	events := trResp.GetEvents()
//	fmt.Printf("got events %d\n", len(events))
//	for _, e := range events {
//		fmt.Printf(" type: %s\n", e.GetType())
//		fmt.Printf(" index: %d\n", e.GetEventIndex())
//		fmt.Printf(" payload: %x\n", e.GetPayload())
//		fmt.Printf(" txid: %x\n", e.GetTransactionId())
//		fmt.Printf(" transaction index: %d\n", e.GetTransactionIndex())
//	}
//	fmt.Println("-------------------")
//}

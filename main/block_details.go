package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func mainsf() {
	//
	//valid := "Hellogit, 世界"
	//fmt.Println(valid)
	//invalid := string([]byte{0xff, 0xfe, 0xfd})
	//fmt.Printf("%q", invalid)
	//os.Exit(0)
	// determine the Access node IP address and port
	accessNodeAddress := "access-001.mainnet16.nodes.onflow.org:9000"

	var opts []grpc.CallOption
	opts = append(opts, grpc.UseCompressor(gzip.Name))

	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(opts...))
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	for  i:=0;i<10;i++ {
		getTx(grpcClient)
	}
	//getBlkAtHt(grpcClient)
	//getBlkRange(grpcClient)
	//ctx := context.Background()
	//
	//
	//// https://github.com/dapperlabs/flow/blob/master/docs/access-api-spec.md#getlatestblock
	//
	//// create a GetLatestBlockRequest request
	//request := access.GetLatestBlockRequest{
	//	// If the latest sealed block is needed set IsSealed to true else if the last finalized block
	//	// is needed, set IsSealed to false
	//	IsSealed: true,
	//}
	//
	//// call the GetLatestBlock GRPC method
	//response, err := grpcClient.GetLatestBlock(ctx, &request)
	//if err != nil {
	//	panic(fmt.Sprintf("failed to get block: %w", err))
	//}
	//
	//block := response.GetBlock()
	//fmt.Println(block.GetHeight())
}

func getBlkRange(grpcClient access.AccessAPIClient) {
	ctx := context.Background()
	request := &access.GetLatestBlockHeaderRequest{
		IsSealed: true,
	}
	resp, err := grpcClient.GetLatestBlockHeader(ctx, request)
	if err != nil {
		panic(fmt.Sprintf("failed to get block: %w", err))
	}
	latestHt := resp.Block.Height
	fmt.Println(latestHt)
	start := latestHt - 20000
	end := start + 2
	fmt.Println(start)

	req := &access.GetEventsForHeightRangeRequest{
		Type:        "flow.AccountCreated",
		StartHeight: start,
		EndHeight:   end,
	}
	r, err := grpcClient.GetEventsForHeightRange(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("failed to get block: %w", err))
	}
	fmt.Println(r.Results)
}

func getBlkAtHt(grpcClient access.AccessAPIClient) {
	ctx := context.Background()

	for height := uint64(60170882); height < 60170883; height++ {

		request := &access.GetBlockByHeightRequest{
			Height: height,
		}
		resp, err := grpcClient.GetBlockByHeight(ctx, request)
		if err != nil {
			panic(fmt.Sprintf("failed to get block: %v", err))
		}
		block := resp.GetBlock()
		fmt.Printf("\nBlock Height: %d\n", block.GetHeight())

		blkEventsRequest := &access.GetEventsForHeightRangeRequest{
			Type: "A.04625c28593d9408.nba_NFT.SetMetadataUpdated",
			StartHeight: 60170881,
			EndHeight: 60170883,
		}

		blkEventsResponse, err := grpcClient.GetEventsForHeightRange(ctx, blkEventsRequest)
		if err != nil {
			panic(err)
		}
		for _, r := range blkEventsResponse.GetResults() {
			fmt.Println(hex.EncodeToString(r.GetBlockId()))
			fmt.Println(len(r.Events))
			for _, e := range r.Events {
				fmt.Printf("%s\n", e.Type)
				fmt.Printf("%x\n", e.Payload)
			}
		}
		os.Exit(0)


		cols := block.GetCollectionGuarantees()
		fmt.Printf(" Found %d collection guarantees \n", len(cols))
		for _, c := range cols {
			id := c.GetCollectionId()
			req := &access.GetCollectionByIDRequest{
				Id: id,
			}
			resp, err := grpcClient.GetCollectionByID(ctx, req)
			if err != nil {
				panic(fmt.Sprintf("failed to get collection: %w", err))
			}
			coll := resp.GetCollection()
			txIDs := coll.GetTransactionIds()
		TransactionLoop:
			for _, txID := range txIDs {
				req := &access.GetTransactionRequest{
					Id: txID,
				}
				resp, err := grpcClient.GetTransaction(ctx, req)
				if err != nil {
					panic(fmt.Sprintf("failed to get tx: %w", err))
				}
				fmt.Printf("Transaction ID: %s\n", hex.EncodeToString(txID))
				fmt.Println("Transaction script")
				fmt.Println(string(resp.GetTransaction().GetScript()))
				//fmt.Println(resp.GetTransaction().GetPayer().GetAddress())
				//return

				resp2, err := grpcClient.GetTransactionResult(ctx, req)
				if err != nil {
					panic(fmt.Sprintf("failed to get tx: %v", err))

					continue TransactionLoop
				}
				fmt.Printf("Transaction ID: %s\n", resp2.Status.String())
				fmt.Printf("Transaction ErrorMsg: %s", resp2.ErrorMessage)
				fmt.Printf("Total events: %d\n", len(resp2.Events))
				for _, e:= range resp2.Events {
					fmt.Printf("%s\n", e.Type)
					fmt.Printf("%s\n", string(e.Payload))
				}
			}
		}
		break
	}
}

func getAccount(grpcClient access.AccessAPIClient) {
	//ctx := context.Background()
	//request := &access.GetAccountRequest{
	//	Address:
	//}
	//resp, err := grpcClient.GetAccount(ctx, request)
	//if err != nil {
	//	panic(fmt.Sprintf("failed to get block: %w", err))
	//}
	//fmt.Println(resp.GetAccount())
}

func getTx(grpcClient access.AccessAPIClient) {
	ctx := context.Background()

	txID := flow.HexToID("ee216673a4102198b1fa24ea87ba22043416d8a0eb39caa17c17c5949537818f")
	req := &access.GetTransactionRequest{
		Id: txID[:],
	}

	start := time.Now()
	resp2, err := grpcClient.GetTransactionResult(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("failed to get tx: %v", err))
	}
	end := time.Now()
	fmt.Println(end.Sub(start).Milliseconds())
	//fmt.Printf("Transaction ID: %s\n", resp2.Status.String())
	if "SEALED" !=  resp2.Status.String() {
		panic(fmt.Sprintf("failed to get tx: %v", err))
	}
	//fmt.Printf("Transaction ErrorMsg: %s", resp2.ErrorMessage)
	//fmt.Printf("Total events: %d\n", len(resp2.Events))
}
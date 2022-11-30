package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"google.golang.org/grpc"
)

func mainZZ(){

	// determine the Access node IP address and port
	// access-f612f908ab95e5346c4b609a0dac1ca8df7b35e7887a537c901ce5ffe26825ef@18.177.71.139:3569
	// access-001.mainnet4.nodes.onflow.org
	//accessNodeAddress := "18.177.71.139:9000"
	//accessNodeAddress := "access-001.devnet22.nodes.onflow.org:9000"
	//accessNodeAddress := "access-003.mainnet5.nodes.onflow.org:9000"
	//accessNodeAddress := "access-004.mainnet7.nodes.onflow.org:9000"
	accessNodeAddress := "access-001.canary4.nodes.onflow.org:9000"

	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	ctx := context.Background()


	// create a GetLatestBlockRequest request
	//request := access.GetLatestBlockRequest{
	//	// If the latest sealed block is needed set IsSealed to true else if the last finalized block
	//	// is needed, set IsSealed to false
	//	IsSealed: false,
	//}

	reqScript := access.ExecuteScriptAtLatestBlockRequest{
		Script: []byte("import TopShot from 0x877931736ee77cff \n pub fun main(): [UInt32]? { \n return TopShot.getPlaysInSet(setID: 2) \n }"),
	}


	//pingreq := access.PingRequest{
	//	}
	//now := time.Now()
	//// call the GetLatestBlock GRPC method to get the maximum height of the chain
	for i:=0;i<10;i++ {
		//fmt.Println(i)
		//response, err := grpcClient.GetLatestBlock(ctx, &request)
		//if err != nil {
		//	fmt.Println(err)
		//	panic("failed to get latest block")
		//}
		//
		//block := response.GetBlock()
		//maxHeight := block.GetHeight()

		//fmt.Println(maxHeight)
		//fmt.Printf("%x\n",block.GetId())
		//os.Exit(0)
		//_, err = grpcClient.Ping(ctx, &pingreq)
		//if err != nil {
		//	fmt.Println(err)
		//	panic("failed to get latest block")
		//}


		r := access.GetLatestBlockRequest{
			IsSealed: true,
		}

		rp, err := grpcClient.GetLatestBlock(ctx, &r)
		if err != nil {
			fmt.Println(err)
			//os.Exit(-1)
		} else {
			fmt.Println(rp.Block.Height)
		}

		resp, err := grpcClient.ExecuteScriptAtLatestBlock(ctx, &reqScript)
		if err != nil {
			fmt.Println(err)
			//os.Exit(-1)
		} else {
			fmt.Println(string(resp.Value))
		}

		//time.Sleep(time.Millisecond * 1005)
		//fmt.Println(time.Since(now).Milliseconds())
	}
	os.Exit(0)

	//Type: "flow.AccountCreated",

	//var startHeight = maxHeight

	//request2 := &access.GetEventsForHeightRangeRequest{
	//	Type: "A.877931736ee77cff.TopShot.Deposit",
	//	StartHeight: startHeight,
	//	EndHeight: maxHeight,
	//}

	//request2 := &access.GetEventsForHeightRangeRequest{
	//	//Type: "flow.AccountCreated",
	//	Type: "A.877931736ee77cff.TopShot.Deposit",
	//	StartHeight: 6731018,
	//	EndHeight: 6731018,
	//}

	//request2 := &access.GetEventsForHeightRangeRequest{
	//	Type: "A.1654653399040a61.FlowToken.TokensDeposited",
	//	StartHeight: uint64(12422554),
	//	EndHeight: uint64(12422554),
	//}
	//
	//response2, err := grpcClient.GetEventsForHeightRange(ctx, request2)
	//if err != nil {
	//	panic(err)
	//}
	//
	//results := response2.GetResults()
	//fmt.Printf("got %d\n", len(results))
	//
	//for _, r := range results {
	//	fmt.Println("Block ID: ", r.GetBlockId())
	//	fmt.Println("Block Ht: ", r.GetBlockHeight())
	//	events := r.GetEvents()
	//	fmt.Printf("got events %d\n", len(events))
	//		for _, e := range events {
	//			fmt.Printf(" type: %s\n", e.GetType())
	//			fmt.Printf(" index: %d\n", e.GetEventIndex())
	//			fmt.Printf(" payload: %x\n", e.GetPayload())
	//			fmt.Printf(" txid: %x\n", e.GetTransactionId())
	//			fmt.Printf(" txid: %d\n", e.GetTransactionIndex())
	//		}
	//}
	//


	txID, err := hex.DecodeString("12aa0026f100522a719ff30a90cc5d1006da61d6c78e6a91ba8af6de53375d43")
	if err != nil {
		panic(err)
	}
	request3 := &access.GetTransactionRequest{
		Id: txID,
	}

	response3, err := grpcClient.GetTransactionResult(ctx, request3)
	if err != nil {
		panic(err)
	}

	fmt.Println(response3.GetStatus())
	fmt.Println(len(response3.GetEvents()))
	return


	response4, err := grpcClient.GetTransaction(ctx, request3)
	if err != nil {
		panic(err)
	}
	tx := response4.GetTransaction()
	fmt.Println(tx)
	fmt.Println(string(tx.Script))


	os.Exit(0)
	//
	//// retrieve all blocks from 0 to maxHeight
	//for i:=uint64(0);i<maxHeight;i++ {
	//
	//	// create a GetBlockByHeightRequest request
	//	request := access.GetBlockByHeightRequest{
	//		Height: i,
	//	}
	//	response, err := grpcClient.GetBlockByHeight(ctx, &request)
	//	if err != nil {
	//		panic("failed to get block by height")
	//	}
	//
	//	block := response.GetBlock()
	//
	//	// Do something with the block here
	//	fmt.Println(block.Id)
	//
	//}


	//request5 :=  &access.GetBlockByHeightRequest{
	//	Height: 6731018,
	//}
	//response5, err :=  grpcClient.GetBlockByHeight(ctx, request5)
	//if err != nil {
	//	panic(err)
	//}
	//
	//blockID := hex.EncodeToString(response5.GetBlock().GetId())
	//fmt.Println(blockID)
	//
	//collectionIDs := response5.GetBlock().GetCollectionGuarantees()
	//fmt.Println(len(collectionIDs))
	//for _, c := range collectionIDs {
	//	fmt.Printf(" querying for collection: %s\n", hex.EncodeToString(c.GetCollectionId()))
	//	request6 := &access.GetCollectionByIDRequest{
	//		Id: c.CollectionId,
	//	}
	//	response6, err := grpcClient.GetCollectionByID(ctx, request6)
	//	if err != nil {
	//		fmt.Printf(" not found: %s\n", hex.EncodeToString(c.GetCollectionId()))
	//		continue
	//	}
	//
	//	fmt.Printf("found %d transactions\n", len(response6.GetCollection().GetTransactionIds()))
	//	for _, t := range response6.GetCollection().GetTransactionIds() {
	//		fmt.Println(" Transaction ID: " + hex.EncodeToString(t))
	//	}
	//}

	//id, err := hex.DecodeString("48ca8755e544633b38e319152394287b2a68f74b1dc3121fc32cd22f927887b9")
	//if err != nil {
	//	panic(err)
	//}
	//request7 :=  &access.GetBlockByIDRequest{
	//	Id:id,
	//}
	//response7, err :=  grpcClient.GetBlockByID(ctx, request7)
	//if err != nil {
	//	panic(err)
	//}
	//
	//collectionIDs := response7.GetBlock().GetCollectionGuarantees()
	//fmt.Printf("Block has %d collections\n", len(collectionIDs))
	//for _, c := range collectionIDs {
	//	fmt.Printf(" querying for collection: %s\n", hex.EncodeToString(c.GetCollectionId()))
	//	request6 := &access.GetCollectionByIDRequest{
	//		Id: c.CollectionId,
	//	}
	//	response6, err := grpcClient.GetCollectionByID(ctx, request6)
	//	if err != nil {
	//		fmt.Printf(" not found: %s\n", hex.EncodeToString(c.GetCollectionId()))
	//		continue
	//	}
	//
	//	fmt.Printf("found %d transactions\n", len(response6.GetCollection().GetTransactionIds()))
	//	for _, t := range response6.GetCollection().GetTransactionIds() {
	//		fmt.Println(" Transaction ID: " + hex.EncodeToString(t))
	//	}
	//}
}




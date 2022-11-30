package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"google.golang.org/grpc"
)

func mainrrr() {

	accessNodeAddress := "access.mainnet.nodes.onflow.org:9000"

	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic("failed to connect to access node")
	}

	ctx := context.Background()

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	latesfinalizedBlk := getLatestBlock(ctx, grpcClient, false)

	fmt.Printf(" latest finalized block height is: %d\n", latesfinalizedBlk.Height)

	start := time.Now()
	latestBlk := getLatestBlock(ctx, grpcClient, true)
	end := time.Now()
	timeToGetLatestBlock := end.Sub(start)

	fmt.Printf(" latest sealed block height is: %d\n", latestBlk.Height)

	start = time.Now()
	getCollectionsAndTransactions(latestBlk, ctx, grpcClient, false)
	end = time.Now()
	nonParallelTime := end.Sub(start) + timeToGetLatestBlock
	fmt.Printf("Total time %d ms (non-parallel)", (end.Sub(start) + timeToGetLatestBlock).Milliseconds())

	start = time.Now()
	getCollectionsAndTransactions(latestBlk, ctx, grpcClient, true)
	end = time.Now()
	parallelTime := end.Sub(start) + timeToGetLatestBlock
	fmt.Printf("\n\n Total time %d ms in non-parallel mode, %d ms in parallel mode\n", nonParallelTime.Milliseconds(), parallelTime.Milliseconds())

}

func getCollectionsAndTransactions(latestBlk *entities.Block, ctx context.Context, grpcClient access.AccessAPIClient, runParallel bool) {
	cols := latestBlk.GetCollectionGuarantees()

	fmt.Printf(" block has %d collections\n", len(cols))

	fmt.Println(" fetching all collections")

	var wg sync.WaitGroup
	collections := make(chan *entities.Collection, len(cols))
	for _, cg := range cols {
		coll := cg
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := getCollection(ctx, grpcClient, coll)
			if c != nil {
				collections <- c
			}
		}()
		if !runParallel {
			wg.Wait()
		}
	}
	wg.Wait()
	close(collections)

	fmt.Println(" fetching all transactions")

	for col := range collections {
		c := col
		for _, t := range c.TransactionIds {
			tx := t
			wg.Add(1)
			go func() {
				defer wg.Done()
				getTransaction(ctx, grpcClient, tx)
			}()
			if !runParallel {
				wg.Wait()
			}
		}
	}
	wg.Wait()
}

func getLatestBlock(ctx context.Context, grpcClient access.AccessAPIClient, sealed bool) *entities.Block {
	request := &access.GetLatestBlockRequest{
		IsSealed: sealed,
	}
	resp, err := grpcClient.GetLatestBlock(ctx, request)
	if err != nil {
		panic(fmt.Sprintf("failed to get block: %w", err))
	}
	return resp.Block
}

func getBlockAtHeight(ctx context.Context, grpcClient access.AccessAPIClient, height uint64) *entities.Block {
	request := &access.GetBlockByHeightRequest{
		Height: height,
	}
	resp, err := grpcClient.GetBlockByHeight(ctx, request)
	if err != nil {
		panic(fmt.Sprintf("failed to get block: %w", err))
	}
	return resp.Block
}

func getCollection(ctx context.Context, grpcClient access.AccessAPIClient, col *entities.CollectionGuarantee) *entities.Collection {
	request := &access.GetCollectionByIDRequest{
		Id: col.CollectionId,
	}
	resp, err := grpcClient.GetCollectionByID(ctx, request)
	if err != nil {
		panic(fmt.Sprintf("failed to get block: %w", err))
	}
	return resp.Collection
}

func getTransaction(ctx context.Context, grpcClient access.AccessAPIClient, id []byte) {
	req := &access.GetTransactionRequest{
		Id: id[:],
	}

	resp, err := grpcClient.GetTransactionResult(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("failed to get tx: %v", err))
	}

	fmt.Printf("transaction %s status: %s\n", flow.BytesToID(id).Hex(), resp.Status.String())
}

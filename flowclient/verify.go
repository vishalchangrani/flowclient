package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"google.golang.org/grpc"
)

var mainnet6Address = "access-001.mainnet6.nodes.onflow.org:9000"
var mainnet7Address = "access-001.mainnet7.nodes.onflow.org:9000"
var mainnet8Address = "access-003.mainnet8.nodes.onflow.org:9000"

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func connect(accessNodeAddress string) access.AccessAPIClient{

	// dial to it using GRPC
	conn, err := grpc.Dial(accessNodeAddress, grpc.WithInsecure())
	if err != nil {
		panic("failed to connect to access node")
	}

	// create an AccessAPIClient
	grpcClient := access.NewAccessAPIClient(conn)

	return grpcClient
}

func verify(ctx context.Context, transactionID string, grpcClient access.AccessAPIClient) entities.TransactionStatus {

	id, err := hex.DecodeString(transactionID)
	if err != nil {
		panic(err)
	}

	request := access.GetTransactionRequest{
		Id: id,
	}

	// call the GetTransactionResult GRPC method
	// https://github.com/dapperlabs/flow/blob/master/docs/access-api-spec.md#gettransactionresult
	response, err := grpcClient.GetTransactionResult(ctx, &request)
	checkError("failed to fetch transaction", err)

	return response.Status
}

func main234() {
	records := readCsvFile("/Users/vishalchangrani/huobi/deposits.csv")
	//fmt.Println(records)

	mainnet6Grpc := connect(mainnet6Address)
	mainnet7Grpc := connect(mainnet7Address)
	mainnet8Grpc := connect(mainnet8Address)


	file, err := os.Create("/Users/vishalchangrani/huobi/result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"block hash","transaction hash","to address","amount", "mainnet_6_status", "mainnet_7_status", "mainnet_8_status"})
	checkError("Cannot write to file", err)

	ctx := context.Background()
	var resultOn6 entities.TransactionStatus
	var resultOn7 entities.TransactionStatus
	var resultOn8 entities.TransactionStatus
	for _, record := range records[1:] {
		tx := record[1]
		fmt.Printf("querying for %s\n", tx)
		resultOn6 = verify(ctx, tx,mainnet6Grpc)
		resultOn7 = verify(ctx, tx,mainnet7Grpc)
		resultOn8 = verify(ctx, tx,mainnet8Grpc)
		var result []string
		result = append(result, record...)
		result = append(result, resultOn6.String())
		result = append(result, resultOn7.String())
		result = append(result, resultOn8.String())
		err := writer.Write(result)
		checkError("Cannot write to file", err)
		time.Sleep(10 * time.Millisecond)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
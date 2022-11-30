package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/onflow/flow/protobuf/go/flow/legacy/access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func mainzzz() {

	// determine the Access node IP address and port
	accessNodeAddress := "localhost:9000"

	//tlsCredentials, err := loadTLSCredentials()
	//if err != nil {
	//	log.Fatal("cannot load TLS credentials: ", err)
	//}

	//// dial to it using GRPC
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

	fmt.Println("server grpc response: ")
	fmt.Println(response.Block.String())
}

func loadTLSCredentials() (*custom, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("/Users/vishalchangrani/securegrpcbootstrap/bootstrap/private-root-information/private-node-info_c048c702ceb9b7e142da96e06b54ab6ab39c142b1a433639edefcd3a014a462a/cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:      certPool,
		VerifyPeerCertificate: VerifyPeerCertificate,
	}

	transportCredentials := credentials.NewTLS(config)
	customCredentials := &custom{
		transportCredentials,
	}
	return customCredentials, nil
}

type custom struct {
	credentials.TransportCredentials
}

func (c *custom) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, authInfo, err := c.TransportCredentials.ServerHandshake(conn)
	tlsInfo := authInfo.(credentials.TLSInfo)
	name := tlsInfo.State.PeerCertificates[0].Subject.CommonName
	fmt.Printf("%s\n", name)
	return conn, authInfo, err
}

func (c *custom) ClientHandshake(ctx context.Context,a string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	fmt.Printf("%s\n", a)
	conn, creds, err := c.TransportCredentials.ClientHandshake(ctx, a, conn)
	return conn, creds, err
}

func VerifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return err
	}
	//fmt.Println(cert.PublicKey)
	ecdsaKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("unknown cert type")
	}
	raw := rawEncode(ecdsaKey)
	fmt.Println("the server public key")
	keyHex := fmt.Sprintf("%#x", raw)
	fmt.Println(keyHex)
	return nil
}

func rawEncode(ecdsaPublicKey *ecdsa.PublicKey) []byte {
	xBytes := ecdsaPublicKey.X.Bytes()
	yBytes := ecdsaPublicKey.Y.Bytes()
	Plen := bitsToBytes((ecdsaPublicKey.Curve.Params().P).BitLen())
	pkEncoded := make([]byte, 2*Plen)
	// pad the public key coordinates with zeroes
	copy(pkEncoded[Plen-len(xBytes):], xBytes)
	copy(pkEncoded[2*Plen-len(yBytes):], yBytes)
	return pkEncoded
}

func bitsToBytes(bits int) int {
	return (bits + 7) >> 3
}
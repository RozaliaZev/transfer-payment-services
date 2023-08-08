package transfer_service

import (
	"fmt"
	ts "services/api/transfer_service"
	"google.golang.org/grpc"
)

//временно!
const serverAddr = "localhost:50051"

func CreateTransferServiceClient() (ts.TransferServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the gRPC server: %v", err)
	}

	client := ts.NewTransferServiceClient(conn)
	return client, conn, nil
}
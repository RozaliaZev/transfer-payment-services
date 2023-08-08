package payment_service

import (
	"fmt"
	ps "services/api/payment_service"
	"google.golang.org/grpc"
)

//временно!
const serverAddr = "localhost:50052"

func CreatePaymentServiceClient() (ps.PaymentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the gRPC server: %v", err)
	}

	client := ps.NewPaymentServiceClient(conn)
	return client, conn, nil
}
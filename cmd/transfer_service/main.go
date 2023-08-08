package main

import (
	"log"
	"net"
	"net/http"
	ts "services/api/transfer_service"
	"services/internal/transfer_service"
	"services/internal/transfer_service/handlers"

	"google.golang.org/grpc"
)

func main() {
	router := handlers.NewRouter()

    log.Println("Server started on port 8082")
    go func() {
        log.Fatal(http.ListenAndServe(":8082", router))
    }()

    grpcServer := grpc.NewServer()
    ts.RegisterTransferServiceServer(grpcServer, transfer_service.TransferServiceServer{})

    go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("gRPC server listening on :50051")
		grpcServer.Serve(lis)
	}()

    go func() {
		log.Fatal(http.ListenAndServe(":8081", router))
	}()

	// Чтобы серверы работали бесконечно и не завершались сразу после запуска, можно использовать select{}
	select {}
}

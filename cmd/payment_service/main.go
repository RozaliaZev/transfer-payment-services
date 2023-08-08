package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	ps "services/api/payment_service"
	"services/internal/payment_service"
	"services/internal/payment_service/handlers"
	"syscall"
	"services/pkg/db"
)

func main() {
	err := db.CreateTableBalances()
	if err != nil {
		log.Fatal("impossible to work with database")
	}
	grpcServer := grpc.NewServer()
	ps.RegisterPaymentServiceServer(grpcServer, payment_service.PaymentServiceServer{})

	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("gRPC server listening on :50052")
		grpcServer.Serve(lis)
	}()

	go handlers.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Println("Server stopped")

	os.Exit(0)
	select {}
}

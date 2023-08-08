API_PATH_PAYMENT = api/payment_service
PROTO_API_DIR_PAYMENT = $(API_PATH_PAYMENT)
PROTO_OUT_DIR_PAYMENT = $(API_PATH_PAYMENT)

API_PATH_TRANSFER = api/transfer_service
PROTO_API_DIR_TRANSFER = $(API_PATH_TRANSFER)
PROTO_OUT_DIR_TRANSFER = $(API_PATH_TRANSFER)

.PHONY: gen/proto/payment
gen/proto/payment: $(API_PATH_PAYMENT)
	mkdir -p $(PROTO_OUT_DIR_PAYMENT)
	protoc \
		-I $(API_PATH_PAYMENT) \
		--go_out=$(PROTO_OUT_DIR_PAYMENT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR_PAYMENT) --go-grpc_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		$(PROTO_API_DIR_PAYMENT)/*.proto

.PHONY: gen/proto/transfer
gen/proto/transfer: $(API_PATH_TRANSFER)
	mkdir -p $(PROTO_OUT_DIR_TRANSFER)
	protoc \
		-I $(API_PATH_TRANSFER) \
		--go_out=$(PROTO_OUT_DIR_TRANSFER) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR_TRANSFER) --go-grpc_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		$(PROTO_API_DIR_TRANSFER)/*.proto

run/transfer_service:
	go run cmd/transfer_service/main.go

run/payment_service:
	go run cmd/payment_service/main.go

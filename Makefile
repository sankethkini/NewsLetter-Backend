run:
	go run main.go

# run this command only in proto directory
protoc:
	protoc -I . --go_out=. --go_opt paths=source_relative --go-grpc_out=. --go-grpc_opt paths=source_relative userpb/v1/api.proto

runclient:
	go run client/user_client.go
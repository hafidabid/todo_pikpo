rm ./grpc/proto/*.pb.go || true
# protoc --go_out=grpc --go_opt=Mgrpc/proto/todo.proto=proto/ ./grpc/proto/todo.proto
protoc ./grpc/proto/todo.proto --go_out=grpc --go-grpc_out=grpc
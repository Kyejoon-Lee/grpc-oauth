package entpb

//go:generate protoc -I=.. -I./googleapis/googleapis/ --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --entgrpc_out=.. --entgrpc_opt=paths=source_relative,schema_path=../../schema entpb/entpb.proto

//go:generate protoc -I. -I/usr/local/include/protobuf event.proto --go_out=plugins=grpc:.
package event

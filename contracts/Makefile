include Makefile.deps



eventdatabase.pb.go:
	protoc --dependency_out=Makefile.deps eventdatabase.proto --go_out=plugins=grpc:.

timestamp.proto:
	cp $(GOPATH)/src/github.com/golang/protobuf/ptypes/timestamp/timestamp.proto .

duration.proto:
	cp $(GOPATH)/src/github.com/golang/protobuf/ptypes/duration/duration.proto .

all:	generated
	go install ./...

clean:
	go clean
	rm -rf genproto

APIS=$(shell find proto -name "*.proto")

descriptor:
	protoc ${APIS} \
	--proto_path='proto' \
	--include_imports \
	--descriptor_set_out=descriptor.pb

generated:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	mkdir -p genproto
	protoc ${APIS} \
	--proto_path='proto' \
	--go_opt='module=github.com/agentio/routeguide-sidecar/genproto' \
    --go_opt=Mroute_guide.proto=github.com/agentio/routeguide-sidecar/genproto/routeguidepb \
	--go_out='genproto'

lint:
	golangci-lint run

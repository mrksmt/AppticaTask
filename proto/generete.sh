export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

protoc --go_out=../api/proto --go_opt=paths=source_relative \
    --go-grpc_out=../api/proto --go-grpc_opt=paths=source_relative \
    apptop.proto

protoc --go_out=../api/proto --go_opt=paths=source_relative \
    --go-grpc_out=../api/proto --go-grpc_opt=paths=source_relative \
    data.proto

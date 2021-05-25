gen_proto:
	cd  api/proto/ && chmod +x generete.sh &&  ./generete.sh

dataservice:
	cd cmd/dataservice/ && DATA_SRV_PORT=8877 go run main.go

dataprocessor:
	cd cmd/dataprocessor/ && UPDATE_RATE=60 DATA_HOST=localhost:8877 go run main.go 

httpendpoint:
	cd cmd/httpendpoint/ && HTTP_SRV_PORT=8081 DATA_HOST=localhost:8877 go run main.go

grpcendpoint:
	cd cmd/grpcendpoint/ && GRPC_SRV_PORT=8082 DATA_HOST=localhost:8877 go run main.go

client_unary:
	cd cmd/grpcclient && GRPC_HOST=localhost:8082 REQUEST_TYPE=unary DATES="2021-05-12" go run main.go

client_streaming:
	cd cmd/grpcclient && GRPC_HOST=localhost:8082 REQUEST_TYPE=streaming DATES="2006-01-02 2021-04-12 2021-05-12 2021-05-22 55555" go run main.go
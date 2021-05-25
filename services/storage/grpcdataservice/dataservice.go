package dataservice

import (
	"context"
	"fmt"
	"log"
	"net"
	"task/api"
	"task/api/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedDataServer
}

type DataService struct {
	cancel context.CancelFunc
}

var _ api.Runnable = (*DataService)(nil)
var stor api.Storage
var cfg *Config

func New(st api.Storage, config *Config) *DataService {
	stor = st
	checkConfig(config)
	return &DataService{}
}

type Server struct {
	cancel context.CancelFunc
}

func (s *DataService) Run(mainParams *api.MainParams) error {

	log.Printf("Starting grpc server at :%s", cfg.Port)

	var localCtx context.Context
	localCtx, s.cancel = context.WithCancel(mainParams.Ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.Port, err)
	}

	srv := grpc.NewServer()

	proto.RegisterDataServer(srv, &server{})

	// run server
	mainParams.Wg.Add(1)
	go func() {
		defer mainParams.Wg.Done()
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Fails to serve: %v", err)
		}
		srv.GracefulStop()
	}()

	// stop server
	mainParams.Wg.Add(1)
	go func() {
		defer mainParams.Wg.Done()
		<-localCtx.Done()
		srv.GracefulStop()
	}()

	return nil
}

func (*server) GetData(ctx context.Context, req *proto.GetDataRequest) (*proto.GetDataResponse, error) {

	val, exist, err := stor.Get(req.Key)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	res := &proto.GetDataResponse{
		Value: val,
		Exist: exist,
		Error: errMsg,
	}

	return res, nil
}

func (*server) PutData(ctx context.Context, req *proto.PutDataRequest) (*proto.PutDataResponse, error) {

	err := stor.Put(req.Key, req.Value)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	res := &proto.PutDataResponse{
		Error: errMsg,
	}

	return res, nil
}

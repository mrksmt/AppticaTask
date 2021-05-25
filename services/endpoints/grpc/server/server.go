package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"task/api"
	"task/api/proto"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedApplicationTopPositionServer
}

type Server struct {
	cancel context.CancelFunc
}

var _ api.Runnable = (*Server)(nil)
var stor api.Storage
var cfg *Config

func New(st api.Storage, config *Config) *Server {
	stor = st
	checkConfig(config)
	return &Server{}
}

func (s *Server) Run(mainParams *api.MainParams) error {

	log.Printf("Starting grpc server at :%s", cfg.Port)

	var localCtx context.Context
	localCtx, s.cancel = context.WithCancel(mainParams.Ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.Port, err)
	}

	srv := grpc.NewServer()

	proto.RegisterApplicationTopPositionServer(srv, &server{})

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

// GetApplicationTopPositions
func (*server) GetApplicationTopPositions(ctx context.Context, req *proto.GetPositionsRequest) (*proto.GetPositionsResponse, error) {
	res := getApplicationTopPositionsResponse(req)
	return res, nil
}

// GetMccmncStream function
func (*server) GetApplicationTopPositionsStreaming(stream proto.ApplicationTopPosition_GetApplicationTopPositionsStreamingServer) error {

	for {

		// Receive the request
		req, err := stream.Recv()
		_ = req
		if err == io.EOF {
			return nil
		}
		if err != nil {
			err = fmt.Errorf("Error when reading client request stream: %v", err)
			log.Println(err)
			return err
		}

		// Process the request
		res := getApplicationTopPositionsResponse(req)

		// Send the response
		err = stream.Send(res)
		if err != nil {
			err = fmt.Errorf("Error when response was sent to the client: %v", res)
			log.Println(err)
			return err
		}

	}
}

func getApplicationTopPositionsResponse(req *proto.GetPositionsRequest) *proto.GetPositionsResponse {

	date := req.GetDate()
	_, err := time.Parse("2006-01-02", date)

	res := &proto.GetPositionsResponse{
		Date: date,
	}

	if err != nil {
		err = errors.Wrap(err, "Bad request err")
		res.Status = http.StatusBadRequest
		res.Message = err.Error()
		return res
	}

	if stor == nil {
		res.Status = http.StatusInternalServerError
		res.Message = "No data storage"
		return res
	}

	data, exist, err := stor.Get([]byte(date))

	if err != nil {
		err = errors.Wrap(err, "Storage Get err")
		res.Status = http.StatusInternalServerError
		res.Message = err.Error()
		return res
	}

	if !exist {
		res.Status = http.StatusNotFound
		res.Message = "Data not found"
		return res
	}

	res.Status = http.StatusOK
	res.Message = "OK"
	res.Data = data
	return res
}

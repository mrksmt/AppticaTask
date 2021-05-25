package grpcclientstorage

import (
	"context"
	"fmt"
	"log"
	"task/api"
	"task/api/proto"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GrpcClientStorage struct {
	cancel context.CancelFunc
	client proto.DataClient
}

var _ api.Storage = (*GrpcClientStorage)(nil)
var cfg *Config

func New(config *Config) *GrpcClientStorage {
	checkConfig(config)
	return &GrpcClientStorage{}
}

func (s *GrpcClientStorage) Run(mainParam *api.MainParams) error {

	var localCtx context.Context
	localCtx, s.cancel = context.WithCancel(mainParam.Ctx)

	mainParam.Wg.Add(1)
	go func() {
		defer mainParam.Wg.Done()

		// get connection
		con, err := getClientConn()
		if err != nil {
			err = errors.Wrap(err, "Data server connecting err")
			log.Println(err)
		}
		defer con.Close()
		s.client = proto.NewDataClient(con)

		<-localCtx.Done()
	}()

	return nil
}

func getClientConn() (*grpc.ClientConn, error) {
	opts := grpc.WithInsecure()
	con, err := grpc.Dial(cfg.Host, opts)
	if err != nil {
		return nil, errors.Wrap(err, "Connecting err")
	}
	return con, nil
}

// getGrpcData function
func (st *GrpcClientStorage) Get(key []byte) ([]byte, bool, error) {

	if st.client == nil {
		return nil, false, fmt.Errorf("No connection to data server")
	}

	req := &proto.GetDataRequest{
		Key: key,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	res, err := st.client.GetData(ctx, req)
	if err != nil {
		return nil, false, errors.Wrap(err, "GetData by GRPC err")
	}

	if res.Error != "" {
		err = fmt.Errorf(res.Error)
	}

	return res.Value, res.Exist, err
}

// Put
func (st *GrpcClientStorage) Put(key, value []byte) error {

	if st.client == nil {
		return fmt.Errorf("No connection to data server")
	}

	req := &proto.PutDataRequest{
		Key:   key,
		Value: value,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	res, err := st.client.PutData(ctx, req)
	if err != nil {
		return errors.Wrap(err, "PutData by GRPC err")
	}

	if res.Error != "" {
		err = fmt.Errorf(res.Error)
	}

	return err
}

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"task/api"
	"task/api/proto"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var cfg *Config

func RunClient(config *Config, dates ...string) {

	checkConfig(config)
	if len(dates) < 1 {
		return
	}

	// get connection

	con, err := getClientConn()
	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}
	defer con.Close()
	client := proto.NewApplicationTopPositionClient(con)

	// make request
	switch cfg.RequestType {
	case "unary":
		if err = getApplicationTopPositions(client, dates[0]); err != nil {
			log.Println(errors.Wrap(err, "Unary request err"))
		}
	case "streaming":
		if err = getApplicationTopPositionsStreaming(client, dates...); err != nil {
			log.Println(errors.Wrap(err, "Streaming request err"))
		}
	default:
		log.Printf("Unknown request type: %s", cfg.RequestType)
	}
}

func getClientConn() (*grpc.ClientConn, error) {

	opts := grpc.WithInsecure()

	con, err := grpc.Dial(cfg.Host, opts)
	if err != nil {
		return nil, errors.Wrap(err, "Connecting err")
	}

	return con, nil
}

// getMccmnc function
func getApplicationTopPositions(c proto.ApplicationTopPositionClient, date string) error {

	req := &proto.GetPositionsRequest{
		Date: date,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	res, err := c.GetApplicationTopPositions(ctx, req)
	if err != nil {
		return errors.Wrap(err, "c.GetMccmnc err")
	}
	fmt.Println(res.Data)

	return nil
}

// getApplicationTopPositionsStreaming function
func getApplicationTopPositionsStreaming(c proto.ApplicationTopPositionClient, dates ...string) error {

	// Get the stream
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	stream, err := c.GetApplicationTopPositionsStreaming(ctx)
	if err != nil {
		return errors.Wrap(err, "getApplicationTopPositionsStreaming err")
	}

	go func() {
		// Iterate over the requests slice
		for _, date := range dates {
			// Send request message
			req := &proto.GetPositionsRequest{
				Date: date,
			}

			err = stream.Send(req)
			if err != nil {
				log.Println(err)
				break
			}
		}
		// Close stream
		stream.CloseSend()
	}()

	for {
		// Get response and possible error message from the stream
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "Receiving response err")
		}

		printResponse(res)
	}

}

func printResponse(res *proto.GetPositionsResponse) {
	defer fmt.Println()

	fmt.Printf("Date:    %s\n", res.Date)
	fmt.Printf("Status:  %d\n", res.Status)
	fmt.Printf("Message: %s\n", res.Message)

	if res.Status != http.StatusOK {
		return
	}

	response := &api.ResultData{}
	json.Unmarshal([]byte(res.GetData()), response)
	prettyPrint, _ := json.MarshalIndent(response, "", "   ")
	fmt.Println(string(prettyPrint))

}

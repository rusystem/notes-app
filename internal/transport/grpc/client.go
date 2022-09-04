package grpc_client

import (
	"context"
	"fmt"
	logs "github.com/rusystem/notes-log/pkg/domain"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn       *grpc.ClientConn
	logsClient logs.LogsClient
}

func NewClient(port int) (*Client, error) {
	var conn *grpc.ClientConn

	addr := fmt.Sprintf("localhost:%d", port)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:       conn,
		logsClient: logs.NewLogsClient(conn),
	}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) LogRequest(ctx context.Context, req logs.LogItem) error {
	action, err := logs.ToPbAction(req.Action)
	if err != nil {
		return err
	}

	entity, err := logs.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}

	_, err = c.logsClient.Insert(ctx, &logs.LogRequest{
		Actions:   action,
		Entity:    entity,
		EntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})

	return err
}

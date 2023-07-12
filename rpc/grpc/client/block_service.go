package client

import (
	"context"

	v1 "github.com/cometbft/cometbft/proto/tendermint/services/block/v1"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/gogoproto/grpc"
)

// ResultBlock Single block (with meta)
type ResultBlock struct {
	BlockID types.BlockID `json:"block_id"`
	Block   *types.Block  `json:"block"`
}

// BlockServiceClient provides block information
type BlockServiceClient interface {
	GetBlock(ctx context.Context, height int64) (*ResultBlock, error)
}

type blockServiceClient struct {
	client v1.BlockServiceClient
}

func newBlockServiceClient(conn grpc.ClientConn) BlockServiceClient {
	return &blockServiceClient{
		client: v1.NewBlockServiceClient(conn),
	}
}

// GetBlock implements BlockServiceClient
func (c *blockServiceClient) GetBlock(ctx context.Context, height int64) (*ResultBlock, error) {
	req := v1.GetBlockRequest{
		Height: height,
	}
	res, err := c.client.GetBlock(ctx, &req)
	if err != nil {
		return nil, err
	}

	// convert Block from proto to core type
	block, err := types.BlockFromProto(&res.Block)
	if err != nil {
		return nil, err
	}

	// convert BlockID from proto to core type
	blockID, err := types.BlockIDFromProto(&res.BlockId)
	if err != nil {
		return nil, err
	}

	response := ResultBlock{
		BlockID: *blockID,
		Block:   block,
	}
	return &response, nil
}

type disabledBlockServiceClient struct{}

func newDisabledBlockServiceClient() BlockServiceClient {
	return &disabledBlockServiceClient{}
}

// GetBlock implements BlockServiceClient
func (*disabledBlockServiceClient) GetBlock(context.Context, int64) (*ResultBlock, error) {
	panic("block service client is disabled")
}

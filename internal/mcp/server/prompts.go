package server

import (
	"context"
	"log"

	"github.com/creachadair/jrpc2"
	"github.com/xcoulon/konverse/internal/mcp/types"
)

func listPromptsHander(logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		resp := types.ListPromptsResult{}
		return resp, nil
	}
}

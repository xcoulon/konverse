package server

import (
	"context"
	"log"

	"github.com/creachadair/jrpc2"
	"github.com/xcoulon/konverse/internal/mcp/types"
)

func listResourcesHander(logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		resp := types.ListResourcesResult{}
		// list all namespaces
		return resp, nil
	}
}

func readResourceHander(logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		params := types.ReadResourceRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			logger.Printf("error while unmarshalling '%s' request parameters: %v\n", req.Method(), err)
			return nil, err
		}

		resp := types.ReadResourceResult{}

		return resp, nil
	}
}

package server

import (
	"context"
	"log"
	"os"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/xcoulon/konverse/internal/mcp/types"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var StdioChannel = channel.Line(os.Stdin, os.Stdout)

func New(cl runtimeclient.Client, logger *log.Logger) *jrpc2.Server {
	mux := handler.Map{
		"initialize":     initializeHandler(logger),
		"prompts/list":   listPromptsHander(logger),
		"resources/list": listResourcesHander(logger),
		"resources/read": readResourceHander(logger),
		"tools/list":     listToolsHander(logger),
		"tools/call":     callToolsHander(cl, logger),
	}
	opts := &jrpc2.ServerOptions{
		Logger: jrpc2.StdLogger(logger),
	}
	return jrpc2.NewServer(mux, opts)
}

var serverInfo = types.Implementation{
	Name:    "konverse",
	Version: "0.1",
}

var protocolVersion = "2024-11-05"

func initializeHandler(logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		resp := &types.InitializeResult{
			ProtocolVersion: protocolVersion,
			ServerInfo:      serverInfo,
			Capabilities: types.ServerCapabilities{
				Resources: &types.ServerCapabilitiesResources{
					ListChanged: types.ToBoolPtr(false),
				},
				Prompts: &types.ServerCapabilitiesPrompts{
					ListChanged: types.ToBoolPtr(false),
				},
				Tools: &types.ServerCapabilitiesTools{
					ListChanged: types.ToBoolPtr(false),
				},
			},
		}
		logger.Printf("returned 'initialize' response: %v\n", resp)
		return resp, nil
	}
}

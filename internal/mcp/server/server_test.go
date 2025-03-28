package server_test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/xcoulon/konverse/internal/mcp/server"
	"github.com/xcoulon/konverse/internal/mcp/types"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestInitialize(t *testing.T) {
	// given
	c2s, s2c := channel.Direct()
	cl := newClient(c2s)
	srv := newServer().Start(s2c)
	defer func(cl *jrpc2.Client, srv *jrpc2.Server) {
		// close the streams
		cl.Close()
		srv.Stop()
	}(cl, srv)

	// when
	resp, err := cl.Call(context.Background(), "initialize", types.InitializeRequestParams{})

	// then
	require.NoError(t, err)
	expected := types.InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo: types.Implementation{
			Name:    "konverse",
			Version: "0.1",
		},
		Capabilities: types.ServerCapabilities{
			Prompts: &types.ServerCapabilitiesPrompts{
				ListChanged: types.ToBoolPtr(false),
			},
			Resources: &types.ServerCapabilitiesResources{
				ListChanged: types.ToBoolPtr(false),
			},
			Tools: &types.ServerCapabilitiesTools{
				ListChanged: types.ToBoolPtr(false),
			},
		},
	}
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedJSON), resp.ResultString())
}

func newClient(c channel.Channel) *jrpc2.Client {
	logger := log.Default()
	logger.SetOutput(os.Stderr)
	cl := jrpc2.NewClient(c, &jrpc2.ClientOptions{
		Logger: jrpc2.StdLogger(logger),
	})
	return cl
}

func newServer(initObjs ...runtime.Object) *jrpc2.Server {
	logger := log.Default()
	logger.SetOutput(os.Stderr)
	cl := fakeclient.NewFakeClient(initObjs...)
	return server.New(cl, logger)
}

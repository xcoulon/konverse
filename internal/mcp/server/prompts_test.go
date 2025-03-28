package server_test

import (
	"context"
	"testing"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xcoulon/konverse/internal/mcp/types"
)

func TestPrompts(t *testing.T) {

	// given
	c2s, s2c := channel.Direct()
	cl := newClient(c2s)
	srv := newServer().Start(s2c)

	defer func(cl *jrpc2.Client, srv *jrpc2.Server) {
		// close the streams
		cl.Close()
		srv.Stop()
	}(cl, srv)

	t.Run("list", func(t *testing.T) {
		// when
		resp, err := cl.Call(context.Background(), "prompts/list", types.ListResourcesRequestParams{})

		// then
		require.NoError(t, err)
		assert.JSONEq(t, `{"prompts":null}`, resp.ResultString())
	})
}

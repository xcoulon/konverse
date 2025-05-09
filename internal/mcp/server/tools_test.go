package server_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xcoulon/konverse/internal/mcp/types"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTools(t *testing.T) {

	// given
	pod1a := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "ns-1",
			Name:      "pod-1a",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "pod-1a-c1",
				},
			},
		},
	}
	pod1b := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "ns-1",
			Name:      "pod-1b",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "pod-1b-c1",
				},
				{
					Name: "pod-1b-c2",
				},
			},
		},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "ns-2",
			Name:      "pod-2",
		},
	}

	c2s, s2c := channel.Direct()
	cl := newClient(c2s)
	srv := newServer(pod1a, pod1b, pod2).Start(s2c)
	defer func(cl *jrpc2.Client, srv *jrpc2.Server) {
		// close the streams
		cl.Close()
		srv.Stop()
	}(cl, srv)

	t.Run("list", func(t *testing.T) {
		// when
		resp, err := cl.Call(context.Background(), "tools/list", types.ListResourcesRequestParams{})

		// then
		require.NoError(t, err)
		expected := types.ListToolsResult{
			Tools: []types.Tool{
				{
					Name:        "listPodsInNamespace",
					Description: types.ToStringPtr("list pods in the given namespace"),
					InputSchema: types.ToolInputSchema{
						Type: "object",
						Properties: types.ToolInputSchemaProperties{
							"namespace": map[string]any{
								"type":        "string",
								"description": "the namespace in which pods are running",
							},
						},
						Required: []string{"namespace"},
					},
					Annotations: &types.ToolAnnotations{
						Title:           types.ToStringPtr("List Pods"),
						DestructiveHint: types.ToBoolPtr(false),
						ReadOnlyHint:    types.ToBoolPtr(true),
					},
				},
			},
		}
		expectedJSON, _ := json.Marshal(expected)
		assert.JSONEq(t, string(expectedJSON), resp.ResultString())
	})

	t.Run("call/listPods", func(t *testing.T) {
		// given

		// when
		resp, err := cl.Call(context.Background(), "tools/call", types.CallToolRequestParams{
			Name: "listPodsInNamespace",
			Arguments: map[string]any{
				"namespace": "ns-1",
			},
		})

		// then
		require.NoError(t, err)
		expected := types.CallToolResult{
			Content: []any{
				types.TextContent{
					Type: `text`,
					Text: `{"name":"pod-1a","containers":["pod-1a-c1"]}`,
				},
				types.TextContent{
					Type: "text",
					Text: `{"name":"pod-1b","containers":["pod-1b-c1","pod-1b-c2"]}`,
				},
			},
		}
		expectedJSON, err := json.Marshal(expected)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), resp.ResultString())

	})
}

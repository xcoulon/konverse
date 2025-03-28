package server

import (
	"context"
	"log"

	"github.com/creachadair/jrpc2"
	"github.com/xcoulon/konverse/internal/mcp/types"

	corev1 "k8s.io/api/core/v1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func listToolsHander(logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		resp := &types.ListToolsResult{
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
				},
			},
		}
		logger.Printf("returned 'initialize' response: %v\n", resp)
		return resp, nil
	}
}

func callToolsHander(cl runtimeclient.Client, logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		params := types.CallToolRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			logger.Printf("error while unmarshalling '%s' request parameters: %v\n", req.Method(), err)
			return nil, err
		}

		pods := &corev1.PodList{}
		if err := cl.List(context.Background(), pods, &runtimeclient.ListOptions{
			Namespace: params.Arguments["namespace"].(string),
		}); err != nil {
			logger.Printf("error while listing pods: %v\n", err)
			return nil, err
		}
		podNames := make([]any, 0, len(pods.Items))
		for _, pod := range pods.Items {
			podNames = append(podNames, types.TextContent{
				Type: "text",
				Text: pod.Name,
			})
		}
		resp := &types.CallToolResult{
			Meta:    types.CallToolResultMeta{},
			Content: podNames,
		}
		logger.Printf("returned 'tools/call' response: %v\n", resp)
		return resp, nil
	}
}

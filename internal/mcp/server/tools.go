package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/creachadair/jrpc2"
	"github.com/xcoulon/konverse/internal/mcp/types"

	corev1 "k8s.io/api/core/v1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const listPodsInNamespace = "listPodsInNamespace"

func listToolsHander(logger *log.Logger) jrpc2.Handler {
	return func(_ context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		resp := &types.ListToolsResult{
			Tools: []types.Tool{
				{
					Name:        listPodsInNamespace,
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
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		logger.Printf("received '%s' request\n", req.Method())
		params := types.CallToolRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			logger.Printf("error while unmarshalling '%s' request parameters: %v\n", req.Method(), err)
			return nil, err
		}
		switch params.Name {
		case listPodsInNamespace:
			return callListPodsInNamespace(ctx, cl, logger, params)
		default:
			return nil, fmt.Errorf("tool '%s' does not exist", params.Name)
		}
	}
}

func callListPodsInNamespace(ctx context.Context, cl runtimeclient.Client, logger *log.Logger, params types.CallToolRequestParams) (any, error) {
	pods := &corev1.PodList{}
	if err := cl.List(ctx, pods, &runtimeclient.ListOptions{
		Namespace: params.Arguments["namespace"].(string),
	}); err != nil {
		logger.Printf("error while listing pods: %v\n", err)
		return nil, err
	}
	result := make([]any, 0, len(pods.Items))
	for _, pod := range pods.Items {
		pr := PodResult{
			Name:       pod.Name,
			Containers: make([]string, 0, len(pod.Spec.Containers)),
		}
		for _, c := range pod.Spec.Containers {
			pr.Containers = append(pr.Containers, c.Name)
		}
		prStr, err := json.Marshal(pr)
		if err != nil {
			return nil, err
		}
		result = append(result, types.TextContent{
			Type: "text",
			Text: string(prStr),
		})
	}

	// combine with some Prometheus metrics

	// rate(container_cpu_usage_seconds_total{container!~"POD|"}[5m])
	resp := &types.CallToolResult{
		Meta:    types.CallToolResultMeta{},
		Content: result,
	}
	logger.Printf("returned 'tools/call' response: %v\n", resp)
	return resp, nil
}

type PodResult struct {
	Name       string   `json:"name"`
	Containers []string `json:"containers"`
}

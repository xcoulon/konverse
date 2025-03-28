package main

import (
	"log"
	"os"

	"github.com/xcoulon/konverse/internal/k8s/client"
	"github.com/xcoulon/konverse/internal/mcp/server"

	flag "github.com/spf13/pflag"
)

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stderr)
	logger.Printf("server started")
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file (default to $KUBECONFIG or $HOME/.kube/config)")
	flag.Parse()
	cl, err := client.NewFromConfig(kubeconfig)
	if err != nil {
		logger.Fatalf("failed to init k8s client: %s", err.Error())
	}
	srv := server.New(cl, logger)
	if err := srv.Start(server.StdioChannel).Wait(); err != nil {
		logger.Fatalf("failed to wait for server: %s", err.Error())
	}

	logger.Printf("server stopped")
}

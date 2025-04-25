# Konverse

Konverse is a Model Context Protocol Server for Kubernetes, giving you access to resources, tools and prompts from a UI such as Anthropic's Claude or Block's Goose

## Building and Installing

Requires Go 1.23 and [Task](https://taskfile.dev/)

```
task build
```


### Generating the MCP types

1. Download the latest JSON schema from the [modelcontextprotocol/specification](https://github.com/modelcontextprotocol/specification/blob/main/schema/) repository

2. Install the `go-jsonschema` generator from https://github.com/omissis/go-jsonschema then run:
   ```
   go-jsonschema -p internal/types resources/schema.json > internal/types/types.go
   ```

## Testing the server with Goose CLI or UI

[Install Goose](https://block.github.io/goose/docs/getting-started/installation) then [add the MCP server](https://block.github.io/goose/docs/getting-started/using-extensions#mcp-servers) with the following command line to run:

`konverse --kubeconfig=</path/to/kubeconfig>`

## Testing the server with Claude AI Desktop App

On macOS, run the following command:

```
code ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

and add the following MCP server definition:
```
{
    "mcpServers": {
        "konverse": {
            "command": "konverse",
            "args": [
                "--kubeconfig",
                "</path/to/kubeconfig>"
            ]
        }
    }
}
```

## License

The code is available under the Apache License 2.0

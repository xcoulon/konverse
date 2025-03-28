# Konverse

Konverse is a Model Context Protocol Server for Kubernetes, giving you access to resources, tools and prompts from a UI such as Anthropic's Claude.


## Generating the MCP types

1. Download the latest JSON schema from the [modelcontextprotocol/specification](https://github.com/modelcontextprotocol/specification/blob/main/schema/) repository

2. Install the `go-jsonschema` generator from https://github.com/omissis/go-jsonschema then run:
   ```
   go-jsonschema -p internal/types resources/schema.json > internal/types/types.go
   ```

## Testing the server with Claude AI Desktop App

On macOS, run the following command:

```
code ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

## License

The code is available under the Apache License 2.0

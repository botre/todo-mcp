# TodoMCP

A proof-of-concept Model Context Protocol (MCP) server for task management.

## Scripts

Run codegen:

```bash
sqlc generate
```

Start server:

```bash
go run .
```

Start MCP Inspector:

```bash
npx @modelcontextprotocol/inspector
```

## Stack

Go, SQLite, sqlc, MCP-Go, Docker, fly.io

## Server configuration

- Transport Type: SSE
- URLs:
  - http://localhost:8080/sse
  - https://todo-mcp.fly.dev/sse

# TodoMCP

A proof-of-concept Model Context Protocol (MCP) server for task management using SQLite for storage and sqlc for type-safe database operations.

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

## Server configuration

- Transport Type: SSE
- URLs:
  - http://localhost:8080/sse
  - https://todo-mcp.fly.dev/sse

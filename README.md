# TodoMCP

A proof-of-concept Model Context Protocol (MCP) server for task management using SQLite for storage and sqlc for type-safe database operations.

## Scripts

Run codegen:

```bash
sql generate
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
- URL: http://localhost:8080/sse

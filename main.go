package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"todo-mcp/codegen"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

func formatTodo(todo codegen.Todo) string {
	return fmt.Sprintf("\n\n- **ID**: %d\n- **Title**: %s\n- **Completed**: %v\n- **Created At**: %v", todo.ID, todo.Title, todo.Completed, todo.CreatedAt)
}

func run() error {
	ctx := context.Background()

	// Connect to the database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()

	// Create schema
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Initialize codegen queries
	queries := codegen.New(db)

	// Create MCP server
	s := server.NewMCPServer(
		"TodoMCP 🚀",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Define MCP tools
	tool_completed_todos := mcp.NewTool("completed_todos",
		mcp.WithDescription("Fetches all completed todos"),
	)

	tool_pending_todos := mcp.NewTool("pending_todos",
		mcp.WithDescription("Fetches all pending todos"),
	)

	tool_create_todo := mcp.NewTool("create_todo",
		mcp.WithDescription("Creates a new todo"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title of the todo"),
		),
	)

	tool_delete_todo := mcp.NewTool("delete_todo",
		mcp.WithDescription("Deletes a specific todo by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo to delete"),
		),
	)

	tool_get_todo := mcp.NewTool("get_todo",
		mcp.WithDescription("Gets a specific todo by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo to retrieve"),
		),
	)

	tool_complete_todo := mcp.NewTool("complete_todo",
		mcp.WithDescription("Marks a specific todo as completed"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo to complete"),
		),
	)

	tool_delete_all_todos := mcp.NewTool("delete_all_todos",
		mcp.WithDescription("Deletes all todos"),
	)

	// Create MCP tool handlers
	completedTodosHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		completed, err := queries.GetAllCompletedTodos(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Format and return the completed todos
		result := "Completed Todos:\n"
		for i, todo := range completed {
			result += fmt.Sprintf("%d. %s\n", i+1, formatTodo(todo))
		}
		if len(completed) == 0 {
			result = "No completed todos found."
		}

		return mcp.NewToolResultText(result), nil
	}

	pendingTodosHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pending, err := queries.GetAllPendingTodos(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := "Pending Todos:\n"
		for i, todo := range pending {
			result += fmt.Sprintf("%d. %s\n", i+1, formatTodo(todo))
		}
		if len(pending) == 0 {
			result = "No pending todos found."
		}

		return mcp.NewToolResultText(result), nil
	}

	createTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		todo, err := queries.CreateTodo(ctx, title)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := fmt.Sprintf("Created todo: %s", formatTodo(todo))
		return mcp.NewToolResultText(result), nil
	}

	deleteTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		err = queries.DeleteTodo(ctx, int64(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := fmt.Sprintf("Deleted todo with ID: %d", id)
		return mcp.NewToolResultText(result), nil
	}

	getTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		todo, err := queries.GetTodoById(ctx, int64(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := fmt.Sprintf("Todo: %s", formatTodo(todo))
		return mcp.NewToolResultText(result), nil
	}

	completeTodoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		todo, err := queries.CompleteTodo(ctx, int64(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := fmt.Sprintf("Completed todo: %s", formatTodo(todo))
		return mcp.NewToolResultText(result), nil
	}

	deleteAllTodosHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		err := queries.DeleteAllTodos(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText("Deleted all todos"), nil
	}

	// Add tool handlers
	s.AddTool(tool_completed_todos, completedTodosHandler)
	s.AddTool(tool_pending_todos, pendingTodosHandler)
	s.AddTool(tool_create_todo, createTodoHandler)
	s.AddTool(tool_get_todo, getTodoHandler)
	s.AddTool(tool_complete_todo, completeTodoHandler)
	s.AddTool(tool_delete_todo, deleteTodoHandler)
	s.AddTool(tool_delete_all_todos, deleteAllTodosHandler)

	// Start server
	log.Println("Starting server...")
	sseServer := server.NewSSEServer(s)
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

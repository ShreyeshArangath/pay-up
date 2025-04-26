package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv" // Needed for parsing the database port

	"github.com/ShreyeshArangath/payup/internal"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server instance
	s := server.NewMCPServer(
		"Payup", // Application name
		"1.0.0", // Version
	)

	// --- Configure Database (Required for MySQL tools) ---
	// In a real application, load database configuration from environment variables,
	// a configuration file, or a secrets manager.
	// Replacing these placeholders is essential for the MySQL tools to work.
	// Using environment variables as a slightly better practice than hardcoding.
	dbUser := os.Getenv("DB_USER") // Example: load from environment
	if dbUser == "" {
		dbUser = "root" // Default placeholder
		log.Println("DB_USER environment variable not set, using default 'root'")
	}
	dbPass := os.Getenv("DB_PASSWORD") // Example: load from environment
	if dbPass == "" {
		dbPass = "password" // Default placeholder
		log.Println("DB_PASSWORD environment variable not set, using default 'password'")
	}
	dbHost := os.Getenv("DB_HOST") // Example: load from environment
	if dbHost == "" {
		dbHost = "127.0.0.1" // Default
		log.Println("DB_HOST environment variable not set, using default '127.0.0.1'")
	}
	dbPortStr := os.Getenv("DB_PORT") // Example: load from environment as string
	dbPort := 3306                    // Default
	if dbPortStr != "" {
		p, err := strconv.Atoi(dbPortStr)
		if err == nil {
			dbPort = p
		} else {
			log.Printf("Invalid DB_PORT environment variable '%s', using default %d", dbPortStr, dbPort)
		}
	} else {
		log.Printf("DB_PORT environment variable not set, using default %d", dbPort)
	}

	dbName := os.Getenv("DB_NAME") // Example: load from environment
	if dbName == "" {
		dbName = "mydatabase" // Default placeholder
		log.Println("DB_NAME environment variable not set, using default 'mydatabase'")
	}

	// Construct the DSN string
	// Add parseTime=true which is often necessary for handling MySQL dates/times with Go's time.Time
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	dbConfig := &internal.Database{
		Host: dbHost,
		Port: dbPort,
		User: dbUser,
		Pass: dbPass,
		Db:   dbName,
		DSN:  dsn, // Provide the full DSN
		// ReadOnly:         false, // Configure if needed
		// WithExplainCheck: true, // Configure if needed based on application requirements
	}

	// Optional: Attempt to connect to the database early to verify configuration
	// _, err := dbConfig.GetDB()
	// if err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }
	// log.Println("Database connection configured successfully.")

	// --- Initialize Tools ---

	// Initialize the MySQL database tools, passing the server and database config
	internal.InitializeMySQLMCPTools(s, dbConfig)
	log.Println("MySQL MCP tools initialized.")

	// Initialize the Expense calculation tools, passing the server
	internal.InitializeExpenseMCPTools(s)
	log.Println("Expense MCP tools initialized.")

	// --- Start the MCP Server ---

	// Define the server address to listen on
	// Use environment variable or default to :8080
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080" // Default listen address (e.g., on all interfaces, port 8080)
		log.Printf("LISTEN_ADDR environment variable not set, using default %s", listenAddr)
	}

	log.Printf("Starting MCP server on %s...", listenAddr)

	// Run the server. This is a blocking call.
	// It will return an error if the server cannot start or stops unexpectedly.

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

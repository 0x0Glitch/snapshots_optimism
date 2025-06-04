package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Token represents a Moonwell market token
type Token struct {
	Symbol     string `json:"symbol"`
	MTokenAddr string `json:"mTokenAddr"`
	Decimals   int    `json:"decimals"`
}

func main() {
	// Get database connection string from environment
	pgDSN := os.Getenv("PG_DSN")
	if pgDSN == "" {
		// Use the provided connection string if PG_DSN is not set
		pgDSN = "postgresql://postgres:U5hL96RqRtzaAAR7@db.jbttfyumstzuzsvnmzca.supabase.co:5432/postgres"
	}

	log.Printf("Connecting to database...")
	
	// Connect to database
	db, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Try to ping the database with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v\nPlease check your connection string and ensure PostgreSQL is running.", err)
	}
	log.Println("Successfully connected to database")

	// Load tokens to get all token symbols
	log.Printf("Loading token data...")
	tokensData, err := os.ReadFile("data/tokens.json")
	if err != nil {
		log.Fatalf("Error reading token file: %v", err)
	}

	var tokens []Token
	if err := json.Unmarshal(tokensData, &tokens); err != nil {
		log.Fatalf("Error parsing token JSON: %v", err)
	}
	log.Printf("Loaded %d tokens", len(tokens))

	// Drop existing tables if they exist
	log.Printf("Dropping existing tables if they exist...")
	
	_, err = db.Exec(`DROP TABLE IF EXISTS public.user_balances_optimism`)
	if err != nil {
		log.Printf("Warning: Error dropping user_balances_optimism table: %v", err)
	} else {
		log.Println("Dropped user_balances_optimism table if it existed")
	}
	
	// Check if any legacy tables exist with different naming
	_, err = db.Exec(`DROP TABLE IF EXISTS public.moonwell_user_balances_optimism`)
	if err != nil {
		log.Printf("Warning: Error dropping user_balances_optimism table: %v", err)
	} else {
		log.Println("Dropped user_balances_optimism table if it existed")
	}

	// Create base table
	log.Printf("Creating new user_balances_optimism table...")
	createTableSQL := `
		CREATE TABLE public.user_balances_optimism (
			user_addr TEXT NOT NULL PRIMARY KEY,
			update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating base table: %v", err)
	}

	// Add token columns
	log.Printf("Adding token columns...")
	for _, token := range tokens {
		alterSQL := fmt.Sprintf(`ALTER TABLE public.user_balances_optimism ADD COLUMN "%s" NUMERIC[] DEFAULT ARRAY[0, 0]::NUMERIC[]`, token.Symbol)
		_, err = db.Exec(alterSQL)
		if err != nil {
			log.Printf("Error adding column for token %s: %v", token.Symbol, err)
		} else {
			log.Printf("Added column for token %s (storing [mToken, borrow] values)", token.Symbol)
		}
	}

	// Create indexes
	log.Printf("Creating indexes...")
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_addr ON public.user_balances_optimism(user_addr)`)
	if err != nil {
		log.Printf("Warning: Error creating user_addr index: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_update_time ON public.user_balances_optimism(update_time)`)
	if err != nil {
		log.Printf("Warning: Error creating update_time index: %v", err)
	}

	log.Printf("Database setup completed successfully!")
} 
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// initialize the connection to the database and create the table

func InitDB() *sql.DB {
	
// environment variables defined in docker-compose
	dbUser := getEnvOrDefault("DB_USER", "usuario")
	dbPass := getEnvOrDefault("DB_PASS", "senha123")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "database")

	// connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	var err error

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Erro ao abrir conexão: %v", err)
	}

	// increase the connection timeout
	for i := 0; i < 10; i++ {
		// verify connection
		if err := db.Ping(); err != nil {
			log.Printf("Tentativa %d: Erro ao conectar no banco: %v. Tentando novamente em 5 segundos...", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Conectado ao banco:", dbName)
		break
	}

	query := `
	CREATE TABLE IF NOT EXISTS items  (
		id INT AUTO_INCREMENT PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		preco DECIMAL(10,2) NOT NULL
	)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	_, err = db.Exec(query)

	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	log.Println("A tabela items criada e verificada com sucesso!!!!")

	return db
}

// returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

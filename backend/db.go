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

// inicializa a conexao com o banco e cria a tabela
func InitDB() *sql.DB {
	// variáveis de ambiente definidas no docker-compose
	dbUser := getEnvOrDefault("DB_USER", "usuario")
	dbPass := getEnvOrDefault("DB_PASS", "senha123")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "database")

	// conexao
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

	// aumentar o tempo de espera para a conexão com o banco
	for i := 0; i < 10; i++ {
		// verifica conexão
		if err := db.Ping(); err != nil {
			log.Printf("Tentativa %d: Erro ao conectar no banco: %v. Tentando novamente em 5 segundos...", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Conectado ao banco:", dbName)
		break
	}

	// cria a tabela do banco
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

// o getEnvOrDefault retorna o valor da variável de ambiente ou um valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

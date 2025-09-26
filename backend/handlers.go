// handlers = logic of manipulation

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func createItem(w http.ResponseWriter, r *http.Request) {
	var bill Bill

	// Decode the JSON request body into the bill struct
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	bill.ID = uuid.New()
	//bill.CreatedAt = time.Now()
	//bill.UpdatedAt = time.Now()

	_, err := db.Exec(
		"INSERT INTO bills (id, embasa, coelba, created_at, updated_at) VALUES (UNHEX(REPLACE(?, '-', '')), ?, ?, ?, ?)",
		bill.ID.String(),
		bill.Embasa,
		bill.Coelba,
		bill.CreatedAt,
		bill.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(bill)
}

func readItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT 
			LOWER(CONCAT(
				SUBSTR(HEX(id), 1, 8), '-',
				SUBSTR(HEX(id), 9, 4), '-',
				SUBSTR(HEX(id), 13, 4), '-',
				SUBSTR(HEX(id), 17, 4), '-',
				SUBSTR(HEX(id), 21, 12)
			)) as id,
			embasa, 
			coelba, 
			created_at, 
			updated_at 
		FROM bills`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var bills []Bill

	for rows.Next() {
		var b Bill
		var idStr string
		err := rows.Scan(
			&idStr,
			&b.Embasa,
			&b.Coelba,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b.ID, err = uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bills = append(bills, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bills)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	// the parse make the UUID from the string
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var bill Bill

	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the updated_at for the timestamp
	bill.UpdatedAt = time.Now()

	_, err = db.Exec(
		"UPDATE bills SET embasa = ?, coelba = ?, updated_at = ? WHERE id = UNHEX(REPLACE(?, '-', ''))",
		bill.Embasa,
		bill.Coelba,
		bill.UpdatedAt,
		id.String(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bill.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM bills WHERE id = UNHEX(REPLACE(?, '-', ''))", id.String())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var bill Bill
	var idStrDB string

	err = db.QueryRow(
		`SELECT 
			LOWER(CONCAT(
				SUBSTR(HEX(id), 1, 8), '-',
				SUBSTR(HEX(id), 9, 4), '-',
				SUBSTR(HEX(id), 13, 4), '-',
				SUBSTR(HEX(id), 17, 4), '-',
				SUBSTR(HEX(id), 21, 12)
			)) as id,
			embasa, 
			coelba, 
			created_at, 
			updated_at 
		FROM bills 
		WHERE id = UNHEX(REPLACE(?, '-', ''))`,
		id.String(),
	).Scan(
		&idStrDB,
		&bill.Embasa,
		&bill.Coelba,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bill not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	bill.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

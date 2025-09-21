// handlers = logica de manipulacao

package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func createItem(w http.ResponseWriter, r *http.Request) {

	var item Item

	// decodifica o JSON  da requisicao para a struct item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest) 
		return
	}

	res, err := db.Exec("INSERT INTO items(nome, preco) VALUES(?, ?)", item.Nome, item.Preco)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
		return
	}

	id, _ := res.LastInsertId() //  recupera o id gerado pelo banco p o item que foi inserido

	item.ID = int(id) // atualiza

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item) // converte  de volta para JSON e envia como resposta http
}

func readItems(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, nome, preco FROM items")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close() 

	var items []Item

	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Nome, &i.Preco); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, i)
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func updateItem(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id") // pega o valor do parametro da rota {id} como string

	id, err := strconv.Atoi(idStr) // converte a string para int, para poder usar no banco

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var item Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE items SET nome=?, preco=? WHERE id=?", item.Nome, item.Preco, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = id

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr) 

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM items WHERE id=?", id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getItemByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr) 

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var item Item

	err = db.QueryRow("SELECT id, nome, preco FROM items WHERE id = ?", id).Scan(&item.ID, &item.Nome, &item.Preco)

	if err != nil {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}

	
	w.Header().Set("Content-Type", "application/json") // informa que os dados são JSON

	json.NewEncoder(w).Encode(item)
}

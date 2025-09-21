package main

import (
	"bytes"         // package to handle bytes
	"encoding/json" // encode and decode JSON
	"net/http"
	"net/http/httptest" //  testar HTTP
	"strconv"           //  string conversion
	"testing"           // tests

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5" //  router package
	_ "github.com/go-sql-driver/mysql"
)

// para entendimento do codigo: o handler vai retorna os dados do item em formato JSON no corpo da resposta HTTP
// para poder verificar os valores do banco (id, nome e preco) no teste precisa converter o JSON de volta para uma struct usando o marshal
// sem essa conversao nao daria para fazer as verificações dos campos do item

// obs = qualquer tipo de dados  no go pode ser convertido para []byte através da função marshal
// o marshal = converte qualquer valor (struct, map) em um array de bytes ([]byte) que vira formato JSON

// mockDB = global variable for database mock
var mockDB sqlmock.Sqlmock // define global variable mockDB

func TestMain(m *testing.M) {

	var err error 

	db, mockDB, err = sqlmock.New() // creates database mock

	if err != nil {
		panic("Error creating database mock: " + err.Error())
	}

	// assigns (atribui) the simulated connection to the global variable
	db = db

	m.Run() // runs the tests
}

func TestCreateItem(t *testing.T) {

	t.Run("sucesso", func(t *testing.T) {
		// creates item
		item := Item{
			Nome: "Monitor", 
			Preco: 850.50
		}

		body, _ := json.Marshal(item) // transforms the item into JSON

		mockDB.ExpectExec("INSERT INTO itens"). // expects INSERT execution
							WithArgs(item.Nome, item.Preco).          // expected arguments
							WillReturnResult(sqlmock.NewResult(1, 1)) // returns the id 1

		req := httptest.NewRequest("POST", "/itens", bytes.NewReader(body)) // creates POST request

		rr := httptest.NewRecorder() // creates recorder for response

		createItem(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
		}

		var createdItem Item // creates variable to decode JSON

		err := json.Unmarshal(rr.Body.Bytes(), &createdItem) // decodes response JSON (decodifica a resposta JSON)

		
		if err != nil {
			t.Fatal("Error decoding JSON response:", err) // stops the execution
		}

		// verifies if ID is correct
		if createdItem.ID != 1 {
			t.Errorf("Unexpected ID: received %d, expected %d", createdItem.ID, 1)
		}

		
		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("Mock expectations not met: %s", err)
		}
	})

	t.Run("json_invalido", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/itens", bytes.NewBufferString("json-invalido")) // if the json is invalid, so:

		rr := httptest.NewRecorder()

		createItem(rr, req) 

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusBadRequest)
		}
	})
}

func TestReadItems(t *testing.T) {

	rows := sqlmock.NewRows([]string{"id", "nome", "preco"}).
		AddRow(1, "Monitor", 850.50). 
		AddRow(2, "Teclado", 120.00)  

	mockDB.ExpectQuery("SELECT id, nome, preco FROM itens"). // expects SELECT execution
									WillReturnRows(rows) // returns the rows

	req := httptest.NewRequest("GET", "/itens", nil) // creates GET request

	rr := httptest.NewRecorder() // creates recorder (always create)

	readItems(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusOK)
	}

	var items []Item

	err := json.Unmarshal(rr.Body.Bytes(), &items) // decodes JSON response

	if err != nil { 
		t.Fatal("Error decoding JSON response:", err)
	}

	// verifies quantity of items
	if len(items) != 2 {
		t.Errorf("Unexpected number of items: received %d, expected %d", len(items), 2)
	}

	if items[0].Nome != "Monitor" { 
		t.Errorf("Unexpected item name: received %s, expected %s", items[0].Nome, "Monitor")
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations not met: %s", err)
	}
}

func TestUpdateItem(t *testing.T) {

	t.Run("sucesso", func(t *testing.T) {
		id := 1
		item := Item{
			Nome:  "Mouse Atualizado",
			Preco: 75.00} // cria item

		body, _ := json.Marshal(item)

		r := chi.NewRouter()

		r.Put("/itens/{id}", updateItem)

		req := httptest.NewRequest("PUT", "/itens/"+strconv.Itoa(id), bytes.NewReader(body)) // creates PUT request

		rr := httptest.NewRecorder() 

		mockDB.ExpectExec("UPDATE itens"). // expects UPDATE execution
							WithArgs(item.Nome, item.Preco, id).      // expected arguments
							WillReturnResult(sqlmock.NewResult(0, 1)) // returns result

		r.ServeHTTP(rr, req) // executes request on router

		if status := rr.Code; status != http.StatusOK { 
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil { 
			t.Errorf("Mock expectations not met: %s", err)
		}
	})

	t.Run("id_invalido", func(t *testing.T) { // case of invalid ID
		req := httptest.NewRequest("PUT", "/itens/nao-e-um-id", nil) // creates invalid request

		rr := httptest.NewRecorder() 

		updateItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusBadRequest)
		}
	})
}

func TestDeleteItem(t *testing.T) {

	t.Run("sucesso", func(t *testing.T) {

		id := 1

		// creates new router
		r := chi.NewRouter()

		r.Delete("/itens/{id}", deleteItem) // defines DELETE route

		req := httptest.NewRequest("DELETE", "/itens/"+strconv.Itoa(id), nil) 

		rr := httptest.NewRecorder()

		// expects DELETE execution
		mockDB.ExpectExec("DELETE FROM itens").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r.ServeHTTP(rr, req) // executes request on router

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusOK)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil { 
			t.Errorf("Mock expectations not met: %s", err)
		}
	})

	// if the ID is invalid
	t.Run("id_invalido", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/itens/nao-e-um-id", nil) // creates invalid request

		rr := httptest.NewRecorder()

		deleteItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Unexpected status code: received %v, expected %v", status, http.StatusBadRequest)
		}
	})
}

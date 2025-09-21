package main

import (
	"bytes"         // pacote para manipular bytes
	"encoding/json" // codificar e decodificar JSON
	"net/http"
	"net/http/httptest" //  testar HTTP
	"strconv"           //  conversao de string
	"testing"           // testes

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5" //  pacote para o router
	_ "github.com/go-sql-driver/mysql"
)

// para entendimento do codigo: o handler vai retorna os dados do item em formato JSON no corpo da resposta HTTP
// para poder verificar os valores do banco (id, nome e preco) no teste precisa converter o JSON de volta para uma struct usando o marshal
// sem essa conversao nao daria para fazer as verificações dos campos do item

// obs = qualquer tipo de dados  no go pode ser convertido para []byte através da função marshal
// o marshal = converte qualquer valor (struct, map) em um array de bytes ([]byte) que vira formato JSON

// mockDB = variavel global para mock do banco
var mockDB sqlmock.Sqlmock // define variavel global mockDB

func TestMain(m *testing.M) {

	var err error 

	db, mockDB, err = sqlmock.New() // cria mock do banco de dados

	if err != nil {
		panic("Erro ao criar mock do banco de dados: " + err.Error())
	}

	// atribui conexao simulada a variavel global
	db = db

	m.Run() // executa os testes
}

func TestCreateItem(t *testing.T) {

	t.Run("sucesso", func(t *testing.T) {
		// cria item
		item := Item{Nome: "Monitor", Preco: 850.50}

		body, _ := json.Marshal(item) // transforma o item em JSON

		mockDB.ExpectExec("INSERT INTO itens"). // espera execucao de INSERT
							WithArgs(item.Nome, item.Preco).          // os argumentos esperados
							WillReturnResult(sqlmock.NewResult(1, 1)) // retorna o id 1

		req := httptest.NewRequest("POST", "/itens", bytes.NewReader(body)) // cria requisicao POST

		rr := httptest.NewRecorder() // cria recorder para resposta

		createItem(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusOK)
		}

		var createdItem Item // cria variavel para decodificar JSON

		err := json.Unmarshal(rr.Body.Bytes(), &createdItem) // decodifica resposta JSON

		
		if err != nil {
			t.Fatal("Erro ao decodificar JSON da resposta:", err) //  para a execucao
		}

		// verifica se ID esta correto
		if createdItem.ID != 1 {
			t.Errorf("ID do item criado inesperado: recebido %d, esperado %d", createdItem.ID, 1)
		}

		
		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectativas do mock nao atendidas: %s", err)
		}
	})

	t.Run("json_invalido", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/itens", bytes.NewBufferString("json-invalido")) // se o json ta innvalido entao :

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

	mockDB.ExpectQuery("SELECT id, nome, preco FROM itens"). // espera execucao de SELECT
									WillReturnRows(rows) // retorna as linhas

	req := httptest.NewRequest("GET", "/itens", nil) // cria requisicao GET

	rr := httptest.NewRecorder() // cria recorder (sempre tem q criar)

	readItems(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusOK)
	}

	var items []Item

	err := json.Unmarshal(rr.Body.Bytes(), &items) // decodifica JSON

	if err != nil { 
		t.Fatal("Erro ao decodificar JSON da resposta:", err)
	}

	// verifica quantidade de itens
	if len(items) != 2 {
		t.Errorf("Numero de itens inesperado: recebido %d, esperado %d", len(items), 2)
	}

	if items[0].Nome != "Monitor" { 
		t.Errorf("Nome do item inesperado: recebido %s, esperado %s", items[0].Nome, "Monitor")
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do mock nao atendidas: %s", err)
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

		req := httptest.NewRequest("PUT", "/itens/"+strconv.Itoa(id), bytes.NewReader(body)) // cria requisicao PUT

		rr := httptest.NewRecorder() 

		mockDB.ExpectExec("UPDATE itens"). // espera execucao de UPDATE
							WithArgs(item.Nome, item.Preco, id).      // argumentos esperados
							WillReturnResult(sqlmock.NewResult(0, 1)) // retorna resultado

		r.ServeHTTP(rr, req) // executa requisicao no router

		if status := rr.Code; status != http.StatusOK { 
			t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusOK)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil { 
			t.Errorf("Expectativas do mock nao atendidas: %s", err)
		}
	})

	t.Run("id_invalido", func(t *testing.T) { // caso de ID invalido
		req := httptest.NewRequest("PUT", "/itens/nao-e-um-id", nil) // cria requisicao invalida

		rr := httptest.NewRecorder() 

		updateItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusBadRequest)
		}
	})
}

func TestDeleteItem(t *testing.T) {

	t.Run("sucesso", func(t *testing.T) {

		id := 1

		// cria novo router
		r := chi.NewRouter()

		r.Delete("/itens/{id}", deleteItem) // define a rota DELETE

		req := httptest.NewRequest("DELETE", "/itens/"+strconv.Itoa(id), nil) 

		rr := httptest.NewRecorder()

		// espera execucao do DELETE
		mockDB.ExpectExec("DELETE FROM itens").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r.ServeHTTP(rr, req) // executa requisicao no router

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusOK)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil { 
			t.Errorf("Expectativas do mock nao atendidas: %s", err)
		}
	})

	// se o ID tiver invalido
	t.Run("id_invalido", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/itens/nao-e-um-id", nil) // cria requisicao invalida

		rr := httptest.NewRecorder()

		deleteItem(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Codigo de status inesperado: recebido %v, esperado %v", status, http.StatusBadRequest)
		}
	})
}

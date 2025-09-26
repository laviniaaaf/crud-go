package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var mockDB sqlmock.Sqlmock

func TestMain(m *testing.M) {
	var err error
	db, mockDB, err = sqlmock.New()
	if err != nil {
		log.Fatal("Error in mock the bank:", err)
	}
	m.Run()
}

func TestReadItems(t *testing.T) {
	testID1 := uuid.New()
	testID2 := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "embasa", "coelba", "created_at", "updated_at"}).
		AddRow(testID1.String(), 123.45, 67.89, now, now).
		AddRow(testID2.String(), 654.32, 109.87, now, now)

	mockDB.ExpectQuery(`SELECT .* FROM bills`).WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/bills", nil)
	rr := httptest.NewRecorder()

	readItems(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status unexpected: received %v, expected %v", status, http.StatusOK)
	}

	var bills []Bill
	err := json.Unmarshal(rr.Body.Bytes(), &bills)
	if err != nil {
		t.Fatal("Error decoding JSON response:", err)
	}

	if len(bills) != 2 {
		t.Errorf("Unexpected number of bills: received %d, expected %d", len(bills), 2)
	}

	if bills[0].Embasa != 123.45 || bills[0].Coelba != 67.89 {
		t.Errorf("Unexpected data in the first bill: %+v", bills[0])
	}

	if bills[1].Embasa != 654.32 || bills[1].Coelba != 109.87 {
		t.Errorf("Unexpected data in the second bill: %+v", bills[1])
	}
}

func TestCreateItem(t *testing.T) {
	bill := Bill{
		Embasa: 123.45,
		Coelba: 67.89,
	}

	mockDB.ExpectExec(`INSERT INTO bills .*`).
		WithArgs(sqlmock.AnyArg(), bill.Embasa, bill.Coelba, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(bill)
	req := httptest.NewRequest("POST", "/bills", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	createItem(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status unexpected: received %v, expected %v. Response: %s",
			status, http.StatusOK, rr.Body.String())
	}

	var response Bill
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Error decoding JSON response:", err)
	}

	if response.Embasa != bill.Embasa || response.Coelba != bill.Coelba {
		t.Errorf("Unexpected data in the response: %+v", response)
	}
}

func TestUpdateItem(t *testing.T) {
	id := uuid.New()
	bill := Bill{
		Embasa: 200.00,
		Coelba: 300.00,
	}

	mockDB.ExpectExec(`UPDATE bills .*`).
		WithArgs(bill.Embasa, bill.Coelba, sqlmock.AnyArg(), id.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(bill)
	req := httptest.NewRequest("PUT", "/bills/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	updateItem(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status: received %v, expected %v", status, http.StatusOK)
	}

	var response Bill
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Error response JSON:", err)
	}

	if response.Embasa != bill.Embasa || response.Coelba != bill.Coelba {
		t.Errorf("Unexpected data in the response: %+v", response)
	}
}

func TestDeleteItem(t *testing.T) {
	id := uuid.New()

	mockDB.ExpectExec(`DELETE FROM bills .*`).
		WithArgs(id.String()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/bills/"+id.String(), nil)

	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	deleteItem(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("STATUS: received %v, expected %v", status, http.StatusOK)
	}
}

// models = estruturas de dados

package main

// para enviar os itens da struct no http precisa que ele esteja em formato JSON,
// pq as requisições http trabalham com JSON como texto

type Item struct {
	ID    int     `json:"id"`
	Nome  string  `json:"nome"`
	Preco float64 `json:"preco"`
}

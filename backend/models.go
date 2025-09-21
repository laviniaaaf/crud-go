// models = data structures

package main

// to send struct items to http, it must be in ( deve estar no) JSON format, 
// since HTTP requests work with JSON as text

type Item struct {
	ID    int     `json:"id"`
	Nome  string  `json:"nome"`
	Preco float64 `json:"preco"`
}

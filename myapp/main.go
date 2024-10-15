package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bda-mota/MyFirstCRUD/myapp/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.DeleteProductByID).Methods("DELETE")
	r.HandleFunc("/products/{id}", handlers.UpdateProductByID).Methods("PUT")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	fmt.Println("Host on door 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

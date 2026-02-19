package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	hand := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Добро пожаловать в мое первое приложение Go net/http!"})
	}
	http.HandleFunc("/", hand)
	http.ListenAndServe(":8080", nil)
}

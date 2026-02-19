package main

import (
	"encoding/json"
	"net/http"
)

type Nums struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

func main() {
	calculate := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var nums Nums
		json.NewDecoder(r.Body).Decode(&nums)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"result": nums.Num1 + nums.Num2})
	}
	http.HandleFunc("/calculate", calculate)
	http.ListenAndServe(":8080", nil)
}

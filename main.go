package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "notification sent",
		})
	})
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID   float64 `json:"id"`
	Name string  `json:"name"`
	Age  float64 `json:"age"`
}

func main() {
	r := chi.NewRouter()

	r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idFloat, err := strconv.ParseFloat(id, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := os.ReadFile("main.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dataJson := []User{}
		err = json.Unmarshal(data, &dataJson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		for _, user := range dataJson {
			if user.ID == idFloat {
				jsonData, err := json.Marshal(user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(jsonData)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	})

	http.ListenAndServe(":3002", r)
}

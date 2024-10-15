package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// TODO: mutex
var storeURL = make(map[string]string)

type URL struct {
	Url string `json:"url"`
}

func genStr() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) > 1 {
		w.Header().Add("Location", storeURL[parts[len(parts)-1]]) // TODO: check exist
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func handleGenURL(w http.ResponseWriter, r *http.Request) {
	var t URL
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newUrl = genStr()
	storeURL[newUrl] = t.Url

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(r.Host + "/" + newUrl))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		// В рамках net/http не придумал лучше способа.
		// Прочитал что нативной работы с /{id} нету.
		// Или еще вариант писать в 1 handler проверку на тип запроса.
		if r.Method == http.MethodPost {
			handleGenURL(w, r)
		} else {
			handleRedirect(w, r)
		}
	})

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

// TODO: mutex?
var storeURL = make(map[string]string)

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
		key := parts[len(parts)-1]
		if targetURL, exists := storeURL[key]; exists {
			w.Header().Set("Location", targetURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}

func handleGenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil || string(body) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newURL = genStr()
	storeURL[newURL] = string(body)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + r.Host + "/" + newURL))
}

func main() {

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

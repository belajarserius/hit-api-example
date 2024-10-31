package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Struktur untuk memetakan data dari API eksternal
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Fungsi untuk menghubungi API eksternal
func fetchUsersFromAPI(w http.ResponseWriter, r *http.Request) {
	// URL dari API eksternal
	url := "https://jsonplaceholder.typicode.com/users"

	// Buat request ke API eksternal
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch data from external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// log.Println("hasil : ", resp)

	// Baca data dari API eksternal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external API", http.StatusInternalServerError)
		return
	}

	// Parse JSON response
	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		http.Error(w, "Failed to parse JSON response", http.StatusInternalServerError)
		return
	}

	// Kembalikan hasil dalam bentuk JSON ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Fungsi untuk handle POST request
func createUser(w http.ResponseWriter, r *http.Request) {

	// Baca body dari request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Simulasi pengiriman data ke API eksternal
	url := "https://jsonplaceholder.typicode.com/users"
	jsonData, _ := json.Marshal(user)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to send request to external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Baca respons dari API eksternal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external API", http.StatusInternalServerError)
		return
	}

	// Kembalikan respons ke client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

// Fungsi untuk handle PUT request
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Baca body dari request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Simulasi pengiriman data update ke API eksternal
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%s", userID)
	jsonData, _ := json.Marshal(user)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request to external API", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Baca respons dari API eksternal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external API", http.StatusInternalServerError)
		return
	}

	// Kembalikan respons ke client
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	// Buat router dengan gorilla/mux
	router := mux.NewRouter()

	// Endpoint API lokal kita untuk hit API eksternal
	router.HandleFunc("/api/users", fetchUsersFromAPI).Methods("GET")

	// Endpoint POST untuk menambah data baru
	router.HandleFunc("/api/users", createUser).Methods("POST")

	// Endpoint PUT untuk mengupdate data berdasarkan ID
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")

	// Jalankan server di port 8000
	fmt.Println("Server is running at :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

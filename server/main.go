package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"github.com/gorilla/mux"
)

type FormData struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Age         int      `json:"age"`
	Role        string   `json:"role"`
	Recommend   string   `json:"recommend"`
	Improvements []string `json:"improvements"`
	Comments    string   `json:"comments"`
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Process form data
	age, err := parseAge(r.Form.Get("age"))
	if err != nil {
		http.Error(w, "Failed to parse age", http.StatusBadRequest)
		log.Println("Error parsing age:", err)
		return
	}
	formData := FormData{
		Name:         r.Form.Get("name"),
		Email:        r.Form.Get("email"),
		Age:          age,
		Role:         r.Form.Get("role"),
		Recommend:    r.Form.Get("recommend"),
		Improvements: r.Form["improvements"],
		Comments:     r.Form.Get("comments"),
	}

	// Convert form data to JSON
	response, err := json.Marshal(formData)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		log.Println("Error encoding JSON response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func parseAge(ageStr string) (int, error) {
	// Check if ageStr is empty
	if ageStr == "" {
		// If empty, return an error
		return 0, errors.New("age is empty")
	}

	// Attempt to parse ageStr to an integer
	age, err := strconv.Atoi(strings.TrimSpace(ageStr))
	if err != nil {
		// If parsing fails, return an error
		return 0, errors.New("failed to parse age")
	}

	// Additional age validation can be added here if needed

	return age, nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/submit", formHandler).Methods(http.MethodPost)

	// Start server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

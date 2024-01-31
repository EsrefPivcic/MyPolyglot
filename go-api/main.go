package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type TranslationRequest struct {
	Text string `json:"text"`
}

type TranslationResponse struct {
	Translation string `json:"translation"`
}

func translateHandler(w http.ResponseWriter, r *http.Request) {
	var request TranslationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inputFile := "input.txt"
	absInputFile, err := filepath.Abs(inputFile)
	if err != nil {
		http.Error(w, "Error getting absolute path for input file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(absInputFile, []byte(request.Text), 0644)
	if err != nil {
		http.Error(w, "Error creating input file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = os.Remove(absInputFile)
	}()

	executableDir, err := filepath.Abs("..\\marian\\build\\debug")
	if err != nil {
		http.Error(w, "Error getting absolute path: "+err.Error(), http.StatusInternalServerError)
		return
	}

	encodeCmd := exec.Command(filepath.Join(executableDir, "spm_encode.exe"), "--model=source.spm", "--output_format=piece", "--input="+absInputFile, "--output=input2.txt")
	encodeCmd.Dir = executableDir
	output, err := encodeCmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding input text: %s\nCommand output: %s", err, output), http.StatusInternalServerError)
		return
	}

	translateCmd := exec.Command(filepath.Join(executableDir, "marian-decoder.exe"), "-c", "decoder.yml", "-i", "input2.txt", "-o", "output.txt")
	translateCmd.Dir = executableDir
	if err := translateCmd.Run(); err != nil {
		http.Error(w, "Error translating text", http.StatusInternalServerError)
		return
	}

	decodeCmd := exec.Command(filepath.Join(executableDir, "spm_decode.exe"), "--model=target.spm", "--output_format=string", "--input=output.txt", "--output=output2.txt")
	decodeCmd.Dir = executableDir
	if err := decodeCmd.Run(); err != nil {
		http.Error(w, "Error decoding translated text", http.StatusInternalServerError)
		return
	}

	output2File := filepath.Join(executableDir, "output2.txt")
	translationBytes, err := os.ReadFile(output2File)
	if err != nil {
		http.Error(w, "Error reading translated text", http.StatusInternalServerError)
		return
	}
	translation := strings.TrimSpace(string(translationBytes))

	response := TranslationResponse{Translation: translation}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.HandleFunc("/translate", translateHandler)

	corsRouter := handlers.CORS(headersOk, originsOk, methodsOk)(router)

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", corsRouter)
}

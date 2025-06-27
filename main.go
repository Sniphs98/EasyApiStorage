package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "apiStorage/docs"
)

// @title File Storage API
// @version 1.0
// @description A simple file storage API with streaming support
// @host localhost:8080
// @BasePath /

const storePath = "./store"

func webInterface(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	
	tmpl.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/", webInterface).Methods("GET")
	r.HandleFunc("/upload/{filename}", uploadFile).Methods("POST")
	r.HandleFunc("/download/{filename}", downloadFile).Methods("GET")
	r.HandleFunc("/files", listFiles).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	
	fmt.Println("Server starting on :8080")
	fmt.Println("Web Interface: http://localhost:8080/")
	fmt.Println("Swagger UI: http://localhost:8080/swagger/")
	http.ListenAndServe(":8080", r)
}

// uploadFile uploads a file with streaming support
// @Summary Upload a file
// @Description Upload a file to the server with streaming support for large files
// @Tags files
// @Accept octet-stream
// @Produce plain
// @Param filename path string true "Filename"
// @Success 201 {string} string "File uploaded successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /upload/{filename} [post]
func uploadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}
	
	filename = strings.ReplaceAll(filename, "..", "")
	filePath := filepath.Join(storePath, filename)
	
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	
	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Could not write file", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File %s uploaded successfully", filename)
}

// downloadFile downloads a file with streaming support
// @Summary Download a file
// @Description Download a file from the server with streaming support for large files
// @Tags files
// @Produce octet-stream
// @Param filename path string true "Filename"
// @Success 200 {file} file "File content"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /download/{filename} [get]
func downloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}
	
	filename = strings.ReplaceAll(filename, "..", "")
	filePath := filepath.Join(storePath, filename)
	
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not open file", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()
	
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get file info", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	
	io.Copy(w, file)
}

// listFiles returns a list of all stored files
// @Summary List all files
// @Description Get a JSON array of all filenames stored on the server
// @Tags files
// @Produce json
// @Success 200 {array} string "List of filenames"
// @Failure 500 {string} string "Internal server error"
// @Router /files [get]
func listFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(storePath)
	if err != nil {
		http.Error(w, "Could not read directory", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))
	
	for i, file := range files {
		if !file.IsDir() {
			if i > 0 {
				w.Write([]byte(","))
			}
			fmt.Fprintf(w, "\"%s\"", file.Name())
		}
	}
	
	w.Write([]byte("]"))
}
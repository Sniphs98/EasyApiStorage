package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "SimpelWebFileBrowser/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title File Storage API
// @version 1.0
// @description A simple file storage API with streaming support and folder navigation
// @host localhost:8080
// @BasePath /

const storePath = "./store"

type FileInfo struct {
	Name     string    `json:"name"`
	IsFolder bool      `json:"isFolder"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"modTime"`
	Path     string    `json:"path"`
}

type CreateFolderRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func webInterface(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	// Check if store directory exists and is accessible
	if stat, err := os.Stat(storePath); err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist, create it
			if err := os.MkdirAll(storePath, 0755); err != nil {
				log.Fatal("Could not create store directory:", err)
			}
			log.Printf("Created store directory: %s", storePath)
		} else {
			log.Fatal("Could not access store directory:", err)
		}
	} else {
		log.Printf("Store directory exists: %s (is directory: %v)", storePath, stat.IsDir())
		if !stat.IsDir() {
			log.Fatal("Store path exists but is not a directory")
		}
	}

	// Test if we can read the directory
	if files, err := os.ReadDir(storePath); err != nil {
		log.Fatal("Cannot read store directory:", err)
	} else {
		log.Printf("Store directory is accessible, contains %d items", len(files))
	}

	r := mux.NewRouter()

	r.HandleFunc("/", webInterface).Methods("GET")
	r.HandleFunc("/upload/{filename}", uploadFile).Methods("POST")
	r.HandleFunc("/download/{filename:.*}", downloadFile).Methods("GET")
	r.HandleFunc("/files", listFiles).Methods("GET")
	r.HandleFunc("/create-folder", createFolder).Methods("POST")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	fmt.Println("Server starting on :8080")
	fmt.Println("Web Interface: http://localhost:8080/")
	fmt.Println("Swagger UI: http://localhost:8080/swagger/")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Helper function to validate and clean paths
func validatePath(inputPath string) (string, error) {
	if inputPath == "" {
		return "", nil
	}

	// Clean the path
	cleaned := filepath.Clean(inputPath)

	// Check for path traversal attempts
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("invalid path: contains '..'")
	}

	// Remove leading slash if present
	if strings.HasPrefix(cleaned, "/") {
		cleaned = strings.TrimPrefix(cleaned, "/")
	}

	return cleaned, nil
}

// Helper function to build safe target directory
func buildTargetDir(basePath, subPath string) (string, error) {
	cleanSubPath, err := validatePath(subPath)
	if err != nil {
		return "", err
	}

	if cleanSubPath == "" {
		return basePath, nil
	}

	return filepath.Join(basePath, cleanSubPath), nil
}

// uploadFile uploads a file with streaming support
// @Summary Upload a file
// @Description Upload a file to the server with streaming support for large files
// @Tags files
// @Accept octet-stream
// @Produce plain
// @Param filename path string true "Filename"
// @Param path query string false "Directory path"
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

	// Get and validate path
	path := r.URL.Query().Get("path")
	targetDir, err := buildTargetDir(storePath, path)
	if err != nil {
		http.Error(w, "Invalid path: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		http.Error(w, "Could not create directory", http.StatusInternalServerError)
		return
	}

	// Clean filename and create file
	filename = strings.ReplaceAll(filename, "..", "")
	filePath := filepath.Join(targetDir, filename)

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
// @Param filename path string true "Filename (can include path)"
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

	// Validate and clean the filename/path
	cleanPath, err := validatePath(filename)
	if err != nil {
		http.Error(w, "Invalid filename: "+err.Error(), http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storePath, cleanPath)

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

	// Get just the filename for the download header
	justFilename := filepath.Base(filename)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", justFilename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	io.Copy(w, file)
}

// listFiles returns a list of all files and folders in the specified directory
// @Summary List files and folders
// @Description Get a JSON array of all files and folders in the specified directory
// @Tags files
// @Produce json
// @Param path query string false "Directory path"
// @Success 200 {array} FileInfo "List of files and folders"
// @Failure 500 {string} string "Internal server error"
// @Router /files [get]
func listFiles(w http.ResponseWriter, r *http.Request) {
	// Get and validate path
	path := r.URL.Query().Get("path")
	targetDir, err := buildTargetDir(storePath, path)
	if err != nil {
		log.Printf("Invalid path error: %v", err)
		http.Error(w, "Invalid path: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Listing files in: %s (from path: %s)", targetDir, path)

	// Check if directory exists and is accessible
	if stat, err := os.Stat(targetDir); err != nil {
		log.Printf("Error accessing directory %s: %v", targetDir, err)
		http.Error(w, "Could not access directory", http.StatusInternalServerError)
		return
	} else {
		log.Printf("Directory %s exists, is directory: %v, permissions: %v", targetDir, stat.IsDir(), stat.Mode())
	}

	files, err := os.ReadDir(targetDir)
	if err != nil {
		log.Printf("Error reading directory %s: %v", targetDir, err)
		http.Error(w, "Could not read directory", http.StatusInternalServerError)
		return
	}

	log.Printf("ReadDir returned %d entries", len(files))

	var fileInfos []FileInfo
	for i, file := range files {
		log.Printf("Processing file %d: %s", i, file.Name())

		info, err := file.Info()
		if err != nil {
			log.Printf("Error getting info for file %s: %v", file.Name(), err)
			continue
		}

		var itemPath string
		if path == "" {
			itemPath = file.Name()
		} else {
			itemPath = filepath.Join(path, file.Name())
		}

		fileInfo := FileInfo{
			Name:     file.Name(),
			IsFolder: file.IsDir(),
			Size:     info.Size(),
			ModTime:  info.ModTime(),
			Path:     itemPath,
		}
		fileInfos = append(fileInfos, fileInfo)
		log.Printf("Added file info: %+v", fileInfo)
	}

	log.Printf("Found %d items total", len(fileInfos))

	w.Header().Set("Content-Type", "application/json")

	// Ensure we always return a valid JSON array, even if empty
	if fileInfos == nil {
		fileInfos = []FileInfo{}
	}

	if err := json.NewEncoder(w).Encode(fileInfos); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully sent response with %d files", len(fileInfos))
}

// createFolder creates a new folder
// @Summary Create a folder
// @Description Create a new folder in the specified directory
// @Tags files
// @Accept json
// @Produce plain
// @Param request body CreateFolderRequest true "Folder creation request"
// @Success 201 {string} string "Folder created successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /create-folder [post]
func createFolder(w http.ResponseWriter, r *http.Request) {
	var request CreateFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "Folder name is required", http.StatusBadRequest)
		return
	}

	// Clean and validate folder name
	request.Name = strings.ReplaceAll(request.Name, "..", "")
	if strings.Contains(request.Name, "/") || strings.Contains(request.Name, "\\") {
		http.Error(w, "Invalid folder name", http.StatusBadRequest)
		return
	}

	// Get and validate base path
	baseDir, err := buildTargetDir(storePath, request.Path)
	if err != nil {
		http.Error(w, "Invalid path: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create the target directory
	targetDir := filepath.Join(baseDir, request.Name)

	log.Printf("Creating folder: %s", targetDir)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Printf("Error creating folder %s: %v", targetDir, err)
		http.Error(w, "Could not create folder", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Folder %s created successfully", request.Name)
}

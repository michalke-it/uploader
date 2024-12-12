package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// FormHandler serves the HTML file upload form
func FormHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "upload.html")
}

// UploadHandler handles file upload to the /upload endpoint
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file from form", http.StatusBadRequest)
		log.Println("Error retrieving file:", err)
		return
	}
	defer file.Close()

	uploadsDir := "./uploads"
	os.MkdirAll(uploadsDir, os.ModePerm)

	dst, err := os.Create(filepath.Join(uploadsDir, handler.Filename))
	if err != nil {
		http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
		log.Println("Error creating file on server:", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to save file on server", http.StatusInternalServerError)
		log.Println("Error saving file on server:", err)
		return
	}

	log.Printf("File uploaded successfully: %s\n", handler.Filename)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func main() {
	logFile, err := os.OpenFile("uploads.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("Server started")

	http.HandleFunc("/", FormHandler)
	http.HandleFunc("/upload", UploadHandler)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

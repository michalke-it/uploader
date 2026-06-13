package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	if strings.HasSuffix(strings.ToLower(handler.Filename), ".zip") {
		// implement Unzip here
		log.Printf("Zip file uploaded successfully: %s\n", handler.Filename)
		fmt.Fprintf(w, "Zip file uploaded successfully: %s\n", handler.Filename)
	} else {
		log.Printf("File uploaded successfully: %s\n", handler.Filename)
		fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
	}
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
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

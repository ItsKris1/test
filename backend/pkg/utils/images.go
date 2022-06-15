package utils

import (
	"io/ioutil"
	"net/http"
	"os"
)

// Patht to default image location
const defaultImage = "imageUpload\\default.png"

// Creates new file and reads image bytes into it
// returns path to new image
func SaveImage(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("avatar")
	if errRead != nil {
		return defaultImage
	}
	defer file.Close()
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not ecognized return default image path
	if err != nil {
		return defaultImage
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return defaultImage
	}
	localFile.Write(fileData)
	return localFile.Name()
}

// creates empty local file based on filt type
func createTempFile(fileType string) (*os.File, error) {
	var localFile *os.File
	var err error

	switch fileType {
	case "image/jpeg":
		localFile, err = ioutil.TempFile("imageUpload", "*.jpg")
	case "image/png":
		localFile, err = ioutil.TempFile("imageUpload", "*.png")
	case "image/gif":
		localFile, err = ioutil.TempFile("imageUpload", "*.gif")
	default:
		return localFile, err
	}
	return localFile, nil
}
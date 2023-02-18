package QueryService

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, UID string) {
	// Set a limit on the size of the uploaded file
	r.ParseMultipartForm(10 << 20) // 10 MB
	// Get a handle on the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new multipart request to send the file to the destination API
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", handler.Filename)
	if err != nil {
		fmt.Println("Error Creating Form File")
		fmt.Println(err)
		return
	}
	io.Copy(part, file)
	writer.Close()

	// Create a new HTTP request to send the multipart data to the destination API
	req, err := http.NewRequest("POST", "http://localhost:5001/api/v0/add", requestBody)
	if err != nil {
		fmt.Println("Error Creating Request")
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request to the destination API and get the response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Making Request")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error Reading Response Body")
		fmt.Println(err)
		return
	}

	// Write the response back to the client
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	w.Write([]byte(UID))
}

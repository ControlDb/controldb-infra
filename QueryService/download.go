package QueryService

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "Missing hash parameter", http.StatusBadRequest)
		return
	}

	if err := ipfsGet(hash, true); err != nil {
		http.Error(w, fmt.Sprintf("Error downloading file: %s", err), http.StatusInternalServerError)
		return
	}

	filename := hash + ".json"
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening file: %s", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing file to response: %s", err), http.StatusInternalServerError)
		return
	}
}

func ipfsGet(hash string, archive bool) error {
	url := fmt.Sprintf("http://localhost:5001/api/v0/get?arg=%s", hash)
	if archive {
		url += "&archive=true"
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	// Open a file for writing the response body
	filename := hash
	if archive {
		filename += ".json"
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}
	defer file.Close()

	// Write the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Error writing response body to file: %s", err)
	}

	fmt.Printf("Downloaded file saved as %s\n", filename)
	return nil
}

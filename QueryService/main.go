package QueryService

import (
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	UID := 1
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		UploadHandler(w, r, strconv.Itoa(UID))
	})
	UID += 1
	mux.HandleFunc("/download", DownloadHandler)
	http.ListenAndServe(":8080", mux)
}

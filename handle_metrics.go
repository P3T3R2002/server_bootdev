package main

import(
	"net/http"
	"os"
	"fmt"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type:", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

//------------------------------------------------------------

func (cfg *apiConfig)handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type:", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.fileserverHits.Load())))
}

//------------------------------------------------------------

func (cfg *apiConfig)handleReset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondError(w, "", 403)
		return
	}
	cfg.db.DeleteUsers(r.Context())
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
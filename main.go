package main
import _ "github.com/lib/pq"

import(
	"net/http"
	"encoding/json"
	"log"
	"os"
	"database/sql"
	"fmt"
	"sync/atomic"
	"github.com/joho/godotenv"
	"github.com/P3T3R2002/server_bootdev/internal/database"
)


type apiConfig struct {
	db  *database.Queries
	JWT_secret string
	POLKA_KEY string
	fileserverHits atomic.Int32
}

func main() {
	fmt.Println("Runing...")
	const port = "8080"
	const root = "."

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldnt open Database!")
		return
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	apiCfg := apiConfig{
		db: dbQueries,
		JWT_secret: os.Getenv("JWT_secret"),
		POLKA_KEY: os.Getenv("POLKA_KEY"),
		fileserverHits: atomic.Int32{},
	}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))
	serveMux.HandleFunc("GET /api/healthz", handleReadiness)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.postChirps)
	serveMux.HandleFunc("POST /api/users", apiCfg.registerUsers)
	serveMux.HandleFunc("PUT /api/users", apiCfg.updateUsers)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginUser)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebhooks)
	
	var server = &http.Server{
		Addr: ":"+port,
		Handler: serveMux,
	}
	log.Fatal(server.ListenAndServe())
}

//------------------------------------------------------------

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

//**********

func respondJson(w http.ResponseWriter, r interface{}, statusCode int) {
    dat, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err) 
		w.WriteHeader(500) 
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)	
    w.Write(dat)
}

//**********

func respondError(w http.ResponseWriter, s string, statusCode int) {
    type returnVals struct {
        Error string `json:"error"`
    }

    respBody := returnVals{
		Error: s,
    }

    dat, err := json.Marshal(respBody)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err) 
			w.WriteHeader(500) 
			return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)	
    w.Write(dat)
}
package main
import "fmt"
import "net/http"
import "github.com/Kaushik-K-dev/chirpy/internal/database"

type apiConfig struct {
	fileserverHits int
	DBQueries      *database.Queries
	jwtSecret      string
	polkaKey	   string
	platform       string 
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlresp := fmt.Sprintf(`<html>

		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		
		</html>`, cfg.fileserverHits)
	fmt.Fprintf(w, htmlresp)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respError(w, http.StatusForbidden, "Not a developer")
		return
	}
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: 0\n")
	err := cfg.DBQueries.DeleteAllUsers(req.Context())
    if err != nil {
		respError(w, http.StatusInternalServerError, "Failed to reset users")
        return
    }
	respJson(w, http.StatusOK, nil)
}
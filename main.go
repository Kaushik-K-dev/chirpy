package main
import "fmt" //Change to log afterwards
import "log"
import "net/http"
import "github.com/joho/godotenv"
import "os"
import "database/sql"
import "github.com/Kaushik-K-dev/chirpy/internal/database"
import _ "github.com/lib/pq"

func readiness(w http.ResponseWriter, req *http.Request)  {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	err := godotenv.Load()
	if err !=nil {log.Fatal("Error loading .env file")}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {log.Fatal("JWT_SECRET environment variable not set")}
	polkaKey := os.Getenv("PolkaKey")
	if polkaKey == "" {log.Fatal("PolkaKey environment variable not set")}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {log.Fatal("DB_URL environment variable not set")}
	platform := os.Getenv("PLATFORM")

	mux := http.NewServeMux()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {log.Fatalf("failed to connect to the database: %v", err)}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: 0,
		DBQueries: dbQueries,
		jwtSecret: jwtSecret,
		polkaKey: polkaKey,
		platform: platform,
	}
	defer db.Close()

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpsPOSTHandler)
	//mux.HandleFunc("GET /api/chirps", apiCfg.chirpsGETHandler)
	//mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.chirpGETbyIdHandler)
	//mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.chirpDELETEHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUsersHandler)
	//mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)
	//mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	//mux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	//mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

	//mux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaHandler)

	path:= "/app/"
	fs := http.FileServer(http.Dir("."))

	mux.Handle(path, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))

	myserv := &http.Server{
		Addr:	":8080",
		Handler:mux,
	}

	err = myserv.ListenAndServe()
	if err!=nil {fmt.Printf("Server Listen/Serve Error: %v\n", err)}
}
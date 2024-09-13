package main
import "fmt" //Change to log afterwards
import "log"
import "net/http"
import "github.com/joho/godotenv"
import "os"

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

	mux := http.NewServeMux()

	db, err := NewDB("database.json")
	if err != nil {log.Fatal("Error creating database")}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB: db,
		jwtSecret: jwtSecret,
	}

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpsPOSTHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.chirpsGETHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.chirpGETbyidHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.chirpDELETEHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUsersHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)	

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
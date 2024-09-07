package main
import "fmt" //Change to log afterwards
import "net/http"

func readiness(w http.ResponseWriter, req *http.Request)  {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("/api/validate_chirp", chirpvalidHandler)


	path:= "/app/"
	fs := http.FileServer(http.Dir("."))
	//mux.Handle(path, http.StripPrefix("/app", fs))
	mux.Handle(path, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))

	myserv := &http.Server{
		Addr:	":8080",
		Handler:mux,
	}

	err := myserv.ListenAndServe()
	if err!=nil {fmt.Printf("Server Listen/Serve Error: %v\n", err)}
}
package main
import "fmt" //Change to log afterwards
import "net/http"

func main()  {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))//Currently on Root
	mux.Handle("/", fs)

	myserv := &http.Server{
		Addr:	":8080",
		Handler:mux,
	}

	err := myserv.ListenAndServe()
	if err!=nil {fmt.Printf("Server Listen/Serve Error: %v\n", err)}
}
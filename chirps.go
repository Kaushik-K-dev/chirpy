package main
import "fmt"
import "net/http"
import "encoding/json"
import "strings"

func respJson(w http.ResponseWriter, code int, dataDump interface{}){
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(dataDump)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respError(w http.ResponseWriter, code int, msg string){
	type returnError struct {
        Error string `json:"error"`
    }
	respJson(w, code, returnError{Error: msg,})
}

func chirpvalidHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
        CleanBody string `json:"cleaned_body"`
    }

	decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	const maxlength = 140
	if len(params.Body) > maxlength {
		respError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	finaltxt := profanityCheck(params.Body)

	respJson(w, http.StatusOK, returnVals{CleanBody: finaltxt,})
}

func profanityCheck(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
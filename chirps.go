package main
import ("fmt"; "net/http"; "encoding/json"; "strings"; "strconv")

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error){
	dbStruct, err := db.loadDB()
	if err != nil {return Chirp{}, err}

	createdId := len(dbStruct.Chirps)+1
	createdChirp := Chirp{
		Id: createdId,
		Body: body,
	}

	dbStruct.Chirps[createdId] = createdChirp

	err = db.writeDB(dbStruct)
	if err !=nil {return Chirp{}, err}

	return createdChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return []Chirp{}, err}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {chirps = append(chirps, chirp)}
	return chirps, nil
}

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

func (cfg *apiConfig) chirpsPOSTHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
	db := cfg.DB
	chirp, err := db.CreateChirp(finaltxt)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Couldn't create chirp")
	}

	respJson(w, http.StatusCreated, chirp)
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

func (cfg *apiConfig) chirpsGETHandler(w http.ResponseWriter, req *http.Request) {
	db := cfg.DB
	chirps, err := db.GetChirps()
	if err != nil {
		respError(w, http.StatusInternalServerError, "Could not load chirps")
		return
	}
	respJson(w, http.StatusOK, chirps)	
}

func (cfg *apiConfig) chirpGETbyidHandler(w http.ResponseWriter, req *http.Request) {
	chirpID, err := strconv.Atoi(req.PathValue("chirpID"))
	if err != nil {
		fmt.Printf("chirpID: %v, Error: %v", chirpID, err)
		respError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respError(w, http.StatusNotFound, "Couldn't find Chirp")
		return
	}

	respJson(w, http.StatusOK, Chirp{Id: chirp.Id, Body: chirp.Body,})
}
package main
import ("fmt"; "net/http"; "encoding/json"; "strings"; "time")
import "github.com/Kaushik-K-dev/chirpy/internal/database"
import "github.com/google/uuid"
//import ("strings"; "strconv"; "sort")
// import "github.com/golang-jwt/jwt/v5"

type Chirp struct {
Id uuid.UUID `json:"id"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
Body string `json:"body"`
UserId uuid.UUID `json:"user_id"`
}

// func (db *DB) CreateChirp(body, author_id string) (*Chirp, error){
// 	dbStruct, err := db.loadDB()
// 	if err != nil {return nil, err}

// 	createdId := len(dbStruct.Chirps)+1
// 	userId, err := strconv.Atoi(author_id)
// 	if err != nil {return nil, err}
// 	createdChirp := Chirp{
// 		Id: createdId,
// 		AuthorId: userId,
// 		Body: body,
// 	}

// 	dbStruct.Chirps[createdId] = &createdChirp

// 	err = db.writeDB(dbStruct)
// 	if err !=nil {return nil, err}

// 	return &createdChirp, nil
// }

// func (db *DB) DeleteChirp(chirpID int) error {
//     dbStruct, err := db.loadDB()
// 	if err != nil {return err}

//     _, err = db.GetChirp(chirpID)
// 	if err != nil {return err}

//     delete(dbStruct.Chirps, chirpID)
//     return db.writeDB(dbStruct)
// }

// func (db *DB) GetChirps() ([]*Chirp, error) {
// 	dbStruct, err := db.loadDB()
// 	if err != nil {return nil, err}

// 	chirps := make([]*Chirp, 0, len(dbStruct.Chirps))
// 	for _, chirp := range dbStruct.Chirps {chirps = append(chirps, chirp)}
// 	return chirps, nil
// }

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
	var input struct {
		Body   string `json:"body"`
    	UserId uuid.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		respError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if input.Body == "" || input.UserId == uuid.Nil {
		respError(w, http.StatusBadRequest, "Body or UserId missing")
		return
	}
	const maxlength = 140
	if len(input.Body) > maxlength {
		respError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	finaltxt := profanityCheck(input.Body)
	chirp, err := cfg.DBQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: finaltxt,
		UserID: input.UserId,
	})
	if err != nil {
		errstr := fmt.Sprintf("Error in CreateChirp: %v\n", err)
		respError(w, http.StatusInternalServerError, errstr)
		return
	}
	respJson(w, http.StatusCreated, Chirp{
		Id:			chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body:		chirp.Body,
		UserId:		chirp.UserID,
	})
}
// 	tokenString := TokenfromHeader(w, req)
// 	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
//         return []byte(cfg.jwtSecret), nil
//     })
//     if err != nil || !token.Valid {
//         respError(w, http.StatusUnauthorized, "Invalid or expired token")
//         return
//     }

//     claims, ok := token.Claims.(*jwt.RegisteredClaims)
//     if !ok {
//         respError(w, http.StatusUnauthorized, "Failed to extract claims from token")
//         return
//     }

//     userId := claims.Subject

// 	type parameters struct {
// 		Body string `json:"body"`
// 	}

// 	decoder := json.NewDecoder(req.Body)
//     params := parameters{}
//     err = decoder.Decode(&params)
//     if err != nil {
// 		respError(w, http.StatusInternalServerError, "Error decoding parameters")
// 		return
// 	}

// 	const maxlength = 140
// 	if len(params.Body) > maxlength {
// 		respError(w, http.StatusBadRequest, "Chirp is too long")
// 		return
// 	}
// 	finaltxt := profanityCheck(params.Body)
// 	chirp, err := cfg.DB.CreateChirp(finaltxt, userId)
// 	if err != nil {
// 		respError(w, http.StatusInternalServerError, "Couldn't create chirp")
// 	}

// 	respJson(w, http.StatusCreated, chirp)
// }

func profanityCheck(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

// func (cfg *apiConfig) chirpsGETHandler(w http.ResponseWriter, req *http.Request) {
// 	db := cfg.DB
// 	chirps, err := db.GetChirps()
// 	if err != nil {
// 		respError(w, http.StatusInternalServerError, "Could not load chirps")
// 		return
// 	}
// 	AuthorId := req.URL.Query().Get("author_id")
// 	sortOrder := req.URL.Query().Get("sort")
// 	var chirpsbyAuth []*Chirp
// 	if AuthorId != "" {
// 		authId, err := strconv.Atoi(AuthorId)
// 		if err != nil {
// 			respError(w, http.StatusBadRequest, "Invalid author_id")
// 			return
// 		}

// 		for _, chirp := range chirps {
// 			if chirp.AuthorId == authId {
// 				chirpsbyAuth = append(chirpsbyAuth, chirp)
// 			}
// 		}
// 	} else {chirpsbyAuth = chirps}

// 	if sortOrder == "" || sortOrder == "asc" {
// 		sort.Slice(chirpsbyAuth, func(i, j int) bool {return chirpsbyAuth[i].Id < chirpsbyAuth[j].Id})
// 	} else if sortOrder == "desc" {
// 		sort.Slice(chirpsbyAuth, func(i, j int) bool {return chirpsbyAuth[i].Id > chirpsbyAuth[j].Id})
// 	} else {
// 		respError(w, http.StatusBadRequest, "Invalid sort parameter. Must be 'asc' or 'desc'.")
// 		return
// 	}
// 	respJson(w, http.StatusOK, chirpsbyAuth)
// }

// func (cfg *apiConfig) chirpGETbyIdHandler(w http.ResponseWriter, req *http.Request) {
// 	chirpID, err := strconv.Atoi(req.PathValue("chirpID"))
// 	if err != nil {
// 		fmt.Printf("chirpID: %v, Error: %v", chirpID, err)
// 		respError(w, http.StatusBadRequest, "Invalid Chirp ID")
// 		return
// 	}

// 	chirp, err := cfg.DB.GetChirp(chirpID)
// 	if err != nil {
// 		respError(w, http.StatusNotFound, "Couldn't find Chirp")
// 		return
// 	}

// 	respJson(w, http.StatusOK, Chirp{Id: chirp.Id, AuthorId: chirp.AuthorId, Body: chirp.Body,})
// }

// func (cfg *apiConfig) chirpDELETEHandler(w http.ResponseWriter, req *http.Request) {
// 	tokenString := TokenfromHeader(w, req)
// 	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
//         return []byte(cfg.jwtSecret), nil
//     })
//     if err != nil || !token.Valid {
//         respError(w, http.StatusUnauthorized, "Invalid or expired token")
//         return
//     }

//     claims, ok := token.Claims.(*jwt.RegisteredClaims)
//     if !ok {
//         respError(w, http.StatusUnauthorized, "Failed to extract claims from token")
//         return
//     }

//     userId := claims.Subject

// 	chirpID, err := strconv.Atoi(req.PathValue("chirpID"))
// 	if err != nil {
// 		fmt.Printf("chirpID: %v, Error: %v", chirpID, err)
// 		respError(w, http.StatusBadRequest, "Invalid Chirp ID")
// 		return
// 	}

// 	chirp, err := cfg.DB.GetChirp(chirpID)
// 	if err != nil {
// 		respError(w, http.StatusNotFound, "Couldn't find Chirp")
// 		return
// 	}
// 	authId, _ := strconv.Atoi(userId)
// 	if chirp.AuthorId != authId {
// 			respError(w, http.StatusForbidden, "Unauthorized to delete Chirp")
// 			return
// 		}

// 	w.WriteHeader(http.StatusNoContent)
// }
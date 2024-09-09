package main
import ("fmt"; "encoding/json"; "net/http")

type User struct{
	Id int `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return User{}, err}

	id := len(dbStruct.Users) + 1
	user := User{
		Id:    id,
		Email: email,
	}
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {return User{}, err}
	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return User{}, err}

	user, ok := dbStruct.Users[id]
	if !ok {return User{}, fmt.Errorf("Error: User doesn't exist")}

	return user, nil
}

func (cfg *apiConfig) createUsersHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		respError(w, http.StatusInternalServerError, "Couldn't create User")
		return
	}

	respJson(w, http.StatusCreated, User{
		Id:    user.Id,
		Email: user.Email,
	})
}